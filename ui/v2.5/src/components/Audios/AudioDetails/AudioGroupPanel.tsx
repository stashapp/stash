import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GroupCard } from "src/components/Groups/GroupCard";

interface IAudioGroupPanelProps {
  audio: GQL.AudioDataFragment;
}

export const AudioGroupPanel: React.FC<IAudioGroupPanelProps> = (
  props: IAudioGroupPanelProps
) => {
  const cards = props.audio.groups.map((audioGroup) => (
    <GroupCard key={audioGroup.group.id} group={audioGroup.group} />
  ));

  return <div className="row justify-content-center">{cards}</div>;
};

export default AudioGroupPanel;
