import React, {
  PropsWithChildren,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
import { CollapseButton } from "./CollapseButton";
import { useOnOutsideClick } from "src/hooks/OutsideClick";
import ScreenUtils, { useMediaQuery } from "src/utils/screen";
import { IViewConfig, useInterfaceLocalForage } from "src/hooks/LocalForage";
import { View } from "../List/views";
import cx from "classnames";
import { Button, CollapseProps } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "./Icon";
import { faSliders } from "@fortawesome/free-solid-svg-icons";

export type SidebarSectionStates = Record<string, boolean>;

const fixedSidebarMediaQuery = "only screen and (max-width: 1199px)";

export const Sidebar: React.FC<
  PropsWithChildren<{
    hide?: boolean;
    onHide?: () => void;
  }>
> = ({ hide, onHide, children }) => {
  const ref = React.useRef<HTMLDivElement>(null);

  const closeOnOutsideClick = useMediaQuery(fixedSidebarMediaQuery) && !hide;

  useOnOutsideClick(
    ref,
    !closeOnOutsideClick ? undefined : onHide,
    "ignore-sidebar-outside-click"
  );

  return (
    <div ref={ref} className="sidebar">
      {children}
    </div>
  );
};

// SidebarPane is a container for a Sidebar and content.
// It is expected that the children will be two elements:
// a Sidebar and a content element.
export const SidebarPane: React.FC<
  PropsWithChildren<{
    hideSidebar?: boolean;
  }>
> = ({ hideSidebar = false, children }) => {
  return (
    <div className={cx("sidebar-pane", { "hide-sidebar": hideSidebar })}>
      {children}
    </div>
  );
};

interface IContext {
  sectionOpen: SidebarSectionStates;
  setSectionOpen: (section: string, open: boolean) => void;
}

export const SidebarStateContext = React.createContext<IContext | null>(null);

export const SidebarSection: React.FC<
  PropsWithChildren<{
    text: React.ReactNode;
    className?: string;
    outsideCollapse?: React.ReactNode;
    // used to store open/closed state in SidebarStateContext
    sectionID?: string;
  }>
> = ({ className = "", text, outsideCollapse, sectionID = "", children }) => {
  // this is optional
  const contextState = React.useContext(SidebarStateContext);
  const openState =
    !contextState || !sectionID
      ? undefined
      : contextState.sectionOpen[sectionID] ?? undefined;

  function onOpenInternal(open: boolean) {
    if (contextState && sectionID) {
      contextState.setSectionOpen(sectionID, open);
    }
  }

  const collapseProps: Partial<CollapseProps> = {
    mountOnEnter: true,
    unmountOnExit: true,
  };
  return (
    <CollapseButton
      className={`sidebar-section ${className}`}
      collapseProps={collapseProps}
      text={text}
      outsideCollapse={outsideCollapse}
      onOpenChanged={onOpenInternal}
      open={openState}
    >
      {children}
    </CollapseButton>
  );
};

export const SidebarToggleButton: React.FC<{
  onClick: () => void;
}> = ({ onClick }) => {
  const intl = useIntl();
  return (
    <Button
      className="minimal sidebar-toggle-button ignore-sidebar-outside-click"
      variant="secondary"
      onClick={onClick}
      title={intl.formatMessage({ id: "actions.sidebar.toggle" })}
    >
      <Icon icon={faSliders} />
    </Button>
  );
};

// show sidebar by default if not on mobile
export function defaultShowSidebar() {
  return !ScreenUtils.matchesMediaQuery(fixedSidebarMediaQuery);
}

export function useSidebarState(view?: View) {
  const [interfaceLocalForage, setInterfaceLocalForage] =
    useInterfaceLocalForage();

  const { data: interfaceLocalForageData, loading } = interfaceLocalForage;

  const viewConfig: IViewConfig = useMemo(() => {
    return view ? interfaceLocalForageData?.viewConfig?.[view] || {} : {};
  }, [view, interfaceLocalForageData]);

  const [showSidebar, setShowSidebar] = useState<boolean>();
  const [sectionOpen, setSectionOpen] = useState<SidebarSectionStates>();

  // set initial state once loading is done
  useEffect(() => {
    if (showSidebar !== undefined) return;

    if (!view) {
      setShowSidebar(defaultShowSidebar());
      return;
    }

    if (loading) return;

    // only show sidebar by default on large screens
    setShowSidebar(!!viewConfig.showSidebar && defaultShowSidebar());
    setSectionOpen(viewConfig.sectionOpen || {});
  }, [
    view,
    loading,
    showSidebar,
    viewConfig.showSidebar,
    viewConfig.sectionOpen,
  ]);

  const onSetShowSidebar = useCallback(
    (show: boolean | ((prevState: boolean | undefined) => boolean)) => {
      const nv = typeof show === "function" ? show(showSidebar) : show;
      setShowSidebar(nv);
      if (view === undefined) return;

      setInterfaceLocalForage((prev) => ({
        ...prev,
        viewConfig: {
          ...prev.viewConfig,
          [view]: {
            ...viewConfig,
            showSidebar: nv,
          },
        },
      }));
    },
    [showSidebar, setInterfaceLocalForage, view, viewConfig]
  );

  const onSetSectionOpen = useCallback(
    (section: string, open: boolean) => {
      const newSectionOpen = { ...sectionOpen, [section]: open };
      setSectionOpen(newSectionOpen);
      if (view === undefined) return;

      setInterfaceLocalForage((prev) => ({
        ...prev,
        viewConfig: {
          ...prev.viewConfig,
          [view]: {
            ...viewConfig,
            sectionOpen: newSectionOpen,
          },
        },
      }));
    },
    [sectionOpen, setInterfaceLocalForage, view, viewConfig]
  );

  return {
    showSidebar: showSidebar ?? defaultShowSidebar(),
    sectionOpen: sectionOpen || {},
    setShowSidebar: onSetShowSidebar,
    setSectionOpen: onSetSectionOpen,
    loading: showSidebar === undefined,
  };
}
