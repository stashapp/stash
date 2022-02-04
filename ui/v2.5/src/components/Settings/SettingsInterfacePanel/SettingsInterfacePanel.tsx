import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
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

  if (error) return <h1>{error.message}</h1>;
  if (loading) return <LoadingIndicator />;

  return (
    <>
      <SettingSection headingID="config.ui.basic_settings">
        <SelectSetting
          id="language"
          headingID="config.ui.language.heading"
          value={iface.language ?? undefined}
          onChange={(v) => saveInterface({ language: v })}
        >
          <option value="de-DE">Deutsch (Deutschland)</option>
          <option value="en-GB">English (United Kingdom)</option>
          <option value="en-US">English (United States)</option>
          <option value="es-ES">Español (España)</option>
          <option value="fi-FI">Suomi</option>
          <option value="fr-FR">Français (France)</option>
          <option value="hr-HR">Hrvatski (Preview)</option>
          <option value="it-IT">Italiano</option>
          <option value="ja-JP">日本語 (日本)</option>
          <option value="nl-NL">Nederlands (Nederland)</option>
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

      <SettingSection headingID="config.ui.images.heading">
        <NumberSetting
          headingID="config.ui.slideshow_delay.heading"
          subHeadingID="config.ui.slideshow_delay.description"
          value={iface.slideshowDelay ?? undefined}
          onChange={(v) => saveInterface({ slideshowDelay: v })}
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
