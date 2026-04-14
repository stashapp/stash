import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { studios } from "../api/client";
import type { FindFilter, Studio, StudioCreate, StudioFilterCriteria } from "../api/types";
import { ListPage, type DisplayMode } from "../components/ListPage";
import { RatingBanner, RatingField } from "../components/Rating";
import { EditModal, Field, TextInput, TextArea, SaveButton } from "../components/EditModal";
import { useMultiSelect } from "../hooks/useMultiSelect";
import { Building2, Film, Image, LayoutGrid, Trash2, Loader2, Check, Edit, Merge, Heart, Box, Users, Layers, Tag as TagIcon } from "lucide-react";
import { STUDIO_CRITERIA } from "../components/FilterDialog";
import { BulkEditDialog } from "../components/BulkEditDialog";
import { MergeDialog } from "../components/MergeDialog";
import { StudioTagger } from "../components/StudioTagger";

const SORT_OPTIONS = [
  { value: "name", label: "Name" },
  { value: "scene_count", label: "Scene Count" },
  { value: "random", label: "Random" },
  { value: "created_at", label: "Recently Added" },
];

interface Props {
  onNavigate: (r: any) => void;
}

export function StudiosPage({ onNavigate }: Props) {
  const [filter, setFilter] = useState<FindFilter>({ page: 1, perPage: 40, sort: "name", direction: "asc" });
  const [displayMode, setDisplayMode] = useState<DisplayMode>("grid");
  const [selectionMode, setSelectionMode] = useState(false);
  const [showCreate, setShowCreate] = useState(false);
  const [objectFilter, setObjectFilter] = useState<Record<string, unknown>>({});
  const [showBulkEdit, setShowBulkEdit] = useState(false);
  const [showMerge, setShowMerge] = useState(false);
  const queryClient = useQueryClient();

  const hasObjectFilter = Object.keys(objectFilter).length > 0;
  const { data, isLoading } = useQuery({
    queryKey: ["studios", filter, objectFilter],
    queryFn: () =>
      hasObjectFilter
        ? studios.findFiltered({ findFilter: filter, objectFilter: objectFilter as StudioFilterCriteria })
        : studios.find(filter),
  });

  const items = data?.items ?? [];
  const { selectedIds, toggle, selectAll, selectNone } = useMultiSelect(items);
  const selecting = selectionMode || selectedIds.size > 0;

  const bulkDeleteMut = useMutation({
    mutationFn: () => studios.bulkDelete([...selectedIds]),
    onSuccess: () => { selectNone(); queryClient.invalidateQueries({ queryKey: ["studios"] }); },
  });

  const bulkEditMut = useMutation({
    mutationFn: (values: Record<string, unknown>) =>
      studios.bulkUpdate({ ids: [...selectedIds], ...values } as any),
    onSuccess: () => {
      setShowBulkEdit(false);
      selectNone();
      queryClient.invalidateQueries({ queryKey: ["studios"] });
    },
  });

  return (
    <>
      <StudioCreateModal open={showCreate} onClose={() => setShowCreate(false)} onCreated={(id) => onNavigate({ page: "studio", id })} />
      <ListPage
        title="Studios"
        filterMode="studios"
        filter={filter}
        onFilterChange={setFilter}
        totalCount={data?.totalCount ?? 0}
        isLoading={isLoading}
        sortOptions={SORT_OPTIONS}
        displayMode={displayMode}
        onDisplayModeChange={setDisplayMode}
        availableDisplayModes={["grid", "list", "tagger"]}
        criteriaDefinitions={STUDIO_CRITERIA}
        objectFilter={objectFilter}
        onObjectFilterChange={setObjectFilter}
        onNew={() => setShowCreate(true)}
        renderOperations={() => (
          <button
            onClick={() => {
              if (selecting) {
                selectNone();
                setSelectionMode(false);
                return;
              }
              setSelectionMode(true);
            }}
            className="flex items-center gap-1 px-2 py-1 rounded text-xs border border-plex-border bg-plex-input text-plex-text-secondary hover:text-plex-text"
          >
            <Check className="w-3.5 h-3.5" />
            {selecting ? "Cancel" : "Select"}
          </button>
        )}
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
              onClick={() => { if (confirm(`Delete ${selectedIds.size} studio(s)?`)) bulkDeleteMut.mutate(); }}
              disabled={bulkDeleteMut.isPending}
              className="flex items-center gap-1 px-2 py-0.5 rounded text-xs text-red-400 hover:text-red-300 hover:bg-red-900/20"
            >
              {bulkDeleteMut.isPending ? <Loader2 className="w-3 h-3 animate-spin" /> : <Trash2 className="w-3 h-3" />}
              Delete
            </button>
          </>
        }
      >
      {displayMode === "tagger" ? (
        <StudioTagger studios={items} />
      ) : displayMode === "grid" ? (
        <div className="grid gap-3" style={{ gridTemplateColumns: "repeat(auto-fill, minmax(var(--card-min-width, 200px), 1fr))" }}>
          {items.map((s) => (
            <StudioCard
              key={s.id}
              studio={s}
              onClick={() => selecting ? toggle(s.id) : onNavigate({ page: "studio", id: s.id })}
              selected={selectedIds.has(s.id)}
              onSelect={() => toggle(s.id)}
              selecting={selecting}
            />
          ))}
        </div>
      ) : (
        <StudioListTable studios={items} onNavigate={onNavigate} selectedIds={selectedIds} onToggle={toggle} selecting={selecting} />
      )}
      {items.length === 0 && (
        <div className="text-center text-plex-text-secondary py-16">
          <Building2 className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>No studios found</p>
        </div>
      )}
      </ListPage>
      <BulkEditDialog
        open={showBulkEdit}
        onClose={() => setShowBulkEdit(false)}
        title="Edit Studios"
        selectedCount={selectedIds.size}
        fields={[{ key: "rating", label: "Rating", type: "number" }, { key: "favorite", label: "Favorite", type: "bool" }]}
        onApply={(values) => bulkEditMut.mutate(values)}
        isPending={bulkEditMut.isPending}
      />
      <MergeDialog
        open={showMerge}
        onClose={() => { setShowMerge(false); selectNone(); }}
        entityType="studio"
        items={items.filter((s) => selectedIds.has(s.id)).map((s) => ({ id: s.id, name: s.name }))}
        onMerge={studios.merge}
        queryKey="studios"
      />
    </>
  );
}

