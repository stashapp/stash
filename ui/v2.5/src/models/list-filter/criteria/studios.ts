import {
  CriterionOption,
  IHierarchicalLabeledIdCriterion,
  ILabeledIdCriterion,
} from "./criterion";

export const StudiosCriterionOption = new CriterionOption("studios", "studios");
export class StudiosCriterion extends IHierarchicalLabeledIdCriterion {
  constructor() {
    super(StudiosCriterionOption, false);
  }
}

export const ParentStudiosCriterionOption = new CriterionOption(
  "parent_studios",
  "parent_studios",
  "parents"
);
export class ParentStudiosCriterion extends ILabeledIdCriterion {
  constructor() {
    super(ParentStudiosCriterionOption, false);
  }
}
