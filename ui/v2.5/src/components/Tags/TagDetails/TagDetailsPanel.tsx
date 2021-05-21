import React from "react";
import { Badge } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";

interface ITagDetails {
  tag: Partial<GQL.TagDataFragment>;
}

export const TagDetailsPanel: React.FC<ITagDetails> = ({ tag }) => {
  function renderAliasesField() {
    if (!tag.aliases?.length) {
      return;
    }

    return (
      <dl className="row">
        <dt className="col-3 col-xl-2">Aliases</dt>
        <dd className="col-9 col-xl-10">
          {tag.aliases.map(a => <Badge className="tag-item" variant="secondary">{a}</Badge>)}
        </dd>
      </dl>
    );
  }

  return <>{renderAliasesField()}</>;
};
