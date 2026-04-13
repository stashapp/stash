import { useEffect, useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  ChevronDown,
  ChevronUp,
  Database,
  Download,
  FolderOpen,
  HardDrive,
  Info,
  Loader2,
  Monitor,
  Plug,
  Plus,
  RefreshCw,
  Save,
  SearchCode,
  Server,
  Shield,
  Trash2,
  PlayCircle,
  Radio,
  ScrollText,
  Upload,
  Wrench,
} from "lucide-react";
import { system, jobs, metadata, database, plugins as pluginsApi, dlna as dlnaApi, logs as logsApi } from "../api/client";
import type { ScanOptions, GenerateOptions, CleanGeneratedOptions, ExportOptions, LogEntry } from "../api/client";
import type {
  DlnaStatus,
  PackageSource,
  Plugin,
  RatingStarPrecision,
  RatingSystemType,
  ScraperSummary,
  StashBox,
  StashConfig,
  StashPathConfig,
  StashBoxValidationResult,
} from "../api/types";
import { useAppConfig } from "../state/AppConfigContext";

type SettingsTab = "tasks" | "library" | "interface" | "security" | "metadata-providers" | "dlna" | "plugins" | "logs" | "system" | "about";

const tabs: { key: SettingsTab; label: string; icon: typeof FolderOpen }[] = [
  { key: "tasks", label: "Tasks", icon: PlayCircle },
  { key: "library", label: "Library", icon: FolderOpen },
  { key: "interface", label: "Interface", icon: Monitor },
  { key: "security", label: "Security", icon: Shield },
  { key: "metadata-providers", label: "Metadata Providers", icon: SearchCode },
  { key: "dlna", label: "Services (DLNA)", icon: Radio },
  { key: "plugins", label: "Plugins", icon: Plug },
  { key: "logs", label: "Logs", icon: ScrollText },
  { key: "system", label: "System", icon: Server },
  { key: "about", label: "About", icon: Info },
];

const languageOptions = [
  { value: "en-US", label: "English (United States)" },
  { value: "en-GB", label: "English (United Kingdom)" },
  { value: "de-DE", label: "Deutsch" },
  { value: "fr-FR", label: "Francais" },
  { value: "es-ES", label: "Espanol" },
  { value: "it-IT", label: "Italiano" },
  { value: "ja-JP", label: "Japanese" },
  { value: "ko-KR", label: "Korean" },
  { value: "nl-NL", label: "Nederlands" },
  { value: "pl-PL", label: "Polski" },
  { value: "pt-BR", label: "Portugues (Brasil)" },
  { value: "ru-RU", label: "Russian" },
  { value: "sv-SE", label: "Svenska" },
  { value: "zh-CN", label: "Chinese (Simplified)" },
  { value: "zh-TW", label: "Chinese (Traditional)" },
];

const menuItems = [
  { value: "scenes", label: "Scenes" },
  { value: "images", label: "Images" },
  { value: "performers", label: "Performers" },
  { value: "galleries", label: "Galleries" },
  { value: "studios", label: "Studios" },
  { value: "tags", label: "Tags" },
  { value: "groups", label: "Groups" },
];

const ratingSystemOptions: { value: RatingSystemType; label: string }[] = [
  { value: "stars", label: "Stars" },
  { value: "decimal", label: "Decimal (0-10.0)" },
];

const starPrecisionOptions: { value: RatingStarPrecision; label: string }[] = [
  { value: "full", label: "Full stars" },
  { value: "half", label: "Half stars" },
  { value: "quarter", label: "Quarter stars" },
  { value: "tenth", label: "Tenth stars" },
];

function emptyPath(): StashPathConfig {
  return { path: "", excludeVideo: false, excludeImage: false };
}

function emptyPackageSource(): PackageSource {
  return { name: "", url: "" };
}

function emptyStashBox(): StashBox {
  return { name: "", endpoint: "", apiKey: "", maxRequestsPerMinute: 240 };
}

function cloneConfig(config: StashConfig): StashConfig {
  return JSON.parse(JSON.stringify(config)) as StashConfig;
}

function linesToList(value: string) {
  return value
    .split(/\r?\n/)
    .map((item) => item.trim())
    .filter(Boolean);
}

function listToLines(values: string[]) {
  return values.join("\n");
}

function normalizeConfig(config: StashConfig): StashConfig {
  return {
    ...config,
    stashPaths: config.stashPaths.filter((path) => path.path.trim() !== ""),
    videoExtensions: config.videoExtensions.map((value) => value.trim()).filter(Boolean),
    imageExtensions: config.imageExtensions.map((value) => value.trim()).filter(Boolean),
    galleryExtensions: config.galleryExtensions.map((value) => value.trim()).filter(Boolean),
    excludePatterns: config.excludePatterns.map((value) => value.trim()).filter(Boolean),
    interface: {
      ...config.interface,
      menuItems: config.interface.menuItems.filter(Boolean),
    },
    security: {
      ...config.security,
      username: config.security.username?.trim() || undefined,
      newPassword: config.security.newPassword?.trim() || undefined,
    },
    scraping: {
      scraperDirectories: config.scraping.scraperDirectories.map((value) => value.trim()).filter(Boolean),
      scraperPackageSources: config.scraping.scraperPackageSources
        .map((source) => ({ name: source.name.trim(), url: source.url.trim() }))
        .filter((source) => source.url !== ""),
      stashBoxes: config.scraping.stashBoxes
        .map((box) => ({
          name: box.name.trim(),
          endpoint: box.endpoint.trim(),
          apiKey: box.apiKey.trim(),
          maxRequestsPerMinute: box.maxRequestsPerMinute,
        }))
        .filter((box) => box.endpoint !== ""),
    },
  };
}

