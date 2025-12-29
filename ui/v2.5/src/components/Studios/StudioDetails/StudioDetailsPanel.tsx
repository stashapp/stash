import React from "react";
import { TagLink } from "src/components/Shared/TagLink";
import * as GQL from "src/core/generated-graphql";
import { DetailItem } from "src/components/Shared/DetailItem";
import { StashIDPill } from "src/components/Shared/StashID";
import { PatchComponent } from "src/patch";
import { Link } from "react-router-dom";

interface IStudioDetailsPanel {
  studio: GQL.StudioDataFragment;
  collapsed?: boolean;
  fullWidth?: boolean;
}

export const StudioDetailsPanel: React.FC<IStudioDetailsPanel> = PatchComponent(
  "StudioDetailsPanel",
  ({ studio, fullWidth }) => {
    function renderTagsField() {
      if (!studio.tags.length) {
        return;
      }
      return (
        <ul className="pl-0">
          {(studio.tags ?? []).map((tag) => (
            <TagLink key={tag.id} linkType="studio" tag={tag} />
          ))}
        </ul>
      );
    }

    function renderStashIDs() {
      if (!studio.stash_ids?.length) {
        return;
      }

      return (
        <ul className="pl-0">
          {studio.stash_ids.map((stashID) => {
            return (
              <li key={stashID.stash_id} className="row no-gutters">
                <StashIDPill stashID={stashID} linkType="studios" />
              </li>
            );
          })}
        </ul>
      );
    }

    function renderURLs() {
      if (!studio.urls?.length) {
        return;
      }

      return (
        <ul className="pl-0">
          {studio.urls.map((url) => (
            <li key={url}>
              <a href={url} target="_blank" rel="noreferrer">
                {url}
              </a>
            </li>
          ))}
        </ul>
      );
    }

    return (
      <div className="detail-group">
        <DetailItem id="details" value={studio.details} fullWidth={fullWidth} />
        <DetailItem id="urls" value={renderURLs()} fullWidth={fullWidth} />
        <DetailItem
          id="parent_studios"
          value={
            studio.parent_studio?.name ? (
              <Link to={`/studios/${studio.parent_studio?.id}`}>
                {studio.parent_studio.name}
              </Link>
            ) : (
              ""
            )
          }
          fullWidth={fullWidth}
        />
        <DetailItem id="tags" value={renderTagsField()} fullWidth={fullWidth} />
        <DetailItem
          id="stash_ids"
          value={renderStashIDs()}
          fullWidth={fullWidth}
        />
      </div>
    );
  }
);

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
