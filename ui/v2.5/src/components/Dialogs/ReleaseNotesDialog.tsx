import React from "react";
import { Form } from "react-bootstrap";
import { ModalComponent } from "../Shared/Modal";
import { faCogs } from "@fortawesome/free-solid-svg-icons";
import { useIntl } from "react-intl";
import { MarkdownPage } from "../Shared/MarkdownPage";

interface IReleaseNotesDialog {
  notes: string[];
  onClose: () => void;
}

export const ReleaseNotesDialog: React.FC<IReleaseNotesDialog> = ({
  notes,
  onClose,
}) => {
  const intl = useIntl();

  return (
    <ModalComponent
      show
      icon={faCogs}
      header={intl.formatMessage({ id: "release_notes" })}
      accept={{
        onClick: onClose,
        text: intl.formatMessage({ id: "actions.close" }),
      }}
    >
      <Form>
        {notes.map((n, i) => (
          <MarkdownPage page={n} key={i} />
        ))}
      </Form>
    </ModalComponent>
  );
};

export default ReleaseNotesDialog;
