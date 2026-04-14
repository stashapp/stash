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
public class GalleriesController(IGalleryRepository galleryRepo, Data.StashContext db) : ControllerBase
{
    [HttpGet]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<PaginatedResponse<GalleryDto>>> Find(
        [FromQuery] string? q, [FromQuery] int page = 1, [FromQuery] int perPage = 25,
        [FromQuery] string? sort = null, [FromQuery] string? direction = null,
        [FromQuery] string? title = null, [FromQuery] int? rating = null,
        [FromQuery] bool? organized = null, [FromQuery] int? studioId = null,
        [FromQuery] string? tagIds = null, [FromQuery] string? performerIds = null,
        CancellationToken ct = default)
    {
        var filter = new GalleryFilter
        {
            Title = title, Rating = rating, Organized = organized, StudioId = studioId,
            TagIds = ParseIntList(tagIds), PerformerIds = ParseIntList(performerIds)
        };
        var findFilter = new FindFilter
        {
            Q = q, Page = page, PerPage = perPage, Sort = sort,
            Direction = direction == "desc" ? SortDirection.Desc : SortDirection.Asc
        };

        var (items, totalCount) = await galleryRepo.FindAsync(filter, findFilter, ct);
        var dtos = await MapListToDtos(items, ct);
        return Ok(new PaginatedResponse<GalleryDto>(dtos, totalCount, page, perPage));
    }

    [HttpPost("find")]
    public async Task<ActionResult<PaginatedResponse<GalleryDto>>> FindPost([FromBody] FilteredQueryRequest<GalleryFilter> req, CancellationToken ct)
    {
        var findFilter = req.FindFilter ?? new FindFilter();
        var filter = req.ObjectFilter ?? new GalleryFilter();
        var (items, totalCount) = await galleryRepo.FindAsync(filter, findFilter, ct);
        var dtos = await MapListToDtos(items, ct);
        return Ok(new PaginatedResponse<GalleryDto>(dtos, totalCount, findFilter.Page, findFilter.PerPage));
    }

    [HttpGet("{id:int}")]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<GalleryDto>> GetById(int id, CancellationToken ct)
    {
        var gallery = await galleryRepo.GetByIdWithRelationsAsync(id, ct);
        if (gallery == null) return NotFound();
        return Ok(MapToDto(gallery));
    }

    [HttpPost]
    public async Task<ActionResult<GalleryDto>> Create([FromBody] GalleryCreateDto dto, CancellationToken ct)
    {
        var gallery = new Gallery
        {
            Title = dto.Title, Code = dto.Code, Date = ParseDate(dto.Date),
            Details = dto.Details, Photographer = dto.Photographer,
            Rating = dto.Rating, Organized = dto.Organized, StudioId = dto.StudioId
        };
        if (dto.Urls?.Count > 0) gallery.Urls = dto.Urls.Select(u => new GalleryUrl { Url = u }).ToList();
        if (dto.TagIds?.Count > 0) gallery.GalleryTags = dto.TagIds.Select(id => new GalleryTag { TagId = id }).ToList();
        if (dto.PerformerIds?.Count > 0) gallery.GalleryPerformers = dto.PerformerIds.Select(id => new GalleryPerformer { PerformerId = id }).ToList();

        gallery = await galleryRepo.AddAsync(gallery, ct);
        var result = await galleryRepo.GetByIdWithRelationsAsync(gallery.Id, ct);
        return CreatedAtAction(nameof(GetById), new { id = gallery.Id }, MapToDto(result!));
    }

    [HttpPut("{id:int}")]
    public async Task<ActionResult<GalleryDto>> Update(int id, [FromBody] GalleryUpdateDto dto, CancellationToken ct)
    {
        var gallery = await galleryRepo.GetByIdWithRelationsAsync(id, ct);
        if (gallery == null) return NotFound();

        if (dto.Title != null) gallery.Title = dto.Title;
        if (dto.Code != null) gallery.Code = dto.Code;
        if (dto.Date != null) gallery.Date = ParseDate(dto.Date);
        if (dto.Details != null) gallery.Details = dto.Details;
        if (dto.Photographer != null) gallery.Photographer = dto.Photographer;
        if (dto.Rating.HasValue) gallery.Rating = dto.Rating;
        if (dto.Organized.HasValue) gallery.Organized = dto.Organized.Value;
        if (dto.StudioId.HasValue) gallery.StudioId = dto.StudioId;

        if (dto.Urls != null)
        {
            gallery.Urls.Clear();
            gallery.Urls = dto.Urls.Select(u => new GalleryUrl { Url = u, GalleryId = id }).ToList();
        }
        if (dto.TagIds != null)
        {
            gallery.GalleryTags.Clear();
            gallery.GalleryTags = dto.TagIds.Select(tid => new GalleryTag { TagId = tid, GalleryId = id }).ToList();
        }
        if (dto.PerformerIds != null)
        {
            gallery.GalleryPerformers.Clear();
            gallery.GalleryPerformers = dto.PerformerIds.Select(pid => new GalleryPerformer { PerformerId = pid, GalleryId = id }).ToList();
        }
        if (dto.SceneIds != null)
        {
            gallery.SceneGalleries.Clear();
            gallery.SceneGalleries = dto.SceneIds.Select(sid => new SceneGallery { SceneId = sid, GalleryId = id }).ToList();
        }
        if (dto.CustomFields != null) gallery.CustomFields = dto.CustomFields;

        await galleryRepo.UpdateAsync(gallery, ct);
        var updated = await galleryRepo.GetByIdWithRelationsAsync(id, ct);
        return Ok(MapToDto(updated!));
    }

