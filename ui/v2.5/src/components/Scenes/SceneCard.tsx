import React, { useEffect, useMemo, useRef } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import { Link, useHistory } from "react-router-dom";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import { CodecLink } from "../Shared/CodecLink";
import {
  GalleryLink,
  TagLink,
  MovieLink,
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
import { GridCard } from "../Shared/GridCard";
import { RatingBanner } from "../Shared/RatingBanner";
import { FormattedNumber, FormattedMessage, useIntl } from "react-intl";
import {
  faBox,
  faInfo,
  faCopy,
  faClipboard,
  faFilm,
  faImages,
  faMapMarkerAlt,
  faTag,
} from "@fortawesome/free-solid-svg-icons";
import { objectPath, objectPathPreMapped, objectTitle } from "src/core/files";
import { PreviewScrubber } from "./PreviewScrubber";
import { IUIConfig } from "src/core/config";
import { Options, PositioningStrategy, Placement } from '@popperjs/core'; 

interface ISceneInfoCardProps {
  props?: GQL.SlimSceneDataFragment;
  files?: Array<{
    path: string;
    audio_codec: string;
    video_codec: string;
    size: number;
    width: number;
    height: number;
    duration: number;
    frame_rate: number;
    bit_rate: number;
    /* eslint-disable @typescript-eslint/no-explicit-any */
    fingerprints: any;
  }>;
  maybeRenderSceneSpecsOverlay: (prop: {
    size: number;
    width: number;
    height: number;
    duration: number;
    frame_rate: number;
    bit_rate: number;
  }, flag: boolean) => void;
  maybeRenderDupeCopies: (prop: {
    /* eslint-disable @typescript-eslint/no-explicit-any */
    fingerprints: any;
  }) => void;
}

export const MaybeRenderSceneInfoDetails: React.FC<ISceneInfoCardProps> = (
  props: ISceneInfoCardProps
) => {
  const { configuration: config } = React.useContext(ConfigurationContext);
  const enableSceneInfoDetails =
    (config?.ui as IUIConfig)?.enableSceneInfoDetails ?? true;

  if (!enableSceneInfoDetails) {
    return null;
  }

  const files = props?.files || [];
  if (files.length === 0) return null;

  const { maybeRenderSceneSpecsOverlay } = props;
  const { maybeRenderDupeCopies } = props;

  function CopyFilePathButton(filePath: string) {
    const copyFilePathToClipboard = () => {
      const pathSeparator = "/";
      const pathParts = filePath.split(pathSeparator);
      const directoryPath = pathParts.slice(0, -1).join(pathSeparator);
  
      try {
        navigator.clipboard.writeText(directoryPath)
      } catch (err) {
        console.error("Failed to copy to clipboard:", err);
      }
    };
  
    return (
      <div className="scene-info-overlay-file-path">
        <div>
          <Button onClick={copyFilePathToClipboard} className="minimal">
            <Icon icon={faClipboard} />
          </Button>
        </div>
        <div>{filePath}</div>
      </div>
    );
  }

  const popoverContent = files.map((prop, index) => (
    <>
      <div className="scene-info-popover" onClick={(e) => e.stopPropagation()}>
        {CopyFilePathButton(objectPathPreMapped(prop) || "")}
        <div className="scene-info-overlay-specs-root">
          {maybeRenderSceneSpecsOverlay(prop, false)}
          <span className={"scene-info-overlay-caps"}>
            <span className="scene-info-overlay-divider">|</span>
            <CodecLink
              key={prop.video_codec}
              codec={prop.video_codec}
              codecType={"video_codec"}
            />
          </span>
          <span className="scene-info-overlay-divider">:</span>
          <span className={"scene-info-overlay-caps"}>
            <CodecLink
              key={prop.audio_codec}
              codec={prop.audio_codec}
              codecType={"audio_codec"}
            />
          </span>
          {maybeRenderDupeCopies(prop)}
        </div>
        {index + 1 != files.length && <hr className="scene-info-overlay-hr" />}
      </div>
    </>
  ));

  interface IPopperType {
    width: number;
  }
  const popperConfig: Options = {
    strategy: 'fixed' as PositioningStrategy,
    placement: "bottom" as Placement,
    modifiers: [
      {
        name: "arrow",
        options: {
          x: 0,
          y: 0,
          centerOffset: 0,
        },
      },
      {
        name: "offset",
        options: {
          offset: ({
            placement,
            popper,
          }: {
            placement: Placement;
            popper: IPopperType;
          }) => {
            if (placement === "bottom") {
              return [popper.width / 4, 20];
            } else {
              return [];
            }
          },
        },
      },
    ],
  };

  return (
    <div
      className="scene-specs-overlay-left"
      onClick={(e) => e.stopPropagation()}
    >
      <HoverPopover
        placement="bottom"
        content={popoverContent}
        popperConfig={popperConfig}
      >
        <Icon icon={faInfo} />
      </HoverPopover>
    </div>
  );
};

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
      <img className="scene-card-preview-image" src={image} alt="" />
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
  index?: number;
  queue?: SceneQueue;
  compact?: boolean;
  selecting?: boolean;
  selected?: boolean | undefined;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const SceneCard: React.FC<ISceneCardProps> = (
  props: ISceneCardProps
) => {
  const history = useHistory();
  const { configuration } = React.useContext(ConfigurationContext);
  const intl = useIntl();

  const file = useMemo(
    () => (props.scene.files.length > 0 ? props.scene.files[0] : undefined),
    [props.scene]
  );

  function maybeRenderSceneSpecsOverlay(
    files?: {
      __typename?: "VideoFile" | undefined;
      size: number;
      width: number;
      height: number;
      duration: number;
      frame_rate: number;
      bit_rate: number;
    },
    renderDefault?: boolean
  ) {
    let sizeObj = null;

    const rootClass = renderDefault
      ? "scene-specs-overlay"
      : "scene-info-overlay-specs-node";
    const resolutionClass = renderDefault
      ? "overlay-resolution"
      : "scene-info-overlay-bold";
    const spanClass = renderDefault
      ? "overlay-filesize extra-scene-info"
      : "null";
    const boldClass = renderDefault
      ? "overlay-bold"
      : "scene-info-overlay-bold";

    if (files?.size) {
      sizeObj = TextUtils.fileSize(files.size);
    }

    return (
      <>
        <div className={rootClass}>
          {sizeObj != null ? (
            <span className={spanClass}>
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
          {!renderDefault && (
            <span className="scene-info-overlay-divider">|</span>
          )}
          {files?.width && files?.height ? (
            <span className={resolutionClass}>
              {" "}
              {TextUtils.resolution(files?.width, files?.height)}
            </span>
          ) : (
            ""
          )}
          {!renderDefault && (
            <span className="scene-info-overlay-divider">|</span>
          )}
          <span className="overlay-duration">
            {(files?.duration ?? 0) >= 1
              ? TextUtils.secondsToTimestamp(files?.duration ?? 0)
              : ""}
          </span>
          {!renderDefault && (
            <span className="scene-info-overlay-divider"> | </span>
          )}
          {!renderDefault && (
            <div>
              <span className={boldClass}>
                <FormattedMessage
                  id="frames_per_second_short"
                  values={{
                    value: intl?.formatNumber(files?.frame_rate ?? 0, {
                      maximumFractionDigits: 0,
                    }),
                  }}
                />
              </span>
              <span className="scene-info-overlay-divider">|</span>
              <span className={spanClass}>
                <FormattedMessage
                  id="megabits_per_second_short"
                  values={{
                    value: intl?.formatNumber(
                      (files?.bit_rate ?? 0) / 1000000,
                      {
                        maximumFractionDigits: 2,
                      }
                    ),
                  }}
                />
              </span>
            </div>
          )}
        </div>
      </>
    );
  }

  function maybeRenderInteractiveSpeedOverlay() {
    return (
      <div className="scene-interactive-speed-overlay">
        {props.scene.interactive_speed ?? ""}
      </div>
    );
  }

  function renderStudioThumbnail() {
    const studioImage = props.scene.studio?.image_path;
    const studioName = props.scene.studio?.name;

    if (configuration?.interface.showStudioAsText || !studioImage) {
      return studioName;
    }

    const studioImageURL = new URL(studioImage);
    if (studioImageURL.searchParams.get("default") === "true") {
      return studioName;
    }

    return (
      <img className="image-thumbnail" alt={studioName} src={studioImage} />
    );
  }

  function maybeRenderSceneStudioOverlay() {
    if (!props.scene.studio) return;

    return (
      <div className="scene-studio-overlay">
        <Link to={`/studios/${props.scene.studio.id}`}>
          {renderStudioThumbnail()}
        </Link>
      </div>
    );
  }

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

  function maybeRenderMoviePopoverButton() {
    if (props.scene.movies.length <= 0) return;

    const popoverContent = props.scene.movies.map((sceneMovie) => (
      <div className="movie-tag-container row" key="movie">
        <Link
          to={`/movies/${sceneMovie.movie.id}`}
          className="movie-tag col m-auto zoom-2"
        >
          <img
            className="image-thumbnail"
            alt={sceneMovie.movie.name ?? ""}
            src={sceneMovie.movie.front_image_path ?? ""}
          />
        </Link>
        <MovieLink
          key={sceneMovie.movie.id}
          movie={sceneMovie.movie}
          className="d-block"
        />
      </div>
    ));

    return (
      <HoverPopover
        placement="bottom"
        content={popoverContent}
        className="movie-count tag-tooltip"
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
        <div className="organized">
          <Button className="minimal">
            <Icon icon={faBox} />
          </Button>
        </div>
      );
    }
  }

  function maybeRenderDupeCopies(
    files: {
      __typename?: "VideoFile" | undefined;
      /* eslint-disable @typescript-eslint/no-explicit-any */
      fingerprints: any;
    },
    renderDefault?: boolean
  ) {
    const className = renderDefault
      ? "other-copies extra-scene-info"
      : "other-copies";

    const phash = files
      ? files.fingerprints.find((fp: { type: string }) => fp.type === "phash")
      : undefined;

    if (phash) {
      return (
        <div className={className}>
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
            {maybeRenderMoviePopoverButton()}
            {maybeRenderSceneMarkerPopoverButton()}
            {maybeRenderOCounter()}
            {maybeRenderGallery()}
            {maybeRenderOrganized()}
            {maybeRenderDupeCopies(props.scene.files[0], true)}
          </ButtonGroup>
        </>
      );
    }
  }

  function isPortrait() {
    const width = file?.width ? file.width : 0;
    const height = file?.height ? file.height : 0;
    return height > width;
  }

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

  const cont = configuration?.interface.continuePlaylistDefault ?? false;

  const sceneLink = props.queue
    ? props.queue.makeLink(props.scene.id, {
        sceneIndex: props.index,
        continue: cont,
      })
    : `/scenes/${props.scene.id}`;

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

  return (
    <GridCard
      className={`scene-card ${zoomIndex()} ${filelessClass()}`}
      url={sceneLink}
      title={objectTitle(props.scene)}
      linkClassName="scene-card-link"
      thumbnailSectionClassName="video-section"
      resumeTime={props.scene.resume_time ?? undefined}
      duration={file?.duration ?? undefined}
      interactiveHeatmap={
        props.scene.interactive_speed
          ? props.scene.paths.interactive_heatmap ?? undefined
          : undefined
      }
      image={
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
          <MaybeRenderSceneInfoDetails
            props={props.scene ?? undefined}
            files={props.scene.files}
            maybeRenderSceneSpecsOverlay={maybeRenderSceneSpecsOverlay}
            maybeRenderDupeCopies={maybeRenderDupeCopies}
          />
          {maybeRenderSceneSpecsOverlay(props.scene.files[0], true)}
          {maybeRenderInteractiveSpeedOverlay()}
        </>
      }
      overlays={maybeRenderSceneStudioOverlay()}
      details={
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
      }
      popovers={maybeRenderPopoverButtonGroup()}
      selected={props.selected}
      selecting={props.selecting}
      onSelectedChanged={props.onSelectedChanged}
    />
  );
};
