import React, { useEffect, useRef, useState, useCallback } from "react";
import {
  defineMessages,
  FormattedMessage,
  MessageDescriptor,
  useIntl,
} from "react-intl";
import { Nav, Navbar, Button, Fade } from "react-bootstrap";
import { IconName } from "@fortawesome/fontawesome-svg-core";
import { LinkContainer } from "react-router-bootstrap";
import { Link, NavLink, useLocation, useHistory } from "react-router-dom";
import Mousetrap from "mousetrap";

import { SessionUtils } from "src/utils";
import { Icon } from "src/components/Shared";
import { ConfigurationContext } from "src/hooks/Config";
import { Manual } from "./Help/Manual";

interface IMenuItem {
  name: string;
  message: MessageDescriptor;
  href: string;
  icon: IconName;
}

const messages = defineMessages({
  scenes: {
    id: "scenes",
    defaultMessage: "Scenes",
  },
  images: {
    id: "images",
    defaultMessage: "Images",
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
  sceneTagger: {
    id: "sceneTagger",
    defaultMessage: "Scene Tagger",
  },
  donate: {
    id: "donate",
    defaultMessage: "Donate",
  },
});

const allMenuItems: IMenuItem[] = [
  {
    name: "scenes",
    message: messages.scenes,
    href: "/scenes",
    icon: "play-circle",
  },
  {
    name: "images",
    message: messages.images,
    href: "/images",
    icon: "image",
  },
  {
    name: "movies",
    message: messages.movies,
    href: "/movies",
    icon: "film",
  },
  {
    name: "markers",
    message: messages.markers,
    href: "/scenes/markers",
    icon: "map-marker-alt",
  },
  {
    name: "galleries",
    message: messages.galleries,
    href: "/galleries",
    icon: "images",
  },
  {
    name: "performers",
    message: messages.performers,
    href: "/performers",
    icon: "user",
  },
  {
    name: "studios",
    message: messages.studios,
    href: "/studios",
    icon: "video",
  },
  {
    name: "tags",
    message: messages.tags,
    href: "/tags",
    icon: "tag",
  },
];

export const MainNavbar: React.FC = () => {
  const history = useHistory();
  const location = useLocation();
  const { configuration, loading } = React.useContext(ConfigurationContext);

  // Show all menu items by default, unless config says otherwise
  const [menuItems, setMenuItems] = useState<IMenuItem[]>(allMenuItems);

  const [expanded, setExpanded] = useState(false);
  const [showManual, setShowManual] = useState(false);

  useEffect(() => {
    const iCfg = configuration?.interface;
    if (iCfg?.menuItems) {
      setMenuItems(
        allMenuItems.filter((menuItem) =>
          iCfg.menuItems!.includes(menuItem.name)
        )
      );
    }
  }, [configuration]);

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
      : location.pathname === "/galleries"
      ? "/galleries/new"
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
    Mousetrap.bind("g i", () => goto("/images"));
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

  const handleDismiss = useCallback(() => setExpanded(false), [setExpanded]);

  return (
    <>
      <Manual show={showManual} onClose={() => setShowManual(false)} />
      <Navbar
        collapseOnSelect
        fixed="top"
        variant="dark"
        bg="dark"
        className="top-nav justify-content-start"
        expand="xl"
        expanded={expanded}
        onToggle={setExpanded}
        ref={navbarRef}
      >
        <Navbar.Toggle className="mr-3" />
        <Navbar.Brand as="div" onClick={handleDismiss}>
          <Link to="/">
            <Button className="minimal brand-link d-inline-block">Stash</Button>
          </Link>
        </Navbar.Brand>

        <Nav className="flex-row ml-auto order-xl-1">
          {!!newButton && <div className="mr-2">{newButton}</div>}
          <Nav.Link
            href="https://opencollective.com/stashapp"
            target="_blank"
            onClick={handleDismiss}
          >
            <Button className="minimal donate" title="Donate">
              <Icon icon="heart" />
              <span>{intl.formatMessage(messages.donate)}</span>
            </Button>
          </Nav.Link>
          <NavLink exact to="/settings" onClick={handleDismiss}>
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

        <Navbar.Collapse>
          <Fade in={!loading}>
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
          </Fade>
        </Navbar.Collapse>
      </Navbar>
    </>
  );
};
