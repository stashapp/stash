import React, {
  PropsWithChildren,
  useEffect,
  useLayoutEffect,
  useReducer,
  useRef,
} from "react";
import { Badge, BadgeProps, Button, Overlay, Popover } from "react-bootstrap";
import {
  Criterion,
  IHierarchicalLabeledIdCriterion,
  ILabeledIdCriterion,
  ModifierCriterion,
} from "src/models/list-filter/criteria/criterion";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import {
  faMagnifyingGlass,
  faTimes,
  faTags,
  faUser,
  faVideo,
  faImage,
  faBuilding,
  faLayerGroup,
  faCalendar,
  faStar,
  faGlobe,
  faVenusMars,
  faFolder,
  faFingerprint,
  faLink,
  faFilm,
  faClock,
  faHashtag,
  faFont,
  faToggleOn,
  faRuler,
  faClosedCaptioning,
  faQuestion,
} from "@fortawesome/free-solid-svg-icons";
import { BsPrefixProps, ReplaceProps } from "react-bootstrap/esm/helpers";
import { CustomFieldsCriterion } from "src/models/list-filter/criteria/custom-fields";
import { useDebounce } from "src/hooks/debounce";
import cx from "classnames";
import { CriterionModifier } from "src/core/generated-graphql";
import {
  criterionIsHierarchicalLabelValue,
  IHierarchicalLabelValue,
  ILabeledId,
  ILabeledValueListValue,
} from "src/models/list-filter/types";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";

// Map criterion types to icons
const criterionTypeIcons: Record<string, IconDefinition> = {
  // Entity filters
  tags: faTags,
  performers: faUser,
  studios: faBuilding,
  parent_studios: faBuilding,
  groups: faLayerGroup,
  movies: faFilm,
  galleries: faImage,
  scenes: faVideo,

  // Hierarchical entity filters
  parent_tags: faTags,
  child_tags: faTags,
  performer_tags: faTags,
  scene_tags: faTags,

  // Date/time filters
  date: faCalendar,
  birthdate: faCalendar,
  death_date: faCalendar,
  created_at: faCalendar,
  updated_at: faCalendar,
  last_played_at: faCalendar,
  scene_date: faCalendar,

  // Rating
  rating: faStar,
  rating100: faStar,

  // Location/path
  country: faGlobe,
  path: faFolder,
  url: faLink,

  // Identity
  gender: faVenusMars,
  phash: faFingerprint,
  stash_id: faLink,

  // Media properties
  resolution: faRuler,
  duration: faClock,
  captions: faClosedCaptioning,
  orientation: faRuler,

  // Counts
  scene_count: faHashtag,
  image_count: faHashtag,
  gallery_count: faHashtag,
  performer_count: faHashtag,
  tag_count: faHashtag,
  o_counter: faHashtag,
  play_count: faHashtag,
  file_count: faHashtag,

  // Boolean
  organized: faToggleOn,
  favorite: faStar,
  performer_favorite: faStar,
  interactive: faToggleOn,
  hasMarkers: faToggleOn,
  hasChapters: faToggleOn,
  is_missing: faToggleOn,
};

// Get icon for a criterion type
function getCriterionIcon(type: string): IconDefinition {
  // Check for exact match first
  if (criterionTypeIcons[type]) {
    return criterionTypeIcons[type];
  }

  // Check for partial matches
  if (type.includes("tag")) return faTags;
  if (type.includes("performer")) return faUser;
  if (type.includes("studio")) return faBuilding;
  if (type.includes("group")) return faLayerGroup;
  if (type.includes("gallery")) return faImage;
  if (type.includes("scene")) return faVideo;
  if (type.includes("date") || type.includes("_at")) return faCalendar;
  if (type.includes("rating")) return faStar;
  if (type.includes("count") || type.includes("counter")) return faHashtag;
  if (type.includes("path") || type.includes("folder")) return faFolder;

  // Default icon
  return faQuestion;
}

