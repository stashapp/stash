import React, { useEffect, useState } from "react";
import { Button, Form, Col, Row, Badge, Dropdown } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import Mousetrap from "mousetrap";
import * as GQL from "src/core/generated-graphql";
import * as yup from "yup";
import {
  useListPerformerScrapers,
  queryScrapePerformer,
  mutateReloadScrapers,
  useTagCreate,
  queryScrapePerformerURL,
} from "src/core/StashService";
import { Icon } from "src/components/Shared/Icon";
import { ImageInput } from "src/components/Shared/ImageInput";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { CollapseButton } from "src/components/Shared/CollapseButton";
import { TagSelect } from "src/components/Shared/Select";
import { CountrySelect } from "src/components/Shared/CountrySelect";
import { URLField } from "src/components/Shared/URLField";
import ImageUtils from "src/utils/image";
import { getStashIDs } from "src/utils/stashIds";
import { stashboxDisplayName } from "src/utils/stashbox";
import { useToast } from "src/hooks/Toast";
import { Prompt } from "react-router-dom";
import { useFormik } from "formik";
import {
  genderToString,
  stringGenderMap,
  stringToGender,
} from "src/utils/gender";
import {
  circumcisedToString,
  stringCircumMap,
  stringToCircumcised,
} from "src/utils/circumcised";
import { ConfigurationContext } from "src/hooks/Config";
import { PerformerScrapeDialog } from "./PerformerScrapeDialog";
import PerformerScrapeModal from "./PerformerScrapeModal";
import PerformerStashBoxModal, { IStashBox } from "./PerformerStashBoxModal";
import cx from "classnames";
import {
  faPlus,
  faSyncAlt,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";
import { StringListInput } from "src/components/Shared/StringListInput";
import isEqual from "lodash-es/isEqual";
import { DateInput } from "src/components/Shared/DateInput";

const isScraper = (
  scraper: GQL.Scraper | GQL.StashBox
): scraper is GQL.Scraper => (scraper as GQL.Scraper).id !== undefined;

interface IPerformerDetails {
  performer: Partial<GQL.PerformerDataFragment>;
  isVisible: boolean;
  onSubmit: (performer: GQL.PerformerCreateInput) => Promise<void>;
  onCancel?: () => void;
  setImage: (image?: string | null) => void;
  setEncodingImage: (loading: boolean) => void;
}

export const PerformerEditPanel: React.FC<IPerformerDetails> = ({
  performer,
  isVisible,
  onSubmit,
  onCancel,
  setImage,
  setEncodingImage,
}) => {
  const Toast = useToast();

  const isNew = performer.id === undefined;

  // Editing state
  const [scraper, setScraper] = useState<GQL.Scraper | IStashBox>();
  const [newTags, setNewTags] = useState<GQL.ScrapedTag[]>();
  const [isScraperModalOpen, setIsScraperModalOpen] = useState<boolean>(false);

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const Scrapers = useListPerformerScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);

  const [scrapedPerformer, setScrapedPerformer] =
    useState<GQL.ScrapedPerformer>();
  const { configuration: stashConfig } = React.useContext(ConfigurationContext);

  const [createTag] = useTagCreate();
  const intl = useIntl();

  const labelXS = 3;
  const labelXL = 2;
  const fieldXS = 9;
  const fieldXL = 7;

  const schema = yup.object({
    name: yup.string().required(),
    disambiguation: yup.string().ensure(),
    alias_list: yup
      .array(yup.string().required())
      .defined()
      .test({
        name: "unique",
        test: (value, context) => {
          const aliases = [context.parent.name, ...value];
          const dupes = aliases
            .map((e, i, a) => {
              if (a.indexOf(e) !== i) {
                return String(i - 1);
              } else {
                return null;
              }
            })
            .filter((e) => e !== null) as string[];
          if (dupes.length === 0) return true;
          return new yup.ValidationError(dupes.join(" "), value, "alias_list");
        },
      }),
    gender: yup.string<GQL.GenderEnum | "">().ensure(),
    birthdate: yup
      .string()
      .ensure()
      .test({
        name: "date",
        test: (value) => {
          if (!value) return true;
          if (!value.match(/^\d{4}-\d{2}-\d{2}$/)) return false;
          if (Number.isNaN(Date.parse(value))) return false;
          return true;
        },
        message: intl.formatMessage({ id: "validation.date_invalid_form" }),
      }),
    death_date: yup
      .string()
      .ensure()
      .test({
        name: "date",
        test: (value) => {
          if (!value) return true;
          if (!value.match(/^\d{4}-\d{2}-\d{2}$/)) return false;
          if (Number.isNaN(Date.parse(value))) return false;
          return true;
        },
        message: intl.formatMessage({ id: "validation.date_invalid_form" }),
      }),
    country: yup.string().ensure(),
    ethnicity: yup.string().ensure(),
    hair_color: yup.string().ensure(),
    eye_color: yup.string().ensure(),
    height_cm: yup.number().nullable().defined().default(null),
    weight: yup.number().nullable().defined().default(null),
    measurements: yup.string().ensure(),
    fake_tits: yup.string().ensure(),
    penis_length: yup.number().nullable().defined().default(null),
    circumcised: yup.string<GQL.CircumisedEnum | "">().ensure(),
    tattoos: yup.string().ensure(),
    piercings: yup.string().ensure(),
    career_length: yup.string().ensure(),
    url: yup.string().ensure(),
    twitter: yup.string().ensure(),
    instagram: yup.string().ensure(),
    details: yup.string().ensure(),
    tag_ids: yup.array(yup.string().required()).defined(),
    ignore_auto_tag: yup.boolean().defined(),
    stash_ids: yup.mixed<GQL.StashIdInput[]>().defined(),
    image: yup.string().nullable().optional(),
  });

  const initialValues = {
    name: performer.name ?? "",
    disambiguation: performer.disambiguation ?? "",
    alias_list: performer.alias_list ?? [],
    gender: (performer.gender as GQL.GenderEnum) ?? "",
    birthdate: performer.birthdate ?? "",
    death_date: performer.death_date ?? "",
    country: performer.country ?? "",
    ethnicity: performer.ethnicity ?? "",
    hair_color: performer.hair_color ?? "",
    eye_color: performer.eye_color ?? "",
    height_cm: performer.height_cm ?? null,
    weight: performer.weight ?? null,
    measurements: performer.measurements ?? "",
    fake_tits: performer.fake_tits ?? "",
    penis_length: performer.penis_length ?? null,
    circumcised: (performer.circumcised as GQL.CircumisedEnum) ?? "",
    tattoos: performer.tattoos ?? "",
    piercings: performer.piercings ?? "",
    career_length: performer.career_length ?? "",
    url: performer.url ?? "",
    twitter: performer.twitter ?? "",
    instagram: performer.instagram ?? "",
    details: performer.details ?? "",
    tag_ids: (performer.tags ?? []).map((t) => t.id),
    ignore_auto_tag: performer.ignore_auto_tag ?? false,
    stash_ids: getStashIDs(performer.stash_ids),
  };

  type InputValues = yup.InferType<typeof schema>;

  const formik = useFormik<InputValues>({
    initialValues,
    enableReinitialize: true,
    validationSchema: schema,
    onSubmit: (values) => onSave(values),
  });

  function translateScrapedGender(scrapedGender?: string) {
    if (!scrapedGender) {
      return;
    }

    // try to translate from enum values first
    const upperGender = scrapedGender.toUpperCase();
    const asEnum = genderToString(upperGender);
    if (asEnum) {
      return stringToGender(asEnum);
    } else {
      // try to match against gender strings
      const caseInsensitive = true;
      return stringToGender(scrapedGender, caseInsensitive);
    }
  }

  function translateScrapedCircumcised(scrapedCircumcised?: string) {
    if (!scrapedCircumcised) {
      return;
    }

    const upperCircumcised = scrapedCircumcised.toUpperCase();
    const asEnum = circumcisedToString(upperCircumcised);
    if (asEnum) {
      return stringToCircumcised(asEnum);
    } else {
      const caseInsensitive = true;
      return stringToCircumcised(scrapedCircumcised, caseInsensitive);
    }
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
              <Icon className="fa-fw" icon={faPlus} />
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

  async function createNewTag(toCreate: GQL.ScrapedTag) {
    const tagInput: GQL.TagCreateInput = { name: toCreate.name ?? "" };
    try {
      const result = await createTag({
        variables: {
          input: tagInput,
        },
      });

      if (!result.data?.tagCreate) {
        Toast.error(new Error("Failed to create tag"));
        return;
      }

      // add the new tag to the new tags value
      const newTagIds = formik.values.tag_ids.concat([
        result.data.tagCreate.id,
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
    if (state.disambiguation) {
      formik.setFieldValue("disambiguation", state.disambiguation);
    }
    if (state.aliases) {
      formik.setFieldValue(
        "alias_list",
        state.aliases.split(",").map((a) => a.trim())
      );
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
      formik.setFieldValue("height_cm", parseInt(state.height, 10));
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
      const newGender = translateScrapedGender(state.gender);
      if (newGender) {
        formik.setFieldValue("gender", newGender);
      }
    }
    if (state.circumcised) {
      // circumcised is a string in the scraper data
      const newCircumcised = translateScrapedCircumcised(state.circumcised);
      if (newCircumcised) {
        formik.setFieldValue("circumcised", newCircumcised);
      }
    }
    if (state.tags) {
      // map tags to their ids and filter out those not found
      const newTagIds = state.tags.map((t) => t.stored_id).filter((t) => t);
      formik.setFieldValue("tag_ids", newTagIds);

      setNewTags(state.tags.filter((t) => !t.stored_id));
    }

    // image is a base64 string
    // #404: don't overwrite image if it has been modified by the user
    // overwrite if not new since it came from a dialog
    // overwrite if image is unset
    if (
      (!isNew || !formik.values.image) &&
      state.images &&
      state.images.length > 0
    ) {
      const imageStr = state.images[0];
      formik.setFieldValue("image", imageStr);
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
    if (state.penis_length) {
      formik.setFieldValue("penis_length", state.penis_length);
    }

    const remoteSiteID = state.remote_site_id;
    if (remoteSiteID && (scraper as IStashBox).endpoint) {
      const newIDs =
        formik.values.stash_ids?.filter(
          (s) => s.endpoint !== (scraper as IStashBox).endpoint
        ) ?? [];
      newIDs?.push({
        endpoint: (scraper as IStashBox).endpoint,
        stash_id: remoteSiteID,
      });
      formik.setFieldValue("stash_ids", newIDs);
    }
  }

  const encodingImage = ImageUtils.usePasteImage(onImageLoad);

  useEffect(() => {
    setImage(formik.values.image);
  }, [formik.values.image, setImage]);

  useEffect(() => {
    setEncodingImage(encodingImage);
  }, [setEncodingImage, encodingImage]);

  function onImageLoad(imageData: string | null) {
    formik.setFieldValue("image", imageData);
  }

  function onImageChange(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function valuesToInput(input: InputValues): GQL.PerformerCreateInput {
    return {
      ...input,
      gender: input.gender || null,
      height_cm: input.height_cm || null,
      weight: input.weight || null,
      penis_length: input.penis_length || null,
      circumcised: input.circumcised || null,
    };
  }

  async function onSave(input: InputValues) {
    setIsLoading(true);
    try {
      await onSubmit(valuesToInput(input));
      formik.resetForm();
    } catch (e) {
      Toast.error(e);
    }
    setIsLoading(false);
  }

  // set up hotkeys
  useEffect(() => {
    if (isVisible) {
      Mousetrap.bind("s s", () => {
        if (formik.dirty) {
          formik.submitForm();
        }
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
    const newQueryableScrapers = (
      Scrapers?.data?.listPerformerScrapers ?? []
    ).filter((s) =>
      s.performer?.supported_scrapes.includes(GQL.ScrapeType.Name)
    );

    setQueryableScrapers(newQueryableScrapers);
  }, [Scrapers]);

  if (isLoading) return <LoadingIndicator />;

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
    selectedPerformer: GQL.ScrapedPerformerDataFragment,
    selectedScraper: GQL.Scraper
  ) {
    setIsScraperModalOpen(false);
    try {
      if (!scraper) return;
      setIsLoading(true);

      const {
        __typename,
        images: _image,
        tags: _tags,
        ...ret
      } = selectedPerformer;

      const result = await queryScrapePerformer(selectedScraper.id, ret);
      if (!result?.data?.scrapeSinglePerformer?.length) return;

      // assume one result
      // if this is a new performer, just dump the data
      if (isNew) {
        updatePerformerEditStateFromScraper(
          result.data.scrapeSinglePerformer[0]
        );
        setScraper(undefined);
      } else {
        setScrapedPerformer(result.data.scrapeSinglePerformer[0]);
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

  async function onScrapeStashBox(performerResult: GQL.ScrapedPerformer) {
    setIsScraperModalOpen(false);

    const result: GQL.ScrapedPerformerDataFragment = {
      ...performerResult,
      images: performerResult.images ?? undefined,
      __typename: "ScrapedPerformer",
    };

    // if this is a new performer, just dump the data
    if (isNew) {
      updatePerformerEditStateFromScraper(result);
      setScraper(undefined);
    } else {
      setScrapedPerformer(result);
    }
  }

  function onScraperSelected(s: GQL.Scraper | IStashBox | undefined) {
    setScraper(s);
    setIsScraperModalOpen(true);
  }

  function renderScraperMenu() {
    if (!performer) {
      return;
    }
    const stashBoxes = stashConfig?.general.stashBoxes ?? [];

    const popover = (
      <Dropdown.Menu id="performer-scraper-popover">
        {stashBoxes.map((s, index) => (
          <Dropdown.Item
            as={Button}
            key={s.endpoint}
            className="minimal"
            onClick={() => onScraperSelected({ ...s, index })}
          >
            {stashboxDisplayName(s.name, index)}
          </Dropdown.Item>
        ))}
        {queryableScrapers
          ? queryableScrapers.map((s) => (
              <Dropdown.Item
                as={Button}
                key={s.name}
                className="minimal"
                onClick={() => onScraperSelected(s)}
              >
                {s.name}
              </Dropdown.Item>
            ))
          : ""}
        <Dropdown.Item
          as={Button}
          className="minimal"
          onClick={() => onReloadScrapers()}
        >
          <span className="fa-icon">
            <Icon icon={faSyncAlt} />
          </span>
          <span>
            <FormattedMessage id="actions.reload_scrapers" />
          </span>
        </Dropdown.Item>
      </Dropdown.Menu>
    );

    return (
      <Dropdown drop="up" className="d-inline-block">
        <Dropdown.Toggle variant="secondary" className="mr-2">
          <FormattedMessage id="actions.scrape_with" />
        </Dropdown.Toggle>
        {popover}
      </Dropdown>
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

    const currentPerformer = {
      ...valuesToInput(formik.values),
      image: formik.values.image ?? performer.image_path,
    };

    return (
      <PerformerScrapeDialog
        performer={currentPerformer}
        scraped={scrapedPerformer}
        scraper={scraper}
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
    setScraper(undefined);
  }

  function renderButtons(classNames: string) {
    return (
      <div className={cx("details-edit", "col-xl-9", classNames)}>
        {!isNew && onCancel ? (
          <Button className="mr-2" variant="primary" onClick={onCancel}>
            <FormattedMessage id="actions.cancel" />
          </Button>
        ) : null}
        {renderScraperMenu()}
        <ImageInput
          isEditing
          onImageChange={onImageChange}
          onImageURL={onImageLoad}
        />
        <div>
          <Button
            className="mr-2"
            variant="danger"
            onClick={() => formik.setFieldValue("image", null)}
          >
            <FormattedMessage id="actions.clear_image" />
          </Button>
        </div>
        <Button
          variant="success"
          disabled={(!isNew && !formik.dirty) || !isEqual(formik.errors, {})}
          onClick={() => formik.submitForm()}
        >
          <FormattedMessage id="actions.save" />
        </Button>
      </div>
    );
  }

  const renderScrapeModal = () => {
    if (!isScraperModalOpen) return;

    return scraper !== undefined && isScraper(scraper) ? (
      <PerformerScrapeModal
        scraper={scraper}
        onHide={() => setScraper(undefined)}
        onSelectPerformer={onScrapePerformer}
        name={formik.values.name || ""}
      />
    ) : scraper !== undefined && !isScraper(scraper) ? (
      <PerformerStashBoxModal
        instance={scraper}
        onHide={() => setScraper(undefined)}
        onSelectPerformer={onScrapeStashBox}
        name={formik.values.name || ""}
      />
    ) : undefined;
  };

  function renderTagsField() {
    return (
      <Form.Group controlId="tags" as={Row}>
        <Form.Label column sm={labelXS} xl={labelXL}>
          <FormattedMessage id="tags" defaultMessage="Tags" />
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
          {intl.formatMessage({ id: "stash_ids" })}
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
                    title={intl.formatMessage({ id: "actions.delete_stashid" })}
                    onClick={() => removeStashID(stashID)}
                  >
                    <Icon icon={faTrashAlt} />
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

  function renderField(
    field: string,
    props?: {
      messageID?: string;
      placeholder?: string;
      type?: string;
    }
  ) {
    const title = intl.formatMessage({ id: props?.messageID ?? field });

    return (
      <Form.Group controlId={field} as={Row}>
        <Form.Label column xs={labelXS} xl={labelXL}>
          {title}
        </Form.Label>
        <Col xs={fieldXS} xl={fieldXL}>
          <Form.Control
            type={props?.type ?? "text"}
            className="text-input"
            placeholder={props?.placeholder ?? title}
            {...formik.getFieldProps(field)}
            isInvalid={!!formik.getFieldMeta(field).error}
          />
          <Form.Control.Feedback type="invalid">
            {formik.getFieldMeta(field).error}
          </Form.Control.Feedback>
        </Col>
      </Form.Group>
    );
  }

  const aliasErrors = Array.isArray(formik.errors.alias_list)
    ? formik.errors.alias_list[0]
    : formik.errors.alias_list;
  const aliasErrorMsg = aliasErrors
    ? intl.formatMessage({ id: "validation.aliases_must_be_unique" })
    : undefined;
  const aliasErrorIdx = aliasErrors?.split(" ").map((e) => parseInt(e));

  return (
    <>
      {renderScrapeModal()}
      {maybeRenderScrapeDialog()}

      <Prompt
        when={formik.dirty}
        message={intl.formatMessage({ id: "dialogs.unsaved_changes" })}
      />
      {renderButtons("mb-3")}

      <Form noValidate onSubmit={formik.handleSubmit} id="performer-edit">
        <Form.Group controlId="name" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            <FormattedMessage id="name" />
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <Form.Control
              className="text-input"
              placeholder={intl.formatMessage({ id: "name" })}
              {...formik.getFieldProps("name")}
              isInvalid={!!formik.errors.name}
            />
            <Form.Control.Feedback type="invalid">
              {formik.errors.name}
            </Form.Control.Feedback>
          </Col>
        </Form.Group>

        <Form.Group controlId="disambiguation" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            <FormattedMessage id="disambiguation" />
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <Form.Control
              className="text-input"
              placeholder={intl.formatMessage({ id: "disambiguation" })}
              {...formik.getFieldProps("disambiguation")}
              isInvalid={!!formik.errors.disambiguation}
            />
            <Form.Control.Feedback type="invalid">
              {formik.errors.disambiguation}
            </Form.Control.Feedback>
          </Col>
        </Form.Group>

        <Form.Group controlId="aliases" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            <FormattedMessage id="aliases" />
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <StringListInput
              value={formik.values.alias_list ?? []}
              setValue={(value) => formik.setFieldValue("alias_list", value)}
              errors={aliasErrorMsg}
              errorIdx={aliasErrorIdx}
            />
          </Col>
        </Form.Group>

        <Form.Group as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            <FormattedMessage id="gender" />
          </Form.Label>
          <Col xs="auto">
            <Form.Control
              as="select"
              className="input-control"
              {...formik.getFieldProps("gender")}
            >
              <option value="" key=""></option>
              {Array.from(stringGenderMap.entries()).map(([name, value]) => (
                <option value={value} key={value}>
                  {name}
                </option>
              ))}
            </Form.Control>
          </Col>
        </Form.Group>

        <Form.Group controlId="birthdate" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            <FormattedMessage id="birthdate" />
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <DateInput
              value={formik.values.birthdate}
              onValueChange={(value) =>
                formik.setFieldValue("birthdate", value)
              }
              error={formik.errors.birthdate}
            />
          </Col>
        </Form.Group>

        <Form.Group controlId="death_date" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            <FormattedMessage id="death_date" />
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <DateInput
              value={formik.values.death_date}
              onValueChange={(value) =>
                formik.setFieldValue("death_date", value)
              }
              error={formik.errors.death_date}
            />
          </Col>
        </Form.Group>

        <Form.Group as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            <FormattedMessage id="country" />
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <CountrySelect
              value={formik.getFieldProps("country").value}
              onChange={(value) => formik.setFieldValue("country", value)}
            />
          </Col>
        </Form.Group>

        {renderField("ethnicity")}
        {renderField("hair_color")}
        {renderField("eye_color")}
        {renderField("height_cm", {
          type: "number",
        })}
        {renderField("weight", {
          type: "number",
          messageID: "weight_kg",
        })}
        {renderField("penis_length", {
          type: "number",
          messageID: "penis_length_cm",
        })}

        <Form.Group as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            <FormattedMessage id="circumcised" />
          </Form.Label>
          <Col xs="auto">
            <Form.Control
              as="select"
              className="input-control"
              {...formik.getFieldProps("circumcised")}
            >
              <option value="" key=""></option>
              {Array.from(stringCircumMap.entries()).map(([name, value]) => (
                <option value={value} key={value}>
                  {name}
                </option>
              ))}
            </Form.Control>
          </Col>
        </Form.Group>

        {renderField("measurements")}
        {renderField("fake_tits")}

        <Form.Group controlId="tattoos" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            <FormattedMessage id="tattoos" />
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <Form.Control
              as="textarea"
              className="text-input"
              placeholder={intl.formatMessage({ id: "tattoos" })}
              {...formik.getFieldProps("tattoos")}
            />
          </Col>
        </Form.Group>

        <Form.Group controlId="piercings" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            <FormattedMessage id="piercings" />
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <Form.Control
              as="textarea"
              className="text-input"
              placeholder={intl.formatMessage({ id: "piercings" })}
              {...formik.getFieldProps("piercings")}
            />
          </Col>
        </Form.Group>

        {renderField("career_length")}

        <Form.Group controlId="url" as={Row}>
          <Form.Label column xs={labelXS} xl={labelXL}>
            <FormattedMessage id="url" />
          </Form.Label>
          <Col xs={fieldXS} xl={fieldXL}>
            <URLField
              {...formik.getFieldProps("url")}
              onScrapeClick={onScrapePerformerURL}
              urlScrapable={urlScrapable}
            />
          </Col>
        </Form.Group>

        {renderField("twitter")}
        {renderField("instagram")}
        <Form.Group controlId="details" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            <FormattedMessage id="details" />
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <Form.Control
              as="textarea"
              className="text-input"
              placeholder={intl.formatMessage({ id: "details" })}
              {...formik.getFieldProps("details")}
            />
          </Col>
        </Form.Group>
        {renderTagsField()}

        {renderStashIDs()}

        <hr />

        <Form.Group controlId="ignore-auto-tag" as={Row}>
          <Form.Label column sm={labelXS} xl={labelXL}>
            <FormattedMessage id="ignore_auto_tag" />
          </Form.Label>
          <Col sm={fieldXS} xl={fieldXL}>
            <Form.Check
              {...formik.getFieldProps({
                name: "ignore_auto_tag",
                type: "checkbox",
              })}
            />
          </Col>
        </Form.Group>

        {renderButtons("mt-3")}
      </Form>
    </>
  );
};
