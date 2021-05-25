import { CriterionOption, ILabeledIdCriterion } from "./criterion";

abstract class AbstractStudiosCriterion extends ILabeledIdCriterion {
  constructor(type: CriterionOption) {
    super(type, false);
  }
}

export const StudiosCriterionOption = new CriterionOption("studios", "studios");
export class StudiosCriterion extends AbstractStudiosCriterion {
  constructor() {
    super(StudiosCriterionOption);
  }
}

export const ParentStudiosCriterionOption = new CriterionOption(
  "parent_studios",
  "parent_studios"
);
export class ParentStudiosCriterion extends AbstractStudiosCriterion {
  constructor() {
    super(ParentStudiosCriterionOption);
  }
}
