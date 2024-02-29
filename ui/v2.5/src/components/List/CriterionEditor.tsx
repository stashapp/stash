import cloneDeep from "lodash-es/cloneDeep";
import React, { useCallback, useEffect, useMemo } from "react";
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
  emptyCriterion: Criterion<CriterionValue>;
  criterion?: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
}

const GenericCriterionEditor: React.FC<IGenericCriterionEditor> = ({
  emptyCriterion,
  criterion: inputCriterion,
  setCriterion: setInputCriterion,
}) => {
  const intl = useIntl();

  const [criterionState, setCriterionState] = React.useState(
    inputCriterion ?? emptyCriterion
  );
  const { options, modifierOptions } = emptyCriterion.criterionOption;

  const setCriterion = useCallback(
    (c: Criterion<CriterionValue>) => {
      setCriterionState(c);

      if (c.isValid()) {
        setInputCriterion(c);
      } else if (c.equals(emptyCriterion)) {
        // remove the criterion
        setInputCriterion(emptyCriterion);
      }
    },
    [setInputCriterion, emptyCriterion]
  );

  useEffect(() => {
    setCriterionState(inputCriterion ?? emptyCriterion);
  }, [inputCriterion, emptyCriterion]);

  const onChangedModifierSelect = useCallback(
    (m: CriterionModifier) => {
      const newCriterion = cloneDeep(criterionState);
      newCriterion.modifier = m;
      setCriterion(newCriterion);
    },
    [criterionState, setCriterion]
  );

  const modifierSelector = useMemo(() => {
    if (!modifierOptions || modifierOptions.length === 0) {
      return;
    }

    return (
      <Form.Group className="modifier-options">
        {modifierOptions.map((m) => (
          <Button
            className={cx("modifier-option", {
              selected: criterionState.modifier === m,
            })}
            key={m}
            onClick={() => onChangedModifierSelect(m)}
          >
            {Criterion.getModifierLabel(intl, m)}
          </Button>
        ))}
      </Form.Group>
    );
  }, [modifierOptions, onChangedModifierSelect, criterionState.modifier, intl]);

  const valueControl = useMemo(() => {
    function onValueChanged(value: CriterionValue) {
      const newCriterion = cloneDeep(criterionState);
      newCriterion.value = value;
      setCriterion(newCriterion);
    }

    // always show stashID filter
    if (criterionState instanceof StashIDCriterion) {
      return (
        <StashIDFilter
          criterion={criterionState}
          onValueChanged={onValueChanged}
        />
      );
    }

    // Hide the value select if the modifier is "IsNull" or "NotNull"
    if (
      criterionState.modifier === CriterionModifier.IsNull ||
      criterionState.modifier === CriterionModifier.NotNull
    ) {
      return;
    }

    if (criterionState instanceof PerformersCriterion) {
      return (
        <PerformersFilter
          criterion={criterionState}
          setCriterion={(c) => setCriterion(c)}
        />
      );
    }

    if (criterionState instanceof StudiosCriterion) {
      return (
        <StudiosFilter
          criterion={criterionState}
          setCriterion={(c) => setCriterion(c)}
        />
      );
    }

    if (criterionState instanceof TagsCriterion) {
      return (
        <TagsFilter
          criterion={criterionState}
          setCriterion={(c) => setCriterion(c)}
        />
      );
    }

    if (criterionState instanceof ILabeledIdCriterion) {
      return (
        <LabeledIdFilter
          criterion={criterionState}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (criterionState instanceof IHierarchicalLabeledIdCriterion) {
      return (
        <HierarchicalLabelValueFilter
          criterion={criterionState}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (
      options &&
      !criterionIsHierarchicalLabelValue(criterionState.value) &&
      !criterionIsNumberValue(criterionState.value) &&
      !criterionIsStashIDValue(criterionState.value) &&
      !criterionIsDateValue(criterionState.value) &&
      !criterionIsTimestampValue(criterionState.value)
    ) {
      if (!Array.isArray(criterionState.value)) {
        return (
          <OptionFilter
            criterion={criterionState}
            setCriterion={setCriterion}
          />
        );
      } else {
        return (
          <OptionListFilter
            criterion={criterionState}
            setCriterion={setCriterion}
          />
        );
      }
    }
    if (criterionState instanceof PathCriterion) {
      return (
        <PathFilter
          criterion={criterionState}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (criterionState instanceof DurationCriterion) {
      return (
        <DurationFilter
          criterion={criterionState}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (criterionState instanceof DateCriterion) {
      return (
        <DateFilter
          criterion={criterionState}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (criterionState instanceof TimestampCriterion) {
      return (
        <TimestampFilter
          criterion={criterionState}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (criterionState instanceof NumberCriterion) {
      return (
        <NumberFilter
          criterion={criterionState}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (criterionState instanceof RatingCriterion) {
      return (
        <RatingFilter
          criterion={criterionState}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (criterionState instanceof PhashCriterion) {
      return (
        <PhashFilter
          criterion={criterionState}
          onValueChanged={onValueChanged}
        />
      );
    }
    if (
      criterionState instanceof CountryCriterion &&
      (criterionState.modifier === CriterionModifier.Equals ||
        criterionState.modifier === CriterionModifier.NotEquals)
    ) {
      return (
        <CountrySelect
          value={criterionState.value}
          onChange={(v) => onValueChanged(v)}
          menuPortalTarget={document.body}
        />
      );
    }
    return (
      <InputFilter criterion={criterionState} onValueChanged={onValueChanged} />
    );
  }, [criterionState, setCriterion, options]);

  return (
    <div>
      {modifierSelector}
      {valueControl}
    </div>
  );
};

interface ICriterionEditor {
  emptyCriterion: Criterion<CriterionValue>;
  criterion?: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
}

export const CriterionEditor: React.FC<ICriterionEditor> = ({
  emptyCriterion,
  criterion,
  setCriterion,
}) => {
  const filterControl = useMemo(() => {
    if (
      emptyCriterion instanceof BooleanCriterion &&
      criterion instanceof BooleanCriterion
    ) {
      return (
        <BooleanFilter
          criterion={criterion ?? emptyCriterion}
          setCriterion={setCriterion}
        />
      );
    }

    return (
      <GenericCriterionEditor
        emptyCriterion={emptyCriterion}
        criterion={criterion}
        setCriterion={setCriterion}
      />
    );
  }, [emptyCriterion, criterion, setCriterion]);

  return <div className="criterion-editor">{filterControl}</div>;
};
