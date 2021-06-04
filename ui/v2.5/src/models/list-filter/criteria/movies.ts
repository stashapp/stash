import { CriterionOption, ILabeledIdCriterion } from "./criterion";

export const MoviesCriterionOption = new CriterionOption("movies", "movies");

export class MoviesCriterion extends ILabeledIdCriterion {
  constructor() {
    super(MoviesCriterionOption, false);
  }
}
