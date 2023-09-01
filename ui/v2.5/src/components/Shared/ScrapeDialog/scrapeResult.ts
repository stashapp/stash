import isEqual from "lodash-es/isEqual";
import clone from "lodash-es/clone";

export class ScrapeResult<T> {
  public newValue?: T;
  public originalValue?: T;
  public scraped: boolean = false;
  public useNewValue: boolean = false;

  public constructor(
    originalValue?: T | null,
    newValue?: T | null,
    useNewValue?: boolean
  ) {
    this.originalValue = originalValue ?? undefined;
    this.newValue = newValue ?? undefined;
    // NOTE: this means that zero values are treated as null
    // this is incorrect for numbers and booleans, but correct for strings
    const hasNewValue = !!this.newValue;

    const valuesEqual = isEqual(originalValue, newValue);
    this.useNewValue = useNewValue ?? (hasNewValue && !valuesEqual);
    this.scraped = hasNewValue && !valuesEqual;
  }

  public setOriginalValue(value?: T) {
    this.originalValue = value;
    this.newValue = value;
  }

  public cloneWithValue(value?: T) {
    const ret = clone(this);

    ret.newValue = value;
    ret.useNewValue = !isEqual(ret.newValue, ret.originalValue);

    // #2691 - if we're setting the value, assume it should be treated as
    // scraped
    ret.scraped = true;

    return ret;
  }

  public getNewValue() {
    if (this.useNewValue) {
      return this.newValue;
    }
  }
}

// for types where !!value is a valid value (boolean and number)
export class ZeroableScrapeResult<T> extends ScrapeResult<T> {
  public constructor(
    originalValue?: T | null,
    newValue?: T | null,
    useNewValue?: boolean
  ) {
    super(originalValue, newValue, useNewValue);

    const hasNewValue = this.newValue !== undefined;

    const valuesEqual = isEqual(originalValue, newValue);
    this.useNewValue = useNewValue ?? (hasNewValue && !valuesEqual);
    this.scraped = hasNewValue && !valuesEqual;
  }
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function hasScrapedValues(values: ScrapeResult<any>[]) {
  return values.some((r) => r.scraped);
}
