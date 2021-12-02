import React from "react";
import { useIntl } from "react-intl";
import { Button, Form } from "react-bootstrap";
import { mutateRunPluginTask, usePlugins } from "src/core/StashService";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { SettingSection } from "../SettingSection";
import { Setting } from "../Inputs";

type Plugin = Pick<GQL.Plugin, "id">;
type PluginTask = Pick<GQL.PluginTask, "name" | "description">;

export const PluginTasks: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();

  const plugins = usePlugins();

  function renderPlugins() {
    if (!plugins.data || !plugins.data.plugins) {
      return;
    }

    const taskPlugins = plugins.data.plugins.filter(
      (p) => p.tasks && p.tasks.length > 0
    );

    return (
      <SettingSection headingID="config.tasks.plugin_tasks">
        {taskPlugins.map((o) => {
          return (
            <div key={`${o.id}`} className="setting-group">
              <Setting heading={o.name}></Setting>
              {renderPluginTasks(o, o.tasks ?? [])}
            </div>
          );
        })}
      </SettingSection>
    );
  }

  function renderPluginTasks(plugin: Plugin, pluginTasks: PluginTask[]) {
    if (!pluginTasks) {
      return;
    }

    return pluginTasks.map((o) => {
      return (
        <Setting
          headingID={o.name}
          subHeadingID={o.description ?? undefined}
          key={o.name}
        >
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
    Toast.success({
      content: intl.formatMessage(
        { id: "config.tasks.added_job_to_queue" },
        { operation_name: operation.name }
      ),
    });
  }

  return <Form.Group>{renderPlugins()}</Form.Group>;
};
