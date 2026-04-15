import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { FilteredImageList } from "src/components/Images/ImageList";
import { View } from "src/components/List/views";

interface ITagImagesPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
  showSubTagContent?: boolean;
}

export const TagImagesPanel: React.FC<ITagImagesPanel> = ({
  active,
  tag,
  showSubTagContent,
}) => {
  const filterHook = useTagFilterHook(tag, showSubTagContent);
  return (
    <FilteredImageList
      filterHook={filterHook}
      alterQuery={active}
      view={View.TagImages}
    />
  );
};
