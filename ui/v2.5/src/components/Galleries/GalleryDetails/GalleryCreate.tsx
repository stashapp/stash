import React, { useMemo } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory, useLocation } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { useGalleryCreate } from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { GalleryEditPanel } from "./GalleryEditPanel";

const GalleryCreate: React.FC = () => {
  const history = useHistory();
  const intl = useIntl();
  const Toast = useToast();

  const location = useLocation();
  const query = useMemo(() => new URLSearchParams(location.search), [location]);
  const gallery = {
    title: query.get("q") ?? undefined,
  };

  const [createGallery] = useGalleryCreate();

  async function onSave(input: GQL.GalleryCreateInput) {
    const result = await createGallery({
      variables: { input },
    });
    if (result.data?.galleryCreate) {
      history.push(`/galleries/${result.data.galleryCreate.id}`);
      Toast.success(
        intl.formatMessage(
          { id: "toast.created_entity" },
          { entity: intl.formatMessage({ id: "gallery" }).toLocaleLowerCase() }
        )
      );
    }
  }

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
          gallery={gallery}
          isVisible
          onSubmit={onSave}
          onDelete={() => {}}
        />
      </div>
    </div>
  );
};

export default GalleryCreate;
