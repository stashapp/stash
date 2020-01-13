import {
  Spinner,
  Tabs,
  Tab,
  Button,
  AnchorButton,
  IconName,
} from "@blueprintjs/core";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { IBaseProps } from "../../../models";
import { ErrorUtils } from "../../../utils/errors";
import { PerformerDetailsPanel } from "./PerformerDetailsPanel";
import { PerformerOperationsPanel } from "./PerformerOperationsPanel";
import { PerformerScenesPanel } from "./PerformerScenesPanel";
import { TextUtils } from "../../../utils/text";
import Lightbox from "react-images";

interface IPerformerProps extends IBaseProps {}

export const Performer: FunctionComponent<IPerformerProps> = (props: IPerformerProps) => {
  const isNew = props.match.params.id === "new";

  // Performer state
  const [performer, setPerformer] = useState<Partial<GQL.PerformerDataFragment>>({});
  const [imagePreview, setImagePreview] = useState<string | undefined>(undefined);
  const [lightboxIsOpen, setLightboxIsOpen] = useState<boolean>(false);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const { data, error, loading } = StashService.useFindPerformer(props.match.params.id);
  const updatePerformer = StashService.usePerformerUpdate();
  const createPerformer = StashService.usePerformerCreate();
  const deletePerformer = StashService.usePerformerDestroy();

  useEffect(() => {
    setIsLoading(loading);
    if (!data || !data.findPerformer || !!error) { return; }
    setPerformer(data.findPerformer);
  }, [data]);

  useEffect(() => {
    setImagePreview(performer.image_path);
  }, [performer]);

  function onImageChange(image: string) {
    setImagePreview(image);
  }

  if ((!isNew && (!data || !data.findPerformer)) || isLoading) {
    return <Spinner size={Spinner.SIZE_LARGE} />; 
  }
  if (!!error) { return <>error...</>; }

  async function onSave(performer : Partial<GQL.PerformerCreateInput> | Partial<GQL.PerformerUpdateInput>) {
    setIsLoading(true);
    try {
      if (!isNew) {
        const result = await updatePerformer({variables: performer as GQL.PerformerUpdateInput});
        setPerformer(result.data.performerUpdate);
      } else {
        const result = await createPerformer({variables: performer as GQL.PerformerCreateInput});
        setPerformer(result.data.performerCreate);
        props.history.push(`/performers/${result.data.performerCreate.id}`);
      }
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
  }

  async function onDelete() {
    setIsLoading(true);
    try {
      await deletePerformer({variables: {id: props.match.params.id}});
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
    
    // redirect to performers page
    props.history.push(`/performers`);
  }

  function renderTabs() {
    function renderEditPanel() {
      return (
        <PerformerDetailsPanel 
          performer={performer} 
          isEditing={true} 
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
        <>
          <Tabs
            renderActiveTabPanelOnly={true}
            large={true}
          >
            <Tab id="performer-details-panel" title="Details" panel={<PerformerDetailsPanel performer={performer} isEditing={false}/>} />
            <Tab id="performer-scenes-panel" title="Scenes" panel={<PerformerScenesPanel performer={performer} base={props} />} />
            <Tab id="performer-edit-panel" title="Edit" panel={renderEditPanel()} />
            <Tab id="performer-operations-panel" title="Operations" panel={<PerformerOperationsPanel performer={performer} />} />
          </Tabs>
        </>
      );
    } else {
      return renderEditPanel();
    }
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

  function setFavorite(v : boolean) {
    performer.favorite = v;
    onSave(performer);
  }

  function renderIcons() {
    function maybeRenderURL(url?: string, icon?: IconName) {
      if (performer.url) {
        if (!icon) {
          icon = "link";
        }

        return (
          <>
            <AnchorButton
              icon={icon}
              href={performer.url}
              minimal={true}
            />
          </>
        )
      }
    }

    return (
      <>
        <span className="name-icons">
          <Button
            icon="heart"
            className={performer.favorite ? "favorite" : "not-favorite"}
            onClick={() => setFavorite(!performer.favorite)}
            minimal={true}
          />
          {maybeRenderURL(performer.url)}
          {/* TODO - render instagram and twitter links with icons */}
        </span>
      </>
    );
  }

  function renderNewView() {
    return (
      <div className="columns is-multiline no-spacing">
        <div className="column is-half details-image-container">
          <img alt="Performer" className="performer" src={imagePreview} />
        </div>
        <div className="column is-half details-detail-container">
          {renderTabs()}
        </div>
      </div>
    );
  }

  const photos = [{src: imagePreview || "", caption: "Image"}];

  function openLightbox() {
    setLightboxIsOpen(true);
  }

  function closeLightbox() {
    setLightboxIsOpen(false);
  }

  if (isNew) {
    return renderNewView();
  }

  return (
    <>
      <div id="performer-page">
        <div className="details-image-container">
          <img alt={performer.name} className="performer" src={imagePreview} onClick={openLightbox} />
        </div>
        <div className="performer-head">
          <h1 className="bp3-heading">
            {performer.name}
            {renderIcons()}
          </h1>
          {maybeRenderAliases()}
          {maybeRenderAge()}
        </div>
        
        <div className="performer-body">
          <div className="details-detail-container">
            {renderTabs()}
          </div>
        </div>
      </div>
      <Lightbox
        images={photos}
        onClose={closeLightbox}
        currentImage={0}
        isOpen={lightboxIsOpen}
        onClickImage={() => window.open(imagePreview, "_blank")}
        width={9999}
      />
    </>
  );
};
