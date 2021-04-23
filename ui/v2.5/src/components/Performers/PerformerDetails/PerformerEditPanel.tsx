import React, { useEffect, useState } from "react";
import {
  Button,
  Popover,
  OverlayTrigger,
  Form,
  Col,
  InputGroup,
  Row,
  Badge,
} from "react-bootstrap";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import {
  getGenderStrings,
  useListPerformerScrapers,
  genderToString,
  stringToGender,
  queryScrapePerformer,
  mutateReloadScrapers,
  usePerformerUpdate,
  usePerformerCreate,
  useTagCreate,
  queryScrapePerformerURL,
} from "src/core/StashService";
import {
  Icon,
  ImageInput,
  LoadingIndicator,
  CollapseButton,
  Modal,
  TagSelect,
} from "src/components/Shared";
import { ImageUtils } from "src/utils";
import { useToast } from "src/hooks";
import { Prompt, useHistory } from "react-router-dom";
import { useFormik } from "formik";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { PerformerScrapeDialog } from "./PerformerScrapeDialog";
import PerformerScrapeModal from "./PerformerScrapeModal";

interface IPerformerDetails {
  performer: Partial<GQL.PerformerDataFragment>;
  isNew?: boolean;
  isVisible: boolean;
  onDelete?: () => void;
  onImageChange?: (image?: string | null) => void;
  onImageEncoding?: (loading?: boolean) => void;
}

