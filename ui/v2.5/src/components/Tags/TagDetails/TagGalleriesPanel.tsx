import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { EnhancedGalleryList as GalleryList } from "src/extensions/facets/enhanced";
import { View } from "src/components/List/views";

interface ITagGalleriesPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
  showSubTagContent?: boolean;
}

export const TagGalleriesPanel: React.FC<ITagGalleriesPanel> = ({
  active,
  tag,
  showSubTagContent,
}) => {
  const filterHook = useTagFilterHook(tag, showSubTagContent);
  return (
    <GalleryList
      filterHook={filterHook}
      alterQuery={active}
      view={View.TagGalleries}
    />
  );
};
