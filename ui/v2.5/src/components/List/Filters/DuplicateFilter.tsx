import React, { useMemo, useState } from "react";
import { useIntl } from "react-intl";
import { BooleanCriterion } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SelectedList } from "./SidebarListFilter";
import { DuplicatedCriterionOption } from "src/models/list-filter/criteria/phash";
import { DuplicatedStashIDCriterionOption } from "src/models/list-filter/criteria/stash-ids";
import { DuplicatedTitleCriterionOption } from "src/models/list-filter/criteria/title";
import { SidebarSection } from "src/components/Shared/Sidebar";
import { Icon } from "src/components/Shared/Icon";
import { faPlus } from "@fortawesome/free-solid-svg-icons";
import { keyboardClickHandler } from "src/utils/keyboard";

// Mapping of duplicate type IDs to their criterion options
const DUPLICATE_TYPES = {
  phash: DuplicatedCriterionOption,
  stash_id: DuplicatedStashIDCriterionOption,
  title: DuplicatedTitleCriterionOption,
} as const;

type DuplicateTypeId = keyof typeof DUPLICATE_TYPES;

interface ISidebarDuplicateFilterProps {
  title?: React.ReactNode;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

export const SidebarDuplicateFilter: React.FC<ISidebarDuplicateFilterProps> = ({
  title,
  filter,
  setFilter,
  sectionID,
}) => {
  const intl = useIntl();
  const [expandedType, setExpandedType] = useState<string | null>(null);

  const trueLabel = intl.formatMessage({ id: "true" });
  const falseLabel = intl.formatMessage({ id: "false" });
  const phashLabel = intl.formatMessage({ id: "media_info.phash" });
  const stashIdLabel = intl.formatMessage({ id: "stash_id" });
  const titleLabel = intl.formatMessage({ id: "title" });

  // Labels for each duplicate type
  const labels: Record<DuplicateTypeId, string> = {
    phash: phashLabel,
    stash_id: stashIdLabel,
    title: titleLabel,
  };

  // Get criterion for a given type
  function getCriterion(typeId: DuplicateTypeId): BooleanCriterion | null {
    const criteria = filter.criteriaFor(
      DUPLICATE_TYPES[typeId].type
    ) as BooleanCriterion[];
    return criteria.length > 0 ? criteria[0] : null;
  }

  // Build selected items list
  const selected: Option[] = useMemo(() => {
    const result: Option[] = [];

    for (const typeId of Object.keys(DUPLICATE_TYPES) as DuplicateTypeId[]) {
      const criterion = getCriterion(typeId);
      if (criterion) {
        const valueLabel = criterion.value === "true" ? trueLabel : falseLabel;
        result.push({
          id: typeId,
          label: `${labels[typeId]}: ${valueLabel}`,
        });
      }
    }

    return result;
  }, [filter, trueLabel, falseLabel, labels]);

  // Available options - show options that aren't already selected
  const options = useMemo(() => {
    const result: { id: DuplicateTypeId; label: string }[] = [];

    for (const typeId of Object.keys(DUPLICATE_TYPES) as DuplicateTypeId[]) {
      if (!getCriterion(typeId)) {
        result.push({ id: typeId, label: labels[typeId] });
      }
    }

    return result;
  }, [filter, labels]);

  function onToggleExpand(id: string) {
    setExpandedType(expandedType === id ? null : id);
  }

  function onUnselect(item: Option) {
    const typeId = item.id as DuplicateTypeId;
    const criterionOption = DUPLICATE_TYPES[typeId];
    if (criterionOption) {
      setFilter(filter.removeCriterion(criterionOption.type));
    }
    setExpandedType(null);
  }

  function onSelectValue(typeId: string, value: "true" | "false") {
    const criterionOption = DUPLICATE_TYPES[typeId as DuplicateTypeId];
    if (!criterionOption) return;

    const existingCriterion = getCriterion(typeId as DuplicateTypeId);
    const newCriterion = existingCriterion
      ? existingCriterion.clone()
      : criterionOption.makeCriterion();
    newCriterion.value = value;
    setFilter(filter.replaceCriteria(criterionOption.type, [newCriterion]));
    setExpandedType(null);
  }

  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      sectionID={sectionID}
      outsideCollapse={
        <SelectedList
          items={selected}
          onUnselect={(i) => onUnselect(i)}
        />
      }
    >
      <div className="queryable-candidate-list">
        <ul>
          {options.map((opt) => (
            <React.Fragment key={opt.id}>
              <li className="unselected-object">
                <a
                  onClick={() => onToggleExpand(opt.id)}
                  onKeyDown={keyboardClickHandler(() => onToggleExpand(opt.id))}
                  tabIndex={0}
                >
                  <div className="label-group">
                    <Icon className="fa-fw include-button single-value" icon={faPlus} />
                    <span className="unselected-object-label">{opt.label}</span>
                  </div>
                </a>
              </li>
              {expandedType === opt.id && (
                <div className="duplicate-sub-options">
                  <div
                    className="duplicate-sub-option"
                    onClick={() => onSelectValue(opt.id, "true")}
                  >
                    {trueLabel}
                  </div>
                  <div
                    className="duplicate-sub-option"
                    onClick={() => onSelectValue(opt.id, "false")}
                  >
                    {falseLabel}
                  </div>
                </div>
              )}
            </React.Fragment>
          ))}
        </ul>
      </div>
    </SidebarSection>
  );
};