import { useState, useRef, useEffect, useCallback, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { scenes, performers, studios, tags, galleries, groups, savedFilters } from "../api/client";
import type { Scene, Performer, Studio, Tag, Gallery, Group, SavedFilter } from "../api/types";
import { formatDuration, formatFileSize, getResolutionLabel, RatingBadge } from "../components/shared";
import { RatingBanner } from "../components/Rating";
import { ChevronLeft, ChevronRight, Settings2, Plus, Trash2, Film, User, Building2, Tag as TagIcon, Images, Clapperboard, GripVertical } from "lucide-react";

// ─── Types ───────────────────────────────────────────────────────────────────

type FilterMode = "scenes" | "performers" | "studios" | "tags" | "galleries" | "groups";

interface CustomFilter {
  type: "custom";
  mode: FilterMode;
  sortBy: string;
  direction: "asc" | "desc";
  header: string;
}

interface SavedFilterRow {
  type: "saved";
  savedFilterId: number;
}

type FrontPageContent = CustomFilter | SavedFilterRow;

// ─── Default content (matches original Stash defaults) ───────────────────────

const DEFAULT_CONTENT: FrontPageContent[] = [
  { type: "custom", mode: "scenes", sortBy: "date", direction: "desc", header: "Recently Released Scenes" },
  { type: "custom", mode: "studios", sortBy: "created_at", direction: "desc", header: "Recently Added Studios" },
  { type: "custom", mode: "groups", sortBy: "date", direction: "desc", header: "Recently Released Groups" },
  { type: "custom", mode: "performers", sortBy: "created_at", direction: "desc", header: "Recently Added Performers" },
  { type: "custom", mode: "galleries", sortBy: "date", direction: "desc", header: "Recently Released Galleries" },
];

// ─── Premade filter options (for adding new rows) ────────────────────────────

const PREMADE_FILTERS: CustomFilter[] = [
  { type: "custom", mode: "scenes", sortBy: "date", direction: "desc", header: "Recently Released Scenes" },
  { type: "custom", mode: "scenes", sortBy: "created_at", direction: "desc", header: "Recently Added Scenes" },
  { type: "custom", mode: "galleries", sortBy: "date", direction: "desc", header: "Recently Released Galleries" },
  { type: "custom", mode: "galleries", sortBy: "created_at", direction: "desc", header: "Recently Added Galleries" },
  { type: "custom", mode: "groups", sortBy: "date", direction: "desc", header: "Recently Released Groups" },
  { type: "custom", mode: "groups", sortBy: "created_at", direction: "desc", header: "Recently Added Groups" },
  { type: "custom", mode: "studios", sortBy: "created_at", direction: "desc", header: "Recently Added Studios" },
  { type: "custom", mode: "performers", sortBy: "created_at", direction: "desc", header: "Recently Added Performers" },
];

const STORAGE_KEY = "stash-front-page-content";

function loadContent(): FrontPageContent[] {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored) return JSON.parse(stored);
  } catch { /* ignore */ }
  return DEFAULT_CONTENT;
}

function saveContent(content: FrontPageContent[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(content));
}

// ─── Home Page Component ─────────────────────────────────────────────────────

interface Props {
  onNavigate: (r: any) => void;
}

export function HomePage({ onNavigate }: Props) {
  const [content, setContent] = useState<FrontPageContent[]>(loadContent);
  const [isEditing, setIsEditing] = useState(false);

  const updateContent = useCallback((newContent: FrontPageContent[]) => {
    setContent(newContent);
    saveContent(newContent);
  }, []);

  if (isEditing) {
    return (
      <FrontPageEditor
        content={content}
        onSave={(c) => { updateContent(c); setIsEditing(false); }}
        onCancel={() => setIsEditing(false)}
      />
    );
  }

  return (
    <div className="space-y-6">
      {content.map((item, i) => (
        <RecommendationRow key={i} content={item} onNavigate={onNavigate} />
      ))}
      <div className="flex justify-end pb-4">
        <button
          onClick={() => setIsEditing(true)}
          className="px-4 py-2 text-sm bg-plex-card border border-plex-border rounded hover:bg-plex-surface text-plex-text"
        >
          Customize
        </button>
      </div>
    </div>
  );
}

