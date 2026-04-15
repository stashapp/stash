import { faBan, faMinus, faThumbsUp } from "@fortawesome/free-solid-svg-icons";
import React, { useState } from "react";
import { Button, ButtonGroup, Dropdown, DropdownButton } from "react-bootstrap";
import { useIntl } from "react-intl";
import { Icon } from "src/components/Shared/Icon";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { SweatDrops } from "src/components/Shared/SweatDrops";
import { useConfigurationContext } from "src/hooks/Config";

export interface IOCounterButtonProps {
  value: number;
  onIncrement: () => Promise<void>;
  onDecrement: () => Promise<void>;
  onReset: () => Promise<void>;
}

export const OCounterButton: React.FC<IOCounterButtonProps> = (
  props: IOCounterButtonProps
) => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();
  const { sfwContentMode } = configuration.interface;

  const icon = !sfwContentMode ? <SweatDrops /> : <Icon icon={faThumbsUp} />;
  const messageID = !sfwContentMode ? "o_count" : "o_count_sfw";

  const [loading, setLoading] = useState(false);

  async function increment() {
    setLoading(true);
    await props.onIncrement();
    setLoading(false);
  }

  async function decrement() {
    setLoading(true);
    await props.onDecrement();
    setLoading(false);
  }

  async function reset() {
    setLoading(true);
    await props.onReset();
    setLoading(false);
  }

  if (loading) return <LoadingIndicator message="" inline small />;

  const renderButton = () => (
    <Button
      className="minimal pr-1"
      onClick={increment}
      variant="secondary"
      title={intl.formatMessage({ id: messageID })}
    >
      {icon}
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
    <ButtonGroup className="o-counter">
      {renderButton()}
      {maybeRenderDropdown()}
    </ButtonGroup>
  );
};
