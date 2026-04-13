import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { scenes, tags, entityImages, performers as performersApi, studios as studiosApi, galleries as galleriesApi, groups as groupsApi } from "../api/client";
import { formatDuration, formatFileSize, formatDate, TagBadge, getResolutionLabel } from "../components/shared";
import { 
  Pencil, Plus, Trash2, Search, Eye, Heart, 
  Check, ChevronLeft, ChevronRight, MoreVertical, PanelLeftClose, PanelLeft,
  Play, Pause, Volume2, VolumeX, Maximize, Minimize,
  SkipBack, SkipForward, Gauge, Clapperboard, Monitor, FolderOpen, Layers
} from "lucide-react";
import { useState, useRef, useEffect, useCallback, Fragment, useMemo } from "react";
import { SceneEditModal } from "./SceneEditModal";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { GenerateDialog } from "../components/GenerateDialog";
import type { Scene, SceneMarkerCreate, SceneUpdate } from "../api/types";
import { ExtensionSlot } from "../router/RouteRegistry";
import { InteractiveRating, RatingField } from "../components/Rating";

interface Props {
  id: number;
  onNavigate: (r: any) => void;
}

type TabKey = "details" | "groups" | "galleries" | "markers" | "filters" | "file-info" | "edit" | "history";

export function SceneDetailPage({ id, onNavigate }: Props) {
  const { data: scene, isLoading } = useQuery({
    queryKey: ["scene", id],
    queryFn: () => scenes.get(id),
  });
  const [editing, setEditing] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [showGenerate, setShowGenerate] = useState(false);
  const [theaterMode, setTheaterMode] = useState(false);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [showOpsMenu, setShowOpsMenu] = useState(false);
  const [activeTab, setActiveTab] = useState<TabKey>("details");
  const queryClient = useQueryClient();
  const seekRef = useRef<((time: number) => void) | null>(null);
  const opsMenuRef = useRef<HTMLDivElement>(null);
  const [videoTime, setVideoTime] = useState(0);
  const [videoFilters, setVideoFilters] = useState({ brightness: 100, contrast: 100, gamma: 100, saturation: 100, hue: 0 });

  useEffect(() => {
    if (scene) document.title = `${scene.title || scene.files?.[0]?.basename || `Scene ${id}`} | Stash`;
    return () => { document.title = "Stash"; };
  }, [scene, id]);

  // Theater mode: hide navbar and expand layout
  useEffect(() => {
    if (theaterMode) {
      document.documentElement.classList.add("theater-mode");
    } else {
      document.documentElement.classList.remove("theater-mode");
    }
    return () => document.documentElement.classList.remove("theater-mode");
  }, [theaterMode]);

  // Keyboard shortcuts: "," for theater mode, a/e/k/i/h for tab navigation, o for O-counter
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const tag = (e.target as HTMLElement).tagName;
      if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return;
      switch (e.key) {
        case ",": setTheaterMode((prev) => !prev); break;
        case "a": setActiveTab("details"); break;
        case "e": setActiveTab("edit"); break;
        case "k": setActiveTab("markers"); break;
        case "i": setActiveTab("file-info"); break;
        case "h": setActiveTab("history"); break;
        case "o": if (scene) incrementOMut.mutate(); break;
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  // Close ops menu on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (opsMenuRef.current && !opsMenuRef.current.contains(e.target as Node)) {
        setShowOpsMenu(false);
      }
    };
    if (showOpsMenu) document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, [showOpsMenu]);

  // Apply CSS filters to video element when videoFilters change
  useEffect(() => {
    const video = document.querySelector('video');
    if (video) {
      const { brightness, contrast, saturation, hue } = videoFilters;
      video.style.filter = `brightness(${brightness}%) contrast(${contrast}%) saturate(${saturation}%) hue-rotate(${hue}deg)`;
    }
    return () => {
      const video = document.querySelector('video');
      if (video) video.style.filter = '';
    };
  }, [videoFilters]);

  const deleteMut = useMutation({
    mutationFn: () => scenes.delete(id),
    onSuccess: () => { 
      queryClient.invalidateQueries({ queryKey: ["scenes"] }); 
      onNavigate({ page: "scenes" }); 
    },
  });

  const incrementPlayMut = useMutation({
    mutationFn: () => scenes.recordPlay(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["scene", id] }),
  });

  const incrementOMut = useMutation({
    mutationFn: () => scenes.incrementO(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["scene", id] }),
  });

  const updateMut = useMutation({
    mutationFn: (data: { organized?: boolean; rating?: number }) => scenes.update(id, data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["scene", id] }),
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-plex-accent" />
      </div>
    );
  }

  if (!scene) return <div className="text-center text-plex-text-secondary py-16">Scene not found</div>;

  const file = scene.files[0];
  const streamUrl = scenes.streamUrl(id);
  const resLabel = file ? getResolutionLabel(file.width, file.height) : null;

  const tabs: { key: TabKey; label: string }[] = [
    { key: "details", label: "Details" },
    { key: "markers", label: "Markers" },
    ...(scene.groups.length > 0 ? [{ key: "groups" as TabKey, label: "Groups" }] : []),
    ...(scene.galleries.length > 0 ? [{ key: "galleries" as TabKey, label: "Galleries" }] : []),
    { key: "filters", label: "Filters" },
    { key: "file-info", label: `File Info${scene.files.length > 1 ? ` (${scene.files.length})` : ""}` },
    { key: "history", label: "History" },
    { key: "edit", label: "Edit" },
  ];

  const studioImageUrl = scene.studioId ? entityImages.studioImageUrl(scene.studioId) : null;

  return (
    <div className={theaterMode ? "-mx-6 -mt-5" : ""}>
      {scene && <SceneEditModal scene={scene} open={editing} onClose={() => setEditing(false)} />}
      <ConfirmDialog
        open={confirmDelete}
        title="Delete Scene"
        message={`Are you sure you want to delete "${scene.title || "Untitled"}"? This cannot be undone.`}
        onConfirm={() => deleteMut.mutate()}
        onCancel={() => setConfirmDelete(false)}
      />
      <GenerateDialog
        open={showGenerate}
        onClose={() => setShowGenerate(false)}
        sceneIds={[id]}
        title={`Generate for "${scene.title || "Untitled"}"`}
      />

      {/* Stash-style layout: left sidebar + right video */}
      <div className={theaterMode ? "flex flex-col" : "flex flex-col xl:flex-row"} style={theaterMode ? undefined : { height: "calc(100vh - 48px)" }}>
        {/* Left sidebar: metadata, tabs, tab content */}
        {!theaterMode && !sidebarCollapsed && (
          <div
            className="w-full xl:border-r border-b xl:border-b-0 border-plex-border overflow-y-auto"
            style={{ flex: "0 0 450px", maxWidth: 450, maxHeight: "calc(100vh - 48px)" }}
          >
            <div className="px-4 pt-4 pb-2">
              {/* Studio logo */}
              {studioImageUrl && scene.studioId && (
                <div className="mb-3 flex items-start gap-4">
                  <button
                    onClick={() => onNavigate({ page: "studio", id: scene.studioId })}
                    className="flex-shrink-0"
                  >
                    <img
                      src={studioImageUrl}
                      alt={scene.studioName || "Studio"}
                      className="max-h-[5rem] max-w-full object-contain"
                      onError={(e) => { (e.target as HTMLImageElement).style.display = "none"; }}
                    />
                  </button>
                </div>
              )}

              {/* Title — large like original's h3 */}
              <h3 className="text-[1.5rem] font-semibold text-plex-text leading-snug line-clamp-2 mt-1">
                {scene.title || file?.path.split(/[\\/]/).pop() || "Untitled"}
              </h3>

              {/* Subheader: date left, resolution+fps right */}
              <div className="flex items-center justify-between mt-2 text-sm text-plex-text-secondary">
                <span>{scene.date ? new Date(scene.date + "T00:00:00").toLocaleDateString(undefined, { year: "numeric", month: "long", day: "numeric" }) : ""}</span>
                <span className="flex items-center gap-1.5">
                  {file && file.frameRate > 0 && <span>{file.frameRate.toFixed(0)} fps</span>}
                  {file && resLabel && <span className="text-plex-accent font-bold">{resLabel}</span>}
                </span>
              </div>

              {/* Studio name text fallback (when no logo) */}
              {scene.studioName && scene.studioId && !studioImageUrl && (
                <button 
                  onClick={() => onNavigate({ page: "studio", id: scene.studioId })}
                  className="text-plex-accent hover:underline text-sm mt-1 block"
                >
                  {scene.studioName}
                </button>
              )}

              {/* Toolbar: rating left, counters + ops right — single row */}
              <div className="flex items-center justify-between mt-3 gap-2">
                <InteractiveRating value={scene.rating} onChange={(value) => updateMut.mutate({ rating: value })} />
                <div className="flex items-center gap-2">
                  <button 
                    onClick={() => incrementPlayMut.mutate()}
                    className="flex items-center gap-1 text-sm text-plex-text-secondary hover:text-plex-text"
                    title="Play count"
                  >
                    <Eye className="w-4 h-4" />
                    <span>{scene.playCount}</span>
                  </button>
                  <button 
                    onClick={() => incrementOMut.mutate()}
                    className="flex items-center gap-1 text-sm text-plex-text-secondary hover:text-plex-accent"
                    title="O counter"
                  >
                    <Heart className={`w-4 h-4 ${scene.oCounter > 0 ? "fill-plex-accent text-plex-accent" : ""}`} />
                    <span>{scene.oCounter}</span>
                  </button>
                  <button 
                    onClick={() => updateMut.mutate({ organized: !scene.organized })}
                    className={`p-1 rounded ${scene.organized ? "bg-green-600 text-white" : "bg-plex-card text-plex-text-muted hover:text-plex-text"}`}
                    title={scene.organized ? "Organized" : "Not organized"}
                  >
                    <Check className="w-4 h-4" />
                  </button>
                  {/* Operations dropdown */}
                  <div className="relative" ref={opsMenuRef}>
                    <button
                      onClick={() => setShowOpsMenu(!showOpsMenu)}
                      className="p-1 rounded text-plex-text-secondary hover:text-plex-text hover:bg-plex-card"
                      title="Operations"
                    >
                      <MoreVertical className="w-4 h-4" />
                    </button>
                    {showOpsMenu && (
                      <div className="absolute right-0 top-full mt-1 z-50 min-w-[160px] bg-plex-card border border-plex-border rounded shadow-lg py-1">
                        <button onClick={() => { setEditing(true); setShowOpsMenu(false); }} className="w-full px-3 py-1.5 text-left text-sm text-plex-text hover:bg-plex-surface flex items-center gap-2"><Pencil className="w-3.5 h-3.5" /> Edit</button>
                        <button onClick={() => { setShowGenerate(true); setShowOpsMenu(false); }} className="w-full px-3 py-1.5 text-left text-sm text-plex-text hover:bg-plex-surface flex items-center gap-2"><Clapperboard className="w-3.5 h-3.5" /> Generate…</button>
                        <button onClick={() => { setTheaterMode(true); setShowOpsMenu(false); }} className="w-full px-3 py-1.5 text-left text-sm text-plex-text hover:bg-plex-surface flex items-center gap-2"><Monitor className="w-3.5 h-3.5" /> Theater Mode</button>
                        <div className="border-t border-plex-border my-1" />
                        <button onClick={() => { setConfirmDelete(true); setShowOpsMenu(false); }} className="w-full px-3 py-1.5 text-left text-sm text-red-400 hover:bg-plex-surface flex items-center gap-2"><Trash2 className="w-3.5 h-3.5" /> Delete</button>
                      </div>
                    )}
                  </div>
                  <ExtensionSlot slot="scene-detail-actions" context={{ scene, onNavigate }} />
                </div>
              </div>
            </div>

            {/* Tab Navigation */}
            <div className="px-4">
              <div className="flex flex-wrap border-b border-plex-border">
                {tabs.map((tab) => (
                  <button
                    key={tab.key}
                    onClick={() => setActiveTab(tab.key)}
                    className={`px-2.5 py-2 text-sm transition-colors border-b-2 ${
                      activeTab === tab.key 
                        ? "border-plex-accent text-plex-accent" 
                        : "border-transparent text-plex-text-secondary hover:text-plex-text"
                    }`}
                  >
                    {tab.label}
                  </button>
                ))}
              </div>
            </div>

            {/* Tab Content */}
            <div className="px-4 py-4">
              {activeTab === "details" && (
                <DetailsTab scene={scene} onNavigate={onNavigate} />
              )}
              {activeTab === "groups" && (
                <GroupsTab scene={scene} onNavigate={onNavigate} />
              )}
              {activeTab === "galleries" && (
                <GalleriesTab scene={scene} onNavigate={onNavigate} />
              )}
              {activeTab === "markers" && (
                <MarkersPanel sceneId={scene.id} markers={scene.markers} onSeek={(t) => seekRef.current?.(t)} />
              )}
              {activeTab === "filters" && (
                <VideoFiltersTab filters={videoFilters} onChange={setVideoFilters} />
              )}
              {activeTab === "file-info" && file && (
                <FileInfoTab file={file} />
              )}
              {activeTab === "history" && (
                <HistoryTab scene={scene} />
              )}
              {activeTab === "edit" && (
                <SceneEditPanel scene={scene} onSaved={() => setActiveTab("details")} />
              )}
            </div>
          </div>
        )}

        {/* Sidebar collapse/expand divider */}
        {!theaterMode && (
          <button
            onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
            className="hidden xl:flex items-center justify-center bg-plex-surface/50 hover:bg-plex-surface border-r border-plex-border transition-colors"
            style={{ flex: "0 0 15px", maxWidth: 15 }}
            title={sidebarCollapsed ? "Show sidebar" : "Hide sidebar"}
          >
            {sidebarCollapsed ? <ChevronRight className="w-4 h-4 text-plex-text-muted" /> : <ChevronLeft className="w-4 h-4 text-plex-text-muted" />}
          </button>
        )}

        {/* Right side: video player + scrubber */}
        <div className="min-w-0 flex flex-col" style={{ flex: sidebarCollapsed ? "0 0 calc(100% - 15px)" : "0 0 calc(100% - 450px - 15px)", maxWidth: sidebarCollapsed ? "calc(100% - 15px)" : "calc(100% - 465px)" }}>
          <div className="bg-black flex-1 flex flex-col">
            {file ? (
              <VideoPlayer
                streamUrl={streamUrl}
                format={file.format}
                duration={file.duration}
                resumeTime={scene.resumeTime}
                sceneId={id}
                onPlay={() => incrementPlayMut.mutate()}
                markers={scene.markers}
                onSeekRegister={(fn) => { seekRef.current = fn; }}
                onTimeUpdate={setVideoTime}
              />
            ) : (
              <div className="flex items-center justify-center h-48 text-plex-text-muted">No video file available</div>
            )}
          </div>
          {/* Scene scrubber */}
          {file && <SceneScrubber sceneId={scene.id} duration={file.duration} markers={scene.markers} onSeek={(t) => seekRef.current?.(t)} currentTime={videoTime} />}

          {/* Theater mode: show metadata below video */}
          {theaterMode && (
            <div className="px-4 pt-3 max-w-5xl mx-auto">
              <h1 className="text-xl font-bold text-plex-text">{scene.title || file?.path.split(/[\\/]/).pop() || "Untitled"}</h1>
              <div className="flex items-center gap-3 mt-2 flex-wrap">
                <InteractiveRating value={scene.rating} onChange={(value) => updateMut.mutate({ rating: value })} />
                <button onClick={() => incrementPlayMut.mutate()} className="flex items-center gap-1 text-sm text-plex-text-secondary hover:text-plex-text"><Eye className="w-4 h-4" />{scene.playCount}</button>
                <button onClick={() => incrementOMut.mutate()} className="flex items-center gap-1 text-sm text-plex-text-secondary hover:text-plex-accent"><Heart className={`w-4 h-4 ${scene.oCounter > 0 ? "fill-plex-accent text-plex-accent" : ""}`} />{scene.oCounter}</button>
                <button onClick={() => setTheaterMode(false)} className="flex items-center gap-1 px-2 py-1 text-xs bg-plex-accent text-white rounded"><Monitor className="w-3 h-3" /> Exit Theater</button>
              </div>
            </div>
          )}
        </div>
      </div>

      <ExtensionSlot slot="scene-detail-main-bottom" context={{ scene, onNavigate }} />
    </div>
  );
}

