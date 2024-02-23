import { Tab, Nav, Dropdown } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { useHistory, Link, RouteComponentProps } from "react-router-dom";
import { Helmet } from "react-helmet";
import {
  useFindImage,
  useImageIncrementO,
  useImageDecrementO,
  useImageResetO,
  useImageUpdate,
  mutateMetadataScan,
} from "src/core/StashService";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { Icon } from "src/components/Shared/Icon";
import { Counter } from "src/components/Shared/Counter";
import { useToast } from "src/hooks/Toast";
import * as Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import { OCounterButton } from "src/components/Scenes/SceneDetails/OCounterButton";
import { OrganizedButton } from "src/components/Scenes/SceneDetails/OrganizedButton";
import { ImageFileInfoPanel } from "./ImageFileInfoPanel";
import { ImageEditPanel } from "./ImageEditPanel";
import { ImageDetailPanel } from "./ImageDetailPanel";
import { DeleteImagesDialog } from "../DeleteImagesDialog";
import { faEllipsisV } from "@fortawesome/free-solid-svg-icons";
import { objectPath, objectTitle } from "src/core/files";
import { isVideo } from "src/utils/visualFile";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";

interface IProps {
  image: GQL.ImageDataFragment;
}

interface IImageParams {
  id: string;
}

const ImagePage: React.FC<IProps> = ({ image }) => {
  const history = useHistory();
  const Toast = useToast();
  const intl = useIntl();

  const [incrementO] = useImageIncrementO(image.id);
  const [decrementO] = useImageDecrementO(image.id);
  const [resetO] = useImageResetO(image.id);

  const [updateImage] = useImageUpdate();

  const [organizedLoading, setOrganizedLoading] = useState(false);

  const [activeTabKey, setActiveTabKey] = useState("image-details-panel");

  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  async function onSave(input: GQL.ImageUpdateInput) {
    await updateImage({
      variables: { input },
    });
    Toast.success(
      intl.formatMessage(
        { id: "toast.updated_entity" },
        { entity: intl.formatMessage({ id: "image" }).toLocaleLowerCase() }
      )
    );
  }

  async function onRescan() {
    if (!image || !image.visual_files.length) {
      return;
    }

    await mutateMetadataScan({
      paths: [objectPath(image)],
    });

    Toast.success(
      intl.formatMessage(
        { id: "toast.rescanning_entity" },
        {
          count: 1,
          singularEntity: intl.formatMessage({ id: "image" }),
        }
      )
    );
  }

  const onOrganizedClick = async () => {
    try {
      setOrganizedLoading(true);
      await updateImage({
        variables: {
          input: {
            id: image.id,
            organized: !image.organized,
          },
        },
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      setOrganizedLoading(false);
    }
  };

  const onIncrementClick = async () => {
    try {
      await incrementO();
    } catch (e) {
      Toast.error(e);
    }
  };

  const onDecrementClick = async () => {
    try {
      await decrementO();
    } catch (e) {
      Toast.error(e);
    }
  };

  const onResetClick = async () => {
    try {
      await resetO();
    } catch (e) {
      Toast.error(e);
    }
  };

  function onDeleteDialogClosed(deleted: boolean) {
    setIsDeleteAlertOpen(false);
    if (deleted) {
      history.push("/images");
    }
  }

  function maybeRenderDeleteDialog() {
    if (isDeleteAlertOpen && image) {
      return (
        <DeleteImagesDialog selected={[image]} onClose={onDeleteDialogClosed} />
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
          <Icon icon={faEllipsisV} />
        </Dropdown.Toggle>
        <Dropdown.Menu className="bg-secondary text-white">
          <Dropdown.Item
            key="rescan"
            className="bg-secondary text-white"
            onClick={() => onRescan()}
          >
            <FormattedMessage id="actions.rescan" />
          </Dropdown.Item>
          <Dropdown.Item
            key="delete-image"
            className="bg-secondary text-white"
            onClick={() => setIsDeleteAlertOpen(true)}
          >
            <FormattedMessage
              id="actions.delete_entity"
              values={{ entityType: intl.formatMessage({ id: "image" }) }}
            />
          </Dropdown.Item>
        </Dropdown.Menu>
      </Dropdown>
    );
  }

  function renderTabs() {
    if (!image) {
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
              <Nav.Link eventKey="image-details-panel">
                <FormattedMessage id="details" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="image-file-info-panel">
                <FormattedMessage id="file_info" />
                <Counter count={image.visual_files.length} hideZero hideOne />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link eventKey="image-edit-panel">
                <FormattedMessage id="actions.edit" />
              </Nav.Link>
            </Nav.Item>
            <Nav.Item className="ml-auto">
              <OCounterButton
                value={image.o_counter || 0}
                onIncrement={onIncrementClick}
                onDecrement={onDecrementClick}
                onReset={onResetClick}
              />
            </Nav.Item>
            <Nav.Item>
              <OrganizedButton
                loading={organizedLoading}
                organized={image.organized}
                onClick={onOrganizedClick}
              />
            </Nav.Item>
            <Nav.Item>{renderOperations()}</Nav.Item>
          </Nav>
        </div>

        <Tab.Content>
          <Tab.Pane eventKey="image-details-panel">
            <ImageDetailPanel image={image} />
          </Tab.Pane>
          <Tab.Pane
            className="file-info-panel"
            eventKey="image-file-info-panel"
          >
            <ImageFileInfoPanel image={image} />
          </Tab.Pane>
          <Tab.Pane eventKey="image-edit-panel" mountOnEnter>
            <ImageEditPanel
              isVisible={activeTabKey === "image-edit-panel"}
              image={image}
              onSubmit={onSave}
              onDelete={() => setIsDeleteAlertOpen(true)}
            />
          </Tab.Pane>
        </Tab.Content>
      </Tab.Container>
    );
  }

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("a", () => setActiveTabKey("image-details-panel"));
    Mousetrap.bind("e", () => setActiveTabKey("image-edit-panel"));
    Mousetrap.bind("f", () => setActiveTabKey("image-file-info-panel"));
    Mousetrap.bind("o", () => {
      onIncrementClick();
    });

    return () => {
      Mousetrap.unbind("a");
      Mousetrap.unbind("e");
      Mousetrap.unbind("f");
      Mousetrap.unbind("o");
    };
  });

  const title = objectTitle(image);
  const ImageView = isVideo(image.visual_files[0]) ? "video" : "img";

  return (
    <div className="row">
      <Helmet>
        <title>{title}</title>
      </Helmet>

      {maybeRenderDeleteDialog()}
      <div className="image-tabs order-xl-first order-last">
        <div className="d-none d-xl-block">
          {image.studio && (
            <h1 className="text-center">
              <Link to={`/studios/${image.studio.id}`}>
                <img
                  src={image.studio.image_path ?? ""}
                  alt={`${image.studio.name} logo`}
                  className="studio-logo"
                />
              </Link>
            </h1>
          )}
          <h3 className="image-header">{title}</h3>
        </div>
        {renderTabs()}
      </div>
      <div className="image-container">
        <ImageView
          loop={image.visual_files[0].__typename == "VideoFile"}
          autoPlay={image.visual_files[0].__typename == "VideoFile"}
          controls={image.visual_files[0].__typename == "VideoFile"}
          className="m-sm-auto no-gutter image-image"
          style={
            image.visual_files[0].__typename == "VideoFile"
              ? { width: "100%", height: "100%" }
              : {}
          }
          alt={title}
          src={image.paths.image ?? ""}
        />
      </div>
    </div>
  );
};

const ImageLoader: React.FC<RouteComponentProps<IImageParams>> = ({
  match,
}) => {
  const { id } = match.params;
  const { data, loading, error } = useFindImage(id);

  useScrollToTopOnMount();

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findImage)
    return <ErrorMessage error={`No image found with id ${id}.`} />;

  return <ImagePage image={data.findImage} />;
};

export default ImageLoader;
