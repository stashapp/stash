import { faEye, faThumbsUp } from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button, ButtonGroup } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { SweatDrops } from "./SweatDrops";
import cx from "classnames";
import { useIntl } from "react-intl";
import { useConfigurationContext } from "src/hooks/Config";

interface ICountButtonProps {
  value: number;
  icon: React.ReactNode;
  onIncrement?: () => void;
  onValueClicked?: () => void;
  title?: string;
  countTitle?: string;
}

export const CountButton: React.FC<ICountButtonProps> = ({
  value,
  icon,
  onIncrement,
  onValueClicked,
  title,
  countTitle,
}) => {
  return (
    <ButtonGroup
      className={cx("count-button", { "increment-only": !onValueClicked })}
    >
      <Button
        className="minimal count-icon"
        variant="secondary"
        onClick={() => onIncrement?.()}
        title={title}
      >
        {icon}
      </Button>
      <Button
        className="minimal count-value"
        variant="secondary"
        onClick={() => (onValueClicked ?? onIncrement)?.()}
        title={!!onValueClicked ? countTitle : undefined}
      >
        <span>{value}</span>
      </Button>
    </ButtonGroup>
  );
};

type CountButtonPropsNoIcon = Omit<ICountButtonProps, "icon">;

export const ViewCountButton: React.FC<CountButtonPropsNoIcon> = (props) => {
  const intl = useIntl();
  return (
    <CountButton
      {...props}
      icon={<Icon icon={faEye} />}
      title={intl.formatMessage({ id: "media_info.play_count" })}
      countTitle={intl.formatMessage({ id: "actions.view_history" })}
    />
  );
};

export const OCounterButton: React.FC<CountButtonPropsNoIcon> = (props) => {
  const intl = useIntl();
  const { configuration } = useConfigurationContext();
  const { sfwContentMode } = configuration.interface;

  const icon = !sfwContentMode ? <SweatDrops /> : <Icon icon={faThumbsUp} />;
  const messageID = !sfwContentMode ? "o_count" : "o_count_sfw";

  return (
    <CountButton
      {...props}
      icon={icon}
      title={intl.formatMessage({ id: messageID })}
      countTitle={intl.formatMessage({ id: "actions.view_history" })}
    />
  );
};
