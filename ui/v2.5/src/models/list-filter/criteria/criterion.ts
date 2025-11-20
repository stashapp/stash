/* eslint @typescript-eslint/no-unused-vars: ["error", { "argsIgnorePattern": "^_" }] */
import { IntlShape } from "react-intl";
import {
  CriterionModifier,
  HierarchicalMultiCriterionInput,
  IntCriterionInput,
  MultiCriterionInput,
  TimestampCriterionInput,
  ConfigDataFragment,
  DateCriterionInput,
} from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import {
  CriterionType,
  IHierarchicalLabelValue,
  ILabeledId,
  INumberValue,
  IOptionType,
  IStashIDValue,
  IDateValue,
  ITimestampValue,
  ILabeledValueListValue,
  IPhashDistanceValue,
  IRangeValue,
} from "../types";

export type Option = string | number | IOptionType;
export type CriterionValue =
  | string
  | boolean
  | string[]
  | ILabeledId[]
  | IHierarchicalLabelValue
  | ILabeledValueListValue
  | INumberValue
  | IStashIDValue
  | IDateValue
  | ITimestampValue
  | IPhashDistanceValue;

export interface ISavedCriterion<T> {
  modifier: CriterionModifier;
  value: T | undefined;
}

const modifierMessageIDs = {
  [CriterionModifier.Equals]: "criterion_modifier.equals",
  [CriterionModifier.NotEquals]: "criterion_modifier.not_equals",
  [CriterionModifier.GreaterThan]: "criterion_modifier.greater_than",
  [CriterionModifier.LessThan]: "criterion_modifier.less_than",
  [CriterionModifier.IsNull]: "criterion_modifier.is_null",
  [CriterionModifier.NotNull]: "criterion_modifier.not_null",
  [CriterionModifier.Includes]: "criterion_modifier.includes",
  [CriterionModifier.IncludesAll]: "criterion_modifier.includes_all",
  [CriterionModifier.Excludes]: "criterion_modifier.excludes",
  [CriterionModifier.MatchesRegex]: "criterion_modifier.matches_regex",
  [CriterionModifier.NotMatchesRegex]: "criterion_modifier.not_matches_regex",
  [CriterionModifier.Between]: "criterion_modifier.between",
  [CriterionModifier.NotBetween]: "criterion_modifier.not_between",
};

export abstract class Criterion {
  public criterionOption: CriterionOption;

  constructor(type: CriterionOption) {
    this.criterionOption = type;
  }

  public isValid(): boolean {
    return true;
  }

  public clone() {
    const ret = Object.assign(Object.create(Object.getPrototypeOf(this)), this);
    ret.cloneValues();
    return ret;
  }

  protected cloneValues() {}

  public abstract getLabel(intl: IntlShape, sfwContentMode?: boolean): string;

  public getId(): string {
    return `${this.criterionOption.type}`;
  }

  public abstract toQueryParams(): Record<string, unknown>;

  // fromDecodedParams is used to set the criterion from the query string
  // i is the decoded parameter object
  public abstract fromDecodedParams(i: Record<string, unknown>): void;

  public abstract applyToCriterionInput(input: Record<string, unknown>): void;

  public abstract applyToSavedCriterion(input: Record<string, unknown>): void;
  public abstract setFromSavedCriterion(criterion: unknown): void;
}

// V = criterion value type
export abstract class ModifierCriterion<
  V extends CriterionValue
