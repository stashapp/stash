import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { scenes } from "../api/client";
import type { Scene, SceneMarkerSummary } from "../api/types";
import { formatDuration } from "../components/shared";
import { Bookmark, Search, ArrowUpDown } from "lucide-react";

interface MarkerWithScene extends SceneMarkerSummary {
  sceneId: number;
  sceneTitle: string;
}

const SORT_OPTIONS = [
  { value: "scene", label: "Scene Title" },
  { value: "time", label: "Timestamp" },
  { value: "tag", label: "Tag" },
  { value: "title", label: "Marker Title" },
] as const;

type SortKey = (typeof SORT_OPTIONS)[number]["value"];

interface Props {
  onNavigate: (r: any) => void;
}

export function SceneMarkersPage({ onNavigate }: Props) {
  const [search, setSearch] = useState("");
  const [sortBy, setSortBy] = useState<SortKey>("scene");
  const [sortDir, setSortDir] = useState<"asc" | "desc">("asc");
  const [tagFilter, setTagFilter] = useState<string>("");
  const [displayMode, setDisplayMode] = useState<"grid" | "wall">("grid");

  // Fetch all scenes (high per-page) and flatten markers
  const { data: scenesData, isLoading } = useQuery({
    queryKey: ["scenes-for-markers"],
    queryFn: () => scenes.find({ page: 1, perPage: 5000 }),
  });

  const markers = useMemo<MarkerWithScene[]>(() => {
    if (!scenesData?.items) return [];
    const result: MarkerWithScene[] = [];
    for (const scene of scenesData.items) {
      for (const marker of scene.markers) {
        result.push({
          ...marker,
          sceneId: scene.id,
          sceneTitle: scene.title || scene.files[0]?.basename || `Scene #${scene.id}`,
        });
      }
    }
    return result;
  }, [scenesData]);

  // Unique tag names for filter dropdown
  const tagNames = useMemo(() => {
    const names = new Set(markers.map((m) => m.primaryTagName));
    return Array.from(names).sort();
  }, [markers]);

  // Filter and sort
  const filtered = useMemo(() => {
    let list = markers;
    if (search) {
      const q = search.toLowerCase();
      list = list.filter(
        (m) =>
          m.title.toLowerCase().includes(q) ||
          m.sceneTitle.toLowerCase().includes(q) ||
          m.primaryTagName.toLowerCase().includes(q)
      );
    }
    if (tagFilter) {
      list = list.filter((m) => m.primaryTagName === tagFilter);
    }
    list = [...list].sort((a, b) => {
      let cmp = 0;
      switch (sortBy) {
        case "scene":
          cmp = a.sceneTitle.localeCompare(b.sceneTitle);
          break;
        case "time":
          cmp = a.seconds - b.seconds;
          break;
        case "tag":
          cmp = a.primaryTagName.localeCompare(b.primaryTagName);
          break;
        case "title":
          cmp = a.title.localeCompare(b.title);
          break;
      }
      return sortDir === "asc" ? cmp : -cmp;
    });
    return list;
  }, [markers, search, tagFilter, sortBy, sortDir]);

  const toggleSortDir = () => setSortDir((d) => (d === "asc" ? "desc" : "asc"));

  return (
    <div>
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          <h1 className="text-xl font-semibold text-plex-text">Scene Markers</h1>
          <span className="text-sm text-plex-text-muted">
            {filtered.length} marker{filtered.length !== 1 ? "s" : ""}
          </span>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setDisplayMode("grid")}
            className={`px-2 py-1 text-xs rounded ${displayMode === "grid" ? "bg-plex-accent text-white" : "bg-plex-input text-plex-text-secondary border border-plex-border"}`}
          >
            Grid
          </button>
          <button
            onClick={() => setDisplayMode("wall")}
            className={`px-2 py-1 text-xs rounded ${displayMode === "wall" ? "bg-plex-accent text-white" : "bg-plex-input text-plex-text-secondary border border-plex-border"}`}
          >
            Wall
          </button>
        </div>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-3 mb-4">
        <div className="relative flex-1 min-w-[200px] max-w-sm">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-plex-text-muted" />
          <input
            type="text"
            placeholder="Filter markers..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full pl-8 pr-3 py-1.5 text-sm bg-plex-input border border-plex-border rounded text-plex-text placeholder:text-plex-text-muted focus:outline-none focus:border-plex-accent"
          />
        </div>
        <select
          value={tagFilter}
          onChange={(e) => setTagFilter(e.target.value)}
          className="px-3 py-1.5 text-sm bg-plex-input border border-plex-border rounded text-plex-text focus:outline-none focus:border-plex-accent"
        >
          <option value="">All Tags</option>
          {tagNames.map((t) => (
            <option key={t} value={t}>
              {t}
            </option>
          ))}
        </select>
        <select
          value={sortBy}
          onChange={(e) => setSortBy(e.target.value as SortKey)}
          className="px-3 py-1.5 text-sm bg-plex-input border border-plex-border rounded text-plex-text focus:outline-none focus:border-plex-accent"
        >
          {SORT_OPTIONS.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
        <button
          onClick={toggleSortDir}
          className="flex items-center gap-1 px-2 py-1.5 text-sm rounded border border-plex-border bg-plex-input text-plex-text-secondary hover:text-plex-text"
          title={sortDir === "asc" ? "Ascending" : "Descending"}
        >
          <ArrowUpDown className="w-3.5 h-3.5" />
          {sortDir === "asc" ? "Asc" : "Desc"}
        </button>
      </div>

      {/* Loading */}
      {isLoading && (
        <div className="text-center py-20 text-plex-text-muted">Loading scenes...</div>
      )}

      {/* Empty state */}
      {!isLoading && filtered.length === 0 && (
        <div className="text-center py-20">
          <Bookmark className="w-16 h-16 mx-auto mb-4 text-plex-text-muted opacity-50" />
          <p className="text-plex-text-secondary text-lg">No markers found</p>
          <p className="text-plex-text-muted text-sm mt-1">
            Add markers to scenes to see them here
          </p>
        </div>
      )}

      {/* Grid mode */}
      {displayMode === "grid" && filtered.length > 0 && (
        <div className="grid gap-2 px-2" style={{ gridTemplateColumns: "repeat(auto-fill, minmax(var(--card-min-width, 200px), 1fr))" }}>
          {filtered.map((marker) => (
            <MarkerCard
              key={`${marker.sceneId}-${marker.id}`}
              marker={marker}
              onClick={() => onNavigate({ page: "scene", id: marker.sceneId })}
            />
          ))}
        </div>
      )}

      {/* Wall mode */}
      {displayMode === "wall" && filtered.length > 0 && (
        <div className="columns-2 sm:columns-3 md:columns-4 lg:columns-5 gap-1 px-2">
          {filtered.map((marker) => (
            <MarkerWallCard
              key={`${marker.sceneId}-${marker.id}`}
              marker={marker}
              onClick={() => onNavigate({ page: "scene", id: marker.sceneId })}
            />
          ))}
        </div>
      )}
    </div>
  );
}

