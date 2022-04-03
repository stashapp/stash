import * as GQL from "src/core/generated-graphql";
import _ from "lodash";

interface IHasRating {
  rating?: GQL.Maybe<number> | undefined;
}

export function getAggregateRating(state: IHasRating[]) {
  let ret: number | undefined;
  let first = true;

  state.forEach((o) => {
    if (first) {
      ret = o.rating ?? undefined;
      first = false;
    } else if (ret !== o.rating) {
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

interface IHasPerformers {
  performers: IHasID[];
}

export function getAggregatePerformerIds(state: IHasPerformers[]) {
  let ret: string[] = [];
  let first = true;

  state.forEach((o) => {
    if (first) {
      ret = o.performers ? o.performers.map((p) => p.id).sort() : [];
      first = false;
    } else {
      const perfIds = o.performers ? o.performers.map((p) => p.id).sort() : [];

      if (!_.isEqual(ret, perfIds)) {
        ret = [];
      }
    }
  });

  return ret;
}

interface IHasTags {
  tags: IHasID[];
}

export function getAggregateTagIds(state: IHasTags[]) {
  let ret: string[] = [];
  let first = true;

  state.forEach((o) => {
    if (first) {
      ret = o.tags ? o.tags.map((t) => t.id).sort() : [];
      first = false;
    } else {
      const tIds = o.tags ? o.tags.map((t) => t.id).sort() : [];

      if (!_.isEqual(ret, tIds)) {
        ret = [];
      }
    }
  });

  return ret;
}

interface IMovie {
  movie: IHasID;
}

interface IHasMovies {
  movies: IMovie[];
}

export function getAggregateMovieIds(state: IHasMovies[]) {
  let ret: string[] = [];
  let first = true;

  state.forEach((o) => {
    if (first) {
      ret = o.movies ? o.movies.map((m) => m.movie.id).sort() : [];
      first = false;
    } else {
      const mIds = o.movies ? o.movies.map((m) => m.movie.id).sort() : [];

      if (!_.isEqual(ret, mIds)) {
        ret = [];
      }
    }
  });

  return ret;
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

export function getAggregateState<T>(
  currentValue: T,
  newValue: T,
  first: boolean
) {
  if (!first && !_.isEqual(currentValue, newValue)) {
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
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const performerValue = getProperty(input, inputKey) as any;

    setProperty(
      output,
      outputKey,
      getAggregateState(currentValue, performerValue, first)
    );
  });
}
