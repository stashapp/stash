import React from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { stringToGender } from "src/utils/gender";
import { getStashboxBase } from "src/utils/stashbox";
import { GridCard } from "src/components/Shared/GridCard/GridCard";
import { CountryFlag } from "src/components/Shared/CountryFlag";
import { PatchComponent } from "src/patch";
import GenderIcon from "src/components/Performers/GenderIcon";

export interface IScrapedPerformerCardProps {
  scrapedPerformer: GQL.ScrapedPerformer;
  endpoint: string;
  cardWidth?: number;
  ageFromDate?: string | null;
  selecting?: boolean;
  selected?: boolean;
  zoomIndex?: number;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

const ScrapedPerformerCardOverlays: React.FC<IScrapedPerformerCardProps> =
  PatchComponent("ScrapedPerformerCard.Overlays", ({ scrapedPerformer }) => {
    function maybeRenderFlag() {
      if (!scrapedPerformer.country) {
        return;
      }
      return (
        <CountryFlag
          className="performer-card__country-flag"
          country={scrapedPerformer.country}
          includeOverlay
        />
      );
    }

    return <>{maybeRenderFlag()}</>;
  });

const ScrapedPerformerCardDetails: React.FC<IScrapedPerformerCardProps> =
  PatchComponent("ScrapedPerformerCard.Details", (props) => {
    const { scrapedPerformer, ageFromDate } = props;
    const intl = useIntl();
    const age = TextUtils.age(
      scrapedPerformer.birthdate,
      ageFromDate ?? scrapedPerformer.death_date
    );
    const ageL10nId = ageFromDate
      ? "media_info.performer_card.age_context"
      : "media_info.performer_card.age";
    const ageL10String = intl.formatMessage({
      id: "years_old",
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
  });

const ScrapedPerformerCardImage: React.FC<IScrapedPerformerCardProps> =
  PatchComponent("ScrapedPerformerCard.Image", ({ scrapedPerformer }) => {
    const intl = useIntl();
    const unknownName = intl.formatMessage({
      id: "component_tagger.results.unnamed",
    });
    const alt = scrapedPerformer.name ?? unknownName;
    return (
      <img
        loading="lazy"
        className="performer-card-image"
        alt={alt}
        src={scrapedPerformer.images?.[0] ?? ""}
      />
    );
  });

const ScrapedPerformerCardTitle: React.FC<IScrapedPerformerCardProps> =
  PatchComponent("ScrapedPerformerCard.Title", ({ scrapedPerformer }) => {
    const intl = useIntl();
    const unknownPerformerName = intl.formatMessage({
      id: "component_tagger.results.unnamed",
    });
    const name = scrapedPerformer.name ?? unknownPerformerName;
    return (
      <div>
        <span className="performer-name">{name}</span>
        {scrapedPerformer.disambiguation && (
          <span className="performer-disambiguation">
            {` (${scrapedPerformer.disambiguation})`}
          </span>
        )}
      </div>
    );
  });

export const ScrapedPerformerCard: React.FC<IScrapedPerformerCardProps> =
  PatchComponent("ScrapedPerformerCard", (props) => {
    const {
      scrapedPerformer,
      endpoint,
      cardWidth,
      selecting,
      selected,
      onSelectedChanged,
      zoomIndex,
    } = props;

    const base = getStashboxBase(endpoint);
    const stashboxUrl =
      base && scrapedPerformer.remote_site_id
        ? `${base}performers/${scrapedPerformer.remote_site_id}`
        : "#";

    return (
      <GridCard
        className={`performer-card zoom-${zoomIndex}`}
        url={stashboxUrl}
        width={cardWidth}
        pretitleIcon={
          <GenderIcon
            className="gender-icon"
            gender={stringToGender(scrapedPerformer.gender, true)}
          />
        }
        title={<ScrapedPerformerCardTitle {...props} />}
        image={<ScrapedPerformerCardImage {...props} />}
        overlays={<ScrapedPerformerCardOverlays {...props} />}
        details={<ScrapedPerformerCardDetails {...props} />}
        selected={selected}
        selecting={selecting}
        onSelectedChanged={onSelectedChanged}
      />
    );
  });
