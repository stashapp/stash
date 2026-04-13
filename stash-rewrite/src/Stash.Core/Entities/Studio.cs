namespace Stash.Core.Entities;

public class Studio : BaseEntity
{
    public string Name { get; set; } = string.Empty;
    public int? ParentId { get; set; }
    public int? Rating { get; set; } // 1-100
    public bool Favorite { get; set; }
    public string? Details { get; set; }
    public bool IgnoreAutoTag { get; set; }
    public bool Organized { get; set; }

    // Image stored as blob reference
    public string? ImageBlobId { get; set; }

    // Navigation properties
    public Studio? Parent { get; set; }
    public ICollection<Studio> Children { get; set; } = [];
    public ICollection<StudioUrl> Urls { get; set; } = [];
    public ICollection<StudioAlias> Aliases { get; set; } = [];
    public ICollection<StudioTag> StudioTags { get; set; } = [];
    public ICollection<StudioStashId> StashIds { get; set; } = [];
    public ICollection<Scene> Scenes { get; set; } = [];
    public ICollection<Gallery> Galleries { get; set; } = [];
    public ICollection<Image> Images { get; set; } = [];
    public ICollection<Group> Groups { get; set; } = [];
    public Dictionary<string, object>? CustomFields { get; set; }
}

public class StudioUrl
{
    public int Id { get; set; }
    public int StudioId { get; set; }
    public string Url { get; set; } = string.Empty;
    public Studio? Studio { get; set; }
}

public class StudioAlias
{
    public int Id { get; set; }
    public int StudioId { get; set; }
    public string Alias { get; set; } = string.Empty;
    public Studio? Studio { get; set; }
}

public class StudioTag
{
    public int StudioId { get; set; }
    public int TagId { get; set; }
    public Studio? Studio { get; set; }
    public Tag? Tag { get; set; }
}

public class StudioStashId
{
    public int Id { get; set; }
    public int StudioId { get; set; }
    public string Endpoint { get; set; } = string.Empty;
    public string StashId { get; set; } = string.Empty;
    public Studio? Studio { get; set; }
}
