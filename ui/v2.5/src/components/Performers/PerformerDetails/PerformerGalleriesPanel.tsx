import React from "react";
import * as GQL from "src/core/generated-graphql";
import { EnhancedGalleryList as GalleryList } from "src/extensions/facets/enhanced";
import { usePerformerFilterHook } from "src/core/performers";
import { View } from "src/components/List/views";
import { PatchComponent } from "src/patch";

interface IPerformerDetailsProps {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerGalleriesPanel: React.FC<IPerformerDetailsProps> =
  PatchComponent("PerformerGalleriesPanel", ({ active, performer }) => {
    const filterHook = usePerformerFilterHook(performer);
    return (
      <GalleryList
        filterHook={filterHook}
        alterQuery={active}
        view={View.PerformerGalleries}
      />
    );
  });
