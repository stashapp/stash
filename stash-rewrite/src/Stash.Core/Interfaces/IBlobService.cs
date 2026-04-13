namespace Stash.Core.Interfaces;

public interface IBlobService
{
    Task<string> StoreBlobAsync(Stream data, string contentType, CancellationToken ct = default);
    Task<(Stream Stream, string ContentType)?> GetBlobAsync(string blobId, CancellationToken ct = default);
    Task DeleteBlobAsync(string blobId, CancellationToken ct = default);
}
