using Microsoft.Extensions.DependencyInjection;
using Stash.Plugins;

namespace Stash.Api.Extensions;

// ============================================================================
// POC 1: Theme Extension — multiple theme definitions + CSS variables
// Proves: theme system, theme selector, CSS variable injection
// ============================================================================
public class ThemeCollectionExtension : IExtension, IUIExtension
{
    public string Id => "com.stash.themes";
    public string Name => "Theme Collection";
    public string Version => "1.0.0";
    public string? Description => "Built-in theme collection with multiple dark and light themes";
    public string? Author => "Stash";
    public string? IconUrl => null;

    public void ConfigureServices(IServiceCollection services, ExtensionContext context) { }

    public UIManifest GetUIManifest() => new()
    {
        Themes =
        [
            new UIThemeDefinition(
                Id: "dark-default",
                Name: "Dark (Default)",
                Description: "The default Stash dark theme"
            ),
            new UIThemeDefinition(
                Id: "dark-midnight",
                Name: "Dark Midnight",
                Description: "Deep midnight blue with purple accents",
                CssVariables: new()
                {
                    ["--color-plex-bg"] = "#0d1117",
                    ["--color-plex-nav"] = "#161b22",
                    ["--color-plex-card"] = "#1c2128",
                    ["--color-plex-surface"] = "#21262d",
                    ["--color-plex-border"] = "#30363d",
                    ["--color-plex-accent"] = "#8b5cf6",
                    ["--color-plex-accent-hover"] = "#7c3aed",
                    ["--color-plex-text"] = "#e6edf3",
                    ["--color-plex-text-secondary"] = "#8b949e",
                    ["--color-plex-text-muted"] = "#484f58",
                }
            ),
            new UIThemeDefinition(
                Id: "dark-emerald",
                Name: "Dark Emerald",
                Description: "Dark theme with emerald green accents",
                CssVariables: new()
                {
                    ["--color-plex-bg"] = "#0f1410",
                    ["--color-plex-nav"] = "#141c16",
                    ["--color-plex-card"] = "#1a241c",
                    ["--color-plex-surface"] = "#1f2b22",
                    ["--color-plex-border"] = "#2d3d30",
                    ["--color-plex-accent"] = "#10b981",
                    ["--color-plex-accent-hover"] = "#059669",
                    ["--color-plex-text"] = "#e6f0e8",
                    ["--color-plex-text-secondary"] = "#7ca38a",
                    ["--color-plex-text-muted"] = "#4a6350",
                }
            ),
            new UIThemeDefinition(
                Id: "dark-rose",
                Name: "Dark Rosé",
                Description: "Dark theme with warm rose accents",
                CssVariables: new()
                {
                    ["--color-plex-bg"] = "#140d0d",
                    ["--color-plex-nav"] = "#1c1414",
                    ["--color-plex-card"] = "#241a1a",
                    ["--color-plex-surface"] = "#2b1f1f",
                    ["--color-plex-border"] = "#3d2d2d",
                    ["--color-plex-accent"] = "#f43f5e",
                    ["--color-plex-accent-hover"] = "#e11d48",
                    ["--color-plex-text"] = "#f0e6e6",
                    ["--color-plex-text-secondary"] = "#a37c7c",
                    ["--color-plex-text-muted"] = "#634a4a",
                }
            ),
            new UIThemeDefinition(
                Id: "dark-ocean",
                Name: "Dark Ocean",
                Description: "Deep ocean blue theme",
                CssVariables: new()
                {
                    ["--color-plex-bg"] = "#0a1628",
                    ["--color-plex-nav"] = "#0f1d32",
                    ["--color-plex-card"] = "#14253d",
                    ["--color-plex-surface"] = "#192c47",
                    ["--color-plex-border"] = "#243a5c",
                    ["--color-plex-accent"] = "#0ea5e9",
                    ["--color-plex-accent-hover"] = "#0284c7",
                    ["--color-plex-text"] = "#e0f2fe",
                    ["--color-plex-text-secondary"] = "#7cacca",
                    ["--color-plex-text-muted"] = "#3b6685",
                }
            ),
        ]
    };
}

