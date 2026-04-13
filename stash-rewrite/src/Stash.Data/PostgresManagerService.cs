using System.Diagnostics;
using System.IO.Compression;
using System.Runtime.InteropServices;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;
using Stash.Core.Interfaces;

namespace Stash.Data;

/// <summary>
/// Manages a self-contained PostgreSQL instance that starts/stops with the app.
/// On first run, downloads portable PostgreSQL binaries automatically.
/// </summary>
public class PostgresManagerService : IHostedService
{
    private readonly PostgresConfig _config;
    private readonly ILogger<PostgresManagerService> _logger;
    private bool _started;

    // PostgreSQL 16.8 portable binaries - stable LTS release
    private const string PgVersion = "16.8-1";
    private const string WinUrl = $"https://get.enterprisedb.com/postgresql/postgresql-{PgVersion}-windows-x64-binaries.zip";
    private const string LinuxUrl = $"https://get.enterprisedb.com/postgresql/postgresql-{PgVersion}-linux-x64-binaries.tar.gz";

    public PostgresManagerService(IOptions<PostgresConfig> config, ILogger<PostgresManagerService> logger)
    {
        _config = config.Value;
        _logger = logger;
    }

    /// <summary>Root directory for all managed postgres files (binaries + data).</summary>
    private string StashDir => _config.DataPath ?? Path.Combine(
        Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData), "stash");

    private string BinDir => Path.Combine(StashDir, "pgsql", "bin");
    private string DataDir => Path.Combine(StashDir, "pgdata");
    private string LogFile => Path.Combine(StashDir, "pg.log");

    private string Exe(string name) =>
        RuntimeInformation.IsOSPlatform(OSPlatform.Windows) ? Path.Combine(BinDir, $"{name}.exe")
                                                            : Path.Combine(BinDir, name);

    // ─── Lifecycle ──────────────────────────────────────────────────

    public async Task StartAsync(CancellationToken ct)
    {
        if (!_config.Managed)
        {
            _logger.LogInformation("Managed PostgreSQL disabled — using external connection string");
            return;
        }

        _logger.LogInformation("Managed PostgreSQL mode enabled");

        // 1. Ensure binaries exist (download if needed)
        if (!File.Exists(Exe("pg_ctl")))
        {
            _logger.LogInformation("PostgreSQL binaries not found — downloading portable {Version}…", PgVersion);
            await DownloadPostgresAsync(ct);
        }
        else
        {
            _logger.LogInformation("PostgreSQL binaries found at {BinDir}", BinDir);
        }

        // 2. Check if a stale instance exists from a previous crash
        await StopStaleInstanceAsync(ct);

        // 3. Init data directory if needed
        if (!File.Exists(Path.Combine(DataDir, "PG_VERSION")))
        {
            _logger.LogInformation("Initializing data directory at {DataDir}", DataDir);
            await InitDbAsync(ct);
        }

        // 4. Start PostgreSQL
        _logger.LogInformation("Starting PostgreSQL on port {Port}", _config.Port);
        await PgCtlAsync($"start -D \"{DataDir}\" -l \"{LogFile}\" -w -o \"-p {_config.Port}\"", ct);
        _started = true;

        // 5. Wait for ready
        await WaitForReadyAsync(ct);

        // 6. Create database if it doesn't exist
        await EnsureDatabaseAsync(ct);

        _logger.LogInformation("Managed PostgreSQL is ready (port {Port}, database '{Db}')", _config.Port, _config.Database);
    }

    public async Task StopAsync(CancellationToken ct)
    {
        if (!_config.Managed || !_started) return;

        _logger.LogInformation("Stopping managed PostgreSQL");
        try
        {
            await PgCtlAsync($"stop -D \"{DataDir}\" -m fast", ct);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Error during PostgreSQL shutdown — it may already be stopped");
        }
        _started = false;
    }

    // ─── Download ───────────────────────────────────────────────────

    private async Task DownloadPostgresAsync(CancellationToken ct)
    {
        Directory.CreateDirectory(StashDir);

        bool isWindows = RuntimeInformation.IsOSPlatform(OSPlatform.Windows);
        string url = isWindows ? WinUrl : LinuxUrl;
        string ext = isWindows ? ".zip" : ".tar.gz";
        string archivePath = Path.Combine(StashDir, $"postgresql{ext}");

        // Download with progress
        using var http = new HttpClient { Timeout = TimeSpan.FromMinutes(10) };
        _logger.LogInformation("Downloading {Url}", url);

        using var response = await http.GetAsync(url, HttpCompletionOption.ResponseHeadersRead, ct);
        response.EnsureSuccessStatusCode();

        var totalBytes = response.Content.Headers.ContentLength ?? -1;
        await using var contentStream = await response.Content.ReadAsStreamAsync(ct);
        await using var fileStream = new FileStream(archivePath, FileMode.Create, FileAccess.Write, FileShare.None, 81920);

        var buffer = new byte[81920];
        long totalRead = 0;
        int lastPct = -1;
        int bytesRead;

        while ((bytesRead = await contentStream.ReadAsync(buffer, ct)) > 0)
        {
            await fileStream.WriteAsync(buffer.AsMemory(0, bytesRead), ct);
            totalRead += bytesRead;
            if (totalBytes > 0)
            {
                int pct = (int)(totalRead * 100 / totalBytes);
                if (pct / 10 > lastPct / 10) // log every 10%
                {
                    _logger.LogInformation("  Download progress: {Pct}% ({MB:F0} MB)",
                        pct, totalRead / 1048576.0);
                    lastPct = pct;
                }
            }
        }
        await fileStream.FlushAsync(ct);
        fileStream.Close();

        _logger.LogInformation("Download complete ({MB:F1} MB). Extracting…", totalRead / 1048576.0);

        // Extract
        if (isWindows)
        {
            ZipFile.ExtractToDirectory(archivePath, StashDir, overwriteFiles: true);
        }
        else
        {
            // tar.gz on Linux/macOS
            var exitCode = await RunAsync("/bin/tar", $"xzf \"{archivePath}\" -C \"{StashDir}\"", StashDir, ct);
            if (exitCode != 0)
                throw new InvalidOperationException("Failed to extract PostgreSQL archive");

            // Make binaries executable
            await RunAsync("/bin/chmod", $"-R +x \"{BinDir}\"", StashDir, ct);
        }

        // Clean up archive
        File.Delete(archivePath);

        if (!File.Exists(Exe("pg_ctl")))
            throw new FileNotFoundException(
                $"Extraction succeeded but pg_ctl not found at expected path: {Exe("pg_ctl")}. " +
                $"Contents of {StashDir}: {string.Join(", ", Directory.GetDirectories(StashDir))}");

        _logger.LogInformation("PostgreSQL {Version} binaries ready at {BinDir}", PgVersion, BinDir);
    }

    // ─── Init / Start / Stop helpers ────────────────────────────────

    private async Task InitDbAsync(CancellationToken ct)
    {
        Directory.CreateDirectory(DataDir);
        var exitCode = await RunAsync(Exe("initdb"),
            $"-D \"{DataDir}\" -U postgres --encoding=UTF8 --locale=C --auth=trust",
            BinDir, ct);

        if (exitCode != 0)
            throw new InvalidOperationException($"initdb failed (exit code {exitCode}). Check {LogFile}");

        // Write pg_hba.conf — local-only trust auth
        await File.WriteAllTextAsync(Path.Combine(DataDir, "pg_hba.conf"),
            """
            # TYPE  DATABASE  USER  ADDRESS       METHOD
            local   all       all                 trust
            host    all       all   127.0.0.1/32  trust
            host    all       all   ::1/128       trust
            """, ct);

        // Append to postgresql.conf
        await File.AppendAllTextAsync(Path.Combine(DataDir, "postgresql.conf"),
            $"""

            # ── Stash managed ──
            port = {_config.Port}
            listen_addresses = '127.0.0.1'
            max_connections = 20
            shared_buffers = 128MB
            log_destination = 'stderr'
            logging_collector = off
            """, ct);
    }

    private async Task PgCtlAsync(string args, CancellationToken ct)
    {
        var exitCode = await RunAsync(Exe("pg_ctl"), args, BinDir, ct);
        if (exitCode != 0)
        {
            // Read log for diagnostics
            var logContent = File.Exists(LogFile) ? await File.ReadAllTextAsync(LogFile, ct) : "(no log file)";
            var lastLines = string.Join('\n', logContent.Split('\n').TakeLast(20));
            throw new InvalidOperationException(
                $"pg_ctl failed (exit code {exitCode}). Last log lines:\n{lastLines}");
        }
    }

    private async Task StopStaleInstanceAsync(CancellationToken ct)
    {
        var pidFile = Path.Combine(DataDir, "postmaster.pid");
        if (!File.Exists(pidFile)) return;

        _logger.LogInformation("Found stale postmaster.pid — stopping previous instance");
        try
        {
            await RunAsync(Exe("pg_ctl"), $"stop -D \"{DataDir}\" -m fast", BinDir, ct);
        }
        catch
        {
            // If it fails (process already dead), just remove the pid file
            try { File.Delete(pidFile); } catch { }
        }
    }

    private async Task WaitForReadyAsync(CancellationToken ct)
    {
        for (int i = 0; i < 30; i++)
        {
            ct.ThrowIfCancellationRequested();
            var exitCode = await RunAsync(Exe("pg_isready"),
                $"-h 127.0.0.1 -p {_config.Port} -U postgres", BinDir, ct);
            if (exitCode == 0)
            {
                _logger.LogDebug("PostgreSQL is accepting connections");
                return;
            }
            await Task.Delay(500, ct);
        }

        var logContent = File.Exists(LogFile) ? await File.ReadAllTextAsync(LogFile, ct) : "(no log)";
        throw new TimeoutException(
            $"PostgreSQL did not become ready within 15 seconds. Log:\n{string.Join('\n', logContent.Split('\n').TakeLast(30))}");
    }

    private async Task EnsureDatabaseAsync(CancellationToken ct)
    {
        // Check if database exists via psql
        var (exitCode, stdout) = await RunWithOutputAsync(Exe("psql"),
            $"-h 127.0.0.1 -p {_config.Port} -U postgres -tAc \"SELECT 1 FROM pg_database WHERE datname='{_config.Database}'\"",
            BinDir, ct);

        if (stdout.Trim() == "1")
        {
            _logger.LogDebug("Database '{Db}' already exists", _config.Database);

            // Ensure pgvector extension is created
            await RunAsync(Exe("psql"),
                $"-h 127.0.0.1 -p {_config.Port} -U postgres -d {_config.Database} -c \"CREATE EXTENSION IF NOT EXISTS vector\"",
                BinDir, ct);
            return;
        }

        _logger.LogInformation("Creating database '{Db}'", _config.Database);
        exitCode = await RunAsync(Exe("createdb"),
            $"-h 127.0.0.1 -p {_config.Port} -U postgres {_config.Database}", BinDir, ct);

        if (exitCode != 0)
            throw new InvalidOperationException($"createdb failed (exit code {exitCode})");

        // Try to create pgvector extension (will fail silently if not available)
        var extResult = await RunAsync(Exe("psql"),
            $"-h 127.0.0.1 -p {_config.Port} -U postgres -d {_config.Database} -c \"CREATE EXTENSION IF NOT EXISTS vector\"",
            BinDir, ct);

        if (extResult != 0)
            _logger.LogWarning("pgvector extension not available — vector search features will be disabled");
    }

    // ─── Process helpers ────────────────────────────────────────────

    private async Task<int> RunAsync(string exe, string args, string workDir, CancellationToken ct)
    {
        _logger.LogDebug("Exec: {Exe} {Args}", Path.GetFileName(exe), args);

        var psi = new ProcessStartInfo
        {
            FileName = exe,
            Arguments = args,
            WorkingDirectory = workDir,
            UseShellExecute = false,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            CreateNoWindow = true,
        };

        // Ensure the PG bin dir is on PATH so sub-processes can find each other
        var path = psi.Environment.TryGetValue("PATH", out var existing) ? existing : "";
        psi.Environment["PATH"] = $"{BinDir}{Path.PathSeparator}{path}";

        using var proc = Process.Start(psi) ?? throw new InvalidOperationException($"Failed to start {exe}");
        await proc.WaitForExitAsync(ct);
        return proc.ExitCode;
    }

    private async Task<(int exitCode, string stdout)> RunWithOutputAsync(
        string exe, string args, string workDir, CancellationToken ct)
    {
        var psi = new ProcessStartInfo
        {
            FileName = exe,
            Arguments = args,
            WorkingDirectory = workDir,
            UseShellExecute = false,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            CreateNoWindow = true,
        };

        var path = psi.Environment.TryGetValue("PATH", out var existing) ? existing : "";
        psi.Environment["PATH"] = $"{BinDir}{Path.PathSeparator}{path}";

        using var proc = Process.Start(psi) ?? throw new InvalidOperationException($"Failed to start {exe}");
        var stdout = await proc.StandardOutput.ReadToEndAsync(ct);
        await proc.WaitForExitAsync(ct);
        return (proc.ExitCode, stdout);
    }
}
