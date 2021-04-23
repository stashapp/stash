/* eslint-disable consistent-return */

import { IntlShape } from "react-intl";
import {
  CriterionModifier,
  MultiCriterionInput,
} from "src/core/generated-graphql";
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
  | "tag_count"
  | "performers"
  | "studios"
  | "movies"
  | "galleries"
  | "birth_year"
  | "age"
  | "ethnicity"
  | "country"
  | "hair_color"
  | "eye_color"
  | "height"
  | "weight"
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
  | "performer_count"
  | "death_year"
  | "url"
  | "stash_id";

type Option = string | number | IOptionType;
export type CriterionValue = string | number | ILabeledId[];

// V = criterion value type
// I = criterion input type
export abstract class Criterion<V extends CriterionValue> {
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
      case "tag_count":
        return "Tag Count";
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
      case "death_year":
        return "Death Year";
      case "age":
        return "Age";
      case "ethnicity":
        return "Ethnicity";
      case "country":
        return "Country";
      case "hair_color":
        return "Hair Color";
      case "eye_color":
        return "Eye Color";
      case "height":
        return "Height";
      case "weight":
        return "Weight";
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
      case "url":
        return "URL";
      case "stash_id":
        return "StashID";
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

export abstract class ILabeledIdCriterion extends Criterion<ILabeledId[]> {
  public getLabelValue(): string {
    return this.value.map((v) => v.label).join(", ");
  }
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
