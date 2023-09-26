import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useStudioFilterHook } from "src/core/studios";
import { ImageList } from "src/components/Images/ImageList";
import { PersistanceLevel } from "src/components/List/ItemList";

interface IStudioImagesPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
}

export const StudioImagesPanel: React.FC<IStudioImagesPanel> = ({
  active,
  studio,
}) => {
  const filterHook = useStudioFilterHook(studio);
  return (
    <ImageList
      filterHook={filterHook}
      alterQuery={active}
      persistState={PersistanceLevel.SAVEDLINKFILTER}
    />
  );
};