// ─── Recommendation Row (dispatcher) ────────────────────────────────────────

function RecommendationRow({ content, onNavigate }: { content: FrontPageContent; onNavigate: (r: any) => void }) {
  if (content.type === "saved") {
    return <SavedFilterRecommendationRow savedFilterId={content.savedFilterId} onNavigate={onNavigate} />;
  }
  return <CustomFilterRecommendationRow filter={content} onNavigate={onNavigate} />;
}

// ─── Custom Filter Row ──────────────────────────────────────────────────────

function CustomFilterRecommendationRow({ filter, onNavigate }: { filter: CustomFilter; onNavigate: (r: any) => void }) {
  const fetchFn = useMemo((): (() => Promise<any>) => {
    const params = { perPage: 25, sort: filter.sortBy, direction: filter.direction };
    switch (filter.mode) {
      case "scenes": return () => scenes.find(params);
      case "performers": return () => performers.find(params);
      case "studios": return () => studios.find(params);
      case "tags": return () => tags.find(params);
      case "galleries": return () => galleries.find(params);
      case "groups": return () => groups.find(params);
    }
  }, [filter]);

  const { data, isLoading } = useQuery<any>({
    queryKey: ["front-page", filter.mode, filter.sortBy, filter.direction],
    queryFn: fetchFn,
  });

  const items = data?.items ?? [];
  if (!isLoading && items.length === 0) return null;

  return (
    <RecommendationRowShell
      header={filter.header}
      viewAllPage={filter.mode}
      onNavigate={onNavigate}
      loading={isLoading}
      count={items.length}
    >
      {items.map((item: any) => (
        <EntityCard key={item.id} item={item} mode={filter.mode} onNavigate={onNavigate} />
      ))}
    </RecommendationRowShell>
  );
}

// ─── Saved Filter Row ───────────────────────────────────────────────────────

function SavedFilterRecommendationRow({ savedFilterId, onNavigate }: { savedFilterId: number; onNavigate: (r: any) => void }) {
  const { data: filter } = useQuery({
    queryKey: ["saved-filter", savedFilterId],
    queryFn: () => savedFilters.get(savedFilterId),
  });

  const mode = filter?.mode as FilterMode | undefined;
  const parsedFilter = useMemo(() => {
    if (!filter?.findFilter) return {};
    try { return JSON.parse(filter.findFilter); } catch { return {}; }
  }, [filter]);

  const parsedObjectFilter = useMemo(() => {
    if (!filter?.objectFilter) return undefined;
    try { return JSON.parse(filter.objectFilter); } catch { return undefined; }
  }, [filter]);

  const fetchFn = useMemo((): (() => Promise<any>) => {
    if (!mode) return () => Promise.resolve({ items: [], totalCount: 0 });
    const findFilter = { perPage: 25, sort: parsedFilter.sort, direction: parsedFilter.direction };
    const fetchMap: Record<string, () => Promise<any>> = {
      scenes: parsedObjectFilter ? () => scenes.findFiltered({ findFilter, objectFilter: parsedObjectFilter }) : () => scenes.find(findFilter),
      performers: parsedObjectFilter ? () => performers.findFiltered({ findFilter, objectFilter: parsedObjectFilter }) : () => performers.find(findFilter),
      studios: parsedObjectFilter ? () => studios.findFiltered({ findFilter, objectFilter: parsedObjectFilter }) : () => studios.find(findFilter),
      tags: () => tags.find(findFilter),
      galleries: parsedObjectFilter ? () => galleries.findFiltered({ findFilter, objectFilter: parsedObjectFilter }) : () => galleries.find(findFilter),
      groups: parsedObjectFilter ? () => groups.findFiltered({ findFilter, objectFilter: parsedObjectFilter }) : () => groups.find(findFilter),
    };
    return fetchMap[mode] ?? (() => Promise.resolve({ items: [], totalCount: 0 }));
  }, [mode, parsedFilter, parsedObjectFilter]);

  const { data, isLoading } = useQuery<any>({
    queryKey: ["front-page-saved", savedFilterId, mode, parsedFilter, parsedObjectFilter],
    queryFn: fetchFn,
    enabled: !!mode,
  });

  const items = (data as any)?.items ?? [];
  if (!filter || (!isLoading && items.length === 0)) return null;

  return (
    <RecommendationRowShell
      header={filter.name}
      viewAllPage={mode ?? "scenes"}
      onNavigate={onNavigate}
      loading={isLoading}
      count={items.length}
    >
      {items.map((item: any) => (
        <EntityCard key={item.id} item={item} mode={mode!} onNavigate={onNavigate} />
      ))}
    </RecommendationRowShell>
  );
}

