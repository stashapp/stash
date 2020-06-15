/* eslint-disable consistent-return */

import { CriterionModifier } from "src/core/generated-graphql";
import DurationUtils from "src/utils/duration";
import { ILabeledId, ILabeledValue, IOptionType } from "../types";

export type CriterionType =
  | "none"
  | "rating"
  | "o_counter"
  | "resolution"
  | "duration"
  | "favorite"
  | "hasMarkers"
  | "sceneIsMissing"
  | "performerIsMissing"
  | "tags"
  | "sceneTags"
  | "performers"
  | "studios"
  | "movies"
  | "birth_year"
  | "age"
  | "ethnicity"
  | "country"
  | "eye_color"
  | "height"
  | "measurements"
  | "fake_tits"
  | "career_length"
  | "tattoos"
  | "piercings"
  | "aliases"
  | "gender"
  | "parent_studios";

type Option = string | number | IOptionType;
export type CriterionValue = string | number | ILabeledId[];

export abstract class Criterion {
  public static getLabel(type: CriterionType = "none") {
    switch (type) {
      case "none":
        return "None";
      case "rating":
        return "Rating";
      case "o_counter":
        return "O-Counter";
      case "resolution":
        return "Resolution";
      case "duration":
        return "Duration";
      case "favorite":
        return "Favorite";
      case "hasMarkers":
        return "Has Markers";
      case "sceneIsMissing":
        return "Is Missing";
      case "performerIsMissing":
        return "Is Missing";
      case "tags":
        return "Tags";
      case "sceneTags":
        return "Scene Tags";
      case "performers":
        return "Performers";
      case "studios":
        return "Studios";
      case "movies":
        return "Movies";
      case "birth_year":
        return "Birth Year";
      case "age":
        return "Age";
      case "ethnicity":
        return "Ethnicity";
      case "country":
        return "Country";
      case "eye_color":
        return "Eye Color";
      case "height":
        return "Height";
      case "measurements":
        return "Measurements";
      case "fake_tits":
        return "Fake Tits";
      case "career_length":
        return "Career Length";
      case "tattoos":
        return "Tattoos";
      case "piercings":
        return "Piercings";
      case "aliases":
        return "Aliases";
      case "gender":
        return "Gender";
      case "parent_studios":
        return "Parent Studios";
    }
  }

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
    }
  }

  public abstract type: CriterionType;
  public abstract parameterName: string;
  public abstract modifier: CriterionModifier;
  public abstract modifierOptions: ILabeledValue[];
  public abstract options: Option[] | undefined;
  public abstract value: CriterionValue;
  public inputType: "number" | "text" | undefined;

  public getLabelValue(): string {
    if (typeof this.value === "string") return this.value;
    if (typeof this.value === "number") return this.value.toString();
    return this.value.map((v) => v.label).join(", ");
  }

  public getLabel(): string {
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

    return `${Criterion.getLabel(this.type)} ${modifierString} ${valueString}`;
  }

  public getId(): string {
    return `${this.parameterName}-${this.modifier.toString()}`; // TODO add values?
  }

  /*
  public set(modifier: CriterionModifier, value: Value) {
    this.modifier = modifier;
    if (Array.isArray(this.value)) {
      this.value.push(value);
    } else {
      this.value = value;
    }
  }
  */

  public encodeValue(): CriterionValue {
    return this.value;
  }
}

export interface ICriterionOption {
  label: string;
  value: CriterionType;
}

export class CriterionOption implements ICriterionOption {
  public label: string;
  public value: CriterionType;

  constructor(label: string, value: CriterionType) {
    this.label = label;
    this.value = value;
  }
}

export class StringCriterion extends Criterion {
  public type: CriterionType;
  public parameterName: string;
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Equals),
    Criterion.getModifierOption(CriterionModifier.NotEquals),
    Criterion.getModifierOption(CriterionModifier.IsNull),
    Criterion.getModifierOption(CriterionModifier.NotNull),
  ];
  public options: string[] | undefined;
  public value: string = "";

  public getLabelValue() {
    return this.value;
  }

  constructor(type: CriterionType, parameterName?: string, options?: string[]) {
    super();

    this.type = type;
    this.options = options;
    this.inputType = "text";

    if (parameterName) {
      this.parameterName = parameterName;
    } else {
      this.parameterName = type;
    }
  }
}

export class NumberCriterion extends Criterion {
  public type: CriterionType;
  public parameterName: string;
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

  constructor(type: CriterionType, parameterName?: string, options?: number[]) {
    super();

    this.type = type;
    this.options = options;
    this.inputType = "number";

    if (parameterName) {
      this.parameterName = parameterName;
    } else {
      this.parameterName = type;
    }
  }
}

export class DurationCriterion extends Criterion {
  public type: CriterionType;
  public parameterName: string;
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Equals),
    Criterion.getModifierOption(CriterionModifier.NotEquals),
    Criterion.getModifierOption(CriterionModifier.GreaterThan),
    Criterion.getModifierOption(CriterionModifier.LessThan),
  ];
  public options: number[] | undefined;
  public value: number = 0;

  constructor(type: CriterionType, parameterName?: string, options?: number[]) {
    super();

    this.type = type;
    this.options = options;
    this.parameterName = parameterName ?? type;
  }

  public getLabelValue() {
    return DurationUtils.secondsToString(this.value);
  }
}
