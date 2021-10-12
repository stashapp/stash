import React, { useState, useEffect, useMemo } from "react";
import { Form, Button, ListGroup, Row, Col } from "react-bootstrap";
import { Icon } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { ThreeStateCheckbox } from "../../Shared/ThreeStateCheckbox";
import { sceneFields } from "../constants";

interface IFieldOptionsEditor {
  availableFields: string[];
  options: GQL.IdentifyFieldOptions;
  editOptions: (o?: GQL.IdentifyFieldOptions) => void;
  removeField: () => void;
  editing: boolean;
}

const FieldOptionsEditor: React.FC<IFieldOptionsEditor> = ({
  availableFields,
  options,
  removeField,
  editOptions,
  editing,
}) => {
  const intl = useIntl();

  const [localOptions, setLocalOptions] = useState(options);

  useEffect(() => {
    setLocalOptions(options);
  }, [options]);

  function renderFieldSelect() {
    return (
      <Form.Group>
        <Form.Label>Field</Form.Label>
        <Form.Control
          disabled={!editing}
          className="w-auto input-control"
          as="select"
          value={localOptions.field}
          onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
            setLocalOptions({ ...localOptions, field: e.currentTarget.value })
          }
        >
          {availableFields.map((f) => (
            <option key={f} value={f}>
              {f}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    );
  }

  function renderStrategySelect() {
    const strategyStrings = Object.keys(GQL.IdentifyFieldStrategy);

    return (
      <Form.Group>
        <Form.Label>Strategy</Form.Label>
        <Form.Control
          disabled={!editing}
          className="w-auto input-control"
          as="select"
          value={localOptions.strategy}
          onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
            setLocalOptions({
              ...localOptions,
              strategy: e.currentTarget.value as GQL.IdentifyFieldStrategy,
            })
          }
        >
          {strategyStrings.map((f) => (
            <option key={f} value={f}>
              {f}
            </option>
          ))}
        </Form.Control>
      </Form.Group>
    );
  }

  function maybeRenderCreateMissing() {
    const createMissingFields = ["studio", "performers", "tags"];

    if (createMissingFields.includes(localOptions.field)) {
      return (
        <ThreeStateCheckbox
          value={
            localOptions.createMissing === null
              ? undefined
              : localOptions.createMissing
          }
          setValue={(v) =>
            setLocalOptions({ ...localOptions, createMissing: v })
          }
          label={intl.formatMessage({
            id: "config.tasks.identify.create_missing",
          })}
        />
      );
    }
  }

  function render() {
    return (
      <Row className="mx-2 align-items-center">
        <Col sm={3}>{renderFieldSelect()}</Col>
        <Col sm={3}>{renderStrategySelect()}</Col>
        <Col sm={3}>{maybeRenderCreateMissing()}</Col>

        <div className="col-3 d-flex justify-content-end">
          {editing ? (
            <>
              <Button
                className="minimal text-success"
                onClick={() => editOptions(localOptions)}
              >
                <Icon icon="check" />
              </Button>
              <Button
                className="minimal text-danger"
                onClick={() => editOptions()}
              >
                <Icon icon="times" />
              </Button>
            </>
          ) : (
            <Button
              className="minimal text-danger"
              onClick={() => removeField()}
            >
              <Icon icon="minus" />
            </Button>
          )}
        </div>
      </Row>
    );
  }

  return (
    <ListGroup.Item as="li" key={options.field}>
      {render()}
    </ListGroup.Item>
  );
};

interface IFieldOptionsList {
  fieldOptions?: GQL.IdentifyFieldOptions[];
  setFieldOptions: (o: GQL.IdentifyFieldOptions[]) => void;
  setEditingField: (v: boolean) => void;
}

export const FieldOptionsList: React.FC<IFieldOptionsList> = ({
  fieldOptions,
  setFieldOptions,
  setEditingField,
}) => {
  const [localFieldOptions, setLocalFieldOptions] = useState<
    GQL.IdentifyFieldOptions[]
  >([]);
  const [editField, setEditField] = useState<
    GQL.IdentifyFieldOptions | undefined
  >();

  useEffect(() => {
    if (fieldOptions) {
      setLocalFieldOptions([...fieldOptions]);
    }
  }, [fieldOptions]);

  const availableFields = useMemo(() => {
    return sceneFields.filter(
      (f) => !localFieldOptions?.some((o) => o !== editField && o.field === f)
    );
  }, [localFieldOptions, editField]);

  function onAdd() {
    const newOptions = [...localFieldOptions];
    const newOption = {
      field: availableFields[0],
      strategy: GQL.IdentifyFieldStrategy.Ignore,
    };
    newOptions.push(newOption);
    setLocalFieldOptions(newOptions);
    setEditField(newOption);
    setEditingField(true);
  }

  function handleEditOptions(o?: GQL.IdentifyFieldOptions) {
    if (!o) {
      if (localFieldOptions.length > (fieldOptions?.length ?? 0)) {
        // must be new field option. remove it
        const newOptions = [...localFieldOptions];
        newOptions.pop();
        setLocalFieldOptions(newOptions);
      }
    } else {
      const newOptions = [...localFieldOptions];
      newOptions.splice(newOptions.indexOf(editField!), 1, o);
      setFieldOptions(newOptions);
    }

    setEditField(undefined);
    setEditingField(false);
  }

  function removeField(index: number) {
    const newOptions = [...localFieldOptions];
    newOptions.splice(index, 1);
    setFieldOptions(newOptions);
  }

  return (
    <Form.Group className="scraper-sources">
      <h5>
        <FormattedMessage id="config.tasks.identify.field_options" />
      </h5>
      <ListGroup as="ul" className="scraper-source-list">
        {localFieldOptions?.map((s, index) => (
          <FieldOptionsEditor
            availableFields={availableFields}
            options={s}
            removeField={() => removeField(index)}
            editOptions={handleEditOptions}
            editing={s === editField}
          />
        ))}
      </ListGroup>
      {!editField && availableFields.length > 0 ? (
        <div className="text-right">
          <Button
            className="minimal add-scraper-source-button"
            onClick={() => onAdd()}
          >
            <Icon icon="plus" />
          </Button>
        </div>
      ) : undefined}
    </Form.Group>
  );
};
