import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class InteractiveCriterion extends Criterion {
  public type: CriterionType = "interactive";
  public parameterName: string = "interactive";
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
  public options: string[] = [true.toString(), false.toString()];
  public value: string = "";
}

export class InteractiveCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("interactive");
  public value: CriterionType = "interactive";
}
