import {
  ResolutionCriterionInput,
  CriterionModifier,
} from "src/core/generated-graphql";
import { stringToResolution, resolutionStrings } from "src/utils/resolution";
import { CriterionType } from "../types";
import {
  Criterion,
  CriterionOption,
  CriterionValue,
  StringCriterion,
} from "./criterion";

abstract class AbstractResolutionCriterion extends StringCriterion {
  protected toCriterionInput(): ResolutionCriterionInput | undefined {
    const value = stringToResolution(this.value);

    if (value !== undefined) {
      return {
        value,
        modifier: this.modifier,
      };
    }
  }
}

class ResolutionCriterionOptionType extends CriterionOption {
  constructor(
    value: CriterionType,
    makeCriterion: () => Criterion<CriterionValue>
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

export const ResolutionCriterionOption = new ResolutionCriterionOptionType(
  "resolution",
  () => new ResolutionCriterion()
);
export class ResolutionCriterion extends AbstractResolutionCriterion {
  constructor() {
    super(ResolutionCriterionOption);
  }
}

export const AverageResolutionCriterionOption =
  new ResolutionCriterionOptionType(
    "average_resolution",
    () => new AverageResolutionCriterion()
  );

export class AverageResolutionCriterion extends AbstractResolutionCriterion {
  constructor() {
    super(AverageResolutionCriterionOption);
  }
}
