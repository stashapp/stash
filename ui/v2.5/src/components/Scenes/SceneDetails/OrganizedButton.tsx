import React from "react";
import cx from "classnames";
import { Button, Spinner } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { defineMessages, useIntl } from "react-intl";
import { faBox } from "@fortawesome/free-solid-svg-icons";
import { useIsAdmin } from "src/hooks/CurrentUser";

export interface IOrganizedButtonProps {
  loading: boolean;
  organized: boolean;
  onClick: () => void;
}

export const OrganizedButton: React.FC<IOrganizedButtonProps> = (
  props: IOrganizedButtonProps
) => {
  const intl = useIntl();
  const isAdmin = useIsAdmin();
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
      // organized may only be changed by administrators
      disabled={!isAdmin}
    >
      <Icon icon={faBox} />
    </Button>
  );
};