// Details Tab Content
function DetailsTab({ scene, onNavigate }: { scene: Scene; onNavigate: (r: any) => void }) {
  return (
    <div className="space-y-4">
      {/* Created/Updated + Code/Director at top like original */}
      <dl className="grid gap-y-1.5 text-sm" style={{ gridTemplateColumns: "auto 1fr" }}>
        <dt className="text-plex-text-muted pr-3">Created</dt>
        <dd className="text-plex-text">{formatDate(scene.createdAt)}</dd>
        <dt className="text-plex-text-muted pr-3">Updated</dt>
        <dd className="text-plex-text">{formatDate(scene.updatedAt)}</dd>
        {scene.code && (
          <>
            <dt className="text-plex-text-muted pr-3">Scene Code</dt>
            <dd className="text-plex-text">{scene.code}</dd>
          </>
        )}
        {scene.director && (
          <>
            <dt className="text-plex-text-muted pr-3">Director</dt>
            <dd className="text-plex-accent">{scene.director}</dd>
          </>
        )}
      </dl>

      {/* Details / Description */}
      {scene.details && (
        <div>
          <p className="text-sm text-plex-text whitespace-pre-wrap">{scene.details}</p>
        </div>
      )}

      {/* Tags */}
      {scene.tags.length > 0 && (
        <div>
          <h6 className="text-sm text-plex-text-muted mb-2">Tags</h6>
          <div className="flex flex-wrap gap-1.5">
            {scene.tags.map((tag: any) => (
              <TagBadge 
                key={tag.id} 
                name={tag.name} 
                onClick={() => onNavigate({ page: "tag", id: tag.id })} 
              />
            ))}
          </div>
        </div>
      )}

      {/* Performers — larger cards matching original's 15rem width */}
      {scene.performers.length > 0 && (
        <div>
          <h6 className="text-sm text-plex-text-muted mb-2">Performer{scene.performers.length > 1 ? "s" : ""}</h6>
          <div className="flex flex-wrap gap-3">
            {scene.performers.map((performer: any) => (
              <PerformerCard 
                key={performer.id} 
                performer={performer}
                sceneDate={scene.date}
                onClick={() => onNavigate({ page: "performer", id: performer.id })}
              />
            ))}
          </div>
        </div>
      )}

      {/* URLs */}
      {scene.urls && scene.urls.length > 0 && (
        <div>
          <h6 className="text-sm text-plex-text-muted mb-2">URLs</h6>
          <div className="space-y-1">
            {scene.urls.map((url: string, i: number) => (
              <a
                key={i}
                href={url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-plex-accent hover:underline text-sm block truncate"
              >
                {url}
              </a>
            ))}
          </div>
        </div>
      )}

      {/* Stash IDs */}
      {scene.stashIds && scene.stashIds.length > 0 && (
        <div>
          <h6 className="text-sm text-plex-text-muted mb-2">Stash IDs</h6>
          <dl className="grid gap-y-1 text-sm" style={{ gridTemplateColumns: "auto 1fr" }}>
            {scene.stashIds.map((sid, i) => (
              <Fragment key={i}>
                <dt className="text-plex-text-muted pr-3 truncate">{sid.endpoint}</dt>
                <dd className="text-plex-text font-mono text-xs break-all">{sid.stashId}</dd>
              </Fragment>
            ))}
          </dl>
        </div>
      )}
    </div>
  );
}

function GroupsTab({ scene, onNavigate }: { scene: Scene; onNavigate: (r: any) => void }) {
  if (scene.groups.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-plex-border bg-plex-card/40 px-4 py-10 text-center text-sm text-plex-text-secondary">
        <Layers className="mx-auto mb-3 h-8 w-8 text-plex-text-muted" />
        No groups linked to this scene.
      </div>
    );
  }

  return (
    <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      {scene.groups.map((group) => (
        <button
          key={group.id}
          onClick={() => onNavigate({ page: "group", id: group.id })}
          className="rounded-xl border border-plex-border bg-plex-card p-4 text-left transition-colors hover:border-plex-accent/60"
        >
          <div className="flex items-center justify-between gap-3">
            <div>
              <div className="text-sm font-medium text-plex-text">{group.name}</div>
              <div className="mt-1 text-xs text-plex-text-secondary">Scene #{group.sceneIndex}</div>
            </div>
            <Layers className="h-5 w-5 text-plex-text-muted" />
          </div>
        </button>
      ))}
    </div>
  );
}

