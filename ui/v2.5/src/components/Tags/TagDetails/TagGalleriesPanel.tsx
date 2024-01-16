import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { GalleryList } from "src/components/Galleries/GalleryList";
import { View } from "src/components/List/views";

interface ITagGalleriesPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
}

export const TagGalleriesPanel: React.FC<ITagGalleriesPanel> = ({
  active,
  tag,
}) => {
  const filterHook = useTagFilterHook(tag);
  return (
    <GalleryList
      filterHook={filterHook}
      alterQuery={active}
      view={View.TagGalleries}
    />
  );
};
