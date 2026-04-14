namespace Stash.Core.Entities.Galleries.Zip;

/// <summary>
/// Service for reading and extracting files from ZIP archives.
/// Provides abstraction over zip file operations for gallery image extraction.
/// </summary>
public interface IZipFileReader
{
    /// <summary>
    /// Lists all entries (files) contained within a zip archive.
    /// </summary>
    /// <param name="zipFilePath">Full path to the zip file on disk</param>
    /// <param name="ct">Cancellation token</param>
    /// <returns>List of entries with metadata (name, size, compressed size)</returns>
    /// <exception cref="FileNotFoundException">When zip file doesn't exist</exception>
    /// <exception cref="InvalidDataException">When file is not a valid zip archive</exception>
    Task<List<ZipEntryInfo>> ListEntriesAsync(string zipFilePath, CancellationToken ct = default);

    /// <summary>
    /// Extracts a specific file from the zip archive to a stream.
    /// </summary>
    /// <param name="zipFilePath">Full path to the zip file on disk</param>
    /// <param name="entryPath">Path of the entry within the zip (e.g., "folder/image.jpg")</param>
    /// <param name="ct">Cancellation token</param>
    /// <returns>Stream containing the extracted file data. Caller is responsible for disposing.</returns>
    /// <exception cref="FileNotFoundException">When zip file or entry doesn't exist</exception>
    Task<Stream> ExtractEntryAsync(string zipFilePath, string entryPath, CancellationToken ct = default);

    /// <summary>
    /// Filters entries by supported image extensions.
    /// </summary>
    /// <param name="entries">List of zip entries to filter</param>
    /// <returns>Only entries that are recognized as image files</returns>
    List<ZipEntryInfo> FilterImageEntries(List<ZipEntryInfo> entries);

    /// <summary>
    /// Detects and converts non-UTF8 encoded filenames in zip entries.
    /// Some zip archives (especially from non-English systems) may have filenames
    /// encoded in legacy character sets. This method attempts to detect and convert them.
    /// </summary>
    /// <param name="entry">Entry with potentially mis-encoded filename</param>
    /// <returns>Entry with corrected filename if encoding was detected and fixed</returns>
    ZipEntryInfo FixEntryEncoding(ZipEntryInfo entry);
}

/// <summary>
/// Represents metadata about a file entry within a zip archive.
/// </summary>
/// <param name="FullName">Full path of the entry within the zip (includes directories)</param>
/// <param name="Name">Filename only (without directory path)</param>
/// <param name="Length">Uncompressed size in bytes</param>
/// <param name="CompressedLength">Compressed size in bytes (as stored in zip)</param>
/// <param name="LastWriteTime">Last modification timestamp of the entry</param>
public record ZipEntryInfo(
    string FullName,
    string Name,
    long Length,
    long CompressedLength,
    DateTimeOffset LastWriteTime
);
