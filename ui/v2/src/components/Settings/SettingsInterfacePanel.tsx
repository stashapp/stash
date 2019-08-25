import {
  Button,
  Checkbox,
  Divider,
  FormGroup,
  H4,
  Spinner,
  TextArea
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent, useEffect, useState } from "react";
import { useInterfaceLocalForage } from "../../hooks/LocalForage";
import { StashService } from "../../core/StashService";
import { ErrorUtils } from "../../utils/errors";
import { ToastUtils } from "../../utils/toasts";

interface IProps {}

export const SettingsInterfacePanel: FunctionComponent<IProps> = () => {
  const {data, setData} = useInterfaceLocalForage();
  const config = StashService.useConfiguration();
  const [css, setCSS] = useState<string>();
  const [cssEnabled, setCSSEnabled] = useState<boolean>();

  const updateInterfaceConfig = StashService.useConfigureInterface({
    css,
    cssEnabled
  });

  useEffect(() => {
    if (!config.data || !config.data.configuration || !!config.error) { return; }
    if (!!config.data.configuration.interface) {
      setCSS(config.data.configuration.interface.css || "");
      setCSSEnabled(config.data.configuration.interface.cssEnabled || false);
    }
  }, [config.data]);

  async function onSave() {
    try {
      const result = await updateInterfaceConfig();
      console.log(result);
      ToastUtils.success("Updated config");
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  return (
    <>
      {!!config.error ? <h1>{config.error.message}</h1> : undefined}
      {(!config.data || !config.data.configuration || config.loading) ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
      <H4>User Interface</H4>
      <FormGroup
        label="Scene / Marker Wall"
        helperText="Configuration for wall items"
      >
        <Checkbox
          checked={!!data ? data.wall.textContainerEnabled : true}
          label="Display title and tags"
          onChange={() => {
            if (!data) { return; }
            const newSettings = _.cloneDeep(data);
            newSettings.wall.textContainerEnabled = !data.wall.textContainerEnabled;
            setData(newSettings);
          }}
        />
        <Checkbox
          checked={!!data ? data.wall.soundEnabled : true}
          label="Enable sound"
          onChange={() => {
            if (!data) { return; }
            const newSettings = _.cloneDeep(data);
            newSettings.wall.soundEnabled = !data.wall.soundEnabled;
            setData(newSettings);
          }}
        />
      </FormGroup>

      <FormGroup
        label="Custom CSS"
        helperText="Page must be reloaded for changes to take effect."
      >
        <Checkbox
          checked={cssEnabled}
          label="Custom CSS enabled"
          onChange={() => {
            setCSSEnabled(!cssEnabled)
          }}
        />

        <TextArea 
          value={css} 
          onChange={(e: any) => setCSS(e.target.value)}
          fill={true}
          rows={16}>
        </TextArea>

        <Divider />
        <Button intent="primary" onClick={() => onSave()}>Save</Button>
      </FormGroup>
    </>
  );
};
