/* eslint-disable react/no-this-in-sfc */

import React, { useEffect, useState } from "react";
import { Button, Tabs, Tab } from "react-bootstrap";
import { useParams, useHistory } from "react-router-dom";
import cx from "classnames";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { Icon, LoadingIndicator } from "src/components/Shared";
import { useToast } from "src/hooks";
import { TextUtils } from "src/utils";
import Lightbox from "react-images";
import { PerformerDetailsPanel } from "./PerformerDetailsPanel";
import { PerformerOperationsPanel } from "./PerformerOperationsPanel";
import { PerformerScenesPanel } from "./PerformerScenesPanel";

export const Performer: React.FC = () => {
  const Toast = useToast();
  const history = useHistory();
  const { id = "new" } = useParams();
  const isNew = id === "new";

  // Performer state
  const [performer, setPerformer] = useState<
    Partial<GQL.PerformerDataFragment>
  >({});
  const [imagePreview, setImagePreview] = useState<string>();
  const [lightboxIsOpen, setLightboxIsOpen] = useState(false);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const { data, error } = StashService.useFindPerformer(id);
  const [updatePerformer] = StashService.usePerformerUpdate();
  const [createPerformer] = StashService.usePerformerCreate();
  const [deletePerformer] = StashService.usePerformerDestroy();

  useEffect(() => {
    setIsLoading(false);
    if (data?.findPerformer) setPerformer(data.findPerformer);
  }, [data]);

  useEffect(() => {
    setImagePreview(performer.image_path ?? undefined);
  }, [performer]);

  function onImageChange(image: string) {
    setImagePreview(image);
  }

  if ((!isNew && (!data || !data.findPerformer)) || isLoading)
    return <LoadingIndicator />;

  if (error) return <div>{error.message}</div>;

  async function onSave(
    performerInput:
      | Partial<GQL.PerformerCreateInput>
      | Partial<GQL.PerformerUpdateInput>
  ) {
    setIsLoading(true);
    try {
      if (!isNew) {
        const result = await updatePerformer({
          variables: performerInput as GQL.PerformerUpdateInput,
        });
        if (result.data?.performerUpdate)
          setPerformer(result.data?.performerUpdate);
      } else {
        const result = await createPerformer({
          variables: performerInput as GQL.PerformerCreateInput,
        });
        if (result.data?.performerCreate) {
          setPerformer(result.data.performerCreate);
          history.push(`/performers/${result.data.performerCreate.id}`);
        }
      }
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  async function onDelete() {
    setIsLoading(true);
    try {
      await deletePerformer({ variables: { id } });
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);

    // redirect to performers page
    history.push("/performers");
  }

  function renderTabs() {
    function renderEditPanel() {
      return (
        <PerformerDetailsPanel
          performer={performer}
          isEditing
          isNew={isNew}
          onDelete={onDelete}
          onSave={onSave}
          onImageChange={onImageChange}
        />
      );
    }

    // render tabs if not new
    if (!isNew) {
      return (
        <Tabs defaultActiveKey="details" id="performer-details">
          <Tab eventKey="details" title="Details">
            <PerformerDetailsPanel performer={performer} isEditing={false} />
          </Tab>
          <Tab eventKey="scenes" title="Scenes">
            <PerformerScenesPanel performer={performer} />
          </Tab>
          <Tab eventKey="edit" title="Edit">
            {renderEditPanel()}
          </Tab>
          <Tab eventKey="operations" title="Operations">
            <PerformerOperationsPanel performer={performer} />
          </Tab>
        </Tabs>
      );
    }
    return renderEditPanel();
  }

  function maybeRenderAge() {
    if (performer && performer.birthdate) {
      // calculate the age from birthdate. In future, this should probably be
      // provided by the server
      return (
        <>
          <div>
            <span className="age">{TextUtils.age(performer.birthdate)}</span>
            <span className="age-tail"> years old</span>
          </div>
        </>
      );
    }
  }

  function maybeRenderAliases() {
    if (performer && performer.aliases) {
      return (
        <>
          <div>
            <span className="alias-head">Also known as </span>
            <span className="alias">{performer.aliases}</span>
          </div>
        </>
      );
    }
  }

  function setFavorite(v: boolean) {
    performer.favorite = v;
    onSave(performer);
  }

  const renderIcons = () => (
    <span className="name-icons d-block d-sm-inline">
      <Button
        className={cx(
          "minimal",
          performer.favorite ? "favorite" : "not-favorite"
        )}
        onClick={() => setFavorite(!performer.favorite)}
      >
        <Icon icon="heart" />
      </Button>
      {performer.url && (
        <Button className="minimal">
          <a
            href={performer.url}
            className="link"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Icon icon="link" />
          </a>
        </Button>
      )}
      {performer.twitter && (
        <Button className="minimal">
          <a
            href={`https://www.twitter.com/${performer.twitter}`}
            className="twitter"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Icon icon="dove" />
          </a>
        </Button>
      )}
      {performer.instagram && (
        <Button className="minimal">
          <a
            href={`https://www.instagram.com/${performer.instagram}`}
            className="instagram"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Icon icon="camera" />
          </a>
        </Button>
      )}
    </span>
  );

  function renderNewView() {
    return (
      <div className="row new-view">
        <div className="col-4">
          <img className="photo" src={imagePreview} alt="Performer" />
        </div>
        <div className="col-6">
          <h2>Create Performer</h2>
          {renderTabs()}
        </div>
      </div>
    );
  }

  const photos = [{ src: imagePreview || "", caption: "Image" }];

  if (isNew) {
    return renderNewView();
  }

  return (
    <div id="performer-page" className="row">
      <div className="image-container col-sm-4 offset-sm-1 d-none d-sm-block">
        <Button variant="link" onClick={() => setLightboxIsOpen(true)}>
          <img className="performer" src={imagePreview} alt="Performer" />
        </Button>
      </div>
      <div className="col col-sm-6">
        <div className="row">
          <div className="image-container col-6 d-block d-sm-none">
            <Button variant="link" onClick={() => setLightboxIsOpen(true)}>
              <img className="performer" src={imagePreview} alt="Performer" />
            </Button>
          </div>
          <div className="performer-head col-6 col-sm-12">
            <h2>
              {performer.name}
              {renderIcons()}
            </h2>
            {maybeRenderAliases()}
            {maybeRenderAge()}
          </div>
        </div>
        <div className="performer-body">
          <div className="performer-tabs">{renderTabs()}</div>
        </div>
      </div>
      <Lightbox
        images={photos}
        onClose={() => setLightboxIsOpen(false)}
        currentImage={0}
        isOpen={lightboxIsOpen}
        onClickImage={() => window.open(imagePreview, "_blank")}
        width={9999}
      />
    </div>
  );
};
