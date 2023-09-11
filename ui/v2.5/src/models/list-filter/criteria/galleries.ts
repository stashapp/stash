import { ILabeledIdCriterion, ILabeledIdCriterionOption } from "./criterion";

const inputType = "galleries";

const galleriesCriterionOption = new ILabeledIdCriterionOption(
  "galleries",
  "galleries",
  true,
  inputType
);

export class GalleriesCriterion extends ILabeledIdCriterion {
  constructor() {
    super(galleriesCriterionOption);
  }
}
