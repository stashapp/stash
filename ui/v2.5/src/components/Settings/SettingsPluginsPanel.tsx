import React, { useMemo } from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  mutateReloadPlugins,
  mutateSetPluginsEnabled,
  usePlugins,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import TextUtils from "src/utils/text";
import { CollapseButton } from "../Shared/CollapseButton";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { SettingSection } from "./SettingSection";
import { Setting, SettingGroup } from "./Inputs";
import { faLink, faSyncAlt } from "@fortawesome/free-solid-svg-icons";

export const SettingsPluginsPanel: React.FC = () => {
  const Toast = useToast();
  const intl = useIntl();

  const [changedPluginID, setChangedPluginID] = React.useState<
    string | undefined
  >();

  const { data, loading, refetch } = usePlugins();

  async function onReloadPlugins() {
    await mutateReloadPlugins().catch((e) => Toast.error(e));
  }

  const pluginElements = useMemo(() => {
    function renderLink(url?: string) {
      if (url) {
        return (
          <Button className="minimal">
            <a
              href={TextUtils.sanitiseURL(url)}
              className="link"
              target="_blank"
              rel="noopener noreferrer"
            >
              <Icon icon={faLink} />
            </a>
          </Button>
        );
      }
    }

    function renderEnableButton(pluginID: string, enabled: boolean) {
      async function onClick() {
        await mutateSetPluginsEnabled({ [pluginID]: !enabled }).catch((e) =>
          Toast.error(e)
        );

        setChangedPluginID(pluginID);
        refetch();
      }

      return (
        <Button size="sm" onClick={onClick}>
          <FormattedMessage
            id={enabled ? "actions.disable" : "actions.enable"}
          />
        </Button>
      );
    }

    function onReloadUI() {
      window.location.reload();
    }

    function maybeRenderReloadUI(pluginID: string) {
      if (pluginID === changedPluginID) {
        return (
          <Button size="sm" onClick={() => onReloadUI()}>
            Reload UI
          </Button>
        );
      }
    }

    function renderPlugins() {
      const elements = (data?.plugins ?? []).map((plugin) => (
        <SettingGroup
          key={plugin.id}
          settingProps={{
            heading: `${plugin.name} ${
              plugin.version ? `(${plugin.version})` : undefined
            }`,
            className: !plugin.enabled ? "disabled" : undefined,
            subHeading: plugin.description,
          }}
          topLevel={
            <>
              {renderLink(plugin.url ?? undefined)}
              {maybeRenderReloadUI(plugin.id)}
              {renderEnableButton(plugin.id, plugin.enabled)}
            </>
          }
        >
          {renderPluginHooks(plugin.hooks ?? undefined)}
        </SettingGroup>
      ));

      return <div>{elements}</div>;
    }

    function renderPluginHooks(
      hooks?: Pick<GQL.PluginHook, "name" | "description" | "hooks">[]
    ) {
      if (!hooks || hooks.length === 0) {
        return;
      }

      return (
        <div className="setting">
          <div>
            <h5>
              <FormattedMessage id="config.plugins.hooks" />
            </h5>
            {hooks.map((h) => (
              <div key={`${h.name}`}>
                <h6>{h.name}</h6>
                <CollapseButton
                  text={intl.formatMessage({
                    id: "config.plugins.triggers_on",
                  })}
                >
                  <ul>
                    {h.hooks?.map((hh) => (
                      <li key={hh}>
                        <code>{hh}</code>
                      </li>
                    ))}
                  </ul>
                </CollapseButton>
                <small className="text-muted">{h.description}</small>
              </div>
            ))}
          </div>
          <div />
        </div>
      );
    }

    return renderPlugins();
  }, [data?.plugins, intl, Toast, changedPluginID, refetch]);

  if (loading) return <LoadingIndicator />;

  return (
    <>
      <SettingSection headingID="config.categories.plugins">
        <Setting headingID="actions.reload_plugins">
          <Button onClick={() => onReloadPlugins()}>
            <span className="fa-icon">
              <Icon icon={faSyncAlt} />
            </span>
            <span>
              <FormattedMessage id="actions.reload_plugins" />
            </span>
          </Button>
        </Setting>
        {pluginElements}
      </SettingSection>
    </>
  );
};
