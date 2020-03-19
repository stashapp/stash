/* eslint-disable */
import gql from 'graphql-tag';
import * as ApolloReactCommon from '@apollo/react-common';
import * as React from 'react';
import * as ApolloReactComponents from '@apollo/react-components';
import * as ApolloReactHooks from '@apollo/react-hooks';
export type Maybe<T> = T | null;
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;

// Generated in 2020-03-05T11:18:15+11:00

/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string,
  String: string,
  Boolean: boolean,
  Int: number,
  Float: number,
  /** Log entries */
  Time: any,
};

export type AutoTagMetadataInput = {
  /** IDs of performers to tag files with, or "*" for all */
  performers?: Maybe<Array<Scalars['String']>>,
  /** IDs of studios to tag files with, or "*" for all */
  studios?: Maybe<Array<Scalars['String']>>,
  /** IDs of tags to tag files with, or "*" for all */
  tags?: Maybe<Array<Scalars['String']>>,
};

export type BulkSceneUpdateInput = {
  clientMutationId?: Maybe<Scalars['String']>,
  ids?: Maybe<Array<Scalars['ID']>>,
  title?: Maybe<Scalars['String']>,
  details?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  date?: Maybe<Scalars['String']>,
  rating?: Maybe<Scalars['Int']>,
  studio_id?: Maybe<Scalars['ID']>,
  gallery_id?: Maybe<Scalars['ID']>,
  performer_ids?: Maybe<Array<Scalars['ID']>>,
  tag_ids?: Maybe<Array<Scalars['ID']>>,
};

export type ConfigGeneralInput = {
  /** Array of file paths to content */
  stashes?: Maybe<Array<Scalars['String']>>,
  /** Path to the SQLite database */
  databasePath?: Maybe<Scalars['String']>,
  /** Path to generated files */
  generatedPath?: Maybe<Scalars['String']>,
  /** Max generated transcode size */
  maxTranscodeSize?: Maybe<StreamingResolutionEnum>,
  /** Max streaming transcode size */
  maxStreamingTranscodeSize?: Maybe<StreamingResolutionEnum>,
  /** Username */
  username?: Maybe<Scalars['String']>,
  /** Password */
  password?: Maybe<Scalars['String']>,
  /** Name of the log file */
  logFile?: Maybe<Scalars['String']>,
  /** Whether to also output to stderr */
  logOut: Scalars['Boolean'],
  /** Minimum log level */
  logLevel: Scalars['String'],
  /** Whether to log http access */
  logAccess: Scalars['Boolean'],
  /** Array of file regexp to exclude from Scan */
  excludes?: Maybe<Array<Scalars['String']>>,
};

export type ConfigGeneralResult = {
   __typename?: 'ConfigGeneralResult',
  /** Array of file paths to content */
  stashes: Array<Scalars['String']>,
  /** Path to the SQLite database */
  databasePath: Scalars['String'],
  /** Path to generated files */
  generatedPath: Scalars['String'],
  /** Max generated transcode size */
  maxTranscodeSize?: Maybe<StreamingResolutionEnum>,
  /** Max streaming transcode size */
  maxStreamingTranscodeSize?: Maybe<StreamingResolutionEnum>,
  /** Username */
  username: Scalars['String'],
  /** Password */
  password: Scalars['String'],
  /** Name of the log file */
  logFile?: Maybe<Scalars['String']>,
  /** Whether to also output to stderr */
  logOut: Scalars['Boolean'],
  /** Minimum log level */
  logLevel: Scalars['String'],
  /** Whether to log http access */
  logAccess: Scalars['Boolean'],
  /** Array of file regexp to exclude from Scan */
  excludes: Array<Scalars['String']>,
};

export type ConfigInterfaceInput = {
  /** Enable sound on mouseover previews */
  soundOnPreview?: Maybe<Scalars['Boolean']>,
  /** Show title and tags in wall view */
  wallShowTitle?: Maybe<Scalars['Boolean']>,
  /** Maximum duration (in seconds) in which a scene video will loop in the scene player */
  maximumLoopDuration?: Maybe<Scalars['Int']>,
  /** If true, video will autostart on load in the scene player */
  autostartVideo?: Maybe<Scalars['Boolean']>,
  /** If true, studio overlays will be shown as text instead of logo images */
  showStudioAsText?: Maybe<Scalars['Boolean']>,
  /** Custom CSS */
  css?: Maybe<Scalars['String']>,
  cssEnabled?: Maybe<Scalars['Boolean']>,
  language?: Maybe<Scalars['String']>,
};

export type ConfigInterfaceResult = {
   __typename?: 'ConfigInterfaceResult',
  /** Enable sound on mouseover previews */
  soundOnPreview?: Maybe<Scalars['Boolean']>,
  /** Show title and tags in wall view */
  wallShowTitle?: Maybe<Scalars['Boolean']>,
  /** Maximum duration (in seconds) in which a scene video will loop in the scene player */
  maximumLoopDuration?: Maybe<Scalars['Int']>,
  /** If true, video will autostart on load in the scene player */
  autostartVideo?: Maybe<Scalars['Boolean']>,
  /** If true, studio overlays will be shown as text instead of logo images */
  showStudioAsText?: Maybe<Scalars['Boolean']>,
  /** Custom CSS */
  css?: Maybe<Scalars['String']>,
  cssEnabled?: Maybe<Scalars['Boolean']>,
  /** Interface language */
  language?: Maybe<Scalars['String']>,
};

/** All configuration settings */
export type ConfigResult = {
   __typename?: 'ConfigResult',
  general: ConfigGeneralResult,
  interface: ConfigInterfaceResult,
};

export enum CriterionModifier {
  /** = */
  Equals = 'EQUALS',
  /** != */
  NotEquals = 'NOT_EQUALS',
  /** > */
  GreaterThan = 'GREATER_THAN',
  /** < */
  LessThan = 'LESS_THAN',
  /** IS NULL */
  IsNull = 'IS_NULL',
  /** IS NOT NULL */
  NotNull = 'NOT_NULL',
  /** INCLUDES ALL */
  IncludesAll = 'INCLUDES_ALL',
  Includes = 'INCLUDES',
  Excludes = 'EXCLUDES'
}

export type FindFilterType = {
  q?: Maybe<Scalars['String']>,
  page?: Maybe<Scalars['Int']>,
  per_page?: Maybe<Scalars['Int']>,
  sort?: Maybe<Scalars['String']>,
  direction?: Maybe<SortDirectionEnum>,
};

export type FindGalleriesResultType = {
   __typename?: 'FindGalleriesResultType',
  count: Scalars['Int'],
  galleries: Array<Gallery>,
};

export type FindPerformersResultType = {
   __typename?: 'FindPerformersResultType',
  count: Scalars['Int'],
  performers: Array<Performer>,
};

export type FindSceneMarkersResultType = {
   __typename?: 'FindSceneMarkersResultType',
  count: Scalars['Int'],
  scene_markers: Array<SceneMarker>,
};

export type FindScenesResultType = {
   __typename?: 'FindScenesResultType',
  count: Scalars['Int'],
  scenes: Array<Scene>,
};

export type FindStudiosResultType = {
   __typename?: 'FindStudiosResultType',
  count: Scalars['Int'],
  studios: Array<Studio>,
};

/** Gallery type */
export type Gallery = {
   __typename?: 'Gallery',
  id: Scalars['ID'],
  checksum: Scalars['String'],
  path: Scalars['String'],
  title?: Maybe<Scalars['String']>,
  /** The files in the gallery */
  files: Array<GalleryFilesType>,
};

export type GalleryFilesType = {
   __typename?: 'GalleryFilesType',
  index: Scalars['Int'],
  name?: Maybe<Scalars['String']>,
  path?: Maybe<Scalars['String']>,
};

export type GenerateMetadataInput = {
  sprites: Scalars['Boolean'],
  previews: Scalars['Boolean'],
  markers: Scalars['Boolean'],
  transcodes: Scalars['Boolean'],
};

export type IntCriterionInput = {
  value: Scalars['Int'],
  modifier: CriterionModifier,
};

export type LogEntry = {
   __typename?: 'LogEntry',
  time: Scalars['Time'],
  level: LogLevel,
  message: Scalars['String'],
};

export enum LogLevel {
  Debug = 'Debug',
  Info = 'Info',
  Progress = 'Progress',
  Warning = 'Warning',
  Error = 'Error'
}

export type MarkerStringsResultType = {
   __typename?: 'MarkerStringsResultType',
  count: Scalars['Int'],
  id: Scalars['ID'],
  title: Scalars['String'],
};

export type MetadataUpdateStatus = {
   __typename?: 'MetadataUpdateStatus',
  progress: Scalars['Float'],
  status: Scalars['String'],
  message: Scalars['String'],
};

export type MultiCriterionInput = {
  value?: Maybe<Array<Scalars['ID']>>,
  modifier: CriterionModifier,
};

export type Mutation = {
   __typename?: 'Mutation',
  sceneUpdate?: Maybe<Scene>,
  bulkSceneUpdate?: Maybe<Array<Scene>>,
  sceneDestroy: Scalars['Boolean'],
  scenesUpdate?: Maybe<Array<Maybe<Scene>>>,
  /** Increments the o-counter for a scene. Returns the new value */
  sceneIncrementO: Scalars['Int'],
  /** Decrements the o-counter for a scene. Returns the new value */
  sceneDecrementO: Scalars['Int'],
  /** Resets the o-counter for a scene to 0. Returns the new value */
  sceneResetO: Scalars['Int'],
  sceneMarkerCreate?: Maybe<SceneMarker>,
  sceneMarkerUpdate?: Maybe<SceneMarker>,
  sceneMarkerDestroy: Scalars['Boolean'],
  performerCreate?: Maybe<Performer>,
  performerUpdate?: Maybe<Performer>,
  performerDestroy: Scalars['Boolean'],
  studioCreate?: Maybe<Studio>,
  studioUpdate?: Maybe<Studio>,
  studioDestroy: Scalars['Boolean'],
  tagCreate?: Maybe<Tag>,
  tagUpdate?: Maybe<Tag>,
  tagDestroy: Scalars['Boolean'],
  /** Change general configuration options */
  configureGeneral: ConfigGeneralResult,
  configureInterface: ConfigInterfaceResult,
};


export type MutationSceneUpdateArgs = {
  input: SceneUpdateInput
};


export type MutationBulkSceneUpdateArgs = {
  input: BulkSceneUpdateInput
};


export type MutationSceneDestroyArgs = {
  input: SceneDestroyInput
};


export type MutationScenesUpdateArgs = {
  input: Array<SceneUpdateInput>
};


export type MutationSceneIncrementOArgs = {
  id: Scalars['ID']
};


export type MutationSceneDecrementOArgs = {
  id: Scalars['ID']
};


export type MutationSceneResetOArgs = {
  id: Scalars['ID']
};


export type MutationSceneMarkerCreateArgs = {
  input: SceneMarkerCreateInput
};


export type MutationSceneMarkerUpdateArgs = {
  input: SceneMarkerUpdateInput
};


export type MutationSceneMarkerDestroyArgs = {
  id: Scalars['ID']
};


export type MutationPerformerCreateArgs = {
  input: PerformerCreateInput
};


export type MutationPerformerUpdateArgs = {
  input: PerformerUpdateInput
};


export type MutationPerformerDestroyArgs = {
  input: PerformerDestroyInput
};


export type MutationStudioCreateArgs = {
  input: StudioCreateInput
};


export type MutationStudioUpdateArgs = {
  input: StudioUpdateInput
};


export type MutationStudioDestroyArgs = {
  input: StudioDestroyInput
};


export type MutationTagCreateArgs = {
  input: TagCreateInput
};


export type MutationTagUpdateArgs = {
  input: TagUpdateInput
};


export type MutationTagDestroyArgs = {
  input: TagDestroyInput
};


export type MutationConfigureGeneralArgs = {
  input: ConfigGeneralInput
};


export type MutationConfigureInterfaceArgs = {
  input: ConfigInterfaceInput
};

export type Performer = {
   __typename?: 'Performer',
  id: Scalars['ID'],
  checksum: Scalars['String'],
  name?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  twitter?: Maybe<Scalars['String']>,
  instagram?: Maybe<Scalars['String']>,
  birthdate?: Maybe<Scalars['String']>,
  ethnicity?: Maybe<Scalars['String']>,
  country?: Maybe<Scalars['String']>,
  eye_color?: Maybe<Scalars['String']>,
  height?: Maybe<Scalars['String']>,
  measurements?: Maybe<Scalars['String']>,
  fake_tits?: Maybe<Scalars['String']>,
  career_length?: Maybe<Scalars['String']>,
  tattoos?: Maybe<Scalars['String']>,
  piercings?: Maybe<Scalars['String']>,
  aliases?: Maybe<Scalars['String']>,
  favorite: Scalars['Boolean'],
  image_path?: Maybe<Scalars['String']>,
  scene_count?: Maybe<Scalars['Int']>,
  scenes: Array<Scene>,
};

export type PerformerCreateInput = {
  name?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  birthdate?: Maybe<Scalars['String']>,
  ethnicity?: Maybe<Scalars['String']>,
  country?: Maybe<Scalars['String']>,
  eye_color?: Maybe<Scalars['String']>,
  height?: Maybe<Scalars['String']>,
  measurements?: Maybe<Scalars['String']>,
  fake_tits?: Maybe<Scalars['String']>,
  career_length?: Maybe<Scalars['String']>,
  tattoos?: Maybe<Scalars['String']>,
  piercings?: Maybe<Scalars['String']>,
  aliases?: Maybe<Scalars['String']>,
  twitter?: Maybe<Scalars['String']>,
  instagram?: Maybe<Scalars['String']>,
  favorite?: Maybe<Scalars['Boolean']>,
  /** This should be base64 encoded */
  image?: Maybe<Scalars['String']>,
};

export type PerformerDestroyInput = {
  id: Scalars['ID'],
};

export type PerformerFilterType = {
  /** Filter by favorite */
  filter_favorites?: Maybe<Scalars['Boolean']>,
  /** Filter by birth year */
  birth_year?: Maybe<IntCriterionInput>,
  /** Filter by age */
  age?: Maybe<IntCriterionInput>,
  /** Filter by ethnicity */
  ethnicity?: Maybe<StringCriterionInput>,
  /** Filter by country */
  country?: Maybe<StringCriterionInput>,
  /** Filter by eye color */
  eye_color?: Maybe<StringCriterionInput>,
  /** Filter by height */
  height?: Maybe<StringCriterionInput>,
  /** Filter by measurements */
  measurements?: Maybe<StringCriterionInput>,
  /** Filter by fake tits value */
  fake_tits?: Maybe<StringCriterionInput>,
  /** Filter by career length */
  career_length?: Maybe<StringCriterionInput>,
  /** Filter by tattoos */
  tattoos?: Maybe<StringCriterionInput>,
  /** Filter by piercings */
  piercings?: Maybe<StringCriterionInput>,
  /** Filter by aliases */
  aliases?: Maybe<StringCriterionInput>,
};

export type PerformerUpdateInput = {
  id: Scalars['ID'],
  name?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  birthdate?: Maybe<Scalars['String']>,
  ethnicity?: Maybe<Scalars['String']>,
  country?: Maybe<Scalars['String']>,
  eye_color?: Maybe<Scalars['String']>,
  height?: Maybe<Scalars['String']>,
  measurements?: Maybe<Scalars['String']>,
  fake_tits?: Maybe<Scalars['String']>,
  career_length?: Maybe<Scalars['String']>,
  tattoos?: Maybe<Scalars['String']>,
  piercings?: Maybe<Scalars['String']>,
  aliases?: Maybe<Scalars['String']>,
  twitter?: Maybe<Scalars['String']>,
  instagram?: Maybe<Scalars['String']>,
  favorite?: Maybe<Scalars['Boolean']>,
  /** This should be base64 encoded */
  image?: Maybe<Scalars['String']>,
};

/** The query root for this schema */
export type Query = {
   __typename?: 'Query',
  /** Find a scene by ID or Checksum */
  findScene?: Maybe<Scene>,
  /** A function which queries Scene objects */
  findScenes: FindScenesResultType,
  findScenesByPathRegex: FindScenesResultType,
  parseSceneFilenames: SceneParserResultType,
  /** A function which queries SceneMarker objects */
  findSceneMarkers: FindSceneMarkersResultType,
  /** Find a performer by ID */
  findPerformer?: Maybe<Performer>,
  /** A function which queries Performer objects */
  findPerformers: FindPerformersResultType,
  /** Find a studio by ID */
  findStudio?: Maybe<Studio>,
  /** A function which queries Studio objects */
  findStudios: FindStudiosResultType,
  findGallery?: Maybe<Gallery>,
  findGalleries: FindGalleriesResultType,
  findTag?: Maybe<Tag>,
  /** Retrieve random scene markers for the wall */
  markerWall: Array<SceneMarker>,
  /** Retrieve random scenes for the wall */
  sceneWall: Array<Scene>,
  /** Get marker strings */
  markerStrings: Array<Maybe<MarkerStringsResultType>>,
  /** Get the list of valid galleries for a given scene ID */
  validGalleriesForScene: Array<Gallery>,
  /** Get stats */
  stats: StatsResultType,
  /** Organize scene markers by tag for a given scene ID */
  sceneMarkerTags: Array<SceneMarkerTag>,
  logs: Array<LogEntry>,
  /** List available scrapers */
  listPerformerScrapers: Array<Scraper>,
  listSceneScrapers: Array<Scraper>,
  /** Scrape a list of performers based on name */
  scrapePerformerList: Array<ScrapedPerformer>,
  /** Scrapes a complete performer record based on a scrapePerformerList result */
  scrapePerformer?: Maybe<ScrapedPerformer>,
  /** Scrapes a complete performer record based on a URL */
  scrapePerformerURL?: Maybe<ScrapedPerformer>,
  /** Scrapes a complete scene record based on an existing scene */
  scrapeScene?: Maybe<ScrapedScene>,
  /** Scrapes a complete performer record based on a URL */
  scrapeSceneURL?: Maybe<ScrapedScene>,
  /** Scrape a performer using Freeones */
  scrapeFreeones?: Maybe<ScrapedPerformer>,
  /** Scrape a list of performers from a query */
  scrapeFreeonesPerformerList: Array<Scalars['String']>,
  /** Returns the current, complete configuration */
  configuration: ConfigResult,
  /** Returns an array of paths for the given path */
  directories: Array<Scalars['String']>,
  /** Start an import. Returns the job ID */
  metadataImport: Scalars['String'],
  /** Start an export. Returns the job ID */
  metadataExport: Scalars['String'],
  /** Start a scan. Returns the job ID */
  metadataScan: Scalars['String'],
  /** Start generating content. Returns the job ID */
  metadataGenerate: Scalars['String'],
  /** Start auto-tagging. Returns the job ID */
  metadataAutoTag: Scalars['String'],
  /** Clean metadata. Returns the job ID */
  metadataClean: Scalars['String'],
  jobStatus: MetadataUpdateStatus,
  stopJob: Scalars['Boolean'],
  allPerformers: Array<Performer>,
  allStudios: Array<Studio>,
  allTags: Array<Tag>,
  /** Version */
  version: Version,
  /** LatestVersion */
  latestversion: ShortVersion,
};


/** The query root for this schema */
export type QueryFindSceneArgs = {
  id?: Maybe<Scalars['ID']>,
  checksum?: Maybe<Scalars['String']>
};


/** The query root for this schema */
export type QueryFindScenesArgs = {
  scene_filter?: Maybe<SceneFilterType>,
  scene_ids?: Maybe<Array<Scalars['Int']>>,
  filter?: Maybe<FindFilterType>
};


/** The query root for this schema */
export type QueryFindScenesByPathRegexArgs = {
  filter?: Maybe<FindFilterType>
};


/** The query root for this schema */
export type QueryParseSceneFilenamesArgs = {
  filter?: Maybe<FindFilterType>,
  config: SceneParserInput
};


/** The query root for this schema */
export type QueryFindSceneMarkersArgs = {
  scene_marker_filter?: Maybe<SceneMarkerFilterType>,
  filter?: Maybe<FindFilterType>
};


/** The query root for this schema */
export type QueryFindPerformerArgs = {
  id: Scalars['ID']
};


/** The query root for this schema */
export type QueryFindPerformersArgs = {
  performer_filter?: Maybe<PerformerFilterType>,
  filter?: Maybe<FindFilterType>
};


/** The query root for this schema */
export type QueryFindStudioArgs = {
  id: Scalars['ID']
};


/** The query root for this schema */
export type QueryFindStudiosArgs = {
  filter?: Maybe<FindFilterType>
};


/** The query root for this schema */
export type QueryFindGalleryArgs = {
  id: Scalars['ID']
};


/** The query root for this schema */
export type QueryFindGalleriesArgs = {
  filter?: Maybe<FindFilterType>
};


/** The query root for this schema */
export type QueryFindTagArgs = {
  id: Scalars['ID']
};


/** The query root for this schema */
export type QueryMarkerWallArgs = {
  q?: Maybe<Scalars['String']>
};


/** The query root for this schema */
export type QuerySceneWallArgs = {
  q?: Maybe<Scalars['String']>
};


/** The query root for this schema */
export type QueryMarkerStringsArgs = {
  q?: Maybe<Scalars['String']>,
  sort?: Maybe<Scalars['String']>
};


/** The query root for this schema */
export type QueryValidGalleriesForSceneArgs = {
  scene_id?: Maybe<Scalars['ID']>
};


/** The query root for this schema */
export type QuerySceneMarkerTagsArgs = {
  scene_id: Scalars['ID']
};


/** The query root for this schema */
export type QueryScrapePerformerListArgs = {
  scraper_id: Scalars['ID'],
  query: Scalars['String']
};


/** The query root for this schema */
export type QueryScrapePerformerArgs = {
  scraper_id: Scalars['ID'],
  scraped_performer: ScrapedPerformerInput
};


/** The query root for this schema */
export type QueryScrapePerformerUrlArgs = {
  url: Scalars['String']
};


/** The query root for this schema */
export type QueryScrapeSceneArgs = {
  scraper_id: Scalars['ID'],
  scene: SceneUpdateInput
};


/** The query root for this schema */
export type QueryScrapeSceneUrlArgs = {
  url: Scalars['String']
};


/** The query root for this schema */
export type QueryScrapeFreeonesArgs = {
  performer_name: Scalars['String']
};


/** The query root for this schema */
export type QueryScrapeFreeonesPerformerListArgs = {
  query: Scalars['String']
};


