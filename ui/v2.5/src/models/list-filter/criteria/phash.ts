import {
  CriterionModifier,
  PhashDistanceCriterionInput,
  DuplicationCriterionInput,
} from "src/core/generated-graphql";
import { IDuplicationValue, IPhashDistanceValue } from "../types";
import {
  Criterion,
  CriterionOption,
  ModifierCriterion,
  ModifierCriterionOption,
} from "./criterion";
import { IntlShape } from "react-intl";

export const PhashCriterionOption = new ModifierCriterionOption({
  messageID: "media_info.phash",
  type: "phash_distance",
  inputType: "text",
  modifierOptions: [
    CriterionModifier.Equals,
    CriterionModifier.NotEquals,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
  makeCriterion: () => new PhashCriterion(),
});

export class PhashCriterion extends ModifierCriterion<IPhashDistanceValue> {
  constructor() {
    super(PhashCriterionOption, { value: "", distance: 0 });
  }

  public cloneValues() {
    this.value = { ...this.value };
  }

  protected getLabelValue() {
    const { value, distance } = this.value;
    if (
      (this.modifier === CriterionModifier.Equals ||
        this.modifier === CriterionModifier.NotEquals) &&
      distance
    ) {
      return `${value} (${distance})`;
    } else {
      return `${value}`;
    }
  }

  public toCriterionInput(): PhashDistanceCriterionInput {
    return {
      value: this.value.value,
      modifier: this.modifier,
      distance: this.value.distance,
    };
  }
}

export const DuplicatedCriterionOption = new CriterionOption({
  messageID: "duplicated",
  type: "duplicated",
  makeCriterion: () => new DuplicatedCriterion(),
});

export class DuplicatedCriterion extends Criterion<IDuplicationValue> {
  public value: IDuplicationValue = {};

  constructor() {
    super(DuplicatedCriterionOption);
  }

  public clone(): DuplicatedCriterion {
    const c = new DuplicatedCriterion();
    c.value = { ...this.value };
    return c;
  }

  public cloneValues() {
    this.value = { ...this.value };
  }

  public getLabel(intl: IntlShape): string {
    const parts: string[] = [];
    const trueLabel = intl.formatMessage({ id: "true" });
    const falseLabel = intl.formatMessage({ id: "false" });

    if (this.value.phash !== undefined) {
      const label = intl.formatMessage({ id: "media_info.phash" });
      parts.push(`${label}: ${this.value.phash ? trueLabel : falseLabel}`);
    }
    if (this.value.stash_id !== undefined) {
      const label = intl.formatMessage({ id: "stash_id" });
      parts.push(`${label}: ${this.value.stash_id ? trueLabel : falseLabel}`);
    }
    if (this.value.title !== undefined) {
      const label = intl.formatMessage({ id: "title" });
      parts.push(`${label}: ${this.value.title ? trueLabel : falseLabel}`);
    }
    if (this.value.url !== undefined) {
      const label = intl.formatMessage({ id: "url" });
      parts.push(`${label}: ${this.value.url ? trueLabel : falseLabel}`);
    }

    // Handle legacy duplicated field
    if (parts.length === 0 && this.value.duplicated !== undefined) {
      const label = intl.formatMessage({ id: "duplicated_phash" });
      return `${label}: ${this.value.duplicated ? trueLabel : falseLabel}`;
    }

    if (parts.length === 0) {
      return intl.formatMessage({ id: "duplicated" });
    }

    return parts.join(", ");
  }

  public toCriterionInput(): DuplicationCriterionInput {
    return {
      duplicated: this.value.duplicated,
      distance: this.value.distance,
      phash: this.value.phash,
      url: this.value.url,
      stash_id: this.value.stash_id,
      title: this.value.title,
    };
  }

  public toQueryParams(): Record<string, unknown> {
    return {
      duplicated: this.value,
    };
  }

  public fromDecodedParams(o: Record<string, unknown>): void {
    const params = o as { duplicated?: IDuplicationValue };
    if (params.duplicated) {
      this.value = params.duplicated;
    }
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public setFromSavedCriterion(criterion: any): void {
    // Handle various saved formats
    if (criterion.value !== undefined) {
      // New format: { value: { phash: true, ... } }
      if (typeof criterion.value === "object") {
        this.value = criterion.value as IDuplicationValue;
      } else if (typeof criterion.value === "string") {
        // Legacy format: { value: "true" } - convert to phash
        this.value = { phash: criterion.value === "true" };
      }
    } else if (typeof criterion === "object") {
      // Direct value format
      this.value = criterion as IDuplicationValue;
    }
  }

  public isValid(): boolean {
    return (
      this.value.phash !== undefined ||
      this.value.url !== undefined ||
      this.value.stash_id !== undefined ||
      this.value.title !== undefined ||
      this.value.duplicated !== undefined
    );
  }
}
