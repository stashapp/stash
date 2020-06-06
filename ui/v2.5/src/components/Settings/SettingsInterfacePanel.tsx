import React, { useEffect, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { DurationInput, LoadingIndicator } from "src/components/Shared";
import { useConfiguration, useConfigureInterface } from "src/core/StashService";
import { useToast } from "src/hooks";

export const SettingsInterfacePanel: React.FC = () => {
  const Toast = useToast();
  const { data: config, error, loading } = useConfiguration();
  const [soundOnPreview, setSoundOnPreview] = useState<boolean>(true);
  const [wallShowTitle, setWallShowTitle] = useState<boolean>(true);
  const [wallPlayback, setWallPlayback] = useState<string>("video");
  const [maximumLoopDuration, setMaximumLoopDuration] = useState<number>(0);
  const [autostartVideo, setAutostartVideo] = useState<boolean>(false);
  const [showStudioAsText, setShowStudioAsText] = useState<boolean>(false);
  const [css, setCSS] = useState<string>();
  const [cssEnabled, setCSSEnabled] = useState<boolean>(false);
  const [language, setLanguage] = useState<string>("en");

  const [updateInterfaceConfig] = useConfigureInterface({
    soundOnPreview,
    wallShowTitle,
    wallPlayback,
    maximumLoopDuration,
    autostartVideo,
    showStudioAsText,
    css,
    cssEnabled,
    language,
  });

  useEffect(() => {
    const iCfg = config?.configuration?.interface;
    setSoundOnPreview(iCfg?.soundOnPreview ?? true);
    setWallShowTitle(iCfg?.wallShowTitle ?? true);
    setWallPlayback(iCfg?.wallPlayback ?? "video");
    setMaximumLoopDuration(iCfg?.maximumLoopDuration ?? 0);
    setAutostartVideo(iCfg?.autostartVideo ?? false);
    setShowStudioAsText(iCfg?.showStudioAsText ?? false);
    setCSS(iCfg?.css ?? "");
    setCSSEnabled(iCfg?.cssEnabled ?? false);
    setLanguage(iCfg?.language ?? "en-US");
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

  if (error) return <h1>{error.message}</h1>;
  if (loading) return <LoadingIndicator />;

  return (
    <>
      <h4>User Interface</h4>
      <Form.Group controlId="language">
        <h6>Language</h6>
        <Form.Control
          as="select"
          className="col-4 input-control"
          value={language}
          onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
            setLanguage(e.currentTarget.value)
          }
        >
          <option value="en-US">English (United States)</option>
          <option value="en-GB">English (United Kingdom)</option>
        </Form.Control>
      </Form.Group>
      <Form.Group>
        <h5>Scene / Marker Wall</h5>
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
        <Form.Label htmlFor="wall-preview">
          <h6>Preview Type</h6>
        </Form.Label>
        <Form.Control
          as="select"
          name="wall-preview"
          className="col-4 input-control"
          value={wallPlayback}
          onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
            setWallPlayback(e.currentTarget.value)
          }
        >
          <option value="video">Video</option>
          <option value="animation">Animated Image</option>
          <option value="image">Static Image</option>
        </Form.Control>
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
        <Form.Group id="auto-start-video">
          <Form.Check
            checked={autostartVideo}
            label="Auto-start video"
            onChange={() => {
              setAutostartVideo(!autostartVideo);
            }}
          />
        </Form.Group>

        <Form.Group id="max-loop-duration">
          <h6>Maximum loop duration</h6>
          <DurationInput
            className="row col col-4"
            numericValue={maximumLoopDuration}
            onValueChange={(duration) => setMaximumLoopDuration(duration ?? 0)}
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
          onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
            setCSS(e.currentTarget.value)
          }
          rows={16}
          className="col col-sm-6 text-input code"
        />
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
