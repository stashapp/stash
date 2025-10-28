import { CriterionValue, ISavedCriterion } from "./criteria/criterion";

export type SavedObjectFilter = {
  [K in CriterionType]?: ISavedCriterion<CriterionValue>;
};

export type SavedUIOptions = {
  display_mode?: DisplayMode;
  zoom_index?: number;
};

// NOTE: add new enum values to the end, to ensure existing data
// is not impacted
export enum DisplayMode {
  Grid,
  List,
  Wall,
  Tagger,
}

export interface ILabeledId {
  id: string;
  label: string;
}

export interface ILabeledValue {
  label: string;
  value: string;
}

export interface ILabeledValueListValue {
  items: ILabeledId[];
  excluded: ILabeledId[];
}

export interface IHierarchicalLabelValue {
  items: ILabeledId[];
  excluded: ILabeledId[];
  depth: number;
}

export interface IRangeValue<V> {
  value: V | undefined;
  value2: V | undefined;
}

export type INumberValue = IRangeValue<number>;
export type IDateValue = IRangeValue<string>;
export type ITimestampValue = IRangeValue<string>;
export interface IPHashDuplicationValue {
  duplicated: boolean;
  distance?: number; // currently not implemented
}

export interface IStashIDValue {
  endpoint: string;
  stashID: string;
}

export interface IPhashDistanceValue {
  value: string;
  distance?: number;
}

export function criterionIsHierarchicalLabelValue(
  value: unknown
): value is IHierarchicalLabelValue {
  return (
    typeof value === "object" && !!value && "items" in value && "depth" in value
  );
}

export function criterionIsNumberValue(value: unknown): value is INumberValue {
  return (
    typeof value === "object" &&
    !!value &&
    "value" in value &&
    "value2" in value
  );
}

export function criterionIsStashIDValue(
  value: unknown
): value is IStashIDValue {
  return (
    typeof value === "object" &&
    !!value &&
    "endpoint" in value &&
    "stashID" in value
  );
}

export function criterionIsDateValue(value: unknown): value is IDateValue {
  return (
    typeof value === "object" &&
    !!value &&
    "value" in value &&
    "value2" in value
  );
}

export function criterionIsTimestampValue(
  value: unknown
): value is ITimestampValue {
  return (
    typeof value === "object" &&
    !!value &&
    "value" in value &&
    "value2" in value
  );
}

export interface IOptionType {
  id: string;
  name?: string;
  image_path?: string;
}

export type CriterionType =
  | "path"
  | "rating100"
  | "organized"
  | "o_counter"
  | "resolution"
  | "average_resolution"
  | "framerate"
  | "bitrate"
  | "video_codec"
  | "audio_codec"
  | "duration"
  | "filter_favorites"
  | "favorite"
  | "has_markers"
  | "is_missing"
  | "tags"
  | "scene_tags"
  | "performer_tags"
  | "studio_tags"
  | "tag_count"
  | "performers"
  | "studios"
  | "scenes"
  | "groups"
  | "movies" // legacy
  | "containing_groups"
  | "containing_group_count"
  | "sub_groups"
  | "sub_group_count"
  | "galleries"
  | "birth_year"
  | "age"
  | "ethnicity"
  | "country"
  | "hair_color"
  | "eye_color"
  | "height_cm"
  | "weight"
  | "measurements"
  | "fake_tits"
  | "penis_length"
  | "circumcised"
  | "career_length"
  | "tattoos"
  | "piercings"
  | "aliases"
  | "gender"
  | "parents"
  | "children"
  | "scene_count"
  | "marker_count"
  | "image_count"
  | "gallery_count"
  | "performer_count"
  | "studio_count"
  | "group_count"
  | "death_year"
  | "url"
  | "interactive"
  | "interactive_speed"
  | "captions"
  | "resume_time"
  | "play_count"
  | "play_duration"
  | "last_played_at"
  | "name"
  | "details"
  | "title"
  | "oshash"
  | "orientation"
  | "checksum"
  | "phash_distance"
  | "director"
  | "synopsis"
  | "parent_count"
  | "child_count"
  | "performer_favorite"
  | "favorite"
  | "performer_age"
  | "duplicated"
  | "ignore_auto_tag"
  | "file_count"
  | "stash_id_endpoint"
  | "date"
  | "created_at"
  | "updated_at"
  | "birthdate"
  | "death_date"
  | "scene_date"
  | "scene_created_at"
  | "scene_updated_at"
  | "description"
  | "code"
  | "photographer"
  | "disambiguation"
  | "has_chapters"
  | "sort_name"
  | "custom_fields";
