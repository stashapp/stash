import { ILabeledIdCriterion, ILabeledIdCriterionOption } from "./criterion";

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
