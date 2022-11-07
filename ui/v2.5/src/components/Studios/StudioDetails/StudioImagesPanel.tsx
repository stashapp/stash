import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useStudioFilterHook } from "src/core/studios";
import { ImageList } from "src/components/Images/ImageList";

interface IStudioImagesPanel {
  studio: GQL.StudioDataFragment;
}

export const StudioImagesPanel: React.FC<IStudioImagesPanel> = ({ studio }) => {
  const filterHook = useStudioFilterHook(studio);
  return <ImageList filterHook={filterHook} />;
};
