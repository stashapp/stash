import { ReactNode } from "react";
import { Placement } from "react-bootstrap/esm/Overlay";
import { IntlShape, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { useConfigurationContext } from "src/hooks/Config";
import { HoverPopover } from "src/components/Shared/HoverPopover";
import { PerformerCard } from "src/components/Performers/PerformerCard";

interface IScrapedPerformerPreviewProps {
  performer: GQL.ScrapedPerformer;
  placement?: Placement;
  children?: ReactNode;
}

const toPerformerCardData = (
  performer: GQL.ScrapedPerformer,
  intl: IntlShape
) => {
  const aliasList = performer.aliases
    ? performer.aliases
        .split(",")
        .map((a) => a.trim())
        .filter(Boolean)
    : [];
  const unknownPerformerName = intl.formatMessage({
    id: "component_tagger.results.unnamed",
    defaultMessage: "Unnamed",
  });
  return {
    id:
      performer.stored_id ??
      performer.remote_site_id ??
      null,
    name: performer.name ?? unknownPerformerName,
    alias_list: aliasList,
    disambiguation: performer.disambiguation ?? null,
    gender: performer.gender ?? null,
    birthdate: performer.birthdate ?? null,
    death_date: performer.death_date ?? null,
    country: performer.country ?? null,
    image_path: performer.images?.[0] ?? null,
    tags: [],
    custom_fields: [],
    stash_ids: [],
    favorite: false,
    ignore_auto_tag: false,
    scene_count: null,
    image_count: null,
    gallery_count: null,
    group_count: null,
    performer_count: null,
    o_counter: null,
    rating100: null,
    urls: performer.urls ?? [],
  } as unknown as GQL.PerformerDataFragment;
};

const ScrapedPerformerCard = ({
  performer,
}: {
  performer: GQL.ScrapedPerformer;
}) => {
  const intl = useIntl();
  return (
    <div className="tag-popover-card tagger-scraped-performer-popover">
      <PerformerCard
        performer={toPerformerCardData(performer, intl)}
        zoomIndex={0}
      />
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
      content={<ScrapedPerformerCard performer={performer} />}
    >
      {children}
    </HoverPopover>
  );
};
