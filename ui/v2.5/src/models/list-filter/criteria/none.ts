import { Criterion, StringCriterionOption } from "./criterion";

export const NoneCriterionOption = new StringCriterionOption("none", "none");
export class NoneCriterion extends Criterion<string> {
  constructor() {
    super(NoneCriterionOption, "none");
  }

  protected getLabelValue(): string {
    return "";
  }
}
