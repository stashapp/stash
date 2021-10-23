import React from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { mutateReloadPlugins, usePlugins } from "src/core/StashService";
import { useToast } from "src/hooks";
import { TextUtils } from "src/utils";
import { CollapseButton, Icon, LoadingIndicator } from "src/components/Shared";

export const SettingsPluginsPanel: React.FC = () => {
  const Toast = useToast();
  const intl = useIntl();

  const { data, loading } = usePlugins();

  async function onReloadPlugins() {
    await mutateReloadPlugins().catch((e) => Toast.error(e));
  }

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
            <Icon icon="link" />
          </a>
        </Button>
      );
    }
  }

  function renderPlugins() {
    const elements = (data?.plugins ?? []).map((plugin) => (
      <div key={plugin.id}>
        <h4>
          {plugin.name} {plugin.version ? `(${plugin.version})` : undefined}{" "}
          {renderLink(plugin.url ?? undefined)}
        </h4>
        {plugin.description ? (
          <small className="text-muted">{plugin.description}</small>
        ) : undefined}
        {renderPluginHooks(plugin.hooks ?? undefined)}
        <hr />
      </div>
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
      <div className="mt-2">
        <h5>
          <FormattedMessage id="config.plugins.hooks" />
        </h5>
        {hooks.map((h) => (
          <div key={`${h.name}`} className="mb-3">
            <h6>{h.name}</h6>
            <CollapseButton
              text={intl.formatMessage({ id: "config.plugins.triggers_on" })}
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
    );
  }

  if (loading) return <LoadingIndicator />;

  return (
    <>
      <h3>
        <FormattedMessage id="config.categories.plugins" />
      </h3>
      <hr />
      {renderPlugins()}
      <Button onClick={() => onReloadPlugins()}>
        <span className="fa-icon">
          <Icon icon="sync-alt" />
        </span>
        <span>
          <FormattedMessage id="actions.reload_plugins" />
        </span>
      </Button>
    </>
  );
};
