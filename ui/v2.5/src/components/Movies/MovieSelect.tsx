import React, { useEffect, useMemo, useState } from "react";
import {
  OptionProps,
  components as reactSelectComponents,
  MultiValueGenericProps,
  SingleValueProps,
} from "react-select";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import {
  queryFindMovies,
  // queryFindMoviesByIDForSelect,
} from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { useIntl } from "react-intl";
import { defaultMaxOptionsShown } from "src/core/config";
import { ListFilterModel } from "src/models/list-filter/filter";
import {
  FilterSelectComponent,
  IFilterIDProps,
  IFilterProps,
  IFilterValueProps,
  Option as SelectOption,
} from "../Shared/FilterSelect";
import { useCompare } from "src/hooks/state";
import { Placement } from "react-bootstrap/esm/Overlay";
import { sortByRelevance } from "src/utils/query";
import { PatchComponent } from "src/pluginApi";

export type Movie = Pick<GQL.Movie, "id" | "name">;
type Option = SelectOption<Movie>;

const _MovieSelect: React.FC<
  IFilterProps &
    IFilterValueProps<Movie> & {
      hoverPlacement?: Placement;
      excludeIds?: string[];
    }
> = (props) => {
  const { configuration } = React.useContext(ConfigurationContext);
  const intl = useIntl();
  const maxOptionsShown =
    configuration?.ui.maxOptionsShown ?? defaultMaxOptionsShown;

  const exclude = useMemo(() => props.excludeIds ?? [], [props.excludeIds]);

  async function loadMovies(input: string): Promise<Option[]> {
    const filter = new ListFilterModel(GQL.FilterMode.Movies);
    filter.searchTerm = input;
    filter.currentPage = 1;
    filter.itemsPerPage = maxOptionsShown;
    filter.sortBy = "title";
    filter.sortDirection = GQL.SortDirectionEnum.Asc;
    const query = await queryFindMovies(filter);
    let ret = query.data.findMovies.movies.filter((movie) => {
      // HACK - we should probably exclude these in the backend query, but
      // this will do in the short-term
      return !exclude.includes(movie.id.toString());
    });

    return sortByRelevance(input, ret, (m) => m.name).map((movie) => ({
      value: movie.id,
      object: movie,
    }));
  }

  const MovieOption: React.FC<OptionProps<Option, boolean>> = (optionProps) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    const title = object.name;

    thisOptionProps = {
      ...optionProps,
      children: <span>{title}</span>,
    };

    return <reactSelectComponents.Option {...thisOptionProps} />;
  };

  const MovieMultiValueLabel: React.FC<
    MultiValueGenericProps<Option, boolean>
  > = (optionProps) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    thisOptionProps = {
      ...optionProps,
      children: object.name,
    };

    return <reactSelectComponents.MultiValueLabel {...thisOptionProps} />;
  };

  const MovieValueLabel: React.FC<SingleValueProps<Option, boolean>> = (
    optionProps
  ) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    thisOptionProps = {
      ...optionProps,
      children: <>{object.name}</>,
    };

    return <reactSelectComponents.SingleValue {...thisOptionProps} />;
  };

  return (
    <FilterSelectComponent<Movie, boolean>
      {...props}
      className={cx(
        "movie-select",
        {
          "movie-select-active": props.active,
        },
        props.className
      )}
      loadOptions={loadMovies}
      components={{
        Option: MovieOption,
        MultiValueLabel: MovieMultiValueLabel,
        SingleValue: MovieValueLabel,
      }}
      isMulti={props.isMulti ?? false}
      placeholder={
        props.noSelectionString ??
        intl.formatMessage(
          { id: "actions.select_entity" },
          {
            entityType: intl.formatMessage({
              id: props.isMulti ? "movies" : "movie",
            }),
          }
        )
      }
      closeMenuOnSelect={!props.isMulti}
    />
  );
};

export const MovieSelect = PatchComponent("MovieSelect", _MovieSelect);

const _MovieIDSelect: React.FC<IFilterProps & IFilterIDProps<Movie>> = (
  props
) => {
  // const { ids, onSelect: onSelectValues } = props;

  // const [values, setValues] = useState<Movie[]>([]);
  // const idsChanged = useCompare(ids);

  // function onSelect(items: Movie[]) {
  //   setValues(items);
  //   onSelectValues?.(items);
  // }

  // async function loadObjectsByID(idsToLoad: string[]): Promise<Movie[]> {
  //   const movieIDs = idsToLoad.map((id) => parseInt(id));
  //   const query = await queryFindMoviesByIDForSelect(movieIDs);
  //   const { movies: loadedMovies } = query.data.findMovies;

  //   return loadedMovies;
  // }

  // useEffect(() => {
  //   if (!idsChanged) {
  //     return;
  //   }

  //   if (!ids || ids?.length === 0) {
  //     setValues([]);
  //     return;
  //   }

  //   // load the values if we have ids and they haven't been loaded yet
  //   const filteredValues = values.filter((v) => ids.includes(v.id.toString()));
  //   if (filteredValues.length === ids.length) {
  //     return;
  //   }

  //   const load = async () => {
  //     const items = await loadObjectsByID(ids);
  //     setValues(items);
  //   };

  //   load();
  // }, [ids, idsChanged, values]);

  // return <MovieSelect {...props} values={values} onSelect={onSelect} />;
  return <></>;
};

export const MovieIDSelect = PatchComponent("MovieIDSelect", _MovieIDSelect);
