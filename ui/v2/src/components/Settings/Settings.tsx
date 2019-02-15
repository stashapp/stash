import {
  Card,
  Tab,
  Tabs,
} from "@blueprintjs/core";
import queryString from "query-string";
import React, { FunctionComponent, useEffect, useState } from "react";
import { IBaseProps } from "../../models";
import { SettingsAboutPanel } from "./SettingsAboutPanel";
import { SettingsConfigurationPanel } from "./SettingsConfigurationPanel";
import { SettingsLogsPanel } from "./SettingsLogsPanel";
import { SettingsTasksPanel } from "./SettingsTasksPanel";

interface IProps extends IBaseProps {}

type TabId = "configuration" | "tasks" | "logs" | "about";

export const Settings: FunctionComponent<IProps> = (props: IProps) => {
  const [tabId, setTabId] = useState<TabId>(getTabId());

  useEffect(() => {
    const location = Object.assign({}, props.history.location);
    location.search = queryString.stringify({tab: tabId}, {encode: false});
    props.history.replace(location);
  }, [tabId]);

  function getTabId(): TabId {
    const queryParams = queryString.parse(props.location.search);
    if (!queryParams.tab || typeof queryParams.tab !== "string") { return "tasks"; }
    return queryParams.tab as TabId;
  }

  return (
    <Card id="details-container">
      <Tabs
        renderActiveTabPanelOnly={true}
        vertical={true}
        onChange={(newId) => setTabId(newId as TabId)}
        defaultSelectedTabId={getTabId()}
      >
        <Tab id="configuration" title="Configuration" panel={<SettingsConfigurationPanel />} />
        <Tab id="tasks" title="Tasks" panel={<SettingsTasksPanel />} />
        <Tab id="logs" title="Logs" panel={<SettingsLogsPanel />} />
        <Tab id="about" title="About" panel={<SettingsAboutPanel />} />
      </Tabs>
    </Card>
  );
};
