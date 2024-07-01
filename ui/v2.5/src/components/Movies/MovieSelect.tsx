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
  queryFindMoviesForSelect,
  queryFindMoviesByIDForSelect,
  useMovieCreate,
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
import { PatchComponent, PatchFunction } from "src/patch";
import { TruncatedText } from "../Shared/TruncatedText";

export type Group = Pick<
  GQL.Movie,
  "id" | "name" | "date" | "front_image_path" | "aliases"
> & {
  studio?: Pick<GQL.Studio, "name"> | null;
};
type Option = SelectOption<Group>;

type FindMoviesResult = Awaited<
  ReturnType<typeof queryFindMoviesForSelect>
>["data"]["findMovies"]["movies"];

function sortMoviesByRelevance(input: string, movies: FindMoviesResult) {
  return sortByRelevance(
    input,
    movies,
    (m) => m.name,
    (m) => (m.aliases ? [m.aliases] : [])
  );
}

const movieSelectSort = PatchFunction(
  "MovieSelect.sort",
  sortMoviesByRelevance
);

const _GroupSelect: React.FC<
  IFilterProps &
    IFilterValueProps<Group> & {
      hoverPlacement?: Placement;
      excludeIds?: string[];
    }
> = (props) => {
  const [createMovie] = useMovieCreate();

  const { configuration } = React.useContext(ConfigurationContext);
  const intl = useIntl();
  const maxOptionsShown =
    configuration?.ui.maxOptionsShown ?? defaultMaxOptionsShown;
  const defaultCreatable =
    !configuration?.interface.disableDropdownCreate.movie ?? true;

  const exclude = useMemo(() => props.excludeIds ?? [], [props.excludeIds]);

  async function loadMovies(input: string): Promise<Option[]> {
    const filter = new ListFilterModel(GQL.FilterMode.Movies);
    filter.searchTerm = input;
    filter.currentPage = 1;
    filter.itemsPerPage = maxOptionsShown;
    filter.sortBy = "name";
    filter.sortDirection = GQL.SortDirectionEnum.Asc;
    const query = await queryFindMoviesForSelect(filter);
    let ret = query.data.findMovies.movies.filter((movie) => {
      // HACK - we should probably exclude these in the backend query, but
      // this will do in the short-term
      return !exclude.includes(movie.id.toString());
    });

    return movieSelectSort(input, ret).map((movie) => ({
      value: movie.id,
      object: movie,
    }));
  }

  const GroupOption: React.FC<OptionProps<Option, boolean>> = (optionProps) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    const title = object.name;

    // if name does not match the input value but an alias does, show the alias
    const { inputValue } = optionProps.selectProps;
    let alias: string | undefined = "";
    if (!title.toLowerCase().includes(inputValue.toLowerCase())) {
      alias = object.aliases || undefined;
    }

    thisOptionProps = {
      ...optionProps,
      children: (
        <span className="group-select-option">
          <span className="group-select-row">
            {object.front_image_path && (
              <img
                className="group-select-image"
                src={object.front_image_path}
                loading="lazy"
              />
            )}

            <span className="group-select-details">
              <TruncatedText
                className="group-select-title"
                text={
                  <span>
                    {title}
                    {alias && (
                      <span className="group-select-alias">{` (${alias})`}</span>
                    )}
                  </span>
                }
                lineCount={1}
              />

              {object.studio?.name && (
                <span className="group-select-studio">
                  {object.studio?.name}
                </span>
              )}

              {object.date && (
                <span className="group-select-date">{object.date}</span>
              )}
            </span>
          </span>
        </span>
      ),
    };

    return <reactSelectComponents.Option {...thisOptionProps} />;
  };

  const GroupMultiValueLabel: React.FC<
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

  const GroupValueLabel: React.FC<SingleValueProps<Option, boolean>> = (
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

  const onCreate = async (name: string) => {
    const result = await createMovie({
      variables: { input: { name } },
    });
    return {
      value: result.data!.movieCreate!.id,
      item: result.data!.movieCreate!,
      message: "Created group",
    };
  };

  const getNamedObject = (id: string, name: string) => {
    return {
      id,
      name,
    };
  };

  const isValidNewOption = (inputValue: string, options: Group[]) => {
    if (!inputValue) {
      return false;
    }

    if (
      options.some((o) => {
        return (
          o.name.toLowerCase() === inputValue.toLowerCase() ||
          o.aliases?.toLowerCase() === inputValue.toLowerCase()
        );
      })
    ) {
      return false;
    }

    return true;
  };

  return (
    <FilterSelectComponent<Group, boolean>
      {...props}
      className={cx(
        "group-select",
        {
          "group-select-active": props.active,
        },
        props.className
      )}
      loadOptions={loadMovies}
      getNamedObject={getNamedObject}
      isValidNewOption={isValidNewOption}
      components={{
        Option: GroupOption,
        MultiValueLabel: GroupMultiValueLabel,
        SingleValue: GroupValueLabel,
      }}
      isMulti={props.isMulti ?? false}
      creatable={props.creatable ?? defaultCreatable}
      onCreate={onCreate}
      placeholder={
        props.noSelectionString ??
        intl.formatMessage(
          { id: "actions.select_entity" },
          {
            entityType: intl.formatMessage({
              id: props.isMulti ? "groups" : "group",
            }),
          }
        )
      }
      closeMenuOnSelect={!props.isMulti}
    />
  );
};

export const GroupSelect = PatchComponent("GroupSelect", _GroupSelect);

const _GroupIDSelect: React.FC<IFilterProps & IFilterIDProps<Group>> = (
  props
) => {
  const { ids, onSelect: onSelectValues } = props;

  const [values, setValues] = useState<Group[]>([]);
  const idsChanged = useCompare(ids);

  function onSelect(items: Group[]) {
    setValues(items);
    onSelectValues?.(items);
  }

  async function loadObjectsByID(idsToLoad: string[]): Promise<Group[]> {
    const query = await queryFindMoviesByIDForSelect(idsToLoad);
    const { movies: loadedMovies } = query.data.findMovies;

    return loadedMovies;
  }

  useEffect(() => {
    if (!idsChanged) {
      return;
    }

    if (!ids || ids?.length === 0) {
      setValues([]);
      return;
    }

    // load the values if we have ids and they haven't been loaded yet
    const filteredValues = values.filter((v) => ids.includes(v.id.toString()));
    if (filteredValues.length === ids.length) {
      return;
    }

    const load = async () => {
      const items = await loadObjectsByID(ids);
      setValues(items);
    };

    load();
  }, [ids, idsChanged, values]);

  return <GroupSelect {...props} values={values} onSelect={onSelect} />;
};

export const GroupIDSelect = PatchComponent("GroupIDSelect", _GroupIDSelect);
