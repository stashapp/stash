import React from "react";
import * as GQL from "src/core/generated-graphql";
import { FilteredAudioList } from "src/components/Audios/AudioList";
import { useTagFilterHook } from "src/core/tags";
import { View } from "src/components/List/views";

interface ITagAudiosPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
  showSubTagContent?: boolean;
}

export const TagAudiosPanel: React.FC<ITagAudiosPanel> = ({
  active,
  tag,
  showSubTagContent,
}) => {
  const filterHook = useTagFilterHook(tag, showSubTagContent);
  return (
    <FilteredAudioList
      filterHook={filterHook}
      alterQuery={active}
      view={View.TagAudios}
    />
  );
};
