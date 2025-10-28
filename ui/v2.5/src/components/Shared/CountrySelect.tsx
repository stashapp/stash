import React from "react";
import Creatable from "react-select/creatable";
import { useIntl } from "react-intl";
import { getCountries } from "src/utils/country";
import { CountryLabel } from "./CountryLabel";
import { PatchComponent } from "src/patch";

interface IProps {
  value?: string;
  onChange?: (value: string) => void;
  disabled?: boolean;
  className?: string;
  showFlag?: boolean;
  isClearable?: boolean;
  menuPortalTarget?: HTMLElement | null;
}

const _CountrySelect: React.FC<IProps> = ({
  value,
  onChange,
  disabled = false,
  isClearable = true,
  showFlag,
  className,
  menuPortalTarget,
}) => {
  const { locale } = useIntl();
  const options = getCountries(locale);
  const selected = options.find((opt) => opt.value === value) ?? {
    label: value,
    value,
  };

  return (
    <Creatable
      classNamePrefix="react-select"
      value={selected}
      isClearable={isClearable}
      formatOptionLabel={(option) => (
        <CountryLabel country={option.value} showFlag={showFlag} />
      )}
      placeholder="Country"
      options={options}
      onChange={(selectedOption) => onChange?.(selectedOption?.value ?? "")}
      isDisabled={disabled || !onChange}
      components={{
        IndicatorSeparator: null,
      }}
      className={`CountrySelect ${className}`}
      menuPortalTarget={menuPortalTarget}
    />
  );
};

export const CountrySelect = PatchComponent("CountrySelect", _CountrySelect);
