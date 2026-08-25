import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleryCard } from "src/components/Galleries/GalleryCard";

interface IAudioGalleriesPanelProps {
  galleries: GQL.SlimGalleryDataFragment[];
}

export const AudioGalleriesPanel: React.FC<IAudioGalleriesPanelProps> = ({
  galleries,
}) => {
  const cards = galleries.map((gallery) => (
    <GalleryCard
      key={gallery.id}
      gallery={gallery}
      selecting={false}
      zoomIndex={2}
    />
  ));

  return <div className="container audio-galleries">{cards}</div>;
};

export default AudioGalleriesPanel;
