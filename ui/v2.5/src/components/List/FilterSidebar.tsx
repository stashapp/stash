import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { CriterionType } from "src/models/list-filter/types";
import { SearchField } from "../List/ListFilter";
import { ListFilterModel } from "src/models/list-filter/filter";
import useFocus from "src/utils/focus";
import {
  Criterion,
  CriterionOption,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { Button } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import { FormattedMessage, useIntl } from "react-intl";
import {
  faFilter,
  faFloppyDisk,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { CriterionEditor } from "../List/CriterionEditor";
import { CollapseButton } from "../Shared/CollapseButton";
import cx from "classnames";
import { EditFilterDialog } from "../List/EditFilterDialog";
import { SaveFilterDialog, SavedFilterList } from "../List/SavedFilterList";
import { ICriterionOption } from "./util";
import { useModal } from "src/hooks/modal";
import { mutateSaveFilter, useSetDefaultFilter } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";

interface ICriterionList {
  filter: ListFilterModel;
  currentCriterion?: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
  criterionOptions: CriterionOption[];
  onRemoveCriterion: (c: string) => void;
}

const CriterionOptionList: React.FC<ICriterionList> = ({
  filter,
  currentCriterion,
  setCriterion,
  criterionOptions,
  onRemoveCriterion,
}) => {
  const intl = useIntl();

  const scrolled = useRef(false);

  const type = currentCriterion?.criterionOption.type;

  const criteriaRefs = useMemo(() => {
    const refs: Record<string, React.RefObject<HTMLDivElement>> = {};
    criterionOptions.forEach((c) => {
      refs[c.type] = React.createRef();
    });
    return refs;
  }, [criterionOptions]);

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
    // find the existing criterion if present
    const existing = filter.criteria.find((c) => c.criterionOption.type === t);
    if (existing) {
      return existing;
    } else {
      const newCriterion = filter.makeCriterion(t);
      return newCriterion;
    }
  }

  function removeClicked(ev: React.MouseEvent, t: string) {
    // needed to prevent the nav item from being selected
    ev.stopPropagation();
    ev.preventDefault();
    onRemoveCriterion(t);
  }

  function renderCard(c: CriterionOption) {
    const isPopulated = filter.criteria.some(
      (cc) => c.type === cc.criterionOption.type
    );

    return (
      <CollapseButton
        key={c.type}
        className={cx({ populated: isPopulated })}
        text={intl.formatMessage({ id: c.messageID })}
        rightControls={
          <span>
            <Button
              className={cx("remove-criterion-button", {
                invisible: !isPopulated,
              })}
              variant="minimal"
              onClick={(e) => removeClicked(e, c.type)}
            >
              <Icon icon={faTimes} />
            </Button>
          </span>
        }
      >
        <CriterionEditor
          criterion={getReleventCriterion(c.type)!}
          setCriterion={setCriterion}
        />
      </CollapseButton>
    );
  }

  return (
    <div className="criterion-list">
      {criterionOptions.map((c) => renderCard(c))}
    </div>
  );
};

export const FilterSidebar: React.FC<{
  className?: string;
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  criterionOptions: ICriterionOption[];
  sidebarOptions: CriterionOption[];
  setCriterionOptions: (v: ICriterionOption[]) => void;
}> = ({
  className,
  filter,
  setFilter,
  criterionOptions,
  sidebarOptions,
  setCriterionOptions,
}) => {
  const intl = useIntl();
  const Toast = useToast();

  const { modal, showModal, closeModal } = useModal();
  const [queryRef, setQueryFocus] = useFocus();

  const [setDefaultFilter] = useSetDefaultFilter();

  const [criterion, setCriterion] = useState<Criterion<CriterionValue>>();

  const [editingCriterion, setEditingCriterion] = useState<string>();
  const [showEditFilter, setShowEditFilter] = useState(false);

  const { criteria } = filter;

  const searchQueryUpdated = useCallback(
    (value: string) => {
      const newFilter = filter.clone();
      newFilter.searchTerm = value;
      newFilter.currentPage = 1;
      setFilter(newFilter);
    },
    [filter, setFilter]
  );

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

  function removeCriterion(c: Criterion<CriterionValue>) {
    const newFilter = filter.clone();

    const newCriteria = criteria.filter((cc) => {
      return cc.getId() !== c.getId();
    });

    newFilter.criteria = newCriteria;

    setFilter(newFilter);
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

  function replaceCriterion(c: Criterion<CriterionValue>) {
    const newFilter = filter.clone();

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

    setFilter(newFilter);
  }

  function onApplyEditFilter(f?: ListFilterModel) {
    setShowEditFilter(false);
    setEditingCriterion(undefined);

    if (!f) return;
    setFilter(f);
  }

  async function onSaveFilterDialogClose(name?: string, id?: string) {
    try {
      if (!name) return;

      await mutateSaveFilter(name, id, filter);

      Toast.success(
        intl.formatMessage(
          {
            id: "toast.saved_entity",
          },
          {
            entity: intl.formatMessage({ id: "filter" }).toLocaleLowerCase(),
          }
        )
      );
    } catch (err) {
      Toast.error(err);
    } finally {
      closeModal();
    }
  }

  async function onSetDefaultFilter() {
    const filterCopy = filter.clone();

    try {
      await setDefaultFilter({
        variables: {
          input: {
            mode: filter.mode,
            find_filter: filterCopy.makeFindFilter(),
            object_filter: filterCopy.makeSavedFindFilter(),
            ui_options: filterCopy.makeUIOptions(),
          },
        },
      });

      Toast.success(
        intl.formatMessage({
          id: "toast.default_filter_set",
        })
      );
    } catch (err) {
      Toast.error(err);
    }
  }

  return (
    <div className={cx("filter-sidebar", className)}>
      <SearchField
        searchTerm={filter.searchTerm}
        setSearchTerm={searchQueryUpdated}
        queryRef={queryRef}
        setQueryFocus={setQueryFocus}
      />
      <div className="saved-filters">
        <CollapseButton
          text={intl.formatMessage({ id: "search_filter.saved_filters" })}
        >
          <SavedFilterList filter={filter} onSetFilter={setFilter} />
        </CollapseButton>
      </div>
      <CriterionOptionList
        filter={filter}
        currentCriterion={criterion}
        setCriterion={replaceCriterion}
        criterionOptions={sidebarOptions}
        onRemoveCriterion={(c) => removeCriterionString(c)}
      />
      <div>
        <Button
          variant="secondary"
          className="edit-filter-button"
          onClick={() => setShowEditFilter(true)}
        >
          <Icon icon={faFilter} />{" "}
          <FormattedMessage id="search_filter.edit_filter" />
        </Button>

        <Button
          variant="secondary"
          className="save-filter-button"
          onClick={() => {
            showModal(
              <SaveFilterDialog
                mode={filter.mode}
                onClose={onSaveFilterDialogClose}
                onSaveAsDefault={() => {
                  closeModal();
                  onSetDefaultFilter();
                }}
              />
            );
          }}
        >
          <Icon icon={faFloppyDisk} />{" "}
          <FormattedMessage id="actions.save_filter" />
        </Button>
      </div>
      {modal}
      {(showEditFilter || editingCriterion) && (
        <EditFilterDialog
          filter={filter}
          criterionOptions={criterionOptions}
          setCriterionOptions={(o) => setCriterionOptions(o)}
          onClose={onApplyEditFilter}
          editingCriterion={editingCriterion}
        />
      )}
    </div>
  );
};