function StudioCard({ studio, onClick, selected, onSelect, selecting }: { studio: Studio; onClick: () => void; selected?: boolean; onSelect?: () => void; selecting?: boolean }) {
  const queryClient = useQueryClient();
  const favMut = useMutation({
    mutationFn: () => studios.update(studio.id, { favorite: !studio.favorite }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["studios"] }),
  });

  return (
    <div className={`bg-plex-card rounded overflow-hidden border hover:border-plex-accent/60 transition-all cursor-pointer relative group ${selected ? "border-plex-accent ring-2 ring-plex-accent" : "border-plex-border"}`} onClick={onClick}>
      <div className="aspect-video bg-plex-surface flex items-center justify-center text-plex-text-muted relative overflow-hidden">
        <div className={`absolute top-1 left-1 z-10 ${selected || selecting ? "opacity-100" : "opacity-0 group-hover:opacity-100"} transition-opacity`}>
          <input type="checkbox" checked={selected} onChange={(e) => { e.stopPropagation(); onSelect?.(); }} onClick={(e) => e.stopPropagation()} className="w-4 h-4 rounded border-plex-border cursor-pointer accent-plex-accent" />
        </div>
        {/* Favorite heart overlay */}
        <button
          onClick={(e) => { e.stopPropagation(); favMut.mutate(); }}
          className={`absolute top-1 right-1 p-1 z-10 transition-opacity ${studio.favorite ? "opacity-100" : "opacity-0 group-hover:opacity-70"}`}
          title={studio.favorite ? "Unfavorite" : "Favorite"}
        >
          <Heart className={`w-5 h-5 ${studio.favorite ? "fill-red-500 text-red-500" : "text-white drop-shadow-md"}`} />
        </button>
        {studio.imagePath ? (
          <img src={studio.imagePath} alt={studio.name} className="w-full h-full object-contain p-4" loading="lazy" />
        ) : (
          <Building2 className="w-10 h-10 opacity-30" />
        )}
        <RatingBanner rating={studio.rating} />
      </div>
      <div className="p-2 text-center">
        <h3 className="font-medium text-sm truncate text-plex-text">{studio.name}</h3>
        {studio.parentName && (
          <div className="text-xs text-plex-text-muted truncate">↑ {studio.parentName}</div>
        )}
      </div>
      {(studio.sceneCount > 0 || studio.imageCount > 0 || studio.galleryCount > 0 || studio.groupCount > 0 || studio.performerCount > 0 || studio.tags.length > 0 || studio.childStudioCount > 0 || studio.organized) && (
        <div className="flex items-center justify-center gap-1 px-2 pb-2 border-t border-plex-border pt-1.5 flex-wrap">
          {studio.sceneCount > 0 && (
            <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Scenes">
              <Film className="w-3 h-3" /> {studio.sceneCount}
            </span>
          )}
          {studio.groupCount > 0 && (
            <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Groups">
              <Layers className="w-3 h-3" /> {studio.groupCount}
            </span>
          )}
          {studio.imageCount > 0 && (
            <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Images">
              <Image className="w-3 h-3" /> {studio.imageCount}
            </span>
          )}
          {studio.galleryCount > 0 && (
            <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Galleries">
              <LayoutGrid className="w-3 h-3" /> {studio.galleryCount}
            </span>
          )}
          {studio.performerCount > 0 && (
            <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Performers">
              <Users className="w-3 h-3" /> {studio.performerCount}
            </span>
          )}
          {studio.tags.length > 0 && (
            <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Tags">
              <TagIcon className="w-3 h-3" /> {studio.tags.length}
            </span>
          )}
          {studio.childStudioCount > 0 && (
            <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Sub-Studios">
              <Building2 className="w-3 h-3" /> {studio.childStudioCount}
            </span>
          )}
          {studio.organized && (
            <span className="text-plex-text-muted" title="Organized">
              <Box className="w-3 h-3" />
            </span>
          )}
        </div>
      )}
    </div>
  );
}

