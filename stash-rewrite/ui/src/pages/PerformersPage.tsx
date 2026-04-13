import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { performers, entityImages } from "../api/client";
import type { FindFilter, Performer, PerformerCreate, PerformerFilterCriteria } from "../api/types";
import { ListPage, type DisplayMode } from "../components/ListPage";
import { RatingBanner, RatingField } from "../components/Rating";
import { EditModal, Field, TextInput, TextArea, SaveButton } from "../components/EditModal";
import { useMultiSelect } from "../hooks/useMultiSelect";
import { PERFORMER_CRITERIA } from "../components/FilterDialog";
import { BulkEditDialog, PERFORMER_BULK_FIELDS } from "../components/BulkEditDialog";
import { Users, Heart, Tag, Film, Image, LayoutGrid, Trash2, Loader2, Check, Edit, Merge } from "lucide-react";
import { MergeDialog } from "../components/MergeDialog";
import { PerformerTagger } from "../components/PerformerTagger";

/** Convert 2-letter ISO country code to flag emoji */
function countryToFlag(code: string): string {
  const upper = code.toUpperCase();
  if (upper.length !== 2) return code;
  return String.fromCodePoint(...[...upper].map(c => 0x1F1E6 + c.charCodeAt(0) - 65));
}

const SORT_OPTIONS = [
  { value: "name", label: "Name" },
  { value: "scene_count", label: "Scene Count" },
  { value: "image_count", label: "Image Count" },
  { value: "gallery_count", label: "Gallery Count" },
  { value: "tag_count", label: "Tag Count" },
  { value: "rating", label: "Rating" },
  { value: "birthdate", label: "Birthdate" },
  { value: "height", label: "Height" },
  { value: "weight", label: "Weight" },
  { value: "created_at", label: "Recently Added" },
  { value: "updated_at", label: "Recently Updated" },
  { value: "random", label: "Random" },
];

interface Props {
  onNavigate: (r: any) => void;
}

