import React, { useMemo, useState } from "react";
import { useIntl } from "react-intl";
import { BooleanCriterion } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { DuplicatedCriterionOption } from "src/models/list-filter/criteria/phash";

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

  const criteria = filter.criteriaFor(
    DuplicatedCriterionOption.type
  ) as BooleanCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  // The main duplicate type option
  const phashOption = useMemo(
    () => ({
      id: "phash",
      label: phashLabel,
    }),
    [phashLabel]
  );

  // Determine if pHash is selected (has a true/false value)
  const phashSelected = criterion !== null;

  // Selected shows "pHash: True" or "pHash: False" when a value is set
  const selected: Option[] = useMemo(() => {
    if (!criterion) return [];

    const valueLabel = criterion.value === "true" ? trueLabel : falseLabel;
    return [
      {
        id: "phash",
        label: `${phashLabel}: ${valueLabel}`,
      },
    ];
  }, [criterion, phashLabel, trueLabel, falseLabel]);

  // Available options - show pHash if not selected
  const options: Option[] = useMemo(() => {
    if (phashSelected) return [];
    return [phashOption];
  }, [phashSelected, phashOption]);

  function onSelect(item: Option) {
    if (item.id === "phash") {
      // Expand to show True/False options
      setExpandedType("phash");
    }
  }

  function onUnselect() {
    setFilter(filter.removeCriterion(DuplicatedCriterionOption.type));
    setExpandedType(null);
  }

  function onSelectValue(value: "true" | "false") {
    const newCriterion = criterion
      ? criterion.clone()
      : DuplicatedCriterionOption.makeCriterion();
    newCriterion.value = value;
    setFilter(
      filter.replaceCriteria(DuplicatedCriterionOption.type, [newCriterion])
    );
    setExpandedType(null);
  }

  // Sub-options shown when pHash is clicked
  const subOptions =
    expandedType === "phash" ? (
      <div className="duplicate-sub-options">
        <div
          className="duplicate-sub-option"
          onClick={() => onSelectValue("true")}
        >
          {trueLabel}
        </div>
        <div
          className="duplicate-sub-option"
          onClick={() => onSelectValue("false")}
        >
          {falseLabel}
        </div>
      </div>
    ) : null;

  return (
    <SidebarListFilter
      title={title}
      candidates={options}
      onSelect={onSelect}
      onUnselect={onUnselect}
      selected={selected}
      singleValue
      postCandidates={subOptions}
      sectionID={sectionID}
    />
  );
};
