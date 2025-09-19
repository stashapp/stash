import React, { useCallback, useMemo, useState } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { FolderSelect } from "src/components/Shared/FolderSelect/FolderSelect";
import { CriterionModifier } from "src/core/generated-graphql";
import { ConfigurationContext } from "src/hooks/Config";
import {
  ModifierCriterion,
  CriterionValue,
  CriterionOption,
} from "../../../models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SidebarSection } from "src/components/Shared/Sidebar";
import { ModifierSelectorButtons } from "../ModifierSelect";
import { cloneDeep } from "lodash-es";
import { SelectedItem, SelectedList } from "./SidebarListFilter";
import { useModifierCriterion, SelectedItems, ModifierControls } from "./StringFilter";

interface IInputFilterProps {
  criterion: ModifierCriterion<CriterionValue>;
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
          type={criterion.modifierCriterionOption().inputType}
          onChange={(v) => onValueChanged(v.target.value)}
          value={criterion.value ? criterion.value.toString() : ""}
        />
      ) : (
        <FolderSelect
          currentDirectory={criterion.value ? criterion.value.toString() : ""}
          onChangeDirectory={onValueChanged}
          collapsible
          quotePath
          hideError
          defaultDirectories={libraryPaths}
        />
      )}
    </Form.Group>
  );
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}

export const SidebarPathFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
}) => {
  const intl = useIntl();
  const { configuration } = React.useContext(ConfigurationContext);
  const libraryPaths = configuration?.general.stashes.map((s) => s.path);

  const {
    criterion,
    defaultModifier,
    modifierOptions,
    onValueChange,
    onChangedModifierSelect
  } = useModifierCriterion(option, filter, setFilter);

  // check if we should show regex input or folder select
  const regex =
    criterion?.modifier === CriterionModifier.MatchesRegex ||
    criterion?.modifier === CriterionModifier.NotMatchesRegex;

  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      outsideCollapse={
        <SelectedItems
          criterion={criterion}
          defaultModifier={defaultModifier}
          onChangedModifierSelect={onChangedModifierSelect}
          onValueChange={onValueChange}
        />
      }
    >
      <div className="path-filter">
        <div className="filter-group">
          <ModifierControls
            modifierOptions={modifierOptions}
            currentModifier={criterion?.modifier || defaultModifier}
            onChangedModifierSelect={onChangedModifierSelect}
          />
          {regex ? (
            <Form.Control
              className="btn-secondary"
              onChange={(v) => onValueChange(v.target.value)}
              value={criterion?.value || ""}
              placeholder={intl.formatMessage({ id: "path" })}
            />
          ) : (
            <FolderSelect
              currentDirectory={criterion?.value || ""}
              onChangeDirectory={onValueChange}
              collapsible
              quotePath
              hideError
              defaultDirectories={libraryPaths}
            />
          )}
        </div>
      </div>
    </SidebarSection>
  );
};
