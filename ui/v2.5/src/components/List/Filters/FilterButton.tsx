import React from "react";
import { Badge, Button } from "react-bootstrap";
import { faFilter } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { useIntl } from "react-intl";

interface IFilterButtonProps {
  count?: number;
  onClick: () => void;
}

export const FilterButton: React.FC<IFilterButtonProps> = ({
  count = 0,
  onClick,
}) => {
  const intl = useIntl();

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
