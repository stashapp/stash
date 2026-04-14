import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { galleries } from "../api/client";
import type { FindFilter, Gallery, GalleryCreate, GalleryFilterCriteria } from "../api/types";
import { ListPage, type DisplayMode } from "../components/ListPage";
import { RatingBanner, RatingField } from "../components/Rating";
import { EditModal, Field, TextInput, TextArea, SaveButton } from "../components/EditModal";
import { useMultiSelect } from "../hooks/useMultiSelect";
import { FolderOpen, Image, Users, Tag, Trash2, Loader2, Edit, Box, Film } from "lucide-react";
import { GALLERY_CRITERIA } from "../components/FilterDialog";
import { BulkEditDialog, GALLERY_BULK_FIELDS } from "../components/BulkEditDialog";

const SORT_OPTIONS = [
  { value: "title", label: "Title" },
  { value: "image_count", label: "Image Count" },
  { value: "rating", label: "Rating" },
  { value: "random", label: "Random" },
  { value: "created_at", label: "Recently Added" },
];

interface Props {
  onNavigate: (r: any) => void;
}

export function GalleriesPage({ onNavigate }: Props) {
  const [filter, setFilter] = useState<FindFilter>({ page: 1, perPage: 40, direction: "desc" });
  const [displayMode, setDisplayMode] = useState<DisplayMode>("grid");
  const [selectionMode, setSelectionMode] = useState(false);
  const [objectFilter, setObjectFilter] = useState<Record<string, unknown>>({});
  const [showBulkEdit, setShowBulkEdit] = useState(false);
  const [showCreate, setShowCreate] = useState(false);
  const queryClient = useQueryClient();

  const hasObjectFilter = Object.keys(objectFilter).length > 0;
  const { data, isLoading } = useQuery({
    queryKey: ["galleries", filter, objectFilter],
    queryFn: () =>
      hasObjectFilter
        ? galleries.findFiltered({ findFilter: filter, objectFilter: objectFilter as GalleryFilterCriteria })
        : galleries.find(filter),
  });

  const items = data?.items ?? [];
  const { selectedIds, toggle, selectAll, selectNone } = useMultiSelect(items);
  const selecting = selectionMode || selectedIds.size > 0;

  const bulkDeleteMut = useMutation({
    mutationFn: () => galleries.bulkDelete([...selectedIds]),
    onSuccess: () => { selectNone(); queryClient.invalidateQueries({ queryKey: ["galleries"] }); },
  });

  const bulkEditMut = useMutation({
    mutationFn: (values: Record<string, unknown>) =>
      galleries.bulkUpdate({ ids: [...selectedIds], ...values } as any),
    onSuccess: () => {
      setShowBulkEdit(false);
      selectNone();
      queryClient.invalidateQueries({ queryKey: ["galleries"] });
    },
  });

  return (
    <>
    <GalleryCreateModal open={showCreate} onClose={() => setShowCreate(false)} onCreated={(id) => onNavigate({ page: "gallery", id })} />
    <ListPage
      title="Galleries"
      filterMode="galleries"
      filter={filter}
      onFilterChange={setFilter}
      totalCount={data?.totalCount ?? 0}
      isLoading={isLoading}
      sortOptions={SORT_OPTIONS}
      displayMode={displayMode}
      onDisplayModeChange={setDisplayMode}
      availableDisplayModes={["grid", "list"]}
      criteriaDefinitions={GALLERY_CRITERIA}
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
          <Edit className="w-3.5 h-3.5" />
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
          <button
            onClick={() => { if (confirm(`Delete ${selectedIds.size} gallery(s)?`)) bulkDeleteMut.mutate(); }}
            disabled={bulkDeleteMut.isPending}
            className="flex items-center gap-1 px-2 py-0.5 rounded text-xs text-red-400 hover:text-red-300 hover:bg-red-900/20"
          >
            {bulkDeleteMut.isPending ? <Loader2 className="w-3 h-3 animate-spin" /> : <Trash2 className="w-3 h-3" />}
            Delete
          </button>
        </>
      }
    >
      {displayMode === "grid" ? (
        <div className="grid gap-3" style={{ gridTemplateColumns: "repeat(auto-fill, minmax(var(--card-min-width, 200px), 1fr))" }}>
          {items.map((g) => (
            <GalleryCard key={g.id} gallery={g} onClick={() => selecting ? toggle(g.id) : onNavigate({ page: "gallery", id: g.id })} selected={selectedIds.has(g.id)} onSelect={() => toggle(g.id)} selecting={selecting} />
          ))}
        </div>
      ) : (
        <GalleryListTable galleries={items} onNavigate={onNavigate} selectedIds={selectedIds} onToggle={toggle} selecting={selecting} />
      )}
      {items.length === 0 && (
        <div className="text-center text-plex-text-secondary py-16">
          <FolderOpen className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>No galleries found</p>
        </div>
      )}
    </ListPage>
    <BulkEditDialog
      open={showBulkEdit}
      onClose={() => setShowBulkEdit(false)}
      title="Edit Galleries"
      selectedCount={selectedIds.size}
      fields={GALLERY_BULK_FIELDS}
      onApply={(values) => bulkEditMut.mutate(values)}
      isPending={bulkEditMut.isPending}
    />
    </>
  );
}

