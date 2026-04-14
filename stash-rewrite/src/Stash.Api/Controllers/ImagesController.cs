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
public class ImagesController(IImageRepository imageRepo, Data.StashContext db) : ControllerBase
{
    [HttpGet]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<PaginatedResponse<ImageDto>>> Find(
        [FromQuery] string? q, [FromQuery] int page = 1, [FromQuery] int perPage = 25,
        [FromQuery] string? sort = null, [FromQuery] string? direction = null,
        [FromQuery] string? title = null, [FromQuery] int? rating = null,
        [FromQuery] bool? organized = null, [FromQuery] int? studioId = null,
        [FromQuery] string? tagIds = null, [FromQuery] string? performerIds = null,
        [FromQuery] int? galleryId = null,
        CancellationToken ct = default)
    {
        var filter = new ImageFilter
        {
            Title = title, Rating = rating, Organized = organized, StudioId = studioId,
            TagIds = ParseIntList(tagIds), PerformerIds = ParseIntList(performerIds),
            GalleryId = galleryId
        };
        var findFilter = new FindFilter
        {
            Q = q, Page = page, PerPage = perPage, Sort = sort,
            Direction = direction == "desc" ? SortDirection.Desc : SortDirection.Asc
        };

        var (items, totalCount) = await imageRepo.FindAsync(filter, findFilter, ct);
        var dtos = await MapListToDtos(items, ct);
        return Ok(new PaginatedResponse<ImageDto>(dtos, totalCount, page, perPage));
    }

    [HttpPost("find")]
    public async Task<ActionResult<PaginatedResponse<ImageDto>>> FindPost([FromBody] FilteredQueryRequest<ImageFilter> req, CancellationToken ct)
    {
        var findFilter = req.FindFilter ?? new FindFilter();
        var filter = req.ObjectFilter ?? new ImageFilter();
        var (items, totalCount) = await imageRepo.FindAsync(filter, findFilter, ct);
        var dtos = await MapListToDtos(items, ct);
        return Ok(new PaginatedResponse<ImageDto>(dtos, totalCount, findFilter.Page, findFilter.PerPage));
    }

    [HttpGet("{id:int}")]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<ImageDto>> GetById(int id, CancellationToken ct)
    {
        var image = await imageRepo.GetByIdWithRelationsAsync(id, ct);
        if (image == null) return NotFound();
        return Ok(MapToDto(image));
    }

    [HttpPost]
    public async Task<ActionResult<ImageDto>> Create([FromBody] ImageCreateDto dto, CancellationToken ct)
    {
        var image = new Image
        {
            Title = dto.Title,
            Code = dto.Code,
            Details = dto.Details,
            Photographer = dto.Photographer,
            Rating = dto.Rating,
            Organized = dto.Organized,
            StudioId = dto.StudioId,
            Date = ParseDate(dto.Date)
        };

        if (dto.Urls?.Count > 0)
            image.Urls = dto.Urls.Select(u => new ImageUrl { Url = u }).ToList();
        if (dto.TagIds?.Count > 0)
            image.ImageTags = dto.TagIds.Select(tagId => new ImageTag { TagId = tagId }).ToList();
        if (dto.PerformerIds?.Count > 0)
            image.ImagePerformers = dto.PerformerIds.Select(performerId => new ImagePerformer { PerformerId = performerId }).ToList();

        image = await imageRepo.AddAsync(image, ct);
        var result = await imageRepo.GetByIdWithRelationsAsync(image.Id, ct);
        return CreatedAtAction(nameof(GetById), new { id = image.Id }, MapToDto(result!));
    }

