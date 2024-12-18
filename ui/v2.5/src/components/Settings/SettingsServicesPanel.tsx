import React, { useState } from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import {
  useDisableDLNA,
  useDLNAStatus,
  useEnableDLNA,
  useAddTempDLNAIP,
  useRemoveTempDLNAIP,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { DurationInput } from "../Shared/DurationInput";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { ModalComponent } from "../Shared/Modal";
import { SettingSection } from "./SettingSection";
import {
  BooleanSetting,
  StringListSetting,
  StringSetting,
  SelectSetting,
  NumberSetting,
} from "./Inputs";
import { useSettings } from "./context";
import {
  videoSortOrderIntlMap,
  defaultVideoSort,
} from "src/utils/dlnaVideoSort";
import {
  faClock,
  faTimes,
  faUserClock,
} from "@fortawesome/free-solid-svg-icons";

const defaultDLNAPort = 1338;

export const SettingsServicesPanel: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();

  const { dlna, loading: configLoading, error, saveDLNA } = useSettings();

  // undefined to hide dialog, true for enable, false for disable
  const [enableDisable, setEnableDisable] = useState<boolean>();

  const [enableUntilRestart, setEnableUntilRestart] = useState<boolean>(false);
  const [enableDuration, setEnableDuration] = useState<number>(0);

  const [ipEntry, setIPEntry] = useState<string>("");
  const [tempIP, setTempIP] = useState<string>();

  const { data: statusData, loading, refetch: statusRefetch } = useDLNAStatus();

  const [enableDLNA] = useEnableDLNA();
  const [disableDLNA] = useDisableDLNA();
  const [addTempDLANIP] = useAddTempDLNAIP();
  const [removeTempDLNAIP] = useRemoveTempDLNAIP();

  if (error) return <h1>{error.message}</h1>;
  if (loading || configLoading) return <LoadingIndicator />;

  async function onTempEnable() {
    const input = {
      variables: {
        input: {
          duration: enableUntilRestart ? undefined : enableDuration,
        },
      },
    };

    try {
      if (enableDisable) {
        await enableDLNA(input);
        Toast.success(
          intl.formatMessage({
            id: "config.dlna.enabled_dlna_temporarily",
          })
        );
      } else {
        await disableDLNA(input);
        Toast.success(
          intl.formatMessage({
            id: "config.dlna.disabled_dlna_temporarily",
          })
        );
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setEnableDisable(undefined);
      statusRefetch();
    }
  }

  async function onAllowTempIP() {
    if (!tempIP) {
      return;
    }

    const input = {
      variables: {
        input: {
          duration: enableUntilRestart ? undefined : enableDuration,
          address: tempIP,
        },
      },
    };

    try {
      await addTempDLANIP(input);
      Toast.success(
        intl.formatMessage({
          id: "config.dlna.allowed_ip_temporarily",
        })
      );
    } catch (e) {
      Toast.error(e);
    } finally {
      setTempIP(undefined);
      statusRefetch();
    }
  }

  async function onDisallowTempIP(address: string) {
    const input = {
      variables: {
        input: {
          address,
        },
      },
    };

    try {
      await removeTempDLNAIP(input);
      Toast.success(intl.formatMessage({ id: "config.dlna.disallowed_ip" }));
    } catch (e) {
      Toast.error(e);
    } finally {
      statusRefetch();
    }
  }

  function renderDeadline(until?: string | null) {
    if (until) {
      const deadline = new Date(until);
      return `until ${intl.formatDate(deadline)}`;
    }

    return "";
  }

  function renderStatus() {
    if (!statusData) {
      return "";
    }

    const { dlnaStatus } = statusData;
    const runningText = intl.formatMessage({
      id: dlnaStatus.running ? "actions.running" : "actions.not_running",
    });

    return `${runningText} ${renderDeadline(dlnaStatus.until)}`;
  }

  function renderEnableButton() {
    // if enabled by default, then show the disable temporarily
    // if disabled by default, then show enable temporarily
    if (dlna.enabled) {
      return (
        <Button onClick={() => setEnableDisable(false)} className="mr-1">
          <FormattedMessage id="actions.temp_disable" />
        </Button>
      );
    }

    return (
      <Button onClick={() => setEnableDisable(true)} className="mr-1">
        <FormattedMessage id="actions.temp_enable" />
      </Button>
    );
  }

  function canCancel() {
    if (!statusData || !dlna) {
      return false;
    }

    const { dlnaStatus } = statusData;
    const { enabled } = dlna;

    return dlnaStatus.until || dlnaStatus.running !== enabled;
  }

  async function cancelTempBehaviour() {
    if (!canCancel()) {
      return;
    }

    const running = statusData?.dlnaStatus.running;

    const input = {
      variables: {
        input: {},
      },
    };

    try {
      if (!running) {
        await enableDLNA(input);
      } else {
        await disableDLNA(input);
      }
      Toast.success(
        intl.formatMessage({
          id: "config.dlna.successfully_cancelled_temporary_behaviour",
        })
      );
    } catch (e) {
      Toast.error(e);
    } finally {
      setEnableDisable(undefined);
      statusRefetch();
    }
  }

  function renderTempCancelButton() {
    if (!canCancel()) {
      return;
    }

    return (
      <Button onClick={() => cancelTempBehaviour()} variant="danger">
        Cancel temporary behaviour
      </Button>
    );
  }

  function renderTempEnableDialog() {
    const text: string = enableDisable ? "enable" : "disable";
    const capitalised = `${text[0].toUpperCase()}${text.slice(1)}`;

    return (
      <ModalComponent
        show={enableDisable !== undefined}
        header={capitalised}
        icon={faClock}
        accept={{
          text: capitalised,
          variant: "primary",
          onClick: onTempEnable,
        }}
        cancel={{
          onClick: () => setEnableDisable(undefined),
          variant: "secondary",
        }}
      >
        <h4>{capitalised} temporarily</h4>
        <Form.Group>
          <Form.Check
            checked={enableUntilRestart}
            label={intl.formatMessage({ id: "config.dlna.until_restart" })}
            onChange={() => setEnableUntilRestart(!enableUntilRestart)}
          />
        </Form.Group>

        <Form.Group id="temp-enable-duration">
          <DurationInput
            value={enableDuration}
            setValue={(v) => setEnableDuration(v ?? 0)}
            disabled={enableUntilRestart}
          />
          <Form.Text className="text-muted">
            Duration to {text} for - in minutes.
          </Form.Text>
        </Form.Group>
      </ModalComponent>
    );
  }

  function renderTempWhitelistDialog() {
    return (
      <ModalComponent
        show={tempIP !== undefined}
        header={intl.formatMessage(
          { id: "config.dlna.allow_temp_ip" },
          { tempIP }
        )}
        icon={faClock}
        accept={{
          text: intl.formatMessage({ id: "actions.allow" }),
          variant: "primary",
          onClick: onAllowTempIP,
        }}
        cancel={{
          onClick: () => setTempIP(undefined),
          variant: "secondary",
        }}
      >
        <h4>{`Allow ${tempIP} temporarily`}</h4>
        <Form.Group>
          <Form.Check
            checked={enableUntilRestart}
            label={intl.formatMessage({ id: "config.dlna.until_restart" })}
            onChange={() => setEnableUntilRestart(!enableUntilRestart)}
          />
        </Form.Group>

        <Form.Group id="temp-enable-duration">
          <DurationInput
            value={enableDuration}
            setValue={(v) => setEnableDuration(v ?? 0)}
            disabled={enableUntilRestart}
          />
          <Form.Text className="text-muted">
            Duration to allow for - in minutes.
          </Form.Text>
        </Form.Group>
      </ModalComponent>
    );
  }

  function renderAllowedIPs() {
    if (!statusData || statusData.dlnaStatus.allowedIPAddresses.length === 0) {
      return;
    }

    const { allowedIPAddresses } = statusData.dlnaStatus;
    return (
      <Form.Group className="content">
        <h6>
          {intl.formatMessage({ id: "config.dlna.allowed_ip_addresses" })}
        </h6>

        <ul className="addresses">
          {allowedIPAddresses.map((a) => (
            <li key={a.ipAddress}>
              <div className="address">
                <code>{a.ipAddress}</code>
                <br />
                <span className="deadline">{renderDeadline(a.until)}</span>
              </div>

              <div className="buttons">
                <Button
                  size="sm"
                  title={intl.formatMessage({ id: "actions.disallow" })}
                  variant="danger"
                  onClick={() => onDisallowTempIP(a.ipAddress)}
                >
                  <Icon icon={faTimes} />
                </Button>
              </div>
            </li>
          ))}
        </ul>
      </Form.Group>
    );
  }

  function renderRecentIPs() {
    if (!statusData) {
      return;
    }

    const { recentIPAddresses } = statusData.dlnaStatus;
    return (
      <ul className="addresses">
        {recentIPAddresses.map((a) => (
          <li key={a}>
            <div className="address">
              <code>{a}</code>
            </div>
            <div>
              <Button
                size="sm"
                title={intl.formatMessage({ id: "actions.allow_temporarily" })}
                onClick={() => setTempIP(a)}
              >
                <Icon icon={faUserClock} />
              </Button>
            </div>
          </li>
        ))}
        <li>
          <div className="address">
            <Form.Control
              className="text-input"
              value={ipEntry}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setIPEntry(e.currentTarget.value)
              }
            />
          </div>
          <div className="buttons">
            <Button
              size="sm"
              title={intl.formatMessage({ id: "actions.allow_temporarily" })}
              onClick={() => setTempIP(ipEntry)}
              disabled={!ipEntry}
            >
              <Icon icon={faUserClock} />
            </Button>
          </div>
        </li>
      </ul>
    );
  }

  const DLNASettingsForm: React.FC = () => {
    return (
      <>
        <SettingSection headingID="settings">
          <StringSetting
            headingID="config.dlna.server_display_name"
            subHeading={intl.formatMessage(
              { id: "config.dlna.server_display_name_desc" },
              { server_name: <code>stash</code> }
            )}
            value={dlna.serverName ?? undefined}
            onChange={(v) => saveDLNA({ serverName: v })}
          />

          <NumberSetting
            headingID="config.dlna.server_port"
            subHeading={intl.formatMessage({
              id: "config.dlna.server_port_desc",
            })}
            value={dlna.port ?? undefined}
            onChange={(v) => saveDLNA({ port: v ? v : defaultDLNAPort })}
          />

          <BooleanSetting
            id="dlna-enabled-by-default"
            headingID="config.dlna.enabled_by_default"
            checked={dlna.enabled ?? undefined}
            onChange={(v) => saveDLNA({ enabled: v })}
          />

          <StringListSetting
            id="dlna-network-interfaces"
            headingID="config.dlna.network_interfaces"
            subHeadingID="config.dlna.network_interfaces_desc"
            value={dlna.interfaces ?? undefined}
            onChange={(v) => saveDLNA({ interfaces: v })}
          />

          <StringListSetting
            id="dlna-default-ip-whitelist"
            headingID="config.dlna.default_ip_whitelist"
            subHeading={intl.formatMessage(
              { id: "config.dlna.default_ip_whitelist_desc" },
              { wildcard: <code>*</code> }
            )}
            defaultNewValue="*"
            value={dlna.whitelistedIPs ?? undefined}
            onChange={(v) => saveDLNA({ whitelistedIPs: v })}
          />

          <SelectSetting
            id="video-sort-order"
            headingID="config.dlna.video_sort_order"
            subHeadingID="config.dlna.video_sort_order_desc"
            value={dlna.videoSortOrder ?? defaultVideoSort}
            onChange={(v) => saveDLNA({ videoSortOrder: v })}
          >
            {Array.from(videoSortOrderIntlMap.entries()).map((v) => (
              <option key={v[0]} value={v[0]}>
                {intl.formatMessage({
                  id: v[1],
                })}
              </option>
            ))}
          </SelectSetting>
        </SettingSection>
      </>
    );
  };

  return (
    <div id="settings-dlna">
      {renderTempEnableDialog()}
      {renderTempWhitelistDialog()}

      <h4>DLNA</h4>

      <Form.Group>
        <h5>
          {intl.formatMessage({ id: "status" }, { statusText: renderStatus() })}
        </h5>
      </Form.Group>

      <SettingSection headingID="actions_name">
        <Form.Group className="content">
          {renderEnableButton()}
          {renderTempCancelButton()}
        </Form.Group>

        {renderAllowedIPs()}

        <Form.Group className="content">
          <h6>
            {intl.formatMessage({ id: "config.dlna.recent_ip_addresses" })}
          </h6>
          <Form.Group>{renderRecentIPs()}</Form.Group>
          <Form.Group>
            <Button onClick={() => statusRefetch()}>
              <FormattedMessage id="actions.refresh" />
            </Button>
          </Form.Group>
        </Form.Group>
      </SettingSection>

      <DLNASettingsForm />
    </div>
  );
};
