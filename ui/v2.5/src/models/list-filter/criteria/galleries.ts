import { CriterionOption, ILabeledIdCriterion } from "./criterion";

const galleriesCriterionOption = new CriterionOption("galleries", "galleries");

export class GalleriesCriterion extends ILabeledIdCriterion {
  constructor() {
    super(galleriesCriterionOption, true);
  }
}
