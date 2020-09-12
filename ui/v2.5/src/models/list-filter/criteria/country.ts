import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionType, ICriterionOption } from "./criterion";

export class CountryCriterion extends Criterion {
  public type: CriterionType = "country";
  public parameterName: string = "performers";
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];
  public options: string[] = [true.toString(), false.toString()];
  public value: string = "";
}

export class CountryCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("performers");
  public value: CriterionType = "country";
}