// Get modifier class name
function getModifierClassName(modifier: CriterionModifier): string {
  switch (modifier) {
    case CriterionModifier.Includes:
    case CriterionModifier.IncludesAll:
    case CriterionModifier.Equals:
    case CriterionModifier.GreaterThan:
    case CriterionModifier.LessThan:
    case CriterionModifier.Between:
    case CriterionModifier.MatchesRegex:
      return "modifier-includes";
    case CriterionModifier.Excludes:
    case CriterionModifier.NotEquals:
    case CriterionModifier.NotBetween:
    case CriterionModifier.NotMatchesRegex:
      return "modifier-excludes";
    case CriterionModifier.IsNull:
      return "modifier-null";
    case CriterionModifier.NotNull:
      return "modifier-not-null";
    default:
      return "";
  }
}

// Get compact modifier label
function getCompactModifierLabel(
  modifier: CriterionModifier,
  intl: ReturnType<typeof useIntl>
): string | null {
  switch (modifier) {
    case CriterionModifier.Includes:
    case CriterionModifier.IncludesAll:
      return null; // Don't show for includes (default)
    case CriterionModifier.Excludes:
      return intl.formatMessage({ id: "criterion_modifier.not" });
    case CriterionModifier.Equals:
      return "=";
    case CriterionModifier.NotEquals:
      return "≠";
    case CriterionModifier.GreaterThan:
      return ">";
    case CriterionModifier.LessThan:
      return "<";
    case CriterionModifier.Between:
      return "↔";
    case CriterionModifier.NotBetween:
      return "!↔";
    case CriterionModifier.IsNull:
      return intl.formatMessage({ id: "criterion_modifier.is_null_short", defaultMessage: "none" });
    case CriterionModifier.NotNull:
      return intl.formatMessage({ id: "criterion_modifier.not_null_short", defaultMessage: "any" });
    case CriterionModifier.MatchesRegex:
      return "~";
    case CriterionModifier.NotMatchesRegex:
      return "!~";
    default:
      return null;
  }
}

type TagItemProps = PropsWithChildren<
  ReplaceProps<"span", BsPrefixProps<"span"> & BadgeProps>
>;

export const TagItem: React.FC<TagItemProps> = (props) => {
  const { className, children, ...others } = props;
  return (
    <Badge
      className={cx("tag-item", className)}
      variant="secondary"
      {...others}
    >
      {children}
    </Badge>
  );
};

