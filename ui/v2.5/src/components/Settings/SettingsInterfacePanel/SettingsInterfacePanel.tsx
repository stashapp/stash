import React from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { DurationInput, LoadingIndicator } from "src/components/Shared";
import { CheckboxGroup } from "./CheckboxGroup";
import { SettingSection } from "../SettingSection";
import {
  BooleanSetting,
  ModalSetting,
  NumberSetting,
  SelectSetting,
  StringSetting,
} from "../Inputs";
import { SettingStateContext } from "../context";
import { DurationUtils } from "src/utils";
import * as GQL from "src/core/generated-graphql";
import {
  imageLightboxDisplayModeIntlMap,
  imageLightboxScrollModeIntlMap,
} from "src/core/enums";
import { useInterfaceLocalForage } from "src/hooks";
import {
  ConnectionState,
  connectionStateLabel,
  InteractiveContext,
} from "src/hooks/Interactive/context";

const allMenuItems = [
  { id: "scenes", headingID: "scenes" },
  { id: "images", headingID: "images" },
  { id: "movies", headingID: "movies" },
  { id: "markers", headingID: "markers" },
  { id: "galleries", headingID: "galleries" },
  { id: "performers", headingID: "performers" },
  { id: "studios", headingID: "studios" },
  { id: "tags", headingID: "tags" },
];