function GalleriesTab({ scene, onNavigate }: { scene: Scene; onNavigate: (r: any) => void }) {
  if (scene.galleries.length === 0) {
    return (
      <div className="rounded-xl border border-dashed border-plex-border bg-plex-card/40 px-4 py-10 text-center text-sm text-plex-text-secondary">
        <FolderOpen className="mx-auto mb-3 h-8 w-8 text-plex-text-muted" />
        No galleries linked to this scene.
      </div>
    );
  }

  return (
    <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      {scene.galleries.map((gallery) => (
        <button
          key={gallery.id}
          onClick={() => onNavigate({ page: "gallery", id: gallery.id })}
          className="group overflow-hidden rounded-xl border border-plex-border bg-plex-card text-left transition-colors hover:border-plex-accent/60"
        >
          <div className="flex aspect-video items-center justify-center bg-gradient-to-br from-plex-surface to-plex-card">
            <FolderOpen className="h-10 w-10 text-plex-text-muted" />
          </div>
          <div className="p-3">
            <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">
              {gallery.title || "Untitled"}
            </p>
            {gallery.date && (
              <p className="mt-1 text-xs text-plex-text-secondary">{formatDate(gallery.date)}</p>
            )}
          </div>
        </button>
      ))}
    </div>
  );
}

// Performer Card matching Stash style — ~200px wide with taller image
function PerformerCard({ performer, sceneDate, onClick }: { performer: any; sceneDate?: string; onClick: () => void }) {
  const imageUrl = performer.imagePath;
  // Calculate age at scene date
  const ageAtScene = (() => {
    if (!sceneDate || !performer.birthdate) return null;
    const scene = new Date(sceneDate);
    const birth = new Date(performer.birthdate);
    let age = scene.getFullYear() - birth.getFullYear();
    const m = scene.getMonth() - birth.getMonth();
    if (m < 0 || (m === 0 && scene.getDate() < birth.getDate())) age--;
    return age > 0 ? age : null;
  })();

  return (
    <button
      onClick={onClick}
      className="bg-plex-card border border-plex-border rounded overflow-hidden hover:border-plex-accent/60 transition-colors text-left"
      style={{ width: "200px" }}
    >
      <div className="aspect-[2/3] bg-plex-surface flex items-center justify-center relative">
        {imageUrl ? (
          <img src={imageUrl} alt={performer.name} className="w-full h-full object-cover" />
        ) : (
          <div className="w-full h-full flex items-center justify-center bg-gradient-to-b from-plex-card to-plex-surface">
            <svg viewBox="0 0 100 150" className="w-2/3 h-2/3 opacity-30">
              <ellipse cx="50" cy="35" rx="25" ry="30" fill="currentColor" className="text-plex-text-muted"/>
              <ellipse cx="50" cy="120" rx="40" ry="45" fill="currentColor" className="text-plex-text-muted"/>
            </svg>
          </div>
        )}
      </div>
      <div className="p-2 text-center">
        <div className="text-sm text-plex-text font-medium truncate">{performer.name}</div>
        <div className="text-xs text-plex-text-muted flex items-center justify-center gap-1 mt-0.5">
          {ageAtScene && <span>{ageAtScene} yrs old</span>}
          {ageAtScene && performer.sceneCount !== undefined && <span>·</span>}
          {performer.sceneCount !== undefined && (
            <span className="flex items-center gap-0.5"><Eye className="w-3 h-3" /> {performer.sceneCount}</span>
          )}
        </div>
      </div>
    </button>
  );
}

// File Info Tab — dl grid matching original's details-list
function FileInfoTab({ file }: { file: any }) {
  return (
    <div className="space-y-3 text-sm">
      <dl className="grid gap-y-1.5" style={{ gridTemplateColumns: "minmax(100px, auto) 1fr" }}>
        <dt className="text-plex-text-muted">Path</dt>
        <dd className="text-plex-text break-all font-mono text-xs">{file.path}</dd>

        <dt className="text-plex-text-muted">File Size</dt>
        <dd className="text-plex-text">{formatFileSize(file.size)}</dd>

        <dt className="text-plex-text-muted">Duration</dt>
        <dd className="text-plex-text">{formatDuration(file.duration)}</dd>

        <dt className="text-plex-text-muted">Dimensions</dt>
        <dd className="text-plex-text">{file.width}×{file.height}</dd>

        <dt className="text-plex-text-muted">Frame Rate</dt>
        <dd className="text-plex-text">{file.frameRate.toFixed(2)} fps</dd>

        <dt className="text-plex-text-muted">Bitrate</dt>
        <dd className="text-plex-text">{Math.round(file.bitRate / 1000)} kbps</dd>

        <dt className="text-plex-text-muted">Video Codec</dt>
        <dd className="text-plex-text">{file.videoCodec}</dd>

        <dt className="text-plex-text-muted">Audio Codec</dt>
        <dd className="text-plex-text">{file.audioCodec}</dd>
      </dl>

      {file.fingerprints && file.fingerprints.length > 0 && (
        <div>
          <h6 className="text-sm text-plex-text-muted mb-1 font-medium">Fingerprints</h6>
          <dl className="grid gap-y-1" style={{ gridTemplateColumns: "auto 1fr" }}>
            {file.fingerprints.map((fp: any) => (
              <Fragment key={fp.type}>
                <dt className="text-plex-text-muted text-xs pr-3">{fp.type}</dt>
                <dd className="text-plex-text font-mono text-xs break-all">{fp.value}</dd>
              </Fragment>
            ))}
          </dl>
        </div>
      )}
    </div>
  );
}

