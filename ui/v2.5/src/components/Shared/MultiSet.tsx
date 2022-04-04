import * as React from "react";
import { useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { Button, ButtonGroup } from "react-bootstrap";
import { FilterSelect } from "./Select";

type ValidTypes =
  | GQL.SlimPerformerDataFragment
  | GQL.SlimTagDataFragment
  | GQL.SlimStudioDataFragment
  | GQL.SlimMovieDataFragment;

interface IMultiSetProps {
  type: "performers" | "studios" | "tags" | "movies";
  existingIds?: string[];
  ids?: string[];
  mode: GQL.BulkUpdateIdMode;
  disabled?: boolean;
  onUpdate: (ids: string[]) => void;
  onSetMode: (mode: GQL.BulkUpdateIdMode) => void;
}

const MultiSet: React.FunctionComponent<IMultiSetProps> = (
  props: IMultiSetProps
) => {
  const intl = useIntl();
  const modes = [
    GQL.BulkUpdateIdMode.Set,
    GQL.BulkUpdateIdMode.Add,
    GQL.BulkUpdateIdMode.Remove,
  ];

  function onUpdate(items: ValidTypes[]) {
    props.onUpdate(items.map((i) => i.id));
  }

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
    if (mode === GQL.BulkUpdateIdMode.Set && props.existingIds) {
      props.onUpdate(props.existingIds);
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
      <FilterSelect
        type={props.type}
        isDisabled={props.disabled}
        isMulti
        isClearable={false}
        onSelect={onUpdate}
        ids={props.ids ?? []}
      />
    </div>
  );
};

export default MultiSet;
