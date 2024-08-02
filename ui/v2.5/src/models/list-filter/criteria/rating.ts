import {
  convertFromRatingFormat,
  convertToRatingFormat,
  defaultRatingSystemOptions,
  RatingSystemOptions,
} from "src/utils/rating";
import {
  ConfigDataFragment,
  CriterionModifier,
  IntCriterionInput,
} from "src/core/generated-graphql";
import { INumberValue } from "../types";
import { Criterion, CriterionOption } from "./criterion";

const modifierOptions = [
  CriterionModifier.Equals,
  CriterionModifier.NotEquals,
  CriterionModifier.GreaterThan,
  CriterionModifier.LessThan,
  CriterionModifier.Between,
  CriterionModifier.NotBetween,
  CriterionModifier.IsNull,
  CriterionModifier.NotNull,
];

function getRatingSystemOptions(config?: ConfigDataFragment) {
  return config?.ui.ratingSystemOptions ?? defaultRatingSystemOptions;
}

export const RatingCriterionOption = new CriterionOption({
  messageID: "rating",
  type: "rating100",
  modifierOptions,
  defaultModifier: CriterionModifier.Equals,
  makeCriterion: (o, config) =>
    new RatingCriterion(getRatingSystemOptions(config)),
  inputType: "number",
});

export class RatingCriterion extends Criterion<INumberValue> {
  ratingSystem: RatingSystemOptions;

  constructor(ratingSystem: RatingSystemOptions) {
    super(RatingCriterionOption, { value: 0, value2: undefined });
    this.ratingSystem = ratingSystem;
  }

  public clone() {
    const newCriterion = new RatingCriterion(this.ratingSystem);
    newCriterion.modifier = this.modifier;
    newCriterion.value = { ...this.value };
    return newCriterion;
  }

  public get value(): INumberValue {
    return this._value;
  }
  public set value(newValue: number | INumberValue) {
    // backwards compatibility - if this.value is a number, use that
    if (typeof newValue !== "object") {
      this._value = {
        value: convertFromRatingFormat(newValue, this.ratingSystem.type),
        value2: undefined,
      };
    } else {
      this._value = newValue;
    }
  }

  public toCriterionInput(): IntCriterionInput {
    return {
      modifier: this.modifier,
      value: this.value.value ?? 0,
      value2: this.value.value2,
    };
  }

  protected getLabelValue() {
    const { value, value2 } = this.value;
    if (
      this.modifier === CriterionModifier.Between ||
      this.modifier === CriterionModifier.NotBetween
    ) {
      return `${convertToRatingFormat(value, this.ratingSystem) ?? 0}, ${
        convertToRatingFormat(value2, this.ratingSystem) ?? 0
      }`;
    } else {
      return `${convertToRatingFormat(value, this.ratingSystem) ?? 0}`;
    }
  }
}
