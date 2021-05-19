import { CriterionOption, StringCriterion } from "./criterion";

export const HasMarkersCriterionOption = new CriterionOption(
  "hasMarkers",
  "hasMarkers",
  "has_markers"
);

export class HasMarkersCriterion extends StringCriterion {
  constructor() {
    super(HasMarkersCriterionOption, [true.toString(), false.toString()]);
  }

  protected toCriterionInput(): string {
    return this.value;
  }
}
