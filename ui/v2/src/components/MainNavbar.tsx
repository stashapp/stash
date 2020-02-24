import {
  Navbar,
  NavbarDivider,
  NavbarGroup,
  NavbarHeading,
  Button,
} from "@blueprintjs/core";
import React, { FunctionComponent, useEffect, useState } from "react";
import { Link, NavLink } from "react-router-dom";
import useLocation from "react-use/lib/useLocation";
import { IMenuItem } from "../App";

interface IProps {
  onMenuToggle() : void
  menuItems: IMenuItem[]
}

export const MainNavbar: FunctionComponent<IProps> = (props) => {
  const [newButtonPath, setNewButtonPath] = useState<string | undefined>(undefined);
  const locationState = useLocation();

  useEffect(() => {
    switch (window.location.pathname) {
      case "/performers": {
        setNewButtonPath("/performers/new");
        break;
      }
      case "/studios": {
        setNewButtonPath("/studios/new");
        break;
      }
      default: {
        setNewButtonPath(undefined);
      }
    }
  }, [locationState.pathname]);

  function renderNewButton() {
    if (!newButtonPath) { return; }
    return (
      <>
        <NavLink
          to={newButtonPath}
          className="bp3-button bp3-intent-primary"
        >
          New
        </NavLink>
        <NavbarDivider />
      </>
    );
  }

  return (
    <>
      <Navbar fixedToTop={true}>
        <div>
          <NavbarGroup align="left">
            <Button className="menu-button" icon="menu" onClick={() => props.onMenuToggle()}/>
            <NavbarHeading><Link to="/" className="bp3-button bp3-minimal">Stash</Link></NavbarHeading>
            <NavbarDivider />

            {props.menuItems.map((i) => {
              return (
                <NavLink
                  exact={true}
                  to={i.href}
                  className={"bp3-button bp3-minimal collapsible-navlink bp3-icon-" + i.icon}
                  activeClassName="bp3-active"
                >
                  {i.text}
                </NavLink>
              );
            })}
          </NavbarGroup>
          <NavbarGroup align="right">
            {renderNewButton()}
            <NavLink
              exact={true}
              to="/settings"
              className="bp3-button bp3-minimal bp3-icon-cog"
              activeClassName="bp3-active"
            />
          </NavbarGroup>
        </div>
      </Navbar>
    </>
  );
};
