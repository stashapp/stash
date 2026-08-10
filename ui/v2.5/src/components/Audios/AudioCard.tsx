import React, { useMemo } from "react";
import { Button, ButtonGroup, OverlayTrigger, Tooltip } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import { GalleryLink, TagLink } from "../Shared/TagLink";
import { HoverPopover } from "../Shared/HoverPopover";
import { TruncatedText } from "../Shared/TruncatedText";
import TextUtils from "src/utils/text";
import { PerformerPopoverButton } from "../Shared/PerformerPopoverButton";
import { GridCard } from "../Shared/GridCard/GridCard";
import { RatingBanner } from "../Shared/RatingBanner";
import { FormattedMessage } from "react-intl";
import {
  faBox,
  faFilm,
  faImages,
  faMicrophone,
  faTag,
} from "@fortawesome/free-solid-svg-icons";
import { objectPath, objectTitle } from "src/core/files";
import { PatchComponent } from "src/patch";
import { StudioOverlay } from "../Shared/GridCard/StudioOverlay";
import { GroupTag } from "../Groups/GroupTag";
import { FileSize } from "../Shared/FileSize";
import { OCounterButton } from "../Shared/CountButton";
import { AudioQueue } from "src/models/audioQueue";
import { useConfigurationContext } from "src/hooks/Config";

interface IAudioCardProps {
  audio: GQL.SlimAudioDataFragment;
  width?: number;
  index?: number;
  queue?: AudioQueue;
  compact?: boolean;
  selecting?: boolean;
  selected?: boolean | undefined;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  fromGroupId?: string;
}

const Description: React.FC<{
  audioNumber?: number;
}> = ({ audioNumber }) => {
  if (!audioNumber) return null;

  return (
    <>
      <hr />
      <span className="audio-group-audio-number">
        <FormattedMessage id="audio" /> #{audioNumber}
      </span>
    </>
  );
};

