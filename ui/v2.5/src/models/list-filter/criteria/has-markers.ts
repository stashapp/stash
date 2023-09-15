import {
  StringBooleanCriterion,
  StringBooleanCriterionOption,
} from "./criterion";

export const HasMarkersCriterionOption = new StringBooleanCriterionOption(
  "hasMarkers",
  "has_markers",
  () => new HasMarkersCriterion()
);

export class HasMarkersCriterion extends StringBooleanCriterion {
  constructor() {
    super(HasMarkersCriterionOption);
  }
}
