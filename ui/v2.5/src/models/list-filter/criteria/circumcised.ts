import { CircumcisionCriterionInput } from "src/core/generated-graphql";
import { circumcisedStrings, stringToCircumcised } from "src/utils/circumcised";
import { CriterionOption, StringCriterion } from "./criterion";

export const CircumcisedCriterionOption = new CriterionOption({
  messageID: "circumcised",
  type: "circumcised",
  options: circumcisedStrings,
});

export class CircumcisedCriterion extends StringCriterion {
  constructor() {
    super(CircumcisedCriterionOption);
  }

  protected toCriterionInput(): CircumcisionCriterionInput {
    return {
      value: stringToCircumcised(this.value),
      modifier: this.modifier,
    };
  }
}
