import React from "react";
import { Button, Form } from "react-bootstrap";
import { useIntl } from "react-intl";
import { PropsWithChildren } from "react-router/node_modules/@types/react";

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
  checked?: boolean;
  onChange?: React.ChangeEventHandler<HTMLInputElement>;
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
      <Form.Switch id={id} checked={checked} onChange={onChange} />
    </Setting>
  );
};

interface ISelectSetting extends ISetting {
  value?: string | number | string[] | undefined;
  onChange?: React.ChangeEventHandler | undefined;
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
        onChange={onChange}
      >
        {children}
      </Form.Control>
    </Setting>
  );
};

interface IDialogSetting extends ISetting {
  value?: string;
  onChange: () => void;
}

export const DialogSetting: React.FC<IDialogSetting> = ({
  id,
  headingID,
  subHeadingID,
  value,
  onChange,
}) => {
  const intl = useIntl();

  return (
    <div className="setting" id={id}>
      <div>
        <h3>{intl.formatMessage({ id: headingID })}</h3>

        <div className="value">{value}</div>

        {subHeadingID ? (
          <div className="sub-heading">
            {intl.formatMessage({ id: subHeadingID })}
          </div>
        ) : undefined}
      </div>
      <div>
        <Button onClick={() => onChange()}>Change</Button>
      </div>
    </div>
  );
};
