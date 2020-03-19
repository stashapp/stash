import {
  Button,
  Classes,
  Dialog,
  HTMLTable,
  Spinner,
  Menu,
  MenuItem,
  Popover,
  Alert,
  FileInput,
} from "@blueprintjs/core";
import _ from "lodash";
import React, { FunctionComponent, useEffect, useState } from "react";
import * as GQL from "../../../core/generated-graphql";
import { StashService } from "../../../core/StashService";
import { ErrorUtils } from "../../../utils/errors";
import { TableUtils } from "../../../utils/table";
import { ScrapePerformerSuggest } from "../../select/ScrapePerformerSuggest";
import { EditableTextUtils } from "../../../utils/editabletext";
import { ImageUtils } from "../../../utils/image";

interface IPerformerDetailsProps {
  performer: Partial<GQL.PerformerDataFragment>
  isNew?: boolean
  isEditing?: boolean
  onSave? : (performer : Partial<GQL.PerformerCreateInput> | Partial<GQL.PerformerUpdateInput>) => void
  onDelete? : () => void
  onImageChange? : (image: string) => void
}

export const PerformerDetailsPanel: FunctionComponent<IPerformerDetailsProps> = (props: IPerformerDetailsProps) => {

  // Editing state
  const [isDisplayingScraperDialog, setIsDisplayingScraperDialog] = useState<GQL.ListPerformerScrapersListPerformerScrapers | undefined>(undefined);
  const [scrapePerformerDetails, setScrapePerformerDetails] = useState<GQL.ScrapePerformerListScrapePerformerList | undefined>(undefined);
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

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
  const [gender, setGender] = useState<string | undefined>('female');

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const Scrapers = StashService.useListPerformerScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.ListPerformerScrapersListPerformerScrapers[]>([]);

  function updatePerformerEditState(state: Partial<GQL.PerformerDataFragment | GQL.ScrapedPerformerDataFragment | GQL.ScrapeFreeonesScrapeFreeones>) {
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
    if ((state as GQL.PerformerDataFragment).favorite !== undefined) {
      setGender((state as GQL.PerformerDataFragment).gender)
    }
  }

  function updatePerformerEditStateFromScraper(state: Partial<GQL.ScrapedPerformerDataFragment | GQL.ScrapeFreeonesScrapeFreeones>) {
    updatePerformerEditState(state);

    // image is a base64 string
    if ((state as GQL.ScrapedPerformerDataFragment).image !== undefined) {
      let imageStr = (state as GQL.ScrapedPerformerDataFragment).image;
      setImage(imageStr);
      if (props.onImageChange) {
        props.onImageChange(imageStr!);
      }
    }
  }

  useEffect(() => {
    setImage(undefined);
    updatePerformerEditState(props.performer);
  }, [props.performer]);

  function onImageLoad(this: FileReader) {
    setImage(this.result as string);
    if (props.onImageChange) {
      props.onImageChange(this.result as string);
    }
  }

  if (props.isEditing) {
    ImageUtils.addPasteImageHook(onImageLoad);
  }
  
  useEffect(() => {
    var newQueryableScrapers : GQL.ListPerformerScrapersListPerformerScrapers[] = [];

    if (!!Scrapers.data && Scrapers.data.listPerformerScrapers) {
      newQueryableScrapers = Scrapers.data.listPerformerScrapers.filter((s) => {
        return s.performer && s.performer.supported_scrapes.includes(GQL.ScrapeType.Name);
      });
    }

    setQueryableScrapers(newQueryableScrapers);

  }, [Scrapers.data]);

  if (isLoading) {
    return <Spinner size={Spinner.SIZE_LARGE} />; 
  }

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
      gender
    };

    if (!props.isNew) {
      (performerInput as GQL.PerformerUpdateInput).id = props.performer.id!;
    }
    return performerInput;
  }

  function onSave() {
    if (props.onSave) {
      props.onSave(getPerformerInput());
    }
  }

  function onDelete() {
    if (props.onDelete) {
      props.onDelete();
    }
  }

  function onImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function onDisplayFreeOnesDialog(scraper: GQL.ListPerformerScrapersListPerformerScrapers) {
    setIsDisplayingScraperDialog(scraper);
  }

  function getQueryScraperPerformerInput() {
    if (!scrapePerformerDetails) {
      return {};
    }

    let ret = _.clone(scrapePerformerDetails);
    delete ret.__typename;

    // image is not supported
    delete ret.image;
    
    return ret as GQL.ScrapedPerformerInput;
  }

  async function onScrapePerformer() {
    setIsDisplayingScraperDialog(undefined);
    try {
      if (!scrapePerformerDetails || !isDisplayingScraperDialog) { return; }
      setIsLoading(true);
      const result = await StashService.queryScrapePerformer(isDisplayingScraperDialog.id, getQueryScraperPerformerInput());
      if (!result.data || !result.data.scrapePerformer) { return; }
      updatePerformerEditStateFromScraper(result.data.scrapePerformer);
    } catch (e) {
      ErrorUtils.handle(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function onScrapePerformerURL() {
    if (!url) { return; }
    setIsLoading(true);
    try {
      const result = await StashService.queryScrapePerformerURL(url);
      if (!result.data || !result.data.scrapePerformerURL) { return; }

      // leave URL as is if not set explicitly
      if (!result.data.scrapePerformerURL.url) {
        result.data.scrapePerformerURL.url = url;
      }

      updatePerformerEditStateFromScraper(result.data.scrapePerformerURL);
    } catch (e) {
      ErrorUtils.handle(e);
    } finally {
      setIsLoading(false);
    }
  }

  function renderEthnicity() {
    return TableUtils.renderHtmlSelect({
      title: "Ethnicity",
      value: ethnicity,
      isEditing: !!props.isEditing,
      onChange: (value: string) => setEthnicity(value),
      selectOptions: ["white", "black", "asian", "hispanic"],
    });
  }

  function renderScraperMenu() {
    function renderScraperMenuItem(scraper : GQL.ListPerformerScrapersListPerformerScrapers) {
      return (
        <MenuItem
          text={scraper.name}
          onClick={() => { onDisplayFreeOnesDialog(scraper); }}
        />
      );
    }
    
    if (!props.performer) { return; }
    if (!props.isEditing) { return; }
    const scraperMenu = (
      <Menu>
        {queryableScrapers ? queryableScrapers.map((s) => renderScraperMenuItem(s)) : undefined}
      </Menu>
    );
    return (
      <Popover content={scraperMenu} position="bottom">
        <Button text="Scrape with..."/>
      </Popover>
    );
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

  function urlScrapable(url: string) : boolean {
    return !!url && !!Scrapers.data && Scrapers.data.listPerformerScrapers && Scrapers.data.listPerformerScrapers.some((s) => {
      return !!s.performer && !!s.performer.urls && s.performer.urls.some((u) => { return url.includes(u); });
    });
  }

  function maybeRenderScrapeButton() {
    if (!url || !props.isEditing || !urlScrapable(url)) {
      return undefined;
    }
    return (
      <Button 
        minimal={true} 
        icon="import" 
        id="scrape-url-button"
        onClick={() => onScrapePerformerURL()}/>
    )
  }

  function renderURLField() {
    return (
      <tr>
        <td id="url-field">
          URL 
          {maybeRenderScrapeButton()}
        </td>
        <td>
          {EditableTextUtils.renderInputGroup({
            value: url, asURL: true, isEditing: !!props.isEditing, onChange: setUrl, placeholder: "URL"
          })}
        </td>
      </tr>
    );
  }

  function renderImageInput() {
    if (!props.isEditing) { return; }
    return (
      <>
        <tr>
          <td>Image</td>
          <td><FileInput text="Choose image..." onInputChange={onImageChange} inputProps={{accept: ".jpg,.jpeg"}} /></td>
        </tr>
      </>
    )
  }

  function maybeRenderButtons() {
    if (props.isEditing) {
      return (
        <>
          <Button className="edit-button" text="Save" intent="primary" onClick={() => onSave()}/>
          {!props.isNew ? <Button className="edit-button" text="Delete" intent="danger" onClick={() => setIsDeleteAlertOpen(true)}/> : undefined}
          {renderScraperMenu()}
        </>
      );
    }
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
          Are you sure you want to delete {name}?
        </p>
      </Alert>
    );
  }

  function maybeRenderName() {
    if (props.isEditing) {
      return TableUtils.renderInputGroup(
        {title: "Name", value: name, isEditing: !!props.isEditing, placeholder: "Name", onChange: setName});
    }
  }

  function maybeRenderAliases() {
    if (props.isEditing) {
      return TableUtils.renderInputGroup(
        {title: "Aliases", value: aliases, isEditing: !!props.isEditing, placeholder: "Aliases", onChange: setAliases});
    }
  }

  function renderGender() {
    return TableUtils.renderHtmlSelect({
      title: "Gender",
      value: gender,
      isEditing: !!props.isEditing,
      onChange: (value: string) => setGender(value),
      selectOptions: ["male", "female"],
    });
  }

  const twitterPrefix = "https://twitter.com/";
  const instagramPrefix = "https://www.instagram.com/";

  return (
    <>
      {renderDeleteAlert()}
      {renderScraperDialog()}
      
      <HTMLTable id="performer-details" style={{width: "100%"}}>
        <tbody>
          {maybeRenderName()}
          {maybeRenderAliases()}
          {renderGender()}
          {TableUtils.renderInputGroup(
            {title: "Birthdate (YYYY-MM-DD)", value: birthdate, isEditing: !!props.isEditing, onChange: setBirthdate})}
          {renderEthnicity()}
          {TableUtils.renderInputGroup(
            {title: "Eye Color", value: eyeColor, isEditing: !!props.isEditing, onChange: setEyeColor})}
          {TableUtils.renderInputGroup(
            {title: "Country", value: country, isEditing: !!props.isEditing, onChange: setCountry})}
          {TableUtils.renderInputGroup(
            {title: "Height (CM)", value: height, isEditing: !!props.isEditing, onChange: setHeight})}
          {TableUtils.renderInputGroup(
            {title: "Measurements", value: measurements, isEditing: !!props.isEditing, onChange: setMeasurements})}
          {TableUtils.renderInputGroup(
            {title: "Fake Tits", value: fakeTits, isEditing: !!props.isEditing, onChange: setFakeTits})}
          {TableUtils.renderInputGroup(
            {title: "Career Length", value: careerLength, isEditing: !!props.isEditing, onChange: setCareerLength})}
          {TableUtils.renderInputGroup(
            {title: "Tattoos", value: tattoos, isEditing: !!props.isEditing, onChange: setTattoos})}
          {TableUtils.renderInputGroup(
            {title: "Piercings", value: piercings, isEditing: !!props.isEditing, onChange: setPiercings})}
          {renderURLField()}
          {TableUtils.renderInputGroup(
            {title: "Twitter", value: twitter, asURL: true, urlPrefix: twitterPrefix, isEditing: !!props.isEditing, onChange: setTwitter})}
          {TableUtils.renderInputGroup(
            {title: "Instagram", value: instagram, asURL: true, urlPrefix: instagramPrefix, isEditing: !!props.isEditing, onChange: setInstagram})}
          {renderImageInput()}
        </tbody>
      </HTMLTable>

      {maybeRenderButtons()}
    </>
  );
};