    [HttpDelete("{id:int}")]
    public async Task<IActionResult> Delete(int id, CancellationToken ct)
    {
        var g = await galleryRepo.GetByIdAsync(id, ct);
        if (g == null) return NotFound();
        await galleryRepo.DeleteAsync(id, ct);
        return NoContent();
    }

    private static GalleryDto MapToDto(Gallery g, int? imageCount = null, int? sceneCount = null) => new(
        g.Id, g.Title, g.Code, g.Date?.ToString("yyyy-MM-dd"), g.Details, g.Photographer,
        g.Rating, g.Organized, g.StudioId, g.Studio?.Name,
        g.Urls.Select(u => u.Url).ToList(),
        g.GalleryTags.Where(gt => gt.Tag != null).Select(gt => new TagDto(gt.Tag!.Id, gt.Tag.Name, gt.Tag.Description, gt.Tag.Favorite, gt.Tag.IgnoreAutoTag, [])).ToList(),
        g.GalleryPerformers.Where(gp => gp.Performer != null).Select(gp => new PerformerSummaryDto(gp.Performer!.Id, gp.Performer.Name, gp.Performer.Disambiguation, gp.Performer.Gender?.ToString(), gp.Performer.Favorite, gp.Performer.ImageBlobId != null ? $"/api/performers/{gp.Performer.Id}/image" : null)).ToList(),
        imageCount ?? g.ImageGalleries?.Count ?? 0,
        sceneCount ?? g.SceneGalleries?.Count ?? 0,
        g.SceneGalleries?.Select(sg => sg.SceneId).ToList() ?? [],
        g.Folder?.Path,
        g.Files?.Select(f => new GalleryFileInfoDto(f.Id, f.Path, f.Size, f.ModTime.ToString("o"),
            f.Fingerprints?.Select(fp => new FingerprintDto(fp.Type, fp.Value)).ToList() ?? [])).ToList() ?? [],
        g.CustomFields,
        g.CreatedAt.ToString("o"), g.UpdatedAt.ToString("o")
    );

    private async Task<List<GalleryDto>> MapListToDtos(IReadOnlyList<Gallery> items, CancellationToken ct)
    {
        if (items.Count == 0) return [];
        var ids = items.Select(g => g.Id).ToList();
        var imgCounts = await db.Set<ImageGallery>().Where(ig => ids.Contains(ig.GalleryId))
            .GroupBy(ig => ig.GalleryId).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        var sceneCounts = await db.Set<SceneGallery>().Where(sg => ids.Contains(sg.GalleryId))
            .GroupBy(sg => sg.GalleryId).Select(g => new { Id = g.Key, Count = g.Count() }).ToDictionaryAsync(x => x.Id, x => x.Count, ct);
        return items.Select(g => MapToDto(g,
            imgCounts.GetValueOrDefault(g.Id, 0),
            sceneCounts.GetValueOrDefault(g.Id, 0)
        )).ToList();
    }

    private static DateOnly? ParseDate(string? date) => DateOnly.TryParse(date, out var d) ? d : null;
    private static List<int>? ParseIntList(string? csv) => string.IsNullOrEmpty(csv) ? null : csv.Split(',').Select(int.Parse).ToList();

    // ===== Image Management =====

    [HttpPost("{id:int}/images")]
    public async Task<IActionResult> AddImages(int id, [FromBody] GalleryAddImagesDto dto, CancellationToken ct)
    {
        var gallery = await db.Galleries.Include(g => g.ImageGalleries).FirstOrDefaultAsync(g => g.Id == id, ct);
        if (gallery == null) return NotFound();

        var existing = gallery.ImageGalleries.Select(ig => ig.ImageId).ToHashSet();
        foreach (var imageId in dto.ImageIds.Where(iid => !existing.Contains(iid)))
            gallery.ImageGalleries.Add(new ImageGallery { ImageId = imageId, GalleryId = id });

        await db.SaveChangesAsync(ct);
        return Ok(new { added = dto.ImageIds.Count });
    }

    [HttpDelete("{id:int}/images")]
    public async Task<IActionResult> RemoveImages(int id, [FromBody] GalleryRemoveImagesDto dto, CancellationToken ct)
    {
        var toRemove = await db.Set<ImageGallery>()
            .Where(ig => ig.GalleryId == id && dto.ImageIds.Contains(ig.ImageId))
            .ToListAsync(ct);

        db.Set<ImageGallery>().RemoveRange(toRemove);
        await db.SaveChangesAsync(ct);
        return Ok(new { removed = toRemove.Count });
    }

