import React, { FunctionComponent } from "react";
import { Button, Popover, Menu, MenuItem } from "@blueprintjs/core";
import { Icons } from "../../utils/icons";

export interface IOCounterButtonProps {
  loading: boolean
  value: number
  onIncrement: () => void
  onDecrement: () => void
  onReset: () => void
  onMenuOpened?: () => void
  onMenuClosed?: () => void
}

export const OCounterButton: FunctionComponent<IOCounterButtonProps> = (props: IOCounterButtonProps) => {
  function renderButton() {
    return (
      <Button
        loading={props.loading}
        icon={Icons.sweatDrops()}
        text={props.value}
        minimal={true}
        onClick={props.onIncrement}
        disabled={props.loading}
      />
    );
  }

  if (props.value) {
    // just render the button by itself
    return (
      <Popover 
        interactionKind={"hover"} 
        hoverOpenDelay={1000} 
        position="bottom" 
        disabled={props.loading} 
        onOpening={props.onMenuOpened}
        onClosing={props.onMenuClosed}
      >
        {renderButton()}
        <Menu>
          <MenuItem text="Decrement" icon="minus" onClick={props.onDecrement}/>
          <MenuItem text="Reset" icon="disable" onClick={props.onReset}/>
        </Menu>
      </Popover>
    );
  } else {
    return renderButton();
  }
}