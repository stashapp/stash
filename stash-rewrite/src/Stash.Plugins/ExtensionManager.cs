using Microsoft.AspNetCore.Routing;
using Microsoft.Extensions.DependencyInjection;

namespace Stash.Plugins;

/// <summary>
/// Manages extension discovery, loading, and lifecycle.
/// </summary>
public class ExtensionManager
{
    private readonly List<IExtension> _extensions = [];
    private readonly ExtensionContext _context;

    public IReadOnlyList<IExtension> Extensions => _extensions;
    public ExtensionContext Context => _context;

    public ExtensionManager(ExtensionContext context)
    {
        _context = context;
    }

    /// <summary>Register an extension instance.</summary>
    public void Register(IExtension extension) => _extensions.Add(extension);

    /// <summary>
    /// Call ConfigureServices on all registered extensions.
    /// </summary>
    public void ConfigureServices(IServiceCollection services)
    {
        foreach (var ext in _extensions)
            ext.ConfigureServices(services, _context);
    }

    /// <summary>
    /// Map API endpoints from all IApiExtension instances.
    /// </summary>
    public void MapEndpoints(IEndpointRouteBuilder endpoints)
    {
        foreach (var ext in _extensions.OfType<IApiExtension>())
            ext.MapEndpoints(endpoints);
    }

    /// <summary>
    /// Initialize all extensions after the app is built.
    /// </summary>
    public async Task InitializeAllAsync(IServiceProvider services, CancellationToken ct = default)
    {
        _lastServiceProvider = services;
        foreach (var ext in _extensions)
            await ext.InitializeAsync(services, ct);
    }

    /// <summary>
    /// Get the aggregated UI manifest from all extensions.
    /// </summary>
    public UIManifest GetAggregatedManifest()
    {
        var manifest = _context.UI.ToManifest();
        foreach (var ext in _extensions.OfType<IUIExtension>())
        {
            var extManifest = ext.GetUIManifest();
            manifest.Pages.AddRange(extManifest.Pages);
        }
        manifest.Pages.Sort((a, b) => a.NavOrder.CompareTo(b.NavOrder));
        return manifest;
    }

    /// <summary>Get all registered extensions.</summary>
    public IReadOnlyList<IExtension> GetAllExtensions() => _extensions;

    /// <summary>Enable a specific extension by ID.</summary>
    public void EnableExtension(string id)
    {
        var ext = _extensions.FirstOrDefault(e => e.Id == id);
        if (ext != null) _disabledIds.Remove(id);
    }

    /// <summary>Disable a specific extension by ID.</summary>
    public void DisableExtension(string id)
    {
        var ext = _extensions.FirstOrDefault(e => e.Id == id);
        if (ext != null) _disabledIds.Add(id);
    }

    /// <summary>Check if an extension is enabled.</summary>
    public bool IsEnabled(string id) => !_disabledIds.Contains(id);

    /// <summary>Reload all extensions.</summary>
    public async Task ReloadAllAsync(CancellationToken ct = default)
    {
        if (_lastServiceProvider != null)
        {
            foreach (var ext in _extensions)
                await ext.InitializeAsync(_lastServiceProvider, ct);
        }
    }

    private IServiceProvider? _lastServiceProvider;
    private readonly HashSet<string> _disabledIds = [];
}
