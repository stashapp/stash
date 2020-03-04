/* eslint-disable react/no-this-in-sfc */

import { Table } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { ImageUtils, TableUtils } from "src/utils";
import {
  DetailsEditNavbar,
  Modal,
  LoadingIndicator
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { StudioScenesPanel } from "./StudioScenesPanel";

export const Studio: React.FC = () => {
  const history = useHistory();
  const Toast = useToast();
  const { id = "new" } = useParams();
  const isNew = id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing studio state
  const [image, setImage] = useState<string>();
  const [name, setName] = useState<string>();
  const [url, setUrl] = useState<string>();

  // Studio state
  const [studio, setStudio] = useState<Partial<GQL.StudioDataFragment>>({});
  const [imagePreview, setImagePreview] = useState<string>();

  const { data, error, loading } = StashService.useFindStudio(id);
  const [updateStudio] = StashService.useStudioUpdate(
    getStudioInput() as GQL.StudioUpdateInput
  );
  const [createStudio] = StashService.useStudioCreate(
    getStudioInput() as GQL.StudioCreateInput
  );
  const [deleteStudio] = StashService.useStudioDestroy(
    getStudioInput() as GQL.StudioDestroyInput
  );

  function updateStudioEditState(state: Partial<GQL.StudioDataFragment>) {
    setName(state.name);
    setUrl(state.url ?? undefined);
  }

  function updateStudioData(studioData: Partial<GQL.StudioDataFragment>) {
    setImage(undefined);
    updateStudioEditState(studioData);
    setImagePreview(studioData.image_path ?? undefined);
    setStudio(studioData);
  }

  useEffect(() => {
    if (data && data.findStudio) {
      setImage(undefined);
      updateStudioEditState(data.findStudio);
      setImagePreview(data.findStudio.image_path ?? undefined);
      setStudio(data.findStudio);
    }
  }, [data]);

  function onImageLoad(this: FileReader) {
    setImagePreview(this.result as string);
    setImage(this.result as string);
  }

  ImageUtils.usePasteImage(onImageLoad);

  if (!isNew && !isEditing) {
    if (!data?.findStudio || loading) return <LoadingIndicator />;
    if (error) return <div>{error.message}</div>;
  }

  function getStudioInput() {
    const input: Partial<GQL.StudioCreateInput | GQL.StudioUpdateInput> = {
      name,
      url,
      image
    };

    if (!isNew) {
      (input as GQL.StudioUpdateInput).id = id;
    }
    return input;
  }

  async function onSave() {
    try {
      if (!isNew) {
        const result = await updateStudio();
        if (result.data?.studioUpdate) {
          updateStudioData(result.data.studioUpdate);
          setIsEditing(false);
        }
      } else {
        const result = await createStudio();
        if (result.data?.studioCreate?.id) {
          history.push(`/studios/${result.data.studioCreate.id}`);
          setIsEditing(false);
        }
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onAutoTag() {
    if (!studio.id) return;
    try {
      await StashService.queryMetadataAutoTag({ studios: [studio.id] });
      Toast.success({ content: "Started auto tagging" });
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

  function onImageChangeHandler(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function renderDeleteAlert() {
    return (
      <Modal
        show={isDeleteAlertOpen}
        icon="trash-alt"
        accept={{ text: "Delete", variant: "danger", onClick: onDelete }}
        cancel={{ onClick: () => setIsDeleteAlertOpen(false) }}
      >
        <p>Are you sure you want to delete {studio.name ?? "studio"}?</p>
      </Modal>
    );
  }

  return (
    <div className="row">
      <div
        className={cx("studio-details", {
          "col ml-sm-5": !isNew,
          "col-8": isNew
        })}
      >
        {isNew && <h2>Add Studio</h2>}
        <img className="logo w-100" alt={name} src={imagePreview} />
        <Table>
          <tbody>
            {TableUtils.renderInputGroup({
              title: "Name",
              value: studio.name ?? "",
              isEditing: !!isEditing,
              onChange: setName
            })}
            {TableUtils.renderInputGroup({
              title: "URL",
              value: url,
              isEditing: !!isEditing,
              onChange: setUrl
            })}
          </tbody>
        </Table>
        <DetailsEditNavbar
          studio={studio}
          isNew={isNew}
          isEditing={isEditing}
          onToggleEdit={() => setIsEditing(!isEditing)}
          onSave={onSave}
          onImageChange={onImageChangeHandler}
          onAutoTag={onAutoTag}
          onDelete={onDelete}
        />
      </div>
      {!isNew && (
        <div className="col-12 col-sm-8">
          <StudioScenesPanel studio={studio} />
        </div>
      )}
      {renderDeleteAlert()}
    </div>
  );
};
