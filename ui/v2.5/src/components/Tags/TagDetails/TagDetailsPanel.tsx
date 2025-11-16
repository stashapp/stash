import React from "react";
import { TagLink } from "src/components/Shared/TagLink";
import { DetailItem } from "src/components/Shared/DetailItem";
import { StashIDPill } from "src/components/Shared/StashID";
import * as GQL from "src/core/generated-graphql";

interface ITagDetails {
  tag: GQL.TagDataFragment;
  fullWidth?: boolean;
}

export const TagDetailsPanel: React.FC<ITagDetails> = ({ tag, fullWidth }) => {
  function renderParentsField() {
    if (!tag.parents?.length) {
      return;
    }

    return (
      <>
        {tag.parents.map((p) => (
          <TagLink
            key={p.id}
            tag={p}
            hoverPlacement="bottom"
            linkType="details"
            showHierarchyIcon={p.parent_count !== 0}
            hierarchyTooltipID="tag_parent_tooltip"
          />
        ))}
      </>
    );
  }

  function renderChildrenField() {
    if (!tag.children?.length) {
      return;
    }

    return (
      <>
        {tag.children.map((c) => (
          <TagLink
            key={c.id}
            tag={c}
            hoverPlacement="bottom"
            linkType="details"
            showHierarchyIcon={c.child_count !== 0}
            hierarchyTooltipID="tag_sub_tag_tooltip"
          />
        ))}
      </>
    );
  }

  function renderStashIDs() {
    if (!tag.stash_ids?.length) {
      return;
    }

    return (
      <ul className="pl-0">
        {tag.stash_ids.map((stashID) => (
          <li key={stashID.stash_id} className="row no-gutters">
            <StashIDPill stashID={stashID} linkType="tags" />
          </li>
        ))}
      </ul>
    );
  }

  return (
    <div className="detail-group">
      <DetailItem
        id="description"
        value={tag.description}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="parent_tags"
        value={renderParentsField()}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="sub_tags"
        value={renderChildrenField()}
        fullWidth={fullWidth}
      />
      <DetailItem
        id="stash_ids"
        value={renderStashIDs()}
        fullWidth={fullWidth}
      />
    </div>
  );
};

export const CompressedTagDetailsPanel: React.FC<ITagDetails> = ({ tag }) => {
  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  return (
    <div className="sticky detail-header">
      <div className="sticky detail-header-group">
        <a className="tag-name" onClick={() => scrollToTop()}>
          {tag.name}
        </a>
        {tag.description ? (
          <>
            <span className="detail-divider">/</span>
            <span className="tag-desc">{tag.description}</span>
          </>
        ) : (
          ""
        )}
      </div>
    </div>
  );
};
