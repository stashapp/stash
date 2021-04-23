import { BooleanCriterion, CriterionOption } from "./criterion";

export class HasMarkersCriterionOption extends CriterionOption {
  constructor() {
    super("hasMarkers", "hasMarkers", "has_markers");
  }
}

export class HasMarkersCriterion extends BooleanCriterion {
  constructor() {
    super(new HasMarkersCriterionOption());
  }
}
