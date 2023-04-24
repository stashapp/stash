import { Button, Tabs, Tab } from "react-bootstrap";
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
import { Counter } from "src/components/Shared/Counter";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { ModalComponent } from "src/components/Shared/Modal";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { useToast } from "src/hooks/Toast";
import { ConfigurationContext } from "src/hooks/Config";
import { Icon } from "src/components/Shared/Icon";
import { StudioScenesPanel } from "./StudioScenesPanel";
import { StudioGalleriesPanel } from "./StudioGalleriesPanel";
import { StudioImagesPanel } from "./StudioImagesPanel";
import { StudioChildrenPanel } from "./StudioChildrenPanel";
import { StudioPerformersPanel } from "./StudioPerformersPanel";
import { StudioEditPanel } from "./StudioEditPanel";
import { StudioDetailsPanel } from "./StudioDetailsPanel";
import { StudioMoviesPanel } from "./StudioMoviesPanel";
import {
  faTrashAlt,
  faChevronRight,
  faChevronLeft,
} from "@fortawesome/free-solid-svg-icons";
import { IUIConfig } from "src/core/config";

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

  const [collapsed, setCollapsed] = useState(false);

  // Configuration settings
  const { configuration } = React.useContext(ConfigurationContext);
  const abbreviateCounter =
    (configuration?.ui as IUIConfig)?.abbreviateCounters ?? false;

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing studio state
  const [image, setImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const [updateStudio] = useStudioUpdate();
  const [deleteStudio] = useStudioDestroy({ id: studio.id });

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => setIsEditing(true));
    Mousetrap.bind("d d", () => {
      onDelete();
    });
    Mousetrap.bind(",", () => setCollapsed(!collapsed));

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
      Mousetrap.unbind(",");
    };
  });

  async function onSave(input: GQL.StudioCreateInput) {
    try {
      const result = await updateStudio({
        variables: {
          input: {
            id: studio.id,
            ...input,
          },
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
      <ModalComponent
        show={isDeleteAlertOpen}
        icon={faTrashAlt}
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
      </ModalComponent>
    );
  }

  function onToggleEdit() {
    setIsEditing(!isEditing);
  }

  function renderImage() {
    let studioImage = studio.image_path;
    if (isEditing) {
      if (image === null && studioImage) {
        const studioImageURL = new URL(studioImage);
        studioImageURL.searchParams.set("default", "true");
        studioImage = studioImageURL.toString();
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

  function getCollapseButtonIcon() {
    return collapsed ? faChevronRight : faChevronLeft;
  }

  return (
    <div className="row">
      <div
        className={`studio-details details-tab ${collapsed ? "collapsed" : ""}`}
      >
        <div className="text-center">
          {encodingImage ? (
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
            setImage={setImage}
            setEncodingImage={setEncodingImage}
          />
        )}
      </div>
      <div className="details-divider d-none d-xl-block">
        <Button onClick={() => setCollapsed(!collapsed)}>
          <Icon className="fa-fw" icon={getCollapseButtonIcon()} />
        </Button>
      </div>
      <div className={`col content-container ${collapsed ? "expanded" : ""}`}>
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
                <Counter
                  abbreviateCounter={abbreviateCounter}
                  count={studio.scene_count ?? 0}
                />
              </React.Fragment>
            }
          >
            <StudioScenesPanel
              active={activeTabKey == "scenes"}
              studio={studio}
            />
          </Tab>
          <Tab
            eventKey="galleries"
            title={
              <React.Fragment>
                {intl.formatMessage({ id: "galleries" })}
                <Counter
                  abbreviateCounter={abbreviateCounter}
                  count={studio.gallery_count ?? 0}
                />
              </React.Fragment>
            }
          >
            <StudioGalleriesPanel
              active={activeTabKey == "galleries"}
              studio={studio}
            />
          </Tab>
          <Tab
            eventKey="images"
            title={
              <React.Fragment>
                {intl.formatMessage({ id: "images" })}
                <Counter
                  abbreviateCounter={abbreviateCounter}
                  count={studio.image_count ?? 0}
                />
              </React.Fragment>
            }
          >
            <StudioImagesPanel
              active={activeTabKey == "images"}
              studio={studio}
            />
          </Tab>
          <Tab
            eventKey="performers"
            title={
              <React.Fragment>
                {intl.formatMessage({ id: "performers" })}
                <Counter
                  abbreviateCounter={abbreviateCounter}
                  count={studio.performer_count ?? 0}
                />
              </React.Fragment>
            }
          >
            <StudioPerformersPanel
              active={activeTabKey == "performers"}
              studio={studio}
            />
          </Tab>
          <Tab
            eventKey="movies"
            title={
              <React.Fragment>
                {intl.formatMessage({ id: "movies" })}
                <Counter
                  abbreviateCounter={abbreviateCounter}
                  count={studio.movie_count ?? 0}
                />
              </React.Fragment>
            }
          >
            <StudioMoviesPanel
              active={activeTabKey == "movies"}
              studio={studio}
            />
          </Tab>
          <Tab
            eventKey="childstudios"
            title={
              <React.Fragment>
                {intl.formatMessage({ id: "subsidiary_studios" })}
                <Counter
                  abbreviateCounter={false}
                  count={studio.child_studios?.length ?? 0}
                />
              </React.Fragment>
            }
          >
            <StudioChildrenPanel
              active={activeTabKey == "childstudios"}
              studio={studio}
            />
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
