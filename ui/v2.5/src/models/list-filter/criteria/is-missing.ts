import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionType } from "../types";
import { ModifierCriterionOption, StringCriterion, Option } from "./criterion";

export class IsMissingCriterion extends StringCriterion {
  public toCriterionInput(): string {
    return this.value;
  }
}

class IsMissingCriterionOption extends ModifierCriterionOption {
  constructor(messageID: string, type: CriterionType, options: Option[]) {
    super({
      messageID,
      type,
      options,
      modifierOptions: [],
      defaultModifier: CriterionModifier.Equals,
      makeCriterion: () => new IsMissingCriterion(this),
    });
  }
}

export const SceneIsMissingCriterionOption = new IsMissingCriterionOption(
  "isMissing",
  "is_missing",
  [
    "title",
    "code",
    "details",
    "director",
    "url",
    "date",
    "rating",
    "cover",
    "galleries",
    "studio",
    "group",
    "performers",
    "tags",
    "stash_id",
  ]
);

export const ImageIsMissingCriterionOption = new IsMissingCriterionOption(
  "isMissing",
  "is_missing",
  [
    "title",
    "details",
    "photographer",
    "url",
    "date",
    "code",
    "rating",
    "galleries",
    "studio",
    "performers",
    "tags",
  ]
);

export const PerformerIsMissingCriterionOption = new IsMissingCriterionOption(
  "isMissing",
  "is_missing",
  [
    "url",
    "ethnicity",
    "country",
    "hair_color",
    "eye_color",
    "height",
    "weight",
    "measurements",
    "fake_tits",
    "penis_length",
    "circumcised",
    "career_start",
    "career_end",
    "tattoos",
    "piercings",
    "aliases",
    "gender",
    "birthdate",
    "death_date",
    "disambiguation",
    "tags",
    "image",
    "details",
    "rating",
    "stash_id",
  ]
);

export const GalleryIsMissingCriterionOption = new IsMissingCriterionOption(
  "isMissing",
  "is_missing",
  [
    "title",
    "code",
    "details",
    "photographer",
    "url",
    "date",
    "rating",
    "cover",
    "studio",
    "performers",
    "tags",
    "scenes",
  ]
);

export const TagIsMissingCriterionOption = new IsMissingCriterionOption(
  "isMissing",
  "is_missing",
  ["image", "aliases", "description", "stash_id"]
);

export const StudioIsMissingCriterionOption = new IsMissingCriterionOption(
  "isMissing",
  "is_missing",
  ["image", "stash_id", "details", "url", "aliases", "tags", "rating"]
);

export const GroupIsMissingCriterionOption = new IsMissingCriterionOption(
  "isMissing",
  "is_missing",
  [
    "aliases",
    "description",
    "director",
    "date",
    "url",
    "rating",
    "studio",
    "performers",
    "tags",
    "front_image",
    "back_image",
    "scenes",
  ]
);
