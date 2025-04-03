import React, { useEffect } from "react";
import { FormattedMessage } from "react-intl";
import { SidebarSection, SidebarToolbar } from "src/components/Shared/Sidebar";
import { ListFilterModel } from "src/models/list-filter/filter";
import { FilterButton } from "./FilterButton";
import { SearchTermInput } from "../ListFilter";
import { SidebarSavedFilterList } from "../SavedFilterList";
import { View } from "../views";
import useFocus from "src/utils/focus";
import ScreenUtils from "src/utils/screen";
import Mousetrap from "mousetrap";

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
  sidebarOpen: boolean;
  onClose?: () => void;
  showEditFilter: () => void;
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  view?: View;
}> = ({ sidebarOpen, onClose, showEditFilter, filter, setFilter, view }) => {
  const focus = useFocus();
  const [, setFocus] = focus;

  // Set the focus on the input field when the sidebar is opened
  // Don't do this on mobile devices
  useEffect(() => {
    if (sidebarOpen && !ScreenUtils.isMobile()) {
      setFocus();
    }
  }, [sidebarOpen, setFocus]);

  return (
    <>
      <FilteredSidebarToolbar
        onClose={onClose}
        showEditFilter={showEditFilter}
        filter={filter}
      />
      <SearchTermInput
        filter={filter}
        onFilterUpdate={setFilter}
        focus={focus}
      />
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

export function useFilteredSidebarKeybinds(props: {
  showSidebar: boolean;
  setShowSidebar: (show: boolean) => void;
}) {
  const { showSidebar, setShowSidebar } = props;

  // Show the sidebar when the user presses the "/" key
  useEffect(() => {
    Mousetrap.bind("/", (e) => {
      if (!showSidebar) {
        setShowSidebar(true);
        e.preventDefault();
      }
    });

    return () => {
      Mousetrap.unbind("/");
    };
  }, [showSidebar, setShowSidebar]);

  // Hide the sidebar when the user presses the "Esc" key
  useEffect(() => {
    Mousetrap.bind("esc", (e) => {
      if (showSidebar) {
        setShowSidebar(false);
        e.preventDefault();
      }
    });

    return () => {
      Mousetrap.unbind("esc");
    };
  }, [showSidebar, setShowSidebar]);
}
