import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Card, Row, Col, Tab, Nav } from "react-bootstrap";
import { LoadingIndicator, TITLE_SUFFIX } from "src/components/Shared";
import { JobTable } from "./JobTable";
import { Helmet } from "react-helmet";
import { LibraryTasks } from "./LibraryTasks";
import { DataManagementTasks } from "./DataManagementTasks";
import { PluginTasks } from "./PluginTasks";

export const Tasks: React.FC = () => {
  const intl = useIntl();
  const [activeTabKey, setActiveTabKey] = useState("library");

  const [isBackupRunning, setIsBackupRunning] = useState<boolean>(false);

  function renderTabNav() {
    function renderNavItem(eventKey: string, msgId: string) {
      return (
        <Nav.Item>
          <Nav.Link eventKey={eventKey}>
            <FormattedMessage id={msgId} />
          </Nav.Link>
        </Nav.Item>
      );
    }

    return (
      <Nav variant="pills" className="flex-column">
        {renderNavItem("library", "library")}
        {renderNavItem("dataManagement", "config.tasks.data_management")}
        {renderNavItem("plugins", "config.categories.plugins")}
      </Nav>
    );
  }

  function renderContentPane() {
    return (
      <Tab.Content>
        <Tab.Pane eventKey="library">
          <LibraryTasks />
        </Tab.Pane>
        <Tab.Pane eventKey="dataManagement">
          <DataManagementTasks setIsBackupRunning={setIsBackupRunning} />
        </Tab.Pane>
        <Tab.Pane eventKey="plugins">
          <PluginTasks />
        </Tab.Pane>
      </Tab.Content>
    );
  }

  if (isBackupRunning) {
    return (
      <LoadingIndicator
        message={intl.formatMessage({ id: "config.tasks.backing_up_database" })}
      />
    );
  }

  const title_template = `${intl.formatMessage({
    id: "config.categories.tasks",
  })} ${TITLE_SUFFIX}`;

  return (
    <Card className="col col-lg-9 mx-auto">
      <Helmet
        defaultTitle={title_template}
        titleTemplate={`%s | ${title_template}`}
      />

      <Tab.Container
        activeKey={activeTabKey}
        id="tasks-tabs"
        onSelect={(tab) => tab && setActiveTabKey(tab)}
      >
        <Row>
          <Col sm={3} md={3}>
            {renderTabNav()}
          </Col>
          <Col sm={9} md={9}>
            <h4>{intl.formatMessage({ id: "config.tasks.job_queue" })}</h4>

            <JobTable />

            <hr />

            {renderContentPane()}
          </Col>
        </Row>
      </Tab.Container>
    </Card>
  );
};
