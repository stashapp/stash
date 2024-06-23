import React from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { DurationInput } from "src/components/Shared/DurationInput";
import { PercentInput } from "src/components/Shared/PercentInput";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { CheckboxGroup } from "./CheckboxGroup";
import { SettingSection } from "../SettingSection";
import {
  BooleanSetting,
  ModalSetting,
  NumberSetting,
  SelectSetting,
  StringSetting,
} from "../Inputs";
import { useSettings } from "../context";
import TextUtils from "src/utils/text";
import * as GQL from "src/core/generated-graphql";
import {
  imageLightboxDisplayModeIntlMap,
  imageLightboxScrollModeIntlMap,
} from "src/core/enums";
import { useInterfaceLocalForage } from "src/hooks/LocalForage";
import {
  ConnectionState,
  connectionStateLabel,
  InteractiveContext,
} from "src/hooks/Interactive/context";
import {
  defaultRatingStarPrecision,
  defaultRatingSystemOptions,
  defaultRatingSystemType,
  RatingStarPrecision,
  ratingStarPrecisionIntlMap,
  ratingSystemIntlMap,
  RatingSystemType,
} from "src/utils/rating";
import {
  imageWallDirectionIntlMap,
  ImageWallDirection,
  defaultImageWallOptions,
  defaultImageWallDirection,
  defaultImageWallMargin,
} from "src/utils/imageWall";
import { defaultMaxOptionsShown } from "src/core/config";

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

  const {
    interface: iface,
    saveInterface,
    ui,
    saveUI,
    loading,
    error,
  } = useSettings();

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

  function saveImageWallMargin(m: number) {
    saveUI({
      imageWallOptions: {
        ...(ui.imageWallOptions ?? defaultImageWallOptions),
        margin: m,
      },
    });
  }

  function saveImageWallDirection(d: ImageWallDirection) {
    saveUI({
      imageWallOptions: {
        ...(ui.imageWallOptions ?? defaultImageWallOptions),
        direction: d,
      },
    });
  }

  function saveRatingSystemType(t: RatingSystemType) {
    saveUI({
      ratingSystemOptions: {
        ...ui.ratingSystemOptions,
        type: t,
      },
    });
  }

  function saveRatingSystemStarPrecision(p: RatingStarPrecision) {
    saveUI({
      ratingSystemOptions: {
        ...(ui.ratingSystemOptions ?? defaultRatingSystemOptions),
        starPrecision: p,
      },
    });
  }

  function validateLocaleString(v: string) {
    if (!v) return;
    try {
      JSON.parse(v);
    } catch (e) {
      throw new Error(
        intl.formatMessage(
          { id: "errors.invalid_json_string" },
          {
            error: (e as SyntaxError).message,
          }
        )
      );
    }
  }

  function validateJavascriptString(v: string) {
    if (!v) return;
    try {
      // creates a function from the string to validate it but does not execute it
      // eslint-disable-next-line @typescript-eslint/no-implied-eval
      new Function(v);
    } catch (e) {
      throw new Error(
        intl.formatMessage(
          { id: "errors.invalid_javascript_string" },
          {
            error: (e as SyntaxError).message,
          }
        )
      );
    }
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
          <option value="af-ZA">Afrikaans (Preview)</option>
          <option value="bn-BD">বাংলা (বাংলাদেশ) (Preview)</option>
          <option value="ca-ES">Catalan (Preview)</option>
          <option value="cs-CZ">Čeština (Česko)</option>
          <option value="da-DK">Dansk (Danmark)</option>
          <option value="de-DE">Deutsch (Deutschland)</option>
          <option value="en-GB">English (United Kingdom)</option>
          <option value="en-US">English (United States)</option>
          <option value="et-EE">Eesti</option>
          <option value="fa-IR">فارسی (ایران) (Preview)</option>
          <option value="fi-FI">Suomi</option>
          <option value="fr-FR">Français (France)</option>
          <option value="hr-HR">Hrvatski (Preview)</option>
          <option value="id-ID">Indonesian (Preview)</option>
          <option value="hu-HU">Magyar (Preview)</option>
          <option value="it-IT">Italiano</option>
          <option value="ja-JP">日本語 (日本)</option>
          <option value="ko-KR">한국어 (대한민국)</option>
          <option value="nl-NL">Nederlands (Nederland)</option>
          <option value="pl-PL">Polski</option>
          <option value="pt-BR">Português (Brasil)</option>
          <option value="ro-RO">Română (Preview)</option>
          <option value="ru-RU">Русский (Россия)</option>
          <option value="es-ES">Español (España)</option>
          <option value="sv-SE">Svenska</option>
          <option value="tr-TR">Türkçe (Türkiye)</option>
          <option value="th-TH">ภาษาไทย (Preview)</option>
          <option value="uk-UA">Ukrainian (Preview)</option>
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

        <BooleanSetting
          id="abbreviate-counters"
          headingID="config.ui.abbreviate_counters.heading"
          subHeadingID="config.ui.abbreviate_counters.description"
          checked={ui.abbreviateCounters ?? undefined}
          onChange={(v) => saveUI({ abbreviateCounters: v })}
        />
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
          advanced
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
          id="enable-chromecast"
          headingID="config.ui.scene_player.options.enable_chromecast"
          checked={ui.enableChromecast ?? undefined}
          onChange={(v) => saveUI({ enableChromecast: v })}
        />
        <BooleanSetting
          id="disable-mobile-media-auto-rotate"
          headingID="config.ui.scene_player.options.disable_mobile_media_auto_rotate"
          checked={ui.disableMobileMediaAutoRotateEnabled ?? undefined}
          onChange={(v) => saveUI({ disableMobileMediaAutoRotateEnabled: v })}
        />
        <BooleanSetting
          id="show-scrubber"
          headingID="config.ui.scene_player.options.show_scrubber"
          checked={iface.showScrubber ?? undefined}
          onChange={(v) => saveInterface({ showScrubber: v })}
        />
        <BooleanSetting
          id="always-start-from-beginning"
          headingID="config.ui.scene_player.options.always_start_from_beginning"
          checked={ui.alwaysStartFromBeginning ?? undefined}
          onChange={(v) => saveUI({ alwaysStartFromBeginning: v })}
        />
        <BooleanSetting
          id="track-activity"
          headingID="config.ui.scene_player.options.track_activity"
          checked={ui.trackActivity ?? true}
          onChange={(v) => saveUI({ trackActivity: v })}
        />
        <StringSetting
          id="vr-tag"
          headingID="config.ui.scene_player.options.vr_tag.heading"
          subHeadingID="config.ui.scene_player.options.vr_tag.description"
          value={ui.vrTag ?? undefined}
          onChange={(v) => saveUI({ vrTag: v })}
        />
        <ModalSetting<number>
          id="ignore-interval"
          headingID="config.ui.minimum_play_percent.heading"
          subHeadingID="config.ui.minimum_play_percent.description"
          value={ui.minimumPlayPercent ?? 0}
          onChange={(v) => saveUI({ minimumPlayPercent: v })}
          disabled={!ui.trackActivity}
          renderField={(value, setValue) => (
            <PercentInput
              numericValue={value}
              onValueChange={(interval) => setValue(interval ?? 0)}
            />
          )}
          renderValue={(v) => {
            return <span>{v}%</span>;
          }}
        />
        <NumberSetting
          headingID="config.ui.slideshow_delay.heading"
          subHeadingID="config.ui.slideshow_delay.description"
          value={iface.imageLightbox?.slideshowDelay ?? undefined}
          onChange={(v) => saveLightboxSettings({ slideshowDelay: v })}
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
              value={value}
              setValue={(duration) => setValue(duration ?? 0)}
            />
          )}
          renderValue={(v) => {
            return <span>{TextUtils.secondsToTimestamp(v ?? 0)}</span>;
          }}
        />

        <BooleanSetting
          id="show-ab-loop"
          headingID="config.ui.scene_player.options.show_ab_loop_controls"
          checked={ui.showAbLoopControls ?? undefined}
          onChange={(v) => saveUI({ showAbLoopControls: v })}
        />
      </SettingSection>
      <SettingSection headingID="config.ui.tag_panel.heading">
        <BooleanSetting
          id="show-tag-card-on-hover"
          headingID="config.ui.show_tag_card_on_hover.heading"
          subHeadingID="config.ui.show_tag_card_on_hover.description"
          checked={ui.showTagCardOnHover ?? true}
          onChange={(v) => saveUI({ showTagCardOnHover: v })}
        />
        <BooleanSetting
          id="show-child-tagged-content"
          headingID="config.ui.tag_panel.options.show_child_tagged_content.heading"
          subHeadingID="config.ui.tag_panel.options.show_child_tagged_content.description"
          checked={ui.showChildTagContent ?? undefined}
          onChange={(v) => saveUI({ showChildTagContent: v })}
        />
      </SettingSection>
      <SettingSection headingID="config.ui.studio_panel.heading">
        <BooleanSetting
          id="show-child-studio-content"
          headingID="config.ui.studio_panel.options.show_child_studio_content.heading"
          subHeadingID="config.ui.studio_panel.options.show_child_studio_content.description"
          checked={ui.showChildStudioContent ?? undefined}
          onChange={(v) => saveUI({ showChildStudioContent: v })}
        />
      </SettingSection>

      <SettingSection headingID="config.ui.image_wall.heading">
        <NumberSetting
          headingID="config.ui.image_wall.margin"
          subHeadingID="dialogs.imagewall.margin_desc"
          value={ui.imageWallOptions?.margin ?? defaultImageWallMargin}
          onChange={(v) => saveImageWallMargin(v)}
        />

        <SelectSetting
          id="image_wall_direction"
          headingID="config.ui.image_wall.direction"
          subHeadingID="dialogs.imagewall.direction.description"
          value={ui.imageWallOptions?.direction ?? defaultImageWallDirection}
          onChange={(v) => saveImageWallDirection(v as ImageWallDirection)}
        >
          {Array.from(imageWallDirectionIntlMap.entries()).map((v) => (
            <option key={v[0]} value={v[0]}>
              {intl.formatMessage({
                id: v[1],
              })}
            </option>
          ))}
        </SelectSetting>
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

      <SettingSection headingID="config.ui.detail.heading">
        <div className="setting-group">
          <div className="setting">
            <div>
              <h3>
                {intl.formatMessage({
                  id: "config.ui.detail.enable_background_image.heading",
                })}
              </h3>
              <div className="sub-heading">
                {intl.formatMessage({
                  id: "config.ui.detail.enable_background_image.description",
                })}
              </div>
            </div>
            <div />
          </div>
          <BooleanSetting
            id="enableMovieBackgroundImage"
            headingID="movie"
            checked={ui.enableMovieBackgroundImage ?? undefined}
            onChange={(v) => saveUI({ enableMovieBackgroundImage: v })}
          />
          <BooleanSetting
            id="enablePerformerBackgroundImage"
            headingID="performer"
            checked={ui.enablePerformerBackgroundImage ?? undefined}
            onChange={(v) => saveUI({ enablePerformerBackgroundImage: v })}
          />
          <BooleanSetting
            id="enableStudioBackgroundImage"
            headingID="studio"
            checked={ui.enableStudioBackgroundImage ?? undefined}
            onChange={(v) => saveUI({ enableStudioBackgroundImage: v })}
          />
          <BooleanSetting
            id="enableTagBackgroundImage"
            headingID="tag"
            checked={ui.enableTagBackgroundImage ?? undefined}
            onChange={(v) => saveUI({ enableTagBackgroundImage: v })}
          />
        </div>
        <BooleanSetting
          id="show_all_details"
          headingID="config.ui.detail.show_all_details.heading"
          subHeadingID="config.ui.detail.show_all_details.description"
          checked={ui.showAllDetails ?? true}
          onChange={(v) => saveUI({ showAllDetails: v })}
        />
        <BooleanSetting
          id="compact_expanded_details"
          headingID="config.ui.detail.compact_expanded_details.heading"
          subHeadingID="config.ui.detail.compact_expanded_details.description"
          checked={ui.compactExpandedDetails ?? undefined}
          onChange={(v) => saveUI({ compactExpandedDetails: v })}
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
          <BooleanSetting
            id="disableDropdownCreate_movie"
            headingID="movie"
            checked={iface.disableDropdownCreate?.movie ?? undefined}
            onChange={(v) =>
              saveInterface({
                disableDropdownCreate: {
                  ...iface.disableDropdownCreate,
                  movie: v,
                },
              })
            }
          />
        </div>
        <NumberSetting
          id="max_options_shown"
          headingID="config.ui.editing.max_options_shown.label"
          value={ui.maxOptionsShown ?? defaultMaxOptionsShown}
          onChange={(v) => saveUI({ maxOptionsShown: v })}
        />
        <SelectSetting
          id="rating_system"
          headingID="config.ui.editing.rating_system.type.label"
          value={ui.ratingSystemOptions?.type ?? defaultRatingSystemType}
          onChange={(v) => saveRatingSystemType(v as RatingSystemType)}
        >
          {Array.from(ratingSystemIntlMap.entries()).map((v) => (
            <option key={v[0]} value={v[0]}>
              {intl.formatMessage({
                id: v[1],
              })}
            </option>
          ))}
        </SelectSetting>
        {(ui.ratingSystemOptions?.type ?? defaultRatingSystemType) ===
          RatingSystemType.Stars && (
          <SelectSetting
            id="rating_system_star_precision"
            headingID="config.ui.editing.rating_system.star_precision.label"
            value={
              ui.ratingSystemOptions?.starPrecision ??
              defaultRatingStarPrecision
            }
            onChange={(v) =>
              saveRatingSystemStarPrecision(v as RatingStarPrecision)
            }
          >
            {Array.from(ratingStarPrecisionIntlMap.entries()).map((v) => (
              <option key={v[0]} value={v[0]}>
                {intl.formatMessage({
                  id: v[1],
                })}
              </option>
            ))}
          </SelectSetting>
        )}
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
      <SettingSection headingID="config.ui.custom_javascript.heading">
        <BooleanSetting
          id="custom-javascript-enabled"
          headingID="config.ui.custom_javascript.option_label"
          checked={iface.javascriptEnabled ?? undefined}
          onChange={(v) => saveInterface({ javascriptEnabled: v })}
        />

        <ModalSetting<string>
          id="custom-javascript"
          headingID="config.ui.custom_javascript.heading"
          subHeadingID="config.ui.custom_javascript.description"
          value={iface.javascript ?? undefined}
          onChange={(v) => saveInterface({ javascript: v })}
          validateChange={validateJavascriptString}
          renderField={(value, setValue, err) => (
            <>
              <Form.Control
                as="textarea"
                value={value}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                  setValue(e.currentTarget.value)
                }
                rows={16}
                className="text-input code"
                isInvalid={!!err}
              />
              <Form.Control.Feedback type="invalid">
                {err}
              </Form.Control.Feedback>
            </>
          )}
          renderValue={() => {
            return <></>;
          }}
        />
      </SettingSection>
      <SettingSection headingID="config.ui.custom_locales.heading">
        <BooleanSetting
          id="custom-locales-enabled"
          headingID="config.ui.custom_locales.option_label"
          checked={iface.customLocalesEnabled ?? undefined}
          onChange={(v) => saveInterface({ customLocalesEnabled: v })}
        />

        <ModalSetting<string>
          id="custom-locales"
          headingID="config.ui.custom_locales.heading"
          subHeadingID="config.ui.custom_locales.description"
          value={iface.customLocales ?? undefined}
          onChange={(v) => saveInterface({ customLocales: v })}
          validateChange={validateLocaleString}
          renderField={(value, setValue, err) => (
            <>
              <Form.Control
                as="textarea"
                value={value}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) =>
                  setValue(e.currentTarget.value)
                }
                rows={16}
                className="text-input code"
                isInvalid={!!err}
              />
              <Form.Control.Feedback type="invalid">
                {err}
              </Form.Control.Feedback>
            </>
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

        <BooleanSetting
          id="use-stash-hosted-funscript"
          headingID="config.ui.use_stash_hosted_funscript.heading"
          subHeadingID="config.ui.use_stash_hosted_funscript.description"
          checked={iface.useStashHostedFunscript ?? false}
          onChange={(v) => saveInterface({ useStashHostedFunscript: v })}
        />
      </SettingSection>
    </>
  );
};
