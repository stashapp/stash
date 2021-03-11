import React from "react";
import * as GQL from "src/core/generated-graphql";
import { tagFilterHook } from "src/core/tags";
import { GalleryList } from "src/components/Galleries/GalleryList";

interface ITagGalleriesPanel {
  tag: GQL.TagDataFragment;
}

export const TagGalleriesPanel: React.FC<ITagGalleriesPanel> = ({ tag }) => {
  return <GalleryList filterHook={tagFilterHook(tag)} />;
};
