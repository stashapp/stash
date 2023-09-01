import { CriterionOption, StringCriterion } from "./criterion";

export const HasMarkersCriterionOption = new CriterionOption({
  messageID: "hasMarkers",
  type: "has_markers",
  options: [true.toString(), false.toString()],
  makeCriterion: () => new HasMarkersCriterion(),
});

export class HasMarkersCriterion extends StringCriterion {
  constructor() {
    super(HasMarkersCriterionOption);
  }

  protected toCriterionInput(): string {
    return this.value;
  }
}
