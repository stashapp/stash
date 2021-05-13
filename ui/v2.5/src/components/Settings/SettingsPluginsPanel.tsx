import React from "react";
import { Button } from "react-bootstrap";
import { mutateReloadPlugins, usePlugins } from "src/core/StashService";
import { useToast } from "src/hooks";
import { TextUtils } from "src/utils";
import { Icon, LoadingIndicator } from "src/components/Shared";

export const SettingsPluginsPanel: React.FC = () => {
  const Toast = useToast();
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
        <h5>
          {plugin.name} {plugin.version ? `(${plugin.version})` : undefined}{" "}
          {renderLink(plugin.url ?? undefined)}
        </h5>
        {plugin.description ? (
          <small className="text-muted">{plugin.description}</small>
        ) : undefined}
        <hr />
      </div>
    ));

    return <div>{elements}</div>;
  }

  if (loading) return <LoadingIndicator />;

  return (
    <>
      <h4>Plugins</h4>
      <hr />
      {renderPlugins()}
      <Button onClick={() => onReloadPlugins()}>
        <span className="fa-icon">
          <Icon icon="sync-alt" />
        </span>
        <span>Reload plugins</span>
      </Button>
    </>
  );
};
