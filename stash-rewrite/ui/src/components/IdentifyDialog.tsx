import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { metadata, system } from "../api/client";
import type { StashBox } from "../api/types";
import {
  Search,
  X,
  Loader2,
  Check,
  GripVertical,
  ChevronDown,
  ChevronUp,
  Info,
} from "lucide-react";

interface Props {
  open: boolean;
  onClose: () => void;
  sceneIds?: number[];
}

interface IdentifySource {
  id: string;
  name: string;
  type: "stash-box" | "scraper" | "auto-tag";
  enabled: boolean;
}

export function IdentifyDialog({ open, onClose, sceneIds }: Props) {
  const queryClient = useQueryClient();
  const { data: config } = useQuery({
    queryKey: ["system-config"],
    queryFn: system.getConfig,
  });

  const stashBoxes = config?.scraping?.stashBoxes ?? [];

  const [sources, setSources] = useState<IdentifySource[]>(() => {
    const src: IdentifySource[] = [];
    stashBoxes.forEach((box, i) => {
      src.push({
        id: `stash-box-${i}`,
        name: box.name || box.endpoint,
        type: "stash-box",
        enabled: true,
      });
    });
    src.push({
      id: "auto-tag",
      name: "Auto Tag (built-in)",
      type: "auto-tag",
      enabled: true,
    });
    return src;
  });

  const [showOptions, setShowOptions] = useState(false);
  const [setCoverImage, setSetCoverImage] = useState(true);
  const [setOrganized, setSetOrganized] = useState(false);
  const [skipMultipleMatches, setSkipMultipleMatches] = useState(true);
  const [skipSingleNamePerformers, setSkipSingleNamePerformers] = useState(true);

  const identifyMut = useMutation({
    mutationFn: () => {
      const enabledSources = sources.filter((s) => s.enabled).map((s) => s.name);
      return metadata.identify({
        sources: enabledSources.length > 0 ? enabledSources : undefined,
        sceneIds,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["jobs"] });
      onClose();
    },
  });

  const toggleSource = (id: string) => {
    setSources(sources.map((s) => (s.id === id ? { ...s, enabled: !s.enabled } : s)));
  };

  const moveSource = (index: number, direction: "up" | "down") => {
    const newSources = [...sources];
    const target = direction === "up" ? index - 1 : index + 1;
    if (target < 0 || target >= newSources.length) return;
    [newSources[index], newSources[target]] = [newSources[target], newSources[index]];
    setSources(newSources);
  };

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
      <div className="bg-plex-surface border border-plex-border rounded-2xl shadow-2xl w-full max-w-lg max-h-[85vh] overflow-hidden flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-plex-border">
          <div>
            <h2 className="text-lg font-bold text-plex-text flex items-center gap-2">
              <Search className="w-5 h-5 text-plex-accent" />
              Identify
            </h2>
            <p className="text-xs text-plex-text-secondary mt-0.5">
              {sceneIds
                ? `Identifying ${sceneIds.length} scene${sceneIds.length !== 1 ? "s" : ""}`
                : "Identifying all scenes"}
            </p>
          </div>
          <button onClick={onClose} className="text-plex-text-muted hover:text-plex-text">
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Body */}
        <div className="flex-1 overflow-y-auto px-5 py-4 space-y-5">
          {/* Sources */}
          <div>
            <h3 className="text-sm font-medium text-plex-text mb-2">Sources (first match wins)</h3>
            <div className="space-y-1.5">
              {sources.map((source, i) => (
                <div
                  key={source.id}
                  className={`flex items-center gap-3 px-3 py-2.5 rounded-lg border transition-colors ${
                    source.enabled
                      ? "bg-plex-card border-plex-border"
                      : "bg-plex-card/50 border-plex-border/50 opacity-60"
                  }`}
                >
                  <div className="flex flex-col gap-0.5 flex-shrink-0">
                    <button
                      onClick={() => moveSource(i, "up")}
                      disabled={i === 0}
                      className="text-plex-text-muted hover:text-plex-text disabled:opacity-30"
                    >
                      <ChevronUp className="w-3 h-3" />
                    </button>
                    <button
                      onClick={() => moveSource(i, "down")}
                      disabled={i === sources.length - 1}
                      className="text-plex-text-muted hover:text-plex-text disabled:opacity-30"
                    >
                      <ChevronDown className="w-3 h-3" />
                    </button>
                  </div>
                  <label className="flex items-center gap-2.5 flex-1 cursor-pointer min-w-0">
                    <input
                      type="checkbox"
                      checked={source.enabled}
                      onChange={() => toggleSource(source.id)}
                      className="h-4 w-4 rounded border-plex-border bg-plex-card text-plex-accent focus:ring-0"
                    />
                    <div className="min-w-0">
                      <div className="text-sm font-medium text-plex-text truncate">{source.name}</div>
                      <div className="text-xs text-plex-text-muted capitalize">{source.type.replace("-", " ")}</div>
                    </div>
                  </label>
                </div>
              ))}
            </div>
            {sources.length === 0 && (
              <div className="text-sm text-plex-text-muted text-center py-4">
                No sources available. Configure StashBox endpoints in Settings &gt; Metadata Providers.
              </div>
            )}
          </div>

          {/* Options */}
          <div>
            <button
              onClick={() => setShowOptions(!showOptions)}
              className="flex items-center gap-1.5 text-sm font-medium text-plex-text-secondary hover:text-plex-text"
            >
              {showOptions ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
              Options
            </button>
            {showOptions && (
              <div className="mt-3 space-y-2 pl-1">
                <label className="flex items-center gap-2 text-sm text-plex-text-secondary">
                  <input
                    type="checkbox"
                    checked={setCoverImage}
                    onChange={(e) => setSetCoverImage(e.target.checked)}
                    className="h-4 w-4 rounded border-plex-border bg-plex-card text-plex-accent focus:ring-0"
                  />
                  Set cover image from scraper
                </label>
                <label className="flex items-center gap-2 text-sm text-plex-text-secondary">
                  <input
                    type="checkbox"
                    checked={setOrganized}
                    onChange={(e) => setSetOrganized(e.target.checked)}
                    className="h-4 w-4 rounded border-plex-border bg-plex-card text-plex-accent focus:ring-0"
                  />
                  Mark identified scenes as organized
                </label>
                <label className="flex items-center gap-2 text-sm text-plex-text-secondary">
                  <input
                    type="checkbox"
                    checked={skipMultipleMatches}
                    onChange={(e) => setSkipMultipleMatches(e.target.checked)}
                    className="h-4 w-4 rounded border-plex-border bg-plex-card text-plex-accent focus:ring-0"
                  />
                  Skip scenes with multiple matches
                </label>
                <label className="flex items-center gap-2 text-sm text-plex-text-secondary">
                  <input
                    type="checkbox"
                    checked={skipSingleNamePerformers}
                    onChange={(e) => setSkipSingleNamePerformers(e.target.checked)}
                    className="h-4 w-4 rounded border-plex-border bg-plex-card text-plex-accent focus:ring-0"
                  />
                  Skip single-name performers
                </label>
                <div className="flex items-start gap-2 mt-2 p-2 bg-blue-900/20 border border-blue-700/30 rounded-lg">
                  <Info className="w-4 h-4 text-blue-400 mt-0.5 flex-shrink-0" />
                  <p className="text-xs text-blue-300">
                    Identified data will be merged with existing data by default. Fields that
                    already have values won't be overwritten.
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Footer */}
        <div className="flex justify-end gap-2 px-5 py-4 border-t border-plex-border">
          <button
            onClick={onClose}
            className="px-4 py-2 text-sm text-plex-text-secondary hover:text-plex-text transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={() => identifyMut.mutate()}
            disabled={identifyMut.isPending || sources.filter((s) => s.enabled).length === 0}
            className="inline-flex items-center gap-2 px-5 py-2 bg-plex-accent hover:bg-plex-accent-hover text-white rounded-lg font-medium disabled:opacity-50 transition-colors"
          >
            {identifyMut.isPending ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Search className="w-4 h-4" />
            )}
            Identify
          </button>
        </div>
      </div>
    </div>
  );
}
