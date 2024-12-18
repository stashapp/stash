import {
  CriterionModifier,
  GenderCriterionInput,
  GenderEnum,
} from "src/core/generated-graphql";
import { genderStrings, stringToGender } from "src/utils/gender";
import {
  CriterionOption,
  ISavedCriterion,
  MultiStringCriterion,
} from "./criterion";

export const GenderCriterionOption = new CriterionOption({
  messageID: "gender",
  type: "gender",
  options: genderStrings,
  modifierOptions: [
    CriterionModifier.Includes,
    CriterionModifier.Excludes,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
  defaultModifier: CriterionModifier.Includes,
  makeCriterion: () => new GenderCriterion(),
});

export class GenderCriterion extends MultiStringCriterion {
  constructor(value: string[] = []) {
    super(GenderCriterionOption, value);
  }

  public toCriterionInput(): GenderCriterionInput {
    const value = this.value.map((v) => stringToGender(v)) as GenderEnum[];

    return {
      value_list: value,
      modifier: this.modifier,
    };
  }

  public setFromSavedCriterion(criterion: ISavedCriterion<string[]>) {
    // backwards compatibility - if the value is a string, convert it to an array
    if (typeof criterion.value === "string") {
      criterion = {
        ...criterion,
        value: [criterion.value],
      };
    }

    super.setFromSavedCriterion(criterion);
  }
}
