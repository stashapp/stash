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
      <div className="row">
        <div className="col-12">
          <h5>
            <FormattedMessage id="history" />{" "}
          </h5>
          <div className="date-row">
            <h6>
              <FormattedMessage id="release_date" />:
              {props.scene.date && (
                <>
                  &nbsp;
                  {TextUtils.formatDateTime(intl, props.scene.date)}
                </>
              )}
            </h6>
            <h6>
              <FormattedMessage id="file_mod_time" />:
              {props.scene.files[0].mod_time && (
                <>
                  &nbsp;
                  {TextUtils.formatDateTime(
                    intl,
                    props.scene.files[0].mod_time
                  )}
                </>
              )}
            </h6>
          </div>
        </div>
      </div>
      <div className="row">
        <div className="col-12">
          <h5>
            <FormattedMessage id="scene_history" />{" "}
          </h5>
          <div className="date-row">
            <h6>
              <FormattedMessage id="created_at" />:
              {props.scene.created_at && (
                <>
                  &nbsp;
                  {TextUtils.formatDateTime(intl, props.scene.created_at)}
                </>
              )}
            </h6>
            <h6>
              <FormattedMessage id="updated_at" />:
              {props.scene.updated_at && (
                <>
                  &nbsp;
                  {TextUtils.formatDateTime(intl, props.scene.updated_at)}
                </>
              )}
            </h6>
          </div>
        </div>
      </div>
      {/* Could replace these play/ocount displays with data from the scenes_playdates/odates table later*/}
      <div className="row">
        <div className="col-12">
          <h5>
            <FormattedMessage id="play_history" />{" "}
          </h5>            
            {props.scene.play_count != null && props.scene.play_count !== 0 ? (
                <>
                  {Array.from({ length: props.scene.play_count - 1 }).map((_, index) => (
                    <h6 key={index}>Play Date Recorded</h6>
                  ))}
                  {props.scene.last_played_at && (
                    <h6>
                      Play Date Recorded:{" "}
                      {TextUtils.formatDateTime(intl, props.scene.last_played_at)}
                    </h6>
                  )}
                </>
              ) : props.scene.play_count === 0 ? (
                <h6>No Play Dates Recorded</h6>
              ) : (
                <h6>N/A</h6>
            )}
          <h6>
            {/* Could make this a toggle if Track Activity (Automatically) is off*/}
            <FormattedMessage id="media_info.play_duration" />:&nbsp;
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
              {Array.from({ length: props.scene.o_counter }).map(
                (_, index) => (
                  <h6 key={index}>O Date Recorded</h6>
                )
              )}
            </div>
          ) : props.scene.o_counter === 0 ? (
            <h6>No O Dates Recorded</h6>
          ) : (
            <h6>N/A</h6>
          )}
        </div>
      </div>
    </>
  );
};

export default SceneHistoryPanel;