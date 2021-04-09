import {
  BooleanCriterion,
  Criterion,
  CriterionType,
  ICriterionOption,
} from "./criterion";

export class OrganizedCriterion extends BooleanCriterion {
  constructor() {
    super("organized", "organized");
  }
}

export class OrganizedCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("organized");
  public value: CriterionType = "organized";
}
