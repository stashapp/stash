import { BooleanCriterion, BooleanCriterionOption } from "./criterion";

export const CaptionedCriterionOption = new BooleanCriterionOption(
  "captioned",
  "captioned"
);

export class CaptionedCriterion extends BooleanCriterion {
  constructor() {
    super(CaptionedCriterionOption);
  }
}
