import { useState, useMemo, useCallback } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { scenes } from "../api/client";
import type { Scene, FindFilter } from "../api/types";
import { Search, ChevronDown, Loader2, Check } from "lucide-react";

interface Props {
  onNavigate: (r: any) => void;
}

// ===== Token / Regex helpers =====

const RECIPES = [
  "{title}.{ext}",
  "{performer} - {title}.{ext}",
  "{studio} - {title}.{ext}",
  "{date}.{title}.{ext}",
  "{performer} - {date} - {title}.{ext}",
  "{studio}.{date}.{performer}.{title}.{ext}",
];

type TokenKind = "title" | "date" | "performer" | "tag" | "studio" | "ext" | "ignore";

interface PatternToken {
  kind: TokenKind;
  groupIndex: number; // -1 for non-capturing
}

const TOKEN_REGEX_MAP: Record<TokenKind, { regex: string; capturing: boolean }> = {
  title: { regex: "(.+?)", capturing: true },
  date: { regex: String.raw`(\d{4}[\.\-]\d{2}[\.\-]\d{2}|\d{2}[\.\-]\d{2}[\.\-]\d{4}|\d{6,8})`, capturing: true },
  performer: { regex: "(.+?)", capturing: true },
  tag: { regex: "(.+?)", capturing: true },
  studio: { regex: "(.+?)", capturing: true },
  ext: { regex: String.raw`\w+`, capturing: false },
  ignore: { regex: ".+?", capturing: false },
};

interface ParseResult {
  title?: string;
  date?: string;
  performer?: string;
  tag?: string;
  studio?: string;
}

function compilePattern(pattern: string): { regex: RegExp; tokens: PatternToken[] } | null {
  if (!pattern.trim()) return null;

  const tokens: PatternToken[] = [];
  let regexStr = "^";
  let groupIdx = 0;

  // Split pattern into segments of tokens and literal separators
  const parts = pattern.split(/(\{[^}]*\})/);

  for (const part of parts) {
    if (!part) continue;

    const tokenMatch = part.match(/^\{(\w*)\}$/);
    if (tokenMatch) {
      const raw = tokenMatch[1].toLowerCase();
      let kind: TokenKind;
      if (raw === "" || raw === "ignore") {
        kind = "ignore";
      } else if (raw in TOKEN_REGEX_MAP) {
        kind = raw as TokenKind;
      } else {
        kind = "ignore";
      }
      const spec = TOKEN_REGEX_MAP[kind];
      regexStr += spec.regex;
      if (spec.capturing) {
        groupIdx++;
        tokens.push({ kind, groupIndex: groupIdx });
      } else {
        tokens.push({ kind, groupIndex: -1 });
      }
    } else {
      // Literal text – treat separator chars as flexible separators
      const escaped = part.replace(/[\.\-_ ]+/g, String.raw`[\.\-_ ]+`);
      regexStr += escaped;
    }
  }

  regexStr += "$";

  try {
    return { regex: new RegExp(regexStr, "i"), tokens };
  } catch {
    return null;
  }
}

function applyPattern(
  basename: string,
  compiled: { regex: RegExp; tokens: PatternToken[] },
  ignoredWords: string[],
  capitalizeTitle: boolean
): ParseResult | null {
  const match = basename.match(compiled.regex);
  if (!match) return null;

  const result: ParseResult = {};

  for (const token of compiled.tokens) {
    if (token.groupIndex < 0) continue;
    let value = match[token.groupIndex];
    if (value === undefined) continue;

    // Clean up separators for readability
    value = value.replace(/[\.\-_]+/g, " ").trim();

    if (token.kind === "title") {
      // Strip ignored words
      if (ignoredWords.length > 0) {
        const lowerIgnored = ignoredWords.map((w) => w.toLowerCase());
        value = value
          .split(/\s+/)
          .filter((w) => !lowerIgnored.includes(w.toLowerCase()))
          .join(" ");
      }
      if (capitalizeTitle) {
        value = value
          .split(/\s+/)
          .map((w) => (w.length > 0 ? w[0].toUpperCase() + w.slice(1).toLowerCase() : w))
          .join(" ");
      }
      result.title = value;
    } else if (token.kind === "date") {
      // Normalise date separators to dashes
      result.date = value.replace(/[\.\- ]+/g, "-");
    } else if (token.kind === "performer") {
      result.performer = value;
    } else if (token.kind === "tag") {
      result.tag = value;
    } else if (token.kind === "studio") {
      result.studio = value;
    }
  }

  return result;
}

// ===== Per-row edit state =====

