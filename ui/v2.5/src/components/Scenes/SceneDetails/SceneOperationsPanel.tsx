import { Button } from "react-bootstrap";
import React, { FunctionComponent } from "react";
import * as GQL from "src/core/generated-graphql";
import { useSceneGenerateScreenshot } from "src/core/StashService";
import { useToast } from "src/hooks";
import { JWUtils } from "src/utils";

interface IOperationsPanelProps {
  scene: GQL.SceneDataFragment;
}

export const SceneOperationsPanel: FunctionComponent<IOperationsPanelProps> = (
  props: IOperationsPanelProps
) => {
  const Toast = useToast();
  const [generateScreenshot] = useSceneGenerateScreenshot();

  async function onGenerateScreenshot(at?: number) {
    await generateScreenshot({
      variables: {
        id: props.scene.id,
        at,
      },
    });
    Toast.success({ content: "Generating screenshot" });
  }

  return (
    <>
      <Button
        className="edit-button"
        onClick={() => onGenerateScreenshot(JWUtils.getPlayer().getPosition())}
      >
        Generate thumbnail from current
      </Button>
      <Button className="edit-button" onClick={() => onGenerateScreenshot()}>
        Generate default thumbnail
      </Button>
    </>
  );
};
