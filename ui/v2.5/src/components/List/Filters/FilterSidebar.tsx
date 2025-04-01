import React from "react";
import { FormattedMessage } from "react-intl";
import { SidebarSection, SidebarToolbar } from "src/components/Shared/Sidebar";
import { ListFilterModel } from "src/models/list-filter/filter";
import { FilterButton } from "./FilterButton";
import { SearchTermInput } from "../ListFilter";
import { SidebarSavedFilterList } from "../SavedFilterList";
import { View } from "../views";

export const FilteredSidebarToolbar: React.FC<{
  onClose?: () => void;
  showEditFilter: () => void;
  filter: ListFilterModel;
}> = ({ onClose, showEditFilter, filter, children }) => {
  return (
    <SidebarToolbar onClose={onClose}>
      {children}
      <FilterButton onClick={() => showEditFilter()} filter={filter} />
    </SidebarToolbar>
  );
};

export const FilteredSidebarHeader: React.FC<{
  onClose?: () => void;
  showEditFilter: () => void;
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  view?: View;
}> = ({ onClose, showEditFilter, filter, setFilter, view }) => {
  return (
    <>
      <FilteredSidebarToolbar
        onClose={onClose}
        showEditFilter={showEditFilter}
        filter={filter}
      />
      <SearchTermInput filter={filter} onFilterUpdate={setFilter} />
      <SidebarSection
        text={<FormattedMessage id="search_filter.saved_filters" />}
      >
        <SidebarSavedFilterList
          filter={filter}
          onSetFilter={setFilter}
          view={view}
        />
      </SidebarSection>
    </>
  );
};