/** The query root for this schema */
export type QueryDirectoriesArgs = {
  path?: Maybe<Scalars['String']>
};


/** The query root for this schema */
export type QueryMetadataScanArgs = {
  input: ScanMetadataInput
};


/** The query root for this schema */
export type QueryMetadataGenerateArgs = {
  input: GenerateMetadataInput
};


/** The query root for this schema */
export type QueryMetadataAutoTagArgs = {
  input: AutoTagMetadataInput
};

export enum ResolutionEnum {
  /** 240p */
  Low = 'LOW',
  /** 480p */
  Standard = 'STANDARD',
  /** 720p */
  StandardHd = 'STANDARD_HD',
  /** 1080p */
  FullHd = 'FULL_HD',
  /** 4k */
  FourK = 'FOUR_K'
}

export type ScanMetadataInput = {
  useFileMetadata: Scalars['Boolean'],
};

export type Scene = {
   __typename?: 'Scene',
  id: Scalars['ID'],
  checksum: Scalars['String'],
  title?: Maybe<Scalars['String']>,
  details?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  date?: Maybe<Scalars['String']>,
  rating?: Maybe<Scalars['Int']>,
  o_counter?: Maybe<Scalars['Int']>,
  path: Scalars['String'],
  file: SceneFileType,
  paths: ScenePathsType,
  is_streamable: Scalars['Boolean'],
  scene_markers: Array<SceneMarker>,
  gallery?: Maybe<Gallery>,
  studio?: Maybe<Studio>,
  tags: Array<Tag>,
  performers: Array<Performer>,
};

export type SceneDestroyInput = {
  id: Scalars['ID'],
  delete_file?: Maybe<Scalars['Boolean']>,
  delete_generated?: Maybe<Scalars['Boolean']>,
};

export type SceneFileType = {
   __typename?: 'SceneFileType',
  size?: Maybe<Scalars['String']>,
  duration?: Maybe<Scalars['Float']>,
  video_codec?: Maybe<Scalars['String']>,
  audio_codec?: Maybe<Scalars['String']>,
  width?: Maybe<Scalars['Int']>,
  height?: Maybe<Scalars['Int']>,
  framerate?: Maybe<Scalars['Float']>,
  bitrate?: Maybe<Scalars['Int']>,
};

export type SceneFilterType = {
  /** Filter by rating */
  rating?: Maybe<IntCriterionInput>,
  /** Filter by o-counter */
  o_counter?: Maybe<IntCriterionInput>,
  /** Filter by resolution */
  resolution?: Maybe<ResolutionEnum>,
  /** Filter by duration (in seconds) */
  duration?: Maybe<IntCriterionInput>,
  /** Filter to only include scenes which have markers. `true` or `false` */
  has_markers?: Maybe<Scalars['String']>,
  /** Filter to only include scenes missing this property */
  is_missing?: Maybe<Scalars['String']>,
  /** Filter to only include scenes with this studio */
  studios?: Maybe<MultiCriterionInput>,
  /** Filter to only include scenes with these tags */
  tags?: Maybe<MultiCriterionInput>,
  /** Filter to only include scenes with these performers */
  performers?: Maybe<MultiCriterionInput>,
};

export type SceneMarker = {
   __typename?: 'SceneMarker',
  id: Scalars['ID'],
  scene: Scene,
  title: Scalars['String'],
  seconds: Scalars['Float'],
  primary_tag: Tag,
  tags: Array<Tag>,
  /** The path to stream this marker */
  stream: Scalars['String'],
  /** The path to the preview image for this marker */
  preview: Scalars['String'],
};

export type SceneMarkerCreateInput = {
  title: Scalars['String'],
  seconds: Scalars['Float'],
  scene_id: Scalars['ID'],
  primary_tag_id: Scalars['ID'],
  tag_ids?: Maybe<Array<Scalars['ID']>>,
};

export type SceneMarkerFilterType = {
  /** Filter to only include scene markers with this tag */
  tag_id?: Maybe<Scalars['ID']>,
  /** Filter to only include scene markers with these tags */
  tags?: Maybe<MultiCriterionInput>,
  /** Filter to only include scene markers attached to a scene with these tags */
  scene_tags?: Maybe<MultiCriterionInput>,
  /** Filter to only include scene markers with these performers */
  performers?: Maybe<MultiCriterionInput>,
};

export type SceneMarkerTag = {
   __typename?: 'SceneMarkerTag',
  tag: Tag,
  scene_markers: Array<SceneMarker>,
};

export type SceneMarkerUpdateInput = {
  id: Scalars['ID'],
  title: Scalars['String'],
  seconds: Scalars['Float'],
  scene_id: Scalars['ID'],
  primary_tag_id: Scalars['ID'],
  tag_ids?: Maybe<Array<Scalars['ID']>>,
};

export type SceneParserInput = {
  ignoreWords?: Maybe<Array<Scalars['String']>>,
  whitespaceCharacters?: Maybe<Scalars['String']>,
  capitalizeTitle?: Maybe<Scalars['Boolean']>,
};

export type SceneParserResult = {
   __typename?: 'SceneParserResult',
  scene: Scene,
  title?: Maybe<Scalars['String']>,
  details?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  date?: Maybe<Scalars['String']>,
  rating?: Maybe<Scalars['Int']>,
  studio_id?: Maybe<Scalars['ID']>,
  gallery_id?: Maybe<Scalars['ID']>,
  performer_ids?: Maybe<Array<Scalars['ID']>>,
  tag_ids?: Maybe<Array<Scalars['ID']>>,
};

export type SceneParserResultType = {
   __typename?: 'SceneParserResultType',
  count: Scalars['Int'],
  results: Array<SceneParserResult>,
};

export type ScenePathsType = {
   __typename?: 'ScenePathsType',
  screenshot?: Maybe<Scalars['String']>,
  preview?: Maybe<Scalars['String']>,
  stream?: Maybe<Scalars['String']>,
  webp?: Maybe<Scalars['String']>,
  vtt?: Maybe<Scalars['String']>,
  chapters_vtt?: Maybe<Scalars['String']>,
};

export type SceneUpdateInput = {
  clientMutationId?: Maybe<Scalars['String']>,
  id: Scalars['ID'],
  title?: Maybe<Scalars['String']>,
  details?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  date?: Maybe<Scalars['String']>,
  rating?: Maybe<Scalars['Int']>,
  studio_id?: Maybe<Scalars['ID']>,
  gallery_id?: Maybe<Scalars['ID']>,
  performer_ids?: Maybe<Array<Scalars['ID']>>,
  tag_ids?: Maybe<Array<Scalars['ID']>>,
  /** This should be base64 encoded */
  cover_image?: Maybe<Scalars['String']>,
};

/** A performer from a scraping operation... */
export type ScrapedPerformer = {
   __typename?: 'ScrapedPerformer',
  name?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  twitter?: Maybe<Scalars['String']>,
  instagram?: Maybe<Scalars['String']>,
  birthdate?: Maybe<Scalars['String']>,
  ethnicity?: Maybe<Scalars['String']>,
  country?: Maybe<Scalars['String']>,
  eye_color?: Maybe<Scalars['String']>,
  height?: Maybe<Scalars['String']>,
  measurements?: Maybe<Scalars['String']>,
  fake_tits?: Maybe<Scalars['String']>,
  career_length?: Maybe<Scalars['String']>,
  tattoos?: Maybe<Scalars['String']>,
  piercings?: Maybe<Scalars['String']>,
  aliases?: Maybe<Scalars['String']>,
};

export type ScrapedPerformerInput = {
  name?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  twitter?: Maybe<Scalars['String']>,
  instagram?: Maybe<Scalars['String']>,
  birthdate?: Maybe<Scalars['String']>,
  ethnicity?: Maybe<Scalars['String']>,
  country?: Maybe<Scalars['String']>,
  eye_color?: Maybe<Scalars['String']>,
  height?: Maybe<Scalars['String']>,
  measurements?: Maybe<Scalars['String']>,
  fake_tits?: Maybe<Scalars['String']>,
  career_length?: Maybe<Scalars['String']>,
  tattoos?: Maybe<Scalars['String']>,
  piercings?: Maybe<Scalars['String']>,
  aliases?: Maybe<Scalars['String']>,
};

export type ScrapedScene = {
   __typename?: 'ScrapedScene',
  title?: Maybe<Scalars['String']>,
  details?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  date?: Maybe<Scalars['String']>,
  file?: Maybe<SceneFileType>,
  studio?: Maybe<ScrapedSceneStudio>,
  tags?: Maybe<Array<ScrapedSceneTag>>,
  performers?: Maybe<Array<ScrapedScenePerformer>>,
};

export type ScrapedScenePerformer = {
   __typename?: 'ScrapedScenePerformer',
  /** Set if performer matched */
  id?: Maybe<Scalars['ID']>,
  name: Scalars['String'],
  url?: Maybe<Scalars['String']>,
  twitter?: Maybe<Scalars['String']>,
  instagram?: Maybe<Scalars['String']>,
  birthdate?: Maybe<Scalars['String']>,
  ethnicity?: Maybe<Scalars['String']>,
  country?: Maybe<Scalars['String']>,
  eye_color?: Maybe<Scalars['String']>,
  height?: Maybe<Scalars['String']>,
  measurements?: Maybe<Scalars['String']>,
  fake_tits?: Maybe<Scalars['String']>,
  career_length?: Maybe<Scalars['String']>,
  tattoos?: Maybe<Scalars['String']>,
  piercings?: Maybe<Scalars['String']>,
  aliases?: Maybe<Scalars['String']>,
};

export type ScrapedSceneStudio = {
   __typename?: 'ScrapedSceneStudio',
  /** Set if studio matched */
  id?: Maybe<Scalars['ID']>,
  name: Scalars['String'],
  url?: Maybe<Scalars['String']>,
};

export type ScrapedSceneTag = {
   __typename?: 'ScrapedSceneTag',
  /** Set if tag matched */
  id?: Maybe<Scalars['ID']>,
  name: Scalars['String'],
};

export type Scraper = {
   __typename?: 'Scraper',
  id: Scalars['ID'],
  name: Scalars['String'],
  /** Details for performer scraper */
  performer?: Maybe<ScraperSpec>,
  /** Details for scene scraper */
  scene?: Maybe<ScraperSpec>,
};

export type ScraperSpec = {
   __typename?: 'ScraperSpec',
  /** URLs matching these can be scraped with */
  urls?: Maybe<Array<Scalars['String']>>,
  supported_scrapes: Array<ScrapeType>,
};

export enum ScrapeType {
  /** From text query */
  Name = 'NAME',
  /** From existing object */
  Fragment = 'FRAGMENT',
  /** From URL */
  Url = 'URL'
}

export type ShortVersion = {
   __typename?: 'ShortVersion',
  shorthash: Scalars['String'],
  url: Scalars['String'],
};

export enum SortDirectionEnum {
  Asc = 'ASC',
  Desc = 'DESC'
}

export type StatsResultType = {
   __typename?: 'StatsResultType',
  scene_count: Scalars['Int'],
  gallery_count: Scalars['Int'],
  performer_count: Scalars['Int'],
  studio_count: Scalars['Int'],
  tag_count: Scalars['Int'],
};

export enum StreamingResolutionEnum {
  /** 240p */
  Low = 'LOW',
  /** 480p */
  Standard = 'STANDARD',
  /** 720p */
  StandardHd = 'STANDARD_HD',
  /** 1080p */
  FullHd = 'FULL_HD',
  /** 4k */
  FourK = 'FOUR_K',
  /** Original */
  Original = 'ORIGINAL'
}

export type StringCriterionInput = {
  value: Scalars['String'],
  modifier: CriterionModifier,
};

export type Studio = {
   __typename?: 'Studio',
  id: Scalars['ID'],
  checksum: Scalars['String'],
  name: Scalars['String'],
  url?: Maybe<Scalars['String']>,
  image_path?: Maybe<Scalars['String']>,
  scene_count?: Maybe<Scalars['Int']>,
};

export type StudioCreateInput = {
  name: Scalars['String'],
  url?: Maybe<Scalars['String']>,
  /** This should be base64 encoded */
  image?: Maybe<Scalars['String']>,
};

export type StudioDestroyInput = {
  id: Scalars['ID'],
};

export type StudioUpdateInput = {
  id: Scalars['ID'],
  name?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  /** This should be base64 encoded */
  image?: Maybe<Scalars['String']>,
};

export type Subscription = {
   __typename?: 'Subscription',
  /** Update from the metadata manager */
  metadataUpdate: MetadataUpdateStatus,
  loggingSubscribe: Array<LogEntry>,
};

export type Tag = {
   __typename?: 'Tag',
  id: Scalars['ID'],
  name: Scalars['String'],
  scene_count?: Maybe<Scalars['Int']>,
  scene_marker_count?: Maybe<Scalars['Int']>,
};

export type TagCreateInput = {
  name: Scalars['String'],
};

export type TagDestroyInput = {
  id: Scalars['ID'],
};

export type TagUpdateInput = {
  id: Scalars['ID'],
  name: Scalars['String'],
};


export type Version = {
   __typename?: 'Version',
  version?: Maybe<Scalars['String']>,
  hash: Scalars['String'],
  build_time: Scalars['String'],
};

export type ConfigGeneralDataFragment = (
  { __typename?: 'ConfigGeneralResult' }
  & Pick<ConfigGeneralResult, 'stashes' | 'databasePath' | 'generatedPath' | 'maxTranscodeSize' | 'maxStreamingTranscodeSize' | 'username' | 'password' | 'logFile' | 'logOut' | 'logLevel' | 'logAccess' | 'excludes'>
);

export type ConfigInterfaceDataFragment = (
  { __typename?: 'ConfigInterfaceResult' }
  & Pick<ConfigInterfaceResult, 'soundOnPreview' | 'wallShowTitle' | 'maximumLoopDuration' | 'autostartVideo' | 'showStudioAsText' | 'css' | 'cssEnabled' | 'language'>
);

export type ConfigDataFragment = (
  { __typename?: 'ConfigResult' }
  & { general: (
    { __typename?: 'ConfigGeneralResult' }
    & ConfigGeneralDataFragment
  ), interface: (
    { __typename?: 'ConfigInterfaceResult' }
    & ConfigInterfaceDataFragment
  ) }
);

export type GalleryDataFragment = (
  { __typename?: 'Gallery' }
  & Pick<Gallery, 'id' | 'checksum' | 'path' | 'title'>
  & { files: Array<(
    { __typename?: 'GalleryFilesType' }
    & Pick<GalleryFilesType, 'index' | 'name' | 'path'>
  )> }
);

export type LogEntryDataFragment = (
  { __typename?: 'LogEntry' }
  & Pick<LogEntry, 'time' | 'level' | 'message'>
);

export type SlimPerformerDataFragment = (
  { __typename?: 'Performer' }
  & Pick<Performer, 'id' | 'name' | 'image_path'>
);

export type PerformerDataFragment = (
  { __typename?: 'Performer' }
  & Pick<Performer, 'id' | 'checksum' | 'name' | 'url' | 'twitter' | 'instagram' | 'birthdate' | 'ethnicity' | 'country' | 'eye_color' | 'height' | 'measurements' | 'fake_tits' | 'career_length' | 'tattoos' | 'piercings' | 'aliases' | 'favorite' | 'image_path' | 'scene_count'>
);

export type SceneMarkerDataFragment = (
  { __typename?: 'SceneMarker' }
  & Pick<SceneMarker, 'id' | 'title' | 'seconds' | 'stream' | 'preview'>
  & { scene: (
    { __typename?: 'Scene' }
    & Pick<Scene, 'id'>
  ), primary_tag: (
    { __typename?: 'Tag' }
    & Pick<Tag, 'id' | 'name'>
  ), tags: Array<(
    { __typename?: 'Tag' }
    & Pick<Tag, 'id' | 'name'>
  )> }
);

export type SlimSceneDataFragment = (
  { __typename?: 'Scene' }
  & Pick<Scene, 'id' | 'checksum' | 'title' | 'details' | 'url' | 'date' | 'rating' | 'o_counter' | 'path'>
  & { file: (
    { __typename?: 'SceneFileType' }
    & Pick<SceneFileType, 'size' | 'duration' | 'video_codec' | 'audio_codec' | 'width' | 'height' | 'framerate' | 'bitrate'>
  ), paths: (
    { __typename?: 'ScenePathsType' }
    & Pick<ScenePathsType, 'screenshot' | 'preview' | 'stream' | 'webp' | 'vtt' | 'chapters_vtt'>
  ), scene_markers: Array<(
    { __typename?: 'SceneMarker' }
    & Pick<SceneMarker, 'id' | 'title' | 'seconds'>
  )>, gallery: Maybe<(
    { __typename?: 'Gallery' }
    & Pick<Gallery, 'id' | 'path' | 'title'>
  )>, studio: Maybe<(
    { __typename?: 'Studio' }
    & Pick<Studio, 'id' | 'name' | 'image_path'>
  )>, tags: Array<(
    { __typename?: 'Tag' }
    & Pick<Tag, 'id' | 'name'>
  )>, performers: Array<(
    { __typename?: 'Performer' }
    & Pick<Performer, 'id' | 'name' | 'favorite' | 'image_path'>
  )> }
);

export type SceneDataFragment = (
  { __typename?: 'Scene' }
  & Pick<Scene, 'id' | 'checksum' | 'title' | 'details' | 'url' | 'date' | 'rating' | 'o_counter' | 'path' | 'is_streamable'>
  & { file: (
    { __typename?: 'SceneFileType' }
    & Pick<SceneFileType, 'size' | 'duration' | 'video_codec' | 'audio_codec' | 'width' | 'height' | 'framerate' | 'bitrate'>
  ), paths: (
    { __typename?: 'ScenePathsType' }
    & Pick<ScenePathsType, 'screenshot' | 'preview' | 'stream' | 'webp' | 'vtt' | 'chapters_vtt'>
  ), scene_markers: Array<(
    { __typename?: 'SceneMarker' }
    & SceneMarkerDataFragment
  )>, gallery: Maybe<(
    { __typename?: 'Gallery' }
    & GalleryDataFragment
  )>, studio: Maybe<(
    { __typename?: 'Studio' }
    & StudioDataFragment
  )>, tags: Array<(
    { __typename?: 'Tag' }
    & TagDataFragment
  )>, performers: Array<(
    { __typename?: 'Performer' }
    & PerformerDataFragment
  )> }
);

export type ScrapedPerformerDataFragment = (
  { __typename?: 'ScrapedPerformer' }
  & Pick<ScrapedPerformer, 'name' | 'url' | 'twitter' | 'instagram' | 'birthdate' | 'ethnicity' | 'country' | 'eye_color' | 'height' | 'measurements' | 'fake_tits' | 'career_length' | 'tattoos' | 'piercings' | 'aliases'>
);

export type ScrapedScenePerformerDataFragment = (
  { __typename?: 'ScrapedScenePerformer' }
  & Pick<ScrapedScenePerformer, 'id' | 'name' | 'url' | 'twitter' | 'instagram' | 'birthdate' | 'ethnicity' | 'country' | 'eye_color' | 'height' | 'measurements' | 'fake_tits' | 'career_length' | 'tattoos' | 'piercings' | 'aliases'>
);

export type ScrapedSceneStudioDataFragment = (
  { __typename?: 'ScrapedSceneStudio' }
  & Pick<ScrapedSceneStudio, 'id' | 'name' | 'url'>
);

export type ScrapedSceneTagDataFragment = (
  { __typename?: 'ScrapedSceneTag' }
  & Pick<ScrapedSceneTag, 'id' | 'name'>
);

export type ScrapedSceneDataFragment = (
  { __typename?: 'ScrapedScene' }
  & Pick<ScrapedScene, 'title' | 'details' | 'url' | 'date'>
  & { file: Maybe<(
    { __typename?: 'SceneFileType' }
    & Pick<SceneFileType, 'size' | 'duration' | 'video_codec' | 'audio_codec' | 'width' | 'height' | 'framerate' | 'bitrate'>
  )>, studio: Maybe<(
    { __typename?: 'ScrapedSceneStudio' }
    & ScrapedSceneStudioDataFragment
  )>, tags: Maybe<Array<(
    { __typename?: 'ScrapedSceneTag' }
    & ScrapedSceneTagDataFragment
  )>>, performers: Maybe<Array<(
    { __typename?: 'ScrapedScenePerformer' }
    & ScrapedScenePerformerDataFragment
  )>> }
);

export type SlimStudioDataFragment = (
  { __typename?: 'Studio' }
  & Pick<Studio, 'id' | 'name' | 'image_path'>
);

export type StudioDataFragment = (
  { __typename?: 'Studio' }
  & Pick<Studio, 'id' | 'checksum' | 'name' | 'url' | 'image_path' | 'scene_count'>
);

export type TagDataFragment = (
  { __typename?: 'Tag' }
  & Pick<Tag, 'id' | 'name' | 'scene_count' | 'scene_marker_count'>
);

export type ConfigureGeneralMutationVariables = {
  input: ConfigGeneralInput
};


export type ConfigureGeneralMutation = (
  { __typename?: 'Mutation' }
  & { configureGeneral: (
    { __typename?: 'ConfigGeneralResult' }
    & ConfigGeneralDataFragment
  ) }
);

export type ConfigureInterfaceMutationVariables = {
  input: ConfigInterfaceInput
};


export type ConfigureInterfaceMutation = (
  { __typename?: 'Mutation' }
  & { configureInterface: (
    { __typename?: 'ConfigInterfaceResult' }
    & ConfigInterfaceDataFragment
  ) }
);

export type PerformerCreateMutationVariables = {
  name?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  birthdate?: Maybe<Scalars['String']>,
  ethnicity?: Maybe<Scalars['String']>,
  country?: Maybe<Scalars['String']>,
  eye_color?: Maybe<Scalars['String']>,
  height?: Maybe<Scalars['String']>,
  measurements?: Maybe<Scalars['String']>,
  fake_tits?: Maybe<Scalars['String']>,
  career_length?: Maybe<Scalars['String']>,
  tattoos?: Maybe<Scalars['String']>,
  piercings?: Maybe<Scalars['String']>,
  aliases?: Maybe<Scalars['String']>,
  twitter?: Maybe<Scalars['String']>,
  instagram?: Maybe<Scalars['String']>,
  favorite?: Maybe<Scalars['Boolean']>,
  image?: Maybe<Scalars['String']>
};


export type PerformerCreateMutation = (
  { __typename?: 'Mutation' }
  & { performerCreate: Maybe<(
    { __typename?: 'Performer' }
    & PerformerDataFragment
  )> }
);

