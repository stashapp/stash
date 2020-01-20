/* tslint:disable */
/* eslint-disable */
// Generated in 2020-01-02T16:26:49+01:00
export type Maybe<T> = T | undefined;

export interface SceneFilterType {
  /** Filter by rating */
  rating?: Maybe<IntCriterionInput>;
  /** Filter by resolution */
  resolution?: Maybe<ResolutionEnum>;
  /** Filter to only include scenes which have markers. `true` or `false` */
  has_markers?: Maybe<string>;
  /** Filter to only include scenes missing this property */
  is_missing?: Maybe<string>;
  /** Filter to only include scenes with this studio */
  studios?: Maybe<MultiCriterionInput>;
  /** Filter to only include scenes with these tags */
  tags?: Maybe<MultiCriterionInput>;
  /** Filter to only include scenes with these performers */
  performers?: Maybe<MultiCriterionInput>;
}

export interface IntCriterionInput {
  value: number;

  modifier: CriterionModifier;
}

export interface MultiCriterionInput {
  value?: Maybe<string[]>;

  modifier: CriterionModifier;
}

export interface FindFilterType {
  q?: Maybe<string>;

  page?: Maybe<number>;

  per_page?: Maybe<number>;

  sort?: Maybe<string>;

  direction?: Maybe<SortDirectionEnum>;
}

export interface SceneParserInput {
  ignoreWords?: Maybe<string[]>;

  whitespaceCharacters?: Maybe<string>;

  capitalizeTitle?: Maybe<boolean>;
}

export interface SceneMarkerFilterType {
  /** Filter to only include scene markers with this tag */
  tag_id?: Maybe<string>;
  /** Filter to only include scene markers with these tags */
  tags?: Maybe<MultiCriterionInput>;
  /** Filter to only include scene markers attached to a scene with these tags */
  scene_tags?: Maybe<MultiCriterionInput>;
  /** Filter to only include scene markers with these performers */
  performers?: Maybe<MultiCriterionInput>;
}

export interface PerformerFilterType {
  /** Filter by favorite */
  filter_favorites?: Maybe<boolean>;
  /** Filter by birth year */
  birth_year?: Maybe<IntCriterionInput>;
  /** Filter by age */
  age?: Maybe<IntCriterionInput>;
  /** Filter by ethnicity */
  ethnicity?: Maybe<StringCriterionInput>;
  /** Filter by country */
  country?: Maybe<StringCriterionInput>;
  /** Filter by eye color */
  eye_color?: Maybe<StringCriterionInput>;
  /** Filter by height */
  height?: Maybe<StringCriterionInput>;
  /** Filter by measurements */
  measurements?: Maybe<StringCriterionInput>;
  /** Filter by fake tits value */
  fake_tits?: Maybe<StringCriterionInput>;
  /** Filter by career length */
  career_length?: Maybe<StringCriterionInput>;
  /** Filter by tattoos */
  tattoos?: Maybe<StringCriterionInput>;
  /** Filter by piercings */
  piercings?: Maybe<StringCriterionInput>;
  /** Filter by aliases */
  aliases?: Maybe<StringCriterionInput>;
}

export interface StringCriterionInput {
  value: string;

  modifier: CriterionModifier;
}

export interface ScrapedPerformerInput {
  name?: Maybe<string>;

  url?: Maybe<string>;

  twitter?: Maybe<string>;

  instagram?: Maybe<string>;

  birthdate?: Maybe<string>;

  ethnicity?: Maybe<string>;

  country?: Maybe<string>;

  eye_color?: Maybe<string>;

  height?: Maybe<string>;

  measurements?: Maybe<string>;

  fake_tits?: Maybe<string>;

  career_length?: Maybe<string>;

  tattoos?: Maybe<string>;

  piercings?: Maybe<string>;

  aliases?: Maybe<string>;
}

export interface SceneUpdateInput {
  clientMutationId?: Maybe<string>;

  id: string;

  title?: Maybe<string>;

  details?: Maybe<string>;

  url?: Maybe<string>;

  date?: Maybe<string>;

  rating?: Maybe<number>;

  studio_id?: Maybe<string>;

  gallery_id?: Maybe<string>;

  performer_ids?: Maybe<string[]>;

  tag_ids?: Maybe<string[]>;
  /** This should be base64 encoded */
  cover_image?: Maybe<string>;
}

export interface ScanMetadataInput {
  useFileMetadata: boolean;
}

export interface GenerateMetadataInput {
  sprites: boolean;

  previews: boolean;

  markers: boolean;

  transcodes: boolean;
}

export interface AutoTagMetadataInput {
  /** IDs of performers to tag files with, or "*" for all */
  performers?: Maybe<string[]>;
  /** IDs of studios to tag files with, or "*" for all */
  studios?: Maybe<string[]>;
  /** IDs of tags to tag files with, or "*" for all */
  tags?: Maybe<string[]>;
}

export interface BulkSceneUpdateInput {
  clientMutationId?: Maybe<string>;

  ids?: Maybe<string[]>;

  title?: Maybe<string>;

  details?: Maybe<string>;

  url?: Maybe<string>;

  date?: Maybe<string>;

  rating?: Maybe<number>;

  studio_id?: Maybe<string>;

  gallery_id?: Maybe<string>;

  performer_ids?: Maybe<string[]>;

  tag_ids?: Maybe<string[]>;
}

export interface SceneDestroyInput {
  id: string;

  delete_file?: Maybe<boolean>;

  delete_generated?: Maybe<boolean>;
}

export interface SceneMarkerCreateInput {
  title: string;

  seconds: number;

  scene_id: string;

  primary_tag_id: string;

  tag_ids?: Maybe<string[]>;
}

export interface SceneMarkerUpdateInput {
  id: string;

  title: string;

  seconds: number;

  scene_id: string;

  primary_tag_id: string;

  tag_ids?: Maybe<string[]>;
}

export interface PerformerCreateInput {
  name?: Maybe<string>;

  url?: Maybe<string>;

  birthdate?: Maybe<string>;

  ethnicity?: Maybe<string>;

  country?: Maybe<string>;

  eye_color?: Maybe<string>;

  height?: Maybe<string>;

  measurements?: Maybe<string>;

  fake_tits?: Maybe<string>;

  career_length?: Maybe<string>;

  tattoos?: Maybe<string>;

  piercings?: Maybe<string>;

  aliases?: Maybe<string>;

  twitter?: Maybe<string>;

  instagram?: Maybe<string>;

  favorite?: Maybe<boolean>;
  /** This should be base64 encoded */
  image?: Maybe<string>;
}

export interface PerformerUpdateInput {
  id: string;

  name?: Maybe<string>;

  url?: Maybe<string>;

  birthdate?: Maybe<string>;

  ethnicity?: Maybe<string>;

  country?: Maybe<string>;

  eye_color?: Maybe<string>;

  height?: Maybe<string>;

  measurements?: Maybe<string>;

  fake_tits?: Maybe<string>;

  career_length?: Maybe<string>;

  tattoos?: Maybe<string>;

  piercings?: Maybe<string>;

  aliases?: Maybe<string>;

  twitter?: Maybe<string>;

  instagram?: Maybe<string>;

  favorite?: Maybe<boolean>;
  /** This should be base64 encoded */
  image?: Maybe<string>;
}

export interface PerformerDestroyInput {
  id: string;
}

export interface StudioCreateInput {
  name: string;

  url?: Maybe<string>;
  /** This should be base64 encoded */
  image?: Maybe<string>;
}

export interface StudioUpdateInput {
  id: string;

  name?: Maybe<string>;

  url?: Maybe<string>;
  /** This should be base64 encoded */
  image?: Maybe<string>;
}

export interface StudioDestroyInput {
  id: string;
}

export interface TagCreateInput {
  name: string;
}

export interface TagUpdateInput {
  id: string;

  name: string;
}

export interface TagDestroyInput {
  id: string;
}

export interface ConfigGeneralInput {
  /** Array of file paths to content */
  stashes?: Maybe<string[]>;
  /** Path to the SQLite database */
  databasePath?: Maybe<string>;
  /** Path to generated files */
  generatedPath?: Maybe<string>;
  /** Max generated transcode size */
  maxTranscodeSize?: Maybe<StreamingResolutionEnum>;
  /** Max streaming transcode size */
  maxStreamingTranscodeSize?: Maybe<StreamingResolutionEnum>;
  /** Username */
  username?: Maybe<string>;
  /** Password */
  password?: Maybe<string>;
  /** Name of the log file */
  logFile?: Maybe<string>;
  /** Whether to also output to stderr */
  logOut: boolean;
  /** Minimum log level */
  logLevel: string;
  /** Whether to log http access */
  logAccess: boolean;
  /** Array of file regexp to exclude from Scan */
  excludes?: Maybe<string[]>;
}

export interface ConfigInterfaceInput {
  /** Enable sound on mouseover previews */
  soundOnPreview?: Maybe<boolean>;
  /** Show title and tags in wall view */
  wallShowTitle?: Maybe<boolean>;
  /** Maximum duration (in seconds) in which a scene video will loop in the scene player */
  maximumLoopDuration?: Maybe<number>;
  /** If true, video will autostart on load in the scene player */
  autostartVideo?: Maybe<boolean>;
  /** If true, studio overlays will be shown as text instead of logo images */
  showStudioAsText?: Maybe<boolean>;
  /** Custom CSS */
  css?: Maybe<string>;

  cssEnabled?: Maybe<boolean>;
}

export enum CriterionModifier {
  Equals = "EQUALS",
  NotEquals = "NOT_EQUALS",
  GreaterThan = "GREATER_THAN",
  LessThan = "LESS_THAN",
  IsNull = "IS_NULL",
  NotNull = "NOT_NULL",
  IncludesAll = "INCLUDES_ALL",
  Includes = "INCLUDES",
  Excludes = "EXCLUDES"
}

export enum ResolutionEnum {
  Low = "LOW",
  Standard = "STANDARD",
  StandardHd = "STANDARD_HD",
  FullHd = "FULL_HD",
  FourK = "FOUR_K"
}

export enum SortDirectionEnum {
  Asc = "ASC",
  Desc = "DESC"
}

export enum LogLevel {
  Debug = "Debug",
  Info = "Info",
  Progress = "Progress",
  Warning = "Warning",
  Error = "Error"
}

export enum ScrapeType {
  Name = "NAME",
  Fragment = "FRAGMENT",
  Url = "URL"
}