    [HttpPut("{id:int}")]
    public async Task<ActionResult<ImageDto>> Update(int id, [FromBody] ImageUpdateDto dto, CancellationToken ct)
    {
        var image = await imageRepo.GetByIdWithRelationsAsync(id, ct);
        if (image == null) return NotFound();

        if (dto.Title != null) image.Title = dto.Title;
        if (dto.Code != null) image.Code = dto.Code;
        if (dto.Details != null) image.Details = dto.Details;
        if (dto.Photographer != null) image.Photographer = dto.Photographer;
        if (dto.Rating.HasValue) image.Rating = dto.Rating;
        if (dto.Organized.HasValue) image.Organized = dto.Organized.Value;
        if (dto.StudioId.HasValue) image.StudioId = dto.StudioId;
        if (dto.Date != null) image.Date = ParseDate(dto.Date);

        if (dto.Urls != null)
        {
            image.Urls.Clear();
            image.Urls = dto.Urls.Select(u => new ImageUrl { Url = u, ImageId = id }).ToList();
        }
        if (dto.TagIds != null)
        {
            image.ImageTags.Clear();
            image.ImageTags = dto.TagIds.Select(tid => new ImageTag { TagId = tid, ImageId = id }).ToList();
        }
        if (dto.PerformerIds != null)
        {
            image.ImagePerformers.Clear();
            image.ImagePerformers = dto.PerformerIds.Select(pid => new ImagePerformer { PerformerId = pid, ImageId = id }).ToList();
        }
        if (dto.GalleryIds != null)
        {
            image.ImageGalleries.Clear();
            image.ImageGalleries = dto.GalleryIds.Select(gid => new ImageGallery { GalleryId = gid, ImageId = id }).ToList();
        }
        if (dto.CustomFields != null) image.CustomFields = dto.CustomFields;

        await imageRepo.UpdateAsync(image, ct);
        var updated = await imageRepo.GetByIdWithRelationsAsync(id, ct);
        return Ok(MapToDto(updated!));
    }

    [HttpDelete("{id:int}")]
    public async Task<IActionResult> Delete(int id, CancellationToken ct)
    {
        var img = await imageRepo.GetByIdAsync(id, ct);
        if (img == null) return NotFound();
        await imageRepo.DeleteAsync(id, ct);
        return NoContent();
    }

    private static ImageDto MapToDto(Image i, int? galleryCount = null) => new(
        i.Id, i.Title, i.Code, i.Details, i.Photographer,
        i.Rating, i.Organized, i.OCounter, i.StudioId, i.Studio?.Name,
        i.Date?.ToString("yyyy-MM-dd"),
        i.Urls.Select(u => u.Url).ToList(),
        i.ImageTags.Where(it => it.Tag != null).Select(it => new TagDto(it.Tag!.Id, it.Tag.Name, it.Tag.Description, it.Tag.Favorite, it.Tag.IgnoreAutoTag, [])).ToList(),
        i.ImagePerformers.Where(ip => ip.Performer != null).Select(ip => new PerformerSummaryDto(ip.Performer!.Id, ip.Performer.Name, ip.Performer.Disambiguation, ip.Performer.Gender?.ToString(), ip.Performer.Favorite, ip.Performer.ImageBlobId != null ? $"/api/performers/{ip.Performer.Id}/image" : null)).ToList(),
        galleryCount ?? i.ImageGalleries?.Count ?? 0,
        i.ImageGalleries?.Select(ig => ig.GalleryId).ToList() ?? [],
        i.CustomFields,
        i.CreatedAt.ToString("o"), i.UpdatedAt.ToString("o")
    );

    private async Task<List<ImageDto>> MapListToDtos(IReadOnlyList<Image> items, CancellationToken ct)
    {
        if (items.Count == 0) return [];
        var ids = items.Select(i => i.Id).ToList();
        var galCounts = await db.Set<ImageGallery>().Where(ig => ids.Contains(ig.ImageId))
            .GroupBy(ig => ig.ImageId).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        return items.Select(i => MapToDto(i, galCounts.GetValueOrDefault(i.Id, 0))).ToList();
    }

    // ===== Activity Tracking =====

