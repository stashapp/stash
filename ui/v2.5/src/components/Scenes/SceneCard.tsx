import React, { useEffect, useMemo, useRef, useState } from "react";
import { Button, ButtonGroup, OverlayTrigger, Tooltip } from "react-bootstrap";
import { Link, useHistory } from "react-router-dom";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import {
  GalleryLink,
  TagLink,
  GroupLink,
  SceneMarkerLink,
} from "../Shared/TagLink";
import { HoverPopover } from "../Shared/HoverPopover";
import { SweatDrops } from "../Shared/SweatDrops";
import { TruncatedText } from "../Shared/TruncatedText";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { SceneQueue } from "src/models/sceneQueue";
import { ConfigurationContext } from "src/hooks/Config";
import { PerformerPopoverButton } from "../Shared/PerformerPopoverButton";
import { GridCard, calculateCardWidth } from "../Shared/GridCard/GridCard";
import { RatingBanner } from "../Shared/RatingBanner";
import { FormattedNumber } from "react-intl";
import {
  faBox,
  faCopy,
  faFilm,
  faImages,
  faMapMarkerAlt,
  faTag,
} from "@fortawesome/free-solid-svg-icons";
import { objectPath, objectTitle } from "src/core/files";
import { PreviewScrubber } from "./PreviewScrubber";
import { PatchComponent } from "src/patch";
import ScreenUtils from "src/utils/screen";
import { StudioOverlay } from "../Shared/GridCard/StudioOverlay";

interface IScenePreviewProps {
  isPortrait: boolean;
  image?: string;
  video?: string;
  soundActive: boolean;
  vttPath?: string;
  onScrubberClick?: (timestamp: number) => void;
}

export const ScenePreview: React.FC<IScenePreviewProps> = ({
  image,
  video,
  isPortrait,
  soundActive,
  vttPath,
  onScrubberClick,
}) => {
  const videoEl = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.intersectionRatio > 0)
          // Catch is necessary due to DOMException if user hovers before clicking on page
          videoEl.current?.play()?.catch(() => {});
        else videoEl.current?.pause();
      });
    });

    if (videoEl.current) observer.observe(videoEl.current);
  });

  useEffect(() => {
    if (videoEl?.current?.volume)
      videoEl.current.volume = soundActive ? 0.05 : 0;
  }, [soundActive]);

  return (
    <div className={cx("scene-card-preview", { portrait: isPortrait })}>
      <img
        className="scene-card-preview-image"
        loading="lazy"
        src={image}
        alt=""
      />
      <video
        disableRemotePlayback
        playsInline
        muted={!soundActive}
        className="scene-card-preview-video"
        loop
        preload="none"
        ref={videoEl}
        src={video}
      />
      <PreviewScrubber vttPath={vttPath} onClick={onScrubberClick} />
    </div>
  );
};

