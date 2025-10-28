import { Button } from "react-bootstrap";
import React, { useState } from "react";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { SubmitStashBoxDraft } from "src/components/Dialogs/SubmitDraft";

interface IPerformerOperationsProps {
  performer: GQL.PerformerDataFragment;
}

export const PerformerSubmitButton: React.FC<IPerformerOperationsProps> = ({
  performer,
}) => {
  const [showDraftModal, setShowDraftModal] = useState(false);

  const { data } = GQL.useConfigurationQuery();
  const boxes = data?.configuration?.general?.stashBoxes ?? [];

  if (boxes.length === 0) return null;

  return (
    <>
      <Button onClick={() => setShowDraftModal(true)}>
        <FormattedMessage id="actions.submit_stash_box" />
      </Button>
      <SubmitStashBoxDraft
        type="performer"
        boxes={boxes}
        entity={performer}
        show={showDraftModal}
        onHide={() => setShowDraftModal(false)}
      />
    </>
  );
};
