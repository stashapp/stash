import { BooleanCriterion, BooleanCriterionOption } from "./criterion";

export const InteractiveCriterionOption = new BooleanCriterionOption(
  "interactive",
  "interactive"
);

export class InteractiveCriterion extends BooleanCriterion {
  constructor() {
    super(InteractiveCriterionOption);
  }
}
