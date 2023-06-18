import { faBan, faMinus } from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, ButtonGroup, Dropdown, DropdownButton } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "src/components/Shared/Icon";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { EyeBall } from "src/components/Shared/EyeBall";

export interface IPlayCounterButtonProps {
  value: number;
  onPlayIncrement: () => Promise<void>;
  onPlayDecrement: () => Promise<void>;
  onPlayReset: () => Promise<void>;
}

export const PlayCounterButton: React.FC<IPlayCounterButtonProps> = (
  props: IPlayCounterButtonProps
) => {
  const intl = useIntl();
  const [loading, setLoading] = useState(false);

  async function increment() {
    setLoading(true);
    await props.onPlayIncrement();
    setLoading(false);
  }

  async function decrement() {
    setLoading(true);
    await props.onPlayDecrement();
    setLoading(false);
  }

  async function reset() {
    setLoading(true);
    await props.onPlayReset();
    setLoading(false);
  }

  if (loading) return <LoadingIndicator message="" inline small />;

  const renderButton = () => (
    <Button
      className="minimal pr-1"
      onClick={increment}
      variant="secondary"
      title={intl.formatMessage({ id: "play_counter" })}
    >
      <EyeBall />
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
          <Dropdown.Item onClick={decrement}>
            <Icon icon={faMinus} />
            <span>Decrement</span>
          </Dropdown.Item>
          <Dropdown.Item onClick={reset}>
            <Icon icon={faBan} />
            <span>Reset</span>
          </Dropdown.Item>
        </DropdownButton>
      );
    }
  };

  return (
    <ButtonGroup className="play-counter">
      {renderButton()}
      {maybeRenderDropdown()}
    </ButtonGroup>
  );
};
