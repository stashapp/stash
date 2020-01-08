import React, { useState } from "react";
import Select, { ValueType } from 'react-select';
import CreatableSelect from 'react-select/creatable';

import { ErrorUtils } from "../../utils/errors";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import useToast from '../Shared/Toast';

type ValidTypes =
  GQL.AllPerformersForFilterAllPerformers |
  GQL.AllTagsForFilterAllTags |
  GQL.AllStudiosForFilterAllStudios;
type Option = { value:string, label:string };

interface ITypeProps {
  type: 'performers' | 'studios' | 'tags';
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
  noSelectionString?: string;
  isMulti?: boolean;
}

export const FilterSelect: React.FC<IFilterProps & ITypeProps> = (props) => (
    props.type === 'performers' ? <PerformerSelect {...props as IFilterProps} />
    : props.type === 'studios' ? <StudioSelect {...props as IFilterProps} />
    : <TagSelect {...props as IFilterProps} />
);

export const PerformerSelect: React.FC<IFilterProps> = (props) => {
  const { data, loading } = StashService.useAllPerformersForFilter();
  
  const normalizedData = data?.allPerformers ?? [];
  const items:Option[] = normalizedData.map(item => ({ value: item.id, label: item.name ?? '' }));

  const onChange = (selectedItems:ValueType<Option>) => {
    const selectedIds = (Array.isArray(selectedItems) ? selectedItems : [selectedItems])
      .map(item => item.value);
    props.onSelect(normalizedData.filter(item => selectedIds.indexOf(item.id) !== -1));
  };

  return <SelectComponent {...props} onChange={onChange} type="performers" isLoading={loading} items={items} />
}

export const StudioSelect: React.FC<IFilterProps> = (props) => {
  const { data, loading } = StashService.useAllStudiosForFilter();

  const normalizedData = data?.allStudios ?? [];
  const items:Option[] = normalizedData.map(item => ({ value: item.id, label: item.name }));

  const onChange = (selectedItems:ValueType<Option>) => {
    const selectedIds = (Array.isArray(selectedItems) ? selectedItems : [selectedItems])
      .map(item => item.value);
    props.onSelect(normalizedData.filter(item => selectedIds.indexOf(item.id) !== -1));
  };

  return <SelectComponent {...props} onChange={onChange} type="studios" isLoading={loading} items={items} />
}

export const TagSelect: React.FC<IFilterProps> = (props) => {
  const [loading, setLoading] = useState(false);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const { data, loading: dataLoading } = StashService.useAllTagsForFilter();
  const createTag = StashService.useTagCreate({name: ''});
  const Toast = useToast();

  const tags = data?.allTags ?? [];

  const onCreate = async (tagName: string) => {
    try {
      setLoading(true);
      const result = await createTag({
        variables: { name: tagName },
      });

      setSelectedIds([...selectedIds, result.data.tagCreate.id]);
      props.onSelect([...tags, result.data.tagCreate].filter(item => selected.indexOf(item.id) !== -1));
      setLoading(false);

      Toast({ content: (<span>Created tag: <b>{tagName}</b></span>) });
    } catch (e) {
      ErrorUtils.handle(e);
    }
  };

  const onChange = (selectedItems:ValueType<Option>) => {
    debugger;
    const selected = (Array.isArray(selectedItems) ? selectedItems : [selectedItems])
      .map(item => item.value);
    setSelectedIds(selected);
    props.onSelect(tags.filter(item => selected.indexOf(item.id) !== -1));
  };

  const selected = tags.filter(tag => selectedIds.indexOf(tag.id) !== -1).map(tag => ({value: tag.id, label: tag.name}));
  const items:Option[] = tags.map(item => ({ value: item.id, label: item.name }));
  return <SelectComponent {...props} onChange={onChange} creatable={true} type="tags" 
    isLoading={loading || dataLoading} items={items} onCreateOption={onCreate} selectedOptions={selected}  />
}

const SelectComponent: React.FC<ISelectProps & ITypeProps> = ({
    type,
    initialIds,
    noSelectionString,
    onChange,
    className,
    items,
    selectedOptions,
    isLoading,
    onCreateOption,
    creatable = false,
    isMulti = false,
}) => {
  const defaultValue = items.filter(item => initialIds?.indexOf(item.value) !== -1) ?? null;

  const props = {
    className: className,
    options: items,
    value: selectedOptions,
    onChange: onChange,
    isMulti: isMulti,
    defaultValue: defaultValue,
    noOptionsMessage: () => (type !== 'tags' ? 'None' : null),
    placeholder: noSelectionString ?? "(No selection)"
  }
  
  return (
    creatable
      ? <CreatableSelect {...props} isLoading={isLoading} isDisabled={isLoading} onCreateOption={onCreateOption} />
      : <Select {...props} isLoading={isLoading} />
  );
};
