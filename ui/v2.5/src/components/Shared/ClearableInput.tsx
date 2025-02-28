import React from "react";
import { Button, FormControl } from "react-bootstrap";
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import { useIntl } from "react-intl";
import { Icon } from "./Icon";
import useFocus from "src/utils/focus";
import cx from "classnames";

interface IClearableInput {
  className?: string;
  value: string;
  setValue: (value: string) => void;
  focus?: ReturnType<typeof useFocus>;
  placeholder?: string;
}

export const ClearableInput: React.FC<IClearableInput> = ({
  className,
  value,
  setValue,
  focus,
  placeholder,
}) => {
  const intl = useIntl();

  const [defaultQueryRef, setQueryFocusDefault] = useFocus();
  const [queryRef, setQueryFocus] = focus || [
    defaultQueryRef,
    setQueryFocusDefault,
  ];
  const queryClearShowing = !!value;

  function onChangeQuery(event: React.FormEvent<HTMLInputElement>) {
    setValue(event.currentTarget.value);
  }

  function onClearQuery() {
    setValue("");
    setQueryFocus();
  }

  return (
    <div className={cx("clearable-input-group", className)}>
      <FormControl
        ref={queryRef}
        placeholder={placeholder}
        value={value}
        onInput={onChangeQuery}
        className="clearable-text-field"
      />
      {queryClearShowing && (
        <Button
          variant="secondary"
          onClick={onClearQuery}
          title={intl.formatMessage({ id: "actions.clear" })}
          className="clearable-text-field-clear"
        >
          <Icon icon={faTimes} />
        </Button>
      )}
    </div>
  );
};

export default ClearableInput;
