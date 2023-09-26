import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneList } from "src/components/Scenes/SceneList";
import { useTagFilterHook } from "src/core/tags";
import { PersistanceLevel } from "src/components/List/ItemList";

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
      persistState={PersistanceLevel.SAVEDLINKFILTER}
    />
  );
};
