import cloneDeep from "lodash-es/cloneDeep";
import React, { useContext, useMemo } from "react";
import { Form } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faHeart,
  faCheck,
  faXmark,
  faGamepad,
  faMapMarkerAlt,
  faClone,
  faBookmark,
  faRobot,
  faStar,
} from "@fortawesome/free-solid-svg-icons";
import { 
  faHeart as faHeartRegular,
  faBookmark as faBookmarkRegular,
  faStar as faStarRegular,
} from "@fortawesome/free-regular-svg-icons";
import {
  BooleanCriterion,
  CriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { FormattedMessage, useIntl } from "react-intl";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { FacetCountsContext } from "src/extensions/hooks/useFacetCounts";

interface IBooleanFilter {
  criterion: BooleanCriterion;
  setCriterion: (c: BooleanCriterion) => void;
}

export const BooleanFilter: React.FC<IBooleanFilter> = ({
  criterion,
  setCriterion,
}) => {
  function onSelect(v: boolean) {
    const c = cloneDeep(criterion);
    if ((v && c.value === "true") || (!v && c.value === "false")) {
      c.value = "";
    } else {
      c.value = v ? "true" : "false";
    }

    setCriterion(c);
  }

  return (
    <div className="boolean-filter">
      <Form.Check
        id={`${criterion.getId()}-true`}
        onChange={() => onSelect(true)}
        checked={criterion.value === "true"}
        type="radio"
        label={<FormattedMessage id="true" />}
      />
      <Form.Check
        id={`${criterion.getId()}-false`}
        onChange={() => onSelect(false)}
        checked={criterion.value === "false"}
        type="radio"
        label={<FormattedMessage id="false" />}
      />
    </div>
  );
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

export const SidebarBooleanFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const intl = useIntl();
  
  // Get facet counts from context
  const { counts: facetCounts } = useContext(FacetCountsContext);
  
  // Map criterion type to facet boolean type, labels, and icons
  const { facetBooleanKey, trueLabel, falseLabel, trueIcon, falseIcon } = useMemo(() => {
    switch (option.type) {
      case "organized":
        return {
          facetBooleanKey: "organized" as const,
          trueLabel: intl.formatMessage({ id: "organized" }),
          falseLabel: intl.formatMessage({ id: "not_organized", defaultMessage: "Not Organized" }),
          trueIcon: <FontAwesomeIcon icon={faCheck} style={{ color: "#6c9" }} fixedWidth />,
          falseIcon: <FontAwesomeIcon icon={faXmark} style={{ color: "#c66" }} fixedWidth />,
        };
      case "interactive":
        return {
          facetBooleanKey: "interactive" as const,
          trueLabel: intl.formatMessage({ id: "interactive" }),
          falseLabel: intl.formatMessage({ id: "not_interactive", defaultMessage: "Not Interactive" }),
          trueIcon: <FontAwesomeIcon icon={faGamepad} style={{ color: "#6cf" }} fixedWidth />,
          falseIcon: <FontAwesomeIcon icon={faGamepad} style={{ color: "#666", opacity: 0.5 }} fixedWidth />,
        };
      case "favorite":
        return {
          facetBooleanKey: "favorite" as const,
          trueLabel: intl.formatMessage({ id: "favourite" }),
          falseLabel: intl.formatMessage({ id: "not_favourite", defaultMessage: "Not Favourite" }),
          trueIcon: <FontAwesomeIcon icon={faHeart} style={{ color: "#e66" }} fixedWidth />,
          falseIcon: <FontAwesomeIcon icon={faHeartRegular} style={{ color: "#888" }} fixedWidth />,
        };
      case "has_markers":
        return {
          facetBooleanKey: null, // No facet support for has_markers
          trueLabel: intl.formatMessage({ id: "has_markers_true", defaultMessage: "Has Markers" }),
          falseLabel: intl.formatMessage({ id: "has_markers_false", defaultMessage: "No Markers" }),
          trueIcon: <FontAwesomeIcon icon={faMapMarkerAlt} style={{ color: "#f93" }} fixedWidth />,
          falseIcon: <FontAwesomeIcon icon={faMapMarkerAlt} style={{ color: "#666", opacity: 0.5 }} fixedWidth />,
        };
      case "duplicated":
        return {
          facetBooleanKey: null, // No facet support for duplicated
          trueLabel: intl.formatMessage({ id: "duplicated_phash_true", defaultMessage: "Duplicated" }),
          falseLabel: intl.formatMessage({ id: "duplicated_phash_false", defaultMessage: "Not Duplicated" }),
          trueIcon: <FontAwesomeIcon icon={faClone} style={{ color: "#c9f" }} fixedWidth />,
          falseIcon: <FontAwesomeIcon icon={faClone} style={{ color: "#666", opacity: 0.5 }} fixedWidth />,
        };
      case "filter_favorites": // For performers, tags, studios (favourite filter)
        return {
          facetBooleanKey: "favorite" as const,
          trueLabel: intl.formatMessage({ id: "favourite" }),
          falseLabel: intl.formatMessage({ id: "not_favourite", defaultMessage: "Not Favourite" }),
          trueIcon: <FontAwesomeIcon icon={faHeart} style={{ color: "#e66" }} fixedWidth />,
          falseIcon: <FontAwesomeIcon icon={faHeartRegular} style={{ color: "#888" }} fixedWidth />,
        };
      case "performer_favorite":
        return {
          facetBooleanKey: null, // No facet support
          trueLabel: intl.formatMessage({ id: "performer_favorite_true", defaultMessage: "Performer Favourite" }),
          falseLabel: intl.formatMessage({ id: "performer_favorite_false", defaultMessage: "Performer Not Favourite" }),
          trueIcon: <FontAwesomeIcon icon={faStar} style={{ color: "#fc6" }} fixedWidth />,
          falseIcon: <FontAwesomeIcon icon={faStarRegular} style={{ color: "#888" }} fixedWidth />,
        };
      case "has_chapters":
        return {
          facetBooleanKey: null, // No facet support
          trueLabel: intl.formatMessage({ id: "has_chapters_true", defaultMessage: "Has Chapters" }),
          falseLabel: intl.formatMessage({ id: "has_chapters_false", defaultMessage: "No Chapters" }),
          trueIcon: <FontAwesomeIcon icon={faBookmark} style={{ color: "#6c9" }} fixedWidth />,
          falseIcon: <FontAwesomeIcon icon={faBookmarkRegular} style={{ color: "#888" }} fixedWidth />,
        };
      case "ignore_auto_tag":
        return {
          facetBooleanKey: null, // No facet support
          trueLabel: intl.formatMessage({ id: "ignore_auto_tag_true", defaultMessage: "Ignored" }),
          falseLabel: intl.formatMessage({ id: "ignore_auto_tag_false", defaultMessage: "Not Ignored" }),
          trueIcon: <FontAwesomeIcon icon={faRobot} style={{ color: "#c66" }} fixedWidth />,
          falseIcon: <FontAwesomeIcon icon={faRobot} style={{ color: "#6c9" }} fixedWidth />,
        };
      default:
        return {
          facetBooleanKey: null,
          trueLabel: intl.formatMessage({ id: "true" }),
          falseLabel: intl.formatMessage({ id: "false" }),
          trueIcon: undefined,
          falseIcon: undefined,
        };
    }
  }, [option.type, intl]);

  const trueOption = useMemo(() => {
    const count = facetBooleanKey
      ? facetCounts.booleans[facetBooleanKey].true
      : undefined;
    return {
      id: "true",
      label: trueLabel,
      count,
      icon: trueIcon,
    };
  }, [trueLabel, facetCounts, facetBooleanKey, trueIcon]);

  const falseOption = useMemo(() => {
    const count = facetBooleanKey
      ? facetCounts.booleans[facetBooleanKey].false
      : undefined;
    return {
      id: "false",
      label: falseLabel,
      count,
      icon: falseIcon,
    };
  }, [falseLabel, facetCounts, facetBooleanKey, falseIcon]);

  const criteria = filter.criteriaFor(option.type) as BooleanCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  const selected: Option[] = useMemo(() => {
    if (!criterion) return [];

    if (criterion.value === "true") {
      return [trueOption];
    } else if (criterion.value === "false") {
      return [falseOption];
    }

    return [];
  }, [trueOption, falseOption, criterion]);

  const options: Option[] = useMemo(() => {
    // Boolean filters always show both options (with zero counts dimmed)
    return [trueOption, falseOption].filter((o) => !selected.includes(o));
  }, [selected, trueOption, falseOption]);

  function onSelect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();
    newCriterion.value = item.id;
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselect() {
    setFilter(filter.removeCriterion(option.type));
  }

  return (
    <SidebarListFilter
      title={title}
      candidates={options}
      onSelect={onSelect}
      onUnselect={onUnselect}
      selected={selected}
      singleValue
      sectionID={sectionID}
    />
  );
};
