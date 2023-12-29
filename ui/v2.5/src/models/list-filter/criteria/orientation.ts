import { orientationStrings } from "src/utils/orientation";
import { CriterionType } from "../types";
import { CriterionOption, StringCriterion } from "./criterion";
import { CriterionModifier } from "../../../core/generated-graphql";

export class OrientationCriterion extends StringCriterion {
  protected toCriterionInput(): string {
    return this.value;
  }
}

class BaseOrientationCriterionOption extends CriterionOption {
  constructor(value: CriterionType) {
    super({
      messageID: value,
      type: value,
      modifierOptions: [],
      options: orientationStrings,
      defaultModifier: CriterionModifier.Equals,
      makeCriterion: () => new OrientationCriterion(this),
    });
  }
}

export const OrientationCriterionOption = new BaseOrientationCriterionOption(
  "orientation"
);
