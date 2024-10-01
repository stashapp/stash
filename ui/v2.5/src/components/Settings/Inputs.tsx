import { faChevronDown, faChevronUp } from "@fortawesome/free-solid-svg-icons";
import React, { PropsWithChildren, useState } from "react";
import { Button, Collapse, Form, Modal, ModalProps } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "../Shared/Icon";
import { StringListInput } from "../Shared/StringListInput";
import { PatchComponent } from "src/patch";
import { useSettings, useSettingsOptional } from "./context";
import { NumberField } from "src/utils/form";

interface ISetting {
  id?: string;
  advanced?: boolean;
  className?: string;
  heading?: React.ReactNode;
  headingID?: string;
  subHeadingID?: string;
  subHeading?: React.ReactNode;
  tooltipID?: string;
  onClick?: React.MouseEventHandler<HTMLDivElement>;
  disabled?: boolean;
}

export const Setting: React.FC<PropsWithChildren<ISetting>> = PatchComponent(
  "Setting",
  (props: PropsWithChildren<ISetting>) => {
    const {
      id,
      className,
      heading,
      headingID,
      subHeadingID,
      subHeading,
      children,
      tooltipID,
      onClick,
      disabled,
      advanced,
    } = props;

    // these components can be used in the setup wizard, where advanced mode is not available
    const { advancedMode } = useSettingsOptional();

    const intl = useIntl();

    function renderHeading() {
      if (headingID) {
        return intl.formatMessage({ id: headingID });
      }
      return heading;
    }

    function renderSubHeading() {
      if (subHeadingID) {
        return (
          <div className="sub-heading">
            {intl.formatMessage({ id: subHeadingID })}
          </div>
        );
      }
      if (subHeading) {
        return <div className="sub-heading">{subHeading}</div>;
      }
    }

    const tooltip = tooltipID
      ? intl.formatMessage({ id: tooltipID })
      : undefined;
    const disabledClassName = disabled ? "disabled" : "";

    if (advanced && !advancedMode) return null;

    return (
      <div
        className={`setting ${className ?? ""} ${disabledClassName}`}
        id={id}
        onClick={onClick}
      >
        <div>
          <h3 title={tooltip}>{renderHeading()}</h3>
          {renderSubHeading()}
        </div>
        <div>{children}</div>
      </div>
    );
  }
) as React.FC<PropsWithChildren<ISetting>>;

interface ISettingGroup {
  settingProps?: ISetting;
  topLevel?: JSX.Element;
  collapsible?: boolean;
  collapsedDefault?: boolean;
}

export const SettingGroup: React.FC<PropsWithChildren<ISettingGroup>> =
  PatchComponent(
    "SettingGroup",
    ({ settingProps, topLevel, collapsible, collapsedDefault, children }) => {
      const [open, setOpen] = useState(!collapsedDefault);

      function renderCollapseButton() {
        if (!collapsible) return;

        return (
          <Button
            className="setting-group-collapse-button"
            variant="minimal"
            onClick={() => setOpen(!open)}
          >
            <Icon className="fa-fw" icon={open ? faChevronUp : faChevronDown} />
          </Button>
        );
      }

      function onDivClick(e: React.MouseEvent<HTMLDivElement>) {
        if (!collapsible) return;

        // ensure button was not clicked
        let target: HTMLElement | null = e.target as HTMLElement;
        while (target && target !== e.currentTarget) {
          if (
            target.nodeName.toLowerCase() === "button" ||
            target.nodeName.toLowerCase() === "a"
          ) {
            // button clicked, swallow event
            return;
          }
          target = target.parentElement;
        }

        setOpen(!open);
      }

      return (
        <div className={`setting-group ${collapsible ? "collapsible" : ""}`}>
          <Setting {...settingProps} onClick={onDivClick}>
            {topLevel}
            {renderCollapseButton()}
          </Setting>
          <Collapse in={open}>
            <div className="collapsible-section">{children}</div>
          </Collapse>
        </div>
      );
    }
  );

interface IBooleanSetting extends ISetting {
  id: string;
  checked?: boolean;
  onChange: (v: boolean) => void;
}

export const BooleanSetting: React.FC<IBooleanSetting> = PatchComponent(
  "BooleanSetting",
  (props) => {
    const { id, disabled, checked, onChange, ...settingProps } = props;

    return (
      <Setting {...settingProps} disabled={disabled}>
        <Form.Switch
          id={id}
          disabled={disabled}
          checked={checked ?? false}
          onChange={() => onChange(!checked)}
        />
      </Setting>
    );
  }
);

