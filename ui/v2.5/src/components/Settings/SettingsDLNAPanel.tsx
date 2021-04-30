import React, { useEffect, useState } from "react";
import { Button, Form } from "react-bootstrap";
import {
  useConfiguration,
  useConfigureDLNA,
  useDLNAStatus,
} from "src/core/StashService";
import { useToast } from "src/hooks";

export const SettingsDLNAPanel: React.FC = () => {
  const Toast = useToast();

  const [enabled, setEnabled] = useState<boolean>(false);

  const { data } = useConfiguration();
  const status = useDLNAStatus();

  const [updateDLNAConfig] = useConfigureDLNA({
    dlnaEnabled: enabled,
  });

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

  function renderStatus() {
    if (!status.data?.dlnaStatus) {
      return "";
    }

    const { dlnaStatus } = status.data;
    return dlnaStatus.running ? "running" : "not running";
  }

  return (
    <>
      <h4>DLNA</h4>

      <Form.Group>
        <h5>Status: {renderStatus()}</h5>
      </Form.Group>

      <Form.Group>
        <h5>Actions</h5>

        {/* TODO - temporarily enable/disable */}
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