> extends Criterion {
  protected _modifier!: CriterionModifier;
  public get modifier(): CriterionModifier {
    return this._modifier;
  }
  public set modifier(value: CriterionModifier) {
    this._modifier = value;
  }

  protected _value!: V;
  public get value(): V {
    return this._value;
  }
  public set value(newValue: V) {
    this._value = newValue;
  }

  public isValid(): boolean {
    return true;
  }

  protected abstract getLabelValue(intl: IntlShape): string;

  constructor(type: ModifierCriterionOption, value: V) {
    super(type);
    this.modifier = type.defaultModifier;
    this.value = value;
  }

  public modifierCriterionOption() {
    return this.criterionOption as ModifierCriterionOption;
  }

  public clone() {
    const ret = Object.assign(Object.create(Object.getPrototypeOf(this)), this);
    ret.cloneValues();
    return ret;
  }

  protected cloneValues() {}

  public static getModifierLabel(intl: IntlShape, modifier: CriterionModifier) {
    const modifierMessageID = modifierMessageIDs[modifier];

    return modifierMessageID
      ? intl.formatMessage({ id: modifierMessageID })
      : "";
  }

  public getLabel(intl: IntlShape, sfwContentMode: boolean = false): string {
    const modifierString = ModifierCriterion.getModifierLabel(
      intl,
      this.modifier
    );
    let valueString = "";

    if (
      this.modifier !== CriterionModifier.IsNull &&
      this.modifier !== CriterionModifier.NotNull
    ) {
      valueString = this.getLabelValue(intl);
    }

    const messageID = !sfwContentMode
      ? this.criterionOption.messageID
      : this.criterionOption.sfwMessageID ?? this.criterionOption.messageID;

    return intl.formatMessage(
      { id: "criterion_modifier.format_string" },
      {
        criterion: intl.formatMessage({ id: messageID }),
        modifierString,
        valueString,
      }
    );
  }

  public toQueryParams(): Record<string, unknown> {
    let encodedCriterion: Record<string, unknown> = {
      type: this.criterionOption.type,
      modifier: this.modifier,
    };

    if (
      this.modifier !== CriterionModifier.IsNull &&
      this.modifier !== CriterionModifier.NotNull
    ) {
      encodedCriterion.value = this.encodeValue();
    }

    return encodedCriterion;
  }

  protected encodeValue(): unknown {
    return this.value;
  }

  protected decodeValue(v: unknown) {
    if (v !== undefined && v !== null) {
      this.value = v as V;
    }
  }

  public fromDecodedParams(i: unknown): void {
    // use same logic as from saved criterion by default
    const c = i as ISavedCriterion<V>;
    this.modifier = c.modifier;
    this.decodeValue(c.value);
  }

  public setFromSavedCriterion(criterion: unknown) {
    const c = criterion as ISavedCriterion<V>;
    if (c.value !== undefined && c.value !== null) {
      this.value = c.value;
    }
    this.modifier = c.modifier;
  }

  public applyToCriterionInput(input: Record<string, unknown>) {
    input[this.criterionOption.type] = this.toCriterionInput();
  }

  // TODO - saved criterion _should_ be criterion input
  // kicking this can down the road a little further
  public applyToSavedCriterion(input: Record<string, unknown>): void {
    input[this.criterionOption.type] = {
      value: this.value,
      modifier: this.modifier,
    };
  }

  protected toCriterionInput(): unknown {
    return {
      value: this.value,
      modifier: this.modifier,
    };
  }
}

export type InputType =
  | "number"
  | "text"
  | "performers"
  | "studios"
  | "tags"
  | "performer_tags"
  | "scenes"
  | "scene_tags"
  | "groups"
  | "galleries"
  | undefined;

type MakeCriterionFn = (
  o: CriterionOption,
  config?: ConfigDataFragment
) => Criterion;

interface ICriterionOptionParams {
  messageID: string;
  type: CriterionType;
  makeCriterion: MakeCriterionFn;
  hidden?: boolean;
  sfwMessageID?: string;
}

export class CriterionOption {
  public readonly type: CriterionType;
  public readonly messageID: string;
  public readonly makeCriterionFn: MakeCriterionFn;
  public readonly sfwMessageID?: string;

  // used for legacy criteria that are not shown in the UI
  public readonly hidden: boolean = false;

  constructor(options: ICriterionOptionParams) {
    this.type = options.type;
    this.messageID = options.messageID;
    this.makeCriterionFn = options.makeCriterion;
    this.hidden = options.hidden ?? false;
    this.sfwMessageID = options.sfwMessageID;
  }

  public makeCriterion(config?: ConfigDataFragment) {
    return this.makeCriterionFn(this, config);
  }
}

interface IModifierCriterionOptionParams extends ICriterionOptionParams {
  inputType?: InputType;
  modifierOptions?: CriterionModifier[];
  defaultModifier?: CriterionModifier;
  options?: Option[];
}

export class ModifierCriterionOption extends CriterionOption {
  public readonly modifierOptions: CriterionModifier[];
  public readonly defaultModifier: CriterionModifier;
  public readonly options: Option[] | undefined;
  public readonly inputType: InputType;

