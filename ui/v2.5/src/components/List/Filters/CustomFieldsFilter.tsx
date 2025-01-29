import React, { useEffect, useMemo, useState } from "react";
import { CustomFieldsCriterion } from "src/models/list-filter/criteria/custom-fields";
import { Button, Col, Form, Row } from "react-bootstrap";
import {
  CriterionModifier,
  CustomFieldCriterionInput,
} from "src/core/generated-graphql";
import { cloneDeep } from "@apollo/client/utilities";
import { ModifierSelect } from "../ModifierSelect";
import { useIntl } from "react-intl";
import { Icon } from "src/components/Shared/Icon";
import { faCheck, faPencil, faTimes } from "@fortawesome/free-solid-svg-icons";
import { FilterTag } from "../FilterTags";
import { ModifierCriterion } from "src/models/list-filter/criteria/criterion";

interface ICustomFieldCriterionEditor {
  criterion?: CustomFieldCriterionInput;
  setCriterion: (c: CustomFieldCriterionInput) => void;
  cancel: () => void;
  editing?: boolean;
}

function getValue(v: string) {
  // if the value is numeric, convert it to a number
  const num = Number(v);
  if (!isNaN(num)) {
    return num;
  } else {
    return v;
  }
}

const CustomFieldCriterionEditor: React.FC<ICustomFieldCriterionEditor> = ({
  criterion,
  setCriterion,
  editing = false,
  cancel,
}) => {
  const intl = useIntl();

  const [field, setField] = React.useState(criterion?.field ?? "");
  const [value, setValue] = React.useState(criterion?.value);
  const [modifier, setModifier] = React.useState(
    criterion?.modifier ?? CriterionModifier.Equals
  );

  const firstValue = value && value.length > 0 ? (value[0] as string) : "";
  const secondValue = value && value.length > 1 ? (value[1] as string) : "";

  useEffect(() => {
    setField((criterion?.field as string) ?? "");
    setValue(criterion?.value ?? []);
    setModifier(criterion?.modifier ?? CriterionModifier.Equals);
  }, [criterion]);

  function setFirstValue(v: string) {
    // convert to numeric if possible
    const nv = getValue(v);

    if (
      modifier === CriterionModifier.Between ||
      modifier === CriterionModifier.NotBetween
    ) {
      setValue([nv, secondValue]);
    } else {
      setValue([nv]);
    }
  }

  function setSecondValue(v: string) {
    setValue([firstValue, getValue(v)]);
  }

  function onChangeModifier(m: CriterionModifier) {
    setModifier(m);
    if (m === CriterionModifier.IsNull || m === CriterionModifier.NotNull) {
      setValue(undefined);
    }
  }

  function onConfirm() {
    setCriterion({
      field,
      value,
      modifier,
    });
  }

  const firstPlaceholder =
    modifier === CriterionModifier.Between ||
    modifier === CriterionModifier.NotBetween
      ? intl.formatMessage({ id: "criterion.greater_than" })
      : intl.formatMessage({ id: "custom_fields.value" });

  const hasTwoValues =
    modifier === CriterionModifier.Between ||
    modifier === CriterionModifier.NotBetween;

  return (
    <Form.Group className="custom-field-filter">
      <div>
        <Row noGutters>
          <Col xs={6}>
            <Form.Control
              className="btn-secondary"
              type="text"
              placeholder={intl.formatMessage({ id: "custom_fields.field" })}
              onChange={(e) => setField(e.target.value)}
              value={field}
            />
          </Col>
          <Col xs={6}>
            <ModifierSelect
              value={modifier}
              onChanged={(m) => onChangeModifier(m)}
            />
          </Col>
        </Row>
        <Row noGutters>
          {modifier !== CriterionModifier.IsNull &&
            modifier !== CriterionModifier.NotNull && (
              <Col xs={hasTwoValues ? 6 : 12}>
                <Form.Control
                  placeholder={firstPlaceholder}
                  className="btn-secondary"
                  type="text"
                  onChange={(e) => setFirstValue(e.target.value)}
                  value={firstValue}
                />
              </Col>
            )}
          {(modifier === CriterionModifier.Between ||
            modifier === CriterionModifier.NotBetween) && (
            <Col xs={6}>
              <Form.Control
                placeholder={intl.formatMessage({ id: "criterion.less_than" })}
                className="btn-secondary"
                type="text"
                onChange={(e) => setSecondValue(e.target.value)}
                value={secondValue}
              />
            </Col>
          )}
        </Row>
      </div>
      <div className="custom-field-filter-buttons">
        <Button variant="success" onClick={() => onConfirm()} disabled={!field}>
          <Icon icon={faCheck} />
        </Button>
        {editing && (
          <Button variant="secondary" onClick={() => cancel()}>
            <Icon icon={faTimes} />
          </Button>
        )}
      </div>
    </Form.Group>
  );
};

