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
import { IconProp } from "@fortawesome/fontawesome-svg-core";
import { Icon } from "../Shared";

interface IListFilterOperation {
  text: string;
  onClick: () => void;
  isDisplayed?: () => boolean;
  icon?: IconProp;
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
          icon: "pencil-alt",
          text: intl.formatMessage({ id: "actions.edit" }),
          onClick: onEdit,
        });
      }
      if (onDelete) {
        buttons.push({
          icon: "trash",
          text: intl.formatMessage({ id: "actions.delete" }),
          onClick: onDelete,
          buttonVariant: "danger",
        });
      }
    }

    if (buttons.length > 0) {
      return (
        <ButtonGroup className="ml-2 mb-1">
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
                  <Icon icon={button.icon as IconProp} />
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
      {maybeRenderButtons()}

      <div className="mx-2">{renderMore()}</div>
    </>
  );
};
