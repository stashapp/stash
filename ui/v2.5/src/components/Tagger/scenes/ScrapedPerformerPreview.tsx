import { ReactNode } from "react";
import { Placement } from "react-bootstrap/esm/Overlay";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { useConfigurationContext } from "src/hooks/Config";
import { HoverPopover } from "src/components/Shared/HoverPopover";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import GenderIcon from "src/components/Performers/GenderIcon";
import { CountryFlag } from "src/components/Shared/CountryFlag";
import TextUtils from "src/utils/text";
import { stringToGender } from "src/utils/gender";

interface IScrapedPerformerPreviewProps {
  performer: GQL.ScrapedPerformer;
  placement?: Placement;
  children?: ReactNode;
}

export const LocalPerformerCard = ({
  performer,
}: {
  performer: GQL.PerformerDataFragment;
}) => (
  <div className="tag-popover-card tagger-matched-performer-popover">
    <PerformerCard performer={performer} zoomIndex={0} />
  </div>
);

export const RemotePerformerCard = ({
  performer,
}: {
  performer: GQL.ScrapedPerformer;
}) => {
  const intl = useIntl();
  const unknownPerformerName = intl.formatMessage({
    id: "component_tagger.results.unnamed",
    defaultMessage: "Unnamed",
  });
  const name = performer.name ?? unknownPerformerName;
  const age = TextUtils.age(performer.birthdate, performer.death_date);
  const ageString = intl.formatMessage(
    { id: "media_info.performer_card.age" },
    {
      age,
      years_old: intl.formatMessage({
        id: "years_old",
        defaultMessage: "years old",
      }),
    }
  );

  return (
    <div className="tag-popover-card tagger-scraped-performer-popover">
      <div className="card performer-card zoom-0">
        <div className="thumbnail-section">
          <img
            loading="lazy"
            className="performer-card-image"
            alt={name}
            src={performer.images?.[0] ?? ""}
          />
          {performer.country && (
            <CountryFlag
              className="performer-card__country-flag"
              country={performer.country}
              includeOverlay
            />
          )}
        </div>
        <div className="card-section">
          <div className="performer-card__title">
            <GenderIcon
              className="gender-icon"
              gender={stringToGender(performer.gender, true)}
            />
            <span className="performer-name">{name}</span>
            {performer.disambiguation && (
              <span className="performer-disambiguation">
                {` (${performer.disambiguation})`}
              </span>
            )}
          </div>
          {age !== 0 && <div className="performer-card__age">{ageString}</div>}
        </div>
      </div>
    </div>
  );
};

export const ScrapedPerformerPreview = ({
  performer,
  placement = "bottom",
  children,
}: IScrapedPerformerPreviewProps) => {
  const { configuration: config } = useConfigurationContext();
  const showPerformerCardOnHover = config?.ui.showTagCardOnHover ?? true;

  if (!showPerformerCardOnHover) {
    return <>{children}</>;
  }

  return (
    <HoverPopover
      className="d-inline-block"
      placement={placement}
      enterDelay={500}
      leaveDelay={100}
      content={<RemotePerformerCard performer={performer} />}
    >
      {children}
    </HoverPopover>
  );
};
