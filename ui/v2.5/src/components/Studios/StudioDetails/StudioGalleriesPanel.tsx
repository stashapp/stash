import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleryList } from "src/components/Galleries/GalleryList";
import { studioFilterHook } from "src/core/studios";

interface IStudioGalleriesPanel {
  studio: GQL.StudioDataFragment;
}

export const StudioGalleriesPanel: React.FC<IStudioGalleriesPanel> = ({
  studio,
}) => {
  return <GalleryList filterHook={studioFilterHook(studio)} />;
};
