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

const fixedSidebarMediaQuery = "only screen and (max-width: 991px)";

export const Sidebar: React.FC<
  PropsWithChildren<{
    hide?: boolean;
    onHide?: () => void;
  }>
> = ({ hide, onHide, children }) => {
  const hideClass = hide ? "hide" : "";
  const ref = React.useRef<HTMLDivElement>(null);

  const closeOnOutsideClick = useMediaQuery(fixedSidebarMediaQuery) && !hide;

  useOnOutsideClick(
    ref,
    !closeOnOutsideClick ? undefined : onHide,
    "ignore-sidebar-outside-click"
  );

  return (
    <div ref={ref} className={`sidebar ${hideClass}`}>
      {children}
    </div>
  );
};

// SidebarPane is a container for a Sidebar and content.
// It is expected that the children will be two elements:
// a Sidebar and a content element.
export const SidebarPane: React.FC<PropsWithChildren<{}>> = ({ children }) => {
  return <div className="sidebar-pane">{children}</div>;
};

export const SidebarSection: React.FC<
  PropsWithChildren<{
    text: React.ReactNode;
    className?: string;
    outsideCollapse?: React.ReactNode;
  }>
> = ({ className, text, outsideCollapse, children }) => {
  const collapseProps = {
    mountOnEnter: true,
    unmountOnExit: true,
  };
  return (
    <CollapseButton
      className={`sidebar-section ${className}`}
      collapseProps={collapseProps}
      text={text}
      outsideCollapse={outsideCollapse}
    >
      {children}
    </CollapseButton>
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

  // set initial state once loading is done
  useEffect(() => {
    if (showSidebar !== undefined) return;

    if (!view) {
      setShowSidebar(defaultShowSidebar());
      return;
    }

    if (loading) return;

    setShowSidebar(viewConfig.showSidebar ?? defaultShowSidebar());
  }, [view, loading, showSidebar, viewConfig.showSidebar]);

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

  return {
    showSidebar: showSidebar ?? defaultShowSidebar(),
    setShowSidebar: onSetShowSidebar,
    loading: showSidebar === undefined,
  };
}
