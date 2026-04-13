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

    /// <summary>Returns a list of all registered extensions.</summary>
    [HttpGet]
    public ActionResult<IEnumerable<ExtensionInfo>> GetExtensions() =>
        Ok(extensionManager.Extensions.Select(e => new ExtensionInfo(e.Id, e.Name, e.Version, e.Description)));
}

public record ExtensionInfo(string Id, string Name, string Version, string? Description);
