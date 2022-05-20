import React, { useEffect, useMemo, useState } from "react";
import Select, {
  ValueType,
  Styles,
  OptionProps,
  components as reactSelectComponents,
  GroupedOptionsType,
  OptionsType,
} from "react-select";
import CreatableSelect from "react-select/creatable";
import { debounce } from "lodash";

import * as GQL from "src/core/generated-graphql";
import {
  useAllTagsForFilter,
  useAllMoviesForFilter,
  useAllStudiosForFilter,
  useAllPerformersForFilter,
  useMarkerStrings,
  useTagCreate,
  useStudioCreate,
  usePerformerCreate,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { TextUtils } from "src/utils";
import { SelectComponents } from "react-select/src/components";
import { ConfigurationContext } from "src/hooks/Config";
import { useIntl } from "react-intl";

export type ValidTypes =
  | GQL.SlimPerformerDataFragment
  | GQL.SlimTagDataFragment
  | GQL.SlimStudioDataFragment
  | GQL.SlimMovieDataFragment;
type Option = { value: string; label: string };

interface ITypeProps {
  type?:
    | "performers"
    | "studios"
    | "parent_studios"
    | "tags"
    | "sceneTags"
    | "performerTags"
    | "parentTags"
    | "childTags"
    | "movies";
}
interface IFilterProps {
  ids?: string[];
  initialIds?: string[];
  onSelect?: (item: ValidTypes[]) => void;
  noSelectionString?: string;
  className?: string;
  isMulti?: boolean;
  isClearable?: boolean;
  isDisabled?: boolean;
  creatable?: boolean;
  menuPortalTarget?: HTMLElement | null;
}
interface ISelectProps<T extends boolean> {
  className?: string;
  items: Option[];
  selectedOptions?: ValueType<Option, T>;
  creatable?: boolean;
  onCreateOption?: (value: string) => void;
  isLoading: boolean;
  isDisabled?: boolean;
  onChange: (item: ValueType<Option, T>) => void;
  initialIds?: string[];
  isMulti: T;
  isClearable?: boolean;
  onInputChange?: (input: string) => void;
  components?: Partial<SelectComponents<Option, T>>;
  filterOption?: (option: Option, rawInput: string) => boolean;
  isValidNewOption?: (
    inputValue: string,
    value: ValueType<Option, T>,
    options: OptionsType<Option> | GroupedOptionsType<Option>
  ) => boolean;
  placeholder?: string;
  showDropdown?: boolean;
  groupHeader?: string;
  menuPortalTarget?: HTMLElement | null;
  closeMenuOnSelect?: boolean;
  noOptionsMessage?: string | null;
}
interface IFilterComponentProps extends IFilterProps {
  items: Array<ValidTypes>;
  onCreate?: (name: string) => Promise<{ item: ValidTypes; message: string }>;
}
interface IFilterSelectProps<T extends boolean>
  extends Omit<ISelectProps<T>, "onChange" | "items" | "onCreateOption"> {}

type Gallery = { id: string; title: string };
interface IGallerySelect {
  galleries: Gallery[];
  onSelect: (items: Gallery[]) => void;
}

type Scene = { id: string; title: string };
interface ISceneSelect {
  scenes: Scene[];
  onSelect: (items: Scene[]) => void;
}

const getSelectedItems = (selectedItems: ValueType<Option, boolean>) =>
  selectedItems
    ? Array.isArray(selectedItems)
      ? selectedItems
      : [selectedItems]
    : [];

const getSelectedValues = (selectedItems: ValueType<Option, boolean>) =>
  getSelectedItems(selectedItems).map((item) => item.value);

const SelectComponent = <T extends boolean>({
  type,
  initialIds,
  onChange,
  className,
  items,
  selectedOptions,
  isLoading,
  isDisabled = false,
  onCreateOption,
  isClearable = true,
  creatable = false,
  isMulti,
  onInputChange,
  filterOption,
  isValidNewOption,
  components,
  placeholder,
  showDropdown = true,
  groupHeader,
  menuPortalTarget,
  closeMenuOnSelect = true,
  noOptionsMessage = type !== "tags" ? "None" : null,
}: ISelectProps<T> & ITypeProps) => {
  const values = items.filter((item) => initialIds?.indexOf(item.value) !== -1);
  const defaultValue = (isMulti ? values : values[0] ?? null) as ValueType<
    Option,
    T
  >;

  const options = groupHeader
    ? [
        {
          label: groupHeader,
          options: items,
        },
      ]
    : items;

  const styles: Partial<Styles<Option, T>> = {
    option: (base) => ({
      ...base,
      color: "#000",
    }),
    container: (base, props) => ({
      ...base,
      zIndex: props.selectProps.isFocused ? 10 : base.zIndex,
    }),
    multiValueRemove: (base, props) => ({
      ...base,
      color: props.selectProps.isFocused ? base.color : "#333333",
    }),
  };

  const props = {
    options,
    value: selectedOptions,
    className,
    classNamePrefix: "react-select",
    onChange,
    isMulti,
    isClearable,
    defaultValue: defaultValue ?? undefined,
    noOptionsMessage: () => noOptionsMessage,
    placeholder: isDisabled ? "" : placeholder,
    onInputChange,
    filterOption,
    isValidNewOption,
    isDisabled,
    isLoading,
    styles,
    closeMenuOnSelect,
    menuPortalTarget,
    components: {
      ...components,
      IndicatorSeparator: () => null,
      ...((!showDropdown || isDisabled) && { DropdownIndicator: () => null }),
      ...(isDisabled && { MultiValueRemove: () => null }),
    },
  };

  return creatable ? (
    <CreatableSelect
      {...props}
      isDisabled={isLoading || isDisabled}
      onCreateOption={onCreateOption}
    />
  ) : (
    <Select {...props} />
  );
};

const FilterSelectComponent = <T extends boolean>(
  props: IFilterComponentProps & ITypeProps & IFilterSelectProps<T>
) => {
  const { items, ids, isMulti, onSelect } = props;
  const [loading, setLoading] = useState(false);
  const selectedIds = ids ?? [];
  const Toast = useToast();

  const options = items.map((i) => ({
    value: i.id,
    label: i.name ?? "",
  }));

  const selected = options.filter((option) =>
    selectedIds.includes(option.value)
  );
  const selectedOptions = (isMulti
    ? selected
    : selected[0] ?? null) as ValueType<Option, T>;

  const onChange = (selectedItems: ValueType<Option, boolean>) => {
    const selectedValues = getSelectedValues(selectedItems);
    onSelect?.(items.filter((item) => selectedValues.includes(item.id)));
  };

  const onCreate = async (name: string) => {
    try {
      setLoading(true);
      const { item: newItem, message } = await props.onCreate!(name);
      props.onSelect?.([
        ...items.filter((item) => selectedIds.includes(item.id)),
        newItem,
      ]);
      setLoading(false);
      Toast.success({
        content: (
          <span>
            {message}: <b>{name}</b>
          </span>
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  };

  return (
    <SelectComponent<T>
      {...props}
      isLoading={props.isLoading || loading}
      onChange={onChange}
      items={options}
      selectedOptions={selectedOptions}
      onCreateOption={props.creatable ? onCreate : undefined}
    />
  );
};

export const GallerySelect: React.FC<IGallerySelect> = (props) => {
  const [query, setQuery] = useState<string>("");
  const { data, loading } = GQL.useFindGalleriesQuery({
    skip: query === "",
    variables: {
      filter: {
        q: query,
      },
    },
  });

  const galleries = data?.findGalleries.galleries ?? [];
  const items = galleries.map((g) => ({
    label: g.title ?? TextUtils.fileNameFromPath(g.path ?? ""),
    value: g.id,
  }));

  const onInputChange = debounce((input: string) => {
    setQuery(input);
  }, 500);

  const onChange = (selectedItems: ValueType<Option, boolean>) => {
    const selected = getSelectedItems(selectedItems);
    props.onSelect(
      selected.map((s) => ({
        id: s.value,
        title: s.label,
      }))
    );
  };

  const options = props.galleries.map((g) => ({
    value: g.id,
    label: g.title ?? "Unknown",
  }));

  return (
    <SelectComponent
      onChange={onChange}
      onInputChange={onInputChange}
      isLoading={loading}
      items={items}
      selectedOptions={options}
      isMulti
      placeholder="Search for gallery..."
      noOptionsMessage={query === "" ? null : "No galleries found."}
      showDropdown={false}
    />
  );
};

export const SceneSelect: React.FC<ISceneSelect> = (props) => {
  const [query, setQuery] = useState<string>("");
  const { data, loading } = GQL.useFindScenesQuery({
    skip: query === "",
    variables: {
      filter: {
        q: query,
      },
    },
  });

  const scenes = data?.findScenes.scenes ?? [];
  const items = scenes.map((s) => ({
    label: s.title ?? TextUtils.fileNameFromPath(s.path ?? ""),
    value: s.id,
  }));

  const onInputChange = debounce((input: string) => {
    setQuery(input);
  }, 500);

  const onChange = (selectedItems: ValueType<Option, true>) => {
    const selected = getSelectedItems(selectedItems);
    props.onSelect(
      (selected ?? []).map((s) => ({
        id: s.value,
        title: s.label,
      }))
    );
  };

  const options = props.scenes.map((s) => ({
    value: s.id,
    label: s.title,
  }));

  return (
    <SelectComponent
      onChange={onChange}
      onInputChange={onInputChange}
      isLoading={loading}
      items={items}
      selectedOptions={options}
      isMulti
      placeholder="Search for scene..."
      noOptionsMessage={query === "" ? null : "No scenes found."}
      showDropdown={false}
    />
  );
};

interface IMarkerSuggestProps {
  initialMarkerTitle?: string;
  onChange: (title: string) => void;
}
export const MarkerTitleSuggest: React.FC<IMarkerSuggestProps> = (props) => {
  const { data, loading } = useMarkerStrings();
  const suggestions = data?.markerStrings ?? [];

  const onChange = (selectedItem: ValueType<Option, false>) =>
    props.onChange(selectedItem?.value ?? "");

  const items = suggestions.map((item) => ({
    label: item?.title ?? "",
    value: item?.title ?? "",
  }));
  const initialIds = props.initialMarkerTitle ? [props.initialMarkerTitle] : [];

  // add initial value to items if still loading, to ensure existing value
  // is populated
  if (loading && initialIds.length > 0) {
    items.push({
      label: initialIds[0],
      value: initialIds[0],
    });
  }

  return (
    <SelectComponent
      isMulti={false}
      creatable
      onChange={onChange}
      isLoading={loading}
      items={items}
      initialIds={initialIds}
      placeholder="Marker title..."
      className="select-suggest"
      showDropdown={false}
      groupHeader="Previously used titles..."
    />
  );
};

export const PerformerSelect: React.FC<IFilterProps> = (props) => {
  const { data, loading } = useAllPerformersForFilter();
  const [createPerformer] = usePerformerCreate();

  const { configuration } = React.useContext(ConfigurationContext);
  const intl = useIntl();
  const defaultCreatable =
    !configuration?.interface.disableDropdownCreate.performer ?? true;

  const performers = data?.allPerformers ?? [];

  const onCreate = async (name: string) => {
    const result = await createPerformer({
      variables: { input: { name } },
    });
    return {
      item: result.data!.performerCreate!,
      message: "Created performer",
    };
  };

  return (
    <FilterSelectComponent
      {...props}
      isMulti={props.isMulti ?? false}
      creatable={props.creatable ?? defaultCreatable}
      onCreate={onCreate}
      type="performers"
      isLoading={loading}
      items={performers}
      placeholder={
        props.noSelectionString ??
        intl.formatMessage(
          { id: "actions.select_entity" },
          { entityType: intl.formatMessage({ id: "performer" }) }
        )
      }
    />
  );
};

export const StudioSelect: React.FC<
  IFilterProps & { excludeIds?: string[] }
> = (props) => {
  const [studioAliases, setStudioAliases] = useState<Record<string, string[]>>(
    {}
  );
  const [allAliases, setAllAliases] = useState<string[]>([]);
  const { data, loading } = useAllStudiosForFilter();
  const [createStudio] = useStudioCreate();
  const intl = useIntl();

  const { configuration } = React.useContext(ConfigurationContext);
  const defaultCreatable =
    !configuration?.interface.disableDropdownCreate.studio ?? true;

  const exclude = useMemo(() => props.excludeIds ?? [], [props.excludeIds]);
  const studios = useMemo(
    () =>
      (data?.allStudios ?? []).filter((studio) => !exclude.includes(studio.id)),
    [data?.allStudios, exclude]
  );

  useEffect(() => {
    // build the studio aliases map
    const newAliases: Record<string, string[]> = {};
    const newAll: string[] = [];
    studios.forEach((s) => {
      newAliases[s.id] = s.aliases;
      newAll.push(...s.aliases);
    });
    setStudioAliases(newAliases);
    setAllAliases(newAll);
  }, [studios]);

  const StudioOption: React.FC<OptionProps<Option, boolean>> = (
    optionProps
  ) => {
    const { inputValue } = optionProps.selectProps;

    let thisOptionProps = optionProps;
    if (
      inputValue &&
      !optionProps.label.toLowerCase().includes(inputValue.toLowerCase())
    ) {
      // must be alias
      const newLabel = `${optionProps.data.label} (alias)`;
      thisOptionProps = {
        ...optionProps,
        children: newLabel,
      };
    }

    return <reactSelectComponents.Option {...thisOptionProps} />;
  };

  const filterOption = (option: Option, rawInput: string): boolean => {
    if (!rawInput) {
      return true;
    }

    const input = rawInput.toLowerCase();
    const optionVal = option.label.toLowerCase();

    if (optionVal.includes(input)) {
      return true;
    }

    // search for studio aliases
    const aliases = studioAliases[option.value];
    // only match on alias if exact
    if (aliases && aliases.some((a) => a.toLowerCase() === input)) {
      return true;
    }

    return false;
  };

  const onCreate = async (name: string) => {
    const result = await createStudio({
      variables: {
        input: { name },
      },
    });
    return { item: result.data!.studioCreate!, message: "Created studio" };
  };

  const isValidNewOption = (
    inputValue: string,
    value: ValueType<Option, boolean>,
    options: OptionsType<Option> | GroupedOptionsType<Option>
  ) => {
    if (!inputValue) {
      return false;
    }

    if (
      (options as OptionsType<Option>).some((o: Option) => {
        return o.label.toLowerCase() === inputValue.toLowerCase();
      })
    ) {
      return false;
    }

    if (allAliases.some((a) => a.toLowerCase() === inputValue.toLowerCase())) {
      return false;
    }

    return true;
  };

  return (
    <FilterSelectComponent
      {...props}
      filterOption={filterOption}
      isValidNewOption={isValidNewOption}
      components={{ Option: StudioOption }}
      isMulti={props.isMulti ?? false}
      type="studios"
      isLoading={loading}
      items={studios}
      placeholder={
        props.noSelectionString ??
        intl.formatMessage(
          { id: "actions.select_entity" },
          { entityType: intl.formatMessage({ id: "studio" }) }
        )
      }
      creatable={props.creatable ?? defaultCreatable}
      onCreate={onCreate}
    />
  );
};

export const MovieSelect: React.FC<IFilterProps> = (props) => {
  const { data, loading } = useAllMoviesForFilter();
  const items = data?.allMovies ?? [];
  const intl = useIntl();

  return (
    <FilterSelectComponent
      {...props}
      isMulti={props.isMulti ?? false}
      type="movies"
      isLoading={loading}
      items={items}
      placeholder={
        props.noSelectionString ??
        intl.formatMessage(
          { id: "actions.select_entity" },
          { entityType: intl.formatMessage({ id: "movie" }) }
        )
      }
    />
  );
};

export const TagSelect: React.FC<IFilterProps & { excludeIds?: string[] }> = (
  props
) => {
  const [tagAliases, setTagAliases] = useState<Record<string, string[]>>({});
  const [allAliases, setAllAliases] = useState<string[]>([]);
  const { data, loading } = useAllTagsForFilter();
  const [createTag] = useTagCreate();
  const intl = useIntl();
  const placeholder =
    props.noSelectionString ??
    intl.formatMessage(
      { id: "actions.select_entity" },
      { entityType: intl.formatMessage({ id: "tags" }) }
    );

  const { configuration } = React.useContext(ConfigurationContext);
  const defaultCreatable =
    !configuration?.interface.disableDropdownCreate.tag ?? true;

  const exclude = useMemo(() => props.excludeIds ?? [], [props.excludeIds]);
  const tags = useMemo(
    () => (data?.allTags ?? []).filter((tag) => !exclude.includes(tag.id)),
    [data?.allTags, exclude]
  );

  useEffect(() => {
    // build the tag aliases map
    const newAliases: Record<string, string[]> = {};
    const newAll: string[] = [];
    tags.forEach((t) => {
      newAliases[t.id] = t.aliases;
      newAll.push(...t.aliases);
    });
    setTagAliases(newAliases);
    setAllAliases(newAll);
  }, [tags]);

  const TagOption: React.FC<OptionProps<Option, boolean>> = (optionProps) => {
    const { inputValue } = optionProps.selectProps;

    let thisOptionProps = optionProps;
    if (
      inputValue &&
      !optionProps.label.toLowerCase().includes(inputValue.toLowerCase())
    ) {
      // must be alias
      const newLabel = `${optionProps.data.label} (alias)`;
      thisOptionProps = {
        ...optionProps,
        children: newLabel,
      };
    }

    return <reactSelectComponents.Option {...thisOptionProps} />;
  };

  const filterOption = (option: Option, rawInput: string): boolean => {
    if (!rawInput) {
      return true;
    }

    const input = rawInput.toLowerCase();
    const optionVal = option.label.toLowerCase();

    if (optionVal.includes(input)) {
      return true;
    }

    // search for tag aliases
    const aliases = tagAliases[option.value];
    // only match on alias if exact
    if (aliases && aliases.some((a) => a.toLowerCase() === input)) {
      return true;
    }

    return false;
  };

  const onCreate = async (name: string) => {
    const result = await createTag({
      variables: {
        input: {
          name,
        },
      },
    });
    return { item: result.data!.tagCreate!, message: "Created tag" };
  };

  const isValidNewOption = (
    inputValue: string,
    value: ValueType<Option, boolean>,
    options: OptionsType<Option> | GroupedOptionsType<Option>
  ) => {
    if (!inputValue) {
      return false;
    }

    if (
      (options as OptionsType<Option>).some((o: Option) => {
        return o.label.toLowerCase() === inputValue.toLowerCase();
      })
    ) {
      return false;
    }

    if (allAliases.some((a) => a.toLowerCase() === inputValue.toLowerCase())) {
      return false;
    }

    return true;
  };

  return (
    <FilterSelectComponent
      {...props}
      filterOption={filterOption}
      isValidNewOption={isValidNewOption}
      components={{ Option: TagOption }}
      isMulti={props.isMulti ?? false}
      items={tags}
      creatable={props.creatable ?? defaultCreatable}
      type="tags"
      placeholder={placeholder}
      isLoading={loading}
      onCreate={onCreate}
      closeMenuOnSelect={!props.isMulti}
    />
  );
};

export const FilterSelect: React.FC<IFilterProps & ITypeProps> = (props) =>
  props.type === "performers" ? (
    <PerformerSelect {...props} creatable={false} />
  ) : props.type === "studios" || props.type === "parent_studios" ? (
    <StudioSelect {...props} creatable={false} />
  ) : props.type === "movies" ? (
    <MovieSelect {...props} creatable={false} />
  ) : (
    <TagSelect {...props} creatable={false} />
  );
