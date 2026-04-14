import { useMemo, useState, useCallback, useRef, useEffect } from "react";
import { createPortal } from "react-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { scenes, tags, performers, studios, galleries, entityImages } from "../api/client";
import type { FindFilter, Scene, SceneCreate, SceneFilterCriteria } from "../api/types";
import { ListPage, type DisplayMode } from "../components/ListPage";
import { RatingBanner, RatingField } from "../components/Rating";
import { SceneTagger } from "../components/SceneTagger";
import { useMultiSelect } from "../hooks/useMultiSelect";
import { formatDuration, formatFileSize, getResolutionLabel, RatingBadge } from "../components/shared";
import { SCENE_CRITERIA } from "../components/FilterDialog";
import { BulkEditDialog, SCENE_BULK_FIELDS } from "../components/BulkEditDialog";
import { EditModal, Field, TextArea, TextInput, SaveButton } from "../components/EditModal";
import { Film, Tag, User, MapPin, Box, Eye, Trash2, Loader2, Edit, Merge, Search, Play, Images } from "lucide-react";
import { MergeDialog } from "../components/MergeDialog";
import { useSceneQueue } from "../state/SceneQueueContext";
import { IdentifyDialog } from "../components/IdentifyDialog";
import { SceneQueue } from "../components/SceneQueue";

const SORT_OPTIONS = [
  { value: "updated_at", label: "Recently Updated" },
  { value: "created_at", label: "Recently Added" },
  { value: "title", label: "Title" },
  { value: "date", label: "Date" },
  { value: "rating", label: "Rating" },
  { value: "play_count", label: "Play Count" },
  { value: "o_counter", label: "O-Counter" },
  { value: "duration", label: "Duration" },
  { value: "file_size", label: "File Size" },
  { value: "file_count", label: "File Count" },
  { value: "resolution", label: "Resolution" },
  { value: "framerate", label: "Frame Rate" },
  { value: "bitrate", label: "Bitrate" },
  { value: "tag_count", label: "Tag Count" },
  { value: "performer_count", label: "Performer Count" },
  { value: "last_played_at", label: "Last Played" },
  { value: "play_duration", label: "Play Duration" },
  { value: "resume_time", label: "Resume Time" },
  { value: "organized", label: "Organized" },
  { value: "random", label: "Random" },
];

interface Props {
  onNavigate: (r: any) => void;
}

