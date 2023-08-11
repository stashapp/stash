import { ILabeledIdCriterion, ILabeledIdCriterionOption } from "./criterion";

const inputType = "movies";

export const MoviesCriterionOption = new ILabeledIdCriterionOption(
  "movies",
  "movies",
  false,
  inputType
);

export class MoviesCriterion extends ILabeledIdCriterion {
  constructor() {
    super(MoviesCriterionOption);
  }
}
