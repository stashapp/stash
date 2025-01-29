import React, { PropsWithChildren } from "react";
import { Badge, BadgeProps, Button } from "react-bootstrap";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import { BsPrefixProps, ReplaceProps } from "react-bootstrap/esm/helpers";

type TagItemProps = PropsWithChildren<
  ReplaceProps<"span", BsPrefixProps<"span"> & BadgeProps>
>;

export const TagItem: React.FC<TagItemProps> = (props) => {
  const { children } = props;
  return (
    <Badge className="tag-item" variant="secondary" {...props}>
      {children}
    </Badge>
  );
};

export const FilterTag: React.FC<{
  label: React.ReactNode;
  onClick: React.MouseEventHandler<HTMLSpanElement>;
  onRemove: React.MouseEventHandler<HTMLElement>;
}> = ({ label, onClick, onRemove }) => {
  return (
    <TagItem onClick={onClick}>
      {label}
      <Button
        variant="secondary"
        onClick={(e) => {
          onRemove(e);
          e.stopPropagation();
        }}
      >
        <Icon icon={faTimes} />
      </Button>
    </TagItem>
  );
};

interface IFilterTagsProps {
  criteria: Criterion[];
  onEditCriterion: (c: Criterion) => void;
  onRemoveCriterion: (c: Criterion) => void;
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
    criterion: Criterion,
    $event: React.MouseEvent<HTMLElement, MouseEvent>
  ) {
    if (!criterion) {
      return;
    }
    onRemoveCriterion(criterion);
    $event.stopPropagation();
  }

  function onClickCriterionTag(criterion: Criterion) {
    onEditCriterion(criterion);
  }

  function renderFilterTags() {
    return criteria.map((criterion) => (
      <FilterTag
        key={criterion.getId()}
        label={criterion.getLabel(intl)}
        onClick={() => onClickCriterionTag(criterion)}
        onRemove={($event) => onRemoveCriterionTag(criterion, $event)}
      />
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