  constructor(options: IModifierCriterionOptionParams) {
    super(options);
    this.modifierOptions = options.modifierOptions ?? [];
    this.defaultModifier = options.defaultModifier ?? CriterionModifier.Equals;
    this.options = options.options;
    this.inputType = options.inputType;
  }
}

export class ILabeledIdCriterionOption extends ModifierCriterionOption {
  constructor(
    messageID: string,
    value: CriterionType,
    includeAll: boolean,
    inputType: InputType,
    makeCriterion?: () => ModifierCriterion<CriterionValue>
  ) {
    const modifierOptions = [
      CriterionModifier.Includes,
      CriterionModifier.Excludes,
      CriterionModifier.IsNull,
      CriterionModifier.NotNull,
    ];

    let defaultModifier = CriterionModifier.Includes;
    if (includeAll) {
      modifierOptions.unshift(CriterionModifier.IncludesAll);
      defaultModifier = CriterionModifier.IncludesAll;
    }

    super({
      messageID,
      type: value,
      modifierOptions,
      defaultModifier,
      inputType,
      makeCriterion: makeCriterion
        ? makeCriterion
        : () => new ILabeledIdCriterion(this),
    });
  }
}

export class ILabeledIdCriterion extends ModifierCriterion<ILabeledId[]> {
  constructor(type: ModifierCriterionOption, value: ILabeledId[] = []) {
    super(type, value);
  }

  public cloneValues() {
    this.value = this.value.map((v) => ({ ...v }));
  }

  protected getLabelValue(_intl: IntlShape): string {
    return this.value.map((v) => v.label).join(", ");
  }

  public toCriterionInput(): MultiCriterionInput {
    return {
      value: this.value.map((v) => v.id),
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

    return this.value.length > 0;
  }
}

export class IHierarchicalLabeledIdCriterion extends ModifierCriterion<IHierarchicalLabelValue> {
  constructor(
    type: ModifierCriterionOption,
    value: IHierarchicalLabelValue = {
      items: [],
      excluded: [],
      depth: 0,
    }
  ) {
    super(type, value);
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
    // so reset it for other modifiers
    if (
      this.value &&
      value !== CriterionModifier.Includes &&
      value !== CriterionModifier.IncludesAll
    ) {
      this.value.excluded = [];
    }
  }

  public setFromSavedCriterion(
    criterion: ISavedCriterion<IHierarchicalLabelValue>
  ) {
    const { modifier, value } = criterion;

    if (value !== undefined) {
      this.value = {
        items: value.items || [],
        excluded: value.excluded || [],
        depth: value.depth || 0,
      };
    }

    const modifierOptions =
      (this.criterionOption as ModifierCriterionOption).modifierOptions ?? [];

    // if the previous modifier was excludes, replace it with the equivalent includes criterion
    // this is what is done on the backend
    // only replace if excludes is not a valid modifierOption
    if (
      modifier === CriterionModifier.Excludes &&
      modifierOptions.find((m) => m === CriterionModifier.Excludes) ===
        undefined
    ) {
      this.modifier = CriterionModifier.Includes;
      this.value.excluded = [...this.value.excluded, ...this.value.items];
      this.value.items = [];
    } else {
      this.modifier = modifier;
    }
  }

  protected getLabelValue(_intl: IntlShape): string {
    const labels = (this.value.items ?? []).map((v) => v.label).join(", ");

    if (this.value.depth === 0) {
      return labels;
    }

    return `${labels} (+${this.value.depth > 0 ? this.value.depth : "all"})`;
  }

