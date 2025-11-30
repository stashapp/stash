import React, { useEffect } from "react";
import { FormattedMessage } from "react-intl";
import { SidebarSection } from "src/components/Shared/MySidebar";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SearchTermInput } from "../ListFilter";
import { SidebarSavedFilterList } from "../SavedFilterList";
import { View } from "../views";
import useFocus from "src/utils/focus";
import ScreenUtils from "src/utils/screen";
import Mousetrap from "mousetrap";
import { Button } from "react-bootstrap";

export const FilteredSidebarHeader: React.FC<{
  sidebarOpen: boolean;
  showEditFilter: () => void;
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  view?: View;
  focus?: ReturnType<typeof useFocus>;
}> = ({
  sidebarOpen,
  showEditFilter,
  filter,
  setFilter,
  view,
  focus: providedFocus,
}) => {
  const localFocus = useFocus();
  const focus = providedFocus ?? localFocus;
  const [, setFocus] = focus;

  // Set the focus on the input field when the sidebar is opened
  // Don't do this on touch devices
  useEffect(() => {
    if (sidebarOpen && !ScreenUtils.isTouch()) {
      setFocus();
    }
  }, [sidebarOpen, setFocus]);

  return (
    <>
      <div className="sidebar-search-container">
        <SearchTermInput
          filter={filter}
          onFilterUpdate={setFilter}
          focus={focus}
        />
      </div>

      <div>
        <Button
          className="edit-filter-button"
          size="sm"
          onClick={() => showEditFilter()}
        >
          <FormattedMessage id="search_filter.edit_filter" />
        </Button>
      </div>

      <SidebarSection
        className="sidebar-saved-filters"
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
