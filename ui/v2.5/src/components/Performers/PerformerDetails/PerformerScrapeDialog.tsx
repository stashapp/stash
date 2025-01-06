import React, { useState } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  ScrapeDialog,
  ScrapedInputGroupRow,
  ScrapedImagesRow,
  ScrapeDialogRow,
  ScrapedTextAreaRow,
  ScrapedCountryRow,
  ScrapedStringListRow,
} from "src/components/Shared/ScrapeDialog/ScrapeDialog";
import { Form } from "react-bootstrap";
import {
  genderStrings,
  genderToString,
  stringToGender,
} from "src/utils/gender";
import {
  circumcisedStrings,
  circumcisedToString,
  stringToCircumcised,
} from "src/utils/circumcised";
import { IStashBox } from "./PerformerStashBoxModal";
import { ScrapeResult } from "src/components/Shared/ScrapeDialog/scrapeResult";
import { Tag } from "src/components/Tags/TagSelect";
import { uniq } from "lodash-es";
import { useScrapedTags } from "src/components/Shared/ScrapeDialog/scrapedTags";

function renderScrapedGender(
  result: ScrapeResult<string>,
  isNew?: boolean,
  onChange?: (value: string) => void
) {
  const selectOptions = [""].concat(genderStrings);

  return (
    <Form.Control
      as="select"
      className="input-control"
      disabled={!isNew}
      plaintext={!isNew}
      value={isNew ? result.newValue : result.originalValue}
      onChange={(e) => {
        if (isNew && onChange) {
          onChange(e.currentTarget.value);
        }
      }}
    >
      {selectOptions.map((opt) => (
        <option value={opt} key={opt}>
          {opt}
        </option>
      ))}
    </Form.Control>
  );
}