// ─── Recommendation Row Shell (horizontal carousel) ─────────────────────────

function RecommendationRowShell({
  header,
  viewAllPage,
  onNavigate,
  loading,
  count,
  children,
}: {
  header: string;
  viewAllPage: string;
  onNavigate: (r: any) => void;
  loading: boolean;
  count: number;
  children: React.ReactNode;
}) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const [canScrollLeft, setCanScrollLeft] = useState(false);
  const [canScrollRight, setCanScrollRight] = useState(false);
  const [currentPage, setCurrentPage] = useState(0);
  const [totalPages, setTotalPages] = useState(1);

  const updateScrollState = useCallback(() => {
    const el = scrollRef.current;
    if (!el) return;
    setCanScrollLeft(el.scrollLeft > 5);
    setCanScrollRight(el.scrollLeft + el.clientWidth < el.scrollWidth - 5);
    // Calculate pages
    if (el.clientWidth > 0) {
      const pages = Math.ceil(el.scrollWidth / el.clientWidth);
      setTotalPages(pages);
      setCurrentPage(Math.round(el.scrollLeft / el.clientWidth));
    }
  }, []);

  useEffect(() => {
    updateScrollState();
    const el = scrollRef.current;
    if (el) {
      el.addEventListener("scroll", updateScrollState);
      const resizeObserver = new ResizeObserver(updateScrollState);
      resizeObserver.observe(el);
      return () => { el.removeEventListener("scroll", updateScrollState); resizeObserver.disconnect(); };
    }
  }, [updateScrollState, count]);

  const scroll = (dir: "left" | "right") => {
    const el = scrollRef.current;
    if (!el) return;
    const scrollAmount = el.clientWidth * 0.85;
    el.scrollBy({ left: dir === "left" ? -scrollAmount : scrollAmount, behavior: "smooth" });
  };

  return (
    <div className="recommendation-row">
      {/* Header */}
      <div className="flex items-center justify-between mb-2 px-1">
        <h2 className="text-base font-semibold text-plex-text">{header}</h2>
        <button
          onClick={() => onNavigate({ page: viewAllPage })}
          className="text-xs text-plex-text-muted hover:text-plex-accent"
        >
          View All
        </button>
      </div>

      {/* Scrollable cards */}
      <div className="relative group">
        {/* Left arrow */}
        {canScrollLeft && (
          <button
            onClick={() => scroll("left")}
            className="absolute left-0 top-0 bottom-0 z-20 w-8 flex items-center justify-center bg-gradient-to-r from-plex-bg/90 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <ChevronLeft className="w-6 h-6 text-white" />
          </button>
        )}

        <div
          ref={scrollRef}
          className="flex gap-2 overflow-x-auto scrollbar-hide scroll-smooth px-1"
          style={{ scrollSnapType: "x mandatory" }}
        >
          {loading
            ? Array.from({ length: 6 }).map((_, i) => (
                <div key={i} className="flex-shrink-0 w-[200px] aspect-video bg-plex-card rounded animate-pulse" />
              ))
            : children}
        </div>

        {/* Right arrow */}
        {canScrollRight && (
          <button
            onClick={() => scroll("right")}
            className="absolute right-0 top-0 bottom-0 z-20 w-8 flex items-center justify-center bg-gradient-to-l from-plex-bg/90 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <ChevronRight className="w-6 h-6 text-white" />
          </button>
        )}
      </div>

      {/* Page dots */}
      {totalPages > 1 && (
        <div className="flex justify-center gap-1.5 mt-2">
          {Array.from({ length: totalPages }).map((_, i) => (
            <button
              key={i}
              onClick={() => {
                const el = scrollRef.current;
                if (el) el.scrollTo({ left: i * el.clientWidth, behavior: "smooth" });
              }}
              className={`w-6 h-1 rounded-full transition-colors ${i === currentPage ? "bg-plex-text" : "bg-plex-text-muted/40"}`}
            />
          ))}
        </div>
      )}
    </div>
  );
}