export enum StreamingResolutionEnum {
  Low = "LOW",
  Standard = "STANDARD",
  StandardHd = "STANDARD_HD",
  FullHd = "FULL_HD",
  FourK = "FOUR_K",
  Original = "ORIGINAL"
}

/** Log entries */
export type Time = any;

// ====================================================
// Documents
// ====================================================

export type ConfigureGeneralVariables = {
  input: ConfigGeneralInput;
};

export type ConfigureGeneralMutation = {
  __typename?: "Mutation";

  configureGeneral: ConfigureGeneralConfigureGeneral;
};

export type ConfigureGeneralConfigureGeneral = ConfigGeneralDataFragment;

export type ConfigureInterfaceVariables = {
  input: ConfigInterfaceInput;
};

export type ConfigureInterfaceMutation = {
  __typename?: "Mutation";

  configureInterface: ConfigureInterfaceConfigureInterface;
};

export type ConfigureInterfaceConfigureInterface = ConfigInterfaceDataFragment;

export type PerformerCreateVariables = {
  name?: Maybe<string>;
  url?: Maybe<string>;
  birthdate?: Maybe<string>;
  ethnicity?: Maybe<string>;
  country?: Maybe<string>;
  eye_color?: Maybe<string>;
  height?: Maybe<string>;
  measurements?: Maybe<string>;
  fake_tits?: Maybe<string>;
  career_length?: Maybe<string>;
  tattoos?: Maybe<string>;
  piercings?: Maybe<string>;
  aliases?: Maybe<string>;
  twitter?: Maybe<string>;
  instagram?: Maybe<string>;
  favorite?: Maybe<boolean>;
  image?: Maybe<string>;
};

export type PerformerCreateMutation = {
  __typename?: "Mutation";

  performerCreate: Maybe<PerformerCreatePerformerCreate>;
};

export type PerformerCreatePerformerCreate = PerformerDataFragment;

export type PerformerUpdateVariables = {
  id: string;
  name?: Maybe<string>;
  url?: Maybe<string>;
  birthdate?: Maybe<string>;
  ethnicity?: Maybe<string>;
  country?: Maybe<string>;
  eye_color?: Maybe<string>;
  height?: Maybe<string>;
  measurements?: Maybe<string>;
  fake_tits?: Maybe<string>;
  career_length?: Maybe<string>;
  tattoos?: Maybe<string>;
  piercings?: Maybe<string>;
  aliases?: Maybe<string>;
  twitter?: Maybe<string>;
  instagram?: Maybe<string>;
  favorite?: Maybe<boolean>;
  image?: Maybe<string>;
};

export type PerformerUpdateMutation = {
  __typename?: "Mutation";

  performerUpdate: Maybe<PerformerUpdatePerformerUpdate>;
};

export type PerformerUpdatePerformerUpdate = PerformerDataFragment;

export type PerformerDestroyVariables = {
  id: string;
};

export type PerformerDestroyMutation = {
  __typename?: "Mutation";

  performerDestroy: boolean;
};

export type SceneMarkerCreateVariables = {
  title: string;
  seconds: number;
  scene_id: string;
  primary_tag_id: string;
  tag_ids?: Maybe<string[]>;
};

export type SceneMarkerCreateMutation = {
  __typename?: "Mutation";

  sceneMarkerCreate: Maybe<SceneMarkerCreateSceneMarkerCreate>;
};

export type SceneMarkerCreateSceneMarkerCreate = SceneMarkerDataFragment;

export type SceneMarkerUpdateVariables = {
  id: string;
  title: string;
  seconds: number;
  scene_id: string;
  primary_tag_id: string;
  tag_ids?: Maybe<string[]>;
};

export type SceneMarkerUpdateMutation = {
  __typename?: "Mutation";

  sceneMarkerUpdate: Maybe<SceneMarkerUpdateSceneMarkerUpdate>;
};

export type SceneMarkerUpdateSceneMarkerUpdate = SceneMarkerDataFragment;

export type SceneMarkerDestroyVariables = {
  id: string;
};

export type SceneMarkerDestroyMutation = {
  __typename?: "Mutation";

  sceneMarkerDestroy: boolean;
};

export type SceneUpdateVariables = {
  id: string;
  title?: Maybe<string>;
  details?: Maybe<string>;
  url?: Maybe<string>;
  date?: Maybe<string>;
  rating?: Maybe<number>;
  studio_id?: Maybe<string>;
  gallery_id?: Maybe<string>;
  performer_ids?: Maybe<string[]>;
  tag_ids?: Maybe<string[]>;
  cover_image?: Maybe<string>;
};

export type SceneUpdateMutation = {
  __typename?: "Mutation";

  sceneUpdate: Maybe<SceneUpdateSceneUpdate>;
};

export type SceneUpdateSceneUpdate = SceneDataFragment;

export type BulkSceneUpdateVariables = {
  ids?: Maybe<string[]>;
  title?: Maybe<string>;
  details?: Maybe<string>;
  url?: Maybe<string>;
  date?: Maybe<string>;
  rating?: Maybe<number>;
  studio_id?: Maybe<string>;
  gallery_id?: Maybe<string>;
  performer_ids?: Maybe<string[]>;
  tag_ids?: Maybe<string[]>;
};

export type BulkSceneUpdateMutation = {
  __typename?: "Mutation";

  bulkSceneUpdate: Maybe<BulkSceneUpdateBulkSceneUpdate[]>;
};

export type BulkSceneUpdateBulkSceneUpdate = SceneDataFragment;

export type ScenesUpdateVariables = {
  input: SceneUpdateInput[];
};

export type ScenesUpdateMutation = {
  __typename?: "Mutation";

  scenesUpdate: Maybe<Maybe<ScenesUpdateScenesUpdate>[]>;
};

export type ScenesUpdateScenesUpdate = SceneDataFragment;

export type SceneDestroyVariables = {
  id: string;
  delete_file?: Maybe<boolean>;
  delete_generated?: Maybe<boolean>;
};

export type SceneDestroyMutation = {
  __typename?: "Mutation";

  sceneDestroy: boolean;
};

export type StudioCreateVariables = {
  name: string;
  url?: Maybe<string>;
  image?: Maybe<string>;
};

export type StudioCreateMutation = {
  __typename?: "Mutation";

  studioCreate: Maybe<StudioCreateStudioCreate>;
};

export type StudioCreateStudioCreate = StudioDataFragment;

export type StudioUpdateVariables = {
  id: string;
  name?: Maybe<string>;
  url?: Maybe<string>;
  image?: Maybe<string>;
};

export type StudioUpdateMutation = {
  __typename?: "Mutation";

  studioUpdate: Maybe<StudioUpdateStudioUpdate>;
};

export type StudioUpdateStudioUpdate = StudioDataFragment;

export type StudioDestroyVariables = {
  id: string;
};

export type StudioDestroyMutation = {
  __typename?: "Mutation";

  studioDestroy: boolean;
};

export type TagCreateVariables = {
  name: string;
};

export type TagCreateMutation = {
  __typename?: "Mutation";

  tagCreate: Maybe<TagCreateTagCreate>;
};

export type TagCreateTagCreate = TagDataFragment;

export type TagDestroyVariables = {
  id: string;
};

export type TagDestroyMutation = {
  __typename?: "Mutation";

  tagDestroy: boolean;
};

export type TagUpdateVariables = {
  id: string;
  name: string;
};

export type TagUpdateMutation = {
  __typename?: "Mutation";

  tagUpdate: Maybe<TagUpdateTagUpdate>;
};

export type TagUpdateTagUpdate = TagDataFragment;

export type FindGalleriesVariables = {
  filter?: Maybe<FindFilterType>;
};

export type FindGalleriesQuery = {
  __typename?: "Query";

  findGalleries: FindGalleriesFindGalleries;
};

export type FindGalleriesFindGalleries = {
  __typename?: "FindGalleriesResultType";

  count: number;

  galleries: FindGalleriesGalleries[];
};

export type FindGalleriesGalleries = GalleryDataFragment;

export type FindGalleryVariables = {
  id: string;
};

export type FindGalleryQuery = {
  __typename?: "Query";

  findGallery: Maybe<FindGalleryFindGallery>;
};

export type FindGalleryFindGallery = GalleryDataFragment;

export type SceneWallVariables = {
  q?: Maybe<string>;
};

export type SceneWallQuery = {
  __typename?: "Query";

  sceneWall: SceneWallSceneWall[];
};

export type SceneWallSceneWall = SceneDataFragment;

export type MarkerWallVariables = {
  q?: Maybe<string>;
};

export type MarkerWallQuery = {
  __typename?: "Query";

  markerWall: MarkerWallMarkerWall[];
};

export type MarkerWallMarkerWall = SceneMarkerDataFragment;

export type FindTagVariables = {
  id: string;
};

export type FindTagQuery = {
  __typename?: "Query";

  findTag: Maybe<FindTagFindTag>;
};

export type FindTagFindTag = TagDataFragment;

export type MarkerStringsVariables = {
  q?: Maybe<string>;
  sort?: Maybe<string>;
};

export type MarkerStringsQuery = {
  __typename?: "Query";

  markerStrings: Maybe<MarkerStringsMarkerStrings>[];
};

export type MarkerStringsMarkerStrings = {
  __typename?: "MarkerStringsResultType";

  id: string;

  count: number;

  title: string;
};

export type AllTagsVariables = {};

export type AllTagsQuery = {
  __typename?: "Query";

  allTags: AllTagsAllTags[];
};

export type AllTagsAllTags = TagDataFragment;

export type AllPerformersForFilterVariables = {};

export type AllPerformersForFilterQuery = {
  __typename?: "Query";

  allPerformers: AllPerformersForFilterAllPerformers[];
};

export type AllPerformersForFilterAllPerformers = SlimPerformerDataFragment;

export type AllStudiosForFilterVariables = {};

export type AllStudiosForFilterQuery = {
  __typename?: "Query";

  allStudios: AllStudiosForFilterAllStudios[];
};

export type AllStudiosForFilterAllStudios = SlimStudioDataFragment;

export type AllTagsForFilterVariables = {};

export type AllTagsForFilterQuery = {
  __typename?: "Query";

  allTags: AllTagsForFilterAllTags[];
};

export type AllTagsForFilterAllTags = {
  __typename?: "Tag";

  id: string;

  name: string;
};

export type ValidGalleriesForSceneVariables = {
  scene_id: string;
};

