import React, { useCallback, useMemo, useState } from "react";
import { useIntl } from "react-intl";
import { CriterionOption } from "../../../models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { IsMissingCriterion } from "src/models/list-filter/criteria/is-missing";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import {
  faQuestion,
  faImage,
  faLink,
  faCalendar,
  faBuilding,
  faUser,
  faTags,
  faLayerGroup,
  faFileAlt,
  faImages,
  faFingerprint,
  faFont,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";

// Map options to icons for better visual identification
const optionIcons: Record<string, IconDefinition> = {
  title: faFont,
  cover: faImage,
  details: faFileAlt,
  url: faLink,
  date: faCalendar,
  galleries: faImages,
  studio: faBuilding,
  group: faLayerGroup,
  performers: faUser,
  tags: faTags,
  stash_id: faFingerprint,
  // Gallery-specific
  scenes: faLayerGroup,
  // Performer-specific
  ethnicity: faUser,
  country: faUser,
  hair_color: faUser,
  eye_color: faUser,
  height: faUser,
  weight: faUser,
  measurements: faUser,
  fake_tits: faUser,
  career_length: faUser,
  tattoos: faUser,
  piercings: faUser,
  aliases: faUser,
  gender: faUser,
  image: faImage,
  // Studio-specific (image already covered)
  // Group-specific
  front_image: faImage,
  back_image: faImage,
};

interface ISidebarIsMissingFilterProps {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

export const SidebarIsMissingFilter: React.FC<ISidebarIsMissingFilterProps> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const intl = useIntl();
  const [query, setQuery] = useState("");

  // Get the available options from the criterion option
  const availableOptions = useMemo(() => {
    return option.options ?? [];
  }, [option]);

  // Create option objects with icons
  const options: Option[] = useMemo(() => {
    return availableOptions.map((opt) => {
      const optValue = typeof opt === "string" ? opt : opt.value;
      const optLabel = typeof opt === "string" ? opt : opt.messageID;
      const icon = optionIcons[optValue] ?? faQuestion;
      
      return {
        id: optValue,
        label: intl.formatMessage({ id: optLabel, defaultMessage: optLabel }),
        icon: <FontAwesomeIcon icon={icon} className="fa-fw" />,
      };
    });
  }, [availableOptions, intl]);

  // Get current criterion
  const criteria = filter.criteriaFor(option.type) as IsMissingCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  // Selected item based on current criterion value
  const selected: Option[] = useMemo(() => {
    if (!criterion || !criterion.value) return [];
    
    const selectedOption = options.find((o) => o.id === criterion.value);
    return selectedOption ? [selectedOption] : [];
  }, [criterion, options]);

  // Filter candidates based on search query and exclude selected
  const candidates: Option[] = useMemo(() => {
    const lowerQuery = query.toLowerCase();
    return options
      .filter((o) => !selected.find((s) => s.id === o.id))
      .filter((o) => o.label.toLowerCase().includes(lowerQuery));
  }, [options, selected, query]);

  const onSelect = useCallback(
    (item: Option) => {
      const newCriterion = option.makeCriterion() as IsMissingCriterion;
      newCriterion.value = item.id;
      setFilter(filter.replaceCriteria(option.type, [newCriterion]));
    },
    [filter, setFilter, option]
  );

  const onUnselect = useCallback(
    (_item: Option) => {
      setFilter(filter.removeCriterion(option.type));
    },
    [filter, setFilter, option]
  );

  return (
    <SidebarListFilter
      title={title}
      candidates={candidates}
      onSelect={onSelect}
      onUnselect={onUnselect}
      selected={selected}
      singleValue={true}
      sectionID={sectionID}
      query={query}
      setQuery={setQuery}
    />
  );
};

