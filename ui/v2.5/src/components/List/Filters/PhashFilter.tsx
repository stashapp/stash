import React, { useCallback, useMemo } from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { IPhashDistanceValue } from "../../../models/list-filter/types";
import {
  CriterionOption,
  ModifierCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { PhashCriterion } from "../../../models/list-filter/criteria/phash";
import { CriterionModifier } from "src/core/generated-graphql";
import { NumberField } from "src/utils/form";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SidebarSection } from "src/components/Shared/MySidebar";
import { SelectedItem } from "./SidebarListFilter";
import { cloneDeep } from "lodash-es";
import { ModifierControls } from "./StringFilter";

interface IPhashFilterProps {
  criterion: ModifierCriterion<IPhashDistanceValue>;
  onValueChanged: (value: IPhashDistanceValue) => void;
}

// Hook for phash-based sidebar filters
export function usePhashCriterion(
  option: CriterionOption,
  filter: ListFilterModel,
  setFilter: (f: ListFilterModel) => void
) {
  const criterion = useMemo(() => {
    const ret = filter.criteria.find(
      (c) => c.criterionOption.type === option.type
    );
    if (ret) return ret as PhashCriterion;

    const newCriterion = filter.makeCriterion(option.type) as PhashCriterion;
    return newCriterion;
  }, [filter, option]);

  const setCriterion = useCallback(
    (c: PhashCriterion) => {
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
    (value: IPhashDistanceValue) => {
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

export const PhashFilter: React.FC<IPhashFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();
  const { value } = criterion;

  function valueChanged(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged({
      value: event.target.value,
      distance: criterion.value.distance,
    });
  }

  function distanceChanged(event: React.ChangeEvent<HTMLInputElement>) {
    let distance = parseInt(event.target.value);
    if (distance < 0 || isNaN(distance)) {
      distance = 0;
    }

    onValueChanged({
      distance,
      value: criterion.value.value,
    });
  }

  return (
    <div>
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          onChange={valueChanged}
          value={value ? value.value : ""}
          placeholder={intl.formatMessage({ id: "media_info.phash" })}
        />
      </Form.Group>
      {criterion.modifier !== CriterionModifier.IsNull &&
        criterion.modifier !== CriterionModifier.NotNull && (
          <Form.Group>
            <NumberField
              className="btn-secondary"
              onChange={distanceChanged}
              value={value ? value.distance : ""}
              placeholder={intl.formatMessage({ id: "distance" })}
            />
          </Form.Group>
        )}
    </div>
  );
};

const PhashSelectedItems: React.FC<{
  criterion: PhashCriterion;
  defaultModifier: CriterionModifier;
  onChangedModifierSelect: (modifier: CriterionModifier) => void;
  onClear: () => void;
}> = ({ criterion, defaultModifier, onChangedModifierSelect, onClear }) => {
  if (
    criterion?.modifier !== CriterionModifier.IsNull &&
    criterion?.modifier !== CriterionModifier.NotNull &&
    criterion?.value.value === ""
  )
    return null;
  const intl = useIntl();

  const valueLabel = useMemo(() => {
    if (!criterion.value) return null;

    const { value, distance } = criterion.value;
    if (!value && distance === undefined) return null;

    if (value && distance !== undefined) {
      return `${value} (distance: ${distance})`;
    } else if (value) {
      return `Phash: ${value}`;
    } else if (
      criterion?.modifier !== CriterionModifier.IsNull &&
      criterion?.modifier !== CriterionModifier.NotNull &&
      distance !== undefined
    ) {
      return `Distance: ${distance}`;
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
}

export const SidebarPhashFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
}) => {
  const {
    criterion,
    defaultModifier,
    modifierOptions,
    onValueChanged,
    onChangedModifierSelect,
  } = usePhashCriterion(option, filter, setFilter);

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
      <PhashFilter criterion={criterion} onValueChanged={onValueChanged} />
    );
  }, [criterion, onValueChanged]);

  const onClear = useCallback(() => {
    setFilter(filter.removeCriterion(option.type));
  }, [filter, setFilter, option.type]);

  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      outsideCollapse={
        <PhashSelectedItems
          criterion={criterion}
          defaultModifier={defaultModifier}
          onChangedModifierSelect={onChangedModifierSelect}
          onClear={onClear}
        />
      }
    >
      <div className="phash-filter">
        <div className="filter-group">
          {modifierSelector}
          {valueControl}
        </div>
      </div>
    </SidebarSection>
  );
};
