import React from "react";
import { Button, Spinner } from "react-bootstrap";
import { Icon, HoverPopover, SweatDrops } from "src/components/Shared";

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
  if (props.loading) return <Spinner animation="border" role="status" />;

  const renderButton = () => (
    <Button
      className="minimal"
      onClick={props.onIncrement}
      variant="secondary"
      title="O-Counter"
    >
      <SweatDrops />
      <span className="ml-2">{props.value}</span>
    </Button>
  );

  if (props.value) {
    return (
      <HoverPopover
        content={
          <div>
            <div>
              <Button
                className="minimal"
                onClick={props.onDecrement}
                variant="secondary"
              >
                <Icon icon="minus" />
                <span>Decrement</span>
              </Button>
            </div>
            <div>
              <Button
                className="minimal"
                onClick={props.onReset}
                variant="secondary"
              >
                <Icon icon="ban" />
                <span>Reset</span>
              </Button>
            </div>
          </div>
        }
        enterDelay={1000}
        placement="bottom"
        onOpen={props.onMenuOpened}
        onClose={props.onMenuClosed}
      >
        {renderButton()}
      </HoverPopover>
    );
  }
  return renderButton();
};