export type PerformerUpdateMutationVariables = {
  id: Scalars['ID'],
  name?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  birthdate?: Maybe<Scalars['String']>,
  ethnicity?: Maybe<Scalars['String']>,
  country?: Maybe<Scalars['String']>,
  eye_color?: Maybe<Scalars['String']>,
  height?: Maybe<Scalars['String']>,
  measurements?: Maybe<Scalars['String']>,
  fake_tits?: Maybe<Scalars['String']>,
  career_length?: Maybe<Scalars['String']>,
  tattoos?: Maybe<Scalars['String']>,
  piercings?: Maybe<Scalars['String']>,
  aliases?: Maybe<Scalars['String']>,
  twitter?: Maybe<Scalars['String']>,
  instagram?: Maybe<Scalars['String']>,
  favorite?: Maybe<Scalars['Boolean']>,
  image?: Maybe<Scalars['String']>
};


export type PerformerUpdateMutation = (
  { __typename?: 'Mutation' }
  & { performerUpdate: Maybe<(
    { __typename?: 'Performer' }
    & PerformerDataFragment
  )> }
);

export type PerformerDestroyMutationVariables = {
  id: Scalars['ID']
};


export type PerformerDestroyMutation = (
  { __typename?: 'Mutation' }
  & Pick<Mutation, 'performerDestroy'>
);

export type SceneMarkerCreateMutationVariables = {
  title: Scalars['String'],
  seconds: Scalars['Float'],
  scene_id: Scalars['ID'],
  primary_tag_id: Scalars['ID'],
  tag_ids?: Maybe<Array<Scalars['ID']>>
};


export type SceneMarkerCreateMutation = (
  { __typename?: 'Mutation' }
  & { sceneMarkerCreate: Maybe<(
    { __typename?: 'SceneMarker' }
    & SceneMarkerDataFragment
  )> }
);

export type SceneMarkerUpdateMutationVariables = {
  id: Scalars['ID'],
  title: Scalars['String'],
  seconds: Scalars['Float'],
  scene_id: Scalars['ID'],
  primary_tag_id: Scalars['ID'],
  tag_ids?: Maybe<Array<Scalars['ID']>>
};


export type SceneMarkerUpdateMutation = (
  { __typename?: 'Mutation' }
  & { sceneMarkerUpdate: Maybe<(
    { __typename?: 'SceneMarker' }
    & SceneMarkerDataFragment
  )> }
);

export type SceneMarkerDestroyMutationVariables = {
  id: Scalars['ID']
};


export type SceneMarkerDestroyMutation = (
  { __typename?: 'Mutation' }
  & Pick<Mutation, 'sceneMarkerDestroy'>
);

export type SceneUpdateMutationVariables = {
  id: Scalars['ID'],
  title?: Maybe<Scalars['String']>,
  details?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  date?: Maybe<Scalars['String']>,
  rating?: Maybe<Scalars['Int']>,
  studio_id?: Maybe<Scalars['ID']>,
  gallery_id?: Maybe<Scalars['ID']>,
  performer_ids?: Maybe<Array<Scalars['ID']>>,
  tag_ids?: Maybe<Array<Scalars['ID']>>,
  cover_image?: Maybe<Scalars['String']>
};


export type SceneUpdateMutation = (
  { __typename?: 'Mutation' }
  & { sceneUpdate: Maybe<(
    { __typename?: 'Scene' }
    & SceneDataFragment
  )> }
);

export type BulkSceneUpdateMutationVariables = {
  ids?: Maybe<Array<Scalars['ID']>>,
  title?: Maybe<Scalars['String']>,
  details?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  date?: Maybe<Scalars['String']>,
  rating?: Maybe<Scalars['Int']>,
  studio_id?: Maybe<Scalars['ID']>,
  gallery_id?: Maybe<Scalars['ID']>,
  performer_ids?: Maybe<Array<Scalars['ID']>>,
  tag_ids?: Maybe<Array<Scalars['ID']>>
};


export type BulkSceneUpdateMutation = (
  { __typename?: 'Mutation' }
  & { bulkSceneUpdate: Maybe<Array<(
    { __typename?: 'Scene' }
    & SceneDataFragment
  )>> }
);

export type ScenesUpdateMutationVariables = {
  input: Array<SceneUpdateInput>
};


export type ScenesUpdateMutation = (
  { __typename?: 'Mutation' }
  & { scenesUpdate: Maybe<Array<Maybe<(
    { __typename?: 'Scene' }
    & SceneDataFragment
  )>>> }
);

export type SceneIncrementOMutationVariables = {
  id: Scalars['ID']
};


export type SceneIncrementOMutation = (
  { __typename?: 'Mutation' }
  & Pick<Mutation, 'sceneIncrementO'>
);

export type SceneDecrementOMutationVariables = {
  id: Scalars['ID']
};


export type SceneDecrementOMutation = (
  { __typename?: 'Mutation' }
  & Pick<Mutation, 'sceneDecrementO'>
);

export type SceneResetOMutationVariables = {
  id: Scalars['ID']
};


export type SceneResetOMutation = (
  { __typename?: 'Mutation' }
  & Pick<Mutation, 'sceneResetO'>
);

export type SceneDestroyMutationVariables = {
  id: Scalars['ID'],
  delete_file?: Maybe<Scalars['Boolean']>,
  delete_generated?: Maybe<Scalars['Boolean']>
};


export type SceneDestroyMutation = (
  { __typename?: 'Mutation' }
  & Pick<Mutation, 'sceneDestroy'>
);

export type StudioCreateMutationVariables = {
  name: Scalars['String'],
  url?: Maybe<Scalars['String']>,
  image?: Maybe<Scalars['String']>
};


export type StudioCreateMutation = (
  { __typename?: 'Mutation' }
  & { studioCreate: Maybe<(
    { __typename?: 'Studio' }
    & StudioDataFragment
  )> }
);

export type StudioUpdateMutationVariables = {
  id: Scalars['ID'],
  name?: Maybe<Scalars['String']>,
  url?: Maybe<Scalars['String']>,
  image?: Maybe<Scalars['String']>
};


export type StudioUpdateMutation = (
  { __typename?: 'Mutation' }
  & { studioUpdate: Maybe<(
    { __typename?: 'Studio' }
    & StudioDataFragment
  )> }
);

export type StudioDestroyMutationVariables = {
  id: Scalars['ID']
};


export type StudioDestroyMutation = (
  { __typename?: 'Mutation' }
  & Pick<Mutation, 'studioDestroy'>
);

export type TagCreateMutationVariables = {
  name: Scalars['String']
};


export type TagCreateMutation = (
  { __typename?: 'Mutation' }
  & { tagCreate: Maybe<(
    { __typename?: 'Tag' }
    & TagDataFragment
  )> }
);

export type TagDestroyMutationVariables = {
  id: Scalars['ID']
};


export type TagDestroyMutation = (
  { __typename?: 'Mutation' }
  & Pick<Mutation, 'tagDestroy'>
);

export type TagUpdateMutationVariables = {
  id: Scalars['ID'],
  name: Scalars['String']
};


export type TagUpdateMutation = (
  { __typename?: 'Mutation' }
  & { tagUpdate: Maybe<(
    { __typename?: 'Tag' }
    & TagDataFragment
  )> }
);

export type FindGalleriesQueryVariables = {
  filter?: Maybe<FindFilterType>
};


export type FindGalleriesQuery = (
  { __typename?: 'Query' }
  & { findGalleries: (
    { __typename?: 'FindGalleriesResultType' }
    & Pick<FindGalleriesResultType, 'count'>
    & { galleries: Array<(
      { __typename?: 'Gallery' }
      & GalleryDataFragment
    )> }
  ) }
);

export type FindGalleryQueryVariables = {
  id: Scalars['ID']
};


export type FindGalleryQuery = (
  { __typename?: 'Query' }
  & { findGallery: Maybe<(
    { __typename?: 'Gallery' }
    & GalleryDataFragment
  )> }
);

export type SceneWallQueryVariables = {
  q?: Maybe<Scalars['String']>
};


export type SceneWallQuery = (
  { __typename?: 'Query' }
  & { sceneWall: Array<(
    { __typename?: 'Scene' }
    & SceneDataFragment
  )> }
);

export type MarkerWallQueryVariables = {
  q?: Maybe<Scalars['String']>
};


export type MarkerWallQuery = (
  { __typename?: 'Query' }
  & { markerWall: Array<(
    { __typename?: 'SceneMarker' }
    & SceneMarkerDataFragment
  )> }
);

export type FindTagQueryVariables = {
  id: Scalars['ID']
};


export type FindTagQuery = (
  { __typename?: 'Query' }
  & { findTag: Maybe<(
    { __typename?: 'Tag' }
    & TagDataFragment
  )> }
);

export type MarkerStringsQueryVariables = {
  q?: Maybe<Scalars['String']>,
  sort?: Maybe<Scalars['String']>
};


export type MarkerStringsQuery = (
  { __typename?: 'Query' }
  & { markerStrings: Array<Maybe<(
    { __typename?: 'MarkerStringsResultType' }
    & Pick<MarkerStringsResultType, 'id' | 'count' | 'title'>
  )>> }
);

export type AllTagsQueryVariables = {};


export type AllTagsQuery = (
  { __typename?: 'Query' }
  & { allTags: Array<(
    { __typename?: 'Tag' }
    & TagDataFragment
  )> }
);

export type AllPerformersForFilterQueryVariables = {};


export type AllPerformersForFilterQuery = (
  { __typename?: 'Query' }
  & { allPerformers: Array<(
    { __typename?: 'Performer' }
    & SlimPerformerDataFragment
  )> }
);

export type AllStudiosForFilterQueryVariables = {};


export type AllStudiosForFilterQuery = (
  { __typename?: 'Query' }
  & { allStudios: Array<(
    { __typename?: 'Studio' }
    & SlimStudioDataFragment
  )> }
);

export type AllTagsForFilterQueryVariables = {};


export type AllTagsForFilterQuery = (
  { __typename?: 'Query' }
  & { allTags: Array<(
    { __typename?: 'Tag' }
    & Pick<Tag, 'id' | 'name'>
  )> }
);

export type ValidGalleriesForSceneQueryVariables = {
  scene_id: Scalars['ID']
};


export type ValidGalleriesForSceneQuery = (
  { __typename?: 'Query' }
  & { validGalleriesForScene: Array<(
    { __typename?: 'Gallery' }
    & Pick<Gallery, 'id' | 'path'>
  )> }
);

export type StatsQueryVariables = {};


export type StatsQuery = (
  { __typename?: 'Query' }
  & { stats: (
    { __typename?: 'StatsResultType' }
    & Pick<StatsResultType, 'scene_count' | 'gallery_count' | 'performer_count' | 'studio_count' | 'tag_count'>
  ) }
);

export type LogsQueryVariables = {};


export type LogsQuery = (
  { __typename?: 'Query' }
  & { logs: Array<(
    { __typename?: 'LogEntry' }
    & LogEntryDataFragment
  )> }
);

export type VersionQueryVariables = {};


export type VersionQuery = (
  { __typename?: 'Query' }
  & { version: (
    { __typename?: 'Version' }
    & Pick<Version, 'version' | 'hash' | 'build_time'>
  ) }
);

export type LatestVersionQueryVariables = {};


export type LatestVersionQuery = (
  { __typename?: 'Query' }
  & { latestversion: (
    { __typename?: 'ShortVersion' }
    & Pick<ShortVersion, 'shorthash' | 'url'>
  ) }
);

export type FindPerformersQueryVariables = {
  filter?: Maybe<FindFilterType>,
  performer_filter?: Maybe<PerformerFilterType>
};


export type FindPerformersQuery = (
  { __typename?: 'Query' }
  & { findPerformers: (
    { __typename?: 'FindPerformersResultType' }
    & Pick<FindPerformersResultType, 'count'>
    & { performers: Array<(
      { __typename?: 'Performer' }
      & PerformerDataFragment
    )> }
  ) }
);

export type FindPerformerQueryVariables = {
  id: Scalars['ID']
};


export type FindPerformerQuery = (
  { __typename?: 'Query' }
  & { findPerformer: Maybe<(
    { __typename?: 'Performer' }
    & PerformerDataFragment
  )> }
);

export type FindSceneMarkersQueryVariables = {
  filter?: Maybe<FindFilterType>,
  scene_marker_filter?: Maybe<SceneMarkerFilterType>
};


export type FindSceneMarkersQuery = (
  { __typename?: 'Query' }
  & { findSceneMarkers: (
    { __typename?: 'FindSceneMarkersResultType' }
    & Pick<FindSceneMarkersResultType, 'count'>
    & { scene_markers: Array<(
      { __typename?: 'SceneMarker' }
      & SceneMarkerDataFragment
    )> }
  ) }
);

export type FindScenesQueryVariables = {
  filter?: Maybe<FindFilterType>,
  scene_filter?: Maybe<SceneFilterType>,
  scene_ids?: Maybe<Array<Scalars['Int']>>
};


export type FindScenesQuery = (
  { __typename?: 'Query' }
  & { findScenes: (
    { __typename?: 'FindScenesResultType' }
    & Pick<FindScenesResultType, 'count'>
    & { scenes: Array<(
      { __typename?: 'Scene' }
      & SlimSceneDataFragment
    )> }
  ) }
);

export type FindScenesByPathRegexQueryVariables = {
  filter?: Maybe<FindFilterType>
};


export type FindScenesByPathRegexQuery = (
  { __typename?: 'Query' }
  & { findScenesByPathRegex: (
    { __typename?: 'FindScenesResultType' }
    & Pick<FindScenesResultType, 'count'>
    & { scenes: Array<(
      { __typename?: 'Scene' }
      & SlimSceneDataFragment
    )> }
  ) }
);

export type FindSceneQueryVariables = {
  id: Scalars['ID'],
  checksum?: Maybe<Scalars['String']>
};


export type FindSceneQuery = (
  { __typename?: 'Query' }
  & { findScene: Maybe<(
    { __typename?: 'Scene' }
    & SceneDataFragment
  )>, sceneMarkerTags: Array<(
    { __typename?: 'SceneMarkerTag' }
    & { tag: (
      { __typename?: 'Tag' }
      & Pick<Tag, 'id' | 'name'>
    ), scene_markers: Array<(
      { __typename?: 'SceneMarker' }
      & SceneMarkerDataFragment
    )> }
  )> }
);

export type ParseSceneFilenamesQueryVariables = {
  filter: FindFilterType,
  config: SceneParserInput
};


export type ParseSceneFilenamesQuery = (
  { __typename?: 'Query' }
  & { parseSceneFilenames: (
    { __typename?: 'SceneParserResultType' }
    & Pick<SceneParserResultType, 'count'>
    & { results: Array<(
      { __typename?: 'SceneParserResult' }
      & Pick<SceneParserResult, 'title' | 'details' | 'url' | 'date' | 'rating' | 'studio_id' | 'gallery_id' | 'performer_ids' | 'tag_ids'>
      & { scene: (
        { __typename?: 'Scene' }
        & SlimSceneDataFragment
      ) }
    )> }
  ) }
);

export type ScrapeFreeonesQueryVariables = {
  performer_name: Scalars['String']
};


export type ScrapeFreeonesQuery = (
  { __typename?: 'Query' }
  & { scrapeFreeones: Maybe<(
    { __typename?: 'ScrapedPerformer' }
    & Pick<ScrapedPerformer, 'name' | 'url' | 'twitter' | 'instagram' | 'birthdate' | 'ethnicity' | 'country' | 'eye_color' | 'height' | 'measurements' | 'fake_tits' | 'career_length' | 'tattoos' | 'piercings' | 'aliases'>
  )> }
);

export type ScrapeFreeonesPerformersQueryVariables = {
  q: Scalars['String']
};


export type ScrapeFreeonesPerformersQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'scrapeFreeonesPerformerList'>
);

export type ListPerformerScrapersQueryVariables = {};


export type ListPerformerScrapersQuery = (
  { __typename?: 'Query' }
  & { listPerformerScrapers: Array<(
    { __typename?: 'Scraper' }
    & Pick<Scraper, 'id' | 'name'>
    & { performer: Maybe<(
      { __typename?: 'ScraperSpec' }
      & Pick<ScraperSpec, 'urls' | 'supported_scrapes'>
    )> }
  )> }
);

export type ListSceneScrapersQueryVariables = {};


export type ListSceneScrapersQuery = (
  { __typename?: 'Query' }
  & { listSceneScrapers: Array<(
    { __typename?: 'Scraper' }
    & Pick<Scraper, 'id' | 'name'>
    & { scene: Maybe<(
      { __typename?: 'ScraperSpec' }
      & Pick<ScraperSpec, 'urls' | 'supported_scrapes'>
    )> }
  )> }
);

export type ScrapePerformerListQueryVariables = {
  scraper_id: Scalars['ID'],
  query: Scalars['String']
};


export type ScrapePerformerListQuery = (
  { __typename?: 'Query' }
  & { scrapePerformerList: Array<(
    { __typename?: 'ScrapedPerformer' }
    & ScrapedPerformerDataFragment
  )> }
);

export type ScrapePerformerQueryVariables = {
  scraper_id: Scalars['ID'],
  scraped_performer: ScrapedPerformerInput
};


export type ScrapePerformerQuery = (
  { __typename?: 'Query' }
  & { scrapePerformer: Maybe<(
    { __typename?: 'ScrapedPerformer' }
    & ScrapedPerformerDataFragment
  )> }
);

export type ScrapePerformerUrlQueryVariables = {
  url: Scalars['String']
};


export type ScrapePerformerUrlQuery = (
  { __typename?: 'Query' }
  & { scrapePerformerURL: Maybe<(
    { __typename?: 'ScrapedPerformer' }
    & ScrapedPerformerDataFragment
  )> }
);

export type ScrapeSceneQueryVariables = {
  scraper_id: Scalars['ID'],
  scene: SceneUpdateInput
};


export type ScrapeSceneQuery = (
  { __typename?: 'Query' }
  & { scrapeScene: Maybe<(
    { __typename?: 'ScrapedScene' }
    & ScrapedSceneDataFragment
  )> }
);

export type ScrapeSceneUrlQueryVariables = {
  url: Scalars['String']
};


export type ScrapeSceneUrlQuery = (
  { __typename?: 'Query' }
  & { scrapeSceneURL: Maybe<(
    { __typename?: 'ScrapedScene' }
    & ScrapedSceneDataFragment
  )> }
);

export type ConfigurationQueryVariables = {};


export type ConfigurationQuery = (
  { __typename?: 'Query' }
  & { configuration: (
    { __typename?: 'ConfigResult' }
    & ConfigDataFragment
  ) }
);

export type DirectoriesQueryVariables = {
  path?: Maybe<Scalars['String']>
};


export type DirectoriesQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'directories'>
);

export type MetadataImportQueryVariables = {};


export type MetadataImportQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'metadataImport'>
);

export type MetadataExportQueryVariables = {};


export type MetadataExportQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'metadataExport'>
);

export type MetadataScanQueryVariables = {
  input: ScanMetadataInput
};


export type MetadataScanQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'metadataScan'>
);

export type MetadataGenerateQueryVariables = {
  input: GenerateMetadataInput
};


export type MetadataGenerateQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'metadataGenerate'>
);

export type MetadataAutoTagQueryVariables = {
  input: AutoTagMetadataInput
};


export type MetadataAutoTagQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'metadataAutoTag'>
);

export type MetadataCleanQueryVariables = {};


export type MetadataCleanQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'metadataClean'>
);

export type JobStatusQueryVariables = {};


export type JobStatusQuery = (
  { __typename?: 'Query' }
  & { jobStatus: (
    { __typename?: 'MetadataUpdateStatus' }
    & Pick<MetadataUpdateStatus, 'progress' | 'status' | 'message'>
  ) }
);

export type StopJobQueryVariables = {};


export type StopJobQuery = (
  { __typename?: 'Query' }
  & Pick<Query, 'stopJob'>
);

export type FindStudiosQueryVariables = {
  filter?: Maybe<FindFilterType>
};


export type FindStudiosQuery = (
  { __typename?: 'Query' }
  & { findStudios: (
    { __typename?: 'FindStudiosResultType' }
    & Pick<FindStudiosResultType, 'count'>
    & { studios: Array<(
      { __typename?: 'Studio' }
      & StudioDataFragment
    )> }
  ) }
);

export type FindStudioQueryVariables = {
  id: Scalars['ID']
};


export type FindStudioQuery = (
  { __typename?: 'Query' }
  & { findStudio: Maybe<(
    { __typename?: 'Studio' }
    & StudioDataFragment
  )> }
);

export type MetadataUpdateSubscriptionVariables = {};


export type MetadataUpdateSubscription = (
  { __typename?: 'Subscription' }
  & { metadataUpdate: (
    { __typename?: 'MetadataUpdateStatus' }
    & Pick<MetadataUpdateStatus, 'progress' | 'status' | 'message'>
  ) }
);

export type LoggingSubscribeSubscriptionVariables = {};


export type LoggingSubscribeSubscription = (
  { __typename?: 'Subscription' }
  & { loggingSubscribe: Array<(
    { __typename?: 'LogEntry' }
    & LogEntryDataFragment
  )> }
);

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
  language
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
${ConfigInterfaceDataFragmentDoc}`;
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
  o_counter
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
  o_counter
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
${PerformerDataFragmentDoc}`;
export const ScrapedPerformerDataFragmentDoc = gql`
    fragment ScrapedPerformerData on ScrapedPerformer {
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
${ScrapedScenePerformerDataFragmentDoc}`;
export const SlimStudioDataFragmentDoc = gql`
    fragment SlimStudioData on Studio {
  id
  name
  image_path
}
    `;
export const ConfigureGeneralDocument = gql`
    mutation ConfigureGeneral($input: ConfigGeneralInput!) {
  configureGeneral(input: $input) {
    ...ConfigGeneralData
  }
}
    ${ConfigGeneralDataFragmentDoc}`;
export type ConfigureGeneralMutationFn = ApolloReactCommon.MutationFunction<ConfigureGeneralMutation, ConfigureGeneralMutationVariables>;
export type ConfigureGeneralComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<ConfigureGeneralMutation, ConfigureGeneralMutationVariables>, 'mutation'>;

    export const ConfigureGeneralComponent = (props: ConfigureGeneralComponentProps) => (
      <ApolloReactComponents.Mutation<ConfigureGeneralMutation, ConfigureGeneralMutationVariables> mutation={ConfigureGeneralDocument} {...props} />
    );
    

