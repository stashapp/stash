import React from "react";
import { Tab, Nav, Row, Col } from "react-bootstrap";
import { useHistory, useLocation } from "react-router-dom";
import { FormattedMessage } from "react-intl";
import { Helmet } from "react-helmet";
import { useTitleProps } from "src/hooks/title";
import { SettingsAboutPanel } from "./SettingsAboutPanel";
import { SettingsConfigurationPanel } from "./SettingsSystemPanel";
import { SettingsInterfacePanel } from "./SettingsInterfacePanel/SettingsInterfacePanel";
import { SettingsLogsPanel } from "./SettingsLogsPanel";
import { SettingsTasksPanel } from "./Tasks/SettingsTasksPanel";
import { SettingsPluginsPanel } from "./SettingsPluginsPanel";
import { SettingsScrapingPanel } from "./SettingsScrapingPanel";
import { SettingsToolsPanel } from "./SettingsToolsPanel";
import { SettingsServicesPanel } from "./SettingsServicesPanel";
import { SettingsContext } from "./context";
import { SettingsLibraryPanel } from "./SettingsLibraryPanel";
import { SettingsSecurityPanel } from "./SettingsSecurityPanel";
import Changelog from "../Changelog/Changelog";

export const Settings: React.FC = () => {
  const location = useLocation();
  const history = useHistory();
  const defaultTab = new URLSearchParams(location.search).get("tab") ?? "tasks";

  const onSelect = (val: string) => history.push(`?tab=${val}`);

  const titleProps = useTitleProps({ id: "settings" });
  return (
    <Tab.Container
      activeKey={defaultTab}
      id="configuration-tabs"
      onSelect={(tab) => tab && onSelect(tab)}
    >
      <Helmet {...titleProps} />
      <Row>
        <Col id="settings-menu-container" sm={3} md={3} xl={2}>
          <Nav variant="pills" className="flex-column">
            <Nav.Item>
              <Nav.Link eventKey="tasks">
                <FormattedMessage id="config.categories.tasks" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="library">
                <FormattedMessage id="library" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="interface">
                <FormattedMessage id="config.categories.interface" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="security">
                <FormattedMessage id="config.categories.security" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="metadata-providers">
                <FormattedMessage id="config.categories.metadata_providers" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="services">
                <FormattedMessage id="config.categories.services" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="system">
                <FormattedMessage id="config.categories.system" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="plugins">
                <FormattedMessage id="config.categories.plugins" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="logs">
                <FormattedMessage id="config.categories.logs" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="tools">
                <FormattedMessage id="config.categories.tools" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="changelog">
                <FormattedMessage id="config.categories.changelog" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="about">
                <FormattedMessage id="config.categories.about" />
              </Nav.Link>
            </Nav.Item>
            <hr className="d-sm-none" />
          </Nav>
        </Col>
        <Col
          id="settings-container"
          sm={{ offset: 3 }}
          md={{ offset: 3 }}
          xl={{ offset: 2 }}
        >
          <SettingsContext>
            <Tab.Content className="mx-auto">
              <Tab.Pane eventKey="library">
                <SettingsLibraryPanel />
              </Tab.Pane>
              <Tab.Pane eventKey="interface">
                <SettingsInterfacePanel />
              </Tab.Pane>
              <Tab.Pane eventKey="security">
                <SettingsSecurityPanel />
              </Tab.Pane>
              <Tab.Pane eventKey="tasks">
                <SettingsTasksPanel />
              </Tab.Pane>
              <Tab.Pane eventKey="services" unmountOnExit>
                <SettingsServicesPanel />
              </Tab.Pane>
              <Tab.Pane eventKey="tools" unmountOnExit>
                <SettingsToolsPanel />
              </Tab.Pane>
              <Tab.Pane eventKey="metadata-providers" unmountOnExit>
                <SettingsScrapingPanel />
              </Tab.Pane>
              <Tab.Pane eventKey="system">
                <SettingsConfigurationPanel />
              </Tab.Pane>
              <Tab.Pane eventKey="plugins" unmountOnExit>
                <SettingsPluginsPanel />
              </Tab.Pane>
              <Tab.Pane eventKey="logs" unmountOnExit>
                <SettingsLogsPanel />
              </Tab.Pane>
              <Tab.Pane eventKey="changelog" unmountOnExit>
                <Changelog />
              </Tab.Pane>
              <Tab.Pane eventKey="about" unmountOnExit>
                <SettingsAboutPanel />
              </Tab.Pane>
            </Tab.Content>
          </SettingsContext>
        </Col>
      </Row>
    </Tab.Container>
  );
};

export default Settings;
