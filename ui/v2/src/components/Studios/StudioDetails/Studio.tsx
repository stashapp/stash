import {
  Button,
  Classes,
  Dialog,
  EditableText,
  HTMLSelect,
  HTMLTable,
  Spinner,
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { IBaseProps } from "../../../models";
import { ErrorUtils } from "../../../utils/errors";
import { TableUtils } from "../../../utils/table";
import { DetailsEditNavbar } from "../../Shared/DetailsEditNavbar";
import { ToastUtils } from "../../../utils/toasts";

interface IProps extends IBaseProps {}

export const Studio: FunctionComponent<IProps> = (props: IProps) => {
  const isNew = props.match.params.id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);

  // Editing studio state
  const [image, setImage] = useState<string | undefined>(undefined);
  const [name, setName] = useState<string | undefined>(undefined);
  const [url, setUrl] = useState<string | undefined>(undefined);

  // Studio state
  const [studio, setStudio] = useState<Partial<GQL.StudioDataFragment>>({});
  const [imagePreview, setImagePreview] = useState<string | undefined>(undefined);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const { data, error, loading } = StashService.useFindStudio(props.match.params.id);
  const updateStudio = StashService.useStudioUpdate(getStudioInput() as GQL.StudioUpdateInput);
  const createStudio = StashService.useStudioCreate(getStudioInput() as GQL.StudioCreateInput);
  const deleteStudio = StashService.useStudioDestroy(getStudioInput() as GQL.StudioDestroyInput);

  function updateStudioEditState(state: Partial<GQL.StudioDataFragment>) {
    setName(state.name);
    setUrl(state.url);
  }

  useEffect(() => {
    setIsLoading(loading);
    if (!data || !data.findStudio || !!error) { return; }
    setStudio(data.findStudio);
  }, [data]);

  useEffect(() => {
    setImagePreview(studio.image_path);
    setImage(undefined);
    updateStudioEditState(studio);
    if (!isNew) {
      setIsEditing(false);
    }
  }, [studio]);

  if (!isNew && !isEditing) {
    if (!data || !data.findStudio || isLoading) { return <Spinner size={Spinner.SIZE_LARGE} />; }
    if (!!error) { return <>error...</>; }
  }

  function getStudioInput() {
    const input: Partial<GQL.StudioCreateInput | GQL.StudioUpdateInput> = {
      name,
      url,
      image,
    };

    if (!isNew) {
      (input as GQL.StudioUpdateInput).id = props.match.params.id;
    }
    return input;
  }

  async function onSave() {
    setIsLoading(true);
    try {
      if (!isNew) {
        const result = await updateStudio();
        setStudio(result.data.studioUpdate);
      } else {
        const result = await createStudio();
        setStudio(result.data.studioCreate);
        props.history.push(`/studios/${result.data.studioCreate.id}`);
      }
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
  }

  async function onAutoTag() {
    if (!studio || !studio.id) {
      return;
    }
    try {
      await StashService.queryMetadataAutoTag({ studios: [studio.id]});
      ToastUtils.success("Started auto tagging");
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  async function onDelete() {
    setIsLoading(true);
    try {
      const result = await deleteStudio();
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
    
    // redirect to studios page
    props.history.push(`/studios`);
  }

  function onImageChange(event: React.FormEvent<HTMLInputElement>) {
    const file: File = (event.target as any).files[0];
    const reader: FileReader = new FileReader();

    reader.onloadend = (e) => {
      setImagePreview(reader.result as string);
      setImage(reader.result as string);
    };
    reader.readAsDataURL(file);
  }

  // TODO: CSS class
  return (
    <>
      <div className="columns is-multiline no-spacing">
        <div className="column is-half details-image-container">
          <img className="studio" src={imagePreview} />
        </div>
        <div className="column is-half details-detail-container">
          <DetailsEditNavbar
            studio={studio}
            isNew={isNew}
            isEditing={isEditing}
            onToggleEdit={() => { setIsEditing(!isEditing); updateStudioEditState(studio); }}
            onSave={onSave}
            onDelete={onDelete}
            onAutoTag={onAutoTag}
            onImageChange={onImageChange}
          />
          <h1 className="bp3-heading">
            <EditableText
              disabled={!isEditing}
              value={name}
              placeholder="Name"
              onChange={(value) => setName(value)}
            />
          </h1>

          <HTMLTable style={{width: "100%"}}>
            <tbody>
              {TableUtils.renderEditableTextTableRow({title: "URL", value: url, isEditing, onChange: setUrl})}
            </tbody>
          </HTMLTable>
        </div>
      </div>
    </>
  );
};