export type ValidGalleriesForSceneQuery = {
  __typename?: "Query";

  validGalleriesForScene: ValidGalleriesForSceneValidGalleriesForScene[];
};

export type ValidGalleriesForSceneValidGalleriesForScene = {
  __typename?: "Gallery";

  id: string;

  path: string;
};

export type StatsVariables = {};

export type StatsQuery = {
  __typename?: "Query";

  stats: StatsStats;
};

export type StatsStats = {
  __typename?: "StatsResultType";

  scene_count: number;

  gallery_count: number;

  performer_count: number;

  studio_count: number;

  tag_count: number;
};

export type LogsVariables = {};

export type LogsQuery = {
  __typename?: "Query";

  logs: LogsLogs[];
};

export type LogsLogs = LogEntryDataFragment;

export type VersionVariables = {};

export type VersionQuery = {
  __typename?: "Query";

  version: VersionVersion;
};

export type VersionVersion = {
  __typename?: "Version";

  version: Maybe<string>;

  hash: string;

  build_time: string;
};

export type FindPerformersVariables = {
  filter?: Maybe<FindFilterType>;
  performer_filter?: Maybe<PerformerFilterType>;
};

export type FindPerformersQuery = {
  __typename?: "Query";

  findPerformers: FindPerformersFindPerformers;
};

export type FindPerformersFindPerformers = {
  __typename?: "FindPerformersResultType";

  count: number;

  performers: FindPerformersPerformers[];
};

export type FindPerformersPerformers = PerformerDataFragment;

export type FindPerformerVariables = {
  id: string;
};

export type FindPerformerQuery = {
  __typename?: "Query";

  findPerformer: Maybe<FindPerformerFindPerformer>;
};

export type FindPerformerFindPerformer = PerformerDataFragment;

export type FindSceneMarkersVariables = {
  filter?: Maybe<FindFilterType>;
  scene_marker_filter?: Maybe<SceneMarkerFilterType>;
};

export type FindSceneMarkersQuery = {
  __typename?: "Query";

  findSceneMarkers: FindSceneMarkersFindSceneMarkers;
};

export type FindSceneMarkersFindSceneMarkers = {
  __typename?: "FindSceneMarkersResultType";

  count: number;

  scene_markers: FindSceneMarkersSceneMarkers[];
};

export type FindSceneMarkersSceneMarkers = SceneMarkerDataFragment;

export type FindScenesVariables = {
  filter?: Maybe<FindFilterType>;
  scene_filter?: Maybe<SceneFilterType>;
  scene_ids?: Maybe<number[]>;
};

export type FindScenesQuery = {
  __typename?: "Query";

  findScenes: FindScenesFindScenes;
};

export type FindScenesFindScenes = {
  __typename?: "FindScenesResultType";

  count: number;

  scenes: FindScenesScenes[];
};

export type FindScenesScenes = SlimSceneDataFragment;

export type FindScenesByPathRegexVariables = {
  filter?: Maybe<FindFilterType>;
};

export type FindScenesByPathRegexQuery = {
  __typename?: "Query";

  findScenesByPathRegex: FindScenesByPathRegexFindScenesByPathRegex;
};

export type FindScenesByPathRegexFindScenesByPathRegex = {
  __typename?: "FindScenesResultType";

  count: number;

  scenes: FindScenesByPathRegexScenes[];
};

export type FindScenesByPathRegexScenes = SlimSceneDataFragment;

export type FindSceneVariables = {
  id: string;
  checksum?: Maybe<string>;
};

export type FindSceneQuery = {
  __typename?: "Query";

  findScene: Maybe<FindSceneFindScene>;

  sceneMarkerTags: FindSceneSceneMarkerTags[];
};

export type FindSceneFindScene = SceneDataFragment;

export type FindSceneSceneMarkerTags = {
  __typename?: "SceneMarkerTag";

  tag: FindSceneTag;

  scene_markers: FindSceneSceneMarkers[];
};

export type FindSceneTag = {
  __typename?: "Tag";

  id: string;

  name: string;
};

export type FindSceneSceneMarkers = SceneMarkerDataFragment;

export type ParseSceneFilenamesVariables = {
  filter: FindFilterType;
  config: SceneParserInput;
};

export type ParseSceneFilenamesQuery = {
  __typename?: "Query";

  parseSceneFilenames: ParseSceneFilenamesParseSceneFilenames;
};

export type ParseSceneFilenamesParseSceneFilenames = {
  __typename?: "SceneParserResultType";

  count: number;

  results: ParseSceneFilenamesResults[];
};

export type ParseSceneFilenamesResults = {
  __typename?: "SceneParserResult";

  scene: ParseSceneFilenamesScene;

  title: Maybe<string>;

  details: Maybe<string>;

  url: Maybe<string>;

  date: Maybe<string>;

  rating: Maybe<number>;

  studio_id: Maybe<string>;

  gallery_id: Maybe<string>;

  performer_ids: Maybe<string[]>;

  tag_ids: Maybe<string[]>;
};

export type ParseSceneFilenamesScene = SlimSceneDataFragment;

export type ScrapeFreeonesVariables = {
  performer_name: string;
};

export type ScrapeFreeonesQuery = {
  __typename?: "Query";

  scrapeFreeones: Maybe<ScrapeFreeonesScrapeFreeones>;
};

export type ScrapeFreeonesScrapeFreeones = {
  __typename?: "ScrapedPerformer";

  name: Maybe<string>;

  url: Maybe<string>;

  twitter: Maybe<string>;

  instagram: Maybe<string>;

  birthdate: Maybe<string>;

  ethnicity: Maybe<string>;

  country: Maybe<string>;

  eye_color: Maybe<string>;

  height: Maybe<string>;

  measurements: Maybe<string>;

  fake_tits: Maybe<string>;

  career_length: Maybe<string>;

  tattoos: Maybe<string>;

  piercings: Maybe<string>;

  aliases: Maybe<string>;
};

export type ScrapeFreeonesPerformersVariables = {
  q: string;
};

export type ScrapeFreeonesPerformersQuery = {
  __typename?: "Query";

  scrapeFreeonesPerformerList: string[];
};

export type ListPerformerScrapersVariables = {};

export type ListPerformerScrapersQuery = {
  __typename?: "Query";

  listPerformerScrapers: ListPerformerScrapersListPerformerScrapers[];
};

export type ListPerformerScrapersListPerformerScrapers = {
  __typename?: "Scraper";

  id: string;

  name: string;

  performer: Maybe<ListPerformerScrapersPerformer>;
};

export type ListPerformerScrapersPerformer = {
  __typename?: "ScraperSpec";

  urls: Maybe<string[]>;

  supported_scrapes: ScrapeType[];
};

export type ListSceneScrapersVariables = {};

export type ListSceneScrapersQuery = {
  __typename?: "Query";

  listSceneScrapers: ListSceneScrapersListSceneScrapers[];
};

export type ListSceneScrapersListSceneScrapers = {
  __typename?: "Scraper";

  id: string;

  name: string;

  scene: Maybe<ListSceneScrapersScene>;
};

export type ListSceneScrapersScene = {
  __typename?: "ScraperSpec";

  urls: Maybe<string[]>;

  supported_scrapes: ScrapeType[];
};

export type ScrapePerformerListVariables = {
  scraper_id: string;
  query: string;
};

export type ScrapePerformerListQuery = {
  __typename?: "Query";

  scrapePerformerList: ScrapePerformerListScrapePerformerList[];
};

export type ScrapePerformerListScrapePerformerList = ScrapedPerformerDataFragment;

export type ScrapePerformerVariables = {
  scraper_id: string;
  scraped_performer: ScrapedPerformerInput;
};

export type ScrapePerformerQuery = {
  __typename?: "Query";

  scrapePerformer: Maybe<ScrapePerformerScrapePerformer>;
};

export type ScrapePerformerScrapePerformer = ScrapedPerformerDataFragment;

export type ScrapePerformerUrlVariables = {
  url: string;
};

export type ScrapePerformerUrlQuery = {
  __typename?: "Query";

  scrapePerformerURL: Maybe<ScrapePerformerUrlScrapePerformerUrl>;
};

export type ScrapePerformerUrlScrapePerformerUrl = ScrapedPerformerDataFragment;

export type ScrapeSceneVariables = {
  scraper_id: string;
  scene: SceneUpdateInput;
};

export type ScrapeSceneQuery = {
  __typename?: "Query";

  scrapeScene: Maybe<ScrapeSceneScrapeScene>;
};

export type ScrapeSceneScrapeScene = ScrapedSceneDataFragment;

export type ScrapeSceneUrlVariables = {
  url: string;
};

export type ScrapeSceneUrlQuery = {
  __typename?: "Query";

  scrapeSceneURL: Maybe<ScrapeSceneUrlScrapeSceneUrl>;
};

export type ScrapeSceneUrlScrapeSceneUrl = ScrapedSceneDataFragment;

export type ConfigurationVariables = {};

export type ConfigurationQuery = {
  __typename?: "Query";

  configuration: ConfigurationConfiguration;
};

export type ConfigurationConfiguration = ConfigDataFragment;

export type DirectoriesVariables = {
  path?: Maybe<string>;
};

export type DirectoriesQuery = {
  __typename?: "Query";

  directories: string[];
};

export type MetadataImportVariables = {};

export type MetadataImportQuery = {
  __typename?: "Query";

  metadataImport: string;
};

export type MetadataExportVariables = {};

export type MetadataExportQuery = {
  __typename?: "Query";

  metadataExport: string;
};

export type MetadataScanVariables = {
  input: ScanMetadataInput;
};

export type MetadataScanQuery = {
  __typename?: "Query";

  metadataScan: string;
};

export type MetadataGenerateVariables = {
  input: GenerateMetadataInput;
};

export type MetadataGenerateQuery = {
  __typename?: "Query";

  metadataGenerate: string;
};

export type MetadataAutoTagVariables = {
  input: AutoTagMetadataInput;
};

export type MetadataAutoTagQuery = {
  __typename?: "Query";

  metadataAutoTag: string;
};

export type MetadataCleanVariables = {};

export type MetadataCleanQuery = {
  __typename?: "Query";

  metadataClean: string;
};

export type JobStatusVariables = {};

export type JobStatusQuery = {
  __typename?: "Query";

  jobStatus: JobStatusJobStatus;
};

export type JobStatusJobStatus = {
  __typename?: "MetadataUpdateStatus";

  progress: number;

  status: string;

  message: string;
};

export type StopJobVariables = {};