export function SettingsPage() {
  const { config, status, configLoading, statusLoading } = useAppConfig();
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<SettingsTab>("library");
  const [draft, setDraft] = useState<StashConfig | null>(null);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [stashBoxValidation, setStashBoxValidation] = useState<Record<string, StashBoxValidationResult>>({});

  useEffect(() => {
    if (!config) {
      return;
    }

    const nextDraft = cloneConfig(config);
    if (nextDraft.stashPaths.length === 0) {
      nextDraft.stashPaths = [emptyPath()];
    }
    if (nextDraft.scraping.scraperPackageSources.length === 0) {
      nextDraft.scraping.scraperPackageSources = [emptyPackageSource()];
    }
    if (nextDraft.scraping.scraperDirectories.length === 0) {
      nextDraft.scraping.scraperDirectories = [""];
    }

    setDraft(nextDraft);
  }, [config]);

  const saveMutation = useMutation({
    mutationFn: (nextConfig: StashConfig) => system.saveConfig(nextConfig),
    onSuccess: (savedConfig) => {
      queryClient.setQueryData(["system-config"], savedConfig);
      queryClient.invalidateQueries({ queryKey: ["system-scrapers"] });
      setSaved(true);
      setError(null);
      setTimeout(() => setSaved(false), 2000);
    },
    onError: (err: Error) => setError(err.message),
  });

  const { data: scrapers = [], isLoading: scrapersLoading } = useQuery({
    queryKey: ["system-scrapers"],
    queryFn: system.listScrapers,
    enabled: activeTab === "metadata-providers",
  });

  const reloadScrapersMutation = useMutation({
    mutationFn: system.reloadScrapers,
    onSuccess: (nextScrapers) => {
      queryClient.setQueryData(["system-scrapers"], nextScrapers);
    },
  });

  const validateStashBoxMutation = useMutation({
    mutationFn: ({ index, stashBox }: { index: number; stashBox: StashBox }) => system.validateStashBox(stashBox),
    onSuccess: (result, variables) => {
      setStashBoxValidation((current) => ({ ...current, [String(variables.index)]: result }));
    },
    onError: (err: Error, variables) => {
      setStashBoxValidation((current) => ({
        ...current,
        [String(variables.index)]: { valid: false, status: err.message },
      }));
    },
  });

  const groupedScrapers = useMemo(() => {
    return scrapers.reduce<Record<string, ScraperSummary[]>>((acc, scraper) => {
      if (!acc[scraper.entityType]) {
        acc[scraper.entityType] = [];
      }
      acc[scraper.entityType].push(scraper);
      return acc;
    }, {});
  }, [scrapers]);

  if (configLoading || !draft) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-6 w-6 animate-spin text-plex-text-muted" />
      </div>
    );
  }

  const updateDraft = (updater: (current: StashConfig) => StashConfig) => {
    setDraft((current) => (current ? updater(current) : current));
  };

  const handleSave = () => {
    saveMutation.mutate(normalizeConfig(draft));
  };

  return (
    <div className="grid gap-6 lg:grid-cols-[240px_minmax(0,1fr)]">
      <aside className="h-fit rounded-2xl border border-plex-border bg-plex-surface p-2 lg:sticky lg:top-16">
        <div className="mb-2 px-3 py-2">
          <h1 className="text-lg font-semibold text-plex-text">Settings</h1>
          <p className="mt-1 text-sm text-plex-text-secondary">Stock Stash-style categories, backed by the rewrite config.</p>
        </div>
        <nav className="space-y-1">
          {tabs.map(({ key, label, icon: Icon }) => (
            <button
              key={key}
              onClick={() => setActiveTab(key)}
              className={`flex w-full items-center gap-2 rounded-xl px-3 py-2 text-left text-sm transition-colors ${
                activeTab === key
                  ? "bg-plex-card text-plex-text shadow-[inset_0_0_0_1px_var(--color-plex-border)]"
                  : "text-plex-text-secondary hover:bg-plex-card hover:text-plex-text"
              }`}
            >
              <Icon className="h-4 w-4" />
              <span>{label}</span>
            </button>
          ))}
        </nav>
      </aside>

      <div className="space-y-5">
        <section className="rounded-2xl border border-plex-border bg-[linear-gradient(180deg,rgba(55,60,63,0.92),rgba(39,43,46,0.98))] p-5 shadow-lg shadow-black/20">
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div>
              <h2 className="text-xl font-semibold text-plex-text">{tabs.find((tab) => tab.key === activeTab)?.label}</h2>
              <p className="mt-1 text-sm text-plex-text-secondary">
                {activeTab === "tasks" && "Scan, generate, and maintenance operations."}
                {activeTab === "library" && "Content locations, generated assets, and scan rules."}
                {activeTab === "interface" && "Language, custom title, navigation, and rating presentation."}
                {activeTab === "security" && "Authentication and session settings. Password changes are persisted immediately."}
                {activeTab === "metadata-providers" && "Scraper directories, package source URLs, configured StashBox endpoints, and discovered Stash-compatible scrapers."}
                {activeTab === "dlna" && "DLNA media server for streaming to compatible devices on your local network."}
                {activeTab === "plugins" && "Manage installed plugins and extensions."}
                {activeTab === "system" && "Host, port, and task concurrency. Server changes take effect after restart."}
                {activeTab === "about" && "Runtime status and effective config locations."}
              </p>
            </div>
            <div className="flex flex-wrap items-center gap-3">
              {error && <span className="text-sm text-red-300">{error}</span>}
              {saved && <span className="text-sm text-emerald-300">Settings saved.</span>}
              <button
                onClick={handleSave}
                disabled={saveMutation.isPending}
                className="inline-flex items-center gap-2 rounded-xl bg-plex-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-plex-accent-hover disabled:cursor-not-allowed disabled:opacity-60"
              >
                {saveMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
                {saveMutation.isPending ? "Saving..." : "Save Settings"}
              </button>
            </div>
          </div>
        </section>

        {activeTab === "tasks" && <TasksPanel />}

        {activeTab === "library" && (
          <>
            <SectionCard title="Library Paths" description="Add the content roots the scanner should process.">
              <div className="space-y-3">
                {draft.stashPaths.map((path, index) => (
                  <div key={index} className="rounded-xl border border-plex-border bg-plex-card p-3">
                    <div className="flex flex-col gap-3 xl:flex-row xl:items-center">
                      <input
                        type="text"
                        value={path.path}
                        onChange={(event) =>
                          updateDraft((current) => ({
                            ...current,
                            stashPaths: current.stashPaths.map((item, itemIndex) =>
                              itemIndex === index ? { ...item, path: event.target.value } : item,
                            ),
                          }))
                        }
                        placeholder="D:\\Media\\Scenes"
                        className="flex-1 rounded-xl border border-plex-border bg-plex-surface px-3 py-2 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
                      />
                      <div className="flex flex-wrap items-center gap-4">
                        <CheckboxLabel
                          label="Exclude videos"
                          checked={path.excludeVideo}
                          onChange={(checked) =>
                            updateDraft((current) => ({
                              ...current,
                              stashPaths: current.stashPaths.map((item, itemIndex) =>
                                itemIndex === index ? { ...item, excludeVideo: checked } : item,
                              ),
                            }))
                          }
                        />
                        <CheckboxLabel
                          label="Exclude images"
                          checked={path.excludeImage}
                          onChange={(checked) =>
                            updateDraft((current) => ({
                              ...current,
                              stashPaths: current.stashPaths.map((item, itemIndex) =>
                                itemIndex === index ? { ...item, excludeImage: checked } : item,
                              ),
                            }))
                          }
                        />
                        <button
                          onClick={() =>
                            updateDraft((current) => ({
                              ...current,
                              stashPaths:
                                current.stashPaths.length > 1
                                  ? current.stashPaths.filter((_, itemIndex) => itemIndex !== index)
                                  : [emptyPath()],
                            }))
                          }
                          className="inline-flex items-center gap-1 rounded-lg border border-plex-border px-2 py-1 text-xs text-red-300 hover:border-red-500 hover:text-red-200"
                        >
                          <Trash2 className="h-3.5 w-3.5" /> Remove
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
                <button
                  onClick={() => updateDraft((current) => ({ ...current, stashPaths: [...current.stashPaths, emptyPath()] }))}
                  className="inline-flex items-center gap-2 rounded-xl border border-dashed border-plex-border px-3 py-2 text-sm text-plex-text-secondary hover:text-plex-text"
                >
                  <Plus className="h-4 w-4" /> Add path
                </button>
              </div>
            </SectionCard>

            <SectionCard title="Generated Assets" description="Control where generated and cached media artifacts are written.">
              <div className="grid gap-4 md:grid-cols-2">
                <TextField
                  label="Generated path"
                  value={draft.generatedPath ?? ""}
                  onChange={(value) => updateDraft((current) => ({ ...current, generatedPath: value || undefined }))}
                  placeholder="D:\\Stash\\generated"
                />
                <TextField
                  label="Cache path"
                  value={draft.cachePath ?? ""}
                  onChange={(value) => updateDraft((current) => ({ ...current, cachePath: value || undefined }))}
                  placeholder="D:\\Stash\\cache"
                />
              </div>
            </SectionCard>

            <SectionCard title="Extensions" description="One extension per line. These values are persisted directly into the backend config.">
              <div className="grid gap-4 lg:grid-cols-3">
                <TextAreaField
                  label="Video extensions"
                  value={listToLines(draft.videoExtensions)}
                  onChange={(value) => updateDraft((current) => ({ ...current, videoExtensions: linesToList(value) }))}
                  rows={7}
                />
                <TextAreaField
                  label="Image extensions"
                  value={listToLines(draft.imageExtensions)}
                  onChange={(value) => updateDraft((current) => ({ ...current, imageExtensions: linesToList(value) }))}
                  rows={7}
                />
                <TextAreaField
                  label="Gallery extensions"
                  value={listToLines(draft.galleryExtensions)}
                  onChange={(value) => updateDraft((current) => ({ ...current, galleryExtensions: linesToList(value) }))}
                  rows={7}
                />
              </div>
            </SectionCard>

            <SectionCard title="Scan Rules" description="Hashing and exclude patterns applied during scan operations.">
              <div className="space-y-4">
                <CheckboxLabel
                  label="Calculate MD5 checksums during scan"
                  checked={draft.calculateMd5}
                  onChange={(checked) => updateDraft((current) => ({ ...current, calculateMd5: checked }))}
                />
                <TextAreaField
                  label="Exclude patterns"
                  value={listToLines(draft.excludePatterns)}
                  onChange={(value) => updateDraft((current) => ({ ...current, excludePatterns: linesToList(value) }))}
                  rows={5}
                  placeholder="**/._*&#10;**/.DS_Store"
                />
              </div>
            </SectionCard>
          </>
        )}

        {activeTab === "interface" && (
          <>
            <SectionCard title="Basic Interface" description="Persisted UI preferences used across the app shell.">
              <div className="grid gap-4 md:grid-cols-2">
                <SelectField
                  label="Language"
                  value={draft.interface.language ?? "en-US"}
                  onChange={(value) => updateDraft((current) => ({ ...current, interface: { ...current.interface, language: value } }))}
                  options={languageOptions}
                />
                <TextField
                  label="Custom title"
                  value={draft.ui.title ?? ""}
                  onChange={(value) => updateDraft((current) => ({ ...current, ui: { ...current.ui, title: value || undefined } }))}
                  placeholder="Stash"
                />
              </div>
            </SectionCard>

            <SectionCard title="Navigation" description="These menu items are reflected in the rewrite navbar immediately after save.">
              <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
                {menuItems.map((item) => (
                  <CheckboxLabel
                    key={item.value}
                    label={item.label}
                    checked={draft.interface.menuItems.includes(item.value)}
                    onChange={(checked) =>
                      updateDraft((current) => ({
                        ...current,
                        interface: {
                          ...current.interface,
                          menuItems: checked
                            ? [...new Set([...current.interface.menuItems, item.value])]
                            : current.interface.menuItems.filter((value) => value !== item.value),
                        },
                      }))
                    }
                  />
                ))}
              </div>
            </SectionCard>

            <SectionCard title="Ratings" description="Stored ratings remain 1-100 internally. This changes how they are displayed and edited in the UI.">
              <div className="grid gap-4 md:grid-cols-2">
                <SelectField
                  label="Rating system"
                  value={draft.ui.ratingSystemOptions.type}
                  onChange={(value) =>
                    updateDraft((current) => ({
                      ...current,
                      ui: {
                        ...current.ui,
                        ratingSystemOptions: {
                          ...current.ui.ratingSystemOptions,
                          type: value as RatingSystemType,
                        },
                      },
                    }))
                  }
                  options={ratingSystemOptions}
                />
                {draft.ui.ratingSystemOptions.type === "stars" && (
                  <SelectField
                    label="Star precision"
                    value={draft.ui.ratingSystemOptions.starPrecision}
                    onChange={(value) =>
                      updateDraft((current) => ({
                        ...current,
                        ui: {
                          ...current.ui,
                          ratingSystemOptions: {
                            ...current.ui.ratingSystemOptions,
                            starPrecision: value as RatingStarPrecision,
                          },
                        },
                      }))
                    }
                    options={starPrecisionOptions}
                  />
                )}
              </div>
            </SectionCard>

            <SectionCard title="Scene Player" description="Playback behavior for the built-in video player.">
              <div className="space-y-3">
                <CheckboxLabel
                  label="Autostart video"
                  checked={draft.ui.autostartVideo}
                  onChange={(checked) => updateDraft((d) => ({ ...d, ui: { ...d.ui, autostartVideo: checked } }))}
                />
                <CheckboxLabel
                  label="Autostart video on play selected"
                  checked={draft.ui.autostartVideoOnPlaySelected}
                  onChange={(checked) => updateDraft((d) => ({ ...d, ui: { ...d.ui, autostartVideoOnPlaySelected: checked } }))}
                />
                <CheckboxLabel
                  label="Continue playlist default"
                  checked={draft.ui.continuePlaylistDefault}
                  onChange={(checked) => updateDraft((d) => ({ ...d, ui: { ...d.ui, continuePlaylistDefault: checked } }))}
                />
                <CheckboxLabel
                  label="Show A-B loop controls"
                  checked={draft.ui.showAbLoopControls}
                  onChange={(checked) => updateDraft((d) => ({ ...d, ui: { ...d.ui, showAbLoopControls: checked } }))}
                />
                <CheckboxLabel
                  label="Track activity"
                  checked={draft.ui.trackActivity}
                  onChange={(checked) => updateDraft((d) => ({ ...d, ui: { ...d.ui, trackActivity: checked } }))}
                />
              </div>
            </SectionCard>

            <SectionCard title="Preview" description="Preview generation and playback settings.">
              <div className="space-y-4">
                <CheckboxLabel
                  label="Sound on preview"
                  checked={draft.ui.soundOnPreview}
                  onChange={(checked) => updateDraft((d) => ({ ...d, ui: { ...d.ui, soundOnPreview: checked } }))}
                />
                <div className="grid gap-4 md:grid-cols-2">
                  <NumberField
                    label="Preview segment duration (seconds)"
                    value={draft.ui.previewSegmentDuration}
                    min={0}
                    onChange={(value) => updateDraft((d) => ({ ...d, ui: { ...d.ui, previewSegmentDuration: value ?? d.ui.previewSegmentDuration } }))}
                  />
                  <NumberField
                    label="Preview segments"
                    value={draft.ui.previewSegments}
                    min={0}
                    onChange={(value) => updateDraft((d) => ({ ...d, ui: { ...d.ui, previewSegments: value ?? d.ui.previewSegments } }))}
                  />
                  <TextField
                    label="Preview exclude start"
                    value={draft.ui.previewExcludeStart}
                    onChange={(value) => updateDraft((d) => ({ ...d, ui: { ...d.ui, previewExcludeStart: value } }))}
                  />
                  <TextField
                    label="Preview exclude end"
                    value={draft.ui.previewExcludeEnd}
                    onChange={(value) => updateDraft((d) => ({ ...d, ui: { ...d.ui, previewExcludeEnd: value } }))}
                  />
                </div>
              </div>
            </SectionCard>

            <SectionCard title="Wall" description="Wall view display options.">
              <div className="space-y-4">
                <CheckboxLabel
                  label="Wall show title"
                  checked={draft.ui.wallShowTitle}
                  onChange={(checked) => updateDraft((d) => ({ ...d, ui: { ...d.ui, wallShowTitle: checked } }))}
                />
                <SelectField
                  label="Wall playback"
                  value={String(draft.ui.wallPlayback)}
                  onChange={(value) => updateDraft((d) => ({ ...d, ui: { ...d.ui, wallPlayback: Number(value) } }))}
                  options={[
                    { value: "0", label: "Audio" },
                    { value: "1", label: "Silent" },
                  ]}
                />
              </div>
            </SectionCard>

            <SectionCard title="Lightbox" description="Lightbox and slideshow behavior.">
              <div className="space-y-4">
                <CheckboxLabel
                  label="Delete file default"
                  checked={draft.ui.deleteFileDefault}
                  onChange={(checked) => updateDraft((d) => ({ ...d, ui: { ...d.ui, deleteFileDefault: checked } }))}
                />
                <NumberField
                  label="Slideshow delay (ms)"
                  value={draft.ui.slideshowDelay}
                  min={500}
                  onChange={(value) => updateDraft((d) => ({ ...d, ui: { ...d.ui, slideshowDelay: value ?? d.ui.slideshowDelay } }))}
                />
              </div>
            </SectionCard>

            <SectionCard title="Custom CSS" description="Inject custom CSS into the application.">
              <div className="space-y-4">
                <CheckboxLabel
                  label="Enable CSS customization"
                  checked={draft.ui.enableCSSCustomization}
                  onChange={(checked) => updateDraft((d) => ({ ...d, ui: { ...d.ui, enableCSSCustomization: checked } }))}
                />
                {draft.ui.enableCSSCustomization && (
                  <TextAreaField
                    label="Custom CSS"
                    value={draft.ui.customCss ?? ""}
                    onChange={(value) => updateDraft((d) => ({ ...d, ui: { ...d.ui, customCss: value || undefined } }))}
                    rows={8}
                    placeholder="/* Enter custom CSS here */"
                  />
                )}
              </div>
            </SectionCard>

            <SectionCard title="Custom JavaScript" description="Inject custom JavaScript into the application.">
              <div className="space-y-4">
                <CheckboxLabel
                  label="Enable JavaScript customization"
                  checked={draft.ui.enableJSCustomization}
                  onChange={(checked) => updateDraft((d) => ({ ...d, ui: { ...d.ui, enableJSCustomization: checked } }))}
                />
                {draft.ui.enableJSCustomization && (
                  <TextAreaField
                    label="Custom JavaScript"
                    value={draft.ui.customJs ?? ""}
                    onChange={(value) => updateDraft((d) => ({ ...d, ui: { ...d.ui, customJs: value || undefined } }))}
                    rows={8}
                    placeholder="// Enter custom JavaScript here"
                  />
                )}
              </div>
            </SectionCard>
          </>
        )}

        {activeTab === "security" && (
          <>
            <SectionCard title="Authentication" description="These values persist to config immediately. Enabling or disabling auth may still require a restart for middleware changes.">
              <div className="space-y-4">
                <CheckboxLabel
                  label="Require authentication"
                  checked={draft.security.enabled}
                  onChange={(checked) => updateDraft((current) => ({ ...current, security: { ...current.security, enabled: checked } }))}
                />
                <div className="grid gap-4 md:grid-cols-2">
                  <TextField
                    label="Username"
                    value={draft.security.username ?? ""}
                    onChange={(value) => updateDraft((current) => ({ ...current, security: { ...current.security, username: value || undefined } }))}
                    placeholder="stash"
                  />
                  <NumberField
                    label="Maximum session age (minutes)"
                    value={draft.security.maxSessionAgeMinutes}
                    min={1}
                    onChange={(value) =>
                      updateDraft((current) => ({
                        ...current,
                        security: {
                          ...current.security,
                          maxSessionAgeMinutes: value ?? current.security.maxSessionAgeMinutes,
                        },
                      }))
                    }
                  />
                </div>
                <TextField
                  label="New password"
                  type="password"
                  value={draft.security.newPassword ?? ""}
                  onChange={(value) => updateDraft((current) => ({ ...current, security: { ...current.security, newPassword: value || undefined } }))}
                  placeholder="Leave blank to keep the current password"
                />
              </div>
            </SectionCard>
          </>
        )}

        {activeTab === "metadata-providers" && (
          <>
            <SectionCard title="Scraper Directories" description="Directories are scanned recursively for Stash-compatible YAML scraper definitions.">
              <div className="space-y-3">
                {draft.scraping.scraperDirectories.map((directory, index) => (
                  <div key={index} className="flex flex-col gap-2 md:flex-row md:items-center">
                    <input
                      type="text"
                      value={directory}
                      onChange={(event) =>
                        updateDraft((current) => ({
                          ...current,
                          scraping: {
                            ...current.scraping,
                            scraperDirectories: current.scraping.scraperDirectories.map((item, itemIndex) =>
                              itemIndex === index ? event.target.value : item,
                            ),
                          },
                        }))
                      }
                      placeholder="C:\\Users\\you\\AppData\\Local\\stash\\scrapers"
                      className="flex-1 rounded-xl border border-plex-border bg-plex-surface px-3 py-2 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
                    />
                    <button
                      onClick={() =>
                        updateDraft((current) => ({
                          ...current,
                          scraping: {
                            ...current.scraping,
                            scraperDirectories:
                              current.scraping.scraperDirectories.length > 1
                                ? current.scraping.scraperDirectories.filter((_, itemIndex) => itemIndex !== index)
                                : [""],
                          },
                        }))
                      }
                      className="inline-flex items-center gap-1 rounded-lg border border-plex-border px-2 py-2 text-xs text-red-300 hover:border-red-500 hover:text-red-200"
                    >
                      <Trash2 className="h-3.5 w-3.5" /> Remove
                    </button>
                  </div>
                ))}
                <button
                  onClick={() =>
                    updateDraft((current) => ({
                      ...current,
                      scraping: {
                        ...current.scraping,
                        scraperDirectories: [...current.scraping.scraperDirectories, ""],
                      },
                    }))
                  }
                  className="inline-flex items-center gap-2 rounded-xl border border-dashed border-plex-border px-3 py-2 text-sm text-plex-text-secondary hover:text-plex-text"
                >
                  <Plus className="h-4 w-4" /> Add scraper directory
                </button>
              </div>
            </SectionCard>

            <SectionCard title="Package Sources" description="Source URLs are stored now so scraper package installation can layer on top of the same config later.">
              <div className="space-y-3">
                {draft.scraping.scraperPackageSources.map((source, index) => (
                  <div key={index} className="grid gap-3 rounded-xl border border-plex-border bg-plex-card p-3 lg:grid-cols-[1fr_2fr_auto]">
                    <TextField
                      label="Source name"
                      value={source.name}
                      onChange={(value) =>
                        updateDraft((current) => ({
                          ...current,
                          scraping: {
                            ...current.scraping,
                            scraperPackageSources: current.scraping.scraperPackageSources.map((item, itemIndex) =>
                              itemIndex === index ? { ...item, name: value } : item,
                            ),
                          },
                        }))
                      }
                      placeholder="Official"
                    />
                    <TextField
                      label="Source URL"
                      value={source.url}
                      onChange={(value) =>
                        updateDraft((current) => ({
                          ...current,
                          scraping: {
                            ...current.scraping,
                            scraperPackageSources: current.scraping.scraperPackageSources.map((item, itemIndex) =>
                              itemIndex === index ? { ...item, url: value } : item,
                            ),
                          },
                        }))
                      }
                      placeholder="https://example.com/packages.yaml"
                    />
                    <div className="flex items-end">
                      <button
                        onClick={() =>
                          updateDraft((current) => ({
                            ...current,
                            scraping: {
                              ...current.scraping,
                              scraperPackageSources:
                                current.scraping.scraperPackageSources.length > 1
                                  ? current.scraping.scraperPackageSources.filter((_, itemIndex) => itemIndex !== index)
                                  : [emptyPackageSource()],
                            },
                          }))
                        }
                        className="inline-flex items-center gap-1 rounded-lg border border-plex-border px-2 py-2 text-xs text-red-300 hover:border-red-500 hover:text-red-200"
                      >
                        <Trash2 className="h-3.5 w-3.5" /> Remove
                      </button>
                    </div>
                  </div>
                ))}
                <button
                  onClick={() =>
                    updateDraft((current) => ({
                      ...current,
                      scraping: {
                        ...current.scraping,
                        scraperPackageSources: [...current.scraping.scraperPackageSources, emptyPackageSource()],
                      },
                    }))
                  }
                  className="inline-flex items-center gap-2 rounded-xl border border-dashed border-plex-border px-3 py-2 text-sm text-plex-text-secondary hover:text-plex-text"
                >
                  <Plus className="h-4 w-4" /> Add package source
                </button>
              </div>
            </SectionCard>

            <SectionCard title="StashBox Instances" description="Configure remote stash-box GraphQL endpoints, validate credentials, and use them from performer detail pages.">
              <div className="space-y-3">
                {draft.scraping.stashBoxes.length === 0 && (
                  <div className="rounded-xl border border-dashed border-plex-border p-4 text-sm text-plex-text-secondary">
                    No StashBox instances configured yet.
                  </div>
                )}

                {draft.scraping.stashBoxes.map((stashBox, index) => {
                  const validation = stashBoxValidation[String(index)];

                  return (
                    <div key={index} className="rounded-xl border border-plex-border bg-plex-card p-3">
                      <div className="grid gap-3 xl:grid-cols-[minmax(0,1fr)_minmax(0,2fr)_minmax(0,2fr)_160px_auto_auto]">
                        <TextField
                          label="Name"
                          value={stashBox.name}
                          onChange={(value) =>
                            updateDraft((current) => ({
                              ...current,
                              scraping: {
                                ...current.scraping,
                                stashBoxes: current.scraping.stashBoxes.map((item, itemIndex) =>
                                  itemIndex === index ? { ...item, name: value } : item,
                                ),
                              },
                            }))
                          }
                          placeholder="StashDB"
                        />
                        <TextField
                          label="Endpoint"
                          value={stashBox.endpoint}
                          onChange={(value) =>
                            updateDraft((current) => ({
                              ...current,
                              scraping: {
                                ...current.scraping,
                                stashBoxes: current.scraping.stashBoxes.map((item, itemIndex) =>
                                  itemIndex === index ? { ...item, endpoint: value } : item,
                                ),
                              },
                            }))
                          }
                          placeholder="https://stashdb.org/graphql"
                        />
                        <TextField
                          label="API key"
                          type="password"
                          value={stashBox.apiKey}
                          onChange={(value) =>
                            updateDraft((current) => ({
                              ...current,
                              scraping: {
                                ...current.scraping,
                                stashBoxes: current.scraping.stashBoxes.map((item, itemIndex) =>
                                  itemIndex === index ? { ...item, apiKey: value } : item,
                                ),
                              },
                            }))
                          }
                          placeholder="Paste API key"
                        />
                        <NumberField
                          label="Max req/min"
                          value={stashBox.maxRequestsPerMinute}
                          min={1}
                          onChange={(value) =>
                            updateDraft((current) => ({
                              ...current,
                              scraping: {
                                ...current.scraping,
                                stashBoxes: current.scraping.stashBoxes.map((item, itemIndex) =>
                                  itemIndex === index
                                    ? { ...item, maxRequestsPerMinute: value ?? item.maxRequestsPerMinute }
                                    : item,
                                ),
                              },
                            }))
                          }
                        />
                        <div className="flex items-end">
                          <button
                            onClick={() => validateStashBoxMutation.mutate({ index, stashBox })}
                            disabled={validateStashBoxMutation.isPending || !stashBox.endpoint.trim()}
                            className="inline-flex items-center gap-2 rounded-xl border border-plex-border px-3 py-2 text-sm text-plex-text hover:border-plex-accent hover:text-plex-accent disabled:opacity-60"
                          >
                            {validateStashBoxMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
                            Validate
                          </button>
                        </div>
                        <div className="flex items-end">
                          <button
                            onClick={() =>
                              updateDraft((current) => ({
                                ...current,
                                scraping: {
                                  ...current.scraping,
                                  stashBoxes: current.scraping.stashBoxes.filter((_, itemIndex) => itemIndex !== index),
                                },
                              }))
                            }
                            className="inline-flex items-center gap-1 rounded-lg border border-plex-border px-2 py-2 text-xs text-red-300 hover:border-red-500 hover:text-red-200"
                          >
                            <Trash2 className="h-3.5 w-3.5" /> Remove
                          </button>
                        </div>
                      </div>
                      {validation && (
                        <p className={`mt-3 text-sm ${validation.valid ? "text-emerald-300" : "text-red-300"}`}>
                          {validation.status}
                        </p>
                      )}
                    </div>
                  );
                })}

                <button
                  onClick={() =>
                    updateDraft((current) => ({
                      ...current,
                      scraping: {
                        ...current.scraping,
                        stashBoxes: [...current.scraping.stashBoxes, emptyStashBox()],
                      },
                    }))
                  }
                  className="inline-flex items-center gap-2 rounded-xl border border-dashed border-plex-border px-3 py-2 text-sm text-plex-text-secondary hover:text-plex-text"
                >
                  <Plus className="h-4 w-4" /> Add StashBox instance
                </button>
              </div>
            </SectionCard>

            <SectionCard title="Discovered Scrapers" description="Scraper definitions are loaded from the configured directories using the same YAML field names Stash expects.">
              <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
                <p className="text-sm text-plex-text-secondary">Reload after changing directories or adding new scraper files.</p>
                <button
                  onClick={() => reloadScrapersMutation.mutate()}
                  disabled={reloadScrapersMutation.isPending}
                  className="inline-flex items-center gap-2 rounded-xl border border-plex-border px-3 py-2 text-sm text-plex-text hover:border-plex-accent hover:text-plex-accent disabled:opacity-60"
                >
                  {reloadScrapersMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
                  Reload scrapers
                </button>
              </div>

              {scrapersLoading ? (
                <div className="flex items-center gap-2 text-sm text-plex-text-secondary">
                  <Loader2 className="h-4 w-4 animate-spin" /> Loading scrapers...
                </div>
              ) : scrapers.length === 0 ? (
                <div className="rounded-xl border border-dashed border-plex-border p-4 text-sm text-plex-text-secondary">
                  No scraper definitions were found in the configured directories.
                </div>
              ) : (
                <div className="space-y-4">
                  {Object.entries(groupedScrapers).map(([entityType, entityScrapers]) => (
                    <ScraperTable key={entityType} entityType={entityType} scrapers={entityScrapers} />
                  ))}
                </div>
              )}
            </SectionCard>
          </>
        )}

        {activeTab === "system" && (
          <>
            <SectionCard title="Server" description="Host and port are persisted immediately but require a restart to rebind the listener.">
              <div className="grid gap-4 md:grid-cols-3">
                <TextField
                  label="Host"
                  value={draft.host}
                  onChange={(value) => updateDraft((current) => ({ ...current, host: value }))}
                />
                <NumberField
                  label="Port"
                  value={draft.port}
                  min={1}
                  onChange={(value) => updateDraft((current) => ({ ...current, port: value ?? current.port }))}
                />
                <NumberField
                  label="Max parallel tasks"
                  value={draft.maxParallelTasks}
                  min={1}
                  max={32}
                  onChange={(value) => updateDraft((current) => ({ ...current, maxParallelTasks: value ?? current.maxParallelTasks }))}
                />
              </div>
            </SectionCard>

            <SectionCard title="FFmpeg" description="Paths to FFmpeg and FFprobe binaries. Leave blank to use system PATH.">
              <div className="grid gap-4 md:grid-cols-2">
                <TextField
                  label="FFmpeg path"
                  value={draft.ffmpegPath ?? ""}
                  onChange={(value) => updateDraft((d) => ({ ...d, ffmpegPath: value || undefined }))}
                  placeholder="C:\\ffmpeg\\bin\\ffmpeg.exe"
                />
                <TextField
                  label="FFprobe path"
                  value={draft.ffprobePath ?? ""}
                  onChange={(value) => updateDraft((d) => ({ ...d, ffprobePath: value || undefined }))}
                  placeholder="C:\\ffmpeg\\bin\\ffprobe.exe"
                />
              </div>
            </SectionCard>

            <SectionCard title="Transcoding" description="Hardware acceleration and transcode size limits. 0 means original resolution.">
              <div className="space-y-4">
                <div className="grid gap-4 md:grid-cols-2">
                  <NumberField
                    label="Max transcode size"
                    value={draft.maxTranscodeSize}
                    min={0}
                    onChange={(value) => updateDraft((d) => ({ ...d, maxTranscodeSize: value ?? d.maxTranscodeSize }))}
                  />
                  <NumberField
                    label="Max streaming transcode size"
                    value={draft.maxStreamingTranscodeSize}
                    min={0}
                    onChange={(value) => updateDraft((d) => ({ ...d, maxStreamingTranscodeSize: value ?? d.maxStreamingTranscodeSize }))}
                  />
                </div>
                <SelectField
                  label="Hardware acceleration"
                  value={draft.transcodeHardwareAcceleration}
                  onChange={(value) => updateDraft((d) => ({ ...d, transcodeHardwareAcceleration: value }))}
                  options={[
                    { value: "none", label: "None" },
                    { value: "nvenc", label: "NVENC" },
                    { value: "vaapi", label: "VAAPI" },
                    { value: "qsv", label: "QSV" },
                  ]}
                />
                <div className="grid gap-4 md:grid-cols-2">
                  <TextField
                    label="Transcode input args"
                    value={draft.transcodeInputArgs ?? ""}
                    onChange={(value) => updateDraft((d) => ({ ...d, transcodeInputArgs: value || undefined }))}
                  />
                  <TextField
                    label="Transcode output args"
                    value={draft.transcodeOutputArgs ?? ""}
                    onChange={(value) => updateDraft((d) => ({ ...d, transcodeOutputArgs: value || undefined }))}
                  />
                  <TextField
                    label="Live transcode input args"
                    value={draft.liveTranscodeInputArgs ?? ""}
                    onChange={(value) => updateDraft((d) => ({ ...d, liveTranscodeInputArgs: value || undefined }))}
                  />
                  <TextField
                    label="Live transcode output args"
                    value={draft.liveTranscodeOutputArgs ?? ""}
                    onChange={(value) => updateDraft((d) => ({ ...d, liveTranscodeOutputArgs: value || undefined }))}
                  />
                </div>
              </div>
            </SectionCard>

            <SectionCard title="Preview Generation" description="Settings for preview video generation during scanning.">
              <div className="space-y-4">
                <SelectField
                  label="Preview preset"
                  value={draft.previewPreset}
                  onChange={(value) => updateDraft((d) => ({ ...d, previewPreset: value }))}
                  options={[
                    { value: "ultrafast", label: "Ultrafast" },
                    { value: "veryfast", label: "Very Fast" },
                    { value: "fast", label: "Fast" },
                    { value: "medium", label: "Medium" },
                    { value: "slow", label: "Slow" },
                    { value: "slower", label: "Slower" },
                    { value: "veryslow", label: "Very Slow" },
                  ]}
                />
                <CheckboxLabel
                  label="Include audio in previews"
                  checked={draft.previewAudio === "true"}
                  onChange={(checked) => updateDraft((d) => ({ ...d, previewAudio: checked ? "true" : "false" }))}
                />
              </div>
            </SectionCard>

            <SectionCard title="Logging" description="Log level and output configuration. Changes take effect after restart.">
              <div className="space-y-4">
                <div className="grid gap-4 md:grid-cols-2">
                  <SelectField
                    label="Log level"
                    value={draft.logLevel}
                    onChange={(value) => updateDraft((d) => ({ ...d, logLevel: value }))}
                    options={[
                      { value: "Trace", label: "Trace" },
                      { value: "Debug", label: "Debug" },
                      { value: "Info", label: "Info" },
                      { value: "Warning", label: "Warning" },
                      { value: "Error", label: "Error" },
                    ]}
                  />
                  <TextField
                    label="Log file"
                    value={draft.logFile ?? ""}
                    onChange={(value) => updateDraft((d) => ({ ...d, logFile: value || undefined }))}
                    placeholder="Leave blank for no file logging"
                  />
                </div>
                <CheckboxLabel
                  label="Log to stdout"
                  checked={draft.logOut}
                  onChange={(checked) => updateDraft((d) => ({ ...d, logOut: checked }))}
                />
                <CheckboxLabel
                  label="Log access requests"
                  checked={draft.logAccess}
                  onChange={(checked) => updateDraft((d) => ({ ...d, logAccess: checked }))}
                />
              </div>
            </SectionCard>
          </>
        )}

        {activeTab === "plugins" && <PluginsPanel />}

        {activeTab === "dlna" && <DlnaPanel />}

        {activeTab === "logs" && <LogsPanel />}

        {activeTab === "about" && (
          <>
            <SectionCard title="About Stash" description="An organizer for your media library.">
              <div className="flex items-start gap-6">
                <div className="w-16 h-16 rounded-xl bg-plex-accent/20 flex items-center justify-center shrink-0">
                  <span className="text-3xl font-bold text-plex-accent">S</span>
                </div>
                <div className="space-y-2">
                  <h2 className="text-2xl font-bold text-plex-text">Stash</h2>
                  {status && <p className="text-sm text-plex-text-secondary">Version {status.version}</p>}
                  <p className="text-sm text-plex-text-muted max-w-lg">
                    A self-hosted media organizer and video streaming app. Organize, tag, and browse your media library with ease.
                  </p>
                  <div className="flex gap-3 pt-1">
                    <a href="https://github.com/stashapp/stash" target="_blank" rel="noopener noreferrer" className="text-xs text-plex-accent hover:underline">GitHub</a>
                    <a href="https://docs.stashapp.cc" target="_blank" rel="noopener noreferrer" className="text-xs text-plex-accent hover:underline">Documentation</a>
                    <a href="https://discord.gg/2TsNFKt" target="_blank" rel="noopener noreferrer" className="text-xs text-plex-accent hover:underline">Discord</a>
                  </div>
                </div>
              </div>
            </SectionCard>

            <SectionCard title="Runtime Status" description="Effective values reported by the running backend instance.">
              {statusLoading && !status ? (
                <div className="flex items-center gap-2 text-sm text-plex-text-secondary">
                  <Loader2 className="h-4 w-4 animate-spin" /> Loading status...
                </div>
              ) : status ? (
                <dl className="grid gap-4 md:grid-cols-2">
                  <InfoPair label="Version" value={status.version} />
                  <InfoPair label="Database" value={status.databasePath} />
                  <InfoPair label="Config file" value={status.configFile} />
                  <InfoPair label="App directory" value={status.appDir} />
                </dl>
              ) : (
                <div className="text-sm text-plex-text-secondary">Runtime status is unavailable.</div>
              )}
            </SectionCard>

            <SectionCard title="System Information" description="Browser and environment details.">
              <dl className="grid gap-4 md:grid-cols-2">
                <InfoPair label="Browser" value={navigator.userAgent.split(/[()]/)[1] || navigator.userAgent.substring(0, 60)} />
                <InfoPair label="Platform" value={navigator.platform} />
                <InfoPair label="Screen resolution" value={`${screen.width}×${screen.height}`} />
                <InfoPair label="Language" value={navigator.language} />
              </dl>
            </SectionCard>

            <SectionCard title="Current Config Summary" description="High-level values from the effective client-side config object.">
              <dl className="grid gap-4 md:grid-cols-2">
                <InfoPair label="Library paths" value={String(draft.stashPaths.filter((path) => path.path.trim() !== "").length)} />
                <InfoPair label="Scraper directories" value={String(draft.scraping.scraperDirectories.filter(Boolean).length)} />
                <InfoPair label="Stash boxes" value={String(draft.scraping.stashBoxes.filter((box) => box.endpoint.trim() !== "").length)} />
                <InfoPair label="Rating system" value={draft.ui.ratingSystemOptions.type} />
                <InfoPair label="Authentication" value={draft.security.enabled ? "enabled" : "disabled"} />
              </dl>
            </SectionCard>

            <SectionCard title="Keyboard Shortcuts" description="Press ? anywhere to view the full shortcut reference.">
              <div className="grid gap-3 md:grid-cols-2 text-sm">
                <div className="flex justify-between"><span className="text-plex-text-secondary">Global navigation</span><span className="font-mono text-xs bg-gray-800 px-2 py-0.5 rounded">g</span></div>
                <div className="flex justify-between"><span className="text-plex-text-secondary">Theater mode (scenes)</span><span className="font-mono text-xs bg-gray-800 px-2 py-0.5 rounded">,</span></div>
                <div className="flex justify-between"><span className="text-plex-text-secondary">Show all shortcuts</span><span className="font-mono text-xs bg-gray-800 px-2 py-0.5 rounded">?</span></div>
                <div className="flex justify-between"><span className="text-plex-text-secondary">Search / filter</span><span className="font-mono text-xs bg-gray-800 px-2 py-0.5 rounded">/</span></div>
              </div>
            </SectionCard>
          </>
        )}
      </div>
    </div>
  );
}

function LogsPanel() {
  const [logLevel, setLogLevel] = useState("");
  const { data: logEntries, isLoading, refetch } = useQuery({
    queryKey: ["logs", logLevel],
    queryFn: () => logsApi.recent(logLevel || undefined, 200),
    refetchInterval: 5000,
  });

  const levelColor = (level: string) => {
    switch (level.toLowerCase()) {
      case "error": case "critical": return "text-red-400";
      case "warning": return "text-yellow-400";
      case "debug": return "text-blue-400";
      case "trace": return "text-gray-500";
      default: return "text-plex-text-secondary";
    }
  };

  return (
    <SectionCard title="Logs" description="Recent log entries from the server.">
      <div className="flex items-center gap-3 mb-4">
        <label className="text-sm text-plex-text-secondary">Log Level</label>
        <select
          value={logLevel}
          onChange={(e) => setLogLevel(e.target.value)}
          className="rounded border border-plex-border bg-plex-surface px-3 py-1.5 text-sm text-plex-text"
        >
          <option value="">All</option>
          <option value="Trace">Trace</option>
          <option value="Debug">Debug</option>
          <option value="Information">Info</option>
          <option value="Warning">Warning</option>
          <option value="Error">Error</option>
        </select>
        <button onClick={() => refetch()} className="flex items-center gap-1 rounded border border-plex-border bg-plex-surface px-3 py-1.5 text-sm text-plex-text-secondary hover:text-plex-text">
          <RefreshCw className="h-3.5 w-3.5" /> Refresh
        </button>
      </div>
      {isLoading ? (
        <div className="flex items-center gap-2 text-sm text-plex-text-secondary">
          <Loader2 className="h-4 w-4 animate-spin" /> Loading logs...
        </div>
      ) : logEntries && logEntries.length > 0 ? (
        <div className="max-h-[600px] overflow-y-auto rounded border border-plex-border bg-gray-950 font-mono text-xs">
          {logEntries.map((entry, i) => (
            <div key={i} className="flex gap-3 border-b border-gray-800/50 px-3 py-1.5 hover:bg-gray-900">
              <span className="shrink-0 text-plex-text-muted">{entry.timestamp}</span>
              <span className={`shrink-0 w-14 font-semibold ${levelColor(entry.level)}`}>{entry.level}</span>
              <span className="text-plex-text break-all">{entry.message}</span>
            </div>
          ))}
        </div>
      ) : (
        <p className="text-sm text-plex-text-muted">No log entries found.</p>
      )}
    </SectionCard>
  );
}

function TasksPanel() {
  const queryClient = useQueryClient();
  const { data: activeJobs, refetch: refetchJobs } = useQuery({
    queryKey: ["jobs"],
    queryFn: () => jobs.list(),
    refetchInterval: 2000,
  });

  // ---- Job Queue ----
  const jobQueue = activeJobs && activeJobs.length > 0 ? (
    <SectionCard title="Job Queue" description="Currently running or queued jobs.">
      <div className="space-y-2">
        {activeJobs.map((job) => (
          <div key={job.id} className="flex items-center justify-between rounded-xl border border-plex-border bg-plex-card p-3">
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <span className="text-sm font-medium text-plex-text">{job.description}</span>
                <span className={`text-xs px-1.5 py-0.5 rounded ${
                  job.status === "Running" ? "bg-green-600/20 text-green-300" :
                  job.status === "Pending" ? "bg-yellow-600/20 text-yellow-300" :
                  "bg-plex-card text-plex-text-muted"
                }`}>
                  {job.status}
                </span>
              </div>
              {job.progress != null && job.progress >= 0 && (
                <div className="mt-2 h-1.5 w-full rounded-full bg-plex-card overflow-hidden">
                  <div className="h-full rounded-full bg-plex-accent transition-all" style={{ width: `${Math.min(job.progress * 100, 100)}%` }} />
                </div>
              )}
            </div>
            <button
              onClick={() => jobs.cancel(job.id).then(() => refetchJobs())}
              className="ml-3 text-xs text-plex-text-muted hover:text-red-300 flex-shrink-0"
            >
              Cancel
            </button>
          </div>
        ))}
      </div>
    </SectionCard>
  ) : null;

  return (
    <>
      {jobQueue}
      <LibraryTasksSection refetchJobs={refetchJobs} />
      <DataManagementSection refetchJobs={refetchJobs} />
      <PluginTasksSection refetchJobs={refetchJobs} />
    </>
  );
}

// ---- Library Tasks ----
function LibraryTasksSection({ refetchJobs }: { refetchJobs: () => void }) {
  const [showScanOpts, setShowScanOpts] = useState(false);
  const [scanOpts, setScanOpts] = useState<ScanOptions>({
    scanGenerateCovers: true,
    scanGeneratePreviews: false,
    scanGenerateSprites: false,
    scanGeneratePhashes: false,
    scanGenerateThumbnails: false,
    scanGenerateImagePhashes: false,
    rescan: false,
  });

  const [showGenOpts, setShowGenOpts] = useState(false);
  const [genOpts, setGenOpts] = useState<GenerateOptions>({
    thumbnails: true,
    previews: false,
    sprites: false,
    markers: false,
    phashes: false,
    imageThumbnails: false,
    imagePhashes: false,
    overwrite: false,
  });

  const scanMut = useMutation({ mutationFn: () => metadata.scan(scanOpts), onSuccess: () => refetchJobs() });
  const genMut = useMutation({ mutationFn: () => metadata.generate(genOpts), onSuccess: () => refetchJobs() });
  const autoTagMut = useMutation({ mutationFn: () => metadata.autoTag(), onSuccess: () => refetchJobs() });

  return (
    <SectionCard title="Library Tasks" description="Scan for new content, generate supporting files, and auto-tag your library.">
      <div className="space-y-4">
        {/* Scan */}
        <TaskCard
          label="Scan"
          description="Scan library paths for new content and update metadata."
          onRun={() => scanMut.mutate()}
          isPending={scanMut.isPending}
          expandable
          expanded={showScanOpts}
          onToggleExpand={() => setShowScanOpts(!showScanOpts)}
        >
          <div className="grid gap-2 sm:grid-cols-2 pt-3 border-t border-plex-border/50">
            <CheckboxLabel label="Generate covers" checked={!!scanOpts.scanGenerateCovers} onChange={(c) => setScanOpts({ ...scanOpts, scanGenerateCovers: c })} />
            <CheckboxLabel label="Generate previews" checked={!!scanOpts.scanGeneratePreviews} onChange={(c) => setScanOpts({ ...scanOpts, scanGeneratePreviews: c })} />
            <CheckboxLabel label="Generate sprites" checked={!!scanOpts.scanGenerateSprites} onChange={(c) => setScanOpts({ ...scanOpts, scanGenerateSprites: c })} />
            <CheckboxLabel label="Generate perceptual hashes" checked={!!scanOpts.scanGeneratePhashes} onChange={(c) => setScanOpts({ ...scanOpts, scanGeneratePhashes: c })} />
            <CheckboxLabel label="Generate image thumbnails" checked={!!scanOpts.scanGenerateThumbnails} onChange={(c) => setScanOpts({ ...scanOpts, scanGenerateThumbnails: c })} />
            <CheckboxLabel label="Generate image phashes" checked={!!scanOpts.scanGenerateImagePhashes} onChange={(c) => setScanOpts({ ...scanOpts, scanGenerateImagePhashes: c })} />
            <CheckboxLabel label="Force rescan (ignore mtime)" checked={!!scanOpts.rescan} onChange={(c) => setScanOpts({ ...scanOpts, rescan: c })} />
          </div>
        </TaskCard>

        {/* Auto Tag */}
        <TaskCard
          label="Auto Tag"
          description="Automatically tag content based on filenames and path patterns."
          onRun={() => autoTagMut.mutate()}
          isPending={autoTagMut.isPending}
        />

        {/* Generate */}
        <TaskCard
          label="Generate"
          description="Generate thumbnails, previews, sprites, markers, and perceptual hashes."
          onRun={() => genMut.mutate()}
          isPending={genMut.isPending}
          expandable
          expanded={showGenOpts}
          onToggleExpand={() => setShowGenOpts(!showGenOpts)}
        >
          <div className="space-y-3 pt-3 border-t border-plex-border/50">
            <p className="text-xs text-plex-text-muted font-medium uppercase tracking-wide">Scene options</p>
            <div className="grid gap-2 sm:grid-cols-2">
              <CheckboxLabel label="Thumbnails / screenshots" checked={!!genOpts.thumbnails} onChange={(c) => setGenOpts({ ...genOpts, thumbnails: c })} />
              <CheckboxLabel label="Video previews" checked={!!genOpts.previews} onChange={(c) => setGenOpts({ ...genOpts, previews: c })} />
              <CheckboxLabel label="Sprite sheets" checked={!!genOpts.sprites} onChange={(c) => setGenOpts({ ...genOpts, sprites: c })} />
              <CheckboxLabel label="Marker previews" checked={!!genOpts.markers} onChange={(c) => setGenOpts({ ...genOpts, markers: c })} />
              <CheckboxLabel label="Perceptual hashes (phash)" checked={!!genOpts.phashes} onChange={(c) => setGenOpts({ ...genOpts, phashes: c })} />
            </div>
            <p className="text-xs text-plex-text-muted font-medium uppercase tracking-wide pt-2">Image options</p>
            <div className="grid gap-2 sm:grid-cols-2">
              <CheckboxLabel label="Image thumbnails" checked={!!genOpts.imageThumbnails} onChange={(c) => setGenOpts({ ...genOpts, imageThumbnails: c })} />
              <CheckboxLabel label="Image phashes" checked={!!genOpts.imagePhashes} onChange={(c) => setGenOpts({ ...genOpts, imagePhashes: c })} />
            </div>
            <div className="pt-2">
              <CheckboxLabel label="Overwrite existing generated files" checked={!!genOpts.overwrite} onChange={(c) => setGenOpts({ ...genOpts, overwrite: c })} />
            </div>
          </div>
        </TaskCard>
      </div>
    </SectionCard>
  );
}

// ---- Data Management ----
function DataManagementSection({ refetchJobs }: { refetchJobs: () => void }) {
  const [cleanDryRun, setCleanDryRun] = useState(false);
  const [showCleanGenOpts, setShowCleanGenOpts] = useState(false);
  const [cleanGenOpts, setCleanGenOpts] = useState<CleanGeneratedOptions>({
    screenshots: true,
    sprites: true,
    transcodes: true,
    markers: true,
    imageThumbnails: true,
    dryRun: false,
  });

  const [showExportOpts, setShowExportOpts] = useState(false);
  const [exportOpts, setExportOpts] = useState<ExportOptions>({
    includeScenes: true,
    includePerformers: true,
    includeStudios: true,
    includeTags: true,
    includeGalleries: true,
    includeGroups: true,
  });

  const cleanMut = useMutation({ mutationFn: () => metadata.clean({ dryRun: cleanDryRun }), onSuccess: () => refetchJobs() });
  const cleanGenMut = useMutation({ mutationFn: () => metadata.cleanGenerated(cleanGenOpts), onSuccess: () => refetchJobs() });
  const exportMut = useMutation({ mutationFn: () => metadata.export(exportOpts), onSuccess: () => refetchJobs() });
  const backupMut = useMutation({ mutationFn: () => database.backup(), onSuccess: () => refetchJobs() });
  const optimizeMut = useMutation({ mutationFn: () => database.optimize(), onSuccess: () => refetchJobs() });

  return (
    <SectionCard title="Data Management" description="Clean orphaned data, manage generated files, export, and database operations.">
      <div className="space-y-4">
        {/* Clean */}
        <TaskCard
          label="Clean"
          description="Find and remove database entries for files that no longer exist on disk."
          onRun={() => cleanMut.mutate()}
          isPending={cleanMut.isPending}
        >
          <div className="pt-3 border-t border-plex-border/50">
            <CheckboxLabel label="Dry run (report only, don't delete)" checked={cleanDryRun} onChange={setCleanDryRun} />
          </div>
        </TaskCard>

        {/* Clean Generated Files */}
        <TaskCard
          label="Clean Generated Files"
          description="Remove generated files (screenshots, sprites, transcodes, etc.) that are no longer needed."
          onRun={() => cleanGenMut.mutate()}
          isPending={cleanGenMut.isPending}
          expandable
          expanded={showCleanGenOpts}
          onToggleExpand={() => setShowCleanGenOpts(!showCleanGenOpts)}
        >
          <div className="grid gap-2 sm:grid-cols-2 pt-3 border-t border-plex-border/50">
            <CheckboxLabel label="Screenshots" checked={!!cleanGenOpts.screenshots} onChange={(c) => setCleanGenOpts({ ...cleanGenOpts, screenshots: c })} />
            <CheckboxLabel label="Sprites" checked={!!cleanGenOpts.sprites} onChange={(c) => setCleanGenOpts({ ...cleanGenOpts, sprites: c })} />
            <CheckboxLabel label="Transcodes" checked={!!cleanGenOpts.transcodes} onChange={(c) => setCleanGenOpts({ ...cleanGenOpts, transcodes: c })} />
            <CheckboxLabel label="Markers" checked={!!cleanGenOpts.markers} onChange={(c) => setCleanGenOpts({ ...cleanGenOpts, markers: c })} />
            <CheckboxLabel label="Image thumbnails" checked={!!cleanGenOpts.imageThumbnails} onChange={(c) => setCleanGenOpts({ ...cleanGenOpts, imageThumbnails: c })} />
            <CheckboxLabel label="Dry run" checked={!!cleanGenOpts.dryRun} onChange={(c) => setCleanGenOpts({ ...cleanGenOpts, dryRun: c })} />
          </div>
        </TaskCard>

        {/* Export */}
        <TaskCard
          label="Full Export"
          description="Export database content to JSON metadata files."
          onRun={() => exportMut.mutate()}
          isPending={exportMut.isPending}
          expandable
          expanded={showExportOpts}
          onToggleExpand={() => setShowExportOpts(!showExportOpts)}
        >
          <div className="grid gap-2 sm:grid-cols-2 pt-3 border-t border-plex-border/50">
            <CheckboxLabel label="Scenes" checked={!!exportOpts.includeScenes} onChange={(c) => setExportOpts({ ...exportOpts, includeScenes: c })} />
            <CheckboxLabel label="Performers" checked={!!exportOpts.includePerformers} onChange={(c) => setExportOpts({ ...exportOpts, includePerformers: c })} />
            <CheckboxLabel label="Studios" checked={!!exportOpts.includeStudios} onChange={(c) => setExportOpts({ ...exportOpts, includeStudios: c })} />
            <CheckboxLabel label="Tags" checked={!!exportOpts.includeTags} onChange={(c) => setExportOpts({ ...exportOpts, includeTags: c })} />
            <CheckboxLabel label="Galleries" checked={!!exportOpts.includeGalleries} onChange={(c) => setExportOpts({ ...exportOpts, includeGalleries: c })} />
            <CheckboxLabel label="Groups" checked={!!exportOpts.includeGroups} onChange={(c) => setExportOpts({ ...exportOpts, includeGroups: c })} />
          </div>
        </TaskCard>

        {/* Database Operations */}
        <div className="grid gap-3 sm:grid-cols-2">
          <TaskCard
            label="Backup Database"
            description="Create a pg_dump backup of the PostgreSQL database."
            onRun={() => backupMut.mutate()}
            isPending={backupMut.isPending}
          />
          <TaskCard
            label="Optimise Database"
            description="Run VACUUM ANALYSE to reclaim space and update query planner statistics."
            onRun={() => optimizeMut.mutate()}
            isPending={optimizeMut.isPending}
          />
        </div>
      </div>
    </SectionCard>
  );
}

// ---- Plugin Tasks ----
function PluginTasksSection({ refetchJobs }: { refetchJobs: () => void }) {
  const { data: pluginList } = useQuery({ queryKey: ["plugins"], queryFn: pluginsApi.list });
  const runTaskMut = useMutation({
    mutationFn: pluginsApi.runTask,
    onSuccess: () => refetchJobs(),
  });

  const enabledPlugins = pluginList?.filter((p) => p.enabled && p.tasks.length > 0) ?? [];

  if (enabledPlugins.length === 0) return null;

  return (
    <SectionCard title="Plugin Tasks" description="Run tasks provided by enabled plugins.">
      <div className="space-y-4">
        {enabledPlugins.map((plugin) => (
          <div key={plugin.id} className="rounded-xl border border-plex-border bg-plex-card overflow-hidden">
            <div className="px-4 py-2.5 border-b border-plex-border bg-black/10 flex items-center gap-2">
              <Plug className="h-3.5 w-3.5 text-plex-text-muted" />
              <span className="text-sm font-medium text-plex-text">{plugin.name}</span>
              <span className="text-xs text-plex-text-muted">v{plugin.version}</span>
            </div>
            <div className="divide-y divide-plex-border/50">
              {plugin.tasks.map((task) => (
                <div key={task.name} className="flex items-center justify-between px-4 py-3">
                  <div>
                    <h4 className="text-sm font-medium text-plex-text">{task.name}</h4>
                    {task.description && <p className="text-xs text-plex-text-secondary mt-0.5">{task.description}</p>}
                  </div>
                  <button
                    onClick={() => runTaskMut.mutate({ pluginId: plugin.id, taskName: task.name })}
                    disabled={runTaskMut.isPending}
                    className="inline-flex items-center gap-1.5 rounded-lg bg-plex-accent px-3 py-1.5 text-xs font-medium text-white hover:bg-plex-accent-hover disabled:opacity-60"
                  >
                    {runTaskMut.isPending ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <PlayCircle className="h-3.5 w-3.5" />}
                    Run
                  </button>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </SectionCard>
  );
}

// ---- Task Card (reusable) ----
function TaskCard({
  label,
  description,
  onRun,
  isPending,
  expandable,
  expanded,
  onToggleExpand,
  children,
}: {
  label: string;
  description: string;
  onRun: () => void;
  isPending: boolean;
  expandable?: boolean;
  expanded?: boolean;
  onToggleExpand?: () => void;
  children?: React.ReactNode;
}) {
  return (
    <div className="rounded-xl border border-plex-border bg-plex-card p-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 min-w-0 flex-1">
          {expandable && onToggleExpand && (
            <button onClick={onToggleExpand} className="text-plex-text-muted hover:text-plex-text flex-shrink-0">
              {expanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
            </button>
          )}
          <div>
            <h4 className="text-sm font-medium text-plex-text">{label}</h4>
            <p className="text-xs text-plex-text-secondary mt-0.5">{description}</p>
          </div>
        </div>
        <button
          onClick={onRun}
          disabled={isPending}
          className="inline-flex items-center gap-2 rounded-lg bg-plex-accent px-4 py-2 text-sm font-medium text-white hover:bg-plex-accent-hover disabled:opacity-60 flex-shrink-0 ml-3"
        >
          {isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <PlayCircle className="h-4 w-4" />}
          Run
        </button>
      </div>
      {/* Always show children (e.g. Clean dry-run checkbox), or show only when expanded */}
      {children && (!expandable || expanded) && (
        <div className="mt-3">{children}</div>
      )}
    </div>
  );
}

function SectionCard({ title, description, children }: { title: string; description: string; children: React.ReactNode }) {
  return (
    <section className="rounded-2xl border border-plex-border bg-plex-surface p-5 shadow-[0_12px_30px_-20px_rgba(0,0,0,0.7)]">
      <div className="mb-4">
        <h3 className="text-base font-semibold text-plex-text">{title}</h3>
        <p className="mt-1 text-sm text-plex-text-secondary">{description}</p>
      </div>
      {children}
    </section>
  );
}

function TextField({
  label,
  value,
  onChange,
  placeholder,
  type = "text",
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  type?: string;
}) {
  return (
    <label className="block text-sm">
      <span className="mb-1 block text-xs font-medium uppercase tracking-wide text-plex-text-muted">{label}</span>
      <input
        type={type}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        placeholder={placeholder}
        className="w-full rounded-xl border border-plex-border bg-plex-card px-3 py-2 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
      />
    </label>
  );
}

function NumberField({
  label,
  value,
  onChange,
  min,
  max,
}: {
  label: string;
  value?: number;
  onChange: (value: number | undefined) => void;
  min?: number;
  max?: number;
}) {
  return (
    <label className="block text-sm">
      <span className="mb-1 block text-xs font-medium uppercase tracking-wide text-plex-text-muted">{label}</span>
      <input
        type="number"
        value={value ?? ""}
        min={min}
        max={max}
        onChange={(event) => onChange(event.target.value ? Number(event.target.value) : undefined)}
        className="w-full rounded-xl border border-plex-border bg-plex-card px-3 py-2 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
      />
    </label>
  );
}

function TextAreaField({
  label,
  value,
  onChange,
  rows,
  placeholder,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  rows: number;
  placeholder?: string;
}) {
  return (
    <label className="block text-sm">
      <span className="mb-1 block text-xs font-medium uppercase tracking-wide text-plex-text-muted">{label}</span>
      <textarea
        value={value}
        onChange={(event) => onChange(event.target.value)}
        rows={rows}
        placeholder={placeholder}
        className="w-full rounded-xl border border-plex-border bg-plex-card px-3 py-2 font-mono text-sm text-plex-text focus:border-plex-accent focus:outline-none"
      />
    </label>
  );
}

function SelectField({
  label,
  value,
  onChange,
  options,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  options: { value: string; label: string }[];
}) {
  return (
    <label className="block text-sm">
      <span className="mb-1 block text-xs font-medium uppercase tracking-wide text-plex-text-muted">{label}</span>
      <select
        value={value}
        onChange={(event) => onChange(event.target.value)}
        className="w-full rounded-xl border border-plex-border bg-plex-card px-3 py-2 text-sm text-plex-text focus:border-plex-accent focus:outline-none"
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </label>
  );
}

function CheckboxLabel({ label, checked, onChange }: { label: string; checked: boolean; onChange: (checked: boolean) => void }) {
  return (
    <label className="flex items-center gap-2 text-sm text-plex-text-secondary">
      <input
        type="checkbox"
        checked={checked}
        onChange={(event) => onChange(event.target.checked)}
        className="h-4 w-4 rounded border-plex-border bg-plex-card text-plex-accent focus:ring-0"
      />
      <span>{label}</span>
    </label>
  );
}

function InfoPair({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-xl border border-plex-border bg-plex-card p-3">
      <dt className="text-xs font-medium uppercase tracking-wide text-plex-text-muted">{label}</dt>
      <dd className="mt-1 break-all text-sm text-plex-text">{value}</dd>
    </div>
  );
}

function ScraperTable({ entityType, scrapers }: { entityType: string; scrapers: ScraperSummary[] }) {
  return (
    <div className="overflow-hidden rounded-xl border border-plex-border bg-plex-card">
      <div className="border-b border-plex-border px-4 py-3 text-sm font-semibold capitalize text-plex-text">
        {entityType} scrapers <span className="text-plex-text-muted">({scrapers.length})</span>
      </div>
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-plex-border text-sm">
          <thead className="bg-black/10 text-left text-xs uppercase tracking-wide text-plex-text-muted">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Supported types</th>
              <th className="px-4 py-3">Supported URLs</th>
              <th className="px-4 py-3">Source</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-plex-border/70">
            {scrapers.map((scraper) => (
              <tr key={scraper.id}>
                <td className="px-4 py-3 font-medium text-plex-text">{scraper.name}</td>
                <td className="px-4 py-3 text-plex-text-secondary">{scraper.supportedScrapes.join(", ")}</td>
                <td className="px-4 py-3 text-plex-text-secondary">
                  {scraper.urls.length > 0 ? scraper.urls.join(", ") : <span className="text-plex-text-muted">No URL matchers</span>}
                </td>
                <td className="px-4 py-3 text-xs text-plex-text-muted">{scraper.sourcePath}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

// ===== DLNA Panel =====
function DlnaPanel() {
  const queryClient = useQueryClient();
  const [newIp, setNewIp] = useState("");

  const { data: status, isLoading } = useQuery({
    queryKey: ["dlna-status"],
    queryFn: () => dlnaApi.status(),
  });

  const enableMut = useMutation({
    mutationFn: (durationMinutes?: number) => dlnaApi.enable(durationMinutes),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["dlna-status"] }),
  });

  const disableMut = useMutation({
    mutationFn: () => dlnaApi.disable(),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["dlna-status"] }),
  });

  const allowIpMut = useMutation({
    mutationFn: (ip: string) => dlnaApi.allowIp(ip),
    onSuccess: () => { setNewIp(""); queryClient.invalidateQueries({ queryKey: ["dlna-status"] }); },
  });

  const removeIpMut = useMutation({
    mutationFn: (ip: string) => dlnaApi.removeIp(ip),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["dlna-status"] }),
  });

  if (isLoading || !status) {
    return (
      <SectionCard title="DLNA Server" description="Loading...">
        <Loader2 className="w-5 h-5 animate-spin text-plex-text-muted" />
      </SectionCard>
    );
  }

  return (
    <>
      <SectionCard title="DLNA Server" description="Enable or disable the DLNA media server for streaming to compatible devices.">
        <div className="space-y-4">
          <div className="flex items-center gap-4">
            <div className={`w-3 h-3 rounded-full ${status.running ? "bg-green-500" : "bg-gray-500"}`} />
            <span className="text-sm font-medium text-plex-text">{status.running ? "Running" : "Stopped"}</span>
            {status.untilDisabled && (
              <span className="text-xs text-plex-text-muted">
                (enabled until {new Date(status.untilDisabled).toLocaleTimeString()})
              </span>
            )}
          </div>
          <div className="flex flex-wrap gap-2">
            {!status.running ? (
              <>
                <button onClick={() => enableMut.mutate(undefined)} className="px-3 py-1.5 text-sm bg-green-600 hover:bg-green-500 text-white rounded">Enable</button>
                <button onClick={() => enableMut.mutate(120)} className="px-3 py-1.5 text-sm bg-plex-card border border-plex-border text-plex-text-secondary hover:text-plex-text rounded">Enable for 2 hours</button>
                <button onClick={() => enableMut.mutate(1440)} className="px-3 py-1.5 text-sm bg-plex-card border border-plex-border text-plex-text-secondary hover:text-plex-text rounded">Enable for 24 hours</button>
              </>
            ) : (
              <button onClick={() => disableMut.mutate()} className="px-3 py-1.5 text-sm bg-red-600 hover:bg-red-500 text-white rounded">Disable</button>
            )}
          </div>
        </div>
      </SectionCard>

      <SectionCard title="Allowed IP Addresses" description="Configure which IP addresses are allowed to access the DLNA server.">
        <div className="space-y-3">
          <div className="flex gap-2">
            <input
              type="text"
              value={newIp}
              onChange={(e) => setNewIp(e.target.value)}
              placeholder="IP address (e.g. 192.168.1.100)"
              className="flex-1 bg-plex-input border border-plex-border rounded px-3 py-1.5 text-sm text-plex-text"
              onKeyDown={(e) => { if (e.key === "Enter" && newIp.trim()) allowIpMut.mutate(newIp.trim()); }}
            />
            <button
              onClick={() => newIp.trim() && allowIpMut.mutate(newIp.trim())}
              className="px-3 py-1.5 text-sm bg-plex-accent hover:bg-plex-accent-hover text-white rounded flex items-center gap-1"
            >
              <Plus className="w-3.5 h-3.5" /> Add
            </button>
          </div>
          {status.allowedIps.length === 0 ? (
            <p className="text-xs text-plex-text-muted">No IP addresses allowed. All devices can connect when the server is running.</p>
          ) : (
            <div className="space-y-1">
              {status.allowedIps.map((ip) => (
                <div key={ip} className="flex items-center justify-between bg-plex-card border border-plex-border rounded px-3 py-2 text-sm">
                  <span className="text-plex-text font-mono">{ip}</span>
                  <button onClick={() => removeIpMut.mutate(ip)} className="text-plex-text-muted hover:text-red-400"><Trash2 className="w-3.5 h-3.5" /></button>
                </div>
              ))}
            </div>
          )}
        </div>
      </SectionCard>

      {status.recentIps.length > 0 && (
        <SectionCard title="Recent Connections" description="IP addresses that have recently connected to the DLNA server.">
          <div className="space-y-1">
            {status.recentIps.map((ip) => (
              <div key={ip} className="flex items-center justify-between bg-plex-card border border-plex-border rounded px-3 py-2 text-sm">
                <span className="text-plex-text font-mono">{ip}</span>
                <button
                  onClick={() => allowIpMut.mutate(ip)}
                  className="text-xs text-plex-accent hover:underline"
                >
                  Allow
                </button>
              </div>
            ))}
          </div>
        </SectionCard>
      )}
    </>
  );
}

// ===== Plugins Panel =====
function PluginSettingsForm({ pluginId, schema }: { pluginId: string; schema: import("../api/types").PluginSettingSchema[] }) {
  const queryClient = useQueryClient();
  const { data: configValues, isLoading } = useQuery({
    queryKey: ["plugin-config", pluginId],
    queryFn: () => pluginsApi.getConfig(pluginId),
  });
  const [localValues, setLocalValues] = useState<Record<string, unknown>>({});
  const [initialized, setInitialized] = useState(false);

  // Initialize local values once config loads
  if (configValues && !initialized) {
    setLocalValues(configValues);
    setInitialized(true);
  }

  const saveMut = useMutation({
    mutationFn: (values: Record<string, unknown>) => pluginsApi.setConfig(pluginId, values),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["plugin-config", pluginId] }),
  });

  const updateValue = (name: string, value: unknown) => {
    setLocalValues((prev) => ({ ...prev, [name]: value }));
  };

  const isDirty = JSON.stringify(localValues) !== JSON.stringify(configValues ?? {});

  if (isLoading) return <Loader2 className="w-4 h-4 animate-spin text-plex-text-secondary" />;

  return (
    <div>
      <div className="text-xs font-medium text-plex-text-secondary mb-2">Settings</div>
      <div className="space-y-2">
        {schema.map((s) => (
          <div key={s.name} className="flex items-center gap-3 bg-gray-900/50 rounded px-3 py-2">
            <label className="text-sm min-w-[140px] shrink-0">
              {s.displayName || s.name}
              {s.description && <div className="text-xs text-plex-text-muted mt-0.5">{s.description}</div>}
            </label>
            {s.type === "BOOLEAN" ? (
              <button
                onClick={() => updateValue(s.name, !localValues[s.name])}
                className={`px-3 py-1 text-xs rounded font-medium transition-colors ${
                  localValues[s.name]
                    ? "bg-green-600/20 text-green-400 hover:bg-green-600/30"
                    : "bg-gray-600/30 text-gray-400 hover:bg-gray-600/40"
                }`}
              >
                {localValues[s.name] ? "On" : "Off"}
              </button>
            ) : s.type === "NUMBER" ? (
              <input
                type="number"
                value={(localValues[s.name] as number) ?? ""}
                onChange={(e) => updateValue(s.name, e.target.value ? Number(e.target.value) : null)}
                className="flex-1 bg-gray-800 border border-gray-600 rounded px-2 py-1 text-sm focus:border-plex-accent outline-none"
              />
            ) : (
              <input
                type="text"
                value={(localValues[s.name] as string) ?? ""}
                onChange={(e) => updateValue(s.name, e.target.value || null)}
                className="flex-1 bg-gray-800 border border-gray-600 rounded px-2 py-1 text-sm focus:border-plex-accent outline-none"
              />
            )}
          </div>
        ))}
      </div>
      {isDirty && (
        <div className="flex justify-end mt-2 gap-2">
          <button
            onClick={() => { setLocalValues(configValues ?? {}); }}
            className="px-3 py-1 text-xs bg-gray-700 hover:bg-gray-600 rounded transition-colors"
          >
            Reset
          </button>
          <button
            onClick={() => saveMut.mutate(localValues)}
            disabled={saveMut.isPending}
            className="px-3 py-1 text-xs bg-plex-accent hover:bg-plex-accent-hover rounded transition-colors disabled:opacity-50"
          >
            {saveMut.isPending ? "Saving..." : "Save Settings"}
          </button>
        </div>
      )}
    </div>
  );
}

function PluginsPanel() {
  const queryClient = useQueryClient();
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const { data: pluginList, isLoading } = useQuery({
    queryKey: ["plugins"],
    queryFn: pluginsApi.list,
  });

  const reloadMut = useMutation({
    mutationFn: pluginsApi.reload,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["plugins"] }),
  });

  const settingsMut = useMutation({
    mutationFn: pluginsApi.saveSettings,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["plugins"] }),
  });

  const runTaskMut = useMutation({
    mutationFn: pluginsApi.runTask,
  });

  const togglePlugin = (id: string, enabled: boolean) => {
    settingsMut.mutate({ enabledMap: { [id]: enabled } });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-10">
        <Loader2 className="w-6 h-6 animate-spin text-plex-text-secondary" />
      </div>
    );
  }

  const sortedPlugins = [...(pluginList ?? [])].sort((a, b) => a.name.localeCompare(b.name));

  return (
    <>
      <SectionCard title="Installed Plugins" description="Manage extensions loaded into this instance.">
        <div className="flex justify-between items-center mb-4">
          <span className="text-sm text-plex-text-secondary">{sortedPlugins.length} plugin{sortedPlugins.length !== 1 ? "s" : ""} found</span>
          <button
            onClick={() => reloadMut.mutate()}
            disabled={reloadMut.isPending}
            className="flex items-center gap-1.5 px-3 py-1.5 text-sm bg-gray-700 hover:bg-gray-600 rounded transition-colors disabled:opacity-50"
          >
            <RefreshCw className={`w-3.5 h-3.5 ${reloadMut.isPending ? "animate-spin" : ""}`} />
            Reload Plugins
          </button>
        </div>

        {sortedPlugins.length === 0 && (
          <div className="text-sm text-plex-text-muted py-6 text-center">No plugins installed</div>
        )}

        <div className="space-y-2">
          {sortedPlugins.map((plugin) => {
            const isExpanded = expandedId === plugin.id;
            return (
              <div key={plugin.id} className="bg-gray-800/50 rounded-lg border border-gray-700/50">
                <div
                  className="flex items-center justify-between px-4 py-3 cursor-pointer hover:bg-gray-700/30 transition-colors"
                  onClick={() => setExpandedId(isExpanded ? null : plugin.id)}
                >
                  <div className="flex items-center gap-3 min-w-0">
                    <div className={`w-2 h-2 rounded-full ${plugin.enabled ? "bg-green-400" : "bg-gray-500"}`} />
                    <div className="min-w-0">
                      <div className="font-medium text-sm flex items-center gap-2">
                        {plugin.name}
                        <span className="text-xs text-plex-text-muted">v{plugin.version}</span>
                      </div>
                      {plugin.description && (
                        <div className="text-xs text-plex-text-secondary truncate">{plugin.description}</div>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center gap-3 shrink-0">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        togglePlugin(plugin.id, !plugin.enabled);
                      }}
                      className={`px-3 py-1 text-xs rounded font-medium transition-colors ${
                        plugin.enabled
                          ? "bg-green-600/20 text-green-400 hover:bg-green-600/30"
                          : "bg-gray-600/30 text-gray-400 hover:bg-gray-600/40"
                      }`}
                    >
                      {plugin.enabled ? "Enabled" : "Disabled"}
                    </button>
                    <span className="text-gray-500 text-xs">{isExpanded ? "▲" : "▼"}</span>
                  </div>
                </div>

                {isExpanded && (
                  <div className="px-4 pb-4 border-t border-gray-700/50 pt-3 space-y-3">
                    <div className="text-xs text-plex-text-muted">
                      <span className="font-medium">ID:</span> {plugin.id}
                      {plugin.url && (
                        <> · <a href={plugin.url} target="_blank" rel="noopener noreferrer" className="text-plex-accent hover:underline">{plugin.url}</a></>
                      )}
                    </div>

                    {plugin.settings && plugin.settings.length > 0 && (
                      <PluginSettingsForm pluginId={plugin.id} schema={plugin.settings} />
                    )}

                    {plugin.tasks.length > 0 && (
                      <div>
                        <div className="text-xs font-medium text-plex-text-secondary mb-2">Tasks</div>
                        <div className="space-y-1.5">
                          {plugin.tasks.map((task) => (
                            <div key={task.name} className="flex items-center justify-between bg-gray-900/50 rounded px-3 py-2">
                              <div>
                                <div className="text-sm font-medium">{task.name}</div>
                                {task.description && (
                                  <div className="text-xs text-plex-text-muted">{task.description}</div>
                                )}
                              </div>
                              <button
                                onClick={() => runTaskMut.mutate({ pluginId: plugin.id, taskName: task.name })}
                                disabled={runTaskMut.isPending}
                                className="px-2 py-1 text-xs bg-blue-600 hover:bg-blue-500 rounded transition-colors disabled:opacity-50"
                              >
                                Run
                              </button>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}

                    {plugin.tasks.length === 0 && (
                      <div className="text-xs text-plex-text-muted">No tasks defined for this plugin.</div>
                    )}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </SectionCard>

      <PackageManagerSection />
    </>
  );
}

// ===== Package Manager =====
function PackageManagerSection() {
  const queryClient = useQueryClient();
  const [pkgType, setPkgType] = useState<"plugin" | "scraper">("plugin");

  const { data: availablePackages, isLoading: availLoading } = useQuery({
    queryKey: ["available-packages", pkgType],
    queryFn: () => pluginsApi.availablePackages(pkgType),
  });

  const { data: installedPackages } = useQuery({
    queryKey: ["installed-packages", pkgType],
    queryFn: () => pluginsApi.installedPackages(pkgType),
  });

  const installMut = useMutation({
    mutationFn: (pkgs: { id: string; sourceUrl: string }[]) => pluginsApi.installPackages(pkgs),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["installed-packages"] });
      queryClient.invalidateQueries({ queryKey: ["available-packages"] });
      queryClient.invalidateQueries({ queryKey: ["plugins"] });
    },
  });

  const uninstallMut = useMutation({
    mutationFn: (ids: string[]) => pluginsApi.uninstallPackages(ids),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["installed-packages"] });
      queryClient.invalidateQueries({ queryKey: ["available-packages"] });
      queryClient.invalidateQueries({ queryKey: ["plugins"] });
    },
  });

  const installedIds = new Set((installedPackages ?? []).map(p => p.name));

  return (
    <SectionCard title="Package Manager" description="Browse and install plugins and scrapers from configured package sources.">
      <div className="flex gap-2 mb-4">
        <button
          onClick={() => setPkgType("plugin")}
          className={`px-3 py-1.5 text-sm rounded transition-colors ${pkgType === "plugin" ? "bg-plex-accent text-white" : "bg-plex-card border border-plex-border text-plex-text-secondary hover:text-plex-text"}`}
        >
          <Plug className="w-3.5 h-3.5 inline mr-1.5" />
          Plugins
        </button>
        <button
          onClick={() => setPkgType("scraper")}
          className={`px-3 py-1.5 text-sm rounded transition-colors ${pkgType === "scraper" ? "bg-plex-accent text-white" : "bg-plex-card border border-plex-border text-plex-text-secondary hover:text-plex-text"}`}
        >
          <SearchCode className="w-3.5 h-3.5 inline mr-1.5" />
          Scrapers
        </button>
      </div>

      {availLoading ? (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="w-5 h-5 animate-spin text-plex-text-muted" />
        </div>
      ) : !availablePackages || availablePackages.length === 0 ? (
        <div className="text-sm text-plex-text-muted text-center py-6">
          No available packages. Add package sources in Settings &gt; Library.
        </div>
      ) : (
        <div className="space-y-2">
          {availablePackages.map((pkg) => {
            const isInstalled = installedIds.has(pkg.name) || pkg.installed;
            const hasUpdate = isInstalled && pkg.installedVersion && pkg.installedVersion !== pkg.version;
            return (
              <div key={`${pkg.name}-${pkg.sourceUrl}`} className="flex items-center justify-between bg-plex-card border border-plex-border rounded-xl px-4 py-3">
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-medium text-plex-text">{pkg.name}</span>
                    <span className="text-xs text-plex-text-muted">v{pkg.version}</span>
                    {isInstalled && (
                      <span className="text-xs px-1.5 py-0.5 rounded bg-green-600/20 text-green-400">Installed</span>
                    )}
                    {hasUpdate && (
                      <span className="text-xs px-1.5 py-0.5 rounded bg-yellow-600/20 text-yellow-400">
                        Update available (v{pkg.installedVersion} → v{pkg.version})
                      </span>
                    )}
                  </div>
                  {pkg.description && <p className="text-xs text-plex-text-secondary mt-0.5 truncate">{pkg.description}</p>}
                </div>
                <div className="flex gap-2 ml-3 flex-shrink-0">
                  {!isInstalled ? (
                    <button
                      onClick={() => installMut.mutate([{ id: pkg.name, sourceUrl: pkg.sourceUrl }])}
                      disabled={installMut.isPending}
                      className="px-3 py-1.5 text-xs bg-plex-accent hover:bg-plex-accent-hover text-white rounded disabled:opacity-50 flex items-center gap-1"
                    >
                      {installMut.isPending ? <Loader2 className="w-3 h-3 animate-spin" /> : <Download className="w-3 h-3" />}
                      Install
                    </button>
                  ) : hasUpdate ? (
                    <button
                      onClick={() => installMut.mutate([{ id: pkg.name, sourceUrl: pkg.sourceUrl }])}
                      disabled={installMut.isPending}
                      className="px-3 py-1.5 text-xs bg-yellow-600 hover:bg-yellow-500 text-white rounded disabled:opacity-50 flex items-center gap-1"
                    >
                      {installMut.isPending ? <Loader2 className="w-3 h-3 animate-spin" /> : <RefreshCw className="w-3 h-3" />}
                      Update
                    </button>
                  ) : (
                    <button
                      onClick={() => uninstallMut.mutate([pkg.name])}
                      disabled={uninstallMut.isPending}
                      className="px-3 py-1.5 text-xs bg-plex-card border border-plex-border text-plex-text-muted hover:text-red-400 hover:border-red-500 rounded disabled:opacity-50 flex items-center gap-1"
                    >
                      {uninstallMut.isPending ? <Loader2 className="w-3 h-3 animate-spin" /> : <Trash2 className="w-3 h-3" />}
                      Uninstall
                    </button>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </SectionCard>
  );
}
