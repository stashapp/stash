using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Stash.Core.Interfaces;
using Stash.Data;

namespace Stash.Api.Services;

public class StreamService(IServiceScopeFactory scopeFactory, IThumbnailService thumbnailService) : IStreamService
{
    private static readonly Dictionary<string, string> MimeTypes = new(StringComparer.OrdinalIgnoreCase)
    {
        [".mp4"] = "video/mp4",
        [".mkv"] = "video/x-matroska",
        [".avi"] = "video/x-msvideo",
        [".webm"] = "video/webm",
        [".mov"] = "video/quicktime",
        [".wmv"] = "video/x-ms-wmv",
        [".flv"] = "video/x-flv",
        [".m4v"] = "video/x-m4v",
        [".mpg"] = "video/mpeg",
        [".mpeg"] = "video/mpeg",
        [".ts"] = "video/mp2t",
        [".rmvb"] = "application/vnd.rn-realmedia-vbr",
        [".rm"] = "application/vnd.rn-realmedia",
    };

    public async Task<(Stream stream, string contentType, long? fileSize)?> GetSceneStream(int sceneId, CancellationToken ct = default)
    {
        using var scope = scopeFactory.CreateScope();
        var db = scope.ServiceProvider.GetRequiredService<StashContext>();

        var videoFile = await db.VideoFiles
            .Include(f => f.ParentFolder)
            .FirstOrDefaultAsync(f => f.SceneId == sceneId, ct);

        if (videoFile == null) return null;

        var filePath = videoFile.ParentFolder != null
            ? Path.Combine(videoFile.ParentFolder.Path, videoFile.Basename)
            : videoFile.Basename;

        if (!File.Exists(filePath)) return null;

        var ext = Path.GetExtension(filePath);
        var contentType = MimeTypes.GetValueOrDefault(ext, "application/octet-stream");
        var fileInfo = new FileInfo(filePath);

        var stream = new FileStream(filePath, FileMode.Open, FileAccess.Read, FileShare.Read, 81920, useAsync: true);
        return (stream, contentType, fileInfo.Length);
    }

    public async Task<(Stream stream, string contentType)?> GetSceneScreenshot(int sceneId, double? seconds, CancellationToken ct = default)
    {
        // For timestamped thumbnails, check cache first, then generate
        if (seconds.HasValue)
        {
            var tsPath = thumbnailService.GetTimestampedThumbnailPath(sceneId, seconds.Value);
            if (!File.Exists(tsPath))
                await thumbnailService.GenerateSceneThumbnailAsync(sceneId, seconds, ct);
            if (File.Exists(tsPath))
            {
                var stream = new FileStream(tsPath, FileMode.Open, FileAccess.Read, FileShare.Read, 8192, useAsync: true);
                return (stream, "image/jpeg");
            }
            return null;
        }

        // Default thumbnail (no timestamp)
        var thumbPath = await thumbnailService.GetSceneThumbnailPathAsync(sceneId, ct);
        if (thumbPath == null) return null;

        var defaultStream = new FileStream(thumbPath, FileMode.Open, FileAccess.Read, FileShare.Read, 8192, useAsync: true);
        return (defaultStream, "image/jpeg");
    }
}
