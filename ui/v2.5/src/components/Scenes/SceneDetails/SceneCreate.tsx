import React, { useEffect, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { SceneEditPanel } from "./SceneEditPanel";
import queryString from "query-string";
import { useFindScene } from "src/core/StashService";
import { ImageUtils } from "src/utils";
import { LoadingIndicator } from "src/components/Shared";

const SceneCreate: React.FC = () => {
  const intl = useIntl();

  // create scene from provided scene id if applicable
  const queryParams = queryString.parse(location.search);

  const fromSceneID = (queryParams?.from_scene_id ?? "") as string;
  const { data, loading } = useFindScene(fromSceneID ?? "");
  const [loadingCoverImage, setLoadingCoverImage] = useState(false);
  const [coverImage, setCoverImage] = useState<string | undefined>(undefined);

  const scene = useMemo(() => {
    if (data?.findScene) {
      return {
        ...data.findScene,
        paths: undefined,
        id: undefined,
      };
    }

    return {};
  }, [data?.findScene]);

  useEffect(() => {
    async function fetchCoverImage() {
      const srcScene = data?.findScene;
      if (srcScene?.paths.screenshot) {
        setLoadingCoverImage(true);
        const imageData = await ImageUtils.imageToDataURL(
          srcScene.paths.screenshot
        );
        setCoverImage(imageData);
        setLoadingCoverImage(false);
      } else {
        setCoverImage(undefined);
      }
    }

    fetchCoverImage();
  }, [data?.findScene]);

  if (loading || loadingCoverImage) {
    return <LoadingIndicator />;
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
        <SceneEditPanel
          scene={scene}
          initialCoverImage={coverImage}
          isVisible
          isNew
        />
      </div>
    </div>
  );
};

export default SceneCreate;
