import React from "react";
import { Badge } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { TagLink } from "src/components/Shared/TagLink";
import * as GQL from "src/core/generated-graphql";

interface ITagDetails {
  tag: GQL.TagDataFragment;
}

export const TagDetailsPanel: React.FC<ITagDetails> = ({ tag }) => {
  function renderAliasesField() {
    if (!tag.aliases.length) {
      return;
    }

    return (
      <dl className="row">
        <dt className="col-3 col-xl-2">
          <FormattedMessage id="aliases" />
        </dt>
        <dd className="col-9 col-xl-10">
          {tag.aliases.map((a) => (
            <Badge className="tag-item" variant="secondary" key={a}>
              {a}
            </Badge>
          ))}
        </dd>
      </dl>
    );
  }

  function renderParentsField() {
    if (!tag.parents?.length) {
      return;
    }

    return (
      <dl className="row">
        <dt className="col-3 col-xl-2">
          <FormattedMessage id="parent_tags" />
        </dt>
        <dd className="col-9 col-xl-10">
          {tag.parents.map((p) => (
            <TagLink key={p.id} tag={p} />
          ))}
        </dd>
      </dl>
    );
  }

  function renderChildrenField() {
    if (!tag.children?.length) {
      return;
    }

    return (
      <dl className="row">
        <dt className="col-3 col-xl-2">
          <FormattedMessage id="sub_tags" />
        </dt>
        <dd className="col-9 col-xl-10">
          {tag.children.map((c) => (
            <TagLink key={c.id} tag={c} />
          ))}
        </dd>
      </dl>
    );
  }

  return (
    <>
      {renderAliasesField()}
      {renderParentsField()}
      {renderChildrenField()}
    </>
  );
};
