import { useCallback, useEffect, useMemo, useState, type ReactNode } from "react";
import { Search, ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight, ArrowUpDown, LayoutGrid, List, Columns3, Grid3X3, ZoomIn, ZoomOut, SlidersHorizontal, X } from "lucide-react";
import type { FindFilter } from "../api/types";
import { ExtensionSlot } from "../router/RouteRegistry";
import { SavedFilterMenu } from "./SavedFilterMenu";
import { FilterDialog, FilterButton, type CriterionDefinition } from "./FilterDialog";
import { useKeySequence } from "../hooks/useKeySequence";

export type DisplayMode = "grid" | "list" | "wall" | "tagger";

interface ListPageProps {
  title: string;
  pageKey?: string;
  filter: FindFilter;
  onFilterChange: (f: FindFilter) => void;
  totalCount: number;
  isLoading: boolean;
  children: ReactNode;
  sortOptions?: { value: string; label: string }[];
  displayMode?: DisplayMode;
  onDisplayModeChange?: (mode: DisplayMode) => void;
  availableDisplayModes?: DisplayMode[];
  selectedIds?: Set<number>;
  onSelectAll?: () => void;
  onSelectNone?: () => void;
  selectionActions?: ReactNode;
  metadataByline?: ReactNode;
  onNew?: () => void;
  renderOperations?: () => ReactNode;
  filterMode?: string;
  // Advanced filtering
  criteriaDefinitions?: CriterionDefinition[];
  objectFilter?: Record<string, unknown>;
  onObjectFilterChange?: (filter: Record<string, unknown>) => void;
  // Quick filter buttons (like original Stash's criterion shortcut row)
  quickFilterIds?: string[];
}

const PER_PAGE_OPTIONS = [20, 40, 60, 120, 250, 500, 1000];