export function ScenesPage({ onNavigate }: Props) {
  const [filter, setFilter] = useState<FindFilter>({ page: 1, perPage: 40, direction: "desc" });
  const [objectFilter, setObjectFilter] = useState<Record<string, unknown>>({});
  const [displayMode, setDisplayMode] = useState<DisplayMode>("grid");
  const [showCreate, setShowCreate] = useState(false);
  const [showBulkEdit, setShowBulkEdit] = useState(false);
  const [showMerge, setShowMerge] = useState(false);
  const [showIdentify, setShowIdentify] = useState(false);
  const [showQueue, setShowQueue] = useState(false);
  const queryClient = useQueryClient();
  const { setQueue } = useSceneQueue();

  const hasObjectFilter = Object.keys(objectFilter).length > 0;

  const { data, isLoading } = useQuery({
    queryKey: ["scenes", filter, objectFilter],
    queryFn: () =>
      hasObjectFilter
        ? scenes.findFiltered({ findFilter: filter, objectFilter: objectFilter as SceneFilterCriteria })
        : scenes.find(filter),
  });

  const items = data?.items ?? [];
  const { selectedIds, toggle, selectAll, selectNone } = useMultiSelect(items);
  const selecting = selectedIds.size > 0;

  const navigateToScene = useCallback((sceneId: number) => {
    const ids = items.map((s) => s.id);
    if (ids.length > 0) setQueue(ids, sceneId);
    onNavigate({ page: "scene", id: sceneId });
  }, [items, setQueue, onNavigate]);

  // Bulk delete
  const bulkDeleteMut = useMutation({
    mutationFn: () => scenes.bulkDelete([...selectedIds]),
    onSuccess: () => {
      selectNone();
      queryClient.invalidateQueries({ queryKey: ["scenes"] });
    },
  });

  // Bulk edit
  const bulkEditMut = useMutation({
    mutationFn: (values: Record<string, unknown>) =>
      scenes.bulkUpdate({
        ids: [...selectedIds],
        ...values,
      } as any),
    onSuccess: () => {
      setShowBulkEdit(false);
      selectNone();
      queryClient.invalidateQueries({ queryKey: ["scenes"] });
    },
  });

  // Metadata byline like original Stash: (1:23:45 - 2.5 GB)
  const byline = useMemo(() => {
    const items = data?.items ?? [];
    const totalDur = items.reduce((sum, s) => sum + (s.files[0]?.duration ?? 0), 0);
    const totalSize = items.reduce((sum, s) => sum + (s.files[0]?.size ?? 0), 0);
    if (!totalDur && !totalSize) return null;
    const parts: string[] = [];
    if (totalDur) parts.push(formatDuration(totalDur));
    if (totalSize) parts.push(formatFileSize(totalSize));
    return <span className="text-xs text-plex-text-muted">({parts.join(" — ")})</span>;
  }, [data?.items]);

  return (
    <>
    <SceneCreateModal open={showCreate} onClose={() => setShowCreate(false)} onCreated={(id) => onNavigate({ page: "scene", id })} />
    <ListPage
      title="Scenes"
      pageKey="scenes"
      filterMode="scenes"
      filter={filter}
      onFilterChange={setFilter}
      totalCount={data?.totalCount ?? 0}
      isLoading={isLoading}
      sortOptions={SORT_OPTIONS}
      displayMode={displayMode}
      onDisplayModeChange={setDisplayMode}
      availableDisplayModes={["grid", "list", "wall", "tagger"]}
      metadataByline={byline}
      criteriaDefinitions={SCENE_CRITERIA}
      objectFilter={objectFilter}
      onObjectFilterChange={setObjectFilter}
      onNew={() => setShowCreate(true)}
      selectedIds={selectedIds}
      onSelectAll={selectAll}
      onSelectNone={selectNone}
      selectionActions={
        <>
          <button
            onClick={() => setShowBulkEdit(true)}
            className="flex items-center gap-1 px-2 py-0.5 rounded text-xs text-plex-accent hover:text-plex-accent-hover hover:bg-plex-accent/10"
          >
            <Edit className="w-3 h-3" />
            Edit
          </button>
          <button
            onClick={() => setShowIdentify(true)}
            className="flex items-center gap-1 px-2 py-0.5 rounded text-xs text-blue-400 hover:text-blue-300 hover:bg-blue-900/20"
          >
            <Search className="w-3 h-3" />
            Identify
          </button>
          {selectedIds.size >= 2 && (
            <button
              onClick={() => setShowMerge(true)}
              className="flex items-center gap-1 px-2 py-0.5 rounded text-xs text-yellow-400 hover:text-yellow-300 hover:bg-yellow-900/20"
            >
              <Merge className="w-3 h-3" />
              Merge
            </button>
          )}
          <button
            onClick={() => setShowQueue(true)}
            className="flex items-center gap-1 px-2 py-0.5 rounded text-xs text-green-400 hover:text-green-300 hover:bg-green-900/20"
          >
            <Play className="w-3 h-3" />
            Play
          </button>
          <button
            onClick={() => { if (confirm(`Delete ${selectedIds.size} scene(s)?`)) bulkDeleteMut.mutate(); }}
            disabled={bulkDeleteMut.isPending}
            className="flex items-center gap-1 px-2 py-0.5 rounded text-xs text-red-400 hover:text-red-300 hover:bg-red-900/20"
          >
            {bulkDeleteMut.isPending ? <Loader2 className="w-3 h-3 animate-spin" /> : <Trash2 className="w-3 h-3" />}
            Delete
          </button>
        </>
      }
    >
      {displayMode === "grid" && (
        <div className="grid gap-2 px-2" style={{ gridTemplateColumns: "repeat(auto-fill, minmax(var(--card-min-width, 200px), 1fr))" }}>
          {items.map((scene) => (
            <SceneCard
              key={scene.id}
              scene={scene}
              onClick={() => navigateToScene(scene.id)}
              onNavigate={onNavigate}
              selected={selectedIds.has(scene.id)}
              onSelect={() => toggle(scene.id)}
              selecting={selecting}
            />
          ))}
        </div>
      )}
      {displayMode === "list" && (
        <SceneListTable scenes={items} onNavigate={onNavigate} selectedIds={selectedIds} onToggle={toggle} selecting={selecting} />
      )}
      {displayMode === "wall" && (
        <div className="columns-2 sm:columns-3 md:columns-4 lg:columns-5 gap-1 px-2">
          {items.map((scene) => (
            <SceneWallCard key={scene.id} scene={scene} onClick={() => navigateToScene(scene.id)} />
          ))}
        </div>
      )}
      {displayMode === "tagger" && (
        <SceneTagger scenes={items} />
      )}
      {items.length === 0 && !isLoading && (
        <div className="text-center py-20">
          <Film className="w-16 h-16 mx-auto mb-4 text-plex-text-muted opacity-50" />
          <p className="text-plex-text-secondary text-lg">No scenes found</p>
          <p className="text-plex-text-muted text-sm mt-1">Try scanning your library to discover content</p>
        </div>
      )}
    </ListPage>

    {/* Bulk Edit Dialog */}
    <BulkEditDialog
      open={showBulkEdit}
      onClose={() => setShowBulkEdit(false)}
      title="Edit Scenes"
      selectedCount={selectedIds.size}
      fields={SCENE_BULK_FIELDS}
      onApply={(values) => bulkEditMut.mutate(values)}
      isPending={bulkEditMut.isPending}
    />
    <MergeDialog
      open={showMerge}
      onClose={() => { setShowMerge(false); selectNone(); }}
      entityType="scene"
      items={items.filter((s) => selectedIds.has(s.id)).map((s) => ({ id: s.id, name: s.title || s.files[0]?.basename || `Scene ${s.id}` }))}
      onMerge={scenes.merge}
      queryKey="scenes"
    />
    <IdentifyDialog
      open={showIdentify}
      onClose={() => { setShowIdentify(false); selectNone(); }}
      sceneIds={[...selectedIds]}
    />
    {showQueue && (
      <SceneQueue
        scenes={items.filter((s) => selectedIds.has(s.id)).map((s) => ({
          id: s.id,
          title: s.title || s.files[0]?.basename,
          duration: s.files[0]?.duration,
          screenshotUrl: scenes.screenshotUrl(s.id),
        }))}
        onClose={() => { setShowQueue(false); selectNone(); }}
        onNavigate={onNavigate}
      />
    )}
    </>
  );
}

