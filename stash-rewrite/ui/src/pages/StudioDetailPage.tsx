import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { galleries, groups, images, metadata, performers, scenes, studios, entityImages } from "../api/client";
import type { FindFilter, Gallery, Group, Image, Performer, Scene, StashBox, StashBoxStudioMatch, Studio } from "../api/types";
import { formatDate, formatDuration, getResolutionLabel, TagBadge, CustomFieldsDisplay } from "../components/shared";
import { ArrowLeft, Building2, CloudDownload, Film, FolderOpen, GitMerge, Heart, ImageIcon, Layers, Link as LinkIcon, Link2, Loader2, Pencil, Search, Trash2, UserRound, Wand2 } from "lucide-react";
import { useEffect, useState } from "react";
import { StudioEditModal } from "./StudioEditModal";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { DetailMergeDialog } from "../components/DetailMergeDialog";
import { ExtensionSlot } from "../router/RouteRegistry";
import { InteractiveRating } from "../components/Rating";
import { useAppConfig } from "../state/AppConfigContext";

interface Props {
  id: number;
  onNavigate: (r: any) => void;
}

type TabKey = "scenes" | "performers" | "galleries" | "images" | "studios" | "groups";

export function StudioDetailPage({ id, onNavigate }: Props) {
  const { config } = useAppConfig();
  const { data: studio, isLoading } = useQuery({
    queryKey: ["studio", id],
    queryFn: () => studios.get(id),
  });
  const [editing, setEditing] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [mergeOpen, setMergeOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<TabKey>("scenes");
  const [sceneFilter, setSceneFilter] = useState<FindFilter>({ page: 1, perPage: 24, direction: "desc" });
  const [galleryFilter, setGalleryFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "desc" });
  const [imageFilter, setImageFilter] = useState<FindFilter>({ page: 1, perPage: 30, direction: "desc" });
  const [performerFilter, setPerformerFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "asc" });
  const [childFilter, setChildFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "asc" });
  const [groupFilter, setGroupFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "asc" });
  const queryClient = useQueryClient();

  useEffect(() => {
    if (studio) document.title = `${studio.name} | Stash`;
    return () => { document.title = "Stash"; };
  }, [studio]);

  // Keyboard shortcuts
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const tag = (e.target as HTMLElement).tagName;
      if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return;
      switch (e.key) {
        case "e": setEditing((v) => !v); break;
        case "f": if (studio) updateMut.mutate({ favorite: !studio.favorite }); break;
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [studio]);

  const deleteMut = useMutation({
    mutationFn: () => studios.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["studios"] });
      onNavigate({ page: "studios" });
    },
  });

  const updateMut = useMutation({
    mutationFn: (data: { favorite?: boolean; rating?: number }) => studios.update(id, data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["studio", id] }),
  });

  const autoTagMut = useMutation({
    mutationFn: () => {
      if (!studio) throw new Error("Studio not loaded");
      return metadata.autoTag({ studios: [studio.name] });
    },
  });

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-plex-accent" />
      </div>
    );
  }

  if (!studio) {
    return <div className="py-16 text-center text-plex-text-secondary">Studio not found</div>;
  }

  const studioImageUrl = studio.imagePath || entityImages.studioImageUrl(studio.id);

  return (
    <div className="min-h-screen">
      <div className="relative overflow-hidden border-b border-plex-border bg-[radial-gradient(circle_at_top_right,_rgba(204,123,25,0.16),_transparent_30%),linear-gradient(180deg,_rgba(50,54,57,0.95),_rgba(31,35,38,1))]">
        {/* Background studio image */}
        <img
          src={entityImages.studioImageUrl(studio.id)}
          alt=""
          className="absolute inset-0 h-full w-full object-cover opacity-10 blur-md scale-110"
          onError={(e) => { (e.target as HTMLImageElement).style.display = "none"; }}
        />
        <div className="absolute inset-0 bg-gradient-to-t from-plex-bg via-plex-bg/70 to-transparent" />
        <div className="relative mx-auto max-w-7xl px-4 py-8">
          <div className="mb-5 flex items-center justify-between gap-4">
            <button
              onClick={() => onNavigate({ page: "studios" })}
              className="flex items-center gap-1 text-sm text-plex-text-secondary hover:text-plex-text"
            >
              <ArrowLeft className="h-4 w-4" /> Back to studios
            </button>
            <div className="flex items-center gap-2">
              <ExtensionSlot slot="studio-detail-actions" context={{ studio, onNavigate }} />
              <button
                onClick={() => setEditing(true)}
                className="flex items-center gap-1.5 rounded bg-plex-accent px-3 py-1.5 text-sm text-white hover:bg-plex-accent-hover"
              >
                <Pencil className="h-3.5 w-3.5" /> Edit
              </button>
              <button
                onClick={() => autoTagMut.mutate()}
                disabled={autoTagMut.isPending}
                className="flex items-center gap-1.5 rounded border border-plex-border bg-plex-card px-3 py-1.5 text-sm text-plex-text-secondary hover:text-plex-text"
              >
                {autoTagMut.isPending ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <Wand2 className="h-3.5 w-3.5" />} Auto Tag
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
                src={studioImageUrl}
                alt={studio.name}
                className="h-full w-full object-contain p-3"
                onError={(e) => {
                  (e.target as HTMLImageElement).style.display = "none";
                  const fallback = (e.target as HTMLImageElement).nextElementSibling as HTMLElement | null;
                  if (fallback) fallback.style.display = "flex";
                }}
              />
              <div className="hidden h-full w-full items-center justify-center bg-plex-card">
                <Building2 className="h-14 w-14 text-plex-accent" />
              </div>
            </div>
            <div className="min-w-0 flex-1">
              <div className="mb-2 flex items-start gap-4">
                <div className="min-w-0 flex-1">
                  <h1 className="truncate text-4xl font-bold text-plex-text">{studio.name}</h1>
                  {studio.parentName && studio.parentId && (
                    <button
                      onClick={() => onNavigate({ page: "studio", id: studio.parentId })}
                      className="mt-1 text-sm text-plex-accent hover:underline"
                    >
                      Part of {studio.parentName}
                    </button>
                  )}
                  {studio.aliases.length > 0 && (
                    <p className="mt-1 text-sm text-plex-text-secondary">Also known as: {studio.aliases.join(", ")}</p>
                  )}
                </div>
                <button
                  onClick={() => updateMut.mutate({ favorite: !studio.favorite })}
                  className={`rounded-full p-2 transition-colors ${
                    studio.favorite
                      ? "bg-red-500/15 text-red-500"
                      : "bg-plex-card text-plex-text-muted hover:text-red-400"
                  }`}
                  title={studio.favorite ? "Remove from favorites" : "Add to favorites"}
                >
                  <Heart className={`h-6 w-6 ${studio.favorite ? "fill-current" : ""}`} />
                </button>
              </div>

              <div className="mb-4 flex flex-wrap items-center gap-3">
                <InteractiveRating value={studio.rating} onChange={(value) => updateMut.mutate({ rating: value })} />
                <span className={`rounded px-2 py-1 text-xs font-medium ${studio.organized ? "bg-green-600 text-white" : "bg-plex-card text-plex-text-muted border border-plex-border"}`}>
                  {studio.organized ? "Organized" : "Unorganized"}
                </span>
                <span className="rounded border border-plex-border bg-plex-card px-2 py-1 text-xs text-plex-text-secondary">
                  {studio.sceneCount} scenes
                </span>
              </div>

              {studio.details && (
                <p className="max-w-4xl whitespace-pre-wrap text-sm leading-6 text-plex-text-secondary">{studio.details}</p>
              )}
              {autoTagMut.isSuccess && (
                <p className="mt-3 text-sm text-emerald-300">Auto-tag job queued.</p>
              )}
            </div>
          </div>
        </div>
      </div>

      <StudioEditModal studio={studio} open={editing} onClose={() => setEditing(false)} />
      <ConfirmDialog
        open={confirmDelete}
        title="Delete Studio"
        message={`Delete "${studio.name}"? This cannot be undone.`}
        onConfirm={() => deleteMut.mutate()}
        onCancel={() => setConfirmDelete(false)}
      />
      <DetailMergeDialog
        open={mergeOpen}
        onClose={() => setMergeOpen(false)}
        entityType="studio"
        targetItem={{ id: studio.id, name: studio.name, imagePath: studioImageUrl, subtitle: studio.parentName }}
        searchItems={async (term) => {
          const response = await studios.find({ page: 1, perPage: 20, sort: "name", direction: "asc", q: term || undefined });
          return response.items.map((item) => ({
            id: item.id,
            name: item.name,
            imagePath: item.imagePath,
            subtitle: item.parentName,
          }));
        }}
        onMerge={(targetId, sourceIds) => studios.merge(targetId, sourceIds)}
        invalidateQueryKeys={[["studio", id], ["studios"]]}
      />

      <div className="mx-auto max-w-7xl px-4 py-6">
        <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
          <div className="min-w-0">
            {studio.tags.length > 0 && (
              <div className="mb-6 rounded-xl border border-plex-border bg-plex-card p-4">
                <h2 className="mb-2 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">Tags</h2>
                <div className="flex flex-wrap gap-1.5">
                  {studio.tags.map((tag) => (
                    <TagBadge key={tag.id} name={tag.name} onClick={() => onNavigate({ page: "tag", id: tag.id })} />
                  ))}
                </div>
              </div>
            )}

            <div className="border-b border-plex-border">
              <div className="flex gap-1 overflow-x-auto">
                {[
                  { key: "scenes", label: "Scenes" },
                  { key: "performers", label: "Performers" },
                  { key: "galleries", label: "Galleries" },
                  { key: "images", label: "Images" },
                  { key: "studios", label: "Sub-studios" },
                  { key: "groups", label: "Groups" },
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
                  </button>
                ))}
              </div>
            </div>

            <div className="py-6">
              {activeTab === "scenes" && (
                <StudioScenesPanel studioId={id} filter={sceneFilter} setFilter={setSceneFilter} onNavigate={onNavigate} />
              )}
              {activeTab === "performers" && (
                <StudioPerformersPanel studioId={id} filter={performerFilter} setFilter={setPerformerFilter} onNavigate={onNavigate} />
              )}
              {activeTab === "galleries" && (
                <StudioGalleriesPanel studioId={id} filter={galleryFilter} setFilter={setGalleryFilter} onNavigate={onNavigate} />
              )}
              {activeTab === "images" && (
                <StudioImagesPanel studioId={id} filter={imageFilter} setFilter={setImageFilter} onNavigate={onNavigate} />
              )}
              {activeTab === "studios" && (
                <ChildStudiosPanel studioId={id} filter={childFilter} setFilter={setChildFilter} onNavigate={onNavigate} />
              )}
              {activeTab === "groups" && (
                <StudioGroupsPanel studioId={id} filter={groupFilter} setFilter={setGroupFilter} onNavigate={onNavigate} />
              )}
            </div>
          </div>

          <aside className="space-y-4">
            <div className="rounded-xl border border-plex-border bg-plex-card p-4">
              <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">Metadata</h2>
              <dl className="space-y-2 text-sm">
                <div>
                  <dt className="text-plex-text-muted">Organized</dt>
                  <dd className="text-plex-text">{studio.organized ? "Yes" : "No"}</dd>
                </div>
                <div>
                  <dt className="text-plex-text-muted">Ignore Auto-Tag</dt>
                  <dd className="text-plex-text">{studio.ignoreAutoTag ? "Yes" : "No"}</dd>
                </div>
                <div>
                  <dt className="text-plex-text-muted">Created</dt>
                  <dd className="text-plex-text">{formatDate(studio.createdAt)}</dd>
                </div>
                <div>
                  <dt className="text-plex-text-muted">Updated</dt>
                  <dd className="text-plex-text">{formatDate(studio.updatedAt)}</dd>
                </div>
              </dl>
            </div>

            {studio.urls.length > 0 && (
              <div className="rounded-xl border border-plex-border bg-plex-card p-4">
                <h2 className="mb-3 flex items-center gap-1.5 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">
                  <LinkIcon className="h-4 w-4" /> URLs
                </h2>
                <div className="space-y-1 text-sm">
                  {studio.urls.map((url, index) => (
                    <a key={index} href={url} target="_blank" rel="noopener noreferrer" className="block truncate text-plex-accent hover:underline">
                      {url}
                    </a>
                  ))}
                </div>
              </div>
            )}

            <CustomFieldsDisplay customFields={studio.customFields} />
            <ExtensionSlot slot="studio-detail-sidebar-bottom" context={{ studio, onNavigate }} />
          </aside>
        </div>

        {/* StashBox Integration */}
        <StudioStashBoxPanel studio={studio} stashBoxes={config?.scraping.stashBoxes ?? []} onNavigate={onNavigate} />

        {/* StashIDs */}
        {studio.stashIds && studio.stashIds.length > 0 && (
          <div className="mt-4 rounded-xl border border-plex-border bg-plex-card p-4">
            <h2 className="mb-3 flex items-center gap-1.5 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">
              <Link2 className="h-4 w-4" /> Stash IDs
            </h2>
            <div className="flex flex-wrap gap-2">
              {studio.stashIds.map((sid) => (
                <span key={`${sid.endpoint}:${sid.stashId}`} className="inline-flex items-center gap-1.5 rounded-full border border-plex-border px-3 py-1 text-xs text-plex-text-secondary">
                  <Link2 className="h-3 w-3 text-plex-accent" />
                  {sid.stashId.slice(0, 12)}…
                </span>
              ))}
            </div>
          </div>
        )}

        <ExtensionSlot slot="studio-detail-bottom" context={{ studio, onNavigate }} />
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

function StudioScenesPanel({ studioId, filter, setFilter, onNavigate }: {
  studioId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["studio-scenes", studioId, filter],
    queryFn: () => scenes.find(filter, { studioId: String(studioId) }),
  });

  if (isLoading) return <LoadingPanel icon={<Film className="h-10 w-10" />} message="Loading scenes..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<Film className="h-12 w-12" />} message="No scenes from this studio" />;

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

function StudioGalleriesPanel({ studioId, filter, setFilter, onNavigate }: {
  studioId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["studio-galleries", studioId, filter],
    queryFn: () => galleries.find(filter, { studioId: String(studioId) }),
  });

  if (isLoading) return <LoadingPanel icon={<FolderOpen className="h-10 w-10" />} message="Loading galleries..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<FolderOpen className="h-12 w-12" />} message="No galleries from this studio" />;

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

function StudioImagesPanel({ studioId, filter, setFilter, onNavigate }: {
  studioId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["studio-images", studioId, filter],
    queryFn: () => images.find(filter, { studioId: String(studioId) }),
  });

  if (isLoading) return <LoadingPanel icon={<ImageIcon className="h-10 w-10" />} message="Loading images..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<ImageIcon className="h-12 w-12" />} message="No images from this studio" />;

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

function ChildStudiosPanel({ studioId, filter, setFilter, onNavigate }: {
  studioId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["child-studios", studioId, filter],
    queryFn: () => studios.find(filter, { parentId: String(studioId) }),
  });

  if (isLoading) return <LoadingPanel icon={<Building2 className="h-10 w-10" />} message="Loading sub-studios..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<Building2 className="h-12 w-12" />} message="No sub-studios" />;

  return (
    <>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
        {data.items.map((childStudio) => (
          <StudioTile key={childStudio.id} studio={childStudio} onClick={() => onNavigate({ page: "studio", id: childStudio.id })} />
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={data.totalCount} />
    </>
  );
}

function StudioPerformersPanel({ studioId, filter, setFilter, onNavigate }: {
  studioId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["studio-performers", studioId, filter],
    queryFn: () => performers.find(filter, { studioId: String(studioId) }),
  });

  if (isLoading) return <LoadingPanel icon={<UserRound className="h-10 w-10" />} message="Loading performers..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<UserRound className="h-12 w-12" />} message="No performers from this studio" />;

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

function StudioGroupsPanel({ studioId, filter, setFilter, onNavigate }: {
  studioId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["studio-groups", studioId, filter],
    queryFn: () => groups.find(filter, { studioId: String(studioId) }),
  });

  if (isLoading) return <LoadingPanel icon={<Layers className="h-10 w-10" />} message="Loading groups..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<Layers className="h-12 w-12" />} message="No groups from this studio" />;

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
        <p className="mt-0.5 truncate text-xs text-plex-text-secondary">{scene.date || ""}</p>
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

function StudioTile({ studio, onClick }: { studio: Studio; onClick: () => void }) {
  return (
    <button onClick={onClick} className="group flex flex-col rounded-xl border border-plex-border bg-plex-card p-4 text-left transition-colors hover:border-plex-accent/60">
      <div className="mb-4 flex h-20 items-center justify-center rounded-lg bg-plex-surface">
        <Building2 className="h-10 w-10 text-plex-text-muted" />
      </div>
      <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">{studio.name}</p>
      <p className="mt-1 text-xs text-plex-text-secondary">{studio.sceneCount} scenes</p>
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

function StudioStashBoxPanel({ studio, stashBoxes, onNavigate }: { studio: Studio; stashBoxes: StashBox[]; onNavigate: (r: any) => void }) {
  const queryClient = useQueryClient();
  const [term, setTerm] = useState(studio.name);
  const [selectedEndpoint, setSelectedEndpoint] = useState("");

  useEffect(() => {
    setTerm(studio.name);
  }, [studio.id, studio.name]);

  useEffect(() => {
    if (selectedEndpoint && !stashBoxes.some((box) => box.endpoint === selectedEndpoint)) {
      setSelectedEndpoint("");
    }
  }, [selectedEndpoint, stashBoxes]);

  const searchMutation = useMutation({
    mutationFn: (variables: { term?: string; endpoint?: string }) => studios.searchStashBox(studio.id, variables.term, variables.endpoint),
  });

  const importMutation = useMutation({
    mutationFn: (match: StashBoxStudioMatch) =>
      studios.importFromStashBox(studio.id, { endpoint: match.endpoint, studioId: match.id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["studio", studio.id] });
      queryClient.invalidateQueries({ queryKey: ["studios"] });
    },
  });

  if (stashBoxes.length === 0) return null;

  return (
    <div className="mt-6 rounded-xl border border-plex-border bg-plex-card p-4">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div>
          <h2 className="text-base font-semibold text-plex-text">StashBox</h2>
          <p className="mt-1 max-w-2xl text-sm text-plex-text-secondary">
            Search configured StashBox instances and merge remote studio metadata into this local studio.
          </p>
        </div>
      </div>

      <div className="mt-4 grid gap-3 xl:grid-cols-[minmax(0,2fr)_minmax(0,1fr)_auto]">
        <label className="block text-sm">
          <span className="mb-1 block text-xs font-medium uppercase tracking-wide text-plex-text-muted">Search term</span>
          <input
            value={term}
            onChange={(event) => setTerm(event.target.value)}
            placeholder={studio.name}
            className="w-full rounded-xl border border-plex-border bg-plex-surface px-3 py-2 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
          />
        </label>
        <label className="block text-sm">
          <span className="mb-1 block text-xs font-medium uppercase tracking-wide text-plex-text-muted">Endpoint</span>
          <select
            value={selectedEndpoint}
            onChange={(event) => setSelectedEndpoint(event.target.value)}
            className="w-full rounded-xl border border-plex-border bg-plex-surface px-3 py-2 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
          >
            <option value="">All configured endpoints</option>
            {stashBoxes.map((box) => (
              <option key={box.endpoint} value={box.endpoint}>
                {box.name || box.endpoint}
              </option>
            ))}
          </select>
        </label>
        <div className="flex items-end">
          <button
            onClick={() => searchMutation.mutate({ term: term.trim() || undefined, endpoint: selectedEndpoint || undefined })}
            disabled={searchMutation.isPending}
            className="inline-flex items-center gap-2 rounded-xl border border-plex-border px-4 py-2 text-sm text-plex-text hover:border-plex-accent hover:text-plex-accent disabled:opacity-60"
          >
            {searchMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Search className="h-4 w-4" />}
            Search StashBox
          </button>
        </div>
      </div>

      {searchMutation.error && (
        <p className="mt-3 text-sm text-red-300">{(searchMutation.error as Error).message}</p>
      )}

      {importMutation.isSuccess && (
        <p className="mt-3 text-sm text-emerald-300">Studio metadata imported from StashBox.</p>
      )}

      {searchMutation.data && (
        <div className="mt-4 space-y-3">
          {searchMutation.data.length === 0 ? (
            <div className="rounded-xl border border-dashed border-plex-border p-4 text-sm text-plex-text-secondary">
              No StashBox studio matches were found.
            </div>
          ) : (
            searchMutation.data.map((match) => (
              <button
                key={`${match.endpoint}:${match.id}`}
                onClick={() => importMutation.mutate(match)}
                disabled={importMutation.isPending}
                className="flex w-full flex-col gap-4 rounded-xl border border-plex-border bg-plex-surface p-4 text-left transition-colors hover:border-plex-accent/60 disabled:opacity-60 md:flex-row"
              >
                <div className="h-20 w-20 flex-shrink-0 overflow-hidden rounded-lg border border-plex-border bg-black/20">
                  {match.imageUrl ? (
                    <img src={match.imageUrl} alt={match.name} className="h-full w-full object-contain p-2" />
                  ) : (
                    <div className="flex h-full w-full items-center justify-center bg-gradient-to-b from-plex-card to-plex-surface">
                      <Building2 className="h-10 w-10 text-plex-text-muted/50" />
                    </div>
                  )}
                </div>

                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="text-base font-semibold text-plex-text">{match.name}</span>
                    <span className="rounded-full border border-plex-border px-2 py-0.5 text-xs text-plex-text-secondary">
                      {match.stashBoxName}
                    </span>
                  </div>
                  {match.parentName && <p className="mt-1 text-sm text-plex-text-secondary">Parent: {match.parentName}</p>}
                  <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs text-plex-text-muted">
                    <span>ID {match.id}</span>
                  </div>
                  {match.aliases.length > 0 && <p className="mt-2 text-xs text-plex-text-secondary">Aliases: {match.aliases.join(", ")}</p>}
                  {match.urls.length > 0 && <p className="mt-1 truncate text-xs text-plex-text-muted">{match.urls[0]}</p>}
                </div>

                <div className="flex items-end">
                  <span className="inline-flex items-center gap-2 rounded-lg bg-plex-accent px-3 py-2 text-sm font-medium text-white">
                    {importMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <CloudDownload className="h-4 w-4" />}
                    Import
                  </span>
                </div>
              </button>
            ))
          )}
        </div>
      )}
    </div>
  );
}