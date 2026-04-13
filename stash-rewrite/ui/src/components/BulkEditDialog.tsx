import { useState, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { X, Plus, Minus } from "lucide-react";
import { tags as tagsApi, performers as performersApi, studios as studiosApi, groups as groupsApi } from "../api/client";
import type { BulkUpdateMode } from "../api/types";

// ===== Generic Bulk Edit Dialog =====

interface BulkEditField {
  key: string;
  label: string;
  type: "number" | "bool" | "string" | "date" | "select" | "multiId";
  entityType?: "tags" | "performers" | "studios" | "groups";
  options?: { label: string; value: string | number }[];
  nullable?: boolean;
}

interface BulkEditDialogProps {
  open: boolean;
  onClose: () => void;
  title: string;
  selectedCount: number;
  fields: BulkEditField[];
  onApply: (values: Record<string, unknown>) => void;
  isPending?: boolean;
}

export function BulkEditDialog({ open, onClose, title, selectedCount, fields, onApply, isPending }: BulkEditDialogProps) {
  const [values, setValues] = useState<Record<string, unknown>>({});
  const [enabledFields, setEnabledFields] = useState<Set<string>>(new Set());

  const toggleField = (key: string) => {
    setEnabledFields((prev) => {
      const next = new Set(prev);
      if (next.has(key)) {
        next.delete(key);
        setValues((v) => { const { [key]: _, [`${key}Mode`]: __, ...rest } = v; return rest; });
      } else {
        next.add(key);
      }
      return next;
    });
  };

  const updateValue = (key: string, val: unknown) => {
    setValues((prev) => ({ ...prev, [key]: val }));
  };

  const handleApply = () => {
    const result: Record<string, unknown> = {};
    for (const f of fields) {
      if (enabledFields.has(f.key)) {
        result[f.key] = values[f.key];
        if (f.type === "multiId") {
          result[`${f.key}Mode`] = values[`${f.key}Mode`] ?? "ADD";
        }
      }
    }
    onApply(result);
  };

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60" onClick={onClose}>
      <div className="bg-plex-surface border border-plex-border rounded-lg shadow-xl w-full max-w-md max-h-[80vh] flex flex-col" onClick={(e) => e.stopPropagation()}>
        <div className="flex items-center justify-between px-4 py-3 border-b border-plex-border">
          <h2 className="text-sm font-semibold text-plex-text">
            {title} <span className="text-plex-text-muted font-normal">({selectedCount} selected)</span>
          </h2>
          <button onClick={onClose} className="p-1 hover:bg-plex-card rounded text-plex-text-muted hover:text-plex-text">
            <X className="w-4 h-4" />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto px-4 py-3 space-y-3">
          {fields.map((field) => (
            <BulkFieldEditor
              key={field.key}
              field={field}
              enabled={enabledFields.has(field.key)}
              onToggle={() => toggleField(field.key)}
              value={values[field.key]}
              mode={(values[`${field.key}Mode`] as BulkUpdateMode) ?? "ADD"}
              onValueChange={(v) => updateValue(field.key, v)}
              onModeChange={(m) => updateValue(`${field.key}Mode`, m)}
            />
          ))}
        </div>

        <div className="flex items-center justify-end gap-2 px-4 py-3 border-t border-plex-border">
          <button onClick={onClose} className="px-3 py-1 rounded text-xs text-plex-text-secondary hover:text-plex-text border border-plex-border">
            Cancel
          </button>
          <button
            onClick={handleApply}
            disabled={isPending || enabledFields.size === 0}
            className="px-4 py-1 rounded text-xs font-medium bg-plex-accent hover:bg-plex-accent-hover text-white disabled:opacity-50"
          >
            {isPending ? "Applying..." : "Apply"}
          </button>
        </div>
      </div>
    </div>
  );
}

