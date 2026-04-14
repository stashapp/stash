// ===== Entity Types =====

export interface Scene {
  id: number;
  title?: string;
  code?: string;
  details?: string;
  director?: string;
  date?: string;
  rating?: number;
  organized: boolean;
  studioId?: number;
  studioName?: string;
  resumeTime: number;
  playDuration: number;
  playCount: number;
  lastPlayedAt?: string;
  oCounter: number;
  urls: string[];
  tags: Tag[];
  performers: PerformerSummary[];
  files: VideoFile[];
  markers: SceneMarkerSummary[];
  groups: GroupSummary[];
  galleries: GallerySummary[];
  stashIds: SceneStashId[];
  customFields?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface SceneStashId {
  endpoint: string;
  stashId: string;
}

export interface SceneCreate {
  title?: string;
  code?: string;
  details?: string;
  director?: string;
  date?: string;
  rating?: number;
  organized?: boolean;
  studioId?: number;
  urls?: string[];
  tagIds?: number[];
  performerIds?: number[];
  galleryIds?: number[];
  groups?: { groupId: number; sceneIndex: number }[];
  customFields?: Record<string, unknown>;
}

export interface SceneUpdate extends Partial<SceneCreate> {}

export interface Performer {
  id: number;
  name: string;
  imagePath?: string;
  disambiguation?: string;
  gender?: string;
  birthdate?: string;
  deathDate?: string;
  ethnicity?: string;
  country?: string;
  eyeColor?: string;
  hairColor?: string;
  heightCm?: number;
  weight?: number;
  measurements?: string;
  fakeTits?: string;
  penisLength?: number;
  circumcised?: string;
  careerStart?: string;
  careerEnd?: string;
  tattoos?: string;
  piercings?: string;
  favorite: boolean;
  rating?: number;
  details?: string;
  ignoreAutoTag: boolean;
  urls: string[];
  aliases: string[];
  tags: Tag[];
  stashIds: PerformerStashId[];
  sceneCount: number;
  imageCount: number;
  galleryCount: number;
  groupCount: number;
  customFields?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface PerformerStashId {
  endpoint: string;
  stashId: string;
}

export interface PerformerSummary {
  id: number;
  name: string;
  disambiguation?: string;
  gender?: string;
  favorite: boolean;
  imagePath?: string;
}

export interface PerformerCreate {
  name: string;
  disambiguation?: string;
  gender?: string;
  birthdate?: string;
  deathDate?: string;
  ethnicity?: string;
  country?: string;
  eyeColor?: string;
  hairColor?: string;
  heightCm?: number;
  weight?: number;
  measurements?: string;
  fakeTits?: string;
  penisLength?: number;
  circumcised?: string;
  careerStart?: string;
  careerEnd?: string;
  tattoos?: string;
  piercings?: string;
  favorite?: boolean;
  rating?: number;
  details?: string;
  ignoreAutoTag?: boolean;
  urls?: string[];
  aliases?: string[];
  tagIds?: number[];
  customFields?: Record<string, unknown>;
}

export interface PerformerUpdate extends Partial<PerformerCreate> {}

export interface Tag {
  id: number;
  name: string;
  description?: string;
  imagePath?: string;
  favorite: boolean;
  ignoreAutoTag: boolean;
  aliases: string[];
  sceneCount?: number;
  sceneMarkerCount?: number;
  imageCount?: number;
  galleryCount?: number;
  groupCount?: number;
  performerCount?: number;
  studioCount?: number;
}

export interface TagDetail extends Tag {
  sortName?: string;
  parents: Tag[];
  children: Tag[];
  sceneCount: number;
  performerCount: number;
  imageCount: number;
  galleryCount: number;
  studioCount: number;
  groupCount: number;
  markerCount: number;
  customFields?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface TagCreate {
  name: string;
  sortName?: string;
  description?: string;
  favorite?: boolean;
  ignoreAutoTag?: boolean;
  aliases?: string[];
  parentIds?: number[];
  childIds?: number[];
  customFields?: Record<string, unknown>;
}

export interface TagUpdate extends Partial<TagCreate> {}

export interface Studio {
  id: number;
  name: string;
  imagePath?: string;
  parentId?: number;
  parentName?: string;
  rating?: number;
  favorite: boolean;
  details?: string;
  ignoreAutoTag: boolean;
  organized: boolean;
  urls: string[];
  aliases: string[];
  tags: Tag[];
  stashIds: StudioStashId[];
  sceneCount: number;
  imageCount: number;
  galleryCount: number;
  groupCount: number;
  performerCount: number;
  childStudioCount: number;
  customFields?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface StudioStashId {
  endpoint: string;
  stashId: string;
}

export interface StudioCreate {
  name: string;
  parentId?: number;
  rating?: number;
  favorite?: boolean;
  details?: string;
  ignoreAutoTag?: boolean;
  organized?: boolean;
  urls?: string[];
  aliases?: string[];
  tagIds?: number[];
  customFields?: Record<string, unknown>;
}

export interface StudioUpdate extends Partial<StudioCreate> {}

export interface Gallery {
  id: number;
  title?: string;
  code?: string;
  date?: string;
  details?: string;
  photographer?: string;
  rating?: number;
  organized: boolean;
  coverPath?: string;
  studioId?: number;
  studioName?: string;
  urls: string[];
  tags: Tag[];
  performers: PerformerSummary[];
  imageCount: number;
  sceneCount: number;
  sceneIds: number[];
  folderPath?: string;
  files: GalleryFileInfo[];
  customFields?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface GalleryFileInfo {
  id: number;
  path: string;
  size: number;
  modTime: string;
  fingerprints: { type: string; value: string }[];
}

export interface GalleryChapter {
  id: number;
  title: string;
  imageIndex: number;
  galleryId: number;
  createdAt: string;
  updatedAt: string;
}

export interface GalleryChapterCreate {
  title: string;
  imageIndex: number;
}

export interface GalleryChapterUpdate {
  title?: string;
  imageIndex?: number;
}

export interface GalleryCreate {
  title?: string;
  code?: string;
  date?: string;
  details?: string;
  photographer?: string;
  rating?: number;
  organized?: boolean;
  studioId?: number;
  urls?: string[];
  tagIds?: number[];
  performerIds?: number[];
  sceneIds?: number[];
  customFields?: Record<string, unknown>;
}

export interface GalleryUpdate extends Partial<GalleryCreate> {}

export interface Image {
  id: number;
  title?: string;
  code?: string;
  details?: string;
  photographer?: string;
  rating?: number;
  organized: boolean;
  oCounter: number;
  studioId?: number;
  studioName?: string;
  date?: string;
  urls: string[];
  tags: Tag[];
  performers: PerformerSummary[];
  galleryCount: number;
  galleryIds: number[];
  customFields?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface ImageCreate {
  title?: string;
  code?: string;
  details?: string;
  photographer?: string;
  rating?: number;
  organized?: boolean;
  studioId?: number;
  date?: string;
  urls?: string[];
  tagIds?: number[];
  performerIds?: number[];
  galleryIds?: number[];
  customFields?: Record<string, unknown>;
}

export interface ImageUpdate {
  title?: string;
  code?: string;
  details?: string;
  photographer?: string;
  rating?: number;
  organized?: boolean;
  studioId?: number;
  date?: string;
  urls?: string[];
  tagIds?: number[];
  performerIds?: number[];
  galleryIds?: number[];
  customFields?: Record<string, unknown>;
}

export interface Group {
  id: number;
  name: string;
  aliases?: string;
  duration?: number;
  date?: string;
  rating?: number;
  studioId?: number;
  studioName?: string;
  director?: string;
  synopsis?: string;
  frontImagePath?: string;
  backImagePath?: string;
  urls: string[];
  tags: Tag[];
  sceneCount: number;
  subGroupCount: number;
  containingGroupCount: number;
  customFields?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface GroupSummary {
  id: number;
  name: string;
  sceneIndex: number;
}

export interface GallerySummary {
  id: number;
  title?: string;
  date?: string;
}

export interface GroupCreate {
  name: string;
  aliases?: string;
  duration?: number;
  date?: string;
  rating?: number;
  studioId?: number;
  director?: string;
  synopsis?: string;
  urls?: string[];
  tagIds?: number[];
  customFields?: Record<string, unknown>;
}

export interface GroupUpdate extends Partial<GroupCreate> {}

export interface VideoFile {
  id: number;
  path: string;
  basename: string;
  format: string;
  width: number;
  height: number;
  duration: number;
  videoCodec: string;
  audioCodec: string;
  frameRate: number;
  bitRate: number;
  size: number;
  fingerprints: Fingerprint[];
}

export interface Fingerprint {
  type: string;
  value: string;
}

export interface SceneMarkerSummary {
  id: number;
  title: string;
  seconds: number;
  endSeconds?: number;
  primaryTagId: number;
  primaryTagName: string;
}

export interface SceneMarkerWall {
  id: number;
  title: string;
  seconds: number;
  endSeconds?: number;
  primaryTagId: number;
  primaryTagName: string;
  sceneId: number;
  sceneTitle: string;
  scenePath: string;
  tags: { id: number; name: string }[];
}

export interface SceneMarkerCreate {
  title: string;
  seconds: number;
  endSeconds?: number;
  primaryTagId: number;
  tagIds?: number[];
}

export interface SceneMarkerUpdate {
  title?: string;
  seconds?: number;
  endSeconds?: number;
  primaryTagId?: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  totalCount: number;
  page: number;
  perPage: number;
}

export interface Stats {
  sceneCount: number;
  imageCount: number;
  galleryCount: number;
  performerCount: number;
  studioCount: number;
  tagCount: number;
  groupCount: number;
  totalFileSize: number;
  totalPlayDuration: number;
}

export interface SystemStatus {
  version: string;
  appDir: string;
  configFile: string;
  databasePath: string;
}

export type RatingSystemType = "stars" | "decimal";
export type RatingStarPrecision = "full" | "half" | "quarter" | "tenth";

export interface RatingSystemOptions {
  type: RatingSystemType;
  starPrecision: RatingStarPrecision;
}

export interface InterfaceConfig {
  language?: string;
  menuItems: string[];
  handyConnectionEnabled: boolean;
  handyKey?: string;
  defaultDurationForImages?: number;
  disableDropdownCreatePerformer: boolean;
  disableDropdownCreateStudio: boolean;
  disableDropdownCreateTag: boolean;
}

export interface UiConfig {
  title?: string;
  abbreviateCounters: boolean;
  ratingSystemOptions: RatingSystemOptions;
  showStudioAsText: boolean;
  customCss?: string;
  customJs?: string;
  enableCSSCustomization: boolean;
  enableJSCustomization: boolean;
  customLocalesPath?: string;
  autostartVideo: boolean;
  autostartVideoOnPlaySelected: boolean;
  continuePlaylistDefault: boolean;
  showAbLoopControls: boolean;
  trackActivity: boolean;
  soundOnPreview: boolean;
  previewSegmentDuration: number;
  previewSegments: number;
  previewExcludeStart: string;
  previewExcludeEnd: string;
  wallShowTitle: boolean;
  wallPlayback: number;
  deleteFileDefault: boolean;
  slideshowDelay: number;
  noBrowser: boolean;
  notificationsEnabled: boolean;
}

export interface SecurityConfig {
  enabled: boolean;
  username?: string;
  maxSessionAgeMinutes: number;
  newPassword?: string;
}

export interface PackageSource {
  name: string;
  url: string;
}

export interface StashBox {
  endpoint: string;
  apiKey: string;
  name: string;
  maxRequestsPerMinute: number;
}

export interface ScrapingConfig {
  scraperDirectories: string[];
  scraperPackageSources: PackageSource[];
  stashBoxes: StashBox[];
}

export interface StashConfig {
  stashPaths: StashPathConfig[];
  generatedPath?: string;
  cachePath?: string;
  host: string;
  port: number;
  maxParallelTasks: number;
  calculateMd5: boolean;
  videoExtensions: string[];
  imageExtensions: string[];
  galleryExtensions: string[];
  excludePatterns: string[];
  maxTranscodeSize: number;
  maxStreamingTranscodeSize: number;
  transcodeHardwareAcceleration: string;
  transcodeInputArgs?: string;
  transcodeOutputArgs?: string;
  liveTranscodeInputArgs?: string;
  liveTranscodeOutputArgs?: string;
  drawFunscriptHeatmapRange: boolean;
  previewPreset: string;
  previewAudio: string;
  logLevel: string;
  logFile?: string;
  logOut: boolean;
  logAccess: boolean;
  ffmpegPath?: string;
  ffprobePath?: string;
  interface: InterfaceConfig;
  ui: UiConfig;
  security: SecurityConfig;
  scraping: ScrapingConfig;
}

export interface StashPathConfig {
  path: string;
  excludeVideo: boolean;
  excludeImage: boolean;
}

export interface JobInfo {
  id: string;
  type: string;
  description: string;
  status: "Pending" | "Running" | "Completed" | "Failed" | "Cancelled";
  progress: number;
  subTask?: string;
  startedAt: string;
  completedAt?: string;
  error?: string;
}

export interface FindFilter {
  q?: string;
  page?: number;
  perPage?: number;
  sort?: string;
  direction?: "asc" | "desc";
}

export interface SavedFilter {
  id: number;
  mode: string;
  name: string;
  findFilter?: string;
  objectFilter?: string;
  uiOptions?: string;
}

export interface SavedFilterCreate {
  mode: string;
  name: string;
  findFilter?: string;
  objectFilter?: string;
  uiOptions?: string;
}

export interface SavedFilterUpdate {
  mode?: string;
  name?: string;
  findFilter?: string;
  objectFilter?: string;
  uiOptions?: string;
}

export interface DlnaStatus {
  running: boolean;
  untilDisabled?: string;
  recentIps: string[];
  allowedIps: string[];
}

export interface ScraperSummary {
  id: string;
  name: string;
  entityType: string;
  supportedScrapes: string[];
  urls: string[];
  sourcePath: string;
}

export interface StashBoxValidationResult {
  valid: boolean;
  status: string;
  username?: string;
}

export interface StashBoxPerformerMatch {
  endpoint: string;
  stashBoxName: string;
  id: string;
  name: string;
  disambiguation?: string;
  gender?: string;
  birthDate?: string;
  country?: string;
  imageUrl?: string;
  deleted: boolean;
  mergedIntoId?: string;
  aliases: string[];
  urls: string[];
}

export interface StashBoxPerformerImportRequest {
  endpoint: string;
  performerId: string;
}

export interface StashBoxStudioMatch {
  endpoint: string;
  stashBoxName: string;
  id: string;
  name: string;
  imageUrl?: string;
  aliases: string[];
  urls: string[];
  parentName?: string;
}

export interface StashBoxStudioImportRequest {
  endpoint: string;
  studioId: string;
}

export interface StashBoxEntityCandidate {
  remoteId: string;
  name: string;
  existsLocally: boolean;
  localId?: number;
}

export interface StashBoxSceneEntityOverride {
  remoteId: string;
  name: string;
  action: string;
  localId?: number;
}

export interface StashBoxSceneMatch {
  endpoint: string;
  stashBoxName: string;
  id: string;
  title?: string;
  code?: string;
  date?: string;
  director?: string;
  details?: string;
  studioName?: string;
  imageUrl?: string;
  duration?: number;
  performerNames: string[];
  tagNames: string[];
  urls: string[];
  fingerprintAlgorithms: string[];
  fingerprints: StashBoxFingerprint[];
  studioCandidate?: StashBoxEntityCandidate;
  performerCandidates: StashBoxEntityCandidate[];
  tagCandidates: StashBoxEntityCandidate[];
}

export interface StashBoxFingerprint {
  algorithm: string;
  hash: string;
  duration?: number;
}

export interface StashBoxSceneImportRequest {
  endpoint: string;
  sceneId: string;
  setCoverImage?: boolean;
  setTags?: boolean;
  setPerformers?: boolean;
  setStudio?: boolean;
  onlyExistingTags?: boolean;
  onlyExistingPerformers?: boolean;
  onlyExistingStudio?: boolean;
  markOrganized?: boolean;
  excludedTagNames?: string[];
  excludedPerformerNames?: string[];
  studioOverride?: StashBoxSceneEntityOverride;
  performerOverrides?: StashBoxSceneEntityOverride[];
  tagOverrides?: StashBoxSceneEntityOverride[];
}

// ===== Filter Criteria =====

export type CriterionModifier =
  | "EQUALS" | "NOT_EQUALS" | "GREATER_THAN" | "LESS_THAN"
  | "INCLUDES" | "EXCLUDES" | "INCLUDES_ALL" | "EXCLUDES_ALL"
  | "IS_NULL" | "NOT_NULL" | "BETWEEN" | "NOT_BETWEEN"
  | "MATCHES_REGEX" | "NOT_MATCHES_REGEX";

export interface IntCriterion {
  value: number;
  value2?: number;
  modifier: CriterionModifier;
}

export interface StringCriterion {
  value: string;
  modifier: CriterionModifier;
}

export interface BoolCriterion {
  value: boolean;
}

export interface MultiIdCriterion {
  value: number[];
  modifier: CriterionModifier;
  excludes?: number[];
}

export interface DateCriterion {
  value: string;
  value2?: string;
  modifier: CriterionModifier;
}

export interface TimestampCriterion {
  value: string;
  value2?: string;
  modifier: CriterionModifier;
}

export interface SceneFilterCriteria {
  title?: string;
  code?: string;
  path?: string;
  rating?: number;
  organized?: boolean;
  studioId?: number;
  groupId?: number;
  tagIds?: number[];
  performerIds?: number[];
  ratingCriterion?: IntCriterion;
  oCounterCriterion?: IntCriterion;
  durationCriterion?: IntCriterion;
  resolutionCriterion?: IntCriterion;
  playCountCriterion?: IntCriterion;
  performerCountCriterion?: IntCriterion;
  tagsCriterion?: MultiIdCriterion;
  performersCriterion?: MultiIdCriterion;
  studiosCriterion?: MultiIdCriterion;
  groupsCriterion?: MultiIdCriterion;
  organizedCriterion?: BoolCriterion;
  hasMarkersCriterion?: BoolCriterion;
  interactiveCriterion?: BoolCriterion;
  pathCriterion?: StringCriterion;
  urlCriterion?: StringCriterion;
  dateCriterion?: DateCriterion;
  createdAtCriterion?: TimestampCriterion;
  updatedAtCriterion?: TimestampCriterion;
  performerFavoriteCriterion?: BoolCriterion;
  videoCodecCriterion?: StringCriterion;
  audioCodecCriterion?: StringCriterion;
  frameRateCriterion?: IntCriterion;
  bitrateInterval?: IntCriterion;
  fileCountCriterion?: IntCriterion;
  stashIdCriterion?: StringCriterion;
  isMissingCriterion?: BoolCriterion;
  duplicatedCriterion?: StringCriterion;
  titleCriterion?: StringCriterion;
  codeCriterion?: StringCriterion;
  detailsCriterion?: StringCriterion;
  directorCriterion?: StringCriterion;
  tagCountCriterion?: IntCriterion;
  resumeTimeCriterion?: IntCriterion;
  playDurationCriterion?: IntCriterion;
  lastPlayedAtCriterion?: TimestampCriterion;
  galleriesCriterion?: MultiIdCriterion;
}

export interface PerformerFilterCriteria {
  name?: string;
  favorite?: boolean;
  rating?: number;
  tagIds?: number[];
  ratingCriterion?: IntCriterion;
  ageCriterion?: IntCriterion;
  genderCriterion?: StringCriterion;
  ethnicityCriterion?: StringCriterion;
  countryCriterion?: StringCriterion;
  favoriteCriterion?: BoolCriterion;
  tagsCriterion?: MultiIdCriterion;
  studiosCriterion?: MultiIdCriterion;
  sceneCountCriterion?: IntCriterion;
  imageCountCriterion?: IntCriterion;
  galleryCountCriterion?: IntCriterion;
  birthdateCriterion?: DateCriterion;
  createdAtCriterion?: TimestampCriterion;
  updatedAtCriterion?: TimestampCriterion;
  pathCriterion?: StringCriterion;
  urlCriterion?: StringCriterion;
  weightCriterion?: IntCriterion;
  heightCriterion?: IntCriterion;
  isMissingCriterion?: BoolCriterion;
  stashIdCriterion?: StringCriterion;
}

export interface TagFilterCriteria {
  name?: string;
  favorite?: boolean;
  favoriteCriterion?: BoolCriterion;
  sceneCountCriterion?: IntCriterion;
  markerCountCriterion?: IntCriterion;
  performerCountCriterion?: IntCriterion;
  parentsCriterion?: MultiIdCriterion;
  childrenCriterion?: MultiIdCriterion;
  isMissingCriterion?: BoolCriterion;
  createdAtCriterion?: TimestampCriterion;
  updatedAtCriterion?: TimestampCriterion;
}

export interface StudioFilterCriteria {
  name?: string;
  favorite?: boolean;
  parentId?: number;
  tagIds?: number[];
  ratingCriterion?: IntCriterion;
  favoriteCriterion?: BoolCriterion;
  tagsCriterion?: MultiIdCriterion;
  sceneCountCriterion?: IntCriterion;
  urlCriterion?: StringCriterion;
  stashIdCriterion?: StringCriterion;
  isMissingCriterion?: BoolCriterion;
  createdAtCriterion?: TimestampCriterion;
  updatedAtCriterion?: TimestampCriterion;
}

export interface GalleryFilterCriteria {
  title?: string;
  rating?: number;
  organized?: boolean;
  studioId?: number;
  tagIds?: number[];
  performerIds?: number[];
  ratingCriterion?: IntCriterion;
  organizedCriterion?: BoolCriterion;
  tagsCriterion?: MultiIdCriterion;
  performersCriterion?: MultiIdCriterion;
  studiosCriterion?: MultiIdCriterion;
  imageCountCriterion?: IntCriterion;
  dateCriterion?: DateCriterion;
  pathCriterion?: StringCriterion;
  urlCriterion?: StringCriterion;
  createdAtCriterion?: TimestampCriterion;
  updatedAtCriterion?: TimestampCriterion;
  performerFavoriteCriterion?: BoolCriterion;
  isMissingCriterion?: BoolCriterion;
}

export interface ImageFilterCriteria {
  title?: string;
  rating?: number;
  organized?: boolean;
  studioId?: number;
  galleryId?: number;
  tagIds?: number[];
  performerIds?: number[];
  ratingCriterion?: IntCriterion;
  organizedCriterion?: BoolCriterion;
  tagsCriterion?: MultiIdCriterion;
  performersCriterion?: MultiIdCriterion;
  studiosCriterion?: MultiIdCriterion;
  galleriesCriterion?: MultiIdCriterion;
  oCounterCriterion?: IntCriterion;
  resolutionCriterion?: IntCriterion;
  pathCriterion?: StringCriterion;
  createdAtCriterion?: TimestampCriterion;
  updatedAtCriterion?: TimestampCriterion;
  performerFavoriteCriterion?: BoolCriterion;
  isMissingCriterion?: BoolCriterion;
}

export interface GroupFilterCriteria {
  name?: string;
  rating?: number;
  studioId?: number;
  ratingCriterion?: IntCriterion;
  durationCriterion?: IntCriterion;
  studiosCriterion?: MultiIdCriterion;
  tagsCriterion?: MultiIdCriterion;
  dateCriterion?: DateCriterion;
  urlCriterion?: StringCriterion;
  createdAtCriterion?: TimestampCriterion;
  updatedAtCriterion?: TimestampCriterion;
  isMissingCriterion?: BoolCriterion;
}

export interface FilteredQueryRequest<T = Record<string, unknown>> {
  findFilter?: FindFilter;
  objectFilter?: T;
}

// ===== Bulk Edit Types =====

export type BulkUpdateMode = "SET" | "ADD" | "REMOVE";

export interface BulkSceneUpdate {
  ids: number[];
  rating?: number;
  organized?: boolean;
  studioId?: number | null;
  date?: string;
  code?: string;
  director?: string;
  tagIds?: number[];
  tagMode?: BulkUpdateMode;
  performerIds?: number[];
  performerMode?: BulkUpdateMode;
  groupIds?: number[];
  groupMode?: BulkUpdateMode;
}

export interface BulkPerformerUpdate {
  ids: number[];
  rating?: number;
  favorite?: boolean;
  tagIds?: number[];
  tagMode?: BulkUpdateMode;
}

export interface BulkTagUpdate {
  ids: number[];
  favorite?: boolean;
  ignoreAutoTag?: boolean;
}

export interface BulkStudioUpdate {
  ids: number[];
  rating?: number;
  favorite?: boolean;
}

export interface BulkGalleryUpdate {
  ids: number[];
  rating?: number;
  organized?: boolean;
  studioId?: number | null;
  tagIds?: number[];
  tagMode?: BulkUpdateMode;
  performerIds?: number[];
  performerMode?: BulkUpdateMode;
}

export interface BulkImageUpdate {
  ids: number[];
  rating?: number;
  organized?: boolean;
  studioId?: number | null;
  tagIds?: number[];
  tagMode?: BulkUpdateMode;
  performerIds?: number[];
  performerMode?: BulkUpdateMode;
}

export interface BulkGroupUpdate {
  ids: number[];
  rating?: number;
  studioId?: number | null;
}

// ===== Plugin Types =====
export interface Plugin {
  id: string;
  name: string;
  description: string;
  version: string;
  enabled: boolean;
  tasks: PluginTask[];
  settings?: PluginSettingSchema[];
  url?: string;
}

export interface PluginSettingSchema {
  name: string;
  type: "STRING" | "NUMBER" | "BOOLEAN";
  displayName?: string;
  description?: string;
}

export interface PluginTask {
  name: string;
  description: string;
}

export interface RunPluginTaskRequest {
  pluginId: string;
  taskName: string;
  args?: Record<string, string>;
}

export interface PluginSettings {
  enabledMap: Record<string, boolean>;
}

export interface Package {
  name: string;
  description: string;
  version: string;
  sourceUrl: string;
  type: string;
  installed: boolean;
  installedVersion?: string;
}

// ===== Extension System Types =====
export interface ExtensionManifest {
  pages: ExtensionPageDef[];
  slots: ExtensionSlotContribution[];
  tabs: ExtensionTabContribution[];
  themes: ExtensionThemeDef[];
  settingsPanels: ExtensionSettingsPanel[];
  pageOverrides: ExtensionPageOverride[];
  dialogOverrides: ExtensionDialogOverride[];
  jsBundleUrl?: string;
  cssBundleUrl?: string;
}

export interface ExtensionPageDef {
  route: string;
  label: string;
  icon?: string;
  detailRoute?: string;
  showInNav: boolean;
  navOrder: number;
  requiredPermission?: string;
  componentName?: string;
  extensionId?: string;
}

export interface ExtensionSlotContribution {
  id: string;
  slot: string;
  extensionId: string;
  contentType: "component" | "html";
  componentName?: string;
  html?: string;
  order: number;
}

export interface ExtensionTabContribution {
  key: string;
  label: string;
  pageType: string;
  extensionId: string;
  componentName: string;
  order: number;
}

export interface ExtensionThemeDef {
  id: string;
  name: string;
  description?: string;
  cssVariables?: Record<string, string>;
  cssUrl?: string;
}

export interface ExtensionSettingsPanel {
  id: string;
  label: string;
  extensionId: string;
  componentName: string;
  order: number;
}

export interface ExtensionPageOverride {
  targetPage: string;
  extensionId: string;
  componentName: string;
  priority: number;
}

export interface ExtensionDialogOverride {
  dialogId: string;
  extensionId: string;
  componentName: string;
  priority: number;
}

export interface ExtensionInfo {
  id: string;
  name: string;
  version: string;
  description?: string;
  author?: string;
  iconUrl?: string;
  enabled: boolean;
  hasUI: boolean;
  hasApi: boolean;
  hasState: boolean;
  hasJobs: boolean;
  hasEvents: boolean;
  jobs: { id: string; name: string; description?: string }[];
}
