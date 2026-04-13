using Microsoft.AspNetCore.Mvc;
using Stash.Core.DTOs;

namespace Stash.Api.Controllers;

/// <summary>
/// DLNA (Digital Living Network Alliance) server endpoints for media streaming to compatible devices.
/// </summary>
[ApiController]
[Route("api/[controller]")]
public class DlnaController(ILogger<DlnaController> logger) : ControllerBase
{
    private static bool _running;
    private static DateTime? _untilDisabled;
    private static readonly List<string> RecentIps = [];
    private static readonly HashSet<string> AllowedIps = [];
    private static readonly HashSet<string> TempIps = [];
    private static readonly Lock StatusLock = new();

    [HttpGet("status")]
    public ActionResult<DlnaStatusDto> GetStatus()
    {
        lock (StatusLock)
        {
            return Ok(new DlnaStatusDto(
                _running,
                _untilDisabled?.ToString("o"),
                [.. RecentIps],
                [.. AllowedIps]
            ));
        }
    }

    [HttpPost("enable")]
    public ActionResult<DlnaStatusDto> Enable([FromBody] DlnaToggleDto? dto)
    {
        lock (StatusLock)
        {
            _running = true;
            _untilDisabled = dto?.DurationMinutes.HasValue == true
                ? DateTime.UtcNow.AddMinutes(dto.DurationMinutes!.Value)
                : null;
            logger.LogInformation("DLNA server enabled{Duration}", _untilDisabled != null ? $" until {_untilDisabled}" : "");
        }
        return GetStatus();
    }

    [HttpPost("disable")]
    public ActionResult<DlnaStatusDto> Disable([FromBody] DlnaToggleDto? dto)
    {
        lock (StatusLock)
        {
            _running = false;
            _untilDisabled = null;
            logger.LogInformation("DLNA server disabled");
        }
        return GetStatus();
    }

    [HttpPost("allow-ip")]
    public ActionResult<DlnaStatusDto> AllowIp([FromBody] DlnaIpDto dto)
    {
        if (string.IsNullOrWhiteSpace(dto.IpAddress))
            return BadRequest("IP address is required");

        lock (StatusLock)
        {
            if (dto.DurationMinutes.HasValue)
                TempIps.Add(dto.IpAddress);
            else
                AllowedIps.Add(dto.IpAddress);

            logger.LogInformation("DLNA: Allowed IP {Ip}", dto.IpAddress);
        }
        return GetStatus();
    }

    [HttpPost("remove-ip")]
    public ActionResult<DlnaStatusDto> RemoveIp([FromBody] DlnaIpDto dto)
    {
        if (string.IsNullOrWhiteSpace(dto.IpAddress))
            return BadRequest("IP address is required");

        lock (StatusLock)
        {
            AllowedIps.Remove(dto.IpAddress);
            TempIps.Remove(dto.IpAddress);
            logger.LogInformation("DLNA: Removed IP {Ip}", dto.IpAddress);
        }
        return GetStatus();
    }
}
