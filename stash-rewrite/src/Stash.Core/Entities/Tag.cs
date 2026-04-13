namespace Stash.Core.Entities;

public class Tag : BaseEntity
{
    public string Name { get; set; } = string.Empty;
    public string? SortName { get; set; }
    public string? Description { get; set; }
    public bool Favorite { get; set; }
    public bool IgnoreAutoTag { get; set; }

    // Image stored as blob reference
    public string? ImageBlobId { get; set; }

    // Navigation properties
    public ICollection<TagAlias> Aliases { get; set; } = [];
    public ICollection<TagParent> ParentRelations { get; set; } = [];
    public ICollection<TagParent> ChildRelations { get; set; } = [];
    public ICollection<TagStashId> StashIds { get; set; } = [];

    // Reverse nav for many-to-many
    public ICollection<SceneTag> SceneTags { get; set; } = [];
    public ICollection<PerformerTag> PerformerTags { get; set; } = [];
    public ICollection<ImageTag> ImageTags { get; set; } = [];
    public ICollection<GalleryTag> GalleryTags { get; set; } = [];
    public ICollection<StudioTag> StudioTags { get; set; } = [];
    public ICollection<GroupTag> GroupTags { get; set; } = [];
    public Dictionary<string, object>? CustomFields { get; set; }
}

public class TagAlias
{
    public int Id { get; set; }
    public int TagId { get; set; }
    public string Alias { get; set; } = string.Empty;
    public Tag? Tag { get; set; }
}

public class TagParent
{
    public int ParentId { get; set; }
    public int ChildId { get; set; }
    public Tag? Parent { get; set; }
    public Tag? Child { get; set; }
}

public class TagStashId
{
    public int Id { get; set; }
    public int TagId { get; set; }
    public string Endpoint { get; set; } = string.Empty;
    public string StashId { get; set; } = string.Empty;
    public Tag? Tag { get; set; }
}
