import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { ImageList } from "src/components/Images/ImageList";

interface ITagImagesPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
}

export const TagImagesPanel: React.FC<ITagImagesPanel> = ({ active, tag }) => {
  const filterHook = useTagFilterHook(tag);
  return <ImageList filterHook={filterHook} alterQuery={active} />;
};
