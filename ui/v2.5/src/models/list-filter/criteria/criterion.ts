/* eslint-disable consistent-return */

import { IntlShape } from "react-intl";
import {
  CriterionModifier,
  HierarchicalMultiCriterionInput,
  MultiCriterionInput,
} from "src/core/generated-graphql";
import DurationUtils from "src/utils/duration";
import {
  CriterionType,
  encodeILabeledId,
  ILabeledId,
  ILabeledValue,
  IOptionType,
  IHierarchicalLabelValue,
} from "../types";

type Option = string | number | IOptionType;
export type CriterionValue =
  | string
  | number
  | ILabeledId[]
  | IHierarchicalLabelValue;

// V = criterion value type
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
    let modifierMessageID: string;
    switch (this.modifier) {
      case CriterionModifier.Equals:
        modifierMessageID = "criterion_modifier.equals";
        break;
      case CriterionModifier.NotEquals:
        modifierMessageID = "criterion_modifier.not_equals";
        break;
      case CriterionModifier.GreaterThan:
        modifierMessageID = "criterion_modifier.greater_than";
        break;
      case CriterionModifier.LessThan:
        modifierMessageID = "criterion_modifier.less_than";
        break;
      case CriterionModifier.IsNull:
        modifierMessageID = "criterion_modifier.is_null";
        break;
      case CriterionModifier.NotNull:
        modifierMessageID = "criterion_modifier.not_null";
        break;
      case CriterionModifier.Includes:
        modifierMessageID = "criterion_modifier.includes";
        break;
      case CriterionModifier.IncludesAll:
        modifierMessageID = "criterion_modifier.includes_all";
        break;
      case CriterionModifier.Excludes:
        modifierMessageID = "criterion_modifier.excludes";
        break;
      case CriterionModifier.MatchesRegex:
        modifierMessageID = "criterion_modifier.matches_regex";
        break;
      case CriterionModifier.NotMatchesRegex:
        modifierMessageID = "criterion_modifier.not_matches_regex";
        break;
      default:
        modifierMessageID = "";
    }

    const modifierString = modifierMessageID
      ? intl.formatMessage({ id: modifierMessageID })
      : "";
    let valueString = "";

    if (
      this.modifier !== CriterionModifier.IsNull &&
      this.modifier !== CriterionModifier.NotNull
    ) {
      valueString = this.getLabelValue();
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
  public readonly messageID: string;
  public readonly value: CriterionType;
  public readonly parameterName: string;

  constructor(messageID: string, value: CriterionType, parameterName?: string) {
    this.messageID = messageID;
    this.value = value;
    this.parameterName = parameterName ?? value;
  }
}

export function createCriterionOption(value: CriterionType) {
  return new CriterionOption(value, value);
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
  public modifier = CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.IncludesAll),
    Criterion.getModifierOption(CriterionModifier.Includes),
    Criterion.getModifierOption(CriterionModifier.Excludes),
  ];

  public options: IOptionType[] = [];
  public value: ILabeledId[] = [];

  public getLabelValue(): string {
    return this.value.map((v) => v.label).join(", ");
  }

  protected toCriterionInput(): MultiCriterionInput {
    return {
      value: this.value.map((v) => v.id),
      modifier: this.modifier,
    };
  }

  public encodeValue() {
    return this.value.map((o) => {
      return encodeILabeledId(o);
    });
  }

  constructor(type: CriterionOption, includeAll: boolean) {
    super(type);

    if (!includeAll) {
      this.modifier = CriterionModifier.Includes;
      this.modifierOptions = [
        Criterion.getModifierOption(CriterionModifier.Includes),
        Criterion.getModifierOption(CriterionModifier.Excludes),
      ];
    }
  }
}

export abstract class IHierarchicalLabeledIdCriterion extends Criterion<IHierarchicalLabelValue> {
  public modifier = CriterionModifier.IncludesAll;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.IncludesAll),
    Criterion.getModifierOption(CriterionModifier.Includes),
    Criterion.getModifierOption(CriterionModifier.Excludes),
  ];

  public options: IOptionType[] = [];
  public value: IHierarchicalLabelValue = {
    items: [],
    depth: 0,
  };

  public encodeValue() {
    return {
      items: this.value.items.map((o) => {
        return encodeILabeledId(o);
      }),
      depth: this.value.depth,
    };
  }

  protected toCriterionInput(): HierarchicalMultiCriterionInput {
    return {
      value: this.value.items.map((v) => v.id),
      modifier: this.modifier,
      depth: this.value.depth,
    };
  }

  public getLabelValue(): string {
    const labels = this.value.items.map((v) => v.label).join(", ");

    if (this.value.depth === 0) {
      return labels;
    }

    return `${labels} (+${this.value.depth > 0 ? this.value.depth : "all"})`;
  }

  public toJSON() {
    const encodedCriterion = {
      type: this.criterionOption.value,
      value: this.encodeValue(),
      modifier: this.modifier,
    };
    return JSON.stringify(encodedCriterion);
  }

  constructor(type: CriterionOption, includeAll: boolean) {
    super(type);

    if (!includeAll) {
      this.modifier = CriterionModifier.Includes;
      this.modifierOptions = [
        Criterion.getModifierOption(CriterionModifier.Includes),
        Criterion.getModifierOption(CriterionModifier.Excludes),
      ];
    }
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
