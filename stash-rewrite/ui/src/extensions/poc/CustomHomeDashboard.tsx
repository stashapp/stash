/**
 * POC: Custom Home Dashboard — replaces the built-in home page via page override.
 * Proves: page replacement capability without modifying core logic.
 */
import { useState, useEffect } from "react";

interface SystemInfo {
  runtime: string;
  os: string;
  cpuCount: number;
  workingSet: number;
  gcMemory: number;
}

export function CustomHomeDashboard() {
  const [systemInfo, setSystemInfo] = useState<SystemInfo | null>(null);

  useEffect(() => {
    fetch("/api/ext/system-tools/info")
      .then((r) => r.json())
      .then(setSystemInfo)
      .catch(() => {});
  }, []);

  return (
    <div className="p-6 max-w-6xl mx-auto space-y-6">
      <div className="flex items-center gap-3 mb-6">
        <h1 className="text-2xl font-bold text-plex-text">Dashboard</h1>
        <span className="text-xs px-2 py-0.5 rounded bg-green-600/20 text-green-300 border border-green-600/30">
          Extension: Custom Home
        </span>
      </div>

      <p className="text-plex-text-secondary">
        This page replaces the default home page via{" "}
        <code className="text-xs bg-plex-surface px-1 py-0.5 rounded">UIPageOverride</code>.
        The built-in home page component was not modified — this extension's component
        takes priority.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <DashboardCard title="Quick Stats" icon="📊">
          <p className="text-plex-text-secondary text-sm">
            Extension-provided dashboard widgets can show any data from extension APIs.
          </p>
        </DashboardCard>

        <DashboardCard title="Recent Activity" icon="🕐">
          <p className="text-plex-text-secondary text-sm">
            Extensions with IEventExtension can track entity changes and display them here.
          </p>
        </DashboardCard>

        <DashboardCard title="System" icon="⚙️">
          {systemInfo ? (
            <div className="text-sm space-y-1">
              <div className="flex justify-between">
                <span className="text-plex-text-secondary">Runtime</span>
                <span className="text-plex-text text-xs">{systemInfo.runtime}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-plex-text-secondary">CPUs</span>
                <span className="text-plex-text">{systemInfo.cpuCount}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-plex-text-secondary">Memory</span>
                <span className="text-plex-text">
                  {(systemInfo.workingSet / 1024 / 1024).toFixed(0)} MB
                </span>
              </div>
            </div>
          ) : (
            <p className="text-plex-text-muted text-sm">Loading...</p>
          )}
        </DashboardCard>
      </div>

      <div className="bg-plex-card rounded-lg border border-plex-border p-6">
        <h2 className="text-lg font-semibold text-plex-text mb-3">How Page Override Works</h2>
        <div className="text-sm text-plex-text-secondary space-y-2">
          <p>
            The <code className="bg-plex-surface px-1 py-0.5 rounded">CustomHomeExtension</code> backend
            registers a <code className="bg-plex-surface px-1 py-0.5 rounded">UIPageOverride</code> targeting
            the "home" page.
          </p>
          <p>
            The frontend <code className="bg-plex-surface px-1 py-0.5 rounded">ExtensionLoader</code> resolves
            the override and the <code className="bg-plex-surface px-1 py-0.5 rounded">App.tsx</code> router
            renders this component instead of the default home page.
          </p>
          <p>
            No core files were modified to support this — the extension system handles the routing.
          </p>
        </div>
      </div>
    </div>
  );
}

function DashboardCard({ title, icon, children }: { title: string; icon: string; children: React.ReactNode }) {
  return (
    <div className="bg-plex-card rounded-lg border border-plex-border p-4">
      <div className="flex items-center gap-2 mb-3">
        <span className="text-lg">{icon}</span>
        <h3 className="font-medium text-plex-text">{title}</h3>
      </div>
      {children}
    </div>
  );
}