// History Tab
function HistoryTab({ scene }: { scene: any }) {
  return (
    <div className="space-y-3 text-sm">
      <div className="grid grid-cols-2 gap-3">
        <div>
          <div className="text-plex-text-muted">Play Count</div>
          <div className="text-plex-text">{scene.playCount}</div>
        </div>
        <div>
          <div className="text-plex-text-muted">Play Duration</div>
          <div className="text-plex-text">{formatDuration(scene.playDuration)}</div>
        </div>
        <div>
          <div className="text-plex-text-muted">O-Counter</div>
          <div className="text-plex-text">{scene.oCounter}</div>
        </div>
        <div>
          <div className="text-plex-text-muted">Organized</div>
          <div className="text-plex-text">{scene.organized ? "Yes" : "No"}</div>
        </div>
      </div>
      <div>
        <div className="text-plex-text-muted">Created</div>
        <div className="text-plex-text">{formatDate(scene.createdAt)}</div>
      </div>
      <div>
        <div className="text-plex-text-muted">Updated</div>
        <div className="text-plex-text">{formatDate(scene.updatedAt)}</div>
      </div>
      {scene.urls && scene.urls.length > 0 && (
        <div>
          <div className="text-plex-text-muted mb-1">URLs</div>
          {scene.urls.map((url: string, i: number) => (
            <a
              key={i}
              href={url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-plex-accent hover:underline text-sm block truncate"
            >
              {url}
            </a>
          ))}
        </div>
      )}
    </div>
  );
}

// Video Filters Tab — matches original Stash's brightness/contrast/gamma/saturation/hue
interface VideoFilters {
  brightness: number;
  contrast: number;
  gamma: number;
  saturation: number;
  hue: number;
}

function VideoFiltersTab({ filters, onChange }: { filters: VideoFilters; onChange: (f: VideoFilters) => void }) {
  const sliders: { key: keyof VideoFilters; label: string; min: number; max: number; default: number; unit: string; formatValue?: (v: number) => string }[] = [
    { key: "brightness", label: "Brightness", min: 0, max: 200, default: 100, unit: "%" },
    { key: "contrast", label: "Contrast", min: 0, max: 200, default: 100, unit: "%" },
    { key: "gamma", label: "Gamma", min: 0, max: 200, default: 100, unit: "", formatValue: (v) => String(v - 100) },
    { key: "saturation", label: "Saturation", min: 0, max: 200, default: 100, unit: "%" },
    { key: "hue", label: "Hue", min: -180, max: 180, default: 0, unit: "°" },
  ];

  const handleReset = () => onChange({ brightness: 100, contrast: 100, gamma: 100, saturation: 100, hue: 0 });

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h5 className="text-sm font-medium text-plex-text">Filters</h5>
        <button onClick={handleReset} className="text-xs text-plex-accent hover:underline">Reset All</button>
      </div>
      {sliders.map(({ key, label, min, max, default: def, unit, formatValue }) => (
        <div key={key} className="flex items-center gap-3">
          <span className="text-sm text-plex-text-muted w-24 flex-shrink-0">{label}</span>
          <input
            type="range"
            min={min}
            max={max}
            value={filters[key]}
            onChange={(e) => onChange({ ...filters, [key]: Number(e.target.value) })}
            className="flex-1 h-1 accent-plex-accent cursor-pointer"
          />
          <button
            onClick={() => onChange({ ...filters, [key]: def })}
            className="text-xs text-plex-text-secondary hover:text-plex-text w-12 text-right cursor-pointer"
            title="Click to reset"
          >
            {formatValue ? formatValue(filters[key]) : `${filters[key]}${unit}`}
          </button>
        </div>
      ))}
    </div>
  );
}

/* ── Video Player with custom controls ── */
const PLAYBACK_RATES = [0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2];
const VOLUME_KEY = "stash-player-volume";
const MUTED_KEY = "stash-player-muted";

