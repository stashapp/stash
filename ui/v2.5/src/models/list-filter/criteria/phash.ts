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

// Shared mapping of duplication field IDs to their i18n message IDs
export const DUPLICATION_FIELD_MESSAGE_IDS = {
  phash: "media_info.phash",
  stash_id: "stash_id",
  title: "title",
  url: "url",
} as const;

export type DuplicationFieldId = keyof typeof DUPLICATION_FIELD_MESSAGE_IDS;

export const DUPLICATION_FIELD_IDS: DuplicationFieldId[] = [
  "phash",
  "stash_id",
  "title",
  "url",
];

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

    for (const fieldId of DUPLICATION_FIELD_IDS) {
      const fieldValue = this.value[fieldId];
      if (fieldValue !== undefined) {
        const label = intl.formatMessage({
          id: DUPLICATION_FIELD_MESSAGE_IDS[fieldId],
        });
        parts.push(`${label}: ${fieldValue ? trueLabel : falseLabel}`);
      }
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
    // Check if any duplication field is set
    const hasFieldSet = DUPLICATION_FIELD_IDS.some(
      (fieldId) => this.value[fieldId] !== undefined
    );
    return hasFieldSet || this.value.duplicated !== undefined;
  }
}
