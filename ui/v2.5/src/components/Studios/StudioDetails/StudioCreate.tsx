import React, { useMemo, useState } from "react";
import { useHistory, useLocation } from "react-router-dom";
import { useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { useStudioCreate } from "src/core/StashService";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { StudioEditPanel } from "./StudioEditPanel";

const StudioCreate: React.FC = () => {
  const history = useHistory();
  const location = useLocation();
  const Toast = useToast();

  const query = useMemo(() => new URLSearchParams(location.search), [location]);
  const studio = {
    name: query.get("q") ?? undefined,
  };

  const intl = useIntl();

  // Editing studio state
  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const [createStudio] = useStudioCreate();

  async function onSave(input: GQL.StudioCreateInput) {
    const result = await createStudio({
      variables: { input },
    });
    if (result.data?.studioCreate?.id) {
      history.push(`/studios/${result.data.studioCreate.id}`);
      Toast.success(
        intl.formatMessage(
          { id: "toast.created_entity" },
          { entity: intl.formatMessage({ id: "studio" }).toLocaleLowerCase() }
        )
      );
    }
  }

  function renderImage() {
    if (image) {
      return <img className="logo" alt="" src={image} />;
    }
  }

  return (
    <div className="row">
      <div className="studio-details col-md-8">
        <h2>
          {intl.formatMessage(
            { id: "actions.add_entity" },
            { entityType: intl.formatMessage({ id: "studio" }) }
          )}
        </h2>
        <div className="text-center">
          {encodingImage ? (
            <LoadingIndicator
              message={intl.formatMessage({ id: "actions.encoding_image" })}
            />
          ) : (
            renderImage()
          )}
        </div>
        <StudioEditPanel
          studio={studio}
          onSubmit={onSave}
          onCancel={() => history.push("/studios")}
          onDelete={() => {}}
          setImage={setImage}
          setEncodingImage={setEncodingImage}
        />
      </div>
    </div>
  );
};

export default StudioCreate;
