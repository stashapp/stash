import React from "react";
import { Badge } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { DetailItem } from "src/components/Shared/DetailItem";

interface IStudioDetailsPanel {
  studio: GQL.StudioDataFragment;
}

export const StudioDetailsPanel: React.FC<IStudioDetailsPanel> = ({
  studio,
}) => {
  function renderTagsList() {
    if (!studio.aliases?.length) {
      return;
    }

    return (
      <>
        {studio.aliases.map((a) => (
          <Badge className="tag-item" variant="secondary" key={a}>
            {a}
          </Badge>
        ))}
      </>
    );
  }

  function renderStashIDs() {
    if (!studio.stash_ids?.length) {
      return;
    }

    return (
      <ul className="pl-0">
        {studio.stash_ids.map((stashID) => {
          const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
          const link = base ? (
            <a
              href={`${base}studios/${stashID.stash_id}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              {stashID.stash_id}
            </a>
          ) : (
            stashID.stash_id
          );
          return (
            <li key={stashID.stash_id} className="row no-gutters">
              {link}
            </li>
          );
        })}
      </ul>
    );
  }

  return (
    <div className="detail-group">
      <DetailItem id="details" value={studio.details} />
      <DetailItem
        id="parent_studios"
        value={
          studio.parent_studio?.name ? (
            <a href={`/studios/${studio.parent_studio?.id}`} target="_self">
              {studio.parent_studio.name}
            </a>
          ) : (
            ""
          )
        }
      />
      <DetailItem id="aliases" value={renderTagsList()} />
      <DetailItem id="StashIDs" value={renderStashIDs()} />
    </div>
  );
};

export const CompressedStudioDetailsPanel: React.FC<IStudioDetailsPanel> = ({
  studio,
}) => {
  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div className="sticky detail-header">
      <div className="sticky detail-header-group">
        <a className="studio-name" onClick={() => scrollToTop()}>
          {studio.name}
        </a>
        <span className="studio-parent">{studio?.parent_studio?.name}</span>
      </div>
    </div>
  );
};