  public toCriterionInput(): HierarchicalMultiCriterionInput {
    let excludes: string[] = [];

    // if modifier is equals, depth must be 0
    const depth =
      this.modifier === CriterionModifier.Equals ? 0 : this.value.depth;

    if (this.value.excluded) {
      excludes = this.value.excluded.map((v) => v.id);
    }
    return {
      value: this.value.items.map((v) => v.id),
      excludes: excludes,
      modifier: this.modifier,
      depth,
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

  public getLabel(intl: IntlShape, sfwContentMode?: boolean): string {
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

      if (this.value.depth !== 0) {
        id += "_depth";
      }
    }

    const messageID = !sfwContentMode
      ? this.criterionOption.messageID
      : this.criterionOption.sfwMessageID ?? this.criterionOption.messageID;

    return intl.formatMessage(
      { id },
      {
        criterion: intl.formatMessage({ id: messageID }),
        modifierString,
        valueString,
        excludedString,
        depth: this.value.depth,
      }
    );
  }
}

export class StringCriterionOption extends ModifierCriterionOption {
  constructor(
    options: Partial<
      Omit<IModifierCriterionOptionParams, "messageID" | "type">
    > &
      Pick<IModifierCriterionOptionParams, "messageID" | "type">
  ) {
    super({
      modifierOptions: [
        CriterionModifier.Equals,
        CriterionModifier.NotEquals,
        CriterionModifier.Includes,
        CriterionModifier.Excludes,
        CriterionModifier.IsNull,
        CriterionModifier.NotNull,
        CriterionModifier.MatchesRegex,
        CriterionModifier.NotMatchesRegex,
      ],
      defaultModifier: CriterionModifier.Equals,
      inputType: "text",
      makeCriterion: () => new StringCriterion(this),
      ...options,
    });
  }
}

export function createStringCriterionOption(
  type: CriterionType,
  messageID?: string,
  options?: { nsfw?: boolean }
) {
  return new StringCriterionOption({
    messageID: messageID ?? type,
    type,
    ...options,
  });
}

export class MandatoryStringCriterionOption extends ModifierCriterionOption {
  constructor(messageID: string, value: CriterionType) {
    super({
      messageID,
      type: value,
      modifierOptions: [
        CriterionModifier.Equals,
        CriterionModifier.NotEquals,
        CriterionModifier.Includes,
        CriterionModifier.Excludes,
        CriterionModifier.MatchesRegex,
        CriterionModifier.NotMatchesRegex,
      ],
      defaultModifier: CriterionModifier.Equals,
      inputType: "text",
      makeCriterion: () => new StringCriterion(this),
    });
  }
}

export function createMandatoryStringCriterionOption(
  value: CriterionType,
  messageID?: string
) {
  return new MandatoryStringCriterionOption(messageID ?? value, value);
}

export class StringCriterion extends ModifierCriterion<string> {
  constructor(type: ModifierCriterionOption) {
    super(type, "");
  }

  protected getLabelValue(_intl: IntlShape) {
    return this.value;
  }

  public isValid(): boolean {
    return (
      this.modifier === CriterionModifier.IsNull ||
      this.modifier === CriterionModifier.NotNull ||
      this.value.length > 0
    );
  }
}

export abstract class MultiStringCriterion extends ModifierCriterion<string[]> {
  constructor(type: ModifierCriterionOption, value: string[] = []) {
    super(type, value);
  }

  public cloneValues() {
    this.value = this.value.slice();
  }

  protected getLabelValue(_intl: IntlShape) {
    return this.value.join(", ");
  }

  public isValid(): boolean {
    return (
      this.modifier === CriterionModifier.IsNull ||
      this.modifier === CriterionModifier.NotNull ||
      this.value.length > 0
    );
  }
}

export class BooleanCriterionOption extends ModifierCriterionOption {
  constructor(
    messageID: string,
    value: CriterionType,
    makeCriterion?: () => ModifierCriterion<CriterionValue>
  ) {
    super({
      messageID,
      type: value,
      modifierOptions: [],
      defaultModifier: CriterionModifier.Equals,
      options: ["true", "false"],
      makeCriterion: makeCriterion
        ? makeCriterion
        : () => new BooleanCriterion(this),
    });
  }
}

export function createBooleanCriterionOption(
  value: CriterionType,
  messageID?: string
) {
  return new BooleanCriterionOption(messageID ?? value, value);
}

export class BooleanCriterion extends StringCriterion {
  public toCriterionInput(): boolean {
    return this.value === "true";
  }

  public isValid() {
    return this.value === "true" || this.value === "false";
  }
}

export class StringBooleanCriterionOption extends ModifierCriterionOption {
  constructor(
    messageID: string,
    value: CriterionType,
    makeCriterion?: () => ModifierCriterion<CriterionValue>
  ) {
    super({
      messageID,
      type: value,
      options: ["true", "false"],
      makeCriterion: makeCriterion
        ? makeCriterion
        : () => new StringBooleanCriterion(this),
    });
  }
}

export class StringBooleanCriterion extends StringCriterion {
  public toCriterionInput(): string {
    return this.value;
  }

