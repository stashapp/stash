import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import { FormattedTime, FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { TextField, URLField } from "src/utils/field";

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
          <h6>
            <FormattedMessage id="release_date" />:&nbsp;{" "}
            {props.scene.date ? TextUtils.formatDateTime(intl, props.scene.date) : ""}
          </h6>
          <h6>
            <FormattedMessage id="file_mod_time" />:&nbsp;{" "}
            {TextUtils.formatDateTime(intl, props.scene.files[0].mod_time)}{" "}
          </h6>          
        </div>
      </div>
      <div className="row">
        <div className="col-12">
          <h5>
            <FormattedMessage id="scene_history" />{" "}
          </h5>
          <h6>
            <FormattedMessage id="created_at" />:&nbsp;{" "}
            {TextUtils.formatDateTime(intl, props.scene.created_at)}{" "}
          </h6>
          <h6>
            <FormattedMessage id="updated_at" />:&nbsp;{" "}
            {TextUtils.formatDateTime(intl, props.scene.updated_at)}{" "}
          </h6>    
        </div>
      </div>
      <div className="row">
        <div className="col-12">
        <h5>
            <FormattedMessage id="play_history" />{" "}
          </h5>
          <h6>
            <FormattedMessage id="media_info.play_duration" />:&nbsp;{TextUtils.secondsToTimestamp(props.scene.play_duration ?? 0)}
          </h6>
          <h6>
            <FormattedMessage id="media_info.play_count" />:&nbsp;{(props.scene.play_count ?? 0).toString()}
          </h6>
        </div>
      </div>
      <div className="row">
        <div className="col-12">
        <h5>
            <FormattedMessage id="o_history" />{" "}
          </h5>
          <h6>
            <FormattedMessage id="media_info.o_count" />:&nbsp;{(props.scene.o_counter ?? 0).toString()}
          </h6>
        </div>
      </div>
    </>
  );
};

export default SceneHistoryPanel;