interface ISceneCardProps {
  scene: GQL.SlimSceneDataFragment;
  containerWidth?: number;
  previewHeight?: number;
  index?: number;
  queue?: SceneQueue;
  compact?: boolean;
  selecting?: boolean;
  selected?: boolean | undefined;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

const SceneCardPopovers = PatchComponent(
  "SceneCard.Popovers",
  (props: ISceneCardProps) => {
    const file = useMemo(
      () => (props.scene.files.length > 0 ? props.scene.files[0] : undefined),
      [props.scene]
    );

    function maybeRenderTagPopoverButton() {
      if (props.scene.tags.length <= 0) return;

      const popoverContent = props.scene.tags.map((tag) => (
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
            <span>{props.scene.tags.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderPerformerPopoverButton() {
      if (props.scene.performers.length <= 0) return;

      return <PerformerPopoverButton performers={props.scene.performers} />;
    }

    function maybeRenderGroupPopoverButton() {
      if (props.scene.movies.length <= 0) return;

      const popoverContent = props.scene.movies.map((sceneGroup) => (
        <div className="group-tag-container row" key={sceneGroup.movie.id}>
          <Link
            to={`/groups/${sceneGroup.movie.id}`}
            className="group-tag col m-auto zoom-2"
          >
            <img
              className="image-thumbnail"
              alt={sceneGroup.movie.name ?? ""}
              src={sceneGroup.movie.front_image_path ?? ""}
            />
          </Link>
          <GroupLink
            key={sceneGroup.movie.id}
            group={sceneGroup.movie}
            className="d-block"
          />
        </div>
      ));

      return (
        <HoverPopover
          placement="bottom"
          content={popoverContent}
          className="group-count tag-tooltip"
        >
          <Button className="minimal">
            <Icon icon={faFilm} />
            <span>{props.scene.movies.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderSceneMarkerPopoverButton() {
      if (props.scene.scene_markers.length <= 0) return;

      const popoverContent = props.scene.scene_markers.map((marker) => {
        const markerWithScene = { ...marker, scene: { id: props.scene.id } };
        return <SceneMarkerLink key={marker.id} marker={markerWithScene} />;
      });

      return (
        <HoverPopover
          className="marker-count"
          placement="bottom"
          content={popoverContent}
        >
          <Button className="minimal">
            <Icon icon={faMapMarkerAlt} />
            <span>{props.scene.scene_markers.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderOCounter() {
      if (props.scene.o_counter) {
        return (
          <div className="o-count">
            <Button className="minimal">
              <span className="fa-icon">
                <SweatDrops />
              </span>
              <span>{props.scene.o_counter}</span>
            </Button>
          </div>
        );
      }
    }

    function maybeRenderGallery() {
      if (props.scene.galleries.length <= 0) return;

      const popoverContent = props.scene.galleries.map((gallery) => (
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
            <span>{props.scene.galleries.length}</span>
          </Button>
        </HoverPopover>
      );
    }

    function maybeRenderOrganized() {
      if (props.scene.organized) {
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

    function maybeRenderDupeCopies() {
      const phash = file
        ? file.fingerprints.find((fp) => fp.type === "phash")
        : undefined;

      if (phash) {
        return (
          <div className="other-copies extra-scene-info">
            <Button
              href={NavUtils.makeScenesPHashMatchUrl(phash.value)}
              className="minimal"
            >
              <Icon icon={faCopy} />
            </Button>
          </div>
        );
      }
    }

    function maybeRenderPopoverButtonGroup() {
      if (
        !props.compact &&
        (props.scene.tags.length > 0 ||
          props.scene.performers.length > 0 ||
          props.scene.movies.length > 0 ||
          props.scene.scene_markers.length > 0 ||
          props.scene?.o_counter ||
          props.scene.galleries.length > 0 ||
          props.scene.organized)
      ) {
        return (
          <>
            <hr />
            <ButtonGroup className="card-popovers">
              {maybeRenderTagPopoverButton()}
              {maybeRenderPerformerPopoverButton()}
              {maybeRenderGroupPopoverButton()}
              {maybeRenderSceneMarkerPopoverButton()}
              {maybeRenderOCounter()}
              {maybeRenderGallery()}
              {maybeRenderOrganized()}
              {maybeRenderDupeCopies()}
            </ButtonGroup>
          </>
        );
      }
    }

    return <>{maybeRenderPopoverButtonGroup()}</>;
  }
);

const SceneCardDetails = PatchComponent(
  "SceneCard.Details",
  (props: ISceneCardProps) => {
    return (
      <div className="scene-card__details">
        <span className="scene-card__date">{props.scene.date}</span>
        <span className="file-path extra-scene-info">
          {objectPath(props.scene)}
        </span>
        <TruncatedText
          className="scene-card__description"
          text={props.scene.details}
          lineCount={3}
        />
      </div>
    );
  }
);

const SceneCardOverlays = PatchComponent(
  "SceneCard.Overlays",
  (props: ISceneCardProps) => {
    return <StudioOverlay studio={props.scene.studio} />;
  }
);

const SceneCardImage = PatchComponent(
  "SceneCard.Image",
  (props: ISceneCardProps) => {
    const history = useHistory();
    const { configuration } = React.useContext(ConfigurationContext);
    const cont = configuration?.interface.continuePlaylistDefault ?? false;

    const file = useMemo(
      () => (props.scene.files.length > 0 ? props.scene.files[0] : undefined),
      [props.scene]
    );

    function maybeRenderSceneSpecsOverlay() {
      let sizeObj = null;
      if (file?.size) {
        sizeObj = TextUtils.fileSize(file.size);
      }
      return (
        <div className="scene-specs-overlay">
          {sizeObj != null ? (
            <span className="overlay-filesize extra-scene-info">
              <FormattedNumber
                value={sizeObj.size}
                maximumFractionDigits={TextUtils.fileSizeFractionalDigits(
                  sizeObj.unit
                )}
              />
              {TextUtils.formatFileSizeUnit(sizeObj.unit)}
            </span>
          ) : (
            ""
          )}
          {file?.width && file?.height ? (
            <span className="overlay-resolution">
              {" "}
              {TextUtils.resolution(file?.width, file?.height)}
            </span>
          ) : (
            ""
          )}
          {(file?.duration ?? 0) >= 1
            ? TextUtils.secondsToTimestamp(file?.duration ?? 0)
            : ""}
        </div>
      );
    }

    function maybeRenderInteractiveSpeedOverlay() {
      return (
        <div className="scene-interactive-speed-overlay">
          {props.scene.interactive_speed ?? ""}
        </div>
      );
    }

    function onScrubberClick(timestamp: number) {
      const link = props.queue
        ? props.queue.makeLink(props.scene.id, {
            sceneIndex: props.index,
            continue: cont,
            start: timestamp,
          })
        : `/scenes/${props.scene.id}?t=${timestamp}`;

      history.push(link);
    }

    function isPortrait() {
      const width = file?.width ? file.width : 0;
      const height = file?.height ? file.height : 0;
      return height > width;
    }

    return (
      <>
        <ScenePreview
          image={props.scene.paths.screenshot ?? undefined}
          video={props.scene.paths.preview ?? undefined}
          isPortrait={isPortrait()}
          soundActive={configuration?.interface?.soundOnPreview ?? false}
          vttPath={props.scene.paths.vtt ?? undefined}
          onScrubberClick={onScrubberClick}
        />
        <RatingBanner rating={props.scene.rating100} />
        {maybeRenderSceneSpecsOverlay()}
        {maybeRenderInteractiveSpeedOverlay()}
      </>
    );
  }
);

export const SceneCard = PatchComponent(
  "SceneCard",
  (props: ISceneCardProps) => {
    const { configuration } = React.useContext(ConfigurationContext);
    const [cardWidth, setCardWidth] = useState<number>();

    const file = useMemo(
      () => (props.scene.files.length > 0 ? props.scene.files[0] : undefined),
      [props.scene]
    );

    function zoomIndex() {
      if (!props.compact && props.zoomIndex !== undefined) {
        return `zoom-${props.zoomIndex}`;
      }

      return "";
    }

    function filelessClass() {
      if (!props.scene.files.length) {
        return "fileless";
      }

      return "";
    }

    useEffect(() => {
      if (
        !props.containerWidth ||
        props.zoomIndex === undefined ||
        ScreenUtils.isMobile()
      )
        return;

      let zoomValue = props.zoomIndex;
      let preferredCardWidth: number;
      switch (zoomValue) {
        case 0:
          preferredCardWidth = 240;
          break;
        case 1:
          preferredCardWidth = 340; // this value is intentionally higher than 320
          break;
        case 2:
          preferredCardWidth = 480;
          break;
        case 3:
          preferredCardWidth = 640;
      }
      let fittedCardWidth = calculateCardWidth(
        props.containerWidth,
        preferredCardWidth!
      );
      setCardWidth(fittedCardWidth);
    }, [props, props.containerWidth, props.zoomIndex]);

    const cont = configuration?.interface.continuePlaylistDefault ?? false;

    const sceneLink = props.queue
      ? props.queue.makeLink(props.scene.id, {
          sceneIndex: props.index,
          continue: cont,
        })
      : `/scenes/${props.scene.id}`;

    return (
      <GridCard
        className={`scene-card ${zoomIndex()} ${filelessClass()}`}
        url={sceneLink}
        title={objectTitle(props.scene)}
        width={cardWidth}
        linkClassName="scene-card-link"
        thumbnailSectionClassName="video-section"
        resumeTime={props.scene.resume_time ?? undefined}
        duration={file?.duration ?? undefined}
        interactiveHeatmap={
          props.scene.interactive_speed
            ? props.scene.paths.interactive_heatmap ?? undefined
            : undefined
        }
        image={<SceneCardImage {...props} />}
        overlays={<SceneCardOverlays {...props} />}
        details={<SceneCardDetails {...props} />}
        popovers={<SceneCardPopovers {...props} />}
        selected={props.selected}
        selecting={props.selecting}
        onSelectedChanged={props.onSelectedChanged}
      />
    );
  }
);
