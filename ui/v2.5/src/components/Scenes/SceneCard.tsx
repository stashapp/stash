import React, { useEffect, useRef } from "react";
import { Button, ButtonGroup, Card, Form } from "react-bootstrap";
import { Link } from "react-router-dom";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import { useConfiguration } from "src/core/StashService";
import { Icon, TagLink, HoverPopover, SweatDrops } from "src/components/Shared";
import { TextUtils } from "src/utils";

interface IScenePreviewProps {
  isPortrait: boolean;
  image?: string;
  video?: string;
  soundActive: boolean;
}

const ScenePreview: React.FC<IScenePreviewProps> = ({
  image,
  video,
  isPortrait,
  soundActive,
}) => {
  const videoEl = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.intersectionRatio > 0)
            // Catch is necessary due to DOMException if user hovers before clicking on page
            videoEl.current?.play().catch(() => {});
          else videoEl.current?.pause();
        });
      },
      { root: document.documentElement }
    );

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
  selecting?: boolean;
  selected: boolean | undefined;
  zoomIndex: number;
  onSelectedChanged: (selected: boolean, shiftKey: boolean) => void;
}

export const SceneCard: React.FC<ISceneCardProps> = (
  props: ISceneCardProps
) => {
  const config = useConfiguration();
  const showStudioAsText =
    config?.data?.configuration.interface.showStudioAsText ?? false;

  function maybeRenderRatingBanner() {
    if (!props.scene.rating) {
      return;
    }
    return (
      <div
        className={`rating-banner ${
          props.scene.rating ? `rating-${props.scene.rating}` : ""
        }`}
      >
        RATING: {props.scene.rating}
      </div>
    );
  }

  function maybeRenderSceneSpecsOverlay() {
    return (
      <div className="scene-specs-overlay">
        {props.scene.file.height ? (
          <span className="overlay-resolution">
            {" "}
            {TextUtils.resolution(props.scene.file.height)}
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
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon="tag" />
          <span>{props.scene.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPerformerPopoverButton() {
    if (props.scene.performers.length <= 0) return;

    const popoverContent = props.scene.performers.map((performer) => (
      <div className="performer-tag-container row" key={performer.id}>
        <Link
          to={`/performers/${performer.id}`}
          className="performer-tag col m-auto zoom-2"
        >
          <img
            className="image-thumbnail"
            alt={performer.name ?? ""}
            src={performer.image_path ?? ""}
          />
        </Link>
        <TagLink key={performer.id} performer={performer} className="d-block" />
      </div>
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal">
          <Icon icon="user" />
          <span>{props.scene.performers.length}</span>
        </Button>
      </HoverPopover>
    );
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
        className="tag-tooltip"
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
      <HoverPopover placement="bottom" content={popoverContent}>
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
        <div>
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
    if (props.scene.gallery) {
      return (
        <div>
          <Link to={`/galleries/${props.scene.gallery.id}`}>
            <Button className="minimal">
              <Icon icon="image" />
            </Button>
          </Link>
        </div>
      );
    }
  }

  function maybeRenderPopoverButtonGroup() {
    if (
      props.scene.tags.length > 0 ||
      props.scene.performers.length > 0 ||
      props.scene.movies.length > 0 ||
      props.scene.scene_markers.length > 0 ||
      props.scene?.o_counter ||
      props.scene.gallery
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
          </ButtonGroup>
        </>
      );
    }
  }

  function handleSceneClick(
    event: React.MouseEvent<HTMLAnchorElement, MouseEvent>
  ) {
    const { shiftKey } = event;

    if (props.selecting) {
      props.onSelectedChanged(!props.selected, shiftKey);
      event.preventDefault();
    }
  }

  function handleDrag(event: React.DragEvent<HTMLAnchorElement>) {
    if (props.selecting) {
      event.dataTransfer.setData("text/plain", "");
      event.dataTransfer.setDragImage(new Image(), 0, 0);
    }
  }

  function handleDragOver(event: React.DragEvent<HTMLAnchorElement>) {
    const ev = event;
    const shiftKey = false;

    if (props.selecting && !props.selected) {
      props.onSelectedChanged(true, shiftKey);
    }

    ev.dataTransfer.dropEffect = "move";
    ev.preventDefault();
  }

  function isPortrait() {
    const { file } = props.scene;
    const width = file.width ? file.width : 0;
    const height = file.height ? file.height : 0;
    return height > width;
  }

  let shiftKey = false;

  return (
    <Card className={`scene-card zoom-${props.zoomIndex}`}>
      <Form.Control
        type="checkbox"
        className="scene-card-check"
        checked={props.selected}
        onChange={() => props.onSelectedChanged(!props.selected, shiftKey)}
        onClick={(event: React.MouseEvent<HTMLInputElement, MouseEvent>) => {
          // eslint-disable-next-line prefer-destructuring
          shiftKey = event.shiftKey;
          event.stopPropagation();
        }}
      />

      <div className="video-section">
        <Link
          to={`/scenes/${props.scene.id}`}
          className="scene-card-link"
          onClick={handleSceneClick}
          onDragStart={handleDrag}
          onDragOver={handleDragOver}
          draggable={props.selecting}
        >
          <ScenePreview
            image={props.scene.paths.screenshot ?? undefined}
            video={props.scene.paths.preview ?? undefined}
            isPortrait={isPortrait()}
            soundActive={
              config.data?.configuration?.interface?.soundOnPreview ?? false
            }
          />
          {maybeRenderRatingBanner()}
          {maybeRenderSceneSpecsOverlay()}
        </Link>
        {maybeRenderSceneStudioOverlay()}
      </div>
      <div className="card-section">
        <h5 className="card-section-title">
          {props.scene.title
            ? props.scene.title
            : TextUtils.fileNameFromPath(props.scene.path)}
        </h5>
        <span>{props.scene.date}</span>
        <p>
          {props.scene.details &&
            TextUtils.truncate(props.scene.details, 100, "... (continued)")}
        </p>
      </div>

      {maybeRenderPopoverButtonGroup()}
    </Card>
  );
};
