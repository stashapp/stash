import React, { PropsWithChildren, useEffect, useMemo } from "react";
import { Button, ButtonGroup, Dropdown } from "react-bootstrap";
import Mousetrap from "mousetrap";
import { FormattedMessage, useIntl } from "react-intl";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { Icon } from "../Shared/Icon";
import {
  faEllipsisH,
  faPencil,
  faPencilAlt,
  faPlay,
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
  className?: string;
}

interface IListOperationButtonsProps {
  onSelectAll?: () => void;
  onSelectNone?: () => void;
  onInvertSelection?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  itemsSelected?: boolean;
  otherOperations?: IListFilterOperation[];
}

export const ListOperationButtons: React.FC<IListOperationButtonsProps> = ({
  onSelectAll,
  onSelectNone,
  onInvertSelection,
  onEdit,
  onDelete,
  itemsSelected,
  otherOperations,
}) => {
  const intl = useIntl();

  useEffect(() => {
    Mousetrap.bind("s a", () => onSelectAll?.());
    Mousetrap.bind("s n", () => onSelectNone?.());
    Mousetrap.bind("s i", () => onInvertSelection?.());

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
      Mousetrap.unbind("s i");
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  }, [
    onSelectAll,
    onSelectNone,
    onInvertSelection,
    itemsSelected,
    onEdit,
    onDelete,
  ]);

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

    function renderInvertSelection() {
      if (onInvertSelection) {
        return (
          <Dropdown.Item
            key="invert-selection"
            className="bg-secondary text-white"
            onClick={() => onInvertSelection?.()}
          >
            <FormattedMessage id="actions.invert_selection" />
          </Dropdown.Item>
        );
      }
    }

    const options = [
      renderSelectAll(),
      renderSelectNone(),
      renderInvertSelection(),
    ].filter((o) => o);

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
  }, [otherOperations, onSelectAll, onSelectNone, onInvertSelection]);

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

export const ListOperations: React.FC<{
  items: number;
  hasSelection?: boolean;
  operations?: IListFilterOperation[];
  onEdit?: () => void;
  onDelete?: () => void;
  onPlay?: () => void;
  operationsClassName?: string;
  operationsMenuClassName?: string;
}> = ({
  items,
  hasSelection = false,
  operations = [],
  onEdit,
  onDelete,
  onPlay,
  operationsClassName = "list-operations",
  operationsMenuClassName,
}) => {
  const intl = useIntl();

  const dropdownOperations = useMemo(() => {
    return operations.filter((o) => {
      if (o.icon) {
        return false;
      }

      if (!o.isDisplayed) {
        return true;
      }

      return o.isDisplayed();
    });
  }, [operations]);

  const buttons = useMemo(() => {
    const otherButtons = (operations ?? []).filter((o) => {
      if (!o.icon) {
        return false;
      }

      if (!o.isDisplayed) {
        return true;
      }

      return o.isDisplayed();
    });

    const ret: React.ReactNode[] = [];

    function addButton(b: React.ReactNode | null) {
      if (b) {
        ret.push(b);
      }
    }

    const playButton =
      !!items && onPlay ? (
        <Button
          className="play-button"
          variant="secondary"
          onClick={() => onPlay()}
          title={intl.formatMessage({ id: "actions.play" })}
        >
          <Icon icon={faPlay} />
        </Button>
      ) : null;

    const editButton =
      hasSelection && onEdit ? (
        <Button
          className="edit-existing-button"
          variant="secondary"
          onClick={() => onEdit()}
        >
          <Icon icon={faPencil} />
        </Button>
      ) : null;

    const deleteButton =
      hasSelection && onDelete ? (
        <Button
          variant="danger"
          className="delete-button btn-danger-minimal"
          onClick={() => onDelete()}
        >
          <Icon icon={faTrash} />
        </Button>
      ) : null;

    addButton(playButton);
    addButton(editButton);
    addButton(deleteButton);

    otherButtons.forEach((button) => {
      addButton(
        <Button
          key={button.text}
          variant={button.buttonVariant ?? "secondary"}
          onClick={button.onClick}
          title={button.text}
          className={button.className}
        >
          <Icon icon={button.icon!} />
        </Button>
      );
    });

    if (ret.length === 0) {
      return null;
    }

    return ret;
  }, [operations, hasSelection, onDelete, onEdit, onPlay, items, intl]);

  if (dropdownOperations.length === 0 && !buttons) {
    return null;
  }

  return (
    <div className="list-operations">
      <ButtonGroup>
        {buttons}

        {dropdownOperations.length > 0 && (
          <OperationDropdown
            className={operationsClassName}
            menuClassName={operationsMenuClassName}
            menuPortalTarget={document.body}
          >
            {dropdownOperations.map((o) => (
              <OperationDropdownItem
                key={o.text}
                onClick={o.onClick}
                text={o.text}
                className={o.className}
              />
            ))}
          </OperationDropdown>
        )}
      </ButtonGroup>
    </div>
  );
};