function BulkFieldEditor({
  field,
  enabled,
  onToggle,
  value,
  mode,
  onValueChange,
  onModeChange,
}: {
  field: BulkEditField;
  enabled: boolean;
  onToggle: () => void;
  value: unknown;
  mode: BulkUpdateMode;
  onValueChange: (v: unknown) => void;
  onModeChange: (m: BulkUpdateMode) => void;
}) {
  return (
    <div>
      <label className="flex items-center gap-2 cursor-pointer">
        <input
          type="checkbox"
          checked={enabled}
          onChange={onToggle}
          className="w-3.5 h-3.5 rounded border-plex-border accent-plex-accent"
        />
        <span className={`text-xs font-medium ${enabled ? "text-plex-text" : "text-plex-text-muted"}`}>
          {field.label}
        </span>
      </label>
      {enabled && (
        <div className="ml-6 mt-1">
          {field.type === "number" && (
            <input
              type="number"
              value={(value as number) ?? ""}
              onChange={(e) => onValueChange(e.target.value ? Number(e.target.value) : undefined)}
              className="w-24 bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent"
            />
          )}
          {field.type === "bool" && (
            <div className="flex gap-2">
              <button
                onClick={() => onValueChange(true)}
                className={`px-3 py-1 rounded text-xs border ${value === true ? "bg-plex-accent text-white border-plex-accent" : "border-plex-border text-plex-text-secondary"}`}
              >
                True
              </button>
              <button
                onClick={() => onValueChange(false)}
                className={`px-3 py-1 rounded text-xs border ${value === false ? "bg-plex-accent text-white border-plex-accent" : "border-plex-border text-plex-text-secondary"}`}
              >
                False
              </button>
            </div>
          )}
          {field.type === "string" && (
            <input
              type="text"
              value={(value as string) ?? ""}
              onChange={(e) => onValueChange(e.target.value || undefined)}
              className="w-full bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent"
            />
          )}
          {field.type === "date" && (
            <input
              type="date"
              value={(value as string) ?? ""}
              onChange={(e) => onValueChange(e.target.value || undefined)}
              className="bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent"
            />
          )}
          {field.type === "select" && field.options && (
            <select
              value={String(value ?? "")}
              onChange={(e) => onValueChange(e.target.value || undefined)}
              className="bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent"
            >
              <option value="">—</option>
              {field.options.map((o) => (
                <option key={o.value} value={o.value}>{o.label}</option>
              ))}
            </select>
          )}
          {field.type === "multiId" && field.entityType && (
            <MultiIdBulkEditor
              entityType={field.entityType}
              value={(value as number[]) ?? []}
              mode={mode}
              onValueChange={onValueChange}
              onModeChange={onModeChange}
            />
          )}
        </div>
      )}
    </div>
  );
}

