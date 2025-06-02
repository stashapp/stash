import React, { useRef, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { SettingSection } from "./SettingSection";
import * as GQL from "src/core/generated-graphql";
import { SettingModal } from "./Inputs";

export interface IStashBoxModal {
  value: GQL.StashBoxInput;
  close: (v?: GQL.StashBoxInput) => void;
}

const defaultMaxRequestsPerMinute = 240;

export const StashBoxModal: React.FC<IStashBoxModal> = ({ value, close }) => {
  const intl = useIntl();
  const endpoint = useRef<HTMLInputElement | null>(null);
  const apiKey = useRef<HTMLInputElement | null>(null);

  const [validate, { data, loading }] = GQL.useValidateStashBoxLazyQuery({
    fetchPolicy: "network-only",
  });

  const handleValidate = () => {
    validate({
      variables: {
        input: {
          endpoint: endpoint.current?.value ?? "",
          api_key: apiKey.current?.value ?? "",
          name: "test",
        },
      },
    });
  };

  return (
    <SettingModal<GQL.StashBoxInput>
      headingID="config.stashbox.title"
      value={value}
      renderField={(v, setValue) => (
        <>
          <Form.Group id="stashbox-name">
            <h6>
              {intl.formatMessage({
                id: "config.stashbox.name",
              })}
            </h6>
            <Form.Control
              placeholder={intl.formatMessage({ id: "config.stashbox.name" })}
              className="text-input stash-box-name"
              value={v?.name}
              isValid={(v?.name?.length ?? 0) > 0}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setValue({ ...v!, name: e.currentTarget.value })
              }
            />
          </Form.Group>

          <Form.Group id="stashbox-name">
            <h6>
              {intl.formatMessage({
                id: "config.stashbox.graphql_endpoint",
              })}
            </h6>
            <Form.Control
              placeholder={intl.formatMessage({
                id: "config.stashbox.graphql_endpoint",
              })}
              className="text-input stash-box-endpoint"
              value={v?.endpoint}
              isValid={(v?.endpoint?.length ?? 0) > 0}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setValue({ ...v!, endpoint: e.currentTarget.value.trim() })
              }
              ref={endpoint}
            />
          </Form.Group>

          <Form.Group id="stashbox-name">
            <h6>
              {intl.formatMessage({
                id: "config.stashbox.api_key",
              })}
            </h6>
            <Form.Control
              placeholder={intl.formatMessage({
                id: "config.stashbox.api_key",
              })}
              className="text-input stash-box-apikey"
              value={v?.api_key}
              isValid={(v?.api_key?.length ?? 0) > 0}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setValue({ ...v!, api_key: e.currentTarget.value.trim() })
              }
              ref={apiKey}
            />
          </Form.Group>

          <Form.Group>
            <Button
              disabled={loading}
              onClick={handleValidate}
              className="mr-3"
            >
              Test Credentials
            </Button>
            {data && (
              <b
                className={
                  data.validateStashBoxCredentials?.valid
                    ? "text-success"
                    : "text-danger"
                }
              >
                {data.validateStashBoxCredentials?.status}
              </b>
            )}
          </Form.Group>

          <Form.Group id="stashbox-max-requests-per-minute">
            <h6>
              {intl.formatMessage({
                id: "config.stashbox.max_requests_per_minute",
              })}
            </h6>
            <Form.Control
              placeholder={intl.formatMessage({
                id: "config.stashbox.max_requests_per_minute",
              })}
              className="text-input"
              value={v?.max_requests_per_minute ?? defaultMaxRequestsPerMinute}
              isValid={
                (v?.max_requests_per_minute ?? defaultMaxRequestsPerMinute) >= 0
              }
              type="number"
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setValue({
                  ...v!,
                  max_requests_per_minute: parseInt(e.currentTarget.value),
                })
              }
            />
            <div className="sub-heading">
              <FormattedMessage
                id="config.stashbox.max_requests_per_minute_description"
                values={{ defaultValue: defaultMaxRequestsPerMinute }}
              />
            </div>
          </Form.Group>
        </>
      )}
      close={close}
    />
  );
};

interface IStashBoxSetting {
  value: GQL.StashBoxInput[];
  onChange: (v: GQL.StashBoxInput[]) => void;
}

export const StashBoxSetting: React.FC<IStashBoxSetting> = ({
  value,
  onChange,
}) => {
  const [isCreating, setIsCreating] = useState(false);
  const [editingIndex, setEditingIndex] = useState<number | undefined>();

  function onEdit(index: number) {
    setEditingIndex(index);
  }

  function onDelete(index: number) {
    onChange(value.filter((v, i) => i !== index));
  }

  function onNew() {
    setIsCreating(true);
  }

  return (
    <SettingSection
      id="stash-boxes"
      headingID="config.stashbox.title"
      subHeadingID="config.stashbox.description"
    >
      {isCreating ? (
        <StashBoxModal
          value={{
            endpoint: "",
            api_key: "",
            name: "",
          }}
          close={(v) => {
            if (v) onChange([...value, v]);
            setIsCreating(false);
          }}
        />
      ) : undefined}

      {editingIndex !== undefined ? (
        <StashBoxModal
          value={value[editingIndex]}
          close={(v) => {
            if (v)
              onChange(
                value.map((vv, index) => {
                  if (index === editingIndex) {
                    return v;
                  }
                  return vv;
                })
              );
            setEditingIndex(undefined);
          }}
        />
      ) : undefined}

      {value.map((b, index) => (
        // eslint-disable-next-line react/no-array-index-key
        <div key={index} className="setting">
          <div>
            <h3>{b.name ?? `#${index}`}</h3>
            <div className="value">{b.endpoint ?? ""}</div>
          </div>
          <div>
            <Button onClick={() => onEdit(index)}>
              <FormattedMessage id="actions.edit" />
            </Button>
            <Button variant="danger" onClick={() => onDelete(index)}>
              <FormattedMessage id="actions.delete" />
            </Button>
          </div>
        </div>
      ))}
      <div className="setting">
        <div />
        <div>
          <Button onClick={() => onNew()}>
            <FormattedMessage id="actions.add" />
          </Button>
        </div>
      </div>
    </SettingSection>
  );
};
