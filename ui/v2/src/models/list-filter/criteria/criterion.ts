import { isArray } from "util";
import { ILabeledId } from "../types";

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

export enum CriterionModifier {
  Equals,
  NotEquals,
  GreaterThan,
  LessThan,
  Inclusive,
  Exclusive,
}

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

  public abstract type: CriterionType;
  public abstract parameterName: string;
  public abstract modifier: CriterionModifier;
  public abstract options: Option[];
  public abstract value: Value;

  public getLabel(): string {
    let modifierString: string;
    switch (this.modifier) {
      case CriterionModifier.Equals: modifierString = "is"; break;
      case CriterionModifier.NotEquals: modifierString = "is not"; break;
      case CriterionModifier.GreaterThan: modifierString = "is greater than"; break;
      case CriterionModifier.LessThan: modifierString = "is less than"; break;
      case CriterionModifier.Inclusive: modifierString = "includes"; break;
      case CriterionModifier.Exclusive: modifierString = "exculdes"; break;
      default: modifierString = "";
    }

    let valueString: string;
    if (isArray(this.value) && this.value.length > 0) {
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
