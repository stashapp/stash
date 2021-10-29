import { Button } from "react-bootstrap";
import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { mutateMetadataAutoTag } from "src/core/StashService";
import { useToast } from "src/hooks";

interface IPerformerOperationsProps {
  performer: Partial<GQL.PerformerDataFragment>;
}

export const PerformerOperationsPanel: React.FC<IPerformerOperationsProps> = ({
  performer,
}) => {
  const intl = useIntl();
  const Toast = useToast();

  async function onAutoTag() {
    if (!performer?.id) {
      return;
    }
    try {
      await mutateMetadataAutoTag({ performers: [performer.id] });
      Toast.success({ content: intl.formatMessage({ id: "toast.started_auto_tagging" }) });
    } catch (e) {
      Toast.error(e);
    }
  }

  return (
    <Button onClick={onAutoTag}>
      <FormattedMessage id="actions.auto_tag" />
    </Button>
  );
};
