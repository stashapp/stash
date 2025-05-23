import React, { useMemo } from "react";
import { Badge, Button } from "react-bootstrap";
import { ListFilterModel } from "src/models/list-filter/filter";
import { faFilter } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { useIntl } from "react-intl";

interface IFilterButtonProps {
  filter: ListFilterModel;
  onClick: () => void;
}

export const FilterButton: React.FC<IFilterButtonProps> = ({
  filter,
  onClick,
}) => {
  const intl = useIntl();
  const count = useMemo(() => filter.count(), [filter]);

  return (
    <Button
      variant="secondary"
      className="filter-button"
      onClick={onClick}
      title={intl.formatMessage({ id: "search_filter.edit_filter" })}
    >
      <Icon icon={faFilter} />
      {count ? (
        <Badge pill variant="info">
          {count}
        </Badge>
      ) : undefined}
    </Button>
  );
};