/* ── Marker Card ── */

function MarkerCard({ marker, onClick }: { marker: MarkerWithScene; onClick: () => void }) {
  const screenshotUrl = scenes.screenshotUrl(marker.sceneId);

  return (
    <div
      onClick={onClick}
      className="cursor-pointer group rounded border border-plex-border bg-plex-card overflow-hidden hover:border-plex-accent transition-colors"
    >
      <div className="relative aspect-video bg-plex-card overflow-hidden">
        <img
          src={screenshotUrl}
          alt={marker.title}
          className="w-full h-full object-cover"
          loading="lazy"
          onError={(e) => {
            (e.target as HTMLImageElement).style.display = "none";
          }}
        />
        {/* Time badge */}
        <div className="absolute bottom-1 right-1">
          <span className="text-[10px] font-mono font-medium text-white bg-black/70 px-1.5 py-0.5 rounded">
            {formatDuration(marker.seconds)}
          </span>
        </div>
        {/* Tag badge */}
        <div className="absolute top-1 left-1">
          <span className="text-[10px] font-medium text-white bg-plex-accent/80 px-1.5 py-0.5 rounded">
            {marker.primaryTagName}
          </span>
        </div>
      </div>
      <div className="p-2">
        <p className="text-xs font-medium text-plex-text truncate">
          {marker.title || marker.primaryTagName}
        </p>
        <p className="text-[10px] text-plex-text-muted truncate mt-0.5">
          {marker.sceneTitle}
        </p>
      </div>
    </div>
  );
}

/* ── Marker Wall Card ── */

function MarkerWallCard({ marker, onClick }: { marker: MarkerWithScene; onClick: () => void }) {
  const screenshotUrl = scenes.screenshotUrl(marker.sceneId);

  return (
    <div
      onClick={onClick}
      className="cursor-pointer relative mb-1 break-inside-avoid group"
    >
      <img
        src={screenshotUrl}
        alt={marker.title}
        className="w-full rounded"
        loading="lazy"
        onError={(e) => {
          (e.target as HTMLImageElement).style.display = "none";
        }}
      />
      <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity rounded flex flex-col justify-end p-2">
        <span className="text-[10px] font-medium text-plex-accent mb-0.5">
          {marker.primaryTagName}
        </span>
        <span className="text-xs font-medium text-white truncate">
          {marker.title || marker.primaryTagName}
        </span>
        <span className="text-[10px] text-gray-300 truncate">
          {marker.sceneTitle} — {formatDuration(marker.seconds)}
        </span>
      </div>
    </div>
  );
}
