using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Stash.Api.Services;
using Stash.Core.DTOs;
using Stash.Core.Entities;
using Stash.Core.Enums;
using Stash.Core.Interfaces;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class StudiosController(IStudioRepository studioRepo, StashBoxService stashBoxService, Data.StashContext db) : ControllerBase
{
    [HttpGet]
    public async Task<ActionResult<PaginatedResponse<StudioDto>>> Find(
        [FromQuery] string? q, [FromQuery] int page = 1, [FromQuery] int perPage = 25,
        [FromQuery] string? sort = null, [FromQuery] string? direction = null,
        [FromQuery] string? name = null, [FromQuery] bool? favorite = null,
        [FromQuery] int? parentId = null, [FromQuery] string? tagIds = null,
        CancellationToken ct = default)
    {
        var filter = new StudioFilter { Name = name, Favorite = favorite, ParentId = parentId, TagIds = ParseIntList(tagIds) };
        var findFilter = new FindFilter
        {
            Q = q, Page = page, PerPage = perPage, Sort = sort,
            Direction = direction == "desc" ? SortDirection.Desc : SortDirection.Asc
        };

        var (items, totalCount) = await studioRepo.FindAsync(filter, findFilter, ct);
        var dtos = await MapListToDtos(items, ct);
        return Ok(new PaginatedResponse<StudioDto>(dtos, totalCount, page, perPage));
    }

    [HttpPost("find")]
    public async Task<ActionResult<PaginatedResponse<StudioDto>>> FindPost([FromBody] FilteredQueryRequest<StudioFilter> req, CancellationToken ct)
    {
        var findFilter = req.FindFilter ?? new FindFilter();
        var filter = req.ObjectFilter ?? new StudioFilter();
        var (items, totalCount) = await studioRepo.FindAsync(filter, findFilter, ct);
        var dtos = await MapListToDtos(items, ct);
        return Ok(new PaginatedResponse<StudioDto>(dtos, totalCount, findFilter.Page, findFilter.PerPage));
    }

    [HttpGet("{id:int}")]
    public async Task<ActionResult<StudioDto>> GetById(int id, CancellationToken ct)
    {
        var studio = await studioRepo.GetByIdWithRelationsAsync(id, ct);
        if (studio == null) return NotFound();
        return Ok(MapToDto(studio));
    }

    [HttpPost]
    public async Task<ActionResult<StudioDto>> Create([FromBody] StudioCreateDto dto, CancellationToken ct)
    {
        var studio = new Studio
        {
            Name = dto.Name, ParentId = dto.ParentId, Rating = dto.Rating,
            Favorite = dto.Favorite, Details = dto.Details,
            IgnoreAutoTag = dto.IgnoreAutoTag, Organized = dto.Organized
        };
        if (dto.Urls?.Count > 0) studio.Urls = dto.Urls.Select(u => new StudioUrl { Url = u }).ToList();
        if (dto.Aliases?.Count > 0) studio.Aliases = dto.Aliases.Select(a => new StudioAlias { Alias = a }).ToList();
        if (dto.TagIds?.Count > 0) studio.StudioTags = dto.TagIds.Select(id => new StudioTag { TagId = id }).ToList();

        studio = await studioRepo.AddAsync(studio, ct);
        var result = await studioRepo.GetByIdWithRelationsAsync(studio.Id, ct);
        return CreatedAtAction(nameof(GetById), new { id = studio.Id }, MapToDto(result!));
    }

    [HttpPut("{id:int}")]
    public async Task<ActionResult<StudioDto>> Update(int id, [FromBody] StudioUpdateDto dto, CancellationToken ct)
    {
        var studio = await studioRepo.GetByIdWithRelationsAsync(id, ct);
        if (studio == null) return NotFound();

        if (dto.Name != null) studio.Name = dto.Name;
        if (dto.ParentId.HasValue) studio.ParentId = dto.ParentId;
        if (dto.Rating.HasValue) studio.Rating = dto.Rating;
        if (dto.Favorite.HasValue) studio.Favorite = dto.Favorite.Value;
        if (dto.Details != null) studio.Details = dto.Details;
        if (dto.IgnoreAutoTag.HasValue) studio.IgnoreAutoTag = dto.IgnoreAutoTag.Value;
        if (dto.Organized.HasValue) studio.Organized = dto.Organized.Value;

        if (dto.Urls != null)
        {
            studio.Urls.Clear();
            studio.Urls = dto.Urls.Select(u => new StudioUrl { Url = u, StudioId = id }).ToList();
        }
        if (dto.Aliases != null)
        {
            studio.Aliases.Clear();
            studio.Aliases = dto.Aliases.Select(a => new StudioAlias { Alias = a, StudioId = id }).ToList();
        }
        if (dto.TagIds != null)
        {
            studio.StudioTags.Clear();
            studio.StudioTags = dto.TagIds.Select(tid => new StudioTag { TagId = tid, StudioId = id }).ToList();
        }

        await studioRepo.UpdateAsync(studio, ct);
        var updated = await studioRepo.GetByIdWithRelationsAsync(id, ct);
        return Ok(MapToDto(updated!));
    }

    [HttpDelete("{id:int}")]
    public async Task<IActionResult> Delete(int id, CancellationToken ct)
    {
        var s = await studioRepo.GetByIdAsync(id, ct);
        if (s == null) return NotFound();
        await studioRepo.DeleteAsync(id, ct);
        return NoContent();
    }

    private static StudioDto MapToDto(Studio s, int? sceneCount = null, int? imageCount = null, int? galleryCount = null, int? groupCount = null, int? performerCount = null, int? childStudioCount = null) => new(
        s.Id, s.Name, s.ParentId, s.Parent?.Name, s.Rating, s.Favorite, s.Details, s.IgnoreAutoTag, s.Organized,
        s.Urls.Select(u => u.Url).ToList(),
        s.Aliases.Select(a => a.Alias).ToList(),
        s.StudioTags.Where(st => st.Tag != null).Select(st => new TagDto(st.Tag!.Id, st.Tag.Name, st.Tag.Description, st.Tag.Favorite, st.Tag.IgnoreAutoTag, [])).ToList(),
        s.StashIds.Select(sid => new StudioStashIdDto(sid.Endpoint, sid.StashId)).ToList(),
        sceneCount ?? s.Scenes?.Count ?? 0,
        imageCount ?? s.Images?.Count ?? 0,
        galleryCount ?? s.Galleries?.Count ?? 0,
        groupCount ?? s.Groups?.Count ?? 0,
        performerCount ?? 0,
        childStudioCount ?? s.Children?.Count ?? 0,
        s.ImageBlobId != null ? $"/api/studios/{s.Id}/image" : null,
        s.CreatedAt.ToString("o"), s.UpdatedAt.ToString("o")
    );

    private async Task<List<StudioDto>> MapListToDtos(IReadOnlyList<Studio> items, CancellationToken ct)
    {
        if (items.Count == 0) return [];
        var ids = items.Select(s => s.Id).ToList();
        var sceneCounts = await db.Scenes.Where(sc => sc.StudioId.HasValue && ids.Contains(sc.StudioId.Value))
            .GroupBy(sc => sc.StudioId!.Value).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var imgCounts = await db.Images.Where(i => i.StudioId.HasValue && ids.Contains(i.StudioId.Value))
            .GroupBy(i => i.StudioId!.Value).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var galCounts = await db.Galleries.Where(g => g.StudioId.HasValue && ids.Contains(g.StudioId.Value))
            .GroupBy(g => g.StudioId!.Value).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var grpCounts = await db.Groups.Where(g => g.StudioId.HasValue && ids.Contains(g.StudioId.Value))
            .GroupBy(g => g.StudioId!.Value).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var perfCounts = await db.Set<ScenePerformer>()
            .Where(sp => sp.Scene!.StudioId.HasValue && ids.Contains(sp.Scene.StudioId!.Value))
            .GroupBy(sp => sp.Scene!.StudioId!.Value)
            .Select(g => new { Id = g.Key, Count = g.Select(sp => sp.PerformerId).Distinct().Count() })
            .ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var childCounts = await db.Studios.Where(s => s.ParentId.HasValue && ids.Contains(s.ParentId.Value))
            .GroupBy(s => s.ParentId!.Value).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        return items.Select(s => MapToDto(s,
            sceneCounts.GetValueOrDefault(s.Id, 0),
            imgCounts.GetValueOrDefault(s.Id, 0),
            galCounts.GetValueOrDefault(s.Id, 0),
            grpCounts.GetValueOrDefault(s.Id, 0),
            perfCounts.GetValueOrDefault(s.Id, 0),
            childCounts.GetValueOrDefault(s.Id, 0)
        )).ToList();
    }

    private static List<int>? ParseIntList(string? csv) => string.IsNullOrEmpty(csv) ? null : csv.Split(',').Select(int.Parse).ToList();

    // ===== Merge =====

    [HttpPost("merge")]
    public async Task<ActionResult<StudioDto>> MergeStudios([FromBody] StudioMergeDto dto, CancellationToken ct)
    {
        var target = await studioRepo.GetByIdWithRelationsAsync(dto.TargetId, ct);
        if (target == null) return NotFound("Target studio not found");

        var sources = await db.Studios
            .Include(s => s.Aliases)
            .Include(s => s.Urls)
            .Include(s => s.Children)
            .Include(s => s.Scenes)
            .Include(s => s.Galleries)
            .Include(s => s.Images)
            .Include(s => s.Groups)
            .Where(s => dto.SourceIds.Contains(s.Id))
            .ToListAsync(ct);

        foreach (var source in sources)
        {
            // Move scenes
            foreach (var scene in source.Scenes)
                scene.StudioId = target.Id;
            // Move galleries
            foreach (var gallery in source.Galleries)
                gallery.StudioId = target.Id;
            // Move images
            foreach (var image in source.Images)
                image.StudioId = target.Id;
            // Move groups
            foreach (var group in source.Groups)
                group.StudioId = target.Id;
            // Reparent child studios
            foreach (var child in source.Children)
                child.ParentId = target.Id;
            // Add source name as alias
            if (!target.Aliases.Any(a => a.Alias == source.Name))
                target.Aliases.Add(new StudioAlias { Alias = source.Name, StudioId = target.Id });
            // Delete source
            db.Studios.Remove(source);
        }

        await db.SaveChangesAsync(ct);
        var result = await studioRepo.GetByIdWithRelationsAsync(target.Id, ct);
        return Ok(MapToDto(result!));
    }

    // ===== Stash-Box =====

    [HttpGet("{id:int}/stash-box/search")]
    public async Task<ActionResult<IReadOnlyList<StashBoxStudioMatchDto>>> SearchStashBox(int id, [FromQuery] string? term, [FromQuery] string? endpoint, CancellationToken ct)
    {
        var studio = await studioRepo.GetByIdWithRelationsAsync(id, ct);
        if (studio == null) return NotFound();

        if (string.IsNullOrWhiteSpace(term))
        {
            var existingStashId = studio.StashIds?.FirstOrDefault(s => string.IsNullOrWhiteSpace(endpoint) || string.Equals(s.Endpoint, endpoint, StringComparison.OrdinalIgnoreCase));
            if (existingStashId != null)
            {
                // For studios, just re-search by name (no individual lookup like performers)
                term = studio.Name;
            }
            else
            {
                term = studio.Name;
            }
        }

        return Ok(await stashBoxService.SearchStudiosAsync(term, endpoint, ct));
    }

    [HttpPost("{id:int}/stash-box/import")]
    public async Task<ActionResult<StudioDto>> ImportFromStashBox(int id, [FromBody] StashBoxStudioImportRequestDto dto, CancellationToken ct)
    {
        var studio = await db.Studios
            .Include(s => s.StashIds)
            .Include(s => s.Aliases)
            .Include(s => s.Urls)
            .Include(s => s.StudioTags).ThenInclude(st => st.Tag)
            .Include(s => s.Parent)
            .FirstOrDefaultAsync(s => s.Id == id, ct);
        if (studio == null) return NotFound();

        var imported = await stashBoxService.MergeStudioAsync(studio, dto.Endpoint, dto.StudioId, ct);
        if (!imported) return NotFound();

        await db.SaveChangesAsync(ct);
        var updated = await studioRepo.GetByIdWithRelationsAsync(id, ct);
        return Ok(MapToDto(updated!));
    }
}
