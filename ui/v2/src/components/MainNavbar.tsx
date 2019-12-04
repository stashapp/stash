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

interface IProps {
  onMenuToggle() : void
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

            <NavLink
              exact={true}
              to="/scenes"
              className="bp3-button bp3-minimal bp3-icon-video collapsible-navlink"
              activeClassName="bp3-active"
            >
              Scenes
            </NavLink>

            <NavLink
              exact={true}
              to="/scenes/markers"
              className="bp3-button bp3-minimal bp3-icon-map-marker collapsible-navlink"
              activeClassName="bp3-active"
            >
              Markers
            </NavLink>

            <NavLink
              exact={true}
              to="/galleries"
              className="bp3-button bp3-minimal bp3-icon-media collapsible-navlink"
              activeClassName="bp3-active"
            >
              Galleries
            </NavLink>

            <NavLink
              exact={true}
              to="/performers"
              className="bp3-button bp3-minimal bp3-icon-person collapsible-navlink"
              activeClassName="bp3-active"
            >
              Performers
            </NavLink>

            <NavLink
              exact={true}
              to="/studios"
              className="bp3-button bp3-minimal bp3-icon-mobile-video collapsible-navlink"
              activeClassName="bp3-active"
            >
              Studios
            </NavLink>

            <NavLink
              exact={true}
              to="/tags"
              className="bp3-button bp3-minimal bp3-icon-tag collapsible-navlink"
              activeClassName="bp3-active"
            >
              Tags
            </NavLink>
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
