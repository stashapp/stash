import {
  IHierarchicalLabeledIdCriterion,
  ILabeledIdCriterion,
  ILabeledIdCriterionOption,
} from "./criterion";

export const StudiosCriterionOption = new ILabeledIdCriterionOption(
  "studios",
  "studios",
  "studios",
  false
);

export class StudiosCriterion extends IHierarchicalLabeledIdCriterion {
  constructor() {
    super(StudiosCriterionOption);
  }
}

export const ParentStudiosCriterionOption = new ILabeledIdCriterionOption(
  "parent_studios",
  "parent_studios",
  "parents",
  false
);
export class ParentStudiosCriterion extends ILabeledIdCriterion {
  constructor() {
    super(ParentStudiosCriterionOption);
  }
}
