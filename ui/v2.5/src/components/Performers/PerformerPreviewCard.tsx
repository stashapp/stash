import * as GQL from "src/core/generated-graphql";
import { CountryFlag } from "src/components/Shared/CountryFlag";
import GenderIcon from "./GenderIcon";

export interface IPerformerPreviewData {
  name: string;
  image?: string | null;
  country?: string | null;
  gender?: GQL.Maybe<GQL.GenderEnum>;
  disambiguation?: string | null;
  ageString?: string | null;
}

export const PerformerPreviewCard = ({
  name,
  image,
  country,
  gender,
  disambiguation,
  ageString,
}: IPerformerPreviewData) => (
  <div className="tag-popover-card tagger-performer-popover">
    <div className="card performer-card zoom-0">
      <div className="thumbnail-section">
        <img
          loading="lazy"
          className="performer-card-image"
          alt={name}
          src={image ?? ""}
        />
        {country && (
          <CountryFlag
            className="performer-card__country-flag"
            country={country}
            includeOverlay
          />
        )}
      </div>
      <div className="card-section">
        <div className="performer-card__title">
          <GenderIcon className="gender-icon" gender={gender} />
          <span className="performer-name">{name}</span>
          {disambiguation && (
            <span className="performer-disambiguation">{` (${disambiguation})`}</span>
          )}
        </div>
        {ageString ? <div className="performer-card__age">{ageString}</div> : null}
      </div>
    </div>
  </div>
);
