import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneList } from "src/components/Scenes/SceneList";
import { useTagFilterHook } from "src/core/tags";
import { View } from "src/components/List/views";

interface ITagScenesPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
}

export const TagScenesPanel: React.FC<ITagScenesPanel> = ({ active, tag }) => {
  const filterHook = useTagFilterHook(tag);
  return (
    <SceneList
      filterHook={filterHook}
      alterQuery={active}
      view={View.TagScenes}
    />
  );
};