// ─── Entity Card (renders appropriate card based on mode) ───────────────────

function EntityCard({ item, mode, onNavigate }: { item: any; mode: FilterMode; onNavigate: (r: any) => void }) {
  switch (mode) {
    case "scenes": return <SceneRecommendationCard scene={item} onNavigate={onNavigate} />;
    case "performers": return <PerformerRecommendationCard performer={item} onNavigate={onNavigate} />;
    case "studios": return <StudioRecommendationCard studio={item} onNavigate={onNavigate} />;
    case "tags": return <TagRecommendationCard tag={item} onNavigate={onNavigate} />;
    case "galleries": return <GalleryRecommendationCard gallery={item} onNavigate={onNavigate} />;
    case "groups": return <GroupRecommendationCard group={item} onNavigate={onNavigate} />;
    default: return null;
  }
}

// ─── Scene Card ─────────────────────────────────────────────────────────────

function SceneRecommendationCard({ scene, onNavigate }: { scene: Scene; onNavigate: (r: any) => void }) {
  const file = scene.files[0];
  const duration = file?.duration ?? 0;
  const resLabel = file ? getResolutionLabel(file.width, file.height) : null;
  const screenshotUrl = scenes.screenshotUrl(scene.id);

  return (
    <div
      onClick={() => onNavigate({ page: "scene", id: scene.id })}
      className="flex-shrink-0 w-[200px] cursor-pointer group rounded overflow-hidden bg-plex-card border border-plex-border hover:border-plex-accent/50 transition-colors"
      style={{ scrollSnapAlign: "start" }}
    >
      <div className="relative aspect-video bg-black">
        <img src={screenshotUrl} alt={scene.title || ""} className="w-full h-full object-cover" loading="lazy" />
        {/* Resolution + duration overlay */}
        <div className="absolute bottom-0 right-0 flex items-center gap-0.5 p-1 text-[10px] text-white">
          {resLabel && <span className="bg-black/70 px-1 py-0.5 rounded font-bold">{resLabel}</span>}
          {duration > 0 && <span className="bg-black/70 px-1 py-0.5 rounded">{formatDuration(duration)}</span>}
        </div>
        <RatingBanner rating={scene.rating} />
      </div>
      <div className="px-2 py-1.5">
        <p className="text-[11px] font-medium text-plex-text truncate group-hover:text-plex-accent">
          {scene.title || file?.basename || "Untitled"}
        </p>
        {scene.date && <p className="text-[10px] text-plex-text-muted">{scene.date}</p>}
      </div>
      {/* Bottom stats */}
      <div className="flex items-center gap-2 px-2 pb-1.5 text-[10px] text-plex-text-muted">
        {scene.tags.length > 0 && (
          <span className="flex items-center gap-0.5"><TagIcon className="w-2.5 h-2.5" />{scene.tags.length}</span>
        )}
        {scene.performers.length > 0 && (
          <span className="flex items-center gap-0.5"><User className="w-2.5 h-2.5" />{scene.performers.length}</span>
        )}
      </div>
    </div>
  );
}

