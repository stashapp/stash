import React from "react";
import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { SettingSection } from "./SettingSection";
import {
  BooleanSetting,
  ModalSetting,
  NumberSetting,
  SelectSetting,
  Setting,
  StringListSetting,
  StringSetting,
} from "./Inputs";
import { useSettings } from "./context";
import {
  VideoPreviewInput,
  VideoPreviewSettingsInput,
} from "./GeneratePreviewOptions";
import { FormattedMessage, useIntl } from "react-intl";
import { Button } from "react-bootstrap";
import { useToast } from "src/hooks/Toast";
import { useHistory } from "react-router-dom";

export const SettingsConfigurationPanel: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();
  const history = useHistory();

  const { general, loading, error, saveGeneral } = useSettings();
  const [mutateDownloadFFMpeg] = GQL.useDownloadFfMpegMutation();

  const transcodeQualities = [
    GQL.StreamingResolutionEnum.Low,
    GQL.StreamingResolutionEnum.Standard,
    GQL.StreamingResolutionEnum.StandardHd,
    GQL.StreamingResolutionEnum.FullHd,
    GQL.StreamingResolutionEnum.FourK,
    GQL.StreamingResolutionEnum.Original,
  ].map(resolutionToString);

  function resolutionToString(r: GQL.StreamingResolutionEnum | undefined) {
    switch (r) {
      case GQL.StreamingResolutionEnum.Low:
        return "240p";
      case GQL.StreamingResolutionEnum.Standard:
        return "480p";
      case GQL.StreamingResolutionEnum.StandardHd:
        return "720p";
      case GQL.StreamingResolutionEnum.FullHd:
        return "1080p";
      case GQL.StreamingResolutionEnum.FourK:
        return "4k";
      case GQL.StreamingResolutionEnum.Original:
        return "Original";
    }

    return "Original";
  }

  function translateQuality(quality: string) {
    switch (quality) {
      case "240p":
        return GQL.StreamingResolutionEnum.Low;
      case "480p":
        return GQL.StreamingResolutionEnum.Standard;
      case "720p":
        return GQL.StreamingResolutionEnum.StandardHd;
      case "1080p":
        return GQL.StreamingResolutionEnum.FullHd;
      case "4k":
        return GQL.StreamingResolutionEnum.FourK;
      case "Original":
        return GQL.StreamingResolutionEnum.Original;
    }

    return GQL.StreamingResolutionEnum.Original;
  }

  const namingHashAlgorithms = [
    GQL.HashAlgorithm.Md5,
    GQL.HashAlgorithm.Oshash,
  ].map(namingHashToString);

  function namingHashToString(value: GQL.HashAlgorithm | undefined) {
    switch (value) {
      case GQL.HashAlgorithm.Oshash:
        return "oshash";
      case GQL.HashAlgorithm.Md5:
        return "MD5";
    }

    return "MD5";
  }

  function translateNamingHash(value: string) {
    switch (value) {
      case "oshash":
        return GQL.HashAlgorithm.Oshash;
      case "MD5":
        return GQL.HashAlgorithm.Md5;
    }

    return GQL.HashAlgorithm.Md5;
  }

  function blobStorageTypeToID(value: GQL.BlobsStorageType | undefined) {
    switch (value) {
      case GQL.BlobsStorageType.Database:
        return "blobs_storage_type.database";
      case GQL.BlobsStorageType.Filesystem:
        return "blobs_storage_type.filesystem";
    }

    return "blobs_storage_type.database";
  }

  async function onDownloadFFMpeg() {
    try {
      await mutateDownloadFFMpeg();
      // navigate to tasks page to see the progress
      history.push("/settings?tab=tasks");
    } catch (e) {
      Toast.error(e);
    }
  }

  if (error) return <h1>{error.message}</h1>;
  if (loading) return <LoadingIndicator />;

  return (
    <>
      <SettingSection headingID="config.application_paths.heading">
        <StringSetting
          id="generated-path"
          headingID="config.general.generated_path_head"
          subHeadingID="config.general.generated_files_location"
          value={general.generatedPath ?? undefined}
          onChange={(v) => saveGeneral({ generatedPath: v })}
        />

        <StringSetting
          id="cache-path"
          headingID="config.general.cache_path_head"
          subHeadingID="config.general.cache_location"
          value={general.cachePath ?? undefined}
          onChange={(v) => saveGeneral({ cachePath: v })}
        />

        <StringSetting
          id="scrapers-path"
          headingID="config.general.scrapers_path.heading"
          subHeadingID="config.general.scrapers_path.description"
          value={general.scrapersPath ?? undefined}
          onChange={(v) => saveGeneral({ scrapersPath: v })}
        />

        <StringSetting
          id="plugins-path"
          headingID="config.general.plugins_path.heading"
          subHeadingID="config.general.plugins_path.description"
          value={general.pluginsPath ?? undefined}
          onChange={(v) => saveGeneral({ pluginsPath: v })}
        />

        <StringSetting
          id="metadata-path"
          headingID="config.general.metadata_path.heading"
          subHeadingID="config.general.metadata_path.description"
          value={general.metadataPath ?? undefined}
          onChange={(v) => saveGeneral({ metadataPath: v })}
        />

        <StringSetting
          id="custom-performer-image-location"
          headingID="config.ui.performers.options.image_location.heading"
          subHeadingID="config.ui.performers.options.image_location.description"
          value={general.customPerformerImageLocation ?? undefined}
          onChange={(v) => saveGeneral({ customPerformerImageLocation: v })}
        />

        <StringSetting
          id="ffmpeg-path"
          headingID="config.general.ffmpeg.ffmpeg_path.heading"
          subHeadingID="config.general.ffmpeg.ffmpeg_path.description"
          value={general.ffmpegPath ?? undefined}
          onChange={(v) => saveGeneral({ ffmpegPath: v })}
        />

        <StringSetting
          id="ffprobe-path"
          headingID="config.general.ffmpeg.ffprobe_path.heading"
          subHeadingID="config.general.ffmpeg.ffprobe_path.description"
          value={general.ffprobePath ?? undefined}
          onChange={(v) => saveGeneral({ ffprobePath: v })}
        />

        <Setting
          heading={
            <>
              <FormattedMessage id="config.general.ffmpeg.download_ffmpeg.heading" />
            </>
          }
          subHeadingID="config.general.ffmpeg.download_ffmpeg.description"
        >
          <Button variant="secondary" onClick={() => onDownloadFFMpeg()}>
            <FormattedMessage id="config.general.ffmpeg.download_ffmpeg.heading" />
          </Button>
        </Setting>

        <StringSetting
          id="python-path"
          headingID="config.general.python_path.heading"
          subHeadingID="config.general.python_path.description"
          value={general.pythonPath ?? undefined}
          onChange={(v) => saveGeneral({ pythonPath: v })}
        />

        <StringSetting
          id="backup-directory-path"
          headingID="config.general.backup_directory_path.heading"
          subHeadingID="config.general.backup_directory_path.description"
          value={general.backupDirectoryPath ?? undefined}
          onChange={(v) => saveGeneral({ backupDirectoryPath: v })}
        />

        <StringSetting
          id="delete-trash-path"
          headingID="config.general.delete_trash_path.heading"
          subHeadingID="config.general.delete_trash_path.description"
          value={general.deleteTrashPath ?? undefined}
          onChange={(v) => saveGeneral({ deleteTrashPath: v })}
        />
      </SettingSection>

      <SettingSection headingID="config.general.database">
        <StringSetting
          id="database-path"
          headingID="config.general.db_path_head"
          subHeadingID="config.general.sqlite_location"
          value={general.databasePath ?? undefined}
          onChange={(v) => saveGeneral({ databasePath: v })}
        />
        <SelectSetting
          id="blobs-storage"
          headingID="config.general.blobs_storage.heading"
          subHeadingID="config.general.blobs_storage.description"
          value={general.blobsStorage ?? GQL.BlobsStorageType.Database}
          onChange={(v) =>
            saveGeneral({ blobsStorage: v as GQL.BlobsStorageType })
          }
        >
          {Object.values(GQL.BlobsStorageType).map((q) => (
            <option key={q} value={q}>
              {intl.formatMessage({ id: blobStorageTypeToID(q) })}
            </option>
          ))}
        </SelectSetting>
        <StringSetting
          id="blobs-path"
          headingID="config.general.blobs_path.heading"
          subHeadingID="config.general.blobs_path.description"
          value={general.blobsPath ?? ""}
          onChange={(v) => saveGeneral({ blobsPath: v })}
        />
      </SettingSection>

      <SettingSection advanced headingID="config.general.hashing">
        <BooleanSetting
          id="calculate-md5-and-ohash"
          headingID="config.general.calculate_md5_and_ohash_label"
          subHeadingID="config.general.calculate_md5_and_ohash_desc"
          checked={general.calculateMD5 ?? false}
          onChange={(v) => saveGeneral({ calculateMD5: v })}
        />

        <SelectSetting
          id="generated_file_naming_hash"
          headingID="config.general.generated_file_naming_hash_head"
          subHeadingID="config.general.generated_file_naming_hash_desc"
          value={namingHashToString(
            general.videoFileNamingAlgorithm ?? undefined
          )}
          onChange={(v) =>
            saveGeneral({ videoFileNamingAlgorithm: translateNamingHash(v) })
          }
        >
          {namingHashAlgorithms.map((q) => (
            <option key={q} value={q}>
              {q}
            </option>
          ))}
        </SelectSetting>
      </SettingSection>

      <SettingSection headingID="config.system.transcoding">
        <SelectSetting
          advanced
          id="transcode-size"
          headingID="config.general.maximum_transcode_size_head"
          subHeadingID="config.general.maximum_transcode_size_desc"
          onChange={(v) =>
            saveGeneral({ maxTranscodeSize: translateQuality(v) })
          }
          value={resolutionToString(general.maxTranscodeSize ?? undefined)}
        >
          {transcodeQualities.map((q) => (
            <option key={q} value={q}>
              {q}
            </option>
          ))}
        </SelectSetting>

        <SelectSetting
          id="streaming-transcode-size"
          headingID="config.general.maximum_streaming_transcode_size_head"
          subHeadingID="config.general.maximum_streaming_transcode_size_desc"
          onChange={(v) =>
            saveGeneral({ maxStreamingTranscodeSize: translateQuality(v) })
          }
          value={resolutionToString(
            general.maxStreamingTranscodeSize ?? undefined
          )}
        >
          {transcodeQualities.map((q) => (
            <option key={q} value={q}>
              {q}
            </option>
          ))}
        </SelectSetting>

        <BooleanSetting
          id="hardware-encoding"
          headingID="config.general.ffmpeg.hardware_acceleration.heading"
          subHeadingID="config.general.ffmpeg.hardware_acceleration.desc"
          checked={general.transcodeHardwareAcceleration ?? false}
          onChange={(v) => saveGeneral({ transcodeHardwareAcceleration: v })}
        />

        <StringListSetting
          advanced
          id="transcode-input-args"
          headingID="config.general.ffmpeg.transcode.input_args.heading"
          subHeadingID="config.general.ffmpeg.transcode.input_args.desc"
          onChange={(v) => saveGeneral({ transcodeInputArgs: v })}
          value={general.transcodeInputArgs ?? []}
        />
        <StringListSetting
          advanced
          id="transcode-output-args"
          headingID="config.general.ffmpeg.transcode.output_args.heading"
          subHeadingID="config.general.ffmpeg.transcode.output_args.desc"
          onChange={(v) => saveGeneral({ transcodeOutputArgs: v })}
          value={general.transcodeOutputArgs ?? []}
        />

        <StringListSetting
          advanced
          id="live-transcode-input-args"
          headingID="config.general.ffmpeg.live_transcode.input_args.heading"
          subHeadingID="config.general.ffmpeg.live_transcode.input_args.desc"
          onChange={(v) => saveGeneral({ liveTranscodeInputArgs: v })}
          value={general.liveTranscodeInputArgs ?? []}
        />
        <StringListSetting
          advanced
          id="live-transcode-output-args"
          headingID="config.general.ffmpeg.live_transcode.output_args.heading"
          subHeadingID="config.general.ffmpeg.live_transcode.output_args.desc"
          onChange={(v) => saveGeneral({ liveTranscodeOutputArgs: v })}
          value={general.liveTranscodeOutputArgs ?? []}
        />
      </SettingSection>

      <SettingSection headingID="config.general.parallel_scan_head">
        <NumberSetting
          id="parallel-tasks"
          headingID="config.general.number_of_parallel_task_for_scan_generation_head"
          subHeadingID="config.general.number_of_parallel_task_for_scan_generation_desc"
          value={general.parallelTasks ?? undefined}
          onChange={(v) => saveGeneral({ parallelTasks: v })}
        />
      </SettingSection>

      <SettingSection headingID="config.general.preview_generation">
        <SelectSetting
          id="scene-gen-preview-preset"
          headingID="dialogs.scene_gen.preview_preset_head"
          subHeadingID="dialogs.scene_gen.preview_preset_desc"
          value={general.previewPreset ?? undefined}
          onChange={(v) =>
            saveGeneral({
              previewPreset: (v as GQL.PreviewPreset) ?? undefined,
            })
          }
        >
          {Object.keys(GQL.PreviewPreset).map((p) => (
            <option value={p.toLowerCase()} key={p}>
              {p}
            </option>
          ))}
        </SelectSetting>

        <BooleanSetting
          id="preview-include-audio"
          headingID="config.general.include_audio_head"
          subHeadingID="config.general.include_audio_desc"
          checked={general.previewAudio ?? false}
          onChange={(v) => saveGeneral({ previewAudio: v })}
        />

        <ModalSetting<VideoPreviewSettingsInput>
          id="video-preview-settings"
          headingID="dialogs.scene_gen.preview_generation_options"
          value={{
            previewExcludeEnd: general.previewExcludeEnd,
            previewExcludeStart: general.previewExcludeStart,
            previewSegmentDuration: general.previewSegmentDuration,
            previewSegments: general.previewSegments,
          }}
          onChange={(v) => saveGeneral(v)}
          renderField={(value, setValue) => (
            <VideoPreviewInput value={value ?? {}} setValue={setValue} />
          )}
          renderValue={() => {
            return <></>;
          }}
        />
      </SettingSection>

      <SettingSection headingID="config.general.heatmap_generation">
        <BooleanSetting
          id="heatmap-draw-range"
          headingID="config.general.funscript_heatmap_draw_range"
          subHeadingID="config.general.funscript_heatmap_draw_range_desc"
          checked={general.drawFunscriptHeatmapRange ?? true}
          onChange={(v) => saveGeneral({ drawFunscriptHeatmapRange: v })}
        />
      </SettingSection>

      <SettingSection headingID="config.general.logging">
        <StringSetting
          headingID="config.general.auth.log_file"
          subHeadingID="config.general.auth.log_file_desc"
          value={general.logFile ?? undefined}
          onChange={(v) => saveGeneral({ logFile: v })}
        />

        <BooleanSetting
          id="log-terminal"
          headingID="config.general.auth.log_to_terminal"
          subHeadingID="config.general.auth.log_to_terminal_desc"
          checked={general.logOut ?? false}
          onChange={(v) => saveGeneral({ logOut: v })}
        />

        <SelectSetting
          id="log-level"
          headingID="config.logs.log_level"
          onChange={(v) => saveGeneral({ logLevel: v })}
          value={general.logLevel ?? undefined}
        >
          {["Trace", "Debug", "Info", "Warning", "Error"].map((o) => (
            <option key={o} value={o}>
              {o}
            </option>
          ))}
        </SelectSetting>

        <BooleanSetting
          id="log-http"
          headingID="config.general.auth.log_http"
          subHeadingID="config.general.auth.log_http_desc"
          checked={general.logAccess ?? false}
          onChange={(v) => saveGeneral({ logAccess: v })}
        />

        <NumberSetting
          id="log-file-max-size"
          headingID="config.general.auth.log_file_max_size"
          subHeadingID="config.general.auth.log_file_max_size_desc"
          value={general.logFileMaxSize ?? 10}
          onChange={(v) => saveGeneral({ logFileMaxSize: v })}
        />
      </SettingSection>
    </>
  );
};
