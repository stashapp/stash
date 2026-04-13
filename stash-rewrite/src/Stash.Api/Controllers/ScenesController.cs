using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Stash.Api.Services;
using Stash.Core.DTOs;
using Stash.Core.Entities;
using Stash.Core.Interfaces;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class ScenesController(ISceneRepository sceneRepo, Data.StashContext db, StashBoxService stashBoxService) : ControllerBase
{
    [HttpGet]
    public async Task<ActionResult<PaginatedResponse<SceneDto>>> Find(
        [FromQuery] string? q, [FromQuery] int page = 1, [FromQuery] int perPage = 25,
        [FromQuery] string? sort = null, [FromQuery] string? direction = null,
        [FromQuery] string? title = null, [FromQuery] int? rating = null,
        [FromQuery] bool? organized = null, [FromQuery] int? studioId = null,
        [FromQuery] int? groupId = null, [FromQuery] int? galleryId = null, [FromQuery] string? tagIds = null, [FromQuery] string? performerIds = null,
        CancellationToken ct = default)
    {
        var filter = new SceneFilter
        {
            Title = title, Rating = rating, Organized = organized, StudioId = studioId, GroupId = groupId, GalleryId = galleryId,
            TagIds = ParseIntList(tagIds), PerformerIds = ParseIntList(performerIds)
        };
        var findFilter = new FindFilter
        {
            Q = q, Page = page, PerPage = perPage, Sort = sort,
            Direction = direction == "desc" ? Core.Enums.SortDirection.Desc : Core.Enums.SortDirection.Asc
        };

        var (items, totalCount) = await sceneRepo.FindAsync(filter, findFilter, ct);
        var dtos = items.Select(MapToDto).ToList();
        return Ok(new PaginatedResponse<SceneDto>(dtos, totalCount, page, perPage));
    }

    /// <summary>POST-based filtered query supporting advanced criteria (JSON body).</summary>
    [HttpPost("find")]
    public async Task<ActionResult<PaginatedResponse<SceneDto>>> FindPost([FromBody] FilteredQueryRequest<SceneFilter> req, CancellationToken ct)
    {
        var findFilter = req.FindFilter ?? new FindFilter();
        var filter = req.ObjectFilter ?? new SceneFilter();
        var (items, totalCount) = await sceneRepo.FindAsync(filter, findFilter, ct);
        var dtos = items.Select(MapToDto).ToList();
        return Ok(new PaginatedResponse<SceneDto>(dtos, totalCount, findFilter.Page, findFilter.PerPage));
    }

    [HttpGet("{id:int}")]
    public async Task<ActionResult<SceneDto>> GetById(int id, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdWithRelationsAsync(id, ct);
        if (scene == null) return NotFound();
        return Ok(MapToDto(scene));
    }

    [HttpPost]
    public async Task<ActionResult<SceneDto>> Create([FromBody] SceneCreateDto dto, CancellationToken ct)
    {
        var scene = new Scene
        {
            Title = dto.Title, Code = dto.Code, Details = dto.Details, Director = dto.Director,
            Date = ParseDate(dto.Date), Rating = dto.Rating, Organized = dto.Organized, StudioId = dto.StudioId
        };
        if (dto.Urls?.Count > 0)
            scene.Urls = dto.Urls.Select(u => new SceneUrl { Url = u }).ToList();
        if (dto.TagIds?.Count > 0)
            scene.SceneTags = dto.TagIds.Select(id => new SceneTag { TagId = id }).ToList();
        if (dto.PerformerIds?.Count > 0)
            scene.ScenePerformers = dto.PerformerIds.Select(id => new ScenePerformer { PerformerId = id }).ToList();
        if (dto.GalleryIds?.Count > 0)
            scene.SceneGalleries = dto.GalleryIds.Select(id => new SceneGallery { GalleryId = id }).ToList();

        scene = await sceneRepo.AddAsync(scene, ct);
        var result = await sceneRepo.GetByIdWithRelationsAsync(scene.Id, ct);
        return CreatedAtAction(nameof(GetById), new { id = scene.Id }, MapToDto(result!));
    }

    [HttpPut("{id:int}")]
    public async Task<ActionResult<SceneDto>> Update(int id, [FromBody] SceneUpdateDto dto, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdWithRelationsAsync(id, ct);
        if (scene == null) return NotFound();

        if (dto.Title != null) scene.Title = dto.Title;
        if (dto.Code != null) scene.Code = dto.Code;
        if (dto.Details != null) scene.Details = dto.Details;
        if (dto.Director != null) scene.Director = dto.Director;
        if (dto.Date != null) scene.Date = ParseDate(dto.Date);
        if (dto.Rating.HasValue) scene.Rating = dto.Rating;
        if (dto.Organized.HasValue) scene.Organized = dto.Organized.Value;
        if (dto.StudioId.HasValue) scene.StudioId = dto.StudioId;

        if (dto.Urls != null)
        {
            scene.Urls.Clear();
            scene.Urls = dto.Urls.Select(u => new SceneUrl { Url = u, SceneId = id }).ToList();
        }
        if (dto.TagIds != null)
        {
            scene.SceneTags.Clear();
            scene.SceneTags = dto.TagIds.Select(tid => new SceneTag { TagId = tid, SceneId = id }).ToList();
        }
        if (dto.PerformerIds != null)
        {
            scene.ScenePerformers.Clear();
            scene.ScenePerformers = dto.PerformerIds.Select(pid => new ScenePerformer { PerformerId = pid, SceneId = id }).ToList();
        }
        if (dto.GalleryIds != null)
        {
            scene.SceneGalleries.Clear();
            scene.SceneGalleries = dto.GalleryIds.Select(gid => new SceneGallery { GalleryId = gid, SceneId = id }).ToList();
        }
        if (dto.Groups != null)
        {
            scene.SceneGroups.Clear();
            scene.SceneGroups = dto.Groups.Select(g => new SceneGroup { GroupId = g.GroupId, SceneIndex = g.SceneIndex, SceneId = id }).ToList();
        }

        await sceneRepo.UpdateAsync(scene, ct);
        var updated = await sceneRepo.GetByIdWithRelationsAsync(id, ct);
        return Ok(MapToDto(updated!));
    }

    [HttpDelete("{id:int}")]
    public async Task<IActionResult> Delete(int id, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(id, ct);
        if (scene == null) return NotFound();
        await sceneRepo.DeleteAsync(id, ct);
        return NoContent();
    }

    [HttpGet("{id:int}/stash-box/search")]
    public async Task<ActionResult<IReadOnlyList<StashBoxSceneMatchDto>>> SearchStashBox(int id, [FromQuery] string? term, [FromQuery] string? endpoint, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdWithRelationsAsync(id, ct);
        if (scene == null) return NotFound();

        return Ok(await stashBoxService.SearchScenesAsync(scene, term, endpoint, ct));
    }

    [HttpPost("{id:int}/stash-box/import")]
    public async Task<ActionResult<SceneDto>> ImportFromStashBox(int id, [FromBody] StashBoxSceneImportRequestDto dto, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdWithRelationsAsync(id, ct);
        if (scene == null) return NotFound();

        var imported = await stashBoxService.MergeSceneAsync(scene, dto.Endpoint, dto.SceneId, dto, ct);
        if (!imported) return NotFound();

        await db.SaveChangesAsync(ct);
        var updated = await sceneRepo.GetByIdWithRelationsAsync(id, ct);
        return Ok(MapToDto(updated!));
    }

    private static SceneDto MapToDto(Scene s) => new(
        s.Id, s.Title, s.Code, s.Details, s.Director,
        s.Date?.ToString("yyyy-MM-dd"), s.Rating, s.Organized, s.StudioId, s.Studio?.Name,
        s.ResumeTime, s.PlayDuration, s.PlayCount, s.LastPlayedAt?.ToString("o"),
        s.OCounter,
        s.Urls.Select(u => u.Url).ToList(),
        s.SceneTags.Where(st => st.Tag != null).Select(st => new TagDto(st.Tag!.Id, st.Tag.Name, st.Tag.Description, st.Tag.Favorite, st.Tag.IgnoreAutoTag, [])).ToList(),
        s.ScenePerformers.Where(sp => sp.Performer != null).Select(sp => new PerformerSummaryDto(sp.Performer!.Id, sp.Performer.Name, sp.Performer.Disambiguation, sp.Performer.Gender?.ToString(), sp.Performer.Favorite, sp.Performer.ImageBlobId != null ? $"/api/performers/{sp.Performer.Id}/image" : null)).ToList(),
        s.Files.Select(f => new VideoFileDto(f.Id, f.Path, f.Basename, f.Format, f.Width, f.Height, f.Duration, f.VideoCodec, f.AudioCodec, f.FrameRate, f.BitRate, f.Size,
            f.Fingerprints.Select(fp => new FingerprintDto(fp.Type, fp.Value)).ToList())).ToList(),
        s.SceneMarkers.Select(m => new SceneMarkerSummaryDto(m.Id, m.Title, m.Seconds, m.EndSeconds, m.PrimaryTagId, m.PrimaryTag?.Name ?? "")).ToList(),
        s.SceneGroups.Where(sg => sg.Group != null).Select(sg => new GroupSummaryDto(sg.Group!.Id, sg.Group.Name, sg.SceneIndex)).ToList(),
        s.SceneGalleries.Where(sg => sg.Gallery != null).Select(sg => new GallerySummaryDto(sg.Gallery!.Id, sg.Gallery.Title, sg.Gallery.Date?.ToString("yyyy-MM-dd"))).ToList(),
        s.StashIds.Select(stashId => new SceneStashIdDto(stashId.Endpoint, stashId.StashId)).ToList(),
        s.CreatedAt.ToString("o"), s.UpdatedAt.ToString("o")
    );

    // ===== Scene Markers =====

    [HttpGet("{sceneId:int}/markers")]
    public async Task<ActionResult<List<SceneMarkerSummaryDto>>> GetMarkers(int sceneId, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(sceneId, ct);
        if (scene == null) return NotFound();

        var markers = await db.SceneMarkers
            .Include(m => m.PrimaryTag)
            .Where(m => m.SceneId == sceneId)
            .OrderBy(m => m.Seconds)
            .Select(m => new SceneMarkerSummaryDto(m.Id, m.Title, m.Seconds, m.EndSeconds, m.PrimaryTagId, m.PrimaryTag!.Name))
            .ToListAsync(ct);

        return Ok(markers);
    }

    [HttpPost("{sceneId:int}/markers")]
    public async Task<ActionResult<SceneMarkerSummaryDto>> CreateMarker(int sceneId, [FromBody] SceneMarkerCreateDto dto, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(sceneId, ct);
        if (scene == null) return NotFound();

        var marker = new SceneMarker
        {
            Title = dto.Title,
            Seconds = dto.Seconds,
            EndSeconds = dto.EndSeconds,
            PrimaryTagId = dto.PrimaryTagId,
            SceneId = sceneId
        };

        if (dto.TagIds?.Count > 0)
            marker.SceneMarkerTags = dto.TagIds.Select(tid => new SceneMarkerTag { TagId = tid }).ToList();

        db.SceneMarkers.Add(marker);
        await db.SaveChangesAsync(ct);

        await db.Entry(marker).Reference(m => m.PrimaryTag).LoadAsync(ct);
        return CreatedAtAction(nameof(GetMarkers), new { sceneId },
            new SceneMarkerSummaryDto(marker.Id, marker.Title, marker.Seconds, marker.EndSeconds, marker.PrimaryTagId, marker.PrimaryTag?.Name ?? ""));
    }

    [HttpPut("{sceneId:int}/markers/{markerId:int}")]
    public async Task<ActionResult<SceneMarkerSummaryDto>> UpdateMarker(int sceneId, int markerId, [FromBody] SceneMarkerUpdateDto dto, CancellationToken ct)
    {
        var marker = await db.SceneMarkers.Include(m => m.PrimaryTag).FirstOrDefaultAsync(m => m.Id == markerId && m.SceneId == sceneId, ct);
        if (marker == null) return NotFound();

        if (dto.Title != null) marker.Title = dto.Title;
        if (dto.Seconds.HasValue) marker.Seconds = dto.Seconds.Value;
        if (dto.EndSeconds.HasValue) marker.EndSeconds = dto.EndSeconds;
        if (dto.PrimaryTagId.HasValue) marker.PrimaryTagId = dto.PrimaryTagId.Value;

        await db.SaveChangesAsync(ct);

        if (dto.PrimaryTagId.HasValue) await db.Entry(marker).Reference(m => m.PrimaryTag).LoadAsync(ct);
        return Ok(new SceneMarkerSummaryDto(marker.Id, marker.Title, marker.Seconds, marker.EndSeconds, marker.PrimaryTagId, marker.PrimaryTag?.Name ?? ""));
    }

    [HttpDelete("{sceneId:int}/markers/{markerId:int}")]
    public async Task<IActionResult> DeleteMarker(int sceneId, int markerId, CancellationToken ct)
    {
        var marker = await db.SceneMarkers.FirstOrDefaultAsync(m => m.Id == markerId && m.SceneId == sceneId, ct);
        if (marker == null) return NotFound();

        db.SceneMarkers.Remove(marker);
        await db.SaveChangesAsync(ct);
        return NoContent();
    }

    // ===== Activity Tracking =====

    [HttpPost("{id:int}/play")]
    public async Task<IActionResult> RecordPlay(int id, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(id, ct);
        if (scene == null) return NotFound();

        scene.PlayCount++;
        scene.LastPlayedAt = DateTime.UtcNow;
        db.Set<ScenePlayHistory>().Add(new ScenePlayHistory { SceneId = id, PlayedAt = DateTime.UtcNow });
        await sceneRepo.UpdateAsync(scene, ct);
        return NoContent();
    }

    [HttpDelete("{id:int}/play")]
    public async Task<IActionResult> DeletePlay(int id, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(id, ct);
        if (scene == null) return NotFound();

        var last = await db.Set<ScenePlayHistory>().Where(h => h.SceneId == id).OrderByDescending(h => h.PlayedAt).FirstOrDefaultAsync(ct);
        if (last != null) { db.Set<ScenePlayHistory>().Remove(last); scene.PlayCount = Math.Max(0, scene.PlayCount - 1); }
        await sceneRepo.UpdateAsync(scene, ct);
        return NoContent();
    }

    [HttpPost("{id:int}/play/reset")]
    public async Task<IActionResult> ResetPlayCount(int id, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(id, ct);
        if (scene == null) return NotFound();

        scene.PlayCount = 0;
        scene.PlayDuration = 0;
        scene.LastPlayedAt = null;
        db.Set<ScenePlayHistory>().RemoveRange(db.Set<ScenePlayHistory>().Where(h => h.SceneId == id));
        await sceneRepo.UpdateAsync(scene, ct);
        return NoContent();
    }

    [HttpPost("{id:int}/o")]
    public async Task<IActionResult> IncrementO(int id, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(id, ct);
        if (scene == null) return NotFound();

        scene.OCounter++;
        db.Set<SceneOHistory>().Add(new SceneOHistory { SceneId = id, OccurredAt = DateTime.UtcNow });
        await sceneRepo.UpdateAsync(scene, ct);
        return NoContent();
    }

    [HttpDelete("{id:int}/o")]
    public async Task<IActionResult> DecrementO(int id, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(id, ct);
        if (scene == null) return NotFound();

        var last = await db.Set<SceneOHistory>().Where(h => h.SceneId == id).OrderByDescending(h => h.OccurredAt).FirstOrDefaultAsync(ct);
        if (last != null) { db.Set<SceneOHistory>().Remove(last); scene.OCounter = Math.Max(0, scene.OCounter - 1); }
        await sceneRepo.UpdateAsync(scene, ct);
        return NoContent();
    }

    [HttpPost("{id:int}/o/reset")]
    public async Task<IActionResult> ResetO(int id, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(id, ct);
        if (scene == null) return NotFound();

        scene.OCounter = 0;
        db.Set<SceneOHistory>().RemoveRange(db.Set<SceneOHistory>().Where(h => h.SceneId == id));
        await sceneRepo.UpdateAsync(scene, ct);
        return NoContent();
    }

    [HttpPost("{id:int}/activity")]
    public async Task<IActionResult> SaveActivity(int id, [FromBody] SceneActivityDto dto, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(id, ct);
        if (scene == null) return NotFound();

        if (dto.ResumeTime.HasValue) scene.ResumeTime = dto.ResumeTime.Value;
        if (dto.PlayDuration.HasValue) scene.PlayDuration += dto.PlayDuration.Value;
        await sceneRepo.UpdateAsync(scene, ct);
        return NoContent();
    }

    [HttpPost("{id:int}/activity/reset")]
    public async Task<IActionResult> ResetActivity(int id, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(id, ct);
        if (scene == null) return NotFound();

        scene.ResumeTime = 0;
        scene.PlayDuration = 0;
        await sceneRepo.UpdateAsync(scene, ct);
        return NoContent();
    }

    // ===== Scene Wall/Discovery =====

    [HttpGet("wall")]
    public async Task<ActionResult<List<SceneDto>>> SceneWall([FromQuery] string? q, [FromQuery] int count = 24, CancellationToken ct = default)
    {
        var query = db.Scenes
            .Include(s => s.Files).ThenInclude(f => f.Fingerprints)
            .Include(s => s.SceneTags).ThenInclude(st => st.Tag)
            .Include(s => s.ScenePerformers).ThenInclude(sp => sp.Performer)
            .Include(s => s.Studio)
            .AsNoTracking();

        if (!string.IsNullOrEmpty(q))
            query = query.Where(s => s.Title != null && EF.Functions.ILike(s.Title, $"%{q}%"));

        var scenes = await query.OrderBy(_ => EF.Functions.Random()).Take(count).ToListAsync(ct);
        return Ok(scenes.Select(MapToDto).ToList());
    }

    [HttpGet("duplicates")]
    public async Task<ActionResult<List<List<SceneDto>>>> FindDuplicates([FromQuery] int distance = 0, [FromQuery] double? durationDiff = null, CancellationToken ct = default)
    {
        // Group by oshash fingerprint to find exact duplicates
        var fingerprints = await db.Set<FileFingerprint>()
            .Where(f => f.Type == "oshash")
            .GroupBy(f => f.Value)
            .Where(g => g.Count() > 1)
            .Select(g => g.Key)
            .ToListAsync(ct);

        var result = new List<List<SceneDto>>();
        foreach (var hash in fingerprints)
        {
            var fileIds = await db.Set<FileFingerprint>()
                .Where(f => f.Type == "oshash" && f.Value == hash)
                .Select(f => f.FileId)
                .ToListAsync(ct);

            var scenes = await db.Scenes
                .Include(s => s.Files).ThenInclude(f => f.Fingerprints)
                .Include(s => s.SceneTags).ThenInclude(st => st.Tag)
                .Include(s => s.ScenePerformers).ThenInclude(sp => sp.Performer)
                .Include(s => s.Studio)
                .Where(s => s.Files.Any(f => fileIds.Contains(f.Id)))
                .AsNoTracking()
                .ToListAsync(ct);

            if (scenes.Count > 1)
                result.Add(scenes.Select(MapToDto).ToList());
        }
        return Ok(result);
    }

    // ===== Bulk Operations =====

    [HttpPost("bulk")]
    public async Task<IActionResult> BulkUpdate([FromBody] BulkSceneUpdateDto dto, CancellationToken ct)
    {
        var scenes = await db.Scenes
            .Include(s => s.SceneTags)
            .Include(s => s.ScenePerformers)
            .Where(s => dto.Ids.Contains(s.Id))
            .ToListAsync(ct);

        foreach (var scene in scenes)
        {
            if (dto.Rating.HasValue) scene.Rating = dto.Rating;
            if (dto.Organized.HasValue) scene.Organized = dto.Organized.Value;
            if (dto.StudioId.HasValue) scene.StudioId = dto.StudioId;

            if (dto.TagIds != null && dto.TagMode == BulkUpdateMode.Set)
            {
                scene.SceneTags.Clear();
                scene.SceneTags = dto.TagIds.Select(tid => new SceneTag { TagId = tid, SceneId = scene.Id }).ToList();
            }
            else if (dto.TagIds != null && dto.TagMode == BulkUpdateMode.Add)
            {
                var existing = scene.SceneTags.Select(st => st.TagId).ToHashSet();
                foreach (var tid in dto.TagIds.Where(t => !existing.Contains(t)))
                    scene.SceneTags.Add(new SceneTag { TagId = tid, SceneId = scene.Id });
            }
            else if (dto.TagIds != null && dto.TagMode == BulkUpdateMode.Remove)
            {
                scene.SceneTags = scene.SceneTags.Where(st => !dto.TagIds.Contains(st.TagId)).ToList();
            }

            if (dto.PerformerIds != null && dto.PerformerMode == BulkUpdateMode.Set)
            {
                scene.ScenePerformers.Clear();
                scene.ScenePerformers = dto.PerformerIds.Select(pid => new ScenePerformer { PerformerId = pid, SceneId = scene.Id }).ToList();
            }
            else if (dto.PerformerIds != null && dto.PerformerMode == BulkUpdateMode.Add)
            {
                var existing = scene.ScenePerformers.Select(sp => sp.PerformerId).ToHashSet();
                foreach (var pid in dto.PerformerIds.Where(p => !existing.Contains(p)))
                    scene.ScenePerformers.Add(new ScenePerformer { PerformerId = pid, SceneId = scene.Id });
            }
            else if (dto.PerformerIds != null && dto.PerformerMode == BulkUpdateMode.Remove)
            {
                scene.ScenePerformers = scene.ScenePerformers.Where(sp => !dto.PerformerIds.Contains(sp.PerformerId)).ToList();
            }
        }

        await db.SaveChangesAsync(ct);
        return Ok(new { updated = scenes.Count });
    }

    // ===== Merge =====

    [HttpPost("merge")]
    public async Task<ActionResult<SceneDto>> MergeScenes([FromBody] SceneMergeDto dto, CancellationToken ct)
    {
        var target = await sceneRepo.GetByIdWithRelationsAsync(dto.TargetId, ct);
        if (target == null) return NotFound("Target scene not found");

        var sources = await db.Scenes
            .Include(s => s.Files)
            .Include(s => s.SceneTags)
            .Include(s => s.ScenePerformers)
            .Include(s => s.SceneGalleries)
            .Include(s => s.Urls)
            .Where(s => dto.SourceIds.Contains(s.Id))
            .ToListAsync(ct);

        var existingTagIds = target.SceneTags.Select(st => st.TagId).ToHashSet();
        var existingPerfIds = target.ScenePerformers.Select(sp => sp.PerformerId).ToHashSet();

        foreach (var source in sources)
        {
            // Move files to target
            foreach (var f in source.Files) f.SceneId = target.Id;
            // Merge tags
            foreach (var st in source.SceneTags.Where(st => !existingTagIds.Contains(st.TagId)))
                target.SceneTags.Add(new SceneTag { TagId = st.TagId, SceneId = target.Id });
            // Merge performers
            foreach (var sp in source.ScenePerformers.Where(sp => !existingPerfIds.Contains(sp.PerformerId)))
                target.ScenePerformers.Add(new ScenePerformer { PerformerId = sp.PerformerId, SceneId = target.Id });
            // Accumulate play counts & o-counters
            target.PlayCount += source.PlayCount;
            target.OCounter += source.OCounter;
            target.PlayDuration += source.PlayDuration;
            // Delete source
            db.Scenes.Remove(source);
        }

        await db.SaveChangesAsync(ct);
        var result = await sceneRepo.GetByIdWithRelationsAsync(target.Id, ct);
        return Ok(MapToDto(result!));
    }

    private static DateOnly? ParseDate(string? date) => DateOnly.TryParse(date, out var d) ? d : null;
    private static List<int>? ParseIntList(string? csv) => string.IsNullOrEmpty(csv) ? null : csv.Split(',').Select(int.Parse).ToList();
}
