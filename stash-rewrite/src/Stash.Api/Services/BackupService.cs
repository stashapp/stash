using Microsoft.Extensions.Logging;
using Stash.Core.Interfaces;

namespace Stash.Api.Services;

public class BackupService(
    IJobService jobService,
    StashConfiguration config,
    ILogger<BackupService> logger) : IBackupService
{
    private static string BackupDir => StashDefaultPaths.GetDataSubdirectory("backups");

    public string StartBackup()
    {
        return jobService.Enqueue("backup", "Backing up database", async (progress, ct) =>
        {
            var backupDir = BackupDir;
            Directory.CreateDirectory(backupDir);

            var timestamp = DateTime.UtcNow.ToString("yyyyMMdd_HHmmss");
            var backupPath = Path.Combine(backupDir, $"stash_backup_{timestamp}.sql");

            progress.Report(0.1, "Starting backup...");

            // Use pg_dump if available, otherwise do a simple schema+data export marker
            var connStr = config.DatabaseConnectionString;
            // Parse connection string for pg_dump args
            var parts = new Dictionary<string, string>(StringComparer.OrdinalIgnoreCase);
            foreach (var part in connStr.Split(';', StringSplitOptions.RemoveEmptyEntries))
            {
                var kv = part.Split('=', 2);
                if (kv.Length == 2)
                    parts[kv[0].Trim()] = kv[1].Trim();
            }

            parts.TryGetValue("Host", out var host);
            parts.TryGetValue("Port", out var port);
            parts.TryGetValue("Database", out var database);
            parts.TryGetValue("Username", out var username);
            parts.TryGetValue("Password", out var password);

            host ??= "localhost";
            port ??= "5432";
            database ??= "stash";
            username ??= "stash";

            try
            {
                var psi = new System.Diagnostics.ProcessStartInfo
                {
                    FileName = "pg_dump",
                    Arguments = $"-h {host} -p {port} -U {username} -d {database} -F c -f \"{backupPath}\"",
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    UseShellExecute = false,
                    CreateNoWindow = true,
                };

                if (!string.IsNullOrEmpty(password))
                    psi.Environment["PGPASSWORD"] = password;

                progress.Report(0.3, "Running pg_dump...");

                using var proc = System.Diagnostics.Process.Start(psi);
                if (proc == null)
                    throw new InvalidOperationException("Failed to start pg_dump");

                var stderr = await proc.StandardError.ReadToEndAsync(ct);
                await proc.WaitForExitAsync(ct);

                if (proc.ExitCode != 0)
                {
                    logger.LogError("pg_dump failed: {Error}", stderr);
                    throw new InvalidOperationException($"pg_dump failed: {stderr}");
                }

                progress.Report(1.0, "Backup complete");
                logger.LogInformation("Database backed up to {Path}", backupPath);
            }
            catch (Exception ex) when (ex is not OperationCanceledException)
            {
                logger.LogError(ex, "Backup failed");
                throw;
            }
        });
    }

    public Task<string?> GetLatestBackupPathAsync(CancellationToken ct = default)
    {
        var backupDir = BackupDir;
        if (!Directory.Exists(backupDir))
            return Task.FromResult<string?>(null);

        var latest = Directory.GetFiles(backupDir, "stash_backup_*")
            .OrderByDescending(f => f)
            .FirstOrDefault();

        return Task.FromResult(latest);
    }
}
