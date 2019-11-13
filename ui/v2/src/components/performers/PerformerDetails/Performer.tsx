import {
  Button,
  Classes,
  Dialog,
  EditableText,
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
import { ScrapePerformerSuggest } from "../../select/ScrapePerformerSuggest";
import { DetailsEditNavbar } from "../../Shared/DetailsEditNavbar";

interface IPerformerProps extends IBaseProps {}

export const Performer: FunctionComponent<IPerformerProps> = (props: IPerformerProps) => {
  const isNew = props.match.params.id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);
  const [isDisplayingScraperDialog, setIsDisplayingScraperDialog] = useState<GQL.ListScrapersListScrapers | undefined>(undefined);
  const [scrapePerformerDetails, setScrapePerformerDetails] = useState<GQL.ScrapePerformerListScrapePerformerList | undefined>(undefined);

  // Editing performer state
  const [image, setImage] = useState<string | undefined>(undefined);
  const [name, setName] = useState<string | undefined>(undefined);
  const [aliases, setAliases] = useState<string | undefined>(undefined);
  const [favorite, setFavorite] = useState<boolean | undefined>(undefined);
  const [birthdate, setBirthdate] = useState<string | undefined>(undefined);
  const [ethnicity, setEthnicity] = useState<string | undefined>(undefined);
  const [country, setCountry] = useState<string | undefined>(undefined);
  const [eyeColor, setEyeColor] = useState<string | undefined>(undefined);
  const [height, setHeight] = useState<string | undefined>(undefined);
  const [measurements, setMeasurements] = useState<string | undefined>(undefined);
  const [fakeTits, setFakeTits] = useState<string | undefined>(undefined);
  const [careerLength, setCareerLength] = useState<string | undefined>(undefined);
  const [tattoos, setTattoos] = useState<string | undefined>(undefined);
  const [piercings, setPiercings] = useState<string | undefined>(undefined);
  const [url, setUrl] = useState<string | undefined>(undefined);
  const [twitter, setTwitter] = useState<string | undefined>(undefined);
  const [instagram, setInstagram] = useState<string | undefined>(undefined);

  // Performer state
  const [performer, setPerformer] = useState<Partial<GQL.PerformerDataFragment>>({});
  const [imagePreview, setImagePreview] = useState<string | undefined>(undefined);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const Scrapers = StashService.useListScrapers(GQL.ScraperType.Performer);

  const { data, error, loading } = StashService.useFindPerformer(props.match.params.id);
  const updatePerformer = StashService.usePerformerUpdate(getPerformerInput() as GQL.PerformerUpdateInput);
  const createPerformer = StashService.usePerformerCreate(getPerformerInput() as GQL.PerformerCreateInput);
  const deletePerformer = StashService.usePerformerDestroy(getPerformerInput() as GQL.PerformerDestroyInput);

  function updatePerformerEditState(state: Partial<GQL.PerformerDataFragment | GQL.ScrapeFreeonesScrapeFreeones>) {
    if ((state as GQL.PerformerDataFragment).favorite !== undefined) {
      setFavorite((state as GQL.PerformerDataFragment).favorite);
    }
    setName(state.name);
    setAliases(state.aliases);
    setBirthdate(state.birthdate);
    setEthnicity(state.ethnicity);
    setCountry(state.country);
    setEyeColor(state.eye_color);
    setHeight(state.height);
    setMeasurements(state.measurements);
    setFakeTits(state.fake_tits);
    setCareerLength(state.career_length);
    setTattoos(state.tattoos);
    setPiercings(state.piercings);
    setUrl(state.url);
    setTwitter(state.twitter);
    setInstagram(state.instagram);
  }

  useEffect(() => {
    setIsLoading(loading);
    if (!data || !data.findPerformer || !!error) { return; }
    setPerformer(data.findPerformer);
  }, [data]);

  useEffect(() => {
    setImagePreview(performer.image_path);
    setImage(undefined);
    updatePerformerEditState(performer);
    if (!isNew) {
      setIsEditing(false);
    }
  }, [performer]);

  if ((!isNew && !isEditing && (!data || !data.findPerformer)) || isLoading) {
    return <Spinner size={Spinner.SIZE_LARGE} />; 
  }
  if (!!error) { return <>error...</>; }

  function getPerformerInput() {
    const performerInput: Partial<GQL.PerformerCreateInput | GQL.PerformerUpdateInput> = {
      name,
      aliases,
      favorite,
      birthdate,
      ethnicity,
      country,
      eye_color: eyeColor,
      height,
      measurements,
      fake_tits: fakeTits,
      career_length: careerLength,
      tattoos,
      piercings,
      url,
      twitter,
      instagram,
      image,
    };

    if (!isNew) {
      (performerInput as GQL.PerformerUpdateInput).id = props.match.params.id;
    }
    return performerInput;
  }

  async function onSave() {
    setIsLoading(true);
    try {
      if (!isNew) {
        const result = await updatePerformer();
        setPerformer(result.data.performerUpdate);
      } else {
        const result = await createPerformer();
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
      const result = await deletePerformer();
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
    
    // redirect to performers page
    props.history.push(`/performers`);
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

  function onDisplayFreeOnesDialog(scraper: GQL.ListScrapersListScrapers) {
    setIsDisplayingScraperDialog(scraper);
  }

  function getQueryScraperPerformerInput() {
    if (!scrapePerformerDetails) {
      return {};
    }

    let ret = _.clone(scrapePerformerDetails);
    delete ret.__typename;
    return ret as GQL.ScrapedPerformerInput;
  }

  async function onScrapePerformer() {
    setIsDisplayingScraperDialog(undefined);
    setIsLoading(true);
    try {
      if (!scrapePerformerDetails || !isDisplayingScraperDialog) { return; }
      const result = await StashService.queryScrapePerformer(isDisplayingScraperDialog.id, getQueryScraperPerformerInput());
      if (!result.data || !result.data.scrapePerformer) { return; }
      updatePerformerEditState(result.data.scrapePerformer);
    } catch (e) {
      ErrorUtils.handle(e);
    }
    setIsLoading(false);
  }

  function renderEthnicity() {
    return TableUtils.renderHtmlSelect({
      title: "Ethnicity",
      value: ethnicity,
      isEditing,
      onChange: (value: string) => setEthnicity(value),
      selectOptions: ["white", "black", "asian", "hispanic"],
    });
  }

  function renderScraperDialog() {
    return (
      <Dialog
        isOpen={!!isDisplayingScraperDialog}
        onClose={() => setIsDisplayingScraperDialog(undefined)}
        title="Scrape"
      >
        <div className="dialog-content">
          <ScrapePerformerSuggest
            placeholder="Performer name"
            style={{width: "100%"}}
            scraperId={isDisplayingScraperDialog ? isDisplayingScraperDialog.id : ""}
            onSelectPerformer={(query) => setScrapePerformerDetails(query)}
          />
        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <Button onClick={() => onScrapePerformer()}>Scrape</Button>
          </div>
        </div>
      </Dialog>
    );
  }

  return (
    <>
      {renderScraperDialog()}
      <div className="columns is-multiline no-spacing">
        <div className="column is-half details-image-container">
          <img className="performer" src={imagePreview} />
        </div>
        <div className="column is-half details-detail-container">
          <DetailsEditNavbar
            performer={performer}
            isNew={isNew}
            isEditing={isEditing}
            onToggleEdit={() => { setIsEditing(!isEditing); updatePerformerEditState(performer); }}
            onSave={onSave}
            onDelete={onDelete}
            onImageChange={onImageChange}
            scrapers={Scrapers.data ? Scrapers.data.listScrapers : undefined}
            onDisplayScraperDialog={onDisplayFreeOnesDialog}
          />
          <h1 className="bp3-heading">
            <EditableText
              disabled={!isEditing}
              value={name}
              placeholder="Name"
              onChange={(value) => setName(value)}
            />
          </h1>
          <h6 className="bp3-heading">
            <span style={{fontWeight: 300}}>Aliases: </span>
            <EditableText
              disabled={!isEditing}
              value={aliases}
              placeholder="Aliases"
              onChange={(value) => setAliases(value)}
            />
          </h6>
          <div>
            <span style={{fontWeight: 300}}>Favorite:</span>
            <Button
              icon="heart"
              disabled={!isEditing}
              className={favorite ? "favorite" : undefined}
              onClick={() => setFavorite(!favorite)}
              minimal={true}
            />
          </div>

          <HTMLTable style={{width: "100%"}}>
            <tbody>
              {TableUtils.renderEditableTextTableRow(
                {title: "Birthdate (YYYY-MM-DD)", value: birthdate, isEditing, onChange: setBirthdate})}
              {renderEthnicity()}
              {TableUtils.renderEditableTextTableRow(
                {title: "Eye Color", value: eyeColor, isEditing, onChange: setEyeColor})}
              {TableUtils.renderEditableTextTableRow(
                {title: "Country", value: country, isEditing, onChange: setCountry})}
              {TableUtils.renderEditableTextTableRow(
                {title: "Height (CM)", value: height, isEditing, onChange: setHeight})}
              {TableUtils.renderEditableTextTableRow(
                {title: "Measurements", value: measurements, isEditing, onChange: setMeasurements})}
              {TableUtils.renderEditableTextTableRow(
                {title: "Fake Tits", value: fakeTits, isEditing, onChange: setFakeTits})}
              {TableUtils.renderEditableTextTableRow(
                {title: "Career Length", value: careerLength, isEditing, onChange: setCareerLength})}
              {TableUtils.renderEditableTextTableRow(
                {title: "Tattoos", value: tattoos, isEditing, onChange: setTattoos})}
              {TableUtils.renderEditableTextTableRow(
                {title: "Piercings", value: piercings, isEditing, onChange: setPiercings})}
              {TableUtils.renderEditableTextTableRow(
                {title: "URL", value: url, isEditing, onChange: setUrl})}
              {TableUtils.renderEditableTextTableRow(
                {title: "Twitter", value: twitter, isEditing, onChange: setTwitter})}
              {TableUtils.renderEditableTextTableRow(
                {title: "Instagram", value: instagram, isEditing, onChange: setInstagram})}
            </tbody>
          </HTMLTable>
        </div>
      </div>
    </>
  );
};
