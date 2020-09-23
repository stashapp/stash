import { Tab, Nav, Dropdown } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { useParams, useHistory, Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { useFindGallery } from "src/core/StashService";
import { LoadingIndicator, Icon } from "src/components/Shared";
import { TextUtils } from "src/utils";
import * as Mousetrap from "mousetrap";
import { GalleryEditPanel } from "./GalleryEditPanel";
import { GalleryDetailPanel } from "./GalleryDetailPanel";
import { DeleteGalleriesDialog } from "../DeleteGalleriesDialog";
import { GalleryImagesPanel } from "./GalleryImagesPanel";
import { GalleryAddPanel } from "./GalleryAddPanel";

interface IGalleryParams {
  id?: string;
  tab?: string;
}

export const Gallery: React.FC = () => {
  const { tab = "images", id = "new" } = useParams<IGalleryParams>();
  const history = useHistory();
  const isNew = id === "new";

  const [gallery, setGallery] = useState<Partial<GQL.GalleryDataFragment>>({});
  const { data, error, loading } = useFindGallery(id);

  const [activeTabKey, setActiveTabKey] = useState("gallery-details-panel");
  const activeRightTabKey = tab === "images" || tab === "add" ? tab : "images";
  const setActiveRightTabKey = (newTab: string | null) => {
    if (tab !== newTab) {
      const tabParam = newTab === "images" ? "" : `/${newTab}`;
      history.replace(`/galleries/${id}${tabParam}`);
    }
  };

  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  useEffect(() => {
    if (data?.findGallery) setGallery(data.findGallery);
  }, [data]);

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
          <Dropdown.Item
            key="delete-gallery"
            className="bg-secondary text-white"
            onClick={() => setIsDeleteAlertOpen(true)}
          >
            Delete Gallery
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
              <Nav.Link eventKey="gallery-details-panel">Details</Nav.Link>
            </Nav.Item>
            {/* {gallery.gallery ? (
              <Nav.Item>
                <Nav.Link eventKey="gallery-gallery-panel">Gallery</Nav.Link>
              </Nav.Item>
            ) : (
              ""
            )} */}
            <Nav.Item>
              <Nav.Link eventKey="gallery-edit-panel">Edit</Nav.Link>
            </Nav.Item>
            <Nav.Item className="ml-auto">{renderOperations()}</Nav.Item>
          </Nav>
        </div>

        <Tab.Content>
          <Tab.Pane eventKey="gallery-details-panel" title="Details">
            <GalleryDetailPanel gallery={gallery} />
          </Tab.Pane>
          {/* {gallery.gallery ? (
            <Tab.Pane eventKey="gallery-gallery-panel" title="Gallery">
              <GalleryViewer gallery={gallery.gallery} />
            </Tab.Pane>
          ) : (
            ""
          )} */}
          <Tab.Pane eventKey="gallery-edit-panel" title="Edit">
            <GalleryEditPanel
              isVisible={activeTabKey === "gallery-edit-panel"}
              gallery={gallery}
              onUpdate={(newGallery) => setGallery(newGallery)}
              onDelete={() => setIsDeleteAlertOpen(true)}
            />
          </Tab.Pane>
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
        onSelect={(k) => k && setActiveRightTabKey(k)}
      >
        <div>
          <Nav variant="tabs" className="mr-auto">
            <Nav.Item>
              <Nav.Link eventKey="images">Images</Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="add">Add</Nav.Link>
            </Nav.Item>
          </Nav>
        </div>

        <Tab.Content>
          <Tab.Pane eventKey="images" title="Images">
            {/* <GalleryViewer gallery={gallery} /> */}
            <GalleryImagesPanel gallery={gallery} />
          </Tab.Pane>
          <Tab.Pane eventKey="add" title="Add">
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

  if (isNew)
    return (
      <div className="row new-view">
        <div className="col-6">
          <h2>Create Gallery</h2>
          <GalleryEditPanel
            gallery={gallery}
            isVisible
            isNew={isNew}
            onUpdate={(newGallery) => setGallery(newGallery)}
            onDelete={() => setIsDeleteAlertOpen(true)}
          />
        </div>
      </div>
    );

  if (loading || !gallery || !data?.findGallery) {
    return <LoadingIndicator />;
  }

  if (error) return <div>{error.message}</div>;

  return (
    <div className="row">
      {maybeRenderDeleteDialog()}
      <div className="gallery-tabs order-xl-first order-last">
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
