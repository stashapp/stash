import React from "react";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { StashSetting } from "./StashConfiguration";
import { SettingSection } from "./SettingSection";
import { BooleanSetting, StringListSetting, StringSetting } from "./Inputs";
import { useSettings } from "./context";
import { useIntl } from "react-intl";
import { faQuestionCircle } from "@fortawesome/free-solid-svg-icons";
import { ExternalLink } from "../Shared/ExternalLink";

export const SettingsLibraryPanel: React.FC = () => {
  const intl = useIntl();
  const { general, loading, error, saveGeneral, defaults, saveDefaults } =
    useSettings();

  function commaDelimitedToList(value: string | undefined) {
    if (value) {
      return value.split(",").map((s) => s.trim());
    }
  }

  function listToCommaDelimited(value: string[] | undefined) {
    if (value) {
      return value.join(", ");
    }
  }

  if (error) return <h1>{error.message}</h1>;
  if (loading) return <LoadingIndicator />;

  return (
    <>
      <StashSetting
        value={general.stashes ?? []}
        onChange={(v) => saveGeneral({ stashes: v })}
      />

      <SettingSection headingID="config.library.media_content_extensions">
        <StringSetting
          id="video-extensions"
          headingID="config.general.video_ext_head"
          subHeadingID="config.general.video_ext_desc"
          value={listToCommaDelimited(general.videoExtensions ?? undefined)}
          onChange={(v) =>
            saveGeneral({ videoExtensions: commaDelimitedToList(v) })
          }
        />

        <StringSetting
          id="image-extensions"
          headingID="config.general.image_ext_head"
          subHeadingID="config.general.image_ext_desc"
          value={listToCommaDelimited(general.imageExtensions ?? undefined)}
          onChange={(v) =>
            saveGeneral({ imageExtensions: commaDelimitedToList(v) })
          }
        />

        <StringSetting
          id="gallery-extensions"
          headingID="config.general.gallery_ext_head"
          subHeadingID="config.general.gallery_ext_desc"
          value={listToCommaDelimited(general.galleryExtensions ?? undefined)}
          onChange={(v) =>
            saveGeneral({ galleryExtensions: commaDelimitedToList(v) })
          }
        />
      </SettingSection>

      <SettingSection headingID="config.library.exclusions">
        <StringListSetting
          id="excluded-video-patterns"
          headingID="config.general.excluded_video_patterns_head"
          subHeading={
            <span>
              {intl.formatMessage({
                id: "config.general.excluded_video_patterns_desc",
              })}
              <ExternalLink href="https://docs.stashapp.cc/beginner-guides/exclude-file-configuration">
                <Icon icon={faQuestionCircle} />
              </ExternalLink>
            </span>
          }
          value={general.excludes ?? undefined}
          onChange={(v) => saveGeneral({ excludes: v })}
          defaultNewValue="sample\.mp4$"
        />

        <StringListSetting
          id="excluded-image-gallery-patterns"
          headingID="config.general.excluded_image_gallery_patterns_head"
          subHeading={
            <span>
              {intl.formatMessage({
                id: "config.general.excluded_image_gallery_patterns_desc",
              })}
              <ExternalLink href="https://docs.stashapp.cc/beginner-guides/exclude-file-configuration">
                <Icon icon={faQuestionCircle} />
              </ExternalLink>
            </span>
          }
          value={general.imageExcludes ?? undefined}
          onChange={(v) => saveGeneral({ imageExcludes: v })}
          defaultNewValue="sample\.jpg$"
        />
      </SettingSection>

      <SettingSection headingID="config.library.gallery_and_image_options">
        <BooleanSetting
          id="create-galleries-from-folders"
          headingID="config.general.create_galleries_from_folders_label"
          subHeadingID="config.general.create_galleries_from_folders_desc"
          checked={general.createGalleriesFromFolders ?? false}
          onChange={(v) => saveGeneral({ createGalleriesFromFolders: v })}
        />

        <BooleanSetting
          id="write-image-thumbnails"
          headingID="config.ui.images.options.write_image_thumbnails.heading"
          subHeadingID="config.ui.images.options.write_image_thumbnails.description"
          checked={general.writeImageThumbnails ?? false}
          onChange={(v) => saveGeneral({ writeImageThumbnails: v })}
        />

        <BooleanSetting
          id="create-image-clips-from-videos"
          headingID="config.ui.images.options.create_image_clips_from_videos.heading"
          subHeadingID="config.ui.images.options.create_image_clips_from_videos.description"
          checked={general.createImageClipsFromVideos ?? false}
          onChange={(v) => saveGeneral({ createImageClipsFromVideos: v })}
        />

        <StringSetting
          id="gallery-cover-regex"
          headingID="config.general.gallery_cover_regex_label"
          subHeadingID="config.general.gallery_cover_regex_desc"
          value={general.galleryCoverRegex ?? ""}
          onChange={(v) => saveGeneral({ galleryCoverRegex: v })}
        />
      </SettingSection>

      <SettingSection headingID="config.ui.delete_options.heading">
        <BooleanSetting
          id="delete-file-default"
          headingID="config.ui.delete_options.options.delete_file"
          checked={defaults.deleteFile ?? undefined}
          onChange={(v) => {
            saveDefaults({ deleteFile: v });
          }}
        />
        <BooleanSetting
          id="delete-generated-default"
          headingID="config.ui.delete_options.options.delete_generated_supporting_files"
          subHeadingID="config.ui.delete_options.description"
          checked={defaults.deleteGenerated ?? undefined}
          onChange={(v) => {
            saveDefaults({ deleteGenerated: v });
          }}
        />
      </SettingSection>
    </>
  );
};
