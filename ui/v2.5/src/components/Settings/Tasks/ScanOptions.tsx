import React from "react";
import * as GQL from "src/core/generated-graphql";
import { BooleanSetting } from "../Inputs";

interface IScanOptions {
  options: GQL.ScanMetadataInput;
  setOptions: (s: GQL.ScanMetadataInput) => void;
}

export const ScanOptions: React.FC<IScanOptions> = ({
  options,
  setOptions: setOptionsState,
}) => {
  const {
    useFileMetadata,
    stripFileExtension,
    scanGeneratePreviews,
    scanGenerateImagePreviews,
    scanGenerateSprites,
    scanGeneratePhashes,
    scanGenerateThumbnails,
  } = options;

  function setOptions(input: Partial<GQL.ScanMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  return (
    <>
      <BooleanSetting
        id="scan-generate-previews"
        headingID="config.tasks.generate_video_previews_during_scan"
        tooltipID="config.tasks.generate_video_previews_during_scan_tooltip"
        checked={scanGeneratePreviews ?? false}
        onChange={(v) => setOptions({ scanGeneratePreviews: v })}
      />
      <BooleanSetting
        id="scan-generate-image-previews"
        className="sub-setting"
        headingID="config.tasks.generate_previews_during_scan"
        tooltipID="config.tasks.generate_previews_during_scan_tooltip"
        checked={scanGenerateImagePreviews ?? false}
        disabled={!scanGeneratePreviews}
        onChange={(v) => setOptions({ scanGenerateImagePreviews: v })}
      />

      <BooleanSetting
        id="scan-generate-sprites"
        headingID="config.tasks.generate_sprites_during_scan"
        checked={scanGenerateSprites ?? false}
        onChange={(v) => setOptions({ scanGenerateSprites: v })}
      />
      <BooleanSetting
        id="scan-generate-phashes"
        checked={scanGeneratePhashes ?? false}
        headingID="config.tasks.generate_phashes_during_scan"
        tooltipID="config.tasks.generate_phashes_during_scan_tooltip"
        onChange={(v) => setOptions({ scanGeneratePhashes: v })}
      />
      <BooleanSetting
        id="scan-generate-thumbnails"
        checked={scanGenerateThumbnails ?? false}
        headingID="config.tasks.generate_thumbnails_during_scan"
        onChange={(v) => setOptions({ scanGenerateThumbnails: v })}
      />
      <BooleanSetting
        id="strip-file-extension"
        checked={stripFileExtension ?? false}
        headingID="config.tasks.dont_include_file_extension_as_part_of_the_title"
        onChange={(v) => setOptions({ stripFileExtension: v })}
      />
      <BooleanSetting
        id="use-file-metadata"
        checked={useFileMetadata ?? false}
        headingID="config.tasks.set_name_date_details_from_metadata_if_present"
        onChange={(v) => setOptions({ useFileMetadata: v })}
      />
    </>
  );
};
