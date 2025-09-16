import React from "react";
import { Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { CriterionModifier } from "../../../core/generated-graphql";
import {
  CriterionOption,
  ModifierCriterion,
} from "../../../models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { SidebarSection } from "src/components/Shared/Sidebar";

interface IStringFilterProps {
  criterion: ModifierCriterion<string>;
  onValueChanged: (value: string) => void;
  placeholder?: string;
}

export const StringFilter: React.FC<IStringFilterProps> = ({
  criterion,
  onValueChanged,
  placeholder,
}) => {
  const intl = useIntl();

  function onValueChange(event: React.ChangeEvent<HTMLInputElement>) {
    onValueChanged(event.target.value);
  }

  return (
    <div>
      <Form.Group>
        <Form.Control
          className="btn-secondary"
          onChange={onValueChange}
          value={criterion.value}
          placeholder={placeholder || intl.formatMessage({ id: "search" })}
        />
      </Form.Group>
    </div>
  );
};

interface ISidebarFilter {
  title?: React.ReactNode;
  option: CriterionOption;
  filter: ListFilterModel;
  setFilter: (f: ListFilterModel) => void;
  placeholder?: string;
  modifier?: CriterionModifier;
}

export const SidebarStringFilter: React.FC<ISidebarFilter> = ({
  title,
  option,
  filter,
  setFilter,
  placeholder,
  modifier = CriterionModifier.Includes,
}) => {
  const intl = useIntl();

  const criteria = filter.criteriaFor(
    option.type
  ) as ModifierCriterion<string>[];
  const criterion = criteria.length > 0 ? criteria[0] : null;

  function onValueChange(event: React.ChangeEvent<HTMLInputElement>) {
    const value = event.target.value;

    if (!value.trim()) {
      // Remove criterion if empty
      setFilter(filter.removeCriterion(option.type));
      return;
    }

    const newCriterion = criterion ? criterion.clone() : option.makeCriterion();
    newCriterion.modifier = modifier;
    newCriterion.value = value;
    setFilter(filter.replaceCriteria(option.type, [newCriterion]));
  }

  return (
    <SidebarSection className="sidebar-list-filter" text={title}>
      <div className="string-filter">
        <div className="filter-group">
          <Form.Control
            className="btn-secondary"
            onChange={onValueChange}
            value={criterion?.value || ""}
            placeholder={placeholder || intl.formatMessage({ id: "search" })}
          />
        </div>
      </div>
    </SidebarSection>
  );
};

// Convenience exports for specific string filters
export const SidebarTattoosFilter: React.FC<ISidebarFilter> = (props) => (
  <SidebarStringFilter
    {...props}
    placeholder={props.placeholder || "tattoos"}
  />
);
