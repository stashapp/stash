import React, { PropsWithChildren, useEffect, useMemo } from "react";
import { Button, ButtonGroup, Dropdown } from "react-bootstrap";
import Mousetrap from "mousetrap";
import { FormattedMessage, useIntl } from "react-intl";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { Icon } from "../Shared/Icon";
import {
  faEllipsisH,
  faPencilAlt,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";
import { createPortal } from "react-dom";

export const OperationDropdown: React.FC<
  PropsWithChildren<{
    className?: string;
    menuPortalTarget?: HTMLElement;
    menuClassName?: string;
  }>
> = ({ className, menuPortalTarget, menuClassName, children }) => {
  if (!children) return null;

  const menu = (
    <Dropdown.Menu className={cx("bg-secondary text-white", menuClassName)}>
      {children}
    </Dropdown.Menu>
  );

  return (
    <Dropdown className={className} as={ButtonGroup}>
      <Dropdown.Toggle variant="secondary" id="more-menu">
        <Icon icon={faEllipsisH} />
      </Dropdown.Toggle>
      {menuPortalTarget ? createPortal(menu, menuPortalTarget) : menu}
    </Dropdown>
  );
};

export const OperationDropdownItem: React.FC<{
  text: string;
  onClick: () => void;
  className?: string;
}> = ({ text, onClick, className }) => {
  return (
    <Dropdown.Item
      className={cx("bg-secondary text-white", className)}
      onClick={onClick}
    >
      {text}
    </Dropdown.Item>
  );
};

export interface IListFilterOperation {
  text: string;
  onClick: () => void;
  isDisplayed?: () => boolean;
  icon?: IconDefinition;
  buttonVariant?: string;
}

interface IListOperationButtonsProps {
  onSelectAll?: () => void;
  onSelectNone?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  itemsSelected?: boolean;
  otherOperations?: IListFilterOperation[];
}

export const ListOperationButtons: React.FC<IListOperationButtonsProps> = ({
  onSelectAll,
  onSelectNone,
  onEdit,
  onDelete,
  itemsSelected,
  otherOperations,
}) => {
  const intl = useIntl();

  useEffect(() => {
    Mousetrap.bind("s a", () => onSelectAll?.());
    Mousetrap.bind("s n", () => onSelectNone?.());

    Mousetrap.bind("e", () => {
      if (itemsSelected) {
        onEdit?.();
      }
    });

    Mousetrap.bind("d d", () => {
      if (itemsSelected) {
        onDelete?.();
      }
    });

    return () => {
      Mousetrap.unbind("s a");
      Mousetrap.unbind("s n");
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

  const buttons = useMemo(() => {
    const ret = (otherOperations ?? []).filter((o) => {
      if (!o.icon) {
        return false;
      }

      if (!o.isDisplayed) {
        return true;
      }

      return o.isDisplayed();
    });

    if (itemsSelected) {
      if (onEdit) {
        ret.push({
          icon: faPencilAlt,
          text: intl.formatMessage({ id: "actions.edit" }),
          onClick: onEdit,
        });
      }
      if (onDelete) {
        ret.push({
          icon: faTrash,
          text: intl.formatMessage({ id: "actions.delete" }),
          onClick: onDelete,
          buttonVariant: "danger",
        });
      }
    }

    return ret;
  }, [otherOperations, itemsSelected, onEdit, onDelete, intl]);

  const operationButtons = useMemo(() => {
    return (
      <>
        {buttons.map((button) => {
          return (
            <Button
              key={button.text}
              variant={button.buttonVariant ?? "secondary"}
              onClick={button.onClick}
              title={button.text}
            >
              <Icon icon={button.icon!} />
            </Button>
          );
        })}
      </>
    );
  }, [buttons]);

  const moreDropdown = useMemo(() => {
    function renderSelectAll() {
      if (onSelectAll) {
        return (
          <Dropdown.Item
            key="select-all"
            className="bg-secondary text-white"
            onClick={() => onSelectAll?.()}
          >
            <FormattedMessage id="actions.select_all" />
          </Dropdown.Item>
        );
      }
    }

    function renderSelectNone() {
      if (onSelectNone) {
        return (
          <Dropdown.Item
            key="select-none"
            className="bg-secondary text-white"
            onClick={() => onSelectNone?.()}
          >
            <FormattedMessage id="actions.select_none" />
          </Dropdown.Item>
        );
      }
    }

    const options = [renderSelectAll(), renderSelectNone()].filter((o) => o);

    if (otherOperations) {
      otherOperations
        .filter((o) => {
          // buttons with icons are rendered in the button group
          if (o.icon) {
            return false;
          }

          if (!o.isDisplayed) {
            return true;
          }

          return o.isDisplayed();
        })
        .forEach((o) => {
          options.push(
            <Dropdown.Item
              key={o.text}
              className="bg-secondary text-white"
              onClick={o.onClick}
            >
              {o.text}
            </Dropdown.Item>
          );
        });
    }

    return (
      <OperationDropdown>
        {options.length > 0 ? options : undefined}
      </OperationDropdown>
    );
  }, [otherOperations, onSelectAll, onSelectNone]);

  // don't render anything if there are no buttons or operations
  if (buttons.length === 0 && !moreDropdown) {
    return null;
  }

  return (
    <>
      <ButtonGroup>
        {operationButtons}
        {moreDropdown}
      </ButtonGroup>
    </>
  );
};
