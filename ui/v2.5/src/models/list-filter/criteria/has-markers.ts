import { CriterionOption, StringCriterion } from "./criterion";

export const HasMarkersCriterionOption = new CriterionOption({
  messageID: "hasMarkers",
  type: "hasMarkers",
  parameterName: "has_markers",
  options: [true.toString(), false.toString()],
});

export class HasMarkersCriterion extends StringCriterion {
  constructor() {
    super(HasMarkersCriterionOption);
  }

  protected toCriterionInput(): string {
    return this.value;
  }
}
