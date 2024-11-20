import { IntlShape } from "react-intl";
import { Criterion, CriterionOption, ModifierCriterion } from "./criterion";
import {
  CriterionModifier,
  CustomFieldCriterionInput,
} from "src/core/generated-graphql";
import { cloneDeep } from "@apollo/client/utilities";

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
      valueString = (first.value![0] as string) ?? "";
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

  public toJSON(): string {
    const encodedCriterion = {
      type: this.criterionOption.type,
      value: this.value,
    };
    return JSON.stringify(encodedCriterion);
  }

  public setFromSavedCriterion(criterion: {
    type: string;
    value: CustomFieldCriterionInput[];
  }): void {
    const { value } = criterion;
    this.value = cloneDeep(value);
  }
}
