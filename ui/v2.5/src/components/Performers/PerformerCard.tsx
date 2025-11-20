import React from "react";
import { Link } from "react-router-dom";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import TextUtils from "src/utils/text";
import { GridCard } from "../Shared/GridCard/GridCard";
import { CountryFlag } from "../Shared/CountryFlag";
import { HoverPopover } from "../Shared/HoverPopover";
import { Icon } from "../Shared/Icon";
import { TagLink } from "../Shared/TagLink";
import { Button, ButtonGroup } from "react-bootstrap";
import {
  ModifierCriterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { PopoverCountButton } from "../Shared/PopoverCountButton";
import GenderIcon from "./GenderIcon";
import { faLink, faTag } from "@fortawesome/free-solid-svg-icons";
import { faInstagram, faTwitter } from "@fortawesome/free-brands-svg-icons";
import { RatingBanner } from "../Shared/RatingBanner";
import { usePerformerUpdate } from "src/core/StashService";
import { ILabeledId } from "src/models/list-filter/types";
import { FavoriteIcon } from "../Shared/FavoriteIcon";
import { PatchComponent } from "src/patch";
import { ExternalLinksButton } from "../Shared/ExternalLinksButton";
import { useConfigurationContext } from "src/hooks/Config";
import { OCounterButton } from "../Shared/CountButton";

export interface IPerformerCardExtraCriteria {
  scenes?: ModifierCriterion<CriterionValue>[];
  images?: ModifierCriterion<CriterionValue>[];
  galleries?: ModifierCriterion<CriterionValue>[];
  groups?: ModifierCriterion<CriterionValue>[];
  performer?: ILabeledId;
}

interface IPerformerCardProps {
  performer: GQL.PerformerDataFragment;
  cardWidth?: number;
  ageFromDate?: string;
  selecting?: boolean;
  selected?: boolean;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
  extraCriteria?: IPerformerCardExtraCriteria;
}

const PerformerCardPopovers: React.FC<IPerformerCardProps> = PatchComponent(
  "PerformerCard.Popovers",
  ({ performer, extraCriteria }) => {
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

      return <OCounterButton value={performer.o_counter} />;
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

    function maybeRenderGroupsPopoverButton() {
      if (!performer.group_count) return;

      return (
        <PopoverCountButton
          className="group-count"
          type="group"
          count={performer.group_count}
          url={NavUtils.makePerformerGroupsUrl(
            performer,
            extraCriteria?.performer,
            extraCriteria?.groups
          )}
        />
      );
    }

    if (
      performer.scene_count ||
      performer.image_count ||
      performer.gallery_count ||
      performer.tags.length > 0 ||
      performer.o_counter ||
      performer.group_count
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

    return null;
  }
);

const PerformerCardOverlays: React.FC<IPerformerCardProps> = PatchComponent(
  "PerformerCard.Overlays",
  ({ performer }) => {
    const { configuration } = useConfigurationContext();
    const uiConfig = configuration?.ui;
    const [updatePerformer] = usePerformerUpdate();

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

    function maybeRenderLinks() {
      if (!uiConfig?.showLinksOnPerformerCard) {
        return;
      }

      if (performer.urls && performer.urls.length > 0) {
        const twitter = performer.urls.filter((u) =>
          u.match(/https?:\/\/(?:www\.)?(?:twitter|x).com\//)
        );
        const instagram = performer.urls.filter((u) =>
          u.match(/https?:\/\/(?:www\.)?instagram.com\//)
        );
        const others = performer.urls.filter(
          (u) => !twitter.includes(u) && !instagram.includes(u)
        );

        return (
          <div
            className="performer-card__links"
            style={{
              position: "absolute",
              left: "0",
              bottom: "0",
              display: "flex",
              gap: "0.5rem",
              flexDirection: "column-reverse",
            }}
          >
            {twitter.length > 0 && (
              <ExternalLinksButton
                className="performer-card__link twitter"
                urls={twitter}
                icon={faTwitter}
                openIfSingle={true}
              ></ExternalLinksButton>
            )}
            {instagram.length > 0 && (
              <ExternalLinksButton
                className="performer-card__link instagram"
                urls={instagram}
                icon={faInstagram}
                openIfSingle={true}
              ></ExternalLinksButton>
            )}
            {others.length > 0 && (
              <ExternalLinksButton
                className="performer-card__link"
                icon={faLink}
                urls={others}
                openIfSingle={true}
              />
            )}
          </div>
        );
      }
    }

    return (
      <>
        <FavoriteIcon
          favorite={performer.favorite}
          onToggleFavorite={onToggleFavorite}
          size="2x"
          className="hide-not-favorite"
        />
        {maybeRenderRatingBanner()}
        {maybeRenderLinks()}
        {maybeRenderFlag()}
      </>
    );
  }
);

const PerformerCardDetails: React.FC<IPerformerCardProps> = PatchComponent(
  "PerformerCard.Details",
  ({ performer, ageFromDate }) => {
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

    return (
      <>
        {age !== 0 ? (
          <div className="performer-card__age">{ageString}</div>
        ) : (
          ""
        )}
      </>
    );
  }
);

const PerformerCardImage: React.FC<IPerformerCardProps> = PatchComponent(
  "PerformerCard.Image",
  ({ performer }) => {
    return (
      <>
        <img
          loading="lazy"
          className="performer-card-image"
          alt={performer.name ?? ""}
          src={performer.image_path ?? ""}
        />
      </>
    );
  }
);

const PerformerCardTitle: React.FC<IPerformerCardProps> = PatchComponent(
  "PerformerCard.Title",
  ({ performer }) => {
    return (
      <div>
        <span className="performer-name">{performer.name}</span>
        {performer.disambiguation && (
          <span className="performer-disambiguation">
            {` (${performer.disambiguation})`}
          </span>
        )}
      </div>
    );
  }
);

export const PerformerCard: React.FC<IPerformerCardProps> = PatchComponent(
  "PerformerCard",
  (props) => {
    const {
      performer,
      cardWidth,
      selecting,
      selected,
      onSelectedChanged,
      zoomIndex,
    } = props;

    return (
      <GridCard
        className={`performer-card zoom-${zoomIndex}`}
        url={`/performers/${performer.id}`}
        width={cardWidth}
        pretitleIcon={
          <GenderIcon className="gender-icon" gender={performer.gender} />
        }
        title={<PerformerCardTitle {...props} />}
        image={<PerformerCardImage {...props} />}
        overlays={<PerformerCardOverlays {...props} />}
        details={<PerformerCardDetails {...props} />}
        popovers={<PerformerCardPopovers {...props} />}
        selected={selected}
        selecting={selecting}
        onSelectedChanged={onSelectedChanged}
      />
    );
  }
);
