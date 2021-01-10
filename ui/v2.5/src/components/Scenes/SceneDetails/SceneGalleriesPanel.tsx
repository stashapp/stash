import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleryCard } from "src/components/Galleries/GalleryCard";

interface ISceneGalleriesPanelProps {
  galleries: GQL.GallerySlimDataFragment[];
}

export const SceneGalleriesPanel: React.FC<ISceneGalleriesPanelProps> = ({ galleries }) => {
  const cards = galleries.map((gallery) => (
    <GalleryCard
      key={gallery.id}
      gallery={gallery}
      selecting={false}
    />
  ));

  return (
    <div className="row justify-content-center">{cards}</div>
  );
};
