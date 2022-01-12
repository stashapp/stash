import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionOption, StringCriterion } from "./criterion";

export const PhashCriterionOption = new CriterionOption({
  messageID: "media_info.phash",
  type: "phash",
  parameterName: "phash",
  inputType: "text",
  modifierOptions: [
    CriterionModifier.Equals,
    CriterionModifier.NotEquals,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
});

export class PhashCriterion extends StringCriterion {
  constructor() {
    super(PhashCriterionOption);
  }
}
