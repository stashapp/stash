import React from "react";
import { Spinner } from "react-bootstrap";
import cx from "classnames";
import { useIntl } from "react-intl";
import { PatchComponent } from "src/patch";

interface ILoadingProps {
  message?: JSX.Element | string;
  inline?: boolean;
  small?: boolean;
  card?: boolean;
}

const CLASSNAME = "LoadingIndicator";
const CLASSNAME_MESSAGE = `${CLASSNAME}-message`;

export const LoadingIndicator: React.FC<ILoadingProps> = PatchComponent(
  "LoadingIndicator",
  ({ message, inline = false, small = false, card = false }) => {
    const intl = useIntl();

    const text = intl.formatMessage({ id: "loading.generic" });

    return (
      <div className={cx(CLASSNAME, { inline, small, "card-based": card })}>
        <Spinner
          animation="border"
          role="status"
          size={small ? "sm" : undefined}
        >
          <span className="sr-only">{text}</span>
        </Spinner>
        {message !== "" && (
          <h4 className={CLASSNAME_MESSAGE}>{message ?? text}</h4>
        )}
      </div>
    );
  }
);
