/* eslint @typescript-eslint/no-unused-vars: ["error", { "argsIgnorePattern": "^_" }] */
import { IntlShape } from "react-intl";
import {
  CriterionModifier,
  MultiCriterionInput,
} from "src/core/generated-graphql";
import { ILabeledValueListValue } from "../types";
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

  public getLabelValue(_intl: IntlShape): string {
    return this.value.items.map((v) => v.label).join(", ");
  }

  protected toCriterionInput(): MultiCriterionInput {
    return {
      value: this.value.items.map((v) => v.id),
      modifier: this.modifier,
    };
  }

  public isValid(): boolean {
    if (this.modifier === CriterionModifier.Equals) {
      return true;
    }

    return this.value.items.length > 0 || this.value.excluded.length > 0;
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

    if (this.value.excluded.length > 0) {
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
