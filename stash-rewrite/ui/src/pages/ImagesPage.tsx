import { useState, useMemo } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { images } from "../api/client";
import type { FindFilter, Image, ImageFilterCriteria } from "../api/types";
import { ListPage, type DisplayMode } from "../components/ListPage";
import { RatingBanner } from "../components/Rating";
import { useMultiSelect } from "../hooks/useMultiSelect";
import { ImageIcon, Users, Tag, Trash2, Loader2, Check, Edit, Box, Heart } from "lucide-react";
import { IMAGE_CRITERIA } from "../components/FilterDialog";
import { BulkEditDialog, IMAGE_BULK_FIELDS } from "../components/BulkEditDialog";
import { Lightbox, type LightboxImage } from "../components/Lightbox";
import { ImageCreateModal } from "./ImageEditModal";

const SORT_OPTIONS = [
  { value: "updated_at", label: "Recently Updated" },
  { value: "created_at", label: "Recently Added" },
  { value: "title", label: "Title" },
  { value: "rating", label: "Rating" },
  { value: "random", label: "Random" },
];

interface Props {
  onNavigate: (r: any) => void;
}

export function ImagesPage({ onNavigate }: Props) {
  const [filter, setFilter] = useState<FindFilter>({ page: 1, perPage: 40, direction: "desc" });
  const [displayMode, setDisplayMode] = useState<DisplayMode>("grid");
  const [selectionMode, setSelectionMode] = useState(false);
  const [objectFilter, setObjectFilter] = useState<Record<string, unknown>>({});
  const [showCreate, setShowCreate] = useState(false);
  const [showBulkEdit, setShowBulkEdit] = useState(false);
  const [lightboxOpen, setLightboxOpen] = useState(false);
  const [lightboxIndex, setLightboxIndex] = useState(0);
  const queryClient = useQueryClient();

  const hasObjectFilter = Object.keys(objectFilter).length > 0;
  const { data, isLoading } = useQuery({
    queryKey: ["images", filter, objectFilter],
    queryFn: () =>
      hasObjectFilter
        ? images.findFiltered({ findFilter: filter, objectFilter: objectFilter as ImageFilterCriteria })
        : images.find(filter),
  });

  const items = data?.items ?? [];
  const { selectedIds, toggle, selectAll, selectNone } = useMultiSelect(items);
  const selecting = selectionMode || selectedIds.size > 0;

  const lightboxImages: LightboxImage[] = useMemo(
    () => items.map((img) => ({ id: img.id, src: images.imageUrl(img.id), title: img.title })),
    [items],
  );

  const bulkDeleteMut = useMutation({
    mutationFn: () => images.bulkDelete([...selectedIds]),
    onSuccess: () => { selectNone(); queryClient.invalidateQueries({ queryKey: ["images"] }); },
  });

  const bulkEditMut = useMutation({
    mutationFn: (values: Record<string, unknown>) =>
      images.bulkUpdate({ ids: [...selectedIds], ...values } as any),
    onSuccess: () => {
      setShowBulkEdit(false);
      selectNone();
      queryClient.invalidateQueries({ queryKey: ["images"] });
    },
  });

  return (
    <>
    <ImageCreateModal open={showCreate} onClose={() => setShowCreate(false)} onCreated={(id) => onNavigate({ page: "image", id })} />
    <ListPage
      title="Images"
      filterMode="images"
      filter={filter}
      onFilterChange={setFilter}
      totalCount={data?.totalCount ?? 0}
      isLoading={isLoading}
      sortOptions={SORT_OPTIONS}
      displayMode={displayMode}
      onDisplayModeChange={setDisplayMode}
      availableDisplayModes={["grid", "wall"]}
      onNew={() => setShowCreate(true)}
      criteriaDefinitions={IMAGE_CRITERIA}
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
          <button
            onClick={() => { if (confirm(`Delete ${selectedIds.size} image(s)?`)) bulkDeleteMut.mutate(); }}
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
        <div className="grid gap-2" style={{ gridTemplateColumns: "repeat(auto-fill, minmax(var(--card-min-width, 140px), 1fr))" }}>
          {items.map((img, idx) => (
            <ImageCard
              key={img.id}
              image={img}
              onPreview={() => {
                if (selecting) { toggle(img.id); return; }
                setLightboxIndex(idx);
                setLightboxOpen(true);
              }}
              onDetails={() => {
                if (selecting) { toggle(img.id); return; }
                onNavigate({ page: "image", id: img.id });
              }}
              selected={selectedIds.has(img.id)}
              onSelect={() => toggle(img.id)}
              selecting={selecting}
            />
          ))}
        </div>
      ) : (
        <div className="columns-2 sm:columns-3 md:columns-4 lg:columns-5 xl:columns-6 gap-2 space-y-2">
          {items.map((img) => (
            <ImageWallCard key={img.id} image={img} onClick={() => onNavigate({ page: "image", id: img.id })} />
          ))}
        </div>
      )}
      {items.length === 0 && (
        <div className="text-center text-plex-text-secondary py-16">
          <ImageIcon className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>No images found</p>
        </div>
      )}
    </ListPage>
    <BulkEditDialog
      open={showBulkEdit}
      onClose={() => setShowBulkEdit(false)}
      title="Edit Images"
      selectedCount={selectedIds.size}
      fields={IMAGE_BULK_FIELDS}
      onApply={(values) => bulkEditMut.mutate(values)}
      isPending={bulkEditMut.isPending}
    />
    <Lightbox
      images={lightboxImages}
      initialIndex={lightboxIndex}
      open={lightboxOpen}
      onClose={() => setLightboxOpen(false)}
    />
    </>
  );
}

