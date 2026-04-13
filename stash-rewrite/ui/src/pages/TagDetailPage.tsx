import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { galleries, groups, images, markers, metadata, performers, scenes, studios, tags, entityImages } from "../api/client";
import type { FindFilter, Gallery, Group, Image, Performer, Scene, SceneMarkerWall, Studio } from "../api/types";
import { formatDate, formatDuration, getResolutionLabel, TagBadge } from "../components/shared";
import { ArrowLeft, Bookmark, Building2, Film, FolderOpen, GitMerge, Heart, ImageIcon, Layers, Pencil, Tag as TagIcon, Trash2, UserRound, Wand2 } from "lucide-react";
import { useEffect, useState } from "react";
import { TagEditModal } from "./TagEditModal";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { MergeDialog } from "../components/MergeDialog";
import { ExtensionSlot } from "../router/RouteRegistry";

interface Props {
  id: number;
  onNavigate: (r: any) => void;
}

type TabKey = "scenes" | "performers" | "images" | "galleries" | "markers" | "studios" | "groups";

export function TagDetailPage({ id, onNavigate }: Props) {
  const { data: tag, isLoading } = useQuery({
    queryKey: ["tag", id],
    queryFn: () => tags.get(id),
  });
  const [editing, setEditing] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [mergeOpen, setMergeOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<TabKey>("scenes");
  const [sceneFilter, setSceneFilter] = useState<FindFilter>({ page: 1, perPage: 24, direction: "desc" });
  const [performerFilter, setPerformerFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "asc" });
  const [imageFilter, setImageFilter] = useState<FindFilter>({ page: 1, perPage: 30, direction: "desc" });
  const [galleryFilter, setGalleryFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "desc" });
  const [studioFilter, setStudioFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "asc" });
  const [groupFilter, setGroupFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "asc" });
  const queryClient = useQueryClient();

  useEffect(() => {
    if (tag) document.title = `${tag.name} | Stash`;
    return () => { document.title = "Stash"; };
  }, [tag]);

  // Keyboard shortcuts
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const el = (e.target as HTMLElement).tagName;
      if (el === "INPUT" || el === "TEXTAREA" || el === "SELECT") return;
      switch (e.key) {
        case "e": setEditing((v) => !v); break;
        case "f": if (tag) updateMut.mutate({ favorite: !tag.favorite }); break;
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [tag]);

  const deleteMut = useMutation({
    mutationFn: () => tags.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tags"] });
      onNavigate({ page: "tags" });
    },
  });

  const updateMut = useMutation({
    mutationFn: (data: { favorite?: boolean }) => tags.update(id, data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["tag", id] }),
  });

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-plex-accent" />
      </div>
    );
  }

  if (!tag) {
    return <div className="py-16 text-center text-plex-text-secondary">Tag not found</div>;
  }

  const tagImageUrl = tag.imagePath || entityImages.tagImageUrl(tag.id);

  return (
    <div className="min-h-screen">
      <div className="relative overflow-hidden border-b border-plex-border bg-[radial-gradient(circle_at_top_left,_rgba(204,123,25,0.2),_transparent_35%),linear-gradient(180deg,_rgba(50,54,57,0.95),_rgba(31,35,38,1))]">
        <div className="mx-auto max-w-7xl px-4 py-8">
          <div className="mb-5 flex items-center justify-between gap-4">
            <button
              onClick={() => onNavigate({ page: "tags" })}
              className="flex items-center gap-1 text-sm text-plex-text-secondary hover:text-plex-text"
            >
              <ArrowLeft className="h-4 w-4" /> Back to tags
            </button>
            <div className="flex items-center gap-2">
              <ExtensionSlot slot="tag-detail-actions" context={{ tag, onNavigate }} />
              <button
                onClick={() => setEditing(true)}
                className="flex items-center gap-1.5 rounded bg-plex-accent px-3 py-1.5 text-sm text-white hover:bg-plex-accent-hover"
              >
                <Pencil className="h-3.5 w-3.5" /> Edit
              </button>
              <button
                onClick={() => metadata.autoTag({ tags: [tag.name] })}
                className="flex items-center gap-1.5 rounded border border-plex-border bg-plex-card px-3 py-1.5 text-sm text-plex-text-secondary hover:text-plex-text"
                disabled={tag.ignoreAutoTag}
              >
                <Wand2 className="h-3.5 w-3.5" /> Auto Tag
              </button>
              <button
                onClick={() => setMergeOpen(true)}
                className="flex items-center gap-1.5 rounded border border-plex-border bg-plex-card px-3 py-1.5 text-sm text-plex-text-secondary hover:text-plex-text"
              >
                <GitMerge className="h-3.5 w-3.5" /> Merge...
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
            <div className="flex h-32 w-32 flex-shrink-0 items-center justify-center overflow-hidden rounded-2xl border border-plex-border bg-plex-card shadow-xl shadow-black/35 md:h-36 md:w-36">
              <img
                src={tagImageUrl}
                alt={tag.name}
                className="h-full w-full object-contain p-3"
                onError={(e) => {
                  (e.target as HTMLImageElement).style.display = "none";
                  const fallback = (e.target as HTMLImageElement).nextElementSibling as HTMLElement | null;
                  if (fallback) fallback.style.display = "flex";
                }}
              />
              <div className="hidden h-full w-full items-center justify-center bg-plex-card">
                <TagIcon className="h-14 w-14 text-plex-accent" />
              </div>
            </div>
            <div className="min-w-0 flex-1">
              <div className="mb-2 flex items-start gap-4">
                <div className="min-w-0 flex-1">
                  <h1 className="truncate text-4xl font-bold text-plex-text">{tag.name}</h1>
                  {tag.sortName && tag.sortName !== tag.name && (
                    <p className="mt-1 text-sm text-plex-text-muted">Sort name: {tag.sortName}</p>
                  )}
                  {tag.aliases.length > 0 && (
                    <p className="mt-1 text-sm text-plex-text-secondary">Also known as: {tag.aliases.join(", ")}</p>
                  )}
                </div>
                <button
                  onClick={() => updateMut.mutate({ favorite: !tag.favorite })}
                  className={`rounded-full p-2 transition-colors ${
                    tag.favorite
                      ? "bg-red-500/15 text-red-500"
                      : "bg-plex-card text-plex-text-muted hover:text-red-400"
                  }`}
                  title={tag.favorite ? "Remove from favorites" : "Add to favorites"}
                >
                  <Heart className={`h-6 w-6 ${tag.favorite ? "fill-current" : ""}`} />
                </button>
              </div>

              {tag.description && (
                <p className="max-w-4xl whitespace-pre-wrap text-sm leading-6 text-plex-text-secondary">{tag.description}</p>
              )}

              <div className="mt-4 flex flex-wrap gap-3">
                <CountCard label="Scenes" value={tag.sceneCount} icon={<Film className="h-4 w-4" />} />
                <CountCard label="Performers" value={tag.performerCount} icon={<UserRound className="h-4 w-4" />} />
                <CountCard label="Images" value={tag.imageCount} icon={<ImageIcon className="h-4 w-4" />} />
                <CountCard label="Galleries" value={tag.galleryCount} icon={<FolderOpen className="h-4 w-4" />} />
                <CountCard label="Markers" value={tag.markerCount} icon={<Bookmark className="h-4 w-4" />} />
                <CountCard label="Studios" value={tag.studioCount} icon={<Building2 className="h-4 w-4" />} />
                <CountCard label="Groups" value={tag.groupCount} icon={<Layers className="h-4 w-4" />} />
              </div>
            </div>
          </div>
        </div>
      </div>

      <TagEditModal tag={tag} open={editing} onClose={() => setEditing(false)} />
      <ConfirmDialog
        open={confirmDelete}
        title="Delete Tag"
        message={`Delete "${tag.name}"? This cannot be undone.`}
        onConfirm={() => deleteMut.mutate()}
        onCancel={() => setConfirmDelete(false)}
      />
      <MergeDialog
        open={mergeOpen}
        onClose={() => setMergeOpen(false)}
        entityType="tag"
        items={[{ id: tag.id, name: tag.name }]}
        onMerge={(targetId, sourceIds) => tags.merge(targetId, sourceIds)}
        queryKey={["tags"]}
      />

      <div className="mx-auto max-w-7xl px-4 py-6">
        <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
          <div className="min-w-0">
            {(tag.parents.length > 0 || tag.children.length > 0) && (
              <div className="mb-6 rounded-xl border border-plex-border bg-plex-card p-4">
                {tag.parents.length > 0 && (
                  <div className="mb-4">
                    <h2 className="mb-2 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">Parent Tags</h2>
                    <div className="flex flex-wrap gap-1.5">
                      {tag.parents.map((parent) => (
                        <span key={parent.id} className="inline-flex items-center gap-1 rounded bg-plex-surface px-1.5 py-1">
                          <span className="text-xs text-plex-text-muted">↖</span>
                          <TagBadge name={parent.name} onClick={() => onNavigate({ page: "tag", id: parent.id })} />
                        </span>
                      ))}
                    </div>
                  </div>
                )}
                {tag.children.length > 0 && (
                  <div>
                    <h2 className="mb-2 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">Sub Tags</h2>
                    <div className="flex flex-wrap gap-1.5">
                      {tag.children.map((child) => (
                        <span key={child.id} className="inline-flex items-center gap-1 rounded bg-plex-surface px-1.5 py-1">
                          <span className="text-xs text-plex-text-muted">↳</span>
                          <TagBadge name={child.name} onClick={() => onNavigate({ page: "tag", id: child.id })} />
                        </span>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}

            <div className="border-b border-plex-border">
              <div className="flex gap-1 overflow-x-auto">
                {[
                  { key: "scenes", label: "Scenes", count: tag.sceneCount },
                  { key: "performers", label: "Performers", count: tag.performerCount },
                  { key: "images", label: "Images", count: tag.imageCount },
                  { key: "galleries", label: "Galleries", count: tag.galleryCount },
                  { key: "markers", label: "Markers", count: tag.markerCount },
                  { key: "studios", label: "Studios", count: tag.studioCount },
                  { key: "groups", label: "Groups", count: tag.groupCount },
                ].map((tab) => (
                  <button
                    key={tab.key}
                    onClick={() => setActiveTab(tab.key as TabKey)}
                    className={`border-b-2 px-4 py-3 text-sm font-medium transition-colors ${
                      activeTab === tab.key
                        ? "border-plex-accent text-plex-text"
                        : "border-transparent text-plex-text-secondary hover:border-plex-text-muted hover:text-plex-text"
                    }`}
                  >
                    {tab.label}
                    <span className="ml-2 rounded-full bg-plex-card px-2 py-0.5 text-xs text-plex-text-muted">{tab.count}</span>
                  </button>
                ))}
              </div>
            </div>

            <div className="py-6">
              {activeTab === "scenes" && (
                <TagScenesPanel tagId={id} filter={sceneFilter} setFilter={setSceneFilter} onNavigate={onNavigate} />
              )}
              {activeTab === "performers" && (
                <TagPerformersPanel tagId={id} filter={performerFilter} setFilter={setPerformerFilter} onNavigate={onNavigate} />
              )}
              {activeTab === "images" && (
                <TagImagesPanel tagId={id} filter={imageFilter} setFilter={setImageFilter} onNavigate={onNavigate} />
              )}
              {activeTab === "galleries" && (
                <TagGalleriesPanel tagId={id} filter={galleryFilter} setFilter={setGalleryFilter} onNavigate={onNavigate} />
              )}
              {activeTab === "markers" && (
                <TagMarkersPanel tagId={id} onNavigate={onNavigate} />
              )}
              {activeTab === "studios" && (
                <TagStudiosPanel tagId={id} filter={studioFilter} setFilter={setStudioFilter} onNavigate={onNavigate} />
              )}
              {activeTab === "groups" && (
                <TagGroupsPanel tagId={id} filter={groupFilter} setFilter={setGroupFilter} onNavigate={onNavigate} />
              )}
            </div>
          </div>

          <aside className="space-y-4">
            <div className="rounded-xl border border-plex-border bg-plex-card p-4">
              <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">Metadata</h2>
              <dl className="space-y-2 text-sm">
                <div>
                  <dt className="text-plex-text-muted">Ignore Auto-Tag</dt>
                  <dd className="text-plex-text">{tag.ignoreAutoTag ? "Yes" : "No"}</dd>
                </div>
                <div>
                  <dt className="text-plex-text-muted">Created</dt>
                  <dd className="text-plex-text">{formatDate(tag.createdAt)}</dd>
                </div>
                <div>
                  <dt className="text-plex-text-muted">Updated</dt>
                  <dd className="text-plex-text">{formatDate(tag.updatedAt)}</dd>
                </div>
              </dl>
            </div>
            <ExtensionSlot slot="tag-detail-sidebar-bottom" context={{ tag, onNavigate }} />
          </aside>
        </div>

        <ExtensionSlot slot="tag-detail-bottom" context={{ tag, onNavigate }} />
      </div>
    </div>
  );
}

function CountCard({ label, value, icon }: { label: string; value: number; icon: React.ReactNode }) {
  return (
    <div className="flex items-center gap-2 rounded-lg border border-plex-border bg-plex-card px-3 py-2">
      <span className="text-plex-accent">{icon}</span>
      <div>
        <div className="text-lg font-semibold text-plex-text">{value}</div>
        <div className="text-xs text-plex-text-muted">{label}</div>
      </div>
    </div>
  );
}

function TagScenesPanel({ tagId, filter, setFilter, onNavigate }: {
  tagId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["tag-scenes", tagId, filter],
    queryFn: () => scenes.find(filter, { tagIds: String(tagId) }),
  });

  if (isLoading) return <LoadingPanel icon={<Film className="h-10 w-10" />} message="Loading scenes..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<Film className="h-12 w-12" />} message="No scenes with this tag" />;

  return (
    <>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
        {data.items.map((scene) => (
          <SceneTile key={scene.id} scene={scene} onClick={() => onNavigate({ page: "scene", id: scene.id })} />
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={data.totalCount} />
    </>
  );
}

function TagPerformersPanel({ tagId, filter, setFilter, onNavigate }: {
  tagId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["tag-performers", tagId, filter],
    queryFn: () => performers.find(filter, { tagIds: String(tagId) }),
  });

  if (isLoading) return <LoadingPanel icon={<UserRound className="h-10 w-10" />} message="Loading performers..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<UserRound className="h-12 w-12" />} message="No performers with this tag" />;

  return (
    <>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
        {data.items.map((performer) => (
          <PerformerTile key={performer.id} performer={performer} onClick={() => onNavigate({ page: "performer", id: performer.id })} />
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={data.totalCount} />
    </>
  );
}

function TagImagesPanel({ tagId, filter, setFilter, onNavigate }: {
  tagId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["tag-images", tagId, filter],
    queryFn: () => images.find(filter, { tagIds: String(tagId) }),
  });

  if (isLoading) return <LoadingPanel icon={<ImageIcon className="h-10 w-10" />} message="Loading images..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<ImageIcon className="h-12 w-12" />} message="No images with this tag" />;

  return (
    <>
      <div className="grid grid-cols-3 gap-3 sm:grid-cols-4 lg:grid-cols-6">
        {data.items.map((image) => (
          <ImageTile key={image.id} image={image} onClick={() => onNavigate({ page: "image", id: image.id })} />
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={data.totalCount} />
    </>
  );
}

function TagGalleriesPanel({ tagId, filter, setFilter, onNavigate }: {
  tagId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["tag-galleries", tagId, filter],
    queryFn: () => galleries.find(filter, { tagIds: String(tagId) }),
  });

  if (isLoading) return <LoadingPanel icon={<FolderOpen className="h-10 w-10" />} message="Loading galleries..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<FolderOpen className="h-12 w-12" />} message="No galleries with this tag" />;

  return (
    <>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
        {data.items.map((gallery) => (
          <GalleryTile key={gallery.id} gallery={gallery} onClick={() => onNavigate({ page: "gallery", id: gallery.id })} />
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={data.totalCount} />
    </>
  );
}

function TagMarkersPanel({ tagId, onNavigate }: { tagId: number; onNavigate: (r: any) => void }) {
  const { data, isLoading } = useQuery({
    queryKey: ["tag-markers", tagId],
    queryFn: () => markers.wall({ tagId, count: 100 }),
  });

  if (isLoading) return <LoadingPanel icon={<Bookmark className="h-10 w-10" />} message="Loading markers..." />;
  if (!data || data.length === 0) return <EmptyPanel icon={<Bookmark className="h-12 w-12" />} message="No markers with this tag" />;

  return (
    <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
      {data.map((marker) => (
        <button
          key={marker.id}
          onClick={() => onNavigate({ page: "scene", id: marker.sceneId })}
          className="group text-left"
        >
          <div className="relative aspect-video overflow-hidden rounded-lg border border-plex-border bg-plex-card shadow-md shadow-black/30">
            <img
              src={scenes.screenshotUrl(marker.sceneId)}
              alt={marker.title}
              className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
              loading="lazy"
            />
            <span className="absolute bottom-1.5 right-1.5 rounded bg-black/75 px-1.5 py-0.5 text-[11px] text-white">
              {formatDuration(marker.seconds)}
            </span>
          </div>
          <div className="pt-2">
            <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">{marker.title}</p>
            <p className="mt-0.5 truncate text-xs text-plex-text-secondary">{marker.sceneTitle || "Untitled Scene"}</p>
          </div>
        </button>
      ))}
    </div>
  );
}

function TagStudiosPanel({ tagId, filter, setFilter, onNavigate }: {
  tagId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["tag-studios", tagId, filter],
    queryFn: () => studios.find(filter, { tagIds: String(tagId) }),
  });

  if (isLoading) return <LoadingPanel icon={<Building2 className="h-10 w-10" />} message="Loading studios..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<Building2 className="h-12 w-12" />} message="No studios with this tag" />;

  return (
    <>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
        {data.items.map((studio) => (
          <StudioTile key={studio.id} studio={studio} onClick={() => onNavigate({ page: "studio", id: studio.id })} />
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={data.totalCount} />
    </>
  );
}

function TagGroupsPanel({ tagId, filter, setFilter, onNavigate }: {
  tagId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["tag-groups", tagId, filter],
    queryFn: () => groups.find(filter, { tagIds: String(tagId) }),
  });

  if (isLoading) return <LoadingPanel icon={<Layers className="h-10 w-10" />} message="Loading groups..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<Layers className="h-12 w-12" />} message="No groups with this tag" />;

  return (
    <>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
        {data.items.map((group) => (
          <GroupTile key={group.id} group={group} onClick={() => onNavigate({ page: "group", id: group.id })} />
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={data.totalCount} />
    </>
  );
}

function SceneTile({ scene, onClick }: { scene: Scene; onClick: () => void }) {
  const file = scene.files[0];
  const duration = file?.duration ?? 0;
  const resLabel = file ? getResolutionLabel(file.width, file.height) : null;

  return (
    <button onClick={onClick} className="group text-left">
      <div className="relative aspect-video overflow-hidden rounded-lg border border-plex-border bg-plex-card shadow-md shadow-black/30">
        <img src={scenes.screenshotUrl(scene.id)} alt={scene.title || ""} className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105" loading="lazy" />
        {duration > 0 && <span className="absolute bottom-1.5 right-1.5 rounded bg-black/75 px-1.5 py-0.5 text-[11px] text-white">{formatDuration(duration)}</span>}
        {resLabel && <span className="absolute top-1.5 right-1.5 rounded bg-black/75 px-1.5 py-0.5 text-[10px] font-bold uppercase text-plex-accent">{resLabel}</span>}
      </div>
      <div className="pt-2">
        <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">{scene.title || "Untitled"}</p>
        <p className="mt-0.5 truncate text-xs text-plex-text-secondary">{scene.date || scene.studioName || ""}</p>
      </div>
    </button>
  );
}

function PerformerTile({ performer, onClick }: { performer: Performer; onClick: () => void }) {
  return (
    <button onClick={onClick} className="group overflow-hidden rounded-xl border border-plex-border bg-plex-card text-left transition-colors hover:border-plex-accent/60">
      <div className="aspect-[2/3] bg-gradient-to-b from-plex-card to-plex-surface" />
      <div className="p-3">
        <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">{performer.name}</p>
        <p className="mt-1 text-xs text-plex-text-secondary">{performer.sceneCount} scenes</p>
      </div>
    </button>
  );
}

function ImageTile({ image, onClick }: { image: Image; onClick: () => void }) {
  return (
    <button onClick={onClick} className="group overflow-hidden rounded-lg border border-plex-border bg-plex-card text-left shadow-md shadow-black/20">
      <div className="aspect-square overflow-hidden bg-plex-surface">
        <img src={images.thumbnailUrl(image.id)} alt={image.title || ""} className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105" loading="lazy" />
      </div>
      <div className="p-2">
        <p className="truncate text-xs text-plex-text group-hover:text-plex-accent">{image.title || "Untitled"}</p>
      </div>
    </button>
  );
}

function GalleryTile({ gallery, onClick }: { gallery: Gallery; onClick: () => void }) {
  return (
    <button onClick={onClick} className="group overflow-hidden rounded-xl border border-plex-border bg-plex-card text-left transition-colors hover:border-plex-accent/60">
      <div className="flex aspect-video items-center justify-center bg-gradient-to-br from-plex-surface to-plex-card">
        <FolderOpen className="h-10 w-10 text-plex-text-muted" />
      </div>
      <div className="p-3">
        <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">{gallery.title || "Untitled"}</p>
        <p className="mt-1 text-xs text-plex-text-secondary">{gallery.imageCount} images</p>
      </div>
    </button>
  );
}

function StudioTile({ studio, onClick }: { studio: Studio; onClick: () => void }) {
  return (
    <button onClick={onClick} className="group overflow-hidden rounded-xl border border-plex-border bg-plex-card text-left transition-colors hover:border-plex-accent/60">
      <div className="flex aspect-video items-center justify-center bg-gradient-to-br from-plex-surface to-plex-card">
        <Building2 className="h-10 w-10 text-plex-text-muted" />
      </div>
      <div className="p-3">
        <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">{studio.name}</p>
        <p className="mt-1 text-xs text-plex-text-secondary">{studio.sceneCount} scenes</p>
      </div>
    </button>
  );
}

function GroupTile({ group, onClick }: { group: Group; onClick: () => void }) {
  return (
    <button onClick={onClick} className="group overflow-hidden rounded-xl border border-plex-border bg-plex-card text-left transition-colors hover:border-plex-accent/60">
      <div className="flex aspect-video items-center justify-center bg-gradient-to-br from-plex-surface to-plex-card">
        <Layers className="h-10 w-10 text-plex-text-muted" />
      </div>
      <div className="p-3">
        <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">{group.name}</p>
      </div>
    </button>
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
      <button
        disabled={page <= 1}
        onClick={() => setFilter({ ...filter, page: page - 1 })}
        className="rounded border border-plex-border bg-plex-card px-4 py-2 text-sm text-plex-text-secondary hover:bg-plex-card-hover disabled:cursor-not-allowed disabled:opacity-50"
      >
        Previous
      </button>
      <span className="text-sm text-plex-text-secondary">Page {page} of {totalPages}</span>
      <button
        disabled={page >= totalPages}
        onClick={() => setFilter({ ...filter, page: page + 1 })}
        className="rounded border border-plex-border bg-plex-card px-4 py-2 text-sm text-plex-text-secondary hover:bg-plex-card-hover disabled:cursor-not-allowed disabled:opacity-50"
      >
        Next
      </button>
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
    <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-plex-border bg-plex-card/40 py-12 text-plex-text-muted">
      <div className="mb-3 opacity-60">{icon}</div>
      <p>{message}</p>
    </div>
  );
}