import { useCallback, useState, useRef } from "react";
import { useMutation } from "@tanstack/react-query";
import { scenes } from "../api/client";
import type { Scene, StashBoxSceneMatch, StashBoxSceneImportRequest } from "../api/types";
import { useAppConfig } from "../state/AppConfigContext";
import { formatDuration, getResolutionLabel } from "./shared";
import {
  Search,
  Loader2,
  Check,
  X,
  ChevronDown,
  ChevronUp,
  AlertCircle,
  CloudDownload,
  Fingerprint,
  Settings2,
  EyeOff,
  Eye,
  Plus,
  Ban,
} from "lucide-react";

interface SceneTaggerProps {
  scenes: Scene[];
}

interface TaggerConfig {
  selectedEndpoint: string;
  showUnmatched: boolean;
  setCoverImage: boolean;
  setTags: boolean;
  tagOperation: "merge" | "overwrite";
  setPerformers: boolean;
  setStudio: boolean;
  onlyExistingTags: boolean;
  onlyExistingPerformers: boolean;
  onlyExistingStudio: boolean;
  markOrganized: boolean;
  preferFingerprints: boolean;
  queryMode: "auto" | "filename" | "dir" | "path" | "metadata";
  blacklist: string[];
  createParentStudios: boolean;
  createParentTags: boolean;
}

interface SceneSearchState {
  loading: boolean;
  results?: StashBoxSceneMatch[];
  error?: string;
  selectedIndex?: number;
  saved?: boolean;
  excludedPerformers?: Set<string>;
  excludedTags?: Set<string>;
  skipStudio?: boolean;
}

const CONCURRENCY_LIMIT = 5;

function cleanQueryString(input: string, blacklist: string[]): string {
  let cleaned = input.replace(/[._-]/g, " ");
  for (const pattern of blacklist) {
    try { cleaned = cleaned.replace(new RegExp(pattern, "gi"), " "); } catch { /* invalid regex */ }
  }
  return cleaned.replace(/\s{2,}/g, " ").trim();
}

