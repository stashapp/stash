import { ILabeledIdCriterion, ILabeledIdCriterionOption } from "./criterion";

const inputType = "galleries";

export const GalleriesCriterionOption = new ILabeledIdCriterionOption(
  "galleries",
  "galleries",
  true,
  inputType,
  () => new GalleriesCriterion()
);

export class GalleriesCriterion extends ILabeledIdCriterion {
  constructor() {
    super(GalleriesCriterionOption);
  }
}
