import {
  createMandatoryNumberCriterionOption,
  createMandatoryStringCriterionOption,
  createStringCriterionOption,
  createDateCriterionOption,
  createMandatoryTimestampCriterionOption,
  createDurationCriterionOption,
} from "./criteria/criterion";
import { HasMarkersCriterionOption } from "./criteria/has-markers";
import { SceneIsMissingCriterionOption } from "./criteria/is-missing";
import {
  GroupsCriterionOption,
  LegacyMoviesCriterionOption,
} from "./criteria/groups";
import { GalleriesCriterionOption } from "./criteria/galleries";
import { OrganizedCriterionOption } from "./criteria/organized";
import { PerformersCriterionOption } from "./criteria/performers";
import { ResolutionCriterionOption } from "./criteria/resolution";
import { StudiosCriterionOption } from "./criteria/studios";
import { InteractiveCriterionOption } from "./criteria/interactive";
import {
  PerformerTagsCriterionOption,
  // StudioTagsCriterionOption,
  TagsCriterionOption,
} from "./criteria/tags";
import { ListFilterOptions, MediaSortByOptions } from "./filter-options";
import { DisplayMode } from "./types";
import {
  DuplicatedCriterionOption,
  PhashCriterionOption,
} from "./criteria/phash";
import { PerformerFavoriteCriterionOption } from "./criteria/favorite";
import { CaptionsCriterionOption } from "./criteria/captions";
import { StashIDCriterionOption } from "./criteria/stash-ids";
import { RatingCriterionOption } from "./criteria/rating";
import { PathCriterionOption } from "./criteria/path";
import { OrientationCriterionOption } from "./criteria/orientation";

const defaultSortBy = "date";
const sortByOptions = [
  "organized",
  "date",
  "file_count",
  "filesize",
  "duration",
  "framerate",
  "resolution",
  "bitrate",
  "last_played_at",
  "resume_time",
  "play_duration",
  "play_count",
  "interactive",
  "interactive_speed",
  "perceptual_similarity",
  "performer_age",
  "studio",
  ...MediaSortByOptions,
]
  .map(ListFilterOptions.createSortBy)
  .concat([
    {
      messageID: "o_count",
      value: "o_counter",
      sfwMessageID: "o_count_sfw",
    },
    {
      messageID: "last_o_at",
      value: "last_o_at",
      sfwMessageID: "last_o_at_sfw",
    },
    {
      messageID: "group_scene_number",
      value: "group_scene_number",
    },
    {
      messageID: "scene_code",
      value: "code",
    },
  ]);
const displayModeOptions = [
  DisplayMode.Grid,
  DisplayMode.List,
  DisplayMode.Wall,
  DisplayMode.Tagger,
];

export const PerformerAgeCriterionOption =
  createMandatoryNumberCriterionOption("performer_age");

export const DurationCriterionOption =
  createDurationCriterionOption("duration");

const criterionOptions = [
  createStringCriterionOption("title"),
  createStringCriterionOption("code", "scene_code"),
  PathCriterionOption,
  createStringCriterionOption("details"),
  createStringCriterionOption("director"),
  createMandatoryStringCriterionOption("oshash", "media_info.hash"),
  createStringCriterionOption("checksum", "media_info.checksum"),
  PhashCriterionOption,
  DuplicatedCriterionOption,
  OrganizedCriterionOption,
  RatingCriterionOption,
  createMandatoryNumberCriterionOption("o_counter", "o_count", {
    sfwMessageID: "o_count_sfw",
  }),
  ResolutionCriterionOption,
  OrientationCriterionOption,
  createMandatoryNumberCriterionOption("framerate"),
  createMandatoryNumberCriterionOption("bitrate"),
  createStringCriterionOption("video_codec"),
  createStringCriterionOption("audio_codec"),
  DurationCriterionOption,
  createDurationCriterionOption("resume_time"),
  createDurationCriterionOption("play_duration"),
  createMandatoryNumberCriterionOption("play_count"),
  createMandatoryTimestampCriterionOption("last_played_at"),
  HasMarkersCriterionOption,
  SceneIsMissingCriterionOption,
  TagsCriterionOption,
  createMandatoryNumberCriterionOption("tag_count"),
  PerformerTagsCriterionOption,
  PerformersCriterionOption,
  createMandatoryNumberCriterionOption("performer_count"),
  PerformerAgeCriterionOption,
  PerformerFavoriteCriterionOption,
  // StudioTagsCriterionOption,
  StudiosCriterionOption,
  GroupsCriterionOption,
  LegacyMoviesCriterionOption,
  GalleriesCriterionOption,
  createStringCriterionOption("url"),
  StashIDCriterionOption,
  InteractiveCriterionOption,
  CaptionsCriterionOption,
  createMandatoryNumberCriterionOption("interactive_speed"),
  createMandatoryNumberCriterionOption("file_count"),
  createDateCriterionOption("date"),
  createMandatoryTimestampCriterionOption("created_at"),
  createMandatoryTimestampCriterionOption("updated_at"),
];

export const SceneListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
