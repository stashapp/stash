import {
  Checkbox,
  FormGroup,
  H4,
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent } from "react";
import { useInterfaceLocalForage } from "../../hooks/LocalForage";

interface IProps {}

export const SettingsInterfacePanel: FunctionComponent<IProps> = () => {
  const {data, setData} = useInterfaceLocalForage();

  return (
    <>
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
    </>
  );
};