function StudioListTable({ studios: items, onNavigate, selectedIds, onToggle, selecting }: { studios: Studio[]; onNavigate: (r: any) => void; selectedIds?: Set<number>; onToggle?: (id: number) => void; selecting?: boolean }) {
  return (
    <table className="w-full text-sm">
      <thead>
        <tr className="border-b border-plex-border text-left text-plex-text-muted text-xs">
          {selectedIds && <th className="w-8 py-2 px-3"></th>}
          <th className="py-2 px-3">Name</th>
          <th className="py-2 px-3">Parent</th>
          <th className="py-2 px-3 text-right">Scenes</th>
          <th className="py-2 px-3 text-right">Rating</th>
        </tr>
      </thead>
      <tbody>
        {items.map((s) => (
          <tr
            key={s.id}
            onClick={() => selecting ? onToggle?.(s.id) : onNavigate({ page: "studio", id: s.id })}
            className={`border-b border-plex-border hover:bg-plex-card cursor-pointer ${selectedIds?.has(s.id) ? "bg-plex-accent/10" : ""}`}
          >
            {selectedIds && <td className="py-2 px-3"><input type="checkbox" checked={selectedIds.has(s.id)} onChange={() => onToggle?.(s.id)} onClick={(e) => e.stopPropagation()} className="w-3.5 h-3.5 rounded border-plex-border cursor-pointer accent-plex-accent" /></td>}
            <td className="py-2 px-3 text-plex-text">{s.name}</td>
            <td className="py-2 px-3 text-plex-text-secondary">{s.parentName ?? ""}</td>
            <td className="py-2 px-3 text-plex-text-secondary text-right">{s.sceneCount}</td>
            <td className="py-2 px-3 text-plex-text-secondary text-right">{s.rating ?? ""}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

/* ── Studio Create Modal ── */
function StudioCreateModal({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const qc = useQueryClient();
  const [form, setForm] = useState({
    name: "",
    details: "",
    rating: undefined as number | undefined,
    favorite: false,
    ignoreAutoTag: false,
    organized: false,
  });

  const mutation = useMutation({
    mutationFn: (data: StudioCreate) => studios.create(data),
    onSuccess: (created) => {
      qc.invalidateQueries({ queryKey: ["studios"] });
      setForm({ name: "", details: "", rating: undefined, favorite: false, ignoreAutoTag: false, organized: false });
      onClose();
      if (created?.id) onCreated(created.id);
    },
  });

  const save = () => {
    const name = form.name.trim();
    if (!name) return;
    mutation.mutate({
      name,
      details: form.details || undefined,
      rating: form.rating,
      favorite: form.favorite || undefined,
      ignoreAutoTag: form.ignoreAutoTag || undefined,
      organized: form.organized || undefined,
    });
  };

  return (
    <EditModal title="Create Studio" open={open} onClose={onClose}>
      <Field label="Name">
        <TextInput value={form.name} onChange={(v) => setForm({ ...form, name: v })} />
      </Field>
      <Field label="Details">
        <TextArea value={form.details} onChange={(v) => setForm({ ...form, details: v })} rows={3} />
      </Field>
      <RatingField value={form.rating} onChange={(value) => setForm({ ...form, rating: value })} />
      <div className="flex items-center gap-4 mb-4">
        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            checked={form.favorite}
            onChange={(e) => setForm({ ...form, favorite: e.target.checked })}
            className="rounded bg-gray-800 border-gray-700"
          />
          Favorite
        </label>
        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            checked={form.ignoreAutoTag}
            onChange={(e) => setForm({ ...form, ignoreAutoTag: e.target.checked })}
            className="rounded bg-gray-800 border-gray-700"
          />
          Ignore Auto Tag
        </label>
        <label className="flex items-center gap-2 text-sm">
          <input
            type="checkbox"
            checked={form.organized}
            onChange={(e) => setForm({ ...form, organized: e.target.checked })}
            className="rounded bg-gray-800 border-gray-700"
          />
          Organized
        </label>
      </div>
      <div className="flex justify-end mt-4">
        <SaveButton loading={mutation.isPending} onClick={save} />
      </div>
    </EditModal>
  );
}
