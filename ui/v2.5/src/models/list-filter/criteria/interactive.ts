import { BooleanCriterion, CriterionOption } from "./criterion";

export const InteractiveCriterionOption = new CriterionOption(
  "interactive",
  "interactive"
);

export class InteractiveCriterion extends BooleanCriterion {
  constructor() {
    super(InteractiveCriterionOption);
  }
}