function MultiIdBulkEditor({
  entityType,
  value,
  mode,
  onValueChange,
  onModeChange,
}: {
  entityType: "tags" | "performers" | "studios" | "groups";
  value: number[];
  mode: BulkUpdateMode;
  onValueChange: (v: unknown) => void;
  onModeChange: (m: BulkUpdateMode) => void;
}) {
  const [searchText, setSearchText] = useState("");

  const { data: entities } = useQuery({
    queryKey: [entityType, "all"],
    queryFn: async () => {
      switch (entityType) {
        case "tags": return (await tagsApi.find({ perPage: 1000, sort: "name", direction: "asc" })).items;
        case "performers": return (await performersApi.find({ perPage: 1000, sort: "name", direction: "asc" })).items;
        case "studios": return (await studiosApi.find({ perPage: 1000, sort: "name", direction: "asc" })).items;
        case "groups": return (await groupsApi.find({ perPage: 1000, sort: "name", direction: "asc" })).items;
        default: return [];
      }
    },
    staleTime: 60000,
  });

  const getName = (e: any) => e.name || e.title || `#${e.id}`;

  const filteredEntities = useMemo(() => {
    if (!entities) return [];
    const q = searchText.toLowerCase();
    return q ? entities.filter((e: any) => getName(e).toLowerCase().includes(q)) : entities;
  }, [entities, searchText]);

  const toggleId = (id: number) => {
    const next = value.includes(id) ? value.filter((i) => i !== id) : [...value, id];
    onValueChange(next);
  };

  return (
    <div className="space-y-2">
      {/* Mode selector */}
      <div className="flex gap-1">
        {(["SET", "ADD", "REMOVE"] as BulkUpdateMode[]).map((m) => (
          <button
            key={m}
            onClick={() => onModeChange(m)}
            className={`px-2 py-0.5 rounded text-[10px] border ${
              m === mode ? "bg-plex-accent text-white border-plex-accent" : "border-plex-border text-plex-text-secondary"
            }`}
          >
            {m}
          </button>
        ))}
      </div>

      {/* Selected */}
      {value.length > 0 && (
        <div className="flex flex-wrap gap-1">
          {value.map((id) => {
            const entity = entities?.find((e: any) => e.id === id);
            return (
              <span key={id} className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] bg-plex-card text-plex-text border border-plex-border">
                {entity ? getName(entity) : `#${id}`}
                <button onClick={() => toggleId(id)} className="hover:text-red-400">
                  <X className="w-2.5 h-2.5" />
                </button>
              </span>
            );
          })}
        </div>
      )}

      {/* Search */}
      <input
        type="text"
        value={searchText}
        onChange={(e) => setSearchText(e.target.value)}
        placeholder={`Search ${entityType}...`}
        className="w-full bg-plex-input border border-plex-border rounded px-2 py-1 text-xs text-plex-text focus:outline-none focus:border-plex-accent placeholder:text-plex-text-muted"
      />
      <div className="max-h-32 overflow-y-auto border border-plex-border rounded bg-plex-input">
        {filteredEntities.slice(0, 50).map((entity: any) => {
          const isSelected = value.includes(entity.id);
          return (
            <button
              key={entity.id}
              onClick={() => toggleId(entity.id)}
              className={`w-full text-left px-2 py-1 text-xs hover:bg-plex-card flex items-center gap-1 ${isSelected ? "text-plex-accent" : "text-plex-text"}`}
            >
              {isSelected ? <Minus className="w-3 h-3" /> : <Plus className="w-3 h-3" />}
              {getName(entity)}
            </button>
          );
        })}
      </div>
    </div>
  );
}

// ===== Pre-configured bulk edit field sets =====

export const SCENE_BULK_FIELDS: BulkEditField[] = [
  { key: "rating", label: "Rating", type: "number" },
  { key: "organized", label: "Organized", type: "bool" },
  { key: "studioId", label: "Studio", type: "select", nullable: true },
  { key: "date", label: "Date", type: "date" },
  { key: "tagIds", label: "Tags", type: "multiId", entityType: "tags" },
  { key: "performerIds", label: "Performers", type: "multiId", entityType: "performers" },
  { key: "groupIds", label: "Groups", type: "multiId", entityType: "groups" },
];

export const PERFORMER_BULK_FIELDS: BulkEditField[] = [
  { key: "rating", label: "Rating", type: "number" },
  { key: "favorite", label: "Favorite", type: "bool" },
  { key: "tagIds", label: "Tags", type: "multiId", entityType: "tags" },
];

export const GALLERY_BULK_FIELDS: BulkEditField[] = [
  { key: "rating", label: "Rating", type: "number" },
  { key: "organized", label: "Organized", type: "bool" },
  { key: "studioId", label: "Studio", type: "select", nullable: true },
  { key: "tagIds", label: "Tags", type: "multiId", entityType: "tags" },
  { key: "performerIds", label: "Performers", type: "multiId", entityType: "performers" },
];

export const IMAGE_BULK_FIELDS: BulkEditField[] = [
  { key: "rating", label: "Rating", type: "number" },
  { key: "organized", label: "Organized", type: "bool" },
  { key: "studioId", label: "Studio", type: "select", nullable: true },
  { key: "tagIds", label: "Tags", type: "multiId", entityType: "tags" },
  { key: "performerIds", label: "Performers", type: "multiId", entityType: "performers" },
];
