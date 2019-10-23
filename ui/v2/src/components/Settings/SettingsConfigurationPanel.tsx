import {
  Button,
  Divider,
  FormGroup,
  H1,
  H4,
  H6,
  InputGroup,
  Spinner,
  Tag,
  Checkbox,
  HTMLSelect,
} from "@blueprintjs/core";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../core/generated-graphql";
import { StashService } from "../../core/StashService";
import { ErrorUtils } from "../../utils/errors";
import { ToastUtils } from "../../utils/toasts";
import { FolderSelect } from "../Shared/FolderSelect/FolderSelect";

interface IProps {}

export const SettingsConfigurationPanel: FunctionComponent<IProps> = (props: IProps) => {
  // Editing config state
  const [stashes, setStashes] = useState<string[]>([]);
  const [databasePath, setDatabasePath] = useState<string | undefined>(undefined);
  const [generatedPath, setGeneratedPath] = useState<string | undefined>(undefined);
  const [username, setUsername] = useState<string | undefined>(undefined);
  const [password, setPassword] = useState<string | undefined>(undefined);
  const [logFile, setLogFile] = useState<string | undefined>();
  const [logOut, setLogOut] = useState<boolean>(true);
  const [logLevel, setLogLevel] = useState<string>("Info");
  const [logAccess, setLogAccess] = useState<boolean>(true);

  const { data, error, loading } = StashService.useConfiguration();

  const updateGeneralConfig = StashService.useConfigureGeneral({
    stashes,
    databasePath,
    generatedPath,
    username,
    password,
    logFile,
    logOut,
    logLevel,
    logAccess,
  });

  useEffect(() => {
    if (!data || !data.configuration || !!error) { return; }
    const conf = StashService.nullToUndefined(data.configuration) as GQL.ConfigDataFragment;
    if (!!conf.general) {
      setStashes(conf.general.stashes || []);
      setDatabasePath(conf.general.databasePath);
      setGeneratedPath(conf.general.generatedPath);
      setUsername(conf.general.username);
      setPassword(conf.general.password);
      setLogFile(conf.general.logFile);
      setLogOut(conf.general.logOut);
      setLogLevel(conf.general.logLevel);
      setLogAccess(conf.general.logAccess);
    }
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
      {!!error ? <h1>{error.message}</h1> : undefined}
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
      <FormGroup
        label="Database Path"
        helperText="File location for the SQLite database (requires restart)"
      >
        <InputGroup value={databasePath} onChange={(e: any) => setDatabasePath(e.target.value)} />
      </FormGroup>
      <FormGroup
        label="Generated Path"
        helperText="Directory location for the generated files (scene markers, scene previews, sprites, etc)"
      >
        <InputGroup value={generatedPath} onChange={(e: any) => setGeneratedPath(e.target.value)} />
      </FormGroup>

      <Divider />
      <H4>Authentication</H4>
      <FormGroup
        label="Username"
        helperText="Username to access Stash. Leave blank to disable user authentication"
      >
        <InputGroup value={username} onChange={(e: any) => setUsername(e.target.value)} />
      </FormGroup>
      <FormGroup
        label="Password"
        helperText="Password to access Stash. Leave blank to disable user authentication"
      >
        <InputGroup type="password" value={password} onChange={(e: any) => setPassword(e.target.value)} />
      </FormGroup>

      <Divider />
      <H4>Logging</H4>
      <FormGroup
        label="Log file"
        helperText="Path to the file to output logging to. Blank to disable file logging. Requires restart."
      >
        <InputGroup value={logFile} onChange={(e: any) => setLogFile(e.target.value)} />
      </FormGroup>

      <FormGroup
        helperText="Logs to the terminal in addition to a file. Always true if file logging is disabled. Requires restart."
      >
        <Checkbox
          checked={logOut}
          label="Log to terminal"
          onChange={(e: any) => setLogOut(e.target.value)}
        />
      </FormGroup>

      <FormGroup inline={true} label="Log Level">
        <HTMLSelect
          options={["Debug", "Info", "Warning", "Error"]}
          onChange={(event) => setLogLevel(event.target.value)}
          value={logLevel}
        />
      </FormGroup>

      <FormGroup
        helperText="Logs http access to the terminal. Requires restart."
      >
        <Checkbox
          checked={logAccess}
          label="Log http access"
          onChange={(e: any) => setLogAccess(e.target.value)}
        />
      </FormGroup>

      <Divider />
      <Button intent="primary" onClick={() => onSave()}>Save</Button>
    </>
  );
};
