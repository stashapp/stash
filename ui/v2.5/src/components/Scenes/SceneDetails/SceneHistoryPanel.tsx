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

  const file = useMemo(
    () => (props.scene.files.length > 0 ? props.scene.files[0] : undefined),
    [props.scene]
  );

  return (
    <>
      <div className="row">
        <div className="col-12">
          <h5>
            <FormattedMessage id="file_history" />{" "}
          </h5>
          <h6>
            <FormattedMessage id="created_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.scene.created_at)}{" "}
          </h6>
          <h6>
            <FormattedMessage id="updated_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.scene.updated_at)}{" "}
          </h6>
          {/* <h6>
            <TextField id="file_mod_time">
            <FormattedTime
              dateStyle="medium"
              timeStyle="medium"
              value={props.file.mod_time ?? 0}
            />
          </TextField>
          </h6> */}
        </div>
      </div>
      <div className="row">
        <div className="col-12">
        </div>
      </div>
    </>
  );
};

export default SceneHistoryPanel;
