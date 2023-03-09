import React, { useState, useRef, useEffect } from "react";
import { Button, ButtonProps } from "react-bootstrap";
import { LoadingIndicator } from "./LoadingIndicator";

interface IOperationButton extends ButtonProps {
  operation?: () => Promise<void>;
  loading?: boolean;
  hideChildrenWhenLoading?: boolean;
  setLoading?: (v: boolean) => void;
}

export const OperationButton: React.FC<IOperationButton> = (props) => {
  const [internalLoading, setInternalLoading] = useState(false);
  const mounted = useRef(false);

  const {
    operation,
    loading: externalLoading,
    hideChildrenWhenLoading = false,
    setLoading: setExternalLoading,
    ...withoutExtras
  } = props;

  useEffect(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
    };
  }, []);

  const setLoading = setExternalLoading || setInternalLoading;
  const loading =
    externalLoading !== undefined ? externalLoading : internalLoading;

  async function handleClick() {
    if (operation && !loading) {
      setLoading(true);
      await operation();

      if (mounted.current) {
        setLoading(false);
      }
    }
  }

  return (
    <Button onClick={handleClick} {...withoutExtras}>
      {loading && (
        <span className="mr-2">
          <LoadingIndicator message="" inline small />
        </span>
      )}
      {(!loading || !hideChildrenWhenLoading) && props.children}
    </Button>
  );
};
