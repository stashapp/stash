import React, {
  PropsWithChildren,
  useEffect,
  useLayoutEffect,
  useReducer,
  useRef,
} from "react";
import { Badge, BadgeProps, Button, Overlay, Popover } from "react-bootstrap";
import { Criterion } from "src/models/list-filter/criteria/criterion";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import { BsPrefixProps, ReplaceProps } from "react-bootstrap/esm/helpers";
import { CustomFieldsCriterion } from "src/models/list-filter/criteria/custom-fields";
import { useDebounce } from "src/hooks/debounce";

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

const MoreFilterTags: React.FC<{
  tags: React.ReactNode[];
}> = ({ tags }) => {
  const [showTooltip, setShowTooltip] = React.useState(false);
  const target = useRef(null);

  if (!tags.length) {
    return null;
  }

  function handleMouseEnter() {
    setShowTooltip(true);
  }

  function handleMouseLeave() {
    setShowTooltip(false);
  }

  return (
    <>
      <Overlay target={target.current} placement="bottom" show={showTooltip}>
        <Popover
          id="more-criteria-popover"
          className="hover-popover-content"
          onMouseEnter={handleMouseEnter}
          onMouseLeave={handleMouseLeave}
          onClick={handleMouseLeave}
        >
          {tags}
        </Popover>
      </Overlay>
      <Badge
        ref={target}
        className={"tag-item more-tags"}
        variant="secondary"
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
      >
        <FormattedMessage
          id="search_filter.more_filter_criteria"
          values={{ count: tags.length }}
        />
      </Badge>
    </>
  );
};

interface IFilterTagsProps {
  criteria: Criterion[];
  onEditCriterion: (c: Criterion) => void;
  onRemoveCriterion: (c: Criterion, valueIndex?: number) => void;
  onRemoveAll: () => void;
  truncateOnOverflow?: boolean;
}

export const FilterTags: React.FC<IFilterTagsProps> = ({
  criteria,
  onEditCriterion,
  onRemoveCriterion,
  onRemoveAll,
  truncateOnOverflow = false,
}) => {
  const intl = useIntl();
  const ref = useRef<HTMLDivElement>(null);

  const [cutoff, setCutoff] = React.useState<number | undefined>();
  const elementGap = 10; // Adjust this value based on your CSS gap or margin
  const moreTagWidth = 80; // reserve space for the "more" tag

  const [, forceUpdate] = useReducer((x) => x + 1, 0);

  const debounceResetCutoff = useDebounce(
    () => {
      setCutoff(undefined);
      // setting cutoff won't trigger a re-render if it's already undefined
      // so we force a re-render to recalculate the cutoff
      forceUpdate();
    },
    100 // Adjust the debounce delay as needed
  );

  // trigger recalculation of cutoff when control resizes
  useEffect(() => {
    if (!truncateOnOverflow || !ref.current) {
      return;
    }

    const resizeObserver = new ResizeObserver(() => {
      debounceResetCutoff();
    });

    const { current } = ref;
    resizeObserver.observe(current);

    return () => {
      resizeObserver.disconnect();
    };
  }, [truncateOnOverflow, debounceResetCutoff]);

  // we need to check this on every render, and the call to setCutoff _should_ be safe
  /* eslint-disable-next-line react-hooks/exhaustive-deps */
  useLayoutEffect(() => {
    if (!truncateOnOverflow) {
      setCutoff(undefined);
      return;
    }

    const { current } = ref;

    if (current) {
      // calculate the number of tags that can fit in the container
      const containerWidth = current.clientWidth;
      const children = Array.from(current.children);

      // don't recalculate anything if the more tag is visible and cutoff is already set
      const moreTags = children.find((child) => {
        return (child as HTMLElement).classList.contains("more-tags");
      });

      if (moreTags && !!cutoff) {
        return;
      }

      const childTags = children.filter((child) => {
        return (
          (child as HTMLElement).classList.contains("tag-item") ||
          (child as HTMLElement).classList.contains("clear-all-button")
        );
      });

      const clearAllButton = children.find((child) => {
        return (child as HTMLElement).classList.contains("clear-all-button");
      });

      // calculate the total width without the more tag
      const defaultTotalWidth = childTags.reduce((total, child, idx) => {
        return (
          total +
          ((child as HTMLElement).offsetWidth ?? 0) +
          (idx === childTags.length - 1 ? 0 : elementGap)
        );
      }, 0);

      if (containerWidth >= defaultTotalWidth) {
        // if the container is wide enough to fit all tags, reset cutoff
        setCutoff(undefined);
        return;
      }

      let totalWidth = 0;
      let visibleCount = 0;

      // reserve space for the more tags control
      totalWidth += moreTagWidth;

      // reserve space for the clear all button if present
      if (clearAllButton) {
        totalWidth += (clearAllButton as HTMLElement).offsetWidth ?? 0;
      }

      for (const child of children) {
        totalWidth += ((child as HTMLElement).offsetWidth ?? 0) + elementGap;
        if (totalWidth > containerWidth) {
          break;
        }
        visibleCount++;
      }

      setCutoff(visibleCount);
    }
  });

  function onRemoveCriterionTag(
    criterion: Criterion,
    $event: React.MouseEvent<HTMLElement, MouseEvent>,
    valueIndex?: number
  ) {
    if (!criterion) {
      return;
    }
    onRemoveCriterion(criterion, valueIndex);
    $event.stopPropagation();
  }

  function onClickCriterionTag(criterion: Criterion) {
    onEditCriterion(criterion);
  }

  function getFilterTags(criterion: Criterion) {
    if (
      criterion instanceof CustomFieldsCriterion &&
      criterion.value.length > 1
    ) {
      return criterion.value.map((value, index) => {
        return (
          <FilterTag
            key={index}
            label={criterion.getValueLabel(intl, value)}
            onClick={() => onClickCriterionTag(criterion)}
            onRemove={($event) =>
              onRemoveCriterionTag(criterion, $event, index)
            }
          />
        );
      });
    }

    return (
      <FilterTag
        key={criterion.getId()}
        label={criterion.getLabel(intl)}
        onClick={() => onClickCriterionTag(criterion)}
        onRemove={($event) => onRemoveCriterionTag(criterion, $event)}
      />
    );
  }

  if (criteria.length === 0) {
    return null;
  }

  const className = "wrap-tags filter-tags";

  const filterTags = criteria.map((c) => getFilterTags(c)).flat();

  if (cutoff && filterTags.length > cutoff) {
    const visibleCriteria = filterTags.slice(0, cutoff);
    const hiddenCriteria = filterTags.slice(cutoff);

    return (
      <div className={className} ref={ref}>
        {visibleCriteria}
        <MoreFilterTags tags={hiddenCriteria} />
        {criteria.length >= 3 && (
          <Button
            variant="minimal"
            className="clear-all-button"
            onClick={() => onRemoveAll()}
          >
            <FormattedMessage id="actions.clear" />
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className={className} ref={ref}>
      {filterTags}
      {criteria.length >= 3 && (
        <Button
          variant="minimal"
          className="clear-all-button"
          onClick={() => onRemoveAll()}
        >
          <FormattedMessage id="actions.clear" />
        </Button>
      )}
    </div>
  );
};
