import React from "react";
import { Link } from "react-router-dom";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { NavUtils, TextUtils } from "src/utils";
import {
  GridCard,
  CountryFlag,
  HoverPopover,
  Icon,
  TagLink,
} from "src/components/Shared";
import { Button, ButtonGroup } from "react-bootstrap";
import {
  Criterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { PopoverCountButton } from "../Shared/PopoverCountButton";
import GenderIcon from "./GenderIcon";

export interface IPerformerCardExtraCriteria {
  scenes: Criterion<CriterionValue>[];
  images: Criterion<CriterionValue>[];
  galleries: Criterion<CriterionValue>[];
  movies: Criterion<CriterionValue>[];
}

interface IPerformerCardProps {
  performer: GQL.PerformerDataFragment;
  ageFromDate?: string;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  extraCriteria?: IPerformerCardExtraCriteria;
}

export const PerformerCard: React.FC<IPerformerCardProps> = ({
  performer,
  ageFromDate,
  selecting,
  selected,
  onSelectedChanged,
  extraCriteria,
}) => {
  const intl = useIntl();
  const age = TextUtils.age(
    performer.birthdate,
    ageFromDate ?? performer.death_date
  );
  const ageL10nId = ageFromDate
    ? "media_info.performer_card.age_context"
    : "media_info.performer_card.age";
  const ageL10String = intl.formatMessage({
    id: "years_old",
    defaultMessage: "years old",
  });
  const ageString = intl.formatMessage(
    { id: ageL10nId },
    { age, years_old: ageL10String }
  );

  function maybeRenderFavoriteIcon() {
    if (performer.favorite === false) {
      return;
    }
    return (
      <div className="favorite">
        <Icon icon="heart" size="2x" />
      </div>
    );
  }

  function maybeRenderScenesPopoverButton() {
    if (!performer.scene_count) return;

    return (
      <PopoverCountButton
        className="scene-count"
        type="scene"
        count={performer.scene_count}
        url={NavUtils.makePerformerScenesUrl(performer, extraCriteria?.scenes)}
      />
    );
  }

  function maybeRenderImagesPopoverButton() {
    if (!performer.image_count) return;

    return (
      <PopoverCountButton
        className="image-count"
        type="image"
        count={performer.image_count}
        url={NavUtils.makePerformerImagesUrl(performer, extraCriteria?.images)}
      />
    );
  }

  function maybeRenderGalleriesPopoverButton() {
    if (!performer.gallery_count) return;

    return (
      <PopoverCountButton
        className="gallery-count"
        type="gallery"
        count={performer.gallery_count}
        url={NavUtils.makePerformerGalleriesUrl(
          performer,
          extraCriteria?.galleries
        )}
      />
    );
  }

  function maybeRenderTagPopoverButton() {
    if (performer.tags.length <= 0) return;

    const popoverContent = performer.tags.map((tag) => (
      <TagLink key={tag.id} tagType="performer" tag={tag} />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal tag-count">
          <Icon icon="tag" />
          <span>{performer.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderMoviesPopoverButton() {
    if (!performer.movie_count) return;

    return (
      <PopoverCountButton
        className="movie-count"
        type="movie"
        count={performer.movie_count}
        url={NavUtils.makePerformerMoviesUrl(performer, extraCriteria?.movies)}
      />
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (
      performer.scene_count ||
      performer.image_count ||
      performer.gallery_count ||
      performer.tags.length > 0 ||
      performer.movie_count
    ) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderMoviesPopoverButton()}
            {maybeRenderImagesPopoverButton()}
            {maybeRenderGalleriesPopoverButton()}
            {maybeRenderTagPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  function maybeRenderRatingBanner() {
    if (!performer.rating) {
      return;
    }
    return (
      <div
        className={`rating-banner ${
          performer.rating ? `rating-${performer.rating}` : ""
        }`}
      >
        <FormattedMessage id="rating" />: {performer.rating}
      </div>
    );
  }

  function maybeRenderFlag() {
    if (performer.country) {
      return (
        <Link to={NavUtils.makePerformersCountryUrl(performer)}>
          <CountryFlag
            className="performer-card__country-flag"
            country={performer.country}
          />
          <span className="performer-card__country-string">
            {performer.country}
          </span>
        </Link>
      );
    }
  }

  return (
    <GridCard
      className="performer-card"
      url={`/performers/${performer.id}`}
      pretitleIcon={
        <GenderIcon className="gender-icon" gender={performer.gender} />
      }
      title={performer.name ?? ""}
      image={
        <>
          <img
            className="performer-card-image"
            alt={performer.name ?? ""}
            src={performer.image_path ?? ""}
          />
          {maybeRenderFavoriteIcon()}
          {maybeRenderRatingBanner()}
          {maybeRenderFlag()}
        </>
      }
      details={
        <>
          {age !== 0 ? (
            <div className="performer-card__age">{ageString}</div>
          ) : (
            ""
          )}
          {maybeRenderPopoverButtonGroup()}
        </>
      }
      selected={selected}
      selecting={selecting}
      onSelectedChanged={onSelectedChanged}
    />
  );
};
