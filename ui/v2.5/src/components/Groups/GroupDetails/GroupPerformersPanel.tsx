import React from "react";
import * as GQL from "src/core/generated-graphql";
import { useGroupFilterHook } from "src/core/groups";
import { EnhancedPerformerList as PerformerList } from "src/extensions/facets/enhanced";
import { View } from "src/components/List/views";

interface IGroupPerformersPanel {
  active: boolean;
  group: GQL.GroupDataFragment;
  showChildGroupContent?: boolean;
}

export const GroupPerformersPanel: React.FC<IGroupPerformersPanel> = ({
  active,
  group,
  showChildGroupContent,
}) => {
  const filterHook = useGroupFilterHook(group, showChildGroupContent);

  return (
    <PerformerList
      filterHook={filterHook}
      alterQuery={active}
      view={View.GroupPerformers}
    />
  );
};
