import React from "react";
import cx from "classnames";
import { Button, Spinner } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { defineMessages, useIntl } from "react-intl";
import { faBox } from "@fortawesome/free-solid-svg-icons";

export interface IOrganizedButtonProps {
  loading: boolean;
  organized: boolean;
  onClick: () => void;
}

export const OrganizedButton: React.FC<IOrganizedButtonProps> = (
  props: IOrganizedButtonProps
) => {
  const intl = useIntl();
  const messages = defineMessages({
    organized: {
      id: "organized",
      defaultMessage: "Organized",
    },
  });

  if (props.loading) return <Spinner animation="border" role="status" />;

  return (
    <Button
      variant="secondary"
      title={intl.formatMessage(messages.organized)}
      className={cx(
        "minimal",
        "organized-button",
        props.organized ? "organized" : "not-organized"
      )}
      onClick={props.onClick}
    >
      <Icon icon={faBox} />
    </Button>
  );
};
