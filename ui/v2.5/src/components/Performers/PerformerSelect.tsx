import React, { useEffect, useState } from "react";
import {
  OptionProps,
  components as reactSelectComponents,
  MultiValueGenericProps,
  SingleValueProps,
} from "react-select";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import {
  usePerformerCreate,
  queryFindPerformersByIDForSelect,
  queryFindPerformersForSelect,
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
import { Link } from "react-router-dom";
import { sortByRelevance } from "src/utils/query";
import { PatchComponent, PatchFunction } from "src/patch";

export type SelectObject = {
  id: string;
  name?: string | null;
  title?: string | null;
};

export type Performer = Pick<
  GQL.Performer,
  "id" | "name" | "alias_list" | "disambiguation" | "image_path"
>;
type Option = SelectOption<Performer>;

type FindPerformersResult = Awaited<
  ReturnType<typeof queryFindPerformersForSelect>
>["data"]["findPerformers"]["performers"];

function sortPerformersByRelevance(
  input: string,
  performers: FindPerformersResult
) {
  return sortByRelevance(
    input,
    performers,
    (p) => p.name,
    (p) => p.alias_list
  );
}

const performerSelectSort = PatchFunction(
  "PerformerSelect.sort",
  sortPerformersByRelevance
);

const _PerformerSelect: React.FC<
  IFilterProps & IFilterValueProps<Performer>
> = (props) => {
  const [createPerformer] = usePerformerCreate();

  const { configuration } = React.useContext(ConfigurationContext);
  const intl = useIntl();
  const maxOptionsShown =
    configuration?.ui.maxOptionsShown ?? defaultMaxOptionsShown;
  const defaultCreatable =
    !configuration?.interface.disableDropdownCreate.performer ?? true;

  async function loadPerformers(input: string): Promise<Option[]> {
    const filter = new ListFilterModel(GQL.FilterMode.Performers);
    filter.searchTerm = input;
    filter.currentPage = 1;
    filter.itemsPerPage = maxOptionsShown;
    filter.sortBy = "name";
    filter.sortDirection = GQL.SortDirectionEnum.Asc;
    const query = await queryFindPerformersForSelect(filter);
    return performerSelectSort(
      input,
      query.data.findPerformers.performers.slice()
    ).map((performer) => ({
      value: performer.id,
      object: performer,
    }));
  }

  const PerformerOption: React.FC<OptionProps<Option, boolean>> = (
    optionProps
  ) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    let { name } = object;

    // if name does not match the input value but an alias does, show the alias
    const { inputValue } = optionProps.selectProps;
    let alias: string | undefined = "";
    if (!name.toLowerCase().includes(inputValue.toLowerCase())) {
      alias = object.alias_list?.find((a) =>
        a.toLowerCase().includes(inputValue.toLowerCase())
      );
    }

    thisOptionProps = {
      ...optionProps,
      children: (
        <span className="react-select-image-option performer-select-option">
          <Link
            to={`/performers/${object.id}`}
            target="_blank"
            className="performer-select-image-link"
          >
            <img
              className="performer-select-image"
              src={object.image_path ?? ""}
              loading="lazy"
            />
          </Link>
          <span>{name}</span>
          {object.disambiguation && (
            <span className="performer-disambiguation">{` (${object.disambiguation})`}</span>
          )}
          {alias && <span className="alias">{` (${alias})`}</span>}
        </span>
      ),
    };

    return <reactSelectComponents.Option {...thisOptionProps} />;
  };

  const PerformerMultiValueLabel: React.FC<
    MultiValueGenericProps<Option, boolean>
  > = (optionProps) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    thisOptionProps = {
      ...optionProps,
      children: (
        <span className="performer-select-value">
          <span>{object.name}</span>
          {object.disambiguation && (
            <span className="performer-disambiguation">{` (${object.disambiguation})`}</span>
          )}
        </span>
      ),
    };

    return <reactSelectComponents.MultiValueLabel {...thisOptionProps} />;
  };

  const PerformerValueLabel: React.FC<SingleValueProps<Option, boolean>> = (
    optionProps
  ) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    thisOptionProps = {
      ...optionProps,
      children: (
        <span className="performer-select-value">
          {object.name}
          {object.disambiguation && (
            <span className="performer-disambiguation">{` (${object.disambiguation})`}</span>
          )}
        </span>
      ),
    };

    return <reactSelectComponents.SingleValue {...thisOptionProps} />;
  };

  const onCreate = async (name: string) => {
    const result = await createPerformer({
      variables: { input: { name } },
    });
    return {
      value: result.data!.performerCreate!.id,
      item: result.data!.performerCreate!,
      message: "Created performer",
    };
  };

  const getNamedObject = (id: string, name: string) => {
    return {
      id,
      name,
      alias_list: [],
    };
  };

  const isValidNewOption = (inputValue: string, options: Performer[]) => {
    if (!inputValue) {
      return false;
    }

    if (
      options.some((o) => {
        return (
          o.name.toLowerCase() === inputValue.toLowerCase() ||
          o.alias_list?.some(
            (a) => a.toLowerCase() === inputValue.toLowerCase()
          )
        );
      })
    ) {
      return false;
    }

    return true;
  };

  return (
    <FilterSelectComponent<Performer, boolean>
      {...props}
      className={cx(
        "performer-select",
        {
          "performer-select-active": props.active,
        },
        props.className
      )}
      loadOptions={loadPerformers}
      getNamedObject={getNamedObject}
      isValidNewOption={isValidNewOption}
      components={{
        Option: PerformerOption,
        MultiValueLabel: PerformerMultiValueLabel,
        SingleValue: PerformerValueLabel,
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
              id: props.isMulti ? "performers" : "performer",
            }),
          }
        )
      }
    />
  );
};

export const PerformerSelect = PatchComponent(
  "PerformerSelect",
  _PerformerSelect
);

const _PerformerIDSelect: React.FC<IFilterProps & IFilterIDProps<Performer>> = (
  props
) => {
  const { ids, onSelect: onSelectValues } = props;

  const [values, setValues] = useState<Performer[]>([]);
  const idsChanged = useCompare(ids);

  function onSelect(items: Performer[]) {
    setValues(items);
    onSelectValues?.(items);
  }

  async function loadObjectsByID(idsToLoad: string[]): Promise<Performer[]> {
    const query = await queryFindPerformersByIDForSelect(idsToLoad);
    const { performers: loadedPerformers } = query.data.findPerformers;

    return loadedPerformers;
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

  return <PerformerSelect {...props} values={values} onSelect={onSelect} />;
};

export const PerformerIDSelect = PatchComponent(
  "PerformerIDSelect",
  _PerformerIDSelect
);
