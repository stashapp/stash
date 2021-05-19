import { BooleanCriterion, CriterionOption } from "./criterion";

export const OrganizedCriterionOption = new CriterionOption(
  "organized",
  "organized"
);

export class OrganizedCriterion extends BooleanCriterion {
  constructor() {
    super(OrganizedCriterionOption);
  }
}
