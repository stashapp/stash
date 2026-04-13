using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Stash.Core.DTOs;
using Stash.Core.Entities;
using Stash.Core.Interfaces;
using Stash.Data;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/scenes/{sceneId:int}/markers")]
public class SceneMarkersController(ISceneMarkerRepository markerRepo, ISceneRepository sceneRepo, StashContext db) : ControllerBase
{
    /// <summary>Returns random scene markers for a wall/discovery view.</summary>
    [HttpGet("/api/markers/wall")]
    public async Task<ActionResult<List<SceneMarkerWallDto>>> MarkerWall([FromQuery] string? q, [FromQuery] int? tagId, [FromQuery] int count = 24, CancellationToken ct = default)
    {
        var query = db.SceneMarkers
            .Include(m => m.PrimaryTag)
            .Include(m => m.Scene).ThenInclude(s => s!.Files).ThenInclude(f => f.ParentFolder)
            .Include(m => m.SceneMarkerTags).ThenInclude(mt => mt.Tag)
            .AsNoTracking();

        if (!string.IsNullOrEmpty(q))
            query = query.Where(m => EF.Functions.ILike(m.Title, $"%{q}%"));
        if (tagId.HasValue)
            query = query.Where(m => m.PrimaryTagId == tagId.Value || m.SceneMarkerTags.Any(mt => mt.TagId == tagId.Value));

        var markers = await query.OrderBy(_ => EF.Functions.Random()).Take(count).ToListAsync(ct);
        return Ok(markers.Select(m => new SceneMarkerWallDto(
            m.Id, m.Title, m.Seconds, m.EndSeconds, m.PrimaryTagId, m.PrimaryTag?.Name ?? "",
            m.SceneId, m.Scene?.Title ?? "", m.Scene?.Files.FirstOrDefault()?.Path ?? "",
            m.SceneMarkerTags.Select(mt => new TagSummaryDto(mt.TagId, mt.Tag?.Name ?? "")).ToList()
        )).ToList());
    }

    [HttpGet]
    public async Task<ActionResult<IReadOnlyList<SceneMarkerSummaryDto>>> GetByScene(int sceneId, CancellationToken ct)
    {
        var markers = await markerRepo.GetBySceneIdAsync(sceneId, ct);
        return Ok(markers.Select(MapToDto).ToList());
    }

    [HttpGet("{id:int}")]
    public async Task<ActionResult<SceneMarkerSummaryDto>> GetById(int sceneId, int id, CancellationToken ct)
    {
        var marker = await markerRepo.GetByIdAsync(id, ct);
        if (marker == null || marker.SceneId != sceneId) return NotFound();
        return Ok(MapToDto(marker));
    }

    [HttpPost]
    public async Task<ActionResult<SceneMarkerSummaryDto>> Create(int sceneId, [FromBody] SceneMarkerCreateDto dto, CancellationToken ct)
    {
        var scene = await sceneRepo.GetByIdAsync(sceneId, ct);
        if (scene == null) return NotFound();

        var marker = new SceneMarker
        {
            Title = dto.Title, Seconds = dto.Seconds, EndSeconds = dto.EndSeconds,
            PrimaryTagId = dto.PrimaryTagId, SceneId = sceneId
        };
        if (dto.TagIds?.Count > 0)
            marker.SceneMarkerTags = dto.TagIds.Select(tid => new SceneMarkerTag { TagId = tid }).ToList();

        marker = await markerRepo.AddAsync(marker, ct);
        return CreatedAtAction(nameof(GetById), new { sceneId, id = marker.Id }, MapToDto(marker));
    }

    [HttpPut("{id:int}")]
    public async Task<ActionResult<SceneMarkerSummaryDto>> Update(int sceneId, int id, [FromBody] SceneMarkerUpdateDto dto, CancellationToken ct)
    {
        var marker = await markerRepo.GetByIdAsync(id, ct);
        if (marker == null || marker.SceneId != sceneId) return NotFound();

        if (dto.Title != null) marker.Title = dto.Title;
        if (dto.Seconds.HasValue) marker.Seconds = dto.Seconds.Value;
        if (dto.EndSeconds.HasValue) marker.EndSeconds = dto.EndSeconds;
        if (dto.PrimaryTagId.HasValue) marker.PrimaryTagId = dto.PrimaryTagId.Value;

        await markerRepo.UpdateAsync(marker, ct);
        return Ok(MapToDto(marker));
    }

    [HttpDelete("{id:int}")]
    public async Task<IActionResult> Delete(int sceneId, int id, CancellationToken ct)
    {
        var marker = await markerRepo.GetByIdAsync(id, ct);
        if (marker == null || marker.SceneId != sceneId) return NotFound();
        await markerRepo.DeleteAsync(id, ct);
        return NoContent();
    }

    private static SceneMarkerSummaryDto MapToDto(SceneMarker m) => new(
        m.Id, m.Title, m.Seconds, m.EndSeconds, m.PrimaryTagId, m.PrimaryTag?.Name ?? "");
}

// DTOs for scene markers (add to DTOs file or keep here for now)
public record SceneMarkerCreateDto(string Title, double Seconds, double? EndSeconds, int PrimaryTagId, List<int>? TagIds);
public record SceneMarkerUpdateDto(string? Title, double? Seconds, double? EndSeconds, int? PrimaryTagId);
