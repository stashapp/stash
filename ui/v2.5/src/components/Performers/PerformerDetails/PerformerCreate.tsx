import React, { useState } from "react";
import { LoadingIndicator } from "src/components/Shared";
import { PerformerEditPanel } from "./PerformerEditPanel";

const PerformerCreate: React.FC = () => {
  const [imagePreview, setImagePreview] = useState<string | null>();
  const [imageEncoding, setImageEncoding] = useState<boolean>(false);

  const activeImage = imagePreview ?? "";

  const onImageChange = (image?: string | null) => setImagePreview(image);
  const onImageEncoding = (isEncoding = false) => setImageEncoding(isEncoding);

  function renderPerformerImage() {
    if (imageEncoding) {
      return <LoadingIndicator message="Encoding image..." />;
    }
    if (activeImage) {
      return <img className="performer" src={activeImage} alt="Performer" />;
    }
  }

  return (
    <div className="row new-view" id="performer-page">
      <div className="performer-image-container col-md-4 text-center">
        {renderPerformerImage()}
      </div>
      <div className="col-md-8">
        <h2>Create Performer</h2>
        <PerformerEditPanel
          performer={{}}
          isVisible
          isNew
          onImageChange={onImageChange}
          onImageEncoding={onImageEncoding}
        />
      </div>
    </div>
  );
};

export default PerformerCreate;
