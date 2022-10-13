import React, { useMemo } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { SceneEditPanel } from "./SceneEditPanel";
import queryString from "query-string";
import { useFindScene } from "src/core/StashService";

const SceneCreate: React.FC = () => {
  const intl = useIntl();

  // create scene from provided scene id if applicable
  const queryParams = queryString.parse(location.search);

  const fromSceneID = (queryParams?.from_scene_id ?? "") as string;
  const { data, loading } = useFindScene(fromSceneID ?? "");

  const scene = useMemo(() => {
    if (data?.findScene) {
      return {
        ...data.findScene,
        id: undefined,
      };
    }

    return {};
  }, [data?.findScene]);

  if (loading) {
    return <></>;
  }

  return (
    <div className="row new-view justify-content-center" id="create-scene-page">
      <div className="col-md-8">
        <h2>
          <FormattedMessage
            id="actions.create_entity"
            values={{ entityType: intl.formatMessage({ id: "scene" }) }}
          />
        </h2>
        <SceneEditPanel scene={scene} isVisible isNew />
      </div>
    </div>
  );
};

export default SceneCreate;
