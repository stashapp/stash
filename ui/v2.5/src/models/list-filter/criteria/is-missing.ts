import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class SceneIsMissingCriterion extends Criterion {
  public type: CriterionType = "sceneIsMissing";
  public parameterName: string = "is_missing";
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
  public options: string[] = [
    "title",
    "url",
    "date",
    "gallery",
    "studio",
    "movie",
    "performers",
  ];
  public value: string = "";
}

export class SceneIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("sceneIsMissing");
  public value: CriterionType = "sceneIsMissing";
}

export class PerformerIsMissingCriterion extends Criterion {
  public type: CriterionType = "performerIsMissing";
  public parameterName: string = "is_missing";
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
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
    "scenes"
  ];
  public value: string = "";
}

export class PerformerIsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("performerIsMissing");
  public value: CriterionType = "performerIsMissing";
}