// ============================================================================
// POC 2: Tab Injection Extension — adds custom tabs to scene detail page
// Proves: tab contributions, component-based slots, API endpoints
// ============================================================================
public class SceneAnalyticsExtension : IExtension, IApiExtension, IUIExtension, IStatefulExtension
{
    public string Id => "com.stash.scene-analytics";
    public string Name => "Scene Analytics";
    public string Version => "1.0.0";
    public string? Description => "Adds play count tracking and analytics tab to scene details";
    public string? Author => "Stash";
    public string? IconUrl => null;
    private IExtensionStore? _store;

    public void SetStore(IExtensionStore store) => _store = store;
    public void ConfigureServices(IServiceCollection services, ExtensionContext context) { }

    public void MapEndpoints(IEndpointRouteBuilder endpoints)
    {
        var group = endpoints.MapGroup("/api/ext/analytics");

        group.MapGet("/scene/{sceneId:int}", async (int sceneId) =>
        {
            if (_store == null) return Results.StatusCode(500);
            var views = await _store.GetAsync($"scene:{sceneId}:views") ?? "0";
            var lastViewed = await _store.GetAsync($"scene:{sceneId}:last_viewed");
            return Results.Ok(new { sceneId, views = int.Parse(views), lastViewed });
        });

        group.MapPost("/scene/{sceneId:int}/view", async (int sceneId) =>
        {
            if (_store == null) return Results.StatusCode(500);
            var current = await _store.GetAsync($"scene:{sceneId}:views") ?? "0";
            await _store.SetAsync($"scene:{sceneId}:views", (int.Parse(current) + 1).ToString());
            await _store.SetAsync($"scene:{sceneId}:last_viewed", DateTime.UtcNow.ToString("O"));
            return Results.Ok();
        });
    }

    public UIManifest GetUIManifest() => new()
    {
        Tabs =
        [
            new UITabContribution(
                Key: "analytics",
                Label: "Analytics",
                PageType: "scene",
                ExtensionId: Id,
                ComponentName: "SceneAnalyticsTab",
                Order: 90
            ),
        ],
        Slots =
        [
            new UISlotContribution(
                Id: "analytics-badge",
                Slot: "scene-detail-actions",
                ExtensionId: Id,
                ContentType: "html",
                Html: """<span class="text-[10px] px-1.5 py-0.5 rounded bg-blue-600/20 text-blue-300 border border-blue-600/30" title="Scene Analytics Extension">📊</span>""",
                Order: 200
            ),
        ],
    };
}

// ============================================================================
// POC 3: Page Replacement Extension — replaces the home page
// Proves: page override capability, component-based page replacement
// ============================================================================
public class CustomHomeExtension : IExtension, IUIExtension
{
    public string Id => "com.stash.custom-home";
    public string Name => "Custom Home Page";
    public string Version => "1.0.0";
    public string? Description => "Replaces the default home page with an enhanced dashboard";
    public string? Author => "Stash";
    public string? IconUrl => null;

    public void ConfigureServices(IServiceCollection services, ExtensionContext context) { }

    public UIManifest GetUIManifest() => new()
    {
        PageOverrides =
        [
            new UIPageOverride(
                TargetPage: "home",
                ExtensionId: Id,
                ComponentName: "CustomHomeDashboard",
                Priority: 100
            ),
        ],
    };
}

// ============================================================================
// POC 4: New Page Extension — adds entirely new pages to navigation
// Proves: new page contribution, nav integration, multi-page extension
// ============================================================================
public class SystemToolsExtension : IExtension, IApiExtension, IUIExtension
{
    public string Id => "com.stash.system-tools";
    public string Name => "System Tools";
    public string Version => "1.0.0";
    public string? Description => "Adds a System Tools page with diagnostics and utilities";
    public string? Author => "Stash";
    public string? IconUrl => null;

    public void ConfigureServices(IServiceCollection services, ExtensionContext context) { }

    public void MapEndpoints(IEndpointRouteBuilder endpoints)
    {
        var group = endpoints.MapGroup("/api/ext/system-tools");

        group.MapGet("/info", () => Results.Ok(new
        {
            runtime = System.Runtime.InteropServices.RuntimeInformation.FrameworkDescription,
            os = System.Runtime.InteropServices.RuntimeInformation.OSDescription,
            cpuCount = Environment.ProcessorCount,
            workingSet = Environment.WorkingSet,
            uptime = Environment.TickCount64,
            gcMemory = GC.GetTotalMemory(false),
        }));

        group.MapGet("/extensions", (ExtensionManager mgr) =>
        {
            return Results.Ok(mgr.Extensions.Select(e => new
            {
                e.Id,
                e.Name,
                e.Version,
                e.Description,
                enabled = mgr.IsEnabled(e.Id),
                capabilities = new
                {
                    ui = e is IUIExtension,
                    api = e is IApiExtension,
                    stateful = e is IStatefulExtension,
                    jobs = e is IJobExtension,
                    events = e is IEventExtension,
                }
            }));
        });
    }

