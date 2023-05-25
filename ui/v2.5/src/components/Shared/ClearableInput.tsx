import React from "react";
import { Button, FormControl } from "react-bootstrap";
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import { useIntl } from "react-intl";
import { Icon } from "./Icon";
import useFocus from "src/utils/focus";

interface IClearableInput {
  value: string;
  setValue: (value: string) => void;
}

export const ClearableInput: React.FC<IClearableInput> = ({
  value,
  setValue,
}) => {
  const intl = useIntl();

  const [queryRef, setQueryFocus] = useFocus();
  const queryClearShowing = !!value;

  function onChangeQuery(event: React.FormEvent<HTMLInputElement>) {
    setValue(event.currentTarget.value);
  }

  function onClearQuery() {
    setValue("");
    setQueryFocus();
  }

  return (
    <div className="clearable-input-group">
      <FormControl
        ref={queryRef}
        placeholder={`${intl.formatMessage({ id: "actions.search" })}…`}
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
