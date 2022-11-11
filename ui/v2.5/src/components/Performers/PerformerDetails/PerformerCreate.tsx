import React, { useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { LoadingIndicator } from "src/components/Shared";
import { PerformerEditPanel } from "./PerformerEditPanel";
import { useLocation } from "react-router-dom";

const PerformerCreate: React.FC = () => {
  const [imagePreview, setImagePreview] = useState<string | null>();
  const [imageEncoding, setImageEncoding] = useState<boolean>(false);

  const location = useLocation();
  const query = useMemo(() => new URLSearchParams(location.search), [location]);
  const performer = {
    name: query.get("q") ?? undefined,
  };

  const activeImage = imagePreview ?? "";
  const intl = useIntl();

  const onImageChange = (image?: string | null) => setImagePreview(image);
  const onImageEncoding = (isEncoding = false) => setImageEncoding(isEncoding);

  function renderPerformerImage() {
    if (imageEncoding) {
      return <LoadingIndicator message="Encoding image..." />;
    }
    if (activeImage) {
      return (
        <img
          className="performer"
          src={activeImage}
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
          onImageChange={onImageChange}
          onImageEncoding={onImageEncoding}
        />
      </div>
    </div>
  );
};

export default PerformerCreate;
