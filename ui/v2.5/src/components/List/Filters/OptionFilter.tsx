import cloneDeep from "lodash-es/cloneDeep";
import React, { useMemo } from "react";
import { Form } from "react-bootstrap";
import {
  CriterionValue,
  ModifierCriterion,
  ModifierCriterionOption,
} from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";
import { CriterionModifier } from "src/core/generated-graphql";
import {
  getModifierCandidates,
  ModifierValue,
  modifierValueToModifier,
} from "./LabeledIdFilter";
import { useIntl } from "react-intl";

interface IOptionsFilter {
  criterion: ModifierCriterion<CriterionValue>;
  setCriterion: (c: ModifierCriterion<CriterionValue>) => void;
}

export const OptionFilter: React.FC<IOptionsFilter> = ({
  criterion,
  setCriterion,
}) => {
  function onSelect(v: string) {
    const c = cloneDeep(criterion);
    if (c.value === v) {
      c.value = "";
    } else {
      c.value = v;
    }

    setCriterion(c);
  }

  const { options } = criterion.modifierCriterionOption();

  return (
    <div className="option-list-filter">
      {options?.map((o) => (
        <Form.Check
          id={`${criterion.getId()}-${o.toString()}`}
          key={o.toString()}
          onChange={() => onSelect(o.toString())}
          checked={criterion.value === o.toString()}
          type="radio"
          label={o.toString()}
        />
      ))}
    </div>
  );
};

interface IOptionsListFilter {
  criterion: ModifierCriterion<CriterionValue>;
  setCriterion: (c: ModifierCriterion<CriterionValue>) => void;
}

export const OptionListFilter: React.FC<IOptionsListFilter> = ({
  criterion,
  setCriterion,
}) => {
  function onSelect(v: string) {
    const c = cloneDeep(criterion);
    const cv = c.value as string[];
    if (cv.includes(v)) {
      c.value = cv.filter((x) => x !== v);
    } else {
      c.value = [...cv, v];
    }

    setCriterion(c);
  }

  const { options } = criterion.modifierCriterionOption();
  const value = criterion.value as string[];

  return (
    <div className="option-list-filter">
      {options?.map((o) => (
        <Form.Check
          id={`${criterion.getId()}-${o.toString()}`}
          key={o.toString()}
          onChange={() => onSelect(o.toString())}
          checked={value.includes(o.toString())}
          type="checkbox"
          label={o.toString()}
        />
      ))}
    </div>
  );
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: ModifierCriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  sectionID?: string;
}

export const SidebarOptionFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  sectionID,
}) => {
  const intl = useIntl();

  const criteria = filter.criteriaFor(
    option.type
  ) as ModifierCriterion<CriterionValue>[];
  const criterion = criteria.length > 0 ? criteria[0] : null;
  const { options: criterionOptions = [] } = option;
  const currentValues = criteria.flatMap((c) => c.value as string[]);

  const hasNullModifiers =
    option.modifierOptions.includes(CriterionModifier.IsNull) &&
    option.modifierOptions.includes(CriterionModifier.NotNull);

  const selected: Option[] = useMemo(() => {
    if (!criterion) return [];

    if (criterion.modifier === CriterionModifier.IsNull) {
      return [
        {
          id: "none",
          label: intl.formatMessage({ id: "criterion_modifier_values.none" }),
        },
      ];
    } else if (criterion.modifier === CriterionModifier.NotNull) {
      return [
        {
          id: "any",
          label: intl.formatMessage({ id: "criterion_modifier_values.any" }),
        },
      ];
    }

    return criterionOptions
      .filter((o) => currentValues.includes(o.toString()))
      .map((o) => ({
        id: o.toString(),
        label: o.toLocaleString(),
      }));
  }, [criterion, currentValues, criterionOptions, intl]);

  const modifierCandidates: Option[] = useMemo(() => {
    if (!hasNullModifiers) return [];

    const c = getModifierCandidates({
      modifier: criterion?.modifier ?? option.defaultModifier,
      defaultModifier: option.defaultModifier,
      hasExcluded: false,
      hasSelected: selected.length > 0,
      singleValue: true, // so that it doesn't include any_of
    });

    return c.map((v) => {
      const messageID = `criterion_modifier_values.${v}`;

      return {
        id: v,
        label: `(${intl.formatMessage({
          id: messageID,
        })})`,
        className: "modifier-object",
        canExclude: false,
      };
    });
  }, [criterion, option, selected, hasNullModifiers, intl]);

  const options = useMemo(() => {
    const o = criterionOptions
      .filter((oo) => !currentValues.includes(oo.toString()))
      .map((oo) => ({
        id: oo.toString(),
        label: oo.toString(),
      }));

    return [...modifierCandidates, ...o];
  }, [criterionOptions, currentValues, modifierCandidates]);

  function onSelect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    if (item.className === "modifier-object") {
      newCriterion.modifier = modifierValueToModifier(item.id as ModifierValue);
      newCriterion.value = [];
      setFilter(filter.replaceCriteria(option.type, [newCriterion]));
      return;
    }

    const cv = newCriterion.value as string[];
    if (cv.includes(item.id)) {
      return;
    } else {
      newCriterion.value = [...cv, item.id];
    }

    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselect(item: Option) {
    if (item.className === "modifier-object") {
      const newCriterion = criterion
        ? criterion.clone()
        : option.makeCriterion();
      newCriterion.modifier = option.defaultModifier;
      setFilter(filter.replaceCriteria(option.type, [newCriterion]));
      return;
    }

    setFilter(filter.removeCriterion(option.type));
  }

  return (
    <>
      <SidebarListFilter
        title={title}
        candidates={options}
        onSelect={onSelect}
        onUnselect={onUnselect}
        selected={selected}
        singleValue
        sectionID={sectionID}
      />
    </>
  );
};
