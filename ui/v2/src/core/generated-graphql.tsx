/* tslint:disable */
// Generated in 2019-03-23T12:23:40-07:00
export type Maybe<T> = T | undefined;

export interface SceneFilterType {
  /** Filter by rating */
  rating?: Maybe<number>;
  /** Filter by resolution */
  resolution?: Maybe<ResolutionEnum>;
  /** Filter to only include scenes which have markers. `true` or `false` */
  has_markers?: Maybe<string>;
  /** Filter to only include scenes missing this property */
  is_missing?: Maybe<string>;
  /** Filter to only include scenes with this studio */
  studio_id?: Maybe<string>;
  /** Filter to only include scenes with these tags */
  tags?: Maybe<string[]>;
  /** Filter to only include scenes with this performer */
  performer_id?: Maybe<string>;
}

export interface FindFilterType {
  q?: Maybe<string>;

  page?: Maybe<number>;

  per_page?: Maybe<number>;

  sort?: Maybe<string>;

  direction?: Maybe<SortDirectionEnum>;
}

export interface SceneMarkerFilterType {
  /** Filter to only include scene markers with this tag */
  tag_id?: Maybe<string>;
  /** Filter to only include scene markers with these tags */
  tags?: Maybe<string[]>;
  /** Filter to only include scene markers attached to a scene with these tags */
  scene_tags?: Maybe<string[]>;
  /** Filter to only include scene markers with these performers */
  performers?: Maybe<string[]>;
}

export interface PerformerFilterType {
  /** Filter by favorite */
  filter_favorites?: Maybe<boolean>;
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
  image: string;
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

export interface StudioCreateInput {
  name: string;

  url?: Maybe<string>;
  /** This should be base64 encoded */
  image: string;
}

export interface StudioUpdateInput {
  id: string;

  name?: Maybe<string>;

  url?: Maybe<string>;
  /** This should be base64 encoded */
  image?: Maybe<string>;
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
  image: string;
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
};

export type SceneUpdateMutation = {
  __typename?: "Mutation";

  sceneUpdate: Maybe<SceneUpdateSceneUpdate>;
};

export type SceneUpdateSceneUpdate = SceneDataFragment;

export type StudioCreateVariables = {
  name: string;
  url?: Maybe<string>;
  image: string;
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

  markerStrings: (Maybe<MarkerStringsMarkerStrings>)[];
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

export type MetadataScanVariables = {};

export type MetadataScanQuery = {
  __typename?: "Query";

  metadataScan: string;
};

export type MetadataGenerateVariables = {};

export type MetadataGenerateQuery = {
  __typename?: "Query";

  metadataGenerate: string;
};

export type MetadataCleanVariables = {};

export type MetadataCleanQuery = {
  __typename?: "Query";

  metadataClean: string;
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

  metadataUpdate: string;
};

export type ConfigGeneralDataFragment = {
  __typename?: "ConfigGeneralResult";

  stashes: Maybe<string[]>;
};

export type ConfigDataFragment = {
  __typename?: "ConfigResult";

  general: ConfigDataGeneral;
};

export type ConfigDataGeneral = ConfigGeneralDataFragment;

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
  }
`;

export const ConfigDataFragmentDoc = gql`
  fragment ConfigData on ConfigResult {
    general {
      ...ConfigGeneralData
    }
  }

  ${ConfigGeneralDataFragmentDoc}
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
    $image: String!
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
export const StudioCreateDocument = gql`
  mutation StudioCreate($name: String!, $url: String, $image: String!) {
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
  query MetadataScan {
    metadataScan
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
  query MetadataGenerate {
    metadataGenerate
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
    metadataUpdate
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
