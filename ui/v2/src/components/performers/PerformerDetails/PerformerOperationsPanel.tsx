import {
  Button,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { ErrorUtils } from "../../../utils/errors";
import { ToastUtils } from "../../../utils/toasts";

interface IPerformerOperationsProps {
  performer: Partial<GQL.PerformerDataFragment>
}

export const PerformerOperationsPanel: FunctionComponent<IPerformerOperationsProps> = (props: IPerformerOperationsProps) => {

  async function onAutoTag() {
    if (!props.performer || !props.performer.id) {
        return;
    }
    try {
        await StashService.mutateMetadataAutoTag({ performers: [props.performer.id]});
        ToastUtils.success("Started auto tagging");
    } catch (e) {
        ErrorUtils.handle(e);
    }
  }

  return (
    <>
      <Button text="Auto Tag" onClick={onAutoTag} />
    </>
  );
};
  