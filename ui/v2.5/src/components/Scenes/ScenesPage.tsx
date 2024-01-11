import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { Pagination } from "../List/Pagination";
import { ListViewOptions, ZoomSelect } from "../List/ListViewOptions";
import { CriterionType, DisplayMode } from "src/models/list-filter/types";
import { PageSizeSelect, SearchField, SortBySelect } from "../List/ListFilter";
import { FilterMode, SortDirectionEnum } from "src/core/generated-graphql";
import { getFilterOptions } from "src/models/list-filter/factory";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useFindScenes } from "src/core/StashService";
import { SceneCardsGrid } from "./SceneCardsGrid";
import SceneQueue from "src/models/sceneQueue";
import { SceneListTable } from "./SceneListTable";
import { SceneWallPanel } from "../Wall/WallPanel";
import { Tagger } from "../Tagger/scenes/SceneTagger";
import { TaggerContext } from "../Tagger/context";
import useFocus from "src/utils/focus";
import {
  Criterion,
  CriterionOption,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { Button } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import { FormattedMessage, useIntl } from "react-intl";
import { faFilter, faTimes } from "@fortawesome/free-solid-svg-icons";
import { CriterionEditor } from "../List/CriterionEditor";
import { CollapseButton } from "../Shared/CollapseButton";
import cx from "classnames";
import { EditFilterDialog } from "../List/EditFilterDialog";
import { SavedFilterList } from "../List/SavedFilterList";

const FilterCriteriaList: React.FC<{
  filter: ListFilterModel;
  hiddenOptions: CriterionOption[];
  onRemoveCriterion: (c: Criterion<CriterionValue>) => void;
  onEditCriterion: (c: Criterion<CriterionValue>) => void;
}> = ({ filter, hiddenOptions, onRemoveCriterion, onEditCriterion }) => {
  const intl = useIntl();

  const criteria = useMemo(
    () =>
      filter.criteria.filter((c) => {
        return hiddenOptions.some((h) => h.type === c.criterionOption.type);
      }),
    [filter.criteria, hiddenOptions]
  );

  if (criteria.length === 0) return null;

  function onClickRemoveCriterion(
    criterion: Criterion<CriterionValue>,
    $event: React.MouseEvent<HTMLElement, MouseEvent>
  ) {
    if (!criterion) {
      return;
    }
    onRemoveCriterion(criterion);
    $event.stopPropagation();
  }

  function onClickCriterionTag(criterion: Criterion<CriterionValue>) {
    onEditCriterion(criterion);
  }

  return (
    <div className="filter-criteria-list">
      <ul>
        {criteria.map((c) => {
          return (
            <li className="filter-criteria-list-item" key={c.getId()}>
              <a onClick={() => onClickCriterionTag(c)}>
                <span>{c.getLabel(intl)}</span>
                <Button
                  className="remove-criterion-button"
                  variant="minimal"
                  onClick={($event) => onClickRemoveCriterion(c, $event)}
                >
                  <Icon icon={faTimes} />
                </Button>
              </a>
            </li>
          );
        })}
      </ul>
      <hr />
    </div>
  );
};

interface ICriterionList {
  filter: ListFilterModel;
  currentCriterion?: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
  criterionOptions: CriterionOption[];
  onRemoveCriterion: (c: string) => void;
  onOpenEditFilter: () => void;
}

const CriterionOptionList: React.FC<ICriterionList> = ({
  filter,
  currentCriterion,
  setCriterion,
  criterionOptions,
  onRemoveCriterion,
  onOpenEditFilter,
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
    return (
      <div>
        <CollapseButton
          text={intl.formatMessage({ id: c.messageID })}
          rightControls={
            <span>
              <Button
                className={cx("remove-criterion-button", {
                  invisible: !filter.criteria.some(
                    (cc) => c.type === cc.criterionOption.type
                  ),
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
      </div>
    );
  }

  return (
    <div className="criterion-list">
      {criterionOptions.map((c) => renderCard(c))}
      <div>
        <Button
          className="minimal edit-filter-button"
          onClick={() => onOpenEditFilter()}
        >
          <Icon icon={faFilter} />{" "}
          <FormattedMessage id="search_filter.edit_filter" />
        </Button>
      </div>
    </div>
  );
};

interface ICriterionOption {
  option: CriterionOption;
  showInSidebar: boolean;
}

const SceneFilter: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
}> = ({ filter, setFilter }) => {
  const intl = useIntl();

  const [queryRef, setQueryFocus] = useFocus();

  const getCriterionOptions = useCallback(() => {
    const options = getFilterOptions(filter.mode);

    return options.criterionOptions.map((o) => {
      return {
        option: o,
        showInSidebar: !options.defaultHiddenOptions.some(
          (c) => c.type === o.type
        ),
      } as ICriterionOption;
    });
  }, [filter.mode]);

  const [criterionOptions, setCriterionOptions] = useState(
    getCriterionOptions()
  );

  const sidebarOptions = useMemo(
    () => criterionOptions.filter((o) => o.showInSidebar).map((o) => o.option),
    [criterionOptions]
  );
  const hiddenOptions = useMemo(
    () => criterionOptions.filter((o) => !o.showInSidebar).map((o) => o.option),
    [criterionOptions]
  );

  const [criterion, setCriterion] = useState<Criterion<CriterionValue>>();

  const [editingCriterion, setEditingCriterion] = useState<string>();
  const [showEditFilter, setShowEditFilter] = useState(false);

  useEffect(() => {
    setCriterionOptions(getCriterionOptions());
  }, [getCriterionOptions]);

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

  return (
    <div className="scene-filter">
      <SearchField
        searchTerm={filter.searchTerm}
        setSearchTerm={searchQueryUpdated}
        queryRef={queryRef}
        setQueryFocus={setQueryFocus}
      />
      <hr />
      <div>
        <CollapseButton
          text={intl.formatMessage({ id: "search_filter.saved_filters" })}
        >
          <SavedFilterList filter={filter} onSetFilter={setFilter} />
        </CollapseButton>
      </div>
      <hr />
      <div>
        <FilterCriteriaList
          filter={filter}
          hiddenOptions={hiddenOptions}
          onRemoveCriterion={(c) =>
            removeCriterionString(c.criterionOption.type)
          }
          onEditCriterion={(c) => setEditingCriterion(c.criterionOption.type)}
        />
      </div>
      <div>
        <CriterionOptionList
          filter={filter}
          currentCriterion={criterion}
          setCriterion={replaceCriterion}
          criterionOptions={sidebarOptions}
          onRemoveCriterion={(c) => removeCriterionString(c)}
          onOpenEditFilter={() => setShowEditFilter(true)}
        />
      </div>
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

export const ListHeader: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  totalItems: number;
}> = ({ filter, setFilter, totalItems }) => {
  const filterOptions = getFilterOptions(filter.mode);

  function onChangeZoom(newZoomIndex: number) {
    const newFilter = filter.clone();
    newFilter.zoomIndex = newZoomIndex;
    setFilter(newFilter);
  }

  function onChangeDisplayMode(displayMode: DisplayMode) {
    const newFilter = filter.clone();
    newFilter.displayMode = displayMode;
    setFilter(newFilter);
  }

  function onChangePageSize(val: number) {
    const newFilter = filter.clone();
    newFilter.itemsPerPage = val;
    newFilter.currentPage = 1;
    setFilter(newFilter);
  }

  function onChangeSortDirection(dir: SortDirectionEnum) {
    const newFilter = filter.clone();
    newFilter.sortDirection = dir;
    setFilter(newFilter);
  }

  function onChangeSortBy(eventKey: string | null) {
    const newFilter = filter.clone();
    newFilter.sortBy = eventKey ?? undefined;
    newFilter.currentPage = 1;
    setFilter(newFilter);
  }

  function onReshuffleRandomSort() {
    const newFilter = filter.clone();
    newFilter.currentPage = 1;
    newFilter.randomSeed = -1;
    setFilter(newFilter);
  }

  const onChangePage = useCallback(
    (page: number) => {
      const newFilter = filter.clone();
      newFilter.currentPage = page;
      setFilter(newFilter);

      // if the current page has a detail-header, then
      // scroll up relative to that rather than 0, 0
      const detailHeader = document.querySelector(".detail-header");
      if (detailHeader) {
        window.scrollTo(0, detailHeader.scrollHeight - 50);
      } else {
        window.scrollTo(0, 0);
      }
    },
    [filter, setFilter]
  );

  return (
    <div className="list-header">
      <div>
        <PageSizeSelect
          pageSize={filter.itemsPerPage}
          setPageSize={onChangePageSize}
        />
        <Pagination
          currentPage={filter.currentPage}
          itemsPerPage={filter.itemsPerPage}
          totalItems={totalItems}
          onChangePage={onChangePage}
        />
      </div>
      <div>
        <SortBySelect
          sortBy={filter.sortBy}
          direction={filter.sortDirection}
          options={filterOptions.sortByOptions}
          setSortBy={onChangeSortBy}
          setDirection={onChangeSortDirection}
          onReshuffleRandomSort={onReshuffleRandomSort}
        />
        <div>
          <ZoomSelect
            minZoom={0}
            maxZoom={3}
            zoomIndex={filter.zoomIndex}
            onChangeZoom={onChangeZoom}
          />
        </div>
        <ListViewOptions
          displayMode={filter.displayMode}
          displayModeOptions={filterOptions.displayModeOptions}
          onSetDisplayMode={onChangeDisplayMode}
        />
      </div>
    </div>
  );
};

export const ScenesPage: React.FC = ({}) => {
  const [filter, setFilter] = useState<ListFilterModel>(
    () => new ListFilterModel(FilterMode.Scenes)
  );

  const result = useFindScenes(filter);
  const [selectedIds /* setSelectedIds */] = useState<Set<string>>(new Set());
  const totalCount = useMemo(
    () => result.data?.findScenes.count ?? 0,
    [result.data?.findScenes.count]
  );

  function renderScenes() {
    if (!result.data?.findScenes) return;

    const queue = SceneQueue.fromListFilterModel(filter);

    if (filter.displayMode === DisplayMode.Grid) {
      return (
        <SceneCardsGrid
          scenes={result.data.findScenes.scenes}
          queue={queue}
          zoomIndex={filter.zoomIndex}
          selectedIds={selectedIds}
          onSelectChange={() => {}}
        />
      );
    }
    if (filter.displayMode === DisplayMode.List) {
      return (
        <SceneListTable
          scenes={result.data.findScenes.scenes}
          queue={queue}
          selectedIds={selectedIds}
          onSelectChange={() => {}}
        />
      );
    }
    if (filter.displayMode === DisplayMode.Wall) {
      return (
        <SceneWallPanel
          scenes={result.data.findScenes.scenes}
          sceneQueue={queue}
        />
      );
    }
    if (filter.displayMode === DisplayMode.Tagger) {
      return (
        <TaggerContext>
          <Tagger scenes={result.data.findScenes.scenes} queue={queue} />
        </TaggerContext>
      );
    }
  }

  return (
    <div id="scenes-page">
      <SceneFilter filter={filter} setFilter={(f) => setFilter(f)} />
      <div className="scenes-page-results">
        <ListHeader
          filter={filter}
          setFilter={(f) => setFilter(f)}
          totalItems={totalCount}
        />
        <div className="scenes-page-items">{renderScenes()}</div>
      </div>
    </div>
  );
};

export default ScenesPage;
