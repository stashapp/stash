import React, { useEffect, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { DurationInput, LoadingIndicator } from "src/components/Shared";
import { StashService } from "src/core/StashService";
import { useToast } from "src/hooks";

export const SettingsInterfacePanel: React.FC = () => {
  const Toast = useToast();
  const config = StashService.useConfiguration();
  const [soundOnPreview, setSoundOnPreview] = useState<boolean>(true);
  const [wallShowTitle, setWallShowTitle] = useState<boolean>(true);
  const [maximumLoopDuration, setMaximumLoopDuration] = useState<number>(0);
  const [autostartVideo, setAutostartVideo] = useState<boolean>(false);
  const [showStudioAsText, setShowStudioAsText] = useState<boolean>(false);
  const [css, setCSS] = useState<string>();
  const [cssEnabled, setCSSEnabled] = useState<boolean>(false);

  const [updateInterfaceConfig] = StashService.useConfigureInterface({
    soundOnPreview,
    wallShowTitle,
    maximumLoopDuration,
    autostartVideo,
    showStudioAsText,
    css,
    cssEnabled
  });

  useEffect(() => {
    if (config.error) return;

    const iCfg = config?.data?.configuration?.interface;
    setSoundOnPreview(iCfg?.soundOnPreview ?? true);
    setWallShowTitle(iCfg?.wallShowTitle ?? true);
    setMaximumLoopDuration(iCfg?.maximumLoopDuration ?? 0);
    setAutostartVideo(iCfg?.autostartVideo ?? false);
    setShowStudioAsText(iCfg?.showStudioAsText ?? false);
    setCSS(iCfg?.css ?? "");
    setCSSEnabled(iCfg?.cssEnabled ?? false);
  }, [config]);

  async function onSave() {
    try {
      const result = await updateInterfaceConfig();
      // eslint-disable-next-line no-console
      console.log(result);
      Toast.success({ content: "Updated config" });
    } catch (e) {
      Toast.error(e);
    }
  }

  return (
    <>
      {config.error ? <h1>{config.error.message}</h1> : ""}
      {!config?.data?.configuration || config.loading ? (
        <LoadingIndicator />
      ) : (
        ""
      )}
      <h4>User Interface</h4>
      <Form.Group>
        <Form.Label>Scene / Marker Wall</Form.Label>
        <Form.Check
          id="wall-show-title"
          checked={wallShowTitle}
          label="Display title and tags"
          onChange={() => setWallShowTitle(!wallShowTitle)}
        />
        <Form.Check
          id="wall-sound-enabled"
          checked={soundOnPreview}
          label="Enable sound"
          onChange={() => setSoundOnPreview(!soundOnPreview)}
        />
        <Form.Text className="text-muted">
          Configuration for wall items
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <h5>Scene List</h5>
        <Form.Check
          id="show-text-studios"
          checked={showStudioAsText}
          label="Show Studios as text"
          onChange={() => {
            setShowStudioAsText(!showStudioAsText);
          }}
        />
      </Form.Group>

      <Form.Group>
        <h5>Scene Player</h5>
        <Form.Check
          id="auto-start-video"
          checked={autostartVideo}
          label="Auto-start video"
          onChange={() => {
            setAutostartVideo(!autostartVideo);
          }}
        />

        <Form.Group id="max-loop-duration">
          <Form.Label>Maximum loop duration</Form.Label>
          <DurationInput
            className="col col-sm-4"
            numericValue={maximumLoopDuration}
            onValueChange={duration => setMaximumLoopDuration(duration)}
          />
          <Form.Text className="text-muted">
            Maximum scene duration where scene player will loop the video - 0 to
            disable
          </Form.Text>
        </Form.Group>
      </Form.Group>

      <Form.Group>
        <h5>Custom CSS</h5>
        <Form.Check
          id="custom-css"
          checked={cssEnabled}
          label="Custom CSS enabled"
          onChange={() => {
            setCSSEnabled(!cssEnabled);
          }}
        />

        <Form.Control
          as="textarea"
          value={css}
          onChange={(e: any) => setCSS(e.target.value)}
          rows={16}
          className="col col-sm-6"
        ></Form.Control>
        <Form.Text className="text-muted">
          Page must be reloaded for changes to take effect.
        </Form.Text>
      </Form.Group>

      <hr />
      <Button variant="primary" onClick={() => onSave()}>
        Save
      </Button>
    </>
  );
};
