import { CriterionOption, StringCriterion } from "./criterion";

export const CaptionedCriterionOption = new CriterionOption({
  messageID: "captioned",
  type: "captioned",
  parameterName: "captioned",
  options: [true.toString(), false.toString()],
});

export class CaptionedCriterion extends StringCriterion {
  constructor() {
    super(CaptionedCriterionOption);
  }

  protected toCriterionInput(): string {
    return this.value;
  }
}