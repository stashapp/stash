import React from "react";
import * as GQL from "src/core/generated-graphql";
import { BooleanSetting, ModalSetting } from "../Inputs";
import { DurationInput } from "../../Shared/DurationInput";
import TextUtils from "../../../utils/text";

interface IScanOptions {
  options: GQL.ScanMetadataInput;
  filters: GQL.ScanMetaDataFilterInput;
  setOptions: (s: GQL.ScanMetadataInput) => void;
  setFilters: (s: GQL.ScanMetaDataFilterInput) => void;
}

export const ScanOptions: React.FC<IScanOptions> = ({
  options,
  filters,
  setOptions: setOptionsState,
  setFilters: setFiltersState,
}) => {
  const {
    scanGenerateCovers,
    scanGeneratePreviews,
    scanGenerateImagePreviews,
    scanGenerateSprites,
    scanGeneratePhashes,
    scanGenerateThumbnails,
    scanGenerateClipPreviews,
    rescan,
  } = options;

  const { minDuration } = filters;

  function setOptions(input: Partial<GQL.ScanMetadataInput>) {
    setOptionsState({ ...options, ...input });
  }

  function setFilters(input: Partial<GQL.ScanMetaDataFilterInput>) {
    setFiltersState({ ...filters, ...input });
  }

  return (
    <>
      <BooleanSetting
        id="scan-generate-covers"
        headingID="config.tasks.generate_video_covers_during_scan"
        checked={scanGenerateCovers ?? true}
        onChange={(v) => setOptions({ scanGenerateCovers: v })}
      />
      <BooleanSetting
        id="scan-generate-previews"
        headingID="config.tasks.generate_video_previews_during_scan"
        tooltipID="config.tasks.generate_video_previews_during_scan_tooltip"
        checked={scanGeneratePreviews ?? false}
        onChange={(v) => setOptions({ scanGeneratePreviews: v })}
      />
      <BooleanSetting
        advanced
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
        tooltipID="config.tasks.generate_sprites_during_scan_tooltip"
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
        id="scan-generate-clip-previews"
        checked={scanGenerateClipPreviews ?? false}
        headingID="config.tasks.generate_clip_previews_during_scan"
        onChange={(v) => setOptions({ scanGenerateClipPreviews: v })}
      />
      <BooleanSetting
        id="force-rescan"
        headingID="config.tasks.rescan"
        tooltipID="config.tasks.rescan_tooltip"
        checked={rescan ?? false}
        onChange={(v) => setOptions({ rescan: v })}
      />
      <ModalSetting<number>
        id="duration-filter"
        headingID="config.tasks.filter_duration_during_scan"
        tooltipID="config.tasks.filter_duration_during_scan_tooltip"
        value={minDuration ?? undefined}
        onChange={(v) => setFilters({ minDuration: v })}
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
    </>
  );
};
