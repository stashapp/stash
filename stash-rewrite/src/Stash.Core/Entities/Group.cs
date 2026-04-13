namespace Stash.Core.Entities;

public class Group : BaseEntity
{
    public string Name { get; set; } = string.Empty;
    public string? Aliases { get; set; }
    public int? Duration { get; set; } // seconds
    public DateOnly? Date { get; set; }
    public int? Rating { get; set; } // 1-100
    public int? StudioId { get; set; }
    public string? Director { get; set; }
    public string? Synopsis { get; set; }

    // Image blobs
    public string? FrontImageBlobId { get; set; }
    public string? BackImageBlobId { get; set; }

    // Navigation properties
    public Studio? Studio { get; set; }
    public ICollection<GroupUrl> Urls { get; set; } = [];
    public ICollection<GroupTag> GroupTags { get; set; } = [];
    public ICollection<SceneGroup> SceneGroups { get; set; } = [];
    public ICollection<GroupRelation> ContainingGroupRelations { get; set; } = [];
    public ICollection<GroupRelation> SubGroupRelations { get; set; } = [];
    public Dictionary<string, object>? CustomFields { get; set; }
}

public class GroupUrl
{
    public int Id { get; set; }
    public int GroupId { get; set; }
    public string Url { get; set; } = string.Empty;
    public Group? Group { get; set; }
}

public class GroupTag
{
    public int GroupId { get; set; }
    public int TagId { get; set; }
    public Group? Group { get; set; }
    public Tag? Tag { get; set; }
}

public class GroupRelation
{
    public int ContainingGroupId { get; set; }
    public int SubGroupId { get; set; }
    public int OrderIndex { get; set; }
    public string? Description { get; set; }
    public Group? ContainingGroup { get; set; }
    public Group? SubGroup { get; set; }
}
