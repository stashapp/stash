namespace Stash.Core.Entities;

public class SceneMarker : BaseEntity
{
    public string Title { get; set; } = string.Empty;
    public double Seconds { get; set; }
    public double? EndSeconds { get; set; }
    public int PrimaryTagId { get; set; }
    public int SceneId { get; set; }

    // Navigation
    public Tag? PrimaryTag { get; set; }
    public Scene? Scene { get; set; }
    public ICollection<SceneMarkerTag> SceneMarkerTags { get; set; } = [];
}

public class SceneMarkerTag
{
    public int SceneMarkerId { get; set; }
    public int TagId { get; set; }
    public SceneMarker? SceneMarker { get; set; }
    public Tag? Tag { get; set; }
}
