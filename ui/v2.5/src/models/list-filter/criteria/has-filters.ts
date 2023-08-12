import { CriterionOption, StringCriterion } from "./criterion";

export const HasFiltersCriterionOption = new CriterionOption({
  messageID: "hasFilters",
  type: "hasFilters",
  parameterName: "has_filters",
  options: [true.toString(), false.toString()],
});

export class HasFiltersCriterion extends StringCriterion {
  constructor() {
    super(HasFiltersCriterionOption);
  }

  protected toCriterionInput(): string {
    return this.value;
  }
}
