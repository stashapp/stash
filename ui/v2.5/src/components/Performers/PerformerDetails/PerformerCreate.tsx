import React, { useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { PerformerEditPanel } from "./PerformerEditPanel";
import { useLocation } from "react-router-dom";

const PerformerCreate: React.FC = () => {
  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const location = useLocation();
  const query = useMemo(() => new URLSearchParams(location.search), [location]);
  const performer = {
    name: query.get("q") ?? undefined,
  };

  const intl = useIntl();

  function renderPerformerImage() {
    if (encodingImage) {
      return <LoadingIndicator message="Encoding image..." />;
    }
    if (image) {
      return (
        <img
          className="performer"
          src={image}
          alt={intl.formatMessage({ id: "performer" })}
        />
      );
    }
  }

  return (
    <div className="row new-view" id="performer-page">
      <div className="performer-image-container col-md-4 text-center">
        {renderPerformerImage()}
      </div>
      <div className="col-md-8">
        <h2>
          <FormattedMessage
            id="actions.create_entity"
            values={{ entityType: intl.formatMessage({ id: "performer" }) }}
          />
        </h2>
        <PerformerEditPanel
          performer={performer}
          isVisible
          setImage={setImage}
          setEncodingImage={setEncodingImage}
        />
      </div>
    </div>
  );
};

export default PerformerCreate;
