import React, { createContext, useContext, useState } from "react";
import { Button } from "react-bootstrap";
import { useIntl } from "react-intl";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCog,
  faCheck,
  faEye,
  faEyeSlash,
} from "@fortawesome/free-solid-svg-icons";
import {
  SidebarFilterDefinition,
  useSidebarFilters,
} from "src/extensions/hooks/useSidebarFilters";

// Context for edit mode state
interface ISidebarFilterEditContext {
  isEditMode: boolean;
  isFilterVisible: (filterId: string) => boolean;
  toggleFilterVisibility: (filterId: string) => void;
}

export const SidebarFilterEditContext = createContext<ISidebarFilterEditContext>({
  isEditMode: false,
  isFilterVisible: () => true,
  toggleFilterVisibility: () => {},
});

interface ISidebarFilterSelectorProps {
  viewName: string;
  filterDefinitions: SidebarFilterDefinition[];
  children: React.ReactNode;
  headerContent?: React.ReactNode;
  onEditModeChange?: (isEditMode: boolean) => void;
}

export const SidebarFilterSelector: React.FC<ISidebarFilterSelectorProps> = ({
  viewName,
  filterDefinitions,
  children,
  headerContent,
  onEditModeChange,
}) => {
  const intl = useIntl();
  const [isEditMode, setIsEditMode] = useState(false);

  const handleEditModeChange = (newEditMode: boolean) => {
    setIsEditMode(newEditMode);
    onEditModeChange?.(newEditMode);
  };

  const { isFilterVisible, toggleFilter } = useSidebarFilters(
    viewName,
    filterDefinitions
  );

  const toggleEditMode = () => handleEditModeChange(!isEditMode);

  const toggleButton = (
    <Button
      variant="link"
      className={`sidebar-filter-selector-button ${isEditMode ? "active" : ""}`}
      onClick={toggleEditMode}
      title={
        isEditMode
          ? intl.formatMessage({ id: "actions.done", defaultMessage: "Done" })
          : intl.formatMessage({
              id: "customize_filters",
              defaultMessage: "Customize Filters",
            })
      }
    >
      <FontAwesomeIcon icon={isEditMode ? faCheck : faCog} />
    </Button>
  );

  return (
    <SidebarFilterEditContext.Provider
      value={{
        isEditMode,
        isFilterVisible,
        toggleFilterVisibility: toggleFilter,
      }}
    >
      {headerContent && (
        <div className="sidebar-section-header">
          <span className="sidebar-section-title">{headerContent}</span>
          {toggleButton}
        </div>
      )}
      {children}
    </SidebarFilterEditContext.Provider>
  );
};

// Visibility toggle button component for individual filters
interface IFilterVisibilityToggleProps {
  filterId: string;
}

export const FilterVisibilityToggle: React.FC<IFilterVisibilityToggleProps> = ({
  filterId,
}) => {
  const intl = useIntl();
  const { isEditMode, isFilterVisible, toggleFilterVisibility } = useContext(
    SidebarFilterEditContext
  );

  if (!isEditMode) return null;

  const visible = isFilterVisible(filterId);

  return (
    <Button
      variant="link"
      className={`filter-visibility-toggle ${visible ? "visible" : "hidden"}`}
      onClick={(e) => {
        e.stopPropagation();
        toggleFilterVisibility(filterId);
      }}
      title={
        visible
          ? intl.formatMessage({ id: "actions.hide", defaultMessage: "Hide" })
          : intl.formatMessage({ id: "actions.show", defaultMessage: "Show" })
      }
    >
      <FontAwesomeIcon icon={visible ? faEye : faEyeSlash} />
    </Button>
  );
};

// Hook to check if a filter should be rendered
export const useFilterVisibility = (filterId: string) => {
  const { isEditMode, isFilterVisible } = useContext(SidebarFilterEditContext);

  // In edit mode, always show all filters
  // In normal mode, only show visible filters
  const shouldRender = isEditMode || isFilterVisible(filterId);
  const isVisible = isFilterVisible(filterId);

  return { shouldRender, isVisible, isEditMode };
};

// Wrapper component for filters that handles visibility
interface IFilterWrapperProps {
  filterId: string;
  children: React.ReactNode;
  alwaysShow?: boolean; // For filters that should ignore visibility settings (e.g., studio filter in studio view)
}

export const FilterWrapper: React.FC<IFilterWrapperProps> = ({
  filterId,
  children,
  alwaysShow = false,
}) => {
  const { shouldRender, isVisible, isEditMode } = useFilterVisibility(filterId);

  // If alwaysShow is false and we shouldn't render, return null
  if (!alwaysShow && !shouldRender) return null;

  // In edit mode, wrap with visibility indicator
  if (isEditMode) {
    return (
      <div className={`filter-wrapper ${!isVisible ? "filter-hidden" : ""}`}>
        <div className="filter-visibility-indicator">
          <FilterVisibilityToggle filterId={filterId} />
        </div>
        {children}
      </div>
    );
  }

  return <>{children}</>;
};
