import { BooleanCriterion, CriterionOption } from "./criterion";

export class OrganizedCriterionOption extends CriterionOption {
  constructor() {
    super("organized", "organized");
  }
}

export class OrganizedCriterion extends BooleanCriterion {
  constructor() {
    super(new OrganizedCriterionOption());
  }
}
