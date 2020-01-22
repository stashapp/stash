import React, { useState, useCallback } from "react";
import Select, { ValueType } from "react-select";
import CreatableSelect from "react-select/creatable";
import { debounce } from "lodash";

import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { useToast } from "src/hooks";

type ValidTypes =
  | GQL.SlimPerformerDataFragment
  | GQL.Tag
  | GQL.SlimStudioDataFragment;
type Option = { value: string; label: string };

interface ITypeProps {
  type?: "performers" | "studios" | "tags";
}
interface IFilterProps {
  initialIds: string[];
  onSelect: (item: ValidTypes[]) => void;
  noSelectionString?: string;
  className?: string;
  isMulti?: boolean;
}
interface ISelectProps {
  className?: string;
  items: Option[];
  selectedOptions?: Option[];
  creatable?: boolean;
  onCreateOption?: (value: string) => void;
  isLoading: boolean;
  onChange: (item: ValueType<Option>) => void;
  initialIds: string[];
  isMulti?: boolean;
  onInputChange?: (input: string) => void;
  placeholder?: string;
}

interface ISceneGallerySelect {
  initialId?: string;
  sceneId: string;
  onSelect: (
    item: GQL.ValidGalleriesForSceneQuery["validGalleriesForScene"][0] | undefined
  ) => void;
}

const getSelectedValues = (selectedItems: ValueType<Option>) =>
  (Array.isArray(selectedItems) ? selectedItems : [selectedItems]).map(
    item => item.value
  );

export const SceneGallerySelect: React.FC<ISceneGallerySelect> = props => {
  const { data, loading } = StashService.useValidGalleriesForScene(
    props.sceneId
  );
  const galleries = data?.validGalleriesForScene ?? [];
  const items = (galleries.length > 0
    ? [{ path: "None", id: "0" }, ...galleries]
    : []
  ).map(g => ({ label: g.path, value: g.id }));

  const onChange = (selectedItems: ValueType<Option>) => {
    const selectedItem = getSelectedValues(selectedItems)[0];
    props.onSelect(galleries.find(g => g.id === selectedItem.value));
  };

  const initialId = props.initialId ? [props.initialId] : [];
  return (
    <SelectComponent
      onChange={onChange}
      isLoading={loading}
      items={items}
      initialIds={initialId}
    />
  );
};

interface IScrapePerformerSuggestProps {
  scraperId: string;
  onSelectPerformer: (
    performer: GQL.ScrapedPerformerDataFragment
  ) => void;
  placeholder?: string;
}
export const ScrapePerformerSuggest: React.FC<IScrapePerformerSuggestProps> = props => {
  const [query, setQuery] = React.useState<string>("");
  const { data, loading } = StashService.useScrapePerformerList(
    props.scraperId,
    query
  );

  const performers = data?.scrapePerformerList ?? [];
  const items = performers.map(item => ({
    label: item.name ?? "",
    value: item.name ?? ""
  }));

  const onInputChange = useCallback(
    debounce((input: string) => {
      setQuery(input);
    }, 500),
    []
  );
  const onChange = (selectedItems: ValueType<Option>) => {
    const name = getSelectedValues(selectedItems)[0];
    const performer = performers.find(p => p.name === name);
    if(performer)
      props.onSelectPerformer(performer)
  };

  return (
    <SelectComponent
      onChange={onChange}
      onInputChange={onInputChange}
      isLoading={loading}
      items={items}
      initialIds={[]}
      placeholder={props.placeholder}
    />
  );
};

interface IMarkerSuggestProps {
  initialMarkerTitle?: string;
  onChange: (title: string) => void;
}
export const MarkerTitleSuggest: React.FC<IMarkerSuggestProps> = props => {
  const { data, loading } = StashService.useMarkerStrings();
  const suggestions = data?.markerStrings ?? [];

  const onChange = (selectedItems: ValueType<Option>) =>
    props.onChange(getSelectedValues(selectedItems)[0]);

  const items = suggestions.map(item => ({
    label: item?.title ?? "",
    value: item?.title ?? ""
  }));
  const initialIds = props.initialMarkerTitle ? [props.initialMarkerTitle] : [];
  return (
    <SelectComponent
      creatable
      onChange={onChange}
      isLoading={loading}
      items={items}
      initialIds={initialIds}
    />
  );
};

