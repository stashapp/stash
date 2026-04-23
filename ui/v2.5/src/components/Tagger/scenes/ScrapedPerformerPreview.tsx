import { ReactNode } from "react";
import { Placement } from "react-bootstrap/esm/Overlay";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  IPerformerPreviewData,
} from "src/components/Performers/PerformerPreviewCard";
import { PerformerPopover } from "src/components/Performers/PerformerPopover";
import TextUtils from "src/utils/text";
import { stringToGender } from "src/utils/gender";

interface IScrapedPerformerPreviewProps {
  performer: GQL.ScrapedPerformer;
  placement?: Placement;
  children?: ReactNode;
}

export const localPerformerToPreviewData = (
  performer: GQL.PerformerDataFragment,
  ageString: string | null
): IPerformerPreviewData => ({
  name: performer.name,
  image: performer.image_path,
  country: performer.country,
  gender: performer.gender,
  disambiguation: performer.disambiguation,
  ageString,
});

export const scrapedPerformerToPreviewData = (
  performer: GQL.ScrapedPerformer,
  name: string,
  ageString: string | null
): IPerformerPreviewData => ({
  name,
  image: performer.images?.[0],
  country: performer.country,
  gender: stringToGender(performer.gender, true),
  disambiguation: performer.disambiguation,
  ageString,
});

export const ScrapedPerformerPreview = ({
  performer,
  placement = "bottom",
  children,
}: IScrapedPerformerPreviewProps) => {
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
  const previewData = scrapedPerformerToPreviewData(
    performer,
    name,
    age !== 0 ? ageString : null
  );

  return (
    <PerformerPopover
      previewData={previewData}
      placement={placement}
      triggerClassName="d-inline-block"
    >
      {children}
    </PerformerPopover>
  );
};
