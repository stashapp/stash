import React from "react";
import { useIntl } from "react-intl";
import { Button, Form } from "react-bootstrap";
import { mutateRunPluginTask, usePlugins } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import * as GQL from "src/core/generated-graphql";
import { SettingSection } from "../SettingSection";
import { Setting, SettingGroup } from "../Inputs";

type Plugin = Pick<GQL.Plugin, "id">;
type PluginTask = Pick<GQL.PluginTask, "name" | "description">;

export const PluginTasks: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();

  const plugins = usePlugins();

  function renderPluginTasks(plugin: Plugin, pluginTasks: PluginTask[]) {
    return pluginTasks.map((o) => {
      return (
        <Setting heading={o.name} subHeading={o.description} key={o.name}>
          <Button
            onClick={() => onPluginTaskClicked(plugin, o)}
            variant="secondary"
            size="sm"
          >
            {o.name}
          </Button>
        </Setting>
      );
    });
  }

  async function onPluginTaskClicked(plugin: Plugin, operation: PluginTask) {
    await mutateRunPluginTask(plugin.id, operation.name);
    Toast.success(
      intl.formatMessage(
        { id: "config.tasks.added_job_to_queue" },
        { operation_name: operation.name }
      )
    );
  }

  if (!plugins.data?.plugins) {
    return null;
  }

  const taskPlugins = plugins.data.plugins.filter(
    (p) => p.enabled && p.tasks && p.tasks.length > 0
  );

  if (!taskPlugins.length) {
    return null;
  }

  return (
    <Form.Group>
      <SettingSection headingID="config.tasks.plugin_tasks">
        {taskPlugins.map((o) => {
          return (
            <SettingGroup
              key={o.id}
              settingProps={{
                heading: o.name,
              }}
              collapsible
            >
              {renderPluginTasks(o, o.tasks!)}
            </SettingGroup>
          );
        })}
      </SettingSection>
    </Form.Group>
  );
};
