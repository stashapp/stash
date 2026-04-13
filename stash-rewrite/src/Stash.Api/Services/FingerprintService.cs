using System.Globalization;
using System.Security.Cryptography;
using Microsoft.EntityFrameworkCore;
using SixLabors.ImageSharp;
using SixLabors.ImageSharp.PixelFormats;
using SixLabors.ImageSharp.Processing;
using Stash.Core.Entities;
using Stash.Core.Interfaces;
using Stash.Data;

namespace Stash.Api.Services;

public interface IFingerprintService
{
    Task<string?> ComputeMd5Async(string path, CancellationToken ct = default);
    Task<string?> ComputeImagePhashAsync(string path, CancellationToken ct = default);
    Task<string?> ComputeVideoPhashAsync(string path, double duration, CancellationToken ct = default);
    string StartGenerateScenePhashes();
    string StartGenerateImagePhashes();
}

public class FingerprintService(
    IServiceScopeFactory scopeFactory,
    IJobService jobService,
    StashConfiguration config,
    ILogger<FingerprintService> logger) : IFingerprintService
{
    private const int HashImageSize = 32;
    private const int LowFrequencySize = 8;

    public async Task<string?> ComputeMd5Async(string path, CancellationToken ct = default)
    {
        if (!File.Exists(path))
            return null;

        await using var stream = new FileStream(path, FileMode.Open, FileAccess.Read, FileShare.Read, 81920, useAsync: true);
        var hash = await MD5.HashDataAsync(stream, ct);
        return Convert.ToHexStringLower(hash);
    }

    public async Task<string?> ComputeImagePhashAsync(string path, CancellationToken ct = default)
    {
        if (!File.Exists(path))
            return null;

        try
        {
            using var image = await SixLabors.ImageSharp.Image.LoadAsync<Rgba32>(path, ct);
            image.Mutate(context => context.Resize(HashImageSize, HashImageSize).Grayscale());

            var values = new double[HashImageSize, HashImageSize];
            for (var y = 0; y < HashImageSize; y++)
            {
                for (var x = 0; x < HashImageSize; x++)
                {
                    values[x, y] = image[x, y].R;
                }
            }

            var dct = ComputeDct(values);
            var lowFrequencyValues = new List<double>(LowFrequencySize * LowFrequencySize - 1);
            for (var y = 0; y < LowFrequencySize; y++)
            {
                for (var x = 0; x < LowFrequencySize; x++)
                {
                    if (x == 0 && y == 0)
                        continue;

                    lowFrequencyValues.Add(dct[x, y]);
                }
            }

            if (lowFrequencyValues.Count == 0)
                return null;

            var average = lowFrequencyValues.Average();
            ulong hash = 0;
            var bitIndex = 0;
            for (var y = 0; y < LowFrequencySize; y++)
            {
                for (var x = 0; x < LowFrequencySize; x++)
                {
                    if (x == 0 && y == 0)
                        continue;

                    if (dct[x, y] > average)
                        hash |= 1UL << bitIndex;

                    bitIndex++;
                }
            }

            return hash.ToString("x16", CultureInfo.InvariantCulture);
        }
        catch (Exception ex)
        {
            logger.LogWarning(ex, "Failed to compute image phash for {Path}", path);
            return null;
        }
    }

    public async Task<string?> ComputeVideoPhashAsync(string path, double duration, CancellationToken ct = default)
    {
        if (!File.Exists(path))
            return null;

        var ffmpegPath = FindFfmpeg();
        if (ffmpegPath == null)
        {
            logger.LogWarning("FFmpeg not found. Cannot compute video phash for {Path}", path);
            return null;
        }

        var tempFile = Path.Combine(Path.GetTempPath(), $"stash-phash-{Guid.NewGuid():N}.jpg");
        var seekSeconds = duration > 0 ? Math.Min(duration * 0.1, 5.0) : 1.0;
        if (seekSeconds <= 0)
            seekSeconds = 1.0;

        try
        {
            using var process = new System.Diagnostics.Process
            {
                StartInfo = new System.Diagnostics.ProcessStartInfo
                {
                    FileName = ffmpegPath,
                    Arguments = $"-ss {seekSeconds.ToString("F2", CultureInfo.InvariantCulture)} -i \"{path}\" -vframes 1 -q:v 4 -vf \"scale=320:-1\" -y \"{tempFile}\"",
                    UseShellExecute = false,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    CreateNoWindow = true,
                }
            };

            process.Start();

            // Read stdout/stderr concurrently to prevent deadlock when buffers fill
            var stderrTask = process.StandardError.ReadToEndAsync(ct);
            var stdoutTask = process.StandardOutput.ReadToEndAsync(ct);

            using var timeoutCts = CancellationTokenSource.CreateLinkedTokenSource(ct);
            timeoutCts.CancelAfter(TimeSpan.FromSeconds(30));

            try
            {
                await process.WaitForExitAsync(timeoutCts.Token);
            }
            catch (OperationCanceledException) when (!ct.IsCancellationRequested)
            {
                logger.LogWarning("FFmpeg timed out extracting frame from {Path}", path);
                try { process.Kill(true); } catch { /* best effort */ }
                return null;
            }

            if (process.ExitCode != 0 || !File.Exists(tempFile))
            {
                var error = await stderrTask;
                logger.LogWarning("FFmpeg failed to extract representative frame for {Path}: {Error}", path, error);
                return null;
            }

            return await ComputeImagePhashAsync(tempFile, ct);
        }
        catch (Exception ex)
        {
            logger.LogWarning(ex, "Failed to compute video phash for {Path}", path);
            return null;
        }
        finally
        {
            try
            {
                if (File.Exists(tempFile))
                    File.Delete(tempFile);
            }
            catch
            {
                // Best-effort cleanup only.
            }
        }
    }

    public string StartGenerateScenePhashes()
    {
        return jobService.Enqueue("generate_scene_phashes", "Generating scene pHashes", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var db = scope.ServiceProvider.GetRequiredService<StashContext>();

            // Get IDs of files that already have a phash
            var filesWithPhash = await db.FileFingerprints
                .Where(fp => fp.Type == "phash")
                .Select(fp => fp.FileId)
                .Distinct()
                .ToHashSetAsync(ct);

            // Only load files that need phash generation
            var files = await db.VideoFiles
                .Include(file => file.ParentFolder)
                .Include(file => file.Fingerprints)
                .Where(file => !filesWithPhash.Contains(file.Id))
                .OrderBy(file => file.Id)
                .ToListAsync(ct);

            if (files.Count == 0)
            {
                progress.Report(1.0, "All scenes already have pHashes");
                return;
            }

            logger.LogInformation("Generating pHashes for {Count} video files", files.Count);
            var saveCounter = 0;
            for (var index = 0; index < files.Count; index++)
            {
                ct.ThrowIfCancellationRequested();
                var file = files[index];
                progress.Report((double)(index + 1) / files.Count, $"({index + 1}/{files.Count}) {file.Basename}");
                await EnsureVideoPhashAsync(db, file, ct);
                saveCounter++;

                if (saveCounter >= 50)
                {
                    await db.SaveChangesAsync(ct);
                    saveCounter = 0;
                }
            }

            if (saveCounter > 0)
                await db.SaveChangesAsync(ct);

            logger.LogInformation("Finished generating pHashes for {Count} video files", files.Count);
        });
    }

    public string StartGenerateImagePhashes()
    {
        return jobService.Enqueue("generate_image_phashes", "Generating image pHashes", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var db = scope.ServiceProvider.GetRequiredService<StashContext>();

            var files = await db.ImageFiles
                .Include(file => file.ParentFolder)
                .Include(file => file.Fingerprints)
                .OrderBy(file => file.Id)
                .ToListAsync(ct);

            if (files.Count == 0)
                return;

            var saveCounter = 0;
            for (var index = 0; index < files.Count; index++)
            {
                ct.ThrowIfCancellationRequested();
                var file = files[index];
                progress.Report((double)(index + 1) / files.Count, file.Basename);
                await EnsureImagePhashAsync(db, file, ct);
                saveCounter++;

                if (saveCounter >= 100)
                {
                    await db.SaveChangesAsync(ct);
                    saveCounter = 0;
                }
            }

            if (saveCounter > 0)
                await db.SaveChangesAsync(ct);
        });
    }

    private async Task EnsureVideoPhashAsync(StashContext db, VideoFile file, CancellationToken ct)
    {
        if (file.Fingerprints.Any(fp => fp.Type == "phash" && !string.IsNullOrWhiteSpace(fp.Value)))
            return;

        var path = ResolveFilePath(file);
        if (path == null)
            return;

        var oshash = file.Fingerprints.FirstOrDefault(fp => fp.Type == "oshash")?.Value;
        if (!string.IsNullOrWhiteSpace(oshash))
        {
            var reused = await FindExistingPhashAsync(db, file.Id, "oshash", oshash, ct);
            if (!string.IsNullOrWhiteSpace(reused))
            {
                AddFingerprint(file, "phash", reused);
                return;
            }
        }

        var phash = await ComputeVideoPhashAsync(path, file.Duration, ct);
        if (!string.IsNullOrWhiteSpace(phash))
            AddFingerprint(file, "phash", phash);
    }

    private async Task EnsureImagePhashAsync(StashContext db, ImageFile file, CancellationToken ct)
    {
        if (file.Fingerprints.Any(fp => fp.Type == "phash" && !string.IsNullOrWhiteSpace(fp.Value)))
            return;

        var path = ResolveFilePath(file);
        if (path == null)
            return;

        var md5 = file.Fingerprints.FirstOrDefault(fp => fp.Type == "md5")?.Value;
        if (string.IsNullOrWhiteSpace(md5))
        {
            md5 = await ComputeMd5Async(path, ct);
            if (!string.IsNullOrWhiteSpace(md5))
                AddFingerprint(file, "md5", md5);
        }

        if (!string.IsNullOrWhiteSpace(md5))
        {
            var reused = await FindExistingPhashAsync(db, file.Id, "md5", md5, ct);
            if (!string.IsNullOrWhiteSpace(reused))
            {
                AddFingerprint(file, "phash", reused);
                return;
            }
        }

        var phash = await ComputeImagePhashAsync(path, ct);
        if (!string.IsNullOrWhiteSpace(phash))
            AddFingerprint(file, "phash", phash);
    }

    private static async Task<string?> FindExistingPhashAsync(StashContext db, int fileId, string sourceType, string sourceValue, CancellationToken ct)
    {
        return await db.FileFingerprints
            .Where(fp => fp.Type == sourceType && fp.Value == sourceValue && fp.FileId != fileId)
            .Join(
                db.FileFingerprints.Where(fp => fp.Type == "phash"),
                source => source.FileId,
                phash => phash.FileId,
                (_, phash) => phash.Value)
            .AsNoTracking()
            .FirstOrDefaultAsync(ct);
    }

    private static void AddFingerprint(BaseFileEntity file, string type, string value)
    {
        if (file.Fingerprints.Any(fp => fp.Type == type && string.Equals(fp.Value, value, StringComparison.OrdinalIgnoreCase)))
            return;

        file.Fingerprints.Add(new FileFingerprint
        {
            Type = type,
            Value = value,
            FileId = file.Id,
        });
    }

    private static string? ResolveFilePath(BaseFileEntity file)
    {
        var path = file.ParentFolder != null
            ? Path.Combine(file.ParentFolder.Path, file.Basename)
            : file.Basename;

        return File.Exists(path) ? path : null;
    }

    private string? FindFfmpeg()
    {
        if (!string.IsNullOrWhiteSpace(config.FfmpegPath) && File.Exists(config.FfmpegPath))
            return config.FfmpegPath;

        var pathDirectories = Environment.GetEnvironmentVariable("PATH")?.Split(Path.PathSeparator, StringSplitOptions.RemoveEmptyEntries) ?? [];
        foreach (var directory in pathDirectories)
        {
            var ffmpegPath = Path.Combine(directory, OperatingSystem.IsWindows() ? "ffmpeg.exe" : "ffmpeg");
            if (File.Exists(ffmpegPath))
                return ffmpegPath;
        }

        return null;
    }

    private static double[,] ComputeDct(double[,] values)
    {
        var result = new double[HashImageSize, HashImageSize];
        for (var u = 0; u < HashImageSize; u++)
        {
            for (var v = 0; v < HashImageSize; v++)
            {
                var sum = 0d;
                for (var x = 0; x < HashImageSize; x++)
                {
                    for (var y = 0; y < HashImageSize; y++)
                    {
                        sum += values[x, y]
                            * Math.Cos(((2 * x + 1) * u * Math.PI) / (2 * HashImageSize))
                            * Math.Cos(((2 * y + 1) * v * Math.PI) / (2 * HashImageSize));
                    }
                }

                result[u, v] = 0.25 * Alpha(u) * Alpha(v) * sum;
            }
        }

        return result;
    }

    private static double Alpha(int index) => index == 0 ? 1d / Math.Sqrt(2d) : 1d;
}