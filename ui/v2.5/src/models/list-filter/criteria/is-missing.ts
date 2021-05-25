import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export abstract class IsMissingCriterion extends Criterion {
  public parameterName: string = "is_missing";
  public modifierOptions = [];
  public modifier = CriterionModifier.Equals;
  public value: string = "";
}

export class SceneIsMissingCriterion extends IsMissingCriterion {
  public type: CriterionType = "sceneIsMissing";
  public options: string[] = [
    "title",
    "details",
    "url",
    "date",
    "galleries",
    "studio",
    "movie",
    "performers",
    "tags",
  ];
}

export class SceneIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("sceneIsMissing");
  public value: CriterionType = "sceneIsMissing";
}

export class ImageIsMissingCriterion extends IsMissingCriterion {
  public type: CriterionType = "imageIsMissing";
  public options: string[] = [
    "title",
    "galleries",
    "studio",
    "performers",
    "tags",
  ];
}

export class ImageIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("imageIsMissing");
  public value: CriterionType = "imageIsMissing";
}

export class PerformerIsMissingCriterion extends IsMissingCriterion {
  public type: CriterionType = "performerIsMissing";
  public options: string[] = [
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
  ];
}

export class PerformerIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("performerIsMissing");
  public value: CriterionType = "performerIsMissing";
}

export class GalleryIsMissingCriterion extends IsMissingCriterion {
  public type: CriterionType = "galleryIsMissing";
  public options: string[] = [
    "title",
    "details",
    "url",
    "date",
    "studio",
    "performers",
    "tags",
    "scenes",
  ];
}

export class GalleryIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("galleryIsMissing");
  public value: CriterionType = "galleryIsMissing";
}

export class TagIsMissingCriterion extends IsMissingCriterion {
  public type: CriterionType = "tagIsMissing";
  public options: string[] = ["image"];
}

export class TagIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("tagIsMissing");
  public value: CriterionType = "tagIsMissing";
}

export class StudioIsMissingCriterion extends IsMissingCriterion {
  public type: CriterionType = "studioIsMissing";
  public options: string[] = ["image", "details"];
}

export class StudioIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("studioIsMissing");
  public value: CriterionType = "studioIsMissing";
}

export class MovieIsMissingCriterion extends IsMissingCriterion {
  public type: CriterionType = "movieIsMissing";
  public options: string[] = ["front_image", "back_image", "scenes"];
}

export class MovieIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("movieIsMissing");
  public value: CriterionType = "movieIsMissing";
}
