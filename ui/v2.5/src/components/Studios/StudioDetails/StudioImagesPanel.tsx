import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useStudioFilterHook } from "src/core/studios";
import { ImageList } from "src/components/Images/ImageList";
import { View } from "src/components/List/views";

interface IStudioImagesPanel {
  active: boolean;
  studio: GQL.StudioDataFragment;
  showChildStudioContent?: boolean;
}

export const StudioImagesPanel: React.FC<IStudioImagesPanel> = ({
  active,
  studio,
  showChildStudioContent,
}) => {
  const filterHook = useStudioFilterHook(studio, showChildStudioContent);
  return (
    <ImageList
      filterHook={filterHook}
      alterQuery={active}
      view={View.StudioImages}
    />
  );
};
