import { Form, Col, Row, Button } from "react-bootstrap";
import React, { useEffect, useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import {
  circumcisedToString,
  stringToCircumcised,
} from "src/utils/circumcised";
import * as FormUtils from "src/utils/form";
import { genderToString, stringToGender } from "src/utils/gender";
import ImageUtils from "src/utils/image";
import {
  mutatePerformerMerge,
  queryFindPerformersByID,
} from "src/core/StashService";
import { FormattedMessage, useIntl } from "react-intl";
import { useToast } from "src/hooks/Toast";
import { faExchangeAlt, faSignInAlt } from "@fortawesome/free-solid-svg-icons";
import { ScrapeDialog } from "../Shared/ScrapeDialog/ScrapeDialog";
import {
  ScrapedImageRow,
  ScrapedInputGroupRow,
  ScrapedStringListRow,
  ScrapedTextAreaRow,
} from "../Shared/ScrapeDialog/ScrapeDialogRow";
import { ModalComponent } from "../Shared/Modal";
import { sortStoredIdObjects } from "src/utils/data";
import {
  ObjectListScrapeResult,
  ScrapeResult,
  ZeroableScrapeResult,
  hasScrapedValues,
} from "../Shared/ScrapeDialog/scrapeResult";
import { ScrapedTagsRow } from "../Shared/ScrapeDialog/ScrapedObjectsRow";
import {
  renderScrapedGenderRow,
  renderScrapedCircumcisedRow,
} from "./PerformerDetails/PerformerScrapeDialog";
import { PerformerSelect } from "./PerformerSelect";
import { uniq } from "lodash-es";

/* eslint-disable-next-line @typescript-eslint/no-explicit-any */
type CustomFieldScrapeResults = Map<string, ZeroableScrapeResult<any>>;

// There are a bunch of similar functions in PerformerScrapeDialog, but since we don't support
// scraping custom fields, this one is only needed here. The `renderScraped` naming is kept the same
// for consistency.
function renderScrapedCustomFieldRows(
  results: CustomFieldScrapeResults,
  onChange: (newCustomFields: CustomFieldScrapeResults) => void
) {
  return (
    <>
      {Array.from(results.entries()).map(([field, result]) => {
        const fieldName = `custom_${field}`;
        return (
          <ScrapedInputGroupRow
            className="custom-field"
            title={field}
            field={fieldName}
            key={fieldName}
            result={result}
            onChange={(newResult) => {
              const newResults = new Map(results);
              newResults.set(field, newResult);
              onChange(newResults);
            }}
          />
        );
      })}
    </>
  );
}

type MergeOptions = {
  values: GQL.PerformerUpdateInput;
};

interface IPerformerMergeDetailsProps {
  sources: GQL.PerformerDataFragment[];
  dest: GQL.PerformerDataFragment;
  onClose: (options?: MergeOptions) => void;
}

const PerformerMergeDetails: React.FC<IPerformerMergeDetailsProps> = ({
  sources,
  dest,
  onClose,
}) => {
  const intl = useIntl();

  const [loading, setLoading] = useState(true);

  const [name, setName] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.name)
  );
  const [disambiguation, setDisambiguation] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.disambiguation)
  );
  const [aliases, setAliases] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(dest.alias_list)
  );
  const [birthdate, setBirthdate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.birthdate)
  );
  const [deathDate, setDeathDate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.death_date)
  );
  const [ethnicity, setEthnicity] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.ethnicity)
  );
  const [country, setCountry] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.country)
  );
  const [hairColor, setHairColor] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.hair_color)
  );
  const [eyeColor, setEyeColor] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.eye_color)
  );
  const [height, setHeight] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.height_cm?.toString())
  );
  const [weight, setWeight] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.weight?.toString())
  );
  const [penisLength, setPenisLength] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.penis_length?.toString())
  );
  const [measurements, setMeasurements] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.measurements)
  );
  const [fakeTits, setFakeTits] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.fake_tits)
  );
  const [careerLength, setCareerLength] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.career_length)
  );
  const [tattoos, setTattoos] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.tattoos)
  );
  const [piercings, setPiercings] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.piercings)
  );
  const [urls, setURLs] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(dest.urls)
  );
  const [gender, setGender] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(genderToString(dest.gender))
  );
  const [circumcised, setCircumcised] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(circumcisedToString(dest.circumcised))
  );
  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.details)
  );
  const [tags, setTags] = useState<ObjectListScrapeResult<GQL.ScrapedTag>>(
    new ObjectListScrapeResult<GQL.ScrapedTag>(
      sortStoredIdObjects(dest.tags.map(idToStoredID))
    )
  );

  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.image_path)
  );

  const [customFields, setCustomFields] = useState<CustomFieldScrapeResults>(
    new Map()
  );

  function idToStoredID(o: { id: string; name: string }) {
    return {
      stored_id: o.id,
      name: o.name,
    };
  }

  // calculate the values for everything
  // uses the first set value for single value fields, and combines all
  useEffect(() => {
    async function loadImages() {
      const src = sources.find((s) => s.image_path);
      if (!dest.image_path || !src) return;

      setLoading(true);

      const destData = await ImageUtils.imageToDataURL(dest.image_path);
      const srcData = await ImageUtils.imageToDataURL(src.image_path!);

      // keep destination image by default
      const useNewValue = false;
      setImage(new ScrapeResult(destData, srcData, useNewValue));

      setLoading(false);
    }

    setName(
      new ScrapeResult(dest.name, sources.find((s) => s.name)?.name, !dest.name)
    );
    setDisambiguation(
      new ScrapeResult(
        dest.disambiguation,
        sources.find((s) => s.disambiguation)?.disambiguation,
        !dest.disambiguation
      )
    );

    // default alias list should be the existing aliases, plus the names of all sources,
    // plus all source aliases, deduplicated
    const allAliases = uniq(
      dest.alias_list.concat(
        sources.map((s) => s.name),
        sources.flatMap((s) => s.alias_list)
      )
    );

    setAliases(
      new ScrapeResult(dest.alias_list, allAliases, !!allAliases.length)
    );
    setBirthdate(
      new ScrapeResult(
        dest.birthdate,
        sources.find((s) => s.birthdate)?.birthdate,
        !dest.birthdate
      )
    );
    setDeathDate(
      new ScrapeResult(
        dest.death_date,
        sources.find((s) => s.death_date)?.death_date,
        !dest.death_date
      )
    );
    setEthnicity(
      new ScrapeResult(
        dest.ethnicity,
        sources.find((s) => s.ethnicity)?.ethnicity,
        !dest.ethnicity
      )
    );
    setCountry(
      new ScrapeResult(
        dest.country,
        sources.find((s) => s.country)?.country,
        !dest.country
      )
    );
    setHairColor(
      new ScrapeResult(
        dest.hair_color,
        sources.find((s) => s.hair_color)?.hair_color,
        !dest.hair_color
      )
    );
    setEyeColor(
      new ScrapeResult(
        dest.eye_color,
        sources.find((s) => s.eye_color)?.eye_color,
        !dest.eye_color
      )
    );
    setHeight(
      new ScrapeResult(
        dest.height_cm?.toString(),
        sources.find((s) => s.height_cm)?.height_cm?.toString(),
        !dest.height_cm
      )
    );
    setWeight(
      new ScrapeResult(
        dest.weight?.toString(),
        sources.find((s) => s.weight)?.weight?.toString(),
        !dest.weight
      )
    );

    setPenisLength(
      new ScrapeResult(
        dest.penis_length?.toString(),
        sources.find((s) => s.penis_length)?.penis_length?.toString(),
        !dest.penis_length
      )
    );
    setMeasurements(
      new ScrapeResult(
        dest.measurements,
        sources.find((s) => s.measurements)?.measurements,
        !dest.measurements
      )
    );
    setFakeTits(
      new ScrapeResult(
        dest.fake_tits,
        sources.find((s) => s.fake_tits)?.fake_tits,
        !dest.fake_tits
      )
    );
    setCareerLength(
      new ScrapeResult(
        dest.career_length,
        sources.find((s) => s.career_length)?.career_length,
        !dest.career_length
      )
    );
    setTattoos(
      new ScrapeResult(
        dest.tattoos,
        sources.find((s) => s.tattoos)?.tattoos,
        !dest.tattoos
      )
    );
    setPiercings(
      new ScrapeResult(
        dest.piercings,
        sources.find((s) => s.piercings)?.piercings,
        !dest.piercings
      )
    );
    setURLs(
      new ScrapeResult(
        dest.urls,
        sources.find((s) => s.urls)?.urls,
        !dest.urls?.length
      )
    );
    setGender(
      new ScrapeResult(
        genderToString(dest.gender),
        sources.find((s) => s.gender)?.gender
          ? genderToString(sources.find((s) => s.gender)?.gender)
          : undefined,
        !dest.gender
      )
    );
    setCircumcised(
      new ScrapeResult(
        circumcisedToString(dest.circumcised),
        sources.find((s) => s.circumcised)?.circumcised
          ? circumcisedToString(sources.find((s) => s.circumcised)?.circumcised)
          : undefined,
        !dest.circumcised
      )
    );
    setDetails(
      new ScrapeResult(
        dest.details,
        sources.find((s) => s.details)?.details,
        !dest.details
      )
    );
    setImage(
      new ScrapeResult(
        dest.image_path,
        sources.find((s) => s.image_path)?.image_path,
        !dest.image_path
      )
    );

    const customFieldNames = new Set<string>(Object.keys(dest.custom_fields));

    for (const s of sources) {
      for (const n of Object.keys(s.custom_fields)) {
        customFieldNames.add(n);
      }
    }

    setCustomFields(
      new Map(
        Array.from(customFieldNames)
          .sort()
          .map((field) => {
            return [
              field,
              new ScrapeResult(
                dest.custom_fields?.[field],
                sources.find((s) => s.custom_fields?.[field])?.custom_fields?.[
                  field
                ],
                dest.custom_fields?.[field] === undefined
              ),
            ];
          })
      )
    );

    loadImages();
  }, [sources, dest]);

  const hasCustomFieldValues = useMemo(() => {
    return hasScrapedValues(Array.from(customFields.values()));
  }, [customFields]);

  // ensure this is updated if fields are changed
  const hasValues = useMemo(() => {
    return (
      hasCustomFieldValues ||
      hasScrapedValues([
        name,
        disambiguation,
        aliases,
        birthdate,
        deathDate,
        ethnicity,
        country,
        hairColor,
        eyeColor,
        height,
        weight,
        penisLength,
        measurements,
        fakeTits,
        careerLength,
        tattoos,
        piercings,
        urls,
        gender,
        circumcised,
        details,
        tags,
        image,
      ])
    );
  }, [
    name,
    disambiguation,
    aliases,
    birthdate,
    deathDate,
    ethnicity,
    country,
    hairColor,
    eyeColor,
    height,
    weight,
    penisLength,
    measurements,
    fakeTits,
    careerLength,
    tattoos,
    piercings,
    urls,
    gender,
    circumcised,
    details,
    tags,
    image,
    hasCustomFieldValues,
  ]);

  function renderScrapeRows() {
    if (loading) {
      return (
        <div>
          <LoadingIndicator />
        </div>
      );
    }

    if (!hasValues) {
      return (
        <div>
          <FormattedMessage id="dialogs.merge.empty_results" />
        </div>
      );
    }

    return (
      <>
        <ScrapedInputGroupRow
          field="name"
          title={intl.formatMessage({ id: "name" })}
          result={name}
          onChange={(value) => setName(value)}
        />
        <ScrapedInputGroupRow
          field="disambiguation"
          title={intl.formatMessage({ id: "disambiguation" })}
          result={disambiguation}
          onChange={(value) => setDisambiguation(value)}
        />
        <ScrapedStringListRow
          field="aliases"
          title={intl.formatMessage({ id: "aliases" })}
          result={aliases}
          onChange={(value) => setAliases(value)}
        />
        <ScrapedInputGroupRow
          field="birthdate"
          title={intl.formatMessage({ id: "birthdate" })}
          result={birthdate}
          onChange={(value) => setBirthdate(value)}
        />
        <ScrapedInputGroupRow
          field="death_date"
          title={intl.formatMessage({ id: "death_date" })}
          result={deathDate}
          onChange={(value) => setDeathDate(value)}
        />
        <ScrapedInputGroupRow
          field="ethnicity"
          title={intl.formatMessage({ id: "ethnicity" })}
          result={ethnicity}
          onChange={(value) => setEthnicity(value)}
        />
        <ScrapedInputGroupRow
          field="country"
          title={intl.formatMessage({ id: "country" })}
          result={country}
          onChange={(value) => setCountry(value)}
        />
        <ScrapedInputGroupRow
          field="hair_color"
          title={intl.formatMessage({ id: "hair_color" })}
          result={hairColor}
          onChange={(value) => setHairColor(value)}
        />
        <ScrapedInputGroupRow
          field="eye_color"
          title={intl.formatMessage({ id: "eye_color" })}
          result={eyeColor}
          onChange={(value) => setEyeColor(value)}
        />
        <ScrapedInputGroupRow
          field="height"
          title={intl.formatMessage({ id: "height" })}
          result={height}
          onChange={(value) => setHeight(value)}
        />
        <ScrapedInputGroupRow
          field="weight"
          title={intl.formatMessage({ id: "weight" })}
          result={weight}
          onChange={(value) => setWeight(value)}
        />
        <ScrapedInputGroupRow
          field="penis_length"
          title={intl.formatMessage({ id: "penis_length" })}
          result={penisLength}
          onChange={(value) => setPenisLength(value)}
        />
        <ScrapedInputGroupRow
          field="measurements"
          title={intl.formatMessage({ id: "measurements" })}
          result={measurements}
          onChange={(value) => setMeasurements(value)}
        />
        <ScrapedInputGroupRow
          field="fake_tits"
          title={intl.formatMessage({ id: "fake_tits" })}
          result={fakeTits}
          onChange={(value) => setFakeTits(value)}
        />
        <ScrapedInputGroupRow
          field="career_length"
          title={intl.formatMessage({ id: "career_length" })}
          result={careerLength}
          onChange={(value) => setCareerLength(value)}
        />
        <ScrapedTextAreaRow
          field="tattoos"
          title={intl.formatMessage({ id: "tattoos" })}
          result={tattoos}
          onChange={(value) => setTattoos(value)}
        />
        <ScrapedTextAreaRow
          field="piercings"
          title={intl.formatMessage({ id: "piercings" })}
          result={piercings}
          onChange={(value) => setPiercings(value)}
        />
        <ScrapedStringListRow
          field="urls"
          title={intl.formatMessage({ id: "urls" })}
          result={urls}
          onChange={(value) => setURLs(value)}
        />
        {renderScrapedGenderRow(
          intl.formatMessage({ id: "gender" }),
          gender,
          (value) => setGender(value)
        )}
        {renderScrapedCircumcisedRow(
          intl.formatMessage({ id: "circumcised" }),
          circumcised,
          (value) => setCircumcised(value)
        )}
        <ScrapedTagsRow
          field="tags"
          title={intl.formatMessage({ id: "tags" })}
          result={tags}
          onChange={(value) => setTags(value)}
        />
        <ScrapedTextAreaRow
          field="details"
          title={intl.formatMessage({ id: "details" })}
          result={details}
          onChange={(value) => setDetails(value)}
        />
        <ScrapedImageRow
          field="image"
          title={intl.formatMessage({ id: "performer_image" })}
          className="performer-image"
          result={image}
          onChange={(value) => setImage(value)}
        />
        {hasCustomFieldValues &&
          renderScrapedCustomFieldRows(customFields, (newCustomFields) =>
            setCustomFields(newCustomFields)
          )}
      </>
    );
  }

  function createValues(): MergeOptions {
    // only set the cover image if it's different from the existing cover image
    const coverImage = image.useNewValue ? image.getNewValue() : undefined;

    return {
      values: {
        id: dest.id,
        name: name.getNewValue(),
        disambiguation: disambiguation.getNewValue(),
        alias_list: aliases
          .getNewValue()
          ?.map((s) => s.trim())
          .filter((s) => s.length > 0),
        birthdate: birthdate.getNewValue(),
        death_date: deathDate.getNewValue(),
        ethnicity: ethnicity.getNewValue(),
        country: country.getNewValue(),
        hair_color: hairColor.getNewValue(),
        eye_color: eyeColor.getNewValue(),
        height_cm: height.getNewValue()
          ? parseFloat(height.getNewValue()!)
          : undefined,
        weight: weight.getNewValue()
          ? parseFloat(weight.getNewValue()!)
          : undefined,
        penis_length: penisLength.getNewValue()
          ? parseFloat(penisLength.getNewValue()!)
          : undefined,
        measurements: measurements.getNewValue(),
        fake_tits: fakeTits.getNewValue(),
        career_length: careerLength.getNewValue(),
        tattoos: tattoos.getNewValue(),
        piercings: piercings.getNewValue(),
        urls: urls.getNewValue(),
        gender: stringToGender(gender.getNewValue()),
        circumcised: stringToCircumcised(circumcised.getNewValue()),
        tag_ids: tags.getNewValue()?.map((t) => t.stored_id!),
        details: details.getNewValue(),
        image: coverImage,
        custom_fields: {
          partial: Object.fromEntries(
            Array.from(customFields.entries()).flatMap(([field, v]) =>
              v.useNewValue ? [[field, v.getNewValue()]] : []
            )
          ),
        },
      },
    };
  }

  const dialogTitle = intl.formatMessage({
    id: "actions.merge",
  });

  const destinationLabel = !hasValues
    ? ""
    : intl.formatMessage({ id: "dialogs.merge.destination" });
  const sourceLabel = !hasValues
    ? ""
    : intl.formatMessage({ id: "dialogs.merge.source" });

  return (
    <ScrapeDialog
      className="performer-merge-dialog"
      title={dialogTitle}
      existingLabel={destinationLabel}
      scrapedLabel={sourceLabel}
      onClose={(apply) => {
        if (!apply) {
          onClose();
        } else {
          onClose(createValues());
        }
      }}
    >
      {renderScrapeRows()}
    </ScrapeDialog>
  );
};

