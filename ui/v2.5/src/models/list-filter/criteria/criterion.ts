/* eslint-disable consistent-return */

import { CriterionModifier } from "src/core/generated-graphql";
import DurationUtils from "src/utils/duration";
import { ILabeledId, ILabeledValue, IOptionType } from "../types";

export type CriterionType =
  | "none"
  | "path"
  | "rating"
  | "organized"
  | "o_counter"
  | "resolution"
  | "average_resolution"
  | "duration"
  | "favorite"
  | "hasMarkers"
  | "sceneIsMissing"
  | "imageIsMissing"
  | "performerIsMissing"
  | "galleryIsMissing"
  | "tagIsMissing"
  | "studioIsMissing"
  | "movieIsMissing"
  | "tags"
  | "sceneTags"
  | "performerTags"
  | "performers"
  | "studios"
  | "movies"
  | "galleries"
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
  | "parent_studios"
  | "scene_count"
  | "marker_count"
  | "image_count"
  | "gallery_count"
  | "performer_count";

type Option = string | number | IOptionType;
export type CriterionValue = string | number | ILabeledId[];

export abstract class Criterion {
  public static getLabel(type: CriterionType = "none") {
    switch (type) {
      case "none":
        return "None";
      case "path":
        return "Path";
      case "rating":
        return "Rating";
      case "organized":
        return "Organized";
      case "o_counter":
        return "O-Counter";
      case "resolution":
        return "Resolution";
      case "average_resolution":
        return "Average Resolution";
      case "duration":
        return "Duration";
      case "favorite":
        return "Favorite";
      case "hasMarkers":
        return "Has Markers";
      case "sceneIsMissing":
      case "imageIsMissing":
      case "performerIsMissing":
      case "galleryIsMissing":
      case "tagIsMissing":
      case "studioIsMissing":
      case "movieIsMissing":
        return "Is Missing";
      case "tags":
        return "Tags";
      case "sceneTags":
        return "Scene Tags";
      case "performerTags":
        return "Performer Tags";
      case "performers":
        return "Performers";
      case "studios":
        return "Studios";
      case "movies":
        return "Movies";
      case "galleries":
        return "Galleries";
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
      case "scene_count":
        return "Scene Count";
      case "marker_count":
        return "Marker Count";
      case "image_count":
        return "Image Count";
      case "gallery_count":
        return "Gallery Count";
      case "performer_count":
        return "Performer Count";
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

    return `${Criterion.getLabel(this.type)} ${modifierString} ${valueString}`;
  }

  public getId(): string {
    return `${this.parameterName}-${this.modifier.toString()}`; // TODO add values?
  }

  private static replaceSpecialCharacter(str: string, c: string) {
    return str.replaceAll(c, encodeURIComponent(c));
  }

  public encodeValue(): CriterionValue {
    // replace certain characters
    if (typeof this.value === "string") {
      let ret = this.value;
      ret = Criterion.replaceSpecialCharacter(ret, "&");
      ret = Criterion.replaceSpecialCharacter(ret, "+");
      return ret;
    }
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
