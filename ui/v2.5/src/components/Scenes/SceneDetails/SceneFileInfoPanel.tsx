import React, { useMemo, useState } from "react";
import { Accordion, Button, Card } from "react-bootstrap";
import {
  FormattedMessage,
  FormattedNumber,
  FormattedTime,
  useIntl,
} from "react-intl";
import { useHistory } from "react-router-dom";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { DeleteFilesDialog } from "src/components/Shared/DeleteFilesDialog";
import { RevealInFilesystemButton } from "src/components/Shared/RevealInFilesystemButton";
import { ReassignFilesDialog } from "src/components/Shared/ReassignFilesDialog";
import * as GQL from "src/core/generated-graphql";
import { mutateSceneSetPrimaryFile } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { useConfigurationContext } from "src/hooks/Config";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { IconField, TextField, URLField, URLsField } from "src/utils/field";
import { StashIDPill } from "src/components/Shared/StashID";
import { PatchComponent } from "../../../patch";
import { FileSize } from "src/components/Shared/FileSize";
import {
  faCameraRotate,
  faClipboard,
  faFileCode,
  faFile,
  faLink,
  faPhotoFilm,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIconProps } from "@fortawesome/react-fontawesome";

const pad2 = (value: number) => value.toString().padStart(2, "0");

function formatCoverTimestamp(seconds: number, frameRate?: number | null) {
  const fps = frameRate ?? 0;
  const roundedSeconds =
    Number.isFinite(fps) && fps > 0
      ? Math.round(seconds * fps) / fps
      : Math.round(seconds * 100) / 100;

  const centiseconds = Math.round(roundedSeconds * 100);
  const totalSeconds = Math.floor(centiseconds / 100);
  const cs = centiseconds % 100;
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const secs = totalSeconds % 60;

  if (hours > 0) {
    return `${hours}:${pad2(minutes)}:${pad2(secs)}.${pad2(cs)}`;
  }

  return `${minutes}:${pad2(secs)}.${pad2(cs)}`;
}

interface IFileInfoPanelProps {
  sceneID: string;
  file: GQL.VideoFileDataFragment;
  coverImageSource?: string | null;
  primary?: boolean;
  ofMany?: boolean;
  onSetPrimaryFile?: () => void;
  onDeleteFile?: () => void;
  onReassign?: () => void;
  loading?: boolean;
}

interface ICoverSourceData {
  icon?: FontAwesomeIconProps["icon"];
  value: string;
  url?: string | null;
  stashEndpoint?: string | null;
  tooltip?: string;
}

