import React, { useContext } from "react";
import * as GQL from "src/core/generated-graphql";
import { Button } from "react-bootstrap";
import { useIntl } from "react-intl";
import { faLock, faLockOpen } from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";
import { Icon } from "./Icon";

export interface ILockContext {
  locked: boolean;
  fieldName: string;
  onToggle: () => void;
}

// set by lockable form fields. Controls inside read it to disable themselves
// and to render the lock toggle inside the input.
export const LockContext = React.createContext<ILockContext | undefined>(
  undefined
);

export function useFieldLocked() {
  return !!useContext(LockContext)?.locked;
}

// renders nothing outside of a lockable field
export const LockFieldButton: React.FC = () => {
  const intl = useIntl();
  const lock = useContext(LockContext);
  if (!lock) return null;

  const { locked, fieldName, onToggle } = lock;
  const title = intl.formatMessage(
    { id: locked ? "actions.unlock_field" : "actions.lock_field" },
    { field: fieldName }
  );

  // stop events reaching the containing control (eg react-select opens its
  // menu on mouse down and swallows space/enter)
  const stopPropagation = (e: React.SyntheticEvent) => e.stopPropagation();

  return (
    <Button
      variant="link"
      className={cx("minimal", "lock-field-button", { locked })}
      title={title}
      aria-label={title}
      aria-pressed={locked}
      onClick={onToggle}
      onMouseDown={stopPropagation}
      onTouchEnd={stopPropagation}
      onKeyDown={stopPropagation}
    >
      <Icon icon={locked ? faLock : faLockOpen} />
    </Button>
  );
};

// locale ids of the lockable scene fields
export const lockableFieldMessageIDs: Record<GQL.SceneLockableField, string> = {
  [GQL.SceneLockableField.Title]: "title",
  [GQL.SceneLockableField.Code]: "scene_code",
  [GQL.SceneLockableField.Details]: "details",
  [GQL.SceneLockableField.Director]: "director",
  [GQL.SceneLockableField.Date]: "date",
  [GQL.SceneLockableField.ProductionDate]: "production_date",
  [GQL.SceneLockableField.Studio]: "studio",
  [GQL.SceneLockableField.Urls]: "urls",
  [GQL.SceneLockableField.Performers]: "performers",
  [GQL.SceneLockableField.Tags]: "tags",
  [GQL.SceneLockableField.Groups]: "groups",
  [GQL.SceneLockableField.Galleries]: "galleries",
};
