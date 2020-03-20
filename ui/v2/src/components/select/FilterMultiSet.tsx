import * as React from "react";

import { ControlGroup, Button } from "@blueprintjs/core";
import * as GQL from "../../core/generated-graphql";
import { FilterMultiSelect } from "./FilterMultiSelect";

type ValidTypes =
  GQL.AllPerformersForFilterAllPerformers |
  GQL.AllTagsForFilterAllTags |
  GQL.AllMoviesForFilterAllMovies | 
  GQL.AllStudiosForFilterAllStudios;

interface IFilterMultiSetProps {
  type: "performers" | "studios" | "movies" | "tags";
  initialIds?: string[];
  mode: GQL.BulkUpdateIdMode;
  onUpdate: (items: ValidTypes[]) => void;
  onSetMode: (mode: GQL.BulkUpdateIdMode) => void;
}

export const FilterMultiSet: React.FunctionComponent<IFilterMultiSetProps> = (props: IFilterMultiSetProps) => {
  function onUpdate(items: ValidTypes[]) {
    props.onUpdate(items);
  }

  function getModeIcon() {
    switch(props.mode) {
      case GQL.BulkUpdateIdMode.Set:
        return "edit";
      case GQL.BulkUpdateIdMode.Add:
        return "plus";
      case GQL.BulkUpdateIdMode.Remove:
        return "cross";
    }
  }

  function getModeText() {
    switch(props.mode) {
      case GQL.BulkUpdateIdMode.Set:
        return "Set";
      case GQL.BulkUpdateIdMode.Add:
        return "Add";
      case GQL.BulkUpdateIdMode.Remove:
        return "Remove";
    }
  }

  function nextMode() {
    switch(props.mode) {
      case GQL.BulkUpdateIdMode.Set:
        return GQL.BulkUpdateIdMode.Add;
      case GQL.BulkUpdateIdMode.Add:
        return GQL.BulkUpdateIdMode.Remove;
      case GQL.BulkUpdateIdMode.Remove:
        return GQL.BulkUpdateIdMode.Set;
    }
  }

  return (
    <ControlGroup>
      <Button 
        icon={getModeIcon()} 
        minimal={true} 
        onClick={() => props.onSetMode(nextMode())}
        title={getModeText()}
      />
      <FilterMultiSelect
        type={props.type}
        initialIds={props.initialIds}
        onUpdate={onUpdate}
      />
    </ControlGroup>
  );
};