function GalleryCard({ gallery, onClick, selected, onSelect, selecting }: { gallery: Gallery; onClick: () => void; selected?: boolean; onSelect?: () => void; selecting?: boolean }) {
  return (
    <div className={`bg-plex-card rounded overflow-hidden border hover:border-plex-accent/60 transition-all cursor-pointer group relative ${selected ? "border-plex-accent ring-2 ring-plex-accent" : "border-plex-border"}`} onClick={onClick}>
      <div className="aspect-video bg-plex-surface flex items-center justify-center relative overflow-hidden">
        <div className={`absolute top-1 left-1 z-10 ${selected || selecting ? "opacity-100" : "opacity-0 group-hover:opacity-100"} transition-opacity`}>
          <input type="checkbox" checked={selected} onChange={(e) => { e.stopPropagation(); onSelect?.(); }} onClick={(e) => e.stopPropagation()} className="w-4 h-4 rounded border-plex-border cursor-pointer accent-plex-accent" />
        </div>
        {gallery.coverPath ? (
          <img src={gallery.coverPath} alt={gallery.title || ""} className="w-full h-full object-cover" loading="lazy" />
        ) : (
          <FolderOpen className="w-10 h-10 text-plex-text-muted opacity-30" />
        )}
        <RatingBanner rating={gallery.rating} />
        {gallery.studioName && (
          <div className="absolute top-1 right-1 text-xs bg-black/70 px-1.5 py-0.5 rounded text-white truncate max-w-[80%]">
            {gallery.studioName}
          </div>
        )}
      </div>
      <div className="p-2">
        <h3 className="font-medium text-sm truncate text-plex-text">{gallery.title || "Untitled"}</h3>
        {gallery.date && <div className="text-xs text-plex-text-secondary">{gallery.date}</div>}
      </div>
      <GalleryCardPopovers gallery={gallery} />
    </div>
  );
}

function GalleryCardPopovers({ gallery }: { gallery: Gallery }) {
  const hasAny = gallery.imageCount > 0 || gallery.performers.length > 0 || gallery.tags.length > 0 || gallery.sceneCount > 0 || gallery.organized;
  if (!hasAny) return null;

  return (
    <div className="flex items-center justify-center gap-1 px-2 pb-2 border-t border-plex-border pt-1.5">
      {gallery.imageCount > 0 && (
        <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Images">
          <Image className="w-3 h-3" /> {gallery.imageCount}
        </span>
      )}
      {gallery.tags.length > 0 && (
        <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Tags">
          <Tag className="w-3 h-3" /> {gallery.tags.length}
        </span>
      )}
      {gallery.performers.length > 0 && (
        <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Performers">
          <Users className="w-3 h-3" /> {gallery.performers.length}
        </span>
      )}
      {gallery.sceneCount > 0 && (
        <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Scenes">
          <Film className="w-3 h-3" /> {gallery.sceneCount}
        </span>
      )}
      {gallery.organized && (
        <span className="text-plex-text-muted" title="Organized">
          <Box className="w-3 h-3" />
        </span>
      )}
    </div>
  );
}

