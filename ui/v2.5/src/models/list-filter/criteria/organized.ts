import { BooleanCriterion, BooleanCriterionOption } from "./criterion";

export const OrganizedCriterionOption = new BooleanCriterionOption(
  "organized",
  "organized"
);

export class OrganizedCriterion extends BooleanCriterion {
  constructor() {
    super(OrganizedCriterionOption);
  }
}
