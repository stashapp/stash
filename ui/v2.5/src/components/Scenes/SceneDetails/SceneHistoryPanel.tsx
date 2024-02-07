import {
  faEllipsisV,
  faPlus,
  faTrash,
} from "@fortawesome/free-solid-svg-icons";
import React from "react";
import { Button, Dropdown } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Counter } from "src/components/Shared/Counter";
import { Icon } from "src/components/Shared/Icon";
import * as GQL from "src/core/generated-graphql";
import { TextField } from "src/utils/field";
import TextUtils from "src/utils/text";

const History: React.FC<{
  className?: string;
  history: string[];
  onRemove: (date: string) => void;
  noneID: string;
}> = ({ className, history, noneID, onRemove }) => {
  const intl = useIntl();

  if (history.length === 0) {
    return (
      <div>
        <FormattedMessage id={noneID} />
      </div>
    );
  }

  return (
    <div className="scene-history">
      <ul className={className}>
        {history.map((playdate, index) => (
          <li key={index}>
            <span>{TextUtils.formatDateTime(intl, playdate)}</span>
            <Button
              className="remove-date-button"
              size="sm"
              variant="minimal"
              onClick={() => onRemove(playdate)}
              title={intl.formatMessage({ id: "actions.remove_date" })}
            >
              <Icon icon={faTrash} />
            </Button>
          </li>
        ))}
      </ul>
    </div>
  );
};

const HistoryMenu: React.FC<{
  hasHistory: boolean;
  onAddDate: () => void;
  onClearDates: () => void;
}> = ({ hasHistory, onAddDate, onClearDates }) => {
  const intl = useIntl();

  return (
    <Dropdown className="history-operations-dropdown">
      <Dropdown.Toggle
        variant="secondary"
        className="minimal"
        title={intl.formatMessage({ id: "operations" })}
      >
        <Icon icon={faEllipsisV} />
      </Dropdown.Toggle>
      <Dropdown.Menu className="bg-secondary text-white">
        <Dropdown.Item
          key="generate"
          className="bg-secondary text-white"
          onClick={() => onAddDate()}
        >
          <FormattedMessage id="actions.add_manual_date" />
        </Dropdown.Item>
        {hasHistory && (
          <Dropdown.Item
            key="generate"
            className="bg-secondary text-white"
            onClick={() => onClearDates()}
          >
            <FormattedMessage id="actions.clear_date_data" />
          </Dropdown.Item>
        )}
      </Dropdown.Menu>
    </Dropdown>
  );
};

interface ISceneHistoryProps {
  scene: GQL.SceneDataFragment;
}

export const SceneHistoryPanel: React.FC<ISceneHistoryProps> = ({ scene }) => {
  const intl = useIntl();

  const playHistory = (scene.play_history ?? []).filter(
    (h) => h != null
  ) as string[];
  const oHistory = (scene.o_history ?? []).filter((h) => h != null) as string[];

  return (
    <div>
      <div className="play-history">
        <div className="history-header">
          <h5>
            <span>
              <FormattedMessage id="play_history" />
              <Counter count={playHistory.length} hideZero />
            </span>
            <span>
              <Button
                size="sm"
                variant="minimal"
                className="add-date-button"
                title={intl.formatMessage({ id: "actions.add_play" })}
              >
                <Icon icon={faPlus} />
              </Button>
              <HistoryMenu
                hasHistory={playHistory.length > 0}
                onAddDate={() => {}}
                onClearDates={() => {}}
              />
            </span>
          </h5>
        </div>

        <History
          history={playHistory ?? []}
          noneID="playdate_recorded_no"
          onRemove={() => {}}
        />
        <dl className="details-list">
          <TextField
            id="media_info.play_duration"
            value={TextUtils.secondsToTimestamp(scene.play_duration ?? 0)}
          />
        </dl>
      </div>

      <div className="o-history">
        <div className="history-header">
          <h5>
            <span>
              <FormattedMessage id="o_history" />
              <Counter count={oHistory.length} hideZero />
            </span>
            <span>
              <Button
                size="sm"
                variant="minimal"
                className="add-date-button"
                title={intl.formatMessage({ id: "actions.add_o" })}
              >
                <Icon icon={faPlus} />
              </Button>
              <HistoryMenu
                hasHistory={oHistory.length > 0}
                onAddDate={() => {}}
                onClearDates={() => {}}
              />
            </span>
          </h5>
        </div>
        <History
          history={oHistory}
          noneID="odate_recorded_no"
          onRemove={() => {}}
        />
      </div>
    </div>
  );
};

export default SceneHistoryPanel;
