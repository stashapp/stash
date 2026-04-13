namespace Stash.Core.Entities;

public class Scene : BaseEntity
{
    public string? Title { get; set; }
    public string? Code { get; set; }
    public string? Details { get; set; }
    public string? Director { get; set; }
    public DateOnly? Date { get; set; }
    public int? Rating { get; set; } // 1-100
    public bool Organized { get; set; }
    public int? StudioId { get; set; }
    public double ResumeTime { get; set; }
    public double PlayDuration { get; set; }
    public int PlayCount { get; set; }
    public DateTime? LastPlayedAt { get; set; }
    public int OCounter { get; set; }

    // Navigation properties
    public Studio? Studio { get; set; }
    public ICollection<SceneUrl> Urls { get; set; } = [];
    public ICollection<VideoFile> Files { get; set; } = [];
    public ICollection<SceneMarker> SceneMarkers { get; set; } = [];
    public ICollection<SceneTag> SceneTags { get; set; } = [];
    public ICollection<ScenePerformer> ScenePerformers { get; set; } = [];
    public ICollection<SceneGallery> SceneGalleries { get; set; } = [];
    public ICollection<SceneGroup> SceneGroups { get; set; } = [];
    public ICollection<SceneStashId> StashIds { get; set; } = [];
    public ICollection<ScenePlayHistory> PlayHistory { get; set; } = [];
    public ICollection<SceneOHistory> OHistory { get; set; } = [];
    public Dictionary<string, object>? CustomFields { get; set; }
}

public class SceneUrl
{
    public int Id { get; set; }
    public int SceneId { get; set; }
    public string Url { get; set; } = string.Empty;
    public Scene? Scene { get; set; }
}

public class SceneTag
{
    public int SceneId { get; set; }
    public int TagId { get; set; }
    public Scene? Scene { get; set; }
    public Tag? Tag { get; set; }
}

public class ScenePerformer
{
    public int SceneId { get; set; }
    public int PerformerId { get; set; }
    public Scene? Scene { get; set; }
    public Performer? Performer { get; set; }
}

public class SceneGallery
{
    public int SceneId { get; set; }
    public int GalleryId { get; set; }
    public Scene? Scene { get; set; }
    public Gallery? Gallery { get; set; }
}

public class SceneGroup
{
    public int SceneId { get; set; }
    public int GroupId { get; set; }
    public int SceneIndex { get; set; }
    public Scene? Scene { get; set; }
    public Group? Group { get; set; }
}

public class SceneStashId
{
    public int Id { get; set; }
    public int SceneId { get; set; }
    public string Endpoint { get; set; } = string.Empty;
    public string StashId { get; set; } = string.Empty;
    public Scene? Scene { get; set; }
}

public class ScenePlayHistory
{
    public int Id { get; set; }
    public int SceneId { get; set; }
    public DateTime PlayedAt { get; set; }
    public Scene? Scene { get; set; }
}

public class SceneOHistory
{
    public int Id { get; set; }
    public int SceneId { get; set; }
    public DateTime OccurredAt { get; set; }
    public Scene? Scene { get; set; }
}
