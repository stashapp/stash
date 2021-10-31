import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { GalleryEditPanel } from "./GalleryEditPanel";

const GalleryCreate: React.FC = () => {
  const intl = useIntl();

  return (
    <div className="row new-view">
      <div className="col-md-6">
        <h2>
          <FormattedMessage
            id="actions.create_entity"
            values={{
              entityType: intl.formatMessage(
                { id: "countables.galleries" },
                { count: 1 }
              ),
            }}
          />
        </h2>
        <GalleryEditPanel
          isNew
          gallery={undefined}
          isVisible
          onDelete={() => {}}
        />
      </div>
    </div>
  );
};

export default GalleryCreate;
