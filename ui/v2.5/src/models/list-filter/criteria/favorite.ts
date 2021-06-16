import { BooleanCriterion, BooleanCriterionOption } from "./criterion";

export const FavoriteCriterionOption = new BooleanCriterionOption(
  "favourite",
  "favorite",
  "filter_favorites"
);

export class FavoriteCriterion extends BooleanCriterion {
  constructor() {
    super(FavoriteCriterionOption);
  }
}