function SceneCreateModal({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const qc = useQueryClient();
  const [title, setTitle] = useState("");
  const [code, setCode] = useState("");
  const [date, setDate] = useState("");
  const [details, setDetails] = useState("");
  const [director, setDirector] = useState("");
  const [rating, setRating] = useState<number | undefined>(undefined);
  const [organized, setOrganized] = useState(false);
  const [urls, setUrls] = useState("");
  const [studioId, setStudioId] = useState<number | undefined>(undefined);

  const [tagSearch, setTagSearch] = useState("");
  const [selectedTags, setSelectedTags] = useState<{ id: number; name: string }[]>([]);
  const [performerSearch, setPerformerSearch] = useState("");
  const [selectedPerformers, setSelectedPerformers] = useState<{ id: number; name: string }[]>([]);
  const [gallerySearch, setGallerySearch] = useState("");
  const [selectedGalleries, setSelectedGalleries] = useState<{ id: number; title: string }[]>([]);

  const { data: studioResults } = useQuery({
    queryKey: ["studios-all"],
    queryFn: () => studios.find({ perPage: 500, sort: "name", direction: "asc" }),
  });

  const { data: tagResults } = useQuery({
    queryKey: ["tags-search", tagSearch],
    queryFn: () => tags.find({ q: tagSearch, perPage: 20, sort: "name", direction: "asc" }),
    enabled: tagSearch.length > 0,
  });

  const { data: performerResults } = useQuery({
    queryKey: ["performers-search", performerSearch],
    queryFn: () => performers.find({ q: performerSearch, perPage: 20, sort: "name", direction: "asc" }),
    enabled: performerSearch.length > 0,
  });

  const { data: galleryResults } = useQuery({
    queryKey: ["galleries-search", gallerySearch],
    queryFn: () => galleries.find({ q: gallerySearch, perPage: 20, sort: "title", direction: "asc" }),
    enabled: gallerySearch.length > 0,
  });

  const createMut = useMutation({
    mutationFn: (data: SceneCreate) => scenes.create(data),
    onSuccess: (created) => {
      qc.invalidateQueries({ queryKey: ["scenes"] });
      setTitle("");
      setCode("");
      setDate("");
      setDetails("");
      setDirector("");
      setRating(undefined);
      setOrganized(false);
      setUrls("");
      setStudioId(undefined);
      setSelectedTags([]);
      setSelectedPerformers([]);
      setSelectedGalleries([]);
      onClose();
      if (created?.id) onCreated(created.id);
    },
  });

  const handleSave = () => {
    const urlList = urls.split("\n").map((u) => u.trim()).filter(Boolean);
    createMut.mutate({
      title: title || undefined,
      code: code || undefined,
      date: date || undefined,
      details: details || undefined,
      director: director || undefined,
      rating,
      organized,
      studioId,
      urls: urlList,
      tagIds: selectedTags.map((t) => t.id),
      performerIds: selectedPerformers.map((p) => p.id),
      galleryIds: selectedGalleries.map((g) => g.id),
    });
  };

  return (
    <EditModal title="Create Scene" open={open} onClose={onClose}>
      <div className="grid grid-cols-2 gap-4">
        <Field label="Title">
          <TextInput value={title} onChange={setTitle} placeholder="Scene title" />
        </Field>
        <Field label="Date">
          <input
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          />
        </Field>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Field label="Code">
          <TextInput value={code} onChange={setCode} placeholder="Scene code" />
        </Field>
        <Field label="Director">
          <TextInput value={director} onChange={setDirector} placeholder="Director" />
        </Field>
      </div>

      <Field label="Details">
        <TextArea value={details} onChange={setDetails} placeholder="Scene description" rows={3} />
      </Field>

      <div className="grid grid-cols-2 gap-4">
        <RatingField value={rating} onChange={setRating} />
        <Field label="Studio">
          <select
            value={studioId ?? ""}
            onChange={(e) => setStudioId(e.target.value ? Number(e.target.value) : undefined)}
            className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option value="">None</option>
            {studioResults?.items.map((s) => (
              <option key={s.id} value={s.id}>{s.name}</option>
            ))}
          </select>
        </Field>
      </div>

      <Field label="URLs (one per line)">
        <TextArea value={urls} onChange={setUrls} placeholder="https://..." rows={2} />
      </Field>

      <label className="flex items-center gap-2 text-sm mb-2">
        <input
          type="checkbox"
          checked={organized}
          onChange={(e) => setOrganized(e.target.checked)}
          className="rounded bg-gray-800 border-gray-700"
        />
        Organized
      </label>

      <Field label="Tags">
        <div className="flex flex-wrap gap-1.5 mb-2">
          {selectedTags.map((t) => (
            <span key={t.id} className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-blue-900 text-blue-300">
              {t.name}
              <button onClick={() => setSelectedTags(selectedTags.filter((x) => x.id !== t.id))} className="hover:text-white">×</button>
            </span>
          ))}
        </div>
        <input
          type="text"
          value={tagSearch}
          onChange={(e) => setTagSearch(e.target.value)}
          placeholder="Search tags..."
          className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-1.5 text-sm text-gray-200 focus:outline-none focus:border-blue-500 mb-1"
        />
        {tagSearch && (tagResults?.items?.length ?? 0) > 0 && (
          <div className="max-h-32 overflow-y-auto bg-gray-800 rounded border border-gray-700">
            {(tagResults?.items ?? []).filter((t) => !selectedTags.some((x) => x.id === t.id)).slice(0, 10).map((t) => (
              <button
                key={t.id}
                onClick={() => { setSelectedTags([...selectedTags, { id: t.id, name: t.name }]); setTagSearch(""); }}
                className="block w-full text-left px-3 py-1.5 text-sm text-gray-300 hover:bg-gray-700"
              >
                {t.name}
              </button>
            ))}
          </div>
        )}
      </Field>

      <Field label="Performers">
        <div className="flex flex-wrap gap-1.5 mb-2">
          {selectedPerformers.map((p) => (
            <span key={p.id} className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-purple-900 text-purple-300">
              {p.name}
              <button onClick={() => setSelectedPerformers(selectedPerformers.filter((x) => x.id !== p.id))} className="hover:text-white">×</button>
            </span>
          ))}
        </div>
        <input
          type="text"
          value={performerSearch}
          onChange={(e) => setPerformerSearch(e.target.value)}
          placeholder="Search performers..."
          className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-1.5 text-sm text-gray-200 focus:outline-none focus:border-blue-500 mb-1"
        />
        {performerSearch && (performerResults?.items?.length ?? 0) > 0 && (
          <div className="max-h-32 overflow-y-auto bg-gray-800 rounded border border-gray-700">
            {(performerResults?.items ?? []).filter((p) => !selectedPerformers.some((x) => x.id === p.id)).slice(0, 10).map((p) => (
              <button
                key={p.id}
                onClick={() => { setSelectedPerformers([...selectedPerformers, { id: p.id, name: p.name }]); setPerformerSearch(""); }}
                className="block w-full text-left px-3 py-1.5 text-sm text-gray-300 hover:bg-gray-700"
              >
                {p.name}
              </button>
            ))}
          </div>
        )}
      </Field>

      <Field label="Galleries">
        <div className="flex flex-wrap gap-1.5 mb-2">
          {selectedGalleries.map((g) => (
            <span key={g.id} className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-emerald-900 text-emerald-300">
              {g.title}
              <button onClick={() => setSelectedGalleries(selectedGalleries.filter((x) => x.id !== g.id))} className="hover:text-white">×</button>
            </span>
          ))}
        </div>
        <input
          type="text"
          value={gallerySearch}
          onChange={(e) => setGallerySearch(e.target.value)}
          placeholder="Search galleries..."
          className="w-full bg-gray-800 border border-gray-700 rounded px-3 py-1.5 text-sm text-gray-200 focus:outline-none focus:border-blue-500 mb-1"
        />
        {gallerySearch && (galleryResults?.items?.length ?? 0) > 0 && (
          <div className="max-h-32 overflow-y-auto bg-gray-800 rounded border border-gray-700">
            {(galleryResults?.items ?? []).filter((g) => !selectedGalleries.some((x) => x.id === g.id)).slice(0, 10).map((g) => (
              <button
                key={g.id}
                onClick={() => { setSelectedGalleries([...selectedGalleries, { id: g.id, title: g.title || "Untitled" }]); setGallerySearch(""); }}
                className="block w-full text-left px-3 py-1.5 text-sm text-gray-300 hover:bg-gray-700"
              >
                {g.title || "Untitled"}
              </button>
            ))}
          </div>
        )}
      </Field>

      <div className="flex justify-end gap-3 mt-4">
        <button onClick={onClose} className="px-4 py-2 text-sm text-gray-400 hover:text-white">Cancel</button>
        <SaveButton loading={createMut.isPending} onClick={handleSave} />
      </div>
    </EditModal>
  );
}

/* ── Scene Card (matches Stash GridCard pattern) ── */

function SceneCard({ scene, onClick, selected, onSelect, onNavigate, selecting }: { scene: Scene; onClick: () => void; selected?: boolean; onSelect?: () => void; selecting?: boolean; onNavigate?: (r: any) => void }) {
  const file = scene.files[0];
  const duration = file?.duration ?? 0;
  const resLabel = file ? getResolutionLabel(file.width, file.height) : null;
  const screenshotUrl = scenes.screenshotUrl(scene.id);
  const previewUrl = scenes.previewUrl(scene.id);
  const videoRef = useRef<HTMLVideoElement>(null);
  const progressPercent = duration > 0 && scene.resumeTime ? Math.min(100, (scene.resumeTime / duration) * 100) : 0;

  // IntersectionObserver: when CSS hover moves video from top:-9999px to top:0,
  // intersection changes → play. When hover ends → moves back → pause.
  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;
    const observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.intersectionRatio > 0) video.play().catch(() => {});
        else video.pause();
      });
    });
    observer.observe(video);
    return () => observer.disconnect();
  }, []);

  return (
    <div
      onClick={onClick}
      className={`scene-card cursor-pointer group rounded border bg-plex-card overflow-hidden ${selected ? "ring-2 ring-plex-accent border-plex-accent" : "border-plex-border"}`}
    >
      {/* Image / Preview area */}
      <div className="scene-card-preview relative aspect-video bg-black overflow-hidden">
        {/* Screenshot image — always present as fallback */}
        <img
          src={screenshotUrl}
          alt={scene.title || ""}
          className="scene-card-preview-image w-full h-full object-cover"
          loading="lazy"
          onError={(e) => { (e.target as HTMLImageElement).style.display = "none"; }}
        />
        {/* Video preview — positioned off-screen, CSS hover brings it into view */}
        <video
          ref={videoRef}
          disableRemotePlayback
          playsInline
          muted
          loop
          preload="none"
          src={previewUrl}
          className="scene-card-preview-video"
        />
        {/* Checkbox — hidden by default, visible on hover or when selected */}
        <div className={`absolute top-1 left-1 z-10 ${selected || selecting ? "opacity-100" : "opacity-0 group-hover:opacity-100"} transition-opacity`}>
          <input
            type="checkbox"
            checked={selected}
            onChange={(e) => { e.stopPropagation(); onSelect?.(); }}
            onClick={(e) => e.stopPropagation()}
            className="w-4 h-4 rounded border-plex-border cursor-pointer accent-plex-accent"
          />
        </div>
        {/* Studio overlay (top-right like Stash) */}
        {scene.studioName && scene.studioId && !selecting && (
          <div className="absolute top-0 right-0 p-1 z-[5]">
            <img
              src={entityImages.studioImageUrl(scene.studioId)}
              alt={scene.studioName}
              className="max-h-8 max-w-[120px] object-contain drop-shadow-md"
              onError={(e) => {
                const el = e.target as HTMLImageElement;
                el.style.display = "none";
                if (el.nextElementSibling) (el.nextElementSibling as HTMLElement).style.display = "";
              }}
            />
            <span className="text-[10px] font-medium text-white bg-black/60 px-1.5 py-0.5 rounded" style={{ display: "none" }}>
              {scene.studioName}
            </span>
          </div>
        )}
        {/* Specs overlay (bottom-right: resolution + duration, matching Stash) */}
        {(duration > 0 || resLabel) && (
          <div className="scene-specs-overlay absolute bottom-0 right-0 flex items-center gap-0.5 px-1.5 py-1 text-[11px] text-white z-[5] transition-opacity">
            {file && (
              <span className="bg-black/70 px-1 py-0.5 rounded extra-scene-info hidden">
                {formatFileSize(file.size)}
              </span>
            )}
            {resLabel && (
              <span className="bg-black/70 px-1 py-0.5 rounded font-black uppercase">
                {resLabel}
              </span>
            )}
            {duration > 0 && (
              <span className="bg-black/70 px-1 py-0.5 rounded">
                {formatDuration(duration)}
              </span>
            )}
          </div>
        )}
        {/* Progress bar (resume time) */}
        {progressPercent > 0 && (
          <div className="absolute bottom-0 left-0 right-0 h-[3px] bg-black/40 z-[6]">
            <div className="h-full bg-plex-accent" style={{ width: `${progressPercent}%` }} />
          </div>
        )}
        {/* Rating banner (diagonal, like Stash) */}
        <RatingBanner rating={scene.rating} />
      </div>

      {/* Title section */}
      <div className="bg-plex-card rounded-b px-2 pt-1.5 pb-1 border-t border-plex-border">
        <p className="text-xs font-medium text-plex-text truncate group-hover:text-plex-accent transition-colors">
          {scene.title || file?.basename || "Untitled"}
        </p>
        {/* Date + path */}
        <div className="text-[10px] text-plex-text-muted mt-0.5">
          {scene.date && <span>{scene.date}</span>}
        </div>
        {scene.details && (
          <p className="text-[10px] text-plex-text-muted mt-0.5 line-clamp-2">{scene.details}</p>
        )}
      </div>

      {/* Popover buttons (like Stash's card-popovers) */}
      <SceneCardPopovers scene={scene} onNavigate={onNavigate} />
    </div>
  );
}

