import React from "react";
import { Badge } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlus } from "@fortawesome/free-solid-svg-icons";

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
            <ParentTag key={p.id} tag={p as GQL.TagDataFragment} />
          ))}
        </dd>
      </dl>
    );
  }

  function ParentTag({ tag }: { tag: GQL.TagDataFragment }) {
    const { loading, error, data } = GQL.useFindTagQuery({
      variables: { id: tag.id },
    });

    const fullTagInfo = data?.findTag as GQL.TagDataFragment;
    const hasParent = fullTagInfo.parents?.length ?? 0;
    const iconMargin = hasParent ? '0px 0px 0px 5px' : '0';

    return (
      <Badge key={fullTagInfo.id} className="tag-item" variant="secondary">
        <Link to={`/tags/${fullTagInfo.id}`}>
          {fullTagInfo.name}
          {hasParent ? <FontAwesomeIcon icon={faPlus} style={{ margin: iconMargin }} /> : ''}
        </Link>
      </Badge>
    );
  }

  function renderChildrenField() {
    if (!tag.children?.length) {
      return null;
    }

    return (
      <dl className="row">
        <dt className="col-3 col-xl-2">
          <FormattedMessage id="sub_tags" />
        </dt>
        <dd className="col-9 col-xl-10">
          {tag.children.map((c) => (
            <ChildTag key={c.id} tag={c as GQL.TagDataFragment} />
          ))}
        </dd>
      </dl>
    );
  }

  function ChildTag({ tag }: { tag: GQL.TagDataFragment }) {
    const { loading, error, data } = GQL.useFindTagQuery({
      variables: { id: tag.id },
    });

    const fullTagInfo = data?.findTag as GQL.TagDataFragment;
    const hasChildren = fullTagInfo.children?.length ?? 0;
    const iconMargin = hasChildren ? '0px 0px 0px 5px' : '0';

    return (
      <Badge key={fullTagInfo.id} className="tag-item" variant="secondary">
        <Link to={`/tags/${fullTagInfo.id}`}>
          {fullTagInfo.name}
          {hasChildren ? <FontAwesomeIcon icon={faPlus} style={{ margin: iconMargin }} /> : ''}
        </Link>
      </Badge>
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
