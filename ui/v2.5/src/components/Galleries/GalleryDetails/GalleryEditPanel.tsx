import React, { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { Button, Form, Col, Row } from "react-bootstrap";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import {
  queryScrapeGalleryURL,
  useGalleryCreate,
  useGalleryUpdate,
  useListGalleryScrapers,
} from "src/core/StashService";
import {
  PerformerSelect,
  TagSelect,
  StudioSelect,
  Icon,
  LoadingIndicator,
} from "src/components/Shared";
import { useToast } from "src/hooks";
import { FormUtils, EditableTextUtils } from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { GalleryScrapeDialog } from "./GalleryScrapeDialog";

interface IProps {
  isVisible: boolean;
  onDelete: () => void;
}

interface INewProps {
  isNew: true;
  gallery: undefined;
}

interface IExistingProps {
  isNew: false;
  gallery: GQL.GalleryDataFragment;
}

export const GalleryEditPanel: React.FC<
  IProps & (INewProps | IExistingProps)
> = (props) => {
  const Toast = useToast();
  const history = useHistory();
  const [title, setTitle] = useState<string>();
  const [details, setDetails] = useState<string>();
  const [url, setUrl] = useState<string>();
  const [date, setDate] = useState<string>();
  const [rating, setRating] = useState<number>();
  const [studioId, setStudioId] = useState<string>();
  const [performerIds, setPerformerIds] = useState<string[]>();
  const [tagIds, setTagIds] = useState<string[]>();

  const Scrapers = useListGalleryScrapers();

  const [
    scrapedGallery,
    setScrapedGallery,
  ] = useState<GQL.ScrapedGallery | null>();

  // Network state
  const [isLoading, setIsLoading] = useState(true);

  const [createGallery] = useGalleryCreate(
    getGalleryInput() as GQL.GalleryCreateInput
  );
  const [updateGallery] = useGalleryUpdate();

  useEffect(() => {
    if (props.isVisible) {
      Mousetrap.bind("s s", () => {
        onSave();
      });
      Mousetrap.bind("d d", () => {
        props.onDelete();
      });

      // numeric keypresses get caught by jwplayer, so blur the element
      // if the rating sequence is started
      Mousetrap.bind("r", () => {
        if (document.activeElement instanceof HTMLElement) {
          document.activeElement.blur();
        }

        Mousetrap.bind("0", () => setRating(NaN));
        Mousetrap.bind("1", () => setRating(1));
        Mousetrap.bind("2", () => setRating(2));
        Mousetrap.bind("3", () => setRating(3));
        Mousetrap.bind("4", () => setRating(4));
        Mousetrap.bind("5", () => setRating(5));

        setTimeout(() => {
          Mousetrap.unbind("0");
          Mousetrap.unbind("1");
          Mousetrap.unbind("2");
          Mousetrap.unbind("3");
          Mousetrap.unbind("4");
          Mousetrap.unbind("5");
        }, 1000);
      });

      return () => {
        Mousetrap.unbind("s s");
        Mousetrap.unbind("d d");

        Mousetrap.unbind("r");
      };
    }
  });

  function updateGalleryEditState(state?: GQL.GalleryDataFragment) {
    const perfIds = state?.performers?.map((performer) => performer.id);
    const tIds = state?.tags ? state?.tags.map((tag) => tag.id) : undefined;

    setTitle(state?.title ?? undefined);
    setDetails(state?.details ?? undefined);
    setUrl(state?.url ?? undefined);
    setDate(state?.date ?? undefined);
    setRating(state?.rating === null ? NaN : state?.rating);
    setStudioId(state?.studio?.id ?? undefined);
    setPerformerIds(perfIds);
    setTagIds(tIds);
  }

  useEffect(() => {
    updateGalleryEditState(props.gallery);
    setIsLoading(false);
  }, [props.gallery]);

  function getGalleryInput() {
    return {
      id: props.isNew ? undefined : props.gallery.id,
      title,
      details,
      url,
      date,
      rating: rating ?? null,
      studio_id: studioId ?? null,
      performer_ids: performerIds,
      tag_ids: tagIds,
    };
  }

  async function onSave() {
    setIsLoading(true);
    try {
      if (props.isNew) {
        const result = await createGallery();
        if (result.data?.galleryCreate) {
          history.push(`/galleries/${result.data.galleryCreate.id}`);
          Toast.success({ content: "Created gallery" });
        }
      } else {
        const result = await updateGallery({
          variables: {
            input: getGalleryInput() as GQL.GalleryUpdateInput,
          },
        });
        if (result.data?.galleryUpdate) {
          Toast.success({ content: "Updated gallery" });
        }
      }
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  function onScrapeDialogClosed(gallery?: GQL.ScrapedGalleryDataFragment) {
    if (gallery) {
      updateGalleryFromScrapedGallery(gallery);
    }
    setScrapedGallery(undefined);
  }

  function maybeRenderScrapeDialog() {
    if (!scrapedGallery) {
      return;
    }

    const currentGallery = getGalleryInput();

    return (
      <GalleryScrapeDialog
        gallery={currentGallery}
        scraped={scrapedGallery}
        onClose={(gallery) => {
          onScrapeDialogClosed(gallery);
        }}
      />
    );
  }

  function urlScrapable(scrapedUrl: string): boolean {
    return (Scrapers?.data?.listGalleryScrapers ?? []).some((s) =>
      (s?.gallery?.urls ?? []).some((u) => scrapedUrl.includes(u))
    );
  }

  function updateGalleryFromScrapedGallery(
    gallery: GQL.ScrapedGalleryDataFragment
  ) {
    if (gallery.title) {
      setTitle(gallery.title);
    }

    if (gallery.details) {
      setDetails(gallery.details);
    }

    if (gallery.date) {
      setDate(gallery.date);
    }

    if (gallery.url) {
      setUrl(gallery.url);
    }

    if (gallery.studio && gallery.studio.stored_id) {
      setStudioId(gallery.studio.stored_id);
    }

    if (gallery.performers && gallery.performers.length > 0) {
      const idPerfs = gallery.performers.filter((p) => {
        return p.stored_id !== undefined && p.stored_id !== null;
      });

      if (idPerfs.length > 0) {
        const newIds = idPerfs.map((p) => p.stored_id);
        setPerformerIds(newIds as string[]);
      }
    }

    if (gallery?.tags?.length) {
      const idTags = gallery.tags.filter((p) => {
        return p.stored_id !== undefined && p.stored_id !== null;
      });

      if (idTags.length > 0) {
        const newIds = idTags.map((p) => p.stored_id);
        setTagIds(newIds as string[]);
      }
    }
  }

  async function onScrapeGalleryURL() {
    if (!url) {
      return;
    }
    setIsLoading(true);
    try {
      const result = await queryScrapeGalleryURL(url);
      if (!result || !result.data || !result.data.scrapeGalleryURL) {
        return;
      }
      setScrapedGallery(result.data.scrapeGalleryURL);
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function maybeRenderScrapeButton() {
    if (!url || !urlScrapable(url)) {
      return undefined;
    }
    return (
      <Button
        className="minimal scrape-url-button"
        onClick={onScrapeGalleryURL}
        title="Scrape"
      >
        <Icon className="fa-fw" icon="file-download" />
      </Button>
    );
  }

  if (isLoading) return <LoadingIndicator />;

  return (
    <div id="gallery-edit-details">
      {maybeRenderScrapeDialog()}
      <div className="form-container row px-3 pt-3">
        <div className="col edit-buttons mb-3 pl-0">
          <Button className="edit-button" variant="primary" onClick={onSave}>
            Save
          </Button>
          <Button
            className="edit-button"
            variant="danger"
            onClick={() => props.onDelete()}
          >
            Delete
          </Button>
        </div>
      </div>
      <div className="form-container row px-3">
        <div className="col-12 col-lg-6 col-xl-12">
          {FormUtils.renderInputGroup({
            title: "Title",
            value: title,
            onChange: setTitle,
            isEditing: true,
          })}
          <Form.Group controlId="url" as={Row}>
            <Col xs={3} className="pr-0 url-label">
              <Form.Label className="col-form-label">URL</Form.Label>
              <div className="float-right scrape-button-container">
                {maybeRenderScrapeButton()}
              </div>
            </Col>
            <Col xs={9}>
              {EditableTextUtils.renderInputGroup({
                title: "URL",
                value: url,
                onChange: setUrl,
                isEditing: true,
              })}
            </Col>
          </Form.Group>
          {FormUtils.renderInputGroup({
            title: "Date",
            value: date,
            isEditing: true,
            onChange: setDate,
            placeholder: "YYYY-MM-DD",
          })}
          <Form.Group controlId="rating" as={Row}>
            {FormUtils.renderLabel({
              title: "Rating",
            })}
            <Col xs={9}>
              <RatingStars
                value={rating}
                onSetRating={(value) => setRating(value)}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="studio" as={Row}>
            {FormUtils.renderLabel({
              title: "Studio",
            })}
            <Col xs={9}>
              <StudioSelect
                onSelect={(items) =>
                  setStudioId(items.length > 0 ? items[0]?.id : undefined)
                }
                ids={studioId ? [studioId] : []}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="performers" as={Row}>
            {FormUtils.renderLabel({
              title: "Performers",
              labelProps: {
                column: true,
                sm: 3,
                xl: 12,
              },
            })}
            <Col sm={9} xl={12}>
              <PerformerSelect
                isMulti
                onSelect={(items) =>
                  setPerformerIds(items.map((item) => item.id))
                }
                ids={performerIds}
              />
            </Col>
          </Form.Group>

          <Form.Group controlId="tags" as={Row}>
            {FormUtils.renderLabel({
              title: "Tags",
              labelProps: {
                column: true,
                sm: 3,
                xl: 12,
              },
            })}
            <Col sm={9} xl={12}>
              <TagSelect
                isMulti
                onSelect={(items) => setTagIds(items.map((item) => item.id))}
                ids={tagIds}
              />
            </Col>
          </Form.Group>
        </div>
        <div className="col-12 col-lg-6 col-xl-12">
          <Form.Group controlId="details">
            <Form.Label>Details</Form.Label>
            <Form.Control
              as="textarea"
              className="gallery-description text-input"
              onChange={(newValue: React.ChangeEvent<HTMLTextAreaElement>) =>
                setDetails(newValue.currentTarget.value)
              }
              value={details}
            />
          </Form.Group>
        </div>
      </div>
    </div>
  );
};
