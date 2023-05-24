/* eslint @typescript-eslint/no-unused-vars: ["error", { "argsIgnorePattern": "^_" }] */
import { IntlShape } from "react-intl";
import {
  CriterionModifier,
  MultiCriterionInput,
} from "src/core/generated-graphql";
import { ILabeledId, ILabeledValueListValue } from "../types";
import { Criterion, CriterionOption } from "./criterion";

const modifierOptions = [
  CriterionModifier.IncludesAll,
  CriterionModifier.Includes,
  CriterionModifier.Equals,
];

const defaultModifier = CriterionModifier.IncludesAll;

export const PerformersCriterionOption = new CriterionOption({
  messageID: "performers",
  type: "performers",
  parameterName: "performers",
  modifierOptions,
  defaultModifier,
});

export class PerformersCriterion extends Criterion<ILabeledValueListValue> {
  constructor() {
    super(PerformersCriterionOption, { items: [], excluded: [] });
  }

  public setValueFromQueryString(v: ILabeledId[] | ILabeledValueListValue) {
    // #3619 - the format of performer value was changed from an array
    // to an object. Check for both formats.
    if (Array.isArray(v)) {
      this.value = { items: v, excluded: [] };
    } else {
      this.value = {
        items: v.items || [],
        excluded: v.excluded || [],
      };
    }
  }

  public getLabelValue(_intl: IntlShape): string {
    return this.value.items.map((v) => v.label).join(", ");
  }

  protected toCriterionInput(): MultiCriterionInput {
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
      this.modifier === CriterionModifier.NotNull ||
      this.modifier === CriterionModifier.Equals
    ) {
      return true;
    }

    return (
      this.value.items.length > 0 ||
      (this.value.excluded && this.value.excluded.length > 0)
    );
  }

  public getLabel(intl: IntlShape): string {
    const modifierString = Criterion.getModifierLabel(intl, this.modifier);
    let valueString = "";

    if (
      this.modifier !== CriterionModifier.IsNull &&
      this.modifier !== CriterionModifier.NotNull
    ) {
      valueString = this.value.items.map((v) => v.label).join(", ");
    }

    let id = "criterion_modifier.format_string";
    let excludedString = "";

    if (this.value.excluded && this.value.excluded.length > 0) {
      id = "criterion_modifier.format_string_excludes";
      excludedString = this.value.excluded.map((v) => v.label).join(", ");
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
