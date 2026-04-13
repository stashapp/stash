import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { savedFilters } from "../api/client";
import type { FindFilter } from "../api/types";
import { Bookmark, ChevronDown, Save, Trash2, Loader2 } from "lucide-react";

interface SavedFilterMenuProps {
  mode: string;
  currentFilter: FindFilter;
  currentObjectFilter?: Record<string, unknown>;
  currentUIOptions?: Record<string, unknown>;
  onApplyFilter: (filter: FindFilter) => void;
  onApplyObjectFilter?: (filter: Record<string, unknown>) => void;
  onApplyUIOptions?: (options: Record<string, unknown>) => void;
}

export function SavedFilterMenu({
  mode,
  currentFilter,
  currentObjectFilter,
  currentUIOptions,
  onApplyFilter,
  onApplyObjectFilter,
  onApplyUIOptions,
}: SavedFilterMenuProps) {
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [saveName, setSaveName] = useState("");
  const [showSave, setShowSave] = useState(false);

  const { data: filters } = useQuery({
    queryKey: ["saved-filters", mode],
    queryFn: () => savedFilters.list(mode),
  });

  const createMut = useMutation({
    mutationFn: () =>
      savedFilters.create({
        mode,
        name: saveName,
        findFilter: JSON.stringify(currentFilter),
        objectFilter: currentObjectFilter && Object.keys(currentObjectFilter).length > 0 ? JSON.stringify(currentObjectFilter) : undefined,
        uiOptions: currentUIOptions && Object.keys(currentUIOptions).length > 0 ? JSON.stringify(currentUIOptions) : undefined,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["saved-filters", mode] });
      setSaveName("");
      setShowSave(false);
    },
  });

  const deleteMut = useMutation({
    mutationFn: (id: number) => savedFilters.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["saved-filters", mode] });
    },
  });

  const applyFilter = (findFilterJson: string | undefined, objectFilterJson?: string, uiOptionsJson?: string) => {
    if (!findFilterJson) return;
    try {
      const parsed = JSON.parse(findFilterJson) as FindFilter;
      onApplyFilter(parsed);
    } catch {
      // ignore invalid JSON
    }

    if (onApplyObjectFilter) {
      try {
        onApplyObjectFilter(objectFilterJson ? JSON.parse(objectFilterJson) as Record<string, unknown> : {});
      } catch {
        onApplyObjectFilter({});
      }
    }

    if (onApplyUIOptions && uiOptionsJson) {
      try {
        onApplyUIOptions(JSON.parse(uiOptionsJson) as Record<string, unknown>);
      } catch {
        // ignore invalid JSON
      }
    }

    setOpen(false);
  };

  return (
    <div className="relative">
      <button
        onClick={() => setOpen(!open)}
        className="flex items-center gap-1 px-2 py-1 rounded text-xs border border-plex-border bg-plex-input text-plex-text-secondary hover:text-plex-text"
        title="Saved filters"
      >
        <Bookmark className="w-3.5 h-3.5" />
        <ChevronDown className="w-3 h-3" />
      </button>

      {open && (
        <div className="absolute top-full right-0 mt-1 w-56 bg-plex-card border border-plex-border rounded-lg shadow-lg z-50">
          <div className="p-2 border-b border-plex-border">
            <p className="text-[10px] text-plex-text-muted uppercase tracking-wider font-medium">
              Saved Filters
            </p>
          </div>

          {/* Existing filters */}
          <div className="max-h-48 overflow-y-auto">
            {(!filters || filters.length === 0) && (
              <p className="px-3 py-2 text-xs text-plex-text-muted">No saved filters</p>
            )}
            {filters?.map((f) => (
              <div
                key={f.id}
                className="flex items-center justify-between px-3 py-1.5 hover:bg-plex-surface cursor-pointer group"
              >
                <button
                  onClick={() => applyFilter(f.findFilter, f.objectFilter, f.uiOptions)}
                  className="text-xs text-plex-text hover:text-plex-accent truncate flex-1 text-left"
                >
                  {f.name}
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    deleteMut.mutate(f.id);
                  }}
                  className="p-0.5 text-plex-text-muted hover:text-red-400 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  <Trash2 className="w-3 h-3" />
                </button>
              </div>
            ))}
          </div>

          {/* Save current */}
          <div className="border-t border-plex-border p-2">
            {showSave ? (
              <div className="flex gap-1">
                <input
                  type="text"
                  value={saveName}
                  onChange={(e) => setSaveName(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && saveName && createMut.mutate()}
                  placeholder="Filter name..."
                  className="flex-1 bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent placeholder:text-plex-text-muted"
                  autoFocus
                />
                <button
                  onClick={() => saveName && createMut.mutate()}
                  disabled={!saveName || createMut.isPending}
                  className="px-2 py-1 rounded text-xs bg-plex-accent text-white hover:bg-plex-accent-hover disabled:opacity-60"
                >
                  {createMut.isPending ? <Loader2 className="w-3 h-3 animate-spin" /> : <Save className="w-3 h-3" />}
                </button>
              </div>
            ) : (
              <button
                onClick={() => setShowSave(true)}
                className="flex items-center gap-1.5 text-xs text-plex-text-secondary hover:text-plex-text w-full"
              >
                <Save className="w-3 h-3" />
                Save current filter
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
