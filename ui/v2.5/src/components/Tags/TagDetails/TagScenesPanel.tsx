import React from "react";
import * as GQL from "src/core/generated-graphql";
import { MyFilteredSceneList } from "src/components/Scenes/MySceneList";
import { useTagFilterHook } from "src/core/tags";
import { View } from "src/components/List/views";

interface ITagScenesPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
  showSubTagContent?: boolean;
}

export const TagScenesPanel: React.FC<ITagScenesPanel> = ({
  active,
  tag,
  showSubTagContent,
}) => {
  const filterHook = useTagFilterHook(tag, showSubTagContent);
  return (
    <MyFilteredSceneList
      filterHook={filterHook}
      alterQuery={active}
      view={View.TagScenes}
    />
  );
};