export type StopJobQuery = {
  __typename?: "Query";

  stopJob: boolean;
};

export type FindStudiosVariables = {
  filter?: Maybe<FindFilterType>;
};

export type FindStudiosQuery = {
  __typename?: "Query";

  findStudios: FindStudiosFindStudios;
};

export type FindStudiosFindStudios = {
  __typename?: "FindStudiosResultType";

  count: number;

  studios: FindStudiosStudios[];
};

export type FindStudiosStudios = StudioDataFragment;

export type FindStudioVariables = {
  id: string;
};

export type FindStudioQuery = {
  __typename?: "Query";

  findStudio: Maybe<FindStudioFindStudio>;
};

export type FindStudioFindStudio = StudioDataFragment;

export type MetadataUpdateVariables = {};

export type MetadataUpdateSubscription = {
  __typename?: "Subscription";

  metadataUpdate: MetadataUpdateMetadataUpdate;
};

export type MetadataUpdateMetadataUpdate = {
  __typename?: "MetadataUpdateStatus";

  progress: number;

  status: string;

  message: string;
};

export type LoggingSubscribeVariables = {};

export type LoggingSubscribeSubscription = {
  __typename?: "Subscription";

  loggingSubscribe: LoggingSubscribeLoggingSubscribe[];
};

export type LoggingSubscribeLoggingSubscribe = LogEntryDataFragment;

export type ConfigGeneralDataFragment = {
  __typename?: "ConfigGeneralResult";

  stashes: string[];

  databasePath: string;

  generatedPath: string;

  maxTranscodeSize: Maybe<StreamingResolutionEnum>;

  maxStreamingTranscodeSize: Maybe<StreamingResolutionEnum>;

  username: string;

  password: string;

  logFile: Maybe<string>;

  logOut: boolean;

  logLevel: string;

  logAccess: boolean;

  excludes: string[];
};

export type ConfigInterfaceDataFragment = {
  __typename?: "ConfigInterfaceResult";

  soundOnPreview: Maybe<boolean>;

  wallShowTitle: Maybe<boolean>;

  maximumLoopDuration: Maybe<number>;

  autostartVideo: Maybe<boolean>;

  showStudioAsText: Maybe<boolean>;

  css: Maybe<string>;

  cssEnabled: Maybe<boolean>;
};

export type ConfigDataFragment = {
  __typename?: "ConfigResult";

  general: ConfigDataGeneral;

  interface: ConfigDataInterface;
};

export type ConfigDataGeneral = ConfigGeneralDataFragment;

export type ConfigDataInterface = ConfigInterfaceDataFragment;

export type GalleryDataFragment = {
  __typename?: "Gallery";

  id: string;

  checksum: string;

  path: string;

  title: Maybe<string>;

  files: GalleryDataFiles[];
};

export type GalleryDataFiles = {
  __typename?: "GalleryFilesType";

  index: number;

  name: Maybe<string>;

  path: Maybe<string>;
};

export type LogEntryDataFragment = {
  __typename?: "LogEntry";

  time: Time;

  level: LogLevel;

  message: string;
};

export type SlimPerformerDataFragment = {
  __typename?: "Performer";

  id: string;

  name: Maybe<string>;

  image_path: Maybe<string>;
};

export type PerformerDataFragment = {
  __typename?: "Performer";

  id: string;

  checksum: string;

  name: Maybe<string>;

  url: Maybe<string>;

  twitter: Maybe<string>;

  instagram: Maybe<string>;

  birthdate: Maybe<string>;

  ethnicity: Maybe<string>;

  country: Maybe<string>;

  eye_color: Maybe<string>;

  height: Maybe<string>;

  measurements: Maybe<string>;

  fake_tits: Maybe<string>;

  career_length: Maybe<string>;

  tattoos: Maybe<string>;

  piercings: Maybe<string>;

  aliases: Maybe<string>;

  favorite: boolean;

  image_path: Maybe<string>;

  scene_count: Maybe<number>;
};

export type SceneMarkerDataFragment = {
  __typename?: "SceneMarker";

  id: string;

  title: string;

  seconds: number;

  stream: string;

  preview: string;

  scene: SceneMarkerDataScene;

  primary_tag: SceneMarkerDataPrimaryTag;

  tags: SceneMarkerDataTags[];
};

export type SceneMarkerDataScene = {
  __typename?: "Scene";

  id: string;
};

export type SceneMarkerDataPrimaryTag = {
  __typename?: "Tag";

  id: string;

  name: string;
};

export type SceneMarkerDataTags = {
  __typename?: "Tag";

  id: string;

  name: string;
};

export type SlimSceneDataFragment = {
  __typename?: "Scene";

  id: string;

  checksum: string;

  title: Maybe<string>;

  details: Maybe<string>;

  url: Maybe<string>;

  date: Maybe<string>;

  rating: Maybe<number>;

  path: string;

  file: SlimSceneDataFile;

  paths: SlimSceneDataPaths;

  scene_markers: SlimSceneDataSceneMarkers[];

  gallery: Maybe<SlimSceneDataGallery>;

  studio: Maybe<SlimSceneDataStudio>;

  tags: SlimSceneDataTags[];

  performers: SlimSceneDataPerformers[];
};

export type SlimSceneDataFile = {
  __typename?: "SceneFileType";

  size: Maybe<string>;

  duration: Maybe<number>;

  video_codec: Maybe<string>;

  audio_codec: Maybe<string>;

  width: Maybe<number>;

  height: Maybe<number>;

  framerate: Maybe<number>;

  bitrate: Maybe<number>;
};

export type SlimSceneDataPaths = {
  __typename?: "ScenePathsType";

  screenshot: Maybe<string>;

  preview: Maybe<string>;

  stream: Maybe<string>;

  webp: Maybe<string>;

  vtt: Maybe<string>;

  chapters_vtt: Maybe<string>;
};

export type SlimSceneDataSceneMarkers = {
  __typename?: "SceneMarker";

  id: string;

  title: string;

  seconds: number;
};

export type SlimSceneDataGallery = {
  __typename?: "Gallery";

  id: string;

  path: string;

  title: Maybe<string>;
};

export type SlimSceneDataStudio = {
  __typename?: "Studio";

  id: string;

  name: string;

  image_path: Maybe<string>;
};

export type SlimSceneDataTags = {
  __typename?: "Tag";

  id: string;

  name: string;
};

export type SlimSceneDataPerformers = {
  __typename?: "Performer";

  id: string;

  name: Maybe<string>;

  favorite: boolean;

  image_path: Maybe<string>;
};

export type SceneDataFragment = {
  __typename?: "Scene";

  id: string;

  checksum: string;

  title: Maybe<string>;

  details: Maybe<string>;

  url: Maybe<string>;

  date: Maybe<string>;

  rating: Maybe<number>;

  path: string;

  file: SceneDataFile;

  paths: SceneDataPaths;

  scene_markers: SceneDataSceneMarkers[];

  is_streamable: boolean;

  gallery: Maybe<SceneDataGallery>;

  studio: Maybe<SceneDataStudio>;

  tags: SceneDataTags[];

  performers: SceneDataPerformers[];
};

export type SceneDataFile = {
  __typename?: "SceneFileType";

  size: Maybe<string>;

  duration: Maybe<number>;

  video_codec: Maybe<string>;

  audio_codec: Maybe<string>;

  width: Maybe<number>;

  height: Maybe<number>;

  framerate: Maybe<number>;

  bitrate: Maybe<number>;
};

export type SceneDataPaths = {
  __typename?: "ScenePathsType";

  screenshot: Maybe<string>;

  preview: Maybe<string>;

  stream: Maybe<string>;

  webp: Maybe<string>;

  vtt: Maybe<string>;

  chapters_vtt: Maybe<string>;
};

export type SceneDataSceneMarkers = SceneMarkerDataFragment;

export type SceneDataGallery = GalleryDataFragment;

export type SceneDataStudio = StudioDataFragment;

export type SceneDataTags = TagDataFragment;

export type SceneDataPerformers = PerformerDataFragment;

export type ScrapedPerformerDataFragment = {
  __typename?: "ScrapedPerformer";

  name: Maybe<string>;

  url: Maybe<string>;

  birthdate: Maybe<string>;

  ethnicity: Maybe<string>;

  country: Maybe<string>;

  eye_color: Maybe<string>;

  height: Maybe<string>;

  measurements: Maybe<string>;

  fake_tits: Maybe<string>;

  career_length: Maybe<string>;

  tattoos: Maybe<string>;

  piercings: Maybe<string>;

  aliases: Maybe<string>;
};

export type ScrapedScenePerformerDataFragment = {
  __typename?: "ScrapedScenePerformer";

  id: Maybe<string>;

  name: string;

  url: Maybe<string>;

  twitter: Maybe<string>;

  instagram: Maybe<string>;

  birthdate: Maybe<string>;

  ethnicity: Maybe<string>;

  country: Maybe<string>;

  eye_color: Maybe<string>;

  height: Maybe<string>;

  measurements: Maybe<string>;

  fake_tits: Maybe<string>;

  career_length: Maybe<string>;

  tattoos: Maybe<string>;

  piercings: Maybe<string>;

  aliases: Maybe<string>;
};

export type ScrapedSceneStudioDataFragment = {
  __typename?: "ScrapedSceneStudio";

  id: Maybe<string>;

  name: string;

  url: Maybe<string>;
};

export type ScrapedSceneTagDataFragment = {
  __typename?: "ScrapedSceneTag";

  id: Maybe<string>;

  name: string;
};

export type ScrapedSceneDataFragment = {
  __typename?: "ScrapedScene";

  title: Maybe<string>;

  details: Maybe<string>;

  url: Maybe<string>;

  date: Maybe<string>;

  file: Maybe<ScrapedSceneDataFile>;

  studio: Maybe<ScrapedSceneDataStudio>;

  tags: Maybe<ScrapedSceneDataTags[]>;

  performers: Maybe<ScrapedSceneDataPerformers[]>;
};

export type ScrapedSceneDataFile = {
  __typename?: "SceneFileType";

  size: Maybe<string>;

  duration: Maybe<number>;

  video_codec: Maybe<string>;

  audio_codec: Maybe<string>;

  width: Maybe<number>;

  height: Maybe<number>;

  framerate: Maybe<number>;

  bitrate: Maybe<number>;
};

export type ScrapedSceneDataStudio = ScrapedSceneStudioDataFragment;

export type ScrapedSceneDataTags = ScrapedSceneTagDataFragment;