    public UIManifest GetUIManifest() => new()
    {
        Pages =
        [
            new UIPageDefinition(
                Route: "system-tools",
                Label: "System Tools",
                Icon: "wrench",
                ShowInNav: true,
                NavOrder: 80,
                ComponentName: "SystemToolsPage",
                ExtensionId: Id
            ),
        ],
    };
}

// ============================================================================
// POC 5: Settings Panel Extension — adds settings tab to extension config
// Proves: settings panel contributions, extension settings persistence
// ============================================================================
public class NotificationSettingsExtension : IExtension, IUIExtension, IStatefulExtension
{
    public string Id => "com.stash.notification-settings";
    public string Name => "Notification Settings";
    public string Version => "1.0.0";
    public string? Description => "Adds notification preferences to extension settings";
    public string? Author => "Stash";
    public string? IconUrl => null;
    private IExtensionStore? _store;

    public void SetStore(IExtensionStore store) => _store = store;
    public void ConfigureServices(IServiceCollection services, ExtensionContext context) { }

    public UIManifest GetUIManifest() => new()
    {
        SettingsPanels =
        [
            new UISettingsPanel(
                Id: "notification-prefs",
                Label: "Notifications",
                ExtensionId: Id,
                ComponentName: "NotificationSettingsPanel",
                Order: 10
            ),
        ],
    };
}

// ============================================================================
// POC 6: Dialog Override Extension — overrides the delete confirmation dialog
// Proves: dialog override capability
// ============================================================================
public class EnhancedDeleteDialogExtension : IExtension, IUIExtension
{
    public string Id => "com.stash.enhanced-delete-dialog";
    public string Name => "Enhanced Delete Dialog";
    public string Version => "1.0.0";
    public string? Description => "Replaces the default delete confirmation with an enhanced version showing what will be affected";
    public string? Author => "Stash";
    public string? IconUrl => null;

    public void ConfigureServices(IServiceCollection services, ExtensionContext context) { }

    public UIManifest GetUIManifest() => new()
    {
        DialogOverrides =
        [
            new UIDialogOverride(
                DialogId: "confirm-delete",
                ExtensionId: Id,
                ComponentName: "EnhancedDeleteDialog",
                Priority: 100
            ),
        ],
    };
}

// ============================================================================
// POC 7: Event Extension — demonstrates entity event subscription
// Proves: event hooks, stateful tracking
// ============================================================================
public class AuditLogExtension : IExtension, IEventExtension, IApiExtension, IStatefulExtension
{
    public string Id => "com.stash.audit-log";
    public string Name => "Audit Log";
    public string Version => "1.0.0";
    public string? Description => "Logs entity changes (create/update/delete) to an audit trail";
    public string? Author => "Stash";
    public string? IconUrl => null;
    private IExtensionStore? _store;

    public void SetStore(IExtensionStore store) => _store = store;
    public void ConfigureServices(IServiceCollection services, ExtensionContext context) { }

    public async Task OnEventAsync(ExtensionEvent evt, CancellationToken ct = default)
    {
        if (_store == null) return;
        var timestamp = DateTime.UtcNow.ToString("O");
        var key = $"log:{timestamp}:{evt.EventType}:{evt.EntityId}";
        await _store.SetAsync(key, System.Text.Json.JsonSerializer.Serialize(new
        {
            evt.EventType,
            evt.EntityType,
            evt.EntityId,
            Timestamp = timestamp,
        }), ct);
    }

    public void MapEndpoints(IEndpointRouteBuilder endpoints)
    {
        var group = endpoints.MapGroup("/api/ext/audit");

        group.MapGet("/log", async () =>
        {
            if (_store == null) return Results.StatusCode(500);
            var all = await _store.GetAllAsync();
            var logs = all.Where(kv => kv.Key.StartsWith("log:"))
                .OrderByDescending(kv => kv.Key)
                .Take(100)
                .Select(kv => System.Text.Json.JsonSerializer.Deserialize<object>(kv.Value))
                .ToList();
            return Results.Ok(logs);
        });
    }
}
