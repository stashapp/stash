import lodashIsEqual from "lodash-es/isEqual";
import clone from "lodash-es/clone";
import { IHasStoredID } from "src/utils/data";

export class ScrapeResult<T> {
  public newValue?: T;
  public originalValue?: T;
  public scraped: boolean = false;
  public useNewValue: boolean = false;
  private isEqual: (
    v1: T | undefined | null,
    v2: T | undefined | null
  ) => boolean;

  public constructor(
    originalValue?: T | null,
    newValue?: T | null,
    useNewValue?: boolean,
    isEqual: (
      v1: T | undefined | null,
      v2: T | undefined | null
    ) => boolean = lodashIsEqual
  ) {
    this.originalValue = originalValue ?? undefined;
    this.newValue = newValue ?? undefined;
    this.isEqual = isEqual;

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
    ret.useNewValue = !this.isEqual(ret.newValue, ret.originalValue);

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
    useNewValue?: boolean,
    isEqual: (
      v1: T | undefined | null,
      v2: T | undefined | null
    ) => boolean = lodashIsEqual
  ) {
    super(originalValue, newValue, useNewValue, isEqual);

    const hasNewValue = this.newValue !== undefined;

    const valuesEqual = isEqual(originalValue, newValue);
    this.useNewValue = useNewValue ?? (hasNewValue && !valuesEqual);
    this.scraped = hasNewValue && !valuesEqual;
  }
}

function storedIDsEqual<T extends IHasStoredID>(
  o1: T[] | undefined | null,
  o2: T[] | undefined | null
) {
  return (
    !!o1 &&
    !!o2 &&
    o1.length === o2.length &&
    o1.every((o) => {
      return o2.find((oo) => o.stored_id === oo.stored_id);
    })
  );
}

export class ObjectListScrapeResult<
  T extends IHasStoredID
> extends ScrapeResult<T[]> {
  public constructor(
    originalValue?: T[] | null,
    newValue?: T[] | null,
    useNewValue?: boolean
  ) {
    super(originalValue, newValue, useNewValue, storedIDsEqual);
  }
}

export class ObjectScrapeResult<
  T extends IHasStoredID
> extends ScrapeResult<T> {
  public constructor(
    originalValue?: T | null,
    newValue?: T | null,
    useNewValue?: boolean
  ) {
    super(
      originalValue,
      newValue,
      useNewValue,
      (o1, o2) => o1?.stored_id === o2?.stored_id
    );
  }
}

export function hasScrapedValues(values: { scraped: boolean }[]): boolean {
  return values.some((r) => r.scraped);
}
