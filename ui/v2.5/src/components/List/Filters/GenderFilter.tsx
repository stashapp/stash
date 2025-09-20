import React, { useMemo } from "react";
import { useIntl } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import {
  CriterionOption,
  ModifierCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { GenderCriterion } from "src/models/list-filter/criteria/gender";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Option, SidebarListFilter } from "./SidebarListFilter";

interface IGenderFilterProps {
  criterion: ModifierCriterion<string[]>;
  onValueChanged: (value: string[]) => void;
}

export const GenderFilter: React.FC<IGenderFilterProps> = ({
  criterion,
  onValueChanged,
}) => {
  const intl = useIntl();

  const genderOptions = useMemo(() => {
    return [
      { id: "Male", label: intl.formatMessage({ id: "gender_types.MALE" }) },
      {
        id: "Female",
        label: intl.formatMessage({ id: "gender_types.FEMALE" }),
      },
      {
        id: "Transgender Male",
        label: intl.formatMessage({ id: "gender_types.TRANSGENDER_MALE" }),
      },
      {
        id: "Transgender Female",
        label: intl.formatMessage({ id: "gender_types.TRANSGENDER_FEMALE" }),
      },
      {
        id: "Intersex",
        label: intl.formatMessage({ id: "gender_types.INTERSEX" }),
      },
      {
        id: "Non-Binary",
        label: intl.formatMessage({ id: "gender_types.NON_BINARY" }),
      },
    ];
  }, [intl]);

  const selectedOptions = useMemo(() => {
    return genderOptions.filter((option) =>
      criterion.value.includes(option.id)
    );
  }, [genderOptions, criterion.value]);

  const availableOptions = useMemo(() => {
    return genderOptions.filter(
      (option) => !criterion.value.includes(option.id)
    );
  }, [genderOptions, criterion.value]);

  function onSelect(item: Option) {
    const newValue = [...criterion.value, item.id];
    onValueChanged(newValue);
  }

  function onUnselect(item: Option) {
    const newValue = criterion.value.filter((v) => v !== item.id);
    onValueChanged(newValue);
  }

  return (
    <div className="gender-filter">
      {/* <SidebarListFilter
        candidates={availableOptions}
        onSelect={onSelect}
        onUnselect={onUnselect}
        selected={selectedOptions}
        singleValue={false}
      /> */}
    </div>
  );
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
}

export const SidebarGenderFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
}) => {
  const intl = useIntl();

  const options = useMemo(() => {
    return [
      { id: "Male", label: intl.formatMessage({ id: "gender_types.MALE" }) },
      {
        id: "Female",
        label: intl.formatMessage({ id: "gender_types.FEMALE" }),
      },
      {
        id: "Transgender Male",
        label: intl.formatMessage({ id: "gender_types.TRANSGENDER_MALE" }),
      },
      {
        id: "Transgender Female",
        label: intl.formatMessage({ id: "gender_types.TRANSGENDER_FEMALE" }),
      },
      {
        id: "Intersex",
        label: intl.formatMessage({ id: "gender_types.INTERSEX" }),
      },
      {
        id: "Non-Binary",
        label: intl.formatMessage({ id: "gender_types.NON_BINARY" }),
      },
    ];
  }, [intl]);

  const criteria = filter.criteriaFor(option.type) as GenderCriterion[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  const selected: Option[] = useMemo(() => {
    if (!criterion) return [];

    if (
      criterion.modifier === CriterionModifier.Includes ||
      criterion.modifier === CriterionModifier.Excludes
    ) {
      return options.filter((option) => criterion.value.includes(option.id));
    }

    return [];
  }, [options, criterion]);

  function onSelect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    if (criterion && criterion.modifier === CriterionModifier.Includes) {
      const currentValues = criterion.value;
      if (!currentValues.includes(item.id)) {
        // Add to selection
        newCriterion.value = [...currentValues, item.id];
      }
    } else {
      // Start new selection
      newCriterion.modifier = CriterionModifier.Includes;
      newCriterion.value = [item.id];
    }
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  function onUnselect(item: Option) {
    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();

    if (criterion && criterion.modifier === CriterionModifier.Includes) {
      const currentValues = criterion.value;
      if (currentValues.includes(item.id)) {
        // Remove from selection
        newCriterion.value = currentValues.filter((v) => v !== item.id);
        if (newCriterion.value.length === 0) {
          setFilter(filter.removeCriterion(option.type));
          return;
        }
      }
    }
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  // handle filtering of selected options
  const candidates = useMemo(() => {
    return options.filter(
      (p) => selected.find((s) => s.id === p.id) === undefined
    );
  }, [options, selected]);

  return (
    <>
      <SidebarListFilter
        title={title}
        candidates={candidates}
        onSelect={onSelect}
        onUnselect={onUnselect}
        selected={selected}
        singleValue={false}
      />
    </>
  );
};

// export function useMyLabeledIdFilterState(props: {
//   option: CriterionOption;
//   filter: ListFilterModel;
//   setFilter: (f: ListFilterModel) => void;
//   singleValue?: boolean;
//   hierarchical?: boolean;
//   includeSubMessageID?: string;
//   queryResults: ILabeledId[] | undefined;
// }) {
//   const {
//     option,
//     filter,
//     setFilter,
//     singleValue = false,
//     hierarchical = false,
//     includeSubMessageID,
//     queryResults,
//   } = props;

//   // defer querying until the user opens the filter
//   const [skip, setSkip] = useState(true);

//   const { criterion, setCriterion } = useCriterion(option, filter, setFilter);

//     const { selected, excluded, onSelect, onUnselect, includingOnly } =
//       useSelectionState({
//         criterion,
//         setCriterion,
//         singleValue,
//         hierarchical,
//         includeSubMessageID,
//       });

//     const candidates = useCandidates({
//       criterion,
//       queryResults,
//       selected,
//       excluded,
//       hierarchical,
//       singleValue,
//       includeSubMessageID,
//     });

//     const onOpen = useCallback(() => {
//       setSkip(false);
//     }, []);

//     return {
//       candidates,
//       onSelect,
//       onUnselect,
//       selected,
//       excluded,
//       canExclude: !includingOnly,
//       onOpen,
//     };
// }

// export const SidebarGenderFilter: React.FC<ISidebarFilter> = ({
//   title,
//   option,
//   filter,
//   setFilter,
//   singleValue = false,
//   hierarchical = false,
//   includeSubMessageID,
// }) => {
//   const intl = useIntl();

//   const queryResults = useMemo(() => {
//     return [
//       { id: "Male", label: intl.formatMessage({ id: "gender_types.MALE" }) },
//       { id: "Female", label: intl.formatMessage({ id: "gender_types.FEMALE" }) },
//       { id: "Transgender Male", label: intl.formatMessage({ id: "gender_types.TRANSGENDER_MALE" }) },
//       { id: "Transgender Female", label: intl.formatMessage({ id: "gender_types.TRANSGENDER_FEMALE" }) },
//       { id: "Intersex", label: intl.formatMessage({ id: "gender_types.INTERSEX" }) },
//       { id: "Non-Binary", label: intl.formatMessage({ id: "gender_types.NON_BINARY" }) },
//     ];
//   }, [intl]);

//   const state = useMyLabeledIdFilterState({
//     filter,
//     setFilter,
//     option,
//     queryResults,
//   });

//   return <SidebarListFilter {...state} title={title} />;
// }
