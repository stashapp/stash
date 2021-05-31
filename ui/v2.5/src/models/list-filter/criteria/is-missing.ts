import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption, StringCriterion } from "./criterion";

export abstract class IsMissingCriterion extends StringCriterion {
  public modifierOptions = [];
  public modifier = CriterionModifier.Equals;

  protected toCriterionInput(): string {
    return this.value;
  }
}

export const SceneIsMissingCriterionOption = new CriterionOption(
  "isMissing",
  "sceneIsMissing",
  "is_missing"
);

export class SceneIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(SceneIsMissingCriterionOption, [
      "title",
      "details",
      "url",
      "date",
      "galleries",
      "studio",
      "movie",
      "performers",
      "tags",
      "stash_id",
    ]);
  }
}

export const ImageIsMissingCriterionOption = new CriterionOption(
  "isMissing",
  "imageIsMissing",
  "is_missing"
);

export class ImageIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(ImageIsMissingCriterionOption, [
      "title",
      "galleries",
      "studio",
      "performers",
      "tags",
    ]);
  }
}

export const PerformerIsMissingCriterionOption = new CriterionOption(
  "isMissing",
  "performerIsMissing",
  "is_missing"
);

export class PerformerIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(PerformerIsMissingCriterionOption, [
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
  }
}

export const GalleryIsMissingCriterionOption = new CriterionOption(
  "isMissing",
  "galleryIsMissing",
  "is_missing"
);

export class GalleryIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(GalleryIsMissingCriterionOption, [
      "title",
      "details",
      "url",
      "date",
      "studio",
      "performers",
      "tags",
      "scenes",
    ]);
  }
}

export const TagIsMissingCriterionOption = new CriterionOption(
  "isMissing",
  "tagIsMissing",
  "is_missing"
);

export class TagIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(TagIsMissingCriterionOption, ["image"]);
  }
}

export const StudioIsMissingCriterionOption = new CriterionOption(
  "isMissing",
  "studioIsMissing",
  "is_missing"
);

export class StudioIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(StudioIsMissingCriterionOption, ["image", "stash_id", "details"]);
  }
}

export const MovieIsMissingCriterionOption = new CriterionOption(
  "isMissing",
  "movieIsMissing",
  "is_missing"
);

export class MovieIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(MovieIsMissingCriterionOption, [
      "front_image",
      "back_image",
      "scenes",
    ]);
  }
}
