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
import { useFilterConfig } from "../List/util";
import { useListSelect } from "src/hooks/listSelect";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { CriterionType } from "src/models/list-filter/types";

type ListSelectProps = ReturnType<typeof useListSelect>;

export const ListPage: React.FC<
  PropsWithChildren<{
    id?: string;
    className?: string;
    filter: ListFilterModel;
    setFilter: (filter: ListFilterModel) => void;
    listSelect: ListSelectProps;
    actionButtons?: React.ReactNode;
    selectedButtons?: (selectedIds: Set<string>) => React.ReactNode;
    metadataByline?: JSX.Element;
    totalCount: number;
    loading: boolean;
  }>
> = ({
  id,
  className,
  filter,
  setFilter,
  listSelect,
  actionButtons,
  selectedButtons,
  metadataByline,
  totalCount,
  loading,
  children,
}) => {
  const { selectedIds, onSelectAll, onSelectNone } = listSelect;

  const [filterCollapsed, setFilterCollapsed] = useState(false);

  const { criterionOptions, setCriterionOptions, sidebarOptions } =
    useFilterConfig(filter.mode);

  const { modal, showModal, closeModal } = useModal();

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

  return (
    <div id={id} className={cx("list-page", className)}>
      {modal}

      <div className={cx("sidebar-container", { collapsed: filterCollapsed })}>
        <FilterSidebar
          className={cx({ collapsed: filterCollapsed })}
          filter={filter}
          setFilter={(f) => setFilter(f)}
          criterionOptions={criterionOptions}
          setCriterionOptions={(o) => setCriterionOptions(o)}
          sidebarOptions={sidebarOptions}
        />
      </div>
      <CollapseDivider
        collapsed={filterCollapsed}
        setCollapsed={(v) => setFilterCollapsed(v)}
      />
      <div className={cx("list-page-results", { expanded: filterCollapsed })}>
        <ListHeader
          filter={filter}
          setFilter={setFilter}
          totalItems={totalCount}
          selectedIds={selectedIds}
          onSelectAll={onSelectAll}
          onSelectNone={onSelectNone}
          actionButtons={actionButtons}
          selectedButtons={selectedButtons}
          sidebarCollapsed={filterCollapsed}
          showFilterDialog={() => editFilter()}
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
              {children}
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default ListPage;