export const SettingsInterfacePanel: React.FC = () => {
  const intl = useIntl();

  const { interface: iface, saveInterface, loading, error } = React.useContext(
    SettingStateContext
  );

  const {
    interactive,
    state: interactiveState,
    error: interactiveError,
    serverOffset: interactiveServerOffset,
    initialised: interactiveInitialised,
    initialise: initialiseInteractive,
    sync: interactiveSync,
  } = React.useContext(InteractiveContext);

  const [, setInterfaceLocalForage] = useInterfaceLocalForage();

  function saveLightboxSettings(v: Partial<GQL.ConfigImageLightboxInput>) {
    // save in local forage as well for consistency
    setInterfaceLocalForage((prev) => {
      return {
        ...prev,
        imageLightbox: {
          ...prev.imageLightbox,
          ...v,
        },
      };
    });

    saveInterface({
      imageLightbox: {
        ...iface.imageLightbox,
        ...v,
      },
    });
  }

  if (error) return <h1>{error.message}</h1>;
  if (loading) return <LoadingIndicator />;

  // https://en.wikipedia.org/wiki/List_of_language_names
  return (
    <>
      <SettingSection headingID="config.ui.basic_settings">
        <SelectSetting
          id="language"
          headingID="config.ui.language.heading"
          value={iface.language ?? undefined}
          onChange={(v) => saveInterface({ language: v })}
        >
          <option value="da-DK">Dansk (Danmark)</option>
          <option value="de-DE">Deutsch (Deutschland)</option>
          <option value="en-GB">English (United Kingdom)</option>
          <option value="en-US">English (United States)</option>
          <option value="es-ES">Español (España)</option>
          <option value="fi-FI">Suomi</option>
          <option value="fr-FR">Français (France)</option>
          <option value="hr-HR">Hrvatski (Preview)</option>
          <option value="it-IT">Italiano</option>
          <option value="ja-JP">日本語 (日本)</option>
          <option value="ko-KR">한국어 (대한민국) (Preview)</option>
          <option value="nl-NL">Nederlands (Nederland)</option>
          <option value="pl-PL">Polski</option>
          <option value="pt-BR">Português (Brasil)</option>
          <option value="ru-RU">Русский (Россия) (Preview)</option>
          <option value="sv-SE">Svenska</option>
          <option value="tr-TR">Türkçe (Türkiye)</option>
          <option value="zh-TW">繁體中文 (台灣)</option>
          <option value="zh-CN">简体中文 (中国)</option>
        </SelectSetting>

        <div className="setting-group">
          <div className="setting">
            <div>
              <h3>
                {intl.formatMessage({
                  id: "config.ui.menu_items.heading",
                })}
              </h3>
              <div className="sub-heading">
                {intl.formatMessage({ id: "config.ui.menu_items.description" })}
              </div>
            </div>
            <div />
          </div>
          <CheckboxGroup
            groupId="menu-items"
            items={allMenuItems}
            checkedIds={iface.menuItems ?? undefined}
            onChange={(v) => saveInterface({ menuItems: v })}
          />
        </div>
      </SettingSection>

      <SettingSection headingID="config.ui.desktop_integration.desktop_integration">
        <BooleanSetting
          id="skip-browser"
          headingID="config.ui.desktop_integration.skip_opening_browser"
          subHeadingID="config.ui.desktop_integration.skip_opening_browser_on_startup"
          checked={iface.noBrowser ?? undefined}
          onChange={(v) => saveInterface({ noBrowser: v })}
        />
        <BooleanSetting
          id="notifications-enabled"
          headingID="config.ui.desktop_integration.notifications_enabled"
          subHeadingID="config.ui.desktop_integration.send_desktop_notifications_for_events"
          checked={iface.notificationsEnabled ?? undefined}
          onChange={(v) => saveInterface({ notificationsEnabled: v })}
        />
      </SettingSection>

      <SettingSection headingID="config.ui.scene_wall.heading">
        <BooleanSetting
          id="wall-show-title"
          headingID="config.ui.scene_wall.options.display_title"
          checked={iface.wallShowTitle ?? undefined}
          onChange={(v) => saveInterface({ wallShowTitle: v })}
        />
        <BooleanSetting
          id="wall-sound-enabled"
          headingID="config.ui.scene_wall.options.toggle_sound"
          checked={iface.soundOnPreview ?? undefined}
          onChange={(v) => saveInterface({ soundOnPreview: v })}
        />

        <SelectSetting
          id="wall-preview"
          headingID="config.ui.preview_type.heading"
          subHeadingID="config.ui.preview_type.description"
          value={iface.wallPlayback ?? undefined}
          onChange={(v) => saveInterface({ wallPlayback: v })}
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
        </SelectSetting>
      </SettingSection>

      <SettingSection headingID="config.ui.scene_list.heading">
        <BooleanSetting
          id="show-text-studios"
          headingID="config.ui.scene_list.options.show_studio_as_text"
          checked={iface.showStudioAsText ?? undefined}
          onChange={(v) => saveInterface({ showStudioAsText: v })}
        />
      </SettingSection>

      <SettingSection headingID="config.ui.scene_player.heading">
        <BooleanSetting
          id="show-scrubber"
          headingID="config.ui.scene_player.options.show_scrubber"
          checked={iface.showScrubber ?? undefined}
          onChange={(v) => saveInterface({ showScrubber: v })}
        />
        <BooleanSetting
          id="auto-start-video"
          headingID="config.ui.scene_player.options.auto_start_video"
          checked={iface.autostartVideo ?? undefined}
          onChange={(v) => saveInterface({ autostartVideo: v })}
        />
        <BooleanSetting
          id="auto-start-video-on-play-selected"
          headingID="config.ui.scene_player.options.auto_start_video_on_play_selected.heading"
          subHeadingID="config.ui.scene_player.options.auto_start_video_on_play_selected.description"
          checked={iface.autostartVideoOnPlaySelected ?? undefined}
          onChange={(v) => saveInterface({ autostartVideoOnPlaySelected: v })}
        />

        <BooleanSetting
          id="continue-playlist-default"
          headingID="config.ui.scene_player.options.continue_playlist_default.heading"
          subHeadingID="config.ui.scene_player.options.continue_playlist_default.description"
          checked={iface.continuePlaylistDefault ?? undefined}
          onChange={(v) => saveInterface({ continuePlaylistDefault: v })}
        />

        <ModalSetting<number>
          id="max-loop-duration"
          headingID="config.ui.max_loop_duration.heading"
          subHeadingID="config.ui.max_loop_duration.description"
          value={iface.maximumLoopDuration ?? undefined}
          onChange={(v) => saveInterface({ maximumLoopDuration: v })}
          renderField={(value, setValue) => (
            <DurationInput
              numericValue={value}
              onValueChange={(duration) => setValue(duration ?? 0)}
            />
          )}
          renderValue={(v) => {
            return <span>{DurationUtils.secondsToString(v ?? 0)}</span>;
          }}
        />
      </SettingSection>

      <SettingSection headingID="config.ui.image_lightbox.heading">
        <NumberSetting
          headingID="config.ui.slideshow_delay.heading"
          subHeadingID="config.ui.slideshow_delay.description"
          value={iface.imageLightbox?.slideshowDelay ?? undefined}
          onChange={(v) => saveLightboxSettings({ slideshowDelay: v })}
        />

        <SelectSetting
          id="lightbox_display_mode"
          headingID="dialogs.lightbox.display_mode.label"
          value={
            iface.imageLightbox?.displayMode ??
            GQL.ImageLightboxDisplayMode.FitXy
          }
          onChange={(v) =>
            saveLightboxSettings({
              displayMode: v as GQL.ImageLightboxDisplayMode,
            })
          }
        >
          {Array.from(imageLightboxDisplayModeIntlMap.entries()).map((v) => (
            <option key={v[0]} value={v[0]}>
              {intl.formatMessage({
                id: v[1],
              })}
            </option>
          ))}
        </SelectSetting>

        <BooleanSetting
          id="lightbox_scale_up"
          headingID="dialogs.lightbox.scale_up.label"
          subHeadingID="dialogs.lightbox.scale_up.description"
          checked={iface.imageLightbox?.scaleUp ?? false}
          onChange={(v) => saveLightboxSettings({ scaleUp: v })}
        />

        <BooleanSetting
          id="lightbox_reset_zoom_on_nav"
          headingID="dialogs.lightbox.reset_zoom_on_nav"
          checked={iface.imageLightbox?.resetZoomOnNav ?? false}
          onChange={(v) => saveLightboxSettings({ resetZoomOnNav: v })}
        />

        <SelectSetting
          id="lightbox_scroll_mode"
          headingID="dialogs.lightbox.scroll_mode.label"
          subHeadingID="dialogs.lightbox.scroll_mode.description"
          value={
            iface.imageLightbox?.scrollMode ?? GQL.ImageLightboxScrollMode.Zoom
          }
          onChange={(v) =>
            saveLightboxSettings({
              scrollMode: v as GQL.ImageLightboxScrollMode,
            })
          }
        >
          {Array.from(imageLightboxScrollModeIntlMap.entries()).map((v) => (
            <option key={v[0]} value={v[0]}>
              {intl.formatMessage({
                id: v[1],
              })}
            </option>
          ))}
        </SelectSetting>

        <NumberSetting
          headingID="config.ui.scroll_attempts_before_change.heading"
          subHeadingID="config.ui.scroll_attempts_before_change.description"
          value={iface.imageLightbox?.scrollAttemptsBeforeChange ?? 0}
          onChange={(v) =>
            saveLightboxSettings({ scrollAttemptsBeforeChange: v })
          }
        />
      </SettingSection>

      <SettingSection headingID="config.ui.editing.heading">
        <div className="setting-group">
          <div className="setting">
            <div>
              <h3>
                {intl.formatMessage({
                  id: "config.ui.editing.disable_dropdown_create.heading",
                })}
              </h3>
              <div className="sub-heading">
                {intl.formatMessage({
                  id: "config.ui.editing.disable_dropdown_create.description",
                })}
              </div>
            </div>
            <div />
          </div>
          <BooleanSetting
            id="disableDropdownCreate_performer"
            headingID="performer"
            checked={iface.disableDropdownCreate?.performer ?? undefined}
            onChange={(v) =>
              saveInterface({
                disableDropdownCreate: {
                  ...iface.disableDropdownCreate,
                  performer: v,
                },
              })
            }
          />
          <BooleanSetting
            id="disableDropdownCreate_studio"
            headingID="studio"
            checked={iface.disableDropdownCreate?.studio ?? undefined}
            onChange={(v) =>
              saveInterface({
                disableDropdownCreate: {
                  ...iface.disableDropdownCreate,
                  studio: v,
                },
              })
            }
          />
          <BooleanSetting
            id="disableDropdownCreate_tag"
            headingID="tag"
            checked={iface.disableDropdownCreate?.tag ?? undefined}
            onChange={(v) =>
              saveInterface({
                disableDropdownCreate: {
                  ...iface.disableDropdownCreate,
                  tag: v,
                },
              })
            }
          />
        </div>
      </SettingSection>

      <SettingSection headingID="config.ui.custom_css.heading">
        <BooleanSetting
          id="custom-css-enabled"
          headingID="config.ui.custom_css.option_label"
          checked={iface.cssEnabled ?? undefined}
          onChange={(v) => saveInterface({ cssEnabled: v })}
        />

        <ModalSetting<string>
          id="custom-css"
          headingID="config.ui.custom_css.heading"
          subHeadingID="config.ui.custom_css.description"
          value={iface.css ?? undefined}
          onChange={(v) => saveInterface({ css: v })}
          renderField={(value, setValue) => (
            <Form.Control
              as="textarea"
              value={value}
              onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                setValue(e.currentTarget.value)
              }
              rows={16}
              className="text-input code"
            />
          )}
          renderValue={() => {
            return <></>;
          }}
        />
      </SettingSection>

      <SettingSection headingID="config.ui.interactive_options">
        <StringSetting
          headingID="config.ui.handy_connection_key.heading"
          subHeadingID="config.ui.handy_connection_key.description"
          value={iface.handyKey ?? undefined}
          onChange={(v) => saveInterface({ handyKey: v })}
        />
        {interactive.handyKey && (
          <>
            <div className="setting" id="handy-status">
              <div>
                <h3>
                  {intl.formatMessage({
                    id: "config.ui.handy_connection.status.heading",
                  })}
                </h3>

                <div className="value">
                  <FormattedMessage
                    id={connectionStateLabel(interactiveState)}
                  />
                  {interactiveError && <span>: {interactiveError}</span>}
                </div>
              </div>
              <div>
                {!interactiveInitialised && (
                  <Button
                    disabled={
                      interactiveState === ConnectionState.Connecting ||
                      interactiveState === ConnectionState.Syncing
                    }
                    onClick={() => initialiseInteractive()}
                  >
                    {intl.formatMessage({
                      id: "config.ui.handy_connection.connect",
                    })}
                  </Button>
                )}
              </div>
            </div>
            <div className="setting" id="handy-server-offset">
              <div>
                <h3>
                  {intl.formatMessage({
                    id: "config.ui.handy_connection.server_offset.heading",
                  })}
                </h3>

                <div className="value">
                  {interactiveServerOffset.toFixed()}ms
                </div>
              </div>
              <div>
                {interactiveInitialised && (
                  <Button
                    disabled={
                      !interactiveInitialised ||
                      interactiveState === ConnectionState.Syncing
                    }
                    onClick={() => interactiveSync()}
                  >
                    {intl.formatMessage({
                      id: "config.ui.handy_connection.sync",
                    })}
                  </Button>
                )}
              </div>
            </div>
          </>
        )}

        <NumberSetting
          headingID="config.ui.funscript_offset.heading"
          subHeadingID="config.ui.funscript_offset.description"
          value={iface.funscriptOffset ?? undefined}
          onChange={(v) => saveInterface({ funscriptOffset: v })}
        />
      </SettingSection>
    </>
  );
};
