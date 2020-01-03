import {
  AnchorButton,
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
  const [maxTranscodeSize, setMaxTranscodeSize] = useState<GQL.StreamingResolutionEnum | undefined>(undefined);
  const [maxStreamingTranscodeSize, setMaxStreamingTranscodeSize] = useState<GQL.StreamingResolutionEnum | undefined>(undefined);
  const [username, setUsername] = useState<string | undefined>(undefined);
  const [password, setPassword] = useState<string | undefined>(undefined);
  const [logFile, setLogFile] = useState<string | undefined>();
  const [logOut, setLogOut] = useState<boolean>(true);
  const [logLevel, setLogLevel] = useState<string>("Info");
  const [logAccess, setLogAccess] = useState<boolean>(true);
  const [excludes, setExcludes] = useState<(string)[]>([]);

  const { data, error, loading } = StashService.useConfiguration();

  const updateGeneralConfig = StashService.useConfigureGeneral({
    stashes,
    databasePath,
    generatedPath,
    maxTranscodeSize,
    maxStreamingTranscodeSize,
    username,
    password,
    logFile,
    logOut,
    logLevel,
    logAccess,
    excludes,

  });

  useEffect(() => {
    if (!data || !data.configuration || !!error) { return; }
    const conf = StashService.nullToUndefined(data.configuration) as GQL.ConfigDataFragment;
    if (!!conf.general) {
      setStashes(conf.general.stashes || []);
      setDatabasePath(conf.general.databasePath);
      setGeneratedPath(conf.general.generatedPath);
      setMaxTranscodeSize(conf.general.maxTranscodeSize);
      setMaxStreamingTranscodeSize(conf.general.maxStreamingTranscodeSize);
      setUsername(conf.general.username);
      setPassword(conf.general.password);
      setLogFile(conf.general.logFile);
      setLogOut(conf.general.logOut);
      setLogLevel(conf.general.logLevel);
      setLogAccess(conf.general.logAccess);
      setExcludes(conf.general.excludes);
    }
  }, [data]);

  function onStashesChanged(directories: string[]) {
    setStashes(directories);
  }

  function excludeRegexChanged(idx: number, value: string) {
    const newExcludes = excludes.map((regex, i)=> {
      const ret = ( idx !== i ) ? regex : value ;
      return ret
      })
    setExcludes(newExcludes);
  }

  function excludeRemoveRegex(idx: number) {
    const newExcludes = excludes.filter((regex, i) => i!== idx );

    setExcludes(newExcludes);
  }

  function excludeAddRegex() {
    const demo = "sample\\.mp4$"
    const newExcludes = excludes.concat(demo);

    setExcludes(newExcludes);
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

  const transcodeQualities = [
    GQL.StreamingResolutionEnum.Low,
    GQL.StreamingResolutionEnum.Standard,
    GQL.StreamingResolutionEnum.StandardHd,
    GQL.StreamingResolutionEnum.FullHd,
    GQL.StreamingResolutionEnum.FourK,
    GQL.StreamingResolutionEnum.Original
  ].map(resolutionToString);

  function resolutionToString(r : GQL.StreamingResolutionEnum | undefined) {
    switch (r) {
      case GQL.StreamingResolutionEnum.Low: return "240p";
      case GQL.StreamingResolutionEnum.Standard: return "480p";
      case GQL.StreamingResolutionEnum.StandardHd: return "720p";
      case GQL.StreamingResolutionEnum.FullHd: return "1080p";
      case GQL.StreamingResolutionEnum.FourK: return "4k";
      case GQL.StreamingResolutionEnum.Original: return "Original";
    }

    return "Original";
  }

  function translateQuality(quality : string) {
    switch (quality) {
      case "240p": return GQL.StreamingResolutionEnum.Low;
      case "480p": return GQL.StreamingResolutionEnum.Standard;
      case "720p": return GQL.StreamingResolutionEnum.StandardHd;
      case "1080p": return GQL.StreamingResolutionEnum.FullHd;
      case "4k": return GQL.StreamingResolutionEnum.FourK;
      case "Original": return GQL.StreamingResolutionEnum.Original;
    }

    return GQL.StreamingResolutionEnum.Original;
  }

  return (
    <>
      {!!error ? <h1>{error.message}</h1> : undefined}
      {(!data || !data.configuration || loading) ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
      <H4>Library</H4>
      <FormGroup>
        <FormGroup>
          <FormGroup
            label="Stashes"
            helperText="Directory locations to your content"
          >
            <FolderSelect
              directories={stashes}
              onDirectoriesChanged={onStashesChanged}
            />
          </FormGroup>
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

        <FormGroup
          label="Excluded Patterns"
        >

       { (excludes) ? excludes.map((regexp, i) => {
         return(
           <InputGroup
             value={regexp}
             onChange={(e: any) => excludeRegexChanged(i, e.target.value)}
             rightElement={<Button icon="minus" minimal={true} intent="danger" onClick={(e: any) => excludeRemoveRegex(i)} />}
           />
           );
         }) : null
       }

          <Button icon="plus" minimal={true} onClick={(e: any) => excludeAddRegex()} />
          <div>
            <p>
              <AnchorButton
                href="https://github.com/stashapp/stash/wiki/Exclude-file-configuration"
                rightIcon="help"
                text="Regexps of files/paths to exclude from Scan and add to Clean"
                minimal={true}
                target="_blank" 
              />
            </p>
          </div>
        </FormGroup>
      </FormGroup>
      
      <Divider />
        <FormGroup>
          <H4>Video</H4>
          <FormGroup 
            label="Maximum transcode size"
            helperText="Maximum size for generated transcodes"
          >
            <HTMLSelect
              options={transcodeQualities}
              onChange={(event) => setMaxTranscodeSize(translateQuality(event.target.value))}
              value={resolutionToString(maxTranscodeSize)}
            />
          </FormGroup>
          <FormGroup 
            label="Maximum streaming transcode size"
            helperText="Maximum size for transcoded streams"
          >
            <HTMLSelect
              options={transcodeQualities}
              onChange={(event) => setMaxStreamingTranscodeSize(translateQuality(event.target.value))}
              value={resolutionToString(maxStreamingTranscodeSize)}
            />
          </FormGroup>
        </FormGroup>
      <Divider />

      <FormGroup>
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
          onChange={() => setLogOut(!logOut)}
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
          onChange={() => setLogAccess(!logAccess)}
        />
      </FormGroup>

      <Divider />
      <Button intent="primary" onClick={() => onSave()}>Save</Button>
    </>
  );
};
