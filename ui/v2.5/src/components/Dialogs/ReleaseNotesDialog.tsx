import React from "react";
import { Form } from "react-bootstrap";
import { Modal } from "src/components/Shared";
import { faCogs } from "@fortawesome/free-solid-svg-icons";
import { useIntl } from "react-intl";
import { MarkdownPage } from "../Shared/MarkdownPage";
import { Module } from "src/docs/en/ReleaseNotes";

interface IReleaseNotesDialog {
  notes: Module[];
  onClose: () => void;
}

export const ReleaseNotesDialog: React.FC<IReleaseNotesDialog> = ({
  notes,
  onClose,
}) => {
  const intl = useIntl();

  return (
    <Modal
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
    </Modal>
  );
};

export default ReleaseNotesDialog;