  public isValid() {
    return this.value === "true" || this.value === "false";
  }
}

export class NumberCriterionOption extends ModifierCriterionOption {
  constructor(messageID: string, value: CriterionType) {
    super({
      messageID,
      type: value,
      modifierOptions: [
        CriterionModifier.Equals,
        CriterionModifier.NotEquals,
        CriterionModifier.GreaterThan,
        CriterionModifier.LessThan,
        CriterionModifier.IsNull,
        CriterionModifier.NotNull,
        CriterionModifier.Between,
        CriterionModifier.NotBetween,
      ],
      defaultModifier: CriterionModifier.Equals,
      inputType: "number",
      makeCriterion: () => new NumberCriterion(this),
    });
  }
}

export function createNumberCriterionOption(
  value: CriterionType,
  messageID?: string
) {
  return new NumberCriterionOption(messageID ?? value, value);
}

export class NullNumberCriterionOption extends ModifierCriterionOption {
  constructor(
    messageID: string,
    value: CriterionType,
    makeCriterion?: MakeCriterionFn
  ) {
    super({
      messageID,
      type: value,
      modifierOptions: [
        CriterionModifier.Equals,
        CriterionModifier.NotEquals,
        CriterionModifier.GreaterThan,
        CriterionModifier.LessThan,
        CriterionModifier.Between,
        CriterionModifier.NotBetween,
        CriterionModifier.IsNull,
        CriterionModifier.NotNull,
      ],
      defaultModifier: CriterionModifier.Equals,
      inputType: "number",
      makeCriterion: makeCriterion
        ? makeCriterion
        : () => new NumberCriterion(this),
    });
  }
}

export function createNullNumberCriterionOption(
  value: CriterionType,
  messageID?: string
) {
  return new NullNumberCriterionOption(messageID ?? value, value);
}

export class MandatoryNumberCriterionOption extends ModifierCriterionOption {
  constructor(
    messageID: string,
    value: CriterionType,
    makeCriterion?: () => ModifierCriterion<CriterionValue>,
    options?: { sfwMessageID?: string }
  ) {
    super({
      messageID,
      type: value,
      modifierOptions: [
        CriterionModifier.Equals,
        CriterionModifier.NotEquals,
        CriterionModifier.GreaterThan,
        CriterionModifier.LessThan,
        CriterionModifier.Between,
        CriterionModifier.NotBetween,
      ],
      defaultModifier: CriterionModifier.Equals,
      inputType: "number",
      makeCriterion: makeCriterion
        ? makeCriterion
        : () => new NumberCriterion(this),
      ...options,
    });
  }
}

export function createMandatoryNumberCriterionOption(
  value: CriterionType,
  messageID?: string,
  options?: { sfwMessageID?: string }
) {
  return new MandatoryNumberCriterionOption(
    messageID ?? value,
    value,
    undefined,
    options
  );
}

export function encodeRangeValue<V>(
  modifier: CriterionModifier,
  value: IRangeValue<V>
): unknown {
  // only encode value2 if modifier is between/not between
  if (
    modifier === CriterionModifier.Between ||
    modifier === CriterionModifier.NotBetween
  ) {
    return { value: value.value, value2: value.value2 };
  }

  return { value: value.value };
}

export function decodeRangeValue<V>(v: {
  value: V | IRangeValue<V>;
  value2?: V;
}): IRangeValue<V> {
  // handle backwards compatible value
  if (typeof v.value === "object") {
    return v.value as IRangeValue<V>;
  } else {
    return { value: v.value, value2: v.value2 };
  }
}

export class NumberCriterion extends ModifierCriterion<INumberValue> {
  constructor(type: ModifierCriterionOption) {
    super(type, { value: undefined, value2: undefined });
  }

  public cloneValues() {
    this.value = { ...this.value };
  }

  public get value(): INumberValue {
    return this._value;
  }
  public set value(newValue: number | INumberValue) {
    // backwards compatibility - if this.value is a number, use that
    if (typeof newValue !== "object") {
      this._value = {
        value: newValue,
        value2: undefined,
      };
    } else {
      this._value = newValue;
    }
  }

  public toCriterionInput(): IntCriterionInput {
    return {
      modifier: this.modifier,
      value: this.value?.value ?? 0,
      value2: this.value?.value2,
    };
  }

  public setFromSavedCriterion(c: {
    modifier: CriterionModifier;
    value: number | INumberValue;
    value2?: number;
  }) {
    super.setFromSavedCriterion(c);
    // this.value = decodeRangeValue(c);
  }

