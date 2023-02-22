import React, { useState } from "react";
import { Form } from "react-bootstrap";
import { mutateExportObjects } from "src/core/StashService";
import { ModalComponent } from "./Modal";
import { useToast } from "src/hooks/Toast";
import downloadFile from "src/utils/download";
import { ExportObjectsInput } from "src/core/generated-graphql";
import { useIntl } from "react-intl";
import { faCogs } from "@fortawesome/free-solid-svg-icons";

interface IExportDialogProps {
  exportInput: ExportObjectsInput;
  onClose: () => void;
}

export const ExportDialog: React.FC<IExportDialogProps> = (
  props: IExportDialogProps
) => {
  const [includeDependencies, setIncludeDependencies] = useState(true);

  // Network state
  const [isRunning, setIsRunning] = useState(false);

  const intl = useIntl();
  const Toast = useToast();

  async function onExport() {
    try {
      setIsRunning(true);
      const ret = await mutateExportObjects({
        ...props.exportInput,
        includeDependencies,
      });

      // download the result
      if (ret.data && ret.data.exportObjects) {
        const link = ret.data.exportObjects;
        downloadFile(link);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsRunning(false);
      props.onClose();
    }
  }

  return (
    <ModalComponent
      show
      icon={faCogs}
      header={intl.formatMessage({ id: "dialogs.export_title" })}
      accept={{
        onClick: onExport,
        text: intl.formatMessage({ id: "actions.export" }),
      }}
      cancel={{
        onClick: () => props.onClose(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      isRunning={isRunning}
    >
      <Form>
        <Form.Group>
          <Form.Check
            id="include-dependencies"
            checked={includeDependencies}
            label={intl.formatMessage({
              id: "dialogs.export_include_related_objects",
            })}
            onChange={() => setIncludeDependencies(!includeDependencies)}
          />
        </Form.Group>
      </Form>
    </ModalComponent>
  );
};
