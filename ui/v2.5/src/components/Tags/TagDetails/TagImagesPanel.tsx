import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { ImageList } from "src/components/Images/ImageList";

interface ITagImagesPanel {
  tag: GQL.TagDataFragment;
}

export const TagImagesPanel: React.FC<ITagImagesPanel> = ({ tag }) => {
  const filterHook = useTagFilterHook(tag);
  return <ImageList filterHook={filterHook} />;
};
