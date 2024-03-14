import cloneDeep from "lodash-es/cloneDeep";
import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { Accordion, Button, Card, Form, Modal } from "react-bootstrap";
import cx from "classnames";
import {
  CriterionValue,
  Criterion,
  CriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { FormattedMessage, useIntl } from "react-intl";
import { ListFilterModel } from "src/models/list-filter/filter";
import { FilterTags } from "./FilterTags";
import { CriterionEditor } from "./CriterionEditor";
import { Icon } from "../Shared/Icon";
import {
  faChevronDown,
  faChevronRight,
  faTimes,
  faEye,
  faEyeSlash,
} from "@fortawesome/free-solid-svg-icons";
import { useCompare, usePrevious } from "src/hooks/state";
import { CriterionType } from "src/models/list-filter/types";
import { useFocusOnce } from "src/utils/focus";
import Mousetrap from "mousetrap";
import { useDragReorder } from "src/hooks/dragReorder";
import { ConfigurationContext } from "src/hooks/Config";

export interface ICriterionOption {
  option: CriterionOption;
  showInSidebar: boolean;
}
import ScreenUtils from "src/utils/screen";

interface ICriterionList {
  searchValue: string;
  criteria: string[];
  currentCriterion?: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
  criterionOptions: ICriterionOption[];
  setCriterionOptions?: (c: ICriterionOption[]) => void;
  selected?: CriterionOption;
  optionSelected: (o?: CriterionOption) => void;
  onRemoveCriterion: (c: string) => void;
}

const CriterionOptionList: React.FC<ICriterionList> = ({
  searchValue,
  criteria,
  currentCriterion,
  setCriterion,
  criterionOptions: icriterionOptions,
  setCriterionOptions,
  selected,
  optionSelected,
  onRemoveCriterion,
}) => {
  const intl = useIntl();
  const { configuration: config } = useContext(ConfigurationContext);

  const prevCriterion = usePrevious(currentCriterion);

  const scrolled = useRef(false);

  const type = currentCriterion?.criterionOption.type;
  const prevType = prevCriterion?.criterionOption.type;

  const { stageList, onDragStart, onDragOver, onDrop, onDragOverDefault } =
    useDragReorder(icriterionOptions, setCriterionOptions ?? (() => {}));

  const filteredList = useMemo(() => {
    const trimmedSearch = searchValue.trim().toLowerCase();
    const mapped = stageList.map((c, i) => ({ option: c, index: i }));
    if (!trimmedSearch) {
      return mapped;
    }

    return mapped.filter((c) => {
      return intl
        .formatMessage({ id: c.option.option.messageID })
        .toLowerCase()
        .includes(trimmedSearch);
    });
  }, [intl, searchValue, stageList]);

  const criterionOptions = useMemo(() => {
    return stageList.map((c) => c.option);
  }, [stageList]);

  const emptyCriteria = useMemo(() => {
    const ret: Record<string, Criterion<CriterionValue>> = {};
    criterionOptions.forEach((c) => {
      ret[c.type] = c.makeCriterion(config);
    });
    return ret;
  }, [config, criterionOptions]);

  const criteriaRefs = useMemo(() => {
    const refs: Record<string, React.RefObject<HTMLDivElement>> = {};
    criterionOptions.forEach((c) => {
      refs[c.type] = React.createRef();
    });
    return refs;
  }, [criterionOptions]);

  function onSelect(k: string | null) {
    if (!k) {
      optionSelected(undefined);
      return;
    }

    const option = criterionOptions.find((c) => c.type === k);
    if (option) {
      optionSelected(option);
    }
  }

  useEffect(() => {
    // scrolling to the current criterion doesn't work well when the
    // dialog is already open, so limit to when we click on the
    // criterion from the external tags
    if (!scrolled.current && type && criteriaRefs[type]?.current) {
      criteriaRefs[type].current!.scrollIntoView({
        behavior: "smooth",
        block: "start",
      });
      scrolled.current = true;
    }
  }, [currentCriterion, criteriaRefs, type]);

  function getReleventCriterion(t: CriterionType) {
    if (currentCriterion?.criterionOption.type === t) {
      return currentCriterion;
    }

    return prevCriterion;
  }

  function removeClicked(ev: React.MouseEvent, t: string) {
    // needed to prevent the nav item from being selected
    ev.stopPropagation();
    ev.preventDefault();
    onRemoveCriterion(t);
  }

  function toggleShowInSidebar(ev: React.MouseEvent, o: CriterionOption) {
    // needed to prevent the nav item from being selected
    ev.stopPropagation();
    ev.preventDefault();

    if (!setCriterionOptions) return;

    const newOptions = icriterionOptions.map((c) => {
      if (c.option.type === o.type) {
        return {
          ...c,
          showInSidebar: !c.showInSidebar,
        } as ICriterionOption;
      }
      return c;
    });

    setCriterionOptions(newOptions);
  }

  function renderCard(o: ICriterionOption, index: number) {
    const c = o.option;
    return (
      <Card key={c.type} data-type={c.type} ref={criteriaRefs[c.type]!}>
        <Accordion.Toggle
          className="filter-item-header"
          eventKey={c.type}
          draggable={!!setCriterionOptions}
          onDragStart={(e: React.DragEvent<HTMLElement>) =>
            onDragStart(e, index)
          }
          onDragEnter={(e: React.DragEvent<HTMLElement>) =>
            onDragOver(e, index)
          }
          onDrop={() => onDrop()}
        >
          <span className="mr-auto">
            <Icon
              className="collapse-icon fa-fw"
              icon={type === c.type ? faChevronDown : faChevronRight}
            />
            <FormattedMessage id={c.messageID} />
          </span>
          {criteria.some((cc) => c.type === cc) && (
            <Button
              className="remove-criterion-button"
              variant="minimal"
              onClick={(e) => removeClicked(e, c.type)}
            >
              <Icon icon={faTimes} />
            </Button>
          )}
          {!!setCriterionOptions && (
            <Button
              className="toggle-criterion-sidebar-button"
              variant="minimal"
              onClick={(e) => toggleShowInSidebar(e, c)}
            >
              <Icon icon={!o.showInSidebar ? faEye : faEyeSlash} />
            </Button>
          )}
        </Accordion.Toggle>
        <Accordion.Collapse eventKey={c.type}>
          {(type === c.type && currentCriterion) ||
          (prevType === c.type && prevCriterion) ? (
            <Card.Body>
              <CriterionEditor
                emptyCriterion={emptyCriteria[c.type]}
                criterion={getReleventCriterion(c.type)!}
                setCriterion={setCriterion}
              />
            </Card.Body>
          ) : (
            <Card.Body></Card.Body>
          )}
        </Accordion.Collapse>
      </Card>
    );
  }

  return (
    <Accordion
      className="criterion-list"
      activeKey={selected?.type}
      onSelect={onSelect}
      onDragOver={onDragOverDefault}
    >
      {filteredList.map((c) => renderCard(c.option, c.index))}
    </Accordion>
  );
};

interface IEditFilterProps {
  filter: ListFilterModel;
  criterionOptions: ICriterionOption[];
  setCriterionOptions?: (c: ICriterionOption[]) => void;
  editingCriterion?: string;
  onClose: (filter?: ListFilterModel) => void;
}

export const EditFilterDialog: React.FC<IEditFilterProps> = ({
  filter,
  criterionOptions: icriterionOptions,
  setCriterionOptions,
  editingCriterion,
  onClose,
}) => {
  const intl = useIntl();

  const [searchValue, setSearchValue] = useState("");
  const [currentFilter, setCurrentFilter] = useState<ListFilterModel>(
    cloneDeep(filter)
  );
  const [criterion, setCriterion] = useState<Criterion<CriterionValue>>();

  const [searchRef, setSearchFocus] = useFocusOnce(!ScreenUtils.isTouch());

  const { criteria } = currentFilter;

  const criteriaList = useMemo(() => {
    return criteria.map((c) => c.criterionOption.type);
  }, [criteria]);

  const optionSelected = useCallback(
    (option?: CriterionOption) => {
      if (!option) {
        setCriterion(undefined);
        return;
      }

      // find the existing criterion if present
      const existing = criteria.find(
        (c) => c.criterionOption.type === option.type
      );
      if (existing) {
        setCriterion(existing);
      } else {
        const newCriterion = filter.makeCriterion(option.type);
        setCriterion(newCriterion);
      }
    },
    [filter, criteria]
  );

  const editingCriterionChanged = useCompare(editingCriterion);

  useEffect(() => {
    if (editingCriterionChanged && editingCriterion) {
      const option = icriterionOptions.find(
        (c) => c.option.type === editingCriterion
      );
      if (option) {
        optionSelected(option.option);
      }
    }
  }, [
    editingCriterion,
    icriterionOptions,
    optionSelected,
    editingCriterionChanged,
  ]);

  useEffect(() => {
    Mousetrap.bind("/", (e) => {
      setSearchFocus();
      e.preventDefault();
    });

    return () => {
      Mousetrap.unbind("/");
    };
  });

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
      optionSelected(undefined);
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
      <Modal show onHide={() => onClose()} className="edit-filter-dialog">
        <Modal.Header>
          <div>
            <FormattedMessage id="search_filter.edit_filter" />
          </div>
          <Form.Control
            className="btn-secondary search-input"
            onChange={(e) => setSearchValue(e.target.value)}
            value={searchValue}
            placeholder={`${intl.formatMessage({ id: "actions.search" })}…`}
            ref={searchRef}
          />
        </Modal.Header>
        <Modal.Body>
          <div
            className={cx("dialog-content", {
              "criterion-selected": !!criterion,
            })}
          >
            <CriterionOptionList
              criteria={criteriaList}
              currentCriterion={criterion}
              setCriterion={replaceCriterion}
              searchValue={searchValue}
              criterionOptions={icriterionOptions}
              setCriterionOptions={setCriterionOptions}
              optionSelected={optionSelected}
              selected={criterion?.criterionOption}
              onRemoveCriterion={(c) => removeCriterionString(c)}
            />
            {criteria.length > 0 && (
              <div>
                <FilterTags
                  criteria={criteria}
                  onEditCriterion={(c) => optionSelected(c.criterionOption)}
                  onRemoveCriterion={(c) => removeCriterion(c)}
                  onRemoveAll={() => onClearAll()}
                />
              </div>
            )}
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => onClose()}>
            <FormattedMessage id="actions.cancel" />
          </Button>
          <Button onClick={() => onClose(currentFilter)}>
            <FormattedMessage id="actions.apply" />
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
};
