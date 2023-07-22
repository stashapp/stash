import React from "react";
import { Badge, OverlayTrigger, Tooltip } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { Link } from "react-router-dom";
import { Icon } from "../../Shared/Icon";
import * as GQL from "src/core/generated-graphql";
import { faFolderTree } from "@fortawesome/free-solid-svg-icons";

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
          {tag.children.map((p) => (
            <Badge key={p.id} className="tag-item" variant="secondary">
              <Link to={`/tags/${p.id}`}>
                {p.name}
                {p.child_count !== 0 && (
                  <>
                    <span className="icon-wrapper">
                      <span className="vertical-line">|</span>
                      <OverlayTrigger
                        placement="top"
                        overlay={
                          <Tooltip id="tag-hierarchy-tooltip">
                            Explore tag hierarchy
                          </Tooltip>
                        }
                      >
                        <Icon icon={faFolderTree} className="tag-icon" />
                      </OverlayTrigger>
                    </span>
                  </>
                )}
              </Link>
            </Badge>
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
            <Badge key={c.id} className="tag-item" variant="secondary">
              <Link to={`/tags/${c.id}`}>
                {c.name}
                {c.child_count !== 0 && (
                  <>
                    <span className="icon-wrapper">
                      <span className="vertical-line">|</span>
                      <OverlayTrigger
                        placement="top"
                        overlay={
                          <Tooltip id="tag-hierarchy-tooltip">
                            Explore tag hierarchy
                          </Tooltip>
                        }
                      >
                        <Icon icon={faFolderTree} className="tag-icon" />
                      </OverlayTrigger>
                    </span>
                  </>
                )}
              </Link>
            </Badge>
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
