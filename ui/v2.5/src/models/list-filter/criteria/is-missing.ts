import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class IsMissingCriterion extends Criterion<string, string> {
  public type: CriterionType = "isMissing";
  public parameterName: string = "is_missing";
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
  public options: string[] = [
    "title",
    "url",
    "date",
    "gallery",
    "studio",
    "performers"
  ];
  public value: string = "";
}

export class IsMissingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("isMissing");
  public value: CriterionType = "isMissing";
}
