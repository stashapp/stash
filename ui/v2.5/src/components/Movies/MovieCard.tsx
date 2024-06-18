import React, { useEffect, useState } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { GridCard, calculateCardWidth } from "../Shared/GridCard/GridCard";
import { HoverPopover } from "../Shared/HoverPopover";
import { Icon } from "../Shared/Icon";
import { SceneLink, TagLink } from "../Shared/TagLink";
import { TruncatedText } from "../Shared/TruncatedText";
import { FormattedMessage } from "react-intl";
import { RatingBanner } from "../Shared/RatingBanner";
import { faPlayCircle, faTag } from "@fortawesome/free-solid-svg-icons";
import ScreenUtils from "src/utils/screen";

interface IProps {
  movie: GQL.MovieDataFragment;
  containerWidth?: number;
  sceneIndex?: number;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const MovieCard: React.FC<IProps> = ({
  movie,
  sceneIndex,
  containerWidth,
  selecting,
  selected,
  onSelectedChanged,
}) => {
  const [cardWidth, setCardWidth] = useState<number>();

  useEffect(() => {
    if (!containerWidth || ScreenUtils.isMobile()) return;

    let preferredCardWidth = 250;
    let fittedCardWidth = calculateCardWidth(
      containerWidth,
      preferredCardWidth!
    );
    setCardWidth(fittedCardWidth);
  }, [containerWidth]);

  function maybeRenderSceneNumber() {
    if (!sceneIndex) return;

    return (
      <>
        <hr />
        <span className="movie-scene-number">
          <FormattedMessage id="scene" /> #{sceneIndex}
        </span>
      </>
    );
  }

  function maybeRenderScenesPopoverButton() {
    if (movie.scenes.length === 0) return;

    const popoverContent = movie.scenes.map((scene) => (
      <SceneLink key={scene.id} scene={scene} />
    ));

    return (
      <HoverPopover
        className="scene-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon={faPlayCircle} />
          <span>{movie.scenes.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderTagPopoverButton() {
    if (movie.tags.length <= 0) return;

    const popoverContent = movie.tags.map((tag) => (
      <TagLink key={tag.id} linkType="movie" tag={tag} />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal tag-count">
          <Icon icon={faTag} />
          <span>{movie.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (sceneIndex || movie.scenes.length > 0 || movie.tags.length > 0) {
      return (
        <>
          {maybeRenderSceneNumber()}
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderTagPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <GridCard
      className="movie-card"
      url={`/movies/${movie.id}`}
      width={cardWidth}
      title={movie.name}
      linkClassName="movie-card-header"
      image={
        <>
          <img
            loading="lazy"
            className="movie-card-image"
            alt={movie.name ?? ""}
            src={movie.front_image_path ?? ""}
          />
          <RatingBanner rating={movie.rating100} />
        </>
      }
      details={
        <div className="movie-card__details">
          <span className="movie-card__date">{movie.date}</span>
          <TruncatedText
            className="movie-card__description"
            text={movie.synopsis}
            lineCount={3}
          />
        </div>
      }
      selected={selected}
      selecting={selecting}
      onSelectedChanged={onSelectedChanged}
      popovers={maybeRenderPopoverButtonGroup()}
    />
  );
};
