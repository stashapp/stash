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
import { Button, ButtonToolbar } from "react-bootstrap";
import { useIntl } from "react-intl";

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

export const SidebarSection: React.FC<
  PropsWithChildren<{
    text: React.ReactNode;
    className?: string;
    outsideCollapse?: React.ReactNode;
  }>
> = ({ className = "", text, outsideCollapse, children }) => {
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

export const SidebarIcon: React.FC = () => (
  <>
    {/* From: https://iconduck.com/icons/19707/sidebar
MIT License
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE. */}
    <svg
      className="svg-inline--fa fa-icon"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="3"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
      <line x1="9" y1="3" x2="9" y2="21" />
    </svg>
  </>
);

export const SidebarToolbar: React.FC<{
  onClose?: () => void;
}> = ({ onClose, children }) => {
  const intl = useIntl();

  return (
    <ButtonToolbar className="sidebar-toolbar">
      {onClose ? (
        <Button
          onClick={onClose}
          className="sidebar-close-button"
          variant="secondary"
          title={intl.formatMessage({ id: "actions.sidebar.close" })}
        >
          <SidebarIcon />
        </Button>
      ) : null}
      {children}
    </ButtonToolbar>
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

    // only show sidebar by default on large screens
    setShowSidebar(!!viewConfig.showSidebar && defaultShowSidebar());
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
