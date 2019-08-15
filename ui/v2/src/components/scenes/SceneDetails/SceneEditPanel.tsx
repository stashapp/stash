import {
  Alert,
  Button,
  FormGroup,
  HTMLSelect,
  InputGroup,
  Spinner,
  TextArea,
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { ErrorUtils } from "../../../utils/errors";
import { ToastUtils } from "../../../utils/toasts";
import { FilterMultiSelect } from "../../select/FilterMultiSelect";
import { FilterSelect } from "../../select/FilterSelect";
import { ValidGalleriesSelect } from "../../select/ValidGalleriesSelect";

interface IProps {
  scene: GQL.SceneDataFragment;
  onUpdate: (scene: GQL.SceneDataFragment) => void;
  onDelete: () => void;
}

export const SceneEditPanel: FunctionComponent<IProps> = (props: IProps) => {
  // Editing scene state
  const [title, setTitle] = useState<string | undefined>(undefined);
  const [details, setDetails] = useState<string | undefined>(undefined);
  const [url, setUrl] = useState<string | undefined>(undefined);
  const [date, setDate] = useState<string | undefined>(undefined);
  const [rating, setRating] = useState<number | undefined>(undefined);
  const [galleryId, setGalleryId] = useState<string | undefined>(undefined);
  const [studioId, setStudioId] = useState<string | undefined>(undefined);
  const [performerIds, setPerformerIds] = useState<string[] | undefined>(undefined);
  const [tagIds, setTagIds] = useState<string[] | undefined>(undefined);

  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);
  const [deleteFile, setDeleteFile] = useState<boolean>(false);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const updateScene = StashService.useSceneUpdate(getSceneInput());
  const deleteScene = StashService.useSceneDestroy(getSceneDeleteInput());

  function updateSceneEditState(state: Partial<GQL.SceneDataFragment>) {
    const perfIds = !!state.performers ? state.performers.map((performer) => performer.id) : undefined;
    const tIds = !!state.tags ? state.tags.map((tag) => tag.id) : undefined;

    setTitle(state.title);
    setDetails(state.details);
    setUrl(state.url);
    setDate(state.date);
    setRating(state.rating);
    setGalleryId(state.gallery ? state.gallery.id : undefined);
    setStudioId(state.studio ? state.studio.id : undefined);
    setPerformerIds(perfIds);
    setTagIds(tIds);
  }

  useEffect(() => {
    updateSceneEditState(props.scene);
  }, [props.scene]);

  // if (!isNew && !isEditing) {
  //   if (!data || !data.findPerformer || isLoading) { return <Spinner size={Spinner.SIZE_LARGE} />; }
  //   if (!!error) { return <>error...</>; }
  // }

  function getSceneInput(): GQL.SceneUpdateInput {
    return {
      id: props.scene.id,
      title,
      details,
      url,
      date,
      rating,
      gallery_id: galleryId,
      studio_id: studioId,
      performer_ids: performerIds,
      tag_ids: tagIds,
    };
  }

  async function onSave() {
    setIsLoading(true);
    try {
      const result = await updateScene();
      props.onUpdate(result.data.sceneUpdate);
      ToastUtils.success("Updated scene");
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
  }

  function getSceneDeleteInput(): GQL.SceneDestroyInput {
    return {
      id: props.scene.id,
      delete_file: deleteFile
    };
  }

  async function onDelete() {
    setIsDeleteAlertOpen(false);
    setIsLoading(true);
    try {
      await deleteScene();
      ToastUtils.success("Deleted scene");
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);

    props.onDelete();
  }

  function renderMultiSelect(type: "performers" | "tags", initialIds: string[] | undefined) {
    return (
      <FilterMultiSelect
        type={type}
        onUpdate={(items) => {
          const ids = items.map((i) => i.id);
          switch (type) {
            case "performers": setPerformerIds(ids); break;
            case "tags": setTagIds(ids); break;
          }
        }}
        initialIds={initialIds}
      />
    );
  }

  function renderDeleteAlert() {
    return (
      <Alert
        cancelButtonText="Cancel"
        confirmButtonText="Delete"
        icon="trash"
        intent="danger"
        isOpen={isDeleteAlertOpen}
        onCancel={() => setIsDeleteAlertOpen(false)}
        onConfirm={() => onDelete()}
      >
        <p>
          Are you sure you want to delete this scene? Unless the file is also deleted, this scene will be re-added when scan is performed.
        </p>
      </Alert>
    );
  }

  return (
    <>
      {renderDeleteAlert()}
      {isLoading ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
      <div className="form-container " style={{width: "50%"}}>
        <FormGroup label="Title">
          <InputGroup
            onChange={(newValue: any) => setTitle(newValue.target.value)}
            value={title}
          />
        </FormGroup>

        <FormGroup label="Details">
          <TextArea
            fill={true}
            onChange={(newValue) => setDetails(newValue.target.value)}
            value={details}
          />
        </FormGroup>

        <FormGroup label="URL">
          <InputGroup
            onChange={(newValue: any) => setUrl(newValue.target.value)}
            value={url}
          />
        </FormGroup>

        <FormGroup label="Date" helperText="YYYY-MM-DD">
          <InputGroup
            onChange={(newValue: any) => setDate(newValue.target.value)}
            value={date}
          />
        </FormGroup>

        <FormGroup label="Rating">
          <HTMLSelect
            options={["", 1, 2, 3, 4, 5]}
            onChange={(event) => setRating(parseInt(event.target.value, 10))}
            value={rating}
          />
        </FormGroup>

        <FormGroup label="Gallery">
          <ValidGalleriesSelect
            sceneId={props.scene.id}
            initialId={galleryId}
            onSelectItem={(item) => setGalleryId(item.id)}
          />
        </FormGroup>

        <FormGroup label="Studio">
          <FilterSelect
            type="studios"
            onSelectItem={(item) => setStudioId(item.id)}
            initialId={studioId}
          />
        </FormGroup>

        <FormGroup label="Performers">
          {renderMultiSelect("performers", performerIds)}
        </FormGroup>

        <FormGroup label="Tags">
          {renderMultiSelect("tags", tagIds)}
        </FormGroup>
      </div>
      <Button text="Save" intent="primary" onClick={() => onSave()}/>
      <Button text="Delete" intent="danger" onClick={() => setIsDeleteAlertOpen(true)}/>
    </>
  );
};
