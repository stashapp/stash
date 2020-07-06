import React, { useEffect, useRef, useState } from "react";
import {
  defineMessages,
  FormattedMessage,
  MessageDescriptor,
  useIntl,
} from "react-intl";
import { Nav, Navbar, Button } from "react-bootstrap";
import { IconName } from "@fortawesome/fontawesome-svg-core";
import { LinkContainer } from "react-router-bootstrap";
import { Link, NavLink, useLocation, useHistory } from "react-router-dom";
import { SessionUtils } from "src/utils";

import { Icon } from "src/components/Shared";
import { Manual } from "./Help/Manual";

interface IMenuItem {
  message: MessageDescriptor;
  href: string;
  icon: IconName;
}

const messages = defineMessages({
  scenes: {
    id: "scenes",
    defaultMessage: "Scenes",
  },
  movies: {
    id: "movies",
    defaultMessage: "Movies",
  },
  markers: {
    id: "markers",
    defaultMessage: "Markers",
  },
  performers: {
    id: "performers",
    defaultMessage: "Performers",
  },
  studios: {
    id: "studios",
    defaultMessage: "Studios",
  },
  tags: {
    id: "tags",
    defaultMessage: "Tags",
  },
  galleries: {
    id: "galleries",
    defaultMessage: "Galleries",
  },
});

const menuItems: IMenuItem[] = [
  {
    icon: "play-circle",
    message: messages.scenes,
    href: "/scenes",
  },
  {
    href: "/movies",
    icon: "film",
    message: messages.movies,
  },
  {
    href: "/scenes/markers",
    icon: "map-marker-alt",
    message: messages.markers,
  },
  {
    href: "/galleries",
    icon: "image",
    message: messages.galleries,
  },
  {
    href: "/performers",
    icon: "user",
    message: messages.performers,
  },
  {
    href: "/studios",
    icon: "video",
    message: messages.studios,
  },
  {
    href: "/tags",
    icon: "tag",
    message: messages.tags,
  },
];

export const MainNavbar: React.FC = () => {
  const history = useHistory();
  const location = useLocation();
  const [expanded, setExpanded] = useState(false);
  const [showManual, setShowManual] = useState(false);

  // react-bootstrap typing bug
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const navbarRef = useRef<any>();
  const intl = useIntl();

  const maybeCollapse = (event: Event) => {
    if (
      navbarRef.current &&
      event.target instanceof Node &&
      !navbarRef.current.contains(event.target)
    ) {
      setExpanded(false);
    }
  };

  useEffect(() => {
    if (expanded) {
      document.addEventListener("click", maybeCollapse);
      document.addEventListener("touchstart", maybeCollapse);
    }
    return () => {
      document.removeEventListener("click", maybeCollapse);
      document.removeEventListener("touchstart", maybeCollapse);
    };
  }, [expanded]);

  function goto(page: string) {
    history.push(page);
    if (document.activeElement instanceof HTMLElement) {
      document.activeElement.blur();
    }
  }

  const newPath =
    location.pathname === "/performers"
      ? "/performers/new"
      : location.pathname === "/studios"
      ? "/studios/new"
      : location.pathname === "/movies"
      ? "/movies/new"
      : location.pathname === "/tags"
      ? "/tags/new"
      : null;
  const newButton =
    newPath === null ? (
      ""
    ) : (
      <Link to={newPath}>
        <Button variant="primary">
          <FormattedMessage id="new" defaultMessage="New" />
        </Button>
      </Link>
    );

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("?", () => setShowManual(!showManual));
    Mousetrap.bind("g s", () => goto("/scenes"));
    Mousetrap.bind("g v", () => goto("/movies"));
    Mousetrap.bind("g k", () => goto("/scenes/markers"));
    Mousetrap.bind("g l", () => goto("/galleries"));
    Mousetrap.bind("g p", () => goto("/performers"));
    Mousetrap.bind("g u", () => goto("/studios"));
    Mousetrap.bind("g t", () => goto("/tags"));
    Mousetrap.bind("g z", () => goto("/settings"));

    if (newPath) {
      Mousetrap.bind("n", () => history.push(newPath));
    }

    return () => {
      Mousetrap.unbind("?");
      Mousetrap.unbind("g s");
      Mousetrap.unbind("g v");
      Mousetrap.unbind("g k");
      Mousetrap.unbind("g l");
      Mousetrap.unbind("g p");
      Mousetrap.unbind("g u");
      Mousetrap.unbind("g t");
      Mousetrap.unbind("g z");

      if (newPath) {
        Mousetrap.unbind("n");
      }
    };
  });

  function maybeRenderLogout() {
    if (SessionUtils.isLoggedIn()) {
      return (
        <Button className="minimal logout-button" href="/logout">
          <Icon icon="sign-out-alt" />
        </Button>
      );
    }
  }

  return (
    <>
      <Manual show={showManual} onClose={() => setShowManual(false)} />
      <Navbar
        collapseOnSelect
        fixed="top"
        variant="dark"
        bg="dark"
        className="top-nav"
        expand="lg"
        expanded={expanded}
        onToggle={setExpanded}
        ref={navbarRef}
      >
        <Navbar.Brand
          as="div"
          className="order-1 order-md-0"
          onClick={() => setExpanded(false)}
        >
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
            {menuItems.map((i) => (
              <Nav.Link eventKey={i.href} as="div" key={i.href}>
                <LinkContainer activeClassName="active" exact to={i.href}>
                  <Button className="minimal w-100">
                    <Icon icon={i.icon} />
                    <span>{intl.formatMessage(i.message)}</span>
                  </Button>
                </LinkContainer>
              </Nav.Link>
            ))}
          </Nav>
        </Navbar.Collapse>
        <Nav className="order-2 flex-row">
          <div className="d-none d-sm-block">{newButton}</div>
          <NavLink exact to="/settings" onClick={() => setExpanded(false)}>
            <Button className="minimal settings-button" title="Settings">
              <Icon icon="cog" />
            </Button>
          </NavLink>
          <Button
            className="minimal help-button"
            onClick={() => setShowManual(true)}
            title="Help"
          >
            <Icon icon="question-circle" />
          </Button>
          {maybeRenderLogout()}
        </Nav>
      </Navbar>
    </>
  );
};
