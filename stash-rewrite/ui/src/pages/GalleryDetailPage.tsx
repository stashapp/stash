import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { galleries, images, scenes, entityImages } from "../api/client";
import type { FindFilter, GalleryChapter } from "../api/types";
import { formatDate, formatDuration, getResolutionLabel, TagBadge } from "../components/shared";
import { ArrowLeft, BookOpen, Film, FolderOpen, ImageIcon, Link as LinkIcon, Pencil, Plus, Trash2, UserRound, Check, Loader2 } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { GalleryEditModal } from "./GalleryEditModal";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { ExtensionSlot } from "../router/RouteRegistry";
import { Lightbox, type LightboxImage } from "../components/Lightbox";
import { InteractiveRating } from "../components/Rating";

interface Props {
  id: number;
  onNavigate: (r: any) => void;
}

type TabKey = "images" | "scenes" | "chapters";

export function GalleryDetailPage({ id, onNavigate }: Props) {
  const { data: gallery, isLoading } = useQuery({
    queryKey: ["gallery", id],
    queryFn: () => galleries.get(id),
  });
  const { data: galleryImages } = useQuery({
    queryKey: ["gallery-images", id],
    queryFn: () => images.find({ page: 1, perPage: 60, direction: "desc" }, { galleryId: id }),
    enabled: !!gallery,
  });
  const { data: chaptersData } = useQuery({
    queryKey: ["gallery-chapters", id],
    queryFn: () => galleries.chapters(id),
    enabled: !!gallery,
  });
  const [editing, setEditing] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [activeTab, setActiveTab] = useState<TabKey>("images");
  const [lightboxOpen, setLightboxOpen] = useState(false);
  const [lightboxIndex, setLightboxIndex] = useState(0);
  const [sceneFilter, setSceneFilter] = useState<FindFilter>({ page: 1, perPage: 24, direction: "desc" });
  const [imageSelectMode, setImageSelectMode] = useState(false);
  const [selectedImageIds, setSelectedImageIds] = useState<Set<number>>(new Set());
  const [showAddImages, setShowAddImages] = useState(false);
  const queryClient = useQueryClient();

  useEffect(() => {
    if (gallery) document.title = `${gallery.title || `Gallery ${id}`} | Stash`;
    return () => { document.title = "Stash"; };
  }, [gallery, id]);

  // Keyboard shortcuts
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const el = (e.target as HTMLElement).tagName;
      if (el === "INPUT" || el === "TEXTAREA" || el === "SELECT") return;
      switch (e.key) {
        case "e": setEditing((v) => !v); break;
        case "a": setActiveTab("images"); break;
        case "c": setActiveTab("chapters"); break;
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  const deleteMut = useMutation({
    mutationFn: () => galleries.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["galleries"] });
      onNavigate({ page: "galleries" });
    },
  });

  const galleryUpdateMut = useMutation({
    mutationFn: (data: { rating?: number; organized?: boolean }) => galleries.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["gallery", id] });
      queryClient.invalidateQueries({ queryKey: ["galleries"] });
    },
  });

  const removeImagesMut = useMutation({
    mutationFn: (imageIds: number[]) => galleries.removeImages(id, imageIds),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["gallery-images", id] });
      queryClient.invalidateQueries({ queryKey: ["gallery", id] });
      setSelectedImageIds(new Set());
      setImageSelectMode(false);
    },
  });

  const addImagesMut = useMutation({
    mutationFn: (imageIds: number[]) => galleries.addImages(id, imageIds),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["gallery-images", id] });
      queryClient.invalidateQueries({ queryKey: ["gallery", id] });
      setShowAddImages(false);
    },
  });

  const coverImage = useMemo(() => galleryImages?.items[0], [galleryImages]);

  const lightboxImages: LightboxImage[] = useMemo(
    () => (galleryImages?.items ?? []).map((img) => ({ id: img.id, src: images.imageUrl(img.id), title: img.title })),
    [galleryImages],
  );

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-plex-accent" />
      </div>
    );
  }

  if (!gallery) {
    return <div className="py-16 text-center text-plex-text-secondary">Gallery not found</div>;
  }

  return (
    <div className="min-h-screen">
      <div className="relative overflow-hidden border-b border-plex-border bg-plex-surface">
        {coverImage ? (
          <>
            <img src={images.imageUrl(coverImage.id)} alt="" className="absolute inset-0 h-full w-full object-cover opacity-20 blur-sm" />
            <div className="absolute inset-0 bg-gradient-to-t from-plex-bg via-plex-bg/70 to-plex-bg/30" />
          </>
        ) : (
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,_rgba(204,123,25,0.18),_transparent_30%),linear-gradient(180deg,_rgba(50,54,57,0.95),_rgba(31,35,38,1))]" />
        )}

        <div className="relative mx-auto max-w-7xl px-4 py-8">
          <div className="mb-5 flex items-center justify-between gap-4">
            <button
              onClick={() => onNavigate({ page: "galleries" })}
              className="flex items-center gap-1 text-sm text-plex-text-secondary hover:text-plex-text"
            >
              <ArrowLeft className="h-4 w-4" /> Back to galleries
            </button>
            <div className="flex items-center gap-2">
              <ExtensionSlot slot="gallery-detail-actions" context={{ gallery, onNavigate }} />
              <button
                onClick={() => setEditing(true)}
                className="flex items-center gap-1.5 rounded bg-plex-accent px-3 py-1.5 text-sm text-white hover:bg-plex-accent-hover"
              >
                <Pencil className="h-3.5 w-3.5" /> Edit
              </button>
              <button
                onClick={() => setConfirmDelete(true)}
                className="flex items-center gap-1.5 rounded border border-plex-border bg-plex-card px-3 py-1.5 text-sm text-plex-text-secondary hover:border-red-500 hover:text-red-300"
              >
                <Trash2 className="h-3.5 w-3.5" /> Delete
              </button>
            </div>
          </div>

          <div className="flex flex-col gap-6 md:flex-row md:items-end">
            <div className="flex h-40 w-32 flex-shrink-0 items-center justify-center overflow-hidden rounded-2xl border border-plex-border bg-plex-card shadow-xl shadow-black/35 sm:h-52 sm:w-40">
              {coverImage ? (
                <img src={images.imageUrl(coverImage.id)} alt={gallery.title || ""} className="h-full w-full object-cover" />
              ) : (
                <FolderOpen className="h-14 w-14 text-plex-text-muted" />
              )}
            </div>
            <div className="min-w-0 flex-1">
              <h1 className="truncate text-4xl font-bold text-plex-text">{gallery.title || "Untitled Gallery"}</h1>
              <div className="mt-2 flex flex-wrap items-center gap-3 text-sm text-plex-text-secondary">
                {gallery.date && <span>{formatDate(gallery.date)}</span>}
                {gallery.studioName && gallery.studioId && (
                  <button onClick={() => onNavigate({ page: "studio", id: gallery.studioId })} className="text-plex-accent hover:underline">
                    {gallery.studioName}
                  </button>
                )}
                {gallery.photographer && <span>Photographer: {gallery.photographer}</span>}
                <span className="flex items-center gap-1"><ImageIcon className="h-4 w-4" /> {gallery.imageCount} images</span>
              </div>
              <div className="mt-4 flex flex-wrap items-center justify-between gap-3 rounded-xl border border-plex-border bg-plex-card/60 px-4 py-3">
                <InteractiveRating value={gallery.rating} onChange={(value) => galleryUpdateMut.mutate({ rating: value })} />
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => galleryUpdateMut.mutate({ organized: !gallery.organized })}
                    className={`rounded px-2 py-1 text-xs font-medium transition-colors ${gallery.organized ? "bg-green-600 text-white" : "border border-plex-border bg-plex-card text-plex-text-secondary hover:text-plex-text"}`}
                  >
                    {gallery.organized ? "Organized" : "Mark Organized"}
                  </button>
                </div>
              </div>
              {gallery.details && (
                <p className="mt-3 max-w-4xl whitespace-pre-wrap text-sm leading-6 text-plex-text-secondary">{gallery.details}</p>
              )}
              {gallery.tags.length > 0 && (
                <div className="mt-4 flex flex-wrap gap-1.5">
                  {gallery.tags.map((tag) => (
                    <TagBadge key={tag.id} name={tag.name} onClick={() => onNavigate({ page: "tag", id: tag.id })} />
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Tabs */}
          <div className="mt-6 flex gap-1 border-b border-plex-border">
            {[
              { key: "images" as TabKey, label: "Images", count: gallery.imageCount },
              { key: "scenes" as TabKey, label: "Scenes" },
              { key: "chapters" as TabKey, label: "Chapters", count: chaptersData?.length ?? 0 },
            ].map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`flex items-center gap-2 px-4 py-2.5 text-sm font-medium transition-colors ${
                  activeTab === tab.key
                    ? "border-b-2 border-plex-accent text-plex-accent"
                    : "text-plex-text-secondary hover:text-plex-text"
                }`}
              >
                {tab.label}
                {tab.count !== undefined && (
                  <span className={`min-w-[20px] rounded-full px-1.5 py-0.5 text-center text-xs ${
                    activeTab === tab.key ? "bg-plex-accent/20 text-plex-accent" : "bg-plex-surface text-plex-text-muted"
                  }`}>{tab.count}</span>
                )}
              </button>
            ))}
          </div>
        </div>
      </div>

      <GalleryEditModal gallery={gallery} open={editing} onClose={() => setEditing(false)} />
      <ConfirmDialog
        open={confirmDelete}
        title="Delete Gallery"
        message={`Delete "${gallery.title || "Untitled"}"? This cannot be undone.`}
        onConfirm={() => deleteMut.mutate()}
        onCancel={() => setConfirmDelete(false)}
      />

      <div className="mx-auto max-w-7xl px-4 py-6">
        <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
          <div className="min-w-0">
            <div className="mb-6 rounded-xl border border-plex-border bg-plex-card p-4">
              <dl className="grid gap-y-2 text-sm md:grid-cols-[160px_1fr]">
                <dt className="text-plex-text-muted">Created</dt>
                <dd className="text-plex-text">{formatDate(gallery.createdAt)}</dd>
                <dt className="text-plex-text-muted">Updated</dt>
                <dd className="text-plex-text">{formatDate(gallery.updatedAt)}</dd>
                {gallery.code && (
                  <>
                    <dt className="text-plex-text-muted">Code</dt>
                    <dd className="text-plex-text">{gallery.code}</dd>
                  </>
                )}
                {gallery.photographer && (
                  <>
                    <dt className="text-plex-text-muted">Photographer</dt>
                    <dd className="text-plex-text">{gallery.photographer}</dd>
                  </>
                )}
              </dl>
            </div>

            {gallery.performers.length > 0 && (
              <div className="mb-6 rounded-xl border border-plex-border bg-plex-card p-4">
                <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">Performers</h2>
                <div className="flex flex-wrap justify-center gap-3">
                  {gallery.performers.map((performer) => (
                    <GalleryPerformerCard key={performer.id} performer={performer} onClick={() => onNavigate({ page: "performer", id: performer.id })} />
                  ))}
                </div>
              </div>
            )}

            {activeTab === "images" && (
              galleryImages && galleryImages.items.length > 0 ? (
                <>
                  {/* Image management toolbar */}
                  <div className="flex items-center gap-2 mb-3">
                    <button
                      onClick={() => { setImageSelectMode(!imageSelectMode); setSelectedImageIds(new Set()); }}
                      className={`flex items-center gap-1 px-2 py-1 rounded text-xs border ${imageSelectMode ? "border-plex-accent text-plex-accent" : "border-plex-border text-plex-text-secondary hover:text-plex-text"}`}
                    >
                      <Check className="w-3.5 h-3.5" />
                      {imageSelectMode ? "Cancel" : "Select"}
                    </button>
                    {imageSelectMode && selectedImageIds.size > 0 && (
                      <button
                        onClick={() => { if (confirm(`Remove ${selectedImageIds.size} image(s) from gallery?`)) removeImagesMut.mutate([...selectedImageIds]); }}
                        disabled={removeImagesMut.isPending}
                        className="flex items-center gap-1 px-2 py-1 rounded text-xs text-red-400 hover:text-red-300 hover:bg-red-900/20 border border-red-800"
                      >
                        {removeImagesMut.isPending ? <Loader2 className="w-3 h-3 animate-spin" /> : <Trash2 className="w-3 h-3" />}
                        Remove {selectedImageIds.size}
                      </button>
                    )}
                    <button
                      onClick={() => setShowAddImages(true)}
                      className="flex items-center gap-1 px-2 py-1 rounded text-xs text-blue-400 hover:text-blue-300 hover:bg-blue-900/20 border border-blue-800 ml-auto"
                    >
                      <Plus className="w-3 h-3" />
                      Add Images
                    </button>
                  </div>
                  <div className="grid grid-cols-3 gap-3 sm:grid-cols-4 lg:grid-cols-6">
                    {galleryImages.items.map((image, idx) => (
                      <button
                        key={image.id}
                        onClick={() => {
                          if (imageSelectMode) {
                            setSelectedImageIds((prev) => {
                              const next = new Set(prev);
                              next.has(image.id) ? next.delete(image.id) : next.add(image.id);
                              return next;
                            });
                          } else {
                            setLightboxIndex(idx);
                            setLightboxOpen(true);
                          }
                        }}
                        className={`group overflow-hidden rounded-lg border bg-plex-card text-left shadow-md shadow-black/20 relative ${selectedImageIds.has(image.id) ? "border-plex-accent ring-2 ring-plex-accent" : "border-plex-border"}`}
                      >
                        {(imageSelectMode || selectedImageIds.has(image.id)) && (
                          <div className="absolute top-1 left-1 z-10">
                            <input
                              type="checkbox"
                              checked={selectedImageIds.has(image.id)}
                              onChange={() => {}}
                              className="w-4 h-4 rounded border-plex-border cursor-pointer accent-plex-accent"
                            />
                          </div>
                        )}
                      <div className="aspect-square overflow-hidden bg-plex-surface">
                        <img src={images.thumbnailUrl(image.id)} alt={image.title || ""} className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105" loading="lazy" />
                      </div>
                      <div className="p-2">
                        <p className="truncate text-xs text-plex-text group-hover:text-plex-accent">{image.title || "Untitled"}</p>
                      </div>
                    </button>
                  ))}
                </div>
                </>
              ) : (
                <EmptyPanel icon={<ImageIcon className="h-12 w-12" />} message="No images in this gallery" />
              )
            )}

            {activeTab === "scenes" && (
              <GalleryScenesPanel galleryId={id} filter={sceneFilter} setFilter={setSceneFilter} onNavigate={onNavigate} />
            )}

            {activeTab === "chapters" && (
              <GalleryChaptersPanel galleryId={id} chapters={chaptersData ?? []} />
            )}

            <ExtensionSlot slot="gallery-detail-main-bottom" context={{ gallery, onNavigate }} />
          </div>

          <aside className="space-y-4">
            <div className="rounded-xl border border-plex-border bg-plex-card p-4">
              <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">Metadata</h2>
              <dl className="space-y-2 text-sm">
                {gallery.code && (
                  <div>
                    <dt className="text-plex-text-muted">Code</dt>
                    <dd className="text-plex-text">{gallery.code}</dd>
                  </div>
                )}
                <div>
                  <dt className="text-plex-text-muted">Organized</dt>
                  <dd className="text-plex-text">{gallery.organized ? "Yes" : "No"}</dd>
                </div>
                <div>
                  <dt className="text-plex-text-muted">Created</dt>
                  <dd className="text-plex-text">{formatDate(gallery.createdAt)}</dd>
                </div>
                <div>
                  <dt className="text-plex-text-muted">Updated</dt>
                  <dd className="text-plex-text">{formatDate(gallery.updatedAt)}</dd>
                </div>
              </dl>
            </div>

            {gallery.urls.length > 0 && (
              <div className="rounded-xl border border-plex-border bg-plex-card p-4">
                <h2 className="mb-3 flex items-center gap-1.5 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">
                  <LinkIcon className="h-4 w-4" /> URLs
                </h2>
                <div className="space-y-1 text-sm">
                  {gallery.urls.map((url, index) => (
                    <a key={index} href={url} target="_blank" rel="noopener noreferrer" className="block truncate text-plex-accent hover:underline">
                      {url}
                    </a>
                  ))}
                </div>
              </div>
            )}

            <ExtensionSlot slot="gallery-detail-sidebar-bottom" context={{ gallery, onNavigate }} />
          </aside>
        </div>

        <ExtensionSlot slot="gallery-detail-bottom" context={{ gallery, onNavigate }} />
      </div>

      {/* Add Images Dialog */}
      {showAddImages && (
        <AddImagesDialog
          galleryId={id}
          existingImageIds={new Set(galleryImages?.items.map((i) => i.id) ?? [])}
          onAdd={(ids) => addImagesMut.mutate(ids)}
          onClose={() => setShowAddImages(false)}
          isPending={addImagesMut.isPending}
        />
      )}

      <Lightbox
        images={lightboxImages}
        initialIndex={lightboxIndex}
        open={lightboxOpen}
        onClose={() => setLightboxOpen(false)}
      />
    </div>
  );
}

