import React from "react";
import * as GQL from "src/core/generated-graphql";
import { GalleryList } from "src/components/Galleries/GalleryList";
import { useStudioFilterHook } from "src/core/studios";
import { View } from "src/components/List/views";

interface IStudioGalleriesPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
}

export const StudioGalleriesPanel: React.FC<IStudioGalleriesPanel> = ({
  active,
  studio,
}) => {
  const filterHook = useStudioFilterHook(studio);
  return (
    <GalleryList
      filterHook={filterHook}
      alterQuery={active}
      view={View.StudioGalleries}
    />
  );
};
