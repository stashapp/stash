import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleryList } from "src/components/Galleries/GalleryList";
import { useStudioFilterHook } from "src/core/studios";

interface IStudioGalleriesPanel {
  studio: GQL.StudioDataFragment;
}

export const StudioGalleriesPanel: React.FC<IStudioGalleriesPanel> = ({
  studio,
}) => {
  const filterHook = useStudioFilterHook(studio);
  return <GalleryList filterHook={filterHook} />;
};
