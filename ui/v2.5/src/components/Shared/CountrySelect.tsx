import React from "react";
import Creatable from "react-select/creatable";
import { useIntl } from "react-intl";
import { getCountries } from "src/utils";
import CountryLabel from "./CountryLabel";

interface IProps {
  value?: string | undefined;
  onChange?: (value: string) => void;
  disabled?: boolean;
  className?: string;
  showFlag?: boolean;
  isClearable?: boolean;
}

const CountrySelect: React.FC<IProps> = ({
  value,
  onChange,
  disabled = false,
  isClearable = true,
  showFlag,
  className,
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
    />
  );
};

export default CountrySelect;
