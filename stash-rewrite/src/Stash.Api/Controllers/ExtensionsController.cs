using Microsoft.AspNetCore.Mvc;
using Stash.Plugins;

namespace Stash.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class ExtensionsController(ExtensionManager extensionManager) : ControllerBase
{
    /// <summary>Returns the aggregated UI manifest from all registered extensions.</summary>
    [HttpGet("manifest")]
    public ActionResult<UIManifest> GetManifest() =>
        Ok(extensionManager.GetAggregatedManifest());

    /// <summary>Returns a list of all registered extensions with capability info.</summary>
    [HttpGet]
    public ActionResult<IEnumerable<ExtensionInfo>> GetExtensions() =>
        Ok(extensionManager.Extensions.Select(e => new ExtensionInfo(
            e.Id,
            e.Name,
            e.Version,
            e.Description,
            e.Author,
            e.IconUrl,
            extensionManager.IsEnabled(e.Id),
            e is IUIExtension,
            e is IApiExtension,
            e is IStatefulExtension,
            e is IJobExtension,
            e is IEventExtension,
            e is IJobExtension je ? je.Jobs.Select(j => new JobInfo(j.Id, j.Name, j.Description)).ToList() : [])));

    /// <summary>Enable an extension.</summary>
    [HttpPost("{id}/enable")]
    public IActionResult Enable(string id)
    {
        var ext = extensionManager.Extensions.FirstOrDefault(e => e.Id == id);
        if (ext == null) return NotFound();
        extensionManager.EnableExtension(id);
        return Ok();
    }

    /// <summary>Disable an extension.</summary>
    [HttpPost("{id}/disable")]
    public IActionResult Disable(string id)
    {
        var ext = extensionManager.Extensions.FirstOrDefault(e => e.Id == id);
        if (ext == null) return NotFound();
        extensionManager.DisableExtension(id);
        return Ok();
    }

    /// <summary>Get extension key-value store data.</summary>
    [HttpGet("{id}/data")]
    public async Task<IActionResult> GetData(string id, CancellationToken ct)
    {
        var ext = extensionManager.Extensions.FirstOrDefault(e => e.Id == id) as IStatefulExtension;
        if (ext == null) return NotFound("Extension not found or not stateful");

        var factory = HttpContext.RequestServices.GetService<IExtensionStoreFactory>();
        if (factory == null) return StatusCode(500, "Store not available");

        var store = factory.CreateStore(id);
        var data = await store.GetAllAsync(ct);
        return Ok(data);
    }

    /// <summary>Set a key-value pair in extension store.</summary>
    [HttpPut("{id}/data/{key}")]
    public async Task<IActionResult> SetData(string id, string key, [FromBody] string value, CancellationToken ct)
    {
        var ext = extensionManager.Extensions.FirstOrDefault(e => e.Id == id) as IStatefulExtension;
        if (ext == null) return NotFound("Extension not found or not stateful");

        var factory = HttpContext.RequestServices.GetService<IExtensionStoreFactory>();
        if (factory == null) return StatusCode(500, "Store not available");

        var store = factory.CreateStore(id);
        await store.SetAsync(key, value, ct);
        return Ok();
    }

    /// <summary>Trigger a job defined by an extension.</summary>
    [HttpPost("{id}/jobs/{jobId}/run")]
    public async Task<IActionResult> RunJob(string id, string jobId, [FromBody] Dictionary<string, string>? parameters, CancellationToken ct)
    {
        var ext = extensionManager.Extensions.OfType<IJobExtension>().FirstOrDefault(e => e.Id == id);
        if (ext == null) return NotFound("Extension not found or has no jobs");

        var job = ext.Jobs.FirstOrDefault(j => j.Id == jobId);
        if (job == null) return NotFound($"Job '{jobId}' not found");

        // Run job in background
        _ = Task.Run(async () =>
        {
            var progress = new SimpleJobProgress();
            await ext.RunJobAsync(jobId, parameters, progress, CancellationToken.None);
        }, CancellationToken.None);

        return Accepted(new { message = $"Job '{job.Name}' started" });
    }

    /// <summary>Serve static assets from an extension's data directory.</summary>
    [HttpGet("assets/{extensionId}/{**path}")]
    public IActionResult GetAsset(string extensionId, string path)
    {
        var ext = extensionManager.Extensions.FirstOrDefault(e => e.Id == extensionId);
        if (ext == null) return NotFound();

        var basePath = Path.Combine(extensionManager.Context.DataDirectory, extensionId);
        var fullPath = Path.GetFullPath(Path.Combine(basePath, path));

        // Security: prevent path traversal
        if (!fullPath.StartsWith(Path.GetFullPath(basePath)))
            return BadRequest("Invalid path");

        if (!System.IO.File.Exists(fullPath)) return NotFound();

        var contentType = Path.GetExtension(fullPath).ToLowerInvariant() switch
        {
            ".js" => "application/javascript",
            ".mjs" => "application/javascript",
            ".css" => "text/css",
            ".json" => "application/json",
            ".html" => "text/html",
            ".svg" => "image/svg+xml",
            ".png" => "image/png",
            ".jpg" or ".jpeg" => "image/jpeg",
            ".woff2" => "font/woff2",
            ".woff" => "font/woff",
            _ => "application/octet-stream"
        };

        return PhysicalFile(fullPath, contentType);
    }
}

public record ExtensionInfo(
    string Id,
    string Name,
    string Version,
    string? Description,
    string? Author,
    string? IconUrl,
    bool Enabled,
    bool HasUI,
    bool HasApi,
    bool HasState,
    bool HasJobs,
    bool HasEvents,
    List<JobInfo> Jobs);

public record JobInfo(string Id, string Name, string? Description);

/// <summary>Simple job progress reporter.</summary>
internal class SimpleJobProgress : IJobProgress
{
    public double Percent { get; private set; }
    public string? Message { get; private set; }
    public void Report(double percent, string? message = null) { Percent = percent; Message = message; }
}
