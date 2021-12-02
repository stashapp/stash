import React from "react";
import * as GQL from "src/core/generated-graphql";
import { BooleanSetting, ModalSetting } from "../Inputs";
import {
  VideoPreviewInput,
  VideoPreviewSettingsInput,
} from "../GeneratePreviewOptions";

interface IGenerateOptions {
  options: GQL.GenerateMetadataInput;
  setOptions: (s: GQL.GenerateMetadataInput) => void;
}

export const GenerateOptions: React.FC<IGenerateOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  const previewOptions: GQL.GeneratePreviewOptionsInput =
    options.previewOptions ?? {};

  function setOptions(input: Partial<GQL.GenerateMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <>
      <BooleanSetting
        id="preview-task"
        checked={options.previews ?? false}
        headingID="dialogs.scene_gen.video_previews"
        onChange={(v) => setOptions({ previews: v })}
      />
      <BooleanSetting
        id="image-preview-task"
        checked={options.imagePreviews ?? false}
        disabled={!options.previews}
        headingID="dialogs.scene_gen.image_previews"
        onChange={(v) => setOptions({ imagePreviews: v })}
      />

      <ModalSetting<VideoPreviewSettingsInput>
        id="video-preview-settings"
        buttonTextID="dialogs.scene_gen.preview_generation_options"
        headingID="dialogs.scene_gen.preview_generation_options"
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

      <BooleanSetting
        id="sprite-task"
        checked={options.sprites ?? false}
        headingID="dialogs.scene_gen.sprites"
        onChange={(v) => setOptions({ sprites: v })}
      />
      <BooleanSetting
        id="marker-task"
        checked={options.markers ?? false}
        headingID="dialogs.scene_gen.markers"
        onChange={(v) => setOptions({ markers: v })}
      />
      <BooleanSetting
        id="marker-image-preview-task"
        checked={options.markerImagePreviews ?? false}
        disabled={!options.markers}
        headingID="dialogs.scene_gen.marker_image_previews"
        onChange={(v) =>
          setOptions({
            markerImagePreviews: v,
          })
        }
      />
      <BooleanSetting
        id="marker-screenshot-task"
        checked={options.markerScreenshots ?? false}
        disabled={!options.markers}
        headingID="dialogs.scene_gen.marker_screenshots"
        onChange={(v) => setOptions({ markerScreenshots: v })}
      />

      <BooleanSetting
        id="transcode-task"
        checked={options.transcodes ?? false}
        headingID="dialogs.scene_gen.transcodes"
        onChange={(v) => setOptions({ transcodes: v })}
      />
      <BooleanSetting
        id="phash-task"
        checked={options.phashes ?? false}
        headingID="dialogs.scene_gen.phash"
        onChange={(v) => setOptions({ phashes: v })}
      />

      <BooleanSetting
        id="overwrite"
        checked={options.overwrite ?? false}
        headingID="dialogs.scene_gen.overwrite"
        onChange={(v) => setOptions({ overwrite: v })}
      />
    </>
  );
};
