import cloneDeep from "lodash-es/cloneDeep";
import React, { useCallback, useMemo } from "react";
import { Button, Form } from "react-bootstrap";
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
import { OptionFilter, OptionListFilter } from "./Filters/OptionFilter";
import { PathFilter } from "./Filters/PathFilter";
import { PerformersCriterion } from "src/models/list-filter/criteria/performers";
import PerformersFilter from "./Filters/PerformersFilter";
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";
import StudiosFilter from "./Filters/StudiosFilter";
import { TagsCriterion } from "src/models/list-filter/criteria/tags";
import TagsFilter from "./Filters/TagsFilter";
import { PhashCriterion } from "src/models/list-filter/criteria/phash";
import { PhashFilter } from "./Filters/PhashFilter";
import cx from "classnames";
import { PathCriterion } from "src/models/list-filter/criteria/path";

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

  const showModifierSelector = useMemo(() => {
    if (
      criterion instanceof PerformersCriterion ||
      criterion instanceof StudiosCriterion ||
      criterion instanceof TagsCriterion
    ) {
      return false;
    }

    return modifierOptions && modifierOptions.length > 1;
  }, [criterion, modifierOptions]);

  const alwaysShowFilter = useMemo(() => {
    return (
      criterion instanceof StashIDCriterion ||
      criterion instanceof PerformersCriterion ||
      criterion instanceof StudiosCriterion ||
      criterion instanceof TagsCriterion
    );
  }, [criterion]);

  const onChangedModifierSelect = useCallback(
    (m: CriterionModifier) => {
      const newCriterion = cloneDeep(criterion);
      newCriterion.modifier = m;
      setCriterion(newCriterion);
    },
    [criterion, setCriterion]
  );

  const modifierSelector = useMemo(() => {
    if (!showModifierSelector) {
      return;
    }

    return (
      <Form.Group className="modifier-options">
        {modifierOptions.map((m) => (
          <Button
            className={cx("modifier-option", {
              selected: criterion.modifier === m,
            })}
            key={m}
            onClick={() => onChangedModifierSelect(m)}
          >
            {Criterion.getModifierLabel(intl, m)}
          </Button>
        ))}
      </Form.Group>
    );
  }, [
    showModifierSelector,
    modifierOptions,
    onChangedModifierSelect,
    criterion.modifier,
    intl,
  ]);

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
      !alwaysShowFilter &&
      (criterion.modifier === CriterionModifier.IsNull ||
        criterion.modifier === CriterionModifier.NotNull)
    ) {
      return;
    }

    if (criterion instanceof PerformersCriterion) {
      return (
        <PerformersFilter
          criterion={criterion}
          setCriterion={(c) => setCriterion(c)}
        />
      );
    }

    if (criterion instanceof StudiosCriterion) {
      return (
        <StudiosFilter
          criterion={criterion}
          setCriterion={(c) => setCriterion(c)}
        />
      );
    }

    if (criterion instanceof TagsCriterion) {
      return (
        <TagsFilter
          criterion={criterion}
          setCriterion={(c) => setCriterion(c)}
        />
      );
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
      !criterionIsTimestampValue(criterion.value)
    ) {
      if (!Array.isArray(criterion.value)) {
        return (
          <OptionFilter criterion={criterion} setCriterion={setCriterion} />
        );
      } else {
        return (
          <OptionListFilter criterion={criterion} setCriterion={setCriterion} />
        );
      }
    }
    if (criterion instanceof PathCriterion) {
      return (
        <PathFilter criterion={criterion} onValueChanged={onValueChanged} />
      );
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
    if (criterion instanceof PhashCriterion) {
      return (
        <PhashFilter criterion={criterion} onValueChanged={onValueChanged} />
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
          menuPortalTarget={document.body}
        />
      );
    }
    return (
      <InputFilter criterion={criterion} onValueChanged={onValueChanged} />
    );
  }, [criterion, setCriterion, options, alwaysShowFilter]);

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
