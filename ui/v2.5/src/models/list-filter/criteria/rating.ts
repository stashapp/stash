import { ConvertFromRatingFormat, ConvertToRatingFormat } from "../../../components/Scenes/SceneDetails/RatingSystem";
import { CriterionModifier, IntCriterionInput } from "../../../core/generated-graphql";
import { INumberValue } from "../types";
import { Criterion, CriterionOption } from "./criterion";

export class RatingCriterion extends Criterion<INumberValue> {
  public get value(): INumberValue {
    return this._value;
  }
  public set value(newValue: number | INumberValue) {
    // backwards compatibility - if this.value is a number, use that
    if (typeof newValue !== "object") {
      this._value = {
        value: newValue,
        value2: undefined,
      };
      this._value.value = ConvertFromRatingFormat(this._value.value);
    } else {
      this._value = newValue;
    }
  }

  protected toCriterionInput(): IntCriterionInput {
    return {
      modifier: this.modifier,
      value: this.value.value,
      value2: this.value.value2,
    };
  }

  public getLabelValue() {
    const { value, value2 } = this.value;
    if (
      this.modifier === CriterionModifier.Between ||
      this.modifier === CriterionModifier.NotBetween
    ) {
      return `${ConvertToRatingFormat(value)}, ${ConvertToRatingFormat(value2 ?? 100) ?? 0}`;
    } else {
      return `${ConvertToRatingFormat(value)}`;
    }
  }

  constructor(type: CriterionOption) {
    super(type, { value: 0, value2: undefined });
  }
}
