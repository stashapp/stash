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
import { useStringCriterion } from "./StringFilter";

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

  const {criterion, setCriterion} = useStringCriterion(option, filter, setFilter);
  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption.defaultModifier;
  const modifierOptions = modifierCriterionOption.modifierOptions;

  function onValueChange(value: string) {
    if (!value.trim()) {
      // Remove criterion if empty
      setFilter(filter.removeCriterion(option.type));
      return;
    }

    // const newCriterion = criterion ? criterion.clone() : option.makeCriterion();
    const newCriterion = cloneDeep(criterion);
    newCriterion.modifier = criterion?.modifier ? criterion.modifier : defaultModifier;
    newCriterion.value = value;
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  const onChangedModifierSelect = useCallback(
    (m: CriterionModifier) => {
      console.log("onChangedModifierSelect", m);
      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = m;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  // check if we should show regex input or folder select
  const regex =
    criterion?.modifier === CriterionModifier.MatchesRegex ||
    criterion?.modifier === CriterionModifier.NotMatchesRegex;

  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      outsideCollapse={
        <ul className="selected-list">
          {criterion?.modifier != defaultModifier ? (
            <SelectedItem
              className="modifier-object"
              label={ModifierCriterion.getModifierLabel(intl, criterion.modifier)}
              onClick={() => onChangedModifierSelect(defaultModifier)}
            />
          ) : null}
          {criterion?.value ? (
            <SelectedItem
              label={criterion.value}
              onClick={() => onValueChange("")}
            />
          ) : null}
        </ul>
      }
    >
      <div className="path-filter">
        <div className="filter-group">
          <ModifierSelectorButtons
            options={modifierOptions}
            value={criterion?.modifier ? criterion.modifier : defaultModifier}
            onChanged={onChangedModifierSelect}
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
