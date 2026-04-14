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
public class GroupsController(IGroupRepository groupRepo, Data.StashContext db) : ControllerBase
{
    [HttpGet]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<PaginatedResponse<GroupDto>>> Find(
        [FromQuery] string? q, [FromQuery] int page = 1, [FromQuery] int perPage = 25,
        [FromQuery] string? sort = null, [FromQuery] string? direction = null,
        [FromQuery] string? name = null, [FromQuery] int? rating = null,
        [FromQuery] int? studioId = null,
        CancellationToken ct = default)
    {
        var filter = new GroupFilter { Name = name, Rating = rating, StudioId = studioId };
        var findFilter = new FindFilter
        {
            Q = q, Page = page, PerPage = perPage, Sort = sort,
            Direction = direction == "desc" ? SortDirection.Desc : SortDirection.Asc
        };

        var (items, totalCount) = await groupRepo.FindAsync(filter, findFilter, ct);
        var dtos = items.Select(MapToDto).ToList();
        return Ok(new PaginatedResponse<GroupDto>(dtos, totalCount, page, perPage));
    }

    [HttpPost("find")]
    public async Task<ActionResult<PaginatedResponse<GroupDto>>> FindPost([FromBody] FilteredQueryRequest<GroupFilter> req, CancellationToken ct)
    {
        var findFilter = req.FindFilter ?? new FindFilter();
        var filter = req.ObjectFilter ?? new GroupFilter();
        var (items, totalCount) = await groupRepo.FindAsync(filter, findFilter, ct);
        var dtos = items.Select(MapToDto).ToList();
        return Ok(new PaginatedResponse<GroupDto>(dtos, totalCount, findFilter.Page, findFilter.PerPage));
    }

    [HttpGet("{id:int}")]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<GroupDto>> GetById(int id, CancellationToken ct)
    {
        var group = await groupRepo.GetByIdWithRelationsAsync(id, ct);
        if (group == null) return NotFound();
        return Ok(MapToDto(group));
    }

    [HttpPost]
    public async Task<ActionResult<GroupDto>> Create([FromBody] GroupCreateDto dto, CancellationToken ct)
    {
        var group = new Group
        {
            Name = dto.Name, Aliases = dto.Aliases, Duration = dto.Duration,
            Date = ParseDate(dto.Date), Rating = dto.Rating, StudioId = dto.StudioId,
            Director = dto.Director, Synopsis = dto.Synopsis
        };
        if (dto.Urls?.Count > 0) group.Urls = dto.Urls.Select(u => new GroupUrl { Url = u }).ToList();
        if (dto.TagIds?.Count > 0) group.GroupTags = dto.TagIds.Select(id => new GroupTag { TagId = id }).ToList();

        group = await groupRepo.AddAsync(group, ct);
        var result = await groupRepo.GetByIdWithRelationsAsync(group.Id, ct);
        return CreatedAtAction(nameof(GetById), new { id = group.Id }, MapToDto(result!));
    }

    [HttpPut("{id:int}")]
    public async Task<ActionResult<GroupDto>> Update(int id, [FromBody] GroupUpdateDto dto, CancellationToken ct)
    {
        var group = await groupRepo.GetByIdWithRelationsAsync(id, ct);
        if (group == null) return NotFound();

        if (dto.Name != null) group.Name = dto.Name;
        if (dto.Aliases != null) group.Aliases = dto.Aliases;
        if (dto.Duration.HasValue) group.Duration = dto.Duration;
        if (dto.Date != null) group.Date = ParseDate(dto.Date);
        if (dto.Rating.HasValue) group.Rating = dto.Rating;
        if (dto.StudioId.HasValue) group.StudioId = dto.StudioId;
        if (dto.Director != null) group.Director = dto.Director;
        if (dto.Synopsis != null) group.Synopsis = dto.Synopsis;

        if (dto.Urls != null)
        {
            group.Urls.Clear();
            group.Urls = dto.Urls.Select(u => new GroupUrl { Url = u, GroupId = id }).ToList();
        }
        if (dto.TagIds != null)
        {
            group.GroupTags.Clear();
            group.GroupTags = dto.TagIds.Select(tid => new GroupTag { TagId = tid, GroupId = id }).ToList();
        }
        if (dto.CustomFields != null) group.CustomFields = dto.CustomFields;

        await groupRepo.UpdateAsync(group, ct);
        var updated = await groupRepo.GetByIdWithRelationsAsync(id, ct);
        return Ok(MapToDto(updated!));
    }

    [HttpDelete("{id:int}")]
    public async Task<IActionResult> Delete(int id, CancellationToken ct)
    {
        var g = await groupRepo.GetByIdAsync(id, ct);
        if (g == null) return NotFound();
        await groupRepo.DeleteAsync(id, ct);
        return NoContent();
    }

