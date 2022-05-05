import { CriterionModifier } from "src/core/generated-graphql";
import { languageMap, valueToCode } from "src/utils/caption";
import { CriterionType } from "../types";
import { CriterionOption, StringCriterion } from "./criterion";

const languageStrings = Array.from(languageMap.values());

class CaptionsCriterionOptionType extends CriterionOption {
  constructor(value: CriterionType) {
    super({
      messageID: value,
      type: value,
      parameterName: value,
      modifierOptions: [
        CriterionModifier.Includes,
        CriterionModifier.Excludes,
        CriterionModifier.IsNull,
        CriterionModifier.NotNull,
      ],
      options: languageStrings,
    });
  }
}

export const CaptionsCriterionOption = new CaptionsCriterionOptionType(
  "captions"
);

export class CaptionCriterion extends StringCriterion {
  protected toCriterionInput() {
    const value = valueToCode(this.value);

    return {
      value,
      modifier: this.modifier,
    };
  }

  constructor() {
    super(CaptionsCriterionOption);
  }
}
