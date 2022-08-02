import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { LoadingIndicator } from "src/components/Shared";
import SceneEditPanel from "./SceneEditPanel";
import * as GQL from "src/core/generated-graphql";

const SceneCreate: React.FC = () => {
  const intl = useIntl();

  return (
    <div className="row new-view" id="scene-page">
      {/* <div className="performer-image-container col-md-4 text-center">
        {renderPerformerImage()}
      </div> */}
      <div className="col-md-8">
        <h2>
          <FormattedMessage
            id="actions.create_entity"
            values={{ entityType: intl.formatMessage({ id: "scene" }) }}
          />
        </h2>
        <SceneEditPanel
          scene={{}}
          isVisible
          isNew
        />
      </div>
    </div>
  );
};

export default SceneCreate;
