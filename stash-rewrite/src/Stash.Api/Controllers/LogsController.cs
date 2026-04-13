using Microsoft.AspNetCore.Mvc;
using Stash.Api.Services;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class LogsController : ControllerBase
{
    [HttpGet]
    public ActionResult<IReadOnlyList<LogEntry>> GetRecentLogs([FromQuery] string? level = null, [FromQuery] int limit = 200)
    {
        var logs = SignalRLogSink.GetRecentLogs();

        if (!string.IsNullOrEmpty(level))
            logs = logs.Where(l => l.Level.Equals(level, StringComparison.OrdinalIgnoreCase)).ToList();

        return Ok(logs.TakeLast(limit).ToList());
    }
}
