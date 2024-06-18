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
  useTagCreate,
  queryFindTagsByIDForSelect,
  queryFindTagsForSelect,
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
import { TagPopover } from "./TagPopover";
import { Placement } from "react-bootstrap/esm/Overlay";
import { sortByRelevance } from "src/utils/query";
import { PatchComponent, PatchFunction } from "src/patch";

export type SelectObject = {
  id: string;
  name?: string | null;
  title?: string | null;
};

export type Tag = Pick<GQL.Tag, "id" | "name" | "aliases" | "image_path">;
type Option = SelectOption<Tag>;

type FindTagsResult = Awaited<
  ReturnType<typeof queryFindTagsForSelect>
>["data"]["findTags"]["tags"];

function sortTagsByRelevance(input: string, tags: FindTagsResult) {
  return sortByRelevance(
    input,
    tags,
    (t) => t.name,
    (t) => t.aliases
  );
}

const tagSelectSort = PatchFunction("TagSelect.sort", sortTagsByRelevance);

const _TagSelect: React.FC<
  IFilterProps &
    IFilterValueProps<Tag> & {
      hoverPlacement?: Placement;
      excludeIds?: string[];
    }
> = (props) => {
  const [createTag] = useTagCreate();

  const { configuration } = React.useContext(ConfigurationContext);
  const intl = useIntl();
  const maxOptionsShown =
    configuration?.ui.maxOptionsShown ?? defaultMaxOptionsShown;
  const defaultCreatable =
    !configuration?.interface.disableDropdownCreate.tag ?? true;

  const exclude = useMemo(() => props.excludeIds ?? [], [props.excludeIds]);

  async function loadTags(input: string): Promise<Option[]> {
    const filter = new ListFilterModel(GQL.FilterMode.Tags);
    filter.searchTerm = input;
    filter.currentPage = 1;
    filter.itemsPerPage = maxOptionsShown;
    filter.sortBy = "name";
    filter.sortDirection = GQL.SortDirectionEnum.Asc;
    const query = await queryFindTagsForSelect(filter);
    let ret = query.data.findTags.tags.filter((tag) => {
      // HACK - we should probably exclude these in the backend query, but
      // this will do in the short-term
      return !exclude.includes(tag.id.toString());
    });

    return tagSelectSort(input, ret).map((tag) => ({
      value: tag.id,
      object: tag,
    }));
  }

  const TagOption: React.FC<OptionProps<Option, boolean>> = (optionProps) => {
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
        <TagPopover id={object.id} placement={props.hoverPlacement ?? "right"}>
          <span className="react-select-image-option">
            {/* the following code causes re-rendering issues when selecting tags */}
            {/* <TagPopover
              id={object.id}
              placement={props.hoverPlacement}
              target={targetRef}
            >
              <a
                href={`/tags/${object.id}`}
                target="_blank"
                rel="noreferrer"
                className="tag-select-image-link"
              >
                <img
                  className="tag-select-image"
                  src={object.image_path ?? ""}
                  loading="lazy"
                />
              </a>
            </TagPopover> */}
            <span>{name}</span>
            {alias && <span className="alias">{` (${alias})`}</span>}
          </span>
        </TagPopover>
      ),
    };

    return <reactSelectComponents.Option {...thisOptionProps} />;
  };

  const TagMultiValueLabel: React.FC<
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

  const TagValueLabel: React.FC<SingleValueProps<Option, boolean>> = (
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
    const result = await createTag({
      variables: { input: { name } },
    });
    return {
      value: result.data!.tagCreate!.id,
      item: result.data!.tagCreate!,
      message: "Created tag",
    };
  };

  const getNamedObject = (id: string, name: string) => {
    return {
      id,
      name,
      aliases: [],
    };
  };

  const isValidNewOption = (inputValue: string, options: Tag[]) => {
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
    <FilterSelectComponent<Tag, boolean>
      {...props}
      className={cx(
        "tag-select",
        {
          "tag-select-active": props.active,
        },
        props.className
      )}
      loadOptions={loadTags}
      getNamedObject={getNamedObject}
      isValidNewOption={isValidNewOption}
      components={{
        Option: TagOption,
        MultiValueLabel: TagMultiValueLabel,
        SingleValue: TagValueLabel,
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
              id: props.isMulti ? "tags" : "tag",
            }),
          }
        )
      }
      closeMenuOnSelect={!props.isMulti}
    />
  );
};

export const TagSelect = PatchComponent("TagSelect", _TagSelect);

const _TagIDSelect: React.FC<IFilterProps & IFilterIDProps<Tag>> = (props) => {
  const { ids, onSelect: onSelectValues } = props;

  const [values, setValues] = useState<Tag[]>([]);
  const idsChanged = useCompare(ids);

  function onSelect(items: Tag[]) {
    setValues(items);
    onSelectValues?.(items);
  }

  async function loadObjectsByID(idsToLoad: string[]): Promise<Tag[]> {
    const query = await queryFindTagsByIDForSelect(idsToLoad);
    const { tags: loadedTags } = query.data.findTags;

    return loadedTags;
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

  return <TagSelect {...props} values={values} onSelect={onSelect} />;
};

export const TagIDSelect = PatchComponent("TagIDSelect", _TagIDSelect);
