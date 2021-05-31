import { BooleanCriterion, CriterionOption } from "./criterion";

export const InteractiveCriterionOption = new CriterionOption(
  "organized",
  "organized"
);

export class InteractiveCriterion extends BooleanCriterion {
  constructor() {
    super(InteractiveCriterionOption);
  }
}
