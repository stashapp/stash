using System.IO.Compression;
using System.Text;

namespace Stash.Core.Entities.Galleries.Zip;

/// <summary>
/// Default implementation of IZipFileReader using .NET's System.IO.Compression.
/// Handles reading zip archives and extracting image files for gallery support.
/// </summary>
public class ZipFileReader : IZipFileReader
{
    // Supported image file extensions (case-insensitive)
    // Based on Stash's supported formats: JPEG, PNG, GIF, WebP, BMP
    private static readonly HashSet<string> ImageExtensions = new(StringComparer.OrdinalIgnoreCase)
    {
        ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg"
    };

    /// <inheritdoc/>
    public async Task<List<ZipEntryInfo>> ListEntriesAsync(string zipFilePath, CancellationToken ct = default)
    {
        // Validate that the zip file exists before attempting to open
        if (!File.Exists(zipFilePath))
            throw new FileNotFoundException($"Zip file not found: {zipFilePath}");

        // Open the zip file in read mode
        // Using FileMode.Open ensures we only read existing files
        // FileAccess.Read prevents any accidental modifications
        // FileShare.Read allows other processes to read the file concurrently
        await using var fileStream = new FileStream(
            zipFilePath,
            FileMode.Open,
            FileAccess.Read,
            FileShare.Read
        );

        // Create a ZipArchive from the file stream
        // ZipArchiveMode.Read is efficient for read-only operations
        // leaveOpen: false ensures the stream is disposed when ZipArchive is disposed
        using var archive = new ZipArchive(fileStream, ZipArchiveMode.Read, leaveOpen: false);

        var entries = new List<ZipEntryInfo>();

        // Iterate through all entries in the zip archive
        foreach (var entry in archive.Entries)
        {
            // Check for cancellation before processing each entry
            ct.ThrowIfCancellationRequested();

            // Skip directory entries (they have no content, only metadata)
            // Directory entries typically end with '/' and have 0 length
            if (entry.FullName.EndsWith('/') || entry.Length == 0)
                continue;

            // Create metadata record for this entry
            var entryInfo = new ZipEntryInfo(
                FullName: entry.FullName,
                Name: entry.Name,
                Length: entry.Length,
                CompressedLength: entry.CompressedLength,
                LastWriteTime: entry.LastWriteTime
            );

            // Attempt to fix encoding issues in filenames
            // Many zip tools use non-UTF8 encodings (CP437, Shift-JIS, etc.)
            entryInfo = FixEntryEncoding(entryInfo);

            entries.Add(entryInfo);
        }

        return entries;
    }

    /// <inheritdoc/>
    public async Task<Stream> ExtractEntryAsync(string zipFilePath, string entryPath, CancellationToken ct = default)
    {
        // Validate zip file exists
        if (!File.Exists(zipFilePath))
            throw new FileNotFoundException($"Zip file not found: {zipFilePath}");

        // Open zip archive for reading
        await using var fileStream = new FileStream(
            zipFilePath,
            FileMode.Open,
            FileAccess.Read,
            FileShare.Read
        );

        using var archive = new ZipArchive(fileStream, ZipArchiveMode.Read, leaveOpen: false);

        // Find the requested entry by its full path
        // Entry names in zip files use forward slashes regardless of OS
        var entry = archive.GetEntry(entryPath);
        if (entry == null)
            throw new FileNotFoundException($"Entry '{entryPath}' not found in zip archive");

        // Create a memory stream to hold the extracted data
        // We can't return the entry.Open() stream directly because it depends on
        // the ZipArchive remaining open, but we're disposing it at the end of this method.
        // So we copy the data to a MemoryStream that the caller can use independently.
        var memoryStream = new MemoryStream();

        // Open the entry's data stream and copy to memory
        await using (var entryStream = entry.Open())
        {
            await entryStream.CopyToAsync(memoryStream, ct);
        }

        // Reset the memory stream position to the beginning so the caller can read it
        memoryStream.Position = 0;

        return memoryStream;
    }

    /// <inheritdoc/>
    public List<ZipEntryInfo> FilterImageEntries(List<ZipEntryInfo> entries)
    {
        // Filter entries to only include files with image extensions
        // This prevents processing non-image files (txt, nfo, metadata files, etc.)
        return entries
            .Where(e => IsImageFile(e.Name))
            .ToList();
    }

    /// <inheritdoc/>
    public ZipEntryInfo FixEntryEncoding(ZipEntryInfo entry)
    {
        // NOTE: Basic implementation for now
        // The original Stash Go code uses the 'chardet' library to detect
        // character encodings automatically (Shift-JIS, CP437, ISO-8859-1, etc.)
        //
        // For a complete implementation, we would need to:
        // 1. Check if the filename contains invalid UTF-8 sequences
        // 2. Use a character encoding detection library (like Ude or similar)
        // 3. Re-decode the filename using the detected encoding
        //
        // For now, we'll return the entry as-is, which works fine for
        // properly UTF-8 encoded zip files (most modern archives)
        //
        // TODO: Implement advanced encoding detection if needed for legacy archives
        // See: pkg/file/zip.go lines 50-86 in original Stash for reference

        return entry;
    }

    /// <summary>
    /// Checks if a filename has a supported image extension.
    /// </summary>
    /// <param name="filename">Name of the file to check</param>
    /// <returns>True if the file is a recognized image format</returns>
    private static bool IsImageFile(string filename)
    {
        var extension = Path.GetExtension(filename);
        return !string.IsNullOrEmpty(extension) && ImageExtensions.Contains(extension);
    }
}
