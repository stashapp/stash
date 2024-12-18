import React from "react";
import * as GQL from "src/core/generated-graphql";
import { MultiSetModeButtons } from "../Shared/MultiSet";
import {
  IRelatedGroupEntry,
  RelatedGroupTable,
} from "./GroupDetails/RelatedGroupTable";
import { Group, GroupSelect } from "./GroupSelect";

export const ContainingGroupsMultiSet: React.FC<{
  existingValue?: IRelatedGroupEntry[];
  value: IRelatedGroupEntry[];
  mode: GQL.BulkUpdateIdMode;
  disabled?: boolean;
  onUpdate: (value: IRelatedGroupEntry[]) => void;
  onSetMode: (mode: GQL.BulkUpdateIdMode) => void;
  menuPortalTarget?: HTMLElement | null;
}> = (props) => {
  const { mode, onUpdate, existingValue } = props;

  function onSetMode(m: GQL.BulkUpdateIdMode) {
    if (m === mode) {
      return;
    }

    // if going to Set, set the existing ids
    if (m === GQL.BulkUpdateIdMode.Set && existingValue) {
      onUpdate(existingValue);
      // if going from Set, wipe the ids
    } else if (
      m !== GQL.BulkUpdateIdMode.Set &&
      mode === GQL.BulkUpdateIdMode.Set
    ) {
      onUpdate([]);
    }

    props.onSetMode(m);
  }

  function onRemoveSet(items: Group[]) {
    onUpdate(items.map((group) => ({ group })));
  }

  return (
    <div className="multi-set">
      <MultiSetModeButtons mode={mode} onSetMode={onSetMode} />
      {mode !== GQL.BulkUpdateIdMode.Remove ? (
        <RelatedGroupTable
          value={props.value}
          onUpdate={props.onUpdate}
          disabled={props.disabled}
          menuPortalTarget={props.menuPortalTarget}
        />
      ) : (
        <GroupSelect
          onSelect={(items) => onRemoveSet(items)}
          values={[]}
          isDisabled={props.disabled}
          menuPortalTarget={props.menuPortalTarget}
        />
      )}
    </div>
  );
};
