/* eslint-disable react/no-this-in-sfc */

import React, { useEffect, useState } from "react";
import { Button, Form, Popover, OverlayTrigger, Table } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { StashService } from "src/core/StashService";
import {
  Icon,
  Modal,
  ImageInput,
  ScrapePerformerSuggest,
  LoadingIndicator
} from "src/components/Shared";
import { ImageUtils, TableUtils } from "src/utils";
import { useToast } from "src/hooks";

interface IPerformerDetails {
  performer: Partial<GQL.PerformerDataFragment>;
  isNew?: boolean;
  isEditing?: boolean;
  onSave?: (
    performer:
      | Partial<GQL.PerformerCreateInput>
      | Partial<GQL.PerformerUpdateInput>
  ) => void;
  onDelete?: () => void;
  onImageChange?: (image: string) => void;
}

export const PerformerDetailsPanel: React.FC<IPerformerDetails> = ({
  performer,
  isNew,
  isEditing,
  onSave,
  onDelete,
  onImageChange
}) => {
  const Toast = useToast();

  // Editing state
  const [isDisplayingScraperDialog, setIsDisplayingScraperDialog] = useState<
    GQL.Scraper
  >();
  const [scrapePerformerDetails, setScrapePerformerDetails] = useState<
    GQL.ScrapedPerformerDataFragment
  >();
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing performer state
  const [image, setImage] = useState<string>();
  const [name, setName] = useState<string>();
  const [aliases, setAliases] = useState<string>();
  const [favorite, setFavorite] = useState<boolean>();
  const [birthdate, setBirthdate] = useState<string>();
  const [ethnicity, setEthnicity] = useState<string>();
  const [country, setCountry] = useState<string>();
  const [eyeColor, setEyeColor] = useState<string>();
  const [height, setHeight] = useState<string>();
  const [measurements, setMeasurements] = useState<string>();
  const [fakeTits, setFakeTits] = useState<string>();
  const [careerLength, setCareerLength] = useState<string>();
  const [tattoos, setTattoos] = useState<string>();
  const [piercings, setPiercings] = useState<string>();
  const [url, setUrl] = useState<string>();
  const [twitter, setTwitter] = useState<string>();
  const [instagram, setInstagram] = useState<string>();

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const Scrapers = StashService.useListPerformerScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);

  function updatePerformerEditState(
    state: Partial<GQL.PerformerDataFragment | GQL.ScrapedPerformer>
  ) {
    if ((state as GQL.PerformerDataFragment).favorite !== undefined) {
      setFavorite((state as GQL.PerformerDataFragment).favorite);
    }
    setName(state.name ?? undefined);
    setAliases(state.aliases ?? undefined);
    setBirthdate(state.birthdate ?? undefined);
    setEthnicity(state.ethnicity ?? undefined);
    setCountry(state.country ?? undefined);
    setEyeColor(state.eye_color ?? undefined);
    setHeight(state.height ?? undefined);
    setMeasurements(state.measurements ?? undefined);
    setFakeTits(state.fake_tits ?? undefined);
    setCareerLength(state.career_length ?? undefined);
    setTattoos(state.tattoos ?? undefined);
    setPiercings(state.piercings ?? undefined);
    setUrl(state.url ?? undefined);
    setTwitter(state.twitter ?? undefined);
    setInstagram(state.instagram ?? undefined);
  }

  useEffect(() => {
    setImage(undefined);
    updatePerformerEditState(performer);
  }, [performer]);

  function onImageLoad(this: FileReader) {
    setImage(this.result as string);
    if (onImageChange) {
      onImageChange(this.result as string);
    }
  }

  if (isEditing) ImageUtils.usePasteImage(onImageLoad);

  useEffect(() => {
    const newQueryableScrapers = (
      Scrapers?.data?.listPerformerScrapers ?? []
    ).filter(s => s.performer?.supported_scrapes.includes(GQL.ScrapeType.Name));

    setQueryableScrapers(newQueryableScrapers);
  }, [Scrapers]);

  if (isLoading) return <LoadingIndicator />;

  function getPerformerInput() {
    const performerInput: Partial<
      GQL.PerformerCreateInput | GQL.PerformerUpdateInput
    > = {
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
      image
    };

    if (!isNew) {
      (performerInput as GQL.PerformerUpdateInput).id = performer.id!;
    }
    return performerInput;
  }

  function onImageChangeHandler(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function onDisplayFreeOnesDialog(scraper: GQL.Scraper) {
    setIsDisplayingScraperDialog(scraper);
  }

  function getQueryScraperPerformerInput() {
    if (!scrapePerformerDetails) return {};

    const { __typename, ...ret } = scrapePerformerDetails;
    return ret;
  }

  async function onScrapePerformer() {
    setIsDisplayingScraperDialog(undefined);
    setIsLoading(true);
    try {
      if (!scrapePerformerDetails || !isDisplayingScraperDialog) return;
      const result = await StashService.queryScrapePerformer(
        isDisplayingScraperDialog.id,
        getQueryScraperPerformerInput()
      );
      if (!result?.data?.scrapePerformer) return;
      updatePerformerEditState(result.data.scrapePerformer);
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  async function onScrapePerformerURL() {
    if (!url) return;
    setIsLoading(true);
    try {
      const result = await StashService.queryScrapePerformerURL(url);
      if (!result.data || !result.data.scrapePerformerURL) {
        return;
      }

      // leave URL as is if not set explicitly
      if (!result.data.scrapePerformerURL.url) {
        result.data.scrapePerformerURL.url = url;
      }
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
      isEditing: !!isEditing,
      onChange: (value: string) => setEthnicity(value),
      selectOptions: ["white", "black", "asian", "hispanic"]
    });
  }

  function renderScraperMenu() {
    if (!performer || !isEditing) {
      return;
    }

    const popover = (
      <Popover id="scraper-popover">
        <Popover.Content>
          <div>
            {queryableScrapers
              ? queryableScrapers.map(s => (
                  <Button
                    key={s.name}
                    variant="link"
                    onClick={() => onDisplayFreeOnesDialog(s)}
                  >
                    {s.name}
                  </Button>
                ))
              : ""}
          </div>
        </Popover.Content>
      </Popover>
    );

    return (
      <OverlayTrigger trigger="click" placement="bottom" overlay={popover}>
        <Button variant="secondary">Scrape with...</Button>
      </OverlayTrigger>
    );
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
            scraperId={
              isDisplayingScraperDialog ? isDisplayingScraperDialog.id : ""
            }
            onSelectPerformer={query => setScrapePerformerDetails(query)}
          />
        </div>
      </Modal>
    );
  }

  function urlScrapable(scrapedUrl: string) {
    return (
      !!scrapedUrl &&
      (Scrapers?.data?.listPerformerScrapers ?? []).some(s =>
        (s?.performer?.urls ?? []).some(u => scrapedUrl.includes(u))
      )
    );
  }

  function maybeRenderScrapeButton() {
    if (!url || !isEditing || !urlScrapable(url)) {
      return undefined;
    }
    return (
      <Button
        className="minimal scrape-url-button"
        onClick={() => onScrapePerformerURL()}
      >
        <Icon icon="file-upload" />
      </Button>
    );
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
            value={url ?? ''}
            readOnly={!isEditing}
            plaintext={!isEditing}
            placeholder="URL"
            onChange={(event: React.FormEvent<HTMLInputElement>) =>
              setUrl(event.currentTarget.value)
            }
          />
        </td>
      </tr>
    );
  }

  function maybeRenderButtons() {
    if (isEditing) {
      return (
        <div className="row">
          <Button
            className="edit-button"
            variant="primary"
            onClick={() => onSave?.(getPerformerInput())}
          >
            Save
          </Button>
          {!isNew ? (
            <Button
              className="edit-button"
              variant="danger"
              onClick={() => setIsDeleteAlertOpen(true)}
            >
              Delete
            </Button>
          ) : (
            ""
          )}
          {renderScraperMenu()}
          <ImageInput
            isEditing={!!isEditing}
            onImageChange={onImageChangeHandler}
          />
        </div>
      );
    }
  }

  function renderDeleteAlert() {
    return (
      <Modal
        show={isDeleteAlertOpen}
        icon="trash-alt"
        accept={{ text: "Delete", variant: "danger", onClick: onDelete }}
        cancel={{ onClick: () => setIsDeleteAlertOpen(false) }}
      >
        <p>Are you sure you want to delete {name}?</p>
      </Modal>
    );
  }

  function maybeRenderName() {
    if (isEditing) {
      return TableUtils.renderInputGroup({
        title: "Name",
        value: name,
        isEditing: !!isEditing,
        placeholder: "Name",
        onChange: setName
      });
    }
  }

  function maybeRenderAliases() {
    if (isEditing) {
      return TableUtils.renderInputGroup({
        title: "Aliases",
        value: aliases,
        isEditing: !!isEditing,
        placeholder: "Aliases",
        onChange: setAliases
      });
    }
  }

  return (
    <>
      {renderDeleteAlert()}
      {renderScraperDialog()}

      <Table id="performer-details" className="w-100">
        <tbody>
          {maybeRenderName()}
          {maybeRenderAliases()}
          {TableUtils.renderInputGroup({
            title: "Birthdate (YYYY-MM-DD)",
            value: birthdate,
            isEditing: !!isEditing,
            onChange: setBirthdate
          })}
          {renderEthnicity()}
          {TableUtils.renderInputGroup({
            title: "Eye Color",
            value: eyeColor,
            isEditing: !!isEditing,
            onChange: setEyeColor
          })}
          {TableUtils.renderInputGroup({
            title: "Country",
            value: country,
            isEditing: !!isEditing,
            onChange: setCountry
          })}
          {TableUtils.renderInputGroup({
            title: "Height (CM)",
            value: height,
            isEditing: !!isEditing,
            onChange: setHeight
          })}
          {TableUtils.renderInputGroup({
            title: "Measurements",
            value: measurements,
            isEditing: !!isEditing,
            onChange: setMeasurements
          })}
          {TableUtils.renderInputGroup({
            title: "Fake Tits",
            value: fakeTits,
            isEditing: !!isEditing,
            onChange: setFakeTits
          })}
          {TableUtils.renderInputGroup({
            title: "Career Length",
            value: careerLength,
            isEditing: !!isEditing,
            onChange: setCareerLength
          })}
          {TableUtils.renderInputGroup({
            title: "Tattoos",
            value: tattoos,
            isEditing: !!isEditing,
            onChange: setTattoos
          })}
          {TableUtils.renderInputGroup({
            title: "Piercings",
            value: piercings,
            isEditing: !!isEditing,
            onChange: setPiercings
          })}
          {renderURLField()}
          {TableUtils.renderInputGroup({
            title: "Twitter",
            value: twitter,
            isEditing: !!isEditing,
            onChange: setTwitter
          })}
          {TableUtils.renderInputGroup({
            title: "Instagram",
            value: instagram,
            isEditing: !!isEditing,
            onChange: setInstagram
          })}
        </tbody>
      </Table>

      {maybeRenderButtons()}
    </>
  );
};
