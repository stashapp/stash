import * as React from "react";

import * as GQL from "src/core/generated-graphql";
import { Button, InputGroup } from "react-bootstrap";
import { Icon } from "src/components/Shared";
import { FilterSelect } from "./Select";

type ValidTypes =
  | GQL.SlimPerformerDataFragment
  | GQL.Tag
  | GQL.SlimStudioDataFragment;

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
  function onUpdate(items: ValidTypes[]) {
    props.onUpdate(items);
  }

  function getModeIcon() {
    switch (props.mode) {
      case GQL.BulkUpdateIdMode.Set:
        return "pencil-alt";
      case GQL.BulkUpdateIdMode.Add:
        return "plus";
      case GQL.BulkUpdateIdMode.Remove:
        return "times";
    }
  }

  function getModeText() {
    switch (props.mode) {
      case GQL.BulkUpdateIdMode.Set:
        return "Set";
      case GQL.BulkUpdateIdMode.Add:
        return "Add";
      case GQL.BulkUpdateIdMode.Remove:
        return "Remove";
    }
  }

  function nextMode() {
    switch (props.mode) {
      case GQL.BulkUpdateIdMode.Set:
        return GQL.BulkUpdateIdMode.Add;
      case GQL.BulkUpdateIdMode.Add:
        return GQL.BulkUpdateIdMode.Remove;
      case GQL.BulkUpdateIdMode.Remove:
        return GQL.BulkUpdateIdMode.Set;
    }
  }

  return (
    <InputGroup className="multi-set">
      <InputGroup.Prepend>
        <Button
          size="sm"
          variant="secondary"
          onClick={() => props.onSetMode(nextMode())}
          title={getModeText()}
          disabled={props.disabled}
        >
          <Icon icon={getModeIcon()} className="fa-fw" />
        </Button>
      </InputGroup.Prepend>

      <FilterSelect
        type={props.type}
        isDisabled={props.disabled}
        isMulti
        isClearable={false}
        onSelect={onUpdate}
        ids={props.ids ?? []}
      />
    </InputGroup>
  );
};

export default MultiSet;
