using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Stash.Core.Interfaces;
using Stash.Data;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api")]
public class EntityImageController(StashContext db, IBlobService blobService) : ControllerBase
{
    // ── Performers ──────────────────────────────────────────────

    [HttpPost("performers/{id:int}/image")]
    public async Task<IActionResult> UploadPerformerImage(int id, IFormFile file, CancellationToken ct)
    {
        if (!IsImage(file)) return BadRequest("File must be an image.");

        var entity = await db.Performers.FindAsync([id], ct);
        if (entity == null) return NotFound();

        if (entity.ImageBlobId != null)
            await blobService.DeleteBlobAsync(entity.ImageBlobId, ct);

        await using var stream = file.OpenReadStream();
        entity.ImageBlobId = await blobService.StoreBlobAsync(stream, file.ContentType, ct);
        await db.SaveChangesAsync(ct);

        return Ok(new { blobId = entity.ImageBlobId });
    }

    [HttpGet("performers/{id:int}/image")]
    public async Task<IActionResult> GetPerformerImage(int id, CancellationToken ct)
    {
        var entity = await db.Performers.FindAsync([id], ct);
        if (entity?.ImageBlobId == null) return NotFound();

        return await ServeBlobAsync(entity.ImageBlobId, ct);
    }

    [HttpDelete("performers/{id:int}/image")]
    public async Task<IActionResult> DeletePerformerImage(int id, CancellationToken ct)
    {
        var entity = await db.Performers.FindAsync([id], ct);
        if (entity?.ImageBlobId == null) return NotFound();

        await blobService.DeleteBlobAsync(entity.ImageBlobId, ct);
        entity.ImageBlobId = null;
        await db.SaveChangesAsync(ct);

        return NoContent();
    }

    // ── Studios ─────────────────────────────────────────────────

    [HttpPost("studios/{id:int}/image")]
    public async Task<IActionResult> UploadStudioImage(int id, IFormFile file, CancellationToken ct)
    {
        if (!IsImage(file)) return BadRequest("File must be an image.");

        var entity = await db.Studios.FindAsync([id], ct);
        if (entity == null) return NotFound();

        if (entity.ImageBlobId != null)
            await blobService.DeleteBlobAsync(entity.ImageBlobId, ct);

        await using var stream = file.OpenReadStream();
        entity.ImageBlobId = await blobService.StoreBlobAsync(stream, file.ContentType, ct);
        await db.SaveChangesAsync(ct);

        return Ok(new { blobId = entity.ImageBlobId });
    }

    [HttpGet("studios/{id:int}/image")]
    public async Task<IActionResult> GetStudioImage(int id, CancellationToken ct)
    {
        var entity = await db.Studios.FindAsync([id], ct);
        if (entity?.ImageBlobId == null) return NotFound();

        return await ServeBlobAsync(entity.ImageBlobId, ct);
    }

    [HttpDelete("studios/{id:int}/image")]
    public async Task<IActionResult> DeleteStudioImage(int id, CancellationToken ct)
    {
        var entity = await db.Studios.FindAsync([id], ct);
        if (entity?.ImageBlobId == null) return NotFound();

        await blobService.DeleteBlobAsync(entity.ImageBlobId, ct);
        entity.ImageBlobId = null;
        await db.SaveChangesAsync(ct);

        return NoContent();
    }

    // ── Tags ────────────────────────────────────────────────────

    [HttpPost("tags/{id:int}/image")]
    public async Task<IActionResult> UploadTagImage(int id, IFormFile file, CancellationToken ct)
    {
        if (!IsImage(file)) return BadRequest("File must be an image.");

        var entity = await db.Tags.FindAsync([id], ct);
        if (entity == null) return NotFound();

        if (entity.ImageBlobId != null)
            await blobService.DeleteBlobAsync(entity.ImageBlobId, ct);

        await using var stream = file.OpenReadStream();
        entity.ImageBlobId = await blobService.StoreBlobAsync(stream, file.ContentType, ct);
        await db.SaveChangesAsync(ct);

        return Ok(new { blobId = entity.ImageBlobId });
    }

    [HttpGet("tags/{id:int}/image")]
    public async Task<IActionResult> GetTagImage(int id, CancellationToken ct)
    {
        var entity = await db.Tags.FindAsync([id], ct);
        if (entity?.ImageBlobId == null) return NotFound();

        return await ServeBlobAsync(entity.ImageBlobId, ct);
    }

