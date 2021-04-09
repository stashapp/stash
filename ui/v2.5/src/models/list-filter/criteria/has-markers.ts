import {
  BooleanCriterion,
  Criterion,
  CriterionType,
  ICriterionOption,
} from "./criterion";

export class HasMarkersCriterion extends BooleanCriterion {
  constructor() {
    super("hasMarkers", "has_markers");
  }
}

export class HasMarkersCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("hasMarkers");
  public value: CriterionType = "hasMarkers";
}
