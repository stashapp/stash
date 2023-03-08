import cloneDeep from "lodash-es/cloneDeep";
import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { Button, Modal, Nav } from "react-bootstrap";
import cx from "classnames";
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
import { Icon } from "../Shared/Icon";
import { faChevronLeft, faTimes } from "@fortawesome/free-solid-svg-icons";

interface ICriterionList {
  criteria: string[];
  criterionOptions: CriterionOption[];
  selected?: CriterionOption;
  optionSelected: (o: CriterionOption) => void;
  onRemoveCriterion: (c: string) => void;
}

const CriterionOptionList: React.FC<ICriterionList> = ({
  criteria,
  criterionOptions,
  selected,
  optionSelected,
  onRemoveCriterion,
}) => {
  function onSelect(k: string | null) {
    if (!k) return;

    const option = criterionOptions.find((c) => c.type === k);

    if (option) {
      optionSelected(option);
    }
  }

  function removeClicked(ev: React.MouseEvent, t: string) {
    // needed to prevent the nav item from being selected
    ev.stopPropagation();
    ev.preventDefault();
    onRemoveCriterion(t);
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
            {criteria.some((cc) => c.type === cc) && (
              <Button
                className="remove-criterion-button"
                variant="minimal"
                onClick={(e) => removeClicked(e, c.type)}
              >
                <Icon icon={faTimes} />
              </Button>
            )}
          </Nav.Link>
        </Nav.Item>
      ))}
    </Nav>
  );
};

interface ICriterionTitle {
  messageID: string;
  onBack: () => void;
}

const CriterionTitle: React.FC<ICriterionTitle> = ({ messageID, onBack }) => {
  return (
    <div className="criterion-title">
      <Button onClick={() => onBack()} variant="secondary">
        <Icon icon={faChevronLeft} />
      </Button>
      <FormattedMessage id={messageID} />
    </div>
  );
};

interface IEditFilterProps {
  filter: ListFilterModel;
  editingCriterion?: string;
  onApply: (filter: ListFilterModel) => void;
  onCancel: () => void;
}

export const EditFilterDialog: React.FC<IEditFilterProps> = ({
  filter,
  editingCriterion,
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

  const criteriaList = useMemo(() => {
    return criteria.map((c) => c.criterionOption.type);
  }, [criteria]);

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
    if (editingCriterion) {
      const option = criterionOptions.find((c) => c.type === editingCriterion);
      if (option) {
        optionSelected(option);
      }
    }
  }, [editingCriterion, criterionOptions, optionSelected]);

  function replaceCriterion(c: Criterion<CriterionValue>) {
    const newFilter = cloneDeep(currentFilter);

    if (!c.isValid()) {
      // remove from the filter if present
      const newCriteria = criteria.filter((cc) => {
        return cc.criterionOption.type !== c.criterionOption.type;
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

  function removeCriterionString(c: string) {
    const cc = criteria.find((ccc) => ccc.criterionOption.type === c);
    if (cc) {
      removeCriterion(cc);
    }
  }

  function onClearAll() {
    const newFilter = cloneDeep(currentFilter);
    newFilter.criteria = [];
    setCurrentFilter(newFilter);
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
          <div
            className={cx("dialog-content", {
              "criterion-selected": !!criterion,
            })}
          >
            <CriterionOptionList
              criteria={criteriaList}
              criterionOptions={criterionOptions}
              optionSelected={optionSelected}
              selected={criterion?.criterionOption}
              onRemoveCriterion={(c) => removeCriterionString(c)}
            />
            <div className="edit-filter-right">
              {criterion ? (
                <div>
                  <CriterionTitle
                    messageID={criterion.criterionOption.messageID}
                    onBack={() => setCriterion(undefined)}
                  />
                  <CriterionEditor
                    criterion={criterion}
                    setCriterion={replaceCriterion}
                  />
                </div>
              ) : undefined}
              <div>
                <FilterTags
                  criteria={criteria}
                  onEditCriterion={(c) => optionSelected(c.criterionOption)}
                  onRemoveCriterion={(c) => removeCriterion(c)}
                  onRemoveAll={() => onClearAll()}
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
