import React from "react";
import { FormattedMessage } from "react-intl";
import { SidebarSection, SidebarToolbar } from "src/components/Shared/Sidebar";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
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

export interface ISidebarFilterProps {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}

export interface ISidebarContentProps {
  messageID: string;
  option: CriterionOption;
  component: React.FC<ISidebarFilterProps>;
}

export interface ISidebarFilterSectionProps {
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  section: ISidebarContentProps;
}

export const SidebarFilterSection: React.FC<ISidebarFilterSectionProps> = ({
  filter,
  setFilter,
  section,
}) => {
  const { messageID, option, component: Component } = section;
  return (
    <Component
      title={<FormattedMessage id={messageID} />}
      data-type={option.type}
      option={option}
      filter={filter}
      setFilter={setFilter}
    />
  );
};

export interface ISidebarFilterSectionsProps {
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sections: ISidebarContentProps[];
}

export const SidebarFilterSections: React.FC<ISidebarFilterSectionsProps> = ({
  filter,
  setFilter,
  sections,
}) => {
  return (
    <>
      {sections.map((section) => (
        <SidebarFilterSection
          key={section.messageID}
          filter={filter}
          setFilter={setFilter}
          section={section}
        />
      ))}
    </>
  );
};
