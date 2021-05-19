/* eslint-disable consistent-return */

import { IntlShape } from "react-intl";
import {
  CriterionModifier,
  MultiCriterionInput,
} from "src/core/generated-graphql";
import DurationUtils from "src/utils/duration";
import {
  CriterionType,
  ILabeledId,
  ILabeledValue,
  IOptionType,
} from "../types";

type Option = string | number | IOptionType;
export type CriterionValue = string | number | ILabeledId[];

// V = criterion value type
// I = criterion input type
export abstract class Criterion<V extends CriterionValue> {
  public static getModifierOption(
    modifier: CriterionModifier = CriterionModifier.Equals
  ): ILabeledValue {
    switch (modifier) {
      case CriterionModifier.Equals:
        return { value: CriterionModifier.Equals, label: "Equals" };
      case CriterionModifier.NotEquals:
        return { value: CriterionModifier.NotEquals, label: "Not Equals" };
      case CriterionModifier.GreaterThan:
        return { value: CriterionModifier.GreaterThan, label: "Greater Than" };
      case CriterionModifier.LessThan:
        return { value: CriterionModifier.LessThan, label: "Less Than" };
      case CriterionModifier.IsNull:
        return { value: CriterionModifier.IsNull, label: "Is NULL" };
      case CriterionModifier.NotNull:
        return { value: CriterionModifier.NotNull, label: "Not NULL" };
      case CriterionModifier.IncludesAll:
        return { value: CriterionModifier.IncludesAll, label: "Includes All" };
      case CriterionModifier.Includes:
        return { value: CriterionModifier.Includes, label: "Includes" };
      case CriterionModifier.Excludes:
        return { value: CriterionModifier.Excludes, label: "Excludes" };
      case CriterionModifier.MatchesRegex:
        return {
          value: CriterionModifier.MatchesRegex,
          label: "Matches Regex",
        };
      case CriterionModifier.NotMatchesRegex:
        return {
          value: CriterionModifier.NotMatchesRegex,
          label: "Not Matches Regex",
        };
    }
  }

  public criterionOption: CriterionOption;
  public abstract modifier: CriterionModifier;
  public abstract modifierOptions: ILabeledValue[];
  public abstract options: Option[] | undefined;
  public abstract value: V;
  public inputType: "number" | "text" | undefined;

  public abstract getLabelValue(): string;

  constructor(type: CriterionOption) {
    this.criterionOption = type;
  }

  public getLabel(intl: IntlShape): string {
    let modifierString: string;
    switch (this.modifier) {
      case CriterionModifier.Equals:
        modifierString = "is";
        break;
      case CriterionModifier.NotEquals:
        modifierString = "is not";
        break;
      case CriterionModifier.GreaterThan:
        modifierString = "is greater than";
        break;
      case CriterionModifier.LessThan:
        modifierString = "is less than";
        break;
      case CriterionModifier.IsNull:
        modifierString = "is null";
        break;
      case CriterionModifier.NotNull:
        modifierString = "is not null";
        break;
      case CriterionModifier.Includes:
        modifierString = "includes";
        break;
      case CriterionModifier.IncludesAll:
        modifierString = "includes all";
        break;
      case CriterionModifier.Excludes:
        modifierString = "excludes";
        break;
      case CriterionModifier.MatchesRegex:
        modifierString = "matches regex";
        break;
      case CriterionModifier.NotMatchesRegex:
        modifierString = "not matches regex";
        break;
      default:
        modifierString = "";
    }

    let valueString = "";

    if (
      this.modifier !== CriterionModifier.IsNull &&
      this.modifier !== CriterionModifier.NotNull
    ) {
      valueString = this.getLabelValue();
    }

    return `${intl.formatMessage({
      id: this.criterionOption.messageID,
    })} ${modifierString} ${valueString}`;
  }

  public getId(): string {
    return `${this.criterionOption.parameterName}-${this.modifier.toString()}`; // TODO add values?
  }

  public encodeValue(): V {
    return this.value;
  }

