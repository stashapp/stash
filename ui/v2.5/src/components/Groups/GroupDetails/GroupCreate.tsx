import React, { useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { useGroupCreate } from "src/core/StashService";
import { useHistory, useLocation } from "react-router-dom";
import { useIntl } from "react-intl";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { GroupEditPanel } from "./GroupEditPanel";

const GroupCreate: React.FC = () => {
  const history = useHistory();
  const intl = useIntl();
  const Toast = useToast();

  const location = useLocation();
  const query = useMemo(() => new URLSearchParams(location.search), [location]);
  const group = {
    name: query.get("q") ?? undefined,
  };

  // Editing group state
  const [frontImage, setFrontImage] = useState<string | null>();
  const [backImage, setBackImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const [createGroup] = useGroupCreate();

  async function onSave(input: GQL.GroupCreateInput) {
    const result = await createGroup({
      variables: { input },
    });
    if (result.data?.groupCreate?.id) {
      history.push(`/groups/${result.data.groupCreate.id}`);
      Toast.success(
        intl.formatMessage(
          { id: "toast.created_entity" },
          { entity: intl.formatMessage({ id: "group" }).toLocaleLowerCase() }
        )
      );
    }
  }

  function renderFrontImage() {
    if (frontImage) {
      return (
        <div className="group-image-container">
          <img alt="Front Cover" src={frontImage} />
        </div>
      );
    }
  }

  function renderBackImage() {
    if (backImage) {
      return (
        <div className="group-image-container">
          <img alt="Back Cover" src={backImage} />
        </div>
      );
    }
  }

  // TODO: CSS class
  return (
    <div className="row">
      <div className="group-details mb-3 col">
        <div className="logo w-100">
          {encodingImage ? (
            <LoadingIndicator
              message={intl.formatMessage({ id: "actions.encoding_image" })}
            />
          ) : (
            <div className="group-images">
              {renderFrontImage()}
              {renderBackImage()}
            </div>
          )}
        </div>

        <GroupEditPanel
          group={group}
          onSubmit={onSave}
          onCancel={() => history.push("/groups")}
          onDelete={() => {}}
          setFrontImage={setFrontImage}
          setBackImage={setBackImage}
          setEncodingImage={setEncodingImage}
        />
      </div>
    </div>
  );
};

export default GroupCreate;