  protected encodeValue(): unknown {
    return encodeRangeValue(this.modifier, this.value);
  }

  protected getLabelValue(_intl: IntlShape) {
    const { value, value2 } = this.value;
    if (
      this.modifier === CriterionModifier.Between ||
      this.modifier === CriterionModifier.NotBetween
    ) {
      return `${value}, ${value2 ?? 0}`;
    } else {
      return `${value}`;
    }
  }

  public isValid(): boolean {
    if (
      this.modifier === CriterionModifier.IsNull ||
      this.modifier === CriterionModifier.NotNull
    ) {
      return true;
    }

    const { value, value2 } = this.value;
    if (value === undefined) {
      return false;
    }

    if (
      value2 === undefined &&
      (this.modifier === CriterionModifier.Between ||
        this.modifier === CriterionModifier.NotBetween)
    ) {
      return false;
    }

    return true;
  }
}

export class DurationCriterionOption extends MandatoryNumberCriterionOption {
  constructor(messageID: string, value: CriterionType) {
    super(messageID, value, () => new DurationCriterion(this));
  }
}

export function createDurationCriterionOption(
  value: CriterionType,
  messageID?: string
) {
  return new DurationCriterionOption(messageID ?? value, value);
}

export class NullDurationCriterionOption extends NullNumberCriterionOption {
  constructor(messageID: string, value: CriterionType) {
    super(messageID, value, () => new DurationCriterion(this));
  }
}

export function createNullDurationCriterionOption(
  value: CriterionType,
  messageID?: string
) {
  return new NullDurationCriterionOption(messageID ?? value, value);
}

export class DurationCriterion extends ModifierCriterion<INumberValue> {
  constructor(type: ModifierCriterionOption) {
    super(type, { value: undefined, value2: undefined });
  }

  public cloneValues() {
    this.value = { ...this.value };
  }

  public toCriterionInput(): IntCriterionInput {
    return {
      modifier: this.modifier,
      value: this.value?.value ?? 0,
      value2: this.value?.value2,
    };
  }

  public setFromSavedCriterion(c: {
    modifier: CriterionModifier;
    value: number | INumberValue;
    value2?: number;
  }) {
    super.setFromSavedCriterion(c);
    // this.value = decodeRangeValue(c);
  }

  protected encodeValue(): unknown {
    return encodeRangeValue(this.modifier, this.value);
  }

  protected getLabelValue(_intl: IntlShape) {
    const value = TextUtils.secondsToTimestamp(this.value.value ?? 0);
    const value2 = TextUtils.secondsToTimestamp(this.value.value2 ?? 0);
    if (
      this.modifier === CriterionModifier.Between ||
      this.modifier === CriterionModifier.NotBetween
    ) {
      return `${value}, ${value2}`;
    } else {
      return value;
    }
  }

  public isValid(): boolean {
    if (
      this.modifier === CriterionModifier.IsNull ||
      this.modifier === CriterionModifier.NotNull
    ) {
      return true;
    }

    const { value, value2 } = this.value;
    if (value === undefined) {
      return false;
    }

    if (
      value2 === undefined &&
      (this.modifier === CriterionModifier.Between ||
        this.modifier === CriterionModifier.NotBetween)
    ) {
      return false;
    }

    return true;
  }
}

export class DateCriterionOption extends ModifierCriterionOption {
  constructor(messageID: string, value: CriterionType) {
    super({
      messageID,
      type: value,
      modifierOptions: [
        CriterionModifier.Equals,
        CriterionModifier.NotEquals,
        CriterionModifier.GreaterThan,
        CriterionModifier.LessThan,
        CriterionModifier.IsNull,
        CriterionModifier.NotNull,
        CriterionModifier.Between,
        CriterionModifier.NotBetween,
      ],
      defaultModifier: CriterionModifier.Equals,
      inputType: "text",
      makeCriterion: () => new DateCriterion(this),
    });
  }
}

export function createDateCriterionOption(value: CriterionType) {
  return new DateCriterionOption(value, value);
}

export class DateCriterion extends ModifierCriterion<IDateValue> {
  constructor(type: ModifierCriterionOption) {
    super(type, { value: "", value2: undefined });
  }

  public cloneValues() {
    this.value = { ...this.value };
  }

