import cloneDeep from "lodash-es/cloneDeep";
import React, { useCallback, useMemo } from "react";
import { Form } from "react-bootstrap";
import { CriterionModifier } from "src/core/generated-graphql";
import {
  DurationCriterion,
  CriterionValue,
  Criterion,
  IHierarchicalLabeledIdCriterion,
  NumberCriterion,
  ILabeledIdCriterion,
  DateCriterion,
  TimestampCriterion,
  BooleanCriterion,
} from "src/models/list-filter/criteria/criterion";
import { useIntl } from "react-intl";
import {
  criterionIsHierarchicalLabelValue,
  criterionIsNumberValue,
  criterionIsStashIDValue,
  criterionIsDateValue,
  criterionIsTimestampValue,
} from "src/models/list-filter/types";
import { DurationFilter } from "./Filters/DurationFilter";
import { NumberFilter } from "./Filters/NumberFilter";
import { LabeledIdFilter } from "./Filters/LabeledIdFilter";
import { HierarchicalLabelValueFilter } from "./Filters/HierarchicalLabelValueFilter";
import { InputFilter } from "./Filters/InputFilter";
import { DateFilter } from "./Filters/DateFilter";
import { TimestampFilter } from "./Filters/TimestampFilter";
import { CountryCriterion } from "src/models/list-filter/criteria/country";
import { CountrySelect } from "../Shared/CountrySelect";
import { StashIDCriterion } from "src/models/list-filter/criteria/stash-ids";
import { StashIDFilter } from "./Filters/StashIDFilter";
import { RatingCriterion } from "../../models/list-filter/criteria/rating";
import { RatingFilter } from "./Filters/RatingFilter";
import { BooleanFilter } from "./Filters/BooleanFilter";
import { OptionsListFilter } from "./Filters/OptionsListFilter";

interface IGenericCriterionEditor {
  criterion: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
}

const GenericCriterionEditor: React.FC<IGenericCriterionEditor> = ({
  criterion,
  setCriterion,
}) => {
  const intl = useIntl();

  const { options, modifierOptions } = criterion.criterionOption;

  const onChangedModifierSelect = useCallback(
    (event: React.ChangeEvent<HTMLSelectElement>) => {
      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = event.target.value as CriterionModifier;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const modifierSelector = useMemo(() => {
    if (!modifierOptions || modifierOptions.length === 0) {
      return;
    }

    return (
      <Form.Control
        as="select"
        onChange={onChangedModifierSelect}
        value={criterion.modifier}
        className="btn-secondary modifier-selector"
      >
        {modifierOptions.map((c) => (
          <option key={c.value} value={c.value}>
            {c.label ? intl.formatMessage({ id: c.label }) : ""}
          </option>
        ))}
      </Form.Control>
    );
  }, [modifierOptions, onChangedModifierSelect, criterion.modifier, intl]);

  const valueControl = useMemo(() => {
    function onValueChanged(value: CriterionValue) {
      const newCriterion = cloneDeep(criterion);
      newCriterion.value = value;
      setCriterion(newCriterion);
    }

    // always show stashID filter
    if (criterion instanceof StashIDCriterion) {
      return (
        <StashIDFilter criterion={criterion} onValueChanged={onValueChanged} />
      );
    }

    // Hide the value select if the modifier is "IsNull" or "NotNull"
    if (
      criterion.modifier === CriterionModifier.IsNull ||
      criterion.modifier === CriterionModifier.NotNull
    ) {
      return;
    }

    if (criterion instanceof ILabeledIdCriterion) {
      return (
        <LabeledIdFilter
          criterion={criterion}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (criterion instanceof IHierarchicalLabeledIdCriterion) {
      return (
        <HierarchicalLabelValueFilter
          criterion={criterion}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (
      options &&
      !criterionIsHierarchicalLabelValue(criterion.value) &&
      !criterionIsNumberValue(criterion.value) &&
      !criterionIsStashIDValue(criterion.value) &&
      !criterionIsDateValue(criterion.value) &&
      !criterionIsTimestampValue(criterion.value) &&
      !Array.isArray(criterion.value)
    ) {
      // if (!modifierOptions || modifierOptions.length === 0) {
      return (
        <OptionsListFilter criterion={criterion} setCriterion={setCriterion} />
      );
      // }

      // return (
      //   <OptionsFilter criterion={criterion} onValueChanged={onValueChanged} />
      // );
    }
    if (criterion instanceof DurationCriterion) {
      return (
        <DurationFilter criterion={criterion} onValueChanged={onValueChanged} />
      );
    }
    if (criterion instanceof DateCriterion) {
      return (
        <DateFilter criterion={criterion} onValueChanged={onValueChanged} />
      );
    }
    if (criterion instanceof TimestampCriterion) {
      return (
        <TimestampFilter
          criterion={criterion}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (criterion instanceof NumberCriterion) {
      return (
        <NumberFilter criterion={criterion} onValueChanged={onValueChanged} />
      );
    }
    if (criterion instanceof RatingCriterion) {
      return (
        <RatingFilter criterion={criterion} onValueChanged={onValueChanged} />
      );
    }
    if (
      criterion instanceof CountryCriterion &&
      (criterion.modifier === CriterionModifier.Equals ||
        criterion.modifier === CriterionModifier.NotEquals)
    ) {
      return (
        <CountrySelect
          value={criterion.value}
          onChange={(v) => onValueChanged(v)}
        />
      );
    }
    return (
      <InputFilter criterion={criterion} onValueChanged={onValueChanged} />
    );
  }, [criterion, setCriterion, options]);

  return (
    <div>
      {modifierSelector}
      {valueControl}
    </div>
  );
};

interface ICriterionEditor {
  criterion: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
}

export const CriterionEditor: React.FC<ICriterionEditor> = ({
  criterion,
  setCriterion,
}) => {
  const filterControl = useMemo(() => {
    if (criterion instanceof BooleanCriterion) {
      return (
        <BooleanFilter criterion={criterion} setCriterion={setCriterion} />
      );
    }

    return (
      <GenericCriterionEditor
        criterion={criterion}
        setCriterion={setCriterion}
      />
    );
  }, [criterion, setCriterion]);

  return <div className="criterion-editor">{filterControl}</div>;
};