export const FilterSelect: React.FC<IFilterProps & ITypeProps> = props =>
  props.type === "performers" ? (
    <PerformerSelect {...(props as IFilterProps)} />
  ) : props.type === "studios" ? (
    <StudioSelect {...(props as IFilterProps)} />
  ) : (
    <TagSelect {...(props as IFilterProps)} />
  );

export const PerformerSelect: React.FC<IFilterProps> = props => {
  const { data, loading } = StashService.useAllPerformersForFilter();

  const normalizedData = data?.allPerformers ?? [];
  const items: Option[] = normalizedData.map(item => ({
    value: item.id,
    label: item.name ?? ""
  }));
  const placeholder = props.noSelectionString ?? "Select performer...";

  const onChange = (selectedItems: ValueType<Option>) => {
    const selectedIds = getSelectedValues(selectedItems);
    props.onSelect(
      normalizedData.filter(item => selectedIds.indexOf(item.id) !== -1)
    );
  };

  return (
    <SelectComponent
      {...props}
      onChange={onChange}
      type="performers"
      isLoading={loading}
      items={items}
      placeholder={placeholder}
    />
  );
};

export const StudioSelect: React.FC<IFilterProps> = props => {
  const { data, loading } = StashService.useAllStudiosForFilter();

  const normalizedData = data?.allStudios ?? [];
  const items: Option[] = normalizedData.map(item => ({
    value: item.id,
    label: item.name
  }));
  const placeholder = props.noSelectionString ?? "Select studio...";

  const onChange = (selectedItems: ValueType<Option>) => {
    const selectedIds = getSelectedValues(selectedItems);
    props.onSelect(
      normalizedData.filter(item => selectedIds.indexOf(item.id) !== -1)
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
    />
  );
};

export const TagSelect: React.FC<IFilterProps> = props => {
  const [loading, setLoading] = useState(false);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const { data, loading: dataLoading } = StashService.useAllTagsForFilter();
  const [createTag] = StashService.useTagCreate({ name: "" });
  const Toast = useToast();
  const placeholder = props.noSelectionString ?? "Select tags...";

  const tags = data?.allTags ?? [];
  const selected = tags
    .filter(tag => selectedIds.indexOf(tag.id) !== -1)
    .map(tag => ({ value: tag.id, label: tag.name }));
  const items: Option[] = tags.map(item => ({
    value: item.id,
    label: item.name
  }));

  const onCreate = async (tagName: string) => {
    try {
      setLoading(true);
      const result = await createTag({
        variables: { name: tagName }
      });

      if(result?.data?.tagCreate) {
        setSelectedIds([...selectedIds, result.data.tagCreate.id]);
        props.onSelect(
          [...tags, result.data.tagCreate].filter(
            item => selectedIds.indexOf(item.id) !== -1
          )
        );
        setLoading(false);

        Toast.success({
          content: (
            <span>
              Created tag: <b>{tagName}</b>
            </span>
          )
        });
      }
    } catch (e) {
      Toast.error(e);
    }
  };

  const onChange = (selectedItems: ValueType<Option>) => {
    const selectedValues = getSelectedValues(selectedItems);
    setSelectedIds(selectedValues);
    props.onSelect(tags.filter(item => selectedValues.indexOf(item.id) !== -1));
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
  onCreateOption,
  creatable = false,
  isMulti = false,
  onInputChange,
  placeholder
}) => {
  const defaultValue =
    items.filter(item => initialIds?.indexOf(item.value) !== -1) ?? null;

  const props = {
    className,
    options: items,
    value: selectedOptions,
    onChange,
    isMulti,
    defaultValue,
    noOptionsMessage: () => (type !== "tags" ? "None" : null),
    placeholder,
    onInputChange
  };

  return creatable ? (
    <CreatableSelect
      {...props}
      isLoading={isLoading}
      isDisabled={isLoading}
      onCreateOption={onCreateOption}
    />
  ) : (
    <Select {...props} isLoading={isLoading} />
  );
};