export function PerformersPage({ onNavigate }: Props) {
  const [filter, setFilter] = useState<FindFilter>({ page: 1, perPage: 40, sort: "name", direction: "asc" });
  const [objectFilter, setObjectFilter] = useState<Record<string, unknown>>({});
  const [displayMode, setDisplayMode] = useState<DisplayMode>("grid");
  const [selectionMode, setSelectionMode] = useState(false);
  const [showCreate, setShowCreate] = useState(false);
  const [showBulkEdit, setShowBulkEdit] = useState(false);
  const [showMerge, setShowMerge] = useState(false);
  const queryClient = useQueryClient();

  const hasObjectFilter = Object.keys(objectFilter).length > 0;

  const { data, isLoading } = useQuery({
    queryKey: ["performers", filter, objectFilter],
    queryFn: () =>
      hasObjectFilter
        ? performers.findFiltered({ findFilter: filter, objectFilter: objectFilter as PerformerFilterCriteria })
        : performers.find(filter),
  });

  const items = data?.items ?? [];
  const { selectedIds, toggle, selectAll, selectNone } = useMultiSelect(items);
  const selecting = selectionMode || selectedIds.size > 0;

  const bulkDeleteMut = useMutation({
    mutationFn: () => performers.bulkDelete([...selectedIds]),
    onSuccess: () => { selectNone(); queryClient.invalidateQueries({ queryKey: ["performers"] }); },
  });

  const bulkEditMut = useMutation({
    mutationFn: (values: Record<string, unknown>) =>
      performers.bulkUpdate({ ids: [...selectedIds], ...values } as any),
    onSuccess: () => {
      setShowBulkEdit(false);
      selectNone();
      queryClient.invalidateQueries({ queryKey: ["performers"] });
    },
  });

  return (
    <>
      <PerformerCreateModal open={showCreate} onClose={() => setShowCreate(false)} onCreated={(id) => onNavigate({ page: "performer", id })} />
      <ListPage
        title="Performers"
        filterMode="performers"
        filter={filter}
        onFilterChange={setFilter}
        totalCount={data?.totalCount ?? 0}
        isLoading={isLoading}
        sortOptions={SORT_OPTIONS}
        displayMode={displayMode}
        onDisplayModeChange={setDisplayMode}
        availableDisplayModes={["grid", "list", "tagger"]}
        onNew={() => setShowCreate(true)}
        criteriaDefinitions={PERFORMER_CRITERIA}
        objectFilter={objectFilter}
        onObjectFilterChange={setObjectFilter}
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
              onClick={() => { if (confirm(`Delete ${selectedIds.size} performer(s)?`)) bulkDeleteMut.mutate(); }}
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
        <PerformerTagger performers={items} />
      ) : displayMode === "grid" ? (
        <div className="grid gap-3" style={{ gridTemplateColumns: "repeat(auto-fill, minmax(var(--card-min-width, 160px), 1fr))" }}>
          {items.map((p) => (
            <PerformerCard
              key={p.id}
              performer={p}
              onClick={() => selecting ? toggle(p.id) : onNavigate({ page: "performer", id: p.id })}
              selected={selectedIds.has(p.id)}
              onSelect={() => toggle(p.id)}
              selecting={selecting}
            />
          ))}
        </div>
      ) : (
        <PerformerListTable performers={items} onNavigate={onNavigate} selectedIds={selectedIds} onToggle={toggle} selecting={selecting} />
      )}
      {items.length === 0 && (
        <div className="text-center text-plex-text-secondary py-16">
          <Users className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>No performers found</p>
        </div>
      )}
      </ListPage>

      {/* Bulk Edit Dialog */}
      <BulkEditDialog
        open={showBulkEdit}
        onClose={() => setShowBulkEdit(false)}
        title="Edit Performers"
        selectedCount={selectedIds.size}
        fields={PERFORMER_BULK_FIELDS}
        onApply={(values) => bulkEditMut.mutate(values)}
        isPending={bulkEditMut.isPending}
      />
      <MergeDialog
        open={showMerge}
        onClose={() => { setShowMerge(false); selectNone(); }}
        entityType="performer"
        items={items.filter((p) => selectedIds.has(p.id)).map((p) => ({ id: p.id, name: p.name, imagePath: p.imagePath }))}
        onMerge={performers.merge}
        queryKey="performers"
      />
    </>
  );
}