const FileInfoPanel: React.FC<IFileInfoPanelProps> = (
  props: IFileInfoPanelProps
) => {
  const intl = useIntl();
  const history = useHistory();
  const { configuration } = useConfigurationContext();
  const coverImageLabel = intl.formatMessage({ id: "scene_cover.image" });
  const coverSourceLabel = intl.formatMessage(
    { id: "scene_cover.source" },
    { image: coverImageLabel }
  );

  // TODO - generalise fingerprints
  const oshash = props.file.fingerprints.find((f) => f.type === "oshash");
  const phash = props.file.fingerprints.find((f) => f.type === "phash");
  const checksum = props.file.fingerprints.find((f) => f.type === "md5");

  const coverSourceField = useMemo(() => {
    const getCoverSourceTooltip = (id: string) =>
      intl.formatMessage({ id }, { image: coverImageLabel });

    const source = props.coverImageSource;
    if (!source) {
      return null;
    }

    let coverSourceData: ICoverSourceData;

    if (source === "default") {
      coverSourceData = {
        icon: faPhotoFilm,
        value: intl.formatMessage({ id: "scene_cover.default" }),
        tooltip: getCoverSourceTooltip("scene_cover.default_tooltip"),
      };
    } else if (source === "clipboard") {
      coverSourceData = {
        icon: faClipboard,
        value: intl.formatMessage({ id: "scene_cover.clipboard" }),
        tooltip: getCoverSourceTooltip("scene_cover.clipboard_tooltip"),
      };
    } else if (source === "userscript") {
      coverSourceData = {
        icon: faFileCode,
        value: intl.formatMessage({ id: "scene_cover.userscript" }),
        tooltip: getCoverSourceTooltip("scene_cover.userscript_tooltip"),
      };
    } else if (source.startsWith("file:")) {
      const fileName = source.slice("file:".length);
      coverSourceData = {
        icon: faFile,
        value: fileName || intl.formatMessage({ id: "scene_cover.file" }),
        tooltip: getCoverSourceTooltip("scene_cover.file_tooltip"),
      };
    } else if (source.startsWith("url:")) {
      const urlValue = source.slice("url:".length).trim();
      coverSourceData = {
        icon: faLink,
        value: urlValue || intl.formatMessage({ id: "actions.from_url" }),
        url: urlValue || null,
        tooltip: getCoverSourceTooltip("scene_cover.url_tooltip"),
      };
    } else if (source.startsWith("timestamp:")) {
      const rawTimestamp = source.slice("timestamp:".length).trim();
      const parsedTimestamp = Number.parseFloat(rawTimestamp);
      const timestampLabel = Number.isFinite(parsedTimestamp)
        ? formatCoverTimestamp(parsedTimestamp, props.file.frame_rate)
        : rawTimestamp;
      coverSourceData = {
        icon: faCameraRotate,
        value: timestampLabel,
        tooltip: getCoverSourceTooltip("scene_cover.timestamp_tooltip"),
      };
    } else if (source.startsWith("stash:")) {
      const endpoint = source.slice("stash:".length).trim();
      const endpointName =
        configuration?.general.stashBoxes.find((sb) => sb.endpoint === endpoint)
          ?.name ?? endpoint;
      coverSourceData = {
        value: endpointName,
        stashEndpoint: endpointName,
        tooltip: getCoverSourceTooltip("scene_cover.stash_tooltip"),
      };
    } else {
      coverSourceData = { value: source };
    }

    return (
      <IconField
        name={coverSourceLabel}
        icon={coverSourceData.icon}
        value={coverSourceData.value}
        url={coverSourceData.url}
        stashEndpoint={coverSourceData.stashEndpoint}
        truncate={Boolean(coverSourceData.url)}
        tooltip={coverSourceData.tooltip}
      />
    );
  }, [
    configuration?.general.stashBoxes,
    coverImageLabel,
    coverSourceLabel,
    intl,
    props.coverImageSource,
    props.file.frame_rate,
  ]);

  function onSplit() {
    history.push(
      `/scenes/new?from_scene_id=${props.sceneID}&file_id=${props.file.id}`
    );
  }

  return (
    <div>
      <dl className="container scene-file-info details-list">
        {props.primary && (
          <>
            <dt></dt>
            <dd className="primary-file">
              <FormattedMessage id="primary_file" />
            </dd>
          </>
        )}
        <TextField
          id="media_info.oshash"
          abbr={intl.formatMessage({ id: "media_info.oshash_meaning" })}
          value={oshash?.value}
          truncate
        />
        <TextField id="media_info.md5" value={checksum?.value} truncate />
        <URLField
          id="media_info.phash"
          abbr={intl.formatMessage({ id: "media_info.phash_meaning" })}
          value={phash?.value}
          url={NavUtils.makeScenesPHashMatchUrl(phash?.value)}
          target="_self"
          truncate
          internal
        />
        {coverSourceField}
        <TextField id="path">
          <span className="d-flex align-items-center">
            <TruncatedText text={props.file.path} />
            <RevealInFilesystemButton fileId={props.file.id} />
          </span>
        </TextField>
        <TextField id="filesize">
          <span className="text-truncate">
            <FileSize size={props.file.size} />
          </span>
        </TextField>
        <TextField id="file_mod_time">
          <FormattedTime
            dateStyle="medium"
            timeStyle="medium"
            value={props.file.mod_time ?? 0}
          />
        </TextField>
        <TextField
          id="duration"
          value={TextUtils.secondsToTimestamp(props.file.duration ?? 0)}
          truncate
        />
        <TextField
          id="dimensions"
          value={`${props.file.width} x ${props.file.height}`}
          truncate
        />
        <TextField id="framerate">
          <FormattedMessage
            id="frames_per_second"
            values={{ value: intl.formatNumber(props.file.frame_rate ?? 0) }}
          />
        </TextField>
        <TextField id="bitrate">
          <FormattedMessage
            id="megabits_per_second"
            values={{
              value: intl.formatNumber((props.file.bit_rate ?? 0) / 1000000, {
                maximumFractionDigits: 2,
              }),
            }}
          />
        </TextField>
        <TextField
          id="media_info.video_codec"
          value={props.file.video_codec ?? ""}
          truncate
        />
        <TextField
          id="media_info.audio_codec"
          value={props.file.audio_codec ?? ""}
          truncate
        />
      </dl>
      {props.ofMany && props.onSetPrimaryFile && !props.primary && (
        <div>
          <Button
            className="edit-button"
            disabled={props.loading}
            onClick={props.onSetPrimaryFile}
          >
            <FormattedMessage id="actions.make_primary" />
          </Button>
          <Button
            className="edit-button"
            disabled={props.loading}
            onClick={props.onReassign}
          >
            <FormattedMessage id="actions.reassign" />
          </Button>
          <Button className="edit-button" onClick={onSplit}>
            <FormattedMessage id="actions.split" />
          </Button>
          <Button
            variant="danger"
            disabled={props.loading}
            onClick={props.onDeleteFile}
          >
            <FormattedMessage id="actions.delete_file" />
          </Button>
        </div>
      )}
    </div>
  );
};

