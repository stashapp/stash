import React from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { useLocation } from "react-router-dom";
import { GalleryEditPanel } from "./GalleryEditPanel";

const GalleryCreate: React.FC = () => {
  const intl = useIntl();

  function useQuery() {
    const { search } = useLocation();
    return React.useMemo(() => new URLSearchParams(search), [search]);
  }

  const query = useQuery();
  const nameQuery = query.get("name");

  return (
    <div className="row new-view">
      <div className="col-md-6">
        <h2>
          <FormattedMessage
            id="actions.create_entity"
            values={{ entityType: intl.formatMessage({ id: "gallery" }) }}
          />
        </h2>
        <GalleryEditPanel
          isNew
          gallery={{ title: nameQuery ?? "" }}
          isVisible
          onDelete={() => {}}
        />
      </div>
    </div>
  );
};

export default GalleryCreate;
