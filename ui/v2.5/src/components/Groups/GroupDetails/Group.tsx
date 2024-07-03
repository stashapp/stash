import React, { useEffect, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Helmet } from "react-helmet";
import cx from "classnames";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  useFindGroup,
  useGroupUpdate,
  useGroupDestroy,
} from "src/core/StashService";
import { useHistory, RouteComponentProps } from "react-router-dom";
import { DetailsEditNavbar } from "src/components/Shared/DetailsEditNavbar";
import { ErrorMessage } from "src/components/Shared/ErrorMessage";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { ModalComponent } from "src/components/Shared/Modal";
import { useToast } from "src/hooks/Toast";
import { GroupScenesPanel } from "./GroupScenesPanel";
import {
  CompressedGroupDetailsPanel,
  GroupDetailsPanel,
} from "./GroupDetailsPanel";
import { GroupEditPanel } from "./GroupEditPanel";
import { faTrashAlt } from "@fortawesome/free-solid-svg-icons";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { ConfigurationContext } from "src/hooks/Config";
import { DetailImage } from "src/components/Shared/DetailImage";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { useLoadStickyHeader } from "src/hooks/detailsPanel";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";
import { ExternalLinksButton } from "src/components/Shared/ExternalLinksButton";
import { BackgroundImage } from "src/components/Shared/DetailsPage/BackgroundImage";
import { DetailTitle } from "src/components/Shared/DetailsPage/DetailTitle";
import { ExpandCollapseButton } from "src/components/Shared/CollapseButton";
import { AliasList } from "src/components/Shared/DetailsPage/AliasList";
import { HeaderImage } from "src/components/Shared/DetailsPage/HeaderImage";
import { LightboxLink } from "src/hooks/Lightbox/LightboxLink";

interface IProps {
  group: GQL.GroupDataFragment;
}

interface IGroupParams {
  id: string;
}

