import { ILabeledIdCriterion, ILabeledIdCriterionOption } from "./criterion";

export const MoviesCriterionOption = new ILabeledIdCriterionOption(
  "movies",
  "movies",
  false
);

export class MoviesCriterion extends ILabeledIdCriterion {
  constructor() {
    super(MoviesCriterionOption);
  }
}
