import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption, StringCriterion } from "./criterion";

export abstract class IsMissingCriterion extends StringCriterion {
  public modifierOptions = [];
  public modifier = CriterionModifier.Equals;

  protected toCriterionInput(): string {
    return this.value;
  }
}

export class SceneIsMissingCriterionOption extends CriterionOption {
  constructor() {
    super("isMissing", "sceneIsMissing", "is_missing");
  }
}

export class SceneIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(new SceneIsMissingCriterionOption(), [
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

export class ImageIsMissingCriterionOption extends CriterionOption {
  constructor() {
    super("isMissing", "imageIsMissing", "is_missing");
  }
}

export class ImageIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(new ImageIsMissingCriterionOption(), [
      "title",
      "galleries",
      "studio",
      "performers",
      "tags",
    ]);
  }
}

export class PerformerIsMissingCriterionOption extends CriterionOption {
  constructor() {
    super("isMissing", "performerIsMissing", "is_missing");
  }
}

export class PerformerIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(new PerformerIsMissingCriterionOption(), [
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
      "scenes",
      "image",
      "details",
      "stash_id",
    ]);
  }
}

export class GalleryIsMissingCriterionOption extends CriterionOption {
  constructor() {
    super("isMissing", "galleryIsMissing", "is_missing");
  }
}

export class GalleryIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(new GalleryIsMissingCriterionOption(), [
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

export class TagIsMissingCriterionOption extends CriterionOption {
  constructor() {
    super("isMissing", "tagIsMissing", "is_missing");
  }
}

export class TagIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(new TagIsMissingCriterionOption(), ["image"]);
  }
}

export class StudioIsMissingCriterionOption extends CriterionOption {
  constructor() {
    super("isMissing", "studioIsMissing", "is_missing");
  }
}

export class StudioIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(new StudioIsMissingCriterionOption(), [
      "image",
      "stash_id",
      "details",
    ]);
  }
}

export class MovieIsMissingCriterionOption extends CriterionOption {
  constructor() {
    super("isMissing", "movieIsMissing", "is_missing");
  }
}

export class MovieIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super(new MovieIsMissingCriterionOption(), [
      "front_image",
      "back_image",
      "scenes",
    ]);
  }
}
