import { CriterionModifier } from "src/core/generated-graphql";
import {
  ModifierCriterionOption,
  IHierarchicalLabeledIdCriterion,
} from "./criterion";
import { CriterionType } from "../types";

const inputType = "groups";

const modifierOptions = [
  CriterionModifier.Includes,
  CriterionModifier.Excludes,
  CriterionModifier.IsNull,
  CriterionModifier.NotNull,
];

const defaultModifier = CriterionModifier.Includes;

class BaseGroupsCriterionOption extends ModifierCriterionOption {
  constructor(messageID: string, type: CriterionType) {
    super({
      messageID,
      type,
      modifierOptions,
      defaultModifier,
      inputType,
      makeCriterion: () => new GroupsCriterion(this),
    });
  }
}

export const GroupsCriterionOption = new BaseGroupsCriterionOption(
  "groups",
  "groups"
);

export class GroupsCriterion extends IHierarchicalLabeledIdCriterion {}

export const ContainingGroupsCriterionOption = new BaseGroupsCriterionOption(
  "containing_groups",
  "containing_groups"
);

export const SubGroupsCriterionOption = new BaseGroupsCriterionOption(
  "sub_groups",
  "sub_groups"
);

// redirects to GroupsCriterion
export const LegacyMoviesCriterionOption = new ModifierCriterionOption({
  messageID: "groups",
  type: "movies",
  modifierOptions,
  defaultModifier,
  inputType,
  hidden: true,
  makeCriterion: () => new GroupsCriterion(GroupsCriterionOption),
});
