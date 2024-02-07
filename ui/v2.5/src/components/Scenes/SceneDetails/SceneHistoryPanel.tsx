import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Counter } from "src/components/Shared/Counter";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";

const History: React.FC<{
  className?: string;
  history: string[];
  noneID: string;
}> = ({ className, history, noneID }) => {
  const intl = useIntl();

  if (history.length === 0) {
    return <FormattedMessage id={noneID} />;
  }

  return (
    <ul className={className}>
      {history.map((playdate, index) => (
        <li key={index}>{TextUtils.formatDateTime(intl, playdate)}</li>
      ))}
    </ul>
  );
};

interface ISceneHistoryProps {
  scene: GQL.SceneDataFragment;
}

export const SceneHistoryPanel: React.FC<ISceneHistoryProps> = ({ scene }) => {
  const playHistory = (scene.play_history ?? []).filter(
    (h) => h != null
  ) as string[];
  const oHistory = (scene.o_history ?? []).filter((h) => h != null) as string[];

  return (
    <div>
      <div className="play-history">
        <h5>
          <FormattedMessage id="play_history" />
          <Counter count={playHistory.length} hideZero hideOne />
        </h5>
        <History history={playHistory ?? []} noneID="playdate_recorded_no" />
        <h6>
          <FormattedMessage id="media_info.play_duration" />
          :&nbsp;
          {TextUtils.secondsToTimestamp(scene.play_duration ?? 0)}
        </h6>
      </div>
      <div className="o-history">
        <h5>
          <FormattedMessage id="o_history" />
          <Counter count={oHistory.length} hideZero />
        </h5>
        <History history={oHistory} noneID="odate_recorded_no" />
      </div>
    </div>
  );
};

export default SceneHistoryPanel;
