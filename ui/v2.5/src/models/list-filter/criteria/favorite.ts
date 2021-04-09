import {
  BooleanCriterion,
  Criterion,
  CriterionType,
  ICriterionOption,
} from "./criterion";

export class FavoriteCriterion extends BooleanCriterion {
  constructor() {
    super("favorite", "filter_favorites");
  }
}

export class FavoriteCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("favorite");
  public value: CriterionType = "favorite";
}
