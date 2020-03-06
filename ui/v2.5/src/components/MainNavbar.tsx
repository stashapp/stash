import React, { useEffect, useRef, useState } from "react";
import { FormattedMessage } from "react-intl";
import { Nav, Navbar, Button } from "react-bootstrap";
import { IconName } from "@fortawesome/fontawesome-svg-core";
import { LinkContainer } from "react-router-bootstrap";
import { Link, useLocation } from "react-router-dom";

import { Icon } from "src/components/Shared";

interface IMenuItem {
  messageID: string;
  href: string;
  icon: IconName;
}

const menuItems: IMenuItem[] = [
  {
    icon: "play-circle",
    messageID: "scenes",
    href: "/scenes"
  },
  {
    href: "/scenes/markers",
    icon: "map-marker-alt",
    messageID: "markers"
  },
  {
    href: "/galleries",
    icon: "image",
    messageID: "galleries"
  },
  {
    href: "/performers",
    icon: "user",
    messageID: "performers"
  },
  {
    href: "/studios",
    icon: "video",
    messageID: "studios"
  },
  {
    href: "/tags",
    icon: "tag",
    messageID: "tags"
  }
];

export const MainNavbar: React.FC = () => {
  const location = useLocation();
  const [expanded, setExpanded] = useState(false);
  // react-bootstrap typing bug
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const navbarRef = useRef<any>();

  const maybeCollapse = (event: Event) => {
    if (navbarRef.current && event.target instanceof Node && !navbarRef.current.contains(event.target)) {
      setExpanded(false);
    }
  };

  useEffect(() => {
    if(expanded) {
      document.addEventListener('click', maybeCollapse);
      document.addEventListener('touchstart', maybeCollapse);
    }
    return () => {
      document.removeEventListener('click', maybeCollapse);
      document.removeEventListener('touchstart', maybeCollapse);
    }
  }, [expanded]);

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
        <Button variant="primary">
          <FormattedMessage id="new" defaultMessage="New" />
        </Button>
      </LinkContainer>
    );

  return (
    <Navbar
      collapseOnSelect
      fixed="top"
      variant="dark"
      bg="dark"
      className="top-nav"
      expand="md"
      expanded={expanded}
      onToggle={setExpanded}
      ref={navbarRef}
    >
        <Navbar.Brand as="div" className="order-1 order-md-0" onClick={() => setExpanded(false)}>
        <Link to="/">
          <Button className="minimal brand-link d-none d-md-inline-block">
            Stash
          </Button>
          <Button className="minimal brand-icon d-inline d-md-none">
            <img src="favicon.ico" alt="" />
          </Button>
        </Link>
      </Navbar.Brand>
      <Navbar.Toggle className="order-0" />
      <Navbar.Collapse className="order-3 order-md-1">
        <Nav className="mr-md-auto">
          {menuItems.map(i => (
            <Nav.Link eventKey={i.href} as="div" key={i.href}>
              <LinkContainer
                activeClassName="active"
                exact
                to={i.href}
                key={i.href}
              >
                <Button className="minimal w-100">
                  <Icon icon={i.icon} />
                  <span>
                    <FormattedMessage id={i.messageID} />
                  </span>
                </Button>
              </LinkContainer>
            </Nav.Link>
          ))}
        </Nav>
      </Navbar.Collapse>
      <Nav className="order-2">
        <div className="d-none d-sm-block">{newButton}</div>
        <LinkContainer exact to="/settings" onClick={() => setExpanded(false)}>
          <Button className="minimal settings-button">
            <Icon icon="cog" />
          </Button>
        </LinkContainer>
      </Nav>
    </Navbar>
  );
};
