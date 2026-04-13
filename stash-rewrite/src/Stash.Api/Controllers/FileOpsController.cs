using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Stash.Core.DTOs;
using Stash.Core.Entities;
using Stash.Data;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/files")]
public class FileOpsController(StashContext db, ILogger<FileOpsController> logger) : ControllerBase
{
    [HttpPost("move")]
    public async Task<IActionResult> MoveFiles([FromBody] MoveFilesDto dto, CancellationToken ct)
    {
        if (!Directory.Exists(dto.DestinationPath))
            return BadRequest("Destination directory does not exist");

        var files = await db.Set<BaseFileEntity>()
            .Include(f => f.ParentFolder)
            .Where(f => dto.FileIds.Contains(f.Id))
            .ToListAsync(ct);

        var movedCount = 0;
        foreach (var file in files)
        {
            var oldPath = Path.Combine(file.ParentFolder?.Path ?? "", file.Basename);
            var newPath = Path.Combine(dto.DestinationPath, file.Basename);

            if (!System.IO.File.Exists(oldPath))
            {
                logger.LogWarning("Source file does not exist: {Path}", oldPath);
                continue;
            }

            if (System.IO.File.Exists(newPath))
            {
                logger.LogWarning("Destination file already exists: {Path}", newPath);
                continue;
            }

            System.IO.File.Move(oldPath, newPath);

            // Update folder reference
            var newFolder = await db.Folders.FirstOrDefaultAsync(f => f.Path == dto.DestinationPath, ct);
            if (newFolder == null)
            {
                newFolder = new Folder { Path = dto.DestinationPath, ModTime = DateTime.UtcNow };
                db.Folders.Add(newFolder);
                await db.SaveChangesAsync(ct);
            }
            file.ParentFolderId = newFolder.Id;
            movedCount++;
        }

        await db.SaveChangesAsync(ct);
        return Ok(new { moved = movedCount, total = files.Count });
    }

    [HttpPost("delete")]
    public async Task<IActionResult> DeleteFiles([FromBody] DeleteFilesDto dto, CancellationToken ct)
    {
        var files = await db.Set<BaseFileEntity>()
            .Include(f => f.ParentFolder)
            .Where(f => dto.FileIds.Contains(f.Id))
            .ToListAsync(ct);

        var deletedCount = 0;
        foreach (var file in files)
        {
            if (dto.DeleteFromDisk)
            {
                var path = Path.Combine(file.ParentFolder?.Path ?? "", file.Basename);
                if (System.IO.File.Exists(path))
                {
                    System.IO.File.Delete(path);
                    logger.LogInformation("Deleted file from disk: {Path}", path);
                }
            }

            db.Set<BaseFileEntity>().Remove(file);
            deletedCount++;
        }

        await db.SaveChangesAsync(ct);
        return Ok(new { deleted = deletedCount });
    }

    [HttpGet("browse")]
    public ActionResult<List<DirectoryEntryDto>> Browse([FromQuery] string? path)
    {
        var targetPath = path ?? Environment.GetFolderPath(Environment.SpecialFolder.UserProfile);
        if (!Directory.Exists(targetPath))
            return NotFound("Directory does not exist");

        var entries = new List<DirectoryEntryDto>();
        try
        {
            foreach (var dir in Directory.GetDirectories(targetPath))
                entries.Add(new DirectoryEntryDto(dir, true));
            foreach (var file in Directory.GetFiles(targetPath))
                entries.Add(new DirectoryEntryDto(file, false));
        }
        catch (UnauthorizedAccessException)
        {
            return Forbid();
        }

        return Ok(entries.OrderBy(e => !e.IsDirectory).ThenBy(e => e.Path).ToList());
    }
}