/**
 * __useConfigureGeneralMutation__
 *
 * To run a mutation, you first call `useConfigureGeneralMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useConfigureGeneralMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [configureGeneralMutation, { data, loading, error }] = useConfigureGeneralMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useConfigureGeneralMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<ConfigureGeneralMutation, ConfigureGeneralMutationVariables>) {
        return ApolloReactHooks.useMutation<ConfigureGeneralMutation, ConfigureGeneralMutationVariables>(ConfigureGeneralDocument, baseOptions);
      }
export type ConfigureGeneralMutationHookResult = ReturnType<typeof useConfigureGeneralMutation>;
export type ConfigureGeneralMutationResult = ApolloReactCommon.MutationResult<ConfigureGeneralMutation>;
export type ConfigureGeneralMutationOptions = ApolloReactCommon.BaseMutationOptions<ConfigureGeneralMutation, ConfigureGeneralMutationVariables>;
export const ConfigureInterfaceDocument = gql`
    mutation ConfigureInterface($input: ConfigInterfaceInput!) {
  configureInterface(input: $input) {
    ...ConfigInterfaceData
  }
}
    ${ConfigInterfaceDataFragmentDoc}`;
export type ConfigureInterfaceMutationFn = ApolloReactCommon.MutationFunction<ConfigureInterfaceMutation, ConfigureInterfaceMutationVariables>;
export type ConfigureInterfaceComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<ConfigureInterfaceMutation, ConfigureInterfaceMutationVariables>, 'mutation'>;

    export const ConfigureInterfaceComponent = (props: ConfigureInterfaceComponentProps) => (
      <ApolloReactComponents.Mutation<ConfigureInterfaceMutation, ConfigureInterfaceMutationVariables> mutation={ConfigureInterfaceDocument} {...props} />
    );
    

/**
 * __useConfigureInterfaceMutation__
 *
 * To run a mutation, you first call `useConfigureInterfaceMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useConfigureInterfaceMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [configureInterfaceMutation, { data, loading, error }] = useConfigureInterfaceMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useConfigureInterfaceMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<ConfigureInterfaceMutation, ConfigureInterfaceMutationVariables>) {
        return ApolloReactHooks.useMutation<ConfigureInterfaceMutation, ConfigureInterfaceMutationVariables>(ConfigureInterfaceDocument, baseOptions);
      }
export type ConfigureInterfaceMutationHookResult = ReturnType<typeof useConfigureInterfaceMutation>;
export type ConfigureInterfaceMutationResult = ApolloReactCommon.MutationResult<ConfigureInterfaceMutation>;
export type ConfigureInterfaceMutationOptions = ApolloReactCommon.BaseMutationOptions<ConfigureInterfaceMutation, ConfigureInterfaceMutationVariables>;
export const PerformerCreateDocument = gql`
    mutation PerformerCreate($name: String, $url: String, $birthdate: String, $ethnicity: String, $country: String, $eye_color: String, $height: String, $measurements: String, $fake_tits: String, $career_length: String, $tattoos: String, $piercings: String, $aliases: String, $twitter: String, $instagram: String, $favorite: Boolean, $image: String) {
  performerCreate(input: {name: $name, url: $url, birthdate: $birthdate, ethnicity: $ethnicity, country: $country, eye_color: $eye_color, height: $height, measurements: $measurements, fake_tits: $fake_tits, career_length: $career_length, tattoos: $tattoos, piercings: $piercings, aliases: $aliases, twitter: $twitter, instagram: $instagram, favorite: $favorite, image: $image}) {
    ...PerformerData
  }
}
    ${PerformerDataFragmentDoc}`;
export type PerformerCreateMutationFn = ApolloReactCommon.MutationFunction<PerformerCreateMutation, PerformerCreateMutationVariables>;
export type PerformerCreateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<PerformerCreateMutation, PerformerCreateMutationVariables>, 'mutation'>;

    export const PerformerCreateComponent = (props: PerformerCreateComponentProps) => (
      <ApolloReactComponents.Mutation<PerformerCreateMutation, PerformerCreateMutationVariables> mutation={PerformerCreateDocument} {...props} />
    );
    

/**
 * __usePerformerCreateMutation__
 *
 * To run a mutation, you first call `usePerformerCreateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `usePerformerCreateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [performerCreateMutation, { data, loading, error }] = usePerformerCreateMutation({
 *   variables: {
 *      name: // value for 'name'
 *      url: // value for 'url'
 *      birthdate: // value for 'birthdate'
 *      ethnicity: // value for 'ethnicity'
 *      country: // value for 'country'
 *      eye_color: // value for 'eye_color'
 *      height: // value for 'height'
 *      measurements: // value for 'measurements'
 *      fake_tits: // value for 'fake_tits'
 *      career_length: // value for 'career_length'
 *      tattoos: // value for 'tattoos'
 *      piercings: // value for 'piercings'
 *      aliases: // value for 'aliases'
 *      twitter: // value for 'twitter'
 *      instagram: // value for 'instagram'
 *      favorite: // value for 'favorite'
 *      image: // value for 'image'
 *   },
 * });
 */
export function usePerformerCreateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<PerformerCreateMutation, PerformerCreateMutationVariables>) {
        return ApolloReactHooks.useMutation<PerformerCreateMutation, PerformerCreateMutationVariables>(PerformerCreateDocument, baseOptions);
      }
export type PerformerCreateMutationHookResult = ReturnType<typeof usePerformerCreateMutation>;
export type PerformerCreateMutationResult = ApolloReactCommon.MutationResult<PerformerCreateMutation>;
export type PerformerCreateMutationOptions = ApolloReactCommon.BaseMutationOptions<PerformerCreateMutation, PerformerCreateMutationVariables>;
export const PerformerUpdateDocument = gql`
    mutation PerformerUpdate($id: ID!, $name: String, $url: String, $birthdate: String, $ethnicity: String, $country: String, $eye_color: String, $height: String, $measurements: String, $fake_tits: String, $career_length: String, $tattoos: String, $piercings: String, $aliases: String, $twitter: String, $instagram: String, $favorite: Boolean, $image: String) {
  performerUpdate(input: {id: $id, name: $name, url: $url, birthdate: $birthdate, ethnicity: $ethnicity, country: $country, eye_color: $eye_color, height: $height, measurements: $measurements, fake_tits: $fake_tits, career_length: $career_length, tattoos: $tattoos, piercings: $piercings, aliases: $aliases, twitter: $twitter, instagram: $instagram, favorite: $favorite, image: $image}) {
    ...PerformerData
  }
}
    ${PerformerDataFragmentDoc}`;
export type PerformerUpdateMutationFn = ApolloReactCommon.MutationFunction<PerformerUpdateMutation, PerformerUpdateMutationVariables>;
export type PerformerUpdateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<PerformerUpdateMutation, PerformerUpdateMutationVariables>, 'mutation'>;

    export const PerformerUpdateComponent = (props: PerformerUpdateComponentProps) => (
      <ApolloReactComponents.Mutation<PerformerUpdateMutation, PerformerUpdateMutationVariables> mutation={PerformerUpdateDocument} {...props} />
    );
    

/**
 * __usePerformerUpdateMutation__
 *
 * To run a mutation, you first call `usePerformerUpdateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `usePerformerUpdateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [performerUpdateMutation, { data, loading, error }] = usePerformerUpdateMutation({
 *   variables: {
 *      id: // value for 'id'
 *      name: // value for 'name'
 *      url: // value for 'url'
 *      birthdate: // value for 'birthdate'
 *      ethnicity: // value for 'ethnicity'
 *      country: // value for 'country'
 *      eye_color: // value for 'eye_color'
 *      height: // value for 'height'
 *      measurements: // value for 'measurements'
 *      fake_tits: // value for 'fake_tits'
 *      career_length: // value for 'career_length'
 *      tattoos: // value for 'tattoos'
 *      piercings: // value for 'piercings'
 *      aliases: // value for 'aliases'
 *      twitter: // value for 'twitter'
 *      instagram: // value for 'instagram'
 *      favorite: // value for 'favorite'
 *      image: // value for 'image'
 *   },
 * });
 */
export function usePerformerUpdateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<PerformerUpdateMutation, PerformerUpdateMutationVariables>) {
        return ApolloReactHooks.useMutation<PerformerUpdateMutation, PerformerUpdateMutationVariables>(PerformerUpdateDocument, baseOptions);
      }
export type PerformerUpdateMutationHookResult = ReturnType<typeof usePerformerUpdateMutation>;
export type PerformerUpdateMutationResult = ApolloReactCommon.MutationResult<PerformerUpdateMutation>;
export type PerformerUpdateMutationOptions = ApolloReactCommon.BaseMutationOptions<PerformerUpdateMutation, PerformerUpdateMutationVariables>;
export const PerformerDestroyDocument = gql`
    mutation PerformerDestroy($id: ID!) {
  performerDestroy(input: {id: $id})
}
    `;
export type PerformerDestroyMutationFn = ApolloReactCommon.MutationFunction<PerformerDestroyMutation, PerformerDestroyMutationVariables>;
export type PerformerDestroyComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<PerformerDestroyMutation, PerformerDestroyMutationVariables>, 'mutation'>;

    export const PerformerDestroyComponent = (props: PerformerDestroyComponentProps) => (
      <ApolloReactComponents.Mutation<PerformerDestroyMutation, PerformerDestroyMutationVariables> mutation={PerformerDestroyDocument} {...props} />
    );
    

/**
 * __usePerformerDestroyMutation__
 *
 * To run a mutation, you first call `usePerformerDestroyMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `usePerformerDestroyMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [performerDestroyMutation, { data, loading, error }] = usePerformerDestroyMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function usePerformerDestroyMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<PerformerDestroyMutation, PerformerDestroyMutationVariables>) {
        return ApolloReactHooks.useMutation<PerformerDestroyMutation, PerformerDestroyMutationVariables>(PerformerDestroyDocument, baseOptions);
      }
export type PerformerDestroyMutationHookResult = ReturnType<typeof usePerformerDestroyMutation>;
export type PerformerDestroyMutationResult = ApolloReactCommon.MutationResult<PerformerDestroyMutation>;
export type PerformerDestroyMutationOptions = ApolloReactCommon.BaseMutationOptions<PerformerDestroyMutation, PerformerDestroyMutationVariables>;
export const SceneMarkerCreateDocument = gql`
    mutation SceneMarkerCreate($title: String!, $seconds: Float!, $scene_id: ID!, $primary_tag_id: ID!, $tag_ids: [ID!] = []) {
  sceneMarkerCreate(input: {title: $title, seconds: $seconds, scene_id: $scene_id, primary_tag_id: $primary_tag_id, tag_ids: $tag_ids}) {
    ...SceneMarkerData
  }
}
    ${SceneMarkerDataFragmentDoc}`;
export type SceneMarkerCreateMutationFn = ApolloReactCommon.MutationFunction<SceneMarkerCreateMutation, SceneMarkerCreateMutationVariables>;
export type SceneMarkerCreateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<SceneMarkerCreateMutation, SceneMarkerCreateMutationVariables>, 'mutation'>;

    export const SceneMarkerCreateComponent = (props: SceneMarkerCreateComponentProps) => (
      <ApolloReactComponents.Mutation<SceneMarkerCreateMutation, SceneMarkerCreateMutationVariables> mutation={SceneMarkerCreateDocument} {...props} />
    );
    

/**
 * __useSceneMarkerCreateMutation__
 *
 * To run a mutation, you first call `useSceneMarkerCreateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSceneMarkerCreateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [sceneMarkerCreateMutation, { data, loading, error }] = useSceneMarkerCreateMutation({
 *   variables: {
 *      title: // value for 'title'
 *      seconds: // value for 'seconds'
 *      scene_id: // value for 'scene_id'
 *      primary_tag_id: // value for 'primary_tag_id'
 *      tag_ids: // value for 'tag_ids'
 *   },
 * });
 */
export function useSceneMarkerCreateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<SceneMarkerCreateMutation, SceneMarkerCreateMutationVariables>) {
        return ApolloReactHooks.useMutation<SceneMarkerCreateMutation, SceneMarkerCreateMutationVariables>(SceneMarkerCreateDocument, baseOptions);
      }
export type SceneMarkerCreateMutationHookResult = ReturnType<typeof useSceneMarkerCreateMutation>;
export type SceneMarkerCreateMutationResult = ApolloReactCommon.MutationResult<SceneMarkerCreateMutation>;
export type SceneMarkerCreateMutationOptions = ApolloReactCommon.BaseMutationOptions<SceneMarkerCreateMutation, SceneMarkerCreateMutationVariables>;
export const SceneMarkerUpdateDocument = gql`
    mutation SceneMarkerUpdate($id: ID!, $title: String!, $seconds: Float!, $scene_id: ID!, $primary_tag_id: ID!, $tag_ids: [ID!] = []) {
  sceneMarkerUpdate(input: {id: $id, title: $title, seconds: $seconds, scene_id: $scene_id, primary_tag_id: $primary_tag_id, tag_ids: $tag_ids}) {
    ...SceneMarkerData
  }
}
    ${SceneMarkerDataFragmentDoc}`;
export type SceneMarkerUpdateMutationFn = ApolloReactCommon.MutationFunction<SceneMarkerUpdateMutation, SceneMarkerUpdateMutationVariables>;
export type SceneMarkerUpdateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<SceneMarkerUpdateMutation, SceneMarkerUpdateMutationVariables>, 'mutation'>;

    export const SceneMarkerUpdateComponent = (props: SceneMarkerUpdateComponentProps) => (
      <ApolloReactComponents.Mutation<SceneMarkerUpdateMutation, SceneMarkerUpdateMutationVariables> mutation={SceneMarkerUpdateDocument} {...props} />
    );
    

/**
 * __useSceneMarkerUpdateMutation__
 *
 * To run a mutation, you first call `useSceneMarkerUpdateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSceneMarkerUpdateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [sceneMarkerUpdateMutation, { data, loading, error }] = useSceneMarkerUpdateMutation({
 *   variables: {
 *      id: // value for 'id'
 *      title: // value for 'title'
 *      seconds: // value for 'seconds'
 *      scene_id: // value for 'scene_id'
 *      primary_tag_id: // value for 'primary_tag_id'
 *      tag_ids: // value for 'tag_ids'
 *   },
 * });
 */
export function useSceneMarkerUpdateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<SceneMarkerUpdateMutation, SceneMarkerUpdateMutationVariables>) {
        return ApolloReactHooks.useMutation<SceneMarkerUpdateMutation, SceneMarkerUpdateMutationVariables>(SceneMarkerUpdateDocument, baseOptions);
      }
export type SceneMarkerUpdateMutationHookResult = ReturnType<typeof useSceneMarkerUpdateMutation>;
export type SceneMarkerUpdateMutationResult = ApolloReactCommon.MutationResult<SceneMarkerUpdateMutation>;
export type SceneMarkerUpdateMutationOptions = ApolloReactCommon.BaseMutationOptions<SceneMarkerUpdateMutation, SceneMarkerUpdateMutationVariables>;
export const SceneMarkerDestroyDocument = gql`
    mutation SceneMarkerDestroy($id: ID!) {
  sceneMarkerDestroy(id: $id)
}
    `;
export type SceneMarkerDestroyMutationFn = ApolloReactCommon.MutationFunction<SceneMarkerDestroyMutation, SceneMarkerDestroyMutationVariables>;
export type SceneMarkerDestroyComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<SceneMarkerDestroyMutation, SceneMarkerDestroyMutationVariables>, 'mutation'>;

    export const SceneMarkerDestroyComponent = (props: SceneMarkerDestroyComponentProps) => (
      <ApolloReactComponents.Mutation<SceneMarkerDestroyMutation, SceneMarkerDestroyMutationVariables> mutation={SceneMarkerDestroyDocument} {...props} />
    );
    

/**
 * __useSceneMarkerDestroyMutation__
 *
 * To run a mutation, you first call `useSceneMarkerDestroyMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSceneMarkerDestroyMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [sceneMarkerDestroyMutation, { data, loading, error }] = useSceneMarkerDestroyMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useSceneMarkerDestroyMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<SceneMarkerDestroyMutation, SceneMarkerDestroyMutationVariables>) {
        return ApolloReactHooks.useMutation<SceneMarkerDestroyMutation, SceneMarkerDestroyMutationVariables>(SceneMarkerDestroyDocument, baseOptions);
      }
export type SceneMarkerDestroyMutationHookResult = ReturnType<typeof useSceneMarkerDestroyMutation>;
export type SceneMarkerDestroyMutationResult = ApolloReactCommon.MutationResult<SceneMarkerDestroyMutation>;
export type SceneMarkerDestroyMutationOptions = ApolloReactCommon.BaseMutationOptions<SceneMarkerDestroyMutation, SceneMarkerDestroyMutationVariables>;
export const SceneUpdateDocument = gql`
    mutation SceneUpdate($id: ID!, $title: String, $details: String, $url: String, $date: String, $rating: Int, $studio_id: ID, $gallery_id: ID, $performer_ids: [ID!] = [], $tag_ids: [ID!] = [], $cover_image: String) {
  sceneUpdate(input: {id: $id, title: $title, details: $details, url: $url, date: $date, rating: $rating, studio_id: $studio_id, gallery_id: $gallery_id, performer_ids: $performer_ids, tag_ids: $tag_ids, cover_image: $cover_image}) {
    ...SceneData
  }
}
    ${SceneDataFragmentDoc}`;
export type SceneUpdateMutationFn = ApolloReactCommon.MutationFunction<SceneUpdateMutation, SceneUpdateMutationVariables>;
export type SceneUpdateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<SceneUpdateMutation, SceneUpdateMutationVariables>, 'mutation'>;

    export const SceneUpdateComponent = (props: SceneUpdateComponentProps) => (
      <ApolloReactComponents.Mutation<SceneUpdateMutation, SceneUpdateMutationVariables> mutation={SceneUpdateDocument} {...props} />
    );
    

/**
 * __useSceneUpdateMutation__
 *
 * To run a mutation, you first call `useSceneUpdateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSceneUpdateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [sceneUpdateMutation, { data, loading, error }] = useSceneUpdateMutation({
 *   variables: {
 *      id: // value for 'id'
 *      title: // value for 'title'
 *      details: // value for 'details'
 *      url: // value for 'url'
 *      date: // value for 'date'
 *      rating: // value for 'rating'
 *      studio_id: // value for 'studio_id'
 *      gallery_id: // value for 'gallery_id'
 *      performer_ids: // value for 'performer_ids'
 *      tag_ids: // value for 'tag_ids'
 *      cover_image: // value for 'cover_image'
 *   },
 * });
 */
export function useSceneUpdateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<SceneUpdateMutation, SceneUpdateMutationVariables>) {
        return ApolloReactHooks.useMutation<SceneUpdateMutation, SceneUpdateMutationVariables>(SceneUpdateDocument, baseOptions);
      }
export type SceneUpdateMutationHookResult = ReturnType<typeof useSceneUpdateMutation>;
export type SceneUpdateMutationResult = ApolloReactCommon.MutationResult<SceneUpdateMutation>;
export type SceneUpdateMutationOptions = ApolloReactCommon.BaseMutationOptions<SceneUpdateMutation, SceneUpdateMutationVariables>;
export const BulkSceneUpdateDocument = gql`
    mutation BulkSceneUpdate($ids: [ID!] = [], $title: String, $details: String, $url: String, $date: String, $rating: Int, $studio_id: ID, $gallery_id: ID, $performer_ids: [ID!], $tag_ids: [ID!]) {
  bulkSceneUpdate(input: {ids: $ids, title: $title, details: $details, url: $url, date: $date, rating: $rating, studio_id: $studio_id, gallery_id: $gallery_id, performer_ids: $performer_ids, tag_ids: $tag_ids}) {
    ...SceneData
  }
}
    ${SceneDataFragmentDoc}`;
export type BulkSceneUpdateMutationFn = ApolloReactCommon.MutationFunction<BulkSceneUpdateMutation, BulkSceneUpdateMutationVariables>;
export type BulkSceneUpdateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<BulkSceneUpdateMutation, BulkSceneUpdateMutationVariables>, 'mutation'>;

    export const BulkSceneUpdateComponent = (props: BulkSceneUpdateComponentProps) => (
      <ApolloReactComponents.Mutation<BulkSceneUpdateMutation, BulkSceneUpdateMutationVariables> mutation={BulkSceneUpdateDocument} {...props} />
    );
    

/**
 * __useBulkSceneUpdateMutation__
 *
 * To run a mutation, you first call `useBulkSceneUpdateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useBulkSceneUpdateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [bulkSceneUpdateMutation, { data, loading, error }] = useBulkSceneUpdateMutation({
 *   variables: {
 *      ids: // value for 'ids'
 *      title: // value for 'title'
 *      details: // value for 'details'
 *      url: // value for 'url'
 *      date: // value for 'date'
 *      rating: // value for 'rating'
 *      studio_id: // value for 'studio_id'
 *      gallery_id: // value for 'gallery_id'
 *      performer_ids: // value for 'performer_ids'
 *      tag_ids: // value for 'tag_ids'
 *   },
 * });
 */
export function useBulkSceneUpdateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<BulkSceneUpdateMutation, BulkSceneUpdateMutationVariables>) {
        return ApolloReactHooks.useMutation<BulkSceneUpdateMutation, BulkSceneUpdateMutationVariables>(BulkSceneUpdateDocument, baseOptions);
      }
export type BulkSceneUpdateMutationHookResult = ReturnType<typeof useBulkSceneUpdateMutation>;
export type BulkSceneUpdateMutationResult = ApolloReactCommon.MutationResult<BulkSceneUpdateMutation>;
export type BulkSceneUpdateMutationOptions = ApolloReactCommon.BaseMutationOptions<BulkSceneUpdateMutation, BulkSceneUpdateMutationVariables>;
export const ScenesUpdateDocument = gql`
    mutation ScenesUpdate($input: [SceneUpdateInput!]!) {
  scenesUpdate(input: $input) {
    ...SceneData
  }
}
    ${SceneDataFragmentDoc}`;
export type ScenesUpdateMutationFn = ApolloReactCommon.MutationFunction<ScenesUpdateMutation, ScenesUpdateMutationVariables>;
export type ScenesUpdateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<ScenesUpdateMutation, ScenesUpdateMutationVariables>, 'mutation'>;

    export const ScenesUpdateComponent = (props: ScenesUpdateComponentProps) => (
      <ApolloReactComponents.Mutation<ScenesUpdateMutation, ScenesUpdateMutationVariables> mutation={ScenesUpdateDocument} {...props} />
    );
    

