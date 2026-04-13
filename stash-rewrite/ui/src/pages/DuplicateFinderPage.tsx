import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { scenes } from "../api/client";
import type { Scene } from "../api/types";
import { formatDuration, formatFileSize, getResolutionLabel } from "../components/shared";
import { Copy, Trash2, Loader2, Search, AlertTriangle, Check } from "lucide-react";

interface Props {
  onNavigate: (r: any) => void;
}

export function DuplicateFinderPage({ onNavigate }: Props) {
  const [distance, setDistance] = useState(0);
  const [groups, setGroups] = useState<Scene[][] | null>(null);
  const [selectedPerGroup, setSelectedPerGroup] = useState<Map<number, Set<number>>>(new Map());
  const queryClient = useQueryClient();

  const findMut = useMutation({
    mutationFn: () => scenes.findDuplicates(distance),
    onSuccess: (data) => {
      setGroups(data);
      setSelectedPerGroup(new Map());
    },
  });

  const deleteMut = useMutation({
    mutationFn: (ids: number[]) => scenes.bulkDelete(ids),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["scenes"] });
      // Re-run search to refresh
      findMut.mutate();
    },
  });

  const toggleSelected = (groupIdx: number, sceneId: number) => {
    setSelectedPerGroup((prev) => {
      const next = new Map(prev);
      const set = new Set(next.get(groupIdx) ?? []);
      if (set.has(sceneId)) set.delete(sceneId);
      else set.add(sceneId);
      next.set(groupIdx, set);
      return next;
    });
  };

  const keepSelected = (groupIdx: number) => {
    if (!groups) return;
    const group = groups[groupIdx];
    const kept = selectedPerGroup.get(groupIdx) ?? new Set();
    const toDelete = group.filter((s) => !kept.has(s.id)).map((s) => s.id);
    if (toDelete.length === 0) return;
    if (!confirm(`Delete ${toDelete.length} scene(s) and keep ${kept.size}?`)) return;
    deleteMut.mutate(toDelete);
  };

  return (
    <div>
      {/* Header */}
      <div className="flex items-center gap-3 mb-6">
        <Copy className="w-6 h-6 text-plex-accent" />
        <h1 className="text-xl font-semibold text-plex-text">Duplicate Finder</h1>
      </div>

      {/* Controls */}
      <div className="flex flex-wrap items-end gap-4 mb-6 p-4 bg-plex-card border border-plex-border rounded-lg">
        <div>
          <label className="block text-xs font-medium text-plex-text-secondary mb-1">
            Phash Distance (0 = exact hash match)
          </label>
          <div className="flex items-center gap-3">
            <input
              type="range"
              min={0}
              max={10}
              value={distance}
              onChange={(e) => setDistance(Number(e.target.value))}
              className="w-48 accent-plex-accent"
            />
            <span className="text-sm font-mono text-plex-text w-6 text-center">{distance}</span>
          </div>
        </div>
        <button
          onClick={() => findMut.mutate()}
          disabled={findMut.isPending}
          className="flex items-center gap-2 px-4 py-2 rounded text-sm font-medium bg-plex-accent hover:bg-plex-accent-hover text-white disabled:opacity-50"
        >
          {findMut.isPending ? (
            <Loader2 className="w-4 h-4 animate-spin" />
          ) : (
            <Search className="w-4 h-4" />
          )}
          {findMut.isPending ? "Searching..." : "Find Duplicates"}
        </button>
      </div>

      {/* Error */}
      {findMut.isError && (
        <div className="flex items-center gap-2 p-3 mb-4 bg-red-900/20 border border-red-800 rounded text-red-300 text-sm">
          <AlertTriangle className="w-4 h-4 shrink-0" />
          {(findMut.error as Error).message}
        </div>
      )}

      {/* Results summary */}
      {groups !== null && (
        <div className="mb-4 text-sm text-plex-text-secondary">
          Found <span className="font-semibold text-plex-text">{groups.length}</span> duplicate group{groups.length !== 1 ? "s" : ""}
          {groups.length > 0 && (
            <span>
              {" "}containing{" "}
              <span className="font-semibold text-plex-text">
                {groups.reduce((n, g) => n + g.length, 0)}
              </span>{" "}
              total scenes
            </span>
          )}
        </div>
      )}

      {/* No results */}
      {groups !== null && groups.length === 0 && (
        <div className="text-center py-16">
          <Check className="w-12 h-12 mx-auto mb-3 text-green-400" />
          <p className="text-plex-text-secondary">No duplicates found</p>
        </div>
      )}

      {/* Duplicate groups */}
      {groups && groups.length > 0 && (
        <div className="space-y-4">
          {groups.map((group, gi) => {
            const selected = selectedPerGroup.get(gi) ?? new Set();
            return (
              <div key={gi} className="border border-plex-border rounded-lg overflow-hidden">
                {/* Group header */}
                <div className="flex items-center justify-between px-4 py-2 bg-plex-card border-b border-plex-border">
                  <span className="text-sm font-medium text-plex-text">
                    Group {gi + 1} — {group.length} scenes
                  </span>
                  <div className="flex items-center gap-2">
                    <span className="text-xs text-plex-text-muted">
                      {selected.size > 0
                        ? `${selected.size} selected to keep`
                        : "Select scenes to keep"}
                    </span>
                    <button
                      onClick={() => keepSelected(gi)}
                      disabled={selected.size === 0 || selected.size === group.length || deleteMut.isPending}
                      className="flex items-center gap-1 px-2 py-1 text-xs rounded bg-red-600 hover:bg-red-500 text-white disabled:opacity-30 disabled:cursor-not-allowed"
                    >
                      <Trash2 className="w-3 h-3" />
                      Delete Others
                    </button>
                  </div>
                </div>

                {/* Scene cards */}
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-0">
                  {group.map((scene) => {
                    const file = scene.files[0];
                    const isSelected = selected.has(scene.id);
                    return (
                      <div
                        key={scene.id}
                        className={`relative border-r border-b border-plex-border last:border-r-0 ${
                          isSelected ? "bg-green-900/20 ring-inset ring-2 ring-green-500/50" : "bg-plex-bg"
                        }`}
                      >
                        {/* Selection overlay */}
                        <button
                          onClick={() => toggleSelected(gi, scene.id)}
                          className="absolute top-2 left-2 z-10"
                        >
                          <div
                            className={`w-5 h-5 rounded border-2 flex items-center justify-center ${
                              isSelected
                                ? "bg-green-500 border-green-500"
                                : "border-plex-text-muted bg-black/40"
                            }`}
                          >
                            {isSelected && <Check className="w-3 h-3 text-white" />}
                          </div>
                        </button>

                        {/* Thumbnail */}
                        <div
                          className="aspect-video bg-plex-card cursor-pointer"
                          onClick={() => onNavigate({ page: "scene", id: scene.id })}
                        >
                          <img
                            src={scenes.screenshotUrl(scene.id)}
                            alt={scene.title || ""}
                            className="w-full h-full object-cover"
                            loading="lazy"
                            onError={(e) => {
                              (e.target as HTMLImageElement).style.display = "none";
                            }}
                          />
                        </div>

                        {/* Details */}
                        <div className="p-2 space-y-1">
                          <p
                            className="text-xs font-medium text-plex-text truncate cursor-pointer hover:text-plex-accent"
                            onClick={() => onNavigate({ page: "scene", id: scene.id })}
                          >
                            {scene.title || file?.basename || `Scene #${scene.id}`}
                          </p>
                          {file && (
                            <div className="flex flex-wrap gap-x-3 gap-y-0.5 text-[10px] text-plex-text-muted">
                              <span>{file.width}×{file.height}</span>
                              <span>{getResolutionLabel(file.width, file.height)}</span>
                              <span>{formatDuration(file.duration)}</span>
                              <span>{formatFileSize(file.size)}</span>
                              <span>{file.videoCodec}</span>
                              <span>{Math.round(file.bitRate / 1000)} kbps</span>
                            </div>
                          )}
                          {file?.path && (
                            <p className="text-[9px] text-plex-text-muted truncate" title={file.path}>
                              {file.path}
                            </p>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
