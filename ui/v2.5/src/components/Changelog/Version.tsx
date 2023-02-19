import { faAngleDown, faAngleUp } from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, Card, Collapse } from "react-bootstrap";
import { FormattedDate, FormattedMessage } from "react-intl";
import { Icon } from "src/components/Shared/Icon";

interface IVersionProps {
  version: string;
  date?: string;
  defaultOpen?: boolean;
  setOpenState: (key: string, state: boolean) => void;
  openState: Record<string, boolean>;
}

const Version: React.FC<IVersionProps> = ({
  version,
  date,
  defaultOpen,
  openState,
  setOpenState,
  children,
}) => {
  const [open, setOpen] = useState(
    defaultOpen ?? openState[version + date] ?? false
  );

  const updateState = () => {
    setOpenState(version + date, !open);
    setOpen(!open);
  };

  return (
    <Card className="changelog-version">
      <Card.Header>
        <h4 className="changelog-version-header d-flex align-items-center">
          <Button onClick={updateState} variant="link">
            <Icon icon={open ? faAngleUp : faAngleDown} className="mr-3" />
            {version} (
            {date ? (
              <FormattedDate value={date} timeZone="utc" />
            ) : (
              <FormattedMessage
                defaultMessage="Development Version"
                id="developmentVersion"
              />
            )}
            )
          </Button>
        </h4>
      </Card.Header>
      <Card.Body>
        <Collapse in={open}>
          <div className="changelog-version-body markdown">{children}</div>
        </Collapse>
      </Card.Body>
    </Card>
  );
};

export default Version;
