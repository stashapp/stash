import { BooleanCriterion, CriterionOption } from "./criterion";

export class FavoriteCriterionOption extends CriterionOption {
  constructor() {
    super("favourite", "favorite", "filter_favorites");
  }
}

export class FavoriteCriterion extends BooleanCriterion {
  constructor() {
    super(new FavoriteCriterionOption());
  }
}
