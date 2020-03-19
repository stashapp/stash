import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class NoneCriterion extends Criterion {
  public type: CriterionType = "none";
  public parameterName: string = "";
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
  public options: undefined;
  public value: string = "none";
}

export class NoneCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("none");
  public value: CriterionType = "none";
}
