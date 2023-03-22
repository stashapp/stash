import React from "react";
import { Form } from "react-bootstrap";
import { FolderSelect } from "src/components/Shared/FolderSelect/FolderSelect";
import { ConfigurationContext } from "src/hooks/Config";
import {
  Criterion,
  CriterionValue,
} from "../../../models/list-filter/criteria/criterion";

interface IInputFilterProps {
  criterion: Criterion<CriterionValue>;
  onValueChanged: (value: string) => void;
}

export const PathFilter: React.FC<IInputFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const { configuration } = React.useContext(ConfigurationContext);
  const libraryPaths = configuration?.general.stashes.map((s) => s.path);

  return (
    <Form.Group>
      <FolderSelect
        currentDirectory={criterion.value ? criterion.value.toString() : ""}
        setCurrentDirectory={(v) => onValueChanged(v)}
        collapsible
        defaultDirectories={libraryPaths}
      />
    </Form.Group>
  );
};
