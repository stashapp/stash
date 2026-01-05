import {
  faFilm,
  faArrowUpLong,
  faArrowDownLong,
} from "@fortawesome/free-solid-svg-icons";
import React, { useMemo } from "react";
import { Button, OverlayTrigger, Tooltip } from "react-bootstrap";
import { Count } from "../Shared/PopoverCountButton";
import { Icon } from "../Shared/Icon";
import { HoverPopover } from "../Shared/HoverPopover";
import { Link } from "react-router-dom";
import NavUtils from "src/utils/navigation";
import * as GQL from "src/core/generated-graphql";
import { useIntl } from "react-intl";
import { GroupTag } from "./GroupTag";

interface IProps {
  group: Pick<
    GQL.ListGroupDataFragment,
    "id" | "name" | "containing_groups" | "sub_group_count"
  >;
}

const ContainingGroupsCount: React.FC<IProps> = ({ group }) => {
  const { containing_groups: containingGroups } = group;

  const popoverContent = useMemo(() => {
    if (!containingGroups.length) {
      return [];
    }

    return containingGroups.map((entry) => (
      <GroupTag
        key={entry.group.id}
        linkType="sub_group"
        group={entry.group}
        description={entry.description ?? undefined}
      />
    ));
  }, [containingGroups]);

  if (!containingGroups.length) {
    return null;
  }

  return (
    <HoverPopover
      className="containing-group-count"
      placement="bottom"
      content={popoverContent}
    >
      <Link
        to={NavUtils.makeContainingGroupsUrl(group)}
        className="related-group-count"
      >
        <Count count={containingGroups.length} />
        <Icon icon={faArrowUpLong} transform="shrink-4" />
      </Link>
    </HoverPopover>
  );
};

const SubGroupCount: React.FC<IProps> = ({ group }) => {
  const intl = useIntl();

  const count = group.sub_group_count;

  if (!count) {
    return null;
  }

  function getTitle() {
    const pluralCategory = intl.formatPlural(count);
    const options = {
      one: "sub_group",
      other: "sub_groups",
    };
    const plural = intl.formatMessage({
      id: options[pluralCategory as "one"] || options.other,
    });
    return `${count} ${plural}`;
  }

  return (
    <OverlayTrigger
      overlay={<Tooltip id={`sub-group-count-tooltip`}>{getTitle()}</Tooltip>}
      placement="bottom"
    >
      <Link
        to={NavUtils.makeSubGroupsUrl(group)}
        className="related-group-count"
      >
        <Count count={count} />
        <Icon icon={faArrowDownLong} transform="shrink-4" />
      </Link>
    </OverlayTrigger>
  );
};

export const RelatedGroupPopoverButton: React.FC<IProps> = ({ group }) => {
  return (
    <span className="related-group-popover-button">
      <Button className="minimal">
        <Icon icon={faFilm} />
        <ContainingGroupsCount group={group} />
        <SubGroupCount group={group} />
      </Button>
    </span>
  );
};
