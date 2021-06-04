import React from "react";
import {
  Button,
  ButtonGroup,
  Dropdown,
  DropdownButton,
  Spinner,
} from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon, SweatDrops } from "src/components/Shared";

export interface IOCounterButtonProps {
  loading: boolean;
  value: number;
  onIncrement: () => void;
  onDecrement: () => void;
  onReset: () => void;
  onMenuOpened?: () => void;
  onMenuClosed?: () => void;
}

export const OCounterButton: React.FC<IOCounterButtonProps> = (
  props: IOCounterButtonProps
) => {
  const intl = useIntl();
  if (props.loading) return <Spinner animation="border" role="status" />;

  const renderButton = () => (
    <Button
      className="minimal pr-1"
      onClick={props.onIncrement}
      variant="secondary"
      title={intl.formatMessage({id: "o_counter"})}
    >
      <SweatDrops />
      <span className="ml-2">{props.value}</span>
    </Button>
  );

  const maybeRenderDropdown = () => {
    if (props.value) {
      return (
        <DropdownButton
          as={ButtonGroup}
          title=" "
          variant="secondary"
          className="pl-0 show-carat"
        >
          <Dropdown.Item onClick={props.onDecrement}>
            <Icon icon="minus" />
            <span>Decrement</span>
          </Dropdown.Item>
          <Dropdown.Item onClick={props.onReset}>
            <Icon icon="ban" />
            <span>Reset</span>
          </Dropdown.Item>
        </DropdownButton>
      );
    }
  };

  return (
    <ButtonGroup className="o-counter">
      {renderButton()}
      {maybeRenderDropdown()}
    </ButtonGroup>
  );
};
