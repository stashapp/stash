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
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";
import { Icon } from "../Shared/Icon";
import {
  faEllipsisH,
  faPencilAlt,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";

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
        <ButtonGroup className="ml-2 mb-2">
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
        <Dropdown>
          <Dropdown.Toggle variant="secondary" id="more-menu">
            <Icon icon={faEllipsisH} />
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

      <div className="mx-2 mb-2">{renderMore()}</div>
    </>
  );
};
