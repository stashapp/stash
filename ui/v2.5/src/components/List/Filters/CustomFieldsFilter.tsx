import React, { useEffect } from "react";
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
import { faPlus, faTimes } from "@fortawesome/free-solid-svg-icons";

interface ICustomFieldCriterionEditor {
  criterion?: CustomFieldCriterionInput;
  setCriterion: (c: CustomFieldCriterionInput) => void;
  onRemove: () => void;
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
  onRemove,
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

    if (field === "") {
      return;
    }

    setCriterion({
      field,
      value,
      modifier: m,
    });
  }

  function onBlur() {
    if (field === "") {
      return;
    }

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
        <Row>
          <Col xs={6}>
            <Form.Control
              className="btn-secondary"
              type="text"
              placeholder={intl.formatMessage({ id: "custom_fields.field" })}
              onChange={(e) => setField(e.target.value)}
              value={field}
              onBlur={onBlur}
            />
          </Col>
          <Col xs={6}>
            <ModifierSelect
              value={modifier}
              onChanged={(m) => onChangeModifier(m)}
            />
          </Col>
        </Row>
        <Row>
          {modifier !== CriterionModifier.IsNull &&
            modifier !== CriterionModifier.NotNull && (
              <Col xs={hasTwoValues ? 6 : 12}>
                <Form.Control
                  placeholder={firstPlaceholder}
                  className="btn-secondary"
                  type="text"
                  onChange={(e) => setFirstValue(e.target.value)}
                  value={firstValue}
                  onBlur={onBlur}
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
                onBlur={onBlur}
              />
            </Col>
          )}
        </Row>
      </div>
      <div>
        <Button
          variant="minimal"
          className="text-danger"
          onClick={() => onRemove()}
        >
          <Icon icon={faTimes} />
        </Button>
      </div>
    </Form.Group>
  );
};

interface ICustomFieldsFilter {
  criterion: CustomFieldsCriterion;
  setCriterion: (c: CustomFieldsCriterion) => void;
}

function initCriterion(
  criterion: CustomFieldsCriterion
): CustomFieldsCriterion {
  const c = cloneDeep(criterion);
  if (c.value.length === 0) {
    c.value.push({
      field: "",
      value: [],
      modifier: CriterionModifier.Equals,
    });
  }

  return c;
}

export const CustomFieldsFilter: React.FC<ICustomFieldsFilter> = ({
  criterion,
  setCriterion,
}) => {
  const [localCriterion, setLocalCriterion] = React.useState(
    initCriterion(criterion)
  );

  function updateCriteria(newCriteria: CustomFieldCriterionInput[]) {
    // update the parent - filter out invalid criteria
    const validCriteria = newCriteria.filter((c) => c.field !== "");
    const newValue = cloneDeep(criterion);
    newValue.value = validCriteria;
    setCriterion(newValue);
  }

  function onChange(index: number, nv: CustomFieldCriterionInput) {
    const newValue = cloneDeep(criterion);
    const newCriteria = newValue.value.slice();
    newCriteria[index] = nv;
    newValue.value = newCriteria;

    setLocalCriterion(newValue);
    updateCriteria(newCriteria);
  }

  function onNewCriterion() {
    const c = cloneDeep(localCriterion);

    c.value.push({
      field: "",
      value: [],
      modifier: CriterionModifier.Equals,
    });

    setLocalCriterion(c);
  }

  function onRemove(index: number) {
    const c = cloneDeep(localCriterion);
    c.value.splice(index, 1);
    setLocalCriterion(c);
    updateCriteria(c.value);
  }

  return (
    <Form.Group>
      {localCriterion.value.map((cv, index) => (
        <CustomFieldCriterionEditor
          key={index}
          criterion={cv}
          setCriterion={(c) => onChange(index, c)}
          onRemove={() => onRemove(index)}
        />
      ))}
      <Button variant="success" onClick={() => onNewCriterion()}>
        <Icon icon={faPlus} />
      </Button>
    </Form.Group>
  );
};
