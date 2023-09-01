import { CriterionModifier } from "src/core/generated-graphql";
import {
  CriterionOption,
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

export const StudiosCriterionOption = new CriterionOption({
  messageID: "studios",
  type: "studios",
  modifierOptions,
  defaultModifier,
  makeCriterion: () => new StudiosCriterion(),
  inputType,
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
  inputType
);
export class ParentStudiosCriterion extends ILabeledIdCriterion {
  constructor() {
    super(ParentStudiosCriterionOption);
  }
}
