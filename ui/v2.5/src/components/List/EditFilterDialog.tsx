import cloneDeep from "lodash-es/cloneDeep";
import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { Button, Modal, Nav } from "react-bootstrap";
import {
  CriterionValue,
  Criterion,
  CriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { makeCriteria } from "src/models/list-filter/criteria/factory";
import { FormattedMessage, useIntl } from "react-intl";
import { ConfigurationContext } from "src/hooks/Config";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getFilterOptions } from "src/models/list-filter/factory";
import { FilterTags } from "./FilterTags";
import { CriterionEditor } from "./CriterionEditor";

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
  const [newCriterion, setNewCriterion] = useState(false);

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
        setNewCriterion(false);
      } else {
        const newCriterion = makeCriteria(config, option.type);
        setCriterion(newCriterion);
        setNewCriterion(true);
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
