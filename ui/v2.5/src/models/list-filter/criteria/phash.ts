import {
  CriterionModifier,
  PhashDistanceCriterionInput,
  PHashDuplicationCriterionInput,
} from "src/core/generated-graphql";
import { IPhashDistanceValue } from "../types";
import {
  BooleanCriterionOption,
  Criterion,
  CriterionOption,
  StringCriterion,
} from "./criterion";

export const PhashCriterionOption = new CriterionOption({
  messageID: "media_info.phash",
  type: "phash_distance",
  inputType: "text",
  modifierOptions: [
    CriterionModifier.Equals,
    CriterionModifier.NotEquals,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
  makeCriterion: () => new PhashCriterion(),
});

export class PhashCriterion extends Criterion<IPhashDistanceValue> {
  constructor() {
    super(PhashCriterionOption, { value: "", distance: 0 });
  }

  public clone() {
    const newCriterion = new PhashCriterion();
    newCriterion.modifier = this.modifier;
    newCriterion.value = { ...this.value };
    return newCriterion;
  }

  protected getLabelValue() {
    const { value, distance } = this.value;
    if (
      (this.modifier === CriterionModifier.Equals ||
        this.modifier === CriterionModifier.NotEquals) &&
      distance
    ) {
      return `${value} (${distance})`;
    } else {
      return `${value}`;
    }
  }

  public toCriterionInput(): PhashDistanceCriterionInput {
    return {
      value: this.value.value,
      modifier: this.modifier,
      distance: this.value.distance,
    };
  }
}

export const DuplicatedCriterionOption = new BooleanCriterionOption(
  "duplicated_phash",
  "duplicated",
  () => new DuplicatedCriterion()
);

export class DuplicatedCriterion extends StringCriterion {
  constructor() {
    super(DuplicatedCriterionOption);
  }

  public toCriterionInput(): PHashDuplicationCriterionInput {
    return {
      duplicated: this.value === "true",
    };
  }
}
