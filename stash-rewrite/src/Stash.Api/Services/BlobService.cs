using Stash.Core.Interfaces;

namespace Stash.Api.Services;

public class BlobService(StashConfiguration config, ILogger<BlobService> logger) : IBlobService
{
    private static readonly Dictionary<string, string> ContentTypeToExtension = new(StringComparer.OrdinalIgnoreCase)
    {
        ["image/jpeg"] = ".jpg",
        ["image/png"] = ".png",
        ["image/webp"] = ".webp",
        ["image/gif"] = ".gif",
    };

    private static readonly Dictionary<string, string> ExtensionToContentType = ContentTypeToExtension
        .ToDictionary(kvp => kvp.Value, kvp => kvp.Key, StringComparer.OrdinalIgnoreCase);

    private string BlobDir => Path.Combine(config.GeneratedPath, "blobs");

    public async Task<string> StoreBlobAsync(Stream data, string contentType, CancellationToken ct = default)
    {
        var blobId = Guid.NewGuid().ToString();
        var extension = GetExtension(contentType);
        var path = GetBlobPath(blobId, extension);

        Directory.CreateDirectory(Path.GetDirectoryName(path)!);

        await using var fs = new FileStream(path, FileMode.Create, FileAccess.Write, FileShare.None);
        await data.CopyToAsync(fs, ct);

        logger.LogDebug("Stored blob {BlobId} at {Path}", blobId, path);
        return blobId;
    }

    public Task<(Stream Stream, string ContentType)?> GetBlobAsync(string blobId, CancellationToken ct = default)
    {
        var (path, contentType) = ResolveBlobFile(blobId);
        if (path == null || contentType == null)
            return Task.FromResult<(Stream Stream, string ContentType)?>(null);

        var fs = new FileStream(path, FileMode.Open, FileAccess.Read, FileShare.Read);
        return Task.FromResult<(Stream Stream, string ContentType)?>((fs, contentType));
    }

    public Task DeleteBlobAsync(string blobId, CancellationToken ct = default)
    {
        var (path, _) = ResolveBlobFile(blobId);
        if (path != null)
        {
            File.Delete(path);
            logger.LogDebug("Deleted blob {BlobId} at {Path}", blobId, path);
        }

        return Task.CompletedTask;
    }

    private string GetBlobPath(string blobId, string extension)
    {
        var bucket = blobId[..2];
        return Path.Combine(BlobDir, bucket, $"{blobId}{extension}");
    }

    /// <summary>
    /// Finds the blob file on disk by checking all known extensions in the bucket directory.
    /// </summary>
    private (string? Path, string? ContentType) ResolveBlobFile(string blobId)
    {
        var bucket = blobId[..2];
        var dir = Path.Combine(BlobDir, bucket);

        foreach (var (ext, ct) in ExtensionToContentType)
        {
            var candidate = Path.Combine(dir, $"{blobId}{ext}");
            if (File.Exists(candidate))
                return (candidate, ct);
        }

        return (null, null);
    }

    private static string GetExtension(string contentType)
    {
        if (ContentTypeToExtension.TryGetValue(contentType, out var ext))
            return ext;

        // Fallback: derive from subtype
        var slash = contentType.IndexOf('/');
        return slash >= 0 ? $".{contentType[(slash + 1)..]}" : ".bin";
    }
}
