import * as GQL from "src/core/generated-graphql";
import isEqual from "lodash-es/isEqual";

interface IHasRating {
  rating100?: GQL.Maybe<number> | undefined;
}

export function getAggregateRating(state: IHasRating[]) {
  let ret: number | undefined;
  let first = true;

  state.forEach((o) => {
    if (first) {
      ret = o.rating100 ?? undefined;
      first = false;
    } else if (ret !== o.rating100) {
      ret = undefined;
    }
  });

  return ret;
}

interface IHasID {
  id: string;
}

interface IHasStudio {
  studio?: GQL.Maybe<IHasID> | undefined;
}

export function getAggregateStudioId(state: IHasStudio[]) {
  let ret: string | undefined;
  let first = true;

  state.forEach((o) => {
    if (first) {
      ret = o?.studio?.id;
      first = false;
    } else {
      const studio = o?.studio?.id;
      if (ret !== studio) {
        ret = undefined;
      }
    }
  });

  return ret;
}

export function getAggregateIds(sortedLists: string[][]) {
  let ret: string[] = [];
  let first = true;

  sortedLists.forEach((l) => {
    if (first) {
      ret = l;
      first = false;
    } else {
      if (!isEqual(ret, l)) {
        ret = [];
      }
    }
  });

  return ret;
}

export function getAggregateGalleryIds(state: { galleries: IHasID[] }[]) {
  const sortedLists = state.map((o) => o.galleries.map((oo) => oo.id).sort());
  return getAggregateIds(sortedLists);
}

export function getAggregatePerformerIds(state: { performers: IHasID[] }[]) {
  const sortedLists = state.map((o) => o.performers.map((oo) => oo.id).sort());
  return getAggregateIds(sortedLists);
}

export function getAggregateTagIds(state: { tags: IHasID[] }[]) {
  const sortedLists = state.map((o) => o.tags.map((oo) => oo.id).sort());
  return getAggregateIds(sortedLists);
}

interface IGroup {
  movie: IHasID;
}

export function getAggregateGroupIds(state: { movies: IGroup[] }[]) {
  const sortedLists = state.map((o) =>
    o.movies.map((oo) => oo.movie.id).sort()
  );
  return getAggregateIds(sortedLists);
}

export function makeBulkUpdateIds(
  ids: string[],
  mode: GQL.BulkUpdateIdMode
): GQL.BulkUpdateIds {
  return {
    mode,
    ids,
  };
}

export function getAggregateInputValue<V>(
  inputValue: V | null | undefined,
  aggregateValue: V | null | undefined
) {
  if (inputValue === undefined) {
    // and all objects have the same value, then we are unsetting the value.
    if (aggregateValue !== undefined) {
      // null to unset rating
      return null;
    }
    // otherwise not setting the rating
    return undefined;
  } else {
    // if value is set, then we are setting the value for all
    return inputValue;
  }
}

// TODO - remove - this is incorrect
export function getAggregateInputIDs(
  mode: GQL.BulkUpdateIdMode,
  inputIds: string[] | undefined,
  aggregateIds: string[]
) {
  if (
    mode === GQL.BulkUpdateIdMode.Set &&
    (!inputIds || inputIds.length === 0)
  ) {
    // and all scenes have the same ids,
    if (aggregateIds.length > 0) {
      // then unset the performerIds, otherwise ignore
      return makeBulkUpdateIds(inputIds || [], mode);
    }
  } else {
    // if performerIds non-empty, then we are setting them
    return makeBulkUpdateIds(inputIds || [], mode);
  }

  return undefined;
}

export function getAggregateState<T, U>(
  currentValue: T,
  newValue: U,
  first: boolean
) {
  if (!first && !isEqual(currentValue, newValue)) {
    return undefined;
  }

  return newValue;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function setProperty<T, K extends keyof T>(obj: T, key: K, value: any) {
  obj[key] = value;
}

function getProperty<T, K extends keyof T>(obj: T, key: K) {
  return obj[key];
}

export function getAggregateStateObject<O, I>(
  output: O,
  input: I,
  fields: string[],
  first: boolean
) {
  fields.forEach((key) => {
    const outputKey = key as keyof O;
    const inputKey = key as keyof I;

    const currentValue = getProperty(output, outputKey);
    const performerValue = getProperty(input, inputKey);

    setProperty(
      output,
      outputKey,
      getAggregateState(currentValue, performerValue, first)
    );
  });
}
