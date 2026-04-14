/**
 * POC: Enhanced Delete Dialog — overrides the confirm-delete dialog.
 * Proves: dialog override capability.
 */

interface EnhancedDeleteDialogProps {
  title?: string;
  message?: string;
  entityType?: string;
  entityName?: string;
  onConfirm?: () => void;
  onCancel?: () => void;
  open?: boolean;
}

export function EnhancedDeleteDialog({
  title = "Confirm Delete",
  message,
  entityType,
  entityName,
  onConfirm,
  onCancel,
  open = true,
}: EnhancedDeleteDialogProps) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="bg-plex-card border border-plex-border rounded-lg shadow-2xl max-w-md w-full mx-4 overflow-hidden">
        {/* Header */}
        <div className="p-4 border-b border-plex-border bg-red-600/10">
          <div className="flex items-center gap-2">
            <span className="text-red-400 text-xl">⚠️</span>
            <h2 className="text-lg font-semibold text-plex-text">{title}</h2>
          </div>
          <span className="text-[10px] px-1.5 py-0.5 rounded bg-red-600/20 text-red-300 border border-red-600/30 mt-2 inline-block">
            Extension: Enhanced Delete Dialog
          </span>
        </div>

        {/* Body */}
        <div className="p-4 space-y-3">
          {entityType && entityName && (
            <div className="bg-plex-surface rounded p-3 border border-plex-border">
              <div className="text-xs text-plex-text-muted mb-1">Deleting {entityType}</div>
              <div className="text-plex-text font-medium">{entityName}</div>
            </div>
          )}

          <p className="text-sm text-plex-text-secondary">
            {message || "This action cannot be undone. All associated data will be permanently removed."}
          </p>

          <div className="text-xs text-plex-text-muted bg-plex-surface rounded p-2 border border-plex-border">
            <strong>What will be affected:</strong>
            <ul className="list-disc list-inside mt-1 space-y-0.5">
              <li>The {entityType || "item"} record will be permanently deleted</li>
              <li>Associated metadata and tags will be unlinked</li>
              <li>Generated thumbnails and previews will be removed</li>
              <li>Files on disk will NOT be deleted</li>
            </ul>
          </div>
        </div>

        {/* Footer */}
        <div className="flex justify-end gap-2 p-4 border-t border-plex-border">
          <button
            onClick={onCancel}
            className="px-4 py-2 text-sm text-plex-text-secondary hover:text-plex-text bg-plex-surface border border-plex-border rounded hover:bg-plex-border transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={onConfirm}
            className="px-4 py-2 text-sm text-white bg-red-600 rounded hover:bg-red-700 transition-colors"
          >
            Delete Permanently
          </button>
        </div>

        <div className="px-4 pb-3">
          <p className="text-[10px] text-plex-text-muted">
            This dialog was provided by the Enhanced Delete Dialog extension via UIDialogOverride.
          </p>
        </div>
      </div>
    </div>
  );
}
