import React, { useEffect } from "react";
import {
  Button,
  ButtonGroup,
  Dropdown,
  OverlayTrigger,
  Tooltip,
} from "react-bootstrap";
import Mousetrap from "mousetrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "../Shared";

interface IListFilterOperation {
  text: string;
  onClick: () => void;
  isDisplayed?: () => boolean;
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

    if (itemsSelected) {
      Mousetrap.bind("e", () => {
        onEdit?.();
      });

      Mousetrap.bind("d d", () => {
        onDelete?.();
      });
    }

    return () => {
      Mousetrap.unbind("s a");
      Mousetrap.unbind("s n");

      if (itemsSelected) {
        Mousetrap.unbind("e");
        Mousetrap.unbind("d d");
      }
    };
  });

  function maybeRenderSelectedButtons() {
    if (itemsSelected && (onEdit || onDelete)) {
      return (
        <ButtonGroup className="ml-2 mb-1">
          {onEdit && (
            <OverlayTrigger
              overlay={
                <Tooltip id="edit">
                  {intl.formatMessage({ id: "actions.edit" })}
                </Tooltip>
              }
            >
              <Button variant="secondary" onClick={onEdit}>
                <Icon icon="pencil-alt" />
              </Button>
            </OverlayTrigger>
          )}

          {onDelete && (
            <OverlayTrigger
              overlay={
                <Tooltip id="delete">
                  {intl.formatMessage({ id: "actions.delete" })}
                </Tooltip>
              }
            >
              <Button variant="danger" onClick={onDelete}>
                <Icon icon="trash" />
              </Button>
            </OverlayTrigger>
          )}
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

    if (options.length > 0) {
      return (
        <Dropdown className="mb-1">
          <Dropdown.Toggle variant="secondary" id="more-menu">
            <Icon icon="ellipsis-h" />
          </Dropdown.Toggle>
          <Dropdown.Menu className="bg-secondary text-white">
            {options}
          </Dropdown.Menu>
        </Dropdown>
      );
    }
  }

  return (
    <>
      {maybeRenderSelectedButtons()}

      <div className="mx-2">{renderMore()}</div>
    </>
  );
};
