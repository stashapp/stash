import React, { useState, useEffect, useMemo } from "react";
import { Form, Button, Table } from "react-bootstrap";
import { Icon } from "src/components/Shared";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { multiValueSceneFields, SceneField, sceneFields } from "./constants";
import { ThreeStateBoolean } from "./ThreeStateBoolean";

interface IFieldOptionsEditor {
  availableFields: SceneField[];
  options: GQL.IdentifyFieldOptions;
  editField: () => void;
  editOptions: (o?: GQL.IdentifyFieldOptions) => void;
  removeField: () => void;
  editing: boolean;
  allowSetDefault: boolean;
  defaultOptions?: GQL.IdentifyMetadataOptionsInput;
}

const FieldOptionsEditor: React.FC<IFieldOptionsEditor> = ({
  availableFields,
  options,
  removeField,
  editField,
  editOptions,
  editing,
  allowSetDefault,
  defaultOptions,
}) => {
  const intl = useIntl();

  const [localOptions, setLocalOptions] = useState(options);

  useEffect(() => {
    setLocalOptions(options);
  }, [options]);

  function renderField() {
    if (!editing) {
      return intl.formatMessage({ id: options.field });
    }

    return (
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
            {intl.formatMessage({ id: f })}
          </option>
        ))}
      </Form.Control>
    );
  }

  function renderStrategy() {
    const strategies = Object.entries(GQL.IdentifyFieldStrategy);

    if (!editing) {
      const field = strategies.find((s) => s[1] === options.strategy);
      return intl.formatMessage({
        id: `config.tasks.identify.field_strategies.${field![0].toLowerCase()}`,
      });
    }

    return (
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
        {strategies.map((f) => (
          <option key={f[0]} value={f[1]}>
            {intl.formatMessage({
              id: `config.tasks.identify.field_strategies.${f[0].toLowerCase()}`,
            })}
          </option>
        ))}
      </Form.Control>
    );
  }

  function maybeRenderCreateMissing() {
    if (
      multiValueSceneFields.includes(localOptions.field as SceneField) &&
      localOptions.strategy !== GQL.IdentifyFieldStrategy.Ignore
    ) {
      const value =
        localOptions.createMissing === null
          ? undefined
          : localOptions.createMissing;

      if (!editing) {
        if (value === undefined) {
          return intl.formatMessage({ id: "use_default" });
        }
        if (value) {
          return <Icon icon="check" className="text-success" />;
        }

        return <Icon icon="times" className="text-danger" />;
      }

      const defaultVal = defaultOptions?.fieldOptions?.find(
        (f) => f.field === localOptions.field
      )?.createMissing;

      return (
        <ThreeStateBoolean
          disabled={!editing}
          allowUndefined={allowSetDefault}
          value={value}
          setValue={(v) =>
            setLocalOptions({ ...localOptions, createMissing: v })
          }
          defaultValue={defaultVal ?? undefined}
        />
      );
    }
  }

  return (
    <tr>
      <td>{renderField()}</td>
      <td>{renderStrategy()}</td>
      <td>{maybeRenderCreateMissing()}</td>
      <td className="text-right">
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
          <>
            <Button className="minimal" onClick={() => editField()}>
              <Icon icon="pencil-alt" />
            </Button>
            <Button
              className="minimal text-danger"
              onClick={() => removeField()}
            >
              <Icon icon="minus" />
            </Button>
          </>
        )}
      </td>
    </tr>
  );
};

interface IFieldOptionsList {
  fieldOptions?: GQL.IdentifyFieldOptions[];
  setFieldOptions: (o: GQL.IdentifyFieldOptions[]) => void;
  setEditingField: (v: boolean) => void;
  allowSetDefault?: boolean;
  defaultOptions?: GQL.IdentifyMetadataOptionsInput;
}

export const FieldOptionsList: React.FC<IFieldOptionsList> = ({
  fieldOptions,
  setFieldOptions,
  setEditingField,
  allowSetDefault = true,
  defaultOptions,
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

  const availableFields: SceneField[] = useMemo(() => {
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

  function onEditField(index: number) {
    setEditField(localFieldOptions[index]);
    setEditingField(true);
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
      {!!localFieldOptions.length && (
        <Table responsive className="field-options-table">
          <thead>
            <tr>
              <th className="w-25">Field</th>
              <th className="w-25">Strategy</th>
              <th className="w-25">Create missing</th>
              {/* eslint-disable-next-line jsx-a11y/control-has-associated-label */}
              <th className="w-25" />
            </tr>
          </thead>
          <tbody>
            {localFieldOptions?.map((s, index) => (
              <FieldOptionsEditor
                allowSetDefault={allowSetDefault}
                availableFields={availableFields}
                options={s}
                removeField={() => removeField(index)}
                editField={() => onEditField(index)}
                editOptions={handleEditOptions}
                editing={s === editField}
                defaultOptions={defaultOptions}
              />
            ))}
          </tbody>
        </Table>
      )}
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
