import { isArray } from "util";
import { CriterionModifier } from "../../../core/generated-graphql";
import { ILabeledId, ILabeledValue } from "../types";

export type CriterionType =
  "none" |
  "rating" |
  "resolution" |
  "favorite" |
  "hasMarkers" |
  "isMissing" |
  "tags" |
  "sceneTags" |
  "performers" |
  "studios";

export abstract class Criterion<Option = any, Value = any> {
  public static getLabel(type: CriterionType = "none"): string {
    switch (type) {
      case "none": return "None";
      case "rating": return "Rating";
      case "resolution": return "Resolution";
      case "favorite": return "Favorite";
      case "hasMarkers": return "Has Markers";
      case "isMissing": return "Is Missing";
      case "tags": return "Tags";
      case "sceneTags": return "Scene Tags";
      case "performers": return "Performers";
      case "studios": return "Studios";
    }
  }

  public static getModifierOption(modifier: CriterionModifier = CriterionModifier.Equals): ILabeledValue {
    switch (modifier) {
      case CriterionModifier.Equals: return {value: CriterionModifier.Equals, label: "Equals"};
      case CriterionModifier.NotEquals: return {value: CriterionModifier.NotEquals, label: "Not Equals"};
      case CriterionModifier.GreaterThan: return {value: CriterionModifier.GreaterThan, label: "Greater Than"};
      case CriterionModifier.LessThan: return {value: CriterionModifier.LessThan, label: "Less Than"};
      case CriterionModifier.IsNull: return {value: CriterionModifier.IsNull, label: "Is NULL"};
      case CriterionModifier.NotNull: return {value: CriterionModifier.NotNull, label: "Not NULL"};
      case CriterionModifier.IncludesAll: return {value: CriterionModifier.IncludesAll, label: "Includes All"};
      case CriterionModifier.Includes: return {value: CriterionModifier.Includes, label: "Includes"};
      case CriterionModifier.Excludes: return {value: CriterionModifier.Excludes, label: "Excludes"};
    }
  }

  public abstract type: CriterionType;
  public abstract parameterName: string;
  public abstract modifier: CriterionModifier;
  public abstract modifierOptions: ILabeledValue[];
  public abstract options: Option[];
  public abstract value: Value;

  public getLabel(): string {
    let modifierString: string;
    switch (this.modifier) {
      case CriterionModifier.Equals: modifierString = "is"; break;
      case CriterionModifier.NotEquals: modifierString = "is not"; break;
      case CriterionModifier.GreaterThan: modifierString = "is greater than"; break;
      case CriterionModifier.LessThan: modifierString = "is less than"; break;
      case CriterionModifier.IsNull: modifierString = "is null"; break;
      case CriterionModifier.NotNull: modifierString = "is not null"; break;
      case CriterionModifier.Includes: modifierString = "includes"; break;
      case CriterionModifier.IncludesAll: modifierString = "includes all"; break;
      case CriterionModifier.Excludes: modifierString = "excludes"; break;
      default: modifierString = "";
    }

    let valueString: string;
    if (this.modifier === CriterionModifier.IsNull || this.modifier === CriterionModifier.NotNull) {
      valueString = "";
    } else if (isArray(this.value) && this.value.length > 0) {
      let items = this.value;
      if ((this.value as ILabeledId[])[0].label) {
        items = this.value.map((item) => item.label) as any;
      }
      valueString = items.join(", ");
    } else if (typeof this.value === "string") {
      valueString = this.value;
    } else {
      valueString = this.value.toString();
    }

    return `${Criterion.getLabel(this.type)} ${modifierString} ${valueString}`;
  }

  public getId(): string {
    return `${this.parameterName}-${this.modifier.toString()}`; // TODO add values?
  }

  public set(modifier: CriterionModifier, value: Value) {
    this.modifier = modifier;
    if (isArray(this.value)) {
      this.value.push(value);
    } else {
      this.value = value;
    }
  }
}

export interface ICriterionOption {
  label: string;
  value: CriterionType;
}