function SceneCardPopovers({ scene, onNavigate }: { scene: Scene; onNavigate?: (r: any) => void }) {
  const hasPopovers =
    scene.tags.length > 0 ||
    scene.performers.length > 0 ||
    scene.groups.length > 0 ||
    scene.galleries.length > 0 ||
    scene.markers.length > 0 ||
    scene.oCounter > 0 ||
    scene.organized;

  if (!hasPopovers) return null;

  return (
    <>
      <hr className="border-plex-border my-0" />
      <div className="flex flex-wrap items-center justify-center gap-1 px-2 py-1.5 bg-plex-card rounded-b card-popovers">
      {scene.performers.length > 0 && (
        <PopoverButton icon={<User className="w-3.5 h-3.5" />} count={scene.performers.length} title="Performers" wide>
          <div className="grid grid-cols-2 gap-2">
            {scene.performers.map((p: any) => (
              <button
                key={p.id}
                onClick={(e) => { e.stopPropagation(); onNavigate?.({ page: "performer", id: p.id }); }}
                className="flex flex-col items-center gap-1.5 text-center cursor-pointer rounded hover:bg-plex-card-hover p-1.5 group/perf transition-colors"
              >
                <div className="w-20 h-28 rounded overflow-hidden bg-plex-surface flex-shrink-0">
                  {p.imagePath ? (
                    <img src={p.imagePath} alt="" className="w-full h-full object-cover" />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center">
                      <User className="w-8 h-8 text-plex-text-muted" />
                    </div>
                  )}
                </div>
                <span className="text-xs text-plex-accent group-hover/perf:underline truncate w-full font-medium">{p.name}</span>
              </button>
            ))}
          </div>
        </PopoverButton>
      )}
      {scene.tags.length > 0 && (
        <PopoverButton icon={<Tag className="w-3.5 h-3.5" />} count={scene.tags.length} title="Tags">
          <div className="flex flex-col gap-0.5">
            {scene.tags.map((t: any) => (
              <button
                key={t.id}
                onClick={(e) => { e.stopPropagation(); onNavigate?.({ page: "tag", id: t.id }); }}
                className="text-xs text-plex-accent hover:underline cursor-pointer truncate text-left px-2 py-1 rounded hover:bg-plex-card-hover transition-colors"
              >
                {t.name}
              </button>
            ))}
          </div>
        </PopoverButton>
      )}
      {scene.oCounter > 0 && (
        <span className="flex items-center gap-0.5 p-1 text-plex-text-muted" title={`O-Counter: ${scene.oCounter}`}>
          <span className="text-xs font-bold">O</span>
          <span className="text-xs">{scene.oCounter}</span>
        </span>
      )}
      {scene.groups.length > 0 && (
        <PopoverButton icon={<Film className="w-3.5 h-3.5" />} count={scene.groups.length} title="Groups">
          <div className="flex flex-col gap-0.5">
            {scene.groups.map((g: any) => (
              <button
                key={g.id}
                onClick={(e) => { e.stopPropagation(); onNavigate?.({ page: "group", id: g.id }); }}
                className="text-xs text-plex-accent hover:underline cursor-pointer truncate text-left px-2 py-1 rounded hover:bg-plex-card-hover transition-colors"
              >
                {g.name}
              </button>
            ))}
          </div>
        </PopoverButton>
      )}
      {scene.galleries.length > 0 && (
        <PopoverButton icon={<Images className="w-3.5 h-3.5" />} count={scene.galleries.length} title="Galleries">
          <div className="flex flex-col gap-0.5">
            {scene.galleries.map((g: any) => (
              <button
                key={g.id}
                onClick={(e) => { e.stopPropagation(); onNavigate?.({ page: "gallery", id: g.id }); }}
                className="text-xs text-plex-accent hover:underline cursor-pointer truncate text-left px-2 py-1 rounded hover:bg-plex-card-hover transition-colors"
              >
                {g.title || "Untitled"}
              </button>
            ))}
          </div>
        </PopoverButton>
      )}
      {scene.markers.length > 0 && (
        <PopoverButton icon={<MapPin className="w-3.5 h-3.5" />} count={scene.markers.length} title="Markers">
          <div className="flex flex-col gap-0.5">
            {scene.markers.map((m: any) => (
              <button
                key={m.id}
                onClick={(e) => { e.stopPropagation(); onNavigate?.({ page: "scene", id: scene.id }); }}
                className="text-xs text-plex-accent hover:underline cursor-pointer truncate text-left px-2 py-1 rounded hover:bg-plex-card-hover transition-colors"
              >
                {m.title} ({formatDuration(m.seconds)})
              </button>
            ))}
          </div>
        </PopoverButton>
      )}
      {scene.organized && (
        <span className="p-1 text-plex-text-muted" title="Organized">
          <Box className="w-3.5 h-3.5" />
        </span>
      )}
    </div>
    </>
  );
}

