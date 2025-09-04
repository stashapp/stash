import { IntlShape } from "react-intl";
import { Criterion, CriterionOption, ModifierCriterion } from "./criterion";
import {
  CriterionModifier,
  CustomFieldCriterionInput,
} from "src/core/generated-graphql";
import { cloneDeep } from "@apollo/client/utilities";

function valueToString(value: unknown[] | undefined | null) {
  if (!value) return "";
  return value.map((v) => v as string).join(", ");
}

export const CustomFieldsCriterionOption = new CriterionOption({
  type: "custom_fields",
  messageID: "custom_fields.title",
  makeCriterion: () => new CustomFieldsCriterion(),
});

export class CustomFieldsCriterion extends Criterion {
  public value: CustomFieldCriterionInput[] = [];

  constructor() {
    super(CustomFieldsCriterionOption);
  }

  public isValid(): boolean {
    return this.value.length > 0;
  }

  public applyToCriterionInput(input: Record<string, unknown>): void {
    input.custom_fields = cloneDeep(this.value);
  }

  public applyToSavedCriterion(input: Record<string, unknown>): void {
    input.custom_fields = cloneDeep(this.value);
  }

  public getLabel(intl: IntlShape): string {
    // show first criterion
    if (this.value.length === 0) {
      return "";
    }

    const first = this.value[0];
    let messageID;
    let valueString = "";

    if (
      first.modifier !== CriterionModifier.IsNull &&
      first.modifier !== CriterionModifier.NotNull &&
      (first.value?.length ?? 0) > 0
    ) {
      valueString = valueToString(first.value);
    }

    const modifierString = ModifierCriterion.getModifierLabel(
      intl,
      first.modifier
    );
    const opts = {
      criterion: first.field,
      modifierString,
      valueString,
      others: "",
    };

    if (this.value.length === 1) {
      messageID = "custom_fields.criteria_format_string";
    } else {
      messageID = "custom_fields.criteria_format_string_others";
      opts.others = (this.value.length - 1).toString();
    }

    return intl.formatMessage({ id: messageID }, opts);
  }

  public getValueLabel(intl: IntlShape, v: CustomFieldCriterionInput): string {
    let valueString = "";

    if (
      v.modifier !== CriterionModifier.IsNull &&
      v.modifier !== CriterionModifier.NotNull &&
      (v.value?.length ?? 0) > 0
    ) {
      valueString = valueToString(v.value);
    }

    const modifierString = ModifierCriterion.getModifierLabel(intl, v.modifier);
    const opts = {
      criterion: v.field,
      modifierString,
      valueString,
    };

    return intl.formatMessage(
      { id: "custom_fields.criteria_format_string" },
      opts
    );
  }

  public toQueryParams(): Record<string, unknown> {
    const encodedCriterion = {
      type: this.criterionOption.type,
      value: this.value,
    };
    return encodedCriterion;
  }

  public fromDecodedParams(i: unknown): void {
    const criterion = i as { value: CustomFieldCriterionInput[] };
    this.value = cloneDeep(criterion.value);
  }

  public setFromSavedCriterion(input: CustomFieldCriterionInput[]): void {
    this.value = cloneDeep(input);
  }
}
