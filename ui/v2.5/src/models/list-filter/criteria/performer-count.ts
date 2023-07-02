import { NumberCriterion, NumberCriterionOption } from "./criterion";

const performerCountCriterionOption = new NumberCriterionOption(
  "performer_count",
  "performer_count"
);

export class PerformerCountCriterion extends NumberCriterion {
  constructor() {
    super(performerCountCriterionOption);
  }
}