function GalleryPerformerCard({ performer, onClick }: { performer: { id: number; name: string; disambiguation?: string; imagePath?: string }; onClick: () => void }) {
  const imageUrl = performer.imagePath || entityImages.performerImageUrl(performer.id);

  return (
    <button
      onClick={onClick}
      className="w-[180px] overflow-hidden rounded-lg border border-plex-border bg-plex-surface text-left transition-colors hover:border-plex-accent/60"
    >
      <div className="aspect-[2/3] overflow-hidden bg-plex-card">
        <img
          src={imageUrl}
          alt={performer.name}
          className="h-full w-full object-cover"
          onError={(e) => {
            (e.target as HTMLImageElement).style.display = "none";
            const fallback = (e.target as HTMLImageElement).nextElementSibling as HTMLElement | null;
            if (fallback) fallback.style.display = "flex";
          }}
        />
        <div className="hidden h-full w-full items-center justify-center bg-gradient-to-b from-plex-card to-plex-surface">
          <UserRound className="h-12 w-12 text-plex-text-muted/50" />
        </div>
      </div>
      <div className="p-2 text-center">
        <p className="truncate text-sm font-medium text-plex-text">{performer.name}</p>
        {performer.disambiguation && <p className="truncate text-xs text-plex-text-muted">{performer.disambiguation}</p>}
      </div>
    </button>
  );
}

