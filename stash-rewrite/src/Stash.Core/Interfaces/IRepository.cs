using Stash.Core.Entities;

namespace Stash.Core.Interfaces;

public interface IRepository<T> where T : class
{
    Task<T?> GetByIdAsync(int id, CancellationToken ct = default);
    Task<IReadOnlyList<T>> GetAllAsync(CancellationToken ct = default);
    Task<T> AddAsync(T entity, CancellationToken ct = default);
    Task UpdateAsync(T entity, CancellationToken ct = default);
    Task DeleteAsync(int id, CancellationToken ct = default);
    Task<int> CountAsync(CancellationToken ct = default);
}

public interface ISceneRepository : IRepository<Scene>
{
    Task<(IReadOnlyList<Scene> Items, int TotalCount)> FindAsync(SceneFilter? filter, FindFilter? findFilter, CancellationToken ct = default);
    Task<Scene?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default);
}

public interface IPerformerRepository : IRepository<Performer>
{
    Task<(IReadOnlyList<Performer> Items, int TotalCount)> FindAsync(PerformerFilter? filter, FindFilter? findFilter, CancellationToken ct = default);
    Task<Performer?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default);
}

public interface ITagRepository : IRepository<Tag>
{
    Task<(IReadOnlyList<Tag> Items, int TotalCount)> FindAsync(TagFilter? filter, FindFilter? findFilter, CancellationToken ct = default);
    Task<Tag?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default);
    Task<Tag?> GetByNameAsync(string name, CancellationToken ct = default);
}

public interface IStudioRepository : IRepository<Studio>
{
    Task<(IReadOnlyList<Studio> Items, int TotalCount)> FindAsync(StudioFilter? filter, FindFilter? findFilter, CancellationToken ct = default);
    Task<Studio?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default);
}

public interface IGalleryRepository : IRepository<Gallery>
{
    Task<(IReadOnlyList<Gallery> Items, int TotalCount)> FindAsync(GalleryFilter? filter, FindFilter? findFilter, CancellationToken ct = default);
    Task<Gallery?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default);
}

public interface IImageRepository : IRepository<Image>
{
    Task<(IReadOnlyList<Image> Items, int TotalCount)> FindAsync(ImageFilter? filter, FindFilter? findFilter, CancellationToken ct = default);
    Task<Image?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default);
}

public interface IGroupRepository : IRepository<Group>
{
    Task<(IReadOnlyList<Group> Items, int TotalCount)> FindAsync(GroupFilter? filter, FindFilter? findFilter, CancellationToken ct = default);
    Task<Group?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default);
}

public interface ISavedFilterRepository : IRepository<SavedFilter>
{
    Task<IReadOnlyList<SavedFilter>> GetByModeAsync(Stash.Core.Enums.FilterMode mode, CancellationToken ct = default);
}

public interface ISceneMarkerRepository : IRepository<SceneMarker>
{
    Task<IReadOnlyList<SceneMarker>> GetBySceneIdAsync(int sceneId, CancellationToken ct = default);
}

// Filter models
public class FindFilter
{
    public string? Q { get; set; }
    public int Page { get; set; } = 1;
    public int PerPage { get; set; } = 25;
    public string? Sort { get; set; }
    public Stash.Core.Enums.SortDirection Direction { get; set; } = Stash.Core.Enums.SortDirection.Asc;
}

// Criterion modifier for advanced filters
public enum CriterionModifier
{
    Equals, NotEquals, GreaterThan, LessThan,
    Includes, Excludes, IncludesAll, ExcludesAll,
    IsNull, NotNull, Between, NotBetween,
    MatchesRegex, NotMatchesRegex
}

public class IntCriterion { public int Value { get; set; } public int? Value2 { get; set; } public CriterionModifier Modifier { get; set; } = CriterionModifier.Equals; }
public class StringCriterion { public string Value { get; set; } = ""; public CriterionModifier Modifier { get; set; } = CriterionModifier.Equals; }
public class BoolCriterion { public bool Value { get; set; } }
public class MultiIdCriterion { public List<int> Value { get; set; } = []; public CriterionModifier Modifier { get; set; } = CriterionModifier.Includes; public List<int>? Excludes { get; set; } }
public class DateCriterion { public string Value { get; set; } = ""; public string? Value2 { get; set; } public CriterionModifier Modifier { get; set; } = CriterionModifier.Equals; }
public class TimestampCriterion { public string Value { get; set; } = ""; public string? Value2 { get; set; } public CriterionModifier Modifier { get; set; } = CriterionModifier.Equals; }

