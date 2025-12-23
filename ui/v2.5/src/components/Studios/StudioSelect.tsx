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
  useStudioCreate,
  queryFindStudiosByIDForSelect,
  queryFindStudiosForSelect,
} from "src/core/StashService";
import { useConfigurationContext } from "src/hooks/Config";
import { useIntl } from "react-intl";
import { defaultMaxOptionsShown, IUIConfig } from "src/core/config";
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

export type SelectObject = {
  id: string;
  name?: string | null;
  title?: string | null;
};

export type Studio = Pick<GQL.Studio, "id" | "name" | "aliases" | "image_path">;
type Option = SelectOption<Studio>;

type FindStudiosResult = Awaited<
  ReturnType<typeof queryFindStudiosForSelect>
>["data"]["findStudios"]["studios"];

function sortStudiosByRelevance(input: string, studios: FindStudiosResult) {
  return sortByRelevance(
    input,
    studios,
    (s) => s.name,
    (s) => s.aliases
  );
}

const studioSelectSort = PatchFunction(
  "StudioSelect.sort",
  sortStudiosByRelevance
);

const _StudioSelect: React.FC<
  IFilterProps &
    IFilterValueProps<Studio> & {
      hoverPlacement?: Placement;
      excludeIds?: string[];
    }
> = (props) => {
  const [createStudio] = useStudioCreate();

  const { configuration } = useConfigurationContext();
  const intl = useIntl();
  const maxOptionsShown =
    (configuration?.ui as IUIConfig).maxOptionsShown ?? defaultMaxOptionsShown;
  const defaultCreatable =
    !configuration?.interface.disableDropdownCreate.studio;

  const exclude = useMemo(() => props.excludeIds ?? [], [props.excludeIds]);

  async function loadStudios(input: string): Promise<Option[]> {
    const filter = new ListFilterModel(GQL.FilterMode.Studios);
    filter.searchTerm = input;
    filter.currentPage = 1;
    filter.itemsPerPage = maxOptionsShown;
    filter.sortBy = "name";
    filter.sortDirection = GQL.SortDirectionEnum.Asc;
    const query = await queryFindStudiosForSelect(filter);
    let ret = query.data.findStudios.studios.filter((studio) => {
      // HACK - we should probably exclude these in the backend query, but
      // this will do in the short-term
      return !exclude.includes(studio.id.toString());
    });

    return studioSelectSort(input, ret).map((studio) => ({
      value: studio.id,
      object: studio,
    }));
  }

  const StudioOption: React.FC<OptionProps<Option, boolean>> = (
    optionProps
  ) => {
    let thisOptionProps = optionProps;

    const { object } = optionProps.data;

    let { name } = object;

    // if name does not match the input value but an alias does, show the alias
    const { inputValue } = optionProps.selectProps;
    let alias: string | undefined = "";
    if (!name.toLowerCase().includes(inputValue.toLowerCase())) {
      alias = object.aliases?.find((a) =>
        a.toLowerCase().includes(inputValue.toLowerCase())
      );
    }

    thisOptionProps = {
      ...optionProps,
      children: (
        <span className="react-select-image-option">
          <span>{name}</span>
          {alias && <span className="alias">&nbsp;({alias})</span>}
        </span>
      ),
    };

    return <reactSelectComponents.Option {...thisOptionProps} />;
  };

  const StudioMultiValueLabel: React.FC<
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

  const StudioValueLabel: React.FC<SingleValueProps<Option, boolean>> = (
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
    const result = await createStudio({
      variables: { input: { name } },
    });
    return {
      value: result.data!.studioCreate!.id,
      item: result.data!.studioCreate!,
      message: "Created studio",
    };
  };

  const getNamedObject = (id: string, name: string) => {
    return {
      id,
      name,
      aliases: [],
    };
  };

  const isValidNewOption = (inputValue: string, options: Studio[]) => {
    if (!inputValue) {
      return false;
    }

    if (
      options.some((o) => {
        return (
          o.name.toLowerCase() === inputValue.toLowerCase() ||
          o.aliases?.some((a) => a.toLowerCase() === inputValue.toLowerCase())
        );
      })
    ) {
      return false;
    }

    return true;
  };

  return (
    <FilterSelectComponent<Studio, boolean>
      {...props}
      className={cx(
        "studio-select",
        {
          "studio-select-active": props.active,
        },
        props.className
      )}
      loadOptions={loadStudios}
      getNamedObject={getNamedObject}
      isValidNewOption={isValidNewOption}
      components={{
        Option: StudioOption,
        MultiValueLabel: StudioMultiValueLabel,
        SingleValue: StudioValueLabel,
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
              id: props.isMulti ? "studios" : "studio",
            }),
          }
        )
      }
      closeMenuOnSelect={!props.isMulti}
    />
  );
};

export const StudioSelect = PatchComponent("StudioSelect", _StudioSelect);

const _StudioIDSelect: React.FC<IFilterProps & IFilterIDProps<Studio>> = (
  props
) => {
  const { ids, onSelect: onSelectValues } = props;

  const [values, setValues] = useState<Studio[]>([]);
  const idsChanged = useCompare(ids);

  function onSelect(items: Studio[]) {
    setValues(items);
    onSelectValues?.(items);
  }

  async function loadObjectsByID(idsToLoad: string[]): Promise<Studio[]> {
    const query = await queryFindStudiosByIDForSelect(idsToLoad);
    const { studios: loadedStudios } = query.data.findStudios;

    return loadedStudios;
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

  return <StudioSelect {...props} values={values} onSelect={onSelect} />;
};

export const StudioIDSelect = PatchComponent("StudioIDSelect", _StudioIDSelect);
