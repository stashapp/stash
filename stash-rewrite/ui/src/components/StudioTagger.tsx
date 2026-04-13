import { useCallback, useState, useRef } from "react";
import { useMutation } from "@tanstack/react-query";
import { studios } from "../api/client";
import type { Studio, StashBoxStudioMatch, StashBoxStudioImportRequest } from "../api/types";
import { useAppConfig } from "../state/AppConfigContext";
import {
  Search, Loader2, Check, X, AlertCircle,
  CloudDownload, Fingerprint, Eye, EyeOff,
} from "lucide-react";

interface StudioTaggerProps {
  studios: Studio[];
}

interface TaggerConfig {
  selectedEndpoint: string;
  showTagged: boolean;
}

interface StudioSearchState {
  loading: boolean;
  results?: StashBoxStudioMatch[];
  error?: string;
  selectedIndex?: number;
  saved?: boolean;
}

const CONCURRENCY_LIMIT = 5;
async function runWithConcurrency<T>(items: T[], fn: (item: T) => Promise<void>, limit: number, signal?: AbortSignal): Promise<void> {
  let index = 0;
  const workers = Array.from({ length: Math.min(limit, items.length) }, async () => {
    while (index < items.length) {
      if (signal?.aborted) return;
      const i = index++;
      await fn(items[i]);
    }
  });
  await Promise.all(workers);
}

export function StudioTagger({ studios: studioList }: StudioTaggerProps) {
  const { config } = useAppConfig();
  const stashBoxes = config?.scraping?.stashBoxes ?? [];

  const [taggerConfig, setTaggerConfig] = useState<TaggerConfig>({
    selectedEndpoint: stashBoxes[0]?.endpoint ?? "",
    showTagged: false,
  });

  const [searchStates, setSearchStates] = useState<Record<number, StudioSearchState>>({});
  const [queryOverrides, setQueryOverrides] = useState<Record<number, string>>({});

  const updateSearchState = useCallback(
    (studioId: number, update: Partial<StudioSearchState>) => {
      setSearchStates((prev) => ({ ...prev, [studioId]: { ...prev[studioId], ...update } }));
    },
    []
  );

  const searchStudio = useCallback(
    async (studio: Studio) => {
      const query = queryOverrides[studio.id] ?? studio.name;
      updateSearchState(studio.id, { loading: true, error: undefined, results: undefined, saved: false });
      try {
        const endpoint = taggerConfig.selectedEndpoint || undefined;
        const results = await studios.searchStashBox(studio.id, query, endpoint);
        updateSearchState(studio.id, {
          loading: false,
          results,
          selectedIndex: results.length > 0 ? 0 : undefined,
        });
      } catch (err) {
        updateSearchState(studio.id, {
          loading: false,
          error: err instanceof Error ? err.message : "Search failed",
        });
      }
    },
    [queryOverrides, taggerConfig.selectedEndpoint, updateSearchState]
  );

  const [batchSearching, setBatchSearching] = useState(false);
  const abortRef = useRef<AbortController | null>(null);
  const searchAll = useCallback(async () => {
    setBatchSearching(true);
    const controller = new AbortController();
    abortRef.current = controller;
    const toSearch = studioList.filter((s) => !searchStates[s.id]?.saved);
    await runWithConcurrency(toSearch, (s) => searchStudio(s), CONCURRENCY_LIMIT, controller.signal);
    setBatchSearching(false);
    abortRef.current = null;
  }, [studioList, searchStates, searchStudio]);

  const cancelBatchSearch = useCallback(() => {
    abortRef.current?.abort();
    setBatchSearching(false);
  }, []);

  if (stashBoxes.length === 0) {
    return (
      <div className="px-4 py-12 text-center">
        <AlertCircle className="w-12 h-12 mx-auto mb-3 text-plex-text-muted opacity-50" />
        <p className="text-plex-text-secondary text-lg">No Stash-Box Sources Configured</p>
        <p className="text-plex-text-muted text-sm mt-1">
          Add a Stash-Box endpoint in Settings &gt; Metadata Providers to use the tagger.
        </p>
      </div>
    );
  }

  const visibleStudios = taggerConfig.showTagged
    ? studioList
    : studioList.filter((s) => !s.stashIds || s.stashIds.length === 0);

  return (
    <div className="space-y-0">
      {/* Toolbar */}
      <div className="flex flex-wrap items-center gap-2 bg-plex-surface border-b border-plex-border px-4 py-2">
        <div className="flex items-center gap-2">
          <label className="text-xs text-plex-text-muted whitespace-nowrap">Source:</label>
          <select
            value={taggerConfig.selectedEndpoint}
            onChange={(e) => setTaggerConfig((c) => ({ ...c, selectedEndpoint: e.target.value }))}
            className="bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent"
          >
            {stashBoxes.map((sb) => (
              <option key={sb.endpoint} value={sb.endpoint}>
                {sb.name || sb.endpoint}
              </option>
            ))}
          </select>
        </div>

        <button
          onClick={() => setTaggerConfig((c) => ({ ...c, showTagged: !c.showTagged }))}
          className="flex items-center gap-1 px-2 py-1 rounded text-xs border border-plex-border bg-plex-input text-plex-text-secondary hover:text-plex-text"
        >
          {taggerConfig.showTagged ? <Eye className="w-3.5 h-3.5" /> : <EyeOff className="w-3.5 h-3.5" />}
          {taggerConfig.showTagged ? "Hide Already Tagged" : "Show Already Tagged"}
        </button>

        {batchSearching ? (
          <button
            onClick={cancelBatchSearch}
            className="flex items-center gap-1.5 px-3 py-1 rounded text-xs font-medium bg-red-600 text-white hover:bg-red-500"
          >
            <X className="w-3.5 h-3.5" />
            Cancel
          </button>
        ) : (
          <button
            onClick={searchAll}
            className="flex items-center gap-1.5 px-3 py-1 rounded text-xs font-medium bg-plex-accent text-white hover:bg-plex-accent-hover"
          >
            <CloudDownload className="w-3.5 h-3.5" />
            Scrape All
          </button>
        )}

        <span className="text-xs text-plex-text-muted ml-auto">
          {visibleStudios.length} studio{visibleStudios.length !== 1 ? "s" : ""}
        </span>
      </div>

      {/* Studio list */}
      <div className="divide-y divide-plex-border">
        {visibleStudios.map((studio) => (
          <StudioTaggerRow
            key={studio.id}
            studio={studio}
            state={searchStates[studio.id]}
            query={queryOverrides[studio.id] ?? studio.name}
            onQueryChange={(q) => setQueryOverrides((prev) => ({ ...prev, [studio.id]: q }))}
            onSearch={() => searchStudio(studio)}
            onUpdateState={(update) => updateSearchState(studio.id, update)}
            endpoint={taggerConfig.selectedEndpoint}
          />
        ))}
      </div>
    </div>
  );
}

