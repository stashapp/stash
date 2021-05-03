import React, { useEffect, useState } from "react";
import { Button, Form } from "react-bootstrap";
import {
  useConfiguration,
  useConfigureDLNA,
  useDisableDLNA,
  useDLNAStatus,
  useEnableDLNA,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { DurationInput, Modal } from "../Shared";

export const SettingsDLNAPanel: React.FC = () => {
  const Toast = useToast();

  const [enabled, setEnabled] = useState<boolean>(false);

  // undefined to hide dialog, true for enable, false for disable
  const [enableDisable, setEnableDisable] = useState<boolean | undefined>(
    undefined
  );

  const [enableUntilRestart, setEnableUntilRestart] = useState<boolean>(false);
  const [enableDuration, setEnableDuration] = useState<number | undefined>(
    undefined
  );

  const { data } = useConfiguration();
  const status = useDLNAStatus();

  const [updateDLNAConfig] = useConfigureDLNA({
    dlnaEnabled: enabled,
  });

  const [enableDLNA] = useEnableDLNA();

  const [disableDLNA] = useDisableDLNA();

  useEffect(() => {
    if (data?.configuration.dlna) {
      const { dlnaEnabled } = data.configuration.dlna;
      setEnabled(dlnaEnabled);
      // TODO - whitelist
    }
  }, [data]);

  async function onSave() {
    try {
      await updateDLNAConfig();
      Toast.success({ content: "Updated config" });
    } catch (e) {
      Toast.error(e);
    } finally {
      status.refetch();
    }
  }

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
      status.refetch();
    }
  }

  function renderStatus() {
    if (!status.data?.dlnaStatus) {
      return "";
    }

    const { dlnaStatus } = status.data;
    const runningText = dlnaStatus.running ? "running" : "not running";
    let untilText = "";

    if (dlnaStatus.until) {
      const deadline = new Date(dlnaStatus.until);
      untilText = `until ${deadline.toLocaleString()}`;
    }

    return `${runningText} ${untilText}`;
  }

  function renderEnableButton() {
    if (!data?.configuration.dlna) {
      return;
    }

    // if enabled by default, then show the disable temporarily
    // if disabled by default, then show enable temporarily
    // TODO - also show a cancel button
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
    if (!status.data || !data) {
      return false;
    }

    const { dlnaStatus } = status.data;
    const { dlnaEnabled } = data.configuration.dlna;

    return dlnaStatus.until || dlnaStatus.running !== dlnaEnabled;
  }

  async function cancelTempBehaviour() {
    if (!canCancel()) {
      return;
    }

    const running = status.data?.dlnaStatus.running;

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
      status.refetch();
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

  return (
    <>
      {renderTempEnableDialog()}

      <h4>DLNA</h4>

      <Form.Group>
        <h5>Status: {renderStatus()}</h5>
      </Form.Group>

      <Form.Group>
        <h5>Actions</h5>

        {renderEnableButton()}
        {renderTempCancelButton()}
      </Form.Group>

      <Form.Group>
        <h5>Settings</h5>
        <Form.Check
          checked={enabled}
          label="Enabled by default"
          onChange={() => setEnabled(!enabled)}
        />

        {/* TODO - default IP whitelist */}
      </Form.Group>

      <hr />

      <Button variant="primary" onClick={() => onSave()}>
        Save
      </Button>
    </>
  );
};
