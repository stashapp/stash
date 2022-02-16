import { CriterionModifier } from "src/core/generated-graphql";
import {
  BooleanCriterionOption,
  CriterionOption,
  PhashDuplicateCriterion,
  StringCriterion,
} from "./criterion";

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

export const DuplicatedCriterionOption = new BooleanCriterionOption(
  "duplicated_phash",
  "duplicated",
  "duplicated"
);

export class DuplicatedCriterion extends PhashDuplicateCriterion {
  constructor() {
    super(DuplicatedCriterionOption);
  }
}