export const FilterTag: React.FC<{
  className?: string;
  label: React.ReactNode;
  onClick: React.MouseEventHandler<HTMLSpanElement>;
  onRemove: React.MouseEventHandler<HTMLElement>;
  icon?: IconDefinition;
  modifierLabel?: string | null;
}> = ({ className, label, onClick, onRemove, icon, modifierLabel }) => {
  return (
    <TagItem className={className} onClick={onClick}>
      {icon && <Icon icon={icon} className="filter-tag-icon" />}
      {modifierLabel && <span className="filter-tag-modifier">{modifierLabel}</span>}
      <span className="filter-tag-label">{label}</span>
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

// Depth indicator component
const DepthIndicator: React.FC<{ depth: number }> = ({ depth }) => {
  if (depth === 0) return null;

  const label = depth === -1 ? "∞" : `+${depth}`;
  const title = depth === -1 ? "All sub-items" : `${depth} level${depth > 1 ? "s" : ""} deep`;

  return (
    <span className="depth-indicator" title={title}>
      {label}
    </span>
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
          className="hover-popover-content more-filter-tags-popover"
          onMouseEnter={handleMouseEnter}
          onMouseLeave={handleMouseLeave}
          onClick={handleMouseLeave}
        >
          <div className="more-filter-tags-content">{tags}</div>
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
  searchTerm?: string;
  criteria: Criterion[];
  onEditSearchTerm?: () => void;
  onEditCriterion: (c: Criterion) => void;
  onRemoveCriterion: (c: Criterion, valueIndex?: number) => void;
  onRemoveAll: () => void;
  onRemoveSearchTerm?: () => void;
  truncateOnOverflow?: boolean;
  splitMultiValue?: boolean;
}

export const FilterTags: React.FC<IFilterTagsProps> = ({
  searchTerm,
  criteria,
  onEditCriterion,
  onRemoveCriterion,
  onRemoveAll,
  onEditSearchTerm,
  onRemoveSearchTerm,
  truncateOnOverflow = false,
  splitMultiValue = true,
}) => {
  const intl = useIntl();
  const ref = useRef<HTMLDivElement>(null);

  const [cutoff, setCutoff] = React.useState<number | undefined>();
  const elementGap = 10;
  const moreTagWidth = 80;

  const [, forceUpdate] = useReducer((x) => x + 1, 0);

  const debounceResetCutoff = useDebounce(
    () => {
      setCutoff(undefined);
      forceUpdate();
    },
    100
  );

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

  /* eslint-disable-next-line react-hooks/exhaustive-deps */
  useLayoutEffect(() => {
    if (!truncateOnOverflow) {
      setCutoff(undefined);
      return;
    }

    const { current } = ref;

    if (current) {
      const containerWidth = current.clientWidth;
      const children = Array.from(current.children);

      const moreTags = children.find((child) => {
        return (child as HTMLElement).classList.contains("more-tags");
      });

      if (moreTags && cutoff !== undefined) {
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

      const defaultTotalWidth = childTags.reduce((total, child, idx) => {
        return (
          total +
          ((child as HTMLElement).offsetWidth ?? 0) +
          (idx === childTags.length - 1 ? 0 : elementGap)
        );
      }, 0);

      if (containerWidth >= defaultTotalWidth) {
        setCutoff(undefined);
        return;
      }

      let totalWidth = 0;
      let visibleCount = 0;

      totalWidth += moreTagWidth;

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

  // Create individual tag for a single value in a multi-value criterion
  function createValueTag(
    criterion: Criterion,
    item: ILabeledId,
    isExcluded: boolean,
    depth?: number,
    valueIndex?: number
  ) {
    const type = criterion.criterionOption.type;
    const icon = getCriterionIcon(type);
    const modifier = isExcluded ? CriterionModifier.Excludes : CriterionModifier.Includes;
    const modifierClass = getModifierClassName(modifier);
    const modifierLabel = getCompactModifierLabel(modifier, intl);

    return (
      <FilterTag
        key={`${criterion.getId()}-${item.id}-${isExcluded ? "ex" : "in"}`}
        className={cx("filter-tag-value", modifierClass)}
        icon={icon}
        modifierLabel={modifierLabel}
        label={
          <span className="filter-tag-value-content">
            {item.label}
            {depth !== undefined && <DepthIndicator depth={depth} />}
          </span>
        }
        onClick={() => onClickCriterionTag(criterion)}
        onRemove={($event) => onRemoveCriterionTag(criterion, $event, valueIndex)}
      />
    );
  }

  // Get filter tags for a hierarchical criterion (tags, performers, etc.)
  function getHierarchicalFilterTags(criterion: IHierarchicalLabeledIdCriterion) {
    const value = criterion.value as IHierarchicalLabelValue;
    const tags: React.ReactNode[] = [];

    // Handle IsNull/NotNull modifiers
    if (
      criterion.modifier === CriterionModifier.IsNull ||
      criterion.modifier === CriterionModifier.NotNull
    ) {
      return getDefaultFilterTag(criterion);
    }

    if (splitMultiValue) {
      // Create individual tags for included items
      value.items.forEach((item, index) => {
        tags.push(createValueTag(criterion, item, false, value.depth, index));
      });

      // Create individual tags for excluded items
      value.excluded.forEach((item, index) => {
        tags.push(createValueTag(criterion, item, true, undefined, index));
      });
    } else {
      // Single tag with all values
      return getDefaultFilterTag(criterion);
    }

    return tags;
  }

  // Get filter tags for a labeled ID list criterion
  function getLabeledIdListFilterTags(criterion: ILabeledIdCriterion) {
    const value = criterion.value as ILabeledValueListValue;
    const tags: React.ReactNode[] = [];

    // Handle IsNull/NotNull modifiers
    if (
      criterion.modifier === CriterionModifier.IsNull ||
      criterion.modifier === CriterionModifier.NotNull
    ) {
      return getDefaultFilterTag(criterion);
    }

    if (splitMultiValue && value.items) {
      // Create individual tags for included items
      value.items.forEach((item, index) => {
        tags.push(createValueTag(criterion, item, false, undefined, index));
      });

      // Create individual tags for excluded items
      if (value.excluded) {
        value.excluded.forEach((item, index) => {
          tags.push(createValueTag(criterion, item, true, undefined, index));
        });
      }
    } else {
      return getDefaultFilterTag(criterion);
    }

    return tags;
  }

  // Default filter tag (for non-split or simple criteria)
  function getDefaultFilterTag(criterion: Criterion) {
    const type = criterion.criterionOption.type;
    const icon = getCriterionIcon(type);

    let modifierClass = "";
    let modifierLabel: string | null = null;

    if (criterion instanceof ModifierCriterion) {
      modifierClass = getModifierClassName(criterion.modifier);
      modifierLabel = getCompactModifierLabel(criterion.modifier, intl);
    }

    // Get the value label without the criterion name
    let label: React.ReactNode;

    if (criterion instanceof ModifierCriterion) {
      if (
        criterion.modifier === CriterionModifier.IsNull ||
        criterion.modifier === CriterionModifier.NotNull
      ) {
        // For null modifiers, just show the criterion name with modifier
        label = intl.formatMessage({ id: criterion.criterionOption.messageID });
      } else if (criterionIsHierarchicalLabelValue(criterion.value)) {
        const hierValue = criterion.value as IHierarchicalLabelValue;
        const items = hierValue.items.map((v) => v.label).join(", ");
        const excluded = hierValue.excluded.map((v) => v.label).join(", ");
        label = (
          <span className="filter-tag-value-content">
            {items}
            {excluded && <span className="excluded-values"> (-{excluded})</span>}
            <DepthIndicator depth={hierValue.depth} />
          </span>
        );
      } else {
        // Get just the value portion
        label = criterion.getLabel(intl).replace(
          new RegExp(`^${intl.formatMessage({ id: criterion.criterionOption.messageID })}\\s*`),
          ""
        );
      }
    } else {
      label = criterion.getLabel(intl);
    }

    return (
      <FilterTag
        key={criterion.getId()}
        className={cx("filter-tag-default", modifierClass)}
        icon={icon}
        modifierLabel={modifierLabel}
        label={label}
        onClick={() => onClickCriterionTag(criterion)}
        onRemove={($event) => onRemoveCriterionTag(criterion, $event)}
      />
    );
  }

  function getFilterTags(criterion: Criterion) {
    // Handle custom fields specially
    if (
      criterion instanceof CustomFieldsCriterion &&
      criterion.value.length > 1
    ) {
      return criterion.value.map((value, index) => {
        return (
          <FilterTag
            key={index}
            icon={faFont}
            label={criterion.getValueLabel(intl, value)}
            onClick={() => onClickCriterionTag(criterion)}
            onRemove={($event) =>
              onRemoveCriterionTag(criterion, $event, index)
            }
          />
        );
      });
    }

    // Check if this is a hierarchical labeled ID criterion
    if (criterion instanceof ModifierCriterion) {
      const value = criterion.value;
      
      // Check for hierarchical value (tags, performers, studios with depth)
      if (criterionIsHierarchicalLabelValue(value)) {
        return getHierarchicalFilterTags(criterion as IHierarchicalLabeledIdCriterion);
      }

      // Check for labeled ID list value (without depth)
      if (
        typeof value === "object" &&
        value !== null &&
        "items" in value &&
        Array.isArray((value as ILabeledValueListValue).items)
      ) {
        return getLabeledIdListFilterTags(criterion as ILabeledIdCriterion);
      }
    }

    // Default handling
    return getDefaultFilterTag(criterion);
  }

  if (criteria.length === 0 && !searchTerm) {
    return null;
  }

  const className = "wrap-tags filter-tags";

  const filterTags = criteria.map((c) => getFilterTags(c)).flat();

  if (searchTerm && searchTerm.length > 0) {
    filterTags.unshift(
      <FilterTag
        key="search-term"
        className="search-term-filter-tag"
        icon={faMagnifyingGlass}
        label={searchTerm}
        onClick={() => onEditSearchTerm?.()}
        onRemove={() => onRemoveSearchTerm?.()}
      />
    );
  }

  const visibleCriteria =
    cutoff !== undefined ? filterTags.slice(0, cutoff) : filterTags;
  const hiddenCriteria = cutoff !== undefined ? filterTags.slice(cutoff) : [];

  return (
    <div className={className} ref={ref}>
      {visibleCriteria}
      <MoreFilterTags tags={hiddenCriteria} />
      {filterTags.length >= 3 && (
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
