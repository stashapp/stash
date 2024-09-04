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
import { ConfigurationContext } from "src/hooks/Config";
import { ListFilterModel } from "src/models/list-filter/filter";
import { getFilterOptions } from "src/models/list-filter/factory";
import { FilterTags } from "./FilterTags";
import { CriterionEditor } from "./CriterionEditor";
import { Icon } from "../Shared/Icon";
import {
  faChevronDown,
  faChevronRight,
  faTimes,
  faThumbtack,
} from "@fortawesome/free-solid-svg-icons";
import { useCompare, usePrevious } from "src/hooks/state";
import { CriterionType } from "src/models/list-filter/types";
import { useToast } from "src/hooks/Toast";
import { useConfigureUI } from "src/core/StashService";
import { FilterMode } from "src/core/generated-graphql";
import { useFocusOnce } from "src/utils/focus";
import Mousetrap from "mousetrap";
import ScreenUtils from "src/utils/screen";

interface ICriterionList {
  criteria: string[];
  currentCriterion?: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
  criterionOptions: CriterionOption[];
  pinnedCriterionOptions: CriterionOption[];
  selected?: CriterionOption;
  optionSelected: (o?: CriterionOption) => void;
  onRemoveCriterion: (c: string) => void;
  onTogglePin: (c: CriterionOption) => void;
}

const CriterionOptionList: React.FC<ICriterionList> = ({
  criteria,
  currentCriterion,
  setCriterion,
  criterionOptions,
  pinnedCriterionOptions,
  selected,
  optionSelected,
  onRemoveCriterion,
  onTogglePin,
}) => {
  const prevCriterion = usePrevious(currentCriterion);

  const scrolled = useRef(false);

  const type = currentCriterion?.criterionOption.type;
  const prevType = prevCriterion?.criterionOption.type;

  const criteriaRefs = useMemo(() => {
    const refs: Record<string, React.RefObject<HTMLDivElement>> = {};
    criterionOptions.forEach((c) => {
      refs[c.type] = React.createRef();
    });
    pinnedCriterionOptions.forEach((c) => {
      refs[c.type] = React.createRef();
    });
    return refs;
  }, [criterionOptions, pinnedCriterionOptions]);

  function onSelect(k: string | null) {
    if (!k) {
      optionSelected(undefined);
      return;
    }

    let option = criterionOptions.find((c) => c.type === k);
    if (!option) {
      option = pinnedCriterionOptions.find((c) => c.type === k);
    }

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

  function togglePin(ev: React.MouseEvent, c: CriterionOption) {
    // needed to prevent the nav item from being selected
    ev.stopPropagation();
    ev.preventDefault();
    onTogglePin(c);
  }

  function renderCard(c: CriterionOption, isPin: boolean) {
    return (
      <Card key={c.type} data-type={c.type} ref={criteriaRefs[c.type]!}>
        <Accordion.Toggle className="filter-item-header" eventKey={c.type}>
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
          <Button
            className="pin-criterion-button"
            variant="minimal"
            onClick={(e) => togglePin(e, c)}
          >
            <Icon icon={faThumbtack} className={isPin ? "" : "tilted"} />
          </Button>
        </Accordion.Toggle>
        <Accordion.Collapse eventKey={c.type}>
          {(type === c.type && currentCriterion) ||
          (prevType === c.type && prevCriterion) ? (
            <Card.Body>
              <CriterionEditor
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
    >
      {pinnedCriterionOptions.length !== 0 && (
        <>
          {pinnedCriterionOptions.map((c) => renderCard(c, true))}
          <div className="pinned-criterion-divider" />
        </>
      )}
      {criterionOptions.map((c) => renderCard(c, false))}
    </Accordion>
  );
};

const FilterModeToConfigKey = {
  [FilterMode.Galleries]: "galleries",
  [FilterMode.Images]: "images",
  [FilterMode.Movies]: "groups",
  [FilterMode.Groups]: "groups",
  [FilterMode.Performers]: "performers",
  [FilterMode.SceneMarkers]: "sceneMarkers",
  [FilterMode.Scenes]: "scenes",
  [FilterMode.Studios]: "studios",
  [FilterMode.Tags]: "tags",
};

function filterModeToConfigKey(filterMode: FilterMode) {
  return FilterModeToConfigKey[filterMode];
}

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
  const Toast = useToast();
  const intl = useIntl();

  const { configuration } = useContext(ConfigurationContext);

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

  const filterOptions = useMemo(() => {
    return getFilterOptions(currentFilter.mode);
  }, [currentFilter.mode]);

  const criterionOptions = useMemo(() => {
    return [...filterOptions.criterionOptions].sort((a, b) => {
      return intl
        .formatMessage({ id: a.messageID })
        .localeCompare(intl.formatMessage({ id: b.messageID }));
    });
  }, [intl, filterOptions.criterionOptions]);

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

  const ui = configuration?.ui ?? {};
  const [saveUI] = useConfigureUI();

  const filteredOptions = useMemo(() => {
    const trimmedSearch = searchValue.trim().toLowerCase();
    if (!trimmedSearch) {
      return criterionOptions;
    }

    return criterionOptions.filter((c) => {
      return intl
        .formatMessage({ id: c.messageID })
        .toLowerCase()
        .includes(trimmedSearch);
    });
  }, [intl, searchValue, criterionOptions]);

  const pinnedFilters = useMemo(
    () => ui.pinnedFilters?.[filterModeToConfigKey(currentFilter.mode)] ?? [],
    [currentFilter.mode, ui.pinnedFilters]
  );
  const pinnedElements = useMemo(
    () => filteredOptions.filter((c) => pinnedFilters.includes(c.messageID)),
    [pinnedFilters, filteredOptions]
  );
  const unpinnedElements = useMemo(
    () => filteredOptions.filter((c) => !pinnedFilters.includes(c.messageID)),
    [pinnedFilters, filteredOptions]
  );

  const editingCriterionChanged = useCompare(editingCriterion);

  useEffect(() => {
    if (editingCriterionChanged && editingCriterion) {
      const option = criterionOptions.find((c) => c.type === editingCriterion);
      if (option) {
        optionSelected(option);
      }
    }
  }, [
    editingCriterion,
    criterionOptions,
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

  async function updatePinnedFilters(filters: string[]) {
    const configKey = filterModeToConfigKey(currentFilter.mode);
    try {
      await saveUI({
        variables: {
          input: {
            ...configuration?.ui,
            pinnedFilters: {
              ...ui.pinnedFilters,
              [configKey]: filters,
            },
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onTogglePinFilter(f: CriterionOption) {
    try {
      const existing = pinnedFilters.find((name) => name === f.messageID);
      if (existing) {
        await updatePinnedFilters(
          pinnedFilters.filter((name) => name !== f.messageID)
        );
      } else {
        await updatePinnedFilters([...pinnedFilters, f.messageID]);
      }
    } catch (err) {
      Toast.error(err);
    }
  }

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
      <Modal show onHide={() => onCancel()} className="edit-filter-dialog">
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
              criterionOptions={unpinnedElements}
              pinnedCriterionOptions={pinnedElements}
              optionSelected={optionSelected}
              selected={criterion?.criterionOption}
              onRemoveCriterion={(c) => removeCriterionString(c)}
              onTogglePin={(c) => onTogglePinFilter(c)}
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