    [HttpPost("{id:int}/o")]
    public async Task<IActionResult> IncrementO(int id, CancellationToken ct)
    {
        var image = await imageRepo.GetByIdAsync(id, ct);
        if (image == null) return NotFound();
        image.OCounter++;
        await imageRepo.UpdateAsync(image, ct);
        return NoContent();
    }

    [HttpDelete("{id:int}/o")]
    public async Task<IActionResult> DecrementO(int id, CancellationToken ct)
    {
        var image = await imageRepo.GetByIdAsync(id, ct);
        if (image == null) return NotFound();
        image.OCounter = Math.Max(0, image.OCounter - 1);
        await imageRepo.UpdateAsync(image, ct);
        return NoContent();
    }

    [HttpPost("{id:int}/o/reset")]
    public async Task<IActionResult> ResetO(int id, CancellationToken ct)
    {
        var image = await imageRepo.GetByIdAsync(id, ct);
        if (image == null) return NotFound();
        image.OCounter = 0;
        await imageRepo.UpdateAsync(image, ct);
        return NoContent();
    }

    // ===== Bulk Operations =====

    [HttpPost("bulk")]
    public async Task<IActionResult> BulkUpdate([FromBody] BulkImageUpdateDto dto, CancellationToken ct)
    {
        var images = await db.Images
            .Include(i => i.ImageTags)
            .Include(i => i.ImagePerformers)
            .Where(i => dto.Ids.Contains(i.Id))
            .ToListAsync(ct);

        foreach (var image in images)
        {
            if (dto.Rating.HasValue) image.Rating = dto.Rating;
            if (dto.Organized.HasValue) image.Organized = dto.Organized.Value;
            if (dto.StudioId.HasValue) image.StudioId = dto.StudioId;

            if (dto.TagIds != null && dto.TagMode == BulkUpdateMode.Set)
            {
                image.ImageTags.Clear();
                image.ImageTags = dto.TagIds.Select(tid => new ImageTag { TagId = tid, ImageId = image.Id }).ToList();
            }
            else if (dto.TagIds != null && dto.TagMode == BulkUpdateMode.Add)
            {
                var existing = image.ImageTags.Select(it => it.TagId).ToHashSet();
                foreach (var tid in dto.TagIds.Where(t => !existing.Contains(t)))
                    image.ImageTags.Add(new ImageTag { TagId = tid, ImageId = image.Id });
            }
            else if (dto.TagIds != null && dto.TagMode == BulkUpdateMode.Remove)
            {
                image.ImageTags = image.ImageTags.Where(it => !dto.TagIds.Contains(it.TagId)).ToList();
            }

            if (dto.PerformerIds != null && dto.PerformerMode == BulkUpdateMode.Set)
            {
                image.ImagePerformers.Clear();
                image.ImagePerformers = dto.PerformerIds.Select(pid => new ImagePerformer { PerformerId = pid, ImageId = image.Id }).ToList();
            }
            else if (dto.PerformerIds != null && dto.PerformerMode == BulkUpdateMode.Add)
            {
                var existing = image.ImagePerformers.Select(ip => ip.PerformerId).ToHashSet();
                foreach (var pid in dto.PerformerIds.Where(p => !existing.Contains(p)))
                    image.ImagePerformers.Add(new ImagePerformer { PerformerId = pid, ImageId = image.Id });
            }
            else if (dto.PerformerIds != null && dto.PerformerMode == BulkUpdateMode.Remove)
            {
                image.ImagePerformers = image.ImagePerformers.Where(ip => !dto.PerformerIds.Contains(ip.PerformerId)).ToList();
            }
        }

        await db.SaveChangesAsync(ct);
        return Ok(new { updated = images.Count });
    }

    private static DateOnly? ParseDate(string? date) => DateOnly.TryParse(date, out var d) ? d : null;
    private static List<int>? ParseIntList(string? csv) => string.IsNullOrEmpty(csv) ? null : csv.Split(',').Select(int.Parse).ToList();
}
