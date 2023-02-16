import React, { useEffect, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { useLocation } from "react-router-dom";
import { SceneEditPanel } from "./SceneEditPanel";
import { useFindScene } from "src/core/StashService";
import ImageUtils from "src/utils/image";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";

const SceneCreate: React.FC = () => {
  const intl = useIntl();

  const location = useLocation();
  const query = useMemo(() => new URLSearchParams(location.search), [location]);

  // create scene from provided scene id if applicable
  const { data, loading } = useFindScene(query.get("from_scene_id") ?? "");
  const [loadingCoverImage, setLoadingCoverImage] = useState(false);
  const [coverImage, setCoverImage] = useState<string>();

  const scene = useMemo(() => {
    if (data?.findScene) {
      return {
        ...data.findScene,
        paths: undefined,
        id: undefined,
      };
    }

    return {
      title: query.get("q") ?? undefined,
    };
  }, [data?.findScene, query]);

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
          fileID={query.get("file_id") ?? undefined}
          initialCoverImage={coverImage}
          isVisible
          isNew
        />
      </div>
    </div>
  );
};

export default SceneCreate;
