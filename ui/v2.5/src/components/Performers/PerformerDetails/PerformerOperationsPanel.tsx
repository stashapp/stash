import { Button } from "react-bootstrap";
import React from "react";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { useToast } from "src/hooks";

interface IPerformerOperationsProps {
  performer: Partial<GQL.PerformerDataFragment>;
}

export const PerformerOperationsPanel: React.FC<IPerformerOperationsProps> = ({
  performer
}) => {
  const Toast = useToast();

  async function onAutoTag() {
    if (!performer?.id) {
      return;
    }
    try {
      await StashService.mutateMetadataAutoTag({ performers: [performer.id] });
      Toast.success({ content: "Started auto tagging" });
    } catch (e) {
      Toast.error(e);
    }
  }

  return <Button onClick={onAutoTag}>Auto Tag</Button>;
};
