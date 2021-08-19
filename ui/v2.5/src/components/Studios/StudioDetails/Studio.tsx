import { Tabs, Tab } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import { FormattedMessage, useIntl } from "react-intl";
import cx from "classnames";
import Mousetrap from "mousetrap";

import * as GQL from "src/core/generated-graphql";
import {
  useFindStudio,
  useStudioUpdate,
  useStudioCreate,
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

interface IStudioParams {
  id?: string;
  tab?: string;
}

export const Studio: React.FC = () => {
  const history = useHistory();
  const Toast = useToast();
  const intl = useIntl();
  const { tab = "details", id = "new" } = useParams<IStudioParams>();
  const isNew = id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Studio state
  const [image, setImage] = useState<string | null>();

  const { data, error } = useFindStudio(id);
  const studio = data?.findStudio;

  const [isLoading, setIsLoading] = useState(false);
  const [updateStudio] = useStudioUpdate();
  const [createStudio] = useStudioCreate();
  const [deleteStudio] = useStudioDestroy({ id });

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => setIsEditing(true));
    Mousetrap.bind("d d", () => onDelete());

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

  useEffect(() => {
    if (data && data.findStudio) {
      setImage(undefined);
    }
  }, [data]);

  function onImageLoad(imageData: string) {
    setImage(imageData);
  }

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, isEditing);

  async function onSave(
    input: Partial<GQL.StudioCreateInput | GQL.StudioUpdateInput>
  ) {
    try {
      setIsLoading(true);

      if (!isNew) {
        const result = await updateStudio({
          variables: {
            input: input as GQL.StudioUpdateInput,
          },
        });
        if (result.data?.studioUpdate) {
          setIsEditing(false);
        }
      } else {
        const result = await createStudio({
          variables: {
            input: input as GQL.StudioCreateInput,
          },
        });
        if (result.data?.studioCreate?.id) {
          history.push(`/studios/${result.data.studioCreate.id}`);
          setIsEditing(false);
        }
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function onAutoTag() {
    if (!studio?.id) return;
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
                studio?.name ??
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
    let studioImage = studio?.image_path;
    if (isEditing) {
      if (image === null) {
        studioImage = `${studioImage}&default=true`;
      } else if (image) {
        studioImage = image;
      }
    }

    if (studioImage) {
      return (
        <img className="logo" alt={studio?.name ?? ""} src={studioImage} />
      );
    }
  }

  const activeTabKey =
    tab === "childstudios" ||
    tab === "images" ||
    tab === "galleries" ||
    tab === "performers"
      ? tab
      : "scenes";
  const setActiveTabKey = (newTab: string | null) => {
    if (tab !== newTab) {
      const tabParam = newTab === "scenes" ? "" : `/${newTab}`;
      history.replace(`/studios/${id}${tabParam}`);
    }
  };

  if (isLoading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!studio?.id && !isNew)
    return <ErrorMessage error={`No studio found with id ${id}.`} />;

  return (
    <div className="row">
      <div
        className={cx("studio-details", {
          "col-md-4": !isNew,
          "col-md-8": isNew,
        })}
      >
        {isNew && (
          <h2>
            {intl.formatMessage(
              { id: "actions.add_entity" },
              { entityType: intl.formatMessage({ id: "studio" }) }
            )}
          </h2>
        )}
        <div className="text-center">
          {imageEncoding ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            renderImage()
          )}
        </div>
        {!isEditing && !isNew && studio ? (
          <>
            <StudioDetailsPanel studio={studio} />
            <DetailsEditNavbar
              objectName={studio.name ?? intl.formatMessage({ id: "studio" })}
              isNew={isNew}
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
            studio={studio ?? ({} as Partial<GQL.Studio>)}
            onSubmit={onSave}
            onCancel={onToggleEdit}
            onDelete={onDelete}
            onImageChange={setImage}
          />
        )}
      </div>
      {studio?.id && (
        <div className="col col-md-8">
          <Tabs
            id="studio-tabs"
            mountOnEnter
            unmountOnExit
            activeKey={activeTabKey}
            onSelect={setActiveTabKey}
          >
            <Tab eventKey="scenes" title={intl.formatMessage({ id: "scenes" })}>
              <StudioScenesPanel studio={studio} />
            </Tab>
            <Tab
              eventKey="galleries"
              title={intl.formatMessage({ id: "galleries" })}
            >
              <StudioGalleriesPanel studio={studio} />
            </Tab>
            <Tab eventKey="images" title={intl.formatMessage({ id: "images" })}>
              <StudioImagesPanel studio={studio} />
            </Tab>
            <Tab
              eventKey="performers"
              title={intl.formatMessage({ id: "performers" })}
            >
              <StudioPerformersPanel studio={studio} />
            </Tab>
            <Tab
              eventKey="childstudios"
              title={intl.formatMessage({ id: "subsidiary_studios" })}
            >
              <StudioChildrenPanel studio={studio} />
            </Tab>
          </Tabs>
        </div>
      )}
      {renderDeleteAlert()}
    </div>
  );
};
