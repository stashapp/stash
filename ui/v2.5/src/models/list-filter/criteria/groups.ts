import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption, ILabeledIdCriterion } from "./criterion";
import { CriterionType } from "../types";

const inputType = "groups";

const modifierOptions = [
  CriterionModifier.Includes,
  CriterionModifier.Excludes,
  CriterionModifier.IsNull,
  CriterionModifier.NotNull,
];

class BaseGroupsCriterionOption extends CriterionOption {
  constructor(
    messageID: string,
    type: CriterionType,
  ) {
    super({
      messageID,
      type,
      modifierOptions,
      inputType,
      makeCriterion: () => new GroupsCriterion(this),
    });
  }
}

export const GroupsCriterionOption = new BaseGroupsCriterionOption(
  "groups",
  "groups",
);

export class GroupsCriterion extends ILabeledIdCriterion {}

// export const ContainingGroupsCriterionOption = new BaseGroupsCriterionOption(
//   "containing_groups",
//   "containing_groups",
// );

// export const ChildTagsCriterionOption = new BaseGroupsCriterionOption(
//   "sub_groups",
//   "sub_groups",
// );
