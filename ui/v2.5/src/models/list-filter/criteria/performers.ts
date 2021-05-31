import { CriterionOption, ILabeledIdCriterion } from "./criterion";

export const PerformersCriterionOption = new CriterionOption(
  "performers",
  "performers"
);

export class PerformersCriterion extends ILabeledIdCriterion {
  constructor() {
    super(PerformersCriterionOption, true);
  }
}
