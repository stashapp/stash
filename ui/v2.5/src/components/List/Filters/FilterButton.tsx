import React from "react";
import { Badge, Button } from "react-bootstrap";
import { faFilter } from "@fortawesome/free-solid-svg-icons";
import { Icon } from "src/components/Shared/Icon";
import { useIntl } from "react-intl";

interface IFilterButtonProps {
  count?: number;
  onClick: () => void;
  title?: string;
}

export const FilterButton: React.FC<IFilterButtonProps> = ({
  count = 0,
  onClick,
  title,
}) => {
  const intl = useIntl();

  if (!title) {
    title = intl.formatMessage({ id: "search_filter.edit_filter" });
  }

  return (
    <Button
      variant="secondary"
      className="filter-button"
      onClick={onClick}
      title={title}
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
