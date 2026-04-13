import type {
  Scene, SceneCreate, SceneUpdate,
  Performer, PerformerCreate, PerformerUpdate,
  Tag, TagDetail, TagCreate, TagUpdate,
  Studio, StudioCreate, StudioUpdate,
  Gallery, GalleryCreate, GalleryUpdate, GalleryChapter, GalleryChapterCreate, GalleryChapterUpdate,
  Image, ImageCreate, ImageUpdate,
  Group, GroupCreate, GroupUpdate,
  SceneMarkerSummary, SceneMarkerCreate, SceneMarkerUpdate,
  SceneMarkerWall,
  PaginatedResponse, Stats, SystemStatus, StashConfig, JobInfo,
  ScraperSummary,
  DlnaStatus,
  StashBox,
  StashBoxPerformerImportRequest,
  StashBoxPerformerMatch,
  StashBoxSceneImportRequest,
  StashBoxSceneMatch,
  StashBoxStudioMatch,
  StashBoxStudioImportRequest,
  StashBoxValidationResult,
  FindFilter,
  SavedFilter,
  SavedFilterCreate,
  SavedFilterUpdate,
  FilteredQueryRequest,
  SceneFilterCriteria,
  PerformerFilterCriteria,
  TagFilterCriteria,
  StudioFilterCriteria,
  GalleryFilterCriteria,
  ImageFilterCriteria,
  GroupFilterCriteria,
  BulkSceneUpdate,
  BulkPerformerUpdate,
  BulkTagUpdate,
  BulkStudioUpdate,
  BulkGalleryUpdate,
  BulkImageUpdate,
  BulkGroupUpdate,
  Plugin,
  PluginTask,
  RunPluginTaskRequest,
  PluginSettings,
  Package,
} from "./types";

const API_BASE = "/api";

