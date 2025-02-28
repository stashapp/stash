import { orientationStrings, stringToOrientation } from "src/utils/orientation";
import { CriterionType } from "../types";
import { ModifierCriterionOption, MultiStringCriterion } from "./criterion";
import {
  OrientationCriterionInput,
  OrientationEnum,
} from "src/core/generated-graphql";

export class OrientationCriterion extends MultiStringCriterion {
  public toCriterionInput(): OrientationCriterionInput {
    return {
      value: this.value
        .map((v) => stringToOrientation(v))
        .filter((v) => v) as OrientationEnum[],
    };
  }
}

class BaseOrientationCriterionOption extends ModifierCriterionOption {
  constructor(value: CriterionType) {
    super({
      messageID: value,
      type: value,
      options: orientationStrings,
      makeCriterion: () => new OrientationCriterion(this),
    });
  }
}

export const OrientationCriterionOption = new BaseOrientationCriterionOption(
  "orientation"
);
