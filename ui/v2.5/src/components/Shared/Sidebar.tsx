import React, { PropsWithChildren } from "react";
import { CollapseButton } from "./CollapseButton";

export const Sidebar: React.FC<PropsWithChildren<{
  hide?: boolean;
}>> = ({ hide, children }) => {
  const hideClass = hide ? "hide" : "";
  return <div className={`sidebar ${hideClass}`}>{children}</div>;
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
  }>
> = ({ text, children }) => {
  return (
    <CollapseButton className="sidebar-section" text={text}>
      {children}
    </CollapseButton>
  );
};