public class SceneFilter
{
    public string? Title { get; set; }
    public string? Code { get; set; }
    public string? Path { get; set; }
    public int? Rating { get; set; }
    public bool? Organized { get; set; }
    public int? StudioId { get; set; }
    public int? GroupId { get; set; }
    public int? GalleryId { get; set; }
    public List<int>? TagIds { get; set; }
    public List<int>? PerformerIds { get; set; }
    // Advanced criteria
    public IntCriterion? RatingCriterion { get; set; }
    public IntCriterion? OCounterCriterion { get; set; }
    public IntCriterion? DurationCriterion { get; set; }
    public IntCriterion? ResolutionCriterion { get; set; }
    public IntCriterion? PlayCountCriterion { get; set; }
    public IntCriterion? PerformerCountCriterion { get; set; }
    public MultiIdCriterion? TagsCriterion { get; set; }
    public MultiIdCriterion? PerformersCriterion { get; set; }
    public MultiIdCriterion? StudiosCriterion { get; set; }
    public MultiIdCriterion? GroupsCriterion { get; set; }
    public BoolCriterion? OrganizedCriterion { get; set; }
    public BoolCriterion? HasMarkersCriterion { get; set; }
    public BoolCriterion? InteractiveCriterion { get; set; }
    public StringCriterion? PathCriterion { get; set; }
    public StringCriterion? UrlCriterion { get; set; }
    public DateCriterion? DateCriterion { get; set; }
    public TimestampCriterion? CreatedAtCriterion { get; set; }
    public TimestampCriterion? UpdatedAtCriterion { get; set; }
    public BoolCriterion? PerformerFavoriteCriterion { get; set; }
    public StringCriterion? VideoCodecCriterion { get; set; }
    public StringCriterion? AudioCodecCriterion { get; set; }
    public IntCriterion? FrameRateCriterion { get; set; }
    public IntCriterion? BitrateInterval { get; set; }
    public IntCriterion? FileCountCriterion { get; set; }
    public StringCriterion? StashIdCriterion { get; set; }
    public BoolCriterion? IsMissingCriterion { get; set; }
    public StringCriterion? DuplicatedCriterion { get; set; }
    public StringCriterion? TitleCriterion { get; set; }
    public StringCriterion? CodeCriterion { get; set; }
    public StringCriterion? DetailsCriterion { get; set; }
    public StringCriterion? DirectorCriterion { get; set; }
    public IntCriterion? TagCountCriterion { get; set; }
    public IntCriterion? ResumeTimeCriterion { get; set; }
    public IntCriterion? PlayDurationCriterion { get; set; }
    public TimestampCriterion? LastPlayedAtCriterion { get; set; }
    public MultiIdCriterion? GalleriesCriterion { get; set; }
}

public class PerformerFilter
{
    public string? Name { get; set; }
    public bool? Favorite { get; set; }
    public int? Rating { get; set; }
    public List<int>? TagIds { get; set; }
    // Advanced criteria
    public IntCriterion? RatingCriterion { get; set; }
    public IntCriterion? AgeCriterion { get; set; }
    public StringCriterion? GenderCriterion { get; set; }
    public StringCriterion? EthnicityCriterion { get; set; }
    public StringCriterion? CountryCriterion { get; set; }
    public BoolCriterion? FavoriteCriterion { get; set; }
    public MultiIdCriterion? TagsCriterion { get; set; }
    public MultiIdCriterion? StudiosCriterion { get; set; }
    public IntCriterion? SceneCountCriterion { get; set; }
    public IntCriterion? ImageCountCriterion { get; set; }
    public IntCriterion? GalleryCountCriterion { get; set; }
    public DateCriterion? BirthdateCriterion { get; set; }
    public TimestampCriterion? CreatedAtCriterion { get; set; }
    public TimestampCriterion? UpdatedAtCriterion { get; set; }
    public StringCriterion? PathCriterion { get; set; }
    public StringCriterion? UrlCriterion { get; set; }
    public IntCriterion? WeightCriterion { get; set; }
    public IntCriterion? HeightCriterion { get; set; }
    public BoolCriterion? IsMissingCriterion { get; set; }
    public StringCriterion? StashIdCriterion { get; set; }
}

