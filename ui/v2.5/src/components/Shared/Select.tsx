import React, { useState, CSSProperties } from "react";
import Select, { ValueType } from "react-select";
import CreatableSelect from "react-select/creatable";
import { debounce } from "lodash";

import * as GQL from "src/core/generated-graphql";
import {
  useAllTagsForFilter,
  useAllMoviesForFilter,
  useAllStudiosForFilter,
  useAllPerformersForFilter,
  useMarkerStrings,
  useScrapePerformerList,
  useValidGalleriesForScene,
  useTagCreate,
} from "src/core/StashService";
import { useToast } from "src/hooks";

type ValidTypes =
  | GQL.SlimPerformerDataFragment
  | GQL.Tag
  | GQL.SlimStudioDataFragment;
type Option = { value: string; label: string };

interface ITypeProps {
  type?:
    | "performers"
    | "studios"
    | "parent_studios"
    | "tags"
    | "sceneTags"
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
}
interface ISelectProps {
  className?: string;
  items: Option[];
  selectedOptions?: Option[];
  creatable?: boolean;
  onCreateOption?: (value: string) => void;
  isLoading: boolean;
  isDisabled?: boolean;
  onChange: (item: ValueType<Option>) => void;
  initialIds?: string[];
  isMulti?: boolean;
  isClearable?: boolean;
  onInputChange?: (input: string) => void;
  placeholder?: string;
  showDropdown?: boolean;
  groupHeader?: string;
  closeMenuOnSelect?: boolean;
}

interface ISceneGallerySelect {
  initialId?: string;
  sceneId: string;
  onSelect: (
    item:
      | GQL.ValidGalleriesForSceneQuery["validGalleriesForScene"][0]
      | undefined
  ) => void;
}

const getSelectedValues = (selectedItems: ValueType<Option>) =>
  selectedItems
    ? (Array.isArray(selectedItems) ? selectedItems : [selectedItems]).map(
        (item) => item.value
      )
    : [];

export const SceneGallerySelect: React.FC<ISceneGallerySelect> = (props) => {
  const { data, loading } = useValidGalleriesForScene(props.sceneId);
  const galleries = data?.validGalleriesForScene ?? [];
  const items = (galleries.length > 0
    ? [{ path: "None", id: "0" }, ...galleries]
    : []
  ).map((g) => ({ label: g.path, value: g.id }));

  const onChange = (selectedItems: ValueType<Option>) => {
    const selectedItem = getSelectedValues(selectedItems)[0];
    props.onSelect(
      selectedItem ? galleries.find((g) => g.id === selectedItem) : undefined
    );
  };

  const selectedOptions: Option[] = props.initialId
    ? items.filter((item) => props.initialId?.indexOf(item.value) !== -1)
    : [];

  return (
    <SelectComponent
      className="input-control"
      onChange={onChange}
      isLoading={loading}
      items={items}
      selectedOptions={selectedOptions}
    />
  );
};

