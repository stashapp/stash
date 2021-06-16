import { ILabeledIdCriterion, ILabeledIdCriterionOption } from "./criterion";

export const PerformersCriterionOption = new ILabeledIdCriterionOption(
  "performers",
  "performers",
  "performers",
  true
);

export class PerformersCriterion extends ILabeledIdCriterion {
  constructor() {
    super(PerformersCriterionOption);
  }
}
