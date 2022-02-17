import React, { FunctionComponent } from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import {
  GridCard,
  HoverPopover,
  Icon,
  TagLink,
  TruncatedText,
} from "src/components/Shared";
import { FormattedMessage } from "react-intl";
import { RatingBanner } from "../Shared/RatingBanner";

interface IProps {
  movie: GQL.MovieDataFragment;
  sceneIndex?: number;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const MovieCard: FunctionComponent<IProps> = (props: IProps) => {
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
      <TagLink key={scene.id} scene={scene} />
    ));

    return (
      <HoverPopover
        className="scene-count"
        placement="bottom"
        content={popoverContent}
      >
        <Button className="minimal">
          <Icon icon="play-circle" />
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
      title={props.movie.name}
      linkClassName="movie-card-header"
      image={
        <>
          <img
            className="movie-card-image"
            alt={props.movie.name ?? ""}
            src={props.movie.front_image_path ?? ""}
          />
          <RatingBanner rating={props.movie.rating} />
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
