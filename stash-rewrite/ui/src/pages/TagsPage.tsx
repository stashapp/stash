import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { tags } from "../api/client";
import type { FindFilter, Tag, TagCreate, TagFilterCriteria } from "../api/types";
import { ListPage, type DisplayMode } from "../components/ListPage";
import { useMultiSelect } from "../hooks/useMultiSelect";
import { EditModal, Field, TextInput, TextArea, SaveButton } from "../components/EditModal";
import { Tag as TagIcon, Film, MapPin, Trash2, Loader2, Check, Edit, Merge, Heart, Image, LayoutGrid, Layers, Users, Building2 } from "lucide-react";
import { MergeDialog } from "../components/MergeDialog";
import { TAG_CRITERIA } from "../components/FilterDialog";
import { BulkEditDialog } from "../components/BulkEditDialog";

const SORT_OPTIONS = [
  { value: "name", label: "Name" },
  { value: "scene_count", label: "Scene Count" },
  { value: "random", label: "Random" },
  { value: "created_at", label: "Recently Added" },
];

interface Props {
  onNavigate: (r: any) => void;
}

export function TagsPage({ onNavigate }: Props) {
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
    queryKey: ["tags", filter, objectFilter],
    queryFn: () =>
      hasObjectFilter
        ? tags.findFiltered({ findFilter: filter, objectFilter: objectFilter as TagFilterCriteria })
        : tags.find(filter),
  });

  const items = data?.items ?? [];
  const { selectedIds, toggle, selectAll, selectNone } = useMultiSelect(items);
  const selecting = selectionMode || selectedIds.size > 0;

  const bulkDeleteMut = useMutation({
    mutationFn: () => tags.bulkDelete([...selectedIds]),
    onSuccess: () => { selectNone(); queryClient.invalidateQueries({ queryKey: ["tags"] }); },
  });

  const bulkEditMut = useMutation({
    mutationFn: (values: Record<string, unknown>) =>
      tags.bulkUpdate({ ids: [...selectedIds], ...values } as any),
    onSuccess: () => {
      setShowBulkEdit(false);
      selectNone();
      queryClient.invalidateQueries({ queryKey: ["tags"] });
    },
  });

  return (
    <>
      <TagCreateModal open={showCreate} onClose={() => setShowCreate(false)} onCreated={(id) => onNavigate({ page: "tag", id })} />
      <ListPage
      title="Tags"
      filterMode="tags"
      filter={filter}
      onFilterChange={setFilter}
      totalCount={data?.totalCount ?? 0}
      isLoading={isLoading}
      sortOptions={SORT_OPTIONS}
      displayMode={displayMode}
      onDisplayModeChange={setDisplayMode}
      availableDisplayModes={["grid", "list"]}
      criteriaDefinitions={TAG_CRITERIA}
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
            onClick={() => { if (confirm(`Delete ${selectedIds.size} tag(s)?`)) bulkDeleteMut.mutate(); }}
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
          {items.map((tag) => (
            <TagCard
              key={tag.id}
              tag={tag}
              onClick={() => selecting ? toggle(tag.id) : onNavigate({ page: "tag", id: tag.id })}
              selected={selectedIds.has(tag.id)}
              onSelect={() => toggle(tag.id)}
              selecting={selecting}
            />
          ))}
        </div>
      ) : (
        <TagListTable tags={items} onNavigate={onNavigate} selectedIds={selectedIds} onToggle={toggle} selecting={selecting} />
      )}
      {items.length === 0 && (
        <div className="text-center text-plex-text-secondary py-16">
          <TagIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>No tags found</p>
        </div>
      )}
      </ListPage>
      <BulkEditDialog
        open={showBulkEdit}
        onClose={() => setShowBulkEdit(false)}
        title="Edit Tags"
        selectedCount={selectedIds.size}
        fields={[{ key: "favorite", label: "Favorite", type: "bool" }]}
        onApply={(values) => bulkEditMut.mutate(values)}
        isPending={bulkEditMut.isPending}
      />
      <MergeDialog
        open={showMerge}
        onClose={() => { setShowMerge(false); selectNone(); }}
        entityType="tag"
        items={items.filter((t) => selectedIds.has(t.id)).map((t) => ({ id: t.id, name: t.name, imagePath: t.imagePath }))}
        onMerge={tags.merge}
        queryKey="tags"
      />
    </>
  );
}

/* ── Tag Create Modal ── */
function TagCreateModal({ open, onClose, onCreated }: { open: boolean; onClose: () => void; onCreated: (id: number) => void }) {
  const qc = useQueryClient();
  const [form, setForm] = useState({ name: "", description: "", aliases: "" });
  const mutation = useMutation({
    mutationFn: (data: TagCreate) => tags.create(data),
    onSuccess: (created) => {
      qc.invalidateQueries({ queryKey: ["tags"] });
      setForm({ name: "", description: "", aliases: "" });
      onClose();
      if (created?.id) onCreated(created.id);
    },
  });
  return (
    <EditModal title="Create Tag" open={open} onClose={onClose}>
      <Field label="Name">
        <TextInput value={form.name} onChange={(v) => setForm({ ...form, name: v })} />
      </Field>
      <Field label="Description">
        <TextArea value={form.description} onChange={(v) => setForm({ ...form, description: v })} rows={3} />
      </Field>
      <Field label="Aliases (one per line)">
        <TextArea value={form.aliases} onChange={(v) => setForm({ ...form, aliases: v })} rows={2} />
      </Field>
      <div className="flex justify-end mt-4">
        <SaveButton loading={mutation.isPending} onClick={() => mutation.mutate({
          name: form.name,
          description: form.description || undefined,
          aliases: form.aliases ? form.aliases.split("\n").map((a) => a.trim()).filter(Boolean) : [],
        })} />
      </div>
    </EditModal>
  );
}

