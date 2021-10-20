import React from "react";
import { Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { IScraperSource } from "./constants";
import { FieldOptionsList } from "./FieldOptions";
import { ThreeStateBoolean } from "./ThreeStateBoolean";

interface IOptionsEditor {
  options: GQL.IdentifyMetadataOptionsInput;
  setOptions: (s: GQL.IdentifyMetadataOptionsInput) => void;
  source?: IScraperSource;
  defaultOptions?: GQL.IdentifyMetadataOptionsInput;
  setEditingField: (v: boolean) => void;
}

export const OptionsEditor: React.FC<IOptionsEditor> = ({
  options,
  setOptions: setOptionsState,
  source,
  setEditingField,
  defaultOptions,
}) => {
  const intl = useIntl();

  function setOptions(v: Partial<GQL.IdentifyMetadataOptionsInput>) {
    setOptionsState({ ...options, ...v });
  }

  const headingID = !source
    ? "config.tasks.identify.default_options"
    : "config.tasks.identify.source_options";
  const checkboxProps = {
    allowUndefined: !!source,
    indeterminateClassname: "text-muted",
  };

  return (
    <Form.Group>
      <h5>
        <FormattedMessage
          id={headingID}
          values={{ source: source?.displayName }}
        />
      </h5>
      <Form.Group>
        <ThreeStateBoolean
          id="include-male-performers"
          value={
            options.includeMalePerformers === null
              ? undefined
              : options.includeMalePerformers
          }
          setValue={(v) =>
            setOptions({
              includeMalePerformers: v,
            })
          }
          label={intl.formatMessage({
            id: "config.tasks.identify.include_male_performers",
          })}
          defaultValue={defaultOptions?.includeMalePerformers ?? undefined}
          {...checkboxProps}
        />
        <ThreeStateBoolean
          id="set-cover-image"
          value={
            options.setCoverImage === null ? undefined : options.setCoverImage
          }
          setValue={(v) =>
            setOptions({
              setCoverImage: v,
            })
          }
          label={intl.formatMessage({
            id: "config.tasks.identify.set_cover_images",
          })}
          defaultValue={defaultOptions?.setCoverImage ?? undefined}
          {...checkboxProps}
        />
        <ThreeStateBoolean
          id="set-organized"
          value={
            options.setOrganized === null ? undefined : options.setOrganized
          }
          setValue={(v) =>
            setOptions({
              setOrganized: v,
            })
          }
          label={intl.formatMessage({
            id: "config.tasks.identify.set_organized",
          })}
          defaultValue={defaultOptions?.setOrganized ?? undefined}
          {...checkboxProps}
        />
      </Form.Group>

      <FieldOptionsList
        fieldOptions={options.fieldOptions ?? undefined}
        setFieldOptions={(o) => setOptions({ fieldOptions: o })}
        setEditingField={setEditingField}
        allowSetDefault={!!source}
        defaultOptions={defaultOptions}
      />

      {!source && (
        <Form.Text className="text-muted">
          {intl.formatMessage({
            id: "config.tasks.identify.explicit_set_description",
          })}
        </Form.Text>
      )}
    </Form.Group>
  );
};
