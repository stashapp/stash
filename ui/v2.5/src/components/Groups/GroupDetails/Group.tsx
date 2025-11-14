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
import { useHistory, RouteComponentProps, Redirect } from "react-router-dom";
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
import { faRefresh, faTrashAlt } from "@fortawesome/free-solid-svg-icons";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { useConfigurationContext } from "src/hooks/Config";
import { DetailImage } from "src/components/Shared/DetailImage";
import { useRatingKeybinds } from "src/hooks/keybinds";
import { useLoadStickyHeader } from "src/hooks/detailsPanel";
import { useScrollToTopOnMount } from "src/hooks/scrollToTop";
import { ExternalLinkButtons } from "src/components/Shared/ExternalLinksButton";
import { BackgroundImage } from "src/components/Shared/DetailsPage/BackgroundImage";
import { DetailTitle } from "src/components/Shared/DetailsPage/DetailTitle";
import { ExpandCollapseButton } from "src/components/Shared/CollapseButton";
import { AliasList } from "src/components/Shared/DetailsPage/AliasList";
import { HeaderImage } from "src/components/Shared/DetailsPage/HeaderImage";
import { LightboxLink } from "src/hooks/Lightbox/LightboxLink";
import {
  TabTitleCounter,
  useTabKey,
} from "src/components/Shared/DetailsPage/Tabs";
import { Button, Tab, Tabs } from "react-bootstrap";
import { GroupSubGroupsPanel } from "./GroupSubGroupsPanel";
import { GroupPerformersPanel } from "./GroupPerformersPanel";
import { Icon } from "src/components/Shared/Icon";
import { goBackOrReplace } from "src/utils/history";

const validTabs = ["default", "scenes", "performers", "subgroups"] as const;
type TabKey = (typeof validTabs)[number];

function isTabKey(tab: string): tab is TabKey {
  return validTabs.includes(tab as TabKey);
}

const GroupTabs: React.FC<{
  tabKey?: TabKey;
  group: GQL.GroupDataFragment;
  abbreviateCounter: boolean;
}> = ({ tabKey, group, abbreviateCounter }) => {
  const {
    scene_count: sceneCount,
    performer_count: performerCount,
    sub_group_count: groupCount,
  } = group;

  const populatedDefaultTab = useMemo(() => {
    if (sceneCount == 0) {
      if (performerCount != 0) {
        return "performers";
      } else if (groupCount !== 0) {
        return "subgroups";
      }
    }

    return "scenes";
  }, [sceneCount, performerCount, groupCount]);

  const { setTabKey } = useTabKey({
    tabKey,
    validTabs,
    defaultTabKey: populatedDefaultTab,
    baseURL: `/groups/${group.id}`,
  });

  return (
    <Tabs
      id="group-tabs"
      mountOnEnter
      unmountOnExit
      activeKey={tabKey}
      onSelect={setTabKey}
    >
      <Tab
        eventKey="scenes"
        title={
          <TabTitleCounter
            messageID="scenes"
            count={sceneCount}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <GroupScenesPanel active={tabKey === "scenes"} group={group} />
      </Tab>
      <Tab
        eventKey="performers"
        title={
          <TabTitleCounter
            messageID="performers"
            count={performerCount}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <GroupPerformersPanel active={tabKey === "performers"} group={group} />
      </Tab>
      <Tab
        eventKey="subgroups"
        title={
          <TabTitleCounter
            messageID="sub_groups"
            count={groupCount}
            abbreviateCounter={abbreviateCounter}
          />
        }
      >
        <GroupSubGroupsPanel active={tabKey === "subgroups"} group={group} />
      </Tab>
    </Tabs>
  );
};

interface IProps {
  group: GQL.GroupDataFragment;
  tabKey?: TabKey;
}

interface IGroupParams {
  id: string;
  tab?: string;
}

const GroupPage: React.FC<IProps> = ({ group, tabKey }) => {
  const intl = useIntl();
  const history = useHistory();
  const Toast = useToast();

  // Configuration settings
  const { configuration } = useConfigurationContext();
  const uiConfig = configuration?.ui;
  const enableBackgroundImage = uiConfig?.enableMovieBackgroundImage ?? false;
  const compactExpandedDetails = uiConfig?.compactExpandedDetails ?? false;
  const showAllDetails = uiConfig?.showAllDetails ?? true;
  const abbreviateCounter = uiConfig?.abbreviateCounters ?? false;

  const [focusedOnFront, setFocusedOnFront] = useState<boolean>(true);

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
      return;
    }

    goBackOrReplace(history, "/groups");
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
          show={enableBackgroundImage && !isEditing}
        />
        <div className="detail-container">
          <HeaderImage encodingImage={encodingImage}>
            <div className="group-images">
              {!!activeFrontImage && (
                <LightboxLink images={lightboxImages}>
                  <DetailImage
                    className={`front-cover ${
                      focusedOnFront ? "active" : "inactive"
                    }`}
                    alt="Front Cover"
                    src={activeFrontImage}
                  />
                </LightboxLink>
              )}
              {!!activeBackImage && (
                <LightboxLink
                  images={lightboxImages}
                  index={lightboxImages.length - 1}
                >
                  <DetailImage
                    className={`back-cover ${
                      !focusedOnFront ? "active" : "inactive"
                    }`}
                    alt="Back Cover"
                    src={activeBackImage}
                  />
                </LightboxLink>
              )}
              {!!(activeFrontImage && activeBackImage) && (
                <Button
                  className="flip"
                  onClick={() => setFocusedOnFront(!focusedOnFront)}
                >
                  <Icon icon={faRefresh} />
                </Button>
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
                  <ExternalLinkButtons urls={group.urls ?? undefined} />
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
          <div className="group-tabs">
            {!isEditing && (
              <GroupTabs
                group={group}
                tabKey={tabKey}
                abbreviateCounter={abbreviateCounter}
              />
            )}
          </div>
        </div>
      </div>
      {renderDeleteAlert()}
    </div>
  );
};

const GroupLoader: React.FC<RouteComponentProps<IGroupParams>> = ({
  location,
  match,
}) => {
  const { id, tab } = match.params;
  const { data, loading, error } = useFindGroup(id);

  useScrollToTopOnMount();

  if (tab && !isTabKey(tab)) {
    return (
      <Redirect
        to={{
          ...location,
          pathname: `/groups/${id}`,
        }}
      />
    );
  }

  if (loading) return <LoadingIndicator />;
  if (error) return <ErrorMessage error={error.message} />;
  if (!data?.findGroup)
    return <ErrorMessage error={`No group found with id ${id}.`} />;

  return (
    <GroupPage group={data.findGroup} tabKey={tab as TabKey | undefined} />
  );
};

export default GroupLoader;
