namespace Stash.Core.Entities;

public class Gallery : BaseEntity
{
    public string? Title { get; set; }
    public string? Code { get; set; }
    public DateOnly? Date { get; set; }
    public string? Details { get; set; }
    public string? Photographer { get; set; }
    public int? Rating { get; set; } // 1-100
    public bool Organized { get; set; }
    public int? StudioId { get; set; }
    public int? FolderId { get; set; }

    // Navigation properties
    public Studio? Studio { get; set; }
    public Folder? Folder { get; set; }
    public ICollection<GalleryUrl> Urls { get; set; } = [];
    public ICollection<GalleryFile> Files { get; set; } = [];
    public ICollection<GalleryChapter> Chapters { get; set; } = [];
    public ICollection<GalleryTag> GalleryTags { get; set; } = [];
    public ICollection<GalleryPerformer> GalleryPerformers { get; set; } = [];
    public ICollection<SceneGallery> SceneGalleries { get; set; } = [];
    public ICollection<ImageGallery> ImageGalleries { get; set; } = [];
    public Dictionary<string, object>? CustomFields { get; set; }
}

public class GalleryUrl
{
    public int Id { get; set; }
    public int GalleryId { get; set; }
    public string Url { get; set; } = string.Empty;
    public Gallery? Gallery { get; set; }
}

public class GalleryTag
{
    public int GalleryId { get; set; }
    public int TagId { get; set; }
    public Gallery? Gallery { get; set; }
    public Tag? Tag { get; set; }
}

public class GalleryPerformer
{
    public int GalleryId { get; set; }
    public int PerformerId { get; set; }
    public Gallery? Gallery { get; set; }
    public Performer? Performer { get; set; }
}

public class GalleryChapter : BaseEntity
{
    public string Title { get; set; } = string.Empty;
    public int ImageIndex { get; set; }
    public int GalleryId { get; set; }
    public Gallery? Gallery { get; set; }
}
