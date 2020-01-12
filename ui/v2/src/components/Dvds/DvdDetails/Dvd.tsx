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
import { ImageUtils } from "../../../utils/image";

interface IProps extends IBaseProps {}

export const Dvd: FunctionComponent<IProps> = (props: IProps) => {
  const isNew = props.match.params.id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);

  // Editing dvd state
  const [frontimage, setFrontImage] = useState<string | undefined>(undefined);
  const [backimage, setBackImage] = useState<string | undefined>(undefined);
  const [name, setName] = useState<string | undefined>(undefined);
  const [aliases, setAliases] = useState<string | undefined>(undefined);
  const [durationdvd, setDurationdvd] = useState<string | undefined>(undefined);
  const [year, setYear] = useState<string | undefined>(undefined);
  const [director, setDirector] = useState<string | undefined>(undefined);
  const [synopsis, setSynopsis] = useState<string | undefined>(undefined);
  const [url, setUrl] = useState<string | undefined>(undefined);

  // Dvd state
  const [dvd, setDvd] = useState<Partial<GQL.DvdDataFragment>>({});
  const [imagePreview, setImagePreview] = useState<string | undefined>(undefined);
  const [backimagePreview, setBackImagePreview] = useState<string | undefined>(undefined);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const { data, error, loading } = StashService.useFindDvd(props.match.params.id);
  const updateDvd = StashService.useDvdUpdate(getDvdInput() as GQL.DvdUpdateInput);
  const createDvd = StashService.useDvdCreate(getDvdInput() as GQL.DvdCreateInput);
  const deleteDvd = StashService.useDvdDestroy(getDvdInput() as GQL.DvdDestroyInput);

  function updateDvdEditState(state: Partial<GQL.DvdDataFragment>) {
    setName(state.name);
    setAliases(state.aliases);
    setDurationdvd(state.durationdvd);
    setYear(state.year);
    setDirector(state.director);
    setSynopsis(state.synopsis);
    setUrl(state.url);
  }

  useEffect(() => {
    setIsLoading(loading);
    if (!data || !data.findDvd || !!error) { return; }
    setDvd(data.findDvd);
  }, [data]);

  useEffect(() => {
    setImagePreview(dvd.frontimage_path);
    setBackImagePreview(dvd.backimage_path);
    setFrontImage(undefined);
    setBackImage(undefined);
    updateDvdEditState(dvd);
    if (!isNew) {
      setIsEditing(false);
    }
  }, [dvd]);

  function onImageLoad(this: FileReader) {
    setImagePreview(this.result as string);
    setFrontImage(this.result as string);
    
  }

  function onBackImageLoad(this: FileReader) {
    setBackImagePreview(this.result as string);
    setBackImage(this.result as string);
  }


  ImageUtils.addPasteImageHook(onImageLoad);
  ImageUtils.addPasteImageHook(onBackImageLoad);

  if (!isNew && !isEditing) {
    if (!data || !data.findDvd || isLoading) { return <Spinner size={Spinner.SIZE_LARGE} />; }
    if (!!error) { return <>error...</>; }
  }

  function getDvdInput() {
    const input: Partial<GQL.DvdCreateInput | GQL.DvdUpdateInput> = {
      name,
      aliases,
      durationdvd,
      year,
      director,
      synopsis,
      url,
      frontimage,
      backimage
    };

    if (!isNew) {
      (input as GQL.DvdUpdateInput).id = props.match.params.id;
    }
    return input;
  }

  async function onSave() {
    setIsLoading(true);
    try {
      if (!isNew) {
        const result = await updateDvd();
        setDvd(result.data.dvdUpdate);
      } else {
        const result = await createDvd();
        setDvd(result.data.dvdCreate);
        props.history.push(`/dvds/${result.data.dvdCreate.id}`);
      }
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
  }

  async function onAutoTag() {
    if (!dvd || !dvd.id) {
      return;
    }
    try {
      await StashService.queryMetadataAutoTag({ dvds: [dvd.id]});
      ToastUtils.success("Started auto tagging");
    } catch (e) {
      ErrorUtils.handle(e);
    }
  }

  async function onDelete() {
    setIsLoading(true);
    try {
      const result = await deleteDvd();
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
    
    // redirect to dvds page
    props.history.push(`/dvds`);
  }

  function onImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function onBackImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onBackImageLoad);
  }

  // TODO: CSS class
  return (
    <>
      <div className="columns is-multiline no-spacing">
        <div className="column is-half details-image-container">
          <img className="dvd" src={imagePreview} />
          <img className="dvd" src={backimagePreview} />
       </div>
        <div className="column is-half details-detail-container">
          <DetailsEditNavbar
            dvd={dvd}
            isNew={isNew}
            isEditing={isEditing}
            onToggleEdit={() => { setIsEditing(!isEditing); updateDvdEditState(dvd); }}
            onSave={onSave}
            onDelete={onDelete}
            onAutoTag={onAutoTag}
            onImageChange={onImageChange}
            onBackImageChange={onBackImageChange}
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
              {TableUtils.renderInputGroup({title: "Aliases", value: aliases, isEditing, onChange: setAliases})}
              {TableUtils.renderInputGroup({title: "Duration", value: durationdvd, isEditing, onChange: setDurationdvd})}
              {TableUtils.renderInputGroup({title: "Year", value: year, isEditing, onChange: setYear})}
              {TableUtils.renderInputGroup({title: "Director", value: director, isEditing, onChange: setDirector})}
              {TableUtils.renderInputGroup({title: "URL", value: url, isEditing, onChange: setUrl})}
              {TableUtils.renderTextArea({title: "Synopsis", value: synopsis, isEditing, onChange: setSynopsis})}            
            
            </tbody>
          </HTMLTable>
        </div>
      </div>
    </>
  );
};
