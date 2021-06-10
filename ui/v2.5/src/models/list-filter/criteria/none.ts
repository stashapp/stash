import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionOption } from "./criterion";

export const NoneCriterionOption = new CriterionOption("none", "none");
export class NoneCriterion extends Criterion<string> {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
  public options: undefined;
  public value: string = "none";

  constructor() {
    super(NoneCriterionOption);
  }

  // eslint-disable-next-line class-methods-use-this
  public getLabelValue(): string {
    return "";
  }
}
