import React from "react";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared";
import { SceneDataFragment } from "src/core/generated-graphql";

export interface IVRButtonProps {
  scene: SceneDataFragment;
}

export const VRButton: React.FC<IVRButtonProps> = ({ scene }) => {
  const isAndroid = /(android)/i.test(navigator.userAgent);

  const { paths } = scene;

  if (!paths || !paths.deovr || !isAndroid) {
    return <span />;
  }

  return (
    <Button
      className="minimal"
      variant="secondary"
      title="Open in DeoVR"
    >
      <a href={paths.deovr}>
        <Icon icon="vr-cardboard" color="white" />
      </a>
    </Button>
  );
};
