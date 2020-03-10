import { Button } from "react-bootstrap";
import React, { FunctionComponent } from "react";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { useToast } from "src/hooks";

interface IOperationsPanelProps {
  scene: GQL.SceneDataFragment;
  playerPosition?: number;
}

export const SceneOperationsPanel: FunctionComponent<IOperationsPanelProps> = (props: IOperationsPanelProps) => {
  const Toast = useToast();
  const [generateScreenshot] = StashService.useSceneGenerateScreenshot();

  async function onGenerateScreenshot() {
    await generateScreenshot({
      variables: {
        id: props.scene.id,
        at: props.playerPosition
      }
    });
    Toast.success({ content: "Generating screenshot" });
  }

  async function onGenerateDefaultScreenshot() {
    await generateScreenshot({
      variables: {
        id: props.scene.id,
      }
    });
    Toast.success({ content: "Generating screenshot" });
  }

  return (
    <>
      <Button className="edit-button" onClick={() => onGenerateScreenshot()}>
        Generate thumbnail from current
      </Button>
      <Button className="edit-button" onClick={() => onGenerateDefaultScreenshot()}>
        Generate default thumbnail
      </Button>
    </>
  );
};
