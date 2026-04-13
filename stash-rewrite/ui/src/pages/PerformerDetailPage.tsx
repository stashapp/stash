import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { galleries, images, metadata, performers, scenes, entityImages } from "../api/client";
import type { FindFilter, Gallery, Image, Performer as PerformerModel, Scene, StashBox, StashBoxPerformerMatch } from "../api/types";
import { formatDate, formatDuration, getResolutionLabel, TagBadge } from "../components/shared";
import { ArrowLeft, Calendar, CloudDownload, ExternalLink, Film, FolderOpen, GitMerge, Heart, ImageIcon, Layers, Link2, Loader2, MapPin, MoreVertical, Pencil, Ruler, Scale, Search, Trash2, Users, UserRound, Wand2 } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { PerformerEditModal } from "./PerformerEditModal";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { MergeDialog } from "../components/MergeDialog";
import { ExtensionSlot } from "../router/RouteRegistry";
import { InteractiveRating } from "../components/Rating";
import { useAppConfig } from "../state/AppConfigContext";

interface Props {
  id: number;
  onNavigate: (r: any) => void;
}

type TabKey = "scenes" | "galleries" | "images" | "groups" | "appearsWith";

export function PerformerDetailPage({ id, onNavigate }: Props) {
  const { config } = useAppConfig();
  const { data: performer, isLoading } = useQuery({
    queryKey: ["performer", id],
    queryFn: () => performers.get(id),
  });
  const [editing, setEditing] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const [mergeOpen, setMergeOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<TabKey>("scenes");
  const [sceneFilter, setSceneFilter] = useState<FindFilter>({ page: 1, perPage: 24, direction: "desc" });
  const [galleryFilter, setGalleryFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "desc" });
  const [imageFilter, setImageFilter] = useState<FindFilter>({ page: 1, perPage: 30, direction: "desc" });
  const [groupFilter, setGroupFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "asc" });
  const [appearsWithFilter, setAppearsWithFilter] = useState<FindFilter>({ page: 1, perPage: 18, direction: "asc" });
  const [showOpsMenu, setShowOpsMenu] = useState(false);
  const opsMenuRef = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();

  const deleteMut = useMutation({
    mutationFn: () => performers.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["performers"] });
      onNavigate({ page: "performers" });
    },
  });

  const updateMut = useMutation({
    mutationFn: (data: { favorite?: boolean; rating?: number }) => performers.update(id, data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["performer", id] }),
  });

  useEffect(() => {
    if (performer) document.title = `${performer.name} | Stash`;
    return () => { document.title = "Stash"; };
  }, [performer]);

  // Keyboard shortcuts
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const tag = (e.target as HTMLElement).tagName;
      if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return;
      switch (e.key) {
        case "e": setEditing((v) => !v); break;
        case "f": if (performer) updateMut.mutate({ favorite: !performer.favorite }); break;
        case "c": setActiveTab("scenes"); break;
        case "g": setActiveTab("galleries"); break;
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [performer]);

  useEffect(() => {
    if (!showOpsMenu) return;
    const handler = (e: MouseEvent) => {
      if (opsMenuRef.current && !opsMenuRef.current.contains(e.target as Node)) setShowOpsMenu(false);
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, [showOpsMenu]);

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-plex-accent" />
      </div>
    );
  }

  if (!performer) {
    return <div className="py-16 text-center text-plex-text-secondary">Performer not found</div>;
  }

  const age = performer.birthdate
    ? Math.floor((Date.now() - new Date(performer.birthdate).getTime()) / 31557600000)
    : null;

  return (
    <div className="min-h-screen">
      <div className="relative overflow-hidden border-b border-plex-border bg-[radial-gradient(circle_at_top_right,_rgba(204,123,25,0.18),_transparent_30%),linear-gradient(180deg,_rgba(50,54,57,0.95),_rgba(31,35,38,1))]">
        <div className="mx-auto max-w-7xl px-4 py-8">
          <div className="mb-5 flex items-center justify-between gap-4">
            <button
              onClick={() => onNavigate({ page: "performers" })}
              className="flex items-center gap-1 text-sm text-plex-text-secondary hover:text-plex-text"
            >
              <ArrowLeft className="h-4 w-4" /> Back to performers
            </button>
            <div className="flex items-center gap-2">
              <ExtensionSlot slot="performer-detail-actions" context={{ performer, onNavigate }} />
              <div className="relative" ref={opsMenuRef}>
                <button
                  onClick={() => setShowOpsMenu(!showOpsMenu)}
                  className="rounded border border-plex-border bg-plex-card p-2 text-plex-text-secondary hover:text-plex-text"
                  title="Actions"
                >
                  <MoreVertical className="h-4 w-4" />
                </button>
                {showOpsMenu && (
                  <div className="absolute right-0 z-50 mt-1 min-w-[160px] rounded-lg border border-plex-border bg-plex-card py-1 shadow-xl">
                    <button onClick={() => { setEditing(true); setShowOpsMenu(false); }} className="flex w-full items-center gap-2 px-3 py-2 text-sm text-plex-text hover:bg-plex-surface">
                      <Pencil className="h-3.5 w-3.5" /> Edit
                    </button>
                    <button onClick={() => { metadata.autoTag({ performers: [performer.name] }); setShowOpsMenu(false); }} className="flex w-full items-center gap-2 px-3 py-2 text-sm text-plex-text hover:bg-plex-surface">
                      <Wand2 className="h-3.5 w-3.5" /> Auto Tag
                    </button>
                    <button onClick={() => { setMergeOpen(true); setShowOpsMenu(false); }} className="flex w-full items-center gap-2 px-3 py-2 text-sm text-plex-text hover:bg-plex-surface">
                      <GitMerge className="h-3.5 w-3.5" /> Merge...
                    </button>
                    <div className="my-1 border-t border-plex-border" />
                    <button onClick={() => { setConfirmDelete(true); setShowOpsMenu(false); }} className="flex w-full items-center gap-2 px-3 py-2 text-sm text-red-400 hover:bg-plex-surface">
                      <Trash2 className="h-3.5 w-3.5" /> Delete
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>

          <div className="flex flex-col gap-6 md:flex-row md:items-start">
            <div className="w-48 flex-shrink-0 md:w-64">
              <div className="aspect-[2/3] overflow-hidden rounded-2xl border border-plex-border bg-plex-card shadow-2xl shadow-black/40">
                <img
                  src={performer.imagePath || entityImages.performerImageUrl(performer.id)}
                  alt={performer.name}
                  className="h-full w-full object-cover"
                  onError={(e) => {
                    const el = e.target as HTMLImageElement;
                    el.style.display = "none";
                    if (el.nextElementSibling) (el.nextElementSibling as HTMLElement).style.display = "";
                  }}
                />
                <div className="flex h-full w-full items-center justify-center bg-gradient-to-b from-plex-card to-plex-surface" style={{ display: "none" }}>
                  <UserRound className="h-20 w-20 text-plex-text-muted/50" />
                </div>
              </div>
            </div>

            <div className="min-w-0 flex-1">
              <div className="mb-2 flex items-start gap-4">
                <div className="min-w-0 flex-1">
                  <h1 className="truncate text-4xl font-bold text-plex-text">{performer.name}</h1>
                  {performer.disambiguation && <p className="mt-1 text-sm text-plex-text-secondary">{performer.disambiguation}</p>}
                  {performer.aliases.length > 0 && (
                    <p className="mt-1 text-sm text-plex-text-secondary">Also known as: {performer.aliases.join(", ")}</p>
                  )}
                </div>
                <button
                  onClick={() => updateMut.mutate({ favorite: !performer.favorite })}
                  className={`rounded-full p-2 transition-colors ${
                    performer.favorite
                      ? "bg-red-500/15 text-red-500"
                      : "bg-plex-card text-plex-text-muted hover:text-red-400"
                  }`}
                  title={performer.favorite ? "Remove from favorites" : "Add to favorites"}
                >
                  <Heart className={`h-6 w-6 ${performer.favorite ? "fill-current" : ""}`} />
                </button>
              </div>

              <InteractiveRating value={performer.rating} onChange={(value) => updateMut.mutate({ rating: value })} />

              <div className="mt-4 flex flex-wrap gap-3">
                <CountCard label="Scenes" value={performer.sceneCount} icon={<Film className="h-4 w-4" />} />
                <CountCard label="Galleries" value={performer.galleryCount} icon={<FolderOpen className="h-4 w-4" />} />
                <CountCard label="Images" value={performer.imageCount} icon={<ImageIcon className="h-4 w-4" />} />
              </div>

              <div className="mt-4 grid grid-cols-2 gap-3 md:grid-cols-4">
                {performer.birthdate && (
                  <InfoItem icon={<Calendar className="h-4 w-4" />} label="Born" value={`${formatDate(performer.birthdate)}${age ? ` (${age})` : ""}`} />
                )}
                {performer.country && <InfoItem icon={<MapPin className="h-4 w-4" />} label="Country" value={performer.country} />}
                {performer.heightCm && <InfoItem icon={<Ruler className="h-4 w-4" />} label="Height" value={`${performer.heightCm} cm`} />}
                {performer.weight && <InfoItem icon={<Scale className="h-4 w-4" />} label="Weight" value={`${performer.weight} kg`} />}
              </div>

              {performer.urls.length > 0 && (
                <div className="mt-4 flex flex-wrap gap-2">
                  {performer.urls.map((url, i) => (
                    <a
                      key={i}
                      href={url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 rounded-full border border-plex-border bg-plex-card px-3 py-1 text-xs text-plex-accent hover:border-plex-accent/60 hover:text-plex-accent-hover"
                    >
                      <ExternalLink className="h-3 w-3" />
                      {(() => { try { return new URL(url).hostname.replace("www.", ""); } catch { return url; } })()}
                    </a>
                  ))}
                </div>
              )}

              {performer.tags.length > 0 && (
                <div className="mt-4 flex flex-wrap gap-1.5">
                  {performer.tags.map((tag) => (
                    <TagBadge key={tag.id} name={tag.name} onClick={() => onNavigate({ page: "tag", id: tag.id })} />
                  ))}
                </div>
              )}

              {performer.details && (
                <p className="mt-4 max-w-4xl whitespace-pre-wrap text-sm leading-6 text-plex-text-secondary">{performer.details}</p>
              )}
            </div>
          </div>
        </div>
      </div>

      <PerformerEditModal performer={performer} open={editing} onClose={() => setEditing(false)} />
      <ConfirmDialog
        open={confirmDelete}
        title="Delete Performer"
        message={`Are you sure you want to delete "${performer.name}"? This cannot be undone.`}
        onConfirm={() => deleteMut.mutate()}
        onCancel={() => setConfirmDelete(false)}
      />
      <MergeDialog
        open={mergeOpen}
        onClose={() => setMergeOpen(false)}
        entityType="performer"
        items={[{ id: performer.id, name: performer.name }]}
        onMerge={(targetId, sourceIds) => performers.merge(targetId, sourceIds)}
        queryKey={["performers"]}
      />

      <div className="mx-auto max-w-7xl px-4 py-6">
        <PerformerDetailsPanel performer={performer} />

        {performer.stashIds.length > 0 && (
          <div className="mt-4 rounded-xl border border-plex-border bg-plex-card p-4">
            <h2 className="mb-3 flex items-center gap-1.5 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">
              <Link2 className="h-4 w-4" /> Stash IDs
            </h2>
            <div className="flex flex-wrap gap-2">
              {performer.stashIds.map((sid) => (
                <span key={`${sid.endpoint}:${sid.stashId}`} className="inline-flex items-center gap-1.5 rounded-full border border-plex-border px-3 py-1 text-xs text-plex-text-secondary">
                  <Link2 className="h-3 w-3 text-plex-accent" />
                  {sid.stashId.slice(0, 12)}…
                </span>
              ))}
            </div>
          </div>
        )}

        <PerformerStashBoxPanel performer={performer} stashBoxes={config?.scraping.stashBoxes ?? []} onNavigate={onNavigate} />

        <div className="mt-8 border-b border-plex-border">
          <div className="flex gap-1 overflow-x-auto">
            {[
              { key: "scenes", label: "Scenes", count: performer.sceneCount },
              { key: "galleries", label: "Galleries", count: performer.galleryCount },
              { key: "images", label: "Images", count: performer.imageCount },
              { key: "groups", label: "Groups", count: performer.groupCount },
              { key: "appearsWith", label: "Appears With", count: undefined as number | undefined },
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
                {tab.count !== undefined && <span className="ml-2 rounded-full bg-plex-card px-2 py-0.5 text-xs text-plex-text-muted">{tab.count}</span>}
              </button>
            ))}
          </div>
        </div>

        <div className="py-6">
          {activeTab === "scenes" && (
            <PerformerScenesPanel performerId={id} filter={sceneFilter} setFilter={setSceneFilter} onNavigate={onNavigate} />
          )}
          {activeTab === "galleries" && (
            <PerformerGalleriesPanel performerId={id} filter={galleryFilter} setFilter={setGalleryFilter} onNavigate={onNavigate} />
          )}
          {activeTab === "images" && (
            <PerformerImagesPanel performerId={id} filter={imageFilter} setFilter={setImageFilter} onNavigate={onNavigate} />
          )}
          {activeTab === "groups" && (
            <PerformerGroupsPanel performerId={id} filter={groupFilter} setFilter={setGroupFilter} onNavigate={onNavigate} />
          )}
          {activeTab === "appearsWith" && (
            <PerformerAppearsWithPanel performerId={id} filter={appearsWithFilter} setFilter={setAppearsWithFilter} onNavigate={onNavigate} />
          )}
        </div>

        <div className="text-xs text-plex-text-muted">
          Created {new Date(performer.createdAt).toLocaleString()} · Updated {new Date(performer.updatedAt).toLocaleString()}
        </div>
        <ExtensionSlot slot="performer-detail-bottom" context={{ performer, onNavigate }} />
      </div>
    </div>
  );
}

function PerformerStashBoxPanel({ performer, stashBoxes, onNavigate }: { performer: PerformerModel; stashBoxes: StashBox[]; onNavigate: (r: any) => void }) {
  const queryClient = useQueryClient();
  const [term, setTerm] = useState(performer.name);
  const [selectedEndpoint, setSelectedEndpoint] = useState("");

  useEffect(() => {
    setTerm(performer.name);
  }, [performer.id, performer.name]);

  useEffect(() => {
    if (selectedEndpoint && !stashBoxes.some((box) => box.endpoint === selectedEndpoint)) {
      setSelectedEndpoint("");
    }
  }, [selectedEndpoint, stashBoxes]);

  const searchMutation = useMutation({
    mutationFn: (variables: { term?: string; endpoint?: string }) => performers.searchStashBox(performer.id, variables.term, variables.endpoint),
  });

  const importMutation = useMutation({
    mutationFn: (match: StashBoxPerformerMatch) =>
      performers.importFromStashBox(performer.id, { endpoint: match.endpoint, performerId: match.id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["performer", performer.id] });
      queryClient.invalidateQueries({ queryKey: ["performers"] });
    },
  });

  const linkedNames = new Map(stashBoxes.map((box) => [box.endpoint, box.name || box.endpoint]));

  return (
    <div className="mt-6 rounded-xl border border-plex-border bg-plex-card p-4">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div>
          <h2 className="text-base font-semibold text-plex-text">StashBox</h2>
          <p className="mt-1 max-w-2xl text-sm text-plex-text-secondary">
            Search configured StashBox instances and merge remote performer metadata into this local performer.
          </p>
        </div>
        {performer.stashIds.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {performer.stashIds.map((stashId) => (
              <span key={`${stashId.endpoint}:${stashId.stashId}`} className="inline-flex items-center gap-1 rounded-full border border-plex-border px-3 py-1 text-xs text-plex-text-secondary">
                <Link2 className="h-3.5 w-3.5 text-plex-accent" />
                {linkedNames.get(stashId.endpoint) ?? stashId.endpoint}
              </span>
            ))}
          </div>
        )}
      </div>

      {stashBoxes.length === 0 ? (
        <div className="mt-4 rounded-xl border border-dashed border-plex-border p-4 text-sm text-plex-text-secondary">
          No StashBox endpoints are configured yet. Use Settings and open Metadata Providers to add one.
          <button
            onClick={() => onNavigate({ page: "settings" })}
            className="ml-2 text-plex-accent hover:text-plex-accent-hover"
          >
            Open settings
          </button>
        </div>
      ) : (
        <>
          <div className="mt-4 grid gap-3 xl:grid-cols-[minmax(0,2fr)_minmax(0,1fr)_auto]">
            <label className="block text-sm">
              <span className="mb-1 block text-xs font-medium uppercase tracking-wide text-plex-text-muted">Search term</span>
              <input
                value={term}
                onChange={(event) => setTerm(event.target.value)}
                placeholder={performer.name}
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
            <p className="mt-3 text-sm text-red-300">{searchMutation.error.message}</p>
          )}

          {importMutation.isSuccess && (
            <p className="mt-3 text-sm text-emerald-300">Performer metadata imported from StashBox.</p>
          )}

          {searchMutation.data && (
            <div className="mt-4 space-y-3">
              {searchMutation.data.length === 0 ? (
                <div className="rounded-xl border border-dashed border-plex-border p-4 text-sm text-plex-text-secondary">
                  No StashBox performer matches were found.
                </div>
              ) : (
                searchMutation.data.map((match) => (
                  <button
                    key={`${match.endpoint}:${match.id}`}
                    onClick={() => importMutation.mutate(match)}
                    disabled={importMutation.isPending}
                    className="flex w-full flex-col gap-4 rounded-xl border border-plex-border bg-plex-surface p-4 text-left transition-colors hover:border-plex-accent/60 disabled:opacity-60 md:flex-row"
                  >
                    <div className="h-28 w-20 flex-shrink-0 overflow-hidden rounded-lg border border-plex-border bg-black/20">
                      {match.imageUrl ? (
                        <img src={match.imageUrl} alt={match.name} className="h-full w-full object-cover" />
                      ) : (
                        <div className="flex h-full w-full items-center justify-center bg-gradient-to-b from-plex-card to-plex-surface">
                          <UserRound className="h-10 w-10 text-plex-text-muted/50" />
                        </div>
                      )}
                    </div>

                    <div className="min-w-0 flex-1">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="text-base font-semibold text-plex-text">{match.name}</span>
                        <span className="rounded-full border border-plex-border px-2 py-0.5 text-xs text-plex-text-secondary">
                          {match.stashBoxName}
                        </span>
                        {match.deleted && <span className="rounded-full bg-red-500/15 px-2 py-0.5 text-xs text-red-300">Deleted</span>}
                      </div>
                      {match.disambiguation && <p className="mt-1 text-sm text-plex-text-secondary">{match.disambiguation}</p>}
                      <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs text-plex-text-muted">
                        {match.gender && <span>{match.gender}</span>}
                        {match.birthDate && <span>Born {match.birthDate}</span>}
                        {match.country && <span>{match.country}</span>}
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
        </>
      )}
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

function InfoItem({ icon, label, value }: { icon: React.ReactNode; label: string; value: string }) {
  return (
    <div className="flex items-center gap-2 text-sm">
      <span className="text-plex-text-muted">{icon}</span>
      <div>
        <div className="text-xs text-plex-text-muted">{label}</div>
        <div className="text-plex-text">{value}</div>
      </div>
    </div>
  );
}

function PerformerDetailsPanel({ performer }: { performer: any }) {
  const details = [
    { label: "Gender", value: performer.gender },
    { label: "Ethnicity", value: performer.ethnicity },
    { label: "Hair Color", value: performer.hairColor },
    { label: "Eye Color", value: performer.eyeColor },
    { label: "Measurements", value: performer.measurements },
    { label: "Fake Tits", value: performer.fakeTits },
    { label: "Circumcised", value: performer.circumcised },
    { label: "Penis Length", value: performer.penisLength ? `${performer.penisLength} cm` : null },
    { label: "Career", value: performer.careerStart ? `${performer.careerStart}${performer.careerEnd ? ` - ${performer.careerEnd}` : " - present"}` : null },
    { label: "Tattoos", value: performer.tattoos },
    { label: "Piercings", value: performer.piercings },
    { label: "Death Date", value: performer.deathDate ? formatDate(performer.deathDate) : null },
  ].filter((detail) => detail.value);

  if (details.length === 0) return null;

  return (
    <div className="rounded-xl border border-plex-border bg-plex-card p-4">
      <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-plex-text-muted">Details</h2>
      <dl className="grid grid-cols-2 gap-x-6 gap-y-3 text-sm md:grid-cols-4">
        {details.map((detail) => (
          <div key={detail.label}>
            <dt className="text-plex-text-muted">{detail.label}</dt>
            <dd className="text-plex-text">{detail.value}</dd>
          </div>
        ))}
      </dl>
    </div>
  );
}

function PerformerScenesPanel({ performerId, filter, setFilter, onNavigate }: {
  performerId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["performer-scenes", performerId, filter],
    queryFn: () => scenes.find(filter, { performerIds: String(performerId) }),
  });

  if (isLoading) return <LoadingPanel icon={<Film className="h-10 w-10" />} message="Loading scenes..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<Film className="h-12 w-12" />} message="No scenes found for this performer" />;

  return (
    <>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
        {data.items.map((scene) => (
          <SceneTile key={scene.id} scene={scene} onClick={() => onNavigate({ page: "scene", id: scene.id })} />
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={data.totalCount} />
    </>
  );
}

function PerformerGalleriesPanel({ performerId, filter, setFilter, onNavigate }: {
  performerId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["performer-galleries", performerId, filter],
    queryFn: () => galleries.find(filter, { performerIds: String(performerId) }),
  });

  if (isLoading) return <LoadingPanel icon={<FolderOpen className="h-10 w-10" />} message="Loading galleries..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<FolderOpen className="h-12 w-12" />} message="No galleries found for this performer" />;

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

function PerformerImagesPanel({ performerId, filter, setFilter, onNavigate }: {
  performerId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  const { data, isLoading } = useQuery({
    queryKey: ["performer-images", performerId, filter],
    queryFn: () => images.find(filter, { performerIds: String(performerId) }),
  });

  if (isLoading) return <LoadingPanel icon={<ImageIcon className="h-10 w-10" />} message="Loading images..." />;
  if (!data || data.items.length === 0) return <EmptyPanel icon={<ImageIcon className="h-12 w-12" />} message="No images found for this performer" />;

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

function PerformerGroupsPanel({ performerId, filter, setFilter, onNavigate }: {
  performerId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  // Fetch all scenes for this performer and extract unique groups
  const { data: scenesData, isLoading } = useQuery({
    queryKey: ["performer-scenes-for-groups", performerId],
    queryFn: () => scenes.find({ page: 1, perPage: 200, direction: "desc" }, { performerIds: String(performerId) }),
  });

  const uniqueGroups = useMemo(() => {
    if (!scenesData) return [];
    const seen = new Map<number, { id: number; name: string; sceneCount: number }>();
    for (const scene of scenesData.items) {
      for (const g of scene.groups ?? []) {
        const existing = seen.get(g.id);
        if (existing) {
          existing.sceneCount++;
        } else {
          seen.set(g.id, { id: g.id, name: g.name, sceneCount: 1 });
        }
      }
    }
    return [...seen.values()].sort((a, b) => b.sceneCount - a.sceneCount);
  }, [scenesData]);

  const page = filter.page ?? 1;
  const perPage = filter.perPage ?? 18;
  const paginated = uniqueGroups.slice((page - 1) * perPage, page * perPage);

  if (isLoading) return <LoadingPanel icon={<Layers className="h-10 w-10" />} message="Loading groups..." />;
  if (uniqueGroups.length === 0) return <EmptyPanel icon={<Layers className="h-12 w-12" />} message="No groups for this performer" />;

  return (
    <>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
        {paginated.map((group) => (
          <button key={group.id} onClick={() => onNavigate({ page: "group", id: group.id })} className="group overflow-hidden rounded-xl border border-plex-border bg-plex-card text-left transition-colors hover:border-plex-accent/60">
            <div className="flex aspect-video items-center justify-center bg-gradient-to-br from-plex-surface to-plex-card">
              <Layers className="h-10 w-10 text-plex-text-muted" />
            </div>
            <div className="p-3">
              <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">{group.name}</p>
              <p className="mt-1 text-xs text-plex-text-secondary">{group.sceneCount} scene{group.sceneCount !== 1 ? "s" : ""}</p>
            </div>
          </button>
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={uniqueGroups.length} />
    </>
  );
}

function PerformerAppearsWithPanel({ performerId, filter, setFilter, onNavigate }: {
  performerId: number;
  filter: FindFilter;
  setFilter: (filter: FindFilter) => void;
  onNavigate: (r: any) => void;
}) {
  // Find all scenes this performer is in, collect co-performers
  const { data: scenesData, isLoading } = useQuery({
    queryKey: ["performer-scenes-for-costars", performerId],
    queryFn: () => scenes.find({ page: 1, perPage: 200, direction: "desc" }, { performerIds: String(performerId) }),
  });

  const coStars = useMemo(() => {
    if (!scenesData) return [];
    const counts = new Map<number, { performer: { id: number; name: string; imagePath?: string }; count: number }>();
    for (const scene of scenesData.items) {
      for (const p of scene.performers ?? []) {
        if (p.id === performerId) continue;
        const existing = counts.get(p.id);
        if (existing) {
          existing.count++;
        } else {
          counts.set(p.id, { performer: { id: p.id, name: p.name, imagePath: p.imagePath }, count: 1 });
        }
      }
    }
    return [...counts.values()].sort((a, b) => b.count - a.count);
  }, [scenesData, performerId]);

  // Paginate client-side
  const page = filter.page ?? 1;
  const perPage = filter.perPage ?? 18;
  const paginated = coStars.slice((page - 1) * perPage, page * perPage);

  if (isLoading) return <LoadingPanel icon={<Users className="h-10 w-10" />} message="Loading co-stars..." />;
  if (coStars.length === 0) return <EmptyPanel icon={<Users className="h-12 w-12" />} message="No co-stars found" />;

  return (
    <>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
        {paginated.map(({ performer: p, count }) => (
          <button key={p.id} onClick={() => onNavigate({ page: "performer", id: p.id })} className="group overflow-hidden rounded-xl border border-plex-border bg-plex-card text-left transition-colors hover:border-plex-accent/60">
            <div className="aspect-[2/3] bg-gradient-to-b from-plex-card to-plex-surface overflow-hidden">
              <img
                src={entityImages.performerImageUrl(p.id)}
                alt={p.name}
                className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
                loading="lazy"
                onError={(e) => { (e.target as HTMLImageElement).style.display = "none"; }}
              />
            </div>
            <div className="p-3">
              <p className="truncate text-sm font-medium text-plex-text group-hover:text-plex-accent">{p.name}</p>
              <p className="mt-1 text-xs text-plex-text-secondary">{count} scene{count !== 1 ? "s" : ""} together</p>
            </div>
          </button>
        ))}
      </div>
      <Pager filter={filter} setFilter={setFilter} totalCount={coStars.length} />
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