using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.OutputCaching;
using Microsoft.EntityFrameworkCore;
using Stash.Core.DTOs;
using Stash.Core.Entities;
using Stash.Core.Enums;
using Stash.Core.Interfaces;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class TagsController(ITagRepository tagRepo, Data.StashContext db) : ControllerBase
{
    [HttpGet]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<PaginatedResponse<TagListDto>>> Find(
        [FromQuery] string? q, [FromQuery] int page = 1, [FromQuery] int perPage = 25,
        [FromQuery] string? sort = null, [FromQuery] string? direction = null,
        [FromQuery] string? name = null, [FromQuery] bool? favorite = null,
        CancellationToken ct = default)
    {
        var filter = new TagFilter { Name = name, Favorite = favorite };
        var findFilter = new FindFilter
        {
            Q = q, Page = page, PerPage = perPage, Sort = sort,
            Direction = direction == "desc" ? SortDirection.Desc : SortDirection.Asc
        };

        var (items, totalCount) = await tagRepo.FindAsync(filter, findFilter, ct);
        var dtos = await MapTagListDtos(items, ct);
        return Ok(new PaginatedResponse<TagListDto>(dtos, totalCount, page, perPage));
    }

    [HttpPost("find")]
    public async Task<ActionResult<PaginatedResponse<TagListDto>>> FindPost([FromBody] FilteredQueryRequest<TagFilter> req, CancellationToken ct)
    {
        var findFilter = req.FindFilter ?? new FindFilter();
        var filter = req.ObjectFilter ?? new TagFilter();
        var (items, totalCount) = await tagRepo.FindAsync(filter, findFilter, ct);
        var dtos = await MapTagListDtos(items, ct);
        return Ok(new PaginatedResponse<TagListDto>(dtos, totalCount, findFilter.Page, findFilter.PerPage));
    }

    [HttpGet("{id:int}")]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<TagDetailDto>> GetById(int id, CancellationToken ct)
    {
        var tag = await tagRepo.GetByIdWithRelationsAsync(id, ct);
        if (tag == null) return NotFound();

        var studioCount = await db.Set<StudioTag>().CountAsync(st => st.TagId == id, ct);
        var groupCount = await db.Set<GroupTag>().CountAsync(gt => gt.TagId == id, ct);
        var markerCount = await db.SceneMarkers.CountAsync(m => m.PrimaryTagId == id, ct)
            + await db.Set<SceneMarkerTag>().CountAsync(mt => mt.TagId == id, ct);

        return Ok(MapToDetailDto(tag, studioCount, groupCount, markerCount));
    }

    [HttpPost]
    public async Task<ActionResult<TagDetailDto>> Create([FromBody] TagCreateDto dto, CancellationToken ct)
    {
        var existing = await tagRepo.GetByNameAsync(dto.Name, ct);
        if (existing != null) return Conflict(new { message = $"Tag '{dto.Name}' already exists" });

        var tag = new Tag
        {
            Name = dto.Name, SortName = dto.SortName, Description = dto.Description,
            Favorite = dto.Favorite, IgnoreAutoTag = dto.IgnoreAutoTag
        };
        if (dto.Aliases?.Count > 0) tag.Aliases = dto.Aliases.Select(a => new TagAlias { Alias = a }).ToList();
        if (dto.ParentIds?.Count > 0) tag.ParentRelations = dto.ParentIds.Select(pid => new TagParent { ParentId = pid }).ToList();
        if (dto.ChildIds?.Count > 0) tag.ChildRelations = dto.ChildIds.Select(cid => new TagParent { ChildId = cid }).ToList();

        tag = await tagRepo.AddAsync(tag, ct);
        var result = await tagRepo.GetByIdWithRelationsAsync(tag.Id, ct);
        return CreatedAtAction(nameof(GetById), new { id = tag.Id }, MapToDetailDto(result!));
    }

    [HttpPut("{id:int}")]
    public async Task<ActionResult<TagDetailDto>> Update(int id, [FromBody] TagUpdateDto dto, CancellationToken ct)
    {
        var tag = await tagRepo.GetByIdWithRelationsAsync(id, ct);
        if (tag == null) return NotFound();

        if (dto.Name != null) tag.Name = dto.Name;
        if (dto.SortName != null) tag.SortName = dto.SortName;
        if (dto.Description != null) tag.Description = dto.Description;
        if (dto.Favorite.HasValue) tag.Favorite = dto.Favorite.Value;
        if (dto.IgnoreAutoTag.HasValue) tag.IgnoreAutoTag = dto.IgnoreAutoTag.Value;

        if (dto.Aliases != null)
        {
            tag.Aliases.Clear();
            tag.Aliases = dto.Aliases.Select(a => new TagAlias { Alias = a, TagId = id }).ToList();
        }
        if (dto.ParentIds != null)
        {
            tag.ParentRelations.Clear();
            tag.ParentRelations = dto.ParentIds.Select(pid => new TagParent { ParentId = pid, ChildId = id }).ToList();
        }
        if (dto.ChildIds != null)
        {
            tag.ChildRelations.Clear();
            tag.ChildRelations = dto.ChildIds.Select(cid => new TagParent { ParentId = id, ChildId = cid }).ToList();
        }
        if (dto.CustomFields != null) tag.CustomFields = dto.CustomFields;

        await tagRepo.UpdateAsync(tag, ct);
        var updated = await tagRepo.GetByIdWithRelationsAsync(id, ct);
        return Ok(MapToDetailDto(updated!));
    }

    [HttpDelete("{id:int}")]
    public async Task<IActionResult> Delete(int id, CancellationToken ct)
    {
        var tag = await tagRepo.GetByIdAsync(id, ct);
        if (tag == null) return NotFound();
        await tagRepo.DeleteAsync(id, ct);
        return NoContent();
    }

    private static TagDetailDto MapToDetailDto(Tag t, int studioCount = 0, int groupCount = 0, int markerCount = 0) => new(
        t.Id, t.Name, t.SortName, t.Description, t.Favorite, t.IgnoreAutoTag,
        t.Aliases.Select(a => a.Alias).ToList(),
        t.ParentRelations.Where(pr => pr.Parent != null).Select(pr => new TagDto(pr.Parent!.Id, pr.Parent.Name, pr.Parent.Description, pr.Parent.Favorite, pr.Parent.IgnoreAutoTag, [])).ToList(),
        t.ChildRelations.Where(cr => cr.Child != null).Select(cr => new TagDto(cr.Child!.Id, cr.Child.Name, cr.Child.Description, cr.Child.Favorite, cr.Child.IgnoreAutoTag, [])).ToList(),
        t.SceneTags?.Count ?? 0, t.PerformerTags?.Count ?? 0, t.ImageTags?.Count ?? 0, t.GalleryTags?.Count ?? 0,
        studioCount, groupCount, markerCount,
        t.CustomFields,
        t.CreatedAt.ToString("o"), t.UpdatedAt.ToString("o")
    );

    private async Task<List<TagListDto>> MapTagListDtos(IReadOnlyList<Tag> items, CancellationToken ct)
    {
        if (items.Count == 0) return [];
        var ids = items.Select(t => t.Id).ToList();
        var sceneCounts = await db.Set<SceneTag>().Where(st => ids.Contains(st.TagId))
            .GroupBy(st => st.TagId).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var markerCounts = await db.SceneMarkers.Where(m => ids.Contains(m.PrimaryTagId))
            .GroupBy(m => m.PrimaryTagId).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var imageCounts = await db.Set<ImageTag>().Where(it => ids.Contains(it.TagId))
            .GroupBy(it => it.TagId).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var galleryCounts = await db.Set<GalleryTag>().Where(gt => ids.Contains(gt.TagId))
            .GroupBy(gt => gt.TagId).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var groupCounts = await db.Set<GroupTag>().Where(gt => ids.Contains(gt.TagId))
            .GroupBy(gt => gt.TagId).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var performerCounts = await db.Set<PerformerTag>().Where(pt => ids.Contains(pt.TagId))
            .GroupBy(pt => pt.TagId).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var studioCounts = await db.Set<StudioTag>().Where(st => ids.Contains(st.TagId))
            .GroupBy(st => st.TagId).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        return items.Select(t => new TagListDto(t.Id, t.Name, t.Description, t.Favorite, t.IgnoreAutoTag,
            t.Aliases.Select(a => a.Alias).ToList(),
            sceneCounts.GetValueOrDefault(t.Id, 0),
            markerCounts.GetValueOrDefault(t.Id, 0),
            imageCounts.GetValueOrDefault(t.Id, 0),
            galleryCounts.GetValueOrDefault(t.Id, 0),
            groupCounts.GetValueOrDefault(t.Id, 0),
            performerCounts.GetValueOrDefault(t.Id, 0),
            studioCounts.GetValueOrDefault(t.Id, 0),
            t.ImageBlobId != null ? $"/api/tags/{t.Id}/image" : null
        )).ToList();
    }

    // ===== Bulk Operations =====

    [HttpPost("bulk")]
    public async Task<IActionResult> BulkUpdate([FromBody] BulkTagUpdateDto dto, CancellationToken ct)
    {
        var tags = await db.Tags.Where(t => dto.Ids.Contains(t.Id)).ToListAsync(ct);
        foreach (var tag in tags)
        {
            if (dto.Favorite.HasValue) tag.Favorite = dto.Favorite.Value;
            if (dto.IgnoreAutoTag.HasValue) tag.IgnoreAutoTag = dto.IgnoreAutoTag.Value;
        }
        await db.SaveChangesAsync(ct);
        return Ok(new { updated = tags.Count });
    }

    // ===== Merge =====

    [HttpPost("merge")]
    public async Task<ActionResult<TagDetailDto>> MergeTags([FromBody] TagMergeDto dto, CancellationToken ct)
    {
        var target = await tagRepo.GetByIdWithRelationsAsync(dto.TargetId, ct);
        if (target == null) return NotFound("Target tag not found");

        var sources = await db.Tags
            .Include(t => t.Aliases)
            .Include(t => t.SceneTags)
            .Include(t => t.PerformerTags)
            .Include(t => t.ImageTags)
            .Include(t => t.GalleryTags)
            .Where(t => dto.SourceIds.Contains(t.Id))
            .ToListAsync(ct);

        foreach (var source in sources)
        {
            // Move scene associations
            foreach (var st in source.SceneTags)
                if (!target.SceneTags.Any(t => t.SceneId == st.SceneId))
                    db.Set<SceneTag>().Add(new SceneTag { SceneId = st.SceneId, TagId = target.Id });
            // Move performer associations
            foreach (var pt in source.PerformerTags)
                if (!target.PerformerTags.Any(t => t.PerformerId == pt.PerformerId))
                    db.Set<PerformerTag>().Add(new PerformerTag { PerformerId = pt.PerformerId, TagId = target.Id });
            // Move image associations
            foreach (var it in source.ImageTags)
                if (!target.ImageTags.Any(t => t.ImageId == it.ImageId))
                    db.Set<ImageTag>().Add(new ImageTag { ImageId = it.ImageId, TagId = target.Id });
            // Add source name as alias
            if (!target.Aliases.Any(a => a.Alias == source.Name))
                target.Aliases.Add(new TagAlias { Alias = source.Name, TagId = target.Id });
            // Delete source
            db.Tags.Remove(source);
        }

        await db.SaveChangesAsync(ct);
        var result = await tagRepo.GetByIdWithRelationsAsync(target.Id, ct);
        return Ok(MapToDetailDto(result!));
    }

    // ===== Marker Wall =====

    [HttpGet("marker-strings")]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<List<string>>> GetMarkerStrings([FromQuery] string? q, [FromQuery] string? sort, CancellationToken ct)
    {
        var query = db.SceneMarkers.AsNoTracking().Select(m => m.Title).Distinct();
        if (!string.IsNullOrEmpty(q))
            query = query.Where(t => EF.Functions.ILike(t, $"%{q}%"));

        var result = sort == "count"
            ? await db.SceneMarkers.AsNoTracking().GroupBy(m => m.Title).OrderByDescending(g => g.Count()).Select(g => g.Key).Take(100).ToListAsync(ct)
            : await query.OrderBy(t => t).Take(100).ToListAsync(ct);

        return Ok(result);
    }
}
