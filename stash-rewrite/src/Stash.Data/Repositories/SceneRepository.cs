using Microsoft.EntityFrameworkCore;
using Stash.Core.Entities;
using Stash.Core.Interfaces;

namespace Stash.Data.Repositories;

public class SceneRepository : ISceneRepository
{
    private readonly StashContext _db;
    public SceneRepository(StashContext db) => _db = db;

    public async Task<Scene?> GetByIdAsync(int id, CancellationToken ct = default)
        => await _db.Scenes.FindAsync([id], ct);

    public async Task<Scene?> GetByIdWithRelationsAsync(int id, CancellationToken ct = default)
        => await _db.Scenes
            .Include(s => s.Studio)
            .Include(s => s.Urls)
            .Include(s => s.SceneTags).ThenInclude(st => st.Tag)
            .Include(s => s.ScenePerformers).ThenInclude(sp => sp.Performer)
            .Include(s => s.SceneGalleries).ThenInclude(sg => sg.Gallery)
            .Include(s => s.SceneGroups).ThenInclude(sg => sg.Group)
            .Include(s => s.Files).ThenInclude(f => f.Fingerprints)
            .Include(s => s.SceneMarkers)
            .Include(s => s.StashIds)
            .AsSplitQuery()
            .FirstOrDefaultAsync(s => s.Id == id, ct);

    public async Task<IReadOnlyList<Scene>> GetAllAsync(CancellationToken ct = default)
        => await _db.Scenes.AsNoTracking().ToListAsync(ct);

    public async Task<Scene> AddAsync(Scene entity, CancellationToken ct = default)
    {
        _db.Scenes.Add(entity);
        await _db.SaveChangesAsync(ct);
        return entity;
    }

