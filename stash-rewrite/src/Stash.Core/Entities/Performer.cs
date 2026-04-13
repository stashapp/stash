using Stash.Core.Enums;

namespace Stash.Core.Entities;

public class Performer : BaseEntity
{
    public string Name { get; set; } = string.Empty;
    public string? Disambiguation { get; set; }
    public GenderEnum? Gender { get; set; }
    public DateOnly? Birthdate { get; set; }
    public DateOnly? DeathDate { get; set; }
    public string? Ethnicity { get; set; }
    public string? Country { get; set; }
    public string? EyeColor { get; set; }
    public string? HairColor { get; set; }
    public int? HeightCm { get; set; }
    public int? Weight { get; set; }
    public string? Measurements { get; set; }
    public string? FakeTits { get; set; }
    public double? PenisLength { get; set; }
    public CircumcisedEnum? Circumcised { get; set; }
    public DateOnly? CareerStart { get; set; }
    public DateOnly? CareerEnd { get; set; }
    public string? Tattoos { get; set; }
    public string? Piercings { get; set; }
    public bool Favorite { get; set; }
    public int? Rating { get; set; } // 1-100
    public string? Details { get; set; }
    public bool IgnoreAutoTag { get; set; }

    // Image stored as blob reference
    public string? ImageBlobId { get; set; }

    // Navigation properties
    public ICollection<PerformerUrl> Urls { get; set; } = [];
    public ICollection<PerformerAlias> Aliases { get; set; } = [];
    public ICollection<PerformerTag> PerformerTags { get; set; } = [];
    public ICollection<ScenePerformer> ScenePerformers { get; set; } = [];
    public ICollection<ImagePerformer> ImagePerformers { get; set; } = [];
    public ICollection<GalleryPerformer> GalleryPerformers { get; set; } = [];
    public ICollection<PerformerStashId> StashIds { get; set; } = [];
    public Dictionary<string, object>? CustomFields { get; set; }
}

public class PerformerUrl
{
    public int Id { get; set; }
    public int PerformerId { get; set; }
    public string Url { get; set; } = string.Empty;
    public Performer? Performer { get; set; }
}

public class PerformerAlias
{
    public int Id { get; set; }
    public int PerformerId { get; set; }
    public string Alias { get; set; } = string.Empty;
    public Performer? Performer { get; set; }
}

public class PerformerTag
{
    public int PerformerId { get; set; }
    public int TagId { get; set; }
    public Performer? Performer { get; set; }
    public Tag? Tag { get; set; }
}

public class PerformerStashId
{
    public int Id { get; set; }
    public int PerformerId { get; set; }
    public string Endpoint { get; set; } = string.Empty;
    public string StashId { get; set; } = string.Empty;
    public Performer? Performer { get; set; }
}
