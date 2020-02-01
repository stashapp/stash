import React from "react";
import { Nav, Navbar, Button } from "react-bootstrap";
import { IconName } from "@fortawesome/fontawesome-svg-core";
import { LinkContainer } from "react-router-bootstrap";
import { Link, useLocation } from "react-router-dom";

import { Icon } from "src/components/Shared";

interface IMenuItem {
  text: string;
  href: string;
  icon: IconName;
}

const menuItems: IMenuItem[] = [
  {
    icon: "play-circle",
    text: "Scenes",
    href: "/scenes"
  },
  {
    href: "/scenes/markers",
    icon: "map-marker-alt",
    text: "Markers"
  },
  {
    href: "/galleries",
    icon: "image",
    text: "Galleries"
  },
  {
    href: "/performers",
    icon: "user",
    text: "Performers"
  },
  {
    href: "/studios",
    icon: "video",
    text: "Studios"
  },
  {
    href: "/tags",
    icon: "tag",
    text: "Tags"
  }
];

export const MainNavbar: React.FC = () => {
  const location = useLocation();

  const path =
    location.pathname === "/performers"
      ? "/performers/new"
      : location.pathname === "/studios"
      ? "/studios/new"
      : null;
  const newButton =
    path === null ? (
      ""
    ) : (
      <LinkContainer to={path}>
        <Button variant="primary">New</Button>
      </LinkContainer>
    );

  return (
    <Navbar fixed="top" variant="dark" bg="dark" className="top-nav">
      <Navbar.Brand as="div">
        <Link to="/">
          <Button className="minimal brand-link d-none d-sm-inline-block">
            Stash
          </Button>
          <Button className="minimal brand-icon d-inline d-sm-none">
            <img src="http://192.168.1.65:9999/favicon.ico" alt="" />
          </Button>
        </Link>
      </Navbar.Brand>
      <Nav className="mr-md-auto">
        {menuItems.map(i => (
          <LinkContainer
            activeClassName="active"
            to={i.href}
            key={i.href}
          >
            <Button className="minimal">
              <Icon icon={i.icon} />
              <span className="d-none d-sm-inline">{i.text}</span>
            </Button>
          </LinkContainer>
        ))}
      </Nav>
      <Nav>
        <div className="d-none d-sm-block">{newButton}</div>
        <LinkContainer exact to="/settings">
          <Button className="minimal">
            <Icon icon="cog" />
          </Button>
        </LinkContainer>
      </Nav>
    </Navbar>
  );
};
