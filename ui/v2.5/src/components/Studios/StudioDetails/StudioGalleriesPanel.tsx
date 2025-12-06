import React from "react";
import * as GQL from "src/core/generated-graphql";
import { EnhancedGalleryList as GalleryList } from "src/extensions/facets/enhanced";
import { useStudioFilterHook } from "src/core/studios";
import { View } from "src/components/List/views";

interface IStudioGalleriesPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
  showChildStudioContent?: boolean;
}

export const StudioGalleriesPanel: React.FC<IStudioGalleriesPanel> = ({
  active,
  studio,
  showChildStudioContent,
}) => {
  const filterHook = useStudioFilterHook(studio, showChildStudioContent);
  return (
    <GalleryList
      filterHook={filterHook}
      alterQuery={active}
      view={View.StudioGalleries}
    />
  );
};
