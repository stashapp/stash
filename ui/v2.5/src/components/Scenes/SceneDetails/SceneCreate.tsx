import React, { useEffect, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory, useLocation } from "react-router-dom";
import { SceneEditPanel } from "./SceneEditPanel";
import * as GQL from "src/core/generated-graphql";
import { mutateCreateScene, useFindScene } from "src/core/StashService";
import ImageUtils from "src/utils/image";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { Button } from "react-bootstrap";
import { Icon } from "src/components/Shared/Icon";
import { faUpload } from "@fortawesome/free-solid-svg-icons";

const SceneCreate: React.FC = () => {
  const history = useHistory();
  const intl = useIntl();
  const Toast = useToast();

  const location = useLocation();
  const query = useMemo(() => new URLSearchParams(location.search), [location]);

  // create scene from provided scene id if applicable
  const { data, loading } = useFindScene(query.get("from_scene_id") ?? "new");
  const [loadingCoverImage, setLoadingCoverImage] = useState(false);
  const [coverImage, setCoverImage] = useState<string>();
  const [uploading, setUploading] = useState(false);
  const [uploadedFile, setUploadedFile] = useState<{ path: string } | null>(null);

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

  async function onFileUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const files = e.target.files;
    if (!files || files.length === 0) return;

    setUploading(true);
    const formData = new FormData();
    formData.append("file", files[0]);

    try {
      const resp = await fetch("/scene/upload", {
        method: "POST",
        body: formData,
      });
      const result = await resp.json();

      if (result.success) {
        setUploadedFile(result);
        Toast.success(
          intl.formatMessage(
            { id: "toast.upload_complete" },
            { filename: files[0].name }
          )
        );
        // Navigate to scenes list to see the scan result
        history.push("/scenes");
      } else {
        Toast.error(intl.formatMessage({ id: "toast.upload_failed" }));
      }
    } catch (err) {
      Toast.error(intl.formatMessage({ id: "toast.upload_failed" }));
    } finally {
      setUploading(false);
    }
  }

  async function onSave(input: GQL.SceneCreateInput, andNew?: boolean) {
    const fileID = query.get("file_id") ?? undefined;
    const result = await mutateCreateScene({
      ...input,
      file_ids: fileID ? [fileID] : undefined,
    });
    if (result.data?.sceneCreate?.id) {
      if (!andNew) {
        history.push(`/scenes/${result.data.sceneCreate.id}`);
      }
      Toast.success(
        intl.formatMessage(
          { id: "toast.created_entity" },
          { entity: intl.formatMessage({ id: "scene" }).toLocaleLowerCase() }
        )
      );
    }
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

        {!uploadedFile && (
          <div className="upload-area mb-4 p-4 text-center border rounded">
            <h5>
              <Icon icon={faUpload} className="mr-2" />
              <FormattedMessage id="actions.upload_file" />
            </h5>
            <p className="text-muted">
              <FormattedMessage id="upload_scene.description" />
            </p>
            <Button
              variant="primary"
              disabled={uploading}
              onClick={() => document.getElementById("scene-file-input")?.click()}
            >
              {uploading ? (
                <LoadingIndicator inline small />
              ) : (
                <FormattedMessage id="actions.choose_file" />
              )}
            </Button>
            <input
              id="scene-file-input"
              type="file"
              accept="video/*"
              className="d-none"
              onChange={onFileUpload}
            />
          </div>
        )}

        {uploadedFile && (
          <div className="alert alert-success">
            <FormattedMessage
              id="upload_scene.success"
              values={{ path: uploadedFile.path }}
            />
          </div>
        )}

        <SceneEditPanel
          scene={scene}
          initialCoverImage={coverImage}
          isVisible
          isNew
          onSubmit={onSave}
        />
      </div>
    </div>
  );
};

export default SceneCreate;
