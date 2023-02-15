import cloneDeep from "lodash-es/cloneDeep";
import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { Button, Form, Modal, Nav } from "react-bootstrap";
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
  CriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { makeCriteria } from "src/models/list-filter/criteria/factory";
import { FormattedMessage, useIntl } from "react-intl";
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
import { ConfigurationContext } from "src/hooks/Config";
import { RatingCriterion } from "../../models/list-filter/criteria/rating";
import { RatingFilter } from "./Filters/RatingFilter";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getFilterOptions } from "src/models/list-filter/factory";
import { FilterTags } from "./FilterTags";

interface ICriterionList {
  criterionOptions: CriterionOption[];
  selected?: CriterionOption;
  optionSelected: (o: CriterionOption) => void;
}

const CriterionOptionList: React.FC<ICriterionList> = ({
  criterionOptions,
  selected,
  optionSelected,
}) => {
  function onSelect(k: string | null) {
    if (!k) return;

    const option = criterionOptions.find((c) => c.type === k);

    if (option) {
      optionSelected(option);
    }
  }

  return (
    <Nav
      variant="pills"
      className="criterion-list"
      activeKey={selected?.type}
      onSelect={onSelect}
    >
      {criterionOptions.map((c) => (
        <Nav.Item key={c.type}>
          <Nav.Link eventKey={c.type}>
            <FormattedMessage id={c.messageID} />
          </Nav.Link>
        </Nav.Item>
      ))}
    </Nav>
  );
};

interface IBooleanFilter {
  criterion: BooleanCriterion;
  setCriterion: (c: BooleanCriterion) => void;
}

const BooleanFilter: React.FC<IBooleanFilter> = ({
  criterion,
  setCriterion,
}) => {
  function onSelect(v: boolean) {
    const c = cloneDeep(criterion);
    if ((v && c.value === "true") || (!v && c.value === "false")) {
      c.value = "";
    } else {
      c.value = v ? "true" : "false";
    }

    setCriterion(c);
  }

  return (
    <div className="boolean-filter">
      <Form.Check
        id={`${criterion.getId()}-true`}
        onChange={() => onSelect(true)}
        checked={criterion.value === "true"}
        type="checkbox"
        label={<FormattedMessage id="true" />}
      />
      <Form.Check
        id={`${criterion.getId()}-false`}
        onChange={() => onSelect(false)}
        checked={criterion.value === "false"}
        type="checkbox"
        label={<FormattedMessage id="false" />}
      />
    </div>
  );
};

interface IOptionsListFilter {
  criterion: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
}

const OptionsListFilter: React.FC<IOptionsListFilter> = ({
  criterion,
  setCriterion,
}) => {
  function onSelect(v: string) {
    const c = cloneDeep(criterion);
    if (c.value === v) {
      c.value = "";
    } else {
      c.value = v;
    }

    setCriterion(c);
  }

  const { options } = criterion.criterionOption;

  return (
    <div className="option-list-filter">
      {options?.map((o) => (
        <Form.Check
          id={`${criterion.getId()}-${o.toString()}`}
          key={o.toString()}
          onChange={() => onSelect(o.toString())}
          checked={criterion.value === o.toString()}
          type="checkbox"
          label={o.toString()}
        />
      ))}
    </div>
  );
};

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
        className="btn-secondary"
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

const CriterionEditor: React.FC<ICriterionEditor> = ({
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

interface IEditFilterProps {
  filter: ListFilterModel;
  onApply: (filter: ListFilterModel) => void;
  onCancel: () => void;
}

export const EditFilterDialog: React.FC<IEditFilterProps> = ({
  filter,
  onApply,
  onCancel,
}) => {
  const intl = useIntl();

  const { configuration: config } = useContext(ConfigurationContext);

  const [currentFilter, setCurrentFilter] = useState<ListFilterModel>(
    cloneDeep(filter)
  );
  const [criterion, setCriterion] = useState<Criterion<CriterionValue>>();

  const { criteria } = currentFilter;

  const filterOptions = useMemo(() => {
    return getFilterOptions(currentFilter.mode);
  }, [currentFilter.mode]);

  const criterionOptions = useMemo(() => {
    const filteredOptions = filterOptions.criterionOptions.filter((o) => {
      return o.type !== "none";
    });

    filteredOptions.sort((a, b) => {
      return intl
        .formatMessage({ id: a.messageID })
        .localeCompare(intl.formatMessage({ id: b.messageID }));
    });

    return filteredOptions;
  }, [intl, filterOptions.criterionOptions]);

  const optionSelected = useCallback(
    (option: CriterionOption) => {
      // find the existing criterion if present
      const existing = criteria.find(
        (c) => c.criterionOption.type === option.type
      );
      if (existing) {
        setCriterion(existing);
      } else {
        const newCriterion = makeCriteria(config, option.type);
        setCriterion(newCriterion);
      }
    },
    [criteria, config]
  );

  useEffect(() => {
    if (!criterion && criterionOptions.length > 0) {
      const option = criterionOptions[0];
      optionSelected(option);
    }
  }, [criterion, criterionOptions, optionSelected]);

  function replaceCriterion(c: Criterion<CriterionValue>) {
    const newFilter = cloneDeep(currentFilter);

    if (!c.isValid()) {
      // remove from the filter if present
      const newCriteria = criteria.filter((cc) => {
        return cc.criterionOption.type === c.criterionOption.type;
      });

      newFilter.criteria = newCriteria;
    } else {
      let found = false;

      const newCriteria = criteria.map((cc) => {
        if (cc.criterionOption.type === c.criterionOption.type) {
          found = true;
          return c;
        }

        return cc;
      });

      if (!found) {
        newCriteria.push(c);
      }

      newFilter.criteria = newCriteria;
    }

    setCurrentFilter(newFilter);
    setCriterion(c);
  }

  function removeCriterion(c: Criterion<CriterionValue>) {
    const newFilter = cloneDeep(currentFilter);

    const newCriteria = criteria.filter((cc) => {
      return cc.getId() !== c.getId();
    });

    newFilter.criteria = newCriteria;

    setCurrentFilter(newFilter);
    if (criterion?.getId() === c.getId()) {
      optionSelected(c.criterionOption);
    }
  }

  return (
    <>
      <Modal
        show
        size="lg"
        onHide={() => onCancel()}
        className="edit-filter-dialog"
      >
        <Modal.Header>
          <FormattedMessage id="search_filter.edit_filter" />
        </Modal.Header>
        <Modal.Body>
          <div className="dialog-content">
            <CriterionOptionList
              criterionOptions={criterionOptions}
              optionSelected={optionSelected}
              selected={criterion?.criterionOption}
            />
            <div className="edit-filter-right">
              {criterion ? (
                <CriterionEditor
                  criterion={criterion}
                  setCriterion={replaceCriterion}
                />
              ) : undefined}
              <div>
                <FilterTags
                  criteria={criteria}
                  onEditCriterion={(c) => optionSelected(c.criterionOption)}
                  onRemoveCriterion={(c) => removeCriterion(c)}
                />
              </div>
            </div>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => onCancel()}>
            <FormattedMessage id="actions.cancel" />
          </Button>
          <Button onClick={() => onApply(currentFilter)}>
            <FormattedMessage id="actions.apply" />
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};
