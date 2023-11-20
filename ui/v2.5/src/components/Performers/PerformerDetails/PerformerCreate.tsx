import React, { useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { PerformerEditPanel } from "./PerformerEditPanel";
import { useHistory, useLocation } from "react-router-dom";
import { useToast } from "src/hooks/Toast";
import * as GQL from "src/core/generated-graphql";
import { usePerformerCreate } from "src/core/StashService";

const PerformerCreate: React.FC = () => {
  const Toast = useToast();
  const history = useHistory();
  const intl = useIntl();

  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const location = useLocation();
  const query = useMemo(() => new URLSearchParams(location.search), [location]);
  const performer = {
    name: query.get("q") ?? undefined,
  };

  const [createPerformer] = usePerformerCreate();

  async function onSave(input: GQL.PerformerCreateInput) {
    const result = await createPerformer({
      variables: { input },
    });
    if (result.data?.performerCreate) {
      history.push(`/performers/${result.data.performerCreate.id}`);
      Toast.success(
        intl.formatMessage(
          { id: "toast.created_entity" },
          {
            entity: intl.formatMessage({ id: "performer" }).toLocaleLowerCase(),
          }
        )
      );
    }
  }

  function renderPerformerImage() {
    if (encodingImage) {
      return (
        <LoadingIndicator
          message={intl.formatMessage({ id: "actions.encoding_image" })}
        />
      );
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
          onSubmit={onSave}
          setImage={setImage}
          setEncodingImage={setEncodingImage}
        />
      </div>
    </div>
  );
};

export default PerformerCreate;
