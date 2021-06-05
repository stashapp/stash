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
  type: "performers" | "studios" | "tags";
  ids?: string[];
  mode: GQL.BulkUpdateIdMode;
  disabled?: boolean;
  onUpdate: (items: ValidTypes[]) => void;
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
    props.onUpdate(items);
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

  function renderModeButton(mode: GQL.BulkUpdateIdMode) {
    return (
      <Button
        variant="primary"
        active={props.mode === mode}
        size="sm"
        onClick={() => props.onSetMode(mode)}
        disabled={props.disabled}
      >
        {getModeText(mode)}
      </Button>
    );
  }

  return (
    <div>
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
