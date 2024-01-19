import React, { useMemo } from "react";
import { Badge, Button, OverlayTrigger, Tooltip } from "react-bootstrap";
import { ListFilterModel } from "src/models/list-filter/filter";
import { IconDefinition, faFilter } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { FormattedMessage } from "react-intl";

interface IFilterButtonProps {
  filter: ListFilterModel;
  icon?: IconDefinition;
  onClick: () => void;
}

export const FilterButton: React.FC<IFilterButtonProps> = ({
  filter,
  icon = faFilter,
  onClick,
}) => {
  const count = useMemo(() => filter.count(), [filter]);

  return (
    <OverlayTrigger
      placement="top"
      overlay={
        <Tooltip id="filter-tooltip">
          <FormattedMessage id="search_filter.name" />
        </Tooltip>
      }
    >
      <Button variant="secondary" className="filter-button" onClick={onClick}>
        <Icon icon={icon} />
        {count ? (
          <Badge pill variant="info">
            {count}
          </Badge>
        ) : undefined}
      </Button>
    </OverlayTrigger>
  );
};
