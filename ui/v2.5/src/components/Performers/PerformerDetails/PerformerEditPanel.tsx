import React, { useEffect, useState } from "react";
import { Button, Popover, OverlayTrigger, Form, Col, InputGroup, Row, Badge } from "react-bootstrap";
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
  queryScrapePerformerURL,
} from "src/core/StashService";
import {
  Icon,
  Modal,
  ImageInput,
  ScrapePerformerSuggest,
  LoadingIndicator,
  TagSelect,
  CollapseButton,
} from "src/components/Shared";
import {
  ImageUtils,
} from "src/utils";
import { useToast } from "src/hooks";
import { useHistory } from "react-router-dom";
import { Formik } from "formik";
import { PerformerScrapeDialog } from "./PerformerScrapeDialog";

interface IPerformerDetails {
  performer: Partial<GQL.PerformerDataFragment>;
  isNew?: boolean;
  isEditing?: boolean;
  isVisible: boolean;
  onDelete?: () => void;
  onImageChange?: (image?: string | null) => void;
  onImageEncoding?: (loading?: boolean) => void;
}

export const PerformerEditPanel: React.FC<IPerformerDetails> = ({
  performer,
  isNew,
  isEditing,
  isVisible,
  onDelete,
  onImageChange,
  onImageEncoding,
}) => {
  const Toast = useToast();
  const history = useHistory();

  // Editing state
  const [
    isDisplayingScraperDialog,
    setIsDisplayingScraperDialog,
  ] = useState<GQL.Scraper>();
  const [
    scrapePerformerDetails,
    setScrapePerformerDetails,
  ] = useState<GQL.ScrapedPerformerDataFragment>();
  const [isDeleteAlertOpen, setIsDeleteAlertOpen] = useState<boolean>(false);

  // Editing performer state
  const [image, setImage] = useState<string | null>();
  const [name, setName] = useState<string>(performer?.name ?? "");
  const [aliases, setAliases] = useState<string>(performer.aliases ?? "");
  const [birthdate, setBirthdate] = useState<string>(performer.birthdate ?? "");
  const [ethnicity, setEthnicity] = useState<string>(performer.ethnicity ?? "");
  const [country, setCountry] = useState<string>(performer.country ?? "");
  const [eyeColor, setEyeColor] = useState<string>(performer.eye_color ?? "");
  const [height, setHeight] = useState<string>(performer.height ?? "");
  const [measurements, setMeasurements] = useState<string>(
    performer.measurements ?? ""
  );
  const [fakeTits, setFakeTits] = useState<string>(performer.fake_tits ?? "");
  const [careerLength, setCareerLength] = useState<string>(
    performer.career_length ?? ""
  );
  const [tattoos, setTattoos] = useState<string>(performer.tattoos ?? "");
  const [piercings, setPiercings] = useState<string>(performer.piercings ?? "");
  // const [url, setUrl] = useState<string>(performer.url ?? "");
  const [twitter, setTwitter] = useState<string>(performer.twitter ?? "");
  const [instagram, setInstagram] = useState<string>(performer.instagram ?? "");
  const [gender, setGender] = useState<string | undefined>(
    genderToString(performer.gender ?? undefined)
  );

  // const [tagIds, setTagIds] = useState<string[]>(
  //   (performer.tags ?? []).map((t) => t.id)
  // );
  const [newTags, setNewTags] = useState<GQL.ScrapedSceneTag[]>();

  const [stashIDs /* , setStashIDs*/ ] = useState<GQL.StashIdInput[]>(
    performer.stash_ids ?? []
  );
  const favorite = performer.favorite ?? false;

  // Network state
  const [isLoading, setIsLoading] = useState(false);

  const [updatePerformer] = usePerformerUpdate();
  const [createPerformer] = usePerformerCreate();

  const Scrapers = useListPerformerScrapers();
  const [queryableScrapers, setQueryableScrapers] = useState<GQL.Scraper[]>([]);

  const [scrapedPerformer, setScrapedPerformer] = useState<
    GQL.ScrapedPerformer | undefined
  >();

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, isEditing);

  // const [createTag] = useTagCreate({ name: "" });

  const genderOptions = [""].concat(getGenderStrings());

  const schema = yup.object({
    name: yup.string().required(),
    aliases: yup.string().optional(),
    gender: yup.string().optional().oneOf(genderOptions),
    birthdate: yup.string().optional(),
    ethnicity: yup.string().optional(),
    eyeColor: yup.string().optional(),
    country: yup.string().optional(),
    height: yup.string().optional(),
    measurements: yup.string().optional(),
    fakeTits: yup.string().optional(),
    careerLength: yup.string().optional(),
    tattoos: yup.string().optional(),
    piercings: yup.string().optional(),
    url: yup.string().optional(),
    twitter: yup.string().optional(),
    instagram: yup.string().optional(),
    tagIds: yup.array(yup.string().required()).optional(),
  });

  const initialValues = {
    name: performer.name ?? undefined,
    aliases: performer.aliases ?? undefined,
    gender: genderToString(performer.gender ?? undefined),
    birthdate: performer.birthdate ?? undefined,
    ethnicity: performer.ethnicity ?? undefined,
    eyeColor: performer.eye_color ?? undefined,
    country: performer.country ?? undefined,
    height: performer.height ?? undefined,
    measurements: performer.measurements ?? undefined,
    fakeTits: performer.fake_tits ?? undefined,
    careerLength: performer.career_length ?? undefined,
    tattoos: performer.tattoos ?? undefined,
    piercings: performer.piercings ?? undefined,
    url: performer.url ?? undefined,
    twitter: performer.twitter ?? undefined,
    instagram: performer.instagram ?? undefined,
    tagIds: performer.tags?.map(t => t.id) ?? undefined,
  };

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
    // let tagInput: GQL.TagCreateInput = { name: "" };
    try {
      // tagInput = Object.assign(tagInput, toCreate);
      // const result = await createTag({
      //   variables: tagInput,
      // });

      // add the new tag to the new tags value
      // const newTagIds = tagIds.concat([result.data!.tagCreate!.id]);
      // setTagIds(newTagIds);

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
      setName(state.name);
    }

    if (state.aliases) {
      setAliases(state.aliases ?? undefined);
    }
    if (state.birthdate) {
      setBirthdate(state.birthdate ?? undefined);
    }
    if (state.ethnicity) {
      setEthnicity(state.ethnicity ?? undefined);
    }
    if (state.country) {
      setCountry(state.country ?? undefined);
    }
    if (state.eye_color) {
      setEyeColor(state.eye_color ?? undefined);
    }
    if (state.height) {
      setHeight(state.height ?? undefined);
    }
    if (state.measurements) {
      setMeasurements(state.measurements ?? undefined);
    }
    if (state.fake_tits) {
      setFakeTits(state.fake_tits ?? undefined);
    }
    if (state.career_length) {
      setCareerLength(state.career_length ?? undefined);
    }
    if (state.tattoos) {
      setTattoos(state.tattoos ?? undefined);
    }
    if (state.piercings) {
      setPiercings(state.piercings ?? undefined);
    }
    // if (state.url) {
    //   setUrl(state.url ?? undefined);
    // }
    if (state.twitter) {
      setTwitter(state.twitter ?? undefined);
    }
    if (state.instagram) {
      setInstagram(state.instagram ?? undefined);
    }
    if (state.gender) {
      // gender is a string in the scraper data
      setGender(translateScrapedGender(state.gender ?? undefined));
    }
    if (state.tags) {
      // const newTagIds = state.tags.map((t) => t.stored_id).filter((t) => t);
      // setTagIds(newTagIds as string[]);

      setNewTags(state.tags.filter((t) => !t.stored_id));
    }

    // image is a base64 string
    // #404: don't overwrite image if it has been modified by the user
    // overwrite if not new since it came from a dialog
    // otherwise follow existing behaviour
    if (
      (!isNew || image === undefined) &&
      (state as GQL.ScrapedPerformerDataFragment).image !== undefined
    ) {
      const imageStr = (state as GQL.ScrapedPerformerDataFragment).image;
      setImage(imageStr ?? undefined);
    }
  }

  function onImageLoad(imageData: string) {
    setImage(imageData);
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
        if (performerInput.image) {
          // Refetch image to bust browser cache
          await fetch(`/performer/${performer.id}/image`, { cache: "reload" });
        }
      } else {
        const result = await createPerformer({
          variables: performerInput as GQL.PerformerCreateInput,
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
    if (isEditing && isVisible) {
      Mousetrap.bind("s s", () => {
        onSave?.(getPerformerInput());
      });

      if (!isNew) {
        Mousetrap.bind("d d", () => {
          setIsDeleteAlertOpen(true);
        });
      }

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
      onImageChange(image);
    }
    return () => onImageChange?.();
  }, [image, onImageChange]);

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
      // url,
      twitter,
      instagram,
      image,
      gender: stringToGender(gender),
      // tag_ids: tagIds,
      stash_ids: stashIDs.map((s) => ({
        stash_id: s.stash_id,
        endpoint: s.endpoint,
      })),
    };

    if (!isNew) {
      (performerInput as GQL.PerformerUpdateInput).id = performer.id!;
    }
    return performerInput;
  }

  function onImageChangeHandler(event: React.FormEvent<HTMLInputElement>) {
    ImageUtils.onImageChange(event, onImageLoad);
  }

  function onDisplayScrapeDialog(scraper: GQL.Scraper) {
    setIsDisplayingScraperDialog(scraper);
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

  function getQueryScraperPerformerInput() {
    if (!scrapePerformerDetails) return {};

    // image is not supported
    const { __typename, image: _image, ...ret } = scrapePerformerDetails;
    return ret;
  }

  async function onScrapePerformer() {
    setIsDisplayingScraperDialog(undefined);
    try {
      if (!scrapePerformerDetails || !isDisplayingScraperDialog) return;
      setIsLoading(true);
      const result = await queryScrapePerformer(
        isDisplayingScraperDialog.id,
        getQueryScraperPerformerInput()
      );
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

  async function onScrapePerformerURL(url?: string) {
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
    if (!performer || !isEditing) {
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
                      onClick={() => onDisplayScrapeDialog(s)}
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
      <OverlayTrigger trigger="click" placement="top" overlay={popover}>
        <Button variant="secondary" className="mr-2">
          Scrape with...
        </Button>
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
            onSelectPerformer={(query) => setScrapePerformerDetails(query)}
          />
        </div>
      </Modal>
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

    const currentPerformer: Partial<GQL.PerformerDataFragment> = {
      name,
      aliases,
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
      // url,
      twitter,
      instagram,
      gender: stringToGender(gender),
      image_path: image ?? performer.image_path,
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

  function maybeRenderScrapeButton(url?: string) {
    return (
      <Button
        variant="secondary"
        disabled={!urlScrapable(url)}
        className="scrape-url-button text-input"
        onClick={() => onScrapePerformerURL()}
      >
        <Icon icon="file-upload" />
      </Button>
    );
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  function renderTagsField(tagIds: string[] | undefined, setFieldValue: (field: string, value: any, shouldValidate?: boolean | undefined) => void) {
    return (
      <Form.Row>
        <Form.Group as={Col} sm="12">
          <Form.Label>
            Tags
          </Form.Label>
          <TagSelect
            menuPortalTarget={document.body}
            isMulti
            onSelect={(items) => setFieldValue("tagIds", items.map((item) => item.id))}
            ids={tagIds}
          />
          {renderNewTags()}
        </Form.Group>
      </Form.Row>
    );
  }

  function maybeRenderButtons() {
    if (isEditing) {
      return (
        <Row>
          <Col className="mt-3 pl-0" xs={12}>
          <Button
            className="mr-2"
            variant="primary"
            onClick={() => onSave?.(getPerformerInput())}
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
            isEditing={!!isEditing}
            onImageChange={onImageChangeHandler}
          />
          {isEditing ? (
            <Button
              className="mx-2"
              variant="danger"
              onClick={() => setImage(null)}
            >
              Clear image
            </Button>
          ) : (
            ""
          )}
          </Col>
        </Row>
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

  // const removeStashID = (stashID: GQL.StashIdInput) => {
  //   setStashIDs(
  //     stashIDs.filter(
  //       (s) =>
  //         !(s.endpoint === stashID.endpoint && s.stash_id === stashID.stash_id)
  //     )
  //   );
  // };

  // function renderStashIDs() {
  //   if (!performer.stash_ids?.length) {
  //     return;
  //   }

  //   return (
  //     <tr>
  //       <td>StashIDs</td>
  //       <td>
  //         <ul className="pl-0">
  //           {stashIDs.map((stashID) => {
  //             const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
  //             const link = base ? (
  //               <a
  //                 href={`${base}performers/${stashID.stash_id}`}
  //                 target="_blank"
  //                 rel="noopener noreferrer"
  //               >
  //                 {stashID.stash_id}
  //               </a>
  //             ) : (
  //               stashID.stash_id
  //             );
  //             return (
  //               <li key={stashID.stash_id} className="row no-gutters">
  //                 {isEditing && (
  //                   <Button
  //                     variant="danger"
  //                     className="mr-2 py-0"
  //                     title="Delete StashID"
  //                     onClick={() => removeStashID(stashID)}
  //                   >
  //                     <Icon icon="trash-alt" />
  //                   </Button>
  //                 )}
  //                 {link}
  //               </li>
  //             );
  //           })}
  //         </ul>
  //       </td>
  //     </tr>
  //   );
  // }

  return (
    <>
      {renderDeleteAlert()}
      {renderScraperDialog()}
      {maybeRenderScrapeDialog()}

      <Row>
        <Col xl={8} lg={12}>
        <Formik
          validationSchema={schema}
          onSubmit={console.log}
          initialValues={initialValues}
        >
          {({
            handleSubmit,
            handleChange,
            values,
            setFieldValue,
            touched,
            errors,
          }) => (
            <Form noValidate onSubmit={handleSubmit} id="performer-edit">
              <Form.Row>
                <Form.Group as={Col} sm="4">
                  <Form.Label>Name</Form.Label>
                  <Form.Control
                    name="name"
                    className="text-input"
                    value={values.name}
                    placeholder="Name"
                    onChange={handleChange}
                    isValid={touched.name && !errors.name}
                  />
                </Form.Group>
                <Form.Group as={Col} sm="8">
                  <Form.Label>Aliases</Form.Label>
                  <Form.Control
                    name="aliases"
                    className="text-input"
                    value={values.aliases}
                    placeholder="Aliases"
                    onChange={handleChange}
                    isValid={touched.aliases && !errors.aliases}
                  />
                </Form.Group>
              </Form.Row>

              <Form.Row>
                <Form.Group as={Col} md={8}>
                  <Form.Label>Gender</Form.Label>
                  <Form.Control
                    name="gender"
                    as="select"
                    className="input-control"
                    value={values.gender}
                    onChange={handleChange}
                    isValid={touched.gender && !errors.gender}
                  >
                    {genderOptions.map((opt) => (
                      <option value={opt} key={opt}>
                        {opt}
                      </option>
                    ))}
                  </Form.Control>
                </Form.Group>
              </Form.Row>

              <Form.Row>
                <Form.Group as={Col} sm="4">
                  <Form.Label>Birthdate</Form.Label>
                  <Form.Control
                    name="birthdate"
                    className="text-input"
                    value={values.birthdate}
                    placeholder="Birthdate"
                    onChange={handleChange}
                    isValid={touched.birthdate && !errors.birthdate}
                  />
                </Form.Group>
              </Form.Row>

              <Form.Row>
                <Form.Group as={Col} sm="4">
                  <Form.Label>Country</Form.Label>
                  <Form.Control
                    name="country"
                    className="text-input"
                    value={values.country}
                    placeholder="Country"
                    onChange={handleChange}
                    isValid={touched.country && !errors.country}
                  />
                </Form.Group>
                <Form.Group as={Col} sm="4">
                  <Form.Label>Ethnicity</Form.Label>
                  <Form.Control
                    name="ethnicity"
                    className="text-input"
                    value={values.ethnicity}
                    placeholder="Ethnicity"
                    onChange={handleChange}
                    isValid={touched.ethnicity && !errors.ethnicity}
                  />
                </Form.Group>
              </Form.Row>

              <Form.Row>
                <Form.Group as={Col} sm="4">
                  <Form.Label>Eye Color</Form.Label>
                  <Form.Control
                    name="eyeColor"
                    className="text-input"
                    value={values.eyeColor}
                    placeholder="Eye Color"
                    onChange={handleChange}
                    isValid={touched.eyeColor && !errors.eyeColor}
                  />
                </Form.Group>
              </Form.Row>

              <Form.Row>
                <Form.Group as={Col} sm="4">
                  <Form.Label>Height (cm)</Form.Label>
                  <Form.Control
                    name="height"
                    className="text-input"
                    value={values.height}
                    placeholder="Height (cm)"
                    onChange={handleChange}
                    isValid={touched.height && !errors.height}
                  />
                </Form.Group>
                <Form.Group as={Col} sm="4">
                  <Form.Label>Measurements</Form.Label>
                  <Form.Control
                    name="measurements"
                    className="text-input"
                    value={values.measurements}
                    placeholder="Measurements"
                    onChange={handleChange}
                    isValid={touched.measurements && !errors.measurements}
                  />
                </Form.Group>
                <Form.Group as={Col} sm="4">
                  <Form.Label>Fake Tits</Form.Label>
                  <Form.Control
                    name="fakeTits"
                    className="text-input"
                    value={values.fakeTits}
                    placeholder="Fake Tits"
                    onChange={handleChange}
                    isValid={touched.fakeTits && !errors.fakeTits}
                  />
                </Form.Group>
              </Form.Row>

              <Form.Row>
                <Form.Group as={Col} lg="6">
                  <Form.Label>Tattoos</Form.Label>
                  <Form.Control
                    name="tattoos"
                    className="text-input"
                    value={values.tattoos}
                    placeholder="Tattoos"
                    onChange={handleChange}
                    isValid={touched.tattoos && !errors.tattoos}
                  />
                </Form.Group>
                <Form.Group as={Col} lg="6">
                  <Form.Label>Piercings</Form.Label>
                  <Form.Control
                    name="piercings"
                    className="text-input"
                    value={values.piercings}
                    placeholder="Piercings"
                    onChange={handleChange}
                    isValid={touched.piercings && !errors.piercings}
                  />
                </Form.Group>
              </Form.Row>
              
              <Form.Row>
                <Form.Group as={Col} sm="4">
                  <Form.Label>Career Length</Form.Label>
                  <Form.Control
                    name="careerLength"
                    className="text-input"
                    value={values.careerLength}
                    placeholder="Career Length"
                    onChange={handleChange}
                    isValid={touched.careerLength && !errors.careerLength}
                  />
                </Form.Group>
              </Form.Row>

              <Form.Row>
                <Form.Group as={Col} sm="12">
                  <Form.Label>
                    URL
                  </Form.Label>
                  <InputGroup>
                    <Form.Control
                      name="url"
                      className="text-input"
                      value={values.url}
                      placeholder="URL"
                      onChange={handleChange}
                      isValid={touched.url && !errors.url}
                    />
                    <InputGroup.Append>
                      {maybeRenderScrapeButton(values.url)}
                    </InputGroup.Append>
                  </InputGroup>
                </Form.Group>
              </Form.Row>

              <Form.Row>
                <Form.Group as={Col} lg="6">
                  <Form.Label>Twitter</Form.Label>
                  <Form.Control
                    name="Twitter"
                    className="text-input"
                    value={values.twitter}
                    placeholder="Twitter"
                    onChange={handleChange}
                    isValid={touched.twitter && !errors.twitter}
                  />
                </Form.Group>
                <Form.Group as={Col} lg="6">
                  <Form.Label>Instagram</Form.Label>
                  <Form.Control
                    name="instagram"
                    className="text-input"
                    value={values.instagram}
                    placeholder="Instagram"
                    onChange={handleChange}
                    isValid={touched.instagram && !errors.instagram}
                  />
                </Form.Group>
              </Form.Row>
              {renderTagsField(values.tagIds, setFieldValue)}
              {/* {renderStashIDs()} */}

              {maybeRenderButtons()}
            </Form>
          )}
        </Formik>
      </Col>
      </Row>
    </>
  );
};
