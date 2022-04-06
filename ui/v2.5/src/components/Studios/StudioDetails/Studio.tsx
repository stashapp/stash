import { Tabs, Tab, Badge } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import { FormattedMessage, useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import Mousetrap from "mousetrap";

import * as GQL from "src/core/generated-graphql";
import {
  useFindStudio,
  useStudioUpdate,
  useStudioDestroy,
  mutateMetadataAutoTag,
} from "src/core/StashService";
import { ImageUtils } from "src/utils";
import {
  DetailsEditNavbar,
  Modal,
  LoadingIndicator,
  ErrorMessage,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { StudioScenesPanel } from "./StudioScenesPanel";
import { StudioGalleriesPanel } from "./StudioGalleriesPanel";
import { StudioImagesPanel } from "./StudioImagesPanel";
import { StudioChildrenPanel } from "./StudioChildrenPanel";
import { StudioPerformersPanel } from "./StudioPerformersPanel";
import { StudioEditPanel } from "./StudioEditPanel";
import { StudioDetailsPanel } from "./StudioDetailsPanel";
import { StudioMoviesPanel } from "./StudioMoviesPanel";

interface IProps {
  studio: GQL.StudioDataFragment;
}

interface IStudioParams {
  tab?: string;
}

const StudioPage: React.FC<IProps> = ({ studio }) => {
  const history = useHistory();
  const Toast = useToast();
  const intl = useIntl();
  const { tab = "details" } = useParams<IStudioParams>();

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Studio state
  const [image, setImage] = useState<string | null>();

  const [updateStudio] = useStudioUpdate();
  const [deleteStudio] = useStudioDestroy({ id: studio.id });

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => setIsEditing(true));
    Mousetrap.bind("d d", () => onDelete());

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

  function onImageLoad(imageData: string) {
    setImage(imageData);
  }

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, isEditing);

  async function onSave(input: Partial<GQL.StudioUpdateInput>) {
    try {
      const result = await updateStudio({
        variables: {
          input: input as GQL.StudioUpdateInput,
        },
      });
      if (result.data?.studioUpdate) {
        setIsEditing(false);
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onAutoTag() {
    if (!studio.id) return;
    try {
      await mutateMetadataAutoTag({ studios: [studio.id] });
      Toast.success({
        content: intl.formatMessage({ id: "toast.started_auto_tagging" }),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onDelete() {
    try {
      await deleteStudio();
    } catch (e) {
      Toast.error(e);
    }

    // redirect to studios page
    history.push(`/studios`);
  }

  function renderDeleteAlert() {
    return (
      <Modal
        show={isDeleteAlertOpen}
        icon="trash-alt"
        accept={{
          text: intl.formatMessage({ id: "actions.delete" }),
          variant: "danger",
          onClick: onDelete,
        }}
        cancel={{ onClick: () => setIsDeleteAlertOpen(false) }}
      >
        <p>
          <FormattedMessage
            id="dialogs.delete_confirm"
            values={{
              entityName:
                studio.name ??
                intl.formatMessage({ id: "studio" }).toLocaleLowerCase(),
            }}
          />
        </p>
      </Modal>
    );
  }

  function onToggleEdit() {
    setIsEditing(!isEditing);
  }

  function renderImage() {
    let studioImage = studio.image_path;
    if (isEditing) {
      if (image === null) {
        studioImage = `${studioImage}&default=true`;
      } else if (image) {
        studioImage = image;
      }
    }

    if (studioImage) {
      return <img className="logo" alt={studio.name} src={studioImage} />;
    }
  }

  const activeTabKey =
    tab === "childstudios" ||
    tab === "images" ||
    tab === "galleries" ||
    tab === "performers" ||
    tab === "movies"
      ? tab
      : "scenes";
  const setActiveTabKey = (newTab: string | null) => {
    if (tab !== newTab) {
      const tabParam = newTab === "scenes" ? "" : `/${newTab}`;
      history.replace(`/studios/${studio.id}${tabParam}`);
    }
  };

  return (
    <div className="row">
      <div className="studio-details col-md-4">
        <div className="text-center">
          {imageEncoding ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            renderImage()
          )}
        </div>
        {!isEditing ? (
          <>
            <Helmet>
              <title>
                {studio.name ?? intl.formatMessage({ id: "studio" })}
              </title>
            </Helmet>
            <StudioDetailsPanel studio={studio} />
            <DetailsEditNavbar
              objectName={studio.name ?? intl.formatMessage({ id: "studio" })}
              isNew={false}
              isEditing={isEditing}
              onToggleEdit={onToggleEdit}
              onSave={() => {}}
              onImageChange={() => {}}
              onClearImage={() => {}}
              onAutoTag={onAutoTag}
              onDelete={onDelete}
            />
          </>
        ) : (
          <StudioEditPanel
            studio={studio}
            onSubmit={onSave}
            onCancel={onToggleEdit}
            onDelete={onDelete}
            onImageChange={setImage}
          />
        )}
      </div>
      <div className="col col-md-8">
        <Tabs
          id="studio-tabs"
          mountOnEnter
          unmountOnExit
          activeKey={activeTabKey}
          onSelect={setActiveTabKey}
        >
          <Tab
            eventKey="scenes"
            title={
              <React.Fragment>
                {intl.formatMessage({ id: "scenes" })}
                <Badge className="left-spacing" pill variant="secondary">
                  {intl.formatNumber(studio.scene_count ?? 0)}
                </Badge>
              </React.Fragment>
            }
          >
            <StudioScenesPanel studio={studio} />
          </Tab>
          <Tab
            eventKey="galleries"
            title={
              <React.Fragment>
                {intl.formatMessage({ id: "galleries" })}
                <Badge className="left-spacing" pill variant="secondary">
                  {intl.formatNumber(studio.gallery_count ?? 0)}
                </Badge>
              </React.Fragment>
            }
          >
            <StudioGalleriesPanel studio={studio} />
          </Tab>
          <Tab
            eventKey="images"
            title={
              <React.Fragment>
                {intl.formatMessage({ id: "images" })}
                <Badge className="left-spacing" pill variant="secondary">
                  {intl.formatNumber(studio.image_count ?? 0)}
                </Badge>
              </React.Fragment>
            }
          >
            <StudioImagesPanel studio={studio} />
          </Tab>
          <Tab
            eventKey="performers"
            title={intl.formatMessage({ id: "performers" })}
          >
            <StudioPerformersPanel studio={studio} />
          </Tab>
          <Tab
            eventKey="movies"
            title={
              <React.Fragment>
                {intl.formatMessage({ id: "movies" })}
                <Badge className="left-spacing" pill variant="secondary">
                  {intl.formatNumber(studio.movie_count ?? 0)}
                </Badge>
              </React.Fragment>
            }
          >
            <StudioMoviesPanel studio={studio} />
          </Tab>
          <Tab
            eventKey="childstudios"
            title={
              <React.Fragment>
                {intl.formatMessage({ id: "subsidiary_studios" })}
                <Badge className="left-spacing" pill variant="secondary">
                  {intl.formatNumber(studio.child_studios?.length)}
                </Badge>
              </React.Fragment>
            }
          >
            <StudioChildrenPanel studio={studio} />
          </Tab>
        </Tabs>
      </div>
      {renderDeleteAlert()}
    </div>
  );
};

const StudioLoader: React.FC = () => {
  const { id } = useParams<{ id?: string }>();
  const { data, loading, error } = useFindStudio(id ?? "");

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findStudio)
    return <ErrorMessage error={`No studio found with id ${id}.`} />;

  return <StudioPage studio={data.findStudio} />;
};

export default StudioLoader;