// ─── Performer Card ─────────────────────────────────────────────────────────

function PerformerRecommendationCard({ performer, onNavigate }: { performer: Performer; onNavigate: (r: any) => void }) {
  return (
    <div
      onClick={() => onNavigate({ page: "performer", id: performer.id })}
      className="flex-shrink-0 w-[160px] cursor-pointer group rounded overflow-hidden bg-plex-card border border-plex-border hover:border-plex-accent/50 transition-colors"
      style={{ scrollSnapAlign: "start" }}
    >
      <div className="relative aspect-[2/3] bg-plex-surface">
        {performer.imagePath ? (
          <img src={`/api/performers/${performer.id}/image`} alt={performer.name} className="w-full h-full object-cover" loading="lazy" />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <User className="w-10 h-10 text-plex-text-muted" />
          </div>
        )}
        <RatingBanner rating={performer.rating} />
      </div>
      <div className="px-2 py-1.5">
        <p className="text-[11px] font-medium text-plex-text truncate group-hover:text-plex-accent">{performer.name}</p>
        {performer.disambiguation && <p className="text-[10px] text-plex-text-muted truncate">{performer.disambiguation}</p>}
      </div>
      <div className="flex items-center gap-2 px-2 pb-1.5 text-[10px] text-plex-text-muted">
        {performer.sceneCount > 0 && <span>{performer.sceneCount} scenes</span>}
      </div>
    </div>
  );
}

// ─── Studio Card ────────────────────────────────────────────────────────────

function StudioRecommendationCard({ studio, onNavigate }: { studio: Studio; onNavigate: (r: any) => void }) {
  return (
    <div
      onClick={() => onNavigate({ page: "studio", id: studio.id })}
      className="flex-shrink-0 w-[200px] cursor-pointer group rounded overflow-hidden bg-plex-card border border-plex-border hover:border-plex-accent/50 transition-colors"
      style={{ scrollSnapAlign: "start" }}
    >
      <div className="relative aspect-video bg-plex-surface flex items-center justify-center p-4">
        {studio.imagePath ? (
          <img src={`/api/studios/${studio.id}/image`} alt={studio.name} className="max-w-full max-h-full object-contain" loading="lazy" />
        ) : (
          <Building2 className="w-10 h-10 text-plex-text-muted" />
        )}
        <RatingBanner rating={studio.rating} />
      </div>
      <div className="px-2 py-1.5">
        <p className="text-[11px] font-medium text-plex-text truncate group-hover:text-plex-accent">{studio.name}</p>
      </div>
      <div className="flex items-center gap-2 px-2 pb-1.5 text-[10px] text-plex-text-muted">
        {studio.sceneCount > 0 && <span>{studio.sceneCount} scenes</span>}
      </div>
    </div>
  );
}

// ─── Tag Card ───────────────────────────────────────────────────────────────

function TagRecommendationCard({ tag, onNavigate }: { tag: Tag; onNavigate: (r: any) => void }) {
  return (
    <div
      onClick={() => onNavigate({ page: "tag", id: tag.id })}
      className="flex-shrink-0 w-[160px] cursor-pointer group rounded overflow-hidden bg-plex-card border border-plex-border hover:border-plex-accent/50 transition-colors"
      style={{ scrollSnapAlign: "start" }}
    >
      <div className="relative aspect-video bg-plex-surface flex items-center justify-center">
        {tag.imagePath ? (
          <img src={`/api/tags/${tag.id}/image`} alt={tag.name} className="max-w-full max-h-full object-contain" loading="lazy" />
        ) : (
          <TagIcon className="w-8 h-8 text-plex-text-muted" />
        )}
      </div>
      <div className="px-2 py-1.5">
        <p className="text-[11px] font-medium text-plex-text truncate group-hover:text-plex-accent">{tag.name}</p>
      </div>
      {tag.sceneCount !== undefined && tag.sceneCount > 0 && (
        <div className="px-2 pb-1.5 text-[10px] text-plex-text-muted">
          {tag.sceneCount} scenes
        </div>
      )}
    </div>
  );
}

