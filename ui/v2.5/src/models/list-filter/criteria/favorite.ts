import { BooleanCriterion, CriterionOption } from "./criterion";

export const FavoriteCriterionOption = new CriterionOption(
  "favourite",
  "favorite",
  "filter_favorites"
);

export class FavoriteCriterion extends BooleanCriterion {
  constructor() {
    super(FavoriteCriterionOption);
  }
}
