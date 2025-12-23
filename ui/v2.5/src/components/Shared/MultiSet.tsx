import React from "react";
import { IntlShape, useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { Button, ButtonGroup } from "react-bootstrap";
import { FilterSelect, SelectObject } from "./Select";
import {
  GalleryIDSelect,
  excludeFileBasedGalleries,
} from "../Galleries/GallerySelect";
import { PerformerIDSelect } from "../Performers/PerformerSelect";
import { StudioIDSelect } from "../Studios/StudioSelect";
import { TagIDSelect } from "../Tags/TagSelect";
import { GroupIDSelect } from "../Groups/GroupSelect";

interface IMultiSetProps {
  type: "performers" | "studios" | "tags" | "groups" | "galleries";
  existingIds?: string[];
  ids?: string[];
  mode: GQL.BulkUpdateIdMode;
  disabled?: boolean;
  onUpdate: (ids: string[]) => void;
  onSetMode: (mode: GQL.BulkUpdateIdMode) => void;
  menuPortalTarget?: HTMLElement | null;
}

const Select: React.FC<IMultiSetProps> = (props) => {
  const { type, disabled } = props;

  function onUpdate(items: SelectObject[]) {
    props.onUpdate(items.map((i) => i.id));
  }

  switch (type) {
    case "performers":
      return (
        <PerformerIDSelect
          isDisabled={disabled}
          isMulti
          isClearable={false}
          onSelect={onUpdate}
          ids={props.ids ?? []}
          menuPortalTarget={props.menuPortalTarget}
        />
      );
    case "studios":
      return (
        <StudioIDSelect
          isDisabled={disabled}
          isMulti
          isClearable={false}
          onSelect={onUpdate}
          ids={props.ids ?? []}
          menuPortalTarget={props.menuPortalTarget}
        />
      );
    case "tags":
      return (
        <TagIDSelect
          isDisabled={disabled}
          isMulti
          isClearable={false}
          onSelect={onUpdate}
          ids={props.ids ?? []}
          menuPortalTarget={props.menuPortalTarget}
        />
      );
    case "groups":
      return (
        <GroupIDSelect
          isDisabled={disabled}
          isMulti
          isClearable={false}
          onSelect={onUpdate}
          ids={props.ids ?? []}
          menuPortalTarget={props.menuPortalTarget}
        />
      );
    case "galleries":
      return (
        <GalleryIDSelect
          isDisabled={disabled}
          isMulti
          isClearable={false}
          onSelect={onUpdate}
          ids={props.ids ?? []}
          // exclude file-based galleries when setting galleries
          extraCriteria={excludeFileBasedGalleries}
          menuPortalTarget={props.menuPortalTarget}
        />
      );
    default:
      return (
        <FilterSelect
          type={type}
          isDisabled={disabled}
          isMulti
          isClearable={false}
          onSelect={onUpdate}
          ids={props.ids ?? []}
          menuPortalTarget={props.menuPortalTarget}
        />
      );
  }
};

function getModeText(intl: IntlShape, mode: GQL.BulkUpdateIdMode) {
  switch (mode) {
    case GQL.BulkUpdateIdMode.Set:
      return intl.formatMessage({
        id: "actions.overwrite",
        defaultMessage: "Overwrite",
      });
    case GQL.BulkUpdateIdMode.Add:
      return intl.formatMessage({ id: "actions.add", defaultMessage: "Add" });
    case GQL.BulkUpdateIdMode.Remove:
      return intl.formatMessage({
        id: "actions.remove",
        defaultMessage: "Remove",
      });
  }
}

export const MultiSetModeButton: React.FC<{
  mode: GQL.BulkUpdateIdMode;
  active: boolean;
  onClick: () => void;
  disabled?: boolean;
}> = ({ mode, active, onClick, disabled }) => {
  const intl = useIntl();

  return (
    <Button
      key={mode}
      variant="primary"
      active={active}
      size="sm"
      onClick={onClick}
      disabled={disabled}
    >
      {getModeText(intl, mode)}
    </Button>
  );
};

const modes = [
  GQL.BulkUpdateIdMode.Set,
  GQL.BulkUpdateIdMode.Add,
  GQL.BulkUpdateIdMode.Remove,
];

export const MultiSetModeButtons: React.FC<{
  mode: GQL.BulkUpdateIdMode;
  onSetMode: (mode: GQL.BulkUpdateIdMode) => void;
  disabled?: boolean;
}> = ({ mode, onSetMode, disabled }) => {
  return (
    <ButtonGroup className="button-group-above">
      {modes.map((m) => (
        <MultiSetModeButton
          key={m}
          mode={m}
          active={mode === m}
          onClick={() => onSetMode(m)}
          disabled={disabled}
        />
      ))}
    </ButtonGroup>
  );
};

export const MultiSet: React.FC<IMultiSetProps> = (props) => {
  const { mode, onUpdate, existingIds } = props;

  function onSetMode(m: GQL.BulkUpdateIdMode) {
    if (m === mode) {
      return;
    }

    // if going to Set, set the existing ids
    if (m === GQL.BulkUpdateIdMode.Set && existingIds) {
      onUpdate(existingIds);
      // if going from Set, wipe the ids
    } else if (
      m !== GQL.BulkUpdateIdMode.Set &&
      mode === GQL.BulkUpdateIdMode.Set
    ) {
      onUpdate([]);
    }

    props.onSetMode(m);
  }

  return (
    <div className="multi-set">
      <MultiSetModeButtons mode={mode} onSetMode={onSetMode} />
      <Select {...props} />
    </div>
  );
};
