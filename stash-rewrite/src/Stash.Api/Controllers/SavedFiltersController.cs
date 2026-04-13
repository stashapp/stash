using Microsoft.AspNetCore.Mvc;
using Stash.Core.DTOs;
using Stash.Core.Entities;
using Stash.Core.Enums;
using Stash.Core.Interfaces;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class SavedFiltersController(ISavedFilterRepository filterRepo) : ControllerBase
{
    [HttpGet]
    public async Task<ActionResult<IReadOnlyList<SavedFilterDto>>> GetAll([FromQuery] string? mode, CancellationToken ct)
    {
        IReadOnlyList<SavedFilter> filters;
        if (mode != null && Enum.TryParse<FilterMode>(mode, true, out var filterMode))
            filters = await filterRepo.GetByModeAsync(filterMode, ct);
        else
            filters = await filterRepo.GetAllAsync(ct);

        return Ok(filters.Select(MapToDto).ToList());
    }

    [HttpGet("{id:int}")]
    public async Task<ActionResult<SavedFilterDto>> GetById(int id, CancellationToken ct)
    {
        var filter = await filterRepo.GetByIdAsync(id, ct);
        if (filter == null) return NotFound();
        return Ok(MapToDto(filter));
    }

    [HttpPost]
    public async Task<ActionResult<SavedFilterDto>> Create([FromBody] SavedFilterCreateDto dto, CancellationToken ct)
    {
        if (!Enum.TryParse<FilterMode>(dto.Mode, true, out var filterMode))
            return BadRequest(new { message = $"Invalid filter mode: {dto.Mode}" });

        var filter = new SavedFilter
        {
            Name = dto.Name, Mode = filterMode,
            FindFilter = dto.FindFilter, ObjectFilter = dto.ObjectFilter, UIOptions = dto.UIOptions
        };

        filter = await filterRepo.AddAsync(filter, ct);
        return CreatedAtAction(nameof(GetById), new { id = filter.Id }, MapToDto(filter));
    }

    [HttpPut("{id:int}")]
    public async Task<ActionResult<SavedFilterDto>> Update(int id, [FromBody] SavedFilterUpdateDto dto, CancellationToken ct)
    {
        var filter = await filterRepo.GetByIdAsync(id, ct);
        if (filter == null) return NotFound();

        if (dto.Name != null) filter.Name = dto.Name;
        if (dto.Mode != null && Enum.TryParse<FilterMode>(dto.Mode, true, out var mode)) filter.Mode = mode;
        if (dto.FindFilter != null) filter.FindFilter = dto.FindFilter;
        if (dto.ObjectFilter != null) filter.ObjectFilter = dto.ObjectFilter;
        if (dto.UIOptions != null) filter.UIOptions = dto.UIOptions;

        await filterRepo.UpdateAsync(filter, ct);
        return Ok(MapToDto(filter));
    }

    [HttpDelete("{id:int}")]
    public async Task<IActionResult> Delete(int id, CancellationToken ct)
    {
        var f = await filterRepo.GetByIdAsync(id, ct);
        if (f == null) return NotFound();
        await filterRepo.DeleteAsync(id, ct);
        return NoContent();
    }

    // ===== Default Filters =====

    [HttpGet("default/{mode}")]
    public async Task<ActionResult<SavedFilterDto?>> GetDefault(string mode, CancellationToken ct)
    {
        if (!Enum.TryParse<FilterMode>(mode, true, out var filterMode))
            return BadRequest(new { message = $"Invalid filter mode: {mode}" });

        var filters = await filterRepo.GetByModeAsync(filterMode, ct);
        var defaultFilter = filters.FirstOrDefault(f => f.Name == $"__default_{mode}");
        return defaultFilter != null ? Ok(MapToDto(defaultFilter)) : Ok((SavedFilterDto?)null);
    }

    [HttpPut("default/{mode}")]
    public async Task<ActionResult<SavedFilterDto>> SetDefault(string mode, [FromBody] SetDefaultFilterDto dto, CancellationToken ct)
    {
        if (!Enum.TryParse<FilterMode>(mode, true, out var filterMode))
            return BadRequest(new { message = $"Invalid filter mode: {mode}" });

        var name = $"__default_{mode}";
        var filters = await filterRepo.GetByModeAsync(filterMode, ct);
        var existing = filters.FirstOrDefault(f => f.Name == name);

        if (dto.FilterId.HasValue)
        {
            var source = await filterRepo.GetByIdAsync(dto.FilterId.Value, ct);
            if (source == null) return NotFound("Source filter not found");

            if (existing != null)
            {
                existing.FindFilter = source.FindFilter;
                existing.ObjectFilter = source.ObjectFilter;
                existing.UIOptions = source.UIOptions;
                await filterRepo.UpdateAsync(existing, ct);
                return Ok(MapToDto(existing));
            }

            var newDefault = new SavedFilter
            {
                Name = name, Mode = filterMode,
                FindFilter = source.FindFilter, ObjectFilter = source.ObjectFilter, UIOptions = source.UIOptions
            };
            newDefault = await filterRepo.AddAsync(newDefault, ct);
            return Ok(MapToDto(newDefault));
        }

        // Clear default
        if (existing != null)
            await filterRepo.DeleteAsync(existing.Id, ct);

        return Ok((SavedFilterDto?)null);
    }

    private static SavedFilterDto MapToDto(SavedFilter f) => new(
        f.Id, f.Mode.ToString(), f.Name, f.FindFilter, f.ObjectFilter, f.UIOptions);
}

public record SavedFilterDto(int Id, string Mode, string Name, string? FindFilter, string? ObjectFilter, string? UIOptions);
public record SavedFilterCreateDto(string Mode, string Name, string? FindFilter, string? ObjectFilter, string? UIOptions);
public record SavedFilterUpdateDto(string? Mode, string? Name, string? FindFilter, string? ObjectFilter, string? UIOptions);
