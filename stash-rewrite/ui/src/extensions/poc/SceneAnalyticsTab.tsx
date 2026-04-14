/**
 * POC: Scene Analytics Tab — injected into scene detail page via tab contribution.
 * Proves: tab injection, extension API calls, stateful extension data.
 */
import { useState, useEffect } from "react";

interface AnalyticsData {
  sceneId: number;
  views: number;
  lastViewed: string | null;
}

export function SceneAnalyticsTab({ entityId }: { entityId?: number }) {
  const [data, setData] = useState<AnalyticsData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!entityId) return;
    fetch(`/api/ext/analytics/scene/${entityId}`)
      .then((r) => r.json())
      .then((d) => setData(d))
      .catch(() => setData(null))
      .finally(() => setLoading(false));
  }, [entityId]);

  const recordView = async () => {
    if (!entityId) return;
    await fetch(`/api/ext/analytics/scene/${entityId}/view`, { method: "POST" });
    const r = await fetch(`/api/ext/analytics/scene/${entityId}`);
    setData(await r.json());
  };

  if (loading) {
    return (
      <div className="p-6 text-plex-text-secondary">Loading analytics...</div>
    );
  }

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center gap-2 mb-4">
        <span className="text-xs px-2 py-0.5 rounded bg-blue-600/20 text-blue-300 border border-blue-600/30">
          Extension: Scene Analytics
        </span>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="bg-plex-card rounded-lg p-4 border border-plex-border">
          <div className="text-sm text-plex-text-secondary mb-1">Total Views</div>
          <div className="text-3xl font-bold text-plex-text">{data?.views ?? 0}</div>
        </div>
        <div className="bg-plex-card rounded-lg p-4 border border-plex-border">
          <div className="text-sm text-plex-text-secondary mb-1">Last Viewed</div>
          <div className="text-lg text-plex-text">
            {data?.lastViewed
              ? new Date(data.lastViewed).toLocaleString()
              : "Never"}
          </div>
        </div>
      </div>

      <button
        onClick={recordView}
        className="px-4 py-2 bg-plex-accent text-white rounded hover:bg-plex-accent-hover transition-colors"
      >
        Record View
      </button>

      <p className="text-xs text-plex-text-muted mt-4">
        This tab was injected by the Scene Analytics extension via UITabContribution.
        View data is persisted to the database via IExtensionStore.
      </p>
    </div>
  );
}
