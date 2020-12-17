import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class OrganizedCriterion extends Criterion {
  public type: CriterionType = "organized";
  public parameterName: string = "organized";
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
  public options: string[] = [true.toString(), false.toString()];
  public value: string = "";
}

export class OrganizedCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("organized");
  public value: CriterionType = "organized";
}
