import React from "react";
import { useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { Button, ButtonGroup } from "react-bootstrap";
import { FilterSelect, SelectObject } from "./Select";
import {
  GalleryIDSelect,
  excludeFileBasedGalleries,
} from "../Galleries/GallerySelect";
import { StringListInput } from "./StringListInput";

interface IMultiSetProps {
  existing?: string[];
  mode: GQL.BulkUpdateIdMode;
  disabled?: boolean;
  onUpdate: (ids: string[]) => void;
  onSetMode: (mode: GQL.BulkUpdateIdMode) => void;
}

interface IMultiSelectProps extends IMultiSetProps {
  type: "performers" | "studios" | "tags" | "movies" | "galleries" | "scenes";
  ids?: string[];
}

const Select: React.FC<IMultiSelectProps> = (props) => {
  const { type, disabled } = props;

  function onUpdate(items: SelectObject[]) {
    props.onUpdate(items.map((i) => i.id));
  }

  if (type === "galleries") {
    return (
      <GalleryIDSelect
        isDisabled={disabled}
        isMulti
        isClearable={false}
        onSelect={onUpdate}
        ids={props.ids ?? []}
        // exclude file-based galleries when setting galleries
        extraCriteria={excludeFileBasedGalleries}
      />
    );
  }

  return (
    <FilterSelect
      type={type}
      isDisabled={disabled}
      isMulti
      isClearable={false}
      onSelect={onUpdate}
      ids={props.ids ?? []}
    />
  );
};

const MultiSet: React.FC<IMultiSetProps> = (props) => {
  const intl = useIntl();
  const modes = [
    GQL.BulkUpdateIdMode.Set,
    GQL.BulkUpdateIdMode.Add,
    GQL.BulkUpdateIdMode.Remove,
  ];

  function getModeText(mode: GQL.BulkUpdateIdMode) {
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

  function onSetMode(mode: GQL.BulkUpdateIdMode) {
    if (mode === props.mode) {
      return;
    }

    // if going to Set, set the existing ids
    if (mode === GQL.BulkUpdateIdMode.Set && props.existing) {
      props.onUpdate(props.existing);
      // if going from Set, wipe the ids
    } else if (
      mode !== GQL.BulkUpdateIdMode.Set &&
      props.mode === GQL.BulkUpdateIdMode.Set
    ) {
      props.onUpdate([]);
    }

    props.onSetMode(mode);
  }

  function renderModeButton(mode: GQL.BulkUpdateIdMode) {
    return (
      <Button
        key={mode}
        variant="primary"
        active={props.mode === mode}
        size="sm"
        onClick={() => onSetMode(mode)}
        disabled={props.disabled}
      >
        {getModeText(mode)}
      </Button>
    );
  }

  return (
    <div className="multi-set">
      <ButtonGroup className="button-group-above">
        {modes.map((m) => renderModeButton(m))}
      </ButtonGroup>
    </div>
  );
};

interface IMultiStringProps extends IMultiSetProps {
  strings?: string[];
}

export const MultiString: React.FC<IMultiStringProps> = (props) => {
    return (
    <div className="multi-string">
      <MultiSet
        mode={props.mode}
        existing={props.existing}
        onUpdate={props.onUpdate}
        onSetMode={props.onSetMode}
        disabled={props.disabled}
      />
      <StringListInput
        value={props.strings ?? []}
        setValue={props.onUpdate}
        readOnly={props.disabled}
        className={"input-control"}
      />
    </div>
  );
};

export const MultiSelect: React.FC<IMultiSelectProps> = (props) => {
  return (
    <div className="multi-set">
      <MultiSet
        mode={props.mode}
        existing={props.existing}
        onUpdate={props.onUpdate}
        onSetMode={props.onSetMode}
        disabled={props.disabled}
      />
      <Select {...props} />
    </div>
  );
};
