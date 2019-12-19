import {
  Spinner,
  Tabs,
  Tab,
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { IBaseProps } from "../../../models";
import { ErrorUtils } from "../../../utils/errors";
import { PerformerDetailsPanel } from "./PerformerDetailsPanel";
import { PerformerOperationsPanel } from "./PerformerOperationsPanel";
import { PerformerScenesPanel } from "./PerformerScenesPanel";

interface IPerformerProps extends IBaseProps {}

export const Performer: FunctionComponent<IPerformerProps> = (props: IPerformerProps) => {
  const isNew = props.match.params.id === "new";

  // Performer state
  const [performer, setPerformer] = useState<Partial<GQL.PerformerDataFragment>>({});
  const [imagePreview, setImagePreview] = useState<string | undefined>(undefined);

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
            <Tab id="performer-edit-panel" title="Edit" panel={renderEditPanel()} />
            <Tab id="performer-scenes-panel" title="Scenes" panel={<PerformerScenesPanel performer={performer} base={props} />} />
            <Tab id="performer-operations-panel" title="Operations" panel={<PerformerOperationsPanel performer={performer} />} />
          </Tabs>
        </>
      );
    } else {
      return renderEditPanel();
    }
  }

  return (
    <>
      <div className="columns is-multiline no-spacing">
        <div className="column is-half details-image-container">
          <img className="performer" src={imagePreview} />
        </div>
        <div className="column is-half details-detail-container">
          {renderTabs()}
        </div>
      </div>
    </>
  );
};
