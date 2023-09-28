import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";

interface ISceneHistoryProps {
  scene: GQL.SceneDataFragment;
}

export const SceneHistoryPanel: React.FC<ISceneHistoryProps> = (props) => {
  const intl = useIntl();

  return (
    <>     
      {/* Could replace these play/ocount displays with data from the scenes_playdates/odates table once accessible in GraphQL */}
      {/* Could also check then if the recorded dates are the same as scene.created_at then that they are 'Estimated Play Date' */}
      <div className="row">
        <div className="col-12">
          <h5>
            <FormattedMessage id="play_history" />{" "}
          </h5>
          {props.scene.play_count != null && props.scene.play_count !== 0 ? (
            <>
              {Array.from({ length: props.scene.play_count - 1 }).map(
                (_, index) => (
                  <h6 key={index}>
                    <FormattedMessage id="playdate_recorded" />
                  </h6>
                )
              )}
              {props.scene.last_played_at && (
                <h6>
                  <FormattedMessage id="playdate_recorded" />
                  {": "}
                  {TextUtils.formatDateTime(intl, props.scene.last_played_at)}
                </h6>
              )}
            </>
          ) : props.scene.play_count === 0 ? (
            <h6>
              <FormattedMessage id="playdate_recorded_no" />
            </h6>
          ) : (
            <h6>N/A</h6>
          )}
          <h6>
            {/* Could make this a toggle if Track Activity is off*/}
            <FormattedMessage id="media_info.play_duration" />
            :&nbsp;
            {TextUtils.secondsToTimestamp(props.scene.play_duration ?? 0)}
          </h6>
        </div>
      </div>
      <div className="row">
        <div className="col-12">
          <h5>
            <FormattedMessage id="o_history" />{" "}
          </h5>
          {props.scene.o_counter != null && props.scene.o_counter !== 0 ? (
            <div className="date-row">
              {Array.from({ length: props.scene.o_counter }).map((_, index) => (
                <h6 key={index}>
                  <FormattedMessage id="odate_recorded" />
                </h6>
              ))}
            </div>
          ) : props.scene.o_counter === 0 ? (
            <h6>
              <FormattedMessage id="odate_recorded_no" />
            </h6>
          ) : (
            <h6>N/A</h6>
          )}
        </div>
      </div>
    </>
  );
};

export default SceneHistoryPanel;