interface ISceneFileInfoPanelProps {
  scene: GQL.SceneDataFragment;
}

const _SceneFileInfoPanel: React.FC<ISceneFileInfoPanelProps> = (
  props: ISceneFileInfoPanelProps
) => {
  const Toast = useToast();
  const coverImageSource =
    (props.scene as { cover_image_source?: string | null })
      .cover_image_source ?? null;

  const [loading, setLoading] = useState(false);
  const [deletingFile, setDeletingFile] = useState<GQL.VideoFileDataFragment>();
  const [reassigningFile, setReassigningFile] =
    useState<GQL.VideoFileDataFragment>();

  function renderStashIDs() {
    if (!props.scene.stash_ids.length) {
      return;
    }

    return (
      <>
        <dt>
          <FormattedMessage id="stash_ids" />
        </dt>
        <dd>
          <dl>
            {props.scene.stash_ids.map((stashID) => {
              return (
                <dd key={stashID.stash_id} className="row no-gutters">
                  <StashIDPill stashID={stashID} linkType="scenes" />
                </dd>
              );
            })}
          </dl>
        </dd>
      </>
    );
  }

  function renderFunscript() {
    if (props.scene.interactive) {
      return (
        <URLField
          name="Funscript"
          url={props.scene.paths.funscript}
          value={props.scene.paths.funscript}
          truncate
        />
      );
    }
  }

  function renderInteractiveSpeed() {
    if (props.scene.interactive_speed) {
      return (
        <TextField id="media_info.interactive_speed">
          <FormattedNumber value={props.scene.interactive_speed} />
        </TextField>
      );
    }
  }

  const filesPanel = useMemo(() => {
    if (props.scene.files.length === 0) {
      return;
    }

    if (props.scene.files.length === 1) {
      return (
        <FileInfoPanel
          sceneID={props.scene.id}
          file={props.scene.files[0]}
          coverImageSource={coverImageSource}
        />
      );
    }

    async function onSetPrimaryFile(fileID: string) {
      try {
        setLoading(true);
        await mutateSceneSetPrimaryFile(props.scene.id, fileID);
      } catch (e) {
        Toast.error(e);
      } finally {
        setLoading(false);
      }
    }

    return (
      <Accordion defaultActiveKey={props.scene.files[0].id}>
        {deletingFile && (
          <DeleteFilesDialog
            onClose={() => setDeletingFile(undefined)}
            selected={[deletingFile]}
          />
        )}
        {reassigningFile && (
          <ReassignFilesDialog
            onClose={() => setReassigningFile(undefined)}
            selected={reassigningFile}
          />
        )}
        {props.scene.files.map((file, index) => (
          <Card key={file.id} className="scene-file-card">
            <Accordion.Toggle as={Card.Header} eventKey={file.id}>
              <TruncatedText text={TextUtils.fileNameFromPath(file.path)} />
            </Accordion.Toggle>
            <Accordion.Collapse eventKey={file.id}>
              <Card.Body>
                <FileInfoPanel
                  sceneID={props.scene.id}
                  file={file}
                  coverImageSource={coverImageSource}
                  primary={index === 0}
                  ofMany
                  onSetPrimaryFile={() => onSetPrimaryFile(file.id)}
                  onDeleteFile={() => setDeletingFile(file)}
                  onReassign={() => setReassigningFile(file)}
                  loading={loading}
                />
              </Card.Body>
            </Accordion.Collapse>
          </Card>
        ))}
      </Accordion>
    );
  }, [
    props.scene,
    loading,
    Toast,
    deletingFile,
    reassigningFile,
    coverImageSource,
  ]);

  return (
    <>
      <dl className="container scene-file-info details-list">
        {props.scene.files.length > 0 && (
          <URLField
            id="media_info.stream"
            url={props.scene.paths.stream}
            value={props.scene.paths.stream}
            truncate
          />
        )}
        {renderFunscript()}
        {renderInteractiveSpeed()}
        <URLsField id="urls" urls={props.scene.urls} truncate />
        {renderStashIDs()}
      </dl>

      {filesPanel}
    </>
  );
};

export const SceneFileInfoPanel = PatchComponent(
  "SceneFileInfoPanel",
  _SceneFileInfoPanel
);
export default SceneFileInfoPanel;