export type ScrapedSceneDataPerformers = ScrapedScenePerformerDataFragment;

export type SlimStudioDataFragment = {
  __typename?: "Studio";

  id: string;

  name: string;

  image_path: Maybe<string>;
};

export type StudioDataFragment = {
  __typename?: "Studio";

  id: string;

  checksum: string;

  name: string;

  url: Maybe<string>;

  image_path: Maybe<string>;

  scene_count: Maybe<number>;
};

export type TagDataFragment = {
  __typename?: "Tag";

  id: string;

  name: string;

  scene_count: Maybe<number>;

  scene_marker_count: Maybe<number>;
};

import gql from "graphql-tag";
import * as ReactApolloHooks from "react-apollo-hooks";

// ====================================================
// Fragments
// ====================================================

export const ConfigGeneralDataFragmentDoc = gql`
  fragment ConfigGeneralData on ConfigGeneralResult {
    stashes
    databasePath
    generatedPath
    maxTranscodeSize
    maxStreamingTranscodeSize
    username
    password
    logFile
    logOut
    logLevel
    logAccess
    excludes
  }
`;

export const ConfigInterfaceDataFragmentDoc = gql`
  fragment ConfigInterfaceData on ConfigInterfaceResult {
    soundOnPreview
    wallShowTitle
    maximumLoopDuration
    autostartVideo
    showStudioAsText
    css
    cssEnabled
  }
`;

export const ConfigDataFragmentDoc = gql`
  fragment ConfigData on ConfigResult {
    general {
      ...ConfigGeneralData
    }
    interface {
      ...ConfigInterfaceData
    }
  }

  ${ConfigGeneralDataFragmentDoc}
  ${ConfigInterfaceDataFragmentDoc}
`;

export const LogEntryDataFragmentDoc = gql`
  fragment LogEntryData on LogEntry {
    time
    level
    message
  }
`;

export const SlimPerformerDataFragmentDoc = gql`
  fragment SlimPerformerData on Performer {
    id
    name
    image_path
  }
`;

export const SlimSceneDataFragmentDoc = gql`
  fragment SlimSceneData on Scene {
    id
    checksum
    title
    details
    url
    date
    rating
    path
    file {
      size
      duration
      video_codec
      audio_codec
      width
      height
      framerate
      bitrate
    }
    paths {
      screenshot
      preview
      stream
      webp
      vtt
      chapters_vtt
    }
    scene_markers {
      id
      title
      seconds
    }
    gallery {
      id
      path
      title
    }
    studio {
      id
      name
      image_path
    }
    tags {
      id
      name
    }
    performers {
      id
      name
      favorite
      image_path
    }
  }
`;

export const SceneMarkerDataFragmentDoc = gql`
  fragment SceneMarkerData on SceneMarker {
    id
    title
    seconds
    stream
    preview
    scene {
      id
    }
    primary_tag {
      id
      name
    }
    tags {
      id
      name
    }
  }
`;

export const GalleryDataFragmentDoc = gql`
  fragment GalleryData on Gallery {
    id
    checksum
    path
    title
    files {
      index
      name
      path
    }
  }
`;

export const StudioDataFragmentDoc = gql`
  fragment StudioData on Studio {
    id
    checksum
    name
    url
    image_path
    scene_count
  }
`;

export const TagDataFragmentDoc = gql`
  fragment TagData on Tag {
    id
    name
    scene_count
    scene_marker_count
  }
`;

export const PerformerDataFragmentDoc = gql`
  fragment PerformerData on Performer {
    id
    checksum
    name
    url
    twitter
    instagram
    birthdate
    ethnicity
    country
    eye_color
    height
    measurements
    fake_tits
    career_length
    tattoos
    piercings
    aliases
    favorite
    image_path
    scene_count
  }
`;

export const SceneDataFragmentDoc = gql`
  fragment SceneData on Scene {
    id
    checksum
    title
    details
    url
    date
    rating
    path
    file {
      size
      duration
      video_codec
      audio_codec
      width
      height
      framerate
      bitrate
    }
    paths {
      screenshot
      preview
      stream
      webp
      vtt
      chapters_vtt
    }
    scene_markers {
      ...SceneMarkerData
    }
    is_streamable
    gallery {
      ...GalleryData
    }
    studio {
      ...StudioData
    }
    tags {
      ...TagData
    }
    performers {
      ...PerformerData
    }
  }

  ${SceneMarkerDataFragmentDoc}
  ${GalleryDataFragmentDoc}
  ${StudioDataFragmentDoc}
  ${TagDataFragmentDoc}
  ${PerformerDataFragmentDoc}
`;

export const ScrapedPerformerDataFragmentDoc = gql`
  fragment ScrapedPerformerData on ScrapedPerformer {
    name
    url
    birthdate
    ethnicity
    country
    eye_color
    height
    measurements
    fake_tits
    career_length
    tattoos
    piercings
    aliases
  }
`;

export const ScrapedSceneStudioDataFragmentDoc = gql`
  fragment ScrapedSceneStudioData on ScrapedSceneStudio {
    id
    name
    url
  }
`;

export const ScrapedSceneTagDataFragmentDoc = gql`
  fragment ScrapedSceneTagData on ScrapedSceneTag {
    id
    name
  }
`;

export const ScrapedScenePerformerDataFragmentDoc = gql`
  fragment ScrapedScenePerformerData on ScrapedScenePerformer {
    id
    name
    url
    twitter
    instagram
    birthdate
    ethnicity
    country
    eye_color
    height
    measurements
    fake_tits
    career_length
    tattoos
    piercings
    aliases
  }
`;

export const ScrapedSceneDataFragmentDoc = gql`
  fragment ScrapedSceneData on ScrapedScene {
    title
    details
    url
    date
    file {
      size
      duration
      video_codec
      audio_codec
      width
      height
      framerate
      bitrate
    }
    studio {
      ...ScrapedSceneStudioData
    }
    tags {
      ...ScrapedSceneTagData
    }
    performers {
      ...ScrapedScenePerformerData
    }
  }

  ${ScrapedSceneStudioDataFragmentDoc}
  ${ScrapedSceneTagDataFragmentDoc}
  ${ScrapedScenePerformerDataFragmentDoc}
`;

export const SlimStudioDataFragmentDoc = gql`
  fragment SlimStudioData on Studio {
    id
    name
    image_path
  }
`;

// ====================================================
// Components
// ====================================================