export function ListPage({
  title,
  pageKey,
  filter,
  onFilterChange,
  totalCount,
  isLoading,
  children,
  sortOptions,
  displayMode,
  onDisplayModeChange,
  availableDisplayModes,
  selectedIds,
  onSelectAll,
  onSelectNone,
  selectionActions,
  metadataByline,
  onNew,
  renderOperations,
  filterMode,
  criteriaDefinitions,
  objectFilter,
  onObjectFilterChange,
  quickFilterIds,
}: ListPageProps) {
  const [searchText, setSearchText] = useState(filter.q ?? "");
  const [filterDialogOpen, setFilterDialogOpen] = useState(false);
  const [filterDialogPreselect, setFilterDialogPreselect] = useState<string | undefined>();
  const [zoomLevel, setZoomLevel] = useState(2); // 0-3 range: 0=smallest, 3=largest
  const perPage = filter.perPage ?? 25;
  const page = filter.page ?? 1;
  const totalPages = Math.max(1, Math.ceil(totalCount / perPage));
  const start = (page - 1) * perPage + 1;
  const end = Math.min(page * perPage, totalCount);
  const slotContext = { pageKey, title, filter, onFilterChange, totalCount, isLoading };
  const selecting = selectedIds && selectedIds.size > 0;

  useEffect(() => {
    setSearchText(filter.q ?? "");
  }, [filter.q]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    onFilterChange({ ...filter, q: searchText || undefined, page: 1 });
  };

  const goTo = useCallback(
    (p: number) => onFilterChange({ ...filter, page: Math.max(1, Math.min(totalPages, p)) }),
    [filter, onFilterChange, totalPages]
  );

  // List-page keyboard shortcuts
  const listBindings = useMemo(() => [
    // "/" focuses search
    { keys: "/", action: () => { document.querySelector<HTMLInputElement>("input[placeholder='Filter...']")?.focus(); } },
    // View switching
    ...(onDisplayModeChange && availableDisplayModes ? [
      ...(availableDisplayModes.includes("grid") ? [{ keys: "v g", action: () => onDisplayModeChange("grid") }] : []),
      ...(availableDisplayModes.includes("list") ? [{ keys: "v l", action: () => onDisplayModeChange("list") }] : []),
      ...(availableDisplayModes.includes("wall") ? [{ keys: "v w", action: () => onDisplayModeChange("wall") }] : []),
      ...(availableDisplayModes.includes("tagger") ? [{ keys: "v t", action: () => onDisplayModeChange("tagger") }] : []),
    ] : []),
    // Selection
    ...(onSelectAll ? [{ keys: "s a", action: onSelectAll }] : []),
    ...(onSelectNone ? [{ keys: "s n", action: onSelectNone }] : []),
    // Pagination
    { keys: "ArrowLeft", action: () => goTo(page - 1) },
    { keys: "ArrowRight", action: () => goTo(page + 1) },
    { keys: "Shift+ArrowLeft", action: () => goTo(page - 10) },
    { keys: "Shift+ArrowRight", action: () => goTo(page + 10) },
    { keys: "Ctrl+Home", action: () => goTo(1) },
    { keys: "Ctrl+End", action: () => goTo(totalPages) },
    // Filter dialog
    ...(criteriaDefinitions && onObjectFilterChange ? [{ keys: "f", action: () => setFilterDialogOpen(true) }] : []),
    // Zoom
    { keys: "+", action: () => setZoomLevel((v) => Math.min(3, v + 0.25)) },
    { keys: "-", action: () => setZoomLevel((v) => Math.max(0, v - 0.25)) },
  ], [onDisplayModeChange, availableDisplayModes, onSelectAll, onSelectNone, goTo, page, totalPages, criteriaDefinitions, onObjectFilterChange]);

  useKeySequence(listBindings);

  // Set page title like original Stash (e.g., "Scenes | Stash")
  useEffect(() => {
    document.title = `${title} | Stash`;
    return () => { document.title = "Stash"; };
  }, [title]);

  return (
    <div className="space-y-0">
      {/* Toolbar - matches Stash's FilteredListToolbar */}
      <div className="flex flex-wrap items-center gap-2 bg-plex-surface border-b border-plex-border px-3 py-1.5 sticky top-0 z-30">
        {/* Title + count + byline */}
        <div className="flex items-center gap-2 mr-2">
          <h1 className="text-sm font-semibold text-plex-text whitespace-nowrap">{title}</h1>
          <span className="text-xs text-plex-text-muted">
            {totalCount > 0 ? `${start}-${end} of ${totalCount.toLocaleString()}` : "0 items"}
          </span>
          {metadataByline}
        </div>

        {/* Search */}
        <form onSubmit={handleSearch} className="relative">
          <Search className="absolute left-2 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-plex-text-muted" />
          <input
            type="text"
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            placeholder="Filter..."
            className="bg-plex-input border border-plex-border rounded pl-7 pr-3 py-1 text-xs text-plex-text w-36 focus:outline-none focus:border-plex-accent placeholder:text-plex-text-muted"
          />
        </form>

        {/* Sort */}
        {sortOptions && (
          <select
            value={filter.sort ?? ""}
            onChange={(e) => onFilterChange({ ...filter, sort: e.target.value || undefined, page: 1 })}
            className="bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent"
          >
            {sortOptions.map((o) => (
              <option key={o.value} value={o.value}>{o.label}</option>
            ))}
          </select>
        )}

        {/* Direction toggle */}
        <button
          onClick={() => onFilterChange({ ...filter, direction: filter.direction === "desc" ? "asc" : "desc" })}
          className="p-1 rounded border border-plex-border bg-plex-input text-plex-text-secondary hover:text-plex-text hover:border-plex-accent"
          title={filter.direction === "desc" ? "Sort descending" : "Sort ascending"}
        >
          <ArrowUpDown className="w-3.5 h-3.5" />
        </button>

        {/* Saved filters */}
        {filterMode && (
          <SavedFilterMenu
            mode={filterMode}
            currentFilter={filter}
            currentObjectFilter={objectFilter}
            currentUIOptions={displayMode ? { displayMode } : undefined}
            onApplyFilter={onFilterChange}
            onApplyObjectFilter={onObjectFilterChange}
            onApplyUIOptions={(options) => {
              const mode = typeof options.displayMode === "string" ? options.displayMode : undefined;
              if (mode && onDisplayModeChange) onDisplayModeChange(mode as DisplayMode);
            }}
          />
        )}

        {/* Advanced filter button */}
        {criteriaDefinitions && onObjectFilterChange && (
          <FilterButton
            activeCount={Object.keys(objectFilter ?? {}).length}
            onClick={() => setFilterDialogOpen(true)}
          />
        )}

        {/* Display mode */}
        {onDisplayModeChange && availableDisplayModes && (
          <div className="flex items-center border border-plex-border rounded overflow-hidden">
            {availableDisplayModes.includes("grid") && (
              <button
                onClick={() => onDisplayModeChange("grid")}
                className={`p-1 ${displayMode === "grid" ? "bg-plex-card text-plex-accent" : "bg-plex-input text-plex-text-secondary hover:text-plex-text"}`}
                title="Grid"
              >
                <LayoutGrid className="w-3.5 h-3.5" />
              </button>
            )}
            {availableDisplayModes.includes("list") && (
              <button
                onClick={() => onDisplayModeChange("list")}
                className={`p-1 ${displayMode === "list" ? "bg-plex-card text-plex-accent" : "bg-plex-input text-plex-text-secondary hover:text-plex-text"}`}
                title="List"
              >
                <List className="w-3.5 h-3.5" />
              </button>
            )}
            {availableDisplayModes.includes("wall") && (
              <button
                onClick={() => onDisplayModeChange("wall")}
                className={`p-1 ${displayMode === "wall" ? "bg-plex-card text-plex-accent" : "bg-plex-input text-plex-text-secondary hover:text-plex-text"}`}
                title="Wall"
              >
                <Grid3X3 className="w-3.5 h-3.5" />
              </button>
            )}
            {availableDisplayModes.includes("tagger") && (
              <button
                onClick={() => onDisplayModeChange("tagger")}
                className={`p-1 ${displayMode === "tagger" ? "bg-plex-card text-plex-accent" : "bg-plex-input text-plex-text-secondary hover:text-plex-text"}`}
                title="Tagger"
              >
                <Columns3 className="w-3.5 h-3.5" />
              </button>
            )}
          </div>
        )}

        {/* Per page */}
        <select
          value={perPage}
          onChange={(e) => onFilterChange({ ...filter, perPage: Number(e.target.value), page: 1 })}
          className="bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent"
        >
          {PER_PAGE_OPTIONS.map((n) => (
            <option key={n} value={n}>{n}</option>
          ))}
        </select>

        {/* Zoom slider (like Stash's card size slider) */}
        {displayMode === "grid" && (
          <div className="flex items-center gap-1">
            <ZoomOut className="w-3 h-3 text-plex-text-muted" />
            <input
              type="range"
              min={0}
              max={3}
              step={0.25}
              value={zoomLevel}
              onChange={(e) => setZoomLevel(Number(e.target.value))}
              className="w-16 h-1 accent-plex-accent cursor-pointer"
              title={`Card size: ${Math.round(zoomLevel * 33 + 33)}%`}
            />
            <ZoomIn className="w-3 h-3 text-plex-text-muted" />
          </div>
        )}

        {/* Operations */}
        <div className="ml-auto flex items-center gap-2">
          {renderOperations?.()}
          <ExtensionSlot slot="list-page-toolbar-end" context={slotContext} />
          {pageKey && <ExtensionSlot slot={`${pageKey}-list-toolbar-end`} context={slotContext} />}
          {onNew && (
            <button
              onClick={onNew}
              className="px-3 py-1 rounded text-xs font-medium bg-blue-600 hover:bg-blue-500 text-white"
            >
              + New
            </button>
          )}
        </div>
      </div>

      {/* Active filter tags (criterion badges) */}
      {objectFilter && onObjectFilterChange && criteriaDefinitions && Object.keys(objectFilter).length > 0 && (
        <div className="flex flex-wrap items-center gap-1.5 bg-plex-surface/50 border-b border-plex-border px-3 py-1.5">
          {Object.entries(objectFilter).map(([key, value]) => {
            const def = criteriaDefinitions.find((d) => d.id === key);
            const label = def?.label ?? key;
            const displayValue = typeof value === "object" && value !== null
              ? Array.isArray(value) ? (value as { label?: string; name?: string; id: number }[]).map((v) => v.label ?? v.name ?? v.id).join(", ") : JSON.stringify(value)
              : String(value);
            return (
              <button
                key={key}
                onClick={() => {
                  const next = { ...objectFilter };
                  delete next[key];
                  onObjectFilterChange(next);
                  onFilterChange({ ...filter, page: 1 });
                }}
                className="group flex items-center gap-1 rounded-full bg-plex-card border border-plex-border px-2.5 py-0.5 text-xs text-plex-text hover:border-red-400 hover:text-red-300 transition-colors"
                title={`Remove filter: ${label}`}
              >
                <span className="text-plex-text-muted">{label}:</span>
                <span className="max-w-[200px] truncate">{displayValue}</span>
                <X className="w-3 h-3 opacity-50 group-hover:opacity-100" />
              </button>
            );
          })}
          <button
            onClick={() => { onObjectFilterChange({}); onFilterChange({ ...filter, page: 1 }); }}
            className="text-xs text-plex-text-muted hover:text-red-300"
          >
            Clear all
          </button>
        </div>
      )}

      {/* Selection bar */}
      {selecting && (
        <div className="flex items-center gap-3 bg-plex-card border-b border-plex-border px-3 py-1.5">
          <span className="text-xs text-plex-text-secondary">
            {selectedIds!.size} selected
          </span>
          <button onClick={onSelectAll} className="text-xs text-plex-accent hover:underline">Select all</button>
          <button onClick={onSelectNone} className="text-xs text-plex-text-secondary hover:text-plex-text">Deselect all</button>
          {selectionActions}
        </div>
      )}

      {/* Pagination top */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-1 bg-plex-surface border-b border-plex-border py-1">
          <PaginationControls page={page} totalPages={totalPages} goTo={goTo} />
        </div>
      )}

      {/* Content */}
      {isLoading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-plex-accent" />
        </div>
      ) : (
        <div className="pt-3" style={{ "--card-min-width": `${Math.round(120 + zoomLevel * 80)}px` } as React.CSSProperties}>
          {children}
        </div>
      )}

      {/* Pagination bottom */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-1 py-4">
          <PaginationControls page={page} totalPages={totalPages} goTo={goTo} />
        </div>
      )}

      {/* Filter Dialog */}
      {criteriaDefinitions && onObjectFilterChange && (
        <FilterDialog
          open={filterDialogOpen}
          onClose={() => { setFilterDialogOpen(false); setFilterDialogPreselect(undefined); }}
          criteria={criteriaDefinitions}
          activeFilter={objectFilter ?? {}}
          onApply={(f) => {
            onObjectFilterChange(f);
            onFilterChange({ ...filter, page: 1 });
          }}
          preselectCriterion={filterDialogPreselect}
        />
      )}
    </div>
  );
}

