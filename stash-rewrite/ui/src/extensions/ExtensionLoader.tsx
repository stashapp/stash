/**
 * Extension Runtime - Fetches the extension manifest and integrates all extension
 * contributions into the frontend: routes, slots, tabs, themes, page overrides,
 * dialog overrides, and settings panels.
 *
 * The architecture:
 * - Backend extensions declare UI contributions via UIManifest (JSON)
 * - This loader fetches the manifest once on mount
 * - Declarative contributions (pages, slots, tabs, themes, overrides) are registered
 *   into context-based registries consumed by the UI
 * - Component-based contributions reference built-in POC components (for built-in extensions)
 *   or would load from JS bundles (for external extensions)
 */
import { useEffect, useState, createContext, useContext, useCallback, type ReactNode, type FC } from "react";
import { useRouteRegistry } from "../router/RouteRegistry";
import { extensions } from "../api/client";
import type {
  ExtensionManifest,
  ExtensionThemeDef,
  ExtensionTabContribution,
  ExtensionPageOverride,
  ExtensionDialogOverride,
  ExtensionSettingsPanel,
} from "../api/types";

// ============================================================================
// Built-in POC extension components (registered by component name)
// These prove that the component resolution system works without external JS bundles.
// External extensions would register components via their JS bundle entry point.
// ============================================================================
import { SceneAnalyticsTab } from "./poc/SceneAnalyticsTab";
import { CustomHomeDashboard } from "./poc/CustomHomeDashboard";
import { SystemToolsPage } from "./poc/SystemToolsPage";
import { NotificationSettingsPanel } from "./poc/NotificationSettingsPanel";
import { EnhancedDeleteDialog } from "./poc/EnhancedDeleteDialog";

/** Global registry mapping componentName → React component. */
const componentRegistry = new Map<string, FC<any>>([
  ["SceneAnalyticsTab", SceneAnalyticsTab],
  ["CustomHomeDashboard", CustomHomeDashboard],
  ["SystemToolsPage", SystemToolsPage],
  ["NotificationSettingsPanel", NotificationSettingsPanel],
  ["EnhancedDeleteDialog", EnhancedDeleteDialog],
]);

/** Resolve a component by name from the registry. */
export function resolveComponent(name: string): FC<any> | undefined {
  return componentRegistry.get(name);
}

/** Register a component into the global registry (for external extensions). */
export function registerComponent(name: string, component: FC<any>) {
  componentRegistry.set(name, component);
}

// ============================================================================
// Extension context — everything the UI needs from the extension system
// ============================================================================
interface ExtensionState {
  manifest: ExtensionManifest | null;
  loaded: boolean;
  error?: string;
  activeThemeId: string | null;
  setActiveTheme: (id: string | null) => void;
  availableThemes: ExtensionThemeDef[];
  /** Tab contributions for a specific page type */
  getTabsForPage: (pageType: string) => ExtensionTabContribution[];
  /** Page override for a specific built-in page (highest priority wins) */
  getPageOverride: (targetPage: string) => ExtensionPageOverride | undefined;
  /** Dialog override for a specific dialog ID (highest priority wins) */
  getDialogOverride: (dialogId: string) => ExtensionDialogOverride | undefined;
  /** Settings panels contributed by extensions */
  settingsPanels: ExtensionSettingsPanel[];
  /** Resolve a React component by name */
  resolveComponent: (name: string) => FC<any> | undefined;
}

const ExtensionContext = createContext<ExtensionState>({
  manifest: null,
  loaded: false,
  activeThemeId: null,
  setActiveTheme: () => {},
  availableThemes: [],
  getTabsForPage: () => [],
  getPageOverride: () => undefined,
  getDialogOverride: () => undefined,
  settingsPanels: [],
  resolveComponent: () => undefined,
});

export function useExtensions() {
  return useContext(ExtensionContext);
}

const THEME_STORAGE_KEY = "stash-active-theme";