interface ISelectSetting extends ISetting {
  value?: string | number | string[];
  onChange: (v: string) => void;
}

export const SelectSetting: React.FC<PropsWithChildren<ISelectSetting>> =
  PatchComponent(
    "SelectSetting",
    ({ id, headingID, subHeadingID, value, children, onChange, advanced }) => {
      return (
        <Setting
          advanced={advanced}
          headingID={headingID}
          subHeadingID={subHeadingID}
          id={id}
        >
          <Form.Control
            className="input-control"
            as="select"
            value={value ?? ""}
            onChange={(e) => onChange(e.currentTarget.value)}
          >
            {children}
          </Form.Control>
        </Setting>
      );
    }
  );

interface IDialogSetting<T> extends ISetting {
  buttonText?: string;
  buttonTextID?: string;
  value?: T;
  renderValue?: (v: T | undefined) => JSX.Element;
  onChange: () => void;
}
const _ChangeButtonSetting = <T extends {}>(props: IDialogSetting<T>) => {
  const {
    id,
    className,
    headingID,
    heading,
    tooltipID,
    subHeadingID,
    subHeading,
    value,
    onChange,
    renderValue,
    buttonText,
    buttonTextID,
    disabled,
  } = props;
  const intl = useIntl();

  const tooltip = tooltipID ? intl.formatMessage({ id: tooltipID }) : undefined;
  const disabledClassName = disabled ? "disabled" : "";

  return (
    <div className={`setting ${className ?? ""} ${disabledClassName}`} id={id}>
      <div>
        <h3 title={tooltip}>
          {headingID
            ? intl.formatMessage({ id: headingID })
            : heading
            ? heading
            : undefined}
        </h3>

        <div className="value">
          {renderValue ? renderValue(value) : undefined}
        </div>

        {subHeadingID ? (
          <div className="sub-heading">
            {intl.formatMessage({ id: subHeadingID })}
          </div>
        ) : subHeading ? (
          <div className="sub-heading">{subHeading}</div>
        ) : undefined}
      </div>
      <div>
        <Button onClick={() => onChange()} disabled={disabled}>
          {buttonText ? (
            buttonText
          ) : (
            <FormattedMessage id={buttonTextID ?? "actions.edit"} />
          )}
        </Button>
      </div>
    </div>
  );
};

export const ChangeButtonSetting = PatchComponent(
  "ChangeButtonSetting",
  _ChangeButtonSetting
) as typeof _ChangeButtonSetting;

export interface ISettingModal<T> {
  heading?: React.ReactNode;
  headingID?: string;
  subHeadingID?: string;
  subHeading?: React.ReactNode;
  value: T | undefined;
  close: (v?: T) => void;
  renderField: (
    value: T | undefined,
    setValue: (v?: T) => void,
    error?: string
  ) => JSX.Element;
  modalProps?: ModalProps;
  validate?: (v: T) => boolean | undefined;
  error?: string | undefined;
}

const _SettingModal = <T extends {}>(props: ISettingModal<T>) => {
  const {
    heading,
    headingID,
    subHeading,
    subHeadingID,
    value,
    close,
    renderField,
    modalProps,
    validate,
    error,
  } = props;

  const intl = useIntl();
  const [currentValue, setCurrentValue] = useState<T | undefined>(value);

  return (
    <Modal show onHide={() => close()} id="setting-dialog" {...modalProps}>
      <Form
        onSubmit={(e) => {
          close(currentValue);
          e.preventDefault();
        }}
      >
        <Modal.Header closeButton>
          {headingID ? <FormattedMessage id={headingID} /> : heading}
        </Modal.Header>
        <Modal.Body>
          {renderField(currentValue, setCurrentValue, error)}
          {subHeadingID ? (
            <div className="sub-heading">
              {intl.formatMessage({ id: subHeadingID })}
            </div>
          ) : subHeading ? (
            <div className="sub-heading">{subHeading}</div>
          ) : undefined}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => close()}>
            <FormattedMessage id="actions.cancel" />
          </Button>
          <Button
            type="submit"
            variant="primary"
            onClick={() => close(currentValue)}
            disabled={
              currentValue === undefined ||
              (validate && !validate(currentValue))
            }
          >
            <FormattedMessage id="actions.confirm" />
          </Button>
        </Modal.Footer>
      </Form>
    </Modal>
  );
};

export const SettingModal = PatchComponent(
  "SettingModal",
  _SettingModal
) as typeof _SettingModal;

