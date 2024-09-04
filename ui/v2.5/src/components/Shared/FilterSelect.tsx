import React, { useMemo, useState } from "react";
import {
  OnChangeValue,
  StylesConfig,
  GroupBase,
  OptionsOrGroups,
  Options,
} from "react-select";
import AsyncSelect from "react-select/async";
import AsyncCreatableSelect, {
  AsyncCreatableProps,
} from "react-select/async-creatable";
import cx from "classnames";

import { useToast } from "src/hooks/Toast";
import { useDebounce } from "src/hooks/debounce";
import { IHasID } from "src/utils/data";

export type Option<T> = { value: string; object: T };

interface ISelectProps<T, IsMulti extends boolean>
  extends AsyncCreatableProps<Option<T>, IsMulti, GroupBase<Option<T>>> {
  selectedOptions?: OnChangeValue<Option<T>, IsMulti>;
  creatable?: boolean;
  isLoading?: boolean;
  isDisabled?: boolean;
  placeholder?: string;
  showDropdown?: boolean;
  groupHeader?: string;
  noOptionsMessageText?: string | null;
}

interface IFilterSelectProps<T, IsMulti extends boolean>
  extends Pick<
    ISelectProps<T, IsMulti>,
    | "selectedOptions"
    | "isLoading"
    | "isMulti"
    | "components"
    | "placeholder"
    | "closeMenuOnSelect"
  > {}

const getSelectedItems = <T,>(
  selectedItems: OnChangeValue<Option<T>, boolean>
) => {
  if (Array.isArray(selectedItems)) {
    return selectedItems;
  } else if (selectedItems) {
    return [selectedItems];
  } else {
    return [];
  }
};

const SelectComponent = <T, IsMulti extends boolean>(
  props: ISelectProps<T, IsMulti>
) => {
  const {
    selectedOptions,
    isLoading,
    isDisabled = false,
    creatable = false,
    components,
    placeholder,
    showDropdown = true,
    noOptionsMessageText: noOptionsMessage = "None",
  } = props;

  const styles: StylesConfig<Option<T>, IsMulti> = {
    option: (base) => ({
      ...base,
      color: "#000",
    }),
    container: (base, state) => ({
      ...base,
      zIndex: state.isFocused ? 10 : base.zIndex,
    }),
    multiValueRemove: (base, state) => ({
      ...base,
      color: state.isFocused ? base.color : "#333333",
    }),
  };

  const componentProps = {
    ...props,
    styles,
    defaultOptions: true,
    isClearable: true,
    value: selectedOptions ?? null,
    className: cx("react-select", props.className),
    classNamePrefix: "react-select",
    noOptionsMessage: () => noOptionsMessage,
    placeholder: isDisabled ? "" : placeholder,
    components: {
      ...components,
      IndicatorSeparator: () => null,
      ...((!showDropdown || isDisabled) && { DropdownIndicator: () => null }),
      ...(isDisabled && { MultiValueRemove: () => null }),
    },
  };

  return creatable ? (
    <AsyncCreatableSelect
      {...componentProps}
      isDisabled={isLoading || isDisabled}
    />
  ) : (
    <AsyncSelect {...componentProps} />
  );
};

export interface IFilterValueProps<T> {
  values?: T[];
  onSelect?: (item: T[]) => void;
}

export interface IFilterProps {
  noSelectionString?: string;
  className?: string;
  active?: boolean;
  isMulti?: boolean;
  isClearable?: boolean;
  isDisabled?: boolean;
  creatable?: boolean;
  menuPortalTarget?: HTMLElement | null;
}

export interface IFilterComponentProps<T> extends IFilterProps {
  loadOptions: (inputValue: string) => Promise<Option<T>[]>;
  onCreate?: (
    name: string
  ) => Promise<{ value: string; item: T; message: string }>;
  getNamedObject?: (id: string, name: string) => T;
  isValidNewOption?: (inputValue: string, options: T[]) => boolean;
}

export const FilterSelectComponent = <
  T extends IHasID,
  IsMulti extends boolean
>(
  props: IFilterValueProps<T> &
    IFilterComponentProps<T> &
    IFilterSelectProps<T, IsMulti>
) => {
  const {
    values,
    isMulti,
    onSelect,
    creatable = false,
    isValidNewOption,
    getNamedObject,
    loadOptions,
  } = props;
  const [loading, setLoading] = useState(false);
  const Toast = useToast();

  const selectedOptions = useMemo(() => {
    if (isMulti && values) {
      return values.map(
        (value) =>
          ({
            object: value,
            value: value.id,
          } as Option<T>)
      ) as unknown as OnChangeValue<Option<T>, IsMulti>;
    }

    if (values?.length) {
      return {
        object: values[0],
        value: values[0].id,
      } as OnChangeValue<Option<T>, IsMulti>;
    }
  }, [values, isMulti]);

  const onChange = (selectedItems: OnChangeValue<Option<T>, boolean>) => {
    const selected = getSelectedItems(selectedItems);

    onSelect?.(selected.map((item) => item.object));
  };

  const onCreate =
    creatable && props.onCreate
      ? async (name: string) => {
          try {
            setLoading(true);
            const {
              value,
              item: newItem,
              message,
            } = await props.onCreate!(name);
            const newItemOption = {
              object: newItem,
              value,
            } as Option<T>;
            if (!isMulti) {
              onChange(newItemOption);
            } else {
              const o = (selectedOptions ?? []) as Option<T>[];
              onChange([...o, newItemOption]);
            }

            setLoading(false);
            Toast.success(
              <span>
                {message}: <b>{name}</b>
              </span>
            );
          } catch (e) {
            Toast.error(e);
          }
        }
      : undefined;

  const getNewOptionData =
    creatable && getNamedObject
      ? (inputValue: string, optionLabel: React.ReactNode) => {
          return {
            value: "",
            object: getNamedObject("", optionLabel as string),
          };
        }
      : undefined;

  const validNewOption =
    creatable && isValidNewOption
      ? (
          inputValue: string,
          value: Options<Option<T>>,
          options: OptionsOrGroups<Option<T>, GroupBase<Option<T>>>
        ) => {
          return isValidNewOption(
            inputValue,
            (options as Options<Option<T>>).map((o) => o.object)
          );
        }
      : undefined;

  const debounceDelay = 100;
  const debounceLoadOptions = useDebounce((inputValue, callback) => {
    loadOptions(inputValue).then(callback);
  }, debounceDelay);

  return (
    <SelectComponent<T, IsMulti>
      {...props}
      loadOptions={debounceLoadOptions}
      isLoading={props.isLoading || loading}
      onChange={onChange}
      selectedOptions={selectedOptions}
      onCreateOption={onCreate}
      getNewOptionData={getNewOptionData}
      isValidNewOption={validNewOption}
    />
  );
};

export interface IFilterIDProps<T> {
  ids?: string[];
  onSelect?: (item: T[]) => void;
}