public class TagFilter
{
    public string? Name { get; set; }
    public bool? Favorite { get; set; }
    // Advanced criteria
    public BoolCriterion? FavoriteCriterion { get; set; }
    public IntCriterion? SceneCountCriterion { get; set; }
    public IntCriterion? MarkerCountCriterion { get; set; }
    public IntCriterion? PerformerCountCriterion { get; set; }
    public MultiIdCriterion? ParentsCriterion { get; set; }
    public MultiIdCriterion? ChildrenCriterion { get; set; }
    public BoolCriterion? IsMissingCriterion { get; set; }
    public TimestampCriterion? CreatedAtCriterion { get; set; }
    public TimestampCriterion? UpdatedAtCriterion { get; set; }
}

public class StudioFilter
{
    public string? Name { get; set; }
    public bool? Favorite { get; set; }
    public int? ParentId { get; set; }
    public List<int>? TagIds { get; set; }
    // Advanced criteria
    public IntCriterion? RatingCriterion { get; set; }
    public BoolCriterion? FavoriteCriterion { get; set; }
    public MultiIdCriterion? TagsCriterion { get; set; }
    public IntCriterion? SceneCountCriterion { get; set; }
    public IntCriterion? GalleryCountCriterion { get; set; }
    public IntCriterion? ImageCountCriterion { get; set; }
    public StringCriterion? UrlCriterion { get; set; }
    public StringCriterion? StashIdCriterion { get; set; }
    public BoolCriterion? IsMissingCriterion { get; set; }
    public TimestampCriterion? CreatedAtCriterion { get; set; }
    public TimestampCriterion? UpdatedAtCriterion { get; set; }
}

public class GalleryFilter
{
    public string? Title { get; set; }
    public int? Rating { get; set; }
    public bool? Organized { get; set; }
    public int? StudioId { get; set; }
    public List<int>? TagIds { get; set; }
    public List<int>? PerformerIds { get; set; }
    // Advanced criteria
    public IntCriterion? RatingCriterion { get; set; }
    public BoolCriterion? OrganizedCriterion { get; set; }
    public MultiIdCriterion? TagsCriterion { get; set; }
    public MultiIdCriterion? PerformersCriterion { get; set; }
    public MultiIdCriterion? StudiosCriterion { get; set; }
    public IntCriterion? ImageCountCriterion { get; set; }
    public DateCriterion? DateCriterion { get; set; }
    public StringCriterion? PathCriterion { get; set; }
    public StringCriterion? UrlCriterion { get; set; }
    public TimestampCriterion? CreatedAtCriterion { get; set; }
    public TimestampCriterion? UpdatedAtCriterion { get; set; }
    public BoolCriterion? PerformerFavoriteCriterion { get; set; }
    public BoolCriterion? IsMissingCriterion { get; set; }
}

public class ImageFilter
{
    public string? Title { get; set; }
    public int? Rating { get; set; }
    public bool? Organized { get; set; }
    public int? StudioId { get; set; }
    public int? GalleryId { get; set; }
    public List<int>? TagIds { get; set; }
    public List<int>? PerformerIds { get; set; }
    // Advanced criteria
    public IntCriterion? RatingCriterion { get; set; }
    public BoolCriterion? OrganizedCriterion { get; set; }
    public MultiIdCriterion? TagsCriterion { get; set; }
    public MultiIdCriterion? PerformersCriterion { get; set; }
    public MultiIdCriterion? StudiosCriterion { get; set; }
    public MultiIdCriterion? GalleriesCriterion { get; set; }
    public IntCriterion? OCounterCriterion { get; set; }
    public IntCriterion? ResolutionCriterion { get; set; }
    public StringCriterion? PathCriterion { get; set; }
    public TimestampCriterion? CreatedAtCriterion { get; set; }
    public TimestampCriterion? UpdatedAtCriterion { get; set; }
    public BoolCriterion? PerformerFavoriteCriterion { get; set; }
    public BoolCriterion? IsMissingCriterion { get; set; }
}

public class GroupFilter
{
    public string? Name { get; set; }
    public int? Rating { get; set; }
    public int? StudioId { get; set; }
    // Advanced criteria
    public IntCriterion? RatingCriterion { get; set; }
    public IntCriterion? DurationCriterion { get; set; }
    public MultiIdCriterion? StudiosCriterion { get; set; }
    public MultiIdCriterion? TagsCriterion { get; set; }
    public DateCriterion? DateCriterion { get; set; }
    public StringCriterion? UrlCriterion { get; set; }
    public TimestampCriterion? CreatedAtCriterion { get; set; }
    public TimestampCriterion? UpdatedAtCriterion { get; set; }
    public BoolCriterion? IsMissingCriterion { get; set; }
}
