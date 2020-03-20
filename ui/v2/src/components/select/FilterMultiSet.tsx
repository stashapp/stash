import * as React from "react";

import { MenuItem, ControlGroup, Button } from "@blueprintjs/core";
import { IMultiSelectProps, ItemPredicate, ItemRenderer, MultiSelect } from "@blueprintjs/select";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { HTMLInputProps } from "../../models";
import { ErrorUtils } from "../../utils/errors";
import { ToastUtils } from "../../utils/toasts";
import { FilterMultiSelect } from "./FilterMultiSelect";

const InternalPerformerMultiSelect = MultiSelect.ofType<GQL.AllPerformersForFilterAllPerformers>();
const InternalTagMultiSelect = MultiSelect.ofType<GQL.AllTagsForFilterAllTags>();
const InternalStudioMultiSelect = MultiSelect.ofType<GQL.AllStudiosForFilterAllStudios>();
const InternalMovieMultiSelect = MultiSelect.ofType<GQL.AllMoviesForFilterAllMovies>();

type ValidTypes =
  GQL.AllPerformersForFilterAllPerformers |
  GQL.AllTagsForFilterAllTags |
  GQL.AllMoviesForFilterAllMovies | 
  GQL.AllStudiosForFilterAllStudios;

export enum MultiSetMode {
    SET = "Set",
    ADD = "Add",
    REMOVE = "Remove"
}

interface IFilterMultiSetProps {
  type: "performers" | "studios" | "movies" | "tags";
  initialIds?: string[];
  mode: MultiSetMode;
  onUpdate: (items: ValidTypes[]) => void;
  onSetMode: (mode: MultiSetMode) => void;
}

export const FilterMultiSet: React.FunctionComponent<IFilterMultiSetProps> = (props: IFilterMultiSetProps) => {
  function onUpdate(items: ValidTypes[]) {
    props.onUpdate(items);
  }

  function getModeIcon() {
    switch(props.mode) {
      case MultiSetMode.SET:
        return "edit";
      case MultiSetMode.ADD:
        return "plus";
      case MultiSetMode.REMOVE:
        return "cross";
    }
  }

  function nextMode() {
    switch(props.mode) {
      case MultiSetMode.SET:
        return MultiSetMode.ADD;
      case MultiSetMode.ADD:
        return MultiSetMode.REMOVE;
      case MultiSetMode.REMOVE:
        return MultiSetMode.SET;
    }
  }

  return (
    <ControlGroup>
      <Button 
        icon={getModeIcon()} 
        minimal={true} 
        onClick={() => props.onSetMode(nextMode())}
        title={props.mode}
      />
      <FilterMultiSelect
        type={props.type}
        initialIds={props.initialIds}
        onUpdate={onUpdate}
      />
    </ControlGroup>
  );
};
