import {
  Criterion,
  CriterionModifier,
  CriterionType,
  ICriterionOption,
} from "./criterion";

export class FavoriteCriterion extends Criterion<string, string> {
  public type: CriterionType = "favorite";
  public parameterName: string = "filter_favorites";
  public modifier = CriterionModifier.Equals;
  public options: string[] = [true.toString(), false.toString()];
  public value: string = "";
}

export class FavoriteCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("favorite");
  public value: CriterionType = "favorite";
}
