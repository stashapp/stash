import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class NoneCriterion extends Criterion<string> {
  public type: CriterionType = "none";
  public parameterName: string = "";
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
  public options: undefined;
  public value: string = "none";

  // eslint-disable-next-line class-methods-use-this
  public getLabelValue(): string {
    return "";
  }
}

export class NoneCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("none");
  public value: CriterionType = "none";
}