function TagCard({ tag, onClick, selected, onSelect, selecting }: { tag: Tag; onClick: () => void; selected?: boolean; onSelect?: () => void; selecting?: boolean }) {
  const queryClient = useQueryClient();
  const favMut = useMutation({
    mutationFn: () => tags.update(tag.id, { favorite: !tag.favorite }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["tags"] }),
  });

  return (
    <div className={`bg-plex-card rounded overflow-hidden border hover:border-plex-accent/60 transition-all cursor-pointer group ${selected ? "border-plex-accent ring-2 ring-plex-accent" : "border-plex-border"}`} onClick={onClick}>
      <div className="aspect-video bg-plex-surface flex items-center justify-center relative overflow-hidden">
        <div className={`absolute top-1 left-1 z-10 ${selected || selecting ? "opacity-100" : "opacity-0 group-hover:opacity-100"} transition-opacity`}>
          <input type="checkbox" checked={selected} onChange={(e) => { e.stopPropagation(); onSelect?.(); }} onClick={(e) => e.stopPropagation()} className="w-4 h-4 rounded border-plex-border cursor-pointer accent-plex-accent" />
        </div>
        {/* Favorite heart overlay */}
        <button
          onClick={(e) => { e.stopPropagation(); favMut.mutate(); }}
          className={`absolute top-1 right-1 p-1 z-10 transition-opacity ${tag.favorite ? "opacity-100" : "opacity-0 group-hover:opacity-70"}`}
          title={tag.favorite ? "Unfavorite" : "Favorite"}
        >
          <Heart className={`w-5 h-5 ${tag.favorite ? "fill-red-500 text-red-500" : "text-white drop-shadow-md"}`} />
        </button>
        {tag.imagePath ? (
          <img src={tag.imagePath} alt={tag.name} className="w-full h-full object-cover" loading="lazy" />
        ) : (
          <TagIcon className="w-10 h-10 text-plex-text-muted opacity-30" />
        )}
      </div>
      <div className="p-2 text-center">
        <h3 className="font-medium text-sm text-plex-text truncate">{tag.name}</h3>
        {tag.description && (
          <p className="text-xs text-plex-text-secondary mt-0.5 line-clamp-1">{tag.description}</p>
        )}
      </div>
      {(tag.sceneCount || tag.sceneMarkerCount) && (
        <div className="flex items-center justify-center gap-2 px-2 pb-2 border-t border-plex-border pt-1.5">
          {tag.sceneCount != null && tag.sceneCount > 0 && (
            <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Scenes">
              <Film className="w-3 h-3" /> {tag.sceneCount}
            </span>
          )}
          {tag.sceneMarkerCount != null && tag.sceneMarkerCount > 0 && (
            <span className="flex items-center gap-0.5 text-xs text-plex-text-muted" title="Markers">
              <MapPin className="w-3 h-3" /> {tag.sceneMarkerCount}
            </span>
          )}
        </div>
      )}
    </div>
  );
}

function TagListTable({ tags: items, onNavigate, selectedIds, onToggle, selecting }: { tags: Tag[]; onNavigate: (r: any) => void; selectedIds?: Set<number>; onToggle?: (id: number) => void; selecting?: boolean }) {
  return (
    <table className="w-full text-sm">
      <thead>
        <tr className="border-b border-plex-border text-left text-plex-text-muted text-xs">
          {selectedIds && <th className="w-8 py-2 px-3"></th>}
          <th className="py-2 px-3">Name</th>
          <th className="py-2 px-3">Description</th>
          <th className="py-2 px-3">Aliases</th>
          <th className="py-2 px-3 text-right">Scenes</th>
        </tr>
      </thead>
      <tbody>
        {items.map((t) => (
          <tr
            key={t.id}
            onClick={() => selecting ? onToggle?.(t.id) : onNavigate({ page: "tag", id: t.id })}
            className={`border-b border-plex-border hover:bg-plex-card cursor-pointer ${selectedIds?.has(t.id) ? "bg-plex-accent/10" : ""}`}
          >
            {selectedIds && <td className="py-2 px-3"><input type="checkbox" checked={selectedIds.has(t.id)} onChange={() => onToggle?.(t.id)} onClick={(e) => e.stopPropagation()} className="w-3.5 h-3.5 rounded border-plex-border cursor-pointer accent-plex-accent" /></td>}
            <td className="py-2 px-3 text-plex-text">{t.name}</td>
            <td className="py-2 px-3 text-plex-text-secondary truncate max-w-xs">{t.description ?? ""}</td>
            <td className="py-2 px-3 text-plex-text-muted truncate max-w-xs">{t.aliases.join(", ")}</td>
            <td className="py-2 px-3 text-plex-text-secondary text-right">{t.sceneCount ?? ""}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
