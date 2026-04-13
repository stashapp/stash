using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Routing;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;

namespace Stash.Plugins;

/// <summary>
/// Base interface for all Stash extensions.
/// Extensions register services, add API endpoints, and define UI contributions.
/// </summary>
public interface IExtension
{
    /// <summary>Unique identifier for the extension (e.g. "com.stash.core-scenes")</summary>
    string Id { get; }

    /// <summary>Human-readable display name</summary>
    string Name { get; }

    /// <summary>Semantic version string</summary>
    string Version { get; }

    /// <summary>Short description of what the extension provides</summary>
    string? Description { get; }

    /// <summary>
    /// Called during host building. Register services, repositories, etc.
    /// </summary>
    void ConfigureServices(IServiceCollection services, ExtensionContext context);

    /// <summary>
    /// Called after the app is built but before it starts serving requests.
    /// Use for initialization, seeding, etc.
    /// </summary>
    Task InitializeAsync(IServiceProvider services, CancellationToken ct = default) => Task.CompletedTask;
}

/// <summary>
/// Context passed to extensions during registration, providing access to
/// shared configuration and extension coordination.
/// </summary>
public class ExtensionContext
{
    /// <summary>The app's configuration root</summary>
    public required IConfiguration Configuration { get; init; }

    /// <summary>Extension data directory (for storing extension-specific files)</summary>
    public required string DataDirectory { get; init; }

    /// <summary>Registry for UI contributions (pages, nav items, slots)</summary>
    public UIRegistry UI { get; } = new();
}

/// <summary>
/// Extensions that contribute API endpoints implement this interface.
/// </summary>
public interface IApiExtension : IExtension
{
    /// <summary>
    /// Called to map additional API endpoints (minimal APIs or controller discovery).
    /// </summary>
    void MapEndpoints(IEndpointRouteBuilder endpoints);
}

/// <summary>
/// Extensions that contribute to the UI implement this interface.
/// The UI manifest is serialized to JSON and served to the frontend.
/// </summary>
public interface IUIExtension : IExtension
{
    /// <summary>
    /// Return the UI manifest describing pages, nav items, and slot contributions.
    /// </summary>
    UIManifest GetUIManifest();
}

/// <summary>
/// Describes a UI page contributed by an extension.
/// </summary>
public record UIPageDefinition(
    string Route,
    string Label,
    string? Icon = null,
    string? DetailRoute = null,
    bool ShowInNav = true,
    int NavOrder = 100,
    string? RequiredPermission = null
);

/// <summary>
/// Complete UI manifest for an extension.
/// </summary>
public class UIManifest
{
    public List<UIPageDefinition> Pages { get; set; } = [];
    public Dictionary<string, object> Settings { get; set; } = [];
}

/// <summary>
/// Central registry for UI contributions, aggregated from all extensions.
/// </summary>
public class UIRegistry
{
    private readonly List<UIPageDefinition> _pages = [];

    public IReadOnlyList<UIPageDefinition> Pages => _pages;

    public void RegisterPage(UIPageDefinition page) => _pages.Add(page);

    public UIManifest ToManifest() => new() { Pages = [.. _pages] };
}