// ─── Gallery Card ───────────────────────────────────────────────────────────

function GalleryRecommendationCard({ gallery, onNavigate }: { gallery: Gallery; onNavigate: (r: any) => void }) {
  return (
    <div
      onClick={() => onNavigate({ page: "gallery", id: gallery.id })}
      className="flex-shrink-0 w-[200px] cursor-pointer group rounded overflow-hidden bg-plex-card border border-plex-border hover:border-plex-accent/50 transition-colors"
      style={{ scrollSnapAlign: "start" }}
    >
      <div className="relative aspect-video bg-plex-surface flex items-center justify-center">
        {gallery.coverPath ? (
          <img src={`/api/galleries/${gallery.id}/cover`} alt={gallery.title || ""} className="w-full h-full object-cover" loading="lazy" />
        ) : (
          <Images className="w-8 h-8 text-plex-text-muted" />
        )}
        <RatingBanner rating={gallery.rating} />
      </div>
      <div className="px-2 py-1.5">
        <p className="text-[11px] font-medium text-plex-text truncate group-hover:text-plex-accent">{gallery.title || "Untitled"}</p>
        {gallery.date && <p className="text-[10px] text-plex-text-muted">{gallery.date}</p>}
      </div>
      <div className="flex items-center gap-2 px-2 pb-1.5 text-[10px] text-plex-text-muted">
        {gallery.imageCount > 0 && <span>{gallery.imageCount} images</span>}
      </div>
    </div>
  );
}

// ─── Group Card ─────────────────────────────────────────────────────────────

function GroupRecommendationCard({ group, onNavigate }: { group: Group; onNavigate: (r: any) => void }) {
  return (
    <div
      onClick={() => onNavigate({ page: "group", id: group.id })}
      className="flex-shrink-0 w-[160px] cursor-pointer group rounded overflow-hidden bg-plex-card border border-plex-border hover:border-plex-accent/50 transition-colors"
      style={{ scrollSnapAlign: "start" }}
    >
      <div className="relative aspect-[2/3] bg-plex-surface flex items-center justify-center">
        {group.frontImagePath ? (
          <img src={`/api/groups/${group.id}/frontimage`} alt={group.name} className="w-full h-full object-cover" loading="lazy" />
        ) : (
          <Clapperboard className="w-8 h-8 text-plex-text-muted" />
        )}
        <RatingBanner rating={group.rating} />
      </div>
      <div className="px-2 py-1.5">
        <p className="text-[11px] font-medium text-plex-text truncate group-hover:text-plex-accent">{group.name}</p>
        {group.date && <p className="text-[10px] text-plex-text-muted">{group.date}</p>}
      </div>
      <div className="flex items-center gap-2 px-2 pb-1.5 text-[10px] text-plex-text-muted">
        {group.sceneCount > 0 && <span>{group.sceneCount} scenes</span>}
      </div>
    </div>
  );
}

// ─── Front Page Editor ──────────────────────────────────────────────────────

