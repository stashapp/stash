import { ILabeledIdCriterion, ILabeledIdCriterionOption } from "./criterion";

const inputType = "groups";

export const MoviesCriterionOption = new ILabeledIdCriterionOption(
  "groups",
  "movies",
  false,
  inputType,
  () => new MoviesCriterion()
);

export class MoviesCriterion extends ILabeledIdCriterion {
  constructor() {
    super(MoviesCriterionOption);
  }
}
