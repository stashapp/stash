import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Form, InputGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import {
  useConfiguration,
  useConfigureGeneral,
  useGenerateAPIKey,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { Icon, LoadingIndicator } from "src/components/Shared";
import {
  StashBoxConfiguration,
  IStashBoxInstance,
} from "./StashBoxConfiguration";
import StashConfiguration from "./StashConfiguration";
import { StringListInput } from "../Shared/StringListInput";
import { SettingGroup } from "./SettingGroup";
import {
  BooleanSetting,
  NumberSetting,
  SelectSetting,
  StringSetting,
} from "./Inputs";
import { debounce } from "lodash";

export const SettingsConfigurationPanel: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();
  // Editing config state
  const [stashes, setStashes] = useState<GQL.StashConfig[]>([]);
  const [previewSegments, setPreviewSegments] = useState<number>(0);
  const [previewSegmentDuration, setPreviewSegmentDuration] = useState<number>(
    0
  );
  const [previewExcludeStart, setPreviewExcludeStart] = useState<
    string | undefined
  >(undefined);
  const [previewExcludeEnd, setPreviewExcludeEnd] = useState<
    string | undefined
  >(undefined);
  const [username, setUsername] = useState<string | undefined>(undefined);
  const [password, setPassword] = useState<string | undefined>(undefined);
  const [maxSessionAge, setMaxSessionAge] = useState<number>(0);
  const [trustedProxies, setTrustedProxies] = useState<string[] | undefined>(
    undefined
  );

  const [excludes, setExcludes] = useState<string[]>([]);
  const [imageExcludes, setImageExcludes] = useState<string[]>([]);
  const [stashBoxes, setStashBoxes] = useState<IStashBoxInstance[]>([]);

  const [toSave, setToSave] = useState<GQL.ConfigGeneralInput | undefined>();

  const { data, error, loading } = useConfiguration();
  const [configuration, setConfiguration] = useState<
    GQL.ConfigGeneralInput | undefined
  >();

  const [generateAPIKey] = useGenerateAPIKey();

  const [updateGeneralConfig] = useConfigureGeneral();

  // saves the configuration if no further changes are made after a half second
  const saveConfig = debounce(async (input: GQL.ConfigGeneralInput) => {
    try {
      await updateGeneralConfig({
        variables: {
          input,
        },
      });
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl
              .formatMessage({ id: "configuration" })
              .toLocaleLowerCase(),
          }
        ),
      });
      setToSave(undefined);
    } catch (e) {
      Toast.error(e);
    }
  }, 500);

  useEffect(() => {
    if (!toSave) {
      return;
    }

    saveConfig(toSave);
  }, [toSave, saveConfig]);

  useEffect(() => {
    if (!data?.configuration || error) return;

    const { general } = data.configuration;

    setConfiguration({
      databasePath: general.databasePath,
      generatedPath: general.generatedPath,
      metadataPath: general.metadataPath,
      cachePath: general.cachePath,
      createGalleriesFromFolders: general.createGalleriesFromFolders,
      parallelTasks: general.parallelTasks,
      writeImageThumbnails: general.writeImageThumbnails,
      calculateMD5: general.calculateMD5,
      previewAudio: general.previewAudio,
      logOut: general.logOut,
      logAccess: general.logAccess,
      videoFileNamingAlgorithm: general.videoFileNamingAlgorithm,
      maxTranscodeSize: general.maxTranscodeSize ?? undefined,
      maxStreamingTranscodeSize: general.maxStreamingTranscodeSize ?? undefined,
      previewPreset: general.previewPreset,
      logLevel: general.logLevel,
      videoExtensions: general.videoExtensions,
      imageExtensions: general.imageExtensions,
      galleryExtensions: general.galleryExtensions,
      logFile: general.logFile,
      customPerformerImageLocation: general.customPerformerImageLocation,
    });

    const conf = data.configuration;
    if (conf.general) {
      setStashes(conf.general.stashes ?? []);

      setPreviewSegments(conf.general.previewSegments);
      setPreviewSegmentDuration(conf.general.previewSegmentDuration);
      setPreviewExcludeStart(conf.general.previewExcludeStart);
      setPreviewExcludeEnd(conf.general.previewExcludeEnd);

      setUsername(conf.general.username);
      setPassword(conf.general.password);
      setMaxSessionAge(conf.general.maxSessionAge);
      setTrustedProxies(conf.general.trustedProxies ?? undefined);

      setExcludes(conf.general.excludes);
      setImageExcludes(conf.general.imageExcludes);

      setStashBoxes(
        conf.general.stashBoxes.map((box, i) => ({
          name: box?.name ?? undefined,
          endpoint: box.endpoint,
          api_key: box.api_key,
          index: i,
        })) ?? []
      );
    }
  }, [data, error]);

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

  async function onGenerateAPIKey() {
    try {
      await generateAPIKey({
        variables: {
          input: {},
        },
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onClearAPIKey() {
    try {
      await generateAPIKey({
        variables: {
          input: {
            clear: true,
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  function save(input: Partial<GQL.ConfigGeneralInput>) {
    if (!configuration) {
      return;
    }

    setConfiguration({
      ...configuration,
      ...input,
    });

    setToSave((current) => {
      if (!current) {
        return input;
      }
      return {
        ...current,
        ...input,
      };
    });
  }

  async function onSave() {
    try {
      const result = await updateGeneralConfig({
        variables: {
          input: {
            stashes: stashes.map((s) => ({
              path: s.path,
              excludeVideo: s.excludeVideo,
              excludeImage: s.excludeImage,
            })),
            previewSegments,
            previewSegmentDuration,
            previewExcludeStart,
            previewExcludeEnd,
            username,
            password,
            maxSessionAge,
            trustedProxies,
            excludes,
            imageExcludes,
            stashBoxes: stashBoxes.map(
              (b) =>
                ({
                  name: b?.name ?? "",
                  api_key: b?.api_key ?? "",
                  endpoint: b?.endpoint ?? "",
                } as GQL.StashBoxInput)
            ),
          },
        },
      });
      // eslint-disable-next-line no-console
      console.log(result);
      Toast.success({
        content: intl.formatMessage(
          { id: "toast.updated_entity" },
          {
            entity: intl
              .formatMessage({ id: "configuration" })
              .toLocaleLowerCase(),
          }
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

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
  if (!data?.configuration || !configuration || loading)
    return <LoadingIndicator />;

  const general = configuration;

  return (
    <>
      <SettingGroup
        headingID="library"
        subHeadingID="config.general.directory_locations_to_your_content"
      >
        <StashConfiguration
          stashes={stashes}
          setStashes={(s) => setStashes(s)}
        />
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "config.general.directory_locations_to_your_content",
          })}
        </Form.Text>
      </SettingGroup>

      <SettingGroup headingID="config.application_paths.heading">
        <StringSetting
          id="database-path"
          headingID="config.general.db_path_head"
          subHeadingID="config.general.sqlite_location"
          value={general.databasePath ?? undefined}
          onChange={(v) => save({ databasePath: v })}
        />

        <StringSetting
          id="generated-path"
          headingID="config.general.generated_path_head"
          subHeadingID="config.general.generated_files_location"
          value={general.generatedPath ?? undefined}
          onChange={(v) => save({ generatedPath: v })}
        />

        <StringSetting
          id="metadata-path"
          headingID="config.general.metadata_path.heading"
          subHeadingID="config.general.metadata_path.description"
          value={general.metadataPath ?? undefined}
          onChange={(v) => save({ metadataPath: v })}
        />

        <StringSetting
          id="cache-path"
          headingID="config.general.cache_path_head"
          subHeadingID="config.general.cache_location"
          value={general.cachePath ?? undefined}
          onChange={(v) => save({ cachePath: v })}
        />

        <StringSetting
          id="custom-performer-image-location"
          headingID="config.ui.performers.options.image_location.heading"
          subHeadingID="config.ui.performers.options.image_location.description"
          value={general.customPerformerImageLocation ?? undefined}
          onChange={(v) => save({ customPerformerImageLocation: v })}
        />
      </SettingGroup>

      <SettingGroup headingID="config.library.media_content_extensions">
        <StringSetting
          id="video-extensions"
          headingID="config.general.video_ext_head"
          subHeadingID="config.general.video_ext_desc"
          value={listToCommaDelimited(general.videoExtensions ?? undefined)}
          onChange={(v) => save({ videoExtensions: commaDelimitedToList(v) })}
        />

        <StringSetting
          id="image-extensions"
          headingID="config.general.image_ext_head"
          subHeadingID="config.general.image_ext_desc"
          value={listToCommaDelimited(general.imageExtensions ?? undefined)}
          onChange={(v) => save({ imageExtensions: commaDelimitedToList(v) })}
        />

        <StringSetting
          id="gallery-extensions"
          headingID="config.general.gallery_ext_head"
          subHeadingID="config.general.gallery_ext_desc"
          value={listToCommaDelimited(general.galleryExtensions ?? undefined)}
          onChange={(v) => save({ galleryExtensions: commaDelimitedToList(v) })}
        />
      </SettingGroup>

      <SettingGroup headingID="config.library.exclusions">
        <Form.Group>
          <h6>
            {intl.formatMessage({
              id: "config.general.excluded_video_patterns_head",
            })}
          </h6>
          <StringListInput
            className="w-50"
            value={excludes}
            setValue={setExcludes}
            defaultNewValue="sample\.mp4$"
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id: "config.general.excluded_video_patterns_desc",
            })}
            <a
              href="https://github.com/stashapp/stash/wiki/Exclude-file-configuration"
              rel="noopener noreferrer"
              target="_blank"
            >
              <Icon icon="question-circle" />
            </a>
          </Form.Text>
        </Form.Group>

        <Form.Group>
          <h6>
            {intl.formatMessage({
              id: "config.general.excluded_image_gallery_patterns_head",
            })}
          </h6>
          <StringListInput
            className="w-50"
            value={imageExcludes}
            setValue={setImageExcludes}
            defaultNewValue="sample\.jpg$"
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id: "config.general.excluded_image_gallery_patterns_desc",
            })}
            <a
              href="https://github.com/stashapp/stash/wiki/Exclude-file-configuration"
              rel="noopener noreferrer"
              target="_blank"
            >
              <Icon icon="question-circle" />
            </a>
          </Form.Text>
        </Form.Group>
      </SettingGroup>

      <SettingGroup headingID="config.library.gallery_and_image_options">
        <BooleanSetting
          id="create-galleries-from-folders"
          headingID="config.general.create_galleries_from_folders_label"
          subHeadingID="config.general.create_galleries_from_folders_desc"
          checked={general.createGalleriesFromFolders ?? false}
          onChange={(v) => save({ createGalleriesFromFolders: v })}
        />

        <BooleanSetting
          id="write-image-thumbnails"
          headingID="config.ui.images.options.write_image_thumbnails.heading"
          subHeadingID="config.ui.images.options.write_image_thumbnails.description"
          checked={general.writeImageThumbnails ?? false}
          onChange={(v) => save({ writeImageThumbnails: v })}
        />
      </SettingGroup>

      <SettingGroup headingID="config.general.hashing">
        <BooleanSetting
          id="calculate-md5-and-ohash"
          headingID="config.general.calculate_md5_and_ohash_label"
          subHeadingID="config.general.calculate_md5_and_ohash_desc"
          checked={general.calculateMD5 ?? false}
          onChange={(v) => save({ calculateMD5: v })}
        />

        <SelectSetting
          id="generated_file_naming_hash"
          headingID="config.general.generated_file_naming_hash_head"
          subHeadingID="config.general.generated_file_naming_hash_desc"
          value={namingHashToString(
            general.videoFileNamingAlgorithm ?? undefined
          )}
          onChange={(v) =>
            save({ videoFileNamingAlgorithm: translateNamingHash(v) })
          }
        >
          {namingHashAlgorithms.map((q) => (
            <option key={q} value={q}>
              {q}
            </option>
          ))}
        </SelectSetting>
      </SettingGroup>

      <SettingGroup headingID="config.system.transcoding">
        <SelectSetting
          id="transcode-size"
          headingID="config.general.maximum_transcode_size_head"
          subHeadingID="config.general.maximum_transcode_size_desc"
          onChange={(v) => save({ maxTranscodeSize: translateQuality(v) })}
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
            save({ maxStreamingTranscodeSize: translateQuality(v) })
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
      </SettingGroup>

      <SettingGroup headingID="config.general.parallel_scan_head">
        <NumberSetting
          id="parallel-tasks"
          headingID="config.general.number_of_parallel_task_for_scan_generation_head"
          subHeadingID="config.general.number_of_parallel_task_for_scan_generation_desc"
          value={general.parallelTasks ?? undefined}
          onChange={(v) => save({ parallelTasks: v })}
        />
      </SettingGroup>

      <SettingGroup headingID="config.general.preview_generation">
        <SelectSetting
          id="scene-gen-preview-preset"
          headingID="dialogs.scene_gen.preview_preset_head"
          subHeadingID="dialogs.scene_gen.preview_preset_desc"
          value={general.previewPreset ?? undefined}
          onChange={(v) =>
            save({ previewPreset: (v as GQL.PreviewPreset) ?? undefined })
          }
        >
          {Object.keys(GQL.PreviewPreset).map((p) => (
            <option value={p.toLowerCase()} key={p}>
              {p}
            </option>
          ))}
        </SelectSetting>

        <Form.Group>
          <BooleanSetting
            id="preview-include-audio"
            headingID="config.general.include_audio_head"
            subHeadingID="config.general.include_audio_desc"
            checked={general.previewAudio ?? false}
            onChange={(v) => save({ previewAudio: v })}
          />
        </Form.Group>

        <Form.Group id="preview-segments">
          <h6>
            {intl.formatMessage({
              id: "dialogs.scene_gen.preview_seg_count_head",
            })}
          </h6>
          <Form.Control
            className="col col-sm-6 text-input"
            type="number"
            value={previewSegments.toString()}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setPreviewSegments(
                Number.parseInt(e.currentTarget.value || "0", 10)
              )
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
            className="col col-sm-6 text-input"
            type="number"
            value={previewSegmentDuration.toString()}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setPreviewSegmentDuration(
                Number.parseFloat(e.currentTarget.value || "0")
              )
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
            className="col col-sm-6 text-input"
            defaultValue={previewExcludeStart}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setPreviewExcludeStart(e.currentTarget.value)
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
            className="col col-sm-6 text-input"
            defaultValue={previewExcludeEnd}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setPreviewExcludeEnd(e.currentTarget.value)
            }
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id: "dialogs.scene_gen.preview_exclude_end_time_desc",
            })}
          </Form.Text>
        </Form.Group>
      </SettingGroup>

      <SettingGroup headingID="config.general.auth.stash-box_integration">
        <StashBoxConfiguration boxes={stashBoxes} saveBoxes={setStashBoxes} />
      </SettingGroup>

      <SettingGroup headingID="config.general.auth.authentication">
        <Form.Group id="username">
          <h6>{intl.formatMessage({ id: "config.general.auth.username" })}</h6>
          <Form.Control
            className="col col-sm-6 text-input"
            defaultValue={username}
            onInput={(e: React.FormEvent<HTMLInputElement>) =>
              setUsername(e.currentTarget.value)
            }
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({ id: "config.general.auth.username_desc" })}
          </Form.Text>
        </Form.Group>
        <Form.Group id="password">
          <h6>{intl.formatMessage({ id: "config.general.auth.password" })}</h6>
          <Form.Control
            className="col col-sm-6 text-input"
            type="password"
            defaultValue={password}
            onInput={(e: React.FormEvent<HTMLInputElement>) =>
              setPassword(e.currentTarget.value)
            }
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({ id: "config.general.auth.password_desc" })}
          </Form.Text>
        </Form.Group>

        <Form.Group id="apikey">
          <h6>{intl.formatMessage({ id: "config.general.auth.api_key" })}</h6>
          <InputGroup>
            <Form.Control
              className="col col-sm-6 text-input"
              value={data.configuration.general.apiKey}
              readOnly
            />
            <InputGroup.Append>
              <Button
                className=""
                title={intl.formatMessage({
                  id: "config.general.auth.generate_api_key",
                })}
                onClick={() => onGenerateAPIKey()}
              >
                <Icon icon="redo" />
              </Button>
              <Button
                className=""
                variant="danger"
                title={intl.formatMessage({
                  id: "config.general.auth.clear_api_key",
                })}
                onClick={() => onClearAPIKey()}
              >
                <Icon icon="minus" />
              </Button>
            </InputGroup.Append>
          </InputGroup>
          <Form.Text className="text-muted">
            {intl.formatMessage({ id: "config.general.auth.api_key_desc" })}
          </Form.Text>
        </Form.Group>

        <Form.Group id="maxSessionAge">
          <h6>
            {intl.formatMessage({
              id: "config.general.auth.maximum_session_age",
            })}
          </h6>
          <Form.Control
            className="col col-sm-6 text-input"
            type="number"
            value={maxSessionAge.toString()}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setMaxSessionAge(
                Number.parseInt(e.currentTarget.value || "0", 10)
              )
            }
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id: "config.general.auth.maximum_session_age_desc",
            })}
          </Form.Text>
        </Form.Group>

        <Form.Group id="trusted-proxies">
          <h6>
            {intl.formatMessage({ id: "config.general.auth.trusted_proxies" })}
          </h6>
          <StringListInput
            value={trustedProxies ?? []}
            setValue={(value) => setTrustedProxies(value)}
            defaultNewValue=""
          />
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id: "config.general.auth.trusted_proxies_desc",
            })}
          </Form.Text>
        </Form.Group>
      </SettingGroup>

      <SettingGroup headingID="config.general.logging">
        <StringSetting
          headingID="config.general.auth.log_file"
          subHeadingID="config.general.auth.log_file_desc"
          value={general.logFile ?? undefined}
          onChange={(v) => save({ logFile: v })}
        />

        <BooleanSetting
          id="log-terminal"
          headingID="config.general.auth.log_to_terminal"
          subHeadingID="config.general.auth.log_to_terminal_desc"
          checked={general.logOut ?? false}
          onChange={(v) => save({ logOut: v })}
        />

        <SelectSetting
          id="log-level"
          headingID="config.logs.log_level"
          onChange={(v) => save({ logLevel: v })}
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
          onChange={(v) => save({ logAccess: v })}
        />
      </SettingGroup>

      <hr />

      <Button variant="primary" onClick={() => onSave()}>
        <FormattedMessage id="actions.save" />
      </Button>
    </>
  );
};
