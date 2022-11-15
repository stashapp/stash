import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { GalleryList } from "src/components/Galleries/GalleryList";

interface ITagGalleriesPanel {
  tag: GQL.TagDataFragment;
}

export const TagGalleriesPanel: React.FC<ITagGalleriesPanel> = ({ tag }) => {
  const filterHook = useTagFilterHook(tag);
  return <GalleryList filterHook={filterHook} />;
};