function GalleryListTable({ galleries: items, onNavigate, selectedIds, onToggle, selecting }: { galleries: Gallery[]; onNavigate: (r: any) => void; selectedIds?: Set<number>; onToggle?: (id: number) => void; selecting?: boolean }) {
  return (
    <table className="w-full text-sm">
      <thead>
        <tr className="border-b border-plex-border text-left text-plex-text-muted text-xs">
          {selectedIds && <th className="w-8 py-2 px-3"></th>}
          <th className="py-2 px-3">Title</th>
          <th className="py-2 px-3">Studio</th>
          <th className="py-2 px-3">Date</th>
          <th className="py-2 px-3 text-right">Images</th>
          <th className="py-2 px-3 text-right">Rating</th>
        </tr>
      </thead>
      <tbody>
        {items.map((g) => (
          <tr key={g.id} onClick={() => selecting ? onToggle?.(g.id) : onNavigate({ page: "gallery", id: g.id })} className={`border-b border-plex-border hover:bg-plex-card cursor-pointer ${selectedIds?.has(g.id) ? "bg-plex-accent/10" : ""}`}>
            {selectedIds && <td className="py-2 px-3"><input type="checkbox" checked={selectedIds.has(g.id)} onChange={() => onToggle?.(g.id)} onClick={(e) => e.stopPropagation()} className="w-3.5 h-3.5 rounded border-plex-border cursor-pointer accent-plex-accent" /></td>}
            <td className="py-2 px-3 text-plex-text">{g.title || "Untitled"}</td>
            <td className="py-2 px-3 text-plex-text-secondary">{g.studioName ?? ""}</td>
            <td className="py-2 px-3 text-plex-text-secondary">{g.date ?? ""}</td>
            <td className="py-2 px-3 text-plex-text-secondary text-right">{g.imageCount}</td>
            <td className="py-2 px-3 text-plex-text-secondary text-right">{g.rating ?? ""}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

/* ── Gallery Create Modal ── */
function GalleryCreateModal({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const qc = useQueryClient();
  const [form, setForm] = useState({
    title: "",
    code: "",
    date: "",
    details: "",
    photographer: "",
    rating: undefined as number | undefined,
  });

  const mutation = useMutation({
    mutationFn: (data: GalleryCreate) => galleries.create(data),
    onSuccess: (created) => {
      qc.invalidateQueries({ queryKey: ["galleries"] });
      setForm({ title: "", code: "", date: "", details: "", photographer: "", rating: undefined });
      onClose();
      if (created?.id) onCreated(created.id);
    },
  });

  const save = () => {
    const title = form.title.trim();
    if (!title) return;
    mutation.mutate({
      title,
      code: form.code || undefined,
      date: form.date || undefined,
      details: form.details || undefined,
      photographer: form.photographer || undefined,
      rating: form.rating,
    });
  };

  return (
    <EditModal title="Create Gallery" open={open} onClose={onClose}>
      <Field label="Title">
        <TextInput value={form.title} onChange={(v) => setForm({ ...form, title: v })} />
      </Field>
      <Field label="Code">
        <TextInput value={form.code} onChange={(v) => setForm({ ...form, code: v })} />
      </Field>
      <Field label="Date">
        <TextInput value={form.date} onChange={(v) => setForm({ ...form, date: v })} placeholder="YYYY-MM-DD" />
      </Field>
      <Field label="Photographer">
        <TextInput value={form.photographer} onChange={(v) => setForm({ ...form, photographer: v })} />
      </Field>
      <Field label="Details">
        <TextArea value={form.details} onChange={(v) => setForm({ ...form, details: v })} rows={3} />
      </Field>
      <RatingField value={form.rating} onChange={(value) => setForm({ ...form, rating: value })} />
      <div className="flex justify-end mt-4">
        <SaveButton loading={mutation.isPending} onClick={save} />
      </div>
    </EditModal>
  );
}
