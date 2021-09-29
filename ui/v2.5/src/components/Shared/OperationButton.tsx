import React, { useState } from "react";
import { Button, ButtonProps } from "react-bootstrap";
import { LoadingIndicator } from "src/components/Shared";

interface IOperationButton extends ButtonProps {
  operation?: () => Promise<void>;
  loading?: boolean;
  setLoading?: (v: boolean) => void;
}

export const OperationButton: React.FC<IOperationButton> = (props) => {
  const [internalLoading, setInternalLoading] = useState(false);

  const {
    operation,
    loading: externalLoading,
    setLoading: setExternalLoading,
    ...withoutExtras
  } = props;

  const setLoading = setExternalLoading || setInternalLoading;
  const loading =
    externalLoading !== undefined ? externalLoading : internalLoading;

  async function handleClick() {
    if (operation) {
      setLoading(true);
      await operation();
      setLoading(false);
    }
  }

  return (
    <Button onClick={handleClick} {...withoutExtras}>
      {loading && (
        <span className="mr-2">
          <LoadingIndicator message="" inline small />
        </span>
      )}
      {props.children}
    </Button>
  );
};
