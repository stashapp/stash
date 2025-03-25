import React, { PropsWithChildren } from "react";
import { CollapseButton } from "./CollapseButton";
import { useOnOutsideClick } from "src/hooks/OutsideClick";
import ScreenUtils, { useMediaQuery } from "src/utils/screen";

const fixedSidebarMediaQuery = "only screen and (max-width: 991px)";

// const CloseButton: React.FC<{
//   onClick: () => void;
// }> = ({ onClick }) => {
//   return (
//     <Button
//       variant="minimal"
//       size="lg"
//       className="close-button"
//       onClick={onClick}
//     >
//       <Icon icon={faTimes} />
//     </Button>
//   );
// };

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
      {/* {onHide && <CloseButton onClick={() => onHide()} />} */}
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
