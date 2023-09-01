import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionType } from "../types";
import { CriterionOption, StringCriterion, Option } from "./criterion";

export class IsMissingCriterion extends StringCriterion {
  public modifierOptions = [];

  protected toCriterionInput(): string {
    return this.value;
  }
}

class IsMissingCriterionOptionClass extends CriterionOption {
  constructor(messageID: string, type: CriterionType, options: Option[]) {
    super({
      messageID,
      type,
      options,
      defaultModifier: CriterionModifier.Equals,
      makeCriterion: () => new IsMissingCriterion(this),
    });
  }
}

export const SceneIsMissingCriterionOption = new IsMissingCriterionOptionClass(
  "isMissing",
  "is_missing",
  [
    "title",
    "cover",
    "details",
    "url",
    "date",
    "galleries",
    "studio",
    "movie",
    "performers",
    "tags",
    "stash_id",
  ]
);

export const ImageIsMissingCriterionOption = new IsMissingCriterionOptionClass(
  "isMissing",
  "is_missing",
  ["title", "galleries", "studio", "performers", "tags"]
);

export const PerformerIsMissingCriterionOption =
  new IsMissingCriterionOptionClass("isMissing", "is_missing", [
    "url",
    "twitter",
    "instagram",
    "ethnicity",
    "country",
    "hair_color",
    "eye_color",
    "height",
    "weight",
    "measurements",
    "fake_tits",
    "career_length",
    "tattoos",
    "piercings",
    "aliases",
    "gender",
    "image",
    "details",
    "stash_id",
  ]);

export const GalleryIsMissingCriterionOption =
  new IsMissingCriterionOptionClass("isMissing", "is_missing", [
    "title",
    "details",
    "url",
    "date",
    "studio",
    "performers",
    "tags",
    "scenes",
  ]);

export const TagIsMissingCriterionOption = new IsMissingCriterionOptionClass(
  "isMissing",
  "is_missing",
  ["image"]
);

export const StudioIsMissingCriterionOption = new IsMissingCriterionOptionClass(
  "isMissing",
  "is_missing",
  ["image", "stash_id", "details"]
);

export const MovieIsMissingCriterionOption = new IsMissingCriterionOptionClass(
  "isMissing",
  "is_missing",
  ["front_image", "back_image", "scenes"]
);
