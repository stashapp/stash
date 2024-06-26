import { ILabeledIdCriterion, ILabeledIdCriterionOption } from "./criterion";

const inputType = "groups";

export const GroupsCriterionOption = new ILabeledIdCriterionOption(
  "groups",
  "groups",
  false,
  inputType,
  () => new GroupsCriterion()
);

export class GroupsCriterion extends ILabeledIdCriterion {
  constructor() {
    super(GroupsCriterionOption);
  }
}
