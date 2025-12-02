import React, { useCallback, useMemo } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { IStashIDValue } from "../../../models/list-filter/types";
import {
  ModifierCriterion,
  CriterionOption,
} from "../../../models/list-filter/criteria/criterion";
import { StashIDCriterion } from "../../../models/list-filter/criteria/stash-ids";
import { CriterionModifier } from "src/core/generated-graphql";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SidebarSection } from "src/components/Shared/Sidebar";
import { SelectedItem } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";
import { ModifierControls } from "./StringFilter";

interface IStashIDFilterProps {
  criterion: ModifierCriterion<IStashIDValue>;
  onValueChanged: (value: IStashIDValue) => void;
}

// Hook for stash ID-based sidebar filters
export function useStashIDCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as StashIDCriterion;

    const newCriterion = filter.makeCriterion(option.type) as StashIDCriterion;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: StashIDCriterion) => {
      const newFilter = cloneDeep(filter);

      // replace or add the criterion
      const newCriteria = filter.criteria.filter((cc) => {
        return cc.criterionOption.type !== c.criterionOption.type;
      });
      newCriteria.push(c);
      newFilter.criteria = newCriteria;
      setFilter(newFilter);
    },
    [filter, setFilter]
  );

  const onValueChanged = useCallback(
    (value: IStashIDValue) => {
      const newCriterion = criterion.clone();
      newCriterion.value = value;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const onChangedModifierSelect = useCallback(
    (modifier: CriterionModifier) => {
      const newCriterion = criterion.clone();
      newCriterion.modifier = modifier;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const modifierCriterionOption = criterion?.modifierCriterionOption();
  const defaultModifier = modifierCriterionOption?.defaultModifier;
  const modifierOptions = modifierCriterionOption?.modifierOptions;

  return {
    criterion,
    setCriterion,
    onValueChanged,
    onChangedModifierSelect,
    defaultModifier,
    modifierOptions,
  };
}

export const StashIDFilter: React.FC<IStashIDFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();
  const { value } = criterion;

  function onEndpointChanged(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged({
      endpoint: event.target.value,
      stashID: criterion.value.stashID,
    });
  }

  function onStashIDChanged(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged({
      stashID: event.target.value,
      endpoint: criterion.value.endpoint,
    });
  }

  return (
    <div>
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          onChange={onEndpointChanged}
          value={value ? value.endpoint : ""}
          placeholder={intl.formatMessage({ id: "stash_id_endpoint" })}
        />
      </Form.Group>
      {criterion.modifier !== CriterionModifier.IsNull &&
        criterion.modifier !== CriterionModifier.NotNull && (
          <Form.Group>
            <Form.Control
              className="btn-secondary"
              onChange={onStashIDChanged}
              value={value ? value.stashID : ""}
              placeholder={intl.formatMessage({ id: "stash_id" })}
            />
          </Form.Group>
        )}
    </div>
  );
};

const StashIDSelectedItems: React.FC<{
  criterion: StashIDCriterion;
  defaultModifier: CriterionModifier;
  onChangedModifierSelect: (modifier: CriterionModifier) => void;
  onClear: () => void;
}> = ({ criterion, defaultModifier, onChangedModifierSelect, onClear }) => {
  if (
    criterion?.modifier !== CriterionModifier.IsNull &&
    criterion?.modifier !== CriterionModifier.NotNull &&
    criterion?.value.endpoint === ""
  )
    return null;
  const intl = useIntl();

  const valueLabel = useMemo(() => {
    if (!criterion.value) return null;

    const { endpoint, stashID } = criterion.value;
    if (!endpoint && !stashID) return null;

    if (endpoint && stashID) {
      return `${endpoint}: ${stashID}`;
    } else if (endpoint) {
      return `Endpoint: ${endpoint}`;
    } else if (stashID) {
      return `ID: ${stashID}`;
    }

    return null;
  }, [criterion.value]);

  return (
    <ul className="selected-list">
      {criterion?.modifier != defaultModifier && criterion?.modifier ? (
        <SelectedItem
          className="modifier-object"
          label={ModifierCriterion.getModifierLabel(intl, criterion.modifier)}
          onClick={() => onChangedModifierSelect(defaultModifier)}
        />
      ) : null}
      {valueLabel ? (
        <SelectedItem label={valueLabel} onClick={onClear} />
      ) : null}
    </ul>
  );
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

export const SidebarStashIDFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const {
    criterion,
    defaultModifier,
    modifierOptions,
    onValueChanged,
    onChangedModifierSelect,
  } = useStashIDCriterion(option, filter, setFilter);

  const modifierSelector = useMemo(() => {
    return (
      <ModifierControls
        modifierOptions={modifierOptions}
        currentModifier={criterion.modifier}
        onChangedModifierSelect={onChangedModifierSelect}
      />
    );
  }, [modifierOptions, onChangedModifierSelect, criterion.modifier]);

  const valueControl = useMemo(() => {
    return (
      <StashIDFilter criterion={criterion} onValueChanged={onValueChanged} />
    );
  }, [criterion, onValueChanged]);

  const onClear = useCallback(() => {
    setFilter(filter.removeCriterion(option.type));
  }, [filter, setFilter, option.type]);

  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      sectionID={sectionID}
      outsideCollapse={
        <StashIDSelectedItems
          criterion={criterion}
          defaultModifier={defaultModifier}
          onChangedModifierSelect={onChangedModifierSelect}
          onClear={onClear}
        />
      }
    >
      <div className="stash-id-filter">
        <div className="filter-group">
          {modifierSelector}
          {valueControl}
        </div>
      </div>
    </SidebarSection>
  );
};
