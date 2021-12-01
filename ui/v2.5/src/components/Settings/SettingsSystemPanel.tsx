import React from "react";
import { useIntl } from "react-intl";
import { Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator } from "src/components/Shared";
import { SettingSection } from "./SettingSection";
import {
  BooleanSetting,
  ModalSetting,
  NumberSetting,
  SelectSetting,
  StringSetting,
} from "./Inputs";
import { SettingStateContext } from "./context";

type VideoPreviewSettingsInput = Pick<
  GQL.ConfigGeneralInput,
  | "previewSegments"
  | "previewSegmentDuration"
  | "previewExcludeStart"
  | "previewExcludeEnd"
>;

interface IVideoPreviewInput {
  value: VideoPreviewSettingsInput;
  setValue: (v: VideoPreviewSettingsInput) => void;
}

const VideoPreviewInput: React.FC<IVideoPreviewInput> = ({
  value,
  setValue,
}) => {
  const intl = useIntl();

  function set(v: Partial<VideoPreviewSettingsInput>) {
    setValue({
      ...value,
      ...v,
    });
  }

  const {
    previewSegments,
    previewSegmentDuration,
    previewExcludeStart,
    previewExcludeEnd,
  } = value;

  return (
    <div>
      <Form.Group id="preview-segments">
        <h6>
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_seg_count_head",
          })}
        </h6>
        <Form.Control
          className="text-input"
          type="number"
          value={previewSegments?.toString() ?? 0}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            set({
              previewSegments: Number.parseInt(
                e.currentTarget.value || "0",
                10
              ),
            })
          }
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_seg_count_desc",
          })}
        </Form.Text>
      </Form.Group>

      <Form.Group id="preview-segment-duration">
        <h6>
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_seg_duration_head",
          })}
        </h6>
        <Form.Control
          className="text-input"
          type="number"
          value={previewSegmentDuration?.toString() ?? 0}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            set({
              previewSegmentDuration: Number.parseFloat(
                e.currentTarget.value || "0"
              ),
            })
          }
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_seg_duration_desc",
          })}
        </Form.Text>
      </Form.Group>

      <Form.Group id="preview-exclude-start">
        <h6>
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_exclude_start_time_head",
          })}
        </h6>
        <Form.Control
          className="text-input"
          value={previewExcludeStart ?? ""}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            set({ previewExcludeStart: e.currentTarget.value })
          }
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_exclude_start_time_desc",
          })}
        </Form.Text>
      </Form.Group>

      <Form.Group id="preview-exclude-start">
        <h6>
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_exclude_end_time_head",
          })}
        </h6>
        <Form.Control
          className="text-input"
          value={previewExcludeEnd ?? ""}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            set({ previewExcludeEnd: e.currentTarget.value })
          }
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "dialogs.scene_gen.preview_exclude_end_time_desc",
          })}
        </Form.Text>
      </Form.Group>
    </div>
  );
};

export const SettingsConfigurationPanel: React.FC = () => {
  const { general, loading, error, saveGeneral } = React.useContext(
    SettingStateContext
  );

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

  if (error) return <h1>{error.message}</h1>;
  if (loading) return <LoadingIndicator />;

  return (
    <>
      <SettingSection headingID="config.application_paths.heading">
        <StringSetting
          id="database-path"
          headingID="config.general.db_path_head"
          subHeadingID="config.general.sqlite_location"
          value={general.databasePath ?? undefined}
          onChange={(v) => saveGeneral({ databasePath: v })}
        />

        <StringSetting
          id="generated-path"
          headingID="config.general.generated_path_head"
          subHeadingID="config.general.generated_files_location"
          value={general.generatedPath ?? undefined}
          onChange={(v) => saveGeneral({ generatedPath: v })}
        />

        <StringSetting
          id="metadata-path"
          headingID="config.general.metadata_path.heading"
          subHeadingID="config.general.metadata_path.description"
          value={general.metadataPath ?? undefined}
          onChange={(v) => saveGeneral({ metadataPath: v })}
        />

        <StringSetting
          id="cache-path"
          headingID="config.general.cache_path_head"
          subHeadingID="config.general.cache_location"
          value={general.cachePath ?? undefined}
          onChange={(v) => saveGeneral({ cachePath: v })}
        />

        <StringSetting
          id="custom-performer-image-location"
          headingID="config.ui.performers.options.image_location.heading"
          subHeadingID="config.ui.performers.options.image_location.description"
          value={general.customPerformerImageLocation ?? undefined}
          onChange={(v) => saveGeneral({ customPerformerImageLocation: v })}
        />
      </SettingSection>

      <SettingSection headingID="config.general.hashing">
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
      </SettingSection>
    </>
  );
};