async function runWithConcurrency<T>(
  items: T[],
  fn: (item: T) => Promise<void>,
  limit: number,
  signal?: AbortSignal
): Promise<void> {
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

export function SceneTagger({ scenes: sceneList }: SceneTaggerProps) {
  const { config } = useAppConfig();
  const stashBoxes = config?.scraping?.stashBoxes ?? [];

  const [taggerConfig, setTaggerConfig] = useState<TaggerConfig>({
    selectedEndpoint: stashBoxes[0]?.endpoint ?? "",
    showUnmatched: true,
    setCoverImage: true,
    setTags: true,
    tagOperation: "merge",
    setPerformers: true,
    setStudio: true,
    onlyExistingTags: false,
    onlyExistingPerformers: false,
    onlyExistingStudio: false,
    markOrganized: false,
    preferFingerprints: true,
    queryMode: "auto",
    blacklist: ["\\sXXX\\s", "1080p", "720p", "2160p", "4K", "KTR", "RARBG", "\\smp4\\s"],
    createParentStudios: true,
    createParentTags: true,
  });
  const [showConfig, setShowConfig] = useState(false);
  const [searchStates, setSearchStates] = useState<Record<number, SceneSearchState>>({});
  const [queryOverrides, setQueryOverrides] = useState<Record<number, string>>({});

  const updateSearchState = useCallback(
    (sceneId: number, update: Partial<SceneSearchState>) => {
      setSearchStates((prev) => ({
        ...prev,
        [sceneId]: { ...prev[sceneId], ...update },
      }));
    },
    []
  );

  // Derive search query from scene
  const getSearchQuery = useCallback(
    (scene: Scene): string => {
      if (queryOverrides[scene.id] !== undefined) return queryOverrides[scene.id];
      const file = scene.files[0];
      const mode = taggerConfig.queryMode;

      // metadata mode: use stored metadata only
      if (mode === "metadata") {
        return scene.title || "";
      }

      // filename/dir/path modes: derive from file path
      if (mode === "filename" && file?.basename) {
        return cleanQueryString(file.basename.replace(/\.\w{2,4}$/, ""), taggerConfig.blacklist);
      }
      if (mode === "dir" && file?.path) {
        const parts = file.path.replace(/\\/g, "/").split("/");
        return parts.length > 1 ? cleanQueryString(parts[parts.length - 2], taggerConfig.blacklist) : "";
      }
      if (mode === "path" && file?.path) {
        return cleanQueryString(file.path, taggerConfig.blacklist);
      }

      // auto mode: try title first, then filename
      if (scene.title) return scene.title;
      if (file?.basename) {
        return cleanQueryString(file.basename.replace(/\.\w{2,4}$/, ""), taggerConfig.blacklist);
      }
      return "";
    },
    [queryOverrides, taggerConfig.queryMode, taggerConfig.blacklist]
  );

  const searchScene = useCallback(
    async (scene: Scene) => {
      const query = getSearchQuery(scene);
      updateSearchState(scene.id, { loading: true, error: undefined, results: undefined, saved: false });
      try {
        let results: StashBoxSceneMatch[] = [];
        const endpoint = taggerConfig.selectedEndpoint || undefined;
        const shouldTryFingerprints = taggerConfig.preferFingerprints || !query;

        if (shouldTryFingerprints) {
          results = await scenes.searchStashBox(scene.id, undefined, endpoint);
        }

        if (results.length === 0 && query) {
          results = await scenes.searchStashBox(scene.id, query, endpoint);
        }

        updateSearchState(scene.id, {
          loading: false,
          results,
          selectedIndex: results.length > 0 ? 0 : undefined,
        });
      } catch (err) {
        updateSearchState(scene.id, {
          loading: false,
          error: err instanceof Error ? err.message : "Search failed",
        });
      }
    },
    [getSearchQuery, taggerConfig.preferFingerprints, taggerConfig.selectedEndpoint, updateSearchState]
  );

  // Fingerprint-only search
  const searchSceneFingerprints = useCallback(
    async (scene: Scene) => {
      updateSearchState(scene.id, { loading: true, error: undefined, results: undefined, saved: false });
      try {
        const endpoint = taggerConfig.selectedEndpoint || undefined;
        const results = await scenes.searchStashBox(scene.id, undefined, endpoint);
        updateSearchState(scene.id, {
          loading: false,
          results,
          selectedIndex: results.length > 0 ? 0 : undefined,
        });
      } catch (err) {
        updateSearchState(scene.id, {
          loading: false,
          error: err instanceof Error ? err.message : "Search failed",
        });
      }
    },
    [taggerConfig.selectedEndpoint, updateSearchState]
  );

  // Batch scrape all (concurrent)
  const [batchSearching, setBatchSearching] = useState(false);
  const abortRef = useRef<AbortController | null>(null);
  const searchAll = useCallback(async () => {
    setBatchSearching(true);
    const controller = new AbortController();
    abortRef.current = controller;
    const toSearch = sceneList.filter((s) => !searchStates[s.id]?.saved);
    await runWithConcurrency(toSearch, (scene) => searchScene(scene), CONCURRENCY_LIMIT, controller.signal);
    setBatchSearching(false);
    abortRef.current = null;
  }, [sceneList, searchStates, searchScene]);

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

  const visibleScenes = taggerConfig.showUnmatched
    ? sceneList
    : sceneList.filter((s) => {
        const state = searchStates[s.id];
        return !state || !state.results || state.results.length > 0;
      });

  return (
    <div className="space-y-0">
      {/* Tagger Toolbar */}
      <div className="flex flex-wrap items-center gap-2 bg-plex-surface border-b border-plex-border px-4 py-2">
        {/* Source selector */}
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

        {/* Show/Hide unmatched */}
        <button
          onClick={() => setTaggerConfig((c) => ({ ...c, showUnmatched: !c.showUnmatched }))}
          className="flex items-center gap-1 px-2 py-1 rounded text-xs border border-plex-border bg-plex-input text-plex-text-secondary hover:text-plex-text"
        >
          {taggerConfig.showUnmatched ? <Eye className="w-3.5 h-3.5" /> : <EyeOff className="w-3.5 h-3.5" />}
          {taggerConfig.showUnmatched ? "Hide Unmatched" : "Show Unmatched"}
        </button>

        {/* Scrape All / Cancel */}
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

        {/* Config toggle */}
        <button
          onClick={() => setShowConfig(!showConfig)}
          className="flex items-center gap-1 px-2 py-1 rounded text-xs border border-plex-border bg-plex-input text-plex-text-secondary hover:text-plex-text ml-auto"
        >
          <Settings2 className="w-3.5 h-3.5" />
          Config
        </button>
      </div>

      {/* Config panel */}
      {showConfig && (
        <div className="bg-plex-card border-b border-plex-border px-4 py-3 space-y-4">
          {/* Query mode + Blacklist */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div>
              <label className="block text-xs text-plex-text-muted mb-1">Query Mode</label>
              <select
                value={taggerConfig.queryMode}
                onChange={(e) => setTaggerConfig((c) => ({ ...c, queryMode: e.target.value as TaggerConfig["queryMode"] }))}
                className="w-full bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent"
              >
                <option value="auto">Auto</option>
                <option value="filename">Filename</option>
                <option value="dir">Directory</option>
                <option value="path">Full Path</option>
                <option value="metadata">Metadata</option>
              </select>
            </div>
            <div>
              <label className="block text-xs text-plex-text-muted mb-1">Blacklist (regex, one per line)</label>
              <textarea
                value={taggerConfig.blacklist.join("\n")}
                onChange={(e) => setTaggerConfig((c) => ({ ...c, blacklist: e.target.value.split("\n").filter(Boolean) }))}
                className="w-full bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent h-16 resize-y font-mono"
                placeholder="\\sXXX\\s&#10;1080p&#10;..."
              />
            </div>
          </div>

          {/* Main options */}
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
            <label className="flex items-center gap-2 text-xs text-plex-text">
              <input
                type="checkbox"
                checked={taggerConfig.preferFingerprints}
                onChange={(e) => setTaggerConfig((c) => ({ ...c, preferFingerprints: e.target.checked }))}
                className="rounded border-plex-border"
              />
              Prefer fingerprint match
            </label>
            <label className="flex items-center gap-2 text-xs text-plex-text">
              <input
                type="checkbox"
                checked={taggerConfig.setCoverImage}
                onChange={(e) => setTaggerConfig((c) => ({ ...c, setCoverImage: e.target.checked }))}
                className="rounded border-plex-border"
              />
              Set cover image
            </label>
            <label className="flex items-center gap-2 text-xs text-plex-text">
              <input
                type="checkbox"
                checked={taggerConfig.markOrganized}
                onChange={(e) => setTaggerConfig((c) => ({ ...c, markOrganized: e.target.checked }))}
                className="rounded border-plex-border"
              />
              Mark as organized
            </label>
          </div>

          {/* Tags section */}
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3 pt-2 border-t border-plex-border/50">
            <label className="flex items-center gap-2 text-xs text-plex-text">
              <input
                type="checkbox"
                checked={taggerConfig.setTags}
                onChange={(e) => setTaggerConfig((c) => ({ ...c, setTags: e.target.checked }))}
                className="rounded border-plex-border"
              />
              Set tags
            </label>
            {taggerConfig.setTags && (
              <>
                <div className="flex items-center gap-2 text-xs text-plex-text">
                  <span className="text-plex-text-muted">Tag operation:</span>
                  <select
                    value={taggerConfig.tagOperation}
                    onChange={(e) => setTaggerConfig((c) => ({ ...c, tagOperation: e.target.value as "merge" | "overwrite" }))}
                    className="bg-plex-input border border-plex-border rounded px-2 py-0.5 text-xs text-plex-text focus:outline-none focus:border-plex-accent"
                  >
                    <option value="merge">Merge</option>
                    <option value="overwrite">Overwrite</option>
                  </select>
                </div>
                <label className="flex items-center gap-2 text-xs text-plex-text">
                  <input
                    type="checkbox"
                    checked={taggerConfig.onlyExistingTags}
                    onChange={(e) => setTaggerConfig((c) => ({ ...c, onlyExistingTags: e.target.checked }))}
                    className="rounded border-plex-border"
                  />
                  Only add existing tags
                </label>
                <label className="flex items-center gap-2 text-xs text-plex-text">
                  <input
                    type="checkbox"
                    checked={taggerConfig.createParentTags}
                    onChange={(e) => setTaggerConfig((c) => ({ ...c, createParentTags: e.target.checked }))}
                    className="rounded border-plex-border"
                  />
                  Create parent tags
                </label>
              </>
            )}
          </div>

          {/* Performers section */}
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3 pt-2 border-t border-plex-border/50">
            <label className="flex items-center gap-2 text-xs text-plex-text">
              <input
                type="checkbox"
                checked={taggerConfig.setPerformers}
                onChange={(e) => setTaggerConfig((c) => ({ ...c, setPerformers: e.target.checked }))}
                className="rounded border-plex-border"
              />
              Set performers
            </label>
            {taggerConfig.setPerformers && (
              <label className="flex items-center gap-2 text-xs text-plex-text">
                <input
                  type="checkbox"
                  checked={taggerConfig.onlyExistingPerformers}
                  onChange={(e) => setTaggerConfig((c) => ({ ...c, onlyExistingPerformers: e.target.checked }))}
                  className="rounded border-plex-border"
                />
                Only add existing performers
              </label>
            )}
          </div>

          {/* Studio section */}
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3 pt-2 border-t border-plex-border/50">
            <label className="flex items-center gap-2 text-xs text-plex-text">
              <input
                type="checkbox"
                checked={taggerConfig.setStudio}
                onChange={(e) => setTaggerConfig((c) => ({ ...c, setStudio: e.target.checked }))}
                className="rounded border-plex-border"
              />
              Set studio
            </label>
            {taggerConfig.setStudio && (
              <>
                <label className="flex items-center gap-2 text-xs text-plex-text">
                  <input
                    type="checkbox"
                    checked={taggerConfig.onlyExistingStudio}
                    onChange={(e) => setTaggerConfig((c) => ({ ...c, onlyExistingStudio: e.target.checked }))}
                    className="rounded border-plex-border"
                  />
                  Only add existing studio
                </label>
                <label className="flex items-center gap-2 text-xs text-plex-text">
                  <input
                    type="checkbox"
                    checked={taggerConfig.createParentStudios}
                    onChange={(e) => setTaggerConfig((c) => ({ ...c, createParentStudios: e.target.checked }))}
                    className="rounded border-plex-border"
                  />
                  Create parent studios
                </label>
              </>
            )}
          </div>
        </div>
      )}

      {/* Scene list */}
      <div className="divide-y divide-plex-border">
        {visibleScenes.map((scene) => (
          <TaggerSceneRow
            key={scene.id}
            scene={scene}
            state={searchStates[scene.id]}
            query={getSearchQuery(scene)}
            onQueryChange={(q) => setQueryOverrides((prev) => ({ ...prev, [scene.id]: q }))}
            onSearch={() => searchScene(scene)}
            onSearchFingerprints={() => searchSceneFingerprints(scene)}
            onUpdateState={(update) => updateSearchState(scene.id, update)}
            endpoint={taggerConfig.selectedEndpoint}
            taggerConfig={taggerConfig}
          />
        ))}
      </div>
    </div>
  );
}

/* ── Scene Tagger Row ── */

interface TaggerSceneRowProps {
  scene: Scene;
  state?: SceneSearchState;
  query: string;
  onQueryChange: (q: string) => void;
  onSearch: () => void;
  onSearchFingerprints: () => void;
  onUpdateState: (update: Partial<SceneSearchState>) => void;
  endpoint: string;
  taggerConfig: TaggerConfig;
}

function TaggerSceneRow({
  scene,
  state,
  query,
  onQueryChange,
  onSearch,
  onSearchFingerprints,
  onUpdateState,
  endpoint,
  taggerConfig,
}: TaggerSceneRowProps) {
  const file = scene.files[0];
  const screenshotUrl = scenes.screenshotUrl(scene.id);
  const [showDetails, setShowDetails] = useState(false);

  const selectedResult = state?.results?.[state.selectedIndex ?? 0];

  const importMut = useMutation({
    mutationFn: () => {
      const excludedTags = state?.excludedTags ? Array.from(state.excludedTags) : undefined;
      const excludedPerformers = state?.excludedPerformers ? Array.from(state.excludedPerformers) : undefined;
      const importReq: StashBoxSceneImportRequest = {
        endpoint,
        sceneId: selectedResult?.id ?? "",
        setCoverImage: taggerConfig.setCoverImage,
        setTags: taggerConfig.setTags,
        setPerformers: taggerConfig.setPerformers,
        setStudio: taggerConfig.setStudio && !state?.skipStudio,
        onlyExistingTags: taggerConfig.onlyExistingTags,
        onlyExistingPerformers: taggerConfig.onlyExistingPerformers,
        onlyExistingStudio: taggerConfig.onlyExistingStudio,
        markOrganized: taggerConfig.markOrganized,
        excludedTagNames: excludedTags,
        excludedPerformerNames: excludedPerformers,
      };
      return scenes.importFromStashBox(scene.id, importReq);
    },
    onSuccess: () => {
      onUpdateState({ saved: true });
    },
  });

  return (
    <div className={`px-4 py-3 ${state?.saved ? "opacity-50" : ""}`}>
      <div className="flex gap-4">
        {/* Scene preview */}
        <div className="flex-shrink-0 w-40">
          <div className="relative aspect-video bg-plex-card rounded overflow-hidden">
            <img
              src={screenshotUrl}
              alt=""
              className="w-full h-full object-cover"
              loading="lazy"
              onError={(e) => {
                (e.target as HTMLImageElement).style.display = "none";
              }}
            />
            {file && file.duration > 0 && (
              <span className="absolute bottom-1 right-1 text-[9px] text-white bg-black/70 px-1 rounded">
                {formatDuration(file.duration)}
              </span>
            )}
          </div>
          <p className="text-xs text-plex-text mt-1 truncate font-medium">
            {scene.title || file?.basename || "Untitled"}
          </p>
          {scene.studioName && (
            <p className="text-[10px] text-plex-text-muted truncate">{scene.studioName}</p>
          )}
          {file && (
            <p className="text-[10px] text-plex-text-muted">
              {getResolutionLabel(file.width, file.height)} · {formatDuration(file.duration)}
            </p>
          )}
          {/* Existing stash IDs */}
          {scene.stashIds && scene.stashIds.length > 0 && (
            <div className="flex flex-wrap gap-1 mt-1">
              {scene.stashIds.map((sid) => (
                <span
                  key={`${sid.endpoint}-${sid.stashId}`}
                  className="text-[9px] px-1.5 py-0.5 rounded bg-green-600/20 text-green-300"
                  title={sid.endpoint}
                >
                  <Fingerprint className="w-2.5 h-2.5 inline mr-0.5" />
                  {sid.stashId.substring(0, 8)}…
                </span>
              ))}
            </div>
          )}
        </div>

        {/* Search + Results */}
        <div className="flex-1 min-w-0">
          {/* Search input */}
          <div className="flex gap-2 mb-2">
            <div className="relative flex-1">
              <input
                type="text"
                value={query}
                onChange={(e) => onQueryChange(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && onSearch()}
                placeholder="Search query..."
                className="w-full bg-plex-input border border-plex-border rounded pl-3 pr-3 py-1.5 text-xs text-plex-text focus:outline-none focus:border-plex-accent placeholder:text-plex-text-muted"
              />
            </div>
            <button
              onClick={onSearch}
              disabled={state?.loading}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded text-xs font-medium bg-plex-accent text-white hover:bg-plex-accent-hover disabled:opacity-60"
            >
              {state?.loading ? (
                <Loader2 className="w-3.5 h-3.5 animate-spin" />
              ) : (
                <Search className="w-3.5 h-3.5" />
              )}
              Search
            </button>
            <button
              onClick={onSearchFingerprints}
              disabled={state?.loading}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded text-xs font-medium bg-plex-surface border border-plex-border text-plex-text hover:bg-plex-card disabled:opacity-60"
              title="Search by fingerprint only"
            >
              <Fingerprint className="w-3.5 h-3.5" />
              Fingerprint
            </button>
          </div>

          {/* Current scene details (collapsible) */}
          <button
            onClick={() => setShowDetails(!showDetails)}
            className="flex items-center gap-1 text-[10px] text-plex-text-muted hover:text-plex-text mb-2"
          >
            {showDetails ? <ChevronUp className="w-3 h-3" /> : <ChevronDown className="w-3 h-3" />}
            Scene Details
          </button>
          {showDetails && (
            <div className="bg-plex-surface rounded p-2 mb-2 text-[10px] text-plex-text-muted space-y-0.5">
              {scene.date && <p>Date: {scene.date}</p>}
              {scene.details && <p className="line-clamp-3">Details: {scene.details}</p>}
              {scene.performers.length > 0 && (
                <p>Performers: {scene.performers.map((p) => p.name).join(", ")}</p>
              )}
              {scene.tags.length > 0 && (
                <p>Tags: {scene.tags.map((t) => t.name).join(", ")}</p>
              )}
              {file?.path && <p className="truncate">File: {file.path}</p>}
            </div>
          )}

          {/* Error */}
          {state?.error && (
            <p className="text-xs text-red-400 mb-2">
              <AlertCircle className="w-3 h-3 inline mr-1" />
              {state.error}
            </p>
          )}

          {/* No results */}
          {state?.results && state.results.length === 0 && (
            <p className="text-xs text-plex-text-muted">No matches found.</p>
          )}

          {/* Results */}
          {state?.results && state.results.length > 0 && (
            <TaggerResults
              results={state.results}
              selectedIndex={state.selectedIndex ?? 0}
              onSelect={(i) => onUpdateState({ selectedIndex: i })}
              onSave={() => importMut.mutate()}
              saving={importMut.isPending}
              saved={state.saved}
              localDuration={file?.duration}
              excludedPerformers={state.excludedPerformers ?? new Set()}
              excludedTags={state.excludedTags ?? new Set()}
              skipStudio={state.skipStudio ?? false}
              onTogglePerformer={(name) => {
                const current = new Set(state.excludedPerformers ?? []);
                if (current.has(name)) current.delete(name);
                else current.add(name);
                onUpdateState({ excludedPerformers: current });
              }}
              onToggleTag={(name) => {
                const current = new Set(state.excludedTags ?? []);
                if (current.has(name)) current.delete(name);
                else current.add(name);
                onUpdateState({ excludedTags: current });
              }}
              onToggleStudio={() => onUpdateState({ skipStudio: !state.skipStudio })}
              taggerConfig={taggerConfig}
            />
          )}

          {/* Saved indicator */}
          {state?.saved && (
            <div className="flex items-center gap-1 mt-2 text-xs text-green-400">
              <Check className="w-3.5 h-3.5" />
              Saved successfully
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

/* ── Tagger Results ── */

interface TaggerResultsProps {
  results: StashBoxSceneMatch[];
  selectedIndex: number;
  onSelect: (index: number) => void;
  onSave: () => void;
  saving?: boolean;
  saved?: boolean;
  localDuration?: number;
  excludedPerformers: Set<string>;
  excludedTags: Set<string>;
  skipStudio: boolean;
  onTogglePerformer: (name: string) => void;
  onToggleTag: (name: string) => void;
  onToggleStudio: () => void;
  taggerConfig: TaggerConfig;
}

function TaggerResults({ results, selectedIndex, onSelect, onSave, saving, saved, localDuration, excludedPerformers, excludedTags, skipStudio, onTogglePerformer, onToggleTag, onToggleStudio, taggerConfig }: TaggerResultsProps) {
  return (
    <div className="space-y-1">
      {results.map((result, i) => (
        <TaggerResultRow
          key={`${result.endpoint}-${result.id}`}
          result={result}
          isSelected={i === selectedIndex}
          onClick={() => onSelect(i)}
          onSave={i === selectedIndex ? onSave : undefined}
          saving={i === selectedIndex ? saving : false}
          saved={saved}
          localDuration={localDuration}
          excludedPerformers={excludedPerformers}
          excludedTags={excludedTags}
          skipStudio={skipStudio}
          onTogglePerformer={i === selectedIndex ? onTogglePerformer : undefined}
          onToggleTag={i === selectedIndex ? onToggleTag : undefined}
          onToggleStudio={i === selectedIndex ? onToggleStudio : undefined}
          taggerConfig={taggerConfig}
        />
      ))}
    </div>
  );
}

function TaggerResultRow({
  result,
  isSelected,
  onClick,
  onSave,
  saving,
  saved,
  localDuration,
  excludedPerformers,
  excludedTags,
  skipStudio,
  onTogglePerformer,
  onToggleTag,
  onToggleStudio,
  taggerConfig,
}: {
  result: StashBoxSceneMatch;
  isSelected: boolean;
  onClick: () => void;
  onSave?: () => void;
  saving?: boolean;
  saved?: boolean;
  localDuration?: number;
  excludedPerformers: Set<string>;
  excludedTags: Set<string>;
  skipStudio: boolean;
  onTogglePerformer?: (name: string) => void;
  onToggleTag?: (name: string) => void;
  onToggleStudio?: () => void;
  taggerConfig: TaggerConfig;
}) {
  // Calculate duration difference
  const durationDiff = localDuration != null && result.duration != null
    ? Math.abs(localDuration - result.duration)
    : undefined;
  const durationMatch = durationDiff != null && durationDiff < 5; // within 5 seconds = match
  const durationConfidence = durationDiff == null
    ? undefined
    : durationMatch
      ? "Exact"
      : durationDiff < 30
        ? "Close"
        : "Weak";
  return (
    <div
      onClick={onClick}
      className={`rounded border cursor-pointer transition-colors ${
        isSelected
          ? "border-plex-accent bg-plex-card"
          : "border-plex-border bg-plex-surface hover:border-plex-accent/50"
      }`}
    >
      {/* Collapsed view */}
      <div className="flex items-center gap-3 p-2">
        {/* Cover thumbnail */}
        {result.imageUrl && (
          <img
            src={result.imageUrl}
            alt=""
            className="w-16 h-9 object-cover rounded flex-shrink-0"
            loading="lazy"
          />
        )}
        <div className="flex-1 min-w-0">
          <p className="text-xs font-medium text-plex-text truncate">
            {result.title || "Untitled"}
            {result.code && <span className="text-plex-text-muted ml-1">({result.code})</span>}
          </p>
          <div className="flex items-center gap-2 text-[10px] text-plex-text-muted">
            {result.studioName && <span>{result.studioName}</span>}
            {result.date && <span>{result.date}</span>}
            {result.performerNames.length > 0 && (
              <span className="truncate">{result.performerNames.join(", ")}</span>
            )}
          </div>
        </div>
        {/* Fingerprint match indicators + duration comparison */}
        <div className="flex items-center gap-1 flex-shrink-0">
          {durationDiff != null && (
            <span
              className={`text-[9px] px-1.5 py-0.5 rounded ${
                durationMatch
                  ? "bg-green-600/20 text-green-300"
                  : durationDiff < 30
                    ? "bg-yellow-600/20 text-yellow-300"
                    : "bg-red-600/20 text-red-300"
              }`}
              title={`Duration differs by ${Math.round(durationDiff)}s (local: ${formatDuration(localDuration!)}, remote: ${formatDuration(result.duration!)})`}
            >
              {durationConfidence} {Math.round(durationDiff)}s
            </span>
          )}
          {result.fingerprintAlgorithms.map((algo) => (
            <span
              key={algo}
              className="text-[9px] px-1.5 py-0.5 rounded bg-green-600/20 text-green-300"
              title={`${algo} match`}
            >
              <Fingerprint className="w-2.5 h-2.5 inline mr-0.5" />
              {algo}
            </span>
          ))}
        </div>
      </div>

      {/* Expanded view */}
      {isSelected && (
        <div className="border-t border-plex-border p-3">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {/* Left: Metadata fields */}
            <div className="space-y-1.5 text-xs">
              {result.title && (
                <FieldRow label="Title" value={result.title} />
              )}
              {result.code && (
                <FieldRow label="Code" value={result.code} />
              )}
              {result.date && (
                <FieldRow label="Date" value={result.date} />
              )}
              {result.director && (
                <FieldRow label="Director" value={result.director} />
              )}
              {result.duration != null && (
                <FieldRow
                  label="Duration"
                  value={`${formatDuration(result.duration)}${
                    durationDiff != null
                      ? ` (${durationConfidence?.toLowerCase()} match, ${Math.round(durationDiff)}s difference)`
                      : ""
                  }`}
                />
              )}
              {result.urls.length > 0 && (
                <FieldRow label="URLs" value={result.urls.join(", ")} />
              )}
              {result.stashBoxName && (
                <FieldRow label="Source" value={result.stashBoxName} />
              )}
            </div>

            {/* Right: Studio + Performers + Tags with Create/Skip */}
            <div className="space-y-2">
              {/* Studio with Create/Skip */}
              {result.studioName && taggerConfig.setStudio && (
                <div>
                  <p className="text-[10px] text-plex-text-muted mb-1">Studio</p>
                  <div className="flex items-center gap-1">
                    <span
                      className={`text-[10px] px-1.5 py-0.5 rounded border ${
                        skipStudio
                          ? "bg-red-600/10 text-red-400 border-red-600/30 line-through"
                          : "bg-purple-600/20 text-purple-300 border-purple-600/30"
                      }`}
                    >
                      {result.studioName}
                    </span>
                    {onToggleStudio && (
                      <button
                        onClick={(e) => { e.stopPropagation(); onToggleStudio(); }}
                        className={`p-0.5 rounded text-[10px] ${
                          skipStudio
                            ? "text-green-400 hover:bg-green-600/20"
                            : "text-red-400 hover:bg-red-600/20"
                        }`}
                        title={skipStudio ? "Include studio" : "Skip studio"}
                      >
                        {skipStudio ? <Plus className="w-3 h-3" /> : <Ban className="w-3 h-3" />}
                      </button>
                    )}
                  </div>
                </div>
              )}

              {/* Performers with Create/Skip */}
              {result.performerNames.length > 0 && taggerConfig.setPerformers && (
                <div>
                  <p className="text-[10px] text-plex-text-muted mb-1">Performers</p>
                  <div className="flex flex-wrap gap-1">
                    {result.performerNames.map((name) => {
                      const excluded = excludedPerformers.has(name);
                      return (
                        <span key={name} className="inline-flex items-center gap-0.5">
                          <span
                            className={`text-[10px] px-1.5 py-0.5 rounded ${
                              excluded
                                ? "bg-red-600/10 text-red-400 line-through"
                                : "bg-blue-600/20 text-blue-300"
                            }`}
                          >
                            {name}
                          </span>
                          {onTogglePerformer && (
                            <button
                              onClick={(e) => { e.stopPropagation(); onTogglePerformer(name); }}
                              className={`p-0.5 rounded ${
                                excluded
                                  ? "text-green-400 hover:bg-green-600/20"
                                  : "text-red-400 hover:bg-red-600/20"
                              }`}
                              title={excluded ? "Include performer" : "Skip performer"}
                            >
                              {excluded ? <Plus className="w-2.5 h-2.5" /> : <X className="w-2.5 h-2.5" />}
                            </button>
                          )}
                        </span>
                      );
                    })}
                  </div>
                </div>
              )}

              {/* Tags with include/exclude */}
              {result.tagNames.length > 0 && taggerConfig.setTags && (
                <div>
                  <p className="text-[10px] text-plex-text-muted mb-1">Tags</p>
                  <div className="flex flex-wrap gap-1">
                    {result.tagNames.map((name) => {
                      const excluded = excludedTags.has(name);
                      return (
                        <span key={name} className="inline-flex items-center gap-0.5">
                          <span
                            className={`text-[10px] px-1.5 py-0.5 rounded border ${
                              excluded
                                ? "bg-red-600/10 text-red-400 border-red-600/30 line-through"
                                : "bg-plex-card text-plex-text-muted border-plex-border"
                            }`}
                          >
                            {name}
                          </span>
                          {onToggleTag && (
                            <button
                              onClick={(e) => { e.stopPropagation(); onToggleTag(name); }}
                              className={`p-0.5 rounded ${
                                excluded
                                  ? "text-green-400 hover:bg-green-600/20"
                                  : "text-red-400 hover:bg-red-600/20"
                              }`}
                              title={excluded ? "Include tag" : "Exclude tag"}
                            >
                              {excluded ? <Plus className="w-2.5 h-2.5" /> : <X className="w-2.5 h-2.5" />}
                            </button>
                          )}
                        </span>
                      );
                    })}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Save button */}
          {onSave && !saved && (
            <div className="flex justify-end mt-3">
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onSave();
                }}
                disabled={saving}
                className="flex items-center gap-1.5 px-4 py-1.5 rounded text-xs font-medium bg-green-600 text-white hover:bg-green-500 disabled:opacity-60"
              >
                {saving ? (
                  <Loader2 className="w-3.5 h-3.5 animate-spin" />
                ) : (
                  <Check className="w-3.5 h-3.5" />
                )}
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
