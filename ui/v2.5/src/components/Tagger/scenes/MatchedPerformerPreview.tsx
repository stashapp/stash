import { ReactNode } from "react";
import { Placement } from "react-bootstrap/esm/Overlay";
import * as GQL from "src/core/generated-graphql";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { HoverPopover } from "src/components/Shared/HoverPopover";
import { useConfigurationContext } from "src/hooks/Config";

interface IPerformerDeltaRow {
  label: string;
  value: string;
}

interface IMatchedPerformerPreviewProps {
  performerID?: string | null;
  performer?: GQL.PerformerDataFragment | null;
  placement?: Placement;
  deltaRows?: IPerformerDeltaRow[];
  warningStashID?: Pick<GQL.StashId, "endpoint" | "stash_id">;
  children?: ReactNode;
}

export const MatchedPerformerPreview = ({
  performerID,
  performer,
  placement = "bottom",
  deltaRows = [],
  warningStashID,
  children,
}: IMatchedPerformerPreviewProps) => {
  const { configuration: config } = useConfigurationContext();
  const showPerformerCardOnHover = config?.ui.showTagCardOnHover ?? true;
  const warningEndpointName = warningStashID
    ? config?.general.stashBoxes.find(
        (sb) => sb.endpoint === warningStashID.endpoint
      )?.name ?? warningStashID.endpoint
    : null;

  if (!performerID || !performer || !showPerformerCardOnHover) {
    return <>{children}</>;
  }

  return (
    <HoverPopover
      className="d-inline-block"
      placement={placement}
      enterDelay={500}
      leaveDelay={100}
      content={
        <div className="tag-popover-card tagger-matched-performer-popover">
          <PerformerCard performer={performer} zoomIndex={0} />
          {(warningStashID || deltaRows.length > 0) && (
            <div className="tagger-matched-performer-popover-extra">
              {warningStashID && (
                <div className="tagger-performer-stashid-warning">
                  <span className="stash-id-pill">
                    <span
                      className="tagger-performer-stashid-warning-chip"
                      title={warningStashID.stash_id}
                    >
                      {warningEndpointName}
                    </span>
                  </span>
                </div>
              )}
              {deltaRows.length > 0 && (
                <div className="tagger-performer-delta-rows mt-2">
                  {deltaRows.map((row) => (
                    <div key={row.label}>
                      <span>{row.label}:</span> <span>{row.value}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      }
    >
      {children}
    </HoverPopover>
  );
};