/**
 * __useScenesUpdateMutation__
 *
 * To run a mutation, you first call `useScenesUpdateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useScenesUpdateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [scenesUpdateMutation, { data, loading, error }] = useScenesUpdateMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useScenesUpdateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<ScenesUpdateMutation, ScenesUpdateMutationVariables>) {
        return ApolloReactHooks.useMutation<ScenesUpdateMutation, ScenesUpdateMutationVariables>(ScenesUpdateDocument, baseOptions);
      }
export type ScenesUpdateMutationHookResult = ReturnType<typeof useScenesUpdateMutation>;
export type ScenesUpdateMutationResult = ApolloReactCommon.MutationResult<ScenesUpdateMutation>;
export type ScenesUpdateMutationOptions = ApolloReactCommon.BaseMutationOptions<ScenesUpdateMutation, ScenesUpdateMutationVariables>;
export const SceneIncrementODocument = gql`
    mutation SceneIncrementO($id: ID!) {
  sceneIncrementO(id: $id)
}
    `;
export type SceneIncrementOMutationFn = ApolloReactCommon.MutationFunction<SceneIncrementOMutation, SceneIncrementOMutationVariables>;
export type SceneIncrementOComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<SceneIncrementOMutation, SceneIncrementOMutationVariables>, 'mutation'>;

    export const SceneIncrementOComponent = (props: SceneIncrementOComponentProps) => (
      <ApolloReactComponents.Mutation<SceneIncrementOMutation, SceneIncrementOMutationVariables> mutation={SceneIncrementODocument} {...props} />
    );
    

/**
 * __useSceneIncrementOMutation__
 *
 * To run a mutation, you first call `useSceneIncrementOMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSceneIncrementOMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [sceneIncrementOMutation, { data, loading, error }] = useSceneIncrementOMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useSceneIncrementOMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<SceneIncrementOMutation, SceneIncrementOMutationVariables>) {
        return ApolloReactHooks.useMutation<SceneIncrementOMutation, SceneIncrementOMutationVariables>(SceneIncrementODocument, baseOptions);
      }
export type SceneIncrementOMutationHookResult = ReturnType<typeof useSceneIncrementOMutation>;
export type SceneIncrementOMutationResult = ApolloReactCommon.MutationResult<SceneIncrementOMutation>;
export type SceneIncrementOMutationOptions = ApolloReactCommon.BaseMutationOptions<SceneIncrementOMutation, SceneIncrementOMutationVariables>;
export const SceneDecrementODocument = gql`
    mutation SceneDecrementO($id: ID!) {
  sceneDecrementO(id: $id)
}
    `;
export type SceneDecrementOMutationFn = ApolloReactCommon.MutationFunction<SceneDecrementOMutation, SceneDecrementOMutationVariables>;
export type SceneDecrementOComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<SceneDecrementOMutation, SceneDecrementOMutationVariables>, 'mutation'>;

    export const SceneDecrementOComponent = (props: SceneDecrementOComponentProps) => (
      <ApolloReactComponents.Mutation<SceneDecrementOMutation, SceneDecrementOMutationVariables> mutation={SceneDecrementODocument} {...props} />
    );
    

/**
 * __useSceneDecrementOMutation__
 *
 * To run a mutation, you first call `useSceneDecrementOMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSceneDecrementOMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [sceneDecrementOMutation, { data, loading, error }] = useSceneDecrementOMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useSceneDecrementOMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<SceneDecrementOMutation, SceneDecrementOMutationVariables>) {
        return ApolloReactHooks.useMutation<SceneDecrementOMutation, SceneDecrementOMutationVariables>(SceneDecrementODocument, baseOptions);
      }
export type SceneDecrementOMutationHookResult = ReturnType<typeof useSceneDecrementOMutation>;
export type SceneDecrementOMutationResult = ApolloReactCommon.MutationResult<SceneDecrementOMutation>;
export type SceneDecrementOMutationOptions = ApolloReactCommon.BaseMutationOptions<SceneDecrementOMutation, SceneDecrementOMutationVariables>;
export const SceneResetODocument = gql`
    mutation SceneResetO($id: ID!) {
  sceneResetO(id: $id)
}
    `;
export type SceneResetOMutationFn = ApolloReactCommon.MutationFunction<SceneResetOMutation, SceneResetOMutationVariables>;
export type SceneResetOComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<SceneResetOMutation, SceneResetOMutationVariables>, 'mutation'>;

    export const SceneResetOComponent = (props: SceneResetOComponentProps) => (
      <ApolloReactComponents.Mutation<SceneResetOMutation, SceneResetOMutationVariables> mutation={SceneResetODocument} {...props} />
    );
    

/**
 * __useSceneResetOMutation__
 *
 * To run a mutation, you first call `useSceneResetOMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSceneResetOMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [sceneResetOMutation, { data, loading, error }] = useSceneResetOMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useSceneResetOMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<SceneResetOMutation, SceneResetOMutationVariables>) {
        return ApolloReactHooks.useMutation<SceneResetOMutation, SceneResetOMutationVariables>(SceneResetODocument, baseOptions);
      }
export type SceneResetOMutationHookResult = ReturnType<typeof useSceneResetOMutation>;
export type SceneResetOMutationResult = ApolloReactCommon.MutationResult<SceneResetOMutation>;
export type SceneResetOMutationOptions = ApolloReactCommon.BaseMutationOptions<SceneResetOMutation, SceneResetOMutationVariables>;
export const SceneDestroyDocument = gql`
    mutation SceneDestroy($id: ID!, $delete_file: Boolean, $delete_generated: Boolean) {
  sceneDestroy(input: {id: $id, delete_file: $delete_file, delete_generated: $delete_generated})
}
    `;
export type SceneDestroyMutationFn = ApolloReactCommon.MutationFunction<SceneDestroyMutation, SceneDestroyMutationVariables>;
export type SceneDestroyComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<SceneDestroyMutation, SceneDestroyMutationVariables>, 'mutation'>;

    export const SceneDestroyComponent = (props: SceneDestroyComponentProps) => (
      <ApolloReactComponents.Mutation<SceneDestroyMutation, SceneDestroyMutationVariables> mutation={SceneDestroyDocument} {...props} />
    );
    

/**
 * __useSceneDestroyMutation__
 *
 * To run a mutation, you first call `useSceneDestroyMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSceneDestroyMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [sceneDestroyMutation, { data, loading, error }] = useSceneDestroyMutation({
 *   variables: {
 *      id: // value for 'id'
 *      delete_file: // value for 'delete_file'
 *      delete_generated: // value for 'delete_generated'
 *   },
 * });
 */
export function useSceneDestroyMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<SceneDestroyMutation, SceneDestroyMutationVariables>) {
        return ApolloReactHooks.useMutation<SceneDestroyMutation, SceneDestroyMutationVariables>(SceneDestroyDocument, baseOptions);
      }
export type SceneDestroyMutationHookResult = ReturnType<typeof useSceneDestroyMutation>;
export type SceneDestroyMutationResult = ApolloReactCommon.MutationResult<SceneDestroyMutation>;
export type SceneDestroyMutationOptions = ApolloReactCommon.BaseMutationOptions<SceneDestroyMutation, SceneDestroyMutationVariables>;
export const StudioCreateDocument = gql`
    mutation StudioCreate($name: String!, $url: String, $image: String) {
  studioCreate(input: {name: $name, url: $url, image: $image}) {
    ...StudioData
  }
}
    ${StudioDataFragmentDoc}`;
export type StudioCreateMutationFn = ApolloReactCommon.MutationFunction<StudioCreateMutation, StudioCreateMutationVariables>;
export type StudioCreateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<StudioCreateMutation, StudioCreateMutationVariables>, 'mutation'>;

    export const StudioCreateComponent = (props: StudioCreateComponentProps) => (
      <ApolloReactComponents.Mutation<StudioCreateMutation, StudioCreateMutationVariables> mutation={StudioCreateDocument} {...props} />
    );
    

/**
 * __useStudioCreateMutation__
 *
 * To run a mutation, you first call `useStudioCreateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useStudioCreateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [studioCreateMutation, { data, loading, error }] = useStudioCreateMutation({
 *   variables: {
 *      name: // value for 'name'
 *      url: // value for 'url'
 *      image: // value for 'image'
 *   },
 * });
 */
export function useStudioCreateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<StudioCreateMutation, StudioCreateMutationVariables>) {
        return ApolloReactHooks.useMutation<StudioCreateMutation, StudioCreateMutationVariables>(StudioCreateDocument, baseOptions);
      }
export type StudioCreateMutationHookResult = ReturnType<typeof useStudioCreateMutation>;
export type StudioCreateMutationResult = ApolloReactCommon.MutationResult<StudioCreateMutation>;
export type StudioCreateMutationOptions = ApolloReactCommon.BaseMutationOptions<StudioCreateMutation, StudioCreateMutationVariables>;
export const StudioUpdateDocument = gql`
    mutation StudioUpdate($id: ID!, $name: String, $url: String, $image: String) {
  studioUpdate(input: {id: $id, name: $name, url: $url, image: $image}) {
    ...StudioData
  }
}
    ${StudioDataFragmentDoc}`;
export type StudioUpdateMutationFn = ApolloReactCommon.MutationFunction<StudioUpdateMutation, StudioUpdateMutationVariables>;
export type StudioUpdateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<StudioUpdateMutation, StudioUpdateMutationVariables>, 'mutation'>;

    export const StudioUpdateComponent = (props: StudioUpdateComponentProps) => (
      <ApolloReactComponents.Mutation<StudioUpdateMutation, StudioUpdateMutationVariables> mutation={StudioUpdateDocument} {...props} />
    );
    

/**
 * __useStudioUpdateMutation__
 *
 * To run a mutation, you first call `useStudioUpdateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useStudioUpdateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [studioUpdateMutation, { data, loading, error }] = useStudioUpdateMutation({
 *   variables: {
 *      id: // value for 'id'
 *      name: // value for 'name'
 *      url: // value for 'url'
 *      image: // value for 'image'
 *   },
 * });
 */
export function useStudioUpdateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<StudioUpdateMutation, StudioUpdateMutationVariables>) {
        return ApolloReactHooks.useMutation<StudioUpdateMutation, StudioUpdateMutationVariables>(StudioUpdateDocument, baseOptions);
      }
export type StudioUpdateMutationHookResult = ReturnType<typeof useStudioUpdateMutation>;
export type StudioUpdateMutationResult = ApolloReactCommon.MutationResult<StudioUpdateMutation>;
export type StudioUpdateMutationOptions = ApolloReactCommon.BaseMutationOptions<StudioUpdateMutation, StudioUpdateMutationVariables>;
export const StudioDestroyDocument = gql`
    mutation StudioDestroy($id: ID!) {
  studioDestroy(input: {id: $id})
}
    `;
export type StudioDestroyMutationFn = ApolloReactCommon.MutationFunction<StudioDestroyMutation, StudioDestroyMutationVariables>;
export type StudioDestroyComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<StudioDestroyMutation, StudioDestroyMutationVariables>, 'mutation'>;

    export const StudioDestroyComponent = (props: StudioDestroyComponentProps) => (
      <ApolloReactComponents.Mutation<StudioDestroyMutation, StudioDestroyMutationVariables> mutation={StudioDestroyDocument} {...props} />
    );
    

/**
 * __useStudioDestroyMutation__
 *
 * To run a mutation, you first call `useStudioDestroyMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useStudioDestroyMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [studioDestroyMutation, { data, loading, error }] = useStudioDestroyMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useStudioDestroyMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<StudioDestroyMutation, StudioDestroyMutationVariables>) {
        return ApolloReactHooks.useMutation<StudioDestroyMutation, StudioDestroyMutationVariables>(StudioDestroyDocument, baseOptions);
      }
export type StudioDestroyMutationHookResult = ReturnType<typeof useStudioDestroyMutation>;
export type StudioDestroyMutationResult = ApolloReactCommon.MutationResult<StudioDestroyMutation>;
export type StudioDestroyMutationOptions = ApolloReactCommon.BaseMutationOptions<StudioDestroyMutation, StudioDestroyMutationVariables>;
export const TagCreateDocument = gql`
    mutation TagCreate($name: String!) {
  tagCreate(input: {name: $name}) {
    ...TagData
  }
}
    ${TagDataFragmentDoc}`;
export type TagCreateMutationFn = ApolloReactCommon.MutationFunction<TagCreateMutation, TagCreateMutationVariables>;
export type TagCreateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<TagCreateMutation, TagCreateMutationVariables>, 'mutation'>;

    export const TagCreateComponent = (props: TagCreateComponentProps) => (
      <ApolloReactComponents.Mutation<TagCreateMutation, TagCreateMutationVariables> mutation={TagCreateDocument} {...props} />
    );
    

/**
 * __useTagCreateMutation__
 *
 * To run a mutation, you first call `useTagCreateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useTagCreateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [tagCreateMutation, { data, loading, error }] = useTagCreateMutation({
 *   variables: {
 *      name: // value for 'name'
 *   },
 * });
 */
export function useTagCreateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<TagCreateMutation, TagCreateMutationVariables>) {
        return ApolloReactHooks.useMutation<TagCreateMutation, TagCreateMutationVariables>(TagCreateDocument, baseOptions);
      }
export type TagCreateMutationHookResult = ReturnType<typeof useTagCreateMutation>;
export type TagCreateMutationResult = ApolloReactCommon.MutationResult<TagCreateMutation>;
export type TagCreateMutationOptions = ApolloReactCommon.BaseMutationOptions<TagCreateMutation, TagCreateMutationVariables>;
export const TagDestroyDocument = gql`
    mutation TagDestroy($id: ID!) {
  tagDestroy(input: {id: $id})
}
    `;
export type TagDestroyMutationFn = ApolloReactCommon.MutationFunction<TagDestroyMutation, TagDestroyMutationVariables>;
export type TagDestroyComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<TagDestroyMutation, TagDestroyMutationVariables>, 'mutation'>;

    export const TagDestroyComponent = (props: TagDestroyComponentProps) => (
      <ApolloReactComponents.Mutation<TagDestroyMutation, TagDestroyMutationVariables> mutation={TagDestroyDocument} {...props} />
    );
    

/**
 * __useTagDestroyMutation__
 *
 * To run a mutation, you first call `useTagDestroyMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useTagDestroyMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [tagDestroyMutation, { data, loading, error }] = useTagDestroyMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useTagDestroyMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<TagDestroyMutation, TagDestroyMutationVariables>) {
        return ApolloReactHooks.useMutation<TagDestroyMutation, TagDestroyMutationVariables>(TagDestroyDocument, baseOptions);
      }
export type TagDestroyMutationHookResult = ReturnType<typeof useTagDestroyMutation>;
export type TagDestroyMutationResult = ApolloReactCommon.MutationResult<TagDestroyMutation>;
export type TagDestroyMutationOptions = ApolloReactCommon.BaseMutationOptions<TagDestroyMutation, TagDestroyMutationVariables>;
export const TagUpdateDocument = gql`
    mutation TagUpdate($id: ID!, $name: String!) {
  tagUpdate(input: {id: $id, name: $name}) {
    ...TagData
  }
}
    ${TagDataFragmentDoc}`;
export type TagUpdateMutationFn = ApolloReactCommon.MutationFunction<TagUpdateMutation, TagUpdateMutationVariables>;
export type TagUpdateComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<TagUpdateMutation, TagUpdateMutationVariables>, 'mutation'>;

    export const TagUpdateComponent = (props: TagUpdateComponentProps) => (
      <ApolloReactComponents.Mutation<TagUpdateMutation, TagUpdateMutationVariables> mutation={TagUpdateDocument} {...props} />
    );
    

/**
 * __useTagUpdateMutation__
 *
 * To run a mutation, you first call `useTagUpdateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useTagUpdateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [tagUpdateMutation, { data, loading, error }] = useTagUpdateMutation({
 *   variables: {
 *      id: // value for 'id'
 *      name: // value for 'name'
 *   },
 * });
 */
export function useTagUpdateMutation(baseOptions?: ApolloReactHooks.MutationHookOptions<TagUpdateMutation, TagUpdateMutationVariables>) {
        return ApolloReactHooks.useMutation<TagUpdateMutation, TagUpdateMutationVariables>(TagUpdateDocument, baseOptions);
      }
export type TagUpdateMutationHookResult = ReturnType<typeof useTagUpdateMutation>;
export type TagUpdateMutationResult = ApolloReactCommon.MutationResult<TagUpdateMutation>;
export type TagUpdateMutationOptions = ApolloReactCommon.BaseMutationOptions<TagUpdateMutation, TagUpdateMutationVariables>;
export const FindGalleriesDocument = gql`
    query FindGalleries($filter: FindFilterType) {
  findGalleries(filter: $filter) {
    count
    galleries {
      ...GalleryData
    }
  }
}
    ${GalleryDataFragmentDoc}`;
export type FindGalleriesComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindGalleriesQuery, FindGalleriesQueryVariables>, 'query'>;

    export const FindGalleriesComponent = (props: FindGalleriesComponentProps) => (
      <ApolloReactComponents.Query<FindGalleriesQuery, FindGalleriesQueryVariables> query={FindGalleriesDocument} {...props} />
    );
    

/**
 * __useFindGalleriesQuery__
 *
 * To run a query within a React component, call `useFindGalleriesQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindGalleriesQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindGalleriesQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useFindGalleriesQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindGalleriesQuery, FindGalleriesQueryVariables>) {
        return ApolloReactHooks.useQuery<FindGalleriesQuery, FindGalleriesQueryVariables>(FindGalleriesDocument, baseOptions);
      }
export function useFindGalleriesLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindGalleriesQuery, FindGalleriesQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindGalleriesQuery, FindGalleriesQueryVariables>(FindGalleriesDocument, baseOptions);
        }
export type FindGalleriesQueryHookResult = ReturnType<typeof useFindGalleriesQuery>;
export type FindGalleriesLazyQueryHookResult = ReturnType<typeof useFindGalleriesLazyQuery>;
export type FindGalleriesQueryResult = ApolloReactCommon.QueryResult<FindGalleriesQuery, FindGalleriesQueryVariables>;
export const FindGalleryDocument = gql`
    query FindGallery($id: ID!) {
  findGallery(id: $id) {
    ...GalleryData
  }
}
    ${GalleryDataFragmentDoc}`;
export type FindGalleryComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindGalleryQuery, FindGalleryQueryVariables>, 'query'> & ({ variables: FindGalleryQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const FindGalleryComponent = (props: FindGalleryComponentProps) => (
      <ApolloReactComponents.Query<FindGalleryQuery, FindGalleryQueryVariables> query={FindGalleryDocument} {...props} />
    );
    

/**
 * __useFindGalleryQuery__
 *
 * To run a query within a React component, call `useFindGalleryQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindGalleryQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindGalleryQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useFindGalleryQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindGalleryQuery, FindGalleryQueryVariables>) {
        return ApolloReactHooks.useQuery<FindGalleryQuery, FindGalleryQueryVariables>(FindGalleryDocument, baseOptions);
      }
export function useFindGalleryLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindGalleryQuery, FindGalleryQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindGalleryQuery, FindGalleryQueryVariables>(FindGalleryDocument, baseOptions);
        }
export type FindGalleryQueryHookResult = ReturnType<typeof useFindGalleryQuery>;
export type FindGalleryLazyQueryHookResult = ReturnType<typeof useFindGalleryLazyQuery>;
export type FindGalleryQueryResult = ApolloReactCommon.QueryResult<FindGalleryQuery, FindGalleryQueryVariables>;
export const SceneWallDocument = gql`
    query SceneWall($q: String) {
  sceneWall(q: $q) {
    ...SceneData
  }
}
    ${SceneDataFragmentDoc}`;
export type SceneWallComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<SceneWallQuery, SceneWallQueryVariables>, 'query'>;

    export const SceneWallComponent = (props: SceneWallComponentProps) => (
      <ApolloReactComponents.Query<SceneWallQuery, SceneWallQueryVariables> query={SceneWallDocument} {...props} />
    );
    

/**
 * __useSceneWallQuery__
 *
 * To run a query within a React component, call `useSceneWallQuery` and pass it any options that fit your needs.
 * When your component renders, `useSceneWallQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSceneWallQuery({
 *   variables: {
 *      q: // value for 'q'
 *   },
 * });
 */
export function useSceneWallQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<SceneWallQuery, SceneWallQueryVariables>) {
        return ApolloReactHooks.useQuery<SceneWallQuery, SceneWallQueryVariables>(SceneWallDocument, baseOptions);
      }
export function useSceneWallLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<SceneWallQuery, SceneWallQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<SceneWallQuery, SceneWallQueryVariables>(SceneWallDocument, baseOptions);
        }
export type SceneWallQueryHookResult = ReturnType<typeof useSceneWallQuery>;
export type SceneWallLazyQueryHookResult = ReturnType<typeof useSceneWallLazyQuery>;
export type SceneWallQueryResult = ApolloReactCommon.QueryResult<SceneWallQuery, SceneWallQueryVariables>;
export const MarkerWallDocument = gql`
    query MarkerWall($q: String) {
  markerWall(q: $q) {
    ...SceneMarkerData
  }
}
    ${SceneMarkerDataFragmentDoc}`;
export type MarkerWallComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<MarkerWallQuery, MarkerWallQueryVariables>, 'query'>;

    export const MarkerWallComponent = (props: MarkerWallComponentProps) => (
      <ApolloReactComponents.Query<MarkerWallQuery, MarkerWallQueryVariables> query={MarkerWallDocument} {...props} />
    );
    

/**
 * __useMarkerWallQuery__
 *
 * To run a query within a React component, call `useMarkerWallQuery` and pass it any options that fit your needs.
 * When your component renders, `useMarkerWallQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMarkerWallQuery({
 *   variables: {
 *      q: // value for 'q'
 *   },
 * });
 */
export function useMarkerWallQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<MarkerWallQuery, MarkerWallQueryVariables>) {
        return ApolloReactHooks.useQuery<MarkerWallQuery, MarkerWallQueryVariables>(MarkerWallDocument, baseOptions);
      }
export function useMarkerWallLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<MarkerWallQuery, MarkerWallQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<MarkerWallQuery, MarkerWallQueryVariables>(MarkerWallDocument, baseOptions);
        }
