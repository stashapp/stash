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
import {
  BooleanSetting,
  NumberSetting,
  Setting,
  SettingGroup,
  StringSetting,
} from "./Inputs";
import { faLink, faSyncAlt } from "@fortawesome/free-solid-svg-icons";
import { useSettings } from "./context";
import {
  AvailablePluginPackages,
  InstalledPluginPackages,
} from "./PluginPackageManager";
import { ExternalLink } from "../Shared/ExternalLink";
import { PatchComponent } from "src/patch";

interface IPluginSettingProps {
  pluginID: string;
  setting: GQL.PluginSetting;
  value: unknown;
  onChange: (value: unknown) => void;
}

const PluginSetting: React.FC<IPluginSettingProps> = ({
  pluginID,
  setting,
  value,
  onChange,
}) => {
  const commonProps = {
    heading: setting.display_name ? setting.display_name : setting.name,
    id: `plugin-${pluginID}-${setting.name}`,
    subHeading: setting.description ?? undefined,
  };

  switch (setting.type) {
    case GQL.PluginSettingTypeEnum.Boolean:
      return (
        <BooleanSetting
          {...commonProps}
          checked={(value as boolean) ?? false}
          onChange={() => onChange(!value)}
        />
      );
    case GQL.PluginSettingTypeEnum.String:
      return (
        <StringSetting
          {...commonProps}
          value={(value as string) ?? ""}
          onChange={(v) => onChange(v)}
        />
      );
    case GQL.PluginSettingTypeEnum.Number:
      return (
        <NumberSetting
          {...commonProps}
          value={(value as number) ?? 0}
          onChange={(v) => onChange(v)}
        />
      );
  }
};

const PluginSettings: React.FC<{
  pluginID: string;
  settings: GQL.PluginSetting[];
}> = PatchComponent("PluginSettings", ({ pluginID, settings }) => {
  const { plugins, savePluginSettings } = useSettings();
  const pluginSettings = plugins[pluginID] ?? {};

  return (
    <div className="plugin-settings">
      {settings.map((setting) => (
        <PluginSetting
          key={setting.name}
          pluginID={pluginID}
          setting={setting}
          value={pluginSettings[setting.name]}
          onChange={(v) =>
            savePluginSettings(pluginID, {
              ...pluginSettings,
              [setting.name]: v,
            })
          }
        />
      ))}
    </div>
  );
});

export const SettingsPluginsPanel: React.FC = () => {
  const Toast = useToast();
  const intl = useIntl();

  const { loading: configLoading } = useSettings();
  const { data, loading } = usePlugins();

  const [changedPluginID, setChangedPluginID] = React.useState<
    string | undefined
  >();

  async function onReloadPlugins() {
    try {
      await mutateReloadPlugins();
    } catch (e) {
      Toast.error(e);
    }
  }

  const pluginElements = useMemo(() => {
    function renderLink(url?: string) {
      if (url) {
        return (
          <Button
            as={ExternalLink}
            href={TextUtils.sanitiseURL(url)}
            className="minimal link"
          >
            <Icon icon={faLink} />
          </Button>
        );
      }
    }

    function renderEnableButton(pluginID: string, enabled: boolean) {
      async function onClick() {
        try {
          await mutateSetPluginsEnabled({ [pluginID]: !enabled });
        } catch (e) {
          Toast.error(e);
        }

        setChangedPluginID(pluginID);
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
          <PluginSettings
            pluginID={plugin.id}
            settings={plugin.settings ?? []}
          />
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
  }, [data?.plugins, intl, Toast, changedPluginID]);

  if (loading || configLoading) return <LoadingIndicator />;

  return (
    <>
      <InstalledPluginPackages />
      <AvailablePluginPackages />

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