const GroupPage: React.FC<IProps> = ({ group }) => {
  const intl = useIntl();
  const history = useHistory();
  const Toast = useToast();

  // Configuration settings
  const { configuration } = React.useContext(ConfigurationContext);
  const uiConfig = configuration?.ui;
  const enableBackgroundImage = uiConfig?.enableMovieBackgroundImage ?? false;
  const compactExpandedDetails = uiConfig?.compactExpandedDetails ?? false;
  const showAllDetails = uiConfig?.showAllDetails ?? true;

  const [collapsed, setCollapsed] = useState<boolean>(!showAllDetails);
  const loadStickyHeader = useLoadStickyHeader();

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing group state
  const [frontImage, setFrontImage] = useState<string | null>();
  const [backImage, setBackImage] = useState<string | null>();
  const [encodingImage, setEncodingImage] = useState<boolean>(false);

  const aliases = useMemo(
    () => (group.aliases ? [group.aliases] : []),
    [group.aliases]
  );

  const isDefaultImage =
    group.front_image_path && group.front_image_path.includes("default=true");

  const lightboxImages = useMemo(() => {
    const covers = [];

    if (group.front_image_path && !isDefaultImage) {
      covers.push({
        paths: {
          thumbnail: group.front_image_path,
          image: group.front_image_path,
        },
      });
    }

    if (group.back_image_path) {
      covers.push({
        paths: {
          thumbnail: group.back_image_path,
          image: group.back_image_path,
        },
      });
    }
    return covers;
  }, [group.front_image_path, group.back_image_path, isDefaultImage]);

  const activeFrontImage = useMemo(() => {
    let existingImage = group.front_image_path;
    if (isEditing) {
      if (frontImage === null && existingImage) {
        const imageURL = new URL(existingImage);
        imageURL.searchParams.set("default", "true");
        return imageURL.toString();
      } else if (frontImage) {
        return frontImage;
      }
    }

    return existingImage;
  }, [isEditing, group.front_image_path, frontImage]);

  const activeBackImage = useMemo(() => {
    let existingImage = group.back_image_path;
    if (isEditing) {
      if (backImage === null) {
        return undefined;
      } else if (backImage) {
        return backImage;
      }
    }

    return existingImage;
  }, [isEditing, group.back_image_path, backImage]);

  const [updateGroup, { loading: updating }] = useGroupUpdate();
  const [deleteGroup, { loading: deleting }] = useGroupDestroy({
    id: group.id,
  });

  // set up hotkeys
  useEffect(() => {
    Mousetrap.bind("e", () => toggleEditing());
    Mousetrap.bind("d d", () => {
      setIsDeleteAlertOpen(true);
    });
    Mousetrap.bind(",", () => setCollapsed(!collapsed));

    return () => {
      Mousetrap.unbind("e");
      Mousetrap.unbind("d d");
    };
  });

  useRatingKeybinds(
    true,
    configuration?.ui.ratingSystemOptions?.type,
    setRating
  );

  async function onSave(input: GQL.GroupCreateInput) {
    await updateGroup({
      variables: {
        input: {
          id: group.id,
          ...input,
        },
      },
    });
    toggleEditing(false);
    Toast.success(
      intl.formatMessage(
        { id: "toast.updated_entity" },
        { entity: intl.formatMessage({ id: "group" }).toLocaleLowerCase() }
      )
    );
  }

  async function onDelete() {
    try {
      await deleteGroup();
    } catch (e) {
      Toast.error(e);
    }

    // redirect to groups page
    history.push(`/groups`);
  }

  function toggleEditing(value?: boolean) {
    if (value !== undefined) {
      setIsEditing(value);
    } else {
      setIsEditing((e) => !e);
    }
    setFrontImage(undefined);
    setBackImage(undefined);
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
                group.name ??
                intl.formatMessage({ id: "group" }).toLocaleLowerCase(),
            }}
          />
        </p>
      </ModalComponent>
    );
  }

  function setRating(v: number | null) {
    if (group.id) {
      updateGroup({
        variables: {
          input: {
            id: group.id,
            rating100: v,
          },
        },
      });
    }
  }

  const renderTabs = () => <GroupScenesPanel active={true} group={group} />;

  function maybeRenderTab() {
    if (!isEditing) {
      return renderTabs();
    }
  }

  if (updating || deleting) return <LoadingIndicator />;

  const headerClassName = cx("detail-header", {
    edit: isEditing,
    collapsed,
    "full-width": !collapsed && !compactExpandedDetails,
  });

  return (
    <div id="group-page" className="row">
      <Helmet>
        <title>{group?.name}</title>
      </Helmet>

      <div className={headerClassName}>
        <BackgroundImage
          imagePath={group.front_image_path ?? undefined}
          show={!enableBackgroundImage && !isEditing}
        />
        <div className="detail-container">
          <HeaderImage encodingImage={encodingImage}>
            <div className="group-images">
              {!!activeFrontImage && (
                <LightboxLink images={lightboxImages}>
                  <DetailImage alt="Front Cover" src={activeFrontImage} />
                </LightboxLink>
              )}
              {!!activeBackImage && (
                <LightboxLink
                  images={lightboxImages}
                  index={lightboxImages.length - 1}
                >
                  <DetailImage alt="Back Cover" src={activeBackImage} />
                </LightboxLink>
              )}
            </div>
          </HeaderImage>
          <div className="row">
            <div className="group-head col">
              <DetailTitle name={group.name} classNamePrefix="group">
                {!isEditing && (
                  <ExpandCollapseButton
                    collapsed={collapsed}
                    setCollapsed={(v) => setCollapsed(v)}
                  />
                )}
                <span className="name-icons">
                  <ExternalLinksButton urls={group.urls} />
                </span>
              </DetailTitle>

              <AliasList aliases={aliases} />
              <RatingSystem
                value={group.rating100}
                onSetRating={(value) => setRating(value)}
                clickToRate
                withoutContext
              />
              {!isEditing && (
                <GroupDetailsPanel
                  group={group}
                  collapsed={collapsed}
                  fullWidth={!collapsed && !compactExpandedDetails}
                />
              )}
              {isEditing ? (
                <GroupEditPanel
                  group={group}
                  onSubmit={onSave}
                  onCancel={() => toggleEditing()}
                  onDelete={onDelete}
                  setFrontImage={setFrontImage}
                  setBackImage={setBackImage}
                  setEncodingImage={setEncodingImage}
                />
              ) : (
                <DetailsEditNavbar
                  objectName={group.name}
                  isNew={false}
                  isEditing={isEditing}
                  onToggleEdit={() => toggleEditing()}
                  onSave={() => {}}
                  onImageChange={() => {}}
                  onDelete={onDelete}
                />
              )}
            </div>
          </div>
        </div>
      </div>

      {!isEditing && loadStickyHeader && (
        <CompressedGroupDetailsPanel group={group} />
      )}

      <div className="detail-body">
        <div className="group-body">
          <div className="group-tabs">{maybeRenderTab()}</div>
        </div>
      </div>
      {renderDeleteAlert()}
    </div>
  );
};

const GroupLoader: React.FC<RouteComponentProps<IGroupParams>> = ({
  match,
}) => {
  const { id } = match.params;
  const { data, loading, error } = useFindGroup(id);

  useScrollToTopOnMount();

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findGroup)
    return <ErrorMessage error={`No group found with id ${id}.`} />;

  return <GroupPage group={data.findGroup} />;
};

export default GroupLoader;