export type MarkerWallQueryHookResult = ReturnType<typeof useMarkerWallQuery>;
export type MarkerWallLazyQueryHookResult = ReturnType<typeof useMarkerWallLazyQuery>;
export type MarkerWallQueryResult = ApolloReactCommon.QueryResult<MarkerWallQuery, MarkerWallQueryVariables>;
export const FindTagDocument = gql`
    query FindTag($id: ID!) {
  findTag(id: $id) {
    ...TagData
  }
}
    ${TagDataFragmentDoc}`;
export type FindTagComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindTagQuery, FindTagQueryVariables>, 'query'> & ({ variables: FindTagQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const FindTagComponent = (props: FindTagComponentProps) => (
      <ApolloReactComponents.Query<FindTagQuery, FindTagQueryVariables> query={FindTagDocument} {...props} />
    );
    

/**
 * __useFindTagQuery__
 *
 * To run a query within a React component, call `useFindTagQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindTagQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindTagQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useFindTagQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindTagQuery, FindTagQueryVariables>) {
        return ApolloReactHooks.useQuery<FindTagQuery, FindTagQueryVariables>(FindTagDocument, baseOptions);
      }
export function useFindTagLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindTagQuery, FindTagQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindTagQuery, FindTagQueryVariables>(FindTagDocument, baseOptions);
        }
export type FindTagQueryHookResult = ReturnType<typeof useFindTagQuery>;
export type FindTagLazyQueryHookResult = ReturnType<typeof useFindTagLazyQuery>;
export type FindTagQueryResult = ApolloReactCommon.QueryResult<FindTagQuery, FindTagQueryVariables>;
export const MarkerStringsDocument = gql`
    query MarkerStrings($q: String, $sort: String) {
  markerStrings(q: $q, sort: $sort) {
    id
    count
    title
  }
}
    `;
export type MarkerStringsComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<MarkerStringsQuery, MarkerStringsQueryVariables>, 'query'>;

    export const MarkerStringsComponent = (props: MarkerStringsComponentProps) => (
      <ApolloReactComponents.Query<MarkerStringsQuery, MarkerStringsQueryVariables> query={MarkerStringsDocument} {...props} />
    );
    

/**
 * __useMarkerStringsQuery__
 *
 * To run a query within a React component, call `useMarkerStringsQuery` and pass it any options that fit your needs.
 * When your component renders, `useMarkerStringsQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMarkerStringsQuery({
 *   variables: {
 *      q: // value for 'q'
 *      sort: // value for 'sort'
 *   },
 * });
 */
export function useMarkerStringsQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<MarkerStringsQuery, MarkerStringsQueryVariables>) {
        return ApolloReactHooks.useQuery<MarkerStringsQuery, MarkerStringsQueryVariables>(MarkerStringsDocument, baseOptions);
      }
export function useMarkerStringsLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<MarkerStringsQuery, MarkerStringsQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<MarkerStringsQuery, MarkerStringsQueryVariables>(MarkerStringsDocument, baseOptions);
        }
export type MarkerStringsQueryHookResult = ReturnType<typeof useMarkerStringsQuery>;
export type MarkerStringsLazyQueryHookResult = ReturnType<typeof useMarkerStringsLazyQuery>;
export type MarkerStringsQueryResult = ApolloReactCommon.QueryResult<MarkerStringsQuery, MarkerStringsQueryVariables>;
export const AllTagsDocument = gql`
    query AllTags {
  allTags {
    ...TagData
  }
}
    ${TagDataFragmentDoc}`;
export type AllTagsComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<AllTagsQuery, AllTagsQueryVariables>, 'query'>;

    export const AllTagsComponent = (props: AllTagsComponentProps) => (
      <ApolloReactComponents.Query<AllTagsQuery, AllTagsQueryVariables> query={AllTagsDocument} {...props} />
    );
    

/**
 * __useAllTagsQuery__
 *
 * To run a query within a React component, call `useAllTagsQuery` and pass it any options that fit your needs.
 * When your component renders, `useAllTagsQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useAllTagsQuery({
 *   variables: {
 *   },
 * });
 */
export function useAllTagsQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<AllTagsQuery, AllTagsQueryVariables>) {
        return ApolloReactHooks.useQuery<AllTagsQuery, AllTagsQueryVariables>(AllTagsDocument, baseOptions);
      }
export function useAllTagsLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<AllTagsQuery, AllTagsQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<AllTagsQuery, AllTagsQueryVariables>(AllTagsDocument, baseOptions);
        }
export type AllTagsQueryHookResult = ReturnType<typeof useAllTagsQuery>;
export type AllTagsLazyQueryHookResult = ReturnType<typeof useAllTagsLazyQuery>;
export type AllTagsQueryResult = ApolloReactCommon.QueryResult<AllTagsQuery, AllTagsQueryVariables>;
export const AllPerformersForFilterDocument = gql`
    query AllPerformersForFilter {
  allPerformers {
    ...SlimPerformerData
  }
}
    ${SlimPerformerDataFragmentDoc}`;
export type AllPerformersForFilterComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<AllPerformersForFilterQuery, AllPerformersForFilterQueryVariables>, 'query'>;

    export const AllPerformersForFilterComponent = (props: AllPerformersForFilterComponentProps) => (
      <ApolloReactComponents.Query<AllPerformersForFilterQuery, AllPerformersForFilterQueryVariables> query={AllPerformersForFilterDocument} {...props} />
    );
    

/**
 * __useAllPerformersForFilterQuery__
 *
 * To run a query within a React component, call `useAllPerformersForFilterQuery` and pass it any options that fit your needs.
 * When your component renders, `useAllPerformersForFilterQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useAllPerformersForFilterQuery({
 *   variables: {
 *   },
 * });
 */
export function useAllPerformersForFilterQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<AllPerformersForFilterQuery, AllPerformersForFilterQueryVariables>) {
        return ApolloReactHooks.useQuery<AllPerformersForFilterQuery, AllPerformersForFilterQueryVariables>(AllPerformersForFilterDocument, baseOptions);
      }
export function useAllPerformersForFilterLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<AllPerformersForFilterQuery, AllPerformersForFilterQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<AllPerformersForFilterQuery, AllPerformersForFilterQueryVariables>(AllPerformersForFilterDocument, baseOptions);
        }
export type AllPerformersForFilterQueryHookResult = ReturnType<typeof useAllPerformersForFilterQuery>;
export type AllPerformersForFilterLazyQueryHookResult = ReturnType<typeof useAllPerformersForFilterLazyQuery>;
export type AllPerformersForFilterQueryResult = ApolloReactCommon.QueryResult<AllPerformersForFilterQuery, AllPerformersForFilterQueryVariables>;
export const AllStudiosForFilterDocument = gql`
    query AllStudiosForFilter {
  allStudios {
    ...SlimStudioData
  }
}
    ${SlimStudioDataFragmentDoc}`;
export type AllStudiosForFilterComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<AllStudiosForFilterQuery, AllStudiosForFilterQueryVariables>, 'query'>;

    export const AllStudiosForFilterComponent = (props: AllStudiosForFilterComponentProps) => (
      <ApolloReactComponents.Query<AllStudiosForFilterQuery, AllStudiosForFilterQueryVariables> query={AllStudiosForFilterDocument} {...props} />
    );
    

/**
 * __useAllStudiosForFilterQuery__
 *
 * To run a query within a React component, call `useAllStudiosForFilterQuery` and pass it any options that fit your needs.
 * When your component renders, `useAllStudiosForFilterQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useAllStudiosForFilterQuery({
 *   variables: {
 *   },
 * });
 */
export function useAllStudiosForFilterQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<AllStudiosForFilterQuery, AllStudiosForFilterQueryVariables>) {
        return ApolloReactHooks.useQuery<AllStudiosForFilterQuery, AllStudiosForFilterQueryVariables>(AllStudiosForFilterDocument, baseOptions);
      }
export function useAllStudiosForFilterLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<AllStudiosForFilterQuery, AllStudiosForFilterQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<AllStudiosForFilterQuery, AllStudiosForFilterQueryVariables>(AllStudiosForFilterDocument, baseOptions);
        }
export type AllStudiosForFilterQueryHookResult = ReturnType<typeof useAllStudiosForFilterQuery>;
export type AllStudiosForFilterLazyQueryHookResult = ReturnType<typeof useAllStudiosForFilterLazyQuery>;
export type AllStudiosForFilterQueryResult = ApolloReactCommon.QueryResult<AllStudiosForFilterQuery, AllStudiosForFilterQueryVariables>;
export const AllTagsForFilterDocument = gql`
    query AllTagsForFilter {
  allTags {
    id
    name
  }
}
    `;
export type AllTagsForFilterComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<AllTagsForFilterQuery, AllTagsForFilterQueryVariables>, 'query'>;

    export const AllTagsForFilterComponent = (props: AllTagsForFilterComponentProps) => (
      <ApolloReactComponents.Query<AllTagsForFilterQuery, AllTagsForFilterQueryVariables> query={AllTagsForFilterDocument} {...props} />
    );
    

/**
 * __useAllTagsForFilterQuery__
 *
 * To run a query within a React component, call `useAllTagsForFilterQuery` and pass it any options that fit your needs.
 * When your component renders, `useAllTagsForFilterQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useAllTagsForFilterQuery({
 *   variables: {
 *   },
 * });
 */
export function useAllTagsForFilterQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<AllTagsForFilterQuery, AllTagsForFilterQueryVariables>) {
        return ApolloReactHooks.useQuery<AllTagsForFilterQuery, AllTagsForFilterQueryVariables>(AllTagsForFilterDocument, baseOptions);
      }
export function useAllTagsForFilterLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<AllTagsForFilterQuery, AllTagsForFilterQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<AllTagsForFilterQuery, AllTagsForFilterQueryVariables>(AllTagsForFilterDocument, baseOptions);
        }
export type AllTagsForFilterQueryHookResult = ReturnType<typeof useAllTagsForFilterQuery>;
export type AllTagsForFilterLazyQueryHookResult = ReturnType<typeof useAllTagsForFilterLazyQuery>;
export type AllTagsForFilterQueryResult = ApolloReactCommon.QueryResult<AllTagsForFilterQuery, AllTagsForFilterQueryVariables>;
export const ValidGalleriesForSceneDocument = gql`
    query ValidGalleriesForScene($scene_id: ID!) {
  validGalleriesForScene(scene_id: $scene_id) {
    id
    path
  }
}
    `;
export type ValidGalleriesForSceneComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ValidGalleriesForSceneQuery, ValidGalleriesForSceneQueryVariables>, 'query'> & ({ variables: ValidGalleriesForSceneQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const ValidGalleriesForSceneComponent = (props: ValidGalleriesForSceneComponentProps) => (
      <ApolloReactComponents.Query<ValidGalleriesForSceneQuery, ValidGalleriesForSceneQueryVariables> query={ValidGalleriesForSceneDocument} {...props} />
    );
    

/**
 * __useValidGalleriesForSceneQuery__
 *
 * To run a query within a React component, call `useValidGalleriesForSceneQuery` and pass it any options that fit your needs.
 * When your component renders, `useValidGalleriesForSceneQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useValidGalleriesForSceneQuery({
 *   variables: {
 *      scene_id: // value for 'scene_id'
 *   },
 * });
 */
export function useValidGalleriesForSceneQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ValidGalleriesForSceneQuery, ValidGalleriesForSceneQueryVariables>) {
        return ApolloReactHooks.useQuery<ValidGalleriesForSceneQuery, ValidGalleriesForSceneQueryVariables>(ValidGalleriesForSceneDocument, baseOptions);
      }
export function useValidGalleriesForSceneLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ValidGalleriesForSceneQuery, ValidGalleriesForSceneQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ValidGalleriesForSceneQuery, ValidGalleriesForSceneQueryVariables>(ValidGalleriesForSceneDocument, baseOptions);
        }
export type ValidGalleriesForSceneQueryHookResult = ReturnType<typeof useValidGalleriesForSceneQuery>;
export type ValidGalleriesForSceneLazyQueryHookResult = ReturnType<typeof useValidGalleriesForSceneLazyQuery>;
export type ValidGalleriesForSceneQueryResult = ApolloReactCommon.QueryResult<ValidGalleriesForSceneQuery, ValidGalleriesForSceneQueryVariables>;
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
export type StatsComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<StatsQuery, StatsQueryVariables>, 'query'>;

    export const StatsComponent = (props: StatsComponentProps) => (
      <ApolloReactComponents.Query<StatsQuery, StatsQueryVariables> query={StatsDocument} {...props} />
    );
    

/**
 * __useStatsQuery__
 *
 * To run a query within a React component, call `useStatsQuery` and pass it any options that fit your needs.
 * When your component renders, `useStatsQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useStatsQuery({
 *   variables: {
 *   },
 * });
 */
export function useStatsQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<StatsQuery, StatsQueryVariables>) {
        return ApolloReactHooks.useQuery<StatsQuery, StatsQueryVariables>(StatsDocument, baseOptions);
      }
export function useStatsLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<StatsQuery, StatsQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<StatsQuery, StatsQueryVariables>(StatsDocument, baseOptions);
        }
export type StatsQueryHookResult = ReturnType<typeof useStatsQuery>;
export type StatsLazyQueryHookResult = ReturnType<typeof useStatsLazyQuery>;
export type StatsQueryResult = ApolloReactCommon.QueryResult<StatsQuery, StatsQueryVariables>;
export const LogsDocument = gql`
    query Logs {
  logs {
    ...LogEntryData
  }
}
    ${LogEntryDataFragmentDoc}`;
export type LogsComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<LogsQuery, LogsQueryVariables>, 'query'>;

    export const LogsComponent = (props: LogsComponentProps) => (
      <ApolloReactComponents.Query<LogsQuery, LogsQueryVariables> query={LogsDocument} {...props} />
    );
    

/**
 * __useLogsQuery__
 *
 * To run a query within a React component, call `useLogsQuery` and pass it any options that fit your needs.
 * When your component renders, `useLogsQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useLogsQuery({
 *   variables: {
 *   },
 * });
 */
export function useLogsQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<LogsQuery, LogsQueryVariables>) {
        return ApolloReactHooks.useQuery<LogsQuery, LogsQueryVariables>(LogsDocument, baseOptions);
      }
export function useLogsLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<LogsQuery, LogsQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<LogsQuery, LogsQueryVariables>(LogsDocument, baseOptions);
        }
export type LogsQueryHookResult = ReturnType<typeof useLogsQuery>;
export type LogsLazyQueryHookResult = ReturnType<typeof useLogsLazyQuery>;
export type LogsQueryResult = ApolloReactCommon.QueryResult<LogsQuery, LogsQueryVariables>;
export const VersionDocument = gql`
    query Version {
  version {
    version
    hash
    build_time
  }
}
    `;
export type VersionComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<VersionQuery, VersionQueryVariables>, 'query'>;

    export const VersionComponent = (props: VersionComponentProps) => (
      <ApolloReactComponents.Query<VersionQuery, VersionQueryVariables> query={VersionDocument} {...props} />
    );
    

/**
 * __useVersionQuery__
 *
 * To run a query within a React component, call `useVersionQuery` and pass it any options that fit your needs.
 * When your component renders, `useVersionQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useVersionQuery({
 *   variables: {
 *   },
 * });
 */
export function useVersionQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<VersionQuery, VersionQueryVariables>) {
        return ApolloReactHooks.useQuery<VersionQuery, VersionQueryVariables>(VersionDocument, baseOptions);
      }
export function useVersionLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<VersionQuery, VersionQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<VersionQuery, VersionQueryVariables>(VersionDocument, baseOptions);
        }
export type VersionQueryHookResult = ReturnType<typeof useVersionQuery>;
export type VersionLazyQueryHookResult = ReturnType<typeof useVersionLazyQuery>;
export type VersionQueryResult = ApolloReactCommon.QueryResult<VersionQuery, VersionQueryVariables>;
export const LatestVersionDocument = gql`
    query LatestVersion {
  latestversion {
    shorthash
    url
  }
}
    `;
export type LatestVersionComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<LatestVersionQuery, LatestVersionQueryVariables>, 'query'>;

    export const LatestVersionComponent = (props: LatestVersionComponentProps) => (
      <ApolloReactComponents.Query<LatestVersionQuery, LatestVersionQueryVariables> query={LatestVersionDocument} {...props} />
    );
    

/**
 * __useLatestVersionQuery__
 *
 * To run a query within a React component, call `useLatestVersionQuery` and pass it any options that fit your needs.
 * When your component renders, `useLatestVersionQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useLatestVersionQuery({
 *   variables: {
 *   },
 * });
 */
export function useLatestVersionQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<LatestVersionQuery, LatestVersionQueryVariables>) {
        return ApolloReactHooks.useQuery<LatestVersionQuery, LatestVersionQueryVariables>(LatestVersionDocument, baseOptions);
      }
export function useLatestVersionLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<LatestVersionQuery, LatestVersionQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<LatestVersionQuery, LatestVersionQueryVariables>(LatestVersionDocument, baseOptions);
        }
export type LatestVersionQueryHookResult = ReturnType<typeof useLatestVersionQuery>;
export type LatestVersionLazyQueryHookResult = ReturnType<typeof useLatestVersionLazyQuery>;
export type LatestVersionQueryResult = ApolloReactCommon.QueryResult<LatestVersionQuery, LatestVersionQueryVariables>;
export const FindPerformersDocument = gql`
    query FindPerformers($filter: FindFilterType, $performer_filter: PerformerFilterType) {
  findPerformers(filter: $filter, performer_filter: $performer_filter) {
    count
    performers {
      ...PerformerData
    }
  }
}
    ${PerformerDataFragmentDoc}`;
export type FindPerformersComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindPerformersQuery, FindPerformersQueryVariables>, 'query'>;

    export const FindPerformersComponent = (props: FindPerformersComponentProps) => (
      <ApolloReactComponents.Query<FindPerformersQuery, FindPerformersQueryVariables> query={FindPerformersDocument} {...props} />
    );
    

/**
 * __useFindPerformersQuery__
 *
 * To run a query within a React component, call `useFindPerformersQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindPerformersQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindPerformersQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *      performer_filter: // value for 'performer_filter'
 *   },
 * });
 */
export function useFindPerformersQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindPerformersQuery, FindPerformersQueryVariables>) {
        return ApolloReactHooks.useQuery<FindPerformersQuery, FindPerformersQueryVariables>(FindPerformersDocument, baseOptions);
      }
export function useFindPerformersLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindPerformersQuery, FindPerformersQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindPerformersQuery, FindPerformersQueryVariables>(FindPerformersDocument, baseOptions);
        }
export type FindPerformersQueryHookResult = ReturnType<typeof useFindPerformersQuery>;
export type FindPerformersLazyQueryHookResult = ReturnType<typeof useFindPerformersLazyQuery>;
export type FindPerformersQueryResult = ApolloReactCommon.QueryResult<FindPerformersQuery, FindPerformersQueryVariables>;
export const FindPerformerDocument = gql`
    query FindPerformer($id: ID!) {
  findPerformer(id: $id) {
    ...PerformerData
  }
}
    ${PerformerDataFragmentDoc}`;
export type FindPerformerComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindPerformerQuery, FindPerformerQueryVariables>, 'query'> & ({ variables: FindPerformerQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const FindPerformerComponent = (props: FindPerformerComponentProps) => (
      <ApolloReactComponents.Query<FindPerformerQuery, FindPerformerQueryVariables> query={FindPerformerDocument} {...props} />
    );
    

/**
 * __useFindPerformerQuery__
 *
 * To run a query within a React component, call `useFindPerformerQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindPerformerQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindPerformerQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useFindPerformerQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindPerformerQuery, FindPerformerQueryVariables>) {
        return ApolloReactHooks.useQuery<FindPerformerQuery, FindPerformerQueryVariables>(FindPerformerDocument, baseOptions);
      }
export function useFindPerformerLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindPerformerQuery, FindPerformerQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindPerformerQuery, FindPerformerQueryVariables>(FindPerformerDocument, baseOptions);
        }
export type FindPerformerQueryHookResult = ReturnType<typeof useFindPerformerQuery>;
export type FindPerformerLazyQueryHookResult = ReturnType<typeof useFindPerformerLazyQuery>;
export type FindPerformerQueryResult = ApolloReactCommon.QueryResult<FindPerformerQuery, FindPerformerQueryVariables>;
export const FindSceneMarkersDocument = gql`
    query FindSceneMarkers($filter: FindFilterType, $scene_marker_filter: SceneMarkerFilterType) {
  findSceneMarkers(filter: $filter, scene_marker_filter: $scene_marker_filter) {
    count
    scene_markers {
      ...SceneMarkerData
    }
  }
}
    ${SceneMarkerDataFragmentDoc}`;
export type FindSceneMarkersComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindSceneMarkersQuery, FindSceneMarkersQueryVariables>, 'query'>;

    export const FindSceneMarkersComponent = (props: FindSceneMarkersComponentProps) => (
      <ApolloReactComponents.Query<FindSceneMarkersQuery, FindSceneMarkersQueryVariables> query={FindSceneMarkersDocument} {...props} />
    );
    

/**
 * __useFindSceneMarkersQuery__
 *
 * To run a query within a React component, call `useFindSceneMarkersQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindSceneMarkersQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindSceneMarkersQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *      scene_marker_filter: // value for 'scene_marker_filter'
 *   },
 * });
 */
export function useFindSceneMarkersQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindSceneMarkersQuery, FindSceneMarkersQueryVariables>) {
        return ApolloReactHooks.useQuery<FindSceneMarkersQuery, FindSceneMarkersQueryVariables>(FindSceneMarkersDocument, baseOptions);
      }
export function useFindSceneMarkersLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindSceneMarkersQuery, FindSceneMarkersQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindSceneMarkersQuery, FindSceneMarkersQueryVariables>(FindSceneMarkersDocument, baseOptions);
        }
export type FindSceneMarkersQueryHookResult = ReturnType<typeof useFindSceneMarkersQuery>;
export type FindSceneMarkersLazyQueryHookResult = ReturnType<typeof useFindSceneMarkersLazyQuery>;
export type FindSceneMarkersQueryResult = ApolloReactCommon.QueryResult<FindSceneMarkersQuery, FindSceneMarkersQueryVariables>;
export const FindScenesDocument = gql`
    query FindScenes($filter: FindFilterType, $scene_filter: SceneFilterType, $scene_ids: [Int!]) {
  findScenes(filter: $filter, scene_filter: $scene_filter, scene_ids: $scene_ids) {
    count
    scenes {
      ...SlimSceneData
    }
  }
}
    ${SlimSceneDataFragmentDoc}`;
