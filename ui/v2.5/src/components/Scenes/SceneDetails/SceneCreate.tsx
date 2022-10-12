import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { SceneEditPanel } from "./SceneEditPanel";

const SceneCreate: React.FC = () => {
  const intl = useIntl();

  return (
    <div className="row new-view justify-content-center" id="create-scene-page">
      <div className="col-md-8">
        <h2>
          <FormattedMessage
            id="actions.create_entity"
            values={{ entityType: intl.formatMessage({ id: "scene" }) }}
          />
        </h2>
        <SceneEditPanel scene={{}} isVisible isNew />
      </div>
    </div>
  );
};

export default SceneCreate;
