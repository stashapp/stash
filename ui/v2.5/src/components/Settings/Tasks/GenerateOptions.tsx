import React from "react";
import * as GQL from "src/core/generated-graphql";
import { BooleanSetting, ModalSetting } from "../Inputs";
import {
  VideoPreviewInput,
  VideoPreviewSettingsInput,
} from "../GeneratePreviewOptions";

interface IGenerateOptions {
  type?: "scene" | "image";
  selection?: boolean;
  options: GQL.GenerateMetadataInput;
  setOptions: (s: GQL.GenerateMetadataInput) => void;
}

export const GenerateOptions: React.FC<IGenerateOptions> = ({
  type,
  selection,
  options,
  setOptions: setOptionsState,
}) => {
  const previewOptions: GQL.GeneratePreviewOptionsInput =
    options.previewOptions ?? {};

  function setOptions(input: Partial<GQL.GenerateMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  const showSceneOptions = !type || type === "scene";
  const showImageOptions = !type || type === "image";

  return (
    <>
      {showSceneOptions && (
        <>
          <BooleanSetting
            id="covers-task"
            headingID="dialogs.scene_gen.covers"
            checked={options.covers ?? false}
            onChange={(v) => setOptions({ covers: v })}
          />
          <BooleanSetting
            id="preview-task"
            checked={options.previews ?? false}
            headingID="dialogs.scene_gen.video_previews"
            tooltipID="dialogs.scene_gen.video_previews_tooltip"
            onChange={(v) => setOptions({ previews: v })}
          />
          <BooleanSetting
            advanced
            className="sub-setting"
            id="image-preview-task"
            checked={options.imagePreviews ?? false}
            disabled={!options.previews}
            headingID="dialogs.scene_gen.image_previews"
            tooltipID="dialogs.scene_gen.image_previews_tooltip"
            onChange={(v) => setOptions({ imagePreviews: v })}
          />

          {/* #2251 - only allow preview generation options to be overridden when generating from a selection */}
          {selection ? (
            <ModalSetting<VideoPreviewSettingsInput>
              id="video-preview-settings"
              className="sub-setting"
              disabled={!options.previews}
              headingID="dialogs.scene_gen.override_preview_generation_options"
              tooltipID="dialogs.scene_gen.override_preview_generation_options_desc"
              value={{
                previewExcludeEnd: previewOptions.previewExcludeEnd,
                previewExcludeStart: previewOptions.previewExcludeStart,
                previewSegmentDuration: previewOptions.previewSegmentDuration,
                previewSegments: previewOptions.previewSegments,
              }}
              onChange={(v) => setOptions({ previewOptions: v })}
              renderField={(value, setValue) => (
                <VideoPreviewInput value={value ?? {}} setValue={setValue} />
              )}
              renderValue={() => {
                return <></>;
              }}
            />
          ) : undefined}

          <BooleanSetting
            id="sprite-task"
            checked={options.sprites ?? false}
            headingID="dialogs.scene_gen.sprites"
            tooltipID="dialogs.scene_gen.sprites_tooltip"
            onChange={(v) => setOptions({ sprites: v })}
          />
          <BooleanSetting
            id="marker-task"
            checked={options.markers ?? false}
            headingID="dialogs.scene_gen.markers"
            tooltipID="dialogs.scene_gen.markers_tooltip"
            onChange={(v) => setOptions({ markers: v })}
          />
          <BooleanSetting
            advanced
            id="marker-image-preview-task"
            className="sub-setting"
            checked={options.markerImagePreviews ?? false}
            headingID="dialogs.scene_gen.marker_image_previews"
            tooltipID="dialogs.scene_gen.marker_image_previews_tooltip"
            onChange={(v) =>
              setOptions({
                markerImagePreviews: v,
              })
            }
          />
          <BooleanSetting
            id="marker-screenshot-task"
            checked={options.markerScreenshots ?? false}
            headingID="dialogs.scene_gen.marker_screenshots"
            tooltipID="dialogs.scene_gen.marker_screenshots_tooltip"
            onChange={(v) => setOptions({ markerScreenshots: v })}
          />

          <BooleanSetting
            advanced
            id="transcode-task"
            checked={options.transcodes ?? false}
            headingID="dialogs.scene_gen.transcodes"
            tooltipID="dialogs.scene_gen.transcodes_tooltip"
            onChange={(v) => setOptions({ transcodes: v })}
          />
          {selection ? (
            <BooleanSetting
              advanced
              id="force-transcode"
              className="sub-setting"
              checked={options.forceTranscodes ?? false}
              disabled={!options.transcodes}
              headingID="dialogs.scene_gen.force_transcodes"
              tooltipID="dialogs.scene_gen.force_transcodes_tooltip"
              onChange={(v) => setOptions({ forceTranscodes: v })}
            />
          ) : undefined}

          <BooleanSetting
            id="phash-task"
            checked={options.phashes ?? false}
            headingID="dialogs.scene_gen.phash"
            tooltipID="dialogs.scene_gen.phash_tooltip"
            onChange={(v) => setOptions({ phashes: v })}
          />

          <BooleanSetting
            id="interactive-heatmap-speed-task"
            checked={options.interactiveHeatmapsSpeeds ?? false}
            headingID="dialogs.scene_gen.interactive_heatmap_speed"
            onChange={(v) => setOptions({ interactiveHeatmapsSpeeds: v })}
          />
        </>
      )}
      {showImageOptions && (
        <>
          <BooleanSetting
            id="clip-previews"
            checked={options.clipPreviews ?? false}
            headingID="dialogs.scene_gen.clip_previews"
            onChange={(v) => setOptions({ clipPreviews: v })}
          />
          <BooleanSetting
            id="image-thumbnails"
            checked={options.imageThumbnails ?? false}
            headingID="dialogs.scene_gen.image_thumbnails"
            onChange={(v) => setOptions({ imageThumbnails: v })}
          />
        </>
      )}
      <BooleanSetting
        id="overwrite"
        checked={options.overwrite ?? false}
        headingID="dialogs.scene_gen.overwrite"
        onChange={(v) => setOptions({ overwrite: v })}
      />
    </>
  );
};
