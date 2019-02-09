/* tslint:disable */
// Generated in 2019-02-09T01:48:09-08:00
export type Maybe<T> = T | null;

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

export namespace FindScenes {
  export type Variables = {
    filter?: Maybe<FindFilterType>;
    scene_filter?: Maybe<SceneFilterType>;
    scene_ids?: Maybe<number[]>;
  };

  export type Query = {
    __typename?: "Query";

    findScenes: FindScenes;
  };

  export type FindScenes = {
    __typename?: "FindScenesResultType";

    count: number;

    scenes: Scenes[];
  };

  export type Scenes = SlimSceneData.Fragment;
}

export namespace FindScene {
  export type Variables = {
    id: string;
    checksum?: Maybe<string>;
  };

  export type Query = {
    __typename?: "Query";

    findScene: Maybe<FindScene>;

    sceneMarkerTags: SceneMarkerTags[];
  };

  export type FindScene = SceneData.Fragment;

  export type SceneMarkerTags = {
    __typename?: "SceneMarkerTag";

    tag: Tag;

    scene_markers: SceneMarkers[];
  };

  export type Tag = {
    __typename?: "Tag";

    id: string;

    name: string;
  };

  export type SceneMarkers = SceneMarkerData.Fragment;
}

export namespace FindSceneForEditing {
  export type Variables = {
    id?: Maybe<string>;
  };

  export type Query = {
    __typename?: "Query";

    findScene: Maybe<FindScene>;

    allPerformers: AllPerformers[];

    allTags: AllTags[];

    allStudios: AllStudios[];

    validGalleriesForScene: ValidGalleriesForScene[];
  };

  export type FindScene = SceneData.Fragment;

  export type AllPerformers = {
    __typename?: "Performer";

    id: string;

    name: Maybe<string>;

    birthdate: Maybe<string>;

    image_path: Maybe<string>;
  };

  export type AllTags = {
    __typename?: "Tag";

    id: string;

    name: string;
  };

  export type AllStudios = {
    __typename?: "Studio";

    id: string;

    name: string;
  };

  export type ValidGalleriesForScene = {
    __typename?: "Gallery";

    id: string;

    path: string;
  };
}

export namespace FindSceneMarkers {
  export type Variables = {
    filter?: Maybe<FindFilterType>;
    scene_marker_filter?: Maybe<SceneMarkerFilterType>;
  };

  export type Query = {
    __typename?: "Query";

    findSceneMarkers: FindSceneMarkers;
  };

  export type FindSceneMarkers = {
    __typename?: "FindSceneMarkersResultType";

    count: number;

    scene_markers: SceneMarkers[];
  };

  export type SceneMarkers = SceneMarkerData.Fragment;
}

export namespace SceneWall {
  export type Variables = {
    q?: Maybe<string>;
  };

  export type Query = {
    __typename?: "Query";

    sceneWall: SceneWall[];
  };

  export type SceneWall = SceneData.Fragment;
}

export namespace MarkerWall {
  export type Variables = {
    q?: Maybe<string>;
  };

  export type Query = {
    __typename?: "Query";

    markerWall: MarkerWall[];
  };

  export type MarkerWall = SceneMarkerData.Fragment;
}

export namespace FindPerformers {
  export type Variables = {
    filter?: Maybe<FindFilterType>;
    performer_filter?: Maybe<PerformerFilterType>;
  };

  export type Query = {
    __typename?: "Query";

    findPerformers: FindPerformers;
  };

  export type FindPerformers = {
    __typename?: "FindPerformersResultType";

    count: number;

    performers: Performers[];
  };

  export type Performers = PerformerData.Fragment;
}

export namespace FindPerformer {
  export type Variables = {
    id: string;
  };

  export type Query = {
    __typename?: "Query";

    findPerformer: Maybe<FindPerformer>;
  };

  export type FindPerformer = PerformerData.Fragment;
}

export namespace FindStudios {
  export type Variables = {
    filter?: Maybe<FindFilterType>;
  };

  export type Query = {
    __typename?: "Query";

    findStudios: FindStudios;
  };

  export type FindStudios = {
    __typename?: "FindStudiosResultType";

    count: number;

    studios: Studios[];
  };

  export type Studios = StudioData.Fragment;
}

export namespace FindStudio {
  export type Variables = {
    id: string;
  };

  export type Query = {
    __typename?: "Query";

    findStudio: Maybe<FindStudio>;
  };

  export type FindStudio = StudioData.Fragment;
}

