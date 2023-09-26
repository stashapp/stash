import { BooleanCriterion, BooleanCriterionOption } from "./criterion";

export const FavoriteCriterionOption = new BooleanCriterionOption(
  "favourite",
  "filter_favorites",
  () => new FavoriteCriterion()
);

export class FavoriteCriterion extends BooleanCriterion {
  constructor() {
    super(FavoriteCriterionOption);
  }
}

export const PerformerFavoriteCriterionOption = new BooleanCriterionOption(
  "performer_favorite",
  "performer_favorite",
  () => new PerformerFavoriteCriterion()
);

export class PerformerFavoriteCriterion extends BooleanCriterion {
  constructor() {
    super(PerformerFavoriteCriterionOption);
  }
}