const AudioCardPopovers = React.memo(
  PatchComponent("AudioCard.Popovers", (props: IAudioCardProps) => {
    const audioNumber = useMemo(() => {
      if (!props.fromGroupId) {
        return undefined;
      }

      const group = props.audio.groups.find(
        (g) => g.group.id === props.fromGroupId
      );
      return group?.audio_index ?? undefined;
    }, [props.fromGroupId, props.audio.groups]);

    function maybeRenderTagPopoverButton() {
      if (props.audio.tags.length <= 0) return;

      const popoverContent = props.audio.tags.map((tag) => (
        <TagLink key={tag.id} tag={tag} />
      ));

      return (
        <HoverPopover
          className="tag-count"
          placement="bottom"
          content={popoverContent}
        >
          <Button className="minimal">
            <Icon icon={faTag} />
            <span>{props.audio.tags.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderPerformerPopoverButton() {
      if (props.audio.performers.length <= 0) return;

      return (
        <PerformerPopoverButton
          performers={props.audio.performers}
          linkType="audio"
        />
      );
    }

    function maybeRenderGroupPopoverButton() {
      if (props.audio.groups.length <= 0) return;

      const popoverContent = props.audio.groups.map((audioGroup) => (
        <GroupTag key={audioGroup.group.id} group={audioGroup.group} />
      ));

      return (
        <HoverPopover
          placement="bottom"
          content={popoverContent}
          className="group-count tag-tooltip"
        >
          <Button className="minimal">
            <Icon icon={faFilm} />
            <span>{props.audio.groups.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderOCounter() {
      if (props.audio.o_counter) {
        return <OCounterButton value={props.audio.o_counter} />;
      }
    }

    function maybeRenderGallery() {
      if (props.audio.galleries.length <= 0) return;

      const popoverContent = props.audio.galleries.map((gallery) => (
        <GalleryLink key={gallery.id} gallery={gallery} />
      ));

      return (
        <HoverPopover
          className="gallery-count"
          placement="bottom"
          content={popoverContent}
        >
          <Button className="minimal">
            <Icon icon={faImages} />
            <span>{props.audio.galleries.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderOrganized() {
      if (props.audio.organized) {
        return (
          <OverlayTrigger
            overlay={<Tooltip id="organised-tooltip">{"Organized"}</Tooltip>}
            placement="bottom"
          >
            <div className="organized">
              <Button className="minimal">
                <Icon icon={faBox} />
              </Button>
            </div>
          </OverlayTrigger>
        );
      }
    }

    function maybeRenderPopoverButtonGroup() {
      if (
        !props.compact &&
        (props.audio.tags.length > 0 ||
          props.audio.performers.length > 0 ||
          props.audio.groups.length > 0 ||
          props.audio?.o_counter ||
          props.audio.galleries.length > 0 ||
          props.audio.organized ||
          audioNumber !== undefined)
      ) {
        return (
          <>
            <Description audioNumber={audioNumber ?? undefined} />
            <hr />
            <ButtonGroup className="card-popovers">
              {maybeRenderTagPopoverButton()}
              {maybeRenderPerformerPopoverButton()}
              {maybeRenderGroupPopoverButton()}
              {maybeRenderOCounter()}
              {maybeRenderGallery()}
              {maybeRenderOrganized()}
            </ButtonGroup>
          </>
        );
      }
    }

    return <>{maybeRenderPopoverButtonGroup()}</>;
  })
);

const AudioCardDetails = React.memo(
  PatchComponent("AudioCard.Details", (props: IAudioCardProps) => {
    return (
      <div className="audio-card__details">
        <span className="audio-card__date">{props.audio.date}</span>
        <span className="file-path extra-audio-info">
          {objectPath(props.audio)}
        </span>
        <TruncatedText
          className="audio-card__description"
          text={props.audio.details}
          lineCount={3}
        />
      </div>
    );
  })
);

const AudioCardOverlays = React.memo(
  PatchComponent("AudioCard.Overlays", (props: IAudioCardProps) => {
    const ret = useMemo(() => {
      return (
        <StudioOverlay studio={props.audio.studio} disabled={props.selecting} />
      );
    }, [props.audio.studio, props.selecting]);

    return ret;
  })
);

interface IAudioSpecsOverlay {
  audio: GQL.SlimAudioDataFragment;
}

export const AudioSpecsOverlay: React.FC<IAudioSpecsOverlay> = React.memo(
  PatchComponent("AudioCard.AudioSpecs", ({ audio }) => {
    const file = audio.files?.[0];
    if (!file) return null;
    return (
      <div className="audio-specs-overlay">
        <span className="overlay-filesize extra-audio-info">
          <FileSize size={file.size} />
        </span>
        {file.duration > 0 ? (
          <span className="overlay-duration">
            {TextUtils.secondsToTimestamp(file.duration)}
          </span>
        ) : (
          ""
        )}
      </div>
    );
  })
);

// AudioCoverImage renders the manually uploaded cover, falling back to a
// microphone glyph to match the server-side default cover image. Audio has no
// generated previews or sprites to scrub through.
const AudioCoverImage = React.memo(
  PatchComponent("AudioCard.Image", (props: IAudioCardProps) => {
    const cover = props.audio.paths.screenshot;

    return (
      <>
        <div className="audio-card-preview">
          {cover ? (
            <img
              className="audio-card-preview-image"
              loading="lazy"
              src={cover}
              alt=""
            />
          ) : (
            <div className="audio-card-preview-placeholder">
              <Icon icon={faMicrophone} />
            </div>
          )}
        </div>
        <RatingBanner rating={props.audio.rating100} />
        <AudioSpecsOverlay audio={props.audio} />
      </>
    );
  })
);

export const AudioCard = React.memo(
  PatchComponent("AudioCard", (props: IAudioCardProps) => {
    const { configuration } = useConfigurationContext();

    const file = useMemo(
      () => (props.audio.files.length > 0 ? props.audio.files[0] : undefined),
      [props.audio]
    );

    function zoomIndex() {
      if (!props.compact && props.zoomIndex !== undefined) {
        return `zoom-${props.zoomIndex}`;
      }

      return "";
    }

    function filelessClass() {
      if (!props.audio.files.length) {
        return "fileless";
      }

      return "";
    }

    const cont = configuration?.interface.continuePlaylistDefault ?? false;

    const audioLink = props.queue
      ? props.queue.makeLink(props.audio.id, {
          audioIndex: props.index,
          continue: cont,
        })
      : `/audios/${props.audio.id}`;

    return (
      <GridCard
        className={`audio-card ${zoomIndex()} ${filelessClass()}`}
        url={audioLink}
        title={objectTitle(props.audio)}
        width={props.width}
        linkClassName="audio-card-link"
        thumbnailSectionClassName="audio-section"
        resumeTime={props.audio.resume_time ?? undefined}
        duration={file?.duration ?? undefined}
        image={<AudioCoverImage {...props} />}
        overlays={<AudioCardOverlays {...props} />}
        details={<AudioCardDetails {...props} />}
        popovers={<AudioCardPopovers {...props} />}
        selected={props.selected}
        selecting={props.selecting}
        onSelectedChanged={props.onSelectedChanged}
      />
    );
  })
);
