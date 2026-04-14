/**
 * POC: System Tools Page — a new page added to the navigation via extension.
 * Proves: new page contribution, nav integration, extension API endpoints.
 */
import { useState, useEffect } from "react";
import { useExtensions } from "../ExtensionLoader";

interface SystemInfo {
  runtime: string;
  os: string;
  cpuCount: number;
  workingSet: number;
  uptime: number;
  gcMemory: number;
}

interface ExtInfo {
  id: string;
  name: string;
  version: string;
  description?: string;
  enabled: boolean;
  capabilities: {
    ui: boolean;
    api: boolean;
    stateful: boolean;
    jobs: boolean;
    events: boolean;
  };
}

export function SystemToolsPage() {
  const [systemInfo, setSystemInfo] = useState<SystemInfo | null>(null);
  const [extList, setExtList] = useState<ExtInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const { availableThemes, activeThemeId } = useExtensions();

  useEffect(() => {
    Promise.all([
      fetch("/api/ext/system-tools/info").then((r) => r.json()),
      fetch("/api/ext/system-tools/extensions").then((r) => r.json()),
    ])
      .then(([info, exts]) => {
        setSystemInfo(info);
        setExtList(exts);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="p-6 text-plex-text-secondary">Loading system tools...</div>;
  }

  return (
    <div className="p-6 max-w-5xl mx-auto space-y-6 overflow-y-auto" style={{ height: "calc(100vh - 48px)" }}>
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-bold text-plex-text">🔧 System Tools</h1>
        <span className="text-xs px-2 py-0.5 rounded bg-amber-600/20 text-amber-300 border border-amber-600/30">
          Extension Page
        </span>
      </div>

      <p className="text-plex-text-secondary">
        This entire page was added by the System Tools extension via{" "}
        <code className="text-xs bg-plex-surface px-1 py-0.5 rounded">UIPageDefinition</code>.
        It appears in the nav sidebar and is rendered by the extension runtime.
      </p>

      {/* System Info */}
      {systemInfo && (
        <div className="bg-plex-card rounded-lg border border-plex-border p-5">
          <h2 className="text-lg font-semibold text-plex-text mb-4">System Information</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
            <InfoItem label="Runtime" value={systemInfo.runtime} />
            <InfoItem label="OS" value={systemInfo.os} />
            <InfoItem label="CPU Cores" value={systemInfo.cpuCount.toString()} />
            <InfoItem label="Working Set" value={`${(systemInfo.workingSet / 1024 / 1024).toFixed(1)} MB`} />
            <InfoItem label="GC Memory" value={`${(systemInfo.gcMemory / 1024 / 1024).toFixed(1)} MB`} />
            <InfoItem label="Uptime" value={formatUptime(systemInfo.uptime)} />
          </div>
        </div>
      )}

      {/* Extension Registry */}
      <div className="bg-plex-card rounded-lg border border-plex-border p-5">
        <h2 className="text-lg font-semibold text-plex-text mb-4">
          Registered Extensions ({extList.length})
        </h2>
        <div className="space-y-3">
          {extList.map((ext) => (
            <div
              key={ext.id}
              className="flex items-start gap-3 p-3 rounded bg-plex-surface border border-plex-border"
            >
              <div className={`w-2 h-2 mt-2 rounded-full ${ext.enabled ? "bg-green-500" : "bg-gray-500"}`} />
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 flex-wrap">
                  <span className="font-medium text-plex-text">{ext.name}</span>
                  <span className="text-xs text-plex-text-muted">v{ext.version}</span>
                </div>
                {ext.description && (
                  <p className="text-sm text-plex-text-secondary mt-0.5">{ext.description}</p>
                )}
                <div className="flex gap-1.5 mt-2 flex-wrap">
                  {ext.capabilities.ui && <CapBadge label="UI" color="blue" />}
                  {ext.capabilities.api && <CapBadge label="API" color="purple" />}
                  {ext.capabilities.stateful && <CapBadge label="Stateful" color="green" />}
                  {ext.capabilities.jobs && <CapBadge label="Jobs" color="amber" />}
                  {ext.capabilities.events && <CapBadge label="Events" color="rose" />}
                </div>
              </div>
              <span className="text-xs text-plex-text-muted font-mono">{ext.id}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Theme Info */}
      <div className="bg-plex-card rounded-lg border border-plex-border p-5">
        <h2 className="text-lg font-semibold text-plex-text mb-4">Available Themes ({availableThemes.length})</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
          {availableThemes.map((theme) => (
            <div
              key={theme.id}
              className={`p-3 rounded border ${
                theme.id === activeThemeId
                  ? "border-plex-accent bg-plex-accent/10"
                  : "border-plex-border bg-plex-surface"
              }`}
            >
              <div className="font-medium text-plex-text text-sm">{theme.name}</div>
              {theme.description && (
                <div className="text-xs text-plex-text-muted mt-1">{theme.description}</div>
              )}
              {theme.id === activeThemeId && (
                <div className="text-xs text-plex-accent mt-1">Active</div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function InfoItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <div className="text-xs text-plex-text-muted mb-0.5">{label}</div>
      <div className="text-sm text-plex-text truncate">{value}</div>
    </div>
  );
}

function CapBadge({ label, color }: { label: string; color: string }) {
  return (
    <span className={`text-[10px] px-1.5 py-0.5 rounded bg-${color}-600/20 text-${color}-300 border border-${color}-600/30`}>
      {label}
    </span>
  );
}

function formatUptime(ms: number): string {
  const s = Math.floor(ms / 1000);
  const m = Math.floor(s / 60);
  const h = Math.floor(m / 60);
  const d = Math.floor(h / 24);
  if (d > 0) return `${d}d ${h % 24}h`;
  if (h > 0) return `${h}h ${m % 60}m`;
  return `${m}m ${s % 60}s`;
}
