using Microsoft.AspNetCore.Mvc;
using Stash.Api.Services;
using Stash.Core.Interfaces;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class StreamController(IStreamService streamService, IThumbnailService thumbnailService) : ControllerBase
{
    private static readonly Dictionary<string, string> ImageMimeTypes = new(StringComparer.OrdinalIgnoreCase)
    {
        [".jpg"] = "image/jpeg", [".jpeg"] = "image/jpeg", [".png"] = "image/png",
        [".gif"] = "image/gif", [".webp"] = "image/webp", [".bmp"] = "image/bmp",
        [".tiff"] = "image/tiff", [".tif"] = "image/tiff", [".svg"] = "image/svg+xml",
    };

    [HttpGet("scene/{sceneId:int}")]
    public async Task<IActionResult> StreamScene(int sceneId, CancellationToken ct)
    {
        var result = await streamService.GetSceneStream(sceneId, ct);
        if (result == null) return NotFound();

        var (stream, contentType, fileSize) = result.Value;
        Response.Headers["Accept-Ranges"] = "bytes";

        if (fileSize.HasValue)
            return File(stream, contentType, enableRangeProcessing: true);

        return File(stream, contentType);
    }

    [HttpGet("scene/{sceneId:int}/screenshot")]
    public async Task<IActionResult> GetScreenshot(int sceneId, [FromQuery] double? seconds, CancellationToken ct)
    {
        var result = await streamService.GetSceneScreenshot(sceneId, seconds, ct);
        if (result == null) return NotFound();

        var (stream, contentType) = result.Value;
        Response.Headers["Cache-Control"] = "public, max-age=86400";
        return File(stream, contentType);
    }

    [HttpGet("scene/{sceneId:int}/preview")]
    public IActionResult GetPreview(int sceneId)
    {
        var path = thumbnailService.GetPreviewPath(sceneId);
        if (!System.IO.File.Exists(path)) return NotFound();

        var stream = new FileStream(path, FileMode.Open, FileAccess.Read, FileShare.Read, 81920, useAsync: true);
        Response.Headers["Cache-Control"] = "public, max-age=86400";
        Response.Headers["Accept-Ranges"] = "bytes";
        return File(stream, "video/mp4", enableRangeProcessing: true);
    }

    [HttpGet("scene/{sceneId:int}/sprite")]
    public IActionResult GetSprite(int sceneId)
    {
        var path = thumbnailService.GetSpritePath(sceneId);
        if (!System.IO.File.Exists(path)) return NotFound();

        var stream = new FileStream(path, FileMode.Open, FileAccess.Read, FileShare.Read, 8192, useAsync: true);
        Response.Headers["Cache-Control"] = "public, max-age=86400";
        return File(stream, "image/jpeg");
    }

    [HttpGet("scene/{sceneId:int}/vtt/thumbs")]
    public IActionResult GetSpriteVtt(int sceneId)
    {
        var path = thumbnailService.GetSpriteVttPath(sceneId);
        if (!System.IO.File.Exists(path)) return NotFound();

        var stream = new FileStream(path, FileMode.Open, FileAccess.Read, FileShare.Read, 8192, useAsync: true);
        Response.Headers["Cache-Control"] = "public, max-age=86400";
        return File(stream, "text/vtt");
    }

    [HttpGet("image/{imageId:int}")]
    public async Task<IActionResult> GetImage(int imageId, CancellationToken ct)
    {
        var filePath = await thumbnailService.GetImageFilePathAsync(imageId, ct);
        if (filePath == null) return NotFound();

        var ext = Path.GetExtension(filePath);
        var contentType = ImageMimeTypes.GetValueOrDefault(ext, "application/octet-stream");
        var stream = new FileStream(filePath, FileMode.Open, FileAccess.Read, FileShare.Read, 81920, useAsync: true);
        var fileInfo = new FileInfo(filePath);

        Response.Headers["Accept-Ranges"] = "bytes";
        Response.Headers["Cache-Control"] = "public, max-age=86400";
        return File(stream, contentType, enableRangeProcessing: true);
    }

    [HttpGet("image/{imageId:int}/thumbnail")]
    public async Task<IActionResult> GetImageThumbnail(int imageId, CancellationToken ct)
    {
        // For images, just serve the original file (they're already images)
        return await GetImage(imageId, ct);
    }
}
