import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { system } from "../api/client";
import type { StashConfig, StashPathConfig } from "../api/types";
import {
  FolderOpen,
  Plus,
  Trash2,
  ChevronRight,
  ChevronLeft,
  Check,
  Loader2,
  Play,
  Settings,
} from "lucide-react";

interface Props {
  config: StashConfig;
  onComplete: () => void;
}

type Step = "welcome" | "paths" | "confirm" | "done";

export function SetupWizardPage({ config, onComplete }: Props) {
  const [step, setStep] = useState<Step>("welcome");
  const [paths, setPaths] = useState<StashPathConfig[]>(
    config.stashPaths.length > 0
      ? config.stashPaths
      : [{ path: "", excludeVideo: false, excludeImage: false }]
  );
  const [error, setError] = useState<string | null>(null);
  const queryClient = useQueryClient();

  const saveMut = useMutation({
    mutationFn: (cfg: StashConfig) => system.saveConfig(cfg),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["system-config"] });
      setStep("done");
    },
    onError: (err: Error) => setError(err.message),
  });

  const addPath = () => {
    setPaths([...paths, { path: "", excludeVideo: false, excludeImage: false }]);
  };

  const removePath = (index: number) => {
    setPaths(paths.filter((_, i) => i !== index));
  };

  const updatePath = (index: number, updates: Partial<StashPathConfig>) => {
    setPaths(paths.map((p, i) => (i === index ? { ...p, ...updates } : p)));
  };

  const validPaths = paths.filter((p) => p.path.trim() !== "");

  const handleConfirm = () => {
    const updatedConfig: StashConfig = {
      ...config,
      stashPaths: validPaths,
    };
    saveMut.mutate(updatedConfig);
  };

  return (
    <div className="min-h-screen bg-plex-bg flex items-center justify-center p-4">
      <div className="w-full max-w-2xl">
        {/* Progress indicator */}
        <div className="flex items-center justify-center gap-2 mb-8">
          {(["welcome", "paths", "confirm", "done"] as Step[]).map((s, i) => (
            <div key={s} className="flex items-center gap-2">
              <div
                className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-colors ${
                  s === step
                    ? "bg-plex-accent text-white"
                    : (["welcome", "paths", "confirm", "done"].indexOf(step) > i)
                    ? "bg-green-600 text-white"
                    : "bg-plex-card border border-plex-border text-plex-text-muted"
                }`}
              >
                {(["welcome", "paths", "confirm", "done"].indexOf(step) > i) ? (
                  <Check className="w-4 h-4" />
                ) : (
                  i + 1
                )}
              </div>
              {i < 3 && <div className="w-12 h-0.5 bg-plex-border" />}
            </div>
          ))}
        </div>

        <div className="bg-plex-surface border border-plex-border rounded-2xl shadow-2xl overflow-hidden">
          {step === "welcome" && (
            <div className="p-8 text-center">
              <div className="w-16 h-16 bg-plex-accent/20 rounded-2xl flex items-center justify-center mx-auto mb-6">
                <Play className="w-8 h-8 text-plex-accent" />
              </div>
              <h1 className="text-2xl font-bold text-plex-text mb-3">Welcome to Stash</h1>
              <p className="text-plex-text-secondary mb-6 max-w-md mx-auto">
                Stash is a self-hosted organizer for your media library. Let's get set up
                by configuring your library directories.
              </p>
              <button
                onClick={() => setStep("paths")}
                className="inline-flex items-center gap-2 px-6 py-3 bg-plex-accent hover:bg-plex-accent-hover text-white rounded-lg font-medium transition-colors"
              >
                Get Started <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          )}

          {step === "paths" && (
            <div className="p-8">
              <h2 className="text-xl font-bold text-plex-text mb-2">Library Paths</h2>
              <p className="text-sm text-plex-text-secondary mb-6">
                Add the directories containing your media files. Stash will scan these
                directories for scenes, images, and galleries.
              </p>

              <div className="space-y-3 mb-4">
                {paths.map((p, i) => (
                  <div key={i} className="space-y-2">
                    <div className="flex gap-2">
                      <div className="flex-1 flex items-center gap-2 bg-plex-card border border-plex-border rounded-lg px-3 py-2">
                        <FolderOpen className="w-4 h-4 text-plex-text-muted flex-shrink-0" />
                        <input
                          type="text"
                          value={p.path}
                          onChange={(e) => updatePath(i, { path: e.target.value })}
                          placeholder="Enter directory path (e.g., /media/videos)"
                          className="flex-1 bg-transparent outline-none text-sm text-plex-text"
                        />
                      </div>
                      {paths.length > 1 && (
                        <button
                          onClick={() => removePath(i)}
                          className="p-2 text-plex-text-muted hover:text-red-400 transition-colors"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      )}
                    </div>
                    <div className="flex gap-4 pl-8">
                      <label className="flex items-center gap-1.5 text-xs text-plex-text-secondary">
                        <input
                          type="checkbox"
                          checked={p.excludeVideo}
                          onChange={(e) => updatePath(i, { excludeVideo: e.target.checked })}
                          className="h-3.5 w-3.5 rounded border-plex-border bg-plex-card text-plex-accent focus:ring-0"
                        />
                        Exclude video
                      </label>
                      <label className="flex items-center gap-1.5 text-xs text-plex-text-secondary">
                        <input
                          type="checkbox"
                          checked={p.excludeImage}
                          onChange={(e) => updatePath(i, { excludeImage: e.target.checked })}
                          className="h-3.5 w-3.5 rounded border-plex-border bg-plex-card text-plex-accent focus:ring-0"
                        />
                        Exclude images
                      </label>
                    </div>
                  </div>
                ))}
              </div>

              <button
                onClick={addPath}
                className="flex items-center gap-1.5 text-sm text-plex-accent hover:text-plex-accent-hover transition-colors mb-6"
              >
                <Plus className="w-4 h-4" /> Add another path
              </button>

              <div className="flex justify-between">
                <button
                  onClick={() => setStep("welcome")}
                  className="flex items-center gap-1.5 px-4 py-2 text-sm text-plex-text-secondary hover:text-plex-text transition-colors"
                >
                  <ChevronLeft className="w-4 h-4" /> Back
                </button>
                <button
                  onClick={() => setStep("confirm")}
                  disabled={validPaths.length === 0}
                  className="inline-flex items-center gap-2 px-5 py-2 bg-plex-accent hover:bg-plex-accent-hover text-white rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  Next <ChevronRight className="w-4 h-4" />
                </button>
              </div>
            </div>
          )}

          {step === "confirm" && (
            <div className="p-8">
              <h2 className="text-xl font-bold text-plex-text mb-2">Confirm Configuration</h2>
              <p className="text-sm text-plex-text-secondary mb-6">
                Review your library paths before saving. You can always change these later in Settings.
              </p>

              <div className="bg-plex-card border border-plex-border rounded-xl p-4 mb-6">
                <h3 className="text-xs font-medium uppercase tracking-wide text-plex-text-muted mb-3">Library Paths</h3>
                <div className="space-y-2">
                  {validPaths.map((p, i) => (
                    <div key={i} className="flex items-center gap-3">
                      <FolderOpen className="w-4 h-4 text-plex-accent flex-shrink-0" />
                      <span className="text-sm text-plex-text font-mono">{p.path}</span>
                      {(p.excludeVideo || p.excludeImage) && (
                        <span className="text-xs text-plex-text-muted">
                          (excludes: {[p.excludeVideo && "video", p.excludeImage && "images"].filter(Boolean).join(", ")})
                        </span>
                      )}
                    </div>
                  ))}
                </div>
              </div>

              {error && (
                <div className="bg-red-900/20 border border-red-700/50 rounded-lg p-3 mb-4 text-sm text-red-300">
                  {error}
                </div>
              )}

              <div className="flex justify-between">
                <button
                  onClick={() => { setStep("paths"); setError(null); }}
                  className="flex items-center gap-1.5 px-4 py-2 text-sm text-plex-text-secondary hover:text-plex-text transition-colors"
                >
                  <ChevronLeft className="w-4 h-4" /> Back
                </button>
                <button
                  onClick={handleConfirm}
                  disabled={saveMut.isPending}
                  className="inline-flex items-center gap-2 px-5 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg font-medium disabled:opacity-50 transition-colors"
                >
                  {saveMut.isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : <Check className="w-4 h-4" />}
                  Save & Continue
                </button>
              </div>
            </div>
          )}

          {step === "done" && (
            <div className="p-8 text-center">
              <div className="w-16 h-16 bg-green-600/20 rounded-2xl flex items-center justify-center mx-auto mb-6">
                <Check className="w-8 h-8 text-green-400" />
              </div>
              <h2 className="text-2xl font-bold text-plex-text mb-3">You're all set!</h2>
              <p className="text-plex-text-secondary mb-2 max-w-md mx-auto">
                Your library paths have been configured. Head to Settings &gt; Tasks to start
                scanning for content.
              </p>
              <p className="text-xs text-plex-text-muted mb-6">
                You can add more paths, configure scrapers, and set up StashBox connections in Settings.
              </p>
              <div className="flex justify-center gap-3">
                <button
                  onClick={onComplete}
                  className="inline-flex items-center gap-2 px-6 py-3 bg-plex-accent hover:bg-plex-accent-hover text-white rounded-lg font-medium transition-colors"
                >
                  Go to Dashboard <ChevronRight className="w-4 h-4" />
                </button>
                <button
                  onClick={() => { onComplete(); window.location.hash = "#/settings"; }}
                  className="inline-flex items-center gap-2 px-6 py-3 bg-plex-card border border-plex-border text-plex-text-secondary hover:text-plex-text rounded-lg font-medium transition-colors"
                >
                  <Settings className="w-4 h-4" /> Open Settings
                </button>
              </div>
            </div>
          )}
        </div>

        {/* Skip link */}
        {step !== "done" && (
          <div className="text-center mt-4">
            <button
              onClick={onComplete}
              className="text-xs text-plex-text-muted hover:text-plex-text-secondary transition-colors"
            >
              Skip setup for now
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