interface RowState {
  sceneId: number;
  basename: string;
  parsed: ParseResult | null;
  editedTitle?: string;
  editedDate?: string;
  selected: boolean;
}

// ===== Component =====

export function SceneFilenameParserPage({ onNavigate }: Props) {
  // Pattern config
  const [pattern, setPattern] = useState("{title}.{ext}");
  const [ignoredWordsStr, setIgnoredWordsStr] = useState("");
  const [capitalizeTitle, setCapitalizeTitle] = useState(false);
  const [ignoreOrganized, setIgnoreOrganized] = useState(false);
  const [perPage, setPerPage] = useState(40);
  const [showRecipes, setShowRecipes] = useState(false);

  // Query trigger
  const [queryEnabled, setQueryEnabled] = useState(false);
  const [appliedConfig, setAppliedConfig] = useState<{
    pattern: string;
    ignoredWords: string[];
    capitalizeTitle: boolean;
    ignoreOrganized: boolean;
    perPage: number;
  } | null>(null);

  // Rows
  const [rows, setRows] = useState<RowState[]>([]);

  const queryClient = useQueryClient();

  const filter: FindFilter = useMemo(
    () => ({ page: 1, perPage: appliedConfig?.perPage ?? perPage }),
    [appliedConfig?.perPage, perPage]
  );

  const { data, isLoading, isFetching } = useQuery({
    queryKey: ["scenes-parser", filter, appliedConfig],
    queryFn: () => scenes.find(filter),
    enabled: queryEnabled && appliedConfig !== null,
  });

  // Reparse whenever data or config changes
  useMemo(() => {
    if (!data?.items || !appliedConfig) return;

    const compiled = compilePattern(appliedConfig.pattern);
    if (!compiled) {
      setRows([]);
      return;
    }

    const items = appliedConfig.ignoreOrganized
      ? data.items.filter((s) => !s.organized)
      : data.items;

    const newRows: RowState[] = items.map((scene) => {
      const file = scene.files[0];
      const basename = file?.basename ?? "";
      const parsed = compiled ? applyPattern(basename, compiled, appliedConfig.ignoredWords, appliedConfig.capitalizeTitle) : null;
      return {
        sceneId: scene.id,
        basename,
        parsed,
        selected: false,
      };
    });

    setRows(newRows);
  }, [data, appliedConfig]);

  // ===== Actions =====

  const handleFind = useCallback(() => {
    const words = ignoredWordsStr
      .split(/\s+/)
      .map((w) => w.trim())
      .filter(Boolean);
    setAppliedConfig({
      pattern,
      ignoredWords: words,
      capitalizeTitle,
      ignoreOrganized,
      perPage,
    });
    setQueryEnabled(true);
  }, [pattern, ignoredWordsStr, capitalizeTitle, ignoreOrganized, perPage]);

  const handleSelectAll = useCallback(
    (checked: boolean) => {
      setRows((prev) => prev.map((r) => ({ ...r, selected: checked })));
    },
    []
  );

  const handleToggle = useCallback((idx: number) => {
    setRows((prev) => prev.map((r, i) => (i === idx ? { ...r, selected: !r.selected } : r)));
  }, []);

  const handleEditTitle = useCallback((idx: number, value: string) => {
    setRows((prev) =>
      prev.map((r, i) => (i === idx ? { ...r, editedTitle: value } : r))
    );
  }, []);

  const handleEditDate = useCallback((idx: number, value: string) => {
    setRows((prev) =>
      prev.map((r, i) => (i === idx ? { ...r, editedDate: value } : r))
    );
  }, []);

  const selectedRows = useMemo(() => rows.filter((r) => r.selected && r.parsed), [rows]);
  const allSelected = rows.length > 0 && rows.every((r) => r.selected);

  // Apply mutation
  const applyMut = useMutation({
    mutationFn: async () => {
      const updates = selectedRows.map((r) => {
        const title = r.editedTitle ?? r.parsed?.title;
        const date = r.editedDate ?? r.parsed?.date;
        return scenes.update(r.sceneId, {
          ...(title !== undefined ? { title } : {}),
          ...(date !== undefined ? { date } : {}),
        });
      });
      await Promise.all(updates);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["scenes"] });
      queryClient.invalidateQueries({ queryKey: ["scenes-parser"] });
    },
  });

  // ===== Render helpers =====

  const displayTitle = (r: RowState) => r.editedTitle ?? r.parsed?.title ?? "";
  const displayDate = (r: RowState) => r.editedDate ?? r.parsed?.date ?? "";
  const isTitleEdited = (r: RowState) => r.editedTitle !== undefined && r.editedTitle !== (r.parsed?.title ?? "");
  const isDateEdited = (r: RowState) => r.editedDate !== undefined && r.editedDate !== (r.parsed?.date ?? "");

  return (
    <div className="min-h-screen bg-gray-900 text-gray-200 p-6">
      <h1 className="text-2xl font-bold mb-6">Scene Filename Parser</h1>

      {/* ===== Pattern Input Section ===== */}
      <div className="bg-gray-800 rounded-lg p-5 mb-6 space-y-4">
        {/* Pattern row */}
        <div className="flex flex-wrap items-end gap-4">
          <div className="flex-1 min-w-[300px]">
            <label className="block text-sm font-medium text-gray-400 mb-1">
              Filename Pattern
            </label>
            <div className="flex gap-2">
              <input
                type="text"
                value={pattern}
                onChange={(e) => setPattern(e.target.value)}
                placeholder="{title}.{ext}"
                className="flex-1 bg-gray-700 border border-gray-600 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
              />
              {/* Recipe dropdown */}
              <div className="relative">
                <button
                  type="button"
                  onClick={() => setShowRecipes(!showRecipes)}
                  className="flex items-center gap-1 bg-gray-700 border border-gray-600 rounded px-3 py-2 text-sm text-gray-300 hover:bg-gray-600"
                >
                  Recipes
                  <ChevronDown className="w-4 h-4" />
                </button>
                {showRecipes && (
                  <div className="absolute right-0 mt-1 w-80 bg-gray-700 border border-gray-600 rounded shadow-lg z-10">
                    {RECIPES.map((r) => (
                      <button
                        key={r}
                        type="button"
                        onClick={() => {
                          setPattern(r);
                          setShowRecipes(false);
                        }}
                        className="block w-full text-left px-3 py-2 text-sm text-gray-200 hover:bg-gray-600 font-mono"
                      >
                        {r}
                      </button>
                    ))}
                  </div>
                )}
              </div>
            </div>
            <p className="text-xs text-gray-500 mt-1">
              Tokens: <code className="text-gray-400">{"{title}"}</code>{" "}
              <code className="text-gray-400">{"{date}"}</code>{" "}
              <code className="text-gray-400">{"{performer}"}</code>{" "}
              <code className="text-gray-400">{"{tag}"}</code>{" "}
              <code className="text-gray-400">{"{studio}"}</code>{" "}
              <code className="text-gray-400">{"{ext}"}</code>{" "}
              <code className="text-gray-400">{"{}"}</code> (ignored)
            </p>
          </div>
        </div>

        {/* Options row */}
        <div className="flex flex-wrap items-end gap-4">
          <div className="min-w-[200px]">
            <label className="block text-sm font-medium text-gray-400 mb-1">
              Ignored Words
            </label>
            <input
              type="text"
              value={ignoredWordsStr}
              onChange={(e) => setIgnoredWordsStr(e.target.value)}
              placeholder="e.g. 720p 1080p bluray"
              className="w-full bg-gray-700 border border-gray-600 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
            />
          </div>

          <label className="flex items-center gap-2 cursor-pointer select-none">
            <input
              type="checkbox"
              checked={capitalizeTitle}
              onChange={(e) => setCapitalizeTitle(e.target.checked)}
              className="rounded bg-gray-700 border-gray-600 text-blue-500 focus:ring-blue-500"
            />
            <span className="text-sm text-gray-300">Capitalize Title</span>
          </label>

          <label className="flex items-center gap-2 cursor-pointer select-none">
            <input
              type="checkbox"
              checked={ignoreOrganized}
              onChange={(e) => setIgnoreOrganized(e.target.checked)}
              className="rounded bg-gray-700 border-gray-600 text-blue-500 focus:ring-blue-500"
            />
            <span className="text-sm text-gray-300">Ignore Organized</span>
          </label>

          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">
              Page Size
            </label>
            <select
              value={perPage}
              onChange={(e) => setPerPage(Number(e.target.value))}
              className="bg-gray-700 border border-gray-600 rounded px-3 py-2 text-sm text-gray-200 focus:outline-none focus:border-blue-500"
            >
              <option value={20}>20</option>
              <option value={40}>40</option>
              <option value={100}>100</option>
            </select>
          </div>

          <button
            type="button"
            onClick={handleFind}
            disabled={!pattern.trim()}
            className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded px-4 py-2 text-sm font-medium"
          >
            <Search className="w-4 h-4" />
            Find
          </button>
        </div>
      </div>

      {/* ===== Loading ===== */}
      {(isLoading || isFetching) && queryEnabled && (
        <div className="flex items-center gap-2 text-gray-400 mb-4">
          <Loader2 className="w-5 h-5 animate-spin" />
          Loading scenes…
        </div>
      )}

      {/* ===== Results Table ===== */}
      {rows.length > 0 && (
        <>
          <div className="bg-gray-800 rounded-lg overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-700 text-left text-gray-400">
                    <th className="px-3 py-3 w-10">
                      <input
                        type="checkbox"
                        checked={allSelected}
                        onChange={(e) => handleSelectAll(e.target.checked)}
                        className="rounded bg-gray-700 border-gray-600 text-blue-500 focus:ring-blue-500"
                      />
                    </th>
                    <th className="px-3 py-3">Filename</th>
                    <th className="px-3 py-3">Title</th>
                    <th className="px-3 py-3 w-40">Date</th>
                    <th className="px-3 py-3">Performers</th>
                    <th className="px-3 py-3">Tags</th>
                    <th className="px-3 py-3">Studio</th>
                  </tr>
                </thead>
                <tbody>
                  {rows.map((row, idx) => (
                    <tr
                      key={row.sceneId}
                      className={`border-b border-gray-700/50 hover:bg-gray-750 ${
                        !row.parsed ? "opacity-50" : ""
                      }`}
                    >
                      <td className="px-3 py-2">
                        <input
                          type="checkbox"
                          checked={row.selected}
                          onChange={() => handleToggle(idx)}
                          disabled={!row.parsed}
                          className="rounded bg-gray-700 border-gray-600 text-blue-500 focus:ring-blue-500 disabled:opacity-30"
                        />
                      </td>
                      <td className="px-3 py-2 font-mono text-xs text-gray-400 max-w-[300px] truncate" title={row.basename}>
                        {row.basename}
                      </td>
                      <td className="px-3 py-2">
                        {row.parsed ? (
                          <input
                            type="text"
                            value={displayTitle(row)}
                            onChange={(e) => handleEditTitle(idx, e.target.value)}
                            className={`w-full bg-gray-700 border rounded px-2 py-1 text-sm text-gray-200 focus:outline-none focus:border-blue-500 ${
                              isTitleEdited(row) ? "border-yellow-500" : "border-gray-600"
                            }`}
                          />
                        ) : (
                          <span className="text-gray-600 italic">no match</span>
                        )}
                      </td>
                      <td className="px-3 py-2">
                        {row.parsed?.date !== undefined ? (
                          <input
                            type="text"
                            value={displayDate(row)}
                            onChange={(e) => handleEditDate(idx, e.target.value)}
                            placeholder="YYYY-MM-DD"
                            className={`w-full bg-gray-700 border rounded px-2 py-1 text-sm text-gray-200 focus:outline-none focus:border-blue-500 ${
                              isDateEdited(row) ? "border-yellow-500" : "border-gray-600"
                            }`}
                          />
                        ) : (
                          <span className="text-gray-600">—</span>
                        )}
                      </td>
                      <td className="px-3 py-2 text-gray-300">
                        {row.parsed?.performer ?? <span className="text-gray-600">—</span>}
                      </td>
                      <td className="px-3 py-2 text-gray-300">
                        {row.parsed?.tag ?? <span className="text-gray-600">—</span>}
                      </td>
                      <td className="px-3 py-2 text-gray-300">
                        {row.parsed?.studio ?? <span className="text-gray-600">—</span>}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Footer */}
          <div className="flex items-center justify-between mt-4">
            <span className="text-sm text-gray-400">
              {selectedRows.length} of {rows.length} selected
              {rows.filter((r) => !r.parsed).length > 0 && (
                <> · {rows.filter((r) => !r.parsed).length} unmatched</>
              )}
            </span>
            <button
              type="button"
              onClick={() => applyMut.mutate()}
              disabled={selectedRows.length === 0 || applyMut.isPending}
              className="flex items-center gap-2 bg-green-600 hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded px-5 py-2 text-sm font-medium"
            >
              {applyMut.isPending ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <Check className="w-4 h-4" />
              )}
              Apply
            </button>
          </div>

          {applyMut.isSuccess && (
            <div className="mt-3 text-sm text-green-400">
              Successfully updated {selectedRows.length} scene(s).
            </div>
          )}
          {applyMut.isError && (
            <div className="mt-3 text-sm text-red-400">
              Error applying changes: {(applyMut.error as Error).message}
            </div>
          )}
        </>
      )}

      {/* Empty state after search */}
      {queryEnabled && !isLoading && !isFetching && rows.length === 0 && appliedConfig && (
        <div className="text-center text-gray-500 py-12">
          No scenes matched the pattern. Try a different pattern or check that scenes have files.
        </div>
      )}
    </div>
  );
}
