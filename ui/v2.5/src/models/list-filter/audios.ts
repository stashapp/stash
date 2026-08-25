import {
  createMandatoryNumberCriterionOption,
  createMandatoryStringCriterionOption,
  createStringCriterionOption,
  createDateCriterionOption,
  createMandatoryTimestampCriterionOption,
  createDurationCriterionOption,
} from "./criteria/criterion";
import { AudioIsMissingCriterionOption } from "./criteria/is-missing";
import {
  GroupsCriterionOption,
  LegacyMoviesCriterionOption,
} from "./criteria/groups";
import { GalleriesCriterionOption } from "./criteria/galleries";
import { OrganizedCriterionOption } from "./criteria/organized";
import { PerformersCriterionOption } from "./criteria/performers";
import { StudiosCriterionOption } from "./criteria/studios";
import {
  PerformerTagsCriterionOption,
  TagsCriterionOption,
} from "./criteria/tags";
import { ListFilterOptions, MediaSortByOptions } from "./filter-options";
import { DisplayMode } from "./types";
import { PerformerFavoriteCriterionOption } from "./criteria/favorite";
import { RatingCriterionOption } from "./criteria/rating";
import { PathCriterionOption } from "./criteria/path";
import { CustomFieldsCriterionOption } from "./criteria/custom-fields";
import { FolderCriterionOption } from "./criteria/folder";

const defaultSortBy = "date";

// NOTE - audio has no video-only sorts (framerate, resolution, perceptual
// similarity, interactive) - see pkg/sqlite/audio.go audioSortOptions
const sortByOptions = [
  "organized",
  "date",
  "file_count",
  "filesize",
  "duration",
  "bitrate",
  "sample_rate",
  "last_played_at",
  "resume_time",
  "play_duration",
  "play_count",
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
      messageID: "group_audio_number",
      value: "group_audio_number",
    },
    {
      messageID: "audio_code",
      value: "code",
    },
  ]);

// no Wall - audio has no previews and often no cover image
const displayModeOptions = [DisplayMode.Grid, DisplayMode.List];

export const AudioPerformerAgeCriterionOption =
  createMandatoryNumberCriterionOption("performer_age");

export const AudioDurationCriterionOption =
  createDurationCriterionOption("duration");

const criterionOptions = [
  createStringCriterionOption("title"),
  createStringCriterionOption("code", "audio_code"),
  PathCriterionOption,
  FolderCriterionOption,
  createStringCriterionOption("details"),
  createMandatoryStringCriterionOption("oshash", "media_info.oshash"),
  createStringCriterionOption("checksum", "media_info.md5"),
  OrganizedCriterionOption,
  RatingCriterionOption,
  createMandatoryNumberCriterionOption("o_counter", "o_count", {
    sfwMessageID: "o_count_sfw",
  }),
  createMandatoryNumberCriterionOption("bitrate"),
  createMandatoryNumberCriterionOption("sample_rate"),
  createStringCriterionOption("audio_codec"),
  AudioDurationCriterionOption,
  createDurationCriterionOption("resume_time"),
  createDurationCriterionOption("play_duration"),
  createMandatoryNumberCriterionOption("play_count"),
  createMandatoryTimestampCriterionOption("last_played_at"),
  AudioIsMissingCriterionOption,
  TagsCriterionOption,
  createMandatoryNumberCriterionOption("tag_count"),
  PerformerTagsCriterionOption,
  PerformersCriterionOption,
  createMandatoryNumberCriterionOption("performer_count"),
  AudioPerformerAgeCriterionOption,
  PerformerFavoriteCriterionOption,
  StudiosCriterionOption,
  GroupsCriterionOption,
  LegacyMoviesCriterionOption,
  GalleriesCriterionOption,
  createStringCriterionOption("url"),
  createMandatoryNumberCriterionOption("file_count"),
  createDateCriterionOption("date"),
  createMandatoryTimestampCriterionOption("created_at"),
  createMandatoryTimestampCriterionOption("updated_at"),
  CustomFieldsCriterionOption,
];

export const AudioListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
