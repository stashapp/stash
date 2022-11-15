import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneList } from "src/components/Scenes/SceneList";
import { useTagFilterHook } from "src/core/tags";

interface ITagScenesPanel {
  tag: GQL.TagDataFragment;
}

export const TagScenesPanel: React.FC<ITagScenesPanel> = ({ tag }) => {
  const filterHook = useTagFilterHook(tag);
  return <SceneList filterHook={filterHook} />;
};
