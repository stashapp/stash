import {
  IconDefinition,
  faMinus,
  faPencilAlt,
  faPlus,
} from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, Col, Form, Row } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "src/components/Shared/Icon";
import { ModalComponent } from "src/components/Shared/Modal";
import { FolderSelect } from "src/components/Shared/FolderSelect/FolderSelect";
import { ConfigurationContext } from "src/hooks/Config";

interface IDirectorySelectionDialogProps {
  animation?: boolean;
  initialPaths?: string[];
  allowEmpty?: boolean;
  allowPathSelection?: boolean;
  message?: string;
  header?: string;
  icon?: IconDefinition;
  acceptButtonText?: string;
  acceptButtonVariant?: "danger" | "primary" | "secondary";
  onClose: (paths?: string[]) => void;
}

export const DirectorySelectionDialog: React.FC<
  IDirectorySelectionDialogProps
> = ({
  animation,
  allowEmpty = false,
  initialPaths = [],
  allowPathSelection = true,
  message,
  header,
  icon = faPencilAlt,
  acceptButtonText,
  acceptButtonVariant = "primary",
  onClose,
}) => {
  const intl = useIntl();
  const { configuration } = React.useContext(ConfigurationContext);

  const libraryPaths = configuration?.general.stashes.map((s) => s.path);

  const [paths, setPaths] = useState<string[]>(initialPaths);
  const [currentDirectory, setCurrentDirectory] = useState<string>("");

  function removePath(p: string) {
    setPaths(paths.filter((path) => path !== p));
  }

  function addPath(p: string) {
    if (p && !paths.includes(p)) {
      setPaths(paths.concat(p));
    }
  }

  return (
    <ModalComponent
      show
      modalProps={{ animation }}
      disabled={!allowEmpty && allowPathSelection && paths.length === 0}
      icon={icon}
      header={header ?? intl.formatMessage({ id: "actions.select_folders" })}
      accept={{
        onClick: () => {
          onClose(paths);
        },
        text: acceptButtonText ?? intl.formatMessage({ id: "actions.confirm" }),
        variant: acceptButtonVariant,
      }}
      cancel={{
        onClick: () => onClose(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
    >
      <div className="dialog-container">
        {paths.map((p) => (
          <Row className="align-items-center mb-1" key={p}>
            <Form.Label column xs={10}>
              {p}
            </Form.Label>
            <Col xs={2} className="d-flex justify-content-end">
              <Button
                className="ml-auto"
                size="sm"
                variant="danger"
                title={intl.formatMessage({ id: "actions.delete" })}
                onClick={() => removePath(p)}
              >
                <Icon icon={faMinus} />
              </Button>
            </Col>
          </Row>
        ))}

        {allowPathSelection ? (
          <FolderSelect
            currentDirectory={currentDirectory}
            setCurrentDirectory={(v) => setCurrentDirectory(v)}
            defaultDirectories={libraryPaths}
            appendButton={
              <Button
                variant="secondary"
                onClick={() => addPath(currentDirectory)}
              >
                <Icon icon={faPlus} />
              </Button>
            }
          />
        ) : undefined}
        {message}
      </div>
    </ModalComponent>
  );
};
