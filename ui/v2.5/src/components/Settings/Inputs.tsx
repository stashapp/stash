import React, { useState } from "react";
import { Button, Form, Modal } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { PropsWithChildren } from "react-router/node_modules/@types/react";
import { StringListInput } from "../Shared/StringListInput";

interface ISetting {
  id?: string;
  headingID: string;
  subHeadingID?: string;
}

const Setting: React.FC<PropsWithChildren<ISetting>> = ({
  id,
  headingID,
  subHeadingID,
  children,
}) => {
  const intl = useIntl();
  return (
    <div className="setting" id={id}>
      <div>
        <h3>{intl.formatMessage({ id: headingID })}</h3>
        {subHeadingID ? (
          <div className="sub-heading">
            {intl.formatMessage({ id: subHeadingID })}
          </div>
        ) : undefined}
      </div>
      <div>{children}</div>
    </div>
  );
};

interface IBooleanSetting extends ISetting {
  id: string;
  checked?: boolean;
  onChange: (v: boolean) => void;
}

export const BooleanSetting: React.FC<IBooleanSetting> = ({
  id,
  headingID,
  subHeadingID,
  checked,
  onChange,
}) => {
  return (
    <Setting headingID={headingID} subHeadingID={subHeadingID}>
      <Form.Switch
        id={id}
        checked={checked}
        onChange={() => onChange(!checked)}
      />
    </Setting>
  );
};

interface ISelectSetting extends ISetting {
  value?: string | number | string[] | undefined;
  onChange: (v: string) => void;
}

export const SelectSetting: React.FC<PropsWithChildren<ISelectSetting>> = ({
  id,
  headingID,
  subHeadingID,
  value,
  children,
  onChange,
}) => {
  return (
    <Setting headingID={headingID} subHeadingID={subHeadingID} id={id}>
      <Form.Control
        className="input-control"
        as="select"
        value={value}
        onChange={(e) => onChange(e.currentTarget.value)}
      >
        {children}
      </Form.Control>
    </Setting>
  );
};

interface IDialogSetting<T> extends ISetting {
  value?: T;
  renderValue: (v: T | undefined) => JSX.Element;
  onChange: () => void;
}

export const ChangeButtonSetting = <T extends {}>(props: IDialogSetting<T>) => {
  const { id, headingID, subHeadingID, value, onChange, renderValue } = props;
  const intl = useIntl();

  return (
    <div className="setting" id={id}>
      <div>
        <h3>{intl.formatMessage({ id: headingID })}</h3>

        <div className="value">{renderValue(value)}</div>

        {subHeadingID ? (
          <div className="sub-heading">
            {intl.formatMessage({ id: subHeadingID })}
          </div>
        ) : undefined}
      </div>
      <div>
        <Button onClick={() => onChange()}>
          <FormattedMessage id="actions.edit" />
        </Button>
      </div>
    </div>
  );
};

export interface ISettingModal<T> {
  headingID: string;
  subHeadingID?: string;
  value: T | undefined;
  close: (v?: T) => void;
  renderField: (value: T | undefined, setValue: (v?: T) => void) => JSX.Element;
}

export const SettingModal = <T extends {}>(props: ISettingModal<T>) => {
  const { headingID, subHeadingID, value, close, renderField } = props;

  const intl = useIntl();
  const [currentValue, setCurrentValue] = useState<T | undefined>(value);

  return (
    <Modal show onHide={() => close()} id="setting-dialog">
      <Form
        onSubmit={(e) => {
          close(currentValue);
          e.preventDefault();
        }}
      >
        <Modal.Header closeButton>
          <FormattedMessage id={headingID} />
        </Modal.Header>
        <Modal.Body>
          {renderField(currentValue, setCurrentValue)}
          {subHeadingID ? (
            <div className="sub-heading">
              {intl.formatMessage({ id: subHeadingID })}
            </div>
          ) : undefined}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => close()}>
            Cancel
          </Button>
          <Button
            type="submit"
            variant="primary"
            onClick={() => close(currentValue)}
          >
            Confirm
          </Button>
        </Modal.Footer>
      </Form>
    </Modal>
  );
};

interface IModalSetting<T> extends ISetting {
  value: T | undefined;
  onChange: (v: T) => void;
  renderField: (value: T | undefined, setValue: (v?: T) => void) => JSX.Element;
  renderValue: (v: T | undefined) => JSX.Element;
}

export const ModalSetting = <T extends {}>(props: IModalSetting<T>) => {
  const {
    id,
    value,
    headingID,
    subHeadingID,
    onChange,
    renderField,
    renderValue,
  } = props;
  const [showModal, setShowModal] = useState(false);

  return (
    <>
      {showModal ? (
        <SettingModal<T>
          headingID={headingID}
          subHeadingID={subHeadingID}
          value={value}
          renderField={renderField}
          close={(v) => {
            if (v !== undefined) onChange(v);
            setShowModal(false);
          }}
        />
      ) : undefined}

      <ChangeButtonSetting<T>
        id={id}
        headingID={headingID}
        subHeadingID={subHeadingID}
        value={value}
        onChange={() => setShowModal(true)}
        renderValue={renderValue}
      />
    </>
  );
};

interface IStringSetting extends ISetting {
  value: string | undefined;
  onChange: (v: string) => void;
}

export const StringSetting: React.FC<IStringSetting> = (props) => {
  return (
    <ModalSetting<string>
      {...props}
      renderField={(value, setValue) => (
        <Form.Control
          className="text-input"
          value={value}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setValue(e.currentTarget.value)
          }
        />
      )}
      renderValue={(value) => <span>{value}</span>}
    />
  );
};

interface INumberSetting extends ISetting {
  value: number | undefined;
  onChange: (v: number) => void;
}

export const NumberSetting: React.FC<INumberSetting> = (props) => {
  return (
    <ModalSetting<number>
      {...props}
      renderField={(value, setValue) => (
        <Form.Control
          className="text-input"
          type="number"
          value={value}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setValue(Number.parseInt(e.currentTarget.value || "0", 10))
          }
        />
      )}
      renderValue={(value) => <span>{value}</span>}
    />
  );
};

interface IStringListSetting extends ISetting {
  value: string[] | undefined;
  onChange: (v: string[]) => void;
}

export const StringListSetting: React.FC<IStringListSetting> = (props) => {
  return (
    <ModalSetting<string[]>
      {...props}
      renderField={(value, setValue) => (
        <StringListInput value={value ?? []} setValue={setValue} />
      )}
      renderValue={(value) => (
        <div>
          {value?.map((v, i) => (
            // eslint-disable-next-line react/no-array-index-key
            <div key={i}>{v}</div>
          ))}
        </div>
      )}
    />
  );
};
