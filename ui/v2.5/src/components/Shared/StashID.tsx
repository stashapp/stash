import React, { useMemo } from "react";
import { StashId } from "src/core/generated-graphql";
import { useConfigurationContext } from "src/hooks/Config";
import { getStashboxBase } from "src/utils/stashbox";
import { ExternalLink } from "./ExternalLink";

export type LinkType = "performers" | "scenes" | "studios" | "tags";

export const StashIDPill: React.FC<{
  stashID: Pick<StashId, "endpoint" | "stash_id">;
  linkType: LinkType;
}> = ({ stashID, linkType }) => {
  const { configuration } = useConfigurationContext();

  const { endpoint, stash_id } = stashID;

  const endpointName = useMemo(() => {
    return (
      configuration?.general.stashBoxes.find((sb) => sb.endpoint === endpoint)
        ?.name ?? endpoint
    );
  }, [configuration?.general.stashBoxes, endpoint]);

  const base = getStashboxBase(endpoint);
  const link = `${base}${linkType}/${stash_id}`;

  return (
    <span className="stash-id-pill" data-endpoint={endpointName}>
      <span>{endpointName}</span>
      <ExternalLink href={link}>{stash_id}</ExternalLink>
    </span>
  );
};
