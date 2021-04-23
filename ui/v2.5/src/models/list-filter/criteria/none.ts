import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionOption } from "./criterion";

export class NoneCriterionOption extends CriterionOption {
  constructor() {
    super("none", "none");
  }
}
export class NoneCriterion extends Criterion<string> {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
  public options: undefined;
  public value: string = "none";

  constructor() {
    super(new NoneCriterionOption());
  }

  // eslint-disable-next-line class-methods-use-this
  public getLabelValue(): string {
    return "";
  }
}
