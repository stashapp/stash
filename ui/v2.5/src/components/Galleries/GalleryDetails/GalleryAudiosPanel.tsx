import React from "react";
import * as GQL from "src/core/generated-graphql";
import { AudioCard } from "src/components/Audios/AudioCard";

interface IGalleryAudiosPanelProps {
  audios: GQL.SlimAudioDataFragment[];
}

export const GalleryAudiosPanel: React.FC<IGalleryAudiosPanelProps> = ({
  audios,
}) => (
  <div className="container gallery-audios">
    {audios.map((audio) => (
      <AudioCard audio={audio} key={audio.id} />
    ))}
  </div>
);
