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
    "gallery",
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

export class PerformerIsMissingCriterion extends IsMissingCriterion {
  public type: CriterionType = "performerIsMissing";
  public options: string[] = [
    "url",
    "twitter",
    "instagram",
    "ethnicity",
    "country",
    "eye_color",
    "height",
    "measurements",
    "fake_tits",
    "career_length",
    "tattoos",
    "piercings",
    "aliases",
    "gender",
    "scenes",
  ];
}

export class PerformerIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("performerIsMissing");
  public value: CriterionType = "performerIsMissing";
}

export class GalleryIsMissingCriterion extends IsMissingCriterion {
  public type: CriterionType = "galleryIsMissing";
  public options: string[] = ["scene"];
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