function PerformerCard({ performer, onClick, selected, onSelect, selecting }: { performer: Performer; onClick: () => void; selected?: boolean; onSelect?: () => void; selecting?: boolean }) {
  const queryClient = useQueryClient();
  const age = performer.birthdate
    ? Math.floor((Date.now() - new Date(performer.birthdate).getTime()) / 31557600000)
    : null;

  const favMut = useMutation({
    mutationFn: () => performers.update(performer.id, { favorite: !performer.favorite }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["performers"] }),
  });

  return (
    <div
      onClick={onClick}
      className={`bg-plex-card rounded overflow-hidden cursor-pointer border hover:border-plex-accent/60 transition-all text-left w-full group relative ${selected ? "border-plex-accent ring-2 ring-plex-accent" : "border-plex-border"}`}
    >
      {/* Image */}
      <div className="aspect-[2/3] bg-plex-surface relative overflow-hidden">
        <div className={`absolute top-1 left-1 z-10 ${selected || selecting ? "opacity-100" : "opacity-0 group-hover:opacity-100"} transition-opacity`}>
          <input
            type="checkbox"
            checked={selected}
            onChange={(e) => { e.stopPropagation(); onSelect?.(); }}
            onClick={(e) => e.stopPropagation()}
            className="w-4 h-4 rounded border-plex-border cursor-pointer accent-plex-accent"
          />
        </div>
        <img 
            src={performer.imagePath || entityImages.performerImageUrl(performer.id)} 
            alt={performer.name} 
            className="w-full h-full object-cover" 
            loading="lazy"
            onError={(e) => {
              const el = e.target as HTMLImageElement;
              el.style.display = "none";
              if (el.nextElementSibling) (el.nextElementSibling as HTMLElement).style.display = "";
            }}
          />
          <div className="w-full h-full flex items-center justify-center" style={{ display: "none" }}>
            <svg viewBox="0 0 100 150" className="w-2/3 h-2/3 opacity-20">
              <ellipse cx="50" cy="35" rx="25" ry="30" fill="currentColor" className="text-plex-text-muted"/>
              <ellipse cx="50" cy="120" rx="40" ry="45" fill="currentColor" className="text-plex-text-muted"/>
            </svg>
          </div>

        {/* Favorite star overlay */}
        <button
          onClick={(e) => { e.stopPropagation(); favMut.mutate(); }}
          className={`absolute top-1 right-1 p-1 transition-opacity ${
            performer.favorite ? "opacity-100" : "opacity-0 group-hover:opacity-70"
          }`}
          title={performer.favorite ? "Unfavorite" : "Favorite"}
        >
          <Heart className={`w-5 h-5 ${performer.favorite ? "fill-red-500 text-red-500" : "text-white drop-shadow-md"}`} />
        </button>

        {/* Rating banner */}
        <RatingBanner rating={performer.rating} />

        {/* Country flag overlay - bottom right */}
        {performer.country && (
          <div className="absolute bottom-1 right-1 text-lg leading-none drop-shadow-md" title={performer.country}>
            {countryToFlag(performer.country)}
          </div>
        )}
      </div>

      {/* Title */}
      <div className="p-2 text-center">
        <div className="text-sm text-plex-text font-medium truncate">
          {performer.gender && (
            <span className="text-plex-text-muted mr-1" title={performer.gender}>
              {performer.gender === "FEMALE" ? "♀" : performer.gender === "MALE" ? "♂" : performer.gender === "NON_BINARY" ? "⚧" : "⚧"}
            </span>
          )}
          {performer.name}
          {performer.disambiguation && (
            <span className="text-plex-text-muted font-normal"> ({performer.disambiguation})</span>
          )}
        </div>
        {age !== null && age > 0 && (
          <div className="text-xs text-plex-text-secondary">{age} years old</div>
        )}
      </div>

      {/* Popovers row */}
      <PerformerCardPopovers performer={performer} />
    </div>
  );
}

function PerformerCardPopovers({ performer }: { performer: Performer }) {
  const hasAny = performer.sceneCount > 0 || performer.imageCount > 0 || performer.galleryCount > 0 || performer.groupCount > 0 || performer.tags.length > 0;
  if (!hasAny) return null;

  return (
    <div className="flex items-center justify-center gap-1 px-2 pb-2 border-t border-plex-border pt-1.5">
      {performer.sceneCount > 0 && (
        <span className="flex items-center gap-0.5 text-xs text-plex-text-muted px-1" title="Scenes">
          <Film className="w-3 h-3" /> {performer.sceneCount}
        </span>
      )}
      {performer.imageCount > 0 && (
        <span className="flex items-center gap-0.5 text-xs text-plex-text-muted px-1" title="Images">
          <Image className="w-3 h-3" /> {performer.imageCount}
        </span>
      )}
      {performer.galleryCount > 0 && (
        <span className="flex items-center gap-0.5 text-xs text-plex-text-muted px-1" title="Galleries">
          <LayoutGrid className="w-3 h-3" /> {performer.galleryCount}
        </span>
      )}
      {performer.groupCount > 0 && (
        <span className="flex items-center gap-0.5 text-xs text-plex-text-muted px-1" title="Groups">
          <Film className="w-3 h-3" /> {performer.groupCount}
        </span>
      )}
      {performer.tags.length > 0 && (
        <span className="flex items-center gap-0.5 text-xs text-plex-text-muted px-1" title="Tags">
          <Tag className="w-3 h-3" /> {performer.tags.length}
        </span>
      )}
    </div>
  );
}

