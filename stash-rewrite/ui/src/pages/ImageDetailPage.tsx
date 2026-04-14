import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { images } from "../api/client";
import { formatDate, TagBadge, CustomFieldsDisplay } from "../components/shared";
import { ArrowLeft, Pencil, Trash2, Link as LinkIcon, Heart, Check, Minus, Plus, RotateCcw } from "lucide-react";
import { useEffect, useState } from "react";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { ImageEditModal } from "./ImageEditModal";
import { ExtensionSlot } from "../router/RouteRegistry";
import { InteractiveRating } from "../components/Rating";

interface Props {
  id: number;
  onNavigate: (r: any) => void;
}

export function ImageDetailPage({ id, onNavigate }: Props) {
  const { data: image, isLoading } = useQuery({
    queryKey: ["image", id],
    queryFn: () => images.get(id),
  });
  const [editing, setEditing] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(false);
  const queryClient = useQueryClient();
  const deleteMut = useMutation({
    mutationFn: () => images.delete(id),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ["images"] }); onNavigate({ page: "images" }); },
  });
  const updateMut = useMutation({
    mutationFn: (data: { organized?: boolean; rating?: number }) => images.update(id, data),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["image", id] }),
  });
  const incrementOMut = useMutation({
    mutationFn: () => images.incrementO(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["image", id] }),
  });
  const decrementOMut = useMutation({
    mutationFn: () => images.decrementO(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["image", id] }),
  });

  useEffect(() => {
    if (image) document.title = `${image.title || `Image ${id}`} | Stash`;
    return () => { document.title = "Stash"; };
  }, [image, id]);

  // Keyboard shortcuts
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const el = (e.target as HTMLElement).tagName;
      if (el === "INPUT" || el === "TEXTAREA" || el === "SELECT") return;
      switch (e.key) {
        case "e": setEditing((v) => !v); break;
        case "o": incrementOMut.mutate(); break;
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-plex-accent" />
      </div>
    );
  }

  if (!image) return <div className="text-center text-plex-text-secondary py-16">Image not found</div>;

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <button
          onClick={() => onNavigate({ page: "images" })}
          className="flex items-center gap-1 text-plex-text-secondary hover:text-plex-text text-sm"
        >
          <ArrowLeft className="w-4 h-4" /> Back to images
        </button>
        <div className="flex items-center gap-2">
          <ExtensionSlot slot="image-detail-actions" context={{ image, onNavigate }} />
          <button
            onClick={() => setEditing(true)}
            className="flex items-center gap-1.5 px-3 py-1.5 text-sm bg-plex-card hover:bg-plex-card-hover border border-plex-border rounded"
          >
            <Pencil className="w-3.5 h-3.5" /> Edit
          </button>
          <button
            onClick={() => setConfirmDelete(true)}
            className="flex items-center gap-1.5 px-3 py-1.5 text-sm bg-plex-card border border-plex-border text-plex-text-secondary hover:text-red-300 hover:border-red-500 rounded"
          >
            <Trash2 className="w-3.5 h-3.5" /> Delete
          </button>
        </div>
      </div>
      {image && <ImageEditModal image={image} open={editing} onClose={() => setEditing(false)} />}
      <ConfirmDialog open={confirmDelete} title="Delete Image" message={`Delete "${image.title || 'Untitled'}"? This cannot be undone.`} onConfirm={() => deleteMut.mutate()} onCancel={() => setConfirmDelete(false)} />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Image display */}
        <div className="lg:col-span-2">
          <div className="bg-plex-card border border-plex-border rounded overflow-hidden flex items-center justify-center min-h-[400px] shadow-xl shadow-black/35">
            <img
              src={images.imageUrl(id)}
              alt={image.title || "Image"}
              className="max-w-full max-h-[80vh] object-contain"
              onError={(e) => {
                (e.target as HTMLImageElement).style.display = "none";
                (e.target as HTMLImageElement).parentElement!.innerHTML = '<div class="flex items-center justify-center h-64"><span class="ml-2 text-plex-text-muted">Image unavailable</span></div>';
              }}
            />
          </div>

          <h1 className="text-2xl font-bold mt-4 mb-2 text-plex-text">{image.title || "Untitled"}</h1>

          <div className="flex flex-wrap items-center gap-3 text-sm text-plex-text-secondary mb-3">
            {image.date && <span>{formatDate(image.date)}</span>}
            {image.studioName && image.studioId && <button onClick={() => onNavigate({ page: "studio", id: image.studioId })} className="text-plex-accent hover:underline">{image.studioName}</button>}
            {image.photographer && <span>Photo: {image.photographer}</span>}
          </div>

          {/* Interactive toolbar */}
          <div className="flex items-center gap-4 mb-4 pb-3 border-b border-plex-border">
            <InteractiveRating value={image.rating} onChange={(value) => updateMut.mutate({ rating: value })} />
            <div className="flex items-center gap-2 ml-auto">
              {/* O-counter */}
              <div className="flex items-center gap-1 text-sm text-plex-text-secondary">
                <Heart className={`w-4 h-4 ${image.oCounter > 0 ? "fill-plex-accent text-plex-accent" : ""}`} />
                <span>{image.oCounter}</span>
                <button onClick={() => incrementOMut.mutate()} className="p-0.5 hover:text-plex-accent" title="Increment O"><Plus className="w-3 h-3" /></button>
                <button onClick={() => decrementOMut.mutate()} className="p-0.5 hover:text-plex-accent" title="Decrement O" disabled={image.oCounter === 0}><Minus className="w-3 h-3" /></button>
              </div>
              {/* Organized toggle */}
              <button 
                onClick={() => updateMut.mutate({ organized: !image.organized })}
                className={`p-1 rounded ${image.organized ? "bg-green-600 text-white" : "bg-plex-card text-plex-text-muted"}`}
                title={image.organized ? "Organized" : "Not organized"}
              >
                <Check className="w-4 h-4" />
              </button>
            </div>
          </div>

          {image.details && (
            <p className="text-plex-text mb-4 whitespace-pre-wrap">{image.details}</p>
          )}

          {image.tags.length > 0 && (
            <div className="flex flex-wrap gap-1.5 mb-4">
              {image.tags.map((tag) => (
                <TagBadge key={tag.id} name={tag.name} onClick={() => onNavigate({ page: "tag", id: tag.id })} />
              ))}
            </div>
          )}

          {image.performers.length > 0 && (
            <div className="mb-4">
              <h3 className="text-sm font-semibold text-plex-text-secondary mb-2">Performers</h3>
              <div className="flex flex-wrap gap-2">
                {image.performers.map((p) => (
                  <button
                    key={p.id}
                    onClick={() => onNavigate({ page: "performer", id: p.id })}
                    className="flex items-center gap-2 bg-plex-card border border-plex-border rounded px-3 py-2 hover:border-plex-accent/60 transition-colors"
                  >
                    <div className="w-8 h-8 rounded-full bg-plex-surface flex items-center justify-center text-xs text-plex-text">
                      {p.name[0]}
                    </div>
                    <span className="text-sm">{p.name}</span>
                  </button>
                ))}
              </div>
            </div>
          )}

          <ExtensionSlot slot="image-detail-main-bottom" context={{ image, onNavigate }} />
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          <div className="bg-plex-card rounded p-4 border border-plex-border">
            <h3 className="text-sm font-semibold text-plex-text-secondary mb-3">Details</h3>
            <dl className="space-y-2 text-sm">
              <div><dt className="text-plex-text-muted">O-Counter</dt><dd className="text-plex-text">{image.oCounter}</dd></div>
              <div><dt className="text-plex-text-muted">Organized</dt><dd className="text-plex-text">{image.organized ? "Yes" : "No"}</dd></div>
              <div><dt className="text-plex-text-muted">Created</dt><dd className="text-plex-text">{formatDate(image.createdAt)}</dd></div>
              <div><dt className="text-plex-text-muted">Updated</dt><dd className="text-plex-text">{formatDate(image.updatedAt)}</dd></div>
            </dl>
          </div>

          {image.urls.length > 0 && (
            <div className="bg-plex-card rounded p-4 border border-plex-border">
              <h3 className="text-sm font-semibold text-plex-text-secondary mb-3 flex items-center gap-1.5"><LinkIcon className="w-4 h-4" /> URLs</h3>
              <div className="space-y-1">
                {image.urls.map((url, i) => (
                  <a key={i} href={url} target="_blank" rel="noopener noreferrer"
                    className="text-plex-accent hover:underline text-sm block truncate">{url}</a>
                ))}
              </div>
            </div>
          )}

          <CustomFieldsDisplay customFields={image.customFields} />
          <ExtensionSlot slot="image-detail-sidebar-bottom" context={{ image, onNavigate }} />
        </div>
      </div>
    </div>
  );
}
