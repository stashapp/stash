import React, { useEffect, useState } from "react";
import { Button, Form, Spinner, Table } from 'react-bootstrap';
import { useParams, useHistory } from 'react-router-dom';
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import { DetailsEditNavbar, Icon, Modal, ScrapePerformerSuggest } from "src/components/Shared";
import { ImageUtils, TableUtils } from 'src/utils'
import { useToast } from 'src/hooks';

export const Performer: React.FC = () => {
  const Toast = useToast();
  const history = useHistory();
  const { id = 'new' } = useParams();
  const isNew = id === "new";

  // Editing state
  const [isEditing, setIsEditing] = useState<boolean>(isNew);
  const [isDisplayingScraperDialog, setIsDisplayingScraperDialog] = useState<GQL.ListPerformerScrapersListPerformerScrapers | undefined>(undefined);
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

  const Scrapers = StashService.useListPerformerScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.ListPerformerScrapersListPerformerScrapers[]>([]);

  const { data, error } = StashService.useFindPerformer(id);
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
    setIsLoading(false);
    if(data?.findPerformer)
      setPerformer(data.findPerformer);
  }, [data]);

  useEffect(() => {
    setImagePreview(performer.image_path);
    setImage(undefined);
    updatePerformerEditState(performer);
    setIsEditing(false);
  }, [performer]);

  function onImageLoad(this: FileReader) {
    setImagePreview(this.result as string);
    setImage(this.result as string);
  }

  ImageUtils.usePasteImage(onImageLoad);

  useEffect(() => {
    var newQueryableScrapers : GQL.ListPerformerScrapersListPerformerScrapers[] = [];

    if (!!Scrapers.data && Scrapers.data.listPerformerScrapers) {
      newQueryableScrapers = Scrapers.data.listPerformerScrapers.filter((s) => {
        return s.performer && s.performer.supported_scrapes.includes(GQL.ScrapeType.Name);
      });
    }

    setQueryableScrapers(newQueryableScrapers);

  }, [Scrapers.data]);

  if ((!isNew && !isEditing && !data?.findPerformer) || isLoading)
    return <Spinner animation="border" variant="light" />;
  if (error)
    return <div>{error.message}</div>;

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
      (performerInput as GQL.PerformerUpdateInput).id = id;
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
        history.push(`/performers/${result.data.performerCreate.id}`);
      }
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  async function onDelete() {
    setIsLoading(true);
    try {
      await deletePerformer();
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);

    // redirect to performers page
    history.push('/performers');
  }

  async function onAutoTag() {
    if (!performer || !performer.id) {
      return;
    }
    try {
      await StashService.queryMetadataAutoTag({ performers: [performer.id]});
      Toast.success({ content: "Started auto tagging" });
    } catch (e) {
      Toast.error(e);
    }
  }

  function onImageChangeHandler(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function onDisplayFreeOnesDialog(scraper: GQL.ListPerformerScrapersListPerformerScrapers) {
    setIsDisplayingScraperDialog(scraper);
  }

  function getQueryScraperPerformerInput() {
    if (!scrapePerformerDetails)
      return {};

    const { __typename, ...ret } = scrapePerformerDetails;
    return ret;
  }

  async function onScrapePerformer() {
    setIsDisplayingScraperDialog(undefined);
    setIsLoading(true);
    try {
      if (!scrapePerformerDetails || !isDisplayingScraperDialog)
        return;
      const result = await StashService.queryScrapePerformer(isDisplayingScraperDialog.id, getQueryScraperPerformerInput());
      if (!result?.data?.scrapePerformer)
        return;
      updatePerformerEditState(result.data.scrapePerformer);
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  async function onScrapePerformerURL() {
    if (!url)
      return;
    setIsLoading(true);
    try {
      const result = await StashService.queryScrapePerformerURL(url);
      if (!result.data || !result.data.scrapePerformerURL) { return; }
      updatePerformerEditState(result.data.scrapePerformerURL);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
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
      <Modal
        show={!!isDisplayingScraperDialog}
        onHide={() => setIsDisplayingScraperDialog(undefined)}
        header="Scrape"
        accept={{ onClick: onScrapePerformer, text: "Scrape" }}
      >
        <div className="dialog-content">
          <ScrapePerformerSuggest
            placeholder="Performer name"
            scraperId={isDisplayingScraperDialog ? isDisplayingScraperDialog.id : ""}
            onSelectPerformer={(query) => setScrapePerformerDetails(query)}
          />
        </div>
      </Modal>
    );
  }

  function urlScrapable(url: string) {
    return !!url && (Scrapers?.data?.listPerformerScrapers ?? []).some(s => (
      (s?.performer?.urls ?? []).some(u => url.includes(u))
    ));
  }

  function maybeRenderScrapeButton() {
    if (!url || !isEditing || !urlScrapable(url)) {
      return undefined;
    }
    return (
      <Button
        id="scrape-url-button"
        onClick={() => onScrapePerformerURL()}>
        <Icon icon="file-upload" />
      </Button>
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
          <Form.Control
            value={url}
            readOnly={!isEditing}
            plaintext={!isEditing}
            placeholder="URL"
            onChange={(event: React.FormEvent<HTMLInputElement>) => setUrl(event.currentTarget.value) }
          />
        </td>
      </tr>
    );
  }

  return (
    <>
      {renderScraperDialog()}
      <div className="row is-multiline no-spacing">
        <div className="col-6 details-image-container">
          <img className="performer" alt="" src={imagePreview} />
        </div>
        <div className="col-6 details-detail-container">
          <DetailsEditNavbar
            performer={performer}
            isNew={isNew}
            isEditing={isEditing}
            onToggleEdit={() => { setIsEditing(!isEditing); updatePerformerEditState(performer); }}
            onSave={onSave}
            onDelete={onDelete}
            onImageChange={onImageChangeHandler}
            scrapers={queryableScrapers}
            onDisplayScraperDialog={onDisplayFreeOnesDialog}
            onAutoTag={onAutoTag}
          />
          <h1>
            { <Form.Control
                readOnly={!isEditing}
                plaintext={!isEditing}
                defaultValue={name}
                placeholder="Name"
                onChange={(event: any) => setName(event.target.value)}
              />
            }
          </h1>
          <h6>
            <Form.Group className="aliases-field" controlId="aliases">
              <Form.Label>Aliases:</Form.Label>
              <Form.Control
                value={aliases}
                readOnly={!isEditing}
                plaintext={!isEditing}
                placeholder="Aliases"
                onChange={(event: React.FormEvent<HTMLInputElement>) => setAliases(event.currentTarget.value) }
              />
            </Form.Group>
          </h6>
          <div>
            <span style={{fontWeight: 300}}>Favorite:</span>
            <Button
              disabled={!isEditing}
              className={favorite ? "favorite" : undefined}
              onClick={() => setFavorite(!favorite)}
            >
              <Icon icon="heart" />
            </Button>
          </div>

          <Table id="performer-details" style={{width: "100%"}}>
            <tbody>
              {TableUtils.renderInputGroup(
                {title: "Birthdate (YYYY-MM-DD)", value: birthdate, isEditing, onChange: setBirthdate})}
              {renderEthnicity()}
              {TableUtils.renderInputGroup(
                {title: "Eye Color", value: eyeColor, isEditing, onChange: setEyeColor})}
              {TableUtils.renderInputGroup(
                {title: "Country", value: country, isEditing, onChange: setCountry})}
              {TableUtils.renderInputGroup(
                {title: "Height (CM)", value: height, isEditing, onChange: setHeight})}
              {TableUtils.renderInputGroup(
                {title: "Measurements", value: measurements, isEditing, onChange: setMeasurements})}
              {TableUtils.renderInputGroup(
                {title: "Fake Tits", value: fakeTits, isEditing, onChange: setFakeTits})}
              {TableUtils.renderInputGroup(
                {title: "Career Length", value: careerLength, isEditing, onChange: setCareerLength})}
              {TableUtils.renderInputGroup(
                {title: "Tattoos", value: tattoos, isEditing, onChange: setTattoos})}
              {TableUtils.renderInputGroup(
                {title: "Piercings", value: piercings, isEditing, onChange: setPiercings})}
              {renderURLField()}
              {TableUtils.renderInputGroup(
                {title: "Twitter", value: twitter, isEditing, onChange: setTwitter})}
              {TableUtils.renderInputGroup(
                {title: "Instagram", value: instagram, isEditing, onChange: setInstagram})}
            </tbody>
          </Table>
        </div>
      </div>
    </>
  );
};
