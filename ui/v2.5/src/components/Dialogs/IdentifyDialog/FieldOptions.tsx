import React, { useState, useEffect, useCallback } from "react";
import { Form, Button, Table } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import {
  multiValueSceneFields,
  SceneField,
  sceneFieldMessageID,
  sceneFields,
} from "./constants";
import { ThreeStateBoolean } from "./ThreeStateBoolean";
import {
  faCheck,
  faPencilAlt,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";

interface IFieldOptionsEditor {
  options: GQL.IdentifyFieldOptions | undefined;
  field: SceneField;
  editField: () => void;
  editOptions: (o?: GQL.IdentifyFieldOptions | null) => void;
  editing: boolean;
  allowSetDefault: boolean;
  defaultOptions?: GQL.IdentifyMetadataOptionsInput;
}

interface IFieldOptions {
  field: string;
  strategy: GQL.IdentifyFieldStrategy | undefined;
  createMissing?: GQL.Maybe<boolean> | undefined;
}

const FieldOptionsEditor: React.FC<IFieldOptionsEditor> = ({
  options,
  field,
  editField,
  editOptions,
  editing,
  allowSetDefault,
  defaultOptions,
}) => {
  const intl = useIntl();

  const [localOptions, setLocalOptions] = useState<IFieldOptions>();

  const resetOptions = useCallback(() => {
    let toSet: IFieldOptions;
    if (!options) {
      // unset - use default values
      toSet = {
        field,
        strategy: undefined,
        createMissing: undefined,
      };
    } else {
      toSet = {
        field,
        strategy: options.strategy,
        createMissing: options.createMissing,
      };
    }
    setLocalOptions(toSet);
  }, [options, field]);

  useEffect(() => {
    resetOptions();
  }, [resetOptions]);

  function renderField() {
    return intl.formatMessage({ id: sceneFieldMessageID(field) });
  }

  function renderStrategy() {
    if (!localOptions) {
      return;
    }

    const strategies = Object.entries(GQL.IdentifyFieldStrategy);
    let { strategy } = localOptions;
    if (strategy === undefined) {
      if (!allowSetDefault) {
        strategy = GQL.IdentifyFieldStrategy.Merge;
      }
    }

    if (!editing) {
      if (strategy === undefined) {
        return intl.formatMessage({ id: "actions.use_default" });
      }

      const f = strategies.find((s) => s[1] === strategy);
      return intl.formatMessage({
        id: `actions.${f![0].toLowerCase()}`,
      });
    }

    return (
      <Form.Group>
        {allowSetDefault ? (
          <Form.Check
            type="radio"
            id={`${field}-strategy-default`}
            checked={strategy === undefined}
            onChange={() =>
              setLocalOptions({
                ...localOptions,
                strategy: undefined,
              })
            }
            disabled={!editing}
            label={intl.formatMessage({ id: "actions.use_default" })}
          />
        ) : undefined}
        {strategies.map((f) => (
          <Form.Check
            type="radio"
            key={f[0]}
            id={`${field}-strategy-${f[0]}`}
            checked={strategy === f[1]}
            onChange={() =>
              setLocalOptions({
                ...localOptions,
                strategy: f[1],
              })
            }
            disabled={!editing}
            label={intl.formatMessage({
              id: `actions.${f[0].toLowerCase()}`,
            })}
          />
        ))}
      </Form.Group>
    );
  }

  function maybeRenderCreateMissing() {
    if (!localOptions) {
      return;
    }

    if (
      multiValueSceneFields.includes(localOptions.field as SceneField) &&
      localOptions.strategy !== GQL.IdentifyFieldStrategy.Ignore
    ) {
      const value =
        localOptions.createMissing === null
          ? undefined
          : localOptions.createMissing;

      if (!editing) {
        if (value === undefined && allowSetDefault) {
          return intl.formatMessage({ id: "actions.use_default" });
        }
        if (value) {
          return <Icon icon={faCheck} className="text-success" />;
        }

        return <Icon icon={faTimes} className="text-danger" />;
      }

      const defaultVal = defaultOptions?.fieldOptions?.find(
        (f) => f.field === localOptions.field
      )?.createMissing;

      // if allowSetDefault is false, then strategy is considered merge
      // if its true, then its using the default value and should not be shown here
      if (localOptions.strategy === undefined && allowSetDefault) {
        return;
      }

      return (
        <ThreeStateBoolean
          id="create-missing"
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

  function onEditOptions() {
    if (!localOptions) {
      return;
    }

    const localOptionsCopy = { ...localOptions };
    if (localOptionsCopy.strategy === undefined && !allowSetDefault) {
      localOptionsCopy.strategy = GQL.IdentifyFieldStrategy.Merge;
    }

    // send null if strategy is undefined
    if (localOptionsCopy.strategy === undefined) {
      editOptions(null);
      resetOptions();
    } else {
      let { createMissing } = localOptionsCopy;
      if (createMissing === undefined && !allowSetDefault) {
        createMissing = false;
      }

      editOptions({
        ...localOptionsCopy,
        strategy: localOptionsCopy.strategy,
        createMissing,
      });
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
              onClick={() => onEditOptions()}
            >
              <Icon icon={faCheck} />
            </Button>
            <Button
              className="minimal text-danger"
              onClick={() => {
                editOptions();
                resetOptions();
              }}
            >
              <Icon icon={faTimes} />
            </Button>
          </>
        ) : (
          <>
            <Button className="minimal" onClick={() => editField()}>
              <Icon icon={faPencilAlt} />
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
  const [localFieldOptions, setLocalFieldOptions] =
    useState<GQL.IdentifyFieldOptions[]>();
  const [editField, setEditField] = useState<string | undefined>();

  useEffect(() => {
    if (fieldOptions) {
      setLocalFieldOptions([...fieldOptions]);
    } else {
      setLocalFieldOptions([]);
    }
  }, [fieldOptions]);

  function handleEditOptions(o?: GQL.IdentifyFieldOptions | null) {
    if (!localFieldOptions) {
      return;
    }

    if (o !== undefined) {
      const newOptions = [...localFieldOptions];
      const index = newOptions.findIndex(
        (option) => option.field === editField
      );
      if (index !== -1) {
        // if null, then we're removing
        if (o === null) {
          newOptions.splice(index, 1);
        } else {
          // replace in list
          newOptions.splice(index, 1, o);
        }
      } else if (o !== null) {
        // don't add if null
        newOptions.push(o);
      }

      setFieldOptions(newOptions);
    }

    setEditField(undefined);
    setEditingField(false);
  }

  function onEditField(field: string) {
    setEditField(field);
    setEditingField(true);
  }

  if (!localFieldOptions) {
    return <></>;
  }

  return (
    <Form.Group className="scraper-sources mt-3">
      <h5>
        <FormattedMessage id="config.tasks.identify.field_options" />
      </h5>
      <Table responsive className="field-options-table">
        <thead>
          <tr>
            <th className="w-25">
              <FormattedMessage id="config.tasks.identify.field" />
            </th>
            <th className="w-25">
              <FormattedMessage id="config.tasks.identify.strategy" />
            </th>
            <th className="w-25">
              <FormattedMessage id="config.tasks.identify.create_missing" />
            </th>
            {/* eslint-disable-next-line jsx-a11y/control-has-associated-label */}
            <th className="w-25" />
          </tr>
        </thead>
        <tbody>
          {sceneFields.map((f) => (
            <FieldOptionsEditor
              key={f}
              field={f}
              allowSetDefault={allowSetDefault}
              options={localFieldOptions.find((o) => o.field === f)}
              editField={() => onEditField(f)}
              editOptions={handleEditOptions}
              editing={f === editField}
              defaultOptions={defaultOptions}
            />
          ))}
        </tbody>
      </Table>
    </Form.Group>
  );
};
