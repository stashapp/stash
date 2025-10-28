import { BooleanCriterion, BooleanCriterionOption } from "./criterion";

export const InteractiveCriterionOption = new BooleanCriterionOption(
  "interactive",
  "interactive",
  () => new InteractiveCriterion()
);

export class InteractiveCriterion extends BooleanCriterion {
  constructor() {
    super(InteractiveCriterionOption);
  }
}