function ImageCard({ image, onPreview, onDetails, selected, onSelect, selecting }: { image: Image; onPreview: () => void; onDetails: () => void; selected?: boolean; onSelect?: () => void; selecting?: boolean }) {
  const thumbnailUrl = images.thumbnailUrl(image.id);

  return (
    <div
      className={`bg-plex-card rounded overflow-hidden cursor-pointer border hover:border-plex-accent/60 transition-all group relative ${selected ? "border-plex-accent ring-2 ring-plex-accent" : "border-plex-border"}`}
    >
      <div className="aspect-square bg-plex-surface relative overflow-hidden" onClick={onPreview}>
        <div className={`absolute top-1 left-1 z-10 ${selected || selecting ? "opacity-100" : "opacity-0 group-hover:opacity-100"} transition-opacity`}>
          <input type="checkbox" checked={selected} onChange={(e) => { e.stopPropagation(); onSelect?.(); }} onClick={(e) => e.stopPropagation()} className="w-4 h-4 rounded border-plex-border cursor-pointer accent-plex-accent" />
        </div>
        <img
          src={thumbnailUrl}
          alt={image.title || "Image"}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          loading="lazy"
          onError={(e) => { (e.target as HTMLImageElement).style.display = 'none'; }}
        />
        <RatingBanner rating={image.rating} />
        {image.studioName && (
          <div className="absolute top-1 right-1 text-[10px] bg-black/70 px-1 py-0.5 rounded text-white truncate max-w-[80%]">
            {image.studioName}
          </div>
        )}
      </div>
      <div className="p-1.5" onClick={onDetails}>
        <p className="text-xs font-medium text-plex-text truncate group-hover:text-plex-accent">
          {image.title || "Untitled"}
        </p>
      </div>
      {(image.performers.length > 0 || image.tags.length > 0 || image.oCounter > 0 || image.organized) && (
        <div className="flex items-center justify-center gap-1 px-1.5 pb-1.5 border-t border-plex-border pt-1">
          {image.performers.length > 0 && (
            <span className="flex items-center gap-0.5 text-[10px] text-plex-text-muted">
              <Users className="w-2.5 h-2.5" /> {image.performers.length}
            </span>
          )}
          {image.tags.length > 0 && (
            <span className="flex items-center gap-0.5 text-[10px] text-plex-text-muted">
              <Tag className="w-2.5 h-2.5" /> {image.tags.length}
            </span>
          )}
          {image.oCounter > 0 && (
            <span className="flex items-center gap-0.5 text-[10px] text-plex-text-muted" title="O-counter">
              <Heart className="w-2.5 h-2.5" /> {image.oCounter}
            </span>
          )}
          {image.organized && (
            <span className="text-plex-text-muted" title="Organized">
              <Box className="w-2.5 h-2.5" />
            </span>
          )}
        </div>
      )}
    </div>
  );
}

function ImageWallCard({ image, onClick }: { image: Image; onClick: () => void }) {
  return (
    <div
      onClick={onClick}
      className="break-inside-avoid cursor-pointer rounded overflow-hidden border border-plex-border hover:border-plex-accent/60 transition-all"
    >
      <img
        src={images.thumbnailUrl(image.id)}
        alt={image.title || "Image"}
        className="w-full object-cover"
        loading="lazy"
      />
    </div>
  );
}