function StudioTaggerRow({
  studio,
  state,
  query,
  onQueryChange,
  onSearch,
  onUpdateState,
  endpoint,
}: {
  studio: Studio;
  state?: StudioSearchState;
  query: string;
  onQueryChange: (q: string) => void;
  onSearch: () => void;
  onUpdateState: (update: Partial<StudioSearchState>) => void;
  endpoint: string;
}) {
  const imageUrl = studio.imagePath;

  const importMut = useMutation({
    mutationFn: () => {
      const selectedResult = state?.results?.[state.selectedIndex ?? 0];
      if (!selectedResult) throw new Error("No result selected");
      const importReq: StashBoxStudioImportRequest = {
        endpoint,
        studioId: selectedResult.id,
      };
      return studios.importFromStashBox(studio.id, importReq);
    },
    onSuccess: () => {
      onUpdateState({ saved: true });
    },
  });

  return (
    <div className={`px-4 py-3 ${state?.saved ? "opacity-50" : ""}`}>
      <div className="flex gap-4">
        {/* Studio image */}
        <div className="flex-shrink-0 w-24">
          <div className="relative aspect-video bg-plex-card rounded overflow-hidden">
            {imageUrl ? (
              <img src={imageUrl} alt="" className="w-full h-full object-contain" loading="lazy" />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-plex-text-muted text-xs">No Image</div>
            )}
          </div>
          <p className="text-xs text-plex-text mt-1 truncate font-medium">{studio.name}</p>
          {studio.stashIds && studio.stashIds.length > 0 && (
            <div className="flex flex-wrap gap-1 mt-1">
              {studio.stashIds.map((sid) => (
                <span key={`${sid.endpoint}-${sid.stashId}`} className="text-[9px] px-1.5 py-0.5 rounded bg-green-600/20 text-green-300" title={sid.endpoint}>
                  <Fingerprint className="w-2.5 h-2.5 inline mr-0.5" />
                  {sid.stashId.substring(0, 8)}…
                </span>
              ))}
            </div>
          )}
        </div>

        {/* Search + Results */}
        <div className="flex-1 min-w-0">
          <div className="flex gap-2 mb-2">
            <input
              type="text"
              value={query}
              onChange={(e) => onQueryChange(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && onSearch()}
              placeholder="Search query..."
              className="flex-1 bg-plex-input border border-plex-border rounded px-3 py-1.5 text-xs text-plex-text focus:outline-none focus:border-plex-accent"
            />
            <button
              onClick={onSearch}
              disabled={state?.loading}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded text-xs font-medium bg-plex-accent text-white hover:bg-plex-accent-hover disabled:opacity-60"
            >
              {state?.loading ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Search className="w-3.5 h-3.5" />}
              Search
            </button>
          </div>

          {state?.error && (
            <p className="text-xs text-red-400 mb-2">
              <AlertCircle className="w-3 h-3 inline mr-1" />{state.error}
            </p>
          )}

          {state?.results && state.results.length === 0 && (
            <p className="text-xs text-plex-text-muted">No matches found.</p>
          )}

          {state?.results && state.results.length > 0 && (
            <div className="space-y-1">
              {state.results.map((result, i) => (
                <StudioResultRow
                  key={`${result.endpoint}-${result.id}`}
                  result={result}
                  isSelected={i === (state.selectedIndex ?? 0)}
                  onClick={() => onUpdateState({ selectedIndex: i })}
                  onSave={i === (state.selectedIndex ?? 0) ? () => importMut.mutate() : undefined}
                  saving={i === (state.selectedIndex ?? 0) ? importMut.isPending : false}
                  saved={state.saved}
                />
              ))}
            </div>
          )}

          {state?.saved && (
            <div className="flex items-center gap-1 mt-2 text-xs text-green-400">
              <Check className="w-3.5 h-3.5" />Saved successfully
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function StudioResultRow({
  result,
  isSelected,
  onClick,
  onSave,
  saving,
  saved,
}: {
  result: StashBoxStudioMatch;
  isSelected: boolean;
  onClick: () => void;
  onSave?: () => void;
  saving?: boolean;
  saved?: boolean;
}) {
  return (
    <div
      onClick={onClick}
      className={`rounded border cursor-pointer transition-colors ${
        isSelected ? "border-plex-accent bg-plex-card" : "border-plex-border bg-plex-surface hover:border-plex-accent/50"
      }`}
    >
      <div className="flex items-center gap-3 p-2">
        {result.imageUrl && (
          <img src={result.imageUrl} alt="" className="h-8 w-16 object-contain rounded flex-shrink-0" loading="lazy" />
        )}
        <div className="flex-1 min-w-0">
          <p className="text-xs font-medium text-plex-text truncate">{result.name}</p>
          <div className="flex items-center gap-2 text-[10px] text-plex-text-muted">
            {result.parentName && <span>Parent: {result.parentName}</span>}
            {result.aliases && result.aliases.length > 0 && <span>{result.aliases.length} alias(es)</span>}
          </div>
        </div>
      </div>

      {isSelected && (
        <div className="border-t border-plex-border p-3">
          <div className="space-y-1 text-xs mb-3">
            {result.parentName && <FieldRow label="Parent" value={result.parentName} />}
            {result.aliases && result.aliases.length > 0 && <FieldRow label="Aliases" value={result.aliases.join(", ")} />}
            {result.urls && result.urls.length > 0 && <FieldRow label="URLs" value={result.urls.join(", ")} />}
          </div>

          {onSave && !saved && (
            <div className="flex justify-end">
              <button
                onClick={(e) => { e.stopPropagation(); onSave(); }}
                disabled={saving}
                className="flex items-center gap-1.5 px-4 py-1.5 rounded text-xs font-medium bg-green-600 text-white hover:bg-green-500 disabled:opacity-60"
              >
                {saving ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Check className="w-3.5 h-3.5" />}
                Save
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

function FieldRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex gap-2">
      <span className="text-plex-text-muted w-16 flex-shrink-0 text-right">{label}:</span>
      <span className="text-plex-text truncate">{value}</span>
    </div>
  );
}
