import React from "react";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared";
import { SceneDataFragment } from "src/core/generated-graphql";
import { ConfigurationContext } from "src/hooks/Config";

export interface IVRButtonProps {
  scene: SceneDataFragment;
}

export const VRButton: React.FC<IVRButtonProps> = ({ scene }) => {
  const config = React.useContext(ConfigurationContext);
  const { paths } = scene;

  // DeoVR urls + android app requires https.
  // To use without https, use the desktop app to access stash's website root.
  if (!paths || !paths.deovr || !window.isSecureContext) {
    return <span />;
  }

  if (!config.configuration?.interface.showSceneDeoVRButton) {
    return <span />;
  }

  return (
    <Button className="minimal" variant="secondary" title="Open in DeoVR">
      <a href={paths.deovr}>
        <Icon icon="vr-cardboard" color="white" />
      </a>
    </Button>
  );
};
