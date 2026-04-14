import { QueryKey, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { GitMerge, Loader2, Search, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

export interface DetailMergeCandidate {
  id: number;
  name: string;
  imagePath?: string;
  subtitle?: string;
}

interface Props {
  open: boolean;
  onClose: () => void;
  entityType: "scene" | "performer" | "tag" | "studio";
  targetItem: DetailMergeCandidate;
  searchItems: (term: string) => Promise<DetailMergeCandidate[]>;
  onMerge: (targetId: number, sourceIds: number[]) => Promise<unknown>;
  invalidateQueryKeys: Array<string | QueryKey>;
  onMerged?: () => void;
}

export function DetailMergeDialog({
  open,
  onClose,
  entityType,
  targetItem,
  searchItems,
  onMerge,
  invalidateQueryKeys,
  onMerged,
}: Props) {
  const queryClient = useQueryClient();
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedIds, setSelectedIds] = useState<number[]>([]);

  useEffect(() => {
    if (!open) {
      setSearchTerm("");
      setSelectedIds([]);
    }
  }, [open]);

  const { data, isLoading } = useQuery({
    queryKey: ["detail-merge", entityType, targetItem.id, searchTerm],
    queryFn: () => searchItems(searchTerm),
    enabled: open,
  });

  const candidates = useMemo(
    () => (data ?? []).filter((item) => item.id !== targetItem.id),
    [data, targetItem.id],
  );

  const mergeMut = useMutation({
    mutationFn: () => onMerge(targetItem.id, selectedIds),
    onSuccess: async () => {
      for (const key of invalidateQueryKeys) {
        await queryClient.invalidateQueries({ queryKey: typeof key === "string" ? [key] : key });
      }
      onMerged?.();
      onClose();
    },
  });

  if (!open) return null;

  const toggleSelection = (id: number) => {
    setSelectedIds((prev) => prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id]);
  };

  return (
    <div className="fixed inset-0 z-[110] flex items-center justify-center px-4">
      <div className="absolute inset-0 bg-black/70" onClick={onClose} />
      <div className="relative flex w-full max-w-2xl max-h-[85vh] flex-col overflow-hidden rounded-xl border border-plex-border bg-plex-bg shadow-2xl">
        <div className="flex items-center justify-between border-b border-plex-border px-5 py-4">
          <div>
            <h2 className="flex items-center gap-2 text-lg font-semibold text-plex-text">
              <GitMerge className="h-5 w-5" /> Merge {entityType.charAt(0).toUpperCase() + entityType.slice(1)}s
            </h2>
            <p className="mt-1 text-sm text-plex-text-secondary">Keep the current {entityType} and merge other matching entries into it.</p>
          </div>
          <button onClick={onClose} className="rounded p-1 text-plex-text-secondary hover:bg-plex-surface hover:text-plex-text">
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="space-y-4 overflow-y-auto px-5 py-4">
          <div className="rounded-lg border border-green-600/30 bg-green-600/10 p-3">
            <div className="text-xs uppercase tracking-wide text-green-300">Merge target</div>
            <div className="mt-1 text-sm font-medium text-plex-text">{targetItem.name}</div>
            {targetItem.subtitle && <div className="mt-0.5 text-xs text-plex-text-muted">{targetItem.subtitle}</div>}
          </div>

          <div className="relative">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-plex-text-muted" />
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder={`Search ${entityType}s to merge...`}
              className="w-full rounded-lg border border-plex-border bg-plex-input py-2 pl-9 pr-3 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
            />
          </div>

          <div className="space-y-2">
            {isLoading && (
              <div className="flex items-center gap-2 rounded-lg border border-plex-border bg-plex-card px-3 py-4 text-sm text-plex-text-secondary">
                <Loader2 className="h-4 w-4 animate-spin" /> Loading candidates...
              </div>
            )}

            {!isLoading && candidates.length === 0 && (
              <div className="rounded-lg border border-plex-border bg-plex-card px-3 py-4 text-sm text-plex-text-secondary">
                No merge candidates found.
              </div>
            )}

            {!isLoading && candidates.map((candidate) => {
              const selected = selectedIds.includes(candidate.id);
              return (
                <button
                  key={candidate.id}
                  onClick={() => toggleSelection(candidate.id)}
                  className={`flex w-full items-center gap-3 rounded-lg border px-3 py-2 text-left transition-colors ${
                    selected
                      ? "border-plex-accent bg-plex-accent/10"
                      : "border-plex-border bg-plex-card hover:border-plex-accent/50 hover:bg-plex-surface"
                  }`}
                >
                  <div className={`flex h-4 w-4 items-center justify-center rounded border ${selected ? "border-plex-accent bg-plex-accent" : "border-plex-border bg-plex-surface"}`}>
                    {selected && <div className="h-1.5 w-1.5 rounded-full bg-white" />}
                  </div>
                  {candidate.imagePath ? (
                    <img src={candidate.imagePath} alt="" className="h-10 w-10 rounded object-cover" />
                  ) : (
                    <div className="h-10 w-10 rounded bg-plex-surface" />
                  )}
                  <div className="min-w-0 flex-1">
                    <div className="truncate text-sm font-medium text-plex-text">{candidate.name}</div>
                    {candidate.subtitle && <div className="truncate text-xs text-plex-text-muted">{candidate.subtitle}</div>}
                  </div>
                </button>
              );
            })}
          </div>
        </div>

        <div className="flex items-center justify-between border-t border-plex-border px-5 py-4">
          <div className="text-sm text-plex-text-secondary">{selectedIds.length} selected for merge</div>
          <div className="flex items-center gap-2">
            <button onClick={onClose} className="rounded border border-plex-border px-3 py-1.5 text-sm text-plex-text-secondary hover:text-plex-text">
              Cancel
            </button>
            <button
              onClick={() => mergeMut.mutate()}
              disabled={selectedIds.length === 0 || mergeMut.isPending}
              className="inline-flex items-center gap-2 rounded bg-yellow-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-yellow-500 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {mergeMut.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <GitMerge className="h-4 w-4" />}
              Merge into current {entityType}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}