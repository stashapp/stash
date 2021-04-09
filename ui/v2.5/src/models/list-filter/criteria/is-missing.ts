import { CriterionModifier } from "src/core/generated-graphql";
import {
  Criterion,
  CriterionType,
  ICriterionOption,
  StringCriterion,
} from "./criterion";

export abstract class IsMissingCriterion extends StringCriterion {
  public modifierOptions = [];
  public modifier = CriterionModifier.Equals;

  constructor(type: CriterionType, options?: string[]) {
    super(type, "is_missing", options);
  }
}

export class SceneIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super("sceneIsMissing", [
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

export class SceneIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("sceneIsMissing");
  public value: CriterionType = "sceneIsMissing";
}

export class ImageIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super("imageIsMissing", [
      "title",
      "galleries",
      "studio",
      "performers",
      "tags",
    ]);
  }
}

export class ImageIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("imageIsMissing");
  public value: CriterionType = "imageIsMissing";
}

export class PerformerIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super("performerIsMissing", [
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

export class PerformerIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("performerIsMissing");
  public value: CriterionType = "performerIsMissing";
}

export class GalleryIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super("galleryIsMissing", [
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

export class GalleryIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("galleryIsMissing");
  public value: CriterionType = "galleryIsMissing";
}

export class TagIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super("tagIsMissing", ["image"]);
  }
}

export class TagIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("tagIsMissing");
  public value: CriterionType = "tagIsMissing";
}

export class StudioIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super("studioIsMissing", ["image", "details", "stash_id"]);
  }
}

export class StudioIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("studioIsMissing");
  public value: CriterionType = "studioIsMissing";
}

export class MovieIsMissingCriterion extends IsMissingCriterion {
  constructor() {
    super("movieIsMissing", ["front_image", "back_image", "scenes"]);
  }
}

export class MovieIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("movieIsMissing");
  public value: CriterionType = "movieIsMissing";
}
