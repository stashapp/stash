import {
    Button,
  } from "@blueprintjs/core";
  import React, { FunctionComponent } from "react";
  import * as GQL from "../../../core/generated-graphql";
  import { StashService } from "../../../core/StashService";
  import { SceneHelpers } from "../helpers";
import { ToastUtils } from "../../../utils/toasts";
  
  interface IOperationsPanelProps {
    scene: GQL.SceneDataFragment;
  }
  
  export const SceneOperationsPanel: FunctionComponent<IOperationsPanelProps> = (props: IOperationsPanelProps) => {

    const jwplayer = SceneHelpers.getPlayer();
    const generateScreenshot = StashService.useSceneGenerateScreenshot();

    async function onGenerateScreenshot() {
      let position = jwplayer.getPosition();

      await generateScreenshot({
        variables: {
          id: props.scene.id,
          at: position
        }
      });
      ToastUtils.success("Generating screenshot");
    }
    
    async function onGenerateDefaultScreenshot() {
      await generateScreenshot({
        variables: {
          id: props.scene.id,
        }
      });
      ToastUtils.success("Generating screenshot");
    }

    return (
      <>
        <Button className="edit-button" text="Generate thumbnail from current" onClick={() => onGenerateScreenshot()}/>
        <Button className="edit-button" text="Generate default thumbnail" onClick={() => onGenerateDefaultScreenshot()}/>
      </>
    );
  };
  