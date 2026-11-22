import React from "react";
import { Button } from "react-bootstrap";
import { useIntl } from "react-intl";
import { faLock, faLockOpen } from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";
import { Icon } from "./Icon";

export const LockFieldButton: React.FC<{
  locked: boolean;
  fieldName: string;
  onToggle: () => void;
}> = ({ locked, fieldName, onToggle }) => {
  const intl = useIntl();
  const title = intl.formatMessage(
    { id: locked ? "actions.unlock_field" : "actions.lock_field" },
    { field: fieldName }
  );

  return (
    <Button
      variant="link"
      className={cx("minimal", "lock-field-button", { locked })}
      title={title}
      aria-label={title}
      aria-pressed={locked}
      onClick={onToggle}
    >
      <Icon icon={locked ? faLock : faLockOpen} />
    </Button>
  );
};