interface IPerformerMergeModalProps {
  show: boolean;
  onClose: (mergedId?: string) => void;
  performers: GQL.SelectPerformerDataFragment[];
}

export const PerformerMergeModal: React.FC<IPerformerMergeModalProps> = ({
  show,
  onClose,
  performers,
}) => {
  const [sourcePerformers, setSourcePerformers] = useState<
    GQL.SelectPerformerDataFragment[]
  >([]);
  const [destPerformer, setDestPerformer] = useState<
    GQL.SelectPerformerDataFragment[]
  >([]);

  const [loadedSources, setLoadedSources] = useState<
    GQL.PerformerDataFragment[]
  >([]);
  const [loadedDest, setLoadedDest] = useState<GQL.PerformerDataFragment>();

  const [running, setRunning] = useState(false);
  const [secondStep, setSecondStep] = useState(false);

  const intl = useIntl();
  const Toast = useToast();

  const title = intl.formatMessage({
    id: "actions.merge",
  });

  useEffect(() => {
    if (performers.length > 0) {
      // set the first performer as the destination, others as source
      setDestPerformer([performers[0]]);

      if (performers.length > 1) {
        setSourcePerformers(performers.slice(1));
      }
    }
  }, [performers]);

  async function loadPerformers() {
    const performerIDs = sourcePerformers.map((s) => parseInt(s.id));
    performerIDs.push(parseInt(destPerformer[0].id));
    const query = await queryFindPerformersByID(performerIDs);
    const { performers: loadedPerformers } = query.data.findPerformers;

    setLoadedDest(loadedPerformers.find((s) => s.id === destPerformer[0].id));
    setLoadedSources(
      loadedPerformers.filter((s) => s.id !== destPerformer[0].id)
    );
    setSecondStep(true);
  }

  async function onMerge(options: MergeOptions) {
    const { values } = options;
    try {
      setRunning(true);
      const result = await mutatePerformerMerge(
        destPerformer[0].id,
        sourcePerformers.map((s) => s.id),
        values
      );
      if (result.data?.performerMerge) {
        Toast.success(intl.formatMessage({ id: "toast.merged_performers" }));
        onClose(destPerformer[0].id);
      }
      onClose();
    } catch (e) {
      Toast.error(e);
    } finally {
      setRunning(false);
    }
  }

  function canMerge() {
    return sourcePerformers.length > 0 && destPerformer.length !== 0;
  }

  function switchPerformers() {
    if (sourcePerformers.length && destPerformer.length) {
      const newDest = sourcePerformers[0];
      setSourcePerformers([...sourcePerformers.slice(1), destPerformer[0]]);
      setDestPerformer([newDest]);
    }
  }

  if (secondStep && destPerformer.length > 0) {
    return (
      <PerformerMergeDetails
        sources={loadedSources}
        dest={loadedDest!}
        onClose={(values) => {
          setSecondStep(false);
          if (values) {
            onMerge(values);
          } else {
            onClose();
          }
        }}
      />
    );
  }

  return (
    <ModalComponent
      dialogClassName="performer-merge-dialog"
      show={show}
      header={title}
      icon={faSignInAlt}
      accept={{
        text: intl.formatMessage({ id: "actions.next_action" }),
        onClick: () => loadPerformers(),
      }}
      disabled={!canMerge()}
      cancel={{
        variant: "secondary",
        onClick: () => onClose(),
      }}
      isRunning={running}
    >
      <div className="form-container row px-3">
        <div className="col-12 col-lg-6 col-xl-12">
          <Form.Group controlId="source" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({ id: "dialogs.merge.source" }),
              labelProps: {
                column: true,
                sm: 3,
                xl: 12,
              },
            })}
            <Col sm={9} xl={12}>
              <PerformerSelect
                isMulti
                onSelect={(items) => setSourcePerformers(items)}
                values={sourcePerformers}
                menuPortalTarget={document.body}
              />
            </Col>
          </Form.Group>
          <Form.Group
            controlId="switch"
            as={Row}
            className="justify-content-center"
          >
            <Button
              variant="secondary"
              onClick={() => switchPerformers()}
              disabled={!sourcePerformers.length || !destPerformer.length}
              title={intl.formatMessage({ id: "actions.swap" })}
            >
              <Icon className="fa-fw" icon={faExchangeAlt} />
            </Button>
          </Form.Group>
          <Form.Group controlId="destination" as={Row}>
            {FormUtils.renderLabel({
              title: intl.formatMessage({
                id: "dialogs.merge.destination",
              }),
              labelProps: {
                column: true,
                sm: 3,
                xl: 12,
              },
            })}
            <Col sm={9} xl={12}>
              <PerformerSelect
                onSelect={(items) => setDestPerformer(items)}
                values={destPerformer}
                menuPortalTarget={document.body}
              />
            </Col>
          </Form.Group>
        </div>
      </div>
    </ModalComponent>
  );
};
