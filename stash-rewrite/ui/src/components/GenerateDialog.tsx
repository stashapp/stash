import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { metadata, type GenerateOptions } from "../api/client";
import { X, Loader2, Clapperboard, Check } from "lucide-react";

interface Props {
  open: boolean;
  onClose: () => void;
  onOpenJobDrawer?: () => void;
  /** If provided, generate only for these IDs. Otherwise generates for all. */
  sceneIds?: number[];
  title?: string;
}

export function GenerateDialog({ open, onClose, onOpenJobDrawer, sceneIds, title }: Props) {
  const [opts, setOpts] = useState<GenerateOptions>({
    thumbnails: true,
    previews: false,
    sprites: false,
    markers: false,
    phashes: false,
    overwrite: false,
  });
  const [submitted, setSubmitted] = useState(false);

  const generateMut = useMutation({
    mutationFn: () => metadata.generate({ ...opts, sceneIds }),
    onSuccess: () => {
      setSubmitted(true);
    },
  });

  if (!open) return null;

  const toggle = (key: keyof GenerateOptions) =>
    setOpts((o) => ({ ...o, [key]: !o[key] }));

  const label = sceneIds?.length
    ? `Generate for ${sceneIds.length} scene${sceneIds.length !== 1 ? "s" : ""}`
    : "Generate All";

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70" onClick={onClose}>
      <div
        className="bg-plex-card border border-plex-border rounded-xl shadow-2xl w-full max-w-md"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-plex-border">
          <div className="flex items-center gap-2">
            <Clapperboard className="w-5 h-5 text-plex-accent" />
            <h2 className="text-lg font-semibold text-plex-text">{title ?? label}</h2>
          </div>
          <button onClick={onClose} className="p-1.5 rounded-lg hover:bg-plex-sidebar-hover text-plex-text-muted hover:text-plex-text">
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Options */}
        <div className="p-5 space-y-3">
          <p className="text-sm text-plex-text-secondary mb-4">Select what to generate:</p>

          <h4 className="text-xs font-semibold text-plex-text-muted uppercase tracking-wider mb-2">Scene Content</h4>
          {([
            ["thumbnails", "Thumbnails / Screenshots"],
            ["previews", "Video Previews"],
            ["sprites", "Sprite Sheets"],
            ["markers", "Marker Screenshots"],
            ["phashes", "Perceptual Hashes"],
          ] as const).map(([key, labelText]) => (
            <label key={key} className="flex items-center gap-3 cursor-pointer group">
              <input
                type="checkbox"
                checked={!!opts[key]}
                onChange={() => toggle(key)}
                className="w-4 h-4 rounded border-plex-border accent-plex-accent"
              />
              <span className="text-sm text-plex-text group-hover:text-plex-accent">{labelText}</span>
            </label>
          ))}

          <div className="border-t border-plex-border my-3" />

          <label className="flex items-center gap-3 cursor-pointer group">
            <input
              type="checkbox"
              checked={!!opts.overwrite}
              onChange={() => toggle("overwrite")}
              className="w-4 h-4 rounded border-plex-border accent-orange-500"
            />
            <span className="text-sm text-orange-400 group-hover:text-orange-300">Overwrite existing</span>
          </label>
        </div>

        {/* Actions */}
        <div className="flex items-center justify-end gap-2 px-5 py-4 border-t border-plex-border">
          {submitted ? (
            <>
              <div className="flex items-center gap-2 text-sm text-green-400 mr-auto">
                <Check className="w-4 h-4" />
                Job started
              </div>
              {onOpenJobDrawer && (
                <button
                  onClick={() => { onClose(); setSubmitted(false); onOpenJobDrawer(); }}
                  className="px-4 py-2 rounded-lg text-sm font-medium bg-plex-accent hover:bg-plex-accent-hover text-white"
                >
                  View Progress
                </button>
              )}
              <button
                onClick={() => { onClose(); setSubmitted(false); }}
                className="px-4 py-2 rounded-lg text-sm text-plex-text-secondary hover:text-plex-text hover:bg-plex-surface"
              >
                Close
              </button>
            </>
          ) : (
            <>
              <button
                onClick={onClose}
                className="px-4 py-2 rounded-lg text-sm text-plex-text-secondary hover:text-plex-text hover:bg-plex-surface"
              >
                Cancel
              </button>
              <button
                onClick={() => generateMut.mutate()}
                disabled={generateMut.isPending}
                className="px-4 py-2 rounded-lg text-sm font-medium bg-plex-accent hover:bg-plex-accent-hover text-white disabled:opacity-50 flex items-center gap-2"
              >
                {generateMut.isPending && <Loader2 className="w-4 h-4 animate-spin" />}
                Generate
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
