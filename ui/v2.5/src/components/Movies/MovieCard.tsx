import React, { useEffect, useState } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { GridCard, calculateCardWidth } from "../Shared/GridCard";
import { HoverPopover } from "../Shared/HoverPopover";
import { Icon } from "../Shared/Icon";
import { SceneLink } from "../Shared/TagLink";
import { TruncatedText } from "../Shared/TruncatedText";
import { FormattedMessage } from "react-intl";
import { RatingBanner } from "../Shared/RatingBanner";
import { faPlayCircle } from "@fortawesome/free-solid-svg-icons";
import ScreenUtils from "src/utils/screen";

interface IProps {
  movie: GQL.MovieDataFragment;
  containerWidth?: number;
  sceneIndex?: number;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const MovieCard: React.FC<IProps> = (props: IProps) => {
  const [cardWidth, setCardWidth] = useState<number>();

  useEffect(() => {
    if (!props.containerWidth || ScreenUtils.isMobile()) return;

    let preferredCardWidth = 250;
    let fittedCardWidth = calculateCardWidth(
      props.containerWidth,
      preferredCardWidth!
    );
    setCardWidth(fittedCardWidth);
  }, [props, props.containerWidth]);

  function maybeRenderSceneNumber() {
    if (!props.sceneIndex) return;

    return (
      <>
        <hr />
        <span className="movie-scene-number">
          <FormattedMessage id="scene" /> #{props.sceneIndex}
        </span>
      </>
    );
  }

  function maybeRenderScenesPopoverButton() {
    if (props.movie.scenes.length === 0) return;

    const popoverContent = props.movie.scenes.map((scene) => (
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
          <span>{props.movie.scenes.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (props.sceneIndex || props.movie.scenes.length > 0) {
      return (
        <>
          {maybeRenderSceneNumber()}
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <GridCard
      className="movie-card"
      url={`/movies/${props.movie.id}`}
      width={cardWidth}
      title={props.movie.name}
      linkClassName="movie-card-header"
      image={
        <>
          <img
            loading="lazy"
            className="movie-card-image"
            alt={props.movie.name ?? ""}
            src={props.movie.front_image_path ?? ""}
          />
          <RatingBanner rating={props.movie.rating100} />
        </>
      }
      details={
        <div className="movie-card__details">
          <span className="movie-card__date">{props.movie.date}</span>
          <TruncatedText
            className="movie-card__description"
            text={props.movie.synopsis}
            lineCount={3}
          />
        </div>
      }
      selected={props.selected}
      selecting={props.selecting}
      onSelectedChanged={props.onSelectedChanged}
      popovers={maybeRenderPopoverButtonGroup()}
    />
  );
};
