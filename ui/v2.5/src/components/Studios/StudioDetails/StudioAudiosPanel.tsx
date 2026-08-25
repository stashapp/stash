import React from "react";
import * as GQL from "src/core/generated-graphql";
import { FilteredAudioList } from "src/components/Audios/AudioList";
import { useStudioFilterHook } from "src/core/studios";
import { View } from "src/components/List/views";

interface IStudioAudiosPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
  showChildStudioContent?: boolean;
}

export const StudioAudiosPanel: React.FC<IStudioAudiosPanel> = ({
  active,
  studio,
  showChildStudioContent,
}) => {
  const filterHook = useStudioFilterHook(studio, showChildStudioContent);
  return (
    <FilteredAudioList
      filterHook={filterHook}
      alterQuery={active}
      view={View.StudioAudios}
    />
  );
};
