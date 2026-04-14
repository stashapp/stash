using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Stash.Core.Entities;
using Stash.Core.Entities.Galleries.Zip;
using Stash.Core.Events;
using Stash.Core.Interfaces;
using Stash.Data;

namespace Stash.Api.Services;

public class ScanService(
    IJobService jobService,
    IServiceScopeFactory scopeFactory,
    StashConfiguration config,
    IEventBus eventBus,
    IFingerprintService fingerprintService,
    ZipGalleryReader zipGalleryReader,
    ILogger<ScanService> logger) : IScanService
{
    public string StartScan(bool scanGenerators = false)
    {
        return jobService.Enqueue("scan", "Scanning library", async (progress, ct) =>
        {
            var cfg = config;
            var paths = cfg.StashPaths.Select(p => p.Path).ToList();

            if (paths.Count == 0)
            {
                logger.LogWarning("No stash paths configured. Nothing to scan.");
                return;
            }

            var videoExts = new HashSet<string>(cfg.VideoExtensions, StringComparer.OrdinalIgnoreCase);
            var imageExts = new HashSet<string>(cfg.ImageExtensions, StringComparer.OrdinalIgnoreCase);
            var galleryExts = new HashSet<string>(cfg.GalleryExtensions, StringComparer.OrdinalIgnoreCase);
            var allExts = videoExts.Union(imageExts).Union(galleryExts).ToHashSet(StringComparer.OrdinalIgnoreCase);

            // Phase 1: Discover files
            progress.Report(0, "Discovering files...");
            var files = new List<DiscoveredFile>();
            foreach (var stashPath in cfg.StashPaths)
            {
                if (!Directory.Exists(stashPath.Path))
                {
                    logger.LogWarning("Stash path does not exist: {Path}", stashPath.Path);
                    continue;
                }

                var dirFiles = Directory.EnumerateFiles(stashPath.Path, "*", SearchOption.AllDirectories)
                    .Where(f =>
                    {
                        var ext = Path.GetExtension(f);
                        if (!allExts.Contains(ext)) return false;
                        if (stashPath.ExcludeVideo && videoExts.Contains(ext)) return false;
                        if (stashPath.ExcludeImage && imageExts.Contains(ext)) return false;
                        return !IsExcluded(f, cfg.ExcludePatterns);
                    })
                    .Select(f => new DiscoveredFile(f, Path.GetExtension(f)));

                files.AddRange(dirFiles);
            }

            logger.LogInformation("Discovered {Count} files to scan", files.Count);
            if (files.Count == 0) return;

            // Phase 2: Process files
            using var scope = scopeFactory.CreateScope();
            var db = scope.ServiceProvider.GetRequiredService<StashContext>();

            var processedCount = 0;
            foreach (var file in files)
            {
                ct.ThrowIfCancellationRequested();
                processedCount++;
                progress.Report((double)processedCount / files.Count, Path.GetFileName(file.Path));

                try
                {
                    // Check if file already exists in DB by path
                    var existingFolder = await db.Folders
                        .FirstOrDefaultAsync(f => f.Path == Path.GetDirectoryName(file.Path), ct);

                    if (existingFolder != null)
                    {
                        var basename = Path.GetFileName(file.Path);
                        var existingFile = await db.Set<BaseFileEntity>()
                            .FirstOrDefaultAsync(f => f.ParentFolderId == existingFolder.Id && f.Basename == basename, ct);

                        if (existingFile != null)
                        {
                            // Check if file has been modified — but always re-process videos with missing metadata
                            var fileInfo = new FileInfo(file.Path);
                            var needsMetadata = existingFile is VideoFile vf && vf.Width == 0 && vf.Height == 0 && vf.Duration == 0;
                            if (!needsMetadata && existingFile.ModTime >= fileInfo.LastWriteTimeUtc && existingFile.Size == fileInfo.Length)
                                continue; // Not modified and metadata present, skip
                        }
                    }

                    // Process the file
                    if (videoExts.Contains(file.Extension))
                        await ProcessVideoFileAsync(db, file.Path, ct);
                    else if (imageExts.Contains(file.Extension))
                        await ProcessImageFileAsync(db, file.Path, ct);
                    else if (galleryExts.Contains(file.Extension))
                        await ProcessGalleryFileAsync(db, file.Path, ct);
                }
                catch (Exception ex)
                {
                    logger.LogError(ex, "Error processing file: {Path}", file.Path);
                }
            }

            await db.SaveChangesAsync(ct);

            logger.LogInformation("Scan completed. Processed {Count} files", processedCount);
            eventBus.Publish(new StashEvent(EventType.ScanCompleted));
        });
    }

    private async Task<Folder> EnsureFolderAsync(StashContext db, string dirPath, CancellationToken ct)
    {
        var folder = await db.Folders.FirstOrDefaultAsync(f => f.Path == dirPath, ct);
        if (folder != null) return folder;

        folder = new Folder
        {
            Path = dirPath,
            ModTime = Directory.GetLastWriteTimeUtc(dirPath)
        };

        // Link parent folder if path has a parent
        var parentDir = Path.GetDirectoryName(dirPath);
        if (!string.IsNullOrEmpty(parentDir) && parentDir != dirPath)
        {
            var parentFolder = await db.Folders.FirstOrDefaultAsync(f => f.Path == parentDir, ct);
            if (parentFolder != null)
                folder.ParentFolderId = parentFolder.Id;
        }

        db.Folders.Add(folder);
        await db.SaveChangesAsync(ct);
        return folder;
    }

    private async Task ProcessVideoFileAsync(StashContext db, string path, CancellationToken ct)
    {
        var fileInfo = new FileInfo(path);
        var dirPath = Path.GetDirectoryName(path) ?? path;
        var folder = await EnsureFolderAsync(db, dirPath, ct);

        var basename = Path.GetFileName(path);
        var existing = await db.VideoFiles
            .FirstOrDefaultAsync(f => f.ParentFolderId == folder.Id && f.Basename == basename, ct);

        if (existing != null)
        {
            existing.Size = fileInfo.Length;
            existing.ModTime = fileInfo.LastWriteTimeUtc;

            // Re-probe if metadata is missing (e.g., FFprobe wasn't available during initial scan)
            if (existing.Width == 0 && existing.Height == 0 && existing.Duration == 0)
            {
                await ProbeVideoAsync(existing, path, ct);
            }
            return;
        }

        // Create video file entry
        var videoFile = new VideoFile
        {
            Basename = basename,
            ParentFolderId = folder.Id,
            Size = fileInfo.Length,
            ModTime = fileInfo.LastWriteTimeUtc,
            Format = Path.GetExtension(path).TrimStart('.').ToLowerInvariant()
        };

        // Probe with FFprobe for metadata
        await ProbeVideoAsync(videoFile, path, ct);

        // Create scene for the video file
        var scene = new Scene
        {
            Title = Path.GetFileNameWithoutExtension(path),
            Files = [videoFile]
        };

        db.Scenes.Add(scene);

        // Compute oshash fingerprint
        var oshash = await ComputeOshashAsync(path, ct);
        if (oshash != null)
        {
            videoFile.Fingerprints.Add(new FileFingerprint
            {
                Type = "oshash",
                Value = oshash
            });
        }

        if (config.CalculateMd5)
        {
            var md5 = await fingerprintService.ComputeMd5Async(path, ct);
            if (!string.IsNullOrWhiteSpace(md5))
            {
                videoFile.Fingerprints.Add(new FileFingerprint
                {
                    Type = "md5",
                    Value = md5,
                });
            }
        }

        logger.LogDebug("Added scene for: {Path}", path);
        eventBus.Publish(new EntityEvent(EventType.SceneCreated, "Scene", 0, scene));
    }

    private async Task ProcessImageFileAsync(StashContext db, string path, CancellationToken ct)
    {
        var fileInfo = new FileInfo(path);
        var dirPath = Path.GetDirectoryName(path) ?? path;
        var folder = await EnsureFolderAsync(db, dirPath, ct);

        var basename = Path.GetFileName(path);
        var existing = await db.ImageFiles
            .FirstOrDefaultAsync(f => f.ParentFolderId == folder.Id && f.Basename == basename, ct);

        if (existing != null)
        {
            existing.Size = fileInfo.Length;
            existing.ModTime = fileInfo.LastWriteTimeUtc;
            return;
        }

        var imageFile = new ImageFile
        {
            Basename = basename,
            ParentFolderId = folder.Id,
            Size = fileInfo.Length,
            ModTime = fileInfo.LastWriteTimeUtc,
            Format = Path.GetExtension(path).TrimStart('.').ToLowerInvariant()
        };

        var image = new Image
        {
            Title = Path.GetFileNameWithoutExtension(path),
            Files = [imageFile]
        };

        if (config.CalculateMd5)
        {
            var md5 = await fingerprintService.ComputeMd5Async(path, ct);
            if (!string.IsNullOrWhiteSpace(md5))
            {
                imageFile.Fingerprints.Add(new FileFingerprint
                {
                    Type = "md5",
                    Value = md5,
                });
            }
        }

        db.Images.Add(image);
        logger.LogDebug("Added image for: {Path}", path);
    }

    private async Task ProcessGalleryFileAsync(StashContext db, string path, CancellationToken ct)
    {
        var fileInfo = new FileInfo(path);
        var dirPath = Path.GetDirectoryName(path) ?? path;
        var folder = await EnsureFolderAsync(db, dirPath, ct);

        var basename = Path.GetFileName(path);
        var existing = await db.Set<GalleryFile>()
            .Include(gf => gf.Gallery)
            .ThenInclude(g => g!.ImageGalleries)
            .FirstOrDefaultAsync(f => f.ParentFolderId == folder.Id && f.Basename == basename, ct);

        // If gallery exists and already has images, skip re-processing
        if (existing?.Gallery?.ImageGalleries.Count > 0)
        {
            logger.LogDebug("Gallery already processed with {Count} images: {Path}",
                existing.Gallery.ImageGalleries.Count, path);
            return;
        }

        // Create or update the gallery file entry
        GalleryFile galleryFile;
        Gallery gallery;

        if (existing != null)
        {
            // Update existing file metadata
            galleryFile = existing;
            galleryFile.Size = fileInfo.Length;
            galleryFile.ModTime = fileInfo.LastWriteTimeUtc;
            gallery = existing.Gallery!;
        }
        else
        {
            // Create new gallery file and gallery
            galleryFile = new GalleryFile
            {
                Basename = basename,
                ParentFolderId = folder.Id,
                Size = fileInfo.Length,
                ModTime = fileInfo.LastWriteTimeUtc
            };

            gallery = new Gallery
            {
                Title = Path.GetFileNameWithoutExtension(path),
                Files = [galleryFile]
            };

            db.Galleries.Add(gallery);
        }

        // Save to get the GalleryFile ID (needed for ZipFileId on images)
        await db.SaveChangesAsync(ct);

        // Now extract images from the zip file
        try
        {
            // Get all images from the zip, sorted by path
            var imageEntries = await zipGalleryReader.GetImageEntriesAsync(path, ct);

            if (imageEntries.Count == 0)
            {
                logger.LogWarning("No images found in gallery zip: {Path}", path);
                return;
            }

            logger.LogDebug("Found {Count} images in gallery: {Path}", imageEntries.Count, path);

            // Create a virtual folder for this zip's contents
            // This ensures images from different zips don't conflict on the unique constraint (ParentFolderId + Basename)
            var virtualFolderPath = $"{path}#virtual";
            var virtualFolder = await db.Folders.FirstOrDefaultAsync(f => f.Path == virtualFolderPath, ct);
            if (virtualFolder == null)
            {
                virtualFolder = new Folder { Path = virtualFolderPath };
                db.Folders.Add(virtualFolder);
                await db.SaveChangesAsync(ct);
            }

            // Create Image entities for each image in the zip
            foreach (var entry in imageEntries)
            {
                // Create ImageFile record representing the image within the zip
                // Use FullName to preserve the internal zip path structure and avoid duplicate basenames
                var imageFile = new ImageFile
                {
                    Basename = entry.FullName,  // Use full internal path to avoid collisions
                    ParentFolderId = virtualFolder.Id,  // Use virtual folder specific to this zip
                    ZipFileId = galleryFile.Id,  // Link to parent zip file
                    Size = entry.Length,
                    ModTime = entry.LastWriteTime.UtcDateTime,
                    Format = Path.GetExtension(entry.Name).TrimStart('.').ToLowerInvariant(),
                    // TODO: Extract dimensions using image processing library
                    Width = 0,
                    Height = 0
                };

                // Create Image entity
                var image = new Image
                {
                    Title = Path.GetFileNameWithoutExtension(entry.Name),
                    Files = [imageFile]
                };

                db.Images.Add(image);

                // Link image to gallery via junction table
                // Note: We'll add this after the image is saved and has an ID
                gallery.ImageGalleries.Add(new ImageGallery
                {
                    Image = image,
                    Gallery = gallery
                });
            }

            // Save all images and their gallery associations
            await db.SaveChangesAsync(ct);

            logger.LogDebug("Added gallery with {Count} images: {Path}", imageEntries.Count, path);
        }
        catch (FileNotFoundException)
        {
            logger.LogError("Zip file not found (may have been moved/deleted): {Path}", path);
        }
        catch (InvalidDataException ex)
        {
            logger.LogError("Invalid or corrupt zip file: {Path} - {Error}", path, ex.Message);
        }
        catch (Exception ex)
        {
            logger.LogError(ex, "Error processing gallery zip file: {Path}", path);
        }
    }

    /// <summary>
    /// Compute OpenSubtitles hash (oshash) for a video file.
    /// Same algorithm as the original stash.
    /// </summary>
    private static async Task<string?> ComputeOshashAsync(string path, CancellationToken ct)
    {
        const int chunkSize = 65536; // 64KB
        try
        {
            await using var stream = new FileStream(path, FileMode.Open, FileAccess.Read, FileShare.Read, chunkSize, useAsync: true);
            var fileSize = stream.Length;
            if (fileSize < chunkSize) return null;

            ulong hash = (ulong)fileSize;
            var buf = new byte[chunkSize];

            // Hash first 64KB
            await stream.ReadExactlyAsync(buf, ct);
            for (int i = 0; i < chunkSize; i += 8)
                hash += BitConverter.ToUInt64(buf, i);

            // Hash last 64KB
            stream.Seek(-chunkSize, SeekOrigin.End);
            await stream.ReadExactlyAsync(buf, ct);
            for (int i = 0; i < chunkSize; i += 8)
                hash += BitConverter.ToUInt64(buf, i);

            return hash.ToString("x16");
        }
        catch
        {
            return null;
        }
    }

    private async Task ProbeVideoAsync(VideoFile videoFile, string path, CancellationToken ct)
    {
        var ffprobePath = FindFfprobe();
        if (ffprobePath == null)
        {
            logger.LogDebug("FFprobe not found, skipping metadata probe for {Path}", path);
            return;
        }

        try
        {
            using var process = new System.Diagnostics.Process
            {
                StartInfo = new System.Diagnostics.ProcessStartInfo
                {
                    FileName = ffprobePath,
                    Arguments = $"-v quiet -print_format json -show_format -show_streams \"{path}\"",
                    UseShellExecute = false,
                    RedirectStandardOutput = true,
                    RedirectStandardError = true,
                    CreateNoWindow = true
                }
            };

            process.Start();
            var json = await process.StandardOutput.ReadToEndAsync(ct);
            await process.WaitForExitAsync(ct);

            if (process.ExitCode != 0 || string.IsNullOrEmpty(json)) return;

            using var doc = System.Text.Json.JsonDocument.Parse(json);
            var root = doc.RootElement;

            // Extract format duration
            if (root.TryGetProperty("format", out var format))
            {
                if (format.TryGetProperty("duration", out var dur) && double.TryParse(dur.GetString(), System.Globalization.NumberStyles.Float, System.Globalization.CultureInfo.InvariantCulture, out var duration))
                    videoFile.Duration = duration;
                if (format.TryGetProperty("bit_rate", out var br) && long.TryParse(br.GetString(), out var bitrate))
                    videoFile.BitRate = bitrate;
            }

            // Extract stream info
            if (root.TryGetProperty("streams", out var streams))
            {
                foreach (var stream in streams.EnumerateArray())
                {
                    var codecType = stream.TryGetProperty("codec_type", out var ct2) ? ct2.GetString() : null;
                    if (codecType == "video" && videoFile.Width == 0)
                    {
                        if (stream.TryGetProperty("width", out var w)) videoFile.Width = w.GetInt32();
                        if (stream.TryGetProperty("height", out var h)) videoFile.Height = h.GetInt32();
                        if (stream.TryGetProperty("codec_name", out var cn)) videoFile.VideoCodec = cn.GetString() ?? "";
                        if (stream.TryGetProperty("r_frame_rate", out var rfr))
                        {
                            var frs = rfr.GetString() ?? "";
                            var frParts = frs.Split('/');
                            if (frParts.Length == 2 && double.TryParse(frParts[0], out var num) && double.TryParse(frParts[1], out var den) && den > 0)
                                videoFile.FrameRate = num / den;
                        }
                    }
                    else if (codecType == "audio" && string.IsNullOrEmpty(videoFile.AudioCodec))
                    {
                        if (stream.TryGetProperty("codec_name", out var acn)) videoFile.AudioCodec = acn.GetString() ?? "";
                    }
                }
            }
        }
        catch (Exception ex)
        {
            logger.LogDebug(ex, "FFprobe failed for {Path}", path);
        }
    }

    private string? FindFfprobe()
    {
        // Check configured FFmpeg path directory for ffprobe
        if (!string.IsNullOrEmpty(config.FfmpegPath))
        {
            var dir = Path.GetDirectoryName(config.FfmpegPath);
            if (dir != null)
            {
                var probe = Path.Combine(dir, OperatingSystem.IsWindows() ? "ffprobe.exe" : "ffprobe");
                if (File.Exists(probe)) return probe;
            }
        }

        // Search PATH
        var pathEnv = Environment.GetEnvironmentVariable("PATH") ?? "";
        foreach (var dir in pathEnv.Split(Path.PathSeparator))
        {
            var probe = Path.Combine(dir, OperatingSystem.IsWindows() ? "ffprobe.exe" : "ffprobe");
            if (File.Exists(probe)) return probe;
        }

        return null;
    }

    private static bool IsExcluded(string path, List<string> patterns)
    {
        foreach (var pattern in patterns)
        {
            if (path.Contains(pattern, StringComparison.OrdinalIgnoreCase))
                return true;
        }
        return false;
    }

    private record DiscoveredFile(string Path, string Extension);
}
