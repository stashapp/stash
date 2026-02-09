import React, { useCallback, useMemo, useState } from "react";
import { useIntl } from "react-intl";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SelectedList } from "./SidebarListFilter";
import {
  DuplicatedCriterion,
  DuplicatedCriterionOption,
  DuplicationFieldId,
  DUPLICATION_FIELD_IDS,
  DUPLICATION_FIELD_MESSAGE_IDS,
} from "src/models/list-filter/criteria/phash";
import { IndeterminateCheckbox } from "src/components/Shared/IndeterminateCheckbox";
import { SidebarSection } from "src/components/Shared/Sidebar";
import { Icon } from "src/components/Shared/Icon";
import { faPlus } from "@fortawesome/free-solid-svg-icons";
import { keyboardClickHandler } from "src/utils/keyboard";

interface IDuplicatedFilter {
  criterion: DuplicatedCriterion;
  setCriterion: (c: DuplicatedCriterion) => void;
}

export const DuplicatedFilter: React.FC<IDuplicatedFilter> = ({
  criterion,
  setCriterion,
}) => {
  const intl = useIntl();

  function onFieldChange(
    fieldId: DuplicationFieldId,
    value: boolean | undefined
  ) {
    const c = criterion.clone();
    if (value === undefined) {
      delete c.value[fieldId];
    } else {
      c.value[fieldId] = value;
    }
    setCriterion(c);
  }

  return (
    <div className="duplicated-filter">
      {DUPLICATION_FIELD_IDS.map((fieldId) => (
        <IndeterminateCheckbox
          key={fieldId}
          id={`duplicated-${fieldId}`}
          label={intl.formatMessage({
            id: DUPLICATION_FIELD_MESSAGE_IDS[fieldId],
          })}
          checked={criterion.value[fieldId]}
          setChecked={(v) => onFieldChange(fieldId, v)}
        />
      ))}
    </div>
  );
};

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

  // Get label for a duplicate type
  const getLabel = useCallback(
    (typeId: DuplicationFieldId) =>
      intl.formatMessage({ id: DUPLICATION_FIELD_MESSAGE_IDS[typeId] }),
    [intl]
  );

  // Get the single duplicated criterion from the filter
  const getCriterion = useCallback((): DuplicatedCriterion | null => {
    const criteria = filter.criteriaFor(
      DuplicatedCriterionOption.type
    ) as DuplicatedCriterion[];
    return criteria.length > 0 ? criteria[0] : null;
  }, [filter]);

  // Get value for a specific type from the criterion
  const getTypeValue = useCallback(
    (typeId: DuplicationFieldId): boolean | undefined => {
      const criterion = getCriterion();
      if (!criterion) return undefined;
      return criterion.value[typeId];
    },
    [getCriterion]
  );

  // Build selected items list
  const selected: Option[] = useMemo(() => {
    const result: Option[] = [];
    const criterion = getCriterion();
    if (!criterion) return result;

    for (const typeId of DUPLICATION_FIELD_IDS) {
      const value = criterion.value[typeId];
      if (value !== undefined) {
        const valueLabel = value ? trueLabel : falseLabel;
        result.push({
          id: typeId,
          label: `${getLabel(typeId)}: ${valueLabel}`,
        });
      }
    }

    return result;
  }, [getCriterion, trueLabel, falseLabel, getLabel]);

  // Available options - show options that aren't already selected
  const options = useMemo(() => {
    const result: { id: DuplicationFieldId; label: string }[] = [];

    for (const typeId of DUPLICATION_FIELD_IDS) {
      if (getTypeValue(typeId) === undefined) {
        result.push({ id: typeId, label: getLabel(typeId) });
      }
    }

    return result;
  }, [getTypeValue, getLabel]);

  function onToggleExpand(id: string) {
    setExpandedType(expandedType === id ? null : id);
  }

  function onUnselect(item: Option) {
    const typeId = item.id as DuplicationFieldId;
    const criterion = getCriterion();

    if (!criterion) return;

    const newCriterion = criterion.clone();
    delete newCriterion.value[typeId];

    // If no fields are set, remove the criterion entirely
    const hasAnyValue = DUPLICATION_FIELD_IDS.some(
      (id) => newCriterion.value[id] !== undefined
    );

    if (!hasAnyValue) {
      setFilter(filter.removeCriterion(DuplicatedCriterionOption.type));
    } else {
      setFilter(
        filter.replaceCriteria(DuplicatedCriterionOption.type, [newCriterion])
      );
    }
    setExpandedType(null);
  }

  function onSelectValue(typeId: string, value: boolean) {
    const criterion = getCriterion();
    const newCriterion = criterion
      ? criterion.clone()
      : (DuplicatedCriterionOption.makeCriterion() as DuplicatedCriterion);

    newCriterion.value[typeId as DuplicationFieldId] = value;
    setFilter(
      filter.replaceCriteria(DuplicatedCriterionOption.type, [newCriterion])
    );
    setExpandedType(null);
  }

  return (
    <SidebarSection
      className="sidebar-list-filter"
      text={title}
      sectionID={sectionID}
      outsideCollapse={
        <SelectedList items={selected} onUnselect={(i) => onUnselect(i)} />
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
                    <Icon
                      className="fa-fw include-button single-value"
                      icon={faPlus}
                    />
                    <span className="unselected-object-label">{opt.label}</span>
                  </div>
                </a>
              </li>
              {expandedType === opt.id && (
                <div className="duplicate-sub-options">
                  <div
                    className="duplicate-sub-option"
                    onClick={() => onSelectValue(opt.id, true)}
                  >
                    {trueLabel}
                  </div>
                  <div
                    className="duplicate-sub-option"
                    onClick={() => onSelectValue(opt.id, false)}
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