const CRITERION_MODIFIER_MAP: Record<string, string> = {
  EQUALS: "equals",
  NOT_EQUALS: "notEquals",
  GREATER_THAN: "greaterThan",
  LESS_THAN: "lessThan",
  INCLUDES: "includes",
  EXCLUDES: "excludes",
  INCLUDES_ALL: "includesAll",
  EXCLUDES_ALL: "excludesAll",
  IS_NULL: "isNull",
  NOT_NULL: "notNull",
  BETWEEN: "between",
  NOT_BETWEEN: "notBetween",
  MATCHES_REGEX: "matchesRegex",
  NOT_MATCHES_REGEX: "notMatchesRegex",
};

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`API Error ${res.status}: ${text}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}

function buildQuery(filter?: FindFilter, extra?: Record<string, string | number | boolean | undefined>): string {
  const params = new URLSearchParams();
  if (filter?.q) params.set("q", filter.q);
  if (filter?.page) params.set("page", String(filter.page));
  if (filter?.perPage) params.set("perPage", String(filter.perPage));
  if (filter?.sort) params.set("sort", filter.sort);
  if (filter?.direction) params.set("direction", filter.direction);
  if (extra) {
    for (const [k, v] of Object.entries(extra)) {
      if (v !== undefined) params.set(k, String(v));
    }
  }
  const qs = params.toString();
  return qs ? `?${qs}` : "";
}

function normalizeCriterionPayload<T>(value: T): T {
  if (Array.isArray(value)) {
    return value.map((item) => normalizeCriterionPayload(item)) as T;
  }

  if (value && typeof value === "object") {
    const normalizedEntries = Object.entries(value as Record<string, unknown>).map(([key, entryValue]) => {
      if (key === "modifier" && typeof entryValue === "string") {
        return [key, CRITERION_MODIFIER_MAP[entryValue] ?? entryValue];
      }

      return [key, normalizeCriterionPayload(entryValue)];
    });

    return Object.fromEntries(normalizedEntries) as T;
  }

  return value;
}

// ===== Scenes =====
export const scenes = {
  find: (filter?: FindFilter, extra?: Record<string, string | number | boolean | undefined>) =>
    request<PaginatedResponse<Scene>>(`/scenes${buildQuery(filter, extra)}`),
  findFiltered: (req: FilteredQueryRequest<SceneFilterCriteria>) =>
    request<PaginatedResponse<Scene>>("/scenes/find", { method: "POST", body: JSON.stringify(normalizeCriterionPayload(req)) }),
  get: (id: number) => request<Scene>(`/scenes/${id}`),
  create: (data: SceneCreate) => request<Scene>("/scenes", { method: "POST", body: JSON.stringify(data) }),
  update: (id: number, data: SceneUpdate) => request<Scene>(`/scenes/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  bulkUpdate: (data: BulkSceneUpdate) => request<void>("/scenes/bulk", { method: "PUT", body: JSON.stringify(data) }),
  delete: (id: number) => request<void>(`/scenes/${id}`, { method: "DELETE" }),
  bulkDelete: (ids: number[]) => request<void>("/scenes/bulk", { method: "DELETE", body: JSON.stringify({ ids }) }),
  merge: (targetId: number, sourceIds: number[]) =>
    request<Scene>("/scenes/merge", { method: "POST", body: JSON.stringify({ targetId, sourceIds }) }),
  recordPlay: (id: number) => request<void>(`/scenes/${id}/play`, { method: "POST" }),
  incrementO: (id: number) => request<void>(`/scenes/${id}/o`, { method: "POST" }),
  decrementO: (id: number) => request<void>(`/scenes/${id}/o`, { method: "DELETE" }),
  resetO: (id: number) => request<void>(`/scenes/${id}/o/reset`, { method: "POST" }),
  saveActivity: (id: number, data: { resumeTime?: number; playDuration?: number }) =>
    request<void>(`/scenes/${id}/activity`, { method: "POST", body: JSON.stringify(data) }),
  searchStashBox: (id: number, term?: string, endpoint?: string) =>
    request<StashBoxSceneMatch[]>(`/scenes/${id}/stash-box/search${buildQuery(undefined, { term, endpoint })}`),
  importFromStashBox: (id: number, data: StashBoxSceneImportRequest) =>
    request<Scene>(`/scenes/${id}/stash-box/import`, { method: "POST", body: JSON.stringify(data) }),
  streamUrl: (id: number) => `${API_BASE}/stream/scene/${id}`,
  screenshotUrl: (id: number) => `${API_BASE}/stream/scene/${id}/screenshot`,
  previewUrl: (id: number) => `${API_BASE}/stream/scene/${id}/preview`,
  spriteUrl: (id: number) => `${API_BASE}/stream/scene/${id}/sprite`,
  spriteVttUrl: (id: number) => `${API_BASE}/stream/scene/${id}/vtt/thumbs`,
  markers: {
    list: (sceneId: number) => request<SceneMarkerSummary[]>(`/scenes/${sceneId}/markers`),
    create: (sceneId: number, data: SceneMarkerCreate) =>
      request<SceneMarkerSummary>(`/scenes/${sceneId}/markers`, { method: "POST", body: JSON.stringify(data) }),
    update: (sceneId: number, id: number, data: SceneMarkerUpdate) =>
      request<SceneMarkerSummary>(`/scenes/${sceneId}/markers/${id}`, { method: "PUT", body: JSON.stringify(data) }),
    delete: (sceneId: number, id: number) =>
      request<void>(`/scenes/${sceneId}/markers/${id}`, { method: "DELETE" }),
  },
  findDuplicates: (distance = 0, durationDiff?: number) => {
    const params = new URLSearchParams();
    params.set("distance", String(distance));
    if (durationDiff !== undefined) params.set("durationDiff", String(durationDiff));
    return request<Scene[][]>(`/scenes/duplicates?${params.toString()}`);
  },
};

// ===== Markers (top-level) =====
export const markers = {
  wall: (opts?: { q?: string; tagId?: number; count?: number }) => {
    const params = new URLSearchParams();
    if (opts?.q) params.set("q", opts.q);
    if (opts?.tagId) params.set("tagId", String(opts.tagId));
    if (opts?.count) params.set("count", String(opts.count));
    const qs = params.toString();
    return request<SceneMarkerWall[]>(`/markers/wall${qs ? `?${qs}` : ""}`);
  },
};

// ===== Performers =====
export const performers = {
  find: (filter?: FindFilter, extra?: Record<string, string | number | boolean | undefined>) =>
    request<PaginatedResponse<Performer>>(`/performers${buildQuery(filter, extra)}`),
  findFiltered: (req: FilteredQueryRequest<PerformerFilterCriteria>) =>
    request<PaginatedResponse<Performer>>("/performers/find", { method: "POST", body: JSON.stringify(normalizeCriterionPayload(req)) }),
  get: (id: number) => request<Performer>(`/performers/${id}`),
  create: (data: PerformerCreate) => request<Performer>("/performers", { method: "POST", body: JSON.stringify(data) }),
  update: (id: number, data: PerformerUpdate) => request<Performer>(`/performers/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  bulkUpdate: (data: BulkPerformerUpdate) => request<void>("/performers/bulk", { method: "PUT", body: JSON.stringify(data) }),
  delete: (id: number) => request<void>(`/performers/${id}`, { method: "DELETE" }),
  bulkDelete: (ids: number[]) => request<void>("/performers/bulk", { method: "DELETE", body: JSON.stringify({ ids }) }),
  merge: (targetId: number, sourceIds: number[]) =>
    request<Performer>("/performers/merge", { method: "POST", body: JSON.stringify({ targetId, sourceIds }) }),
  searchStashBox: (id: number, term?: string, endpoint?: string) =>
    request<StashBoxPerformerMatch[]>(`/performers/${id}/stash-box/search${buildQuery(undefined, { term, endpoint })}`),
  importFromStashBox: (id: number, data: StashBoxPerformerImportRequest) =>
    request<Performer>(`/performers/${id}/stash-box/import`, { method: "POST", body: JSON.stringify(data) }),
};

// ===== Tags =====
export const tags = {
  find: (filter?: FindFilter, extra?: Record<string, string | number | boolean | undefined>) =>
    request<PaginatedResponse<Tag>>(`/tags${buildQuery(filter, extra)}`),
  findFiltered: (req: FilteredQueryRequest<TagFilterCriteria>) =>
    request<PaginatedResponse<Tag>>("/tags/find", { method: "POST", body: JSON.stringify(normalizeCriterionPayload(req)) }),
  get: (id: number) => request<TagDetail>(`/tags/${id}`),
  create: (data: TagCreate) => request<TagDetail>("/tags", { method: "POST", body: JSON.stringify(data) }),
  update: (id: number, data: TagUpdate) => request<TagDetail>(`/tags/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  bulkUpdate: (data: BulkTagUpdate) => request<void>("/tags/bulk", { method: "PUT", body: JSON.stringify(data) }),
  delete: (id: number) => request<void>(`/tags/${id}`, { method: "DELETE" }),
  bulkDelete: (ids: number[]) => request<void>("/tags/bulk", { method: "DELETE", body: JSON.stringify({ ids }) }),
  merge: (targetId: number, sourceIds: number[]) =>
    request<TagDetail>("/tags/merge", { method: "POST", body: JSON.stringify({ targetId, sourceIds }) }),
};

// ===== Studios =====
export const studios = {
  find: (filter?: FindFilter, extra?: Record<string, string | number | boolean | undefined>) =>
    request<PaginatedResponse<Studio>>(`/studios${buildQuery(filter, extra)}`),
  findFiltered: (req: FilteredQueryRequest<StudioFilterCriteria>) =>
    request<PaginatedResponse<Studio>>("/studios/find", { method: "POST", body: JSON.stringify(normalizeCriterionPayload(req)) }),
  get: (id: number) => request<Studio>(`/studios/${id}`),
  create: (data: StudioCreate) => request<Studio>("/studios", { method: "POST", body: JSON.stringify(data) }),
  update: (id: number, data: StudioUpdate) => request<Studio>(`/studios/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  bulkUpdate: (data: BulkStudioUpdate) => request<void>("/studios/bulk", { method: "PUT", body: JSON.stringify(data) }),
  delete: (id: number) => request<void>(`/studios/${id}`, { method: "DELETE" }),
  bulkDelete: (ids: number[]) => request<void>("/studios/bulk", { method: "DELETE", body: JSON.stringify({ ids }) }),
  merge: (targetId: number, sourceIds: number[]) =>
    request<Studio>("/studios/merge", { method: "POST", body: JSON.stringify({ targetId, sourceIds }) }),
  searchStashBox: (id: number, term?: string, endpoint?: string) => {
    const params = new URLSearchParams();
    if (term) params.set("term", term);
    if (endpoint) params.set("endpoint", endpoint);
    const qs = params.toString();
    return request<StashBoxStudioMatch[]>(`/studios/${id}/stash-box/search${qs ? `?${qs}` : ""}`);
  },
  importFromStashBox: (id: number, data: StashBoxStudioImportRequest) =>
    request<Studio>(`/studios/${id}/stash-box/import`, { method: "POST", body: JSON.stringify(data) }),
};

// ===== Galleries =====
export const galleries = {
  find: (filter?: FindFilter, extra?: Record<string, string | number | boolean | undefined>) =>
    request<PaginatedResponse<Gallery>>(`/galleries${buildQuery(filter, extra)}`),
  findFiltered: (req: FilteredQueryRequest<GalleryFilterCriteria>) =>
    request<PaginatedResponse<Gallery>>("/galleries/find", { method: "POST", body: JSON.stringify(normalizeCriterionPayload(req)) }),
  get: (id: number) => request<Gallery>(`/galleries/${id}`),
  create: (data: GalleryCreate) => request<Gallery>("/galleries", { method: "POST", body: JSON.stringify(data) }),
  update: (id: number, data: GalleryUpdate) => request<Gallery>(`/galleries/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  bulkUpdate: (data: BulkGalleryUpdate) => request<void>("/galleries/bulk", { method: "PUT", body: JSON.stringify(data) }),
  delete: (id: number) => request<void>(`/galleries/${id}`, { method: "DELETE" }),
  bulkDelete: (ids: number[]) => request<void>("/galleries/bulk", { method: "DELETE", body: JSON.stringify({ ids }) }),
  chapters: (id: number) => request<GalleryChapter[]>(`/galleries/${id}/chapters`),
  createChapter: (id: number, data: GalleryChapterCreate) =>
    request<GalleryChapter>(`/galleries/${id}/chapters`, { method: "POST", body: JSON.stringify(data) }),
  updateChapter: (galleryId: number, chapterId: number, data: GalleryChapterUpdate) =>
    request<GalleryChapter>(`/galleries/${galleryId}/chapters/${chapterId}`, { method: "PUT", body: JSON.stringify(data) }),
  deleteChapter: (galleryId: number, chapterId: number) =>
    request<void>(`/galleries/${galleryId}/chapters/${chapterId}`, { method: "DELETE" }),
  addImages: (id: number, imageIds: number[]) =>
    request<{ added: number }>(`/galleries/${id}/images`, { method: "POST", body: JSON.stringify({ imageIds }) }),
  removeImages: (id: number, imageIds: number[]) =>
    request<{ removed: number }>(`/galleries/${id}/images`, { method: "DELETE", body: JSON.stringify({ imageIds }) }),
};

// ===== Images =====
export const images = {
  find: (filter?: FindFilter, extra?: Record<string, string | number | boolean | undefined>) =>
    request<PaginatedResponse<Image>>(`/images${buildQuery(filter, extra)}`),
  findFiltered: (req: FilteredQueryRequest<ImageFilterCriteria>) =>
    request<PaginatedResponse<Image>>("/images/find", { method: "POST", body: JSON.stringify(normalizeCriterionPayload(req)) }),
  get: (id: number) => request<Image>(`/images/${id}`),
  create: (data: ImageCreate) => request<Image>("/images", { method: "POST", body: JSON.stringify(data) }),
  update: (id: number, data: ImageUpdate) => request<Image>(`/images/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  bulkUpdate: (data: BulkImageUpdate) => request<void>("/images/bulk", { method: "PUT", body: JSON.stringify(data) }),
  delete: (id: number) => request<void>(`/images/${id}`, { method: "DELETE" }),
  bulkDelete: (ids: number[]) => request<void>("/images/bulk", { method: "DELETE", body: JSON.stringify({ ids }) }),
  incrementO: (id: number) => request<void>(`/images/${id}/o`, { method: "POST" }),
  decrementO: (id: number) => request<void>(`/images/${id}/o`, { method: "DELETE" }),
  resetO: (id: number) => request<void>(`/images/${id}/o/reset`, { method: "POST" }),
  imageUrl: (id: number) => `${API_BASE}/stream/image/${id}`,
  thumbnailUrl: (id: number) => `${API_BASE}/stream/image/${id}/thumbnail`,
};

// ===== Groups =====
export const groups = {
  find: (filter?: FindFilter, extra?: Record<string, string | number | boolean | undefined>) =>
    request<PaginatedResponse<Group>>(`/groups${buildQuery(filter, extra)}`),
  findFiltered: (req: FilteredQueryRequest<GroupFilterCriteria>) =>
    request<PaginatedResponse<Group>>("/groups/find", { method: "POST", body: JSON.stringify(normalizeCriterionPayload(req)) }),
  get: (id: number) => request<Group>(`/groups/${id}`),
  create: (data: GroupCreate) => request<Group>("/groups", { method: "POST", body: JSON.stringify(data) }),
  update: (id: number, data: GroupUpdate) => request<Group>(`/groups/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  bulkUpdate: (data: BulkGroupUpdate) => request<void>("/groups/bulk", { method: "PUT", body: JSON.stringify(data) }),
  delete: (id: number) => request<void>(`/groups/${id}`, { method: "DELETE" }),
  bulkDelete: (ids: number[]) => request<void>("/groups/bulk", { method: "DELETE", body: JSON.stringify({ ids }) }),
  subGroups: (id: number) => request<Group[]>(`/groups/${id}/subgroups`),
  containingGroups: (id: number) => request<Group[]>(`/groups/${id}/containinggroups`),
  addSubGroup: (id: number, subGroupId: number, orderIndex?: number) =>
    request<void>(`/groups/${id}/subgroups`, { method: "POST", body: JSON.stringify({ subGroupId, orderIndex }) }),
  removeSubGroup: (id: number, subGroupId: number) =>
    request<void>(`/groups/${id}/subgroups/${subGroupId}`, { method: "DELETE" }),
  reorderSubGroups: (id: number, subGroupIds: number[]) =>
    request<void>(`/groups/${id}/subgroups/reorder`, { method: "PUT", body: JSON.stringify({ subGroupIds }) }),
};

// ===== Entity Images =====
async function uploadImage(path: string, file: File): Promise<{ blobId: string }> {
  const formData = new FormData();
  formData.append("file", file);
  const res = await fetch(`${API_BASE}${path}`, { method: "POST", body: formData });
  if (!res.ok) throw new Error(`Upload failed: ${res.status}`);
  return res.json();
}

async function deleteImage(path: string): Promise<void> {
  const res = await fetch(`${API_BASE}${path}`, { method: "DELETE" });
  if (!res.ok && res.status !== 404) throw new Error(`Delete failed: ${res.status}`);
}

export const entityImages = {
  performerImageUrl: (id: number) => `${API_BASE}/performers/${id}/image`,
  uploadPerformerImage: (id: number, file: File) => uploadImage(`/performers/${id}/image`, file),
  deletePerformerImage: (id: number) => deleteImage(`/performers/${id}/image`),

  studioImageUrl: (id: number) => `${API_BASE}/studios/${id}/image`,
  uploadStudioImage: (id: number, file: File) => uploadImage(`/studios/${id}/image`, file),
  deleteStudioImage: (id: number) => deleteImage(`/studios/${id}/image`),

  tagImageUrl: (id: number) => `${API_BASE}/tags/${id}/image`,
  uploadTagImage: (id: number, file: File) => uploadImage(`/tags/${id}/image`, file),
  deleteTagImage: (id: number) => deleteImage(`/tags/${id}/image`),

  groupFrontImageUrl: (id: number) => `${API_BASE}/groups/${id}/image/front`,
  uploadGroupFrontImage: (id: number, file: File) => uploadImage(`/groups/${id}/image/front`, file),
  deleteGroupFrontImage: (id: number) => deleteImage(`/groups/${id}/image/front`),

  groupBackImageUrl: (id: number) => `${API_BASE}/groups/${id}/image/back`,
  uploadGroupBackImage: (id: number, file: File) => uploadImage(`/groups/${id}/image/back`, file),
  deleteGroupBackImage: (id: number) => deleteImage(`/groups/${id}/image/back`),
};

// ===== System =====
export const system = {
  status: () => request<SystemStatus>("/system/status"),
  stats: () => request<Stats>("/system/stats"),
  getConfig: () => request<StashConfig>("/system/config"),
  saveConfig: (config: StashConfig) =>
    request<StashConfig>("/system/config", { method: "PUT", body: JSON.stringify(config) }),
  listScrapers: () => request<ScraperSummary[]>("/system/scrapers"),
  reloadScrapers: () => request<ScraperSummary[]>("/system/scrapers/reload", { method: "POST" }),
  validateStashBox: (stashBox: StashBox) =>
    request<StashBoxValidationResult>("/system/stash-boxes/validate", { method: "POST", body: JSON.stringify(stashBox) }),
};

// ===== DLNA =====
export const dlna = {
  status: () => request<DlnaStatus>("/dlna/status"),
  enable: (durationMinutes?: number) =>
    request<DlnaStatus>("/dlna/enable", { method: "POST", body: JSON.stringify({ durationMinutes }) }),
  disable: () => request<DlnaStatus>("/dlna/disable", { method: "POST" }),
  allowIp: (ipAddress: string, durationMinutes?: number) =>
    request<DlnaStatus>("/dlna/allow-ip", { method: "POST", body: JSON.stringify({ ipAddress, durationMinutes }) }),
  removeIp: (ipAddress: string) =>
    request<DlnaStatus>("/dlna/remove-ip", { method: "POST", body: JSON.stringify({ ipAddress }) }),
};

// ===== Jobs =====
export const jobs = {
  list: () => request<JobInfo[]>("/jobs"),
  history: () => request<JobInfo[]>("/jobs/history"),
  get: (id: string) => request<JobInfo>(`/jobs/${id}`),
  cancel: (id: string) => request<void>(`/jobs/${id}`, { method: "DELETE" }),
  scan: (opts?: { generatePreviews?: boolean }) =>
    request<{ jobId: string }>(`/jobs/scan${opts?.generatePreviews ? "?generatePreviews=true" : ""}`, { method: "POST" }),
  generateThumbnails: () => request<{ jobId: string }>("/jobs/generate-thumbnails", { method: "POST" }),
  generateScenePhashes: () => request<{ jobId: string }>("/jobs/generate-scene-phashes", { method: "POST" }),
  generateImagePhashes: () => request<{ jobId: string }>("/jobs/generate-image-phashes", { method: "POST" }),
  autoTag: (opts?: { performerIds?: string[]; studioIds?: string[]; tagIds?: string[] }) =>
    request<{ jobId: string }>("/jobs/auto-tag", { method: "POST", body: JSON.stringify(opts ?? {}) }),
  clean: (dryRun = false) =>
    request<{ jobId: string }>(`/jobs/clean${dryRun ? "?dryRun=true" : ""}`, { method: "POST" }),
  backup: () => request<{ jobId: string }>("/jobs/backup", { method: "POST" }),
  latestBackup: () => request<{ path: string }>("/jobs/backup/latest"),
};

// ===== Metadata Tasks =====
export interface ScanOptions {
  paths?: string[];
  scanGenerateCovers?: boolean;
  scanGeneratePreviews?: boolean;
  scanGenerateSprites?: boolean;
  scanGeneratePhashes?: boolean;
  scanGenerateThumbnails?: boolean;
  scanGenerateImagePhashes?: boolean;
  rescan?: boolean;
}

export interface GenerateOptions {
  thumbnails?: boolean;
  previews?: boolean;
  sprites?: boolean;
  markers?: boolean;
  phashes?: boolean;
  imageThumbnails?: boolean;
  imagePhashes?: boolean;
  overwrite?: boolean;
  sceneIds?: number[];
}

export interface CleanOptions {
  paths?: string[];
  dryRun?: boolean;
}

export interface CleanGeneratedOptions {
  screenshots?: boolean;
  sprites?: boolean;
  transcodes?: boolean;
  markers?: boolean;
  imageThumbnails?: boolean;
  dryRun?: boolean;
}

export interface ExportOptions {
  includeScenes?: boolean;
  includePerformers?: boolean;
  includeStudios?: boolean;
  includeTags?: boolean;
  includeGalleries?: boolean;
  includeGroups?: boolean;
}

export const metadata = {
  scan: (opts?: ScanOptions) =>
    request<{ jobId: string }>("/metadata/scan", { method: "POST", body: JSON.stringify(opts ?? {}) }),
  generate: (opts?: GenerateOptions) =>
    request<{ jobId: string }>("/metadata/generate", { method: "POST", body: JSON.stringify(opts ?? {}) }),
  autoTag: (opts?: { performers?: string[]; studios?: string[]; tags?: string[] }) =>
    request<{ jobId: string }>("/metadata/auto-tag", { method: "POST", body: JSON.stringify(opts ?? {}) }),
  clean: (opts?: CleanOptions) =>
    request<{ jobId: string }>("/metadata/clean", { method: "POST", body: JSON.stringify(opts ?? {}) }),
  cleanGenerated: (opts?: CleanGeneratedOptions) =>
    request<{ jobId: string }>("/metadata/clean-generated", { method: "POST", body: JSON.stringify(opts ?? {}) }),
  export: (opts?: ExportOptions) =>
    request<{ jobId: string }>("/metadata/export", { method: "POST", body: JSON.stringify(opts ?? {}) }),
  identify: (opts?: { sources?: string[]; sceneIds?: number[] }) =>
    request<{ jobId: string }>("/metadata/identify", { method: "POST", body: JSON.stringify(opts ?? {}) }),
};

// ===== Database =====
export const database = {
  backup: () => request<{ backupPath: string; sizeBytes: number; timestamp: string }>("/database/backup", { method: "POST" }),
  optimize: () => request<void>("/database/optimize", { method: "POST" }),
};

// ===== Logs =====
export interface LogEntry {
  timestamp: string;
  level: string;
  message: string;
  exception?: string;
}

export const logs = {
  recent: (level?: string, limit?: number) => {
    const params = new URLSearchParams();
    if (level) params.set("level", level);
    if (limit) params.set("limit", String(limit));
    const qs = params.toString();
    return request<LogEntry[]>(`/logs${qs ? `?${qs}` : ""}`);
  },
};

// ===== Saved Filters =====
export const savedFilters = {
  list: (mode?: string) => request<SavedFilter[]>(`/savedfilters${mode ? `?mode=${mode}` : ""}`),
  get: (id: number) => request<SavedFilter>(`/savedfilters/${id}`),
  create: (data: SavedFilterCreate) => request<SavedFilter>("/savedfilters", { method: "POST", body: JSON.stringify(data) }),
  update: (id: number, data: SavedFilterUpdate) => request<SavedFilter>(`/savedfilters/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  delete: (id: number) => request<void>(`/savedfilters/${id}`, { method: "DELETE" }),
};

// ===== Plugins =====
export const plugins = {
  list: () => request<Plugin[]>("/plugins"),
  getTasks: () => request<PluginTask[]>("/plugins/tasks"),
  runTask: (data: RunPluginTaskRequest) => request<{ jobId: string }>("/plugins/run-task", { method: "POST", body: JSON.stringify(data) }),
  saveSettings: (data: PluginSettings) => request<void>("/plugins/settings", { method: "POST", body: JSON.stringify(data) }),
  reload: () => request<{ message: string }>("/plugins/reload", { method: "POST" }),
  getConfig: (pluginId: string) => request<Record<string, unknown>>(`/plugins/${encodeURIComponent(pluginId)}/config`),
  setConfig: (pluginId: string, values: Record<string, unknown>) =>
    request<void>(`/plugins/${encodeURIComponent(pluginId)}/config`, { method: "POST", body: JSON.stringify(values) }),
  installedPackages: (type?: string) => request<Package[]>(`/plugins/packages/installed${type ? `?type=${type}` : ""}`),
  availablePackages: (type?: string, source?: string) => {
    const params = new URLSearchParams();
    if (type) params.set("type", type);
    if (source) params.set("source", source);
    const qs = params.toString();
    return request<Package[]>(`/plugins/packages/available${qs ? `?${qs}` : ""}`);
  },
  installPackages: (packages: { id: string; sourceUrl: string }[]) =>
    request<{ jobId: string }>("/plugins/packages/install", { method: "POST", body: JSON.stringify({ packages }) }),
  uninstallPackages: (ids: string[]) =>
    request<{ uninstalled: string[] }>("/plugins/packages/uninstall", { method: "POST", body: JSON.stringify(ids) }),
};
