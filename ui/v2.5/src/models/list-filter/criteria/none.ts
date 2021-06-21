import { Criterion, StringCriterionOption } from "./criterion";

export const NoneCriterionOption = new StringCriterionOption(
  "none",
  "none",
  "none"
);
export class NoneCriterion extends Criterion<string> {
  constructor() {
    super(NoneCriterionOption, "none");
  }

  // eslint-disable-next-line class-methods-use-this
  public getLabelValue(): string {
    return "";
  }
}
