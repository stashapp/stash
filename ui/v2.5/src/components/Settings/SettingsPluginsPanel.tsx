import React, { useState, useEffect } from "react";
import { Button } from "react-bootstrap";
import {
  mutateReloadPlugins,
  usePlugins,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { Icon, LoadingIndicator } from "../Shared";

export const SettingsPluginsPanel: React.FC = () => {
  const Toast = useToast();
  
  const plugins = usePlugins();

  // Network state
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (plugins) {  
      setIsLoading(false);
    }
  }, [plugins]);

  async function onReloadPlugins() {
    setIsLoading(true);
    try {
      await mutateReloadPlugins();

      // reload the performer scrapers
      await plugins.refetch();
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function renderPlugins() {
    if (!plugins.data || !plugins.data.plugins) {
      return;
    }

    return (
      <ul>
        {plugins.data?.plugins.map(p => {
          return <li key={p.name}>{p.name}</li>;
        })}
      </ul>
    );
  }

  if (isLoading) return <LoadingIndicator />;

  return (
    <>
      <h5>Plugins</h5>
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
