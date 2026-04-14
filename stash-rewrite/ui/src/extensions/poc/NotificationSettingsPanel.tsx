/**
 * POC: Notification Settings Panel — injected into the Extensions settings tab.
 * Proves: settings panel contribution, extension settings persistence via IExtensionStore.
 */
import { useState, useEffect } from "react";
import { extensions as extApi } from "../../api/client";

const EXT_ID = "com.stash.notification-settings";

export function NotificationSettingsPanel() {
  const [settings, setSettings] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    extApi.getData(EXT_ID)
      .then(setSettings)
      .catch(() => setSettings({}))
      .finally(() => setLoading(false));
  }, []);

  const updateSetting = async (key: string, value: string) => {
    setSaving(true);
    try {
      await extApi.setData(EXT_ID, key, value);
      setSettings((prev) => ({ ...prev, [key]: value }));
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <div className="text-plex-text-secondary">Loading notification settings...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2 mb-2">
        <span className="text-xs px-2 py-0.5 rounded bg-teal-600/20 text-teal-300 border border-teal-600/30">
          Extension: Notification Settings
        </span>
      </div>

      <p className="text-sm text-plex-text-secondary">
        This settings panel was contributed by the Notification Settings extension via{" "}
        <code className="text-xs bg-plex-surface px-1 py-0.5 rounded">UISettingsPanel</code>.
        Values are persisted to the database via IExtensionStore.
      </p>

      <div className="space-y-3">
        <ToggleSetting
          label="Enable scan notifications"
          description="Show a notification when a library scan completes"
          value={settings["notify.scan.complete"] === "true"}
          onChange={(v) => updateSetting("notify.scan.complete", v ? "true" : "false")}
          disabled={saving}
        />
        <ToggleSetting
          label="Enable new content alerts"
          description="Alert when new scenes are detected"
          value={settings["notify.new.content"] === "true"}
          onChange={(v) => updateSetting("notify.new.content", v ? "true" : "false")}
          disabled={saving}
        />
        <ToggleSetting
          label="Enable error notifications"
          description="Show notifications for processing errors"
          value={settings["notify.errors"] !== "false"}
          onChange={(v) => updateSetting("notify.errors", v ? "true" : "false")}
          disabled={saving}
        />
      </div>
    </div>
  );
}

function ToggleSetting({
  label,
  description,
  value,
  onChange,
  disabled,
}: {
  label: string;
  description: string;
  value: boolean;
  onChange: (v: boolean) => void;
  disabled: boolean;
}) {
  return (
    <div className="flex items-center justify-between p-3 rounded bg-plex-surface border border-plex-border">
      <div>
        <div className="text-sm font-medium text-plex-text">{label}</div>
        <div className="text-xs text-plex-text-muted">{description}</div>
      </div>
      <button
        onClick={() => onChange(!value)}
        disabled={disabled}
        className={`relative w-10 h-5 rounded-full transition-colors ${
          value ? "bg-plex-accent" : "bg-plex-border"
        } ${disabled ? "opacity-50" : ""}`}
      >
        <span
          className={`absolute top-0.5 left-0.5 w-4 h-4 rounded-full bg-white transition-transform ${
            value ? "translate-x-5" : ""
          }`}
        />
      </button>
    </div>
  );
}
