import {
  CriterionModifier,
  PhashDistanceCriterionInput,
  PHashDuplicationCriterionInput,
} from "src/core/generated-graphql";
import { IPhashDistanceValue } from "../types";
import {
  BooleanCriterionOption,
  ModifierCriterion,
  ModifierCriterionOption,
  StringCriterion,
} from "./criterion";

export const PhashCriterionOption = new ModifierCriterionOption({
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

export class PhashCriterion extends ModifierCriterion<IPhashDistanceValue> {
  constructor() {
    super(PhashCriterionOption, { value: "", distance: 0 });
  }

  public cloneValues() {
    this.value = { ...this.value };
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
