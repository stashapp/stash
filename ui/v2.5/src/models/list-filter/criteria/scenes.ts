import {
  ModifierCriterionOption,
  ILabeledIdCriterion,
  ILabeledIdCriterionOption,
} from "./criterion";
import { CriterionModifier } from "src/core/generated-graphql";

const inputType = "scenes";

export const ScenesCriterionOption = new ILabeledIdCriterionOption(
  "scenes",
  "scenes",
  true,
  inputType,
  () => new ScenesCriterion()
);

export class ScenesCriterion extends ILabeledIdCriterion {
  constructor() {
    super(ScenesCriterionOption);
  }
}

const modifierOptions = [
  CriterionModifier.Includes,
  CriterionModifier.Excludes,
];

const defaultModifier = CriterionModifier.Includes;

export const MarkersScenesCriterionOption = new ModifierCriterionOption({
  messageID: "scenes",
  type: "scenes",
  modifierOptions,
  defaultModifier,
  inputType,
  makeCriterion: () => new MarkersScenesCriterion(),
});

export class MarkersScenesCriterion extends ILabeledIdCriterion {
  constructor() {
    super(MarkersScenesCriterionOption);
  }
}