function GalleryScenesPanel({ galleryId, filter, setFilter, onNavigate }: {
  galleryId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["gallery-scenes", galleryId, filter],
    queryFn: () => scenes.find(filter, { galleryId: String(galleryId) }),
  });

  if (isLoading) return <LoadingPanel icon={<Film className="h-10 w-10" />} message="Loading scenes..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<Film className="h-12 w-12" />} message="No scenes for this gallery" />;

  return (
    <>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
        {data.items.map((scene) => (
          <button key={scene.id} onClick={() => onNavigate({ page: "scene", id: scene.id })} className="group text-left">
            <div className="relative aspect-video overflow-hidden rounded-lg border border-plex-border bg-plex-card shadow-md shadow-black/30">
              <img src={scenes.screenshotUrl(scene.id)} alt={scene.title || ""} className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105" loading="lazy" />
              {scene.files[0]?.duration > 0 && (
                <span className="absolute bottom-1.5 right-1.5 rounded bg-black/75 px-1.5 py-0.5 text-[11px] text-white">{formatDuration(scene.files[0].duration)}</span>
              )}
              {scene.files[0] && getResolutionLabel(scene.files[0].width, scene.files[0].height) && (
                <span className="absolute top-1.5 right-1.5 rounded bg-black/75 px-1.5 py-0.5 text-[10px] font-bold uppercase text-plex-accent">
                  {getResolutionLabel(scene.files[0].width, scene.files[0].height)}
                </span>
              )}
            </div>
            <div className="pt-2">
              <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">{scene.title || "Untitled"}</p>
              <p className="mt-0.5 truncate text-xs text-plex-text-secondary">{scene.date || scene.studioName || ""}</p>
            </div>
          </button>
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={data.totalCount} />
    </>
  );
}

function GalleryChaptersPanel({ galleryId, chapters }: { galleryId: number; chapters: GalleryChapter[] }) {
  const queryClient = useQueryClient();
  const [adding, setAdding] = useState(false);
  const [newTitle, setNewTitle] = useState("");
  const [newIndex, setNewIndex] = useState(0);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editTitle, setEditTitle] = useState("");
  const [editIndex, setEditIndex] = useState(0);

  const createMut = useMutation({
    mutationFn: () => galleries.createChapter(galleryId, { title: newTitle, imageIndex: newIndex }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["gallery-chapters", galleryId] });
      setAdding(false);
      setNewTitle("");
      setNewIndex(0);
    },
  });

  const updateMut = useMutation({
    mutationFn: (chapterId: number) => galleries.updateChapter(galleryId, chapterId, { title: editTitle, imageIndex: editIndex }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["gallery-chapters", galleryId] });
      setEditingId(null);
    },
  });

  const deleteMut = useMutation({
    mutationFn: (chapterId: number) => galleries.deleteChapter(galleryId, chapterId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["gallery-chapters", galleryId] });
    },
  });

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold uppercase tracking-wide text-plex-text-muted">Chapters</h3>
        <button
          onClick={() => setAdding(true)}
          className="flex items-center gap-1.5 rounded bg-plex-accent px-3 py-1.5 text-sm text-white hover:bg-plex-accent-hover"
        >
          <Plus className="h-3.5 w-3.5" /> Add Chapter
        </button>
      </div>

      {adding && (
        <div className="rounded-xl border border-plex-accent/40 bg-plex-card p-4">
          <div className="grid gap-3 sm:grid-cols-[1fr_auto_auto_auto]">
            <input
              type="text"
              placeholder="Chapter title"
              value={newTitle}
              onChange={(e) => setNewTitle(e.target.value)}
              className="rounded border border-plex-border bg-plex-surface px-3 py-1.5 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
            />
            <input
              type="number"
              placeholder="Image index"
              value={newIndex}
              onChange={(e) => setNewIndex(parseInt(e.target.value) || 0)}
              className="w-24 rounded border border-plex-border bg-plex-surface px-3 py-1.5 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
            />
            <button onClick={() => createMut.mutate()} disabled={!newTitle} className="rounded bg-plex-accent px-3 py-1.5 text-sm text-white hover:bg-plex-accent-hover disabled:opacity-50">Save</button>
            <button onClick={() => setAdding(false)} className="rounded border border-plex-border px-3 py-1.5 text-sm text-plex-text-secondary hover:text-plex-text">Cancel</button>
          </div>
        </div>
      )}

      {chapters.length === 0 ? (
        <EmptyPanel icon={<BookOpen className="h-12 w-12" />} message="No chapters" />
      ) : (
        <div className="divide-y divide-plex-border rounded-xl border border-plex-border bg-plex-card">
          {chapters.map((ch) => (
            <div key={ch.id} className="flex items-center justify-between px-4 py-3">
              {editingId === ch.id ? (
                <div className="flex flex-1 items-center gap-3">
                  <input
                    type="text"
                    value={editTitle}
                    onChange={(e) => setEditTitle(e.target.value)}
                    className="flex-1 rounded border border-plex-border bg-plex-surface px-3 py-1 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
                  />
                  <input
                    type="number"
                    value={editIndex}
                    onChange={(e) => setEditIndex(parseInt(e.target.value) || 0)}
                    className="w-20 rounded border border-plex-border bg-plex-surface px-3 py-1 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
                  />
                  <button onClick={() => updateMut.mutate(ch.id)} className="text-sm text-plex-accent hover:underline">Save</button>
                  <button onClick={() => setEditingId(null)} className="text-sm text-plex-text-secondary hover:text-plex-text">Cancel</button>
                </div>
              ) : (
                <>
                  <div>
                    <p className="text-sm font-medium text-plex-text">{ch.title}</p>
                    <p className="text-xs text-plex-text-secondary">Image #{ch.imageIndex}</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => {
                        setEditingId(ch.id);
                        setEditTitle(ch.title);
                        setEditIndex(ch.imageIndex);
                      }}
                      className="rounded p-1 text-plex-text-muted hover:text-plex-accent"
                    >
                      <Pencil className="h-3.5 w-3.5" />
                    </button>
                    <button
                      onClick={() => deleteMut.mutate(ch.id)}
                      className="rounded p-1 text-plex-text-muted hover:text-red-400"
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </button>
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function AddImagesDialog({ galleryId, existingImageIds, onAdd, onClose, isPending }: {
  galleryId: number;
  existingImageIds: Set<number>;
  onAdd: (ids: number[]) => void;
  onClose: () => void;
  isPending: boolean;
}) {
  const [searchFilter, setSearchFilter] = useState<FindFilter>({ page: 1, perPage: 30, direction: "desc" });
  const [selected, setSelected] = useState<Set<number>>(new Set());

  const { data } = useQuery({
    queryKey: ["images-for-gallery", searchFilter],
    queryFn: () => images.find(searchFilter),
  });

  const allImages = data?.items ?? [];
  const available = allImages.filter((i) => !existingImageIds.has(i.id));

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70" onClick={onClose}>
      <div className="bg-plex-card border border-plex-border rounded-xl shadow-2xl w-full max-w-4xl max-h-[80vh] flex flex-col" onClick={(e) => e.stopPropagation()}>
        <div className="flex items-center justify-between px-5 py-4 border-b border-plex-border">
          <h2 className="text-lg font-semibold text-plex-text">Add Images to Gallery</h2>
          <div className="flex items-center gap-3">
            <span className="text-xs text-plex-text-muted">{selected.size} selected</span>
            <button
              onClick={() => onAdd([...selected])}
              disabled={selected.size === 0 || isPending}
              className="px-3 py-1.5 rounded text-sm font-medium bg-plex-accent hover:bg-plex-accent-hover text-white disabled:opacity-50 flex items-center gap-2"
            >
              {isPending && <Loader2 className="w-3.5 h-3.5 animate-spin" />}
              Add {selected.size > 0 ? selected.size : ""}
            </button>
          </div>
        </div>

        <div className="px-5 py-3 border-b border-plex-border">
          <input
            type="text"
            placeholder="Search images..."
            value={searchFilter.q ?? ""}
            onChange={(e) => setSearchFilter((f) => ({ ...f, q: e.target.value || undefined, page: 1 }))}
            className="w-full bg-plex-input border border-plex-border rounded px-3 py-2 text-sm text-plex-text focus:outline-none focus:border-plex-accent"
          />
        </div>

        <div className="flex-1 overflow-y-auto p-5">
          {available.length > 0 ? (
            <div className="grid grid-cols-4 sm:grid-cols-5 lg:grid-cols-6 gap-3">
              {available.map((image) => (
                <button
                  key={image.id}
                  onClick={() => setSelected((prev) => { const n = new Set(prev); n.has(image.id) ? n.delete(image.id) : n.add(image.id); return n; })}
                  className={`group overflow-hidden rounded-lg border text-left relative ${selected.has(image.id) ? "border-plex-accent ring-2 ring-plex-accent" : "border-plex-border"}`}
                >
                  {selected.has(image.id) && (
                    <div className="absolute top-1 left-1 z-10">
                      <div className="w-5 h-5 rounded bg-plex-accent flex items-center justify-center">
                        <Check className="w-3 h-3 text-white" />
                      </div>
                    </div>
                  )}
                  <div className="aspect-square overflow-hidden bg-plex-surface">
                    <img src={images.thumbnailUrl(image.id)} alt={image.title || ""} className="h-full w-full object-cover" loading="lazy" />
                  </div>
                  <div className="p-1.5">
                    <p className="truncate text-xs text-plex-text">{image.title || "Untitled"}</p>
                  </div>
                </button>
              ))}
            </div>
          ) : (
            <div className="text-center py-12 text-plex-text-muted">No images available to add</div>
          )}
        </div>

        <div className="flex items-center justify-between px-5 py-3 border-t border-plex-border">
          <div className="flex items-center gap-2">
            <button onClick={() => setSearchFilter((f) => ({ ...f, page: Math.max(1, (f.page ?? 1) - 1) }))} disabled={(searchFilter.page ?? 1) <= 1} className="px-2 py-1 rounded text-xs text-plex-text-secondary hover:text-plex-text disabled:opacity-30">Prev</button>
            <span className="text-xs text-plex-text-muted">Page {searchFilter.page ?? 1}</span>
            <button onClick={() => setSearchFilter((f) => ({ ...f, page: (f.page ?? 1) + 1 }))} className="px-2 py-1 rounded text-xs text-plex-text-secondary hover:text-plex-text">Next</button>
          </div>
          <button onClick={onClose} className="px-3 py-1.5 rounded text-sm text-plex-text-secondary hover:text-plex-text">Cancel</button>
        </div>
      </div>
    </div>
  );
}

function Pager({ filter, setFilter, totalCount }: {
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  totalCount: number;
}) {
  const perPage = filter.perPage ?? 1;
  const page = filter.page ?? 1;
  const totalPages = Math.max(1, Math.ceil(totalCount / perPage));
  if (totalPages <= 1) return null;
  return (
    <div className="mt-6 flex items-center justify-center gap-4">
      <button disabled={page <= 1} onClick={() => setFilter({ ...filter, page: page - 1 })} className="rounded border border-plex-border bg-plex-card px-4 py-2 text-sm text-plex-text-secondary hover:bg-plex-card-hover disabled:cursor-not-allowed disabled:opacity-50">Previous</button>
      <span className="text-sm text-plex-text-secondary">Page {page} of {totalPages}</span>
      <button disabled={page >= totalPages} onClick={() => setFilter({ ...filter, page: page + 1 })} className="rounded border border-plex-border bg-plex-card px-4 py-2 text-sm text-plex-text-secondary hover:bg-plex-card-hover disabled:cursor-not-allowed disabled:opacity-50">Next</button>
    </div>
  );
}

function LoadingPanel({ icon, message }: { icon: React.ReactNode; message: string }) {
  return (
    <div className="flex flex-col items-center justify-center py-12 text-plex-text-muted">
      <div className="mb-3 animate-pulse">{icon}</div>
      <p>{message}</p>
    </div>
  );
}

function EmptyPanel({ icon, message }: { icon: React.ReactNode; message: string }) {
  return (
    <div className="rounded-xl border border-dashed border-plex-border bg-plex-card/40 py-12 text-center text-plex-text-muted">
      <div className="mx-auto mb-3 flex justify-center opacity-60">{icon}</div>
      <p>{message}</p>
    </div>
  );
}