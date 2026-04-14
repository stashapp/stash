using System.Security.Cryptography;
using System.Text;
using Microsoft.EntityFrameworkCore;
using Stash.Core.Entities.Galleries.Zip;
using Stash.Core.Interfaces;
using Stash.Data;

namespace Stash.Api.Services;

public interface IThumbnailService
{
    Task<string?> GetSceneThumbnailPathAsync(int sceneId, CancellationToken ct = default);
    Task<string?> GetImageFilePathAsync(int imageId, CancellationToken ct = default);
    Task<(Stream stream, string contentType, bool supportsRangeRequests)?> GetImageStreamAsync(int imageId, CancellationToken ct = default);
    Task GenerateSceneThumbnailAsync(int sceneId, double? atSeconds = null, CancellationToken ct = default);
    Task GenerateScenePreviewAsync(int sceneId, CancellationToken ct = default);
    Task GenerateSceneSpriteAsync(int sceneId, CancellationToken ct = default);
    string GetThumbnailPathForScene(int sceneId);
    string GetTimestampedThumbnailPath(int sceneId, double seconds);
    string GetPreviewPath(int sceneId);
    string GetSpritePath(int sceneId);
    string GetSpriteVttPath(int sceneId);
    string StartGenerateAllThumbnails();
}

public class ThumbnailService(
    IServiceScopeFactory scopeFactory,
    IJobService jobService,
    StashConfiguration config,
    IZipFileReader zipFileReader,
    ILogger<ThumbnailService> logger) : IThumbnailService
{
    private string ThumbnailDir => Path.Combine(config.GeneratedPath, "screenshots");
    private string PreviewDir => Path.Combine(config.GeneratedPath, "previews");
    private string VttDir => Path.Combine(config.GeneratedPath, "vtt");
    private static readonly SemaphoreSlim _ffmpegSemaphore = new(4, 4); // max concurrent FFmpeg processes
    private string? _cachedFfmpegPath;
    private bool _ffmpegSearched;
    private string? _hwEncoder;
    private bool _hwEncoderSearched;

    private static readonly Dictionary<string, string> ImageMimeTypes = new(StringComparer.OrdinalIgnoreCase)
    {
        [".jpg"] = "image/jpeg", [".jpeg"] = "image/jpeg", [".png"] = "image/png",
        [".gif"] = "image/gif", [".webp"] = "image/webp", [".bmp"] = "image/bmp",
        [".tiff"] = "image/tiff", [".tif"] = "image/tiff", [".svg"] = "image/svg+xml",
    };

    // Preview generation defaults (matching original Stash)
    private const int PreviewSegments = 12;
    private const double PreviewSegmentDuration = 0.75;
    private const int PreviewWidth = 640;
    private const string PreviewPreset = "fast";
    private const int PreviewCrf = 21;

    // Sprite generation defaults
    private const int SpriteFrameCount = 81; // 9x9 grid
    private const int SpriteFrameSize = 160; // px

    public async Task<string?> GetSceneThumbnailPathAsync(int sceneId, CancellationToken ct)
    {
        var thumbPath = GetThumbnailPath(sceneId);
        if (File.Exists(thumbPath)) return thumbPath;

        // Try to generate on-demand
        await GenerateSceneThumbnailAsync(sceneId, null, ct);
        return File.Exists(thumbPath) ? thumbPath : null;
    }

    public async Task<string?> GetImageFilePathAsync(int imageId, CancellationToken ct)
    {
        using var scope = scopeFactory.CreateScope();
        var db = scope.ServiceProvider.GetRequiredService<StashContext>();

        var imageFile = await db.ImageFiles
            .Include(f => f.ParentFolder)
            .FirstOrDefaultAsync(f => f.ImageId == imageId, ct);

        if (imageFile == null) return null;

        var filePath = imageFile.ParentFolder != null
            ? Path.Combine(imageFile.ParentFolder.Path, imageFile.Basename)
            : imageFile.Basename;

        return File.Exists(filePath) ? filePath : null;
    }

    public async Task<(Stream stream, string contentType, bool supportsRangeRequests)?> GetImageStreamAsync(int imageId, CancellationToken ct)
    {
        using var scope = scopeFactory.CreateScope();
        var db = scope.ServiceProvider.GetRequiredService<StashContext>();

        var imageFile = await db.ImageFiles
            .Include(f => f.ParentFolder)
            .FirstOrDefaultAsync(f => f.ImageId == imageId, ct);

        if (imageFile == null) return null;

        // Check if this image is stored in a zip archive
        if (imageFile.ZipFileId.HasValue)
        {
            // Get the parent zip file
            var zipFile = await db.GalleryFiles
                .Include(gf => gf.ParentFolder)
                .FirstOrDefaultAsync(gf => gf.Id == imageFile.ZipFileId.Value, ct);

            if (zipFile == null) return null;

            var zipFilePath = Path.Combine(zipFile.ParentFolder?.Path ?? "", zipFile.Basename);

            // Extract the image from the zip
            var stream = await zipFileReader.ExtractEntryAsync(zipFilePath, imageFile.Basename, ct);
            var ext = Path.GetExtension(imageFile.Basename);
            var contentType = ImageMimeTypes.GetValueOrDefault(ext, "application/octet-stream");

            // Zip-based streams don't support range requests (they're already in memory)
            return (stream, contentType, false);
        }
        else
        {
            // Regular file-based image
            var filePath = imageFile.ParentFolder != null
                ? Path.Combine(imageFile.ParentFolder.Path, imageFile.Basename)
                : imageFile.Basename;

            if (!File.Exists(filePath)) return null;

            var ext = Path.GetExtension(filePath);
            var contentType = ImageMimeTypes.GetValueOrDefault(ext, "application/octet-stream");
            var stream = new FileStream(filePath, FileMode.Open, FileAccess.Read, FileShare.Read, 81920, useAsync: true);

            // File-based streams support range requests
            return (stream, contentType, true);
        }
    }

    public async Task GenerateSceneThumbnailAsync(int sceneId, double? atSeconds, CancellationToken ct)
    {
        // Check if the target thumbnail already exists (avoid re-generation)
        var thumbPath = atSeconds.HasValue
            ? GetTimestampedThumbnailPath(sceneId, atSeconds.Value)
            : GetThumbnailPath(sceneId);

        if (File.Exists(thumbPath)) return;

        using var scope = scopeFactory.CreateScope();
        var db = scope.ServiceProvider.GetRequiredService<StashContext>();

        var videoFile = await db.VideoFiles
            .Include(f => f.ParentFolder)
            .AsNoTracking()
            .FirstOrDefaultAsync(f => f.SceneId == sceneId, ct);

        if (videoFile == null) return;

        var filePath = videoFile.ParentFolder != null
            ? Path.Combine(videoFile.ParentFolder.Path, videoFile.Basename)
            : videoFile.Basename;

        if (!File.Exists(filePath)) return;

        var ffmpegPath = GetCachedFfmpegPath();
        if (ffmpegPath == null)
        {
            logger.LogWarning("FFmpeg not found. Cannot generate thumbnail for scene {SceneId}", sceneId);
            return;
        }

        var thumbDir = Path.GetDirectoryName(thumbPath)!;
        Directory.CreateDirectory(thumbDir);

        var seekSeconds = atSeconds ?? videoFile.Duration * 0.2;
        if (seekSeconds <= 0) seekSeconds = 1;

        // Limit concurrent FFmpeg processes
        await _ffmpegSemaphore.WaitAsync(ct);
        try
        {
            // Double-check after acquiring semaphore (another request may have generated it)
            if (File.Exists(thumbPath)) return;

            var process = new System.Diagnostics.Process
            {
                StartInfo = new System.Diagnostics.ProcessStartInfo
                {
                    FileName = ffmpegPath,
                    Arguments = $"-v error -ss {seekSeconds:F2} -i \"{filePath}\" -vframes 1 -q:v 2 -f image2 -y \"{thumbPath}\"",
                    UseShellExecute = false,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    CreateNoWindow = true
                }
            };

            process.Start();
            // Drain stderr to prevent buffer deadlocks
            var stderrTask = process.StandardError.ReadToEndAsync(ct);
            using var cts = CancellationTokenSource.CreateLinkedTokenSource(ct);
            cts.CancelAfter(TimeSpan.FromSeconds(15)); // timeout per thumbnail
            try
            {
                await process.WaitForExitAsync(cts.Token);
            }
            catch (OperationCanceledException) when (!ct.IsCancellationRequested)
            {
                try { process.Kill(entireProcessTree: true); } catch { }
                logger.LogWarning("FFmpeg timed out for scene {SceneId} at {Seconds}s", sceneId, seekSeconds);
                return;
            }

            if (process.ExitCode != 0)
            {
                var stderr = await stderrTask;
                logger.LogWarning("FFmpeg failed for scene {SceneId}: {Error}", sceneId, stderr);
            }
        }
        catch (Exception ex) when (ex is not OperationCanceledException)
        {
            logger.LogError(ex, "Error generating thumbnail for scene {SceneId}", sceneId);
        }
        finally
        {
            _ffmpegSemaphore.Release();
        }
    }

    /// <summary>Get the path for a timestamp-specific cached thumbnail.</summary>
    public string GetTimestampedThumbnailPath(int sceneId, double seconds)
    {
        var hash = Convert.ToHexStringLower(SHA256.HashData(BitConverter.GetBytes(sceneId)));
        var subDir = hash[..2];
        var secKey = ((int)seconds).ToString();
        return Path.Combine(ThumbnailDir, subDir, $"{sceneId}_t{secKey}.jpg");
    }

    public string GetThumbnailPathForScene(int sceneId) => GetThumbnailPath(sceneId);

    public string GetPreviewPath(int sceneId)
    {
        var hash = Convert.ToHexStringLower(SHA256.HashData(BitConverter.GetBytes(sceneId)));
        return Path.Combine(PreviewDir, hash[..2], $"{sceneId}.mp4");
    }

    public string GetSpritePath(int sceneId)
    {
        var hash = Convert.ToHexStringLower(SHA256.HashData(BitConverter.GetBytes(sceneId)));
        return Path.Combine(VttDir, hash[..2], $"{sceneId}_sprite.jpg");
    }

    public string GetSpriteVttPath(int sceneId)
    {
        var hash = Convert.ToHexStringLower(SHA256.HashData(BitConverter.GetBytes(sceneId)));
        return Path.Combine(VttDir, hash[..2], $"{sceneId}_thumbs.vtt");
    }

    /// <summary>Generate a multi-segment video preview clip (mp4) for a scene.</summary>
    public async Task GenerateScenePreviewAsync(int sceneId, CancellationToken ct)
    {
        var previewPath = GetPreviewPath(sceneId);
        if (File.Exists(previewPath)) return;

        var (filePath, duration) = await GetSceneFileInfoAsync(sceneId, ct);
        if (filePath == null || duration <= 0) return;

        var ffmpegPath = GetCachedFfmpegPath();
        if (ffmpegPath == null)
        {
            logger.LogWarning("FFmpeg not found, cannot generate preview for scene {SceneId}", sceneId);
            return;
        }

        var previewDir = Path.GetDirectoryName(previewPath)!;
        Directory.CreateDirectory(previewDir);

        var tmpDir = Path.Combine(config.GeneratedPath, "tmp", $"preview_{sceneId}");
        Directory.CreateDirectory(tmpDir);

        try
        {
            var segmentCount = PreviewSegments;
            var segmentDuration = PreviewSegmentDuration;

            // If video is too short for all segments, use a single full-video preview
            if (duration < segmentDuration * segmentCount)
            {
                await RunFfmpegAsync(ffmpegPath,
                    $"-v error -y -i \"{filePath}\" -max_muxing_queue_size 1024 -c:v {GetH264Encoder()} -vf \"scale={PreviewWidth}:-2\" -pix_fmt yuv420p -profile:v high -level 4.2 -preset {PreviewPreset} -crf {PreviewCrf} -threads 4 -an \"{previewPath}\"",
                    TimeSpan.FromMinutes(5), ct);
                return;
            }

            var interval = duration / segmentCount;
            var chunkFiles = new List<string>();

            for (int i = 0; i < segmentCount; i++)
            {
                ct.ThrowIfCancellationRequested();
                var seekTime = interval * i + interval * 0.5;
                if (seekTime >= duration) seekTime = duration - segmentDuration;
                if (seekTime < 0) seekTime = 0;

                var chunkPath = Path.Combine(tmpDir, $"chunk_{i:D3}.mp4");
                chunkFiles.Add(chunkPath);

                await RunFfmpegAsync(ffmpegPath,
                    $"-v error -y -ss {seekTime:F2} -i \"{filePath}\" -t {segmentDuration:F2} -max_muxing_queue_size 1024 -c:v {GetH264Encoder()} -vf \"scale={PreviewWidth}:-2\" -pix_fmt yuv420p -profile:v high -level 4.2 -preset {PreviewPreset} -crf {PreviewCrf} -threads 4 -an \"{chunkPath}\"",
                    TimeSpan.FromSeconds(60), ct);
            }

            // Create concat file — use forward slashes for FFmpeg compatibility on all platforms
            var concatListPath = Path.Combine(tmpDir, "concat.txt");
            var concatLines = chunkFiles
                .Where(File.Exists)
                .Select(f => $"file '{Path.GetFullPath(f).Replace('\\', '/')}'");
            await File.WriteAllTextAsync(concatListPath, string.Join("\n", concatLines), ct);

            // Concatenate chunks into final preview
            await RunFfmpegAsync(ffmpegPath,
                $"-v error -y -f concat -safe 0 -i \"{concatListPath}\" -c:v copy \"{previewPath}\"",
                TimeSpan.FromSeconds(30), ct);

            if (!File.Exists(previewPath))
                logger.LogWarning("Preview generation failed for scene {SceneId} - output not created", sceneId);
        }
        catch (Exception ex) when (ex is not OperationCanceledException)
        {
            logger.LogError(ex, "Error generating preview for scene {SceneId}", sceneId);
        }
        finally
        {
            try { if (Directory.Exists(tmpDir)) Directory.Delete(tmpDir, true); } catch { }
        }
    }

    /// <summary>Generate a sprite sheet (JPEG grid) and VTT timeline file for a scene.</summary>
    public async Task GenerateSceneSpriteAsync(int sceneId, CancellationToken ct)
    {
        var spritePath = GetSpritePath(sceneId);
        var vttPath = GetSpriteVttPath(sceneId);
        if (File.Exists(spritePath) && File.Exists(vttPath)) return;

        var (filePath, duration) = await GetSceneFileInfoAsync(sceneId, ct);
        if (filePath == null || duration <= 0) return;

        var ffmpegPath = GetCachedFfmpegPath();
        if (ffmpegPath == null) return;

        var spriteDir = Path.GetDirectoryName(spritePath)!;
        Directory.CreateDirectory(spriteDir);

        try
        {
            // Calculate grid dimensions
            var frameCount = Math.Min(SpriteFrameCount, Math.Max(1, (int)(duration / 2)));
            var cols = (int)Math.Ceiling(Math.Sqrt(frameCount));
            var rows = (int)Math.Ceiling((double)frameCount / cols);
            var interval = duration / frameCount;

            // Use FFmpeg to generate sprite sheet in one command using fps + tile filters
            await _ffmpegSemaphore.WaitAsync(ct);
            try
            {
                // fps=1/interval extracts frames at the right interval, scale to sprite size, tile into grid
                var fpsRate = 1.0 / interval;
                await RunFfmpegAsync(ffmpegPath,
                    $"-v error -y -i \"{filePath}\" -vf \"fps={fpsRate:F6},scale={SpriteFrameSize}:-2,tile={cols}x{rows}\" -frames:v 1 -q:v 5 \"{spritePath}\"",
                    TimeSpan.FromMinutes(5), ct);
            }
            finally
            {
                _ffmpegSemaphore.Release();
            }

            if (!File.Exists(spritePath))
            {
                logger.LogWarning("Sprite generation failed for scene {SceneId}", sceneId);
                return;
            }

            // Calculate actual frame dimensions from the sprite image
            var thumbWidth = SpriteFrameSize;
            // Estimate height from typical 16:9 aspect ratio; FFmpeg -2 ensures even
            var thumbHeight = (int)Math.Round(thumbWidth * 9.0 / 16.0);
            if (thumbHeight % 2 != 0) thumbHeight++;

            // Generate VTT file
            var vttBuilder = new StringBuilder();
            vttBuilder.AppendLine("WEBVTT");
            vttBuilder.AppendLine();

            var spriteFileName = Path.GetFileName(spritePath);
            for (int i = 0; i < frameCount; i++)
            {
                var startTime = i * interval;
                var endTime = Math.Min((i + 1) * interval, duration);
                var col = i % cols;
                var row = i / cols;
                var x = col * thumbWidth;
                var y = row * thumbHeight;

                vttBuilder.AppendLine($"{FormatVttTime(startTime)} --> {FormatVttTime(endTime)}");
                vttBuilder.AppendLine($"{spriteFileName}#xywh={x},{y},{thumbWidth},{thumbHeight}");
                vttBuilder.AppendLine();
            }

            await File.WriteAllTextAsync(vttPath, vttBuilder.ToString(), ct);
        }
        catch (Exception ex) when (ex is not OperationCanceledException)
        {
            logger.LogError(ex, "Error generating sprite for scene {SceneId}", sceneId);
        }
    }

    private static string FormatVttTime(double seconds)
    {
        var ts = TimeSpan.FromSeconds(seconds);
        return $"{(int)ts.TotalHours:D2}:{ts.Minutes:D2}:{ts.Seconds:D2}.{ts.Milliseconds:D3}";
    }

    private async Task<(string? filePath, double duration)> GetSceneFileInfoAsync(int sceneId, CancellationToken ct)
    {
        using var scope = scopeFactory.CreateScope();
        var db = scope.ServiceProvider.GetRequiredService<StashContext>();

        var videoFile = await db.VideoFiles
            .Include(f => f.ParentFolder)
            .AsNoTracking()
            .FirstOrDefaultAsync(f => f.SceneId == sceneId, ct);

        if (videoFile == null) return (null, 0);

        var filePath = videoFile.ParentFolder != null
            ? Path.Combine(videoFile.ParentFolder.Path, videoFile.Basename)
            : videoFile.Basename;

        return File.Exists(filePath) ? (filePath, videoFile.Duration) : (null, 0);
    }

    private async Task RunFfmpegAsync(string ffmpegPath, string args, TimeSpan timeout, CancellationToken ct)
    {
        var process = new System.Diagnostics.Process
        {
            StartInfo = new System.Diagnostics.ProcessStartInfo
            {
                FileName = ffmpegPath,
                Arguments = args,
                UseShellExecute = false,
                RedirectStandardOutput = true,
                RedirectStandardError = true,
                CreateNoWindow = true
            }
        };

        process.Start();
        var stderrTask = process.StandardError.ReadToEndAsync(ct);
        using var cts = CancellationTokenSource.CreateLinkedTokenSource(ct);
        cts.CancelAfter(timeout);
        try
        {
            await process.WaitForExitAsync(cts.Token);
        }
        catch (OperationCanceledException) when (!ct.IsCancellationRequested)
        {
            try { process.Kill(entireProcessTree: true); } catch { }
            logger.LogWarning("FFmpeg timed out: {Args}", args[..Math.Min(200, args.Length)]);
            return;
        }

        if (process.ExitCode != 0)
        {
            var stderr = await stderrTask;
            logger.LogWarning("FFmpeg failed (exit {Code}): {Error}", process.ExitCode, stderr[..Math.Min(500, stderr.Length)]);
        }
    }

    public string StartGenerateAllThumbnails()
    {
        return jobService.Enqueue("generate_thumbnails", "Generating thumbnails", async (progress, ct) =>
        {
            using var scope = scopeFactory.CreateScope();
            var db = scope.ServiceProvider.GetRequiredService<StashContext>();

            var sceneIds = await db.Scenes.Select(s => s.Id).ToListAsync(ct);
            var total = sceneIds.Count;
            var processed = 0;

            foreach (var sceneId in sceneIds)
            {
                ct.ThrowIfCancellationRequested();
                processed++;
                progress.Report((double)processed / total, $"Scene {processed}/{total}");

                var thumbPath = GetThumbnailPath(sceneId);
                if (File.Exists(thumbPath)) continue;

                await GenerateSceneThumbnailAsync(sceneId, null, ct);
            }

            logger.LogInformation("Generated thumbnails for {Count} scenes", total);
        });
    }

    private string GetThumbnailPath(int sceneId)
    {
        var hash = Convert.ToHexStringLower(SHA256.HashData(BitConverter.GetBytes(sceneId)));
        var subDir = hash[..2];
        return Path.Combine(ThumbnailDir, subDir, $"{sceneId}.jpg");
    }

    private string? GetCachedFfmpegPath()
    {
        if (_ffmpegSearched) return _cachedFfmpegPath;
        _cachedFfmpegPath = FindFfmpeg();
        _ffmpegSearched = true;
        return _cachedFfmpegPath;
    }

    private string? FindFfmpeg()
    {
        if (!string.IsNullOrEmpty(config.FfmpegPath) && File.Exists(config.FfmpegPath))
            return config.FfmpegPath;

        // Search PATH
        var pathDirs = Environment.GetEnvironmentVariable("PATH")?.Split(Path.PathSeparator) ?? [];
        foreach (var dir in pathDirs)
        {
            var ffmpeg = Path.Combine(dir, OperatingSystem.IsWindows() ? "ffmpeg.exe" : "ffmpeg");
            if (File.Exists(ffmpeg)) return ffmpeg;
        }

        return null;
    }

    /// <summary>Get the best available H.264 encoder, preferring HW-accelerated encoders.</summary>
    private string GetH264Encoder()
    {
        if (_hwEncoderSearched) return _hwEncoder ?? "libx264";
        _hwEncoderSearched = true;

        var ffmpegPath = GetCachedFfmpegPath();
        if (ffmpegPath == null) return "libx264";

        try
        {
            var process = new System.Diagnostics.Process
            {
                StartInfo = new System.Diagnostics.ProcessStartInfo
                {
                    FileName = ffmpegPath,
                    Arguments = "-hide_banner -encoders",
                    UseShellExecute = false,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    CreateNoWindow = true
                }
            };
            process.Start();
            var output = process.StandardOutput.ReadToEnd();
            process.WaitForExit(5000);

            // Prefer NVENC > QSV > AMF > libx264
            string[] hwEncoders = ["h264_nvenc", "h264_qsv", "h264_amf", "h264_videotoolbox"];
            foreach (var enc in hwEncoders)
            {
                if (output.Contains(enc, StringComparison.OrdinalIgnoreCase))
                {
                    _hwEncoder = enc;
                    logger.LogInformation("Using HW-accelerated H.264 encoder: {Encoder}", enc);
                    return enc;
                }
            }
        }
        catch (Exception ex)
        {
            logger.LogWarning(ex, "Failed to detect HW encoders, falling back to libx264");
        }

        return "libx264";
    }
}
