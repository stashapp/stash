namespace Stash.Core.Entities.Galleries.Zip;

/// <summary>
/// High-level service for reading images from zip-based galleries.
/// Provides gallery-specific operations built on top of IZipFileReader.
/// </summary>
public class ZipGalleryReader(IZipFileReader zipReader)
{
    /// <summary>
    /// Gets all image entries from a zip file, sorted by path.
    /// This maintains a consistent ordering for gallery image indexing.
    /// </summary>
    /// <param name="zipFilePath">Path to the zip file</param>
    /// <param name="ct">Cancellation token</param>
    /// <returns>Sorted list of image entries found in the zip</returns>
    public async Task<List<ZipEntryInfo>> GetImageEntriesAsync(string zipFilePath, CancellationToken ct = default)
    {
        // Step 1: Get all entries from the zip file
        var allEntries = await zipReader.ListEntriesAsync(zipFilePath, ct);

        // Step 2: Filter to only image files
        var imageEntries = zipReader.FilterImageEntries(allEntries);

        // Step 3: Sort by full path to ensure consistent ordering
        // This is critical! The image index used by the gallery API depends on
        // this exact ordering being maintained. If we sort differently here
        // vs. when serving images, indexes will be incorrect.
        //
        // We use ordinal (case-sensitive) comparison to match the original
        // Stash behavior which orders by file path.
        var sortedEntries = imageEntries
            .OrderBy(e => e.FullName, StringComparer.Ordinal)
            .ToList();

        return sortedEntries;
    }

    /// <summary>
    /// Gets an image entry at a specific index within a zip-based gallery.
    /// Index is 0-based and corresponds to the sorted order of images.
    /// </summary>
    /// <param name="zipFilePath">Path to the zip file</param>
    /// <param name="index">0-based index of the image to retrieve</param>
    /// <param name="ct">Cancellation token</param>
    /// <returns>The image entry at the specified index</returns>
    /// <exception cref="IndexOutOfRangeException">When index is invalid</exception>
    public async Task<ZipEntryInfo> GetImageAtIndexAsync(
        string zipFilePath,
        int index,
        CancellationToken ct = default)
    {
        // Get sorted list of all images
        var images = await GetImageEntriesAsync(zipFilePath, ct);

        // Validate index is within bounds
        if (index < 0 || index >= images.Count)
        {
            throw new IndexOutOfRangeException(
                $"Image index {index} is out of range. Gallery contains {images.Count} images (valid range: 0-{images.Count - 1})");
        }

        return images[index];
    }

    /// <summary>
    /// Extracts an image at a specific index from a zip-based gallery.
    /// </summary>
    /// <param name="zipFilePath">Path to the zip file</param>
    /// <param name="index">0-based index of the image to extract</param>
    /// <param name="ct">Cancellation token</param>
    /// <returns>Stream containing the image data. Caller must dispose.</returns>
    public async Task<Stream> ExtractImageAtIndexAsync(
        string zipFilePath,
        int index,
        CancellationToken ct = default)
    {
        // Step 1: Get the entry info for the image at this index
        var imageEntry = await GetImageAtIndexAsync(zipFilePath, index, ct);

        // Step 2: Extract the actual image data
        return await zipReader.ExtractEntryAsync(zipFilePath, imageEntry.FullName, ct);
    }

    /// <summary>
    /// Gets metadata about a zip-based gallery (image count, total size, etc.).
    /// </summary>
    /// <param name="zipFilePath">Path to the zip file</param>
    /// <param name="ct">Cancellation token</param>
    /// <returns>Gallery metadata</returns>
    public async Task<ZipGalleryInfo> GetGalleryInfoAsync(string zipFilePath, CancellationToken ct = default)
    {
        var images = await GetImageEntriesAsync(zipFilePath, ct);

        // Calculate total uncompressed size of all images
        var totalSize = images.Sum(img => img.Length);

        // Get the zip file size on disk
        var zipFileInfo = new FileInfo(zipFilePath);

        return new ZipGalleryInfo(
            ImageCount: images.Count,
            TotalUncompressedSize: totalSize,
            ZipFileSize: zipFileInfo.Length,
            FirstImageName: images.FirstOrDefault()?.Name ?? "",
            LastImageName: images.LastOrDefault()?.Name ?? ""
        );
    }

    /// <summary>
    /// Finds a cover image within the zip based on a filename pattern.
    /// Useful for selecting a gallery cover automatically (e.g., "cover.jpg", "001.jpg").
    /// </summary>
    /// <param name="zipFilePath">Path to the zip file</param>
    /// <param name="pattern">Pattern to match (e.g., "cover", "001")</param>
    /// <param name="ct">Cancellation token</param>
    /// <returns>Index of the matching image, or null if not found</returns>
    public async Task<int?> FindCoverImageIndexAsync(
        string zipFilePath,
        string pattern,
        CancellationToken ct = default)
    {
        var images = await GetImageEntriesAsync(zipFilePath, ct);

        // Try to find an image whose name contains the pattern (case-insensitive)
        for (int i = 0; i < images.Count; i++)
        {
            if (images[i].Name.Contains(pattern, StringComparison.OrdinalIgnoreCase))
            {
                return i;
            }
        }

        // No match found - caller should fall back to default (usually index 0)
        return null;
    }
}

/// <summary>
/// Metadata about a zip-based gallery.
/// </summary>
/// <param name="ImageCount">Total number of images found in the zip</param>
/// <param name="TotalUncompressedSize">Combined size of all images when extracted</param>
/// <param name="ZipFileSize">Size of the zip file on disk</param>
/// <param name="FirstImageName">Name of the first image (alphabetically)</param>
/// <param name="LastImageName">Name of the last image (alphabetically)</param>
public record ZipGalleryInfo(
    int ImageCount,
    long TotalUncompressedSize,
    long ZipFileSize,
    string FirstImageName,
    string LastImageName
);
