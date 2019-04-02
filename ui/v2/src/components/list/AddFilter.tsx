import {
  Button,
  Classes,
  Dialog,
  FormGroup,
  HTMLSelect,
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent, useEffect, useRef, useState } from "react";
import { isArray } from "util";
import { CriterionModifier } from "../../core/generated-graphql";
import { Criterion, CriterionType } from "../../models/list-filter/criteria/criterion";
import { NoneCriterion } from "../../models/list-filter/criteria/none";
import { PerformersCriterion } from "../../models/list-filter/criteria/performers";
import { StudiosCriterion } from "../../models/list-filter/criteria/studios";
import { TagsCriterion } from "../../models/list-filter/criteria/tags";
import { makeCriteria } from "../../models/list-filter/criteria/utils";
import { ListFilterModel } from "../../models/list-filter/filter";
import { FilterMultiSelect } from "../select/FilterMultiSelect";

interface IAddFilterProps {
  onAddCriterion: (criterion: Criterion, oldId?: string) => void;
  onCancel: () => void;
  filter: ListFilterModel;
  editingCriterion?: Criterion;
}

export const AddFilter: FunctionComponent<IAddFilterProps> = (props: IAddFilterProps) => {
  const singleValueSelect = useRef<HTMLSelect>(null);

  const [isOpen, setIsOpen] = useState(false);
  const [criterion, setCriterion] = useState<Criterion<any, any>>(new NoneCriterion());

  // Configure if we are editing an existing criterion
  useEffect(() => {
    if (!props.editingCriterion) { return; }
    setIsOpen(true);
    setCriterion(props.editingCriterion);
  }, [props.editingCriterion]);

  function onChangedCriteriaType(event: React.ChangeEvent<HTMLSelectElement>) {
    const newCriterionType = event.target.value as CriterionType;
    const newCriterion = makeCriteria(newCriterionType);
    setCriterion(newCriterion);
  }

  function onChangedModifierSelect(event: React.ChangeEvent<HTMLSelectElement>) {
    const newCriterion = _.cloneDeep(criterion);
    newCriterion.modifier = event.target.value as any;
    setCriterion(newCriterion);
  }

  function onChangedSingleSelect(event: React.ChangeEvent<HTMLSelectElement>) {
    const newCriterion = _.cloneDeep(criterion);
    newCriterion.value = event.target.value;
    setCriterion(newCriterion);
  }

  function onAddFilter() {
    if (!isArray(criterion.value) && !!singleValueSelect.current) {
      const value = singleValueSelect.current.props.defaultValue;
      if (value === undefined || value === "" || typeof value === "number") { criterion.value = criterion.options[0]; }
    }
    const oldId = !!props.editingCriterion ? props.editingCriterion.getId() : undefined;
    props.onAddCriterion(criterion, oldId);
    onToggle();
  }

  function onToggle() {
    if (isOpen) {
      props.onCancel();
    }
    setIsOpen(!isOpen);
    setCriterion(makeCriteria());
  }

  const maybeRenderFilterPopoverContents = () => {
    if (criterion.type === "none") { return; }

    function renderModifier() {
      if (criterion.modifierOptions.length === 0) { return; }
      return (
        <div>
          <HTMLSelect
            options={criterion.modifierOptions}
            onChange={onChangedModifierSelect}
            defaultValue={criterion.modifier}
          />
        </div>
      );
    }

    function renderSelect() {
      // Hide the value select if the modifier is "IsNull" or "NotNull"
      if (criterion.modifier === CriterionModifier.IsNull || criterion.modifier === CriterionModifier.NotNull) {
        return;
      }

      if (isArray(criterion.value)) {
        let type: "performers" | "studios" | "tags" | "" = "";
        if (criterion instanceof PerformersCriterion) {
          type = "performers";
        } else if (criterion instanceof StudiosCriterion) {
          type = "studios";
        } else if (criterion instanceof TagsCriterion) {
          type = "tags";
        }

        if (type === "") {
          return (<>todo</>);
        } else {
          return (
            <FilterMultiSelect
              type={type}
              onUpdate={(items) => criterion.value = items.map((i) => ({id: i.id, label: i.name!}))}
              openOnKeyDown={true}
              initialIds={criterion.value.map((labeled: any) => labeled.id)}
            />
          );
        }
      } else {
        return (
          <HTMLSelect
            ref={singleValueSelect}
            options={criterion.options}
            onChange={onChangedSingleSelect}
            defaultValue={criterion.value}
          />
        );
      }
    }
    return (
      <FormGroup>
        {renderModifier()}
        {renderSelect()}
      </FormGroup>
    );
  };

  function maybeRenderFilterSelect() {
    if (!!props.editingCriterion) { return; }
    return (
      <FormGroup label="Filter">
        <HTMLSelect
          style={{flexBasis: "min-content"}}
          options={props.filter.criterionOptions}
          onChange={onChangedCriteriaType}
          defaultValue={criterion.type}
        />
      </FormGroup>
    );
  }

  const title = !props.editingCriterion ? "Add Filter" : "Update Filter";
  return (
    <>
      <Button onClick={() => onToggle()} active={isOpen} large={true}>Filter</Button>
      <Dialog isOpen={isOpen} onClose={() => onToggle()} title={title}>
        <div className="dialog-content">
          {maybeRenderFilterSelect()}
          {maybeRenderFilterPopoverContents()}
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button onClick={onAddFilter} disabled={criterion.type === "none"}>{title}</Button>
          </div>
        </div>
      </Dialog>
    </>
  );
};
