import React, { useEffect, useRef } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import { Link } from "react-router-dom";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import {
  Icon,
  TagLink,
  HoverPopover,
  SweatDrops,
  TruncatedText,
} from "src/components/Shared";
import { NavUtils, TextUtils } from "src/utils";
import { SceneQueue } from "src/models/sceneQueue";
import { ConfigurationContext } from "src/hooks/Config";
import { PerformerPopoverButton } from "../Shared/PerformerPopoverButton";
import { GridCard } from "../Shared/GridCard";
import { RatingBanner } from "../Shared/RatingBanner";
import { FormattedNumber } from "react-intl";

interface IScenePreviewProps {
  isPortrait: boolean;
  image?: string;
  video?: string;
  soundActive: boolean;
}

export const ScenePreview: React.FC<IScenePreviewProps> = ({
  image,
  video,
  isPortrait,
  soundActive,
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
  const { configuration } = React.useContext(ConfigurationContext);

  // studio image is missing if it uses the default
  const missingStudioImage = props.scene.studio?.image_path?.endsWith(
    "?default=true"
  );
  const showStudioAsText =
    missingStudioImage || (configuration?.interface.showStudioAsText ?? false);

  function maybeRenderSceneSpecsOverlay() {
    let sizeObj = null;
    if (props.scene.file.size) {
      sizeObj = TextUtils.fileSize(parseInt(props.scene.file.size));
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
        {props.scene.file.width && props.scene.file.height ? (
          <span className="overlay-resolution">
            {" "}
            {TextUtils.resolution(
              props.scene.file.width,
              props.scene.file.height
            )}
          </span>
        ) : (
          ""
        )}
        {(props.scene.file.duration ?? 0) >= 1
          ? TextUtils.secondsToTimestamp(props.scene.file.duration ?? 0)
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

  function maybeRenderSceneStudioOverlay() {
    if (!props.scene.studio) return;

    return (
      <div className="scene-studio-overlay">
        <Link to={`/studios/${props.scene.studio.id}`}>
          {showStudioAsText ? (
            props.scene.studio.name
          ) : (
            <img
              className="image-thumbnail"
              alt={props.scene.studio.name}
              src={props.scene.studio.image_path ?? ""}
            />
          )}
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
          <Icon icon="tag" />
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
        <TagLink
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
          <Icon icon="film" />
          <span>{props.scene.movies.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderSceneMarkerPopoverButton() {
    if (props.scene.scene_markers.length <= 0) return;

    const popoverContent = props.scene.scene_markers.map((marker) => {
      const markerPopover = { ...marker, scene: { id: props.scene.id } };
      return <TagLink key={marker.id} marker={markerPopover} />;
    });

    return (
      <HoverPopover
        className="marker-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon="map-marker-alt" />
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
      <TagLink key={gallery.id} gallery={gallery} />
    ));

    return (
      <HoverPopover
        className="gallery-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon="images" />
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
            <Icon icon="box" />
          </Button>
        </div>
      );
    }
  }

  function maybeRenderDupeCopies() {
    if (props.scene.phash) {
      return (
        <div className="other-copies extra-scene-info">
          <Button
            href={NavUtils.makeScenesPHashMatchUrl(props.scene.phash)}
            className="minimal"
          >
            <Icon icon="copy" />
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
            {maybeRenderDupeCopies()}
          </ButtonGroup>
        </>
      );
    }
  }

  function isPortrait() {
    const { file } = props.scene;
    const width = file.width ? file.width : 0;
    const height = file.height ? file.height : 0;
    return height > width;
  }

  function zoomIndex() {
    if (!props.compact && props.zoomIndex !== undefined) {
      return `zoom-${props.zoomIndex}`;
    }
  }

  const cont = configuration?.interface.continuePlaylistDefault ?? false;

  const sceneLink = props.queue
    ? props.queue.makeLink(props.scene.id, {
        sceneIndex: props.index,
        continue: cont,
      })
    : `/scenes/${props.scene.id}`;

  return (
    <GridCard
      className={`scene-card ${zoomIndex()}`}
      url={sceneLink}
      title={
        props.scene.title
          ? props.scene.title
          : TextUtils.fileNameFromPath(props.scene.path)
      }
      linkClassName="scene-card-link"
      thumbnailSectionClassName="video-section"
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
          />
          <RatingBanner rating={props.scene.rating} />
          {maybeRenderSceneSpecsOverlay()}
          {maybeRenderInteractiveSpeedOverlay()}
        </>
      }
      overlays={maybeRenderSceneStudioOverlay()}
      details={
        <div className="scene-card__details">
          <span className="scene-card__date">{props.scene.date}</span>
          <span className="file-path extra-scene-info">{props.scene.path}</span>
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