    [HttpGet("{id:int}/subgroups")]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<List<GroupDto>>> GetSubGroups(int id, CancellationToken ct)
    {
        var relations = await db.Set<GroupRelation>()
            .Where(r => r.ContainingGroupId == id)
            .OrderBy(r => r.OrderIndex)
            .Include(r => r.SubGroup!).ThenInclude(g => g.Urls)
            .Include(r => r.SubGroup!).ThenInclude(g => g.GroupTags).ThenInclude(gt => gt.Tag)
            .Include(r => r.SubGroup!).ThenInclude(g => g.SceneGroups)
            .ToListAsync(ct);
        return Ok(relations.Where(r => r.SubGroup != null).Select(r => MapToDto(r.SubGroup!)).ToList());
    }

    [HttpGet("{id:int}/containinggroups")]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<List<GroupDto>>> GetContainingGroups(int id, CancellationToken ct)
    {
        var relations = await db.Set<GroupRelation>()
            .Where(r => r.SubGroupId == id)
            .OrderBy(r => r.OrderIndex)
            .Include(r => r.ContainingGroup!).ThenInclude(g => g.Urls)
            .Include(r => r.ContainingGroup!).ThenInclude(g => g.GroupTags).ThenInclude(gt => gt.Tag)
            .Include(r => r.ContainingGroup!).ThenInclude(g => g.SceneGroups)
            .ToListAsync(ct);
        return Ok(relations.Where(r => r.ContainingGroup != null).Select(r => MapToDto(r.ContainingGroup!)).ToList());
    }

    [HttpPost("{id:int}/subgroups")]
    public async Task<IActionResult> AddSubGroup(int id, [FromBody] AddSubGroupDto dto, CancellationToken ct)
    {
        var group = await groupRepo.GetByIdAsync(id, ct);
        if (group == null) return NotFound();

        var existing = await db.Set<GroupRelation>()
            .Where(r => r.ContainingGroupId == id)
            .ToListAsync(ct);

        if (existing.Any(r => r.SubGroupId == dto.SubGroupId))
            return Conflict("Sub-group already exists");

        var maxOrder = existing.Count > 0 ? existing.Max(r => r.OrderIndex) + 1 : 0;
        db.Set<GroupRelation>().Add(new GroupRelation
        {
            ContainingGroupId = id,
            SubGroupId = dto.SubGroupId,
            OrderIndex = dto.OrderIndex ?? maxOrder,
            Description = dto.Description,
        });
        await db.SaveChangesAsync(ct);
        return Ok();
    }

    [HttpDelete("{id:int}/subgroups/{subGroupId:int}")]
    public async Task<IActionResult> RemoveSubGroup(int id, int subGroupId, CancellationToken ct)
    {
        var relation = await db.Set<GroupRelation>()
            .FirstOrDefaultAsync(r => r.ContainingGroupId == id && r.SubGroupId == subGroupId, ct);
        if (relation == null) return NotFound();
        db.Set<GroupRelation>().Remove(relation);
        await db.SaveChangesAsync(ct);
        return NoContent();
    }

    [HttpPut("{id:int}/subgroups/reorder")]
    public async Task<IActionResult> ReorderSubGroups(int id, [FromBody] ReorderSubGroupsDto dto, CancellationToken ct)
    {
        var relations = await db.Set<GroupRelation>()
            .Where(r => r.ContainingGroupId == id)
            .ToListAsync(ct);

        for (var i = 0; i < dto.SubGroupIds.Count; i++)
        {
            var rel = relations.FirstOrDefault(r => r.SubGroupId == dto.SubGroupIds[i]);
            if (rel != null) rel.OrderIndex = i;
        }
        await db.SaveChangesAsync(ct);
        return Ok();
    }

    private static GroupDto MapToDto(Group g) => new(
        g.Id, g.Name, g.Aliases, g.Duration, g.Date?.ToString("yyyy-MM-dd"),
        g.Rating, g.StudioId, g.Studio?.Name, g.Director, g.Synopsis,
        g.Urls.Select(u => u.Url).ToList(),
        g.GroupTags.Where(gt => gt.Tag != null).Select(gt => new TagDto(gt.Tag!.Id, gt.Tag.Name, gt.Tag.Description, gt.Tag.Favorite, gt.Tag.IgnoreAutoTag, [])).ToList(),
        g.SceneGroups?.Count ?? 0,
        g.SubGroupRelations?.Count ?? 0,
        g.ContainingGroupRelations?.Count ?? 0,
        g.CustomFields,
        g.CreatedAt.ToString("o"), g.UpdatedAt.ToString("o")
    );

    private static DateOnly? ParseDate(string? date) => DateOnly.TryParse(date, out var d) ? d : null;
}