function VideoPlayer({
  streamUrl,
  format,
  duration,
  resumeTime,
  sceneId,
  onPlay,
  markers,
  onSeekRegister,
  onTimeUpdate: onTimeUpdateProp,
}: {
  streamUrl: string;
  format: string;
  duration: number;
  resumeTime?: number;
  sceneId: number;
  onPlay: () => void;
  markers: { id: number; title: string; seconds: number; primaryTagName: string }[];
  onSeekRegister?: (fn: (time: number) => void) => void;
  onTimeUpdate?: (time: number) => void;
}) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [playing, setPlaying] = useState(false);
  const [currentTime, setCurTime] = useState(0);
  const [buffered, setBuffered] = useState(0);
  const [vol, setVol] = useState(() => {
    const saved = localStorage.getItem(VOLUME_KEY);
    return saved ? Number(saved) : 1;
  });
  const [muted, setMuted] = useState(() => localStorage.getItem(MUTED_KEY) === "true");
  const [fullscreen, setFullscreen] = useState(false);
  const [showControls, setShowControls] = useState(true);
  const [showSpeed, setShowSpeed] = useState(false);
  const [rate, setRate] = useState(1);
  const hideTimer = useRef<ReturnType<typeof setTimeout>>(null);
  const playTriggered = useRef(false);
  const activityTimer = useRef<ReturnType<typeof setTimeout>>(null);

  // Restore volume
  useEffect(() => {
    const v = videoRef.current;
    if (!v) return;
    v.volume = vol;
    v.muted = muted;
  }, []);

  // Register seek callback for external components (markers, scrubber)
  useEffect(() => {
    if (onSeekRegister) {
      onSeekRegister((time: number) => {
        const v = videoRef.current;
        if (v) {
          v.currentTime = time;
          v.play().catch(() => {});
        }
      });
    }
  }, [onSeekRegister]);

  // Resume from saved position
  useEffect(() => {
    const v = videoRef.current;
    if (v && resumeTime && resumeTime > 0) {
      v.currentTime = resumeTime;
    }
  }, [resumeTime]);

  // Save activity periodically
  useEffect(() => {
    const saveActivity = () => {
      const v = videoRef.current;
      if (v && !v.paused && v.currentTime > 0) {
        scenes.saveActivity(sceneId, { resumeTime: v.currentTime, playDuration: 5 }).catch(() => {});
      }
    };
    activityTimer.current = setInterval(saveActivity, 5000);
    return () => { if (activityTimer.current) clearInterval(activityTimer.current); };
  }, [sceneId]);

  // Fullscreen change listener
  useEffect(() => {
    const handler = () => setFullscreen(!!document.fullscreenElement);
    document.addEventListener("fullscreenchange", handler);
    return () => document.removeEventListener("fullscreenchange", handler);
  }, []);

  // Auto-hide controls
  const resetHideTimer = useCallback(() => {
    setShowControls(true);
    if (hideTimer.current) clearTimeout(hideTimer.current);
    hideTimer.current = setTimeout(() => {
      if (videoRef.current && !videoRef.current.paused) setShowControls(false);
    }, 3000);
  }, []);

  // Keyboard shortcuts
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const v = videoRef.current;
      if (!v) return;
      const tag = (e.target as HTMLElement).tagName;
      if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return;

      switch (e.key) {
        case " ":
        case "k":
          e.preventDefault();
          v.paused ? v.play() : v.pause();
          break;
        case "ArrowLeft":
          e.preventDefault();
          v.currentTime = Math.max(0, v.currentTime - (e.shiftKey ? 10 : 5));
          break;
        case "ArrowRight":
          e.preventDefault();
          v.currentTime = Math.min(v.duration, v.currentTime + (e.shiftKey ? 10 : 5));
          break;
        case "ArrowUp":
          e.preventDefault();
          v.volume = Math.min(1, v.volume + 0.1);
          setVol(v.volume);
          localStorage.setItem(VOLUME_KEY, String(v.volume));
          break;
        case "ArrowDown":
          e.preventDefault();
          v.volume = Math.max(0, v.volume - 0.1);
          setVol(v.volume);
          localStorage.setItem(VOLUME_KEY, String(v.volume));
          break;
        case "m":
          v.muted = !v.muted;
          setMuted(v.muted);
          localStorage.setItem(MUTED_KEY, String(v.muted));
          break;
        case "f":
          if (document.fullscreenElement) document.exitFullscreen();
          else containerRef.current?.requestFullscreen();
          break;
        case "0": case "1": case "2": case "3": case "4":
        case "5": case "6": case "7": case "8": case "9":
          e.preventDefault();
          v.currentTime = v.duration * (Number(e.key) / 10);
          break;
      }
      resetHideTimer();
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [resetHideTimer]);

  const togglePlay = () => {
    const v = videoRef.current;
    if (!v) return;
    v.paused ? v.play() : v.pause();
  };

  const seekTo = (e: React.MouseEvent<HTMLDivElement>) => {
    const v = videoRef.current;
    if (!v) return;
    const rect = e.currentTarget.getBoundingClientRect();
    const pct = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width));
    v.currentTime = pct * v.duration;
  };

  const changeVolume = (e: React.MouseEvent<HTMLDivElement>) => {
    const v = videoRef.current;
    if (!v) return;
    const rect = e.currentTarget.getBoundingClientRect();
    const pct = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width));
    v.volume = pct;
    v.muted = false;
    setVol(pct);
    setMuted(false);
    localStorage.setItem(VOLUME_KEY, String(pct));
    localStorage.setItem(MUTED_KEY, "false");
  };

  const toggleFullscreen = () => {
    if (document.fullscreenElement) document.exitFullscreen();
    else containerRef.current?.requestFullscreen();
  };

  const changeRate = (r: number) => {
    const v = videoRef.current;
    if (v) v.playbackRate = r;
    setRate(r);
    setShowSpeed(false);
  };

  const fmtTime = (s: number) => {
    if (!isFinite(s)) return "0:00";
    const h = Math.floor(s / 3600);
    const m = Math.floor((s % 3600) / 60);
    const sec = Math.floor(s % 60);
    return h > 0 ? `${h}:${m.toString().padStart(2, "0")}:${sec.toString().padStart(2, "0")}` : `${m}:${sec.toString().padStart(2, "0")}`;
  };

  return (
    <div
      ref={containerRef}
      className="max-w-5xl mx-auto relative group"
      onMouseMove={resetHideTimer}
      onMouseLeave={() => playing && setShowControls(false)}
    >
      <video
        ref={videoRef}
        key={streamUrl}
        className="w-full cursor-pointer"
        preload="metadata"
        onClick={togglePlay}
        onDoubleClick={toggleFullscreen}
        onPlay={() => {
          setPlaying(true);
          if (!playTriggered.current) { playTriggered.current = true; onPlay(); }
        }}
        onPause={() => setPlaying(false)}
        onTimeUpdate={() => { const t = videoRef.current?.currentTime ?? 0; setCurTime(t); onTimeUpdateProp?.(t); }}
        onProgress={() => {
          const v = videoRef.current;
          if (v && v.buffered.length > 0) setBuffered(v.buffered.end(v.buffered.length - 1));
        }}
        onEnded={() => {
          setPlaying(false);
          scenes.saveActivity(sceneId, { resumeTime: 0 }).catch(() => {});
        }}
      >
        <source src={streamUrl} type={`video/${format || "mp4"}`} />
      </video>

      {/* Custom Controls Overlay */}
      <div
        className={`absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/90 via-black/50 to-transparent transition-opacity ${
          showControls ? "opacity-100" : "opacity-0 pointer-events-none"
        }`}
        style={{ padding: "40px 0 0 0" }}
      >
        {/* Seek bar */}
        <div className="px-3">
          <div className="relative h-4 flex items-center cursor-pointer group/seek" onClick={seekTo}>
            <div className="w-full h-1 bg-white/20 rounded-full group-hover/seek:h-1.5 transition-all relative">
              {/* Buffered */}
              <div className="absolute top-0 left-0 h-full bg-white/30 rounded-full" style={{ width: `${(buffered / (duration || 1)) * 100}%` }} />
              {/* Progress */}
              <div className="absolute top-0 left-0 h-full bg-plex-accent rounded-full" style={{ width: `${(currentTime / (duration || 1)) * 100}%` }} />
              {/* Marker dots */}
              {markers.map((m) => (
                <div
                  key={m.id}
                  className="absolute top-1/2 -translate-y-1/2 w-2 h-2 bg-yellow-400 rounded-full cursor-pointer hover:scale-150 transition-transform z-10"
                  style={{ left: `${(m.seconds / (duration || 1)) * 100}%` }}
                  title={m.title || m.primaryTagName}
                  onClick={(e) => {
                    e.stopPropagation();
                    const v = videoRef.current;
                    if (v) { v.currentTime = m.seconds; v.play().catch(() => {}); }
                  }}
                />
              ))}
            </div>
            {/* Seek thumb */}
            <div
              className="absolute top-1/2 -translate-y-1/2 w-3 h-3 bg-plex-accent rounded-full opacity-0 group-hover/seek:opacity-100 transition-opacity"
              style={{ left: `${(currentTime / (duration || 1)) * 100}%`, transform: "translate(-50%, -50%)" }}
            />
          </div>
        </div>

        {/* Controls row */}
        <div className="flex items-center gap-2 px-3 py-2 text-white">
          <button onClick={togglePlay} className="hover:text-plex-accent p-1">
            {playing ? <Pause className="w-5 h-5" /> : <Play className="w-5 h-5" />}
          </button>

          <button onClick={() => { const v = videoRef.current; if (v) v.currentTime = Math.max(0, v.currentTime - 10); }} className="hover:text-plex-accent p-1" title="Back 10s">
            <SkipBack className="w-4 h-4" />
          </button>
          <button onClick={() => { const v = videoRef.current; if (v) v.currentTime = Math.min(v.duration, v.currentTime + 10); }} className="hover:text-plex-accent p-1" title="Forward 10s">
            <SkipForward className="w-4 h-4" />
          </button>

          {/* Volume */}
          <button onClick={() => {
            const v = videoRef.current;
            if (!v) return;
            v.muted = !v.muted;
            setMuted(v.muted);
            localStorage.setItem(MUTED_KEY, String(v.muted));
          }} className="hover:text-plex-accent p-1">
            {muted || vol === 0 ? <VolumeX className="w-4 h-4" /> : <Volume2 className="w-4 h-4" />}
          </button>
          <div className="w-20 h-3 flex items-center cursor-pointer group/vol" onClick={changeVolume}>
            <div className="w-full h-1 bg-white/20 rounded-full relative">
              <div className="absolute top-0 left-0 h-full bg-white rounded-full" style={{ width: `${(muted ? 0 : vol) * 100}%` }} />
            </div>
          </div>

          <span className="text-xs text-white/70 ml-1 select-none tabular-nums">
            {fmtTime(currentTime)} / {fmtTime(duration)}
          </span>

          <div className="ml-auto flex items-center gap-2">
            {/* Playback speed */}
            <div className="relative">
              <button
                onClick={() => setShowSpeed(!showSpeed)}
                className={`hover:text-plex-accent p-1 text-xs font-medium flex items-center gap-1 ${rate !== 1 ? "text-plex-accent" : ""}`}
              >
                <Gauge className="w-4 h-4" />
                {rate !== 1 ? `${rate}x` : ""}
              </button>
              {showSpeed && (
                <div className="absolute bottom-full right-0 mb-2 bg-gray-900 border border-gray-700 rounded shadow-lg py-1 z-10">
                  {PLAYBACK_RATES.map((r) => (
                    <button
                      key={r}
                      onClick={() => changeRate(r)}
                      className={`block w-full text-left px-4 py-1 text-sm hover:bg-gray-800 ${r === rate ? "text-plex-accent" : "text-white"}`}
                    >
                      {r}x
                    </button>
                  ))}
                </div>
              )}
            </div>

            <button onClick={toggleFullscreen} className="hover:text-plex-accent p-1">
              {fullscreen ? <Minimize className="w-4 h-4" /> : <Maximize className="w-4 h-4" />}
            </button>
          </div>
        </div>
      </div>

      {/* Big play button overlay when paused */}
      {!playing && (
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
          <div className="bg-black/40 rounded-full p-4">
            <Play className="w-12 h-12 text-white" />
          </div>
        </div>
      )}
    </div>
  );
}

