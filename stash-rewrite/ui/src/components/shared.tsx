export { RatingBadge } from "./Rating";

export function TagBadge({ name, onClick }: { name: string; onClick?: () => void }) {
  return (
    <span
      onClick={onClick}
      className={`inline-block px-2 py-0.5 rounded text-xs font-medium bg-plex-card text-plex-text-secondary border border-plex-border ${
        onClick ? "cursor-pointer hover:bg-plex-card-hover hover:text-plex-text" : ""
      }`}
    >
      {name}
    </span>
  );
}

export function formatDuration(seconds: number): string {
  if (!seconds || seconds <= 0) return "0:00";
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);
  if (h > 0) return `${h}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
  return `${m}:${s.toString().padStart(2, "0")}`;
}

export function formatFileSize(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
}

export function formatDate(dateStr?: string): string {
  if (!dateStr) return "";
  try {
    return new Date(dateStr).toLocaleDateString();
  } catch {
    return dateStr;
  }
}

export function getResolutionLabel(width: number, height: number): string | null {
  const number = width > height ? height : width;
  if (number >= 6144) return "HUGE";
  if (number >= 3840) return "8K";
  if (number >= 3584) return "7K";
  if (number >= 3000) return "6K";
  if (number >= 2560) return "5K";
  if (number >= 1920) return "4K";
  if (number >= 1440) return "1440p";
  if (number >= 1080) return "1080p";
  if (number >= 720) return "720p";
  if (number >= 540) return "540p";
  if (number >= 480) return "480p";
  if (number >= 360) return "360p";
  if (number >= 240) return "240p";
  if (number >= 144) return "144p";
  return null;
}

export function CustomFieldsDisplay({ customFields }: { customFields?: Record<string, unknown> }) {
  if (!customFields || Object.keys(customFields).length === 0) return null;
  return (
    <div className="bg-plex-card rounded-xl p-4">
      <h3 className="text-sm font-semibold text-plex-text-secondary mb-3">Custom Fields</h3>
      <div className="grid grid-cols-2 gap-2 text-sm">
        {Object.entries(customFields).map(([key, value]) => (
          <div key={key} className="flex flex-col">
            <span className="text-plex-text-muted text-xs">{key}</span>
            <span className="text-plex-text">{String(value ?? "")}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

export function CustomFieldsEditor({ value, onChange }: { value: Record<string, string>; onChange: (v: Record<string, string>) => void }) {
  const entries = Object.entries(value);
  const addField = () => onChange({ ...value, "": "" });
  const removeField = (key: string) => {
    const next = { ...value };
    delete next[key];
    onChange(next);
  };
  const updateKey = (oldKey: string, newKey: string) => {
    const next: Record<string, string> = {};
    for (const [k, v] of Object.entries(value)) {
      next[k === oldKey ? newKey : k] = v;
    }
    onChange(next);
  };
  const updateValue = (key: string, newVal: string) => {
    onChange({ ...value, [key]: newVal });
  };

  return (
    <div className="space-y-2">
      {entries.map(([key, val], i) => (
        <div key={i} className="flex gap-2 items-center">
          <input value={key} onChange={(e) => updateKey(key, e.target.value)} placeholder="Field name" className="flex-1 rounded border border-plex-border bg-plex-surface px-2 py-1 text-sm text-plex-text" />
          <input value={val} onChange={(e) => updateValue(key, e.target.value)} placeholder="Value" className="flex-1 rounded border border-plex-border bg-plex-surface px-2 py-1 text-sm text-plex-text" />
          <button type="button" onClick={() => removeField(key)} className="text-red-400 hover:text-red-300 text-sm px-1">×</button>
        </div>
      ))}
      <button type="button" onClick={addField} className="text-xs text-plex-accent hover:underline">+ Add Field</button>
    </div>
  );
}
