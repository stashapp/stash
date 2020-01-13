import {
  Button,
  Checkbox,
  Divider,
  FormGroup,
  H4,
  Spinner,
  TextArea,
  NumericInput
} from "@blueprintjs/core";
import React, { FunctionComponent, useEffect, useState } from "react";
import { StashService } from "../../core/StashService";
import { ErrorUtils } from "../../utils/errors";
import { ToastUtils } from "../../utils/toasts";

interface IProps {}

export const SettingsInterfacePanel: FunctionComponent<IProps> = () => {
  const config = StashService.useConfiguration();
  const [soundOnPreview, setSoundOnPreview] = useState<boolean>();
  const [wallShowTitle, setWallShowTitle] = useState<boolean>();
  const [maximumLoopDuration, setMaximumLoopDuration] = useState<number>(0);
  const [autostartVideo, setAutostartVideo] = useState<boolean>();
  const [showStudioAsText, setShowStudioAsText] = useState<boolean>();
  const [css, setCSS] = useState<string>();
  const [cssEnabled, setCSSEnabled] = useState<boolean>();

  const updateInterfaceConfig = StashService.useConfigureInterface({
    soundOnPreview,
    wallShowTitle,
    maximumLoopDuration,
    autostartVideo,
    showStudioAsText,
    css,
    cssEnabled
  });

  useEffect(() => {
    if (!config.data || !config.data.configuration || !!config.error) { return; }
    if (!!config.data.configuration.interface) {
      let iCfg = config.data.configuration.interface;
      setSoundOnPreview(iCfg.soundOnPreview !== undefined ? iCfg.soundOnPreview : true);
      setWallShowTitle(iCfg.wallShowTitle !== undefined ? iCfg.wallShowTitle : true);
      setMaximumLoopDuration(iCfg.maximumLoopDuration || 0);
      setAutostartVideo(iCfg.autostartVideo !== undefined ? iCfg.autostartVideo : false);
      setShowStudioAsText(iCfg.showStudioAsText !== undefined ? iCfg.showStudioAsText : false);
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
          checked={wallShowTitle}
          label="Display title and tags"
          onChange={() => setWallShowTitle(!wallShowTitle)}
        />
        <Checkbox
          checked={soundOnPreview}
          label="Enable sound"
          onChange={() => setSoundOnPreview(!soundOnPreview)}
        />
      </FormGroup>

      <FormGroup
        label="Scene List"
      >
        <Checkbox
          checked={showStudioAsText}
          label="Show Studios as text"
          onChange={() => {
            setShowStudioAsText(!showStudioAsText)
          }}
        />
      </FormGroup>
      
      <FormGroup
        label="Scene Player"
      >
        <Checkbox
          checked={autostartVideo}
          label="Auto-start video"
          onChange={() => {
            setAutostartVideo(!autostartVideo)
          }}
        />

        <FormGroup
          label="Maximum loop duration"
          helperText="Maximum scene duration - in seconds - where scene player will loop the video - 0 to disable"
        >
          <NumericInput 
            value={maximumLoopDuration} 
            type="number"
            onValueChange={(value: number) => setMaximumLoopDuration(value)}
            min={0}
            minorStepSize={1}
          />
        </FormGroup>
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
      </FormGroup>

      <Divider />
      <Button intent="primary" onClick={() => onSave()}>Save</Button>
    </>
  );
};
