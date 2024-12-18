/* eslint @typescript-eslint/no-unused-vars: ["error", { "argsIgnorePattern": "^_" }] */
import { IntlShape } from "react-intl";
import {
  CriterionModifier,
  StashIdCriterionInput,
} from "src/core/generated-graphql";
import { IStashIDValue } from "../types";
import { Criterion, CriterionOption } from "./criterion";

export const StashIDCriterionOption = new CriterionOption({
  messageID: "stash_id",
  type: "stash_id_endpoint",
  modifierOptions: [
    CriterionModifier.Equals,
    CriterionModifier.NotEquals,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
  makeCriterion: () => new StashIDCriterion(),
});

export class StashIDCriterion extends Criterion<IStashIDValue> {
  constructor() {
    super(StashIDCriterionOption, {
      endpoint: "",
      stashID: "",
    });
  }

  public cloneValues() {
    this.value = { ...this.value };
  }

  public get value(): IStashIDValue {
    return this._value;
  }

  public set value(newValue: string | IStashIDValue) {
    // backwards compatibility - if this.value is a string, use that as stash_id
    if (typeof newValue !== "object") {
      this._value = {
        endpoint: "",
        stashID: newValue,
      };
    } else {
      this._value = newValue;
    }
  }

  public toCriterionInput(): StashIdCriterionInput {
    return {
      endpoint: this.value.endpoint,
      stash_id: this.value.stashID,
      modifier: this.modifier,
    };
  }

  public getLabel(intl: IntlShape): string {
    const modifierString = Criterion.getModifierLabel(intl, this.modifier);
    let valueString = "";

    if (
      this.modifier !== CriterionModifier.IsNull &&
      this.modifier !== CriterionModifier.NotNull
    ) {
      valueString = this.getLabelValue(intl);
    } else if (this.value.endpoint) {
      valueString = "(" + this.value.endpoint + ")";
    }

    return intl.formatMessage(
      { id: "criterion_modifier.format_string" },
      {
        criterion: intl.formatMessage({ id: this.criterionOption.messageID }),
        modifierString,
        valueString,
      }
    );
  }

  protected getLabelValue(_intl: IntlShape) {
    let ret = this.value.stashID;
    if (this.value.endpoint) {
      ret += " (" + this.value.endpoint + ")";
    }

    return ret;
  }

  public toJSON() {
    let encodedCriterion;
    if (
      (this.modifier === CriterionModifier.IsNull ||
        this.modifier === CriterionModifier.NotNull) &&
      !this.value.endpoint
    ) {
      encodedCriterion = {
        type: this.criterionOption.type,
        modifier: this.modifier,
      };
    } else {
      encodedCriterion = {
        type: this.criterionOption.type,
        value: this.value,
        modifier: this.modifier,
      };
    }
    return JSON.stringify(encodedCriterion);
  }

  public isValid(): boolean {
    return (
      this.modifier === CriterionModifier.IsNull ||
      this.modifier === CriterionModifier.NotNull ||
      this.value.stashID.length > 0
    );
  }
}