export namespace FindGalleries {
  export type Variables = {
    filter?: Maybe<FindFilterType>;
  };

  export type Query = {
    __typename?: "Query";

    findGalleries: FindGalleries;
  };

  export type FindGalleries = {
    __typename?: "FindGalleriesResultType";

    count: number;

    galleries: Galleries[];
  };

  export type Galleries = GalleryData.Fragment;
}

export namespace FindGallery {
  export type Variables = {
    id: string;
  };

  export type Query = {
    __typename?: "Query";

    findGallery: Maybe<FindGallery>;
  };

  export type FindGallery = GalleryData.Fragment;
}

export namespace FindTag {
  export type Variables = {
    id: string;
  };

  export type Query = {
    __typename?: "Query";

    findTag: Maybe<FindTag>;
  };

  export type FindTag = TagData.Fragment;
}

export namespace MarkerStrings {
  export type Variables = {
    q?: Maybe<string>;
    sort?: Maybe<string>;
  };

  export type Query = {
    __typename?: "Query";

    markerStrings: (Maybe<MarkerStrings>)[];
  };

  export type MarkerStrings = {
    __typename?: "MarkerStringsResultType";

    id: string;

    count: number;

    title: string;
  };
}

export namespace ScrapeFreeones {
  export type Variables = {
    performer_name: string;
  };

  export type Query = {
    __typename?: "Query";

    scrapeFreeones: Maybe<ScrapeFreeones>;
  };

  export type ScrapeFreeones = {
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
}

export namespace ScrapeFreeonesPerformers {
  export type Variables = {
    q: string;
  };

  export type Query = {
    __typename?: "Query";

    scrapeFreeonesPerformerList: string[];
  };
}

export namespace AllTags {
  export type Variables = {};

  export type Query = {
    __typename?: "Query";

    allTags: AllTags[];
  };

  export type AllTags = TagData.Fragment;
}

export namespace AllPerformersForFilter {
  export type Variables = {};

  export type Query = {
    __typename?: "Query";

    allPerformers: AllPerformers[];
  };

  export type AllPerformers = SlimPerformerData.Fragment;
}

export namespace AllTagsForFilter {
  export type Variables = {};

  export type Query = {
    __typename?: "Query";

    allTags: AllTags[];
  };

  export type AllTags = {
    __typename?: "Tag";

    id: string;

    name: string;
  };
}

export namespace Stats {
  export type Variables = {};

  export type Query = {
    __typename?: "Query";

    stats: Stats;
  };

  export type Stats = {
    __typename?: "StatsResultType";

    scene_count: number;

    gallery_count: number;

    performer_count: number;

    studio_count: number;

    tag_count: number;
  };
}

export namespace SceneUpdate {
  export type Variables = {
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

  export type Mutation = {
    __typename?: "Mutation";

    sceneUpdate: Maybe<SceneUpdate>;
  };

  export type SceneUpdate = SceneData.Fragment;
}

export namespace PerformerCreate {
  export type Variables = {
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

  export type Mutation = {
    __typename?: "Mutation";

    performerCreate: Maybe<PerformerCreate>;
  };

  export type PerformerCreate = PerformerData.Fragment;
}

export namespace PerformerUpdate {
  export type Variables = {
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

  export type Mutation = {
    __typename?: "Mutation";

    performerUpdate: Maybe<PerformerUpdate>;
  };

  export type PerformerUpdate = PerformerData.Fragment;
}

export namespace StudioCreate {
  export type Variables = {
    name: string;
    url?: Maybe<string>;
    image: string;
  };

  export type Mutation = {
    __typename?: "Mutation";

    studioCreate: Maybe<StudioCreate>;
  };

  export type StudioCreate = StudioData.Fragment;
}

export namespace StudioUpdate {
  export type Variables = {
    id: string;
    name?: Maybe<string>;
    url?: Maybe<string>;
    image?: Maybe<string>;
  };

  export type Mutation = {
    __typename?: "Mutation";

    studioUpdate: Maybe<StudioUpdate>;
  };

  export type StudioUpdate = StudioData.Fragment;
}

export namespace TagCreate {
  export type Variables = {
    name: string;
  };

  export type Mutation = {
    __typename?: "Mutation";

    tagCreate: Maybe<TagCreate>;
  };

  export type TagCreate = TagData.Fragment;
}

export namespace TagDestroy {
  export type Variables = {
    id: string;
  };

  export type Mutation = {
    __typename?: "Mutation";

    tagDestroy: boolean;
  };
}

export namespace TagUpdate {
  export type Variables = {
    id: string;
    name: string;
  };

  export type Mutation = {
    __typename?: "Mutation";

    tagUpdate: Maybe<TagUpdate>;
  };

  export type TagUpdate = TagData.Fragment;
}

export namespace SceneMarkerCreate {
  export type Variables = {
    title: string;
    seconds: number;
    scene_id: string;
    primary_tag_id: string;
    tag_ids?: Maybe<string[]>;
  };

  export type Mutation = {
    __typename?: "Mutation";

    sceneMarkerCreate: Maybe<SceneMarkerCreate>;
  };

  export type SceneMarkerCreate = SceneMarkerData.Fragment;
}

export namespace SceneMarkerUpdate {
  export type Variables = {
    id: string;
    title: string;
    seconds: number;
    scene_id: string;
    primary_tag_id: string;
    tag_ids?: Maybe<string[]>;
  };

  export type Mutation = {
    __typename?: "Mutation";

    sceneMarkerUpdate: Maybe<SceneMarkerUpdate>;
  };

  export type SceneMarkerUpdate = SceneMarkerData.Fragment;
}

export namespace SceneMarkerDestroy {
  export type Variables = {
    id: string;
  };

  export type Mutation = {
    __typename?: "Mutation";

    sceneMarkerDestroy: boolean;
  };
}

export namespace MetadataImport {
  export type Variables = {};

  export type Query = {
    __typename?: "Query";

    metadataImport: string;
  };
}

export namespace MetadataExport {
  export type Variables = {};

  export type Query = {
    __typename?: "Query";

    metadataExport: string;
  };
}

export namespace MetadataScan {
  export type Variables = {};

  export type Query = {
    __typename?: "Query";

    metadataScan: string;
  };
}

export namespace MetadataGenerate {
  export type Variables = {};

  export type Query = {
    __typename?: "Query";

    metadataGenerate: string;
  };
}

export namespace MetadataClean {
  export type Variables = {};

  export type Query = {
    __typename?: "Query";

    metadataClean: string;
  };
}

export namespace MetadataUpdate {
  export type Variables = {};

  export type Subscription = {
    __typename?: "Subscription";

    metadataUpdate: string;
  };
}

export namespace GalleryData {
  export type Fragment = {
    __typename?: "Gallery";

    id: string;

    checksum: string;

    path: string;

    title: Maybe<string>;

    files: Files[];
  };

  export type Files = {
    __typename?: "GalleryFilesType";

    index: number;

    name: Maybe<string>;

    path: Maybe<string>;
  };
}

export namespace SlimPerformerData {
  export type Fragment = {
    __typename?: "Performer";

    id: string;

    name: Maybe<string>;

    image_path: Maybe<string>;
  };
}

export namespace PerformerData {
  export type Fragment = {
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
}

export namespace SceneMarkerData {
  export type Fragment = {
    __typename?: "SceneMarker";

    id: string;

    title: string;

    seconds: number;

    stream: string;

    preview: string;

    scene: Scene;

    primary_tag: PrimaryTag;

    tags: Tags[];
  };

  export type Scene = {
    __typename?: "Scene";

    id: string;
  };

  export type PrimaryTag = {
    __typename?: "Tag";

    id: string;

    name: string;
  };

  export type Tags = {
    __typename?: "Tag";

    id: string;

    name: string;
  };
}

export namespace SlimSceneData {
  export type Fragment = {
    __typename?: "Scene";

    id: string;

    checksum: string;

    title: Maybe<string>;

    details: Maybe<string>;

    url: Maybe<string>;

    date: Maybe<string>;

    rating: Maybe<number>;

    path: string;

    file: File;

    paths: Paths;

    scene_markers: SceneMarkers[];

    gallery: Maybe<Gallery>;

    studio: Maybe<Studio>;

    tags: Tags[];

    performers: Performers[];
  };

  export type File = {
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

  export type Paths = {
    __typename?: "ScenePathsType";

    screenshot: Maybe<string>;

    preview: Maybe<string>;

    stream: Maybe<string>;

    webp: Maybe<string>;

    vtt: Maybe<string>;

    chapters_vtt: Maybe<string>;
  };

  export type SceneMarkers = {
    __typename?: "SceneMarker";

    id: string;

    title: string;

    seconds: number;
  };

  export type Gallery = {
    __typename?: "Gallery";

    id: string;

    path: string;

    title: Maybe<string>;
  };

  export type Studio = {
    __typename?: "Studio";

    id: string;

    name: string;

    image_path: Maybe<string>;
  };

  export type Tags = {
    __typename?: "Tag";

    id: string;

    name: string;
  };

  export type Performers = {
    __typename?: "Performer";

    id: string;

    name: Maybe<string>;

    favorite: boolean;

    image_path: Maybe<string>;
  };
}

export namespace SceneData {
  export type Fragment = {
    __typename?: "Scene";

    id: string;

    checksum: string;

    title: Maybe<string>;

    details: Maybe<string>;

    url: Maybe<string>;

    date: Maybe<string>;

    rating: Maybe<number>;

    path: string;

    file: File;

    paths: Paths;

    scene_markers: SceneMarkers[];

    is_streamable: boolean;

    gallery: Maybe<Gallery>;

    studio: Maybe<Studio>;

    tags: Tags[];

    performers: Performers[];
  };

  export type File = {
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

  export type Paths = {
    __typename?: "ScenePathsType";

    screenshot: Maybe<string>;

    preview: Maybe<string>;

    stream: Maybe<string>;

    webp: Maybe<string>;

    vtt: Maybe<string>;

    chapters_vtt: Maybe<string>;
  };

  export type SceneMarkers = SceneMarkerData.Fragment;

  export type Gallery = GalleryData.Fragment;

  export type Studio = StudioData.Fragment;

  export type Tags = TagData.Fragment;

  export type Performers = PerformerData.Fragment;
}

export namespace StudioData {
  export type Fragment = {
    __typename?: "Studio";

    id: string;

    checksum: string;

    name: string;

    url: Maybe<string>;

    image_path: Maybe<string>;

    scene_count: Maybe<number>;
  };
}

export namespace TagData {
  export type Fragment = {
    __typename?: "Tag";

    id: string;

    name: string;

    scene_count: Maybe<number>;

    scene_marker_count: Maybe<number>;
  };
}

// ====================================================
// START: Apollo Angular template
// ====================================================

import { Injectable } from "@angular/core";
import * as Apollo from "apollo-angular";

import gql from "graphql-tag";

// ====================================================
// GraphQL Fragments
// ====================================================

export const SlimPerformerDataFragment = gql`
  fragment SlimPerformerData on Performer {
    id
    name
    image_path
  }
`;

export const SlimSceneDataFragment = gql`
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

export const SceneMarkerDataFragment = gql`
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

export const GalleryDataFragment = gql`
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

export const StudioDataFragment = gql`
  fragment StudioData on Studio {
    id
    checksum
    name
    url
    image_path
    scene_count
  }
`;

export const TagDataFragment = gql`
  fragment TagData on Tag {
    id
    name
    scene_count
    scene_marker_count
  }
`;

export const PerformerDataFragment = gql`
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

export const SceneDataFragment = gql`
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

  ${SceneMarkerDataFragment}
  ${GalleryDataFragment}
  ${StudioDataFragment}
  ${TagDataFragment}
  ${PerformerDataFragment}
`;

// ====================================================
// Apollo Services
// ====================================================

@Injectable({
  providedIn: "root"
})
export class FindScenesGQL extends Apollo.Query<
  FindScenes.Query,
  FindScenes.Variables
> {
  document: any = gql`
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

    ${SlimSceneDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class FindSceneGQL extends Apollo.Query<
  FindScene.Query,
  FindScene.Variables
> {
  document: any = gql`
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

    ${SceneDataFragment}
    ${SceneMarkerDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class FindSceneForEditingGQL extends Apollo.Query<
  FindSceneForEditing.Query,
  FindSceneForEditing.Variables
> {
  document: any = gql`
    query FindSceneForEditing($id: ID) {
      findScene(id: $id) {
        ...SceneData
      }
      allPerformers {
        id
        name
        birthdate
        image_path
      }
      allTags {
        id
        name
      }
      allStudios {
        id
        name
      }
      validGalleriesForScene(scene_id: $id) {
        id
        path
      }
    }

    ${SceneDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class FindSceneMarkersGQL extends Apollo.Query<
  FindSceneMarkers.Query,
  FindSceneMarkers.Variables
> {
  document: any = gql`
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

    ${SceneMarkerDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class SceneWallGQL extends Apollo.Query<
  SceneWall.Query,
  SceneWall.Variables
> {
  document: any = gql`
    query SceneWall($q: String) {
      sceneWall(q: $q) {
        ...SceneData
      }
    }

    ${SceneDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class MarkerWallGQL extends Apollo.Query<
  MarkerWall.Query,
  MarkerWall.Variables
> {
  document: any = gql`
    query MarkerWall($q: String) {
      markerWall(q: $q) {
        ...SceneMarkerData
      }
    }

    ${SceneMarkerDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class FindPerformersGQL extends Apollo.Query<
  FindPerformers.Query,
  FindPerformers.Variables
> {
  document: any = gql`
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

    ${PerformerDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class FindPerformerGQL extends Apollo.Query<
  FindPerformer.Query,
  FindPerformer.Variables
> {
  document: any = gql`
    query FindPerformer($id: ID!) {
      findPerformer(id: $id) {
        ...PerformerData
      }
    }

    ${PerformerDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class FindStudiosGQL extends Apollo.Query<
  FindStudios.Query,
  FindStudios.Variables
> {
  document: any = gql`
    query FindStudios($filter: FindFilterType) {
      findStudios(filter: $filter) {
        count
        studios {
          ...StudioData
        }
      }
    }

    ${StudioDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class FindStudioGQL extends Apollo.Query<
  FindStudio.Query,
  FindStudio.Variables
> {
  document: any = gql`
    query FindStudio($id: ID!) {
      findStudio(id: $id) {
        ...StudioData
      }
    }

    ${StudioDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class FindGalleriesGQL extends Apollo.Query<
  FindGalleries.Query,
  FindGalleries.Variables
> {
  document: any = gql`
    query FindGalleries($filter: FindFilterType) {
      findGalleries(filter: $filter) {
        count
        galleries {
          ...GalleryData
        }
      }
    }

    ${GalleryDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class FindGalleryGQL extends Apollo.Query<
  FindGallery.Query,
  FindGallery.Variables
> {
  document: any = gql`
    query FindGallery($id: ID!) {
      findGallery(id: $id) {
        ...GalleryData
      }
    }

    ${GalleryDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class FindTagGQL extends Apollo.Query<FindTag.Query, FindTag.Variables> {
  document: any = gql`
    query FindTag($id: ID!) {
      findTag(id: $id) {
        ...TagData
      }
    }

    ${TagDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class MarkerStringsGQL extends Apollo.Query<
  MarkerStrings.Query,
  MarkerStrings.Variables
> {
  document: any = gql`
    query MarkerStrings($q: String, $sort: String) {
      markerStrings(q: $q, sort: $sort) {
        id
        count
        title
      }
    }
  `;
}
@Injectable({
  providedIn: "root"
})
export class ScrapeFreeonesGQL extends Apollo.Query<
  ScrapeFreeones.Query,
  ScrapeFreeones.Variables
> {
  document: any = gql`
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
}
@Injectable({
  providedIn: "root"
})
export class ScrapeFreeonesPerformersGQL extends Apollo.Query<
  ScrapeFreeonesPerformers.Query,
  ScrapeFreeonesPerformers.Variables
> {
  document: any = gql`
    query ScrapeFreeonesPerformers($q: String!) {
      scrapeFreeonesPerformerList(query: $q)
    }
  `;
}
@Injectable({
  providedIn: "root"
})
export class AllTagsGQL extends Apollo.Query<AllTags.Query, AllTags.Variables> {
  document: any = gql`
    query AllTags {
      allTags {
        ...TagData
      }
    }

    ${TagDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class AllPerformersForFilterGQL extends Apollo.Query<
  AllPerformersForFilter.Query,
  AllPerformersForFilter.Variables
> {
  document: any = gql`
    query AllPerformersForFilter {
      allPerformers {
        ...SlimPerformerData
      }
    }

    ${SlimPerformerDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class AllTagsForFilterGQL extends Apollo.Query<
  AllTagsForFilter.Query,
  AllTagsForFilter.Variables
> {
  document: any = gql`
    query AllTagsForFilter {
      allTags {
        id
        name
      }
    }
  `;
}
@Injectable({
  providedIn: "root"
})
export class StatsGQL extends Apollo.Query<Stats.Query, Stats.Variables> {
  document: any = gql`
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
}
@Injectable({
  providedIn: "root"
})
export class SceneUpdateGQL extends Apollo.Mutation<
  SceneUpdate.Mutation,
  SceneUpdate.Variables
> {
  document: any = gql`
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

    ${SceneDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class PerformerCreateGQL extends Apollo.Mutation<
  PerformerCreate.Mutation,
  PerformerCreate.Variables
> {
  document: any = gql`
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

    ${PerformerDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class PerformerUpdateGQL extends Apollo.Mutation<
  PerformerUpdate.Mutation,
  PerformerUpdate.Variables
> {
  document: any = gql`
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

    ${PerformerDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class StudioCreateGQL extends Apollo.Mutation<
  StudioCreate.Mutation,
  StudioCreate.Variables
> {
  document: any = gql`
    mutation StudioCreate($name: String!, $url: String, $image: String!) {
      studioCreate(input: { name: $name, url: $url, image: $image }) {
        ...StudioData
      }
    }

    ${StudioDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class StudioUpdateGQL extends Apollo.Mutation<
  StudioUpdate.Mutation,
  StudioUpdate.Variables
> {
  document: any = gql`
    mutation StudioUpdate(
      $id: ID!
      $name: String
      $url: String
      $image: String
    ) {
      studioUpdate(input: { id: $id, name: $name, url: $url, image: $image }) {
        ...StudioData
      }
    }

    ${StudioDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class TagCreateGQL extends Apollo.Mutation<
  TagCreate.Mutation,
  TagCreate.Variables
> {
  document: any = gql`
    mutation TagCreate($name: String!) {
      tagCreate(input: { name: $name }) {
        ...TagData
      }
    }

    ${TagDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class TagDestroyGQL extends Apollo.Mutation<
  TagDestroy.Mutation,
  TagDestroy.Variables
> {
  document: any = gql`
    mutation TagDestroy($id: ID!) {
      tagDestroy(input: { id: $id })
    }
  `;
}
@Injectable({
  providedIn: "root"
})
export class TagUpdateGQL extends Apollo.Mutation<
  TagUpdate.Mutation,
  TagUpdate.Variables
> {
  document: any = gql`
    mutation TagUpdate($id: ID!, $name: String!) {
      tagUpdate(input: { id: $id, name: $name }) {
        ...TagData
      }
    }

    ${TagDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class SceneMarkerCreateGQL extends Apollo.Mutation<
  SceneMarkerCreate.Mutation,
  SceneMarkerCreate.Variables
> {
  document: any = gql`
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

    ${SceneMarkerDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class SceneMarkerUpdateGQL extends Apollo.Mutation<
  SceneMarkerUpdate.Mutation,
  SceneMarkerUpdate.Variables
> {
  document: any = gql`
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

    ${SceneMarkerDataFragment}
  `;
}
@Injectable({
  providedIn: "root"
})
export class SceneMarkerDestroyGQL extends Apollo.Mutation<
  SceneMarkerDestroy.Mutation,
  SceneMarkerDestroy.Variables
> {
  document: any = gql`
    mutation SceneMarkerDestroy($id: ID!) {
      sceneMarkerDestroy(id: $id)
    }
  `;
}
@Injectable({
  providedIn: "root"
})
export class MetadataImportGQL extends Apollo.Query<
  MetadataImport.Query,
  MetadataImport.Variables
> {
  document: any = gql`
    query MetadataImport {
      metadataImport
    }
  `;
}
@Injectable({
  providedIn: "root"
})
export class MetadataExportGQL extends Apollo.Query<
  MetadataExport.Query,
  MetadataExport.Variables
> {
  document: any = gql`
    query MetadataExport {
      metadataExport
    }
  `;
}
@Injectable({
  providedIn: "root"
})
export class MetadataScanGQL extends Apollo.Query<
  MetadataScan.Query,
  MetadataScan.Variables
> {
  document: any = gql`
    query MetadataScan {
      metadataScan
    }
  `;
}
@Injectable({
  providedIn: "root"
})
export class MetadataGenerateGQL extends Apollo.Query<
  MetadataGenerate.Query,
  MetadataGenerate.Variables
> {
  document: any = gql`
    query MetadataGenerate {
      metadataGenerate
    }
  `;
}
@Injectable({
  providedIn: "root"
})
export class MetadataCleanGQL extends Apollo.Query<
  MetadataClean.Query,
  MetadataClean.Variables
> {
  document: any = gql`
    query MetadataClean {
      metadataClean
    }
  `;
}
@Injectable({
  providedIn: "root"
})
export class MetadataUpdateGQL extends Apollo.Subscription<
  MetadataUpdate.Subscription,
  MetadataUpdate.Variables
> {
  document: any = gql`
    subscription MetadataUpdate {
      metadataUpdate
    }
  `;
}

// ====================================================
// END: Apollo Angular template
// ====================================================
