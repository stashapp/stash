using System.Reflection;
using System.Runtime.Loader;
using Microsoft.AspNetCore.Routing;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;

namespace Stash.Plugins;

/// <summary>
/// Manages extension discovery, loading, lifecycle, and capability wiring.
/// </summary>
public class ExtensionManager
{
    private readonly List<IExtension> _extensions = [];
    private readonly ExtensionContext _context;
    private readonly Dictionary<string, AssemblyLoadContext> _loadContexts = [];
    private readonly HashSet<string> _disabledIds = [];
    private IServiceProvider? _lastServiceProvider;
    private ILogger<ExtensionManager>? _logger;

    public IReadOnlyList<IExtension> Extensions => _extensions;
    public ExtensionContext Context => _context;

    public ExtensionManager(ExtensionContext context)
    {
        _context = context;
    }

    /// <summary>Register an extension instance.</summary>
    public void Register(IExtension extension) => _extensions.Add(extension);

    /// <summary>
    /// Discover and load .NET extension assemblies from a directory.
    /// </summary>
    public void DiscoverExtensions(string extensionsDir)
    {
        if (!Directory.Exists(extensionsDir)) return;

        foreach (var dir in Directory.GetDirectories(extensionsDir))
        {
            var dlls = Directory.GetFiles(dir, "*.dll");
            foreach (var dll in dlls)
            {
                try
                {
                    var loadContext = new AssemblyLoadContext(Path.GetFileNameWithoutExtension(dll), isCollectible: true);
                    var assembly = loadContext.LoadFromAssemblyPath(dll);
                    var extensionTypes = assembly.GetTypes()
                        .Where(t => typeof(IExtension).IsAssignableFrom(t) && !t.IsAbstract && !t.IsInterface);

                    foreach (var type in extensionTypes)
                    {
                        if (Activator.CreateInstance(type) is IExtension ext)
                        {
                            _extensions.Add(ext);
                            _loadContexts[ext.Id] = loadContext;
                        }
                    }
                }
                catch
                {
                    // Skip DLLs that can't be loaded as extensions
                }
            }
        }
    }

    /// <summary>Call ConfigureServices on all registered extensions.</summary>
    public void ConfigureServices(IServiceCollection services)
    {
        foreach (var ext in _extensions)
            ext.ConfigureServices(services, _context);
    }

    /// <summary>Map API endpoints from all IApiExtension instances.</summary>
    public void MapEndpoints(IEndpointRouteBuilder endpoints)
    {
        foreach (var ext in _extensions.OfType<IApiExtension>())
            ext.MapEndpoints(endpoints);
    }

    /// <summary>
    /// Initialize all extensions after the app is built.
    /// Wires up capability interfaces (stores, event subscriptions, etc).
    /// </summary>
    public async Task InitializeAllAsync(IServiceProvider services, CancellationToken ct = default)
    {
        _lastServiceProvider = services;
        _logger = services.GetService<ILogger<ExtensionManager>>();

        // Wire stateful extensions with their DB-backed stores
        WireStatefulExtensions(services);

        // Initialize all extensions
        foreach (var ext in _extensions)
        {
            try
            {
                await ext.InitializeAsync(services, ct);
                _logger?.LogInformation("Extension {Id} ({Name} v{Version}) initialized", ext.Id, ext.Name, ext.Version);
            }
            catch (Exception ex)
            {
                _logger?.LogError(ex, "Failed to initialize extension {Id}", ext.Id);
            }
        }
    }

    /// <summary>
    /// Shut down all extensions gracefully.
    /// </summary>
    public async Task ShutdownAllAsync(CancellationToken ct = default)
    {
        foreach (var ext in _extensions)
        {
            try
            {
                await ext.ShutdownAsync(ct);
            }
            catch (Exception ex)
            {
                _logger?.LogError(ex, "Error shutting down extension {Id}", ext.Id);
            }
        }
    }

    /// <summary>Get the aggregated UI manifest from all enabled extensions.</summary>
    public UIManifest GetAggregatedManifest()
    {
        var manifest = _context.UI.ToManifest();
        foreach (var ext in _extensions.OfType<IUIExtension>())
        {
            if (_disabledIds.Contains(ext.Id)) continue;
            var extManifest = ext.GetUIManifest();
            manifest.Pages.AddRange(extManifest.Pages);
            manifest.Slots.AddRange(extManifest.Slots);
            manifest.Tabs.AddRange(extManifest.Tabs);
            manifest.Themes.AddRange(extManifest.Themes);
            manifest.SettingsPanels.AddRange(extManifest.SettingsPanels);
            manifest.PageOverrides.AddRange(extManifest.PageOverrides);
            manifest.DialogOverrides.AddRange(extManifest.DialogOverrides);
        }
        manifest.Pages.Sort((a, b) => a.NavOrder.CompareTo(b.NavOrder));
        manifest.Slots.Sort((a, b) => a.Order.CompareTo(b.Order));
        manifest.Tabs.Sort((a, b) => a.Order.CompareTo(b.Order));
        return manifest;
    }

    /// <summary>Get all registered extensions.</summary>
    public IReadOnlyList<IExtension> GetAllExtensions() => _extensions;

    /// <summary>Enable an extension by ID.</summary>
    public void EnableExtension(string id) => _disabledIds.Remove(id);

    /// <summary>Disable an extension by ID.</summary>
    public void DisableExtension(string id) => _disabledIds.Add(id);

    /// <summary>Check if an extension is enabled.</summary>
    public bool IsEnabled(string id) => !_disabledIds.Contains(id);

    /// <summary>Dispatch an entity event to all enabled IEventExtension instances.</summary>
    public async Task DispatchEventAsync(ExtensionEvent evt, CancellationToken ct = default)
    {
        foreach (var ext in _extensions.OfType<IEventExtension>())
        {
            if (_disabledIds.Contains(ext.Id)) continue;
            try
            {
                await ext.OnEventAsync(evt, ct);
            }
            catch (Exception ex)
            {
                _logger?.LogError(ex, "Extension {Id} failed handling event {EventType}", ext.Id, evt.EventType);
            }
        }
    }

    /// <summary>Get all job definitions across all enabled IJobExtension instances.</summary>
    public IEnumerable<(IJobExtension Extension, ExtensionJobDefinition Job)> GetAllJobs()
    {
        foreach (var ext in _extensions.OfType<IJobExtension>())
        {
            if (_disabledIds.Contains(ext.Id)) continue;
            foreach (var job in ext.Jobs)
                yield return (ext, job);
        }
    }

    /// <summary>Reload all extensions.</summary>
    public async Task ReloadAllAsync(CancellationToken ct = default)
    {
        if (_lastServiceProvider != null)
        {
            foreach (var ext in _extensions)
                await ext.InitializeAsync(_lastServiceProvider, ct);
        }
    }

    private void WireStatefulExtensions(IServiceProvider services)
    {
        // Look for the ExtensionStoreFactory in the DI container
        var factory = services.GetService<IExtensionStoreFactory>();
        if (factory is null)
        {
            _logger?.LogWarning("No IExtensionStoreFactory registered; stateful extensions won't have stores");
            return;
        }

        foreach (var ext in _extensions.OfType<IStatefulExtension>())
        {
            var store = factory.CreateStore(ext.Id);
            ext.SetStore(store);
            _logger?.LogDebug("Wired store for extension {Id}", ext.Id);
        }
    }
}

/// <summary>
/// Factory interface for creating extension stores. Implemented in Stash.Data.
/// </summary>
public interface IExtensionStoreFactory
{
    IExtensionStore CreateStore(string extensionId);
}
