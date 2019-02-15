import {
  Navbar,
  NavbarDivider,
  NavbarGroup,
  NavbarHeading,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { Link, NavLink } from "react-router-dom";

interface IMainNavbarProps {}

export const MainNavbar: FunctionComponent<IMainNavbarProps> = (props) => {
  let newButtonPath: string | undefined;
  let newButtonElement: JSX.Element | undefined;
  switch (window.location.pathname) {
    case "/performers": {
      newButtonPath = "/performers/new";
      break;
    }
    case "/studios": {
      newButtonPath = "/studios/new";
      break;
    }
  }
  if (!!newButtonPath) {
    newButtonElement = (
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
    <Navbar fixedToTop={true}>
      <div>
        <NavbarGroup align="left">
          <NavbarHeading><Link to="/" className="bp3-button bp3-minimal">Stash</Link></NavbarHeading>
          <NavbarDivider />

          <NavLink
            exact={true}
            to="/scenes"
            className="bp3-button bp3-minimal bp3-icon-video"
            activeClassName="bp3-active"
          >
            Scenes
          </NavLink>

          <NavLink
            exact={true}
            to="/scenes/markers"
            className="bp3-button bp3-minimal bp3-icon-map-marker"
            activeClassName="bp3-active"
          >
            Markers
          </NavLink>

          <NavLink
            exact={true}
            to="/galleries"
            className="bp3-button bp3-minimal bp3-icon-media"
            activeClassName="bp3-active"
          >
            Galleries
          </NavLink>

          <NavLink
            exact={true}
            to="/performers"
            className="bp3-button bp3-minimal bp3-icon-person"
            activeClassName="bp3-active"
          >
            Performers
          </NavLink>

          <NavLink
            exact={true}
            to="/studios"
            className="bp3-button bp3-minimal bp3-icon-mobile-video"
            activeClassName="bp3-active"
          >
            Studios
          </NavLink>

          <NavLink
            exact={true}
            to="/tags"
            className="bp3-button bp3-minimal bp3-icon-tag"
            activeClassName="bp3-active"
          >
            Tags
          </NavLink>
        </NavbarGroup>
        <NavbarGroup align="right">
          {newButtonElement}
          <NavLink
            exact={true}
            to="/settings"
            className="bp3-button bp3-minimal bp3-icon-cog"
            activeClassName="bp3-active"
          />
        </NavbarGroup>
      </div>
    </Navbar>
  );
};
