import { Tab, Nav, Dropdown } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { useParams, useHistory, Link } from "react-router-dom";
import { FormattedMessage, useIntl } from "react-intl";
import {
  mutateMetadataScan,
  useFindGallery,
  useGalleryUpdate,
} from "src/core/StashService";
import { ErrorMessage, LoadingIndicator, Icon } from "src/components/Shared";
import { TextUtils } from "src/utils";
import * as Mousetrap from "mousetrap";
import { useToast } from "src/hooks";
import { OrganizedButton } from "src/components/Scenes/SceneDetails/OrganizedButton";
import { GalleryEditPanel } from "./GalleryEditPanel";
import { GalleryDetailPanel } from "./GalleryDetailPanel";
import { DeleteGalleriesDialog } from "../DeleteGalleriesDialog";
import { GalleryImagesPanel } from "./GalleryImagesPanel";
import { GalleryAddPanel } from "./GalleryAddPanel";
import { GalleryFileInfoPanel } from "./GalleryFileInfoPanel";
import { GalleryScenesPanel } from "./GalleryScenesPanel";

interface IGalleryParams {
  id?: string;
  tab?: string;
}

export const Gallery: React.FC = () => {
  const { tab = "images", id = "new" } = useParams<IGalleryParams>();
  const history = useHistory();
  const Toast = useToast();
  const intl = useIntl();
  const isNew = id === "new";

  const { data, error, loading } = useFindGallery(id);
  const gallery = data?.findGallery;

  const [activeTabKey, setActiveTabKey] = useState("gallery-details-panel");
  const activeRightTabKey = tab === "images" || tab === "add" ? tab : "images";
  const setActiveRightTabKey = (newTab: string | null) => {
    if (tab !== newTab) {
      const tabParam = newTab === "images" ? "" : `/${newTab}`;
      history.replace(`/galleries/${id}${tabParam}`);
    }
  };

  const [updateGallery] = useGalleryUpdate();

  const [organizedLoading, setOrganizedLoading] = useState(false);

  const onOrganizedClick = async () => {
    try {
      setOrganizedLoading(true);
      await updateGallery({
        variables: {
          input: {
            id: gallery?.id ?? "",
            organized: !gallery?.organized,
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      setOrganizedLoading(false);
    }
  };

  async function onRescan() {
    if (!gallery || !gallery.path) {
      return;
    }

    await mutateMetadataScan({
      paths: [gallery.path],
    });

    Toast.success({
      content: intl.formatMessage({ id: "toast.rescanning_gallery" }),
    });
  }

  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  function onDeleteDialogClosed(deleted: boolean) {
    setIsDeleteAlertOpen(false);
    if (deleted) {
      history.push("/galleries");
    }
  }

  function maybeRenderDeleteDialog() {
    if (isDeleteAlertOpen && gallery) {
      return (
        <DeleteGalleriesDialog
          selected={[gallery]}
          onClose={onDeleteDialogClosed}
        />
      );
    }
  }

  function renderOperations() {
    return (
      <Dropdown>
        <Dropdown.Toggle
          variant="secondary"
          id="operation-menu"
          className="minimal"
          title="Operations"
        >
          <Icon icon="ellipsis-v" />
        </Dropdown.Toggle>
        <Dropdown.Menu className="bg-secondary text-white">
          {gallery?.path ? (
            <Dropdown.Item
              key="rescan"
              className="bg-secondary text-white"
              onClick={() => onRescan()}
            >
              <FormattedMessage id="actions.rescan" />
            </Dropdown.Item>
          ) : undefined}
          <Dropdown.Item
            key="delete-gallery"
            className="bg-secondary text-white"
            onClick={() => setIsDeleteAlertOpen(true)}
          >
            <FormattedMessage id="actions.delete_gallery" />
          </Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    );
  }

  function renderTabs() {
    if (!gallery) {
      return;
    }

    return (
      <Tab.Container
        activeKey={activeTabKey}
        onSelect={(k) => k && setActiveTabKey(k)}
      >
        <div>
          <Nav variant="tabs" className="mr-auto">
            <Nav.Item>
              <Nav.Link eventKey="gallery-details-panel">
                <FormattedMessage id="details" />
              </Nav.Link>
            </Nav.Item>
            {gallery.scenes.length > 0 && (
              <Nav.Item>
                <Nav.Link eventKey="gallery-scenes-panel">
                  <FormattedMessage id="scenes" />
                </Nav.Link>
              </Nav.Item>
            )}
            {gallery.path ? (
              <Nav.Item>
                <Nav.Link eventKey="gallery-file-info-panel">
                  File Info
                </Nav.Link>
              </Nav.Item>
            ) : undefined}
            <Nav.Item>
              <Nav.Link eventKey="gallery-edit-panel">
                <FormattedMessage id="actions.edit" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item className="ml-auto">
              <OrganizedButton
                loading={organizedLoading}
                organized={gallery.organized}
                onClick={onOrganizedClick}
              />
            </Nav.Item>
            <Nav.Item>{renderOperations()}</Nav.Item>
          </Nav>
        </div>

        <Tab.Content>
          <Tab.Pane eventKey="gallery-details-panel">
            <GalleryDetailPanel gallery={gallery} />
          </Tab.Pane>
          <Tab.Pane
            className="file-info-panel"
            eventKey="gallery-file-info-panel"
          >
            <GalleryFileInfoPanel gallery={gallery} />
          </Tab.Pane>
          <Tab.Pane eventKey="gallery-edit-panel">
            <GalleryEditPanel
              isVisible={activeTabKey === "gallery-edit-panel"}
              isNew={false}
              gallery={gallery}
              onDelete={() => setIsDeleteAlertOpen(true)}
            />
          </Tab.Pane>
          {gallery.scenes.length > 0 && (
            <Tab.Pane eventKey="gallery-scenes-panel">
              <GalleryScenesPanel scenes={gallery.scenes} />
            </Tab.Pane>
          )}
        </Tab.Content>
      </Tab.Container>
    );
  }

  function renderRightTabs() {
    if (!gallery) {
      return;
    }

    return (
      <Tab.Container
        activeKey={activeRightTabKey}
        unmountOnExit
        onSelect={(k) => k && setActiveRightTabKey(k)}
      >
        <div>
          <Nav variant="tabs" className="mr-auto">
            <Nav.Item>
              <Nav.Link eventKey="images">
                <FormattedMessage id="images" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="add">
                <FormattedMessage id="actions.add" />
              </Nav.Link>
            </Nav.Item>
          </Nav>
        </div>

        <Tab.Content>
          <Tab.Pane eventKey="images">
            <GalleryImagesPanel gallery={gallery} />
          </Tab.Pane>
          <Tab.Pane eventKey="add">
            <GalleryAddPanel gallery={gallery} />
          </Tab.Pane>
        </Tab.Content>
      </Tab.Container>
    );
  }

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("a", () => setActiveTabKey("gallery-details-panel"));
    Mousetrap.bind("e", () => setActiveTabKey("gallery-edit-panel"));
    Mousetrap.bind("f", () => setActiveTabKey("gallery-file-info-panel"));

    return () => {
      Mousetrap.unbind("a");
      Mousetrap.unbind("e");
      Mousetrap.unbind("f");
    };
  });

  if (loading) {
    return <LoadingIndicator />;
  }

  if (error) return <ErrorMessage error={error.message} />;

  if (isNew)
    return (
      <div className="row new-view">
        <div className="col-6">
          <h2>
            <FormattedMessage id="actions.create_gallery" />
          </h2>
          <GalleryEditPanel
            isNew
            gallery={undefined}
            isVisible
            onDelete={() => setIsDeleteAlertOpen(true)}
          />
        </div>
      </div>
    );

  if (!gallery)
    return <ErrorMessage error={`No gallery with id ${id} found.`} />;

  return (
    <div className="row">
      {maybeRenderDeleteDialog()}
      <div className="gallery-tabs">
        <div className="d-none d-xl-block">
          {gallery.studio && (
            <h1 className="text-center">
              <Link to={`/studios/${gallery.studio.id}`}>
                <img
                  src={gallery.studio.image_path ?? ""}
                  alt={`${gallery.studio.name} logo`}
                  className="studio-logo"
                />
              </Link>
            </h1>
          )}
          <h3 className="gallery-header">
            {gallery.title ?? TextUtils.fileNameFromPath(gallery.path ?? "")}
          </h3>
        </div>
        {renderTabs()}
      </div>
      <div className="gallery-container">{renderRightTabs()}</div>
    </div>
  );
};
