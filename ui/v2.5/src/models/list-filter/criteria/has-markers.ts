import { CriterionOption, StringCriterion } from "./criterion";

export class HasMarkersCriterionOption extends CriterionOption {
  constructor() {
    super("hasMarkers", "hasMarkers", "has_markers");
  }
}

export class HasMarkersCriterion extends StringCriterion {
  constructor() {
    super(new HasMarkersCriterionOption(), [true.toString(), false.toString()]);
  }

  protected toCriterionInput(): string {
    return this.value;
  }
}
