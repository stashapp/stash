import React from "react";
import { Badge, Button } from "react-bootstrap";
import {
  Criterion,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import { faTimes } from "@fortawesome/free-solid-svg-icons";

interface IFilterTagsProps {
  criteria: Criterion<CriterionValue>[];
  onEditCriterion: (c: Criterion<CriterionValue>) => void;
  onRemoveCriterion: (c: Criterion<CriterionValue>) => void;
  onRemoveAll: () => void;
}

export const FilterTags: React.FC<IFilterTagsProps> = ({
  criteria,
  onEditCriterion,
  onRemoveCriterion,
  onRemoveAll,
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
          <Icon icon={faTimes} />
        </Button>
      </Badge>
    ));
  }

  function maybeRenderClearAll() {
    if (criteria.length < 3) {
      return;
    }

    return (
      <Button
        variant="minimal"
        className="clear-all-button"
        onClick={() => onRemoveAll()}
      >
        <FormattedMessage id="actions.clear" />
      </Button>
    );
  }

  return (
    <div className="d-flex justify-content-center mb-2 wrap-tags filter-tags">
      {renderFilterTags()}
      {maybeRenderClearAll()}
    </div>
  );
};
