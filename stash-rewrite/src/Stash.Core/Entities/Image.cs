namespace Stash.Core.Entities;

public class Image : BaseEntity
{
    public string? Title { get; set; }
    public string? Code { get; set; }
    public string? Details { get; set; }
    public string? Photographer { get; set; }
    public int? Rating { get; set; } // 1-100
    public bool Organized { get; set; }
    public int OCounter { get; set; }
    public int? StudioId { get; set; }
    public DateOnly? Date { get; set; }

    // Navigation properties
    public Studio? Studio { get; set; }
    public ICollection<ImageUrl> Urls { get; set; } = [];
    public ICollection<ImageFile> Files { get; set; } = [];
    public ICollection<ImageTag> ImageTags { get; set; } = [];
    public ICollection<ImagePerformer> ImagePerformers { get; set; } = [];
    public ICollection<ImageGallery> ImageGalleries { get; set; } = [];
    public Dictionary<string, object>? CustomFields { get; set; }
}

public class ImageUrl
{
    public int Id { get; set; }
    public int ImageId { get; set; }
    public string Url { get; set; } = string.Empty;
    public Image? Image { get; set; }
}

public class ImageTag
{
    public int ImageId { get; set; }
    public int TagId { get; set; }
    public Image? Image { get; set; }
    public Tag? Tag { get; set; }
}

public class ImagePerformer
{
    public int ImageId { get; set; }
    public int PerformerId { get; set; }
    public Image? Image { get; set; }
    public Performer? Performer { get; set; }
}

public class ImageGallery
{
    public int ImageId { get; set; }
    public int GalleryId { get; set; }
    public Image? Image { get; set; }
    public Gallery? Gallery { get; set; }
}