// Scene Scrubber / Timeline Component
function SceneScrubber({ 
  sceneId, 
  duration, 
  markers,
  onSeek,
  currentTime,
}: { 
  sceneId: number; 
  duration: number; 
  markers: { id: number; title: string; seconds: number; primaryTagName: string }[];
  onSeek?: (time: number) => void;
  currentTime?: number;
}) {
  const containerRef = useRef<HTMLDivElement>(null);
  const scrollRef = useRef<HTMLDivElement>(null);
  const [spriteData, setSpriteData] = useState<{ entries: { start: number; end: number; x: number; y: number; w: number; h: number }[]; imageUrl: string } | null>(null);
  const [spriteError, setSpriteError] = useState(false);
  
  const spriteVttUrl = `/api/stream/scene/${sceneId}/vtt/thumbs`;
  const spriteImageUrl = `/api/stream/scene/${sceneId}/sprite`;
  const screenshotUrl = `/api/stream/scene/${sceneId}/screenshot`;
  
  const formatTime = (s: number) => {
    const m = Math.floor(s / 60);
    const sec = Math.floor(s % 60);
    return `${m}:${sec.toString().padStart(2, "0")}`;
  };

  // Load and parse VTT sprite data
  useEffect(() => {
    fetch(spriteVttUrl)
      .then(r => { if (!r.ok) throw new Error("VTT not found"); return r.text(); })
      .then(text => {
        const entries: typeof spriteData extends null ? never : NonNullable<typeof spriteData>["entries"] = [];
        const blocks = text.split(/\n\n+/);
        for (const block of blocks) {
          const lines = block.trim().split("\n");
          for (let i = 0; i < lines.length; i++) {
            const timeMatch = lines[i].match(/(\d{2}:\d{2}:\d{2}\.\d{3})\s*-->\s*(\d{2}:\d{2}:\d{2}\.\d{3})/);
            if (timeMatch && lines[i + 1]) {
              const xywhMatch = lines[i + 1].match(/#xywh=(\d+),(\d+),(\d+),(\d+)/);
              if (xywhMatch) {
                entries.push({
                  start: parseVttTime(timeMatch[1]),
                  end: parseVttTime(timeMatch[2]),
                  x: parseInt(xywhMatch[1]),
                  y: parseInt(xywhMatch[2]),
                  w: parseInt(xywhMatch[3]),
                  h: parseInt(xywhMatch[4]),
                });
              }
            }
          }
        }
        if (entries.length > 0) {
          setSpriteData({ entries, imageUrl: spriteImageUrl });
        } else {
          setSpriteError(true);
        }
      })
      .catch(() => setSpriteError(true));
  }, [sceneId, spriteVttUrl, spriteImageUrl]);

  const thumbCount = spriteData ? spriteData.entries.length : Math.min(Math.ceil(duration / 10), 60);
  const thumbWidth = 160;
  const thumbHeight = spriteData?.entries[0] ? Math.round(thumbWidth * (spriteData.entries[0].h / spriteData.entries[0].w)) : 90;

  // Determine which thumbnail index is active based on current video time
  const activeIndex = useMemo(() => {
    if (currentTime == null || currentTime <= 0) return -1;
    if (spriteData) {
      for (let i = spriteData.entries.length - 1; i >= 0; i--) {
        if (currentTime >= spriteData.entries[i].start) return i;
      }
      return 0;
    }
    // Fallback: evenly-spaced thumbs
    const interval = duration / thumbCount;
    return Math.min(Math.floor(currentTime / interval), thumbCount - 1);
  }, [currentTime, spriteData, duration, thumbCount]);

  // Auto-scroll to active thumbnail
  useEffect(() => {
    if (activeIndex >= 0 && scrollRef.current) {
      const targetLeft = activeIndex * thumbWidth;
      const { scrollLeft, clientWidth } = scrollRef.current;
      if (targetLeft < scrollLeft || targetLeft + thumbWidth > scrollLeft + clientWidth) {
        scrollRef.current.scrollTo({ left: Math.max(0, targetLeft - clientWidth / 2 + thumbWidth / 2), behavior: "smooth" });
      }
    }
  }, [activeIndex, thumbWidth]);

  const scroll = (dir: number) => {
    if (scrollRef.current) scrollRef.current.scrollBy({ left: dir * thumbWidth * 4, behavior: "smooth" });
  };
  
  return (
    <div className="flex-shrink-0 bg-[#1a1a1a] border-t border-plex-border">
      {/* Markers bar */}
      {markers.length > 0 && (
        <div className="relative h-5 bg-[#333]">
          {markers.map((marker) => (
            <div
              key={marker.id}
              className="absolute top-0 h-full px-2 bg-plex-accent/80 text-[10px] text-white flex items-center whitespace-nowrap cursor-pointer hover:bg-plex-accent"
              style={{ 
                left: `${(marker.seconds / duration) * 100}%`,
                transform: "translateX(-50%)"
              }}
              title={`${marker.title} - ${marker.primaryTagName}`}
              onClick={() => onSeek?.(marker.seconds)}
            >
              {marker.title || marker.primaryTagName}
            </div>
          ))}
        </div>
      )}
      
      {/* Thumbnails scrubber - uses sprite sheet if available, falls back to individual screenshots */}
      <div className="relative flex overflow-hidden" ref={containerRef}>
        <button onClick={() => scroll(-1)} className="flex-shrink-0 w-7 bg-[#222] hover:bg-[#333] text-plex-text-muted border-r border-plex-border z-10">
          <ChevronLeft className="w-4 h-4 mx-auto" />
        </button>
        
        <div ref={scrollRef} className="flex-1 flex overflow-x-auto scrollbar-thin scrollbar-thumb-plex-border">
          {Array.from({ length: Math.max(thumbCount, 1) }).map((_, i) => {
            const time = spriteData ? spriteData.entries[i]?.start ?? (i / thumbCount) * duration : (i / thumbCount) * duration;
            const entry = spriteData?.entries[i];
            const isActive = i === activeIndex;
            return (
              <div 
                key={i} 
                className={`flex-shrink-0 relative cursor-pointer hover:ring-2 hover:ring-plex-accent hover:z-10 ${isActive ? "ring-2 ring-plex-accent z-10" : ""}`}
                style={{ width: thumbWidth }}
                onClick={() => onSeek?.(time)}
              >
                <div className="bg-plex-surface" style={{ width: thumbWidth, height: thumbHeight }}>
                  {entry ? (
                    <div
                      style={{
                        width: thumbWidth,
                        height: thumbHeight,
                        backgroundImage: `url(${spriteData!.imageUrl})`,
                        backgroundPosition: `-${entry.x * (thumbWidth / entry.w)}px -${entry.y * (thumbHeight / entry.h)}px`,
                        backgroundSize: `${(spriteData!.entries[0].w * Math.ceil(Math.sqrt(thumbCount))) * (thumbWidth / entry.w)}px auto`,
                      }}
                    />
                  ) : !spriteError ? (
                    <img 
                      src={`${screenshotUrl}?seconds=${Math.floor(time)}`} 
                      alt="" 
                      className="w-full h-full object-cover"
                      loading="lazy"
                      onError={(e) => { (e.target as HTMLImageElement).style.display = 'none'; }}
                    />
                  ) : null}
                </div>
                <div className="absolute bottom-0 left-0 right-0 text-center text-[10px] text-white bg-black/70 py-0.5">
                  {formatTime(time)}
                </div>
              </div>
            );
          })}
        </div>
        
        <button onClick={() => scroll(1)} className="flex-shrink-0 w-7 bg-[#222] hover:bg-[#333] text-plex-text-muted border-l border-plex-border z-10">
          <ChevronRight className="w-4 h-4 mx-auto" />
        </button>
      </div>
    </div>
  );
}

function parseVttTime(timeStr: string): number {
  const parts = timeStr.split(":");
  return parseInt(parts[0]) * 3600 + parseInt(parts[1]) * 60 + parseFloat(parts[2]);
}

// Markers Panel (for Markers tab)
function MarkersPanel({ 
  sceneId, 
  markers,
  onSeek,
}: { 
  sceneId: number; 
  markers: { id: number; title: string; seconds: number; endSeconds?: number; primaryTagId: number; primaryTagName: string }[];
  onSeek?: (time: number) => void;
}) {
  const queryClient = useQueryClient();
  const [adding, setAdding] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [title, setTitle] = useState("");
  const [seconds, setSeconds] = useState(0);
  const [tagSearch, setTagSearch] = useState("");
  const [selectedTagId, setSelectedTagId] = useState<number | null>(null);
  const [selectedTagName, setSelectedTagName] = useState("");

  const { data: tagResults } = useQuery({
    queryKey: ["tags-search", tagSearch],
    queryFn: () => tags.find({ q: tagSearch, perPage: 10 }),
    enabled: tagSearch.length >= 1,
  });

  const createMutation = useMutation({
    mutationFn: (data: SceneMarkerCreate) => scenes.markers.create(sceneId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["scene", sceneId] });
      resetForm();
    },
  });

  const updateMutation = useMutation({
    mutationFn: (data: { id: number; title?: string; seconds?: number; primaryTagId?: number }) =>
      scenes.markers.update(sceneId, data.id, {
        title: data.title,
        seconds: data.seconds,
        primaryTagId: data.primaryTagId,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["scene", sceneId] });
      resetForm();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (markerId: number) => scenes.markers.delete(sceneId, markerId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["scene", sceneId] }),
  });

  const resetForm = () => {
    setAdding(false);
    setEditingId(null);
    setTitle("");
    setSeconds(0);
    setTagSearch("");
    setSelectedTagId(null);
    setSelectedTagName("");
  };

  const startEdit = (marker: { id: number; title: string; seconds: number; primaryTagId: number; primaryTagName: string }) => {
    setAdding(true);
    setEditingId(marker.id);
    setTitle(marker.title || "");
    setSeconds(marker.seconds);
    setTagSearch("");
    setSelectedTagId(marker.primaryTagId);
    setSelectedTagName(marker.primaryTagName);
  };

  const formatTime = (s: number) => {
    const m = Math.floor(s / 60);
    const sec = Math.floor(s % 60);
    return `${m}:${sec.toString().padStart(2, "0")}`;
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <span className="text-sm text-plex-text-secondary">{markers.length} marker{markers.length !== 1 ? "s" : ""}</span>
        <button onClick={() => adding ? resetForm() : setAdding(true)} className="text-plex-accent hover:underline text-sm flex items-center gap-1">
          <Plus className="w-3.5 h-3.5" /> Add
        </button>
      </div>

      {adding && (
        <div className="bg-plex-card border border-plex-border rounded p-3 mb-3 space-y-2">
          <input
            type="text"
            placeholder="Marker title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full bg-plex-input border border-plex-border rounded px-3 py-1.5 text-sm text-plex-text"
          />
          <div className="flex gap-2">
            <input
              type="number"
              placeholder="Seconds"
              value={seconds || ""}
              onChange={(e) => setSeconds(Number(e.target.value))}
              className="w-28 bg-plex-input border border-plex-border rounded px-3 py-1.5 text-sm text-plex-text"
              min={0}
            />
            <div className="relative flex-1">
              <div className="flex items-center bg-plex-input border border-plex-border rounded px-3 py-1.5 text-sm">
                <Search className="w-3.5 h-3.5 text-plex-text-muted mr-2 flex-shrink-0" />
                <input
                  type="text"
                  placeholder={selectedTagName || "Search tag..."}
                  value={tagSearch}
                  onChange={(e) => { setTagSearch(e.target.value); setSelectedTagId(null); setSelectedTagName(""); }}
                  className="bg-transparent w-full outline-none text-plex-text"
                />
              </div>
              {tagSearch && tagResults && tagResults.items.length > 0 && (
                <div className="absolute z-10 mt-1 w-full bg-plex-card border border-plex-border rounded shadow-lg max-h-40 overflow-y-auto">
                  {tagResults.items.map((t: { id: number; name: string }) => (
                    <button
                      key={t.id}
                      onClick={() => { setSelectedTagId(t.id); setSelectedTagName(t.name); setTagSearch(""); }}
                      className="block w-full text-left px-3 py-1.5 text-sm text-plex-text-secondary hover:bg-plex-card-hover hover:text-plex-text"
                    >
                      {t.name}
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>
          <div className="flex gap-2 justify-end">
            <button onClick={resetForm} className="px-3 py-1 text-sm text-plex-text-secondary hover:text-plex-text">Cancel</button>
            <button
              onClick={() => selectedTagId && (editingId
                ? updateMutation.mutate({ id: editingId, title, seconds, primaryTagId: selectedTagId })
                : createMutation.mutate({ title, seconds, primaryTagId: selectedTagId }))}
              disabled={!selectedTagId || createMutation.isPending || updateMutation.isPending}
              className="px-3 py-1 text-sm bg-plex-accent hover:bg-plex-accent-hover text-white rounded disabled:opacity-50"
            >
              {editingId ? "Update" : "Save"}
            </button>
          </div>
        </div>
      )}

      {markers.length === 0 && !adding && (
        <p className="text-plex-text-muted text-sm">No markers yet. Click Add to create one.</p>
      )}

      <div className="space-y-1">
        {markers.map((m) => (
          <div key={m.id} className="flex items-center justify-between bg-plex-card border border-plex-border rounded px-3 py-2 text-sm group">
            <button
              className="flex items-center gap-3 hover:text-plex-accent transition-colors"
              onClick={() => onSeek?.(m.seconds)}
              title="Seek to marker"
            >
              <span className="text-plex-accent font-mono text-xs w-12">{formatTime(m.seconds)}</span>
              <span className="text-plex-text group-hover:text-plex-accent">{m.title || "Untitled"}</span>
              <span className="text-xs text-plex-text-secondary bg-plex-surface px-1.5 py-0.5 rounded">{m.primaryTagName}</span>
            </button>
            <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <button
                onClick={() => startEdit(m)}
                className="text-plex-text-muted hover:text-plex-accent"
                title="Edit marker"
              >
                <Pencil className="w-3.5 h-3.5" />
              </button>
              <button
                onClick={() => deleteMutation.mutate(m.id)}
                className="text-plex-text-muted hover:text-red-400"
                title="Delete marker"
              >
                <Trash2 className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// ===== Inline Scene Edit Panel =====
function SceneEditPanel({ scene, onSaved }: { scene: Scene; onSaved: () => void }) {
  const queryClient = useQueryClient();
  const [title, setTitle] = useState(scene.title || "");
  const [code, setCode] = useState(scene.code || "");
  const [details, setDetails] = useState(scene.details || "");
  const [director, setDirector] = useState(scene.director || "");
  const [date, setDate] = useState(scene.date || "");
  const [rating, setRating] = useState<number | undefined>(scene.rating ?? undefined);
  const [organized, setOrganized] = useState(scene.organized);
  const [urls, setUrls] = useState(scene.urls.join("\n"));
  const [studioId, setStudioId] = useState<number | undefined>(scene.studioId ?? undefined);
  const [selectedTagIds, setSelectedTagIds] = useState<number[]>(scene.tags.map((t) => t.id));
  const [selectedPerformerIds, setSelectedPerformerIds] = useState<number[]>(scene.performers.map((p) => p.id));
  const [selectedGalleryIds, setSelectedGalleryIds] = useState<number[]>(scene.galleries.map((g) => g.id));
  const [selectedGroups, setSelectedGroups] = useState<{ groupId: number; sceneIndex: number }[]>(
    scene.groups.map((g) => ({ groupId: g.id, sceneIndex: g.sceneIndex }))
  );
  const [tagSearch, setTagSearch] = useState("");
  const [perfSearch, setPerfSearch] = useState("");
  const [gallerySearch, setGallerySearch] = useState("");
  const [groupSearch, setGroupSearch] = useState("");

  const { data: allTags } = useQuery({ queryKey: ["tags-all"], queryFn: () => tags.find({ perPage: 500, sort: "name", direction: "asc" }) });
  const { data: allPerformers } = useQuery({ queryKey: ["performers-all"], queryFn: () => performersApi.find({ perPage: 500, sort: "name", direction: "asc" }) });
  const { data: allStudios } = useQuery({ queryKey: ["studios-all"], queryFn: () => studiosApi.find({ perPage: 500, sort: "name", direction: "asc" }) });
  const { data: allGalleries } = useQuery({ queryKey: ["galleries-all"], queryFn: () => galleriesApi.find({ perPage: 500, sort: "title", direction: "asc" }) });
  const { data: allGroups } = useQuery({ queryKey: ["groups-all"], queryFn: () => groupsApi.find({ perPage: 500, sort: "name", direction: "asc" }) });

  useEffect(() => {
    setTitle(scene.title || ""); setCode(scene.code || ""); setDetails(scene.details || "");
    setDirector(scene.director || ""); setDate(scene.date || ""); setRating(scene.rating ?? undefined);
    setOrganized(scene.organized); setUrls(scene.urls.join("\n")); setStudioId(scene.studioId ?? undefined);
    setSelectedTagIds(scene.tags.map((t) => t.id)); setSelectedPerformerIds(scene.performers.map((p) => p.id));
    setSelectedGalleryIds(scene.galleries.map((g) => g.id));
    setSelectedGroups(scene.groups.map((g) => ({ groupId: g.id, sceneIndex: g.sceneIndex })));
  }, [scene]);

  const mutation = useMutation({
    mutationFn: (data: SceneUpdate) => scenes.update(scene.id, data),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ["scene", scene.id] }); queryClient.invalidateQueries({ queryKey: ["scenes"] }); onSaved(); },
  });

  const handleSave = () => {
    const urlList = urls.split("\n").map((u) => u.trim()).filter(Boolean);
    mutation.mutate({ title: title || undefined, code: code || undefined, details: details || undefined,
      director: director || undefined, date: date || undefined, rating, organized, studioId,
      urls: urlList, tagIds: selectedTagIds, performerIds: selectedPerformerIds, galleryIds: selectedGalleryIds, groups: selectedGroups });
  };

  const filteredTags = allTags?.items.filter((t) => !selectedTagIds.includes(t.id) && t.name.toLowerCase().includes(tagSearch.toLowerCase())) ?? [];
  const filteredPerformers = allPerformers?.items.filter((p) => !selectedPerformerIds.includes(p.id) && p.name.toLowerCase().includes(perfSearch.toLowerCase())) ?? [];
  const filteredGalleries = allGalleries?.items.filter((g) => !selectedGalleryIds.includes(g.id) && (g.title || "").toLowerCase().includes(gallerySearch.toLowerCase())) ?? [];
  const selectedGroupIds = selectedGroups.map((g) => g.groupId);
  const filteredGroupsList = allGroups?.items.filter((g) => !selectedGroupIds.includes(g.id) && g.name.toLowerCase().includes(groupSearch.toLowerCase())) ?? [];
  const selectedTags = allTags?.items.filter((t) => selectedTagIds.includes(t.id)) ?? scene.tags;
  const selectedPerformers = allPerformers?.items.filter((p) => selectedPerformerIds.includes(p.id)) ?? scene.performers.map((p) => ({ ...p }));
  const selectedGalleries = allGalleries?.items.filter((g) => selectedGalleryIds.includes(g.id)) ?? scene.galleries;

  const inputCls = "w-full bg-plex-surface border border-plex-border rounded px-3 py-2 text-sm text-plex-text focus:outline-none focus:border-plex-accent";

  return (
    <div className="space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <label className="space-y-1"><span className="text-xs text-plex-text-secondary">Title</span><input value={title} onChange={(e) => setTitle(e.target.value)} className={inputCls} /></label>
        <label className="space-y-1"><span className="text-xs text-plex-text-secondary">Date</span><input type="date" value={date} onChange={(e) => setDate(e.target.value)} className={inputCls} /></label>
      </div>
      <div className="grid grid-cols-2 gap-3">
        <label className="space-y-1"><span className="text-xs text-plex-text-secondary">Code</span><input value={code} onChange={(e) => setCode(e.target.value)} className={inputCls} /></label>
        <label className="space-y-1"><span className="text-xs text-plex-text-secondary">Director</span><input value={director} onChange={(e) => setDirector(e.target.value)} className={inputCls} /></label>
      </div>
      <label className="block space-y-1"><span className="text-xs text-plex-text-secondary">Details</span><textarea value={details} onChange={(e) => setDetails(e.target.value)} rows={3} className={inputCls} /></label>
      <div className="grid grid-cols-2 gap-3">
        <RatingField value={rating} onChange={setRating} />
        <label className="space-y-1">
          <span className="text-xs text-plex-text-secondary">Studio</span>
          <select value={studioId ?? ""} onChange={(e) => setStudioId(e.target.value ? Number(e.target.value) : undefined)} className={inputCls}>
            <option value="">None</option>
            {allStudios?.items.map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
          </select>
        </label>
      </div>
      <label className="block space-y-1"><span className="text-xs text-plex-text-secondary">URLs (one per line)</span><textarea value={urls} onChange={(e) => setUrls(e.target.value)} rows={2} className={inputCls} /></label>

      {/* Tags */}
      <div className="space-y-1">
        <span className="text-xs text-plex-text-secondary">Tags</span>
        <div className="flex flex-wrap gap-1 mb-1">
          {selectedTags.map((t) => <span key={t.id} className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs bg-blue-900 text-blue-300">{t.name}<button onClick={() => setSelectedTagIds(selectedTagIds.filter((id) => id !== t.id))} className="hover:text-white">×</button></span>)}
        </div>
        <input value={tagSearch} onChange={(e) => setTagSearch(e.target.value)} placeholder="Search tags…" className={inputCls} />
        {tagSearch && filteredTags.length > 0 && <div className="max-h-24 overflow-y-auto bg-plex-surface rounded border border-plex-border">{filteredTags.slice(0, 10).map((t) => <button key={t.id} onClick={() => { setSelectedTagIds([...selectedTagIds, t.id]); setTagSearch(""); }} className="block w-full text-left px-3 py-1 text-sm text-plex-text hover:bg-plex-card">{t.name}</button>)}</div>}
      </div>

      {/* Performers */}
      <div className="space-y-1">
        <span className="text-xs text-plex-text-secondary">Performers</span>
        <div className="flex flex-wrap gap-1 mb-1">
          {selectedPerformers.map((p) => <span key={p.id} className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs bg-purple-900 text-purple-300">{p.name}<button onClick={() => setSelectedPerformerIds(selectedPerformerIds.filter((id) => id !== p.id))} className="hover:text-white">×</button></span>)}
        </div>
        <input value={perfSearch} onChange={(e) => setPerfSearch(e.target.value)} placeholder="Search performers…" className={inputCls} />
        {perfSearch && filteredPerformers.length > 0 && <div className="max-h-24 overflow-y-auto bg-plex-surface rounded border border-plex-border">{filteredPerformers.slice(0, 10).map((p) => <button key={p.id} onClick={() => { setSelectedPerformerIds([...selectedPerformerIds, p.id]); setPerfSearch(""); }} className="block w-full text-left px-3 py-1 text-sm text-plex-text hover:bg-plex-card">{p.name}{p.disambiguation ? ` (${p.disambiguation})` : ""}</button>)}</div>}
      </div>

      {/* Galleries */}
      <div className="space-y-1">
        <span className="text-xs text-plex-text-secondary">Galleries</span>
        <div className="flex flex-wrap gap-1 mb-1">
          {selectedGalleries.map((g) => <span key={g.id} className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs bg-emerald-900 text-emerald-300">{g.title || "Untitled"}<button onClick={() => setSelectedGalleryIds(selectedGalleryIds.filter((id) => id !== g.id))} className="hover:text-white">×</button></span>)}
        </div>
        <input value={gallerySearch} onChange={(e) => setGallerySearch(e.target.value)} placeholder="Search galleries…" className={inputCls} />
        {gallerySearch && filteredGalleries.length > 0 && <div className="max-h-24 overflow-y-auto bg-plex-surface rounded border border-plex-border">{filteredGalleries.slice(0, 10).map((g) => <button key={g.id} onClick={() => { setSelectedGalleryIds([...selectedGalleryIds, g.id]); setGallerySearch(""); }} className="block w-full text-left px-3 py-1 text-sm text-plex-text hover:bg-plex-card">{g.title || "Untitled"}</button>)}</div>}
      </div>

      {/* Groups */}
      <div className="space-y-1">
        <span className="text-xs text-plex-text-secondary">Groups</span>
        <div className="space-y-1 mb-1">
          {selectedGroups.map((sg) => {
            const group = allGroups?.items.find((g) => g.id === sg.groupId);
            return (
              <div key={sg.groupId} className="flex items-center gap-2">
                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs bg-orange-900 text-orange-300">
                  {group?.name || `Group #${sg.groupId}`}
                  <button onClick={() => setSelectedGroups(selectedGroups.filter((g) => g.groupId !== sg.groupId))} className="hover:text-white">×</button>
                </span>
                <label className="flex items-center gap-1 text-xs text-plex-text-muted">
                  Scene #
                  <input type="number" min={0} value={sg.sceneIndex}
                    onChange={(e) => setSelectedGroups(selectedGroups.map((g) => g.groupId === sg.groupId ? { ...g, sceneIndex: Number(e.target.value) || 0 } : g))}
                    className="w-16 bg-plex-surface border border-plex-border rounded px-2 py-0.5 text-xs text-plex-text focus:outline-none focus:border-plex-accent" />
                </label>
              </div>
            );
          })}
        </div>
        <input value={groupSearch} onChange={(e) => setGroupSearch(e.target.value)} placeholder="Search groups…" className={inputCls} />
        {groupSearch && filteredGroupsList.length > 0 && <div className="max-h-24 overflow-y-auto bg-plex-surface rounded border border-plex-border">{filteredGroupsList.slice(0, 10).map((g) => <button key={g.id} onClick={() => { setSelectedGroups([...selectedGroups, { groupId: g.id, sceneIndex: 0 }]); setGroupSearch(""); }} className="block w-full text-left px-3 py-1 text-sm text-plex-text hover:bg-plex-card">{g.name}</button>)}</div>}
      </div>

      <label className="flex items-center gap-2 text-sm text-plex-text-secondary cursor-pointer">
        <input type="checkbox" checked={organized} onChange={(e) => setOrganized(e.target.checked)} className="rounded border-plex-border bg-plex-surface" /> Organized
      </label>

      {mutation.error && <div className="bg-red-900/50 border border-red-700 text-red-300 rounded p-2 text-sm">{(mutation.error as Error).message}</div>}

      <div className="flex justify-end gap-3 pt-2">
        <button onClick={onSaved} className="px-4 py-2 text-sm text-plex-text-secondary hover:text-plex-text">Cancel</button>
        <button onClick={handleSave} disabled={mutation.isPending} className="px-4 py-2 text-sm bg-plex-accent hover:bg-plex-accent-hover text-white rounded disabled:opacity-50">
          {mutation.isPending ? "Saving…" : "Save"}
        </button>
      </div>
    </div>
  );
}
