namespace Stash.Core.Entities;

public class Folder
{
    public int Id { get; set; }
    public string Path { get; set; } = string.Empty;
    public int? ParentFolderId { get; set; }
    public int? ZipFileId { get; set; }
    public DateTime ModTime { get; set; }
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;

    // Navigation
    public Folder? ParentFolder { get; set; }
    public ICollection<Folder> SubFolders { get; set; } = [];
    public ICollection<BaseFileEntity> Files { get; set; } = [];
}

public abstract class BaseFileEntity
{
    public int Id { get; set; }
    public string Basename { get; set; } = string.Empty;
    public int ParentFolderId { get; set; }
    public int? ZipFileId { get; set; }
    public long Size { get; set; }
    public DateTime ModTime { get; set; }
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;

    // Navigation
    public Folder? ParentFolder { get; set; }
    public ICollection<FileFingerprint> Fingerprints { get; set; } = [];

    // Computed (not stored)
    public string Path => ParentFolder != null
        ? System.IO.Path.Combine(ParentFolder.Path, Basename)
        : Basename;
}

public class VideoFile : BaseFileEntity
{
    public string Format { get; set; } = string.Empty;
    public int Width { get; set; }
    public int Height { get; set; }
    public double Duration { get; set; }
    public string VideoCodec { get; set; } = string.Empty;
    public string AudioCodec { get; set; } = string.Empty;
    public double FrameRate { get; set; }
    public long BitRate { get; set; }
    public bool Interactive { get; set; }
    public int? InteractiveSpeed { get; set; }

    // FK to Scene
    public int? SceneId { get; set; }
    public Scene? Scene { get; set; }
}

public class ImageFile : BaseFileEntity
{
    public string Format { get; set; } = string.Empty;
    public int Width { get; set; }
    public int Height { get; set; }

    // FK to Image
    public int? ImageId { get; set; }
    public Image? Image { get; set; }
}

public class GalleryFile : BaseFileEntity
{
    // FK to Gallery
    public int? GalleryId { get; set; }
    public Gallery? Gallery { get; set; }
}

public class FileFingerprint
{
    public int Id { get; set; }
    public int FileId { get; set; }
    public string Type { get; set; } = string.Empty; // "oshash", "md5", "phash"
    public string Value { get; set; } = string.Empty;
    public BaseFileEntity? File { get; set; }
}
