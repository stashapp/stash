import React, { useEffect, useState } from "react";
import { useFormik } from "formik";
import { Button, Form } from "react-bootstrap";
import { Prompt } from "react-router";
import * as yup from "yup";
import {
  useConfiguration,
  useConfigureDLNA,
  useDisableDLNA,
  useDLNAStatus,
  useEnableDLNA,
  useAddTempDLNAIP,
  useRemoveTempDLNAIP,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { DurationInput, Icon, LoadingIndicator, Modal } from "../Shared";
import { StringListInput } from "../Shared/StringListInput";

export const SettingsDLNAPanel: React.FC = () => {
  const Toast = useToast();

  // settings
  const schema = yup.object({
    dlnaEnabled: yup.boolean().required(),
    dlnaWhitelistedIPs: yup.array(yup.string().required()).required(),
  });

  interface IConfigValues {
    dlnaEnabled: boolean;
    dlnaWhitelistedIPs: string[];
  }

  const initialValues: IConfigValues = {
    dlnaEnabled: false,
    dlnaWhitelistedIPs: [],
  };

  // undefined to hide dialog, true for enable, false for disable
  const [enableDisable, setEnableDisable] = useState<boolean | undefined>(
    undefined
  );

  const [enableUntilRestart, setEnableUntilRestart] = useState<boolean>(false);
  const [enableDuration, setEnableDuration] = useState<number | undefined>(
    undefined
  );

  const [ipEntry, setIPEntry] = useState<string>("");
  const [tempIP, setTempIP] = useState<string | undefined>();

  const { data } = useConfiguration();
  const { data: statusData, loading, refetch: statusRefetch } = useDLNAStatus();

  const [updateDLNAConfig] = useConfigureDLNA();

  const [enableDLNA] = useEnableDLNA();
  const [disableDLNA] = useDisableDLNA();
  const [addTempDLANIP] = useAddTempDLNAIP();
  const [removeTempDLNAIP] = useRemoveTempDLNAIP();

  async function onSave(input: IConfigValues) {
    try {
      await updateDLNAConfig({
        variables: {
          input,
        },
      });
      Toast.success({ content: "Updated config" });
    } catch (e) {
      Toast.error(e);
    } finally {
      statusRefetch();
    }
  }

  const formik = useFormik({
    initialValues,
    validationSchema: schema,
    onSubmit: (values) => onSave(values),
  });

  useEffect(() => {
    if (data?.configuration.dlna) {
      const { dlnaEnabled, dlnaWhitelistedIPs } = data.configuration.dlna;
      formik.setFieldValue("dlnaEnabled", dlnaEnabled);
      formik.setFieldValue("dlnaWhitelistedIPs", dlnaWhitelistedIPs);
    }
  }, [data, formik]);

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
        Toast.success({ content: "Enabled DLNA temporarily" });
      } else {
        await disableDLNA(input);
        Toast.success({ content: "Disabled DLNA temporarily" });
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
      Toast.success({ content: "Allowed IP temporarily" });
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
      Toast.success({ content: "Disallowed IP" });
    } catch (e) {
      Toast.error(e);
    } finally {
      statusRefetch();
    }
  }

  function renderDeadline(until?: string) {
    if (until) {
      const deadline = new Date(until);
      return `until ${deadline.toLocaleString()}`;
    }

    return "";
  }

  function renderStatus() {
    if (!statusData) {
      return "";
    }

    const { dlnaStatus } = statusData;
    const runningText = dlnaStatus.running ? "running" : "not running";

    return `${runningText} ${renderDeadline(dlnaStatus.until)}`;
  }

  function renderEnableButton() {
    if (!data?.configuration.dlna) {
      return;
    }

    // if enabled by default, then show the disable temporarily
    // if disabled by default, then show enable temporarily
    if (data?.configuration.dlna.dlnaEnabled) {
      return (
        <Button onClick={() => setEnableDisable(false)} className="mr-1">
          Disable temporarily...
        </Button>
      );
    }

    return (
      <Button onClick={() => setEnableDisable(true)} className="mr-1">
        Enable temporarily...
      </Button>
    );
  }

  function canCancel() {
    if (!statusData || !data) {
      return false;
    }

    const { dlnaStatus } = statusData;
    const { dlnaEnabled } = data.configuration.dlna;

    return dlnaStatus.until || dlnaStatus.running !== dlnaEnabled;
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
      Toast.success({ content: "Successfully cancelled temporary behaviour" });
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
      <Modal
        show={enableDisable !== undefined}
        header={capitalised}
        icon="clock"
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
            label="until restart"
            onChange={() => setEnableUntilRestart(!enableUntilRestart)}
          />
        </Form.Group>

        <Form.Group id="temp-enable-duration">
          <DurationInput
            numericValue={enableDuration ?? 0}
            onValueChange={(v) => setEnableDuration(v ?? 0)}
            disabled={enableUntilRestart}
          />
          <Form.Text className="text-muted">
            Duration to {text} for - in minutes.
          </Form.Text>
        </Form.Group>
      </Modal>
    );
  }

  function renderTempWhitelistDialog() {
    return (
      <Modal
        show={tempIP !== undefined}
        header={`Allow ${tempIP}`}
        icon="clock"
        accept={{
          text: "Allow",
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
            label="until restart"
            onChange={() => setEnableUntilRestart(!enableUntilRestart)}
          />
        </Form.Group>

        <Form.Group id="temp-enable-duration">
          <DurationInput
            numericValue={enableDuration ?? 0}
            onValueChange={(v) => setEnableDuration(v ?? 0)}
            disabled={enableUntilRestart}
          />
          <Form.Text className="text-muted">
            Duration to allow for - in minutes.
          </Form.Text>
        </Form.Group>
      </Modal>
    );
  }

  function renderAllowedIPs() {
    if (!statusData || statusData.dlnaStatus.allowedIPAddresses.length === 0) {
      return;
    }

    const { allowedIPAddresses } = statusData.dlnaStatus;
    return (
      <Form.Group>
        <h6>Allowed IP addresses</h6>

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
                  title="Disallow"
                  variant="danger"
                  onClick={() => onDisallowTempIP(a.ipAddress)}
                >
                  <Icon icon="times" />
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
                title="Allow temporarily"
                onClick={() => setTempIP(a)}
              >
                <Icon icon="user-clock" />
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
              title="Allow temporarily"
              onClick={() => setTempIP(ipEntry)}
              disabled={!ipEntry}
            >
              <Icon icon="user-clock" />
            </Button>
          </div>
        </li>
      </ul>
    );
  }

  if (loading) return <LoadingIndicator />;

  return (
    <div id="settings-dlna">
      {renderTempEnableDialog()}
      {renderTempWhitelistDialog()}

      <Prompt
        when={formik.dirty}
        message="Unsaved changes. Are you sure you want to leave?"
      />

      <h4>DLNA</h4>

      <Form.Group>
        <h5>Status: {renderStatus()}</h5>
      </Form.Group>

      <Form.Group>
        <h5>Actions</h5>

        <Form.Group>
          {renderEnableButton()}
          {renderTempCancelButton()}
        </Form.Group>

        {renderAllowedIPs()}

        <Form.Group>
          <h6>Recent IP addresses</h6>
          <Form.Group>{renderRecentIPs()}</Form.Group>
          <Form.Group>
            <Button onClick={() => statusRefetch()}>Refresh</Button>
          </Form.Group>
        </Form.Group>
      </Form.Group>

      <Form noValidate onSubmit={formik.handleSubmit}>
        <Form.Group>
          <h5>Settings</h5>
          <Form.Group>
            <Form.Check
              checked={formik.values.dlnaEnabled}
              label="Enabled by default"
              onChange={() =>
                formik.setFieldValue("dlnaEnabled", !formik.values.dlnaEnabled)
              }
            />
          </Form.Group>

          <Form.Group>
            <h6>Default IP Whitelist</h6>
            <StringListInput
              value={formik.values.dlnaWhitelistedIPs}
              setValue={(value) =>
                formik.setFieldValue("dlnaWhitelistedIPs", value)
              }
              defaultNewValue="*"
              className="ip-whitelist-input"
            />
            <Form.Text className="text-muted">
              Default IP addresses allow to access DLNA. Use <code>*</code> to
              allow all IP addresses.
            </Form.Text>
          </Form.Group>
        </Form.Group>

        <hr />

        <Button variant="primary" onClick={() => formik.submitForm()}>
          Save
        </Button>
      </Form>
    </div>
  );
};