function PaginationControls({ page, totalPages, goTo }: { page: number; totalPages: number; goTo: (p: number) => void }) {
  return (
    <>
      <button onClick={() => goTo(1)} disabled={page <= 1} className="p-1 rounded hover:bg-plex-card disabled:opacity-30 disabled:cursor-not-allowed text-plex-text-secondary hover:text-plex-text">
        <ChevronsLeft className="w-3.5 h-3.5" />
      </button>
      <button onClick={() => goTo(page - 1)} disabled={page <= 1} className="p-1 rounded hover:bg-plex-card disabled:opacity-30 disabled:cursor-not-allowed text-plex-text-secondary hover:text-plex-text">
        <ChevronLeft className="w-3.5 h-3.5" />
      </button>
      {getPageNumbers(page, totalPages).map((p, i) =>
        p === -1 ? (
          <span key={`ellipsis-${i}`} className="px-1 text-plex-text-muted text-xs">…</span>
        ) : (
          <button
            key={p}
            onClick={() => goTo(p)}
            className={`min-w-[28px] h-7 rounded text-xs font-medium ${
              p === page ? "bg-plex-accent text-white" : "text-plex-text-secondary hover:bg-plex-card hover:text-plex-text"
            }`}
          >
            {p}
          </button>
        )
      )}
      <button onClick={() => goTo(page + 1)} disabled={page >= totalPages} className="p-1 rounded hover:bg-plex-card disabled:opacity-30 disabled:cursor-not-allowed text-plex-text-secondary hover:text-plex-text">
        <ChevronRight className="w-3.5 h-3.5" />
      </button>
      <button onClick={() => goTo(totalPages)} disabled={page >= totalPages} className="p-1 rounded hover:bg-plex-card disabled:opacity-30 disabled:cursor-not-allowed text-plex-text-secondary hover:text-plex-text">
        <ChevronsRight className="w-3.5 h-3.5" />
      </button>
    </>
  );
}

function getPageNumbers(current: number, total: number): number[] {
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);
  const pages: number[] = [1];
  if (current > 3) pages.push(-1);
  for (let i = Math.max(2, current - 1); i <= Math.min(total - 1, current + 1); i++) pages.push(i);
  if (current < total - 2) pages.push(-1);
  pages.push(total);
  return pages;
}