export const ConfigureGeneralDocument = gql`
  mutation ConfigureGeneral($input: ConfigGeneralInput!) {
    configureGeneral(input: $input) {
      ...ConfigGeneralData
    }
  }

  ${ConfigGeneralDataFragmentDoc}
`;
export function useConfigureGeneral(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    ConfigureGeneralMutation,
    ConfigureGeneralVariables
  >
) {
  return ReactApolloHooks.useMutation<
    ConfigureGeneralMutation,
    ConfigureGeneralVariables
  >(ConfigureGeneralDocument, baseOptions);
}
export const ConfigureInterfaceDocument = gql`
  mutation ConfigureInterface($input: ConfigInterfaceInput!) {
    configureInterface(input: $input) {
      ...ConfigInterfaceData
    }
  }

  ${ConfigInterfaceDataFragmentDoc}
`;
export function useConfigureInterface(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    ConfigureInterfaceMutation,
    ConfigureInterfaceVariables
  >
) {
  return ReactApolloHooks.useMutation<
    ConfigureInterfaceMutation,
    ConfigureInterfaceVariables
  >(ConfigureInterfaceDocument, baseOptions);
}
export const PerformerCreateDocument = gql`
  mutation PerformerCreate(
    $name: String
    $url: String
    $birthdate: String
    $ethnicity: String
    $country: String
    $eye_color: String
    $height: String
    $measurements: String
    $fake_tits: String
    $career_length: String
    $tattoos: String
    $piercings: String
    $aliases: String
    $twitter: String
    $instagram: String
    $favorite: Boolean
    $image: String
  ) {
    performerCreate(
      input: {
        name: $name
        url: $url
        birthdate: $birthdate
        ethnicity: $ethnicity
        country: $country
        eye_color: $eye_color
        height: $height
        measurements: $measurements
        fake_tits: $fake_tits
        career_length: $career_length
        tattoos: $tattoos
        piercings: $piercings
        aliases: $aliases
        twitter: $twitter
        instagram: $instagram
        favorite: $favorite
        image: $image
      }
    ) {
      ...PerformerData
    }
  }

  ${PerformerDataFragmentDoc}
`;
export function usePerformerCreate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    PerformerCreateMutation,
    PerformerCreateVariables
  >
) {
  return ReactApolloHooks.useMutation<
    PerformerCreateMutation,
    PerformerCreateVariables
  >(PerformerCreateDocument, baseOptions);
}
export const PerformerUpdateDocument = gql`
  mutation PerformerUpdate(
    $id: ID!
    $name: String
    $url: String
    $birthdate: String
    $ethnicity: String
    $country: String
    $eye_color: String
    $height: String
    $measurements: String
    $fake_tits: String
    $career_length: String
    $tattoos: String
    $piercings: String
    $aliases: String
    $twitter: String
    $instagram: String
    $favorite: Boolean
    $image: String
  ) {
    performerUpdate(
      input: {
        id: $id
        name: $name
        url: $url
        birthdate: $birthdate
        ethnicity: $ethnicity
        country: $country
        eye_color: $eye_color
        height: $height
        measurements: $measurements
        fake_tits: $fake_tits
        career_length: $career_length
        tattoos: $tattoos
        piercings: $piercings
        aliases: $aliases
        twitter: $twitter
        instagram: $instagram
        favorite: $favorite
        image: $image
      }
    ) {
      ...PerformerData
    }
  }

  ${PerformerDataFragmentDoc}
`;
export function usePerformerUpdate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    PerformerUpdateMutation,
    PerformerUpdateVariables
  >
) {
  return ReactApolloHooks.useMutation<
    PerformerUpdateMutation,
    PerformerUpdateVariables
  >(PerformerUpdateDocument, baseOptions);
}
export const PerformerDestroyDocument = gql`
  mutation PerformerDestroy($id: ID!) {
    performerDestroy(input: { id: $id })
  }
`;
export function usePerformerDestroy(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    PerformerDestroyMutation,
    PerformerDestroyVariables
  >
) {
  return ReactApolloHooks.useMutation<
    PerformerDestroyMutation,
    PerformerDestroyVariables
  >(PerformerDestroyDocument, baseOptions);
}
export const SceneMarkerCreateDocument = gql`
  mutation SceneMarkerCreate(
    $title: String!
    $seconds: Float!
    $scene_id: ID!
    $primary_tag_id: ID!
    $tag_ids: [ID!] = []
  ) {
    sceneMarkerCreate(
      input: {
        title: $title
        seconds: $seconds
        scene_id: $scene_id
        primary_tag_id: $primary_tag_id
        tag_ids: $tag_ids
      }
    ) {
      ...SceneMarkerData
    }
  }

  ${SceneMarkerDataFragmentDoc}
`;
export function useSceneMarkerCreate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    SceneMarkerCreateMutation,
    SceneMarkerCreateVariables
  >
) {
  return ReactApolloHooks.useMutation<
    SceneMarkerCreateMutation,
    SceneMarkerCreateVariables
  >(SceneMarkerCreateDocument, baseOptions);
}
export const SceneMarkerUpdateDocument = gql`
  mutation SceneMarkerUpdate(
    $id: ID!
    $title: String!
    $seconds: Float!
    $scene_id: ID!
    $primary_tag_id: ID!
    $tag_ids: [ID!] = []
  ) {
    sceneMarkerUpdate(
      input: {
        id: $id
        title: $title
        seconds: $seconds
        scene_id: $scene_id
        primary_tag_id: $primary_tag_id
        tag_ids: $tag_ids
      }
    ) {
      ...SceneMarkerData
    }
  }

  ${SceneMarkerDataFragmentDoc}
`;
export function useSceneMarkerUpdate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    SceneMarkerUpdateMutation,
    SceneMarkerUpdateVariables
  >
) {
  return ReactApolloHooks.useMutation<
    SceneMarkerUpdateMutation,
    SceneMarkerUpdateVariables
  >(SceneMarkerUpdateDocument, baseOptions);
}
export const SceneMarkerDestroyDocument = gql`
  mutation SceneMarkerDestroy($id: ID!) {
    sceneMarkerDestroy(id: $id)
  }
`;
export function useSceneMarkerDestroy(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    SceneMarkerDestroyMutation,
    SceneMarkerDestroyVariables
  >
) {
  return ReactApolloHooks.useMutation<
    SceneMarkerDestroyMutation,
    SceneMarkerDestroyVariables
  >(SceneMarkerDestroyDocument, baseOptions);
}
export const SceneUpdateDocument = gql`
  mutation SceneUpdate(
    $id: ID!
    $title: String
    $details: String
    $url: String
    $date: String
    $rating: Int
    $studio_id: ID
    $gallery_id: ID
    $performer_ids: [ID!] = []
    $tag_ids: [ID!] = []
    $cover_image: String
  ) {
    sceneUpdate(
      input: {
        id: $id
        title: $title
        details: $details
        url: $url
        date: $date
        rating: $rating
        studio_id: $studio_id
        gallery_id: $gallery_id
        performer_ids: $performer_ids
        tag_ids: $tag_ids
        cover_image: $cover_image
      }
    ) {
      ...SceneData
    }
  }

  ${SceneDataFragmentDoc}
`;
export function useSceneUpdate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    SceneUpdateMutation,
    SceneUpdateVariables
  >
) {
  return ReactApolloHooks.useMutation<
    SceneUpdateMutation,
    SceneUpdateVariables
  >(SceneUpdateDocument, baseOptions);
}
export const BulkSceneUpdateDocument = gql`
  mutation BulkSceneUpdate(
    $ids: [ID!] = []
    $title: String
    $details: String
    $url: String
    $date: String
    $rating: Int
    $studio_id: ID
    $gallery_id: ID
    $performer_ids: [ID!]
    $tag_ids: [ID!]
  ) {
    bulkSceneUpdate(
      input: {
        ids: $ids
        title: $title
        details: $details
        url: $url
        date: $date
        rating: $rating
        studio_id: $studio_id
        gallery_id: $gallery_id
        performer_ids: $performer_ids
        tag_ids: $tag_ids
      }
    ) {
      ...SceneData
    }
  }

  ${SceneDataFragmentDoc}
`;
export function useBulkSceneUpdate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    BulkSceneUpdateMutation,
    BulkSceneUpdateVariables
  >
) {
  return ReactApolloHooks.useMutation<
    BulkSceneUpdateMutation,
    BulkSceneUpdateVariables
  >(BulkSceneUpdateDocument, baseOptions);
}
export const ScenesUpdateDocument = gql`
  mutation ScenesUpdate($input: [SceneUpdateInput!]!) {
    scenesUpdate(input: $input) {
      ...SceneData
    }
  }

  ${SceneDataFragmentDoc}
`;
export function useScenesUpdate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    ScenesUpdateMutation,
    ScenesUpdateVariables
  >
) {
  return ReactApolloHooks.useMutation<
    ScenesUpdateMutation,
    ScenesUpdateVariables
  >(ScenesUpdateDocument, baseOptions);
}
export const SceneDestroyDocument = gql`
  mutation SceneDestroy(
    $id: ID!
    $delete_file: Boolean
    $delete_generated: Boolean
  ) {
    sceneDestroy(
      input: {
        id: $id
        delete_file: $delete_file
        delete_generated: $delete_generated
      }
    )
  }
`;
export function useSceneDestroy(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    SceneDestroyMutation,
    SceneDestroyVariables
  >
) {
  return ReactApolloHooks.useMutation<
    SceneDestroyMutation,
    SceneDestroyVariables
  >(SceneDestroyDocument, baseOptions);
}
export const StudioCreateDocument = gql`
  mutation StudioCreate($name: String!, $url: String, $image: String) {
    studioCreate(input: { name: $name, url: $url, image: $image }) {
      ...StudioData
    }
  }

  ${StudioDataFragmentDoc}
`;
export function useStudioCreate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    StudioCreateMutation,
    StudioCreateVariables
  >
) {
  return ReactApolloHooks.useMutation<
    StudioCreateMutation,
    StudioCreateVariables
  >(StudioCreateDocument, baseOptions);
}
export const StudioUpdateDocument = gql`
  mutation StudioUpdate($id: ID!, $name: String, $url: String, $image: String) {
    studioUpdate(input: { id: $id, name: $name, url: $url, image: $image }) {
      ...StudioData
    }
  }

  ${StudioDataFragmentDoc}
`;
export function useStudioUpdate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    StudioUpdateMutation,
    StudioUpdateVariables
  >
) {
  return ReactApolloHooks.useMutation<
    StudioUpdateMutation,
    StudioUpdateVariables
  >(StudioUpdateDocument, baseOptions);
}
export const StudioDestroyDocument = gql`
  mutation StudioDestroy($id: ID!) {
    studioDestroy(input: { id: $id })
  }
`;
export function useStudioDestroy(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    StudioDestroyMutation,
    StudioDestroyVariables
  >
) {
  return ReactApolloHooks.useMutation<
    StudioDestroyMutation,
    StudioDestroyVariables
  >(StudioDestroyDocument, baseOptions);
}
export const TagCreateDocument = gql`
  mutation TagCreate($name: String!) {
    tagCreate(input: { name: $name }) {
      ...TagData
    }
  }

  ${TagDataFragmentDoc}
`;
export function useTagCreate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    TagCreateMutation,
    TagCreateVariables
  >
) {
  return ReactApolloHooks.useMutation<TagCreateMutation, TagCreateVariables>(
    TagCreateDocument,
    baseOptions
  );
}
export const TagDestroyDocument = gql`
  mutation TagDestroy($id: ID!) {
    tagDestroy(input: { id: $id })
  }
`;
export function useTagDestroy(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    TagDestroyMutation,
    TagDestroyVariables
  >
) {
  return ReactApolloHooks.useMutation<TagDestroyMutation, TagDestroyVariables>(
    TagDestroyDocument,
    baseOptions
  );
}
export const TagUpdateDocument = gql`
  mutation TagUpdate($id: ID!, $name: String!) {
    tagUpdate(input: { id: $id, name: $name }) {
      ...TagData
    }
  }

  ${TagDataFragmentDoc}
`;
export function useTagUpdate(
  baseOptions?: ReactApolloHooks.MutationHookOptions<
    TagUpdateMutation,
    TagUpdateVariables
  >
) {
  return ReactApolloHooks.useMutation<TagUpdateMutation, TagUpdateVariables>(
    TagUpdateDocument,
    baseOptions
  );
}
export const FindGalleriesDocument = gql`
  query FindGalleries($filter: FindFilterType) {
    findGalleries(filter: $filter) {
      count
      galleries {
        ...GalleryData
      }
    }
  }

  ${GalleryDataFragmentDoc}
`;
export function useFindGalleries(
  baseOptions?: ReactApolloHooks.QueryHookOptions<FindGalleriesVariables>
) {
  return ReactApolloHooks.useQuery<FindGalleriesQuery, FindGalleriesVariables>(
    FindGalleriesDocument,
    baseOptions
  );
}
export const FindGalleryDocument = gql`
  query FindGallery($id: ID!) {
    findGallery(id: $id) {
      ...GalleryData
    }
  }

  ${GalleryDataFragmentDoc}
`;
export function useFindGallery(
  baseOptions?: ReactApolloHooks.QueryHookOptions<FindGalleryVariables>
) {
  return ReactApolloHooks.useQuery<FindGalleryQuery, FindGalleryVariables>(
    FindGalleryDocument,
    baseOptions
  );
}
export const SceneWallDocument = gql`
  query SceneWall($q: String) {
    sceneWall(q: $q) {
      ...SceneData
    }
  }

  ${SceneDataFragmentDoc}
`;
export function useSceneWall(
  baseOptions?: ReactApolloHooks.QueryHookOptions<SceneWallVariables>
) {
  return ReactApolloHooks.useQuery<SceneWallQuery, SceneWallVariables>(
    SceneWallDocument,
    baseOptions
  );
}
export const MarkerWallDocument = gql`
  query MarkerWall($q: String) {
    markerWall(q: $q) {
      ...SceneMarkerData
    }
  }

  ${SceneMarkerDataFragmentDoc}
`;
export function useMarkerWall(
  baseOptions?: ReactApolloHooks.QueryHookOptions<MarkerWallVariables>
) {
  return ReactApolloHooks.useQuery<MarkerWallQuery, MarkerWallVariables>(
    MarkerWallDocument,
    baseOptions
  );
}
export const FindTagDocument = gql`
  query FindTag($id: ID!) {
    findTag(id: $id) {
      ...TagData
    }
  }

  ${TagDataFragmentDoc}
`;
export function useFindTag(
  baseOptions?: ReactApolloHooks.QueryHookOptions<FindTagVariables>
) {
  return ReactApolloHooks.useQuery<FindTagQuery, FindTagVariables>(
    FindTagDocument,
    baseOptions
  );
}
export const MarkerStringsDocument = gql`
  query MarkerStrings($q: String, $sort: String) {
    markerStrings(q: $q, sort: $sort) {
      id
      count
      title
    }
  }
`;
export function useMarkerStrings(
  baseOptions?: ReactApolloHooks.QueryHookOptions<MarkerStringsVariables>
) {
  return ReactApolloHooks.useQuery<MarkerStringsQuery, MarkerStringsVariables>(
    MarkerStringsDocument,
    baseOptions
  );
}
export const AllTagsDocument = gql`
  query AllTags {
    allTags {
      ...TagData
    }
  }

  ${TagDataFragmentDoc}
`;
export function useAllTags(
  baseOptions?: ReactApolloHooks.QueryHookOptions<AllTagsVariables>
) {
  return ReactApolloHooks.useQuery<AllTagsQuery, AllTagsVariables>(
    AllTagsDocument,
    baseOptions
  );
}
export const AllPerformersForFilterDocument = gql`
  query AllPerformersForFilter {
    allPerformers {
      ...SlimPerformerData
    }
  }

  ${SlimPerformerDataFragmentDoc}
`;
export function useAllPerformersForFilter(
  baseOptions?: ReactApolloHooks.QueryHookOptions<
    AllPerformersForFilterVariables
  >
) {
  return ReactApolloHooks.useQuery<
    AllPerformersForFilterQuery,
    AllPerformersForFilterVariables
  >(AllPerformersForFilterDocument, baseOptions);
}
export const AllStudiosForFilterDocument = gql`
  query AllStudiosForFilter {
    allStudios {
      ...SlimStudioData
    }
  }

  ${SlimStudioDataFragmentDoc}
`;
export function useAllStudiosForFilter(
  baseOptions?: ReactApolloHooks.QueryHookOptions<AllStudiosForFilterVariables>
) {
  return ReactApolloHooks.useQuery<
    AllStudiosForFilterQuery,
    AllStudiosForFilterVariables
  >(AllStudiosForFilterDocument, baseOptions);
}
export const AllTagsForFilterDocument = gql`
  query AllTagsForFilter {
    allTags {
      id
      name
    }
  }
`;
export function useAllTagsForFilter(
  baseOptions?: ReactApolloHooks.QueryHookOptions<AllTagsForFilterVariables>
) {
  return ReactApolloHooks.useQuery<
    AllTagsForFilterQuery,
    AllTagsForFilterVariables
  >(AllTagsForFilterDocument, baseOptions);
}
export const ValidGalleriesForSceneDocument = gql`
  query ValidGalleriesForScene($scene_id: ID!) {
    validGalleriesForScene(scene_id: $scene_id) {
      id
      path
    }
  }
`;
export function useValidGalleriesForScene(
  baseOptions?: ReactApolloHooks.QueryHookOptions<
    ValidGalleriesForSceneVariables
  >
) {
  return ReactApolloHooks.useQuery<
    ValidGalleriesForSceneQuery,
    ValidGalleriesForSceneVariables
  >(ValidGalleriesForSceneDocument, baseOptions);
}
export const StatsDocument = gql`
  query Stats {
    stats {
      scene_count
      gallery_count
      performer_count
      studio_count
      tag_count
    }
  }
`;
export function useStats(
  baseOptions?: ReactApolloHooks.QueryHookOptions<StatsVariables>
) {
  return ReactApolloHooks.useQuery<StatsQuery, StatsVariables>(
    StatsDocument,
    baseOptions
  );
}
export const LogsDocument = gql`
  query Logs {
    logs {
      ...LogEntryData
    }
  }

  ${LogEntryDataFragmentDoc}
`;
export function useLogs(
  baseOptions?: ReactApolloHooks.QueryHookOptions<LogsVariables>
) {
  return ReactApolloHooks.useQuery<LogsQuery, LogsVariables>(
    LogsDocument,
    baseOptions
  );
}
export const VersionDocument = gql`
  query Version {
    version {
      version
      hash
      build_time
    }
  }
`;
export function useVersion(
  baseOptions?: ReactApolloHooks.QueryHookOptions<VersionVariables>
) {
  return ReactApolloHooks.useQuery<VersionQuery, VersionVariables>(
    VersionDocument,
    baseOptions
  );
}
export const FindPerformersDocument = gql`
  query FindPerformers(
    $filter: FindFilterType
    $performer_filter: PerformerFilterType
  ) {
    findPerformers(filter: $filter, performer_filter: $performer_filter) {
      count
      performers {
        ...PerformerData
      }
    }
  }

  ${PerformerDataFragmentDoc}
`;
export function useFindPerformers(
  baseOptions?: ReactApolloHooks.QueryHookOptions<FindPerformersVariables>
) {
  return ReactApolloHooks.useQuery<
    FindPerformersQuery,
    FindPerformersVariables
  >(FindPerformersDocument, baseOptions);
}
export const FindPerformerDocument = gql`
  query FindPerformer($id: ID!) {
    findPerformer(id: $id) {
      ...PerformerData
    }
  }

  ${PerformerDataFragmentDoc}
`;
export function useFindPerformer(
  baseOptions?: ReactApolloHooks.QueryHookOptions<FindPerformerVariables>
) {
  return ReactApolloHooks.useQuery<FindPerformerQuery, FindPerformerVariables>(
    FindPerformerDocument,
    baseOptions
  );
}
export const FindSceneMarkersDocument = gql`
  query FindSceneMarkers(
    $filter: FindFilterType
    $scene_marker_filter: SceneMarkerFilterType
  ) {
    findSceneMarkers(
      filter: $filter
      scene_marker_filter: $scene_marker_filter
    ) {
      count
      scene_markers {
        ...SceneMarkerData
      }
    }
  }

  ${SceneMarkerDataFragmentDoc}
`;
export function useFindSceneMarkers(
  baseOptions?: ReactApolloHooks.QueryHookOptions<FindSceneMarkersVariables>
) {
  return ReactApolloHooks.useQuery<
    FindSceneMarkersQuery,
    FindSceneMarkersVariables
  >(FindSceneMarkersDocument, baseOptions);
}
export const FindScenesDocument = gql`
  query FindScenes(
    $filter: FindFilterType
    $scene_filter: SceneFilterType
    $scene_ids: [Int!]
  ) {
    findScenes(
      filter: $filter
      scene_filter: $scene_filter
      scene_ids: $scene_ids
    ) {
      count
      scenes {
        ...SlimSceneData
      }
    }
  }

  ${SlimSceneDataFragmentDoc}
`;
export function useFindScenes(
  baseOptions?: ReactApolloHooks.QueryHookOptions<FindScenesVariables>
) {
  return ReactApolloHooks.useQuery<FindScenesQuery, FindScenesVariables>(
    FindScenesDocument,
    baseOptions
  );
}
export const FindScenesByPathRegexDocument = gql`
  query FindScenesByPathRegex($filter: FindFilterType) {
    findScenesByPathRegex(filter: $filter) {
      count
      scenes {
        ...SlimSceneData
      }
    }
  }

  ${SlimSceneDataFragmentDoc}
`;
export function useFindScenesByPathRegex(
  baseOptions?: ReactApolloHooks.QueryHookOptions<
    FindScenesByPathRegexVariables
  >
) {
  return ReactApolloHooks.useQuery<
    FindScenesByPathRegexQuery,
    FindScenesByPathRegexVariables
  >(FindScenesByPathRegexDocument, baseOptions);
}
export const FindSceneDocument = gql`
  query FindScene($id: ID!, $checksum: String) {
    findScene(id: $id, checksum: $checksum) {
      ...SceneData
    }
    sceneMarkerTags(scene_id: $id) {
      tag {
        id
        name
      }
      scene_markers {
        ...SceneMarkerData
      }
    }
  }

  ${SceneDataFragmentDoc}
  ${SceneMarkerDataFragmentDoc}
`;
export function useFindScene(
  baseOptions?: ReactApolloHooks.QueryHookOptions<FindSceneVariables>
) {
  return ReactApolloHooks.useQuery<FindSceneQuery, FindSceneVariables>(
    FindSceneDocument,
    baseOptions
  );
}
export const ParseSceneFilenamesDocument = gql`
  query ParseSceneFilenames(
    $filter: FindFilterType!
    $config: SceneParserInput!
  ) {
    parseSceneFilenames(filter: $filter, config: $config) {
      count
      results {
        scene {
          ...SlimSceneData
        }
        title
        details
        url
        date
        rating
        studio_id
        gallery_id
        performer_ids
        tag_ids
      }
    }
  }

  ${SlimSceneDataFragmentDoc}
`;
export function useParseSceneFilenames(
  baseOptions?: ReactApolloHooks.QueryHookOptions<ParseSceneFilenamesVariables>
) {
  return ReactApolloHooks.useQuery<
    ParseSceneFilenamesQuery,
    ParseSceneFilenamesVariables
  >(ParseSceneFilenamesDocument, baseOptions);
}
export const ScrapeFreeonesDocument = gql`
  query ScrapeFreeones($performer_name: String!) {
    scrapeFreeones(performer_name: $performer_name) {
      name
      url
      twitter
      instagram
      birthdate
      ethnicity
      country
      eye_color
      height
      measurements
      fake_tits
      career_length
      tattoos
      piercings
      aliases
    }
  }
`;
export function useScrapeFreeones(
  baseOptions?: ReactApolloHooks.QueryHookOptions<ScrapeFreeonesVariables>
) {
  return ReactApolloHooks.useQuery<
    ScrapeFreeonesQuery,
    ScrapeFreeonesVariables
  >(ScrapeFreeonesDocument, baseOptions);
}
export const ScrapeFreeonesPerformersDocument = gql`
  query ScrapeFreeonesPerformers($q: String!) {
    scrapeFreeonesPerformerList(query: $q)
  }
`;
export function useScrapeFreeonesPerformers(
  baseOptions?: ReactApolloHooks.QueryHookOptions<
    ScrapeFreeonesPerformersVariables
  >
) {
  return ReactApolloHooks.useQuery<
    ScrapeFreeonesPerformersQuery,
    ScrapeFreeonesPerformersVariables
  >(ScrapeFreeonesPerformersDocument, baseOptions);
}
export const ListPerformerScrapersDocument = gql`
  query ListPerformerScrapers {
    listPerformerScrapers {
      id
      name
      performer {
        urls
        supported_scrapes
      }
    }
  }
`;
export function useListPerformerScrapers(
  baseOptions?: ReactApolloHooks.QueryHookOptions<
    ListPerformerScrapersVariables
  >
) {
  return ReactApolloHooks.useQuery<
    ListPerformerScrapersQuery,
    ListPerformerScrapersVariables
  >(ListPerformerScrapersDocument, baseOptions);
}
export const ListSceneScrapersDocument = gql`
  query ListSceneScrapers {
    listSceneScrapers {
      id
      name
      scene {
        urls
        supported_scrapes
      }
    }
  }
`;
export function useListSceneScrapers(
  baseOptions?: ReactApolloHooks.QueryHookOptions<ListSceneScrapersVariables>
) {
  return ReactApolloHooks.useQuery<
    ListSceneScrapersQuery,
    ListSceneScrapersVariables
  >(ListSceneScrapersDocument, baseOptions);
}
export const ScrapePerformerListDocument = gql`
  query ScrapePerformerList($scraper_id: ID!, $query: String!) {
    scrapePerformerList(scraper_id: $scraper_id, query: $query) {
      ...ScrapedPerformerData
    }
  }

  ${ScrapedPerformerDataFragmentDoc}
`;
export function useScrapePerformerList(
  baseOptions?: ReactApolloHooks.QueryHookOptions<ScrapePerformerListVariables>
) {
  return ReactApolloHooks.useQuery<
    ScrapePerformerListQuery,
    ScrapePerformerListVariables
  >(ScrapePerformerListDocument, baseOptions);
}
export const ScrapePerformerDocument = gql`
  query ScrapePerformer(
    $scraper_id: ID!
    $scraped_performer: ScrapedPerformerInput!
  ) {
    scrapePerformer(
      scraper_id: $scraper_id
      scraped_performer: $scraped_performer
    ) {
      ...ScrapedPerformerData
    }
  }

  ${ScrapedPerformerDataFragmentDoc}
`;
export function useScrapePerformer(
  baseOptions?: ReactApolloHooks.QueryHookOptions<ScrapePerformerVariables>
) {
  return ReactApolloHooks.useQuery<
    ScrapePerformerQuery,
    ScrapePerformerVariables
  >(ScrapePerformerDocument, baseOptions);
}
export const ScrapePerformerUrlDocument = gql`
  query ScrapePerformerURL($url: String!) {
    scrapePerformerURL(url: $url) {
      ...ScrapedPerformerData
    }
  }

  ${ScrapedPerformerDataFragmentDoc}
`;
export function useScrapePerformerUrl(
  baseOptions?: ReactApolloHooks.QueryHookOptions<ScrapePerformerUrlVariables>
) {
  return ReactApolloHooks.useQuery<
    ScrapePerformerUrlQuery,
    ScrapePerformerUrlVariables
  >(ScrapePerformerUrlDocument, baseOptions);
}
export const ScrapeSceneDocument = gql`
  query ScrapeScene($scraper_id: ID!, $scene: SceneUpdateInput!) {
    scrapeScene(scraper_id: $scraper_id, scene: $scene) {
      ...ScrapedSceneData
    }
  }

  ${ScrapedSceneDataFragmentDoc}
`;
export function useScrapeScene(
  baseOptions?: ReactApolloHooks.QueryHookOptions<ScrapeSceneVariables>
) {
  return ReactApolloHooks.useQuery<ScrapeSceneQuery, ScrapeSceneVariables>(
    ScrapeSceneDocument,
    baseOptions
  );
}
export const ScrapeSceneUrlDocument = gql`
  query ScrapeSceneURL($url: String!) {
    scrapeSceneURL(url: $url) {
      ...ScrapedSceneData
    }
  }

  ${ScrapedSceneDataFragmentDoc}
`;
export function useScrapeSceneUrl(
  baseOptions?: ReactApolloHooks.QueryHookOptions<ScrapeSceneUrlVariables>
) {
  return ReactApolloHooks.useQuery<
    ScrapeSceneUrlQuery,
    ScrapeSceneUrlVariables
  >(ScrapeSceneUrlDocument, baseOptions);
}
export const ConfigurationDocument = gql`
  query Configuration {
    configuration {
      ...ConfigData
    }
  }

  ${ConfigDataFragmentDoc}
`;
export function useConfiguration(
  baseOptions?: ReactApolloHooks.QueryHookOptions<ConfigurationVariables>
) {
  return ReactApolloHooks.useQuery<ConfigurationQuery, ConfigurationVariables>(
    ConfigurationDocument,
    baseOptions
  );
}
export const DirectoriesDocument = gql`
  query Directories($path: String) {
    directories(path: $path)
  }
`;
export function useDirectories(
  baseOptions?: ReactApolloHooks.QueryHookOptions<DirectoriesVariables>
) {
  return ReactApolloHooks.useQuery<DirectoriesQuery, DirectoriesVariables>(
    DirectoriesDocument,
    baseOptions
  );
}
export const MetadataImportDocument = gql`
  query MetadataImport {
    metadataImport
  }
`;
export function useMetadataImport(
  baseOptions?: ReactApolloHooks.QueryHookOptions<MetadataImportVariables>
) {
  return ReactApolloHooks.useQuery<
    MetadataImportQuery,
    MetadataImportVariables
  >(MetadataImportDocument, baseOptions);
}
export const MetadataExportDocument = gql`
  query MetadataExport {
    metadataExport
  }
`;
export function useMetadataExport(
  baseOptions?: ReactApolloHooks.QueryHookOptions<MetadataExportVariables>
) {
  return ReactApolloHooks.useQuery<
    MetadataExportQuery,
    MetadataExportVariables
  >(MetadataExportDocument, baseOptions);
}
export const MetadataScanDocument = gql`
  query MetadataScan($input: ScanMetadataInput!) {
    metadataScan(input: $input)
  }
`;
export function useMetadataScan(
  baseOptions?: ReactApolloHooks.QueryHookOptions<MetadataScanVariables>
) {
  return ReactApolloHooks.useQuery<MetadataScanQuery, MetadataScanVariables>(
    MetadataScanDocument,
    baseOptions
  );
}
export const MetadataGenerateDocument = gql`
  query MetadataGenerate($input: GenerateMetadataInput!) {
    metadataGenerate(input: $input)
  }
`;
export function useMetadataGenerate(
  baseOptions?: ReactApolloHooks.QueryHookOptions<MetadataGenerateVariables>
) {
  return ReactApolloHooks.useQuery<
    MetadataGenerateQuery,
    MetadataGenerateVariables
  >(MetadataGenerateDocument, baseOptions);
}
export const MetadataAutoTagDocument = gql`
  query MetadataAutoTag($input: AutoTagMetadataInput!) {
    metadataAutoTag(input: $input)
  }
`;
export function useMetadataAutoTag(
  baseOptions?: ReactApolloHooks.QueryHookOptions<MetadataAutoTagVariables>
) {
  return ReactApolloHooks.useQuery<
    MetadataAutoTagQuery,
    MetadataAutoTagVariables
  >(MetadataAutoTagDocument, baseOptions);
}
export const MetadataCleanDocument = gql`
  query MetadataClean {
    metadataClean
  }
`;
export function useMetadataClean(
  baseOptions?: ReactApolloHooks.QueryHookOptions<MetadataCleanVariables>
) {
  return ReactApolloHooks.useQuery<MetadataCleanQuery, MetadataCleanVariables>(
    MetadataCleanDocument,
    baseOptions
  );
}
export const JobStatusDocument = gql`
  query JobStatus {
    jobStatus {
      progress
      status
      message
    }
  }
`;
export function useJobStatus(
  baseOptions?: ReactApolloHooks.QueryHookOptions<JobStatusVariables>
) {
  return ReactApolloHooks.useQuery<JobStatusQuery, JobStatusVariables>(
    JobStatusDocument,
    baseOptions
  );
}
export const StopJobDocument = gql`
  query StopJob {
    stopJob
  }
`;
export function useStopJob(
  baseOptions?: ReactApolloHooks.QueryHookOptions<StopJobVariables>
) {
  return ReactApolloHooks.useQuery<StopJobQuery, StopJobVariables>(
    StopJobDocument,
    baseOptions
  );
}
export const FindStudiosDocument = gql`
  query FindStudios($filter: FindFilterType) {
    findStudios(filter: $filter) {
      count
      studios {
        ...StudioData
      }
    }
  }

  ${StudioDataFragmentDoc}
`;
export function useFindStudios(
  baseOptions?: ReactApolloHooks.QueryHookOptions<FindStudiosVariables>
) {
  return ReactApolloHooks.useQuery<FindStudiosQuery, FindStudiosVariables>(
    FindStudiosDocument,
    baseOptions
  );
}
export const FindStudioDocument = gql`
  query FindStudio($id: ID!) {
    findStudio(id: $id) {
      ...StudioData
    }
  }

  ${StudioDataFragmentDoc}
`;
export function useFindStudio(
  baseOptions?: ReactApolloHooks.QueryHookOptions<FindStudioVariables>
) {
  return ReactApolloHooks.useQuery<FindStudioQuery, FindStudioVariables>(
    FindStudioDocument,
    baseOptions
  );
}
export const MetadataUpdateDocument = gql`
  subscription MetadataUpdate {
    metadataUpdate {
      progress
      status
      message
    }
  }
`;
export function useMetadataUpdate(
  baseOptions?: ReactApolloHooks.SubscriptionHookOptions<
    MetadataUpdateSubscription,
    MetadataUpdateVariables
  >
) {
  return ReactApolloHooks.useSubscription<
    MetadataUpdateSubscription,
    MetadataUpdateVariables
  >(MetadataUpdateDocument, baseOptions);
}
export const LoggingSubscribeDocument = gql`
  subscription LoggingSubscribe {
    loggingSubscribe {
      ...LogEntryData
    }
  }

  ${LogEntryDataFragmentDoc}
`;
export function useLoggingSubscribe(
  baseOptions?: ReactApolloHooks.SubscriptionHookOptions<
    LoggingSubscribeSubscription,
    LoggingSubscribeVariables
  >
) {
  return ReactApolloHooks.useSubscription<
    LoggingSubscribeSubscription,
    LoggingSubscribeVariables
  >(LoggingSubscribeDocument, baseOptions);
}
