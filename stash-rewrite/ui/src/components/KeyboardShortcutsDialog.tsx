import { Keyboard, X } from "lucide-react";

interface ShortcutSection {
  title: string;
  shortcuts: { keys: string; description: string }[];
}

const sections: ShortcutSection[] = [
  {
    title: "Global Navigation",
    shortcuts: [
      { keys: "/", description: "Focus page filter or global search" },
      { keys: "g s", description: "Go to Scenes" },
      { keys: "g i", description: "Go to Images" },
      { keys: "g v", description: "Go to Groups" },
      { keys: "g k", description: "Go to Markers" },
      { keys: "g l", description: "Go to Galleries" },
      { keys: "g p", description: "Go to Performers" },
      { keys: "g u", description: "Go to Studios" },
      { keys: "g t", description: "Go to Tags" },
      { keys: "g z", description: "Go to Settings" },
      { keys: "g d", description: "Go to Dashboard" },
      { keys: "?", description: "Show this help" },
    ],
  },
  {
    title: "List Page",
    shortcuts: [
      { keys: "v g", description: "Grid view" },
      { keys: "v l", description: "List view" },
      { keys: "v w", description: "Wall view" },
      { keys: "v t", description: "Tagger view (scenes)" },
      { keys: "s a", description: "Select all on page" },
      { keys: "s n", description: "Deselect all" },
      { keys: "f", description: "Open advanced filters" },
      { keys: "←", description: "Previous page" },
      { keys: "→", description: "Next page" },
      { keys: "Shift + ←/→", description: "Jump 10 pages" },
      { keys: "Ctrl + Home/End", description: "First / last page" },
    ],
  },
  {
    title: "Scene Detail",
    shortcuts: [
      { keys: ",", description: "Toggle theater mode" },
    ],
  },
  {
    title: "Video Player",
    shortcuts: [
      { keys: "Space", description: "Play / Pause" },
      { keys: "k", description: "Play / Pause" },
      { keys: "←/→", description: "Seek ±5s (Shift: ±10s)" },
      { keys: "↑/↓", description: "Volume ±10%" },
      { keys: "m", description: "Toggle mute" },
      { keys: "f", description: "Toggle fullscreen" },
      { keys: "0-9", description: "Seek to 0-90% duration" },
    ],
  },
];

function KeyCap({ children }: { children: string }) {
  return (
    <kbd className="inline-flex items-center justify-center min-w-[24px] h-6 px-1.5 rounded bg-plex-sidebar-bg border border-plex-border text-xs font-mono text-plex-text">
      {children}
    </kbd>
  );
}

function formatKeys(keys: string) {
  return keys.split(" ").map((k, i) => (
    <span key={i} className="inline-flex items-center gap-0.5">
      {i > 0 && <span className="text-plex-text-muted mx-0.5">then</span>}
      <KeyCap>{k}</KeyCap>
    </span>
  ));
}

export function KeyboardShortcutsDialog({ open, onClose }: { open: boolean; onClose: () => void }) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70" onClick={onClose}>
      <div
        className="bg-plex-card border border-plex-border rounded-xl shadow-2xl w-full max-w-3xl max-h-[80vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-plex-border sticky top-0 bg-plex-card z-10">
          <div className="flex items-center gap-2">
            <Keyboard className="w-5 h-5 text-plex-accent" />
            <h2 className="text-lg font-semibold text-plex-text">Keyboard Shortcuts</h2>
          </div>
          <button onClick={onClose} className="p-1.5 rounded-lg hover:bg-plex-sidebar-hover text-plex-text-muted hover:text-plex-text">
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Content */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-6">
          {sections.map((section) => (
            <div key={section.title}>
              <h3 className="text-sm font-semibold text-plex-accent mb-3 uppercase tracking-wider">{section.title}</h3>
              <div className="space-y-2">
                {section.shortcuts.map((s) => (
                  <div key={s.keys} className="flex items-center justify-between gap-4">
                    <span className="text-sm text-plex-text-secondary">{s.description}</span>
                    <div className="flex items-center gap-1 shrink-0">{formatKeys(s.keys)}</div>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
