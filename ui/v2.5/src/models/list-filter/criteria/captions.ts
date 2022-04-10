import { StringCriterion, StringCriterionOption } from "./criterion";

export const CaptionsCriterionOption = new StringCriterionOption(
  "captions",
  "captions"
);

export class CaptionsCriterion extends StringCriterion {
  constructor() {
    super(CaptionsCriterionOption);
  }
}