    [HttpDelete("tags/{id:int}/image")]
    public async Task<IActionResult> DeleteTagImage(int id, CancellationToken ct)
    {
        var entity = await db.Tags.FindAsync([id], ct);
        if (entity?.ImageBlobId == null) return NotFound();

        await blobService.DeleteBlobAsync(entity.ImageBlobId, ct);
        entity.ImageBlobId = null;
        await db.SaveChangesAsync(ct);

        return NoContent();
    }

    // ── Groups (front) ──────────────────────────────────────────

    [HttpPost("groups/{id:int}/image/front")]
    public async Task<IActionResult> UploadGroupFrontImage(int id, IFormFile file, CancellationToken ct)
    {
        if (!IsImage(file)) return BadRequest("File must be an image.");

        var entity = await db.Groups.FindAsync([id], ct);
        if (entity == null) return NotFound();

        if (entity.FrontImageBlobId != null)
            await blobService.DeleteBlobAsync(entity.FrontImageBlobId, ct);

        await using var stream = file.OpenReadStream();
        entity.FrontImageBlobId = await blobService.StoreBlobAsync(stream, file.ContentType, ct);
        await db.SaveChangesAsync(ct);

        return Ok(new { blobId = entity.FrontImageBlobId });
    }

    [HttpGet("groups/{id:int}/image/front")]
    public async Task<IActionResult> GetGroupFrontImage(int id, CancellationToken ct)
    {
        var entity = await db.Groups.FindAsync([id], ct);
        if (entity?.FrontImageBlobId == null) return NotFound();

        return await ServeBlobAsync(entity.FrontImageBlobId, ct);
    }

    [HttpDelete("groups/{id:int}/image/front")]
    public async Task<IActionResult> DeleteGroupFrontImage(int id, CancellationToken ct)
    {
        var entity = await db.Groups.FindAsync([id], ct);
        if (entity?.FrontImageBlobId == null) return NotFound();

        await blobService.DeleteBlobAsync(entity.FrontImageBlobId, ct);
        entity.FrontImageBlobId = null;
        await db.SaveChangesAsync(ct);

        return NoContent();
    }

    // ── Groups (back) ───────────────────────────────────────────

    [HttpPost("groups/{id:int}/image/back")]
    public async Task<IActionResult> UploadGroupBackImage(int id, IFormFile file, CancellationToken ct)
    {
        if (!IsImage(file)) return BadRequest("File must be an image.");

        var entity = await db.Groups.FindAsync([id], ct);
        if (entity == null) return NotFound();

        if (entity.BackImageBlobId != null)
            await blobService.DeleteBlobAsync(entity.BackImageBlobId, ct);

        await using var stream = file.OpenReadStream();
        entity.BackImageBlobId = await blobService.StoreBlobAsync(stream, file.ContentType, ct);
        await db.SaveChangesAsync(ct);

        return Ok(new { blobId = entity.BackImageBlobId });
    }

    [HttpGet("groups/{id:int}/image/back")]
    public async Task<IActionResult> GetGroupBackImage(int id, CancellationToken ct)
    {
        var entity = await db.Groups.FindAsync([id], ct);
        if (entity?.BackImageBlobId == null) return NotFound();

        return await ServeBlobAsync(entity.BackImageBlobId, ct);
    }

    [HttpDelete("groups/{id:int}/image/back")]
    public async Task<IActionResult> DeleteGroupBackImage(int id, CancellationToken ct)
    {
        var entity = await db.Groups.FindAsync([id], ct);
        if (entity?.BackImageBlobId == null) return NotFound();

        await blobService.DeleteBlobAsync(entity.BackImageBlobId, ct);
        entity.BackImageBlobId = null;
        await db.SaveChangesAsync(ct);

        return NoContent();
    }

    // ── Helpers ──────────────────────────────────────────────────

    private static bool IsImage(IFormFile file) =>
        file.ContentType.StartsWith("image/", StringComparison.OrdinalIgnoreCase);

    private async Task<IActionResult> ServeBlobAsync(string blobId, CancellationToken ct)
    {
        var result = await blobService.GetBlobAsync(blobId, ct);
        if (result == null) return NotFound();

        var (stream, contentType) = result.Value;
        Response.Headers.CacheControl = "public, max-age=86400";
        return File(stream, contentType);
    }
}
