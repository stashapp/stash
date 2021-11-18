import React, { useEffect, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { DurationInput, LoadingIndicator } from "src/components/Shared";
import {
  useConfiguration,
  useConfigureDefaults,
  useConfigureInterface,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { CheckboxGroup } from "./CheckboxGroup";
import { withoutTypename } from "src/utils";

const allMenuItems = [
  { id: "scenes", label: "Scenes" },
  { id: "images", label: "Images" },
  { id: "movies", label: "Movies" },
  { id: "markers", label: "Markers" },
  { id: "galleries", label: "Galleries" },
  { id: "performers", label: "Performers" },
  { id: "studios", label: "Studios" },
  { id: "tags", label: "Tags" },
];

const SECONDS_TO_MS = 1000;

export const SettingsInterfacePanel: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();
  const { data: config, error, loading } = useConfiguration();
  const [menuItemIds, setMenuItemIds] = useState<string[]>(
    allMenuItems.map((item) => item.id)
  );
  const [noBrowser, setNoBrowserFlag] = useState<boolean>(false);
  const [soundOnPreview, setSoundOnPreview] = useState<boolean>(true);
  const [wallShowTitle, setWallShowTitle] = useState<boolean>(true);
  const [wallPlayback, setWallPlayback] = useState<string>("video");
  const [maximumLoopDuration, setMaximumLoopDuration] = useState<number>(0);
  const [autostartVideo, setAutostartVideo] = useState<boolean>(false);
  const [
    autostartVideoOnPlaySelected,
    setAutostartVideoOnPlaySelected,
  ] = useState(true);
  const [continuePlaylistDefault, setContinuePlaylistDefault] = useState(false);
  const [slideshowDelay, setSlideshowDelay] = useState<number>(0);
  const [showStudioAsText, setShowStudioAsText] = useState<boolean>(false);
  const [css, setCSS] = useState<string>();
  const [cssEnabled, setCSSEnabled] = useState<boolean>(false);
  const [language, setLanguage] = useState<string>("en");
  const [handyKey, setHandyKey] = useState<string>();
  const [funscriptOffset, setFunscriptOffset] = useState<number>(0);
  const [deleteFileDefault, setDeleteFileDefault] = useState<boolean>(false);
  const [deleteGeneratedDefault, setDeleteGeneratedDefault] = useState<boolean>(
    true
  );
  const [
    disableDropdownCreate,
    setDisableDropdownCreate,
  ] = useState<GQL.ConfigDisableDropdownCreateInput>({});

  const [updateInterfaceConfig] = useConfigureInterface({
    menuItems: menuItemIds,
    soundOnPreview,
    wallShowTitle,
    wallPlayback,
    maximumLoopDuration,
    noBrowser,
    autostartVideo,
    autostartVideoOnPlaySelected,
    continuePlaylistDefault,
    showStudioAsText,
    css,
    cssEnabled,
    language,
    slideshowDelay,
    handyKey,
    funscriptOffset,
    disableDropdownCreate,
  });

  const [updateDefaultsConfig] = useConfigureDefaults();

  useEffect(() => {
    if (config) {
      const { interface: iCfg, defaults } = config.configuration;
      setMenuItemIds(iCfg.menuItems ?? allMenuItems.map((item) => item.id));
      setSoundOnPreview(iCfg.soundOnPreview ?? true);
      setWallShowTitle(iCfg.wallShowTitle ?? true);
      setWallPlayback(iCfg.wallPlayback ?? "video");
      setMaximumLoopDuration(iCfg.maximumLoopDuration ?? 0);
      setNoBrowserFlag(iCfg?.noBrowser ?? false);
      setAutostartVideo(iCfg.autostartVideo ?? false);
      setAutostartVideoOnPlaySelected(
        iCfg.autostartVideoOnPlaySelected ?? true
      );
      setContinuePlaylistDefault(iCfg.continuePlaylistDefault ?? false);
      setShowStudioAsText(iCfg.showStudioAsText ?? false);
      setCSS(iCfg.css ?? "");
      setCSSEnabled(iCfg.cssEnabled ?? false);
      setLanguage(iCfg.language ?? "en-US");
      setSlideshowDelay(iCfg.slideshowDelay ?? 5000);
      setHandyKey(iCfg.handyKey ?? "");
      setFunscriptOffset(iCfg.funscriptOffset ?? 0);
      setDisableDropdownCreate({
        performer: iCfg.disabledDropdownCreate.performer,
        studio: iCfg.disabledDropdownCreate.studio,
        tag: iCfg.disabledDropdownCreate.tag,
      });

      setDeleteFileDefault(defaults.deleteFile ?? false);
      setDeleteGeneratedDefault(defaults.deleteGenerated ?? true);
    }
  }, [config]);

  async function onSave() {
    const prevCSS = config?.configuration.interface.css;
    const prevCSSenabled = config?.configuration.interface.cssEnabled;
    try {
      if (config?.configuration.defaults) {
        await updateDefaultsConfig({
          variables: {
            input: {
              ...withoutTypename(config?.configuration.defaults),
              deleteFile: deleteFileDefault,
              deleteGenerated: deleteGeneratedDefault,
            },
          },
        });
      }
      const result = await updateInterfaceConfig();

      // Force refetch of custom css if it was changed
      if (
        prevCSS !== result.data?.configureInterface.css ||
        prevCSSenabled !== result.data?.configureInterface.cssEnabled
      ) {
        await fetch("/css", { cache: "reload" });
        window.location.reload();
      }

      Toast.success({
        content: intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl
              .formatMessage({ id: "configuration" })
              .toLocaleLowerCase(),
          }
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  if (error) return <h1>{error.message}</h1>;
  if (loading) return <LoadingIndicator />;

  return (
    <>
      <h4>{intl.formatMessage({ id: "config.ui.title" })}</h4>
      <Form.Group controlId="language">
        <h5>{intl.formatMessage({ id: "config.ui.language.heading" })}</h5>
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
          <option value="es-ES">Spanish (Spain)</option>
          <option value="de-DE">German (Germany)</option>
          <option value="pt-BR">Portuguese (Brazil)</option>
          <option value="fr-FR">French (France)</option>
          <option value="it-IT">Italian (Italy)</option>
          <option value="sv-SE">Swedish (Sweden)</option>
          <option value="zh-TW">繁體中文 (台灣)</option>
          <option value="zh-CN">简体中文 (中国)</option>
        </Form.Control>
      </Form.Group>
      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.ui.menu_items.heading" })}</h5>
        <CheckboxGroup
          groupId="menu-items"
          items={allMenuItems}
          checkedIds={menuItemIds}
          onChange={setMenuItemIds}
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.ui.menu_items.description" })}
        </Form.Text>
      </Form.Group>

      <hr />

      <h4>
        {intl.formatMessage({
          id: "config.ui.desktop_integration.desktop_integration",
        })}
      </h4>
      <Form.Group>
        <Form.Check
          id="skip-browser"
          checked={noBrowser}
          label={intl.formatMessage({
            id: "config.ui.desktop_integration.skip_opening_browser",
          })}
          onChange={() => setNoBrowserFlag(!noBrowser)}
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "config.ui.desktop_integration.skip_opening_browser_on_startup",
          })}
        </Form.Text>
      </Form.Group>
      <hr />

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.ui.scene_wall.heading" })}</h5>
        <Form.Check
          id="wall-show-title"
          checked={wallShowTitle}
          label={intl.formatMessage({
            id: "config.ui.scene_wall.options.display_title",
          })}
          onChange={() => setWallShowTitle(!wallShowTitle)}
        />
        <Form.Check
          id="wall-sound-enabled"
          checked={soundOnPreview}
          label={intl.formatMessage({
            id: "config.ui.scene_wall.options.toggle_sound",
          })}
          onChange={() => setSoundOnPreview(!soundOnPreview)}
        />
        <Form.Label htmlFor="wall-preview">
          <h6>
            {intl.formatMessage({ id: "config.ui.preview_type.heading" })}
          </h6>
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
          <option value="video">
            {intl.formatMessage({ id: "config.ui.preview_type.options.video" })}
          </option>
          <option value="animation">
            {intl.formatMessage({
              id: "config.ui.preview_type.options.animated",
            })}
          </option>
          <option value="image">
            {intl.formatMessage({
              id: "config.ui.preview_type.options.static",
            })}
          </option>
        </Form.Control>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.ui.preview_type.description" })}
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.ui.scene_list.heading" })}</h5>
        <Form.Check
          id="show-text-studios"
          checked={showStudioAsText}
          label={intl.formatMessage({
            id: "config.ui.scene_list.options.show_studio_as_text",
          })}
          onChange={() => {
            setShowStudioAsText(!showStudioAsText);
          }}
        />
      </Form.Group>

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.ui.scene_player.heading" })}</h5>
        <Form.Group>
          <Form.Check
            id="auto-start-video"
            checked={autostartVideo}
            label={intl.formatMessage({
              id: "config.ui.scene_player.options.auto_start_video",
            })}
            onChange={() => {
              setAutostartVideo(!autostartVideo);
            }}
          />
        </Form.Group>
        <Form.Group id="auto-start-video-on-play-selected">
          <Form.Check
            checked={autostartVideoOnPlaySelected}
            label={intl.formatMessage({
              id:
                "config.ui.scene_player.options.auto_start_video_on_play_selected.heading",
            })}
            onChange={() => {
              setAutostartVideoOnPlaySelected(!autostartVideoOnPlaySelected);
            }}
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id:
                "config.ui.scene_player.options.auto_start_video_on_play_selected.description",
            })}
          </Form.Text>
        </Form.Group>

        <Form.Group id="continue-playlist-default">
          <Form.Check
            checked={continuePlaylistDefault}
            label={intl.formatMessage({
              id:
                "config.ui.scene_player.options.continue_playlist_default.heading",
            })}
            onChange={() => {
              setContinuePlaylistDefault(!continuePlaylistDefault);
            }}
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id:
                "config.ui.scene_player.options.continue_playlist_default.description",
            })}
          </Form.Text>
        </Form.Group>

        <Form.Group id="max-loop-duration">
          <h6>
            {intl.formatMessage({ id: "config.ui.max_loop_duration.heading" })}
          </h6>
          <DurationInput
            className="row col col-4"
            numericValue={maximumLoopDuration}
            onValueChange={(duration) => setMaximumLoopDuration(duration ?? 0)}
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id: "config.ui.max_loop_duration.description",
            })}
          </Form.Text>
        </Form.Group>
      </Form.Group>

      <Form.Group id="slideshow-delay">
        <h5>
          {intl.formatMessage({ id: "config.ui.slideshow_delay.heading" })}
        </h5>
        <Form.Control
          className="col col-sm-6 text-input"
          type="number"
          value={slideshowDelay / SECONDS_TO_MS}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setSlideshowDelay(
              Number.parseInt(e.currentTarget.value, 10) * SECONDS_TO_MS
            );
          }}
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.ui.slideshow_delay.description" })}
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.ui.editing.heading" })}</h5>

        <Form.Group>
          <h6>
            {intl.formatMessage({
              id: "config.ui.editing.disable_dropdown_create.heading",
            })}
          </h6>
          <Form.Check
            id="disableDropdownCreate_performer"
            checked={disableDropdownCreate.performer ?? false}
            label={intl.formatMessage({
              id: "performer",
            })}
            onChange={() => {
              setDisableDropdownCreate({
                ...disableDropdownCreate,
                performer: !disableDropdownCreate.performer ?? true,
              });
            }}
          />

          <Form.Check
            id="disableDropdownCreate_studio"
            checked={disableDropdownCreate.studio ?? false}
            label={intl.formatMessage({
              id: "studio",
            })}
            onChange={() => {
              setDisableDropdownCreate({
                ...disableDropdownCreate,
                studio: !disableDropdownCreate.studio ?? true,
              });
            }}
          />

          <Form.Check
            id="disableDropdownCreate_tag"
            checked={disableDropdownCreate.tag ?? false}
            label={intl.formatMessage({
              id: "tag",
            })}
            onChange={() => {
              setDisableDropdownCreate({
                ...disableDropdownCreate,
                tag: !disableDropdownCreate.tag ?? true,
              });
            }}
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id: "config.ui.editing.disable_dropdown_create.description",
            })}
          </Form.Text>
        </Form.Group>
      </Form.Group>

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.ui.custom_css.heading" })}</h5>
        <Form.Check
          id="custom-css"
          checked={cssEnabled}
          label={intl.formatMessage({
            id: "config.ui.custom_css.option_label",
          })}
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
          {intl.formatMessage({ id: "config.ui.custom_css.description" })}
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <h5>
          {intl.formatMessage({ id: "config.ui.handy_connection_key.heading" })}
        </h5>
        <Form.Control
          className="col col-sm-6 text-input"
          value={handyKey}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setHandyKey(e.currentTarget.value);
          }}
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "config.ui.handy_connection_key.description",
          })}
        </Form.Text>
      </Form.Group>
      <Form.Group>
        <h5>
          {intl.formatMessage({ id: "config.ui.funscript_offset.heading" })}
        </h5>
        <Form.Control
          className="col col-sm-6 text-input"
          type="number"
          value={funscriptOffset}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setFunscriptOffset(Number.parseInt(e.currentTarget.value, 10));
          }}
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.ui.funscript_offset.description" })}
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <h5>
          {intl.formatMessage({ id: "config.ui.delete_options.heading" })}
        </h5>
        <Form.Check
          id="delete-file-default"
          checked={deleteFileDefault}
          label={intl.formatMessage({
            id: "config.ui.delete_options.options.delete_file",
          })}
          onChange={() => {
            setDeleteFileDefault(!deleteFileDefault);
          }}
        />
        <Form.Check
          id="delete-generated-default"
          checked={deleteGeneratedDefault}
          label={intl.formatMessage({
            id:
              "config.ui.delete_options.options.delete_generated_supporting_files",
          })}
          onChange={() => {
            setDeleteGeneratedDefault(!deleteGeneratedDefault);
          }}
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "config.ui.delete_options.description",
          })}
        </Form.Text>
      </Form.Group>

      <hr />
      <Button variant="primary" onClick={() => onSave()}>
        {intl.formatMessage({ id: "actions.save" })}
      </Button>
    </>
  );
};
