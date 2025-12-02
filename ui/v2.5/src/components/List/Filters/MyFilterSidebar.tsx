import React, { useEffect } from "react";
import { FormattedMessage } from "react-intl";
import { SidebarSection } from "src/components/Shared/Sidebar";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SearchTermInput } from "../ListFilter";
import { SidebarSavedFilterList } from "../SavedFilterList";
import { View } from "../views";
import useFocus from "src/utils/focus";
import ScreenUtils from "src/utils/screen";
import Mousetrap from "mousetrap";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { faBookmark } from "@fortawesome/free-solid-svg-icons";
import { FilterMode } from "src/core/generated-graphql";

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
  // TODO consider updating this logic to use the actions.search message along with coresponding messages for the content types
  const placeholderMap: Record<FilterMode, string> = {
    [FilterMode.Scenes]: "Search scenes...",
    [FilterMode.Performers]: "Search performers...",
    [FilterMode.Studios]: "Search studios...",
    [FilterMode.Galleries]: "Search galleries",
    [FilterMode.SceneMarkers]: "Search scene markers...",
    [FilterMode.Movies]: "Search movies...",
    [FilterMode.Groups]: "Search groups...",
    [FilterMode.Tags]: "Search tags...",
    [FilterMode.Images]: "Search images...",
  };

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
          placeholder={placeholderMap[filter.mode]}
        />
      </div>

      <div className="sidebar-saved-filters">
        <div className="sidebar-section-header">
          <Icon icon={faBookmark} />
          <FormattedMessage id="search_filter.saved_filters" />
        </div>
        <SidebarSavedFilterList
          filter={filter}
          onSetFilter={setFilter}
          view={view}
        />
      </div>
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
