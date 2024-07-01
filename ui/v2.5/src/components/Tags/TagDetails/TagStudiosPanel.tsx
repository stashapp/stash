import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useTagFilterHook } from "src/core/tags";
import { StudioList } from "src/components/Studios/StudioList";

interface ITagStudiosPanel {
  active: boolean;
  tag: GQL.TagDataFragment;
}

export const TagStudiosPanel: React.FC<ITagStudiosPanel> = ({
  active,
  tag,
}) => {
  const filterHook = useTagFilterHook(tag);
  return <StudioList filterHook={filterHook} alterQuery={active} />;
};