  public setFromSavedCriterion(c: {
    modifier: CriterionModifier;
    value: string | IDateValue;
    value2?: string;
  }) {
    super.setFromSavedCriterion(c);
    // this.value = decodeRangeValue(c);
  }

  protected encodeValue(): unknown {
    return encodeRangeValue(this.modifier, this.value);
  }

  public toCriterionInput(): DateCriterionInput {
    return {
      modifier: this.modifier,
      value: this.value?.value ?? "",
      value2: this.value?.value2,
    };
  }

  protected getLabelValue() {
    const { value } = this.value;
    return this.modifier === CriterionModifier.Between ||
      this.modifier === CriterionModifier.NotBetween
      ? `${value}, ${this.value.value2}`
      : `${value}`;
  }

  public isValid(): boolean {
    if (
      this.modifier === CriterionModifier.IsNull ||
      this.modifier === CriterionModifier.NotNull
    ) {
      return true;
    }

    const { value, value2 } = this.value;
    if (!value) {
      return false;
    }

    if (
      !value2 &&
      (this.modifier === CriterionModifier.Between ||
        this.modifier === CriterionModifier.NotBetween)
    ) {
      return false;
    }

    return true;
  }
}

export class TimestampCriterionOption extends ModifierCriterionOption {
  constructor(messageID: string, value: CriterionType) {
    super({
      messageID,
      type: value,
      modifierOptions: [
        CriterionModifier.GreaterThan,
        CriterionModifier.LessThan,
        CriterionModifier.IsNull,
        CriterionModifier.NotNull,
        CriterionModifier.Between,
        CriterionModifier.NotBetween,
      ],
      defaultModifier: CriterionModifier.GreaterThan,
      inputType: "text",
      makeCriterion: () => new TimestampCriterion(this),
    });
  }
}

export function createTimestampCriterionOption(value: CriterionType) {
  return new TimestampCriterionOption(value, value);
}

export class MandatoryTimestampCriterionOption extends ModifierCriterionOption {
  constructor(messageID: string, value: CriterionType) {
    super({
      messageID,
      type: value,
      modifierOptions: [
        CriterionModifier.GreaterThan,
        CriterionModifier.LessThan,
        CriterionModifier.Between,
        CriterionModifier.NotBetween,
      ],
      defaultModifier: CriterionModifier.GreaterThan,
      inputType: "text",
      makeCriterion: () => new TimestampCriterion(this),
    });
  }
}

export function createMandatoryTimestampCriterionOption(value: CriterionType) {
  return new MandatoryTimestampCriterionOption(value, value);
}

export class TimestampCriterion extends ModifierCriterion<ITimestampValue> {
  constructor(type: ModifierCriterionOption) {
    super(type, { value: "", value2: undefined });
  }

  public cloneValues() {
    this.value = { ...this.value };
  }

  public toCriterionInput(): TimestampCriterionInput {
    return {
      modifier: this.modifier,
      value: this.transformValueToInput(this.value.value ?? ""),
      value2: this.value.value2
        ? this.transformValueToInput(this.value.value2)
        : null,
    };
  }

  public setFromSavedCriterion(c: {
    modifier: CriterionModifier;
    value: string | ITimestampValue;
    value2?: string;
  }) {
    super.setFromSavedCriterion(c);
    this.value = decodeRangeValue(c);
  }

  protected encodeValue(): unknown {
    return encodeRangeValue(this.modifier, this.value);
  }

  protected getLabelValue() {
    const { value } = this.value;
    return this.modifier === CriterionModifier.Between ||
      this.modifier === CriterionModifier.NotBetween
      ? `${value}, ${this.value.value2}`
      : `${value}`;
  }

  private transformValueToInput(value: string): string {
    value = value.trim();
    if (/^\d{4}-\d{2}-\d{2}(( |T)\d{2}:\d{2})?$/.test(value)) {
      return value.replace(" ", "T");
    }

    return "";
  }

  public isValid(): boolean {
    if (
      this.modifier === CriterionModifier.IsNull ||
      this.modifier === CriterionModifier.NotNull
    ) {
      return true;
    }

    const { value, value2 } = this.value;
    if (!value) {
      return false;
    }

    if (
      !value2 &&
      (this.modifier === CriterionModifier.Between ||
        this.modifier === CriterionModifier.NotBetween)
    ) {
      return false;
    }

    return true;
  }
}
