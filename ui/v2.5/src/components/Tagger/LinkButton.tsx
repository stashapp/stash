import React from "react";
import { useIntl } from "react-intl";
import { faLink } from "@fortawesome/free-solid-svg-icons";

import { OperationButton } from "../Shared/OperationButton";
import { Icon } from "../Shared/Icon";

export const LinkButton: React.FC<{
  disabled: boolean;
  onLink: () => Promise<void>;
}> = ({ disabled, onLink }) => {
  const intl = useIntl();

  return (
    <OperationButton
      variant="secondary"
      disabled={disabled}
      operation={onLink}
      hideChildrenWhenLoading
      title={intl.formatMessage({ id: "component_tagger.verb_link_existing" })}
    >
      <Icon icon={faLink} />
    </OperationButton>
  );
};
