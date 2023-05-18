import React from "react";
import { Form } from "react-bootstrap";
import { FolderSelect } from "src/components/Shared/FolderSelect/FolderSelect";
import { CriterionModifier } from "src/core/generated-graphql";
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

  // don't show folder select for regex
  const regex =
    criterion.modifier === CriterionModifier.MatchesRegex ||
    criterion.modifier === CriterionModifier.NotMatchesRegex;

  return (
    <Form.Group>
      {regex ? (
        <Form.Control
          className="btn-secondary"
          type={criterion.criterionOption.inputType}
          onChange={(v) => onValueChanged(v.target.value)}
          value={criterion.value ? criterion.value.toString() : ""}
        />
      ) : (
        <FolderSelect
          currentDirectory={criterion.value ? criterion.value.toString() : ""}
          setCurrentDirectory={(v) => onValueChanged(v)}
          collapsible
          quoteSpaced
          hideError
          defaultDirectories={libraryPaths}
        />
      )}
    </Form.Group>
  );
};
