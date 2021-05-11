import React from "react";

import {SceneDataFragment} from "../../core/generated-graphql";
import {Interactive} from "../../utils/interactive";

export function ScenePlayerInteractiveControls(props: {
  interactive: Interactive,
  scene: SceneDataFragment,
}) {
  if (!props.scene.interactive) {
    return null;
  }
  return <div>
    <p>Current connection key: { props.interactive.handyKey }</p>
  </div>;
}
