import { ReactNode } from "react";
import { Placement } from "react-bootstrap/esm/Overlay";
import * as GQL from "src/core/generated-graphql";
import { useConfigurationContext } from "src/hooks/Config";
import { HoverPopover } from "src/components/Shared/HoverPopover";
import { PerformerCard } from "src/components/Performers/PerformerCard";

interface IScrapedPerformerPreviewProps {
  performer: GQL.ScrapedPerformer;
  placement?: Placement;
  children?: ReactNode;
}

const toPerformerCardData = (performer: GQL.ScrapedPerformer) =>
  ({
    id:
      performer.stored_id ??
      performer.remote_site_id ??
      `scraped-${performer.name?.replace(/\s+/g, "-").toLowerCase() ?? "performer"}`,
    name: performer.name ?? "Unknown performer",
    disambiguation: performer.disambiguation ?? null,
    gender: performer.gender ?? null,
    birthdate: performer.birthdate ?? null,
    death_date: performer.death_date ?? null,
    country: performer.country ?? null,
    image_path: performer.images?.[0] ?? performer.image_path ?? null,
    tags: [],
    stash_ids: [],
    favorite: false,
    scene_count: null,
    image_count: null,
    gallery_count: null,
    group_count: null,
    o_counter: null,
    rating100: null,
    urls: performer.urls ?? [],
  }) as GQL.PerformerDataFragment;

const ScrapedPerformerCard = ({ performer }: { performer: GQL.ScrapedPerformer }) => {
  return (
    <div className="tag-popover-card tagger-scraped-performer-popover">
      <PerformerCard performer={toPerformerCardData(performer)} zoomIndex={0} />
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
