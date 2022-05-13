import React from "react";
import { Badge, Button } from "react-bootstrap";
import {
  Criterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { useIntl } from "react-intl";
import { Icon } from "../Shared";

interface IFilterTagsProps {
  criteria: Criterion<CriterionValue>[];
  onEditCriterion: (c: Criterion<CriterionValue>) => void;
  onRemoveCriterion: (c: Criterion<CriterionValue>) => void;
}

export const FilterTags: React.FC<IFilterTagsProps> = ({
  criteria,
  onEditCriterion,
  onRemoveCriterion,
}) => {
  const intl = useIntl();

  function onRemoveCriterionTag(
    criterion: Criterion<CriterionValue>,
    $event: React.MouseEvent<HTMLElement, MouseEvent>
  ) {
    if (!criterion) {
      return;
    }
    onRemoveCriterion(criterion);
    $event.stopPropagation();
  }

  function onClickCriterionTag(criterion: Criterion<CriterionValue>) {
    onEditCriterion(criterion);
  }

  function renderFilterTags() {
    return criteria.map((criterion) => (
      <Badge
        className="tag-item"
        variant="secondary"
        key={criterion.getId()}
        onClick={() => onClickCriterionTag(criterion)}
      >
        {criterion.getLabel(intl)}
        <Button
          variant="secondary"
          onClick={($event) => onRemoveCriterionTag(criterion, $event)}
        >
          <Icon icon="times" />
        </Button>
      </Badge>
    ));
  }

  return (
    <div className="d-flex justify-content-center mb-2 wrap-tags">
      {renderFilterTags()}
    </div>
  );
};