export function ExtensionLoaderProvider({ children }: { children: ReactNode }) {
  const { register, registerSlot } = useRouteRegistry();
  const [manifest, setManifest] = useState<ExtensionManifest | null>(null);
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState<string | undefined>();
  const [activeThemeId, setActiveThemeIdState] = useState<string | null>(
    () => localStorage.getItem(THEME_STORAGE_KEY)
  );

  const setActiveTheme = useCallback((id: string | null) => {
    setActiveThemeIdState(id);
    if (id) {
      localStorage.setItem(THEME_STORAGE_KEY, id);
    } else {
      localStorage.removeItem(THEME_STORAGE_KEY);
    }
  }, []);

  // Fetch manifest on mount
  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const m = await extensions.getManifest();
        if (cancelled) return;
        setManifest(m);

        // Register extension pages as routes
        for (const page of m.pages) {
          if (page.showInNav) {
            register({
              page: page.route,
              navItem: {
                page: page.route,
                label: page.label,
                icon: undefined,
                order: page.navOrder,
              },
            });
          }
        }

        // Register slot contributions
        for (const slot of m.slots) {
          if (slot.contentType === "html" && slot.html) {
            registerSlot({
              id: slot.id,
              slot: slot.slot,
              // eslint-disable-next-line react/no-danger
              render: () => <div dangerouslySetInnerHTML={{ __html: slot.html! }} />,
              order: slot.order,
            });
          } else if (slot.contentType === "component" && slot.componentName) {
            const Component = resolveComponent(slot.componentName);
            if (Component) {
              registerSlot({
                id: slot.id,
                slot: slot.slot,
                render: (props) => <Component {...props} />,
                order: slot.order,
              });
            }
          }
        }

        setLoaded(true);
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "Failed to load extensions");
          setLoaded(true);
        }
      }
    })();
    return () => { cancelled = true; };
  }, [register, registerSlot]);

  // Apply active theme CSS variables
  useEffect(() => {
    if (!manifest) return;

    const existingStyle = document.getElementById("stash-theme-override");
    if (existingStyle) existingStyle.remove();
    const existingLink = document.getElementById("stash-theme-css");
    if (existingLink) existingLink.remove();

    if (!activeThemeId) return;

    const theme = manifest.themes.find((t) => t.id === activeThemeId);
    if (!theme) return;

    if (theme.cssVariables && Object.keys(theme.cssVariables).length > 0) {
      const style = document.createElement("style");
      style.id = "stash-theme-override";
      const vars = Object.entries(theme.cssVariables)
        .map(([key, val]) => `  ${key}: ${val};`)
        .join("\n");
      style.textContent = `:root {\n${vars}\n}`;
      document.head.appendChild(style);
    }

    if (theme.cssUrl) {
      const link = document.createElement("link");
      link.id = "stash-theme-css";
      link.rel = "stylesheet";
      link.href = theme.cssUrl;
      document.head.appendChild(link);
    }

    return () => {
      document.getElementById("stash-theme-override")?.remove();
      document.getElementById("stash-theme-css")?.remove();
    };
  }, [activeThemeId, manifest]);

  // Derived lookups
  const getTabsForPage = useCallback(
    (pageType: string) =>
      manifest?.tabs.filter((t) => t.pageType === pageType) ?? [],
    [manifest]
  );

  const getPageOverride = useCallback(
    (targetPage: string) => {
      const overrides = manifest?.pageOverrides.filter((o) => o.targetPage === targetPage) ?? [];
      return overrides.sort((a, b) => b.priority - a.priority)[0];
    },
    [manifest]
  );

  const getDialogOverride = useCallback(
    (dialogId: string) => {
      const overrides = manifest?.dialogOverrides.filter((o) => o.dialogId === dialogId) ?? [];
      return overrides.sort((a, b) => b.priority - a.priority)[0];
    },
    [manifest]
  );

  const availableThemes = manifest?.themes ?? [];
  const settingsPanels = manifest?.settingsPanels ?? [];

  return (
    <ExtensionContext.Provider
      value={{
        manifest,
        loaded,
        error,
        activeThemeId,
        setActiveTheme,
        availableThemes,
        getTabsForPage,
        getPageOverride,
        getDialogOverride,
        settingsPanels,
        resolveComponent,
      }}
    >
      {children}
    </ExtensionContext.Provider>
  );
}