export type FindScenesComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindScenesQuery, FindScenesQueryVariables>, 'query'>;

    export const FindScenesComponent = (props: FindScenesComponentProps) => (
      <ApolloReactComponents.Query<FindScenesQuery, FindScenesQueryVariables> query={FindScenesDocument} {...props} />
    );
    

/**
 * __useFindScenesQuery__
 *
 * To run a query within a React component, call `useFindScenesQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindScenesQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindScenesQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *      scene_filter: // value for 'scene_filter'
 *      scene_ids: // value for 'scene_ids'
 *   },
 * });
 */
export function useFindScenesQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindScenesQuery, FindScenesQueryVariables>) {
        return ApolloReactHooks.useQuery<FindScenesQuery, FindScenesQueryVariables>(FindScenesDocument, baseOptions);
      }
export function useFindScenesLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindScenesQuery, FindScenesQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindScenesQuery, FindScenesQueryVariables>(FindScenesDocument, baseOptions);
        }
export type FindScenesQueryHookResult = ReturnType<typeof useFindScenesQuery>;
export type FindScenesLazyQueryHookResult = ReturnType<typeof useFindScenesLazyQuery>;
export type FindScenesQueryResult = ApolloReactCommon.QueryResult<FindScenesQuery, FindScenesQueryVariables>;
export const FindScenesByPathRegexDocument = gql`
    query FindScenesByPathRegex($filter: FindFilterType) {
  findScenesByPathRegex(filter: $filter) {
    count
    scenes {
      ...SlimSceneData
    }
  }
}
    ${SlimSceneDataFragmentDoc}`;
export type FindScenesByPathRegexComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindScenesByPathRegexQuery, FindScenesByPathRegexQueryVariables>, 'query'>;

    export const FindScenesByPathRegexComponent = (props: FindScenesByPathRegexComponentProps) => (
      <ApolloReactComponents.Query<FindScenesByPathRegexQuery, FindScenesByPathRegexQueryVariables> query={FindScenesByPathRegexDocument} {...props} />
    );
    

/**
 * __useFindScenesByPathRegexQuery__
 *
 * To run a query within a React component, call `useFindScenesByPathRegexQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindScenesByPathRegexQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindScenesByPathRegexQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useFindScenesByPathRegexQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindScenesByPathRegexQuery, FindScenesByPathRegexQueryVariables>) {
        return ApolloReactHooks.useQuery<FindScenesByPathRegexQuery, FindScenesByPathRegexQueryVariables>(FindScenesByPathRegexDocument, baseOptions);
      }
export function useFindScenesByPathRegexLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindScenesByPathRegexQuery, FindScenesByPathRegexQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindScenesByPathRegexQuery, FindScenesByPathRegexQueryVariables>(FindScenesByPathRegexDocument, baseOptions);
        }
export type FindScenesByPathRegexQueryHookResult = ReturnType<typeof useFindScenesByPathRegexQuery>;
export type FindScenesByPathRegexLazyQueryHookResult = ReturnType<typeof useFindScenesByPathRegexLazyQuery>;
export type FindScenesByPathRegexQueryResult = ApolloReactCommon.QueryResult<FindScenesByPathRegexQuery, FindScenesByPathRegexQueryVariables>;
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
${SceneMarkerDataFragmentDoc}`;
export type FindSceneComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindSceneQuery, FindSceneQueryVariables>, 'query'> & ({ variables: FindSceneQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const FindSceneComponent = (props: FindSceneComponentProps) => (
      <ApolloReactComponents.Query<FindSceneQuery, FindSceneQueryVariables> query={FindSceneDocument} {...props} />
    );
    

/**
 * __useFindSceneQuery__
 *
 * To run a query within a React component, call `useFindSceneQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindSceneQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindSceneQuery({
 *   variables: {
 *      id: // value for 'id'
 *      checksum: // value for 'checksum'
 *   },
 * });
 */
export function useFindSceneQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindSceneQuery, FindSceneQueryVariables>) {
        return ApolloReactHooks.useQuery<FindSceneQuery, FindSceneQueryVariables>(FindSceneDocument, baseOptions);
      }
export function useFindSceneLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindSceneQuery, FindSceneQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindSceneQuery, FindSceneQueryVariables>(FindSceneDocument, baseOptions);
        }
export type FindSceneQueryHookResult = ReturnType<typeof useFindSceneQuery>;
export type FindSceneLazyQueryHookResult = ReturnType<typeof useFindSceneLazyQuery>;
export type FindSceneQueryResult = ApolloReactCommon.QueryResult<FindSceneQuery, FindSceneQueryVariables>;
export const ParseSceneFilenamesDocument = gql`
    query ParseSceneFilenames($filter: FindFilterType!, $config: SceneParserInput!) {
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
    ${SlimSceneDataFragmentDoc}`;
export type ParseSceneFilenamesComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ParseSceneFilenamesQuery, ParseSceneFilenamesQueryVariables>, 'query'> & ({ variables: ParseSceneFilenamesQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const ParseSceneFilenamesComponent = (props: ParseSceneFilenamesComponentProps) => (
      <ApolloReactComponents.Query<ParseSceneFilenamesQuery, ParseSceneFilenamesQueryVariables> query={ParseSceneFilenamesDocument} {...props} />
    );
    

/**
 * __useParseSceneFilenamesQuery__
 *
 * To run a query within a React component, call `useParseSceneFilenamesQuery` and pass it any options that fit your needs.
 * When your component renders, `useParseSceneFilenamesQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useParseSceneFilenamesQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *      config: // value for 'config'
 *   },
 * });
 */
export function useParseSceneFilenamesQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ParseSceneFilenamesQuery, ParseSceneFilenamesQueryVariables>) {
        return ApolloReactHooks.useQuery<ParseSceneFilenamesQuery, ParseSceneFilenamesQueryVariables>(ParseSceneFilenamesDocument, baseOptions);
      }
export function useParseSceneFilenamesLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ParseSceneFilenamesQuery, ParseSceneFilenamesQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ParseSceneFilenamesQuery, ParseSceneFilenamesQueryVariables>(ParseSceneFilenamesDocument, baseOptions);
        }
export type ParseSceneFilenamesQueryHookResult = ReturnType<typeof useParseSceneFilenamesQuery>;
export type ParseSceneFilenamesLazyQueryHookResult = ReturnType<typeof useParseSceneFilenamesLazyQuery>;
export type ParseSceneFilenamesQueryResult = ApolloReactCommon.QueryResult<ParseSceneFilenamesQuery, ParseSceneFilenamesQueryVariables>;
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
export type ScrapeFreeonesComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ScrapeFreeonesQuery, ScrapeFreeonesQueryVariables>, 'query'> & ({ variables: ScrapeFreeonesQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const ScrapeFreeonesComponent = (props: ScrapeFreeonesComponentProps) => (
      <ApolloReactComponents.Query<ScrapeFreeonesQuery, ScrapeFreeonesQueryVariables> query={ScrapeFreeonesDocument} {...props} />
    );
    

/**
 * __useScrapeFreeonesQuery__
 *
 * To run a query within a React component, call `useScrapeFreeonesQuery` and pass it any options that fit your needs.
 * When your component renders, `useScrapeFreeonesQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useScrapeFreeonesQuery({
 *   variables: {
 *      performer_name: // value for 'performer_name'
 *   },
 * });
 */
export function useScrapeFreeonesQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ScrapeFreeonesQuery, ScrapeFreeonesQueryVariables>) {
        return ApolloReactHooks.useQuery<ScrapeFreeonesQuery, ScrapeFreeonesQueryVariables>(ScrapeFreeonesDocument, baseOptions);
      }
export function useScrapeFreeonesLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ScrapeFreeonesQuery, ScrapeFreeonesQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ScrapeFreeonesQuery, ScrapeFreeonesQueryVariables>(ScrapeFreeonesDocument, baseOptions);
        }
export type ScrapeFreeonesQueryHookResult = ReturnType<typeof useScrapeFreeonesQuery>;
export type ScrapeFreeonesLazyQueryHookResult = ReturnType<typeof useScrapeFreeonesLazyQuery>;
export type ScrapeFreeonesQueryResult = ApolloReactCommon.QueryResult<ScrapeFreeonesQuery, ScrapeFreeonesQueryVariables>;
export const ScrapeFreeonesPerformersDocument = gql`
    query ScrapeFreeonesPerformers($q: String!) {
  scrapeFreeonesPerformerList(query: $q)
}
    `;
export type ScrapeFreeonesPerformersComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ScrapeFreeonesPerformersQuery, ScrapeFreeonesPerformersQueryVariables>, 'query'> & ({ variables: ScrapeFreeonesPerformersQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const ScrapeFreeonesPerformersComponent = (props: ScrapeFreeonesPerformersComponentProps) => (
      <ApolloReactComponents.Query<ScrapeFreeonesPerformersQuery, ScrapeFreeonesPerformersQueryVariables> query={ScrapeFreeonesPerformersDocument} {...props} />
    );
    

/**
 * __useScrapeFreeonesPerformersQuery__
 *
 * To run a query within a React component, call `useScrapeFreeonesPerformersQuery` and pass it any options that fit your needs.
 * When your component renders, `useScrapeFreeonesPerformersQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useScrapeFreeonesPerformersQuery({
 *   variables: {
 *      q: // value for 'q'
 *   },
 * });
 */
export function useScrapeFreeonesPerformersQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ScrapeFreeonesPerformersQuery, ScrapeFreeonesPerformersQueryVariables>) {
        return ApolloReactHooks.useQuery<ScrapeFreeonesPerformersQuery, ScrapeFreeonesPerformersQueryVariables>(ScrapeFreeonesPerformersDocument, baseOptions);
      }
export function useScrapeFreeonesPerformersLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ScrapeFreeonesPerformersQuery, ScrapeFreeonesPerformersQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ScrapeFreeonesPerformersQuery, ScrapeFreeonesPerformersQueryVariables>(ScrapeFreeonesPerformersDocument, baseOptions);
        }
export type ScrapeFreeonesPerformersQueryHookResult = ReturnType<typeof useScrapeFreeonesPerformersQuery>;
export type ScrapeFreeonesPerformersLazyQueryHookResult = ReturnType<typeof useScrapeFreeonesPerformersLazyQuery>;
export type ScrapeFreeonesPerformersQueryResult = ApolloReactCommon.QueryResult<ScrapeFreeonesPerformersQuery, ScrapeFreeonesPerformersQueryVariables>;
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
export type ListPerformerScrapersComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ListPerformerScrapersQuery, ListPerformerScrapersQueryVariables>, 'query'>;

    export const ListPerformerScrapersComponent = (props: ListPerformerScrapersComponentProps) => (
      <ApolloReactComponents.Query<ListPerformerScrapersQuery, ListPerformerScrapersQueryVariables> query={ListPerformerScrapersDocument} {...props} />
    );
    

/**
 * __useListPerformerScrapersQuery__
 *
 * To run a query within a React component, call `useListPerformerScrapersQuery` and pass it any options that fit your needs.
 * When your component renders, `useListPerformerScrapersQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListPerformerScrapersQuery({
 *   variables: {
 *   },
 * });
 */
export function useListPerformerScrapersQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ListPerformerScrapersQuery, ListPerformerScrapersQueryVariables>) {
        return ApolloReactHooks.useQuery<ListPerformerScrapersQuery, ListPerformerScrapersQueryVariables>(ListPerformerScrapersDocument, baseOptions);
      }
export function useListPerformerScrapersLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ListPerformerScrapersQuery, ListPerformerScrapersQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ListPerformerScrapersQuery, ListPerformerScrapersQueryVariables>(ListPerformerScrapersDocument, baseOptions);
        }
export type ListPerformerScrapersQueryHookResult = ReturnType<typeof useListPerformerScrapersQuery>;
export type ListPerformerScrapersLazyQueryHookResult = ReturnType<typeof useListPerformerScrapersLazyQuery>;
export type ListPerformerScrapersQueryResult = ApolloReactCommon.QueryResult<ListPerformerScrapersQuery, ListPerformerScrapersQueryVariables>;
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
export type ListSceneScrapersComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ListSceneScrapersQuery, ListSceneScrapersQueryVariables>, 'query'>;

    export const ListSceneScrapersComponent = (props: ListSceneScrapersComponentProps) => (
      <ApolloReactComponents.Query<ListSceneScrapersQuery, ListSceneScrapersQueryVariables> query={ListSceneScrapersDocument} {...props} />
    );
    

/**
 * __useListSceneScrapersQuery__
 *
 * To run a query within a React component, call `useListSceneScrapersQuery` and pass it any options that fit your needs.
 * When your component renders, `useListSceneScrapersQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListSceneScrapersQuery({
 *   variables: {
 *   },
 * });
 */
export function useListSceneScrapersQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ListSceneScrapersQuery, ListSceneScrapersQueryVariables>) {
        return ApolloReactHooks.useQuery<ListSceneScrapersQuery, ListSceneScrapersQueryVariables>(ListSceneScrapersDocument, baseOptions);
      }
export function useListSceneScrapersLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ListSceneScrapersQuery, ListSceneScrapersQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ListSceneScrapersQuery, ListSceneScrapersQueryVariables>(ListSceneScrapersDocument, baseOptions);
        }
export type ListSceneScrapersQueryHookResult = ReturnType<typeof useListSceneScrapersQuery>;
export type ListSceneScrapersLazyQueryHookResult = ReturnType<typeof useListSceneScrapersLazyQuery>;
export type ListSceneScrapersQueryResult = ApolloReactCommon.QueryResult<ListSceneScrapersQuery, ListSceneScrapersQueryVariables>;
export const ScrapePerformerListDocument = gql`
    query ScrapePerformerList($scraper_id: ID!, $query: String!) {
  scrapePerformerList(scraper_id: $scraper_id, query: $query) {
    ...ScrapedPerformerData
  }
}
    ${ScrapedPerformerDataFragmentDoc}`;
export type ScrapePerformerListComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ScrapePerformerListQuery, ScrapePerformerListQueryVariables>, 'query'> & ({ variables: ScrapePerformerListQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const ScrapePerformerListComponent = (props: ScrapePerformerListComponentProps) => (
      <ApolloReactComponents.Query<ScrapePerformerListQuery, ScrapePerformerListQueryVariables> query={ScrapePerformerListDocument} {...props} />
    );
    

/**
 * __useScrapePerformerListQuery__
 *
 * To run a query within a React component, call `useScrapePerformerListQuery` and pass it any options that fit your needs.
 * When your component renders, `useScrapePerformerListQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useScrapePerformerListQuery({
 *   variables: {
 *      scraper_id: // value for 'scraper_id'
 *      query: // value for 'query'
 *   },
 * });
 */
export function useScrapePerformerListQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ScrapePerformerListQuery, ScrapePerformerListQueryVariables>) {
        return ApolloReactHooks.useQuery<ScrapePerformerListQuery, ScrapePerformerListQueryVariables>(ScrapePerformerListDocument, baseOptions);
      }
export function useScrapePerformerListLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ScrapePerformerListQuery, ScrapePerformerListQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ScrapePerformerListQuery, ScrapePerformerListQueryVariables>(ScrapePerformerListDocument, baseOptions);
        }
export type ScrapePerformerListQueryHookResult = ReturnType<typeof useScrapePerformerListQuery>;
export type ScrapePerformerListLazyQueryHookResult = ReturnType<typeof useScrapePerformerListLazyQuery>;
export type ScrapePerformerListQueryResult = ApolloReactCommon.QueryResult<ScrapePerformerListQuery, ScrapePerformerListQueryVariables>;
export const ScrapePerformerDocument = gql`
    query ScrapePerformer($scraper_id: ID!, $scraped_performer: ScrapedPerformerInput!) {
  scrapePerformer(scraper_id: $scraper_id, scraped_performer: $scraped_performer) {
    ...ScrapedPerformerData
  }
}
    ${ScrapedPerformerDataFragmentDoc}`;
export type ScrapePerformerComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ScrapePerformerQuery, ScrapePerformerQueryVariables>, 'query'> & ({ variables: ScrapePerformerQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const ScrapePerformerComponent = (props: ScrapePerformerComponentProps) => (
      <ApolloReactComponents.Query<ScrapePerformerQuery, ScrapePerformerQueryVariables> query={ScrapePerformerDocument} {...props} />
    );
    

/**
 * __useScrapePerformerQuery__
 *
 * To run a query within a React component, call `useScrapePerformerQuery` and pass it any options that fit your needs.
 * When your component renders, `useScrapePerformerQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useScrapePerformerQuery({
 *   variables: {
 *      scraper_id: // value for 'scraper_id'
 *      scraped_performer: // value for 'scraped_performer'
 *   },
 * });
 */
export function useScrapePerformerQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ScrapePerformerQuery, ScrapePerformerQueryVariables>) {
        return ApolloReactHooks.useQuery<ScrapePerformerQuery, ScrapePerformerQueryVariables>(ScrapePerformerDocument, baseOptions);
      }
export function useScrapePerformerLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ScrapePerformerQuery, ScrapePerformerQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ScrapePerformerQuery, ScrapePerformerQueryVariables>(ScrapePerformerDocument, baseOptions);
        }
export type ScrapePerformerQueryHookResult = ReturnType<typeof useScrapePerformerQuery>;
export type ScrapePerformerLazyQueryHookResult = ReturnType<typeof useScrapePerformerLazyQuery>;
export type ScrapePerformerQueryResult = ApolloReactCommon.QueryResult<ScrapePerformerQuery, ScrapePerformerQueryVariables>;
export const ScrapePerformerUrlDocument = gql`
    query ScrapePerformerURL($url: String!) {
  scrapePerformerURL(url: $url) {
    ...ScrapedPerformerData
  }
}
    ${ScrapedPerformerDataFragmentDoc}`;
export type ScrapePerformerUrlComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ScrapePerformerUrlQuery, ScrapePerformerUrlQueryVariables>, 'query'> & ({ variables: ScrapePerformerUrlQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const ScrapePerformerUrlComponent = (props: ScrapePerformerUrlComponentProps) => (
      <ApolloReactComponents.Query<ScrapePerformerUrlQuery, ScrapePerformerUrlQueryVariables> query={ScrapePerformerUrlDocument} {...props} />
    );
    

/**
 * __useScrapePerformerUrlQuery__
 *
 * To run a query within a React component, call `useScrapePerformerUrlQuery` and pass it any options that fit your needs.
 * When your component renders, `useScrapePerformerUrlQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useScrapePerformerUrlQuery({
 *   variables: {
 *      url: // value for 'url'
 *   },
 * });
 */
export function useScrapePerformerUrlQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ScrapePerformerUrlQuery, ScrapePerformerUrlQueryVariables>) {
        return ApolloReactHooks.useQuery<ScrapePerformerUrlQuery, ScrapePerformerUrlQueryVariables>(ScrapePerformerUrlDocument, baseOptions);
      }
export function useScrapePerformerUrlLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ScrapePerformerUrlQuery, ScrapePerformerUrlQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ScrapePerformerUrlQuery, ScrapePerformerUrlQueryVariables>(ScrapePerformerUrlDocument, baseOptions);
        }
export type ScrapePerformerUrlQueryHookResult = ReturnType<typeof useScrapePerformerUrlQuery>;
export type ScrapePerformerUrlLazyQueryHookResult = ReturnType<typeof useScrapePerformerUrlLazyQuery>;
export type ScrapePerformerUrlQueryResult = ApolloReactCommon.QueryResult<ScrapePerformerUrlQuery, ScrapePerformerUrlQueryVariables>;
export const ScrapeSceneDocument = gql`
    query ScrapeScene($scraper_id: ID!, $scene: SceneUpdateInput!) {
  scrapeScene(scraper_id: $scraper_id, scene: $scene) {
    ...ScrapedSceneData
  }
}
    ${ScrapedSceneDataFragmentDoc}`;
export type ScrapeSceneComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ScrapeSceneQuery, ScrapeSceneQueryVariables>, 'query'> & ({ variables: ScrapeSceneQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const ScrapeSceneComponent = (props: ScrapeSceneComponentProps) => (
      <ApolloReactComponents.Query<ScrapeSceneQuery, ScrapeSceneQueryVariables> query={ScrapeSceneDocument} {...props} />
    );
    

/**
 * __useScrapeSceneQuery__
 *
 * To run a query within a React component, call `useScrapeSceneQuery` and pass it any options that fit your needs.
 * When your component renders, `useScrapeSceneQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useScrapeSceneQuery({
 *   variables: {
 *      scraper_id: // value for 'scraper_id'
 *      scene: // value for 'scene'
 *   },
 * });
 */
export function useScrapeSceneQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ScrapeSceneQuery, ScrapeSceneQueryVariables>) {
        return ApolloReactHooks.useQuery<ScrapeSceneQuery, ScrapeSceneQueryVariables>(ScrapeSceneDocument, baseOptions);
      }
export function useScrapeSceneLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ScrapeSceneQuery, ScrapeSceneQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ScrapeSceneQuery, ScrapeSceneQueryVariables>(ScrapeSceneDocument, baseOptions);
        }
export type ScrapeSceneQueryHookResult = ReturnType<typeof useScrapeSceneQuery>;
export type ScrapeSceneLazyQueryHookResult = ReturnType<typeof useScrapeSceneLazyQuery>;
export type ScrapeSceneQueryResult = ApolloReactCommon.QueryResult<ScrapeSceneQuery, ScrapeSceneQueryVariables>;
export const ScrapeSceneUrlDocument = gql`
    query ScrapeSceneURL($url: String!) {
  scrapeSceneURL(url: $url) {
    ...ScrapedSceneData
  }
}
    ${ScrapedSceneDataFragmentDoc}`;
