import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { GridCard, calculateCardWidth } from "../Shared/GridCard/GridCard";
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
import { faTag } from "@fortawesome/free-solid-svg-icons";
import { RatingBanner } from "../Shared/RatingBanner";
import { usePerformerUpdate } from "src/core/StashService";
import { ILabeledId } from "src/models/list-filter/types";
import ScreenUtils from "src/utils/screen";
import { FavoriteIcon } from "../Shared/FavoriteIcon";

export interface IPerformerCardExtraCriteria {
  scenes?: Criterion<CriterionValue>[];
  images?: Criterion<CriterionValue>[];
  galleries?: Criterion<CriterionValue>[];
  groups?: Criterion<CriterionValue>[];
  performer?: ILabeledId;
}

interface IPerformerCardProps {
  performer: GQL.PerformerDataFragment;
  containerWidth?: number;
  ageFromDate?: string;
  selecting?: boolean;
  selected?: boolean;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  extraCriteria?: IPerformerCardExtraCriteria;
}

export const PerformerCard: React.FC<IPerformerCardProps> = (
  props: IPerformerCardProps
) => {
  const intl = useIntl();
  const age = TextUtils.age(
    props.performer.birthdate,
    props.ageFromDate ?? props.performer.death_date
  );
  const ageL10nId = props.ageFromDate
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
        preferredCardWidth = 300;
        break;
      case 2:
        preferredCardWidth = 375;
        break;
      case 3:
        preferredCardWidth = 470;
    }
    let fittedCardWidth = calculateCardWidth(
      props.containerWidth,
      preferredCardWidth!
    );
    setCardWidth(fittedCardWidth);
  }, [props.containerWidth, props.zoomIndex]);

  function onToggleFavorite(v: boolean) {
    if (props.performer.id) {
      updatePerformer({
        variables: {
          input: {
            id: props.performer.id,
            favorite: v,
          },
        },
      });
    }
  }

  function maybeRenderScenesPopoverButton() {
    if (!props.performer.scene_count) return;

    return (
      <PopoverCountButton
        className="scene-count"
        type="scene"
        count={props.performer.scene_count}
        url={NavUtils.makePerformerScenesUrl(
          props.performer,
          props.extraCriteria?.performer,
          props.extraCriteria?.scenes
        )}
      />
    );
  }

  function maybeRenderImagesPopoverButton() {
    if (!props.performer.image_count) return;

    return (
      <PopoverCountButton
        className="image-count"
        type="image"
        count={props.performer.image_count}
        url={NavUtils.makePerformerImagesUrl(
          props.performer,
          props.extraCriteria?.performer,
          props.extraCriteria?.images
        )}
      />
    );
  }

  function maybeRenderGalleriesPopoverButton() {
    if (!props.performer.gallery_count) return;

    return (
      <PopoverCountButton
        className="gallery-count"
        type="gallery"
        count={props.performer.gallery_count}
        url={NavUtils.makePerformerGalleriesUrl(
          props.performer,
          props.extraCriteria?.performer,
          props.extraCriteria?.galleries
        )}
      />
    );
  }

  function maybeRenderOCounter() {
    if (!props.performer.o_counter) return;

    return (
      <div className="o-counter">
        <Button className="minimal">
          <span className="fa-icon">
            <SweatDrops />
          </span>
          <span>{props.performer.o_counter}</span>
        </Button>
      </div>
    );
  }

  function maybeRenderTagPopoverButton() {
    if (props.performer.tags.length <= 0) return;

    const popoverContent = props.performer.tags.map((tag) => (
      <TagLink key={tag.id} linkType="performer" tag={tag} />
    ));

    return (
      <HoverPopover placement="bottom" content={popoverContent}>
        <Button className="minimal tag-count">
          <Icon icon={faTag} />
          <span>{props.performer.tags.length}</span>
        </Button>
      </HoverPopover>
    );
  }

  function maybeRenderGroupsPopoverButton() {
    if (!props.performer.group_count) return;

    return (
      <PopoverCountButton
        className="group-count"
        type="group"
        count={props.performer.group_count}
        url={NavUtils.makePerformerGroupsUrl(
          props.performer,
          props.extraCriteria?.performer,
          props.extraCriteria?.groups
        )}
      />
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (
      props.performer.scene_count ||
      props.performer.image_count ||
      props.performer.gallery_count ||
      props.performer.tags.length > 0 ||
      props.performer.o_counter ||
      props.performer.group_count
    ) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderGroupsPopoverButton()}
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
    if (!props.performer.rating100) {
      return;
    }
    return <RatingBanner rating={props.performer.rating100} />;
  }

  function maybeRenderFlag() {
    if (props.performer.country) {
      return (
        <Link to={NavUtils.makePerformersCountryUrl(props.performer)}>
          <CountryFlag
            className="performer-card__country-flag"
            country={props.performer.country}
            includeOverlay
          />
          <span className="performer-card__country-string">
            {props.performer.country}
          </span>
        </Link>
      );
    }
  }

  return (
    <GridCard
      className={`performer-card zoom-${props.zoomIndex}`}
      url={`/performers/${props.performer.id}`}
      width={cardWidth}
      pretitleIcon={
        <GenderIcon className="gender-icon" gender={props.performer.gender} />
      }
      title={
        <div>
          <span className="performer-name">{props.performer.name}</span>
          {props.performer.disambiguation && (
            <span className="performer-disambiguation">
              {` (${props.performer.disambiguation})`}
            </span>
          )}
        </div>
      }
      image={
        <>
          <img
            loading="lazy"
            className="performer-card-image"
            alt={props.performer.name ?? ""}
            src={props.performer.image_path ?? ""}
          />
        </>
      }
      overlays={
        <>
          <FavoriteIcon
            favorite={props.performer.favorite}
            onToggleFavorite={onToggleFavorite}
            size="2x"
            className="hide-not-favorite"
          />
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
      selected={props.selected}
      selecting={props.selecting}
      onSelectedChanged={props.onSelectedChanged}
    />
  );
};
