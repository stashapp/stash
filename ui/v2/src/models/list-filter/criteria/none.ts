import {
  Criterion,
  CriterionModifier,
  CriterionType,
  ICriterionOption,
} from "./criterion";

export class NoneCriterion extends Criterion<any, any> {
  public type: CriterionType = "none";
  public parameterName: string = "";
  public modifier = CriterionModifier.Equals;
  public options: any;
  public value: any;
}

export class NoneCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("none");
  public value: CriterionType = "none";
}