function PopoverButton({ icon, count, title, children, wide }: { icon: React.ReactNode; count: number; title: string; children?: React.ReactNode; wide?: boolean }) {
  const [open, setOpen] = useState(false);
  const buttonRef = useRef<HTMLDivElement>(null);
  const popoverRef = useRef<HTMLDivElement>(null);
  const enterTimer = useRef<ReturnType<typeof setTimeout>>(undefined);
  const leaveTimer = useRef<ReturnType<typeof setTimeout>>(undefined);
  const [popoverStyle, setPopoverStyle] = useState<React.CSSProperties>({});

  const handleMouseEnter = useCallback(() => {
    clearTimeout(leaveTimer.current);
    enterTimer.current = setTimeout(() => {
      if (buttonRef.current) {
        const rect = buttonRef.current.getBoundingClientRect();
        const spaceAbove = rect.top;
        const showBelow = spaceAbove < 220;
        const style: React.CSSProperties = {
          position: "fixed",
          zIndex: 9999,
        };
        if (showBelow) {
          style.top = rect.bottom + 4;
        } else {
          style.bottom = window.innerHeight - rect.top + 4;
        }
        // Horizontal: try to center on button, clamp to viewport
        const centerX = rect.left + rect.width / 2;
        const popWidth = wide ? 300 : 220; // approximate min width
        let left = centerX - popWidth / 2;
        if (left < 8) left = 8;
        if (left + popWidth > window.innerWidth - 8) left = window.innerWidth - 8 - popWidth;
        style.left = left;
        setPopoverStyle(style);
      }
      setOpen(true);
    }, 200);
  }, []);

  const handleMouseLeave = useCallback(() => {
    clearTimeout(enterTimer.current);
    leaveTimer.current = setTimeout(() => setOpen(false), 200);
  }, []);

  useEffect(() => () => { clearTimeout(enterTimer.current); clearTimeout(leaveTimer.current); }, []);

  return (
    <div className="relative" ref={buttonRef} onMouseEnter={handleMouseEnter} onMouseLeave={handleMouseLeave}>
      <button
        className="flex items-center gap-1 px-1.5 py-1 text-plex-text-secondary hover:text-plex-accent rounded text-xs transition-colors"
        title={title}
        onClick={(e) => e.stopPropagation()}
      >
        {icon}
        <span className="font-medium">{count}</span>
      </button>
      {open && children && createPortal(
        <div
          ref={popoverRef}
          style={popoverStyle}
          className={`bg-plex-surface border border-plex-border rounded-lg shadow-2xl shadow-black/40 p-2.5 ${wide ? "min-w-[280px] max-w-[360px]" : "min-w-[180px] max-w-[min(280px,calc(100vw-1rem))]"} max-h-[320px] overflow-y-auto`}
          onClick={(e) => e.stopPropagation()}
          onMouseEnter={() => { clearTimeout(leaveTimer.current); }}
          onMouseLeave={handleMouseLeave}
        >
          <div className="text-[10px] uppercase tracking-wider text-plex-text-muted font-semibold mb-1.5 px-1">{title}</div>
          {children}
        </div>,
        document.body
      )}
    </div>
  );
}