function FrontPageEditor({
  content,
  onSave,
  onCancel,
}: {
  content: FrontPageContent[];
  onSave: (content: FrontPageContent[]) => void;
  onCancel: () => void;
}) {
  const [items, setItems] = useState<FrontPageContent[]>([...content]);
  const [showAddModal, setShowAddModal] = useState(false);

  const { data: allSavedFilters } = useQuery({
    queryKey: ["saved-filters-all"],
    queryFn: () => savedFilters.list(),
  });

  const moveItem = (fromIndex: number, toIndex: number) => {
    if (toIndex < 0 || toIndex >= items.length) return;
    const newItems = [...items];
    const [moved] = newItems.splice(fromIndex, 1);
    newItems.splice(toIndex, 0, moved);
    setItems(newItems);
  };

  const removeItem = (index: number) => {
    setItems(items.filter((_, i) => i !== index));
  };

  const addItem = (item: FrontPageContent) => {
    setItems([...items, item]);
    setShowAddModal(false);
  };

  return (
    <div className="max-w-3xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-xl font-semibold text-plex-text">Customize Front Page</h1>
        <div className="flex gap-2">
          <button onClick={onCancel} className="px-4 py-2 text-sm text-plex-text-muted hover:text-plex-text">
            Cancel
          </button>
          <button
            onClick={() => onSave(items)}
            className="px-4 py-2 text-sm bg-plex-accent text-white rounded hover:bg-plex-accent-hover"
          >
            Save
          </button>
        </div>
      </div>

      <div className="space-y-2">
        {items.map((item, i) => (
          <div key={i} className="flex items-center gap-2 p-3 bg-plex-card border border-plex-border rounded">
            <div className="flex flex-col gap-0.5">
              <button
                onClick={() => moveItem(i, i - 1)}
                disabled={i === 0}
                className="text-plex-text-muted hover:text-plex-text disabled:opacity-30"
              >
                <ChevronLeft className="w-4 h-4 rotate-90" />
              </button>
              <button
                onClick={() => moveItem(i, i + 1)}
                disabled={i === items.length - 1}
                className="text-plex-text-muted hover:text-plex-text disabled:opacity-30"
              >
                <ChevronRight className="w-4 h-4 rotate-90" />
              </button>
            </div>
            <GripVertical className="w-4 h-4 text-plex-text-muted" />
            <div className="flex-1">
              <p className="text-sm text-plex-text">
                {item.type === "custom" ? item.header : `Saved Filter #${item.savedFilterId}`}
              </p>
              <p className="text-xs text-plex-text-muted">
                {item.type === "custom" ? `${item.mode} • ${item.sortBy} • ${item.direction}` : "Saved filter"}
              </p>
            </div>
            <button onClick={() => removeItem(i)} className="text-red-400 hover:text-red-300 p-1">
              <Trash2 className="w-4 h-4" />
            </button>
          </div>
        ))}
      </div>

      <button
        onClick={() => setShowAddModal(true)}
        className="mt-4 flex items-center gap-2 px-4 py-2 text-sm text-plex-accent hover:text-plex-accent-hover border border-plex-border rounded hover:border-plex-accent/50"
      >
        <Plus className="w-4 h-4" />
        Add Row
      </button>

      {/* Add Row Modal */}
      {showAddModal && (
        <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50" onClick={() => setShowAddModal(false)}>
          <div className="bg-plex-surface border border-plex-border rounded-lg p-6 max-w-md w-full mx-4 max-h-[80vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <h3 className="text-lg font-semibold text-plex-text mb-4">Add Content Row</h3>

            <h4 className="text-sm font-medium text-plex-text-muted mb-2">Premade Filters</h4>
            <div className="space-y-1 mb-4">
              {PREMADE_FILTERS.map((f, i) => (
                <button
                  key={i}
                  onClick={() => addItem(f)}
                  className="block w-full text-left px-3 py-2 text-sm text-plex-text hover:bg-plex-card rounded"
                >
                  {f.header}
                </button>
              ))}
            </div>

            {allSavedFilters && allSavedFilters.length > 0 && (
              <>
                <h4 className="text-sm font-medium text-plex-text-muted mb-2">Saved Filters</h4>
                <div className="space-y-1">
                  {allSavedFilters.map((sf) => (
                    <button
                      key={sf.id}
                      onClick={() => addItem({ type: "saved", savedFilterId: sf.id })}
                      className="block w-full text-left px-3 py-2 text-sm text-plex-text hover:bg-plex-card rounded"
                    >
                      <span className="text-plex-text-muted text-xs mr-2">{sf.mode}:</span>
                      {sf.name}
                    </button>
                  ))}
                </div>
              </>
            )}

            <div className="flex justify-end mt-4">
              <button onClick={() => setShowAddModal(false)} className="px-4 py-2 text-sm text-plex-text-muted hover:text-plex-text">
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
