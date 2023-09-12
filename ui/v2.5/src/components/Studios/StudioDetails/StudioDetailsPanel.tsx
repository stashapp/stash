import React from "react";
import * as GQL from "src/core/generated-graphql";
import { DetailItem } from "src/components/Shared/DetailItem";

interface IStudioDetailsPanel {
  studio: GQL.StudioDataFragment;
  collapsed?: boolean;
  fullWidth?: boolean;
}

export const StudioDetailsPanel: React.FC<IStudioDetailsPanel> = ({
  studio,
  collapsed,
  fullWidth,
}) => {
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

  function maybeRenderExtraDetails() {
    if (!collapsed) {
      return (
        <DetailItem
          id="stash_ids"
          value={renderStashIDs()}
          fullWidth={fullWidth}
        />
      );
    }
  }

  return (
    <div className="detail-group">
      <DetailItem id="details" value={studio.details} fullWidth={fullWidth} />
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
        fullWidth={fullWidth}
      />
      {maybeRenderExtraDetails()}
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
        {studio?.parent_studio?.name ? (
          <>
            <span className="detail-divider">/</span>
            <span className="studio-parent">{studio?.parent_studio?.name}</span>
          </>
        ) : (
          ""
        )}
      </div>
    </div>
  );
};
