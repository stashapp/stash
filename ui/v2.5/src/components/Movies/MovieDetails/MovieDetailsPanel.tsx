import React from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { DetailItem } from "src/components/Shared/DetailItem";
import { Link } from "react-router-dom";
import { DirectorLink } from "src/components/Shared/Link";
import { TagLink } from "src/components/Shared/TagLink";

interface IGroupDetailsPanel {
  group: GQL.GroupDataFragment;
  collapsed?: boolean;
  fullWidth?: boolean;
}

export const GroupDetailsPanel: React.FC<IGroupDetailsPanel> = ({
  group,
  collapsed,
  fullWidth,
}) => {
  // Network state
  const intl = useIntl();

  function renderTagsField() {
    if (!group.tags.length) {
      return;
    }
    return (
      <ul className="pl-0">
        {(group.tags ?? []).map((tag) => (
          <TagLink key={tag.id} linkType="group" tag={tag} />
        ))}
      </ul>
    );
  }

  function maybeRenderExtraDetails() {
    if (!collapsed) {
      return (
        <>
          <DetailItem
            id="synopsis"
            value={group.synopsis}
            fullWidth={fullWidth}
          />
          <DetailItem
            id="tags"
            value={renderTagsField()}
            fullWidth={fullWidth}
          />
        </>
      );
    }
  }

  return (
    <div className="detail-group">
      <DetailItem
        id="duration"
        value={
          group.duration ? TextUtils.secondsToTimestamp(group.duration) : ""
        }
        fullWidth={fullWidth}
      />
      <DetailItem
        id="date"
        value={group.date ? TextUtils.formatDate(intl, group.date) : ""}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="studio"
        value={
          group.studio?.id ? (
            <Link to={`/studios/${group.studio?.id}`}>
              {group.studio?.name}
            </Link>
          ) : (
            ""
          )
        }
        fullWidth={fullWidth}
      />

      <DetailItem
        id="director"
        value={
          group.director ? (
            <DirectorLink director={group.director} linkType="group" />
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

export const CompressedGroupDetailsPanel: React.FC<IGroupDetailsPanel> = ({
  group,
}) => {
  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div className="sticky detail-header">
      <div className="sticky detail-header-group">
        <a className="group-name" onClick={() => scrollToTop()}>
          {group.name}
        </a>
        {group?.studio?.name ? (
          <>
            <span className="detail-divider">/</span>
            <span className="group-studio">{group?.studio?.name}</span>
          </>
        ) : (
          ""
        )}
      </div>
    </div>
  );
};