/* ── Scene List Table ── */

function SceneListTable({ scenes, onNavigate, selectedIds, onToggle }: { scenes: Scene[]; onNavigate: (r: any) => void; selectedIds?: Set<number>; onToggle?: (id: number) => void; selecting?: boolean }) {
  return (
    <div className="overflow-x-auto px-2">
      <table className="w-full text-xs text-plex-text">
        <thead>
          <tr className="border-b border-plex-border text-plex-text-muted">
            {selectedIds && <th className="w-8 py-2 px-2"></th>}
            <th className="text-left py-2 px-2 font-medium">Title</th>
            <th className="text-left py-2 px-2 font-medium">Date</th>
            <th className="text-left py-2 px-2 font-medium">Rating</th>
            <th className="text-left py-2 px-2 font-medium">Duration</th>
            <th className="text-left py-2 px-2 font-medium">Size</th>
            <th className="text-left py-2 px-2 font-medium">Resolution</th>
            <th className="text-right py-2 px-2 font-medium">Plays</th>
          </tr>
        </thead>
        <tbody>
          {scenes.map((scene) => {
            const file = scene.files[0];
            return (
              <tr
                key={scene.id}
                onClick={() => onNavigate({ page: "scene", id: scene.id })}
                className={`border-b border-plex-border hover:bg-plex-card cursor-pointer ${selectedIds?.has(scene.id) ? "bg-plex-accent/10" : ""}`}
              >
                {selectedIds && (
                  <td className="py-1.5 px-2">
                    <input
                      type="checkbox"
                      checked={selectedIds.has(scene.id)}
                      onChange={() => onToggle?.(scene.id)}
                      onClick={(e) => e.stopPropagation()}
                      className="w-3.5 h-3.5 rounded border-plex-border cursor-pointer accent-plex-accent"
                    />
                  </td>
                )}
                <td className="py-1.5 px-2">
                  <span className="text-plex-text hover:text-plex-accent">
                    {scene.title || file?.basename || "Untitled"}
                  </span>
                  {scene.studioName && (
                    <span className="text-plex-text-muted ml-2">— {scene.studioName}</span>
                  )}
                </td>
                <td className="py-1.5 px-2 text-plex-text-muted">{scene.date || ""}</td>
                <td className="py-1.5 px-2"><RatingBadge rating={scene.rating} /></td>
                <td className="py-1.5 px-2 text-plex-text-muted">{file ? formatDuration(file.duration) : ""}</td>
                <td className="py-1.5 px-2 text-plex-text-muted">{file ? formatFileSize(file.size) : ""}</td>
                <td className="py-1.5 px-2 text-plex-text-muted">{file ? getResolutionLabel(file.width, file.height) : ""}</td>
                <td className="py-1.5 px-2 text-right text-plex-text-muted">{scene.playCount || ""}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

/* ── Scene Wall Card ── */

function SceneWallCard({ scene, onClick }: { scene: Scene; onClick: () => void }) {
  const file = scene.files[0];
  const screenshotUrl = scenes.screenshotUrl(scene.id);

  return (
    <div onClick={onClick} className="mb-1 cursor-pointer group break-inside-avoid">
      <div className="relative rounded overflow-hidden">
        <img
          src={screenshotUrl}
          alt={scene.title || ""}
          className="w-full object-cover"
          loading="lazy"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
        <div className="absolute bottom-0 left-0 right-0 p-1.5 opacity-0 group-hover:opacity-100 transition-opacity">
          <p className="text-[10px] text-white font-medium truncate">
            {scene.title || file?.basename || "Untitled"}
          </p>
        </div>
        {file && file.duration > 0 && (
          <span className="absolute top-1 right-1 text-[10px] text-white bg-black/70 px-1 rounded">
            {formatDuration(file.duration)}
          </span>
        )}
      </div>
    </div>
  );
}
