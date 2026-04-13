import { createContext, useContext, useState, useCallback, type ReactNode, type ComponentType } from "react";

export interface NavItem {
  page: string;
  label: string;
  icon?: ComponentType<{ className?: string }>;
  order?: number;
}

export interface RouteEntry {
  /** Page key used in route state */
  page: string;
  /** Component to render for list/collection views (no id) */
  component?: ComponentType<{ onNavigate: (r: any) => void }>;
  /** Component to render for detail views (with id) */
  detailComponent?: ComponentType<{ id: number; onNavigate: (r: any) => void }>;
  navItem?: NavItem;
}

export interface SlotEntry<TContext = any> {
  /** Unique id for this extension contribution */
  id: string;
  /** Named extension slot (e.g. "scene-detail-sidebar") */
  slot: string;
  /** Render function invoked with the host-provided context */
  render: (context: TContext) => ReactNode;
  /** Optional ordering. Lower values render first. */
  order?: number;
}

interface RouteRegistryContextValue {
  routes: RouteEntry[];
  slots: SlotEntry[];
  register: (entry: RouteEntry) => void;
  registerSlot: (entry: SlotEntry) => void;
}

const RouteRegistryContext = createContext<RouteRegistryContextValue | null>(null);

export function RouteRegistryProvider({ children }: { children: ReactNode }) {
  const [routes, setRoutes] = useState<RouteEntry[]>([]);
  const [slots, setSlots] = useState<SlotEntry[]>([]);

  const register = useCallback((entry: RouteEntry) => {
    setRoutes((prev) => {
      // Replace if same page key already registered
      const idx = prev.findIndex((r) => r.page === entry.page);
      if (idx >= 0) {
        const next = [...prev];
        next[idx] = entry;
        return next;
      }
      return [...prev, entry];
    });
  }, []);

  const registerSlot = useCallback((entry: SlotEntry) => {
    setSlots((prev) => {
      const idx = prev.findIndex((s) => s.id === entry.id);
      if (idx >= 0) {
        const next = [...prev];
        next[idx] = entry;
        return next;
      }
      return [...prev, entry];
    });
  }, []);

  return (
    <RouteRegistryContext.Provider value={{ routes, slots, register, registerSlot }}>
      {children}
    </RouteRegistryContext.Provider>
  );
}

export function useRouteRegistry() {
  const ctx = useContext(RouteRegistryContext);
  if (!ctx) throw new Error("useRouteRegistry must be used inside RouteRegistryProvider");
  return ctx;
}

export function ExtensionSlot<TContext>({ slot, context }: { slot: string; context: TContext }) {
  const { slots } = useRouteRegistry();
  const matching = slots
    .filter((s) => s.slot === slot)
    .sort((a, b) => (a.order ?? 100) - (b.order ?? 100));

  if (matching.length === 0) return null;

  return (
    <>
      {matching.map((entry) => (
        <div key={entry.id}>{entry.render(context)}</div>
      ))}
    </>
  );
}
