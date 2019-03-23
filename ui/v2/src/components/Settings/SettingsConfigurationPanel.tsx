import {
  Button,
  Divider,
  FormGroup,
  H1,
  H4,
  H6,
  Spinner,
  Tag,
} from "@blueprintjs/core";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { ErrorUtils } from "../../utils/errors";
import { TextUtils } from "../../utils/text";
import { ToastUtils } from "../../utils/toasts";
import { FolderSelect } from "../Shared/FolderSelect/FolderSelect";

interface IProps {}

export const SettingsConfigurationPanel: FunctionComponent<IProps> = (props: IProps) => {
  // Editing config state
  const [stashes, setStashes] = useState<string[]>([]);

  // const [config, setConfig] = useState<Partial<GQL.ConfigDataFragment>>({});
  const { data, error, loading } = StashService.useConfiguration();

  const updateGeneralConfig = StashService.useConfigureGeneral({
    stashes,
  });

  useEffect(() => {
    if (!data || !data.configuration || !!error) { return; }
    const conf = StashService.nullToUndefined(data.configuration) as GQL.ConfigDataFragment;
    if (!!conf.general) {
      setStashes(conf.general.stashes || []);
    }
    // setConfig(conf);
  }, [data]);

  function onStashesChanged(directories: string[]) {
    setStashes(directories);
  }

  async function onSave() {
    try {
      const result = await updateGeneralConfig();
      console.log(result);
      ToastUtils.success("Updated config");
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  return (
    <>
      {!!error ? error : undefined}
      {(!data || !data.configuration || loading) ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
      <H4>Library</H4>
      <FormGroup
        label="Stashes"
        helperText="Directory locations to your content"
      >
        <FolderSelect
          directories={stashes}
          onDirectoriesChanged={onStashesChanged}
        />
      </FormGroup>
      <Divider />
      <Button intent="primary" onClick={() => onSave()}>Save</Button>
    </>
  );
};