function valueToString(value: unknown[] | undefined | null) {
  if (!value) return "";
  return value.map((v) => v as string).join(", ");
}

const CustomFieldFilterTag: React.FC<{
  criterion: CustomFieldCriterionInput;
  editing?: boolean;
  onEditCriterion: () => void;
  onRemoveCriterion: () => void;
}> = ({ criterion, editing, onEditCriterion, onRemoveCriterion }) => {
  const intl = useIntl();

  const label = useMemo(() => {
    const { field, modifier, value } = criterion;
    const modifierString = ModifierCriterion.getModifierLabel(intl, modifier);

    const str = intl.formatMessage(
      { id: "criterion_modifier.format_string" },
      {
        criterion: field,
        modifierString,
        valueString: valueToString(value),
      }
    );

    if (editing) {
      return (
        <span>
          <Icon icon={faPencil} />
          {str}
        </span>
      );
    }

    return <>{str}</>;
  }, [criterion, editing, intl]);

  return (
    <FilterTag
      label={label}
      onClick={onEditCriterion}
      onRemove={onRemoveCriterion}
    />
  );
};

const CustomFieldsCriteriaPills: React.FC<{
  criteria: CustomFieldCriterionInput[];
  editIndex?: number;
  onEditCriterion: (index: number) => void;
  onRemoveCriterion: (index: number) => void;
}> = ({ criteria, editIndex, onEditCriterion, onRemoveCriterion }) => {
  return (
    <div className="d-flex justify-content-center mb-2 wrap-tags filter-tags">
      {criteria.map((c, index) => (
        <CustomFieldFilterTag
          key={index}
          editing={index === editIndex}
          criterion={c}
          onEditCriterion={() => onEditCriterion(index)}
          onRemoveCriterion={() => onRemoveCriterion(index)}
        />
      ))}
    </div>
  );
};

interface ICustomFieldsFilter {
  criterion: CustomFieldsCriterion;
  setCriterion: (c: CustomFieldsCriterion) => void;
}

function initCriterion(
  criterion: CustomFieldsCriterion
): CustomFieldsCriterion {
  return cloneDeep(criterion);
}

function createNewCriterion(): CustomFieldCriterionInput {
  return {
    field: "",
    value: [],
    modifier: CriterionModifier.Equals,
  };
}

export const CustomFieldsFilter: React.FC<ICustomFieldsFilter> = ({
  criterion,
  setCriterion,
}) => {
  const [localCriterion, setLocalCriterion] = React.useState(
    initCriterion(criterion)
  );

  const [editCriterion, setEditCriterion] = useState(createNewCriterion());
  const editIndex = useMemo(
    () => localCriterion.value.indexOf(editCriterion),
    [localCriterion, editCriterion]
  );

  function updateCriteria(newCriteria: CustomFieldCriterionInput[]) {
    // update the parent - filter out invalid criteria
    const validCriteria = newCriteria.filter((c) => c.field !== "");
    const newValue = cloneDeep(criterion);
    newValue.value = validCriteria;
    setCriterion(newValue);
  }

  function onChange(nv: CustomFieldCriterionInput) {
    const newValue = cloneDeep(localCriterion);

    // if the criterion is new, add it to the list
    if (editIndex === -1) {
      newValue.value.push(nv);
    } else {
      newValue.value[editIndex] = nv;
    }

    setLocalCriterion(newValue);
    updateCriteria(newValue.value);
    setEditCriterion(createNewCriterion());
  }

  function onRemove(index: number) {
    const c = cloneDeep(localCriterion);
    c.value.splice(index, 1);
    setLocalCriterion(c);
    updateCriteria(c.value);
    if (index === editIndex) {
      setEditCriterion(createNewCriterion());
    }
  }

  return (
    <Form.Group>
      <CustomFieldCriterionEditor
        criterion={editCriterion}
        editing={editCriterion.field !== ""}
        setCriterion={onChange}
        cancel={() => setEditCriterion(createNewCriterion())}
      />
      <CustomFieldsCriteriaPills
        criteria={localCriterion.value}
        editIndex={editIndex !== -1 ? editIndex : undefined}
        onEditCriterion={(index) =>
          setEditCriterion(localCriterion.value[index])
        }
        onRemoveCriterion={(index) => onRemove(index)}
      />
    </Form.Group>
  );
};
