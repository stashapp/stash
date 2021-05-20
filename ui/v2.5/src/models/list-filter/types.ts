// NOTE: add new enum values to the end, to ensure existing data

// is not impacted
export enum DisplayMode {
  Grid,
  List,
  Wall,
  Tagger,
}

export enum FilterMode {
  Scenes,
  Performers,
  Studios,
  Galleries,
  SceneMarkers,
  Movies,
  Tags,
  Images,
}

export interface ILabeledId {
  id: string;
  label: string;
}

export interface ILabeledValue {
  label: string;
  value: string;
}

export function encodeILabeledId(o: ILabeledId) {
  // escape \ to \\ so that it encodes to JSON correctly
  const adjustedLabel = o.label.replaceAll("\\", "\\\\");
  return { ...o, label: encodeURIComponent(adjustedLabel) };
}

export interface IOptionType {
  id: string;
  name?: string;
  image_path?: string;
}

export type CriterionType =
  | "none"
  | "path"
  | "rating"
  | "organized"
  | "o_counter"
  | "resolution"
  | "average_resolution"
  | "duration"
  | "favorite"
  | "hasMarkers"
  | "sceneIsMissing"
  | "imageIsMissing"
  | "performerIsMissing"
  | "galleryIsMissing"
  | "tagIsMissing"
  | "studioIsMissing"
  | "movieIsMissing"
  | "tags"
  | "sceneTags"
  | "performerTags"
  | "tag_count"
  | "performers"
  | "studios"
  | "movies"
  | "galleries"
  | "birth_year"
  | "age"
  | "ethnicity"
  | "country"
  | "hair_color"
  | "eye_color"
  | "height"
  | "weight"
  | "measurements"
  | "fake_tits"
  | "career_length"
  | "tattoos"
  | "piercings"
  | "aliases"
  | "gender"
  | "parent_studios"
  | "scene_count"
  | "marker_count"
  | "image_count"
  | "gallery_count"
  | "performer_count"
  | "death_year"
  | "url"
  | "stash_id";