export type ScrapeSceneUrlComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ScrapeSceneUrlQuery, ScrapeSceneUrlQueryVariables>, 'query'> & ({ variables: ScrapeSceneUrlQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const ScrapeSceneUrlComponent = (props: ScrapeSceneUrlComponentProps) => (
      <ApolloReactComponents.Query<ScrapeSceneUrlQuery, ScrapeSceneUrlQueryVariables> query={ScrapeSceneUrlDocument} {...props} />
    );
    

/**
 * __useScrapeSceneUrlQuery__
 *
 * To run a query within a React component, call `useScrapeSceneUrlQuery` and pass it any options that fit your needs.
 * When your component renders, `useScrapeSceneUrlQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useScrapeSceneUrlQuery({
 *   variables: {
 *      url: // value for 'url'
 *   },
 * });
 */
export function useScrapeSceneUrlQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ScrapeSceneUrlQuery, ScrapeSceneUrlQueryVariables>) {
        return ApolloReactHooks.useQuery<ScrapeSceneUrlQuery, ScrapeSceneUrlQueryVariables>(ScrapeSceneUrlDocument, baseOptions);
      }
export function useScrapeSceneUrlLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ScrapeSceneUrlQuery, ScrapeSceneUrlQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ScrapeSceneUrlQuery, ScrapeSceneUrlQueryVariables>(ScrapeSceneUrlDocument, baseOptions);
        }
export type ScrapeSceneUrlQueryHookResult = ReturnType<typeof useScrapeSceneUrlQuery>;
export type ScrapeSceneUrlLazyQueryHookResult = ReturnType<typeof useScrapeSceneUrlLazyQuery>;
export type ScrapeSceneUrlQueryResult = ApolloReactCommon.QueryResult<ScrapeSceneUrlQuery, ScrapeSceneUrlQueryVariables>;
export const ConfigurationDocument = gql`
    query Configuration {
  configuration {
    ...ConfigData
  }
}
    ${ConfigDataFragmentDoc}`;
export type ConfigurationComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ConfigurationQuery, ConfigurationQueryVariables>, 'query'>;

    export const ConfigurationComponent = (props: ConfigurationComponentProps) => (
      <ApolloReactComponents.Query<ConfigurationQuery, ConfigurationQueryVariables> query={ConfigurationDocument} {...props} />
    );
    

/**
 * __useConfigurationQuery__
 *
 * To run a query within a React component, call `useConfigurationQuery` and pass it any options that fit your needs.
 * When your component renders, `useConfigurationQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useConfigurationQuery({
 *   variables: {
 *   },
 * });
 */
export function useConfigurationQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<ConfigurationQuery, ConfigurationQueryVariables>) {
        return ApolloReactHooks.useQuery<ConfigurationQuery, ConfigurationQueryVariables>(ConfigurationDocument, baseOptions);
      }
export function useConfigurationLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<ConfigurationQuery, ConfigurationQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<ConfigurationQuery, ConfigurationQueryVariables>(ConfigurationDocument, baseOptions);
        }
export type ConfigurationQueryHookResult = ReturnType<typeof useConfigurationQuery>;
export type ConfigurationLazyQueryHookResult = ReturnType<typeof useConfigurationLazyQuery>;
export type ConfigurationQueryResult = ApolloReactCommon.QueryResult<ConfigurationQuery, ConfigurationQueryVariables>;
export const DirectoriesDocument = gql`
    query Directories($path: String) {
  directories(path: $path)
}
    `;
export type DirectoriesComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<DirectoriesQuery, DirectoriesQueryVariables>, 'query'>;

    export const DirectoriesComponent = (props: DirectoriesComponentProps) => (
      <ApolloReactComponents.Query<DirectoriesQuery, DirectoriesQueryVariables> query={DirectoriesDocument} {...props} />
    );
    

/**
 * __useDirectoriesQuery__
 *
 * To run a query within a React component, call `useDirectoriesQuery` and pass it any options that fit your needs.
 * When your component renders, `useDirectoriesQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useDirectoriesQuery({
 *   variables: {
 *      path: // value for 'path'
 *   },
 * });
 */
export function useDirectoriesQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<DirectoriesQuery, DirectoriesQueryVariables>) {
        return ApolloReactHooks.useQuery<DirectoriesQuery, DirectoriesQueryVariables>(DirectoriesDocument, baseOptions);
      }
export function useDirectoriesLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<DirectoriesQuery, DirectoriesQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<DirectoriesQuery, DirectoriesQueryVariables>(DirectoriesDocument, baseOptions);
        }
export type DirectoriesQueryHookResult = ReturnType<typeof useDirectoriesQuery>;
export type DirectoriesLazyQueryHookResult = ReturnType<typeof useDirectoriesLazyQuery>;
export type DirectoriesQueryResult = ApolloReactCommon.QueryResult<DirectoriesQuery, DirectoriesQueryVariables>;
export const MetadataImportDocument = gql`
    query MetadataImport {
  metadataImport
}
    `;
export type MetadataImportComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<MetadataImportQuery, MetadataImportQueryVariables>, 'query'>;

    export const MetadataImportComponent = (props: MetadataImportComponentProps) => (
      <ApolloReactComponents.Query<MetadataImportQuery, MetadataImportQueryVariables> query={MetadataImportDocument} {...props} />
    );
    

/**
 * __useMetadataImportQuery__
 *
 * To run a query within a React component, call `useMetadataImportQuery` and pass it any options that fit your needs.
 * When your component renders, `useMetadataImportQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetadataImportQuery({
 *   variables: {
 *   },
 * });
 */
export function useMetadataImportQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<MetadataImportQuery, MetadataImportQueryVariables>) {
        return ApolloReactHooks.useQuery<MetadataImportQuery, MetadataImportQueryVariables>(MetadataImportDocument, baseOptions);
      }
export function useMetadataImportLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<MetadataImportQuery, MetadataImportQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<MetadataImportQuery, MetadataImportQueryVariables>(MetadataImportDocument, baseOptions);
        }
export type MetadataImportQueryHookResult = ReturnType<typeof useMetadataImportQuery>;
export type MetadataImportLazyQueryHookResult = ReturnType<typeof useMetadataImportLazyQuery>;
export type MetadataImportQueryResult = ApolloReactCommon.QueryResult<MetadataImportQuery, MetadataImportQueryVariables>;
export const MetadataExportDocument = gql`
    query MetadataExport {
  metadataExport
}
    `;
export type MetadataExportComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<MetadataExportQuery, MetadataExportQueryVariables>, 'query'>;

    export const MetadataExportComponent = (props: MetadataExportComponentProps) => (
      <ApolloReactComponents.Query<MetadataExportQuery, MetadataExportQueryVariables> query={MetadataExportDocument} {...props} />
    );
    

/**
 * __useMetadataExportQuery__
 *
 * To run a query within a React component, call `useMetadataExportQuery` and pass it any options that fit your needs.
 * When your component renders, `useMetadataExportQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetadataExportQuery({
 *   variables: {
 *   },
 * });
 */
export function useMetadataExportQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<MetadataExportQuery, MetadataExportQueryVariables>) {
        return ApolloReactHooks.useQuery<MetadataExportQuery, MetadataExportQueryVariables>(MetadataExportDocument, baseOptions);
      }
export function useMetadataExportLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<MetadataExportQuery, MetadataExportQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<MetadataExportQuery, MetadataExportQueryVariables>(MetadataExportDocument, baseOptions);
        }
export type MetadataExportQueryHookResult = ReturnType<typeof useMetadataExportQuery>;
export type MetadataExportLazyQueryHookResult = ReturnType<typeof useMetadataExportLazyQuery>;
export type MetadataExportQueryResult = ApolloReactCommon.QueryResult<MetadataExportQuery, MetadataExportQueryVariables>;
export const MetadataScanDocument = gql`
    query MetadataScan($input: ScanMetadataInput!) {
  metadataScan(input: $input)
}
    `;
export type MetadataScanComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<MetadataScanQuery, MetadataScanQueryVariables>, 'query'> & ({ variables: MetadataScanQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const MetadataScanComponent = (props: MetadataScanComponentProps) => (
      <ApolloReactComponents.Query<MetadataScanQuery, MetadataScanQueryVariables> query={MetadataScanDocument} {...props} />
    );
    

/**
 * __useMetadataScanQuery__
 *
 * To run a query within a React component, call `useMetadataScanQuery` and pass it any options that fit your needs.
 * When your component renders, `useMetadataScanQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetadataScanQuery({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useMetadataScanQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<MetadataScanQuery, MetadataScanQueryVariables>) {
        return ApolloReactHooks.useQuery<MetadataScanQuery, MetadataScanQueryVariables>(MetadataScanDocument, baseOptions);
      }
export function useMetadataScanLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<MetadataScanQuery, MetadataScanQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<MetadataScanQuery, MetadataScanQueryVariables>(MetadataScanDocument, baseOptions);
        }
export type MetadataScanQueryHookResult = ReturnType<typeof useMetadataScanQuery>;
export type MetadataScanLazyQueryHookResult = ReturnType<typeof useMetadataScanLazyQuery>;
export type MetadataScanQueryResult = ApolloReactCommon.QueryResult<MetadataScanQuery, MetadataScanQueryVariables>;
export const MetadataGenerateDocument = gql`
    query MetadataGenerate($input: GenerateMetadataInput!) {
  metadataGenerate(input: $input)
}
    `;
export type MetadataGenerateComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<MetadataGenerateQuery, MetadataGenerateQueryVariables>, 'query'> & ({ variables: MetadataGenerateQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const MetadataGenerateComponent = (props: MetadataGenerateComponentProps) => (
      <ApolloReactComponents.Query<MetadataGenerateQuery, MetadataGenerateQueryVariables> query={MetadataGenerateDocument} {...props} />
    );
    

/**
 * __useMetadataGenerateQuery__
 *
 * To run a query within a React component, call `useMetadataGenerateQuery` and pass it any options that fit your needs.
 * When your component renders, `useMetadataGenerateQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetadataGenerateQuery({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useMetadataGenerateQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<MetadataGenerateQuery, MetadataGenerateQueryVariables>) {
        return ApolloReactHooks.useQuery<MetadataGenerateQuery, MetadataGenerateQueryVariables>(MetadataGenerateDocument, baseOptions);
      }
export function useMetadataGenerateLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<MetadataGenerateQuery, MetadataGenerateQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<MetadataGenerateQuery, MetadataGenerateQueryVariables>(MetadataGenerateDocument, baseOptions);
        }
export type MetadataGenerateQueryHookResult = ReturnType<typeof useMetadataGenerateQuery>;
export type MetadataGenerateLazyQueryHookResult = ReturnType<typeof useMetadataGenerateLazyQuery>;
export type MetadataGenerateQueryResult = ApolloReactCommon.QueryResult<MetadataGenerateQuery, MetadataGenerateQueryVariables>;
export const MetadataAutoTagDocument = gql`
    query MetadataAutoTag($input: AutoTagMetadataInput!) {
  metadataAutoTag(input: $input)
}
    `;
export type MetadataAutoTagComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<MetadataAutoTagQuery, MetadataAutoTagQueryVariables>, 'query'> & ({ variables: MetadataAutoTagQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const MetadataAutoTagComponent = (props: MetadataAutoTagComponentProps) => (
      <ApolloReactComponents.Query<MetadataAutoTagQuery, MetadataAutoTagQueryVariables> query={MetadataAutoTagDocument} {...props} />
    );
    

/**
 * __useMetadataAutoTagQuery__
 *
 * To run a query within a React component, call `useMetadataAutoTagQuery` and pass it any options that fit your needs.
 * When your component renders, `useMetadataAutoTagQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetadataAutoTagQuery({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useMetadataAutoTagQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<MetadataAutoTagQuery, MetadataAutoTagQueryVariables>) {
        return ApolloReactHooks.useQuery<MetadataAutoTagQuery, MetadataAutoTagQueryVariables>(MetadataAutoTagDocument, baseOptions);
      }
export function useMetadataAutoTagLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<MetadataAutoTagQuery, MetadataAutoTagQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<MetadataAutoTagQuery, MetadataAutoTagQueryVariables>(MetadataAutoTagDocument, baseOptions);
        }
export type MetadataAutoTagQueryHookResult = ReturnType<typeof useMetadataAutoTagQuery>;
export type MetadataAutoTagLazyQueryHookResult = ReturnType<typeof useMetadataAutoTagLazyQuery>;
export type MetadataAutoTagQueryResult = ApolloReactCommon.QueryResult<MetadataAutoTagQuery, MetadataAutoTagQueryVariables>;
export const MetadataCleanDocument = gql`
    query MetadataClean {
  metadataClean
}
    `;
export type MetadataCleanComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<MetadataCleanQuery, MetadataCleanQueryVariables>, 'query'>;

    export const MetadataCleanComponent = (props: MetadataCleanComponentProps) => (
      <ApolloReactComponents.Query<MetadataCleanQuery, MetadataCleanQueryVariables> query={MetadataCleanDocument} {...props} />
    );
    

/**
 * __useMetadataCleanQuery__
 *
 * To run a query within a React component, call `useMetadataCleanQuery` and pass it any options that fit your needs.
 * When your component renders, `useMetadataCleanQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetadataCleanQuery({
 *   variables: {
 *   },
 * });
 */
export function useMetadataCleanQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<MetadataCleanQuery, MetadataCleanQueryVariables>) {
        return ApolloReactHooks.useQuery<MetadataCleanQuery, MetadataCleanQueryVariables>(MetadataCleanDocument, baseOptions);
      }
export function useMetadataCleanLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<MetadataCleanQuery, MetadataCleanQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<MetadataCleanQuery, MetadataCleanQueryVariables>(MetadataCleanDocument, baseOptions);
        }
export type MetadataCleanQueryHookResult = ReturnType<typeof useMetadataCleanQuery>;
export type MetadataCleanLazyQueryHookResult = ReturnType<typeof useMetadataCleanLazyQuery>;
export type MetadataCleanQueryResult = ApolloReactCommon.QueryResult<MetadataCleanQuery, MetadataCleanQueryVariables>;
export const JobStatusDocument = gql`
    query JobStatus {
  jobStatus {
    progress
    status
    message
  }
}
    `;
export type JobStatusComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<JobStatusQuery, JobStatusQueryVariables>, 'query'>;

    export const JobStatusComponent = (props: JobStatusComponentProps) => (
      <ApolloReactComponents.Query<JobStatusQuery, JobStatusQueryVariables> query={JobStatusDocument} {...props} />
    );
    

/**
 * __useJobStatusQuery__
 *
 * To run a query within a React component, call `useJobStatusQuery` and pass it any options that fit your needs.
 * When your component renders, `useJobStatusQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useJobStatusQuery({
 *   variables: {
 *   },
 * });
 */
export function useJobStatusQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<JobStatusQuery, JobStatusQueryVariables>) {
        return ApolloReactHooks.useQuery<JobStatusQuery, JobStatusQueryVariables>(JobStatusDocument, baseOptions);
      }
export function useJobStatusLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<JobStatusQuery, JobStatusQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<JobStatusQuery, JobStatusQueryVariables>(JobStatusDocument, baseOptions);
        }
export type JobStatusQueryHookResult = ReturnType<typeof useJobStatusQuery>;
export type JobStatusLazyQueryHookResult = ReturnType<typeof useJobStatusLazyQuery>;
export type JobStatusQueryResult = ApolloReactCommon.QueryResult<JobStatusQuery, JobStatusQueryVariables>;
export const StopJobDocument = gql`
    query StopJob {
  stopJob
}
    `;
export type StopJobComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<StopJobQuery, StopJobQueryVariables>, 'query'>;

    export const StopJobComponent = (props: StopJobComponentProps) => (
      <ApolloReactComponents.Query<StopJobQuery, StopJobQueryVariables> query={StopJobDocument} {...props} />
    );
    

/**
 * __useStopJobQuery__
 *
 * To run a query within a React component, call `useStopJobQuery` and pass it any options that fit your needs.
 * When your component renders, `useStopJobQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useStopJobQuery({
 *   variables: {
 *   },
 * });
 */
export function useStopJobQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<StopJobQuery, StopJobQueryVariables>) {
        return ApolloReactHooks.useQuery<StopJobQuery, StopJobQueryVariables>(StopJobDocument, baseOptions);
      }
export function useStopJobLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<StopJobQuery, StopJobQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<StopJobQuery, StopJobQueryVariables>(StopJobDocument, baseOptions);
        }
export type StopJobQueryHookResult = ReturnType<typeof useStopJobQuery>;
export type StopJobLazyQueryHookResult = ReturnType<typeof useStopJobLazyQuery>;
export type StopJobQueryResult = ApolloReactCommon.QueryResult<StopJobQuery, StopJobQueryVariables>;
export const FindStudiosDocument = gql`
    query FindStudios($filter: FindFilterType) {
  findStudios(filter: $filter) {
    count
    studios {
      ...StudioData
    }
  }
}
    ${StudioDataFragmentDoc}`;
export type FindStudiosComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindStudiosQuery, FindStudiosQueryVariables>, 'query'>;

    export const FindStudiosComponent = (props: FindStudiosComponentProps) => (
      <ApolloReactComponents.Query<FindStudiosQuery, FindStudiosQueryVariables> query={FindStudiosDocument} {...props} />
    );
    

/**
 * __useFindStudiosQuery__
 *
 * To run a query within a React component, call `useFindStudiosQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindStudiosQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindStudiosQuery({
 *   variables: {
 *      filter: // value for 'filter'
 *   },
 * });
 */
export function useFindStudiosQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindStudiosQuery, FindStudiosQueryVariables>) {
        return ApolloReactHooks.useQuery<FindStudiosQuery, FindStudiosQueryVariables>(FindStudiosDocument, baseOptions);
      }
export function useFindStudiosLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindStudiosQuery, FindStudiosQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindStudiosQuery, FindStudiosQueryVariables>(FindStudiosDocument, baseOptions);
        }
export type FindStudiosQueryHookResult = ReturnType<typeof useFindStudiosQuery>;
export type FindStudiosLazyQueryHookResult = ReturnType<typeof useFindStudiosLazyQuery>;
export type FindStudiosQueryResult = ApolloReactCommon.QueryResult<FindStudiosQuery, FindStudiosQueryVariables>;
export const FindStudioDocument = gql`
    query FindStudio($id: ID!) {
  findStudio(id: $id) {
    ...StudioData
  }
}
    ${StudioDataFragmentDoc}`;
export type FindStudioComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<FindStudioQuery, FindStudioQueryVariables>, 'query'> & ({ variables: FindStudioQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const FindStudioComponent = (props: FindStudioComponentProps) => (
      <ApolloReactComponents.Query<FindStudioQuery, FindStudioQueryVariables> query={FindStudioDocument} {...props} />
    );
    

/**
 * __useFindStudioQuery__
 *
 * To run a query within a React component, call `useFindStudioQuery` and pass it any options that fit your needs.
 * When your component renders, `useFindStudioQuery` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useFindStudioQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useFindStudioQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<FindStudioQuery, FindStudioQueryVariables>) {
        return ApolloReactHooks.useQuery<FindStudioQuery, FindStudioQueryVariables>(FindStudioDocument, baseOptions);
      }
export function useFindStudioLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<FindStudioQuery, FindStudioQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<FindStudioQuery, FindStudioQueryVariables>(FindStudioDocument, baseOptions);
        }
export type FindStudioQueryHookResult = ReturnType<typeof useFindStudioQuery>;
export type FindStudioLazyQueryHookResult = ReturnType<typeof useFindStudioLazyQuery>;
export type FindStudioQueryResult = ApolloReactCommon.QueryResult<FindStudioQuery, FindStudioQueryVariables>;
export const MetadataUpdateDocument = gql`
    subscription MetadataUpdate {
  metadataUpdate {
    progress
    status
    message
  }
}
    `;
export type MetadataUpdateComponentProps = Omit<ApolloReactComponents.SubscriptionComponentOptions<MetadataUpdateSubscription, MetadataUpdateSubscriptionVariables>, 'subscription'>;

    export const MetadataUpdateComponent = (props: MetadataUpdateComponentProps) => (
      <ApolloReactComponents.Subscription<MetadataUpdateSubscription, MetadataUpdateSubscriptionVariables> subscription={MetadataUpdateDocument} {...props} />
    );
    

/**
 * __useMetadataUpdateSubscription__
 *
 * To run a query within a React component, call `useMetadataUpdateSubscription` and pass it any options that fit your needs.
 * When your component renders, `useMetadataUpdateSubscription` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMetadataUpdateSubscription({
 *   variables: {
 *   },
 * });
 */
export function useMetadataUpdateSubscription(baseOptions?: ApolloReactHooks.SubscriptionHookOptions<MetadataUpdateSubscription, MetadataUpdateSubscriptionVariables>) {
        return ApolloReactHooks.useSubscription<MetadataUpdateSubscription, MetadataUpdateSubscriptionVariables>(MetadataUpdateDocument, baseOptions);
      }
export type MetadataUpdateSubscriptionHookResult = ReturnType<typeof useMetadataUpdateSubscription>;
export type MetadataUpdateSubscriptionResult = ApolloReactCommon.SubscriptionResult<MetadataUpdateSubscription>;
export const LoggingSubscribeDocument = gql`
    subscription LoggingSubscribe {
  loggingSubscribe {
    ...LogEntryData
  }
}
    ${LogEntryDataFragmentDoc}`;
export type LoggingSubscribeComponentProps = Omit<ApolloReactComponents.SubscriptionComponentOptions<LoggingSubscribeSubscription, LoggingSubscribeSubscriptionVariables>, 'subscription'>;

    export const LoggingSubscribeComponent = (props: LoggingSubscribeComponentProps) => (
      <ApolloReactComponents.Subscription<LoggingSubscribeSubscription, LoggingSubscribeSubscriptionVariables> subscription={LoggingSubscribeDocument} {...props} />
    );
    

/**
 * __useLoggingSubscribeSubscription__
 *
 * To run a query within a React component, call `useLoggingSubscribeSubscription` and pass it any options that fit your needs.
 * When your component renders, `useLoggingSubscribeSubscription` returns an object from Apollo Client that contains loading, error, and data properties 
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useLoggingSubscribeSubscription({
 *   variables: {
 *   },
 * });
 */
export function useLoggingSubscribeSubscription(baseOptions?: ApolloReactHooks.SubscriptionHookOptions<LoggingSubscribeSubscription, LoggingSubscribeSubscriptionVariables>) {
        return ApolloReactHooks.useSubscription<LoggingSubscribeSubscription, LoggingSubscribeSubscriptionVariables>(LoggingSubscribeDocument, baseOptions);
      }
export type LoggingSubscribeSubscriptionHookResult = ReturnType<typeof useLoggingSubscribeSubscription>;
export type LoggingSubscribeSubscriptionResult = ApolloReactCommon.SubscriptionResult<LoggingSubscribeSubscription>;