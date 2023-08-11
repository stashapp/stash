import { ILabeledIdCriterion, ILabeledIdCriterionOption } from "./criterion";

const galleriesCriterionOption = new ILabeledIdCriterionOption(
  "galleries",
  "galleries",
  true
);

export class GalleriesCriterion extends ILabeledIdCriterion {
  constructor() {
    super(galleriesCriterionOption);
  }
}
