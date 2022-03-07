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

export const PerformerFavoriteCriterionOption = new BooleanCriterionOption(
  "performer_favorite",
  "performer_favorite",
  "performer_favorite"
);

export class PerformerFavoriteCriterion extends BooleanCriterion {
  constructor() {
    super(PerformerFavoriteCriterionOption);
  }
}