    public async Task UpdateAsync(Scene entity, CancellationToken ct = default)
    {
        _db.Scenes.Update(entity);
        await _db.SaveChangesAsync(ct);
    }

    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var entity = await _db.Scenes.FindAsync([id], ct);
        if (entity != null)
        {
            _db.Scenes.Remove(entity);
            await _db.SaveChangesAsync(ct);
        }
    }

    public async Task<int> CountAsync(CancellationToken ct = default)
        => await _db.Scenes.CountAsync(ct);

    public async Task<(IReadOnlyList<Scene> Items, int TotalCount)> FindAsync(SceneFilter? filter, FindFilter? findFilter, CancellationToken ct = default)
    {
        var query = _db.Scenes
            .Include(s => s.Studio)
            .Include(s => s.SceneTags).ThenInclude(st => st.Tag)
            .Include(s => s.ScenePerformers).ThenInclude(sp => sp.Performer)
            .Include(s => s.SceneGalleries).ThenInclude(sg => sg.Gallery)
            .Include(s => s.Files).ThenInclude(file => file.Fingerprints)
            .Include(s => s.Files).ThenInclude(file => file.ParentFolder)
            .Include(s => s.SceneMarkers)
            .AsSplitQuery()
            .AsQueryable();

        // Apply basic filters (backward compat)
        if (filter != null)
        {
            if (!string.IsNullOrEmpty(filter.Title))
                query = query.Where(s => s.Title != null && EF.Functions.ILike(s.Title, $"%{filter.Title}%"));
            if (filter.Rating.HasValue)
                query = query.Where(s => s.Rating >= filter.Rating.Value);
            if (filter.Organized.HasValue)
                query = query.Where(s => s.Organized == filter.Organized.Value);
            if (filter.StudioId.HasValue)
                query = query.Where(s => s.StudioId == filter.StudioId.Value);
            if (filter.GroupId.HasValue)
                query = query.Where(s => s.SceneGroups.Any(sg => sg.GroupId == filter.GroupId.Value));
            if (filter.GalleryId.HasValue)
                query = query.Where(s => s.SceneGalleries.Any(sg => sg.GalleryId == filter.GalleryId.Value));
            if (filter.TagIds?.Count > 0)
                query = query.Where(s => s.SceneTags.Any(st => filter.TagIds.Contains(st.TagId)));
            if (filter.PerformerIds?.Count > 0)
                query = query.Where(s => s.ScenePerformers.Any(sp => filter.PerformerIds.Contains(sp.PerformerId)));

            // Advanced criteria
            query = ApplyIntCriterion(query, filter.RatingCriterion, s => s.Rating ?? 0);
            query = ApplyIntCriterion(query, filter.OCounterCriterion, s => s.OCounter);
            query = ApplyIntCriterion(query, filter.PlayCountCriterion, s => s.PlayCount);

            if (filter.PerformerCountCriterion != null)
                query = ApplyIntCriterion(query, filter.PerformerCountCriterion, s => s.ScenePerformers.Count);

            if (filter.DurationCriterion != null)
                query = ApplyIntCriterion(query, filter.DurationCriterion, s => (int)(s.Files.Select(f => f.Duration).Max()));

            if (filter.ResolutionCriterion != null)
                query = ApplyIntCriterion(query, filter.ResolutionCriterion, s => s.Files.Select(f => f.Height).Max());

            if (filter.FrameRateCriterion != null)
                query = ApplyIntCriterion(query, filter.FrameRateCriterion, s => (int)(s.Files.Select(f => f.FrameRate).Max()));

            if (filter.BitrateInterval != null)
                query = ApplyIntCriterion(query, filter.BitrateInterval, s => (int)(s.Files.Select(f => f.BitRate).Max() / 1000));

            if (filter.FileCountCriterion != null)
                query = ApplyIntCriterion(query, filter.FileCountCriterion, s => s.Files.Count);

            query = ApplyMultiIdCriterion(query, filter.TagsCriterion, s => s.SceneTags.Select(st => st.TagId));
            query = ApplyMultiIdCriterion(query, filter.PerformersCriterion, s => s.ScenePerformers.Select(sp => sp.PerformerId));

            if (filter.StudiosCriterion != null)
            {
                var ids = filter.StudiosCriterion.Value;
                query = filter.StudiosCriterion.Modifier switch
                {
                    CriterionModifier.Includes => query.Where(s => s.StudioId.HasValue && ids.Contains(s.StudioId.Value)),
                    CriterionModifier.Excludes => query.Where(s => !s.StudioId.HasValue || !ids.Contains(s.StudioId.Value)),
                    _ => query.Where(s => s.StudioId.HasValue && ids.Contains(s.StudioId.Value)),
                };
            }

            query = ApplyMultiIdCriterion(query, filter.GroupsCriterion, s => s.SceneGroups.Select(sg => sg.GroupId));

            if (filter.OrganizedCriterion != null)
                query = query.Where(s => s.Organized == filter.OrganizedCriterion.Value);

            if (filter.HasMarkersCriterion != null)
                query = filter.HasMarkersCriterion.Value
                    ? query.Where(s => s.SceneMarkers.Count > 0)
                    : query.Where(s => s.SceneMarkers.Count == 0);

            if (filter.InteractiveCriterion != null)
                query = query.Where(s => s.Files.Any(f => f.Interactive == filter.InteractiveCriterion.Value));

            if (filter.PathCriterion != null)
            {
                var val = filter.PathCriterion.Value;
                query = filter.PathCriterion.Modifier switch
                {
                    CriterionModifier.Equals => query.Where(s => s.Files.Any(f => f.Basename == val)),
                    CriterionModifier.NotEquals => query.Where(s => !s.Files.Any(f => f.Basename == val)),
                    CriterionModifier.Includes => query.Where(s => s.Files.Any(f => EF.Functions.ILike(f.Basename, $"%{val}%"))),
                    CriterionModifier.Excludes => query.Where(s => !s.Files.Any(f => EF.Functions.ILike(f.Basename, $"%{val}%"))),
                    CriterionModifier.MatchesRegex => query.Where(s => s.Files.Any(f => EF.Functions.ILike(f.Basename, $"%{val}%"))),
                    CriterionModifier.NotMatchesRegex => query.Where(s => !s.Files.Any(f => EF.Functions.ILike(f.Basename, $"%{val}%"))),
                    _ => query,
                };
            }

            if (filter.VideoCodecCriterion != null)
            {
                var val = filter.VideoCodecCriterion.Value;
                query = filter.VideoCodecCriterion.Modifier switch
                {
                    CriterionModifier.Equals => query.Where(s => s.Files.Any(f => f.VideoCodec == val)),
                    CriterionModifier.NotEquals => query.Where(s => !s.Files.Any(f => f.VideoCodec == val)),
                    _ => query.Where(s => s.Files.Any(f => EF.Functions.ILike(f.VideoCodec, $"%{val}%"))),
                };
            }

            if (filter.AudioCodecCriterion != null)
            {
                var val = filter.AudioCodecCriterion.Value;
                query = filter.AudioCodecCriterion.Modifier switch
                {
                    CriterionModifier.Equals => query.Where(s => s.Files.Any(f => f.AudioCodec == val)),
                    CriterionModifier.NotEquals => query.Where(s => !s.Files.Any(f => f.AudioCodec == val)),
                    _ => query.Where(s => s.Files.Any(f => EF.Functions.ILike(f.AudioCodec, $"%{val}%"))),
                };
            }

            if (filter.DateCriterion != null)
            {
                var crit = filter.DateCriterion;
                if (DateOnly.TryParse(crit.Value, out var d1))
                {
                    DateOnly.TryParse(crit.Value2, out var d2);
                    query = crit.Modifier switch
                    {
                        CriterionModifier.Equals => query.Where(s => s.Date == d1),
                        CriterionModifier.NotEquals => query.Where(s => s.Date != d1),
                        CriterionModifier.GreaterThan => query.Where(s => s.Date > d1),
                        CriterionModifier.LessThan => query.Where(s => s.Date < d1),
                        CriterionModifier.Between => query.Where(s => s.Date >= d1 && s.Date <= d2),
                        CriterionModifier.NotBetween => query.Where(s => s.Date < d1 || s.Date > d2),
                        CriterionModifier.IsNull => query.Where(s => s.Date == null),
                        CriterionModifier.NotNull => query.Where(s => s.Date != null),
                        _ => query,
                    };
                }
            }

            if (filter.PerformerFavoriteCriterion != null)
                query = filter.PerformerFavoriteCriterion.Value
                    ? query.Where(s => s.ScenePerformers.Any(sp => sp.Performer!.Favorite))
                    : query.Where(s => !s.ScenePerformers.Any(sp => sp.Performer!.Favorite));

            if (filter.StashIdCriterion != null)
            {
                query = filter.StashIdCriterion.Modifier switch
                {
                    CriterionModifier.IsNull => query.Where(s => s.StashIds.Count == 0),
                    CriterionModifier.NotNull => query.Where(s => s.StashIds.Count > 0),
                    _ => query.Where(s => s.StashIds.Any(sid => EF.Functions.ILike(sid.Endpoint, $"%{filter.StashIdCriterion.Value}%"))),
                };
            }

            // Title criterion
            if (filter.TitleCriterion != null)
            {
                var val = filter.TitleCriterion.Value;
                query = filter.TitleCriterion.Modifier switch
                {
                    CriterionModifier.Equals => query.Where(s => s.Title == val),
                    CriterionModifier.NotEquals => query.Where(s => s.Title != val),
                    CriterionModifier.Includes => query.Where(s => s.Title != null && EF.Functions.ILike(s.Title, $"%{val}%")),
                    CriterionModifier.Excludes => query.Where(s => s.Title == null || !EF.Functions.ILike(s.Title, $"%{val}%")),
                    CriterionModifier.IsNull => query.Where(s => s.Title == null || s.Title == ""),
                    CriterionModifier.NotNull => query.Where(s => s.Title != null && s.Title != ""),
                    _ => query,
                };
            }

            // Code criterion
            if (filter.CodeCriterion != null)
            {
                var val = filter.CodeCriterion.Value;
                query = filter.CodeCriterion.Modifier switch
                {
                    CriterionModifier.Equals => query.Where(s => s.Code == val),
                    CriterionModifier.NotEquals => query.Where(s => s.Code != val),
                    CriterionModifier.Includes => query.Where(s => s.Code != null && EF.Functions.ILike(s.Code, $"%{val}%")),
                    CriterionModifier.Excludes => query.Where(s => s.Code == null || !EF.Functions.ILike(s.Code, $"%{val}%")),
                    CriterionModifier.IsNull => query.Where(s => s.Code == null || s.Code == ""),
                    CriterionModifier.NotNull => query.Where(s => s.Code != null && s.Code != ""),
                    _ => query,
                };
            }

            // Details criterion
            if (filter.DetailsCriterion != null)
            {
                var val = filter.DetailsCriterion.Value;
                query = filter.DetailsCriterion.Modifier switch
                {
                    CriterionModifier.Includes => query.Where(s => s.Details != null && EF.Functions.ILike(s.Details, $"%{val}%")),
                    CriterionModifier.Excludes => query.Where(s => s.Details == null || !EF.Functions.ILike(s.Details, $"%{val}%")),
                    CriterionModifier.IsNull => query.Where(s => s.Details == null || s.Details == ""),
                    CriterionModifier.NotNull => query.Where(s => s.Details != null && s.Details != ""),
                    _ => query,
                };
            }

            // Director criterion
            if (filter.DirectorCriterion != null)
            {
                var val = filter.DirectorCriterion.Value;
                query = filter.DirectorCriterion.Modifier switch
                {
                    CriterionModifier.Equals => query.Where(s => s.Director == val),
                    CriterionModifier.NotEquals => query.Where(s => s.Director != val),
                    CriterionModifier.Includes => query.Where(s => s.Director != null && EF.Functions.ILike(s.Director, $"%{val}%")),
                    CriterionModifier.Excludes => query.Where(s => s.Director == null || !EF.Functions.ILike(s.Director, $"%{val}%")),
                    CriterionModifier.IsNull => query.Where(s => s.Director == null || s.Director == ""),
                    CriterionModifier.NotNull => query.Where(s => s.Director != null && s.Director != ""),
                    _ => query,
                };
            }

            // Tag count criterion
            if (filter.TagCountCriterion != null)
                query = ApplyIntCriterion(query, filter.TagCountCriterion, s => s.SceneTags.Count);

            // Resume time criterion
            if (filter.ResumeTimeCriterion != null)
                query = ApplyIntCriterion(query, filter.ResumeTimeCriterion, s => (int)s.ResumeTime);

            // Play duration criterion
            if (filter.PlayDurationCriterion != null)
                query = ApplyIntCriterion(query, filter.PlayDurationCriterion, s => (int)s.PlayDuration);

            // Galleries criterion
            if (filter.GalleriesCriterion != null)
                query = ApplyMultiIdCriterion(query, filter.GalleriesCriterion, s => s.SceneGalleries.Select(sg => sg.GalleryId));
        }

        // Apply text search
        if (findFilter != null && !string.IsNullOrEmpty(findFilter.Q))
        {
            var q = findFilter.Q;
            query = query.Where(s =>
                (s.Title != null && EF.Functions.ILike(s.Title, $"%{q}%")) ||
                (s.Details != null && EF.Functions.ILike(s.Details, $"%{q}%")) ||
                (s.Code != null && EF.Functions.ILike(s.Code, $"%{q}%")) ||
                s.Files.Any(f => EF.Functions.ILike(f.Basename, $"%{q}%")));
        }

        var totalCount = await query.CountAsync(ct);

        // Apply sorting
        var sort = findFilter?.Sort ?? "updated_at";
        var desc = findFilter?.Direction == Core.Enums.SortDirection.Desc;
        query = sort switch
        {
            "title" => desc ? query.OrderByDescending(s => s.Title) : query.OrderBy(s => s.Title),
            "date" => desc ? query.OrderByDescending(s => s.Date) : query.OrderBy(s => s.Date),
            "rating" => desc ? query.OrderByDescending(s => s.Rating) : query.OrderBy(s => s.Rating),
            "play_count" => desc ? query.OrderByDescending(s => s.PlayCount) : query.OrderBy(s => s.PlayCount),
            "o_counter" => desc ? query.OrderByDescending(s => s.OCounter) : query.OrderBy(s => s.OCounter),
            "organized" => desc ? query.OrderByDescending(s => s.Organized) : query.OrderBy(s => s.Organized),
            "last_played_at" => desc ? query.OrderByDescending(s => s.LastPlayedAt) : query.OrderBy(s => s.LastPlayedAt),
            "play_duration" => desc ? query.OrderByDescending(s => s.PlayDuration) : query.OrderBy(s => s.PlayDuration),
            "resume_time" => desc ? query.OrderByDescending(s => s.ResumeTime) : query.OrderBy(s => s.ResumeTime),
            "random" => query.OrderBy(_ => EF.Functions.Random()),
            "duration" => desc
                ? query.OrderByDescending(s => s.Files.Select(file => (double?)file.Duration).Max() ?? 0)
                : query.OrderBy(s => s.Files.Select(file => (double?)file.Duration).Max() ?? 0),
            "file_size" => desc
                ? query.OrderByDescending(s => s.Files.Select(file => (long?)file.Size).Max() ?? 0)
                : query.OrderBy(s => s.Files.Select(file => (long?)file.Size).Max() ?? 0),
            "file_count" => desc
                ? query.OrderByDescending(s => s.Files.Count)
                : query.OrderBy(s => s.Files.Count),
            "resolution" => desc
                ? query.OrderByDescending(s => s.Files.Select(file => file.Height).Max())
                : query.OrderBy(s => s.Files.Select(file => file.Height).Max()),
            "framerate" => desc
                ? query.OrderByDescending(s => s.Files.Select(file => file.FrameRate).Max())
                : query.OrderBy(s => s.Files.Select(file => file.FrameRate).Max()),
            "bitrate" => desc
                ? query.OrderByDescending(s => s.Files.Select(file => file.BitRate).Max())
                : query.OrderBy(s => s.Files.Select(file => file.BitRate).Max()),
            "tag_count" => desc
                ? query.OrderByDescending(s => s.SceneTags.Count)
                : query.OrderBy(s => s.SceneTags.Count),
            "performer_count" => desc
                ? query.OrderByDescending(s => s.ScenePerformers.Count)
                : query.OrderBy(s => s.ScenePerformers.Count),
            "created_at" => desc ? query.OrderByDescending(s => s.CreatedAt) : query.OrderBy(s => s.CreatedAt),
            _ => desc ? query.OrderByDescending(s => s.UpdatedAt) : query.OrderBy(s => s.UpdatedAt),
        };

        // Apply pagination
        var page = findFilter?.Page ?? 1;
        var perPage = findFilter?.PerPage ?? 25;
        var items = await query.Skip((page - 1) * perPage).Take(perPage).AsNoTracking().ToListAsync(ct);

        return (items, totalCount);
    }

    // Helper methods for criterion-based filtering
    private static IQueryable<Scene> ApplyIntCriterion(IQueryable<Scene> query, IntCriterion? criterion, System.Linq.Expressions.Expression<Func<Scene, int>> selector)
    {
        if (criterion == null) return query;
        var val = criterion.Value;
        var val2 = criterion.Value2 ?? val;
        var param = selector.Parameters[0];
        var body = selector.Body;

        return criterion.Modifier switch
        {
            CriterionModifier.Equals => query.Where(System.Linq.Expressions.Expression.Lambda<Func<Scene, bool>>(
                System.Linq.Expressions.Expression.Equal(body, System.Linq.Expressions.Expression.Constant(val)), param)),
            CriterionModifier.NotEquals => query.Where(System.Linq.Expressions.Expression.Lambda<Func<Scene, bool>>(
                System.Linq.Expressions.Expression.NotEqual(body, System.Linq.Expressions.Expression.Constant(val)), param)),
            CriterionModifier.GreaterThan => query.Where(System.Linq.Expressions.Expression.Lambda<Func<Scene, bool>>(
                System.Linq.Expressions.Expression.GreaterThan(body, System.Linq.Expressions.Expression.Constant(val)), param)),
            CriterionModifier.LessThan => query.Where(System.Linq.Expressions.Expression.Lambda<Func<Scene, bool>>(
                System.Linq.Expressions.Expression.LessThan(body, System.Linq.Expressions.Expression.Constant(val)), param)),
            CriterionModifier.Between => query.Where(System.Linq.Expressions.Expression.Lambda<Func<Scene, bool>>(
                System.Linq.Expressions.Expression.AndAlso(
                    System.Linq.Expressions.Expression.GreaterThanOrEqual(body, System.Linq.Expressions.Expression.Constant(val)),
                    System.Linq.Expressions.Expression.LessThanOrEqual(body, System.Linq.Expressions.Expression.Constant(val2))), param)),
            CriterionModifier.NotBetween => query.Where(System.Linq.Expressions.Expression.Lambda<Func<Scene, bool>>(
                System.Linq.Expressions.Expression.OrElse(
                    System.Linq.Expressions.Expression.LessThan(body, System.Linq.Expressions.Expression.Constant(val)),
                    System.Linq.Expressions.Expression.GreaterThan(body, System.Linq.Expressions.Expression.Constant(val2))), param)),
            _ => query,
        };
    }

    private static IQueryable<Scene> ApplyMultiIdCriterion(IQueryable<Scene> query, MultiIdCriterion? criterion, System.Linq.Expressions.Expression<Func<Scene, IEnumerable<int>>> idsSelector)
    {
        if (criterion == null || (criterion.Value.Count == 0 && (criterion.Excludes == null || criterion.Excludes.Count == 0))) return query;

        var sceneParam = idsSelector.Parameters[0];
        var sceneIds = idsSelector.Body;

        // Apply inclusion filter on Value list
        if (criterion.Value.Count > 0)
        {
            var ids = criterion.Value;
            var selectedIds = System.Linq.Expressions.Expression.Constant(ids);
            var sceneIdParam = System.Linq.Expressions.Expression.Parameter(typeof(int), "sceneId");
            var selectedIdParam = System.Linq.Expressions.Expression.Parameter(typeof(int), "selectedId");

            var anySelectedInScene = System.Linq.Expressions.Expression.Call(
                typeof(Enumerable),
                nameof(Enumerable.Any),
                [typeof(int)],
                sceneIds,
                System.Linq.Expressions.Expression.Lambda<Func<int, bool>>(
                    System.Linq.Expressions.Expression.Call(
                        typeof(Enumerable),
                        nameof(Enumerable.Contains),
                        [typeof(int)],
                        selectedIds,
                        sceneIdParam),
                    sceneIdParam));

            var allSelectedInScene = System.Linq.Expressions.Expression.Call(
                typeof(Enumerable),
                nameof(Enumerable.All),
                [typeof(int)],
                selectedIds,
                System.Linq.Expressions.Expression.Lambda<Func<int, bool>>(
                    System.Linq.Expressions.Expression.Call(
                        typeof(Enumerable),
                        nameof(Enumerable.Contains),
                        [typeof(int)],
                        sceneIds,
                        selectedIdParam),
                    selectedIdParam));

            System.Linq.Expressions.Expression body = criterion.Modifier switch
            {
                CriterionModifier.Includes => anySelectedInScene,
                CriterionModifier.Excludes => System.Linq.Expressions.Expression.Not(anySelectedInScene),
                CriterionModifier.IncludesAll => allSelectedInScene,
                CriterionModifier.ExcludesAll => System.Linq.Expressions.Expression.Not(allSelectedInScene),
                _ => anySelectedInScene,
            };

            query = query.Where(System.Linq.Expressions.Expression.Lambda<Func<Scene, bool>>(body, sceneParam));
        }

        // Apply exclusion filter on Excludes list (always excludes, regardless of modifier)
        if (criterion.Excludes?.Count > 0)
        {
            var excludeIds = criterion.Excludes;
            var excludeConst = System.Linq.Expressions.Expression.Constant(excludeIds);
            var excludeIdParam = System.Linq.Expressions.Expression.Parameter(typeof(int), "excludeId");

            var anyExcludedInScene = System.Linq.Expressions.Expression.Call(
                typeof(Enumerable),
                nameof(Enumerable.Any),
                [typeof(int)],
                sceneIds,
                System.Linq.Expressions.Expression.Lambda<Func<int, bool>>(
                    System.Linq.Expressions.Expression.Call(
                        typeof(Enumerable),
                        nameof(Enumerable.Contains),
                        [typeof(int)],
                        excludeConst,
                        excludeIdParam),
                    excludeIdParam));

            var notAny = System.Linq.Expressions.Expression.Not(anyExcludedInScene);
            query = query.Where(System.Linq.Expressions.Expression.Lambda<Func<Scene, bool>>(notAny, sceneParam));
        }

        return query;
    }
}
