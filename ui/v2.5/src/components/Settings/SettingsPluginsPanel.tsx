import React, { useMemo } from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { mutateReloadPlugins, usePlugins } from "src/core/StashService";
import { useToast } from "src/hooks";
import { TextUtils } from "src/utils";
import { CollapseButton, Icon, LoadingIndicator } from "src/components/Shared";
import { SettingSection } from "./SettingSection";
import { Setting, SettingGroup } from "./Inputs";
import { faLink, faSyncAlt } from "@fortawesome/free-solid-svg-icons";

export const SettingsPluginsPanel: React.FC = () => {
  const Toast = useToast();
  const intl = useIntl();

  const { data, loading } = usePlugins();

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

    function renderPlugins() {
      const elements = (data?.plugins ?? []).map((plugin) => (
        <SettingGroup
          key={plugin.id}
          settingProps={{
            heading: `${plugin.name} ${
              plugin.version ? `(${plugin.version})` : undefined
            }`,
            subHeading: plugin.description,
          }}
          topLevel={renderLink(plugin.url ?? undefined)}
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
  }, [data?.plugins, intl]);

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
