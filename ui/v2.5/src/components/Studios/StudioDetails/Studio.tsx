/* eslint-disable react/no-this-in-sfc */

import { Form, Spinner, Table } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { useParams, useHistory } from "react-router-dom";

import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { ImageUtils, TableUtils } from "src/utils";
import { DetailsEditNavbar } from "src/components/Shared";
import { useToast } from "src/hooks";

export const Studio: React.FC = () => {
  const history = useHistory();
  const Toast = useToast();
  const { id = "new" } = useParams();
  const isNew = id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);

  // Editing studio state
  const [image, setImage] = useState<string | undefined>(undefined);
  const [name, setName] = useState<string | undefined>(undefined);
  const [url, setUrl] = useState<string | undefined>(undefined);

  // Studio state
  const [studio, setStudio] = useState<Partial<GQL.StudioDataFragment>>({});
  const [imagePreview, setImagePreview] = useState<string | undefined>(
    undefined
  );

  const { data, error, loading } = StashService.useFindStudio(id);
  const updateStudio = StashService.useStudioUpdate(
    getStudioInput() as GQL.StudioUpdateInput
  );
  const createStudio = StashService.useStudioCreate(
    getStudioInput() as GQL.StudioCreateInput
  );
  const deleteStudio = StashService.useStudioDestroy(
    getStudioInput() as GQL.StudioDestroyInput
  );

  function updateStudioEditState(state: Partial<GQL.StudioDataFragment>) {
    setName(state.name);
    setUrl(state.url);
  }

  function updateStudioData(studioData: Partial<GQL.StudioDataFragment>) {
    setImage(undefined);
    updateStudioEditState(studioData);
    setImagePreview(studioData.image_path);
    setStudio(studioData);
  }

  useEffect(() => {
    if (data && data.findStudio) {
      setImage(undefined);
      updateStudioEditState(data.findStudio);
      setImagePreview(data.findStudio.image_path);
      setStudio(data.findStudio);
    }
  }, [data]);

  function onImageLoad(this: FileReader) {
    setImagePreview(this.result as string);
    setImage(this.result as string);
  }

  ImageUtils.usePasteImage(onImageLoad);

  if (!isNew && !isEditing) {
    if (!data?.findStudio || loading)
      return <Spinner animation="border" variant="light" />;
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
        updateStudioData(result.data.studioUpdate);
        setIsEditing(false);
      } else {
        const result = await createStudio();
        history.push(`/studios/${result.data.studioCreate.id}`);
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  async function onAutoTag() {
    if (!studio || !studio.id) {
      return;
    }
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

  // TODO: CSS class
  return (
    <div className="columns is-multiline no-spacing">
      <div className="column is-half details-image-container">
        <img className="studio" alt="" src={imagePreview} />
      </div>
      <div className="column is-half details-detail-container">
        <DetailsEditNavbar
          studio={studio}
          isNew={isNew}
          isEditing={isEditing}
          onToggleEdit={() => {
            setIsEditing(!isEditing);
            updateStudioEditState(studio);
          }}
          onSave={onSave}
          onDelete={onDelete}
          onAutoTag={onAutoTag}
          onImageChange={onImageChangeHandler}
        />
        <h1>
          {!isEditing ? (
            <span>{studio.name}</span>
          ) : (
            <Form.Group controlId="studio-name">
              <Form.Label>Name</Form.Label>
              <Form.Control
                defaultValue={studio.name || ""}
                placeholder="Name"
                onChange={(event: any) => setName(event.target.value)}
              />
            </Form.Group>
          )}
        </h1>

        <Table style={{ width: "100%" }}>
          <tbody>
            {TableUtils.renderInputGroup({
              title: "URL",
              value: studio.url,
              isEditing,
              onChange: (val: string) => setUrl(val)
            })}
          </tbody>
        </Table>
      </div>
    </div>
  );
};
