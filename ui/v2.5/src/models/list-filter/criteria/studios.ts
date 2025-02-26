import { CriterionModifier } from "src/core/generated-graphql";
import {
  ModifierCriterionOption,
  IHierarchicalLabeledIdCriterion,
  ILabeledIdCriterion,
  ILabeledIdCriterionOption,
} from "./criterion";

const modifierOptions = [
  CriterionModifier.Includes,
  CriterionModifier.IsNull,
  CriterionModifier.NotNull,
];

const defaultModifier = CriterionModifier.Includes;
const inputType = "studios";

export const StudiosCriterionOption = new ModifierCriterionOption({
  messageID: "studios",
  type: "studios",
  modifierOptions,
  defaultModifier,
  inputType,
  makeCriterion: () => new StudiosCriterion(),
});

export class StudiosCriterion extends IHierarchicalLabeledIdCriterion {
  constructor() {
    super(StudiosCriterionOption);
  }
}

export const ParentStudiosCriterionOption = new ILabeledIdCriterionOption(
  "parent_studios",
  "parents",
  false,
  inputType,
  () => new ParentStudiosCriterion()
);

export class ParentStudiosCriterion extends ILabeledIdCriterion {
  constructor() {
    super(ParentStudiosCriterionOption);
  }
}