interface IModalSetting<T> extends ISetting {
  value: T | undefined;
  buttonText?: string;
  buttonTextID?: string;
  onChange: (v: T) => void;
  renderField: (
    value: T | undefined,
    setValue: (v?: T) => void,
    error?: string
  ) => JSX.Element;
  renderValue?: (v: T | undefined) => JSX.Element;
  modalProps?: ModalProps;
  validateChange?: (v: T) => void | undefined;
}

export const _ModalSetting = <T extends {}>(props: IModalSetting<T>) => {
  const {
    id,
    className,
    value,
    headingID,
    heading,
    subHeadingID,
    subHeading,
    onChange,
    renderField,
    renderValue,
    tooltipID,
    buttonText,
    buttonTextID,
    modalProps,
    disabled,
    advanced,
    validateChange,
  } = props;
  const [showModal, setShowModal] = useState(false);
  const [error, setError] = useState<string>();
  const { advancedMode } = useSettings();

  function onClose(v: T | undefined) {
    setError(undefined);
    if (v !== undefined) {
      if (validateChange) {
        try {
          validateChange(v);
        } catch (e) {
          setError((e as Error).message);
          return;
        }
      }

      onChange(v);
    }
    setShowModal(false);
  }

  if (advanced && !advancedMode) return null;

  return (
    <>
      {showModal ? (
        <SettingModal<T>
          headingID={headingID}
          subHeadingID={subHeadingID}
          heading={heading}
          subHeading={subHeading}
          value={value}
          renderField={renderField}
          close={onClose}
          error={error}
          {...modalProps}
        />
      ) : undefined}

      <ChangeButtonSetting<T>
        id={id}
        className={className}
        disabled={disabled}
        buttonText={buttonText}
        buttonTextID={buttonTextID}
        headingID={headingID}
        heading={heading}
        tooltipID={tooltipID}
        subHeadingID={subHeadingID}
        subHeading={subHeading}
        value={value}
        onChange={() => setShowModal(true)}
        renderValue={renderValue}
      />
    </>
  );
};

export const ModalSetting = PatchComponent(
  "ModalSetting",
  _ModalSetting
) as typeof _ModalSetting;

interface IStringSetting extends ISetting {
  value: string | undefined;
  onChange: (v: string) => void;
}

export const StringSetting: React.FC<IStringSetting> = PatchComponent(
  "StringSetting",
  (props) => {
    return (
      <ModalSetting<string>
        {...props}
        renderField={(value, setValue) => (
          <Form.Control
            className="text-input"
            value={value ?? ""}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setValue(e.currentTarget.value)
            }
          />
        )}
        renderValue={(value) => <span>{value}</span>}
      />
    );
  }
);

interface INumberSetting extends ISetting {
  value: number | undefined;
  onChange: (v: number) => void;
}

export const NumberSetting: React.FC<INumberSetting> = PatchComponent(
  "NumberSetting",
  (props) => {
    return (
      <ModalSetting<number>
        {...props}
        renderField={(value, setValue) => (
          <NumberField
            className="text-input"
            value={value ?? 0}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setValue(Number.parseInt(e.currentTarget.value || "0", 10))
            }
          />
        )}
        renderValue={(value) => <span>{value}</span>}
      />
    );
  }
);

interface IStringListSetting extends ISetting {
  value: string[] | undefined;
  defaultNewValue?: string;
  onChange: (v: string[]) => void;
}

export const StringListSetting: React.FC<IStringListSetting> = PatchComponent(
  "StringListSetting",
  (props) => {
    return (
      <ModalSetting<string[]>
        {...props}
        renderField={(value, setValue) => (
          <StringListInput
            value={value ?? []}
            setValue={setValue}
            placeholder={props.defaultNewValue}
          />
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
  }
);

interface IConstantSetting<T> extends ISetting {
  value?: T;
  renderValue?: (v: T | undefined) => JSX.Element;
}

export const _ConstantSetting = <T extends {}>(props: IConstantSetting<T>) => {
  const { id, headingID, subHeading, subHeadingID, renderValue, value } = props;
  const intl = useIntl();

  return (
    <div className="setting" id={id}>
      <div>
        <h3>{headingID ? intl.formatMessage({ id: headingID }) : undefined}</h3>

        <div className="value">{renderValue ? renderValue(value) : value}</div>

        {subHeadingID ? (
          <div className="sub-heading">
            {intl.formatMessage({ id: subHeadingID })}
          </div>
        ) : subHeading ? (
          <div className="sub-heading">{subHeading}</div>
        ) : undefined}
      </div>
      <div />
    </div>
  );
};

export const ConstantSetting = PatchComponent(
  "ConstantSetting",
  _ConstantSetting
) as typeof _ConstantSetting;
