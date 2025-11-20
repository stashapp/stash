import {
  createStringCriterionOption,
  createDateCriterionOption,
  createMandatoryTimestampCriterionOption,
  createDurationCriterionOption,
  createMandatoryNumberCriterionOption,
} from "./criteria/criterion";
import { GroupIsMissingCriterionOption } from "./criteria/is-missing";
import { StudiosCriterionOption } from "./criteria/studios";
import { PerformersCriterionOption } from "./criteria/performers";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";
import { RatingCriterionOption } from "./criteria/rating";
// import { StudioTagsCriterionOption } from "./criteria/tags";
import { TagsCriterionOption } from "./criteria/tags";
import {
  ContainingGroupsCriterionOption,
  SubGroupsCriterionOption,
} from "./criteria/groups";

const defaultSortBy = "name";

const sortByOptions = [
  "name",
  "random",
  "date",
  "duration",
  "rating",
  "tag_count",
  "sub_group_order",
]
  .map(ListFilterOptions.createSortBy)
  .concat([
    {
      messageID: "scene_count",
      value: "scenes_count",
    },
    {
      messageID: "o_count",
      value: "o_counter",
      sfwMessageID: "o_count_sfw",
    },
  ]);
const displayModeOptions = [DisplayMode.Grid];
const criterionOptions = [
  // StudioTagsCriterionOption,
  StudiosCriterionOption,
  GroupIsMissingCriterionOption,
  createStringCriterionOption("url"),
  createStringCriterionOption("name"),
  createStringCriterionOption("director"),
  createStringCriterionOption("synopsis"),
  createDurationCriterionOption("duration"),
  RatingCriterionOption,
  PerformersCriterionOption,
  createDateCriterionOption("date"),
  createMandatoryNumberCriterionOption("o_counter", "o_count", {
    sfwMessageID: "o_count_sfw",
  }),
  ContainingGroupsCriterionOption,
  SubGroupsCriterionOption,
  createMandatoryNumberCriterionOption("containing_group_count"),
  createMandatoryNumberCriterionOption("sub_group_count"),
  TagsCriterionOption,
  createMandatoryNumberCriterionOption("tag_count"),
  createMandatoryTimestampCriterionOption("created_at"),
  createMandatoryTimestampCriterionOption("updated_at"),
];

export const GroupListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
