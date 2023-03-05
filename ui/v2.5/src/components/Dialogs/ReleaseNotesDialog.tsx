import React from "react";
import { ModalComponent } from "../Shared/Modal";
import { faCogs } from "@fortawesome/free-solid-svg-icons";
import { useIntl } from "react-intl";
import { MarkdownPage } from "../Shared/MarkdownPage";
import { IReleaseNotes } from "src/docs/en/ReleaseNotes";

interface IReleaseNotesDialog {
  notes: IReleaseNotes[];
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
      <div className="m-n3">
        {notes
          .map((n, i) => (
            <div key={i} className="m-3">
              <h3>{n.version}</h3>
              <MarkdownPage page={n.content} />
            </div>
          ))
          .reduce((accu, curr) => (
            <>
              {accu}
              <hr />
              {curr}
            </>
          ))}
      </div>
    </ModalComponent>
  );
};

export default ReleaseNotesDialog;