  public toJSON() {
    const encodedCriterion = {
      type: this.criterionOption.value,
      // #394 - the presence of a # symbol results in the query URL being
      // malformed. We could set encode: true in the queryString.stringify
      // call below, but this results in a URL that gets pretty long and ugly.
      // Instead, we'll encode the criteria values.
      value: this.encodeValue(),
      modifier: this.modifier,
    };
    return JSON.stringify(encodedCriterion);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public apply(outputFilter: Record<string, any>) {
    // eslint-disable-next-line no-param-reassign
    outputFilter[this.criterionOption.parameterName] = this.toCriterionInput();
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  protected toCriterionInput(): any {
    return {
      value: this.value,
      modifier: this.modifier,
    };
  }
}

export class CriterionOption {
  public messageID: string;
  public value: CriterionType;
  public parameterName: string;

  constructor(messageID: string, value: CriterionType, parameterName?: string) {
    this.messageID = messageID;
    this.value = value;
    this.parameterName = parameterName ?? value;
  }
}

export class StringCriterion extends Criterion<string> {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [
    StringCriterion.getModifierOption(CriterionModifier.Equals),
    StringCriterion.getModifierOption(CriterionModifier.NotEquals),
    StringCriterion.getModifierOption(CriterionModifier.Includes),
    StringCriterion.getModifierOption(CriterionModifier.Excludes),
    StringCriterion.getModifierOption(CriterionModifier.IsNull),
    StringCriterion.getModifierOption(CriterionModifier.NotNull),
    StringCriterion.getModifierOption(CriterionModifier.MatchesRegex),
    StringCriterion.getModifierOption(CriterionModifier.NotMatchesRegex),
  ];
  public options: string[] | undefined;
  public value: string = "";

  public getLabelValue() {
    return this.value;
  }

  public encodeValue() {
    // replace certain characters
    let ret = this.value;
    ret = StringCriterion.replaceSpecialCharacter(ret, "&");
    ret = StringCriterion.replaceSpecialCharacter(ret, "+");
    return ret;
  }

  private static replaceSpecialCharacter(str: string, c: string) {
    return str.replaceAll(c, encodeURIComponent(c));
  }

  constructor(type: CriterionOption, options?: string[]) {
    super(type);

    this.options = options;
    this.inputType = "text";
  }
}

export class MandatoryStringCriterion extends StringCriterion {
  public modifierOptions = [
    StringCriterion.getModifierOption(CriterionModifier.Equals),
    StringCriterion.getModifierOption(CriterionModifier.NotEquals),
    StringCriterion.getModifierOption(CriterionModifier.Includes),
    StringCriterion.getModifierOption(CriterionModifier.Excludes),
    StringCriterion.getModifierOption(CriterionModifier.MatchesRegex),
    StringCriterion.getModifierOption(CriterionModifier.NotMatchesRegex),
  ];
}

export class BooleanCriterion extends StringCriterion {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];

  constructor(type: CriterionOption) {
    super(type, [true.toString(), false.toString()]);
  }

  protected toCriterionInput(): boolean {
    return this.value === "true";
  }
}

export class NumberCriterion extends Criterion<number> {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Equals),
    Criterion.getModifierOption(CriterionModifier.NotEquals),
    Criterion.getModifierOption(CriterionModifier.GreaterThan),
    Criterion.getModifierOption(CriterionModifier.LessThan),
    Criterion.getModifierOption(CriterionModifier.IsNull),
    Criterion.getModifierOption(CriterionModifier.NotNull),
  ];
  public options: number[] | undefined;
  public value: number = 0;

  public getLabelValue() {
    return this.value.toString();
  }

  constructor(type: CriterionOption, options?: number[]) {
    super(type);

    this.options = options;
    this.inputType = "number";
  }
}

export abstract class ILabeledIdCriterion extends Criterion<ILabeledId[]> {
  public getLabelValue(): string {
    return this.value.map((v) => v.label).join(", ");
  }

  protected toCriterionInput(): MultiCriterionInput {
    return {
      value: this.value.map((v) => v.id),
      modifier: this.modifier,
    };
  }
}

export class MandatoryNumberCriterion extends NumberCriterion {
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Equals),
    Criterion.getModifierOption(CriterionModifier.NotEquals),
    Criterion.getModifierOption(CriterionModifier.GreaterThan),
    Criterion.getModifierOption(CriterionModifier.LessThan),
  ];
}

export class DurationCriterion extends Criterion<number> {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Equals),
    Criterion.getModifierOption(CriterionModifier.NotEquals),
    Criterion.getModifierOption(CriterionModifier.GreaterThan),
    Criterion.getModifierOption(CriterionModifier.LessThan),
  ];
  public options: number[] | undefined;
  public value: number = 0;

  constructor(type: CriterionOption, options?: number[]) {
    super(type);

    this.options = options;
  }

  public getLabelValue() {
    return DurationUtils.secondsToString(this.value);
  }
}
