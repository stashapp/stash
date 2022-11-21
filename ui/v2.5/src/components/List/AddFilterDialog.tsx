import cloneDeep from "lodash-es/cloneDeep";
import React, { useContext, useEffect, useRef, useState } from "react";
import { Button, Form, Modal } from "react-bootstrap";
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
} from "src/models/list-filter/criteria/criterion";
import {
  NoneCriterion,
  NoneCriterionOption,
} from "src/models/list-filter/criteria/none";
import { makeCriteria } from "src/models/list-filter/criteria/factory";
import { ListFilterOptions } from "src/models/list-filter/filter-options";
import { FormattedMessage, useIntl } from "react-intl";
import {
  criterionIsHierarchicalLabelValue,
  criterionIsNumberValue,
  criterionIsStashIDValue,
  criterionIsDateValue,
  criterionIsTimestampValue,
  CriterionType,
} from "src/models/list-filter/types";
import { DurationFilter } from "./Filters/DurationFilter";
import { NumberFilter } from "./Filters/NumberFilter";
import { LabeledIdFilter } from "./Filters/LabeledIdFilter";
import { HierarchicalLabelValueFilter } from "./Filters/HierarchicalLabelValueFilter";
import { OptionsFilter } from "./Filters/OptionsFilter";
import { InputFilter } from "./Filters/InputFilter";
import { DateFilter } from "./Filters/DateFilter";
import { TimestampFilter } from "./Filters/TimestampFilter";
import { CountryCriterion } from "src/models/list-filter/criteria/country";
import { CountrySelect } from "../Shared";
import { StashIDCriterion } from "src/models/list-filter/criteria/stash-ids";
import { StashIDFilter } from "./Filters/StashIDFilter";
import { ConfigurationContext } from "src/hooks/Config";
import { RatingCriterion } from "../../models/list-filter/criteria/rating";
import { RatingFilter } from "./Filters/RatingFilter";

interface IAddFilterProps {
  onAddCriterion: (
    criterion: Criterion<CriterionValue>,
    oldId?: string
  ) => void;
  onCancel: () => void;
  filterOptions: ListFilterOptions;
  editingCriterion?: Criterion<CriterionValue>;
  existingCriterions: Criterion<CriterionValue>[];
}

export const AddFilterDialog: React.FC<IAddFilterProps> = ({
  onAddCriterion,
  onCancel,
  filterOptions,
  editingCriterion,
  existingCriterions,
}) => {
  const defaultValue = useRef<string | number | undefined>();

  const [criterion, setCriterion] = useState<Criterion<CriterionValue>>(
    new NoneCriterion()
  );
  const { options, modifierOptions } = criterion.criterionOption;

  const valueStage = useRef<CriterionValue>(criterion.value);
  const { configuration: config } = useContext(ConfigurationContext);

  const intl = useIntl();

  // Configure if we are editing an existing criterion
  useEffect(() => {
    if (!editingCriterion) {
      setCriterion(makeCriteria(config));
    } else {
      setCriterion(editingCriterion);
    }
  }, [config, editingCriterion]);

  useEffect(() => {
    valueStage.current = criterion.value;
  }, [criterion]);

  function onChangedCriteriaType(event: React.ChangeEvent<HTMLSelectElement>) {
    const newCriterionType = event.target.value as CriterionType;
    const newCriterion = makeCriteria(config, newCriterionType);
    setCriterion(newCriterion);
  }

  function onChangedModifierSelect(
    event: React.ChangeEvent<HTMLSelectElement>
  ) {
    const newCriterion = cloneDeep(criterion);
    newCriterion.modifier = event.target.value as CriterionModifier;
    setCriterion(newCriterion);
  }

  function onValueChanged(value: CriterionValue) {
    const newCriterion = cloneDeep(criterion);
    newCriterion.value = value;
    setCriterion(newCriterion);
  }

  function onAddFilter() {
    const oldId = editingCriterion ? editingCriterion.getId() : undefined;
    onAddCriterion(criterion, oldId);
  }

  const maybeRenderFilterPopoverContents = () => {
    if (criterion.criterionOption.type === "none") {
      return;
    }

    function renderModifier() {
      if (modifierOptions.length === 0) {
        return;
      }
      return (
        <Form.Control
          as="select"
          onChange={onChangedModifierSelect}
          value={criterion.modifier}
          className="btn-secondary"
        >
          {modifierOptions.map((c) => (
            <option key={c.value} value={c.value}>
              {c.label ? intl.formatMessage({ id: c.label }) : ""}
            </option>
          ))}
        </Form.Control>
      );
    }

    function renderSelect() {
      // always show stashID filter
      if (criterion instanceof StashIDCriterion) {
        return (
          <StashIDFilter
            criterion={criterion}
            onValueChanged={onValueChanged}
          />
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
        defaultValue.current = criterion.value;
        return (
          <OptionsFilter
            criterion={criterion}
            onValueChanged={onValueChanged}
          />
        );
      }
      if (criterion instanceof DurationCriterion) {
        return (
          <DurationFilter
            criterion={criterion}
            onValueChanged={onValueChanged}
          />
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
    }
    return (
      <>
        <Form.Group>{renderModifier()}</Form.Group>
        {renderSelect()}
      </>
    );
  };

  function maybeRenderFilterCriterion() {
    if (!editingCriterion) {
      return;
    }

    return (
      <Form.Group>
        <strong>
          {intl.formatMessage({
            id: editingCriterion.criterionOption.messageID,
          })}
        </strong>
      </Form.Group>
    );
  }

  function maybeRenderFilterSelect() {
    if (editingCriterion) {
      return;
    }

    const thisOptions = [NoneCriterionOption]
      .concat(filterOptions.criterionOptions)
      .filter(
        (c) =>
          !existingCriterions.find((ec) => ec.criterionOption.type === c.type)
      )
      .map((c) => {
        return {
          value: c.type,
          text: intl.formatMessage({ id: c.messageID }),
        };
      })
      .sort((a, b) => {
        if (a.value === "none") return -1;
        if (b.value === "none") return 1;
        return a.text.localeCompare(b.text);
      });

    return (
      <Form.Group controlId="filter">
        <Form.Label>
          <FormattedMessage id="search_filter.name" />
        </Form.Label>
        <Form.Control
          as="select"
          onChange={onChangedCriteriaType}
          value={criterion.criterionOption.type}
          className="btn-secondary"
        >
          {thisOptions.map((c) => (
            <option key={c.value} value={c.value} disabled={c.value === "none"}>
              {c.text}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    );
  }

  function isValid() {
    if (criterion.criterionOption.type === "none") {
      return false;
    }

    if (criterion instanceof RatingCriterion) {
      switch (criterion.modifier) {
        case CriterionModifier.Equals:
        case CriterionModifier.NotEquals:
        case CriterionModifier.LessThan:
          return !!criterion.value.value;
        case CriterionModifier.Between:
        case CriterionModifier.NotBetween:
          return criterion.value.value < (criterion.value.value2 ?? 0);
      }
    }

    return true;
  }

  const title = !editingCriterion
    ? intl.formatMessage({ id: "search_filter.add_filter" })
    : intl.formatMessage({ id: "search_filter.update_filter" });
  return (
    <>
      <Modal show onHide={() => onCancel()} className="add-filter-dialog">
        <Modal.Header>{title}</Modal.Header>
        <Modal.Body>
          <div className="dialog-content">
            {maybeRenderFilterSelect()}
            {maybeRenderFilterCriterion()}
            {maybeRenderFilterPopoverContents()}
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button onClick={onAddFilter} disabled={!isValid()}>
            {title}
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};
