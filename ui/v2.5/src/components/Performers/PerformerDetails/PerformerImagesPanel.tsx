import React from "react";
import * as GQL from "src/core/generated-graphql";
import { ImageList } from "src/components/Images/ImageList";
import { usePerformerFilterHook } from "src/core/performers";
import { View } from "src/components/List/views";
import { PatchComponent } from "src/patch";

interface IPerformerImagesPanel {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerImagesPanel: React.FC<IPerformerImagesPanel> =
  PatchComponent("PerformerImagesPanel", ({ active, performer }) => {
    const filterHook = usePerformerFilterHook(performer);
    return (
      <ImageList
        filterHook={filterHook}
        alterQuery={active}
        view={View.PerformerImages}
      />
    );
  });
