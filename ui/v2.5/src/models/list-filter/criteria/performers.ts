/* eslint @typescript-eslint/no-unused-vars: ["error", { "argsIgnorePattern": "^_" }] */
import { IntlShape } from "react-intl";
import {
  CriterionModifier,
  MultiCriterionInput,
} from "src/core/generated-graphql";
import { ILabeledId, ILabeledValueListValue } from "../types";
import {
  ModifierCriterion,
  ModifierCriterionOption,
  ISavedCriterion,
} from "./criterion";

const modifierOptions = [
  CriterionModifier.IncludesAll,
  CriterionModifier.Includes,
  CriterionModifier.Equals,
  CriterionModifier.IsNull,
  CriterionModifier.NotNull,
];

const defaultModifier = CriterionModifier.IncludesAll;

const inputType = "performers";

export const PerformersCriterionOption = new ModifierCriterionOption({
  messageID: "performers",
  type: "performers",
  modifierOptions,
  defaultModifier,
  inputType,
  makeCriterion: () => new PerformersCriterion(),
});

export class PerformersCriterion extends ModifierCriterion<ILabeledValueListValue> {
  constructor() {
    super(PerformersCriterionOption, { items: [], excluded: [] });
  }

  public cloneValues() {
    this.value = {
      ...this.value,
      items: this.value.items.map((v) => ({ ...v })),
      excluded: this.value.excluded.map((v) => ({ ...v })),
    };
  }

  override get modifier(): CriterionModifier {
    return this._modifier;
  }
  override set modifier(value: CriterionModifier) {
    this._modifier = value;

    // excluded only makes sense for includes and includes all
    // reset it for other modifiers
    if (
      value !== CriterionModifier.Includes &&
      value !== CriterionModifier.IncludesAll
    ) {
      this.value.excluded = [];
    }
  }

  public setFromSavedCriterion(
    criterion: ISavedCriterion<ILabeledId[] | ILabeledValueListValue>
  ) {
    const { modifier, value } = criterion;

    // #3619 - the format of performer value was changed from an array
    // to an object. Check for both formats.
    if (Array.isArray(value)) {
      this.value = { items: value, excluded: [] };
    } else if (value !== undefined) {
      this.value = {
        items: value.items || [],
        excluded: value.excluded || [],
      };
    }

    // if the previous modifier was excludes, replace it with the equivalent includes criterion
    // this is what is done on the backend
    if (modifier === CriterionModifier.Excludes) {
      this.modifier = CriterionModifier.Includes;
      this.value.excluded = [...this.value.excluded, ...this.value.items];
      this.value.items = [];
    } else {
      this.modifier = modifier;
    }
  }

  protected getLabelValue(_intl: IntlShape): string {
    return this.value.items.map((v) => v.label).join(", ");
  }

  public toCriterionInput(): MultiCriterionInput {
    let excludes: string[] = [];
    if (this.value.excluded) {
      excludes = this.value.excluded.map((v) => v.id);
    }
    return {
      value: this.value.items.map((v) => v.id),
      excludes: excludes,
      modifier: this.modifier,
    };
  }

  public isValid(): boolean {
    if (
      this.modifier === CriterionModifier.IsNull ||
      this.modifier === CriterionModifier.NotNull
    ) {
      return true;
    }

    return (
      this.value.items.length > 0 ||
      (this.value.excluded && this.value.excluded.length > 0)
    );
  }

  public getLabel(intl: IntlShape): string {
    let id = "criterion_modifier.format_string";
    let modifierString = ModifierCriterion.getModifierLabel(
      intl,
      this.modifier
    );
    let valueString = "";
    let excludedString = "";

    if (
      this.modifier !== CriterionModifier.IsNull &&
      this.modifier !== CriterionModifier.NotNull
    ) {
      valueString = this.value.items.map((v) => v.label).join(", ");

      if (this.value.excluded && this.value.excluded.length > 0) {
        if (this.value.items.length === 0) {
          modifierString = ModifierCriterion.getModifierLabel(
            intl,
            CriterionModifier.Excludes
          );
          valueString = this.value.excluded.map((v) => v.label).join(", ");
        } else {
          id = "criterion_modifier.format_string_excludes";
          excludedString = this.value.excluded.map((v) => v.label).join(", ");
        }
      }
    }

    return intl.formatMessage(
      { id },
      {
        criterion: intl.formatMessage({ id: this.criterionOption.messageID }),
        modifierString,
        valueString,
        excludedString,
      }
    );
  }
}
