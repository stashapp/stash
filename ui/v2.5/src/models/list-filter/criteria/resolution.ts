import {
  ResolutionCriterionInput,
  CriterionModifier,
} from "src/core/generated-graphql";
import { stringToResolution, resolutionStrings } from "src/utils/resolution";
import { CriterionType } from "../types";
import {
  ModifierCriterion,
  ModifierCriterionOption,
  CriterionValue,
  StringCriterion,
} from "./criterion";

class BaseResolutionCriterionOption extends ModifierCriterionOption {
  constructor(
    value: CriterionType,
    makeCriterion: () => ModifierCriterion<CriterionValue>
  ) {
    super({
      messageID: value,
      type: value,
      modifierOptions: [
        CriterionModifier.Equals,
        CriterionModifier.NotEquals,
        CriterionModifier.GreaterThan,
        CriterionModifier.LessThan,
      ],
      options: resolutionStrings,
      makeCriterion,
    });
  }
}

class BaseResolutionCriterion extends StringCriterion {
  public toCriterionInput(): ResolutionCriterionInput | undefined {
    const value = stringToResolution(this.value);

    if (value !== undefined) {
      return {
        value,
        modifier: this.modifier,
      };
    }
  }
}

export const ResolutionCriterionOption = new BaseResolutionCriterionOption(
  "resolution",
  () => new ResolutionCriterion()
);

export class ResolutionCriterion extends BaseResolutionCriterion {
  constructor() {
    super(ResolutionCriterionOption);
  }
}

export const AverageResolutionCriterionOption =
  new BaseResolutionCriterionOption(
    "average_resolution",
    () => new AverageResolutionCriterion()
  );

export class AverageResolutionCriterion extends BaseResolutionCriterion {
  constructor() {
    super(AverageResolutionCriterionOption);
  }
}