export const PerformerEditPanel: React.FC<IPerformerDetails> = ({
  performer,
  isNew,
  isVisible,
  onDelete,
  onImageChange,
  onImageEncoding,
}) => {
  const Toast = useToast();
  const history = useHistory();

  // Editing state
  const [scraper, setScraper] = useState<GQL.Scraper | undefined>();
  const [newTags, setNewTags] = useState<GQL.ScrapedSceneTag[]>();
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const [updatePerformer] = usePerformerUpdate();
  const [createPerformer] = usePerformerCreate();

  const Scrapers = useListPerformerScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);

  const [scrapedPerformer, setScrapedPerformer] = useState<
    GQL.ScrapedPerformer | undefined
  >();

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, true);

  const [createTag] = useTagCreate({ name: "" });

  const genderOptions = [""].concat(getGenderStrings());

  const labelXS = 3;
  const labelXL = 2;
  const fieldXS = 9;
  const fieldXL = 7;

  const schema = yup.object({
    name: yup.string().required(),
    aliases: yup.string().optional(),
    gender: yup.string().optional().oneOf(genderOptions),
    birthdate: yup.string().optional(),
    ethnicity: yup.string().optional(),
    eye_color: yup.string().optional(),
    country: yup.string().optional(),
    height: yup.string().optional(),
    measurements: yup.string().optional(),
    fake_tits: yup.string().optional(),
    career_length: yup.string().optional(),
    tattoos: yup.string().optional(),
    piercings: yup.string().optional(),
    url: yup.string().optional(),
    twitter: yup.string().optional(),
    instagram: yup.string().optional(),
    tag_ids: yup.array(yup.string().required()).optional(),
    stash_ids: yup.mixed<GQL.StashIdInput>().optional(),
    image: yup.string().optional().nullable(),
    rating: yup.number().optional().nullable(),
    details: yup.string().optional(),
    death_date: yup.string().optional(),
    hair_color: yup.string().optional(),
    weight: yup.number().optional(),
  });

  const initialValues = {
    name: performer.name ?? "",
    aliases: performer.aliases ?? "",
    gender: genderToString(performer.gender ?? undefined),
    birthdate: performer.birthdate ?? "",
    ethnicity: performer.ethnicity ?? "",
    eye_color: performer.eye_color ?? "",
    country: performer.country ?? "",
    height: performer.height ?? "",
    measurements: performer.measurements ?? "",
    fake_tits: performer.fake_tits ?? "",
    career_length: performer.career_length ?? "",
    tattoos: performer.tattoos ?? "",
    piercings: performer.piercings ?? "",
    url: performer.url ?? "",
    twitter: performer.twitter ?? "",
    instagram: performer.instagram ?? "",
    tag_ids: (performer.tags ?? []).map((t) => t.id),
    stash_ids: performer.stash_ids ?? undefined,
    image: undefined,
    rating: performer.rating ?? undefined,
    details: performer.details ?? "",
    death_date: performer.death_date ?? "",
    hair_color: performer.hair_color ?? "",
    weight: performer.weight ?? "",
  };

  type InputValues = typeof initialValues;

  const formik = useFormik({
    initialValues,
    validationSchema: schema,
    onSubmit: (values) => onSave(getPerformerInput(values)),
  });

  function setRating(v: number) {
    formik.setFieldValue("rating", v);
  }

  function translateScrapedGender(scrapedGender?: string) {
    if (!scrapedGender) {
      return;
    }

    let retEnum: GQL.GenderEnum | undefined;

    // try to translate from enum values first
    const upperGender = scrapedGender?.toUpperCase();
    const asEnum = genderToString(upperGender as GQL.GenderEnum);
    if (asEnum) {
      retEnum = stringToGender(asEnum);
    } else {
      // try to match against gender strings
      const caseInsensitive = true;
      retEnum = stringToGender(scrapedGender, caseInsensitive);
    }

    return genderToString(retEnum);
  }

  function renderNewTags() {
    if (!newTags || newTags.length === 0) {
      return;
    }

    const ret = (
      <>
        {newTags.map((t) => (
          <Badge
            className="tag-item"
            variant="secondary"
            key={t.name}
            onClick={() => createNewTag(t)}
          >
            {t.name}
            <Button className="minimal ml-2">
              <Icon className="fa-fw" icon="plus" />
            </Button>
          </Badge>
        ))}
      </>
    );

    const minCollapseLength = 10;

    if (newTags.length >= minCollapseLength) {
      return (
        <CollapseButton text={`Missing (${newTags.length})`}>
          {ret}
        </CollapseButton>
      );
    }

    return ret;
  }

  async function createNewTag(toCreate: GQL.ScrapedSceneTag) {
    let tagInput: GQL.TagCreateInput = { name: "" };
    try {
      tagInput = Object.assign(tagInput, toCreate);
      const result = await createTag({
        variables: tagInput,
      });

      // add the new tag to the new tags value
      const newTagIds = formik.values.tag_ids.concat([
        result.data!.tagCreate!.id,
      ]);
      formik.setFieldValue("tag_ids", newTagIds);

      // remove the tag from the list
      const newTagsClone = newTags!.concat();
      const pIndex = newTagsClone.indexOf(toCreate);
      newTagsClone.splice(pIndex, 1);

      setNewTags(newTagsClone);

      Toast.success({
        content: (
          <span>
            Created tag: <b>{toCreate.name}</b>
          </span>
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  function updatePerformerEditStateFromScraper(
    state: Partial<GQL.ScrapedPerformerDataFragment>
  ) {
    if (state.name) {
      formik.setFieldValue("name", state.name);
    }

    if (state.aliases) {
      formik.setFieldValue("aliases", state.aliases);
    }
    if (state.birthdate) {
      formik.setFieldValue("birthdate", state.birthdate);
    }
    if (state.ethnicity) {
      formik.setFieldValue("ethnicity", state.ethnicity);
    }
    if (state.country) {
      formik.setFieldValue("country", state.country);
    }
    if (state.eye_color) {
      formik.setFieldValue("eye_color", state.eye_color);
    }
    if (state.height) {
      formik.setFieldValue("height", state.height);
    }
    if (state.measurements) {
      formik.setFieldValue("measurements", state.measurements);
    }
    if (state.fake_tits) {
      formik.setFieldValue("fake_tits", state.fake_tits);
    }
    if (state.career_length) {
      formik.setFieldValue("career_length", state.career_length);
    }
    if (state.tattoos) {
      formik.setFieldValue("tattoos", state.tattoos);
    }
    if (state.piercings) {
      formik.setFieldValue("piercings", state.piercings);
    }
    if (state.url) {
      formik.setFieldValue("url", state.url);
    }
    if (state.twitter) {
      formik.setFieldValue("twitter", state.twitter);
    }
    if (state.instagram) {
      formik.setFieldValue("instagram", state.instagram);
    }
    if (state.gender) {
      // gender is a string in the scraper data
      formik.setFieldValue(
        "gender",
        translateScrapedGender(state.gender ?? undefined)
      );
    }
    if (state.tags) {
      // map tags to their ids and filter out those not found
      const newTagIds = state.tags.map((t) => t.stored_id).filter((t) => t);
      formik.setFieldValue("tag_ids", newTagIds as string[]);

      setNewTags(state.tags.filter((t) => !t.stored_id));
    }

    // image is a base64 string
    // #404: don't overwrite image if it has been modified by the user
    // overwrite if not new since it came from a dialog
    // otherwise follow existing behaviour
    if (
      (!isNew || formik.values.image === undefined) &&
      (state as GQL.ScrapedPerformerDataFragment).image !== undefined
    ) {
      const imageStr = (state as GQL.ScrapedPerformerDataFragment).image;
      formik.setFieldValue("image", imageStr ?? undefined);
    }
    if (state.details) {
      formik.setFieldValue("details", state.details);
    }
    if (state.death_date) {
      formik.setFieldValue("death_date", state.death_date);
    }
    if (state.hair_color) {
      formik.setFieldValue("hair_color", state.hair_color);
    }
    if (state.weight) {
      formik.setFieldValue("weight", state.weight);
    }
  }

  function onImageLoad(imageData: string) {
    formik.setFieldValue("image", imageData);
  }

  async function onSave(
    performerInput:
      | Partial<GQL.PerformerCreateInput>
      | Partial<GQL.PerformerUpdateInput>
  ) {
    setIsLoading(true);
    try {
      if (!isNew) {
        await updatePerformer({
          variables: {
            input: {
              ...performerInput,
              stash_ids: performerInput?.stash_ids?.map((s) => ({
                endpoint: s.endpoint,
                stash_id: s.stash_id,
              })),
            } as GQL.PerformerUpdateInput,
          },
        });
        history.push(`/performers/${performer.id}`);
      } else {
        const result = await createPerformer({
          variables: { input: performerInput as GQL.PerformerCreateInput },
        });
        if (result.data?.performerCreate) {
          history.push(`/performers/${result.data.performerCreate.id}`);
        }
      }
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  // set up hotkeys
  useEffect(() => {
    if (isVisible) {
      Mousetrap.bind("s s", () => {
        onSave?.(getPerformerInput(formik.values));
      });

      if (!isNew) {
        Mousetrap.bind("d d", () => {
          setIsDeleteAlertOpen(true);
        });
      }

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

        if (!isNew) {
          Mousetrap.unbind("d d");
        }
      };
    }
  });

  useEffect(() => {
    if (onImageChange) {
      onImageChange(formik.values.image);
    }
    return () => onImageChange?.();
  }, [formik.values.image, onImageChange]);

  useEffect(() => onImageEncoding?.(imageEncoding), [
    onImageEncoding,
    imageEncoding,
  ]);

  useEffect(() => {
    const newQueryableScrapers = (
      Scrapers?.data?.listPerformerScrapers ?? []
    ).filter((s) =>
      s.performer?.supported_scrapes.includes(GQL.ScrapeType.Name)
    );

    setQueryableScrapers(newQueryableScrapers);
  }, [Scrapers]);

  if (isLoading) return <LoadingIndicator />;

  function getPerformerInput(values: InputValues) {
    const performerInput: Partial<
      GQL.PerformerCreateInput | GQL.PerformerUpdateInput
    > = {
      ...values,
      gender: stringToGender(values.gender),
      rating: values.rating ?? null,
      weight: Number(values.weight),
    };

    if (!isNew) {
      (performerInput as GQL.PerformerUpdateInput).id = performer.id!;
    }
    return performerInput;
  }

  function onImageChangeHandler(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function onImageChangeURL(url: string) {
    formik.setFieldValue("image", url);
  }

  async function onReloadScrapers() {
    setIsLoading(true);
    try {
      await mutateReloadScrapers();

      // reload the performer scrapers
      await Scrapers.refetch();
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function onScrapePerformer(
    selectedPerformer: GQL.ScrapedPerformerDataFragment
  ) {
    setScraper(undefined);
    try {
      if (!scraper) return;
      setIsLoading(true);

      const {
        __typename,
        image: _image,
        tags: _tags,
        ...ret
      } = selectedPerformer;

      const result = await queryScrapePerformer(scraper.id, ret);
      if (!result?.data?.scrapePerformer) return;

      // if this is a new performer, just dump the data
      if (isNew) {
        updatePerformerEditStateFromScraper(result.data.scrapePerformer);
      } else {
        setScrapedPerformer(result.data.scrapePerformer);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  async function onScrapePerformerURL() {
    const { url } = formik.values;
    if (!url) return;
    setIsLoading(true);
    try {
      const result = await queryScrapePerformerURL(url);
      if (!result.data || !result.data.scrapePerformerURL) {
        return;
      }

      // if this is a new performer, just dump the data
      if (isNew) {
        updatePerformerEditStateFromScraper(result.data.scrapePerformerURL);
      } else {
        setScrapedPerformer(result.data.scrapePerformerURL);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsLoading(false);
    }
  }

  function renderScraperMenu() {
    if (!performer) {
      return;
    }

    const popover = (
      <Popover id="performer-scraper-popover">
        <Popover.Content>
          <>
            {queryableScrapers
              ? queryableScrapers.map((s) => (
                  <div key={s.name}>
                    <Button
                      key={s.name}
                      className="minimal"
                      onClick={() => setScraper(s)}
                    >
                      {s.name}
                    </Button>
                  </div>
                ))
              : ""}
            <div>
              <Button className="minimal" onClick={() => onReloadScrapers()}>
                <span className="fa-icon">
                  <Icon icon="sync-alt" />
                </span>
                <span>Reload scrapers</span>
              </Button>
            </div>
          </>
        </Popover.Content>
      </Popover>
    );

    return (
      <OverlayTrigger
        trigger="click"
        placement="top"
        overlay={popover}
        rootClose
      >
        <Button variant="secondary" className="mr-2">
          Scrape with...
        </Button>
      </OverlayTrigger>
    );
  }

  function urlScrapable(scrapedUrl?: string) {
    return (
      !!scrapedUrl &&
      (Scrapers?.data?.listPerformerScrapers ?? []).some((s) =>
        (s?.performer?.urls ?? []).some((u) => scrapedUrl.includes(u))
      )
    );
  }

  function maybeRenderScrapeDialog() {
    if (!scrapedPerformer) {
      return;
    }

    const currentPerformer: Partial<GQL.PerformerUpdateInput> = {
      ...formik.values,
      gender: stringToGender(formik.values.gender),
      image: formik.values.image ?? performer.image_path,
      weight: Number(formik.values.weight),
    };

    return (
      <PerformerScrapeDialog
        performer={currentPerformer}
        scraped={scrapedPerformer}
        onClose={(p) => {
          onScrapeDialogClosed(p);
        }}
      />
    );
  }

  function onScrapeDialogClosed(p?: GQL.ScrapedPerformerDataFragment) {
    if (p) {
      updatePerformerEditStateFromScraper(p);
    }
    setScrapedPerformer(undefined);
  }

  function maybeRenderScrapeButton() {
    return (
      <Button
        variant="secondary"
        disabled={!urlScrapable(formik.values.url)}
        className="scrape-url-button text-input"
        onClick={() => onScrapePerformerURL()}
      >
        <Icon icon="file-upload" />
      </Button>
    );
  }

  function renderButtons() {
    return (
      <Row>
        <Col className="mt-3" xs={12}>
          <Button
            className="mr-2"
            variant="primary"
            onClick={() => formik.submitForm()}
          >
            Save
          </Button>
          {!isNew ? (
            <Button
              className="mr-2"
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
            isEditing
            onImageChange={onImageChangeHandler}
            onImageURL={onImageChangeURL}
          />
          <Button
            className="mx-2"
            variant="danger"
            onClick={() => formik.setFieldValue("image", null)}
          >
            Clear image
          </Button>
        </Col>
      </Row>
    );
  }

  const renderScrapeModal = () =>
    scraper !== undefined && (
      <PerformerScrapeModal
        scraper={scraper}
        onHide={() => setScraper(undefined)}
        onSelectPerformer={onScrapePerformer}
        name={formik.values.name || ""}
      />
    );

  function renderDeleteAlert() {
    return (
      <Modal
        show={isDeleteAlertOpen}
        icon="trash-alt"
        accept={{ text: "Delete", variant: "danger", onClick: onDelete }}
        cancel={{ onClick: () => setIsDeleteAlertOpen(false) }}
      >
        <p>Are you sure you want to delete {performer.name}?</p>
      </Modal>
    );
  }

  function renderTagsField() {
    return (
      <Form.Group controlId="tags" as={Row}>
        <Form.Label column sm={labelXS} xl={labelXL}>
          Tags
        </Form.Label>
        <Col xs={fieldXS} xl={fieldXL}>
          <TagSelect
            menuPortalTarget={document.body}
            isMulti
            onSelect={(items) =>
              formik.setFieldValue(
                "tag_ids",
                items.map((item) => item.id)
              )
            }
            ids={formik.values.tag_ids}
          />
          {renderNewTags()}
        </Col>
      </Form.Group>
    );
  }

  const removeStashID = (stashID: GQL.StashIdInput) => {
    formik.setFieldValue(
      "stash_ids",
      (formik.values.stash_ids ?? []).filter(
        (s) =>
          !(s.endpoint === stashID.endpoint && s.stash_id === stashID.stash_id)
      )
    );
  };

  function renderStashIDs() {
    if (!formik.values.stash_ids?.length) {
      return;
    }

    return (
      <Row>
        <Form.Label column sm={labelXS} xl={labelXL}>
          StashIDs
        </Form.Label>
        <Col sm={fieldXS} xl={fieldXL}>
          <ul className="pl-0">
            {formik.values.stash_ids.map((stashID) => {
              const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
              const link = base ? (
                <a
                  href={`${base}performers/${stashID.stash_id}`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {stashID.stash_id}
                </a>
              ) : (
                stashID.stash_id
              );
              return (
                <li key={stashID.stash_id} className="row no-gutters mb-1">
                  <Button
                    variant="danger"
                    className="mr-2 py-0"
                    title="Delete StashID"
                    onClick={() => removeStashID(stashID)}
                  >
                    <Icon icon="trash-alt" />
                  </Button>
                  {link}
                </li>
              );
            })}
          </ul>
        </Col>
      </Row>
    );
  }

  function renderTextField(field: string, title: string, placeholder?: string) {
    return (
      <Form.Group controlId={field} as={Row}>
        <Form.Label column xs={labelXS} xl={labelXL}>
          {title}
        </Form.Label>
        <Col xs={fieldXS} xl={fieldXL}>
          <Form.Control
            className="text-input"
            placeholder={placeholder ?? title}
            {...formik.getFieldProps(field)}
            isInvalid={!!formik.getFieldMeta(field).error}
          />
        </Col>
      </Form.Group>
    );
  }

  return (
    <>
      {renderDeleteAlert()}
      {renderScrapeModal()}
      {maybeRenderScrapeDialog()}

      <Prompt
        when={formik.dirty}
        message="Unsaved changes. Are you sure you want to leave?"
      />

      <Form noValidate onSubmit={formik.handleSubmit} id="performer-edit">
        <Form.Group controlId="name" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            Name
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <Form.Control
              className="text-input"
              placeholder="Name"
              {...formik.getFieldProps("name")}
              isInvalid={!!formik.errors.name}
            />
            <Form.Control.Feedback type="invalid">
              {formik.errors.name}
            </Form.Control.Feedback>
          </Col>
        </Form.Group>

        <Form.Group controlId="aliases" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            Alias
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <Form.Control
              as="textarea"
              className="text-input"
              placeholder="Alias"
              {...formik.getFieldProps("aliases")}
            />
          </Col>
        </Form.Group>

        <Form.Group as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            Gender
          </Form.Label>
          <Col xs="auto">
            <Form.Control
              as="select"
              className="input-control"
              {...formik.getFieldProps("gender")}
            >
              {genderOptions.map((opt) => (
                <option value={opt} key={opt}>
                  {opt}
                </option>
              ))}
            </Form.Control>
          </Col>
        </Form.Group>

        {renderTextField("birthdate", "Birthdate", "YYYY-MM-DD")}
        {renderTextField("death_date", "Death Date", "YYYY-MM-DD")}
        {renderTextField("country", "Country")}
        {renderTextField("ethnicity", "Ethnicity")}
        {renderTextField("hair_color", "Hair Color")}
        {renderTextField("eye_color", "Eye Color")}
        {renderTextField("height", "Height (cm)")}
        {renderTextField("weight", "Weight (kg)")}
        {renderTextField("measurements", "Measurements")}
        {renderTextField("fake_tits", "Fake Tits")}

        <Form.Group controlId="tattoos" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            Tattoos
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <Form.Control
              as="textarea"
              className="text-input"
              placeholder="Tattoos"
              {...formik.getFieldProps("tattoos")}
            />
          </Col>
        </Form.Group>

        <Form.Group controlId="piercings" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            Piercings
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <Form.Control
              as="textarea"
              className="text-input"
              placeholder="Piercings"
              {...formik.getFieldProps("piercings")}
            />
          </Col>
        </Form.Group>

        {renderTextField("career_length", "Career Length")}

        <Form.Group controlId="name" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            URL
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <InputGroup>
              <Form.Control
                className="text-input"
                placeholder="URL"
                {...formik.getFieldProps("url")}
              />
              <InputGroup.Append>{maybeRenderScrapeButton()}</InputGroup.Append>
            </InputGroup>
          </Col>
        </Form.Group>

        {renderTextField("twitter", "Twitter")}
        {renderTextField("instagram", "Instagram")}
        <Form.Group controlId="details" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            Details
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <Form.Control
              as="textarea"
              className="text-input"
              placeholder="Details"
              {...formik.getFieldProps("details")}
            />
          </Col>
        </Form.Group>
        {renderTagsField()}

        <Form.Group controlId="rating" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            Rating
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <RatingStars
              value={formik.values.rating ?? undefined}
              onSetRating={(value) => formik.setFieldValue("rating", value)}
            />
          </Col>
        </Form.Group>
        {renderStashIDs()}

        {renderButtons()}
      </Form>
    </>
  );
};