function renderScrapedGenderRow(
  title: string,
  result: ScrapeResult<string>,
  onChange: (value: ScrapeResult<string>) => void
) {
  return (
    <ScrapeDialogRow
      title={title}
      result={result}
      renderOriginalField={() => renderScrapedGender(result)}
      renderNewField={() =>
        renderScrapedGender(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
    />
  );
}

function renderScrapedCircumcised(
  result: ScrapeResult<string>,
  isNew?: boolean,
  onChange?: (value: string) => void
) {
  const selectOptions = [""].concat(circumcisedStrings);

  return (
    <Form.Control
      as="select"
      className="input-control"
      disabled={!isNew}
      plaintext={!isNew}
      value={isNew ? result.newValue : result.originalValue}
      onChange={(e) => {
        if (isNew && onChange) {
          onChange(e.currentTarget.value);
        }
      }}
    >
      {selectOptions.map((opt) => (
        <option value={opt} key={opt}>
          {opt}
        </option>
      ))}
    </Form.Control>
  );
}

function renderScrapedCircumcisedRow(
  title: string,
  result: ScrapeResult<string>,
  onChange: (value: ScrapeResult<string>) => void
) {
  return (
    <ScrapeDialogRow
      title={title}
      result={result}
      renderOriginalField={() => renderScrapedCircumcised(result)}
      renderNewField={() =>
        renderScrapedCircumcised(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      onChange={onChange}
    />
  );
}

interface IPerformerScrapeDialogProps {
  performer: Partial<GQL.PerformerUpdateInput>;
  performerTags: Tag[];
  scraped: GQL.ScrapedPerformer;
  scraper?: GQL.Scraper | IStashBox;

  onClose: (scrapedPerformer?: GQL.ScrapedPerformer) => void;
}

export const PerformerScrapeDialog: React.FC<IPerformerScrapeDialogProps> = (
  props: IPerformerScrapeDialogProps
) => {
  const intl = useIntl();

  const endpoint = (props.scraper as IStashBox)?.endpoint ?? undefined;

  function getCurrentRemoteSiteID() {
    if (!endpoint) {
      return;
    }

    return props.performer.stash_ids?.find((s) => s.endpoint === endpoint)
      ?.stash_id;
  }

  function translateScrapedGender(scrapedGender?: string | null) {
    if (!scrapedGender) {
      return;
    }

    let retEnum: GQL.GenderEnum | undefined;

    // try to translate from enum values first
    const upperGender = scrapedGender.toUpperCase();
    const asEnum = genderToString(upperGender);
    if (asEnum) {
      retEnum = stringToGender(asEnum);
    } else {
      // try to match against gender strings
      const caseInsensitive = true;
      retEnum = stringToGender(scrapedGender, caseInsensitive);
    }

    return genderToString(retEnum);
  }

  function translateScrapedCircumcised(scrapedCircumcised?: string | null) {
    if (!scrapedCircumcised) {
      return;
    }

    let retEnum: GQL.CircumisedEnum | undefined;

    // try to translate from enum values first
    const upperCircumcised = scrapedCircumcised.toUpperCase();
    const asEnum = circumcisedToString(upperCircumcised);
    if (asEnum) {
      retEnum = stringToCircumcised(asEnum);
    } else {
      // try to match against circumcised strings
      const caseInsensitive = true;
      retEnum = stringToCircumcised(scrapedCircumcised, caseInsensitive);
    }

    return circumcisedToString(retEnum);
  }

  const [name, setName] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.name, props.scraped.name)
  );
  const [disambiguation, setDisambiguation] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.performer.disambiguation,
      props.scraped.disambiguation
    )
  );
  const [aliases, setAliases] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.performer.alias_list?.join(", "),
      props.scraped.aliases
    )
  );
  const [birthdate, setBirthdate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.birthdate, props.scraped.birthdate)
  );
  const [deathDate, setDeathDate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.performer.death_date,
      props.scraped.death_date
    )
  );
  const [ethnicity, setEthnicity] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.ethnicity, props.scraped.ethnicity)
  );
  const [country, setCountry] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.country, props.scraped.country)
  );
  const [hairColor, setHairColor] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.performer.hair_color,
      props.scraped.hair_color
    )
  );
  const [eyeColor, setEyeColor] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.eye_color, props.scraped.eye_color)
  );
  const [height, setHeight] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.performer.height_cm?.toString(),
      props.scraped.height
    )
  );
  const [weight, setWeight] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.performer.weight?.toString(),
      props.scraped.weight
    )
  );
  const [penisLength, setPenisLength] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.performer.penis_length?.toString(),
      props.scraped.penis_length
    )
  );
  const [measurements, setMeasurements] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.performer.measurements,
      props.scraped.measurements
    )
  );
  const [fakeTits, setFakeTits] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.fake_tits, props.scraped.fake_tits)
  );
  const [careerLength, setCareerLength] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.performer.career_length,
      props.scraped.career_length
    )
  );
  const [tattoos, setTattoos] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.tattoos, props.scraped.tattoos)
  );
  const [piercings, setPiercings] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.piercings, props.scraped.piercings)
  );
  const [urls, setURLs] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      props.performer.urls,
      props.scraped.urls
        ? uniq((props.performer.urls ?? []).concat(props.scraped.urls ?? []))
        : undefined
    )
  );
  const [gender, setGender] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      genderToString(props.performer.gender),
      translateScrapedGender(props.scraped.gender)
    )
  );
  const [circumcised, setCircumcised] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      circumcisedToString(props.performer.circumcised),
      translateScrapedCircumcised(props.scraped.circumcised)
    )
  );
  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.details, props.scraped.details)
  );
  const [remoteSiteID, setRemoteSiteID] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      getCurrentRemoteSiteID(),
      props.scraped.remote_site_id
    )
  );

  const { tags, newTags, scrapedTagsRow } = useScrapedTags(
    props.performerTags,
    props.scraped.tags
  );

  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      props.performer.image,
      props.scraped.images && props.scraped.images.length > 0
        ? props.scraped.images[0]
        : undefined
    )
  );

  const images =
    props.scraped.images && props.scraped.images.length > 0
      ? props.scraped.images
      : [];

  const allFields = [
    name,
    disambiguation,
    aliases,
    birthdate,
    ethnicity,
    country,
    eyeColor,
    height,
    measurements,
    fakeTits,
    penisLength,
    circumcised,
    careerLength,
    tattoos,
    piercings,
    urls,
    gender,
    image,
    tags,
    details,
    deathDate,
    hairColor,
    weight,
    remoteSiteID,
  ];
  // don't show the dialog if nothing was scraped
  if (allFields.every((r) => !r.scraped) && newTags.length === 0) {
    props.onClose();
    return <></>;
  }

  function makeNewScrapedItem(): GQL.ScrapedPerformer {
    const newImage = image.getNewValue();
    return {
      name: name.getNewValue() ?? "",
      disambiguation: disambiguation.getNewValue(),
      aliases: aliases.getNewValue(),
      birthdate: birthdate.getNewValue(),
      ethnicity: ethnicity.getNewValue(),
      country: country.getNewValue(),
      eye_color: eyeColor.getNewValue(),
      height: height.getNewValue(),
      measurements: measurements.getNewValue(),
      fake_tits: fakeTits.getNewValue(),
      career_length: careerLength.getNewValue(),
      tattoos: tattoos.getNewValue(),
      piercings: piercings.getNewValue(),
      urls: urls.getNewValue(),
      gender: gender.getNewValue(),
      tags: tags.getNewValue(),
      images: newImage ? [newImage] : undefined,
      details: details.getNewValue(),
      death_date: deathDate.getNewValue(),
      hair_color: hairColor.getNewValue(),
      weight: weight.getNewValue(),
      penis_length: penisLength.getNewValue(),
      circumcised: circumcised.getNewValue(),
      remote_site_id: remoteSiteID.getNewValue(),
    };
  }

  function renderScrapeRows() {
    return (
      <>
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "name" })}
          result={name}
          onChange={(value) => setName(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "disambiguation" })}
          result={disambiguation}
          onChange={(value) => setDisambiguation(value)}
        />
        <ScrapedTextAreaRow
          title={intl.formatMessage({ id: "aliases" })}
          result={aliases}
          onChange={(value) => setAliases(value)}
        />
        {renderScrapedGenderRow(
          intl.formatMessage({ id: "gender" }),
          gender,
          (value) => setGender(value)
        )}
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "birthdate" })}
          result={birthdate}
          onChange={(value) => setBirthdate(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "death_date" })}
          result={deathDate}
          onChange={(value) => setDeathDate(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "ethnicity" })}
          result={ethnicity}
          onChange={(value) => setEthnicity(value)}
        />
        <ScrapedCountryRow
          title={intl.formatMessage({ id: "country" })}
          result={country}
          onChange={(value) => setCountry(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "hair_color" })}
          result={hairColor}
          onChange={(value) => setHairColor(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "eye_color" })}
          result={eyeColor}
          onChange={(value) => setEyeColor(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "weight" })}
          result={weight}
          onChange={(value) => setWeight(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "height" })}
          result={height}
          onChange={(value) => setHeight(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "penis_length" })}
          result={penisLength}
          onChange={(value) => setPenisLength(value)}
        />
        {renderScrapedCircumcisedRow(
          intl.formatMessage({ id: "circumcised" }),
          circumcised,
          (value) => setCircumcised(value)
        )}
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "measurements" })}
          result={measurements}
          onChange={(value) => setMeasurements(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "fake_tits" })}
          result={fakeTits}
          onChange={(value) => setFakeTits(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "career_length" })}
          result={careerLength}
          onChange={(value) => setCareerLength(value)}
        />
        <ScrapedTextAreaRow
          title={intl.formatMessage({ id: "tattoos" })}
          result={tattoos}
          onChange={(value) => setTattoos(value)}
        />
        <ScrapedTextAreaRow
          title={intl.formatMessage({ id: "piercings" })}
          result={piercings}
          onChange={(value) => setPiercings(value)}
        />
        <ScrapedStringListRow
          title={intl.formatMessage({ id: "urls" })}
          result={urls}
          onChange={(value) => setURLs(value)}
        />
        <ScrapedTextAreaRow
          title={intl.formatMessage({ id: "details" })}
          result={details}
          onChange={(value) => setDetails(value)}
        />
        {scrapedTagsRow}
        <ScrapedImagesRow
          title={intl.formatMessage({ id: "performer_image" })}
          className="performer-image"
          result={image}
          images={images}
          onChange={(value) => setImage(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "stash_id" })}
          result={remoteSiteID}
          locked
          onChange={(value) => setRemoteSiteID(value)}
        />
      </>
    );
  }

  return (
    <ScrapeDialog
      title={intl.formatMessage(
        { id: "dialogs.scrape_entity_title" },
        { entity_type: intl.formatMessage({ id: "performer" }) }
      )}
      renderScrapeRows={renderScrapeRows}
      onClose={(apply) => {
        props.onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    />
  );
};
