import React, { PropsWithChildren, useEffect } from "react";
import {
  Button,
  ButtonGroup,
  Dropdown,
  OverlayTrigger,
  Tooltip,
} from "react-bootstrap";
import Mousetrap from "mousetrap";
import { FormattedMessage, useIntl } from "react-intl";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { Icon } from "../Shared/Icon";
import {
  faEllipsisH,
  faPencil,
  faPencilAlt,
  faPlay,
  faPlus,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import cx from "classnames";

export const OperationDropdown: React.FC<
  PropsWithChildren<{
    className?: string;
  }>
> = ({ className, children }) => {
  if (!children) return null;

  return (
    <Dropdown className={className} as={ButtonGroup}>
      <Dropdown.Toggle variant="secondary" id="more-menu">
        <Icon icon={faEllipsisH} />
      </Dropdown.Toggle>
      <Dropdown.Menu className="bg-secondary text-white">
        {children}
      </Dropdown.Menu>
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

  function maybeRenderButtons() {
    const buttons = (otherOperations ?? []).filter((o) => {
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
        buttons.push({
          icon: faPencilAlt,
          text: intl.formatMessage({ id: "actions.edit" }),
          onClick: onEdit,
        });
      }
      if (onDelete) {
        buttons.push({
          icon: faTrash,
          text: intl.formatMessage({ id: "actions.delete" }),
          onClick: onDelete,
          buttonVariant: "danger",
        });
      }
    }

    if (buttons.length > 0) {
      return (
        <ButtonGroup className="ml-2">
          {buttons.map((button) => {
            return (
              <OverlayTrigger
                overlay={<Tooltip id="edit">{button.text}</Tooltip>}
                key={button.text}
              >
                <Button
                  variant={button.buttonVariant ?? "secondary"}
                  onClick={button.onClick}
                >
                  {button.icon ? <Icon icon={button.icon} /> : undefined}
                </Button>
              </OverlayTrigger>
            );
          })}
        </ButtonGroup>
      );
    }
  }

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

  function renderMore() {
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
  }

  return (
    <>
      {maybeRenderButtons()}

      <ButtonGroup className="ml-2">{renderMore()}</ButtonGroup>
    </>
  );
};

interface IListOperations {
  text: string;
  onClick: () => void;
  isDisplayed?: () => boolean;
  className?: string;
}

export const ListOperations: React.FC<{
  items: number;
  hasSelection?: boolean;
  operations?: IListOperations[];
  onEdit?: () => void;
  onDelete?: () => void;
  onPlay?: () => void;
  onCreateNew?: () => void;
  entityType?: string;
  operationsClassName?: string;
}> = ({
  items,
  hasSelection = false,
  operations = [],
  onEdit,
  onDelete,
  onPlay,
  onCreateNew,
  entityType,
  operationsClassName,
}) => {
  const intl = useIntl();

  return (
    <div>
      <ButtonGroup>
        {!!items && onPlay && (
          <Button
            className="play-button"
            variant="secondary"
            onClick={() => onPlay()}
            title={intl.formatMessage({ id: "actions.play" })}
          >
            <Icon icon={faPlay} />
          </Button>
        )}
        {!hasSelection && onCreateNew && (
          <Button
            className="create-new-button"
            variant="secondary"
            onClick={() => onCreateNew()}
            title={intl.formatMessage(
              { id: "actions.create_entity" },
              { entityType }
            )}
          >
            <Icon icon={faPlus} />
          </Button>
        )}

        {hasSelection && (onEdit || onDelete) && (
          <>
            {onEdit && (
              <Button variant="secondary" onClick={() => onEdit()}>
                <Icon icon={faPencil} />
              </Button>
            )}
            {onDelete && (
              <Button
                variant="danger"
                className="btn-danger-minimal"
                onClick={() => onDelete()}
              >
                <Icon icon={faTrash} />
              </Button>
            )}
          </>
        )}

        {operations.length > 0 && (
          <OperationDropdown className={operationsClassName}>
            {operations.map((o) => {
              if (o.isDisplayed && !o.isDisplayed()) {
                return null;
              }

              return (
                <OperationDropdownItem
                  key={o.text}
                  onClick={o.onClick}
                  text={o.text}
                  className={o.className}
                />
              );
            })}
          </OperationDropdown>
        )}
      </ButtonGroup>
    </div>
  );
};
