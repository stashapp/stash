import React, { PropsWithChildren, useState } from "react";
import { PaginationIndex } from "../List/Pagination";
import { ListFilterModel } from "src/models/list-filter/filter";
import cx from "classnames";
import { FilterSidebar } from "../List/FilterSidebar";
import { ListHeader } from "../List/ListHeader";
import { useModal } from "src/hooks/modal";
import { CollapseDivider } from "../Shared/CollapseDivider";
import { FilterTags } from "../List/FilterTags";
import { EditFilterDialog } from "../List/EditFilterDialog";
import {
  useFilterConfig,
  useListKeyboardShortcuts,
  useSaveLocalFilterState,
} from "../List/util";
import { useListSelect } from "src/hooks/listSelect";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { CriterionType } from "src/models/list-filter/types";
import { ScrollToTopButton, useScrollTop } from "../Shared/ScrollToTopButton";

type ListSelectProps = ReturnType<typeof useListSelect>;

export const ListPage: React.FC<
  PropsWithChildren<{
    id?: string;
    className?: string;
    filter: ListFilterModel;
    initialSidebarCollapsed?: boolean;
    setFilter: (filter: ListFilterModel) => void;
    listSelect: ListSelectProps;
    actionButtons?: React.ReactNode;
    selectedButtons?: (selectedIds: Set<string>) => React.ReactNode;
    metadataByline?: JSX.Element;
    totalCount: number;
    loading: boolean;
    pageView?: string;
  }>
> = ({
  id,
  className,
  filter,
  initialSidebarCollapsed,
  setFilter,
  listSelect,
  actionButtons,
  selectedButtons,
  metadataByline,
  totalCount,
  loading,
  pageView,
  children,
}) => {
  const { selectedIds, onSelectAll, onSelectNone } = listSelect;

  const [sidebarCollapsed, setSidebarCollapsed] = useState(
    initialSidebarCollapsed ?? false
  );

  // scroll to top button only applies when in small viewports where the header is not
  // sticky. useScrollTop returns the window scroll position
  const scrollTop = useScrollTop();

  const { criterionOptions, setCriterionOptions, sidebarOptions } =
    useFilterConfig(filter.mode);

  const { modal, showModal, closeModal } = useModal();

  useSaveLocalFilterState(pageView, filter, sidebarCollapsed);

  function editFilter(editType?: CriterionType) {
    showModal(
      <EditFilterDialog
        filter={filter}
        criterionOptions={criterionOptions}
        setCriterionOptions={(o) => setCriterionOptions(o)}
        onClose={(f) => {
          if (f) setFilter(f);
          closeModal();
        }}
        editingCriterion={editType}
      />
    );
  }

  useListKeyboardShortcuts({
    filter,
    setFilter,
    showEditFilter: editFilter,
    totalCount,
    toggleSidebarCollapsed: () => setSidebarCollapsed(!sidebarCollapsed),
    onSelectAll,
    onSelectNone,
  });

  return (
    <div id={id} className={cx("list-page", className)}>
      {modal}

      <div className={cx("sidebar-container", { collapsed: sidebarCollapsed })}>
        <FilterSidebar
          className={cx({ collapsed: sidebarCollapsed })}
          filter={filter}
          setFilter={(f) => setFilter(f)}
          criterionOptions={criterionOptions}
          setCriterionOptions={(o) => setCriterionOptions(o)}
          sidebarOptions={sidebarOptions}
          hideSidebar={() => setSidebarCollapsed(true)}
          totalCount={loading ? undefined : totalCount}
        />
      </div>
      <CollapseDivider
        collapsed={sidebarCollapsed}
        setCollapsed={(v) => setSidebarCollapsed(v)}
      />
      <div
        className={cx("list-page-results", {
          expanded: sidebarCollapsed,
          "items-selected": selectedIds.size > 0,
        })}
      >
        <ListHeader
          filter={filter}
          setFilter={setFilter}
          totalItems={totalCount}
          selectedIds={selectedIds}
          onSelectAll={onSelectAll}
          onSelectNone={onSelectNone}
          actionButtons={actionButtons}
          selectedButtons={selectedButtons}
          sidebarCollapsed={sidebarCollapsed}
          showFilterDialog={() => editFilter()}
          showSidebar={() => setSidebarCollapsed(false)}
        />
        <div>
          <FilterTags
            criteria={filter.criteria}
            onEditCriterion={(c) => editFilter(c.criterionOption.type)}
            onRemoveAll={() => setFilter(filter.clearCriteria())}
            onRemoveCriterion={(c) =>
              setFilter(filter.removeCriterion(c.criterionOption.type))
            }
          />
        </div>
        <div className="list-page-items">
          {loading ? (
            <LoadingIndicator />
          ) : (
            <>
              <PaginationIndex
                itemsPerPage={filter.itemsPerPage}
                currentPage={filter.currentPage}
                totalItems={totalCount}
                metadataByline={metadataByline}
              />
              <ScrollToTopButton
                scrollTop={scrollTop}
                onClick={() => {
                  window.scrollTo({ top: 0, behavior: "smooth" });
                }}
              />
              {children}
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default ListPage;