interface IScrapePerformerSuggestProps {
  scraperId: string;
  onSelectPerformer: (performer: GQL.ScrapedPerformerDataFragment) => void;
  placeholder?: string;
}
export const ScrapePerformerSuggest: React.FC<IScrapePerformerSuggestProps> = (
  props
) => {
  const [query, setQuery] = React.useState<string>("");
  const { data, loading } = useScrapePerformerList(props.scraperId, query);

  const performers = data?.scrapePerformerList ?? [];
  const items = performers.map((item) => ({
    label: item.name ?? "",
    value: item.name ?? "",
  }));

  const onInputChange = debounce((input: string) => {
    setQuery(input);
  }, 500);

  const onChange = (selectedItems: ValueType<Option>) => {
    const name = getSelectedValues(selectedItems)[0];
    const performer = performers.find((p) => p.name === name);
    if (performer) props.onSelectPerformer(performer);
  };

  return (
    <SelectComponent
      onChange={onChange}
      onInputChange={onInputChange}
      isLoading={loading}
      items={items}
      initialIds={[]}
      placeholder={props.placeholder}
      className="select-suggest"
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

  const onChange = (selectedItems: ValueType<Option>) =>
    props.onChange(getSelectedValues(selectedItems)[0]);

  const items = suggestions.map((item) => ({
    label: item?.title ?? "",
    value: item?.title ?? "",
  }));
  const initialIds = props.initialMarkerTitle ? [props.initialMarkerTitle] : [];
  return (
    <SelectComponent
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
export const FilterSelect: React.FC<IFilterProps & ITypeProps> = (props) =>
  props.type === "performers" ? (
    <PerformerSelect {...(props as IFilterProps)} />
  ) : props.type === "studios" || props.type === "parent_studios" ? (
    <StudioSelect {...(props as IFilterProps)} />
  ) : props.type === "movies" ? (
    <MovieSelect {...(props as IFilterProps)} />
  ) : (
    <TagSelect {...(props as IFilterProps)} />
  );

export const PerformerSelect: React.FC<IFilterProps> = (props) => {
  const { data, loading } = useAllPerformersForFilter();

  const normalizedData = data?.allPerformersSlim ?? [];
  const items: Option[] = normalizedData.map((item) => ({
    value: item.id,
    label: item.name ?? "",
  }));
  const placeholder = props.noSelectionString ?? "Select performer...";
  const selectedOptions: Option[] = props.ids
    ? items.filter((item) => props.ids?.indexOf(item.value) !== -1)
    : [];

  const onChange = (selectedItems: ValueType<Option>) => {
    const selectedIds = getSelectedValues(selectedItems);
    props.onSelect?.(
      normalizedData.filter((item) => selectedIds.indexOf(item.id) !== -1)
    );
  };

  return (
    <SelectComponent
      {...props}
      selectedOptions={selectedOptions}
      onChange={onChange}
      type="performers"
      isLoading={loading}
      items={items}
      placeholder={placeholder}
    />
  );
};

export const StudioSelect: React.FC<IFilterProps> = (props) => {
  const { data, loading } = useAllStudiosForFilter();

  const normalizedData = data?.allStudiosSlim ?? [];

  const items = (normalizedData.length > 0
    ? [{ name: "None", id: "0" }, ...normalizedData]
    : []
  ).map((item) => ({
    value: item.id,
    label: item.name,
  }));

  const placeholder = props.noSelectionString ?? "Select studio...";
  const selectedOptions: Option[] = props.ids
    ? items.filter((item) => props.ids?.indexOf(item.value) !== -1)
    : [];

  const onChange = (selectedItems: ValueType<Option>) => {
    const selectedIds = getSelectedValues(selectedItems);
    props.onSelect?.(
      normalizedData.filter((item) => selectedIds.indexOf(item.id) !== -1)
    );
  };

  return (
    <SelectComponent
      {...props}
      onChange={onChange}
      type="studios"
      isLoading={loading}
      items={items}
      placeholder={placeholder}
      selectedOptions={selectedOptions}
    />
  );
};

export const MovieSelect: React.FC<IFilterProps> = (props) => {
  const { data, loading } = useAllMoviesForFilter();

  const normalizedData = data?.allMoviesSlim ?? [];

  const items = (normalizedData.length > 0
    ? [{ name: "None", id: "0" }, ...normalizedData]
    : []
  ).map((item) => ({
    value: item.id,
    label: item.name,
  }));

  const placeholder = props.noSelectionString ?? "Select movie...";
  const selectedOptions: Option[] = props.ids
    ? items.filter((item) => props.ids?.indexOf(item.value) !== -1)
    : [];

  const onChange = (selectedItems: ValueType<Option>) => {
    const selectedIds = getSelectedValues(selectedItems);
    props.onSelect?.(
      normalizedData.filter((item) => selectedIds.indexOf(item.id) !== -1)
    );
  };

  return (
    <SelectComponent
      {...props}
      onChange={onChange}
      type="studios"
      isLoading={loading}
      items={items}
      placeholder={placeholder}
      selectedOptions={selectedOptions}
    />
  );
};

export const TagSelect: React.FC<IFilterProps> = (props) => {
  const [loading, setLoading] = useState(false);
  const [selectedIds, setSelectedIds] = useState<string[]>(props.ids ?? []);
  const { data, loading: dataLoading } = useAllTagsForFilter();
  const [createTag] = useTagCreate({ name: "" });
  const Toast = useToast();
  const placeholder = props.noSelectionString ?? "Select tags...";

  const selectedTags = props.ids ?? selectedIds;

  const tags = data?.allTagsSlim ?? [];
  const selected = tags
    .filter((tag) => selectedTags.indexOf(tag.id) !== -1)
    .map((tag) => ({ value: tag.id, label: tag.name }));
  const items: Option[] = tags.map((item) => ({
    value: item.id,
    label: item.name,
  }));

  const onCreate = async (tagName: string) => {
    try {
      setLoading(true);
      const result = await createTag({
        variables: { name: tagName },
      });

      if (result?.data?.tagCreate) {
        setSelectedIds([...selectedIds, result.data.tagCreate.id]);
        props.onSelect?.([
          ...tags.filter((item) => selectedIds.indexOf(item.id) !== -1),
          result.data.tagCreate,
        ]);
        setLoading(false);

        Toast.success({
          content: (
            <span>
              Created tag: <b>{tagName}</b>
            </span>
          ),
        });
      }
    } catch (e) {
      Toast.error(e);
    }
  };

  const onChange = (selectedItems: ValueType<Option>) => {
    const selectedValues = getSelectedValues(selectedItems);
    setSelectedIds(selectedValues);
    props.onSelect?.(
      tags.filter((item) => selectedValues.indexOf(item.id) !== -1)
    );
  };

  return (
    <SelectComponent
      {...props}
      onChange={onChange}
      creatable
      type="tags"
      placeholder={placeholder}
      isLoading={loading || dataLoading}
      items={items}
      onCreateOption={onCreate}
      selectedOptions={selected}
      closeMenuOnSelect={false}
    />
  );
};

const SelectComponent: React.FC<ISelectProps & ITypeProps> = ({
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
  isMulti = false,
  onInputChange,
  placeholder,
  showDropdown = true,
  groupHeader,
  closeMenuOnSelect = true,
}) => {
  const defaultValue =
    items.filter((item) => initialIds?.indexOf(item.value) !== -1) ?? null;

  const options = groupHeader
    ? [
        {
          label: groupHeader,
          options: items,
        },
      ]
    : items;

  const styles = {
    option: (base: CSSProperties) => ({
      ...base,
      color: "#000",
    }),
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    container: (base: CSSProperties, state: any) => ({
      ...base,
      zIndex: state.isFocused ? 10 : base.zIndex,
    }),
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    multiValueRemove: (base: CSSProperties, state: any) => ({
      ...base,
      color: state.isFocused ? base.color : "#333333",
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
    defaultValue,
    noOptionsMessage: () => (type !== "tags" ? "None" : null),
    placeholder: isDisabled ? "" : placeholder,
    onInputChange,
    isDisabled,
    isLoading,
    styles,
    closeMenuOnSelect,
    components: {
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