function PerformerListTable({ performers: items, onNavigate, selectedIds, onToggle, selecting }: { performers: Performer[]; onNavigate: (r: any) => void; selectedIds?: Set<number>; onToggle?: (id: number) => void; selecting?: boolean }) {
  return (
    <table className="w-full text-sm">
      <thead>
        <tr className="border-b border-plex-border text-left text-plex-text-muted text-xs">
          {selectedIds && <th className="w-8 py-2 px-3"></th>}
          <th className="py-2 px-3">Name</th>
          <th className="py-2 px-3">Gender</th>
          <th className="py-2 px-3">Age</th>
          <th className="py-2 px-3">Country</th>
          <th className="py-2 px-3 text-right">Scenes</th>
          <th className="py-2 px-3 text-right">Rating</th>
          <th className="py-2 px-3">Favorite</th>
        </tr>
      </thead>
      <tbody>
        {items.map((p) => {
          const age = p.birthdate
            ? Math.floor((Date.now() - new Date(p.birthdate).getTime()) / 31557600000)
            : null;
          return (
            <tr 
              key={p.id} 
              onClick={() => selecting ? onToggle?.(p.id) : onNavigate({ page: "performer", id: p.id })}
              className={`border-b border-plex-border hover:bg-plex-card cursor-pointer ${selectedIds?.has(p.id) ? "bg-plex-accent/10" : ""}`}
            >
              {selectedIds && (
                <td className="py-2 px-3">
                  <input type="checkbox" checked={selectedIds.has(p.id)} onChange={() => onToggle?.(p.id)} onClick={(e) => e.stopPropagation()} className="w-3.5 h-3.5 rounded border-plex-border cursor-pointer accent-plex-accent" />
                </td>
              )}
              <td className="py-2 px-3 text-plex-text">
                {p.name}
                {p.disambiguation && <span className="text-plex-text-muted ml-1">({p.disambiguation})</span>}
              </td>
              <td className="py-2 px-3 text-plex-text-secondary capitalize">{p.gender?.toLowerCase()}</td>
              <td className="py-2 px-3 text-plex-text-secondary">{age ?? ""}</td>
              <td className="py-2 px-3 text-plex-text-secondary">{p.country ?? ""}</td>
              <td className="py-2 px-3 text-plex-text-secondary text-right">{p.sceneCount}</td>
              <td className="py-2 px-3 text-plex-text-secondary text-right">{p.rating ?? ""}</td>
              <td className="py-2 px-3">
                {p.favorite && <Heart className="w-4 h-4 fill-red-500 text-red-500" />}
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}

/* ── Performer Create Modal ── */
function PerformerCreateModal({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const qc = useQueryClient();
  const [form, setForm] = useState({
    name: "",
    disambiguation: "",
    gender: "",
    details: "",
    favorite: false,
    ignoreAutoTag: false,
    rating: undefined as number | undefined,
  });

  const mutation = useMutation({
    mutationFn: (data: PerformerCreate) => performers.create(data),
    onSuccess: (created) => {
      qc.invalidateQueries({ queryKey: ["performers"] });
      setForm({ name: "", disambiguation: "", gender: "", details: "", favorite: false, ignoreAutoTag: false, rating: undefined });
      onClose();
      if (created?.id) onCreated(created.id);
    },
  });

  const save = () => {
    const name = form.name.trim();
    if (!name) return;
    mutation.mutate({
      name,
      disambiguation: form.disambiguation || undefined,
      gender: form.gender || undefined,
      details: form.details || undefined,
      favorite: form.favorite || undefined,
      ignoreAutoTag: form.ignoreAutoTag || undefined,
      rating: form.rating,
    });
  };

  return (
    <EditModal title="Create Performer" open={open} onClose={onClose}>
      <Field label="Name">
        <TextInput value={form.name} onChange={(v) => setForm({ ...form, name: v })} />
      </Field>
      <Field label="Disambiguation">
        <TextInput value={form.disambiguation} onChange={(v) => setForm({ ...form, disambiguation: v })} />
      </Field>
      <Field label="Gender">
        <TextInput value={form.gender} onChange={(v) => setForm({ ...form, gender: v })} />
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
      </div>
      <div className="flex justify-end mt-4">
        <SaveButton loading={mutation.isPending} onClick={save} />
      </div>
    </EditModal>
  );
}