    // ===== Chapters =====

    [HttpGet("{id:int}/chapters")]
    [OutputCache(PolicyName = "ShortCache")]
    public async Task<ActionResult<List<GalleryChapterDto>>> GetChapters(int id, CancellationToken ct)
    {
        var chapters = await db.GalleryChapters
            .Where(c => c.GalleryId == id)
            .OrderBy(c => c.ImageIndex)
            .Select(c => new GalleryChapterDto(c.Id, c.Title, c.ImageIndex, c.GalleryId, c.CreatedAt.ToString("o"), c.UpdatedAt.ToString("o")))
            .ToListAsync(ct);

        return Ok(chapters);
    }

    [HttpPost("{id:int}/chapters")]
    public async Task<ActionResult<GalleryChapterDto>> CreateChapter(int id, [FromBody] GalleryChapterCreateDto dto, CancellationToken ct)
    {
        var gallery = await db.Galleries.FindAsync([id], ct);
        if (gallery == null) return NotFound();

        var chapter = new GalleryChapter { Title = dto.Title, ImageIndex = dto.ImageIndex, GalleryId = id };
        db.GalleryChapters.Add(chapter);
        await db.SaveChangesAsync(ct);
        return CreatedAtAction(nameof(GetChapters), new { id }, new GalleryChapterDto(chapter.Id, chapter.Title, chapter.ImageIndex, chapter.GalleryId, chapter.CreatedAt.ToString("o"), chapter.UpdatedAt.ToString("o")));
    }

    [HttpPut("{galleryId:int}/chapters/{chapterId:int}")]
    public async Task<ActionResult<GalleryChapterDto>> UpdateChapter(int galleryId, int chapterId, [FromBody] GalleryChapterUpdateDto dto, CancellationToken ct)
    {
        var chapter = await db.GalleryChapters.FirstOrDefaultAsync(c => c.Id == chapterId && c.GalleryId == galleryId, ct);
        if (chapter == null) return NotFound();

        if (dto.Title != null) chapter.Title = dto.Title;
        if (dto.ImageIndex.HasValue) chapter.ImageIndex = dto.ImageIndex.Value;
        await db.SaveChangesAsync(ct);
        return Ok(new GalleryChapterDto(chapter.Id, chapter.Title, chapter.ImageIndex, chapter.GalleryId, chapter.CreatedAt.ToString("o"), chapter.UpdatedAt.ToString("o")));
    }

    [HttpDelete("{galleryId:int}/chapters/{chapterId:int}")]
    public async Task<IActionResult> DeleteChapter(int galleryId, int chapterId, CancellationToken ct)
    {
        var chapter = await db.GalleryChapters.FirstOrDefaultAsync(c => c.Id == chapterId && c.GalleryId == galleryId, ct);
        if (chapter == null) return NotFound();
        db.GalleryChapters.Remove(chapter);
        await db.SaveChangesAsync(ct);
        return NoContent();
    }

    // ===== Bulk Operations =====

    [HttpPost("bulk")]
    public async Task<IActionResult> BulkUpdate([FromBody] BulkGalleryUpdateDto dto, CancellationToken ct)
    {
        var galleries = await db.Galleries
            .Include(g => g.GalleryTags)
            .Include(g => g.GalleryPerformers)
            .Where(g => dto.Ids.Contains(g.Id))
            .ToListAsync(ct);

        foreach (var gallery in galleries)
        {
            if (dto.Rating.HasValue) gallery.Rating = dto.Rating;
            if (dto.Organized.HasValue) gallery.Organized = dto.Organized.Value;
            if (dto.StudioId.HasValue) gallery.StudioId = dto.StudioId;

            if (dto.TagIds != null && dto.TagMode == BulkUpdateMode.Set)
            {
                gallery.GalleryTags.Clear();
                gallery.GalleryTags = dto.TagIds.Select(tid => new GalleryTag { TagId = tid, GalleryId = gallery.Id }).ToList();
            }
            else if (dto.TagIds != null && dto.TagMode == BulkUpdateMode.Add)
            {
                var existing = gallery.GalleryTags.Select(gt => gt.TagId).ToHashSet();
                foreach (var tid in dto.TagIds.Where(t => !existing.Contains(t)))
                    gallery.GalleryTags.Add(new GalleryTag { TagId = tid, GalleryId = gallery.Id });
            }
            else if (dto.TagIds != null && dto.TagMode == BulkUpdateMode.Remove)
            {
                gallery.GalleryTags = gallery.GalleryTags.Where(gt => !dto.TagIds.Contains(gt.TagId)).ToList();
            }
        }

        await db.SaveChangesAsync(ct);
        return Ok(new { updated = galleries.Count });
    }
}
