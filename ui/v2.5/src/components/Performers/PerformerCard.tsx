import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { GridCard, calculateCardWidth } from "../Shared/GridCard";
import { CountryFlag } from "../Shared/CountryFlag";
import { SweatDrops } from "../Shared/SweatDrops";
import { HoverPopover } from "../Shared/HoverPopover";
import { Icon } from "../Shared/Icon";
import { TagLink } from "../Shared/TagLink";
import { Button, ButtonGroup } from "react-bootstrap";
import {
  Criterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { PopoverCountButton } from "../Shared/PopoverCountButton";
import GenderIcon from "./GenderIcon";
import { faHeart, faTag } from "@fortawesome/free-solid-svg-icons";
import { RatingBanner } from "../Shared/RatingBanner";
import cx from "classnames";
import { usePerformerUpdate } from "src/core/StashService";
import { ILabeledId } from "src/models/list-filter/types";
import ScreenUtils from "src/utils/screen";

export interface IPerformerCardExtraCriteria {
  scenes?: Criterion<CriterionValue>[];
  images?: Criterion<CriterionValue>[];
  galleries?: Criterion<CriterionValue>[];
  movies?: Criterion<CriterionValue>[];
  performer?: ILabeledId;
}

interface IPerformerCardProps {
  performer: GQL.PerformerDataFragment;
  containerWidth?: number;
  ageFromDate?: string;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  extraCriteria?: IPerformerCardExtraCriteria;
}

export const PerformerCard: React.FC<IPerformerCardProps> = ({
  performer,
  containerWidth,
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

  const [updatePerformer] = usePerformerUpdate();
  const [cardWidth, setCardWidth] = useState<number>();

  useEffect(() => {
    if (!containerWidth || ScreenUtils.isMobile()) return;

    let preferredCardWidth = 300;
    let fittedCardWidth = calculateCardWidth(
      containerWidth,
      preferredCardWidth!
    );
    setCardWidth(fittedCardWidth);
  }, [containerWidth]);

  function renderFavoriteIcon() {
    return (
      <Link to="" onClick={(e) => e.preventDefault()}>
        <Button
          className={cx(
            "minimal",
            "mousetrap",
            "favorite-button",
            performer.favorite ? "favorite" : "not-favorite"
          )}
          onClick={() => onToggleFavorite!(!performer.favorite)}
        >
          <Icon icon={faHeart} size="2x" />
        </Button>
      </Link>
    );
  }

  function onToggleFavorite(v: boolean) {
    if (performer.id) {
      updatePerformer({
        variables: {
          input: {
            id: performer.id,
            favorite: v,
          },
        },
      });
    }
  }

  function maybeRenderScenesPopoverButton() {
    if (!performer.scene_count) return;

    return (
      <PopoverCountButton
        className="scene-count"
        type="scene"
        count={performer.scene_count}
        url={NavUtils.makePerformerScenesUrl(
          performer,
          extraCriteria?.performer,
          extraCriteria?.scenes
        )}
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
        url={NavUtils.makePerformerImagesUrl(
          performer,
          extraCriteria?.performer,
          extraCriteria?.images
        )}
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
          extraCriteria?.performer,
          extraCriteria?.galleries
        )}
      />
    );
  }

  function maybeRenderOCounter() {
    if (!performer.o_counter) return;

    return (
      <div className="o-counter">
        <Button className="minimal">
          <span className="fa-icon">
            <SweatDrops />
          </span>
          <span>{performer.o_counter}</span>
        </Button>
      </div>
    );
  }

  function maybeRenderTagPopoverButton() {
    if (performer.tags.length <= 0) return;

    const popoverContent = performer.tags.map((tag) => (
      <TagLink key={tag.id} linkType="performer" tag={tag} />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal tag-count">
          <Icon icon={faTag} />
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
        url={NavUtils.makePerformerMoviesUrl(
          performer,
          extraCriteria?.performer,
          extraCriteria?.movies
        )}
      />
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (
      performer.scene_count ||
      performer.image_count ||
      performer.gallery_count ||
      performer.tags.length > 0 ||
      performer.o_counter ||
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
            {maybeRenderOCounter()}
          </ButtonGroup>
        </>
      );
    }
  }

  function maybeRenderRatingBanner() {
    if (!performer.rating100) {
      return;
    }
    return <RatingBanner rating={performer.rating100} />;
  }

  function maybeRenderFlag() {
    if (performer.country) {
      return (
        <Link to={NavUtils.makePerformersCountryUrl(performer)}>
          <CountryFlag
            className="performer-card__country-flag"
            country={performer.country}
            includeOverlay
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
      width={cardWidth}
      pretitleIcon={
        <GenderIcon className="gender-icon" gender={performer.gender} />
      }
      title={
        <div>
          <span className="performer-name">{performer.name}</span>
          {performer.disambiguation && (
            <span className="performer-disambiguation">
              {` (${performer.disambiguation})`}
            </span>
          )}
        </div>
      }
      image={
        <>
          <img
            loading="lazy"
            className="performer-card-image"
            alt={performer.name ?? ""}
            src={performer.image_path ?? ""}
          />
        </>
      }
      overlays={
        <>
          {renderFavoriteIcon()}
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
        </>
      }
      popovers={maybeRenderPopoverButtonGroup()}
      selected={selected}
      selecting={selecting}
      onSelectedChanged={onSelectedChanged}
    />
  );
};
