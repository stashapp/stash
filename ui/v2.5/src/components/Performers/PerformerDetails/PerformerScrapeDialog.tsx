import React, { useState } from "react";
import * as GQL from "src/core/generated-graphql";
import {
  ScrapeDialog,
  ScrapeResult,
  ScrapedInputGroupRow,
  ScrapedImageRow,
  ScrapeDialogRow,
} from "src/components/Shared/ScrapeDialog";
import {
  getGenderStrings,
  genderToString,
  stringToGender,
  useTagCreate,
} from "src/core/StashService";
import { Form } from "react-bootstrap";
import { TagSelect } from "src/components/Shared";
import { useToast } from "src/hooks";
import _ from "lodash";

function renderScrapedGender(
  result: ScrapeResult<string>,
  isNew?: boolean,
  onChange?: (value: string) => void
) {
  const selectOptions = [""].concat(getGenderStrings());

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
  result: ScrapeResult<string>,
  onChange: (value: ScrapeResult<string>) => void
) {
  return (
    <ScrapeDialogRow
      title="Gender"
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

function renderScrapedTags(
  result: ScrapeResult<string[]>,
  isNew?: boolean,
  onChange?: (value: string[]) => void
) {
  const resultValue = isNew ? result.newValue : result.originalValue;
  const value = resultValue ?? [];

  return (
    <TagSelect
      isMulti
      className="form-control react-select"
      isDisabled={!isNew}
      onSelect={(items) => {
        if (onChange) {
          onChange(items.map((i) => i.id));
        }
      }}
      ids={value}
    />
  );
}

function renderScrapedTagsRow(
  result: ScrapeResult<string[]>,
  onChange: (value: ScrapeResult<string[]>) => void,
  newTags: GQL.ScrapedSceneTag[],
  onCreateNew?: (value: GQL.ScrapedSceneTag) => void
) {
  return (
    <ScrapeDialogRow
      title="Tags"
      result={result}
      renderOriginalField={() => renderScrapedTags(result)}
      renderNewField={() =>
        renderScrapedTags(result, true, (value) =>
          onChange(result.cloneWithValue(value))
        )
      }
      newValues={newTags}
      onChange={onChange}
      onCreateNew={onCreateNew}
    />
  );
}

interface IPerformerScrapeDialogProps {
  performer: Partial<GQL.PerformerUpdateInput>;
  scraped: GQL.ScrapedPerformer;

  onClose: (scrapedPerformer?: GQL.ScrapedPerformer) => void;
}

export const PerformerScrapeDialog: React.FC<IPerformerScrapeDialogProps> = (
  props: IPerformerScrapeDialogProps
) => {
  function translateScrapedGender(scrapedGender?: string | null) {
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

  const [name, setName] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.name, props.scraped.name)
  );
  const [aliases, setAliases] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.aliases, props.scraped.aliases)
  );
  const [birthdate, setBirthdate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.birthdate, props.scraped.birthdate)
  );
  const [ethnicity, setEthnicity] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.ethnicity, props.scraped.ethnicity)
  );
  const [country, setCountry] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.country, props.scraped.country)
  );
  const [eyeColor, setEyeColor] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.eye_color, props.scraped.eye_color)
  );
  const [height, setHeight] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.height, props.scraped.height)
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
  const [url, setURL] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.url, props.scraped.url)
  );
  const [twitter, setTwitter] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.twitter, props.scraped.twitter)
  );
  const [instagram, setInstagram] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.instagram, props.scraped.instagram)
  );
  const [gender, setGender] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(
      genderToString(props.performer.gender ?? undefined),
      translateScrapedGender(props.scraped.gender)
    )
  );

  const [createTag] = useTagCreate({ name: "" });
  const Toast = useToast();

  interface IHasStoredID {
    stored_id?: string | null;
  }

  function mapStoredIdObjects(
    scrapedObjects?: IHasStoredID[]
  ): string[] | undefined {
    if (!scrapedObjects) {
      return undefined;
    }
    const ret = scrapedObjects
      .map((p) => p.stored_id)
      .filter((p) => {
        return p !== undefined && p !== null;
      }) as string[];

    if (ret.length === 0) {
      return undefined;
    }

    // sort by id numerically
    ret.sort((a, b) => {
      return parseInt(a, 10) - parseInt(b, 10);
    });

    return ret;
  }

  function sortIdList(idList?: string[] | null) {
    if (!idList) {
      return;
    }

    const ret = _.clone(idList);
    // sort by id numerically
    ret.sort((a, b) => {
      return parseInt(a, 10) - parseInt(b, 10);
    });

    return ret;
  }

  const [tags, setTags] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(
      sortIdList(props.performer.tag_ids ?? undefined),
      mapStoredIdObjects(props.scraped.tags ?? undefined)
    )
  );

  const [newTags, setNewTags] = useState<GQL.ScrapedSceneTag[]>(
    props.scraped.tags?.filter((t) => !t.stored_id) ?? []
  );

  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(props.performer.image, props.scraped.image)
  );

  const allFields = [
    name,
    aliases,
    birthdate,
    ethnicity,
    country,
    eyeColor,
    height,
    measurements,
    fakeTits,
    careerLength,
    tattoos,
    piercings,
    url,
    twitter,
    instagram,
    gender,
    image,
    tags,
  ];
  // don't show the dialog if nothing was scraped
  if (allFields.every((r) => !r.scraped)) {
    props.onClose();
    return <></>;
  }

  async function createNewTag(toCreate: GQL.ScrapedSceneTag) {
    let tagInput: GQL.TagCreateInput = { name: "" };
    try {
      tagInput = Object.assign(tagInput, toCreate);
      const result = await createTag({
        variables: tagInput,
      });

      // add the new tag to the new tags value
      const tagClone = tags.cloneWithValue(tags.newValue);
      if (!tagClone.newValue) {
        tagClone.newValue = [];
      }
      tagClone.newValue.push(result.data!.tagCreate!.id);
      setTags(tagClone);

      // remove the tag from the list
      const newTagsClone = newTags.concat();
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

  function makeNewScrapedItem(): GQL.ScrapedPerformer {
    return {
      name: name.getNewValue(),
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
      url: url.getNewValue(),
      twitter: twitter.getNewValue(),
      instagram: instagram.getNewValue(),
      gender: gender.getNewValue(),
      tags: tags.getNewValue()?.map((m) => {
        return {
          stored_id: m,
          name: "",
        };
      }),
      image: image.getNewValue(),
    };
  }

  function renderScrapeRows() {
    return (
      <>
        <ScrapedInputGroupRow
          title="Name"
          result={name}
          onChange={(value) => setName(value)}
        />
        <ScrapedInputGroupRow
          title="Aliases"
          result={aliases}
          onChange={(value) => setAliases(value)}
        />
        {renderScrapedGenderRow(gender, (value) => setGender(value))}
        <ScrapedInputGroupRow
          title="Birthdate"
          result={birthdate}
          onChange={(value) => setBirthdate(value)}
        />
        <ScrapedInputGroupRow
          title="Ethnicity"
          result={ethnicity}
          onChange={(value) => setEthnicity(value)}
        />
        <ScrapedInputGroupRow
          title="Country"
          result={country}
          onChange={(value) => setCountry(value)}
        />
        <ScrapedInputGroupRow
          title="Eye Color"
          result={eyeColor}
          onChange={(value) => setEyeColor(value)}
        />
        <ScrapedInputGroupRow
          title="Height"
          result={height}
          onChange={(value) => setHeight(value)}
        />
        <ScrapedInputGroupRow
          title="Measurements"
          result={measurements}
          onChange={(value) => setMeasurements(value)}
        />
        <ScrapedInputGroupRow
          title="Fake Tits"
          result={fakeTits}
          onChange={(value) => setFakeTits(value)}
        />
        <ScrapedInputGroupRow
          title="Career Length"
          result={careerLength}
          onChange={(value) => setCareerLength(value)}
        />
        <ScrapedInputGroupRow
          title="Tattoos"
          result={tattoos}
          onChange={(value) => setTattoos(value)}
        />
        <ScrapedInputGroupRow
          title="Piercings"
          result={piercings}
          onChange={(value) => setPiercings(value)}
        />
        <ScrapedInputGroupRow
          title="URL"
          result={url}
          onChange={(value) => setURL(value)}
        />
        <ScrapedInputGroupRow
          title="Twitter"
          result={twitter}
          onChange={(value) => setTwitter(value)}
        />
        <ScrapedInputGroupRow
          title="Instagram"
          result={instagram}
          onChange={(value) => setInstagram(value)}
        />
        {renderScrapedTagsRow(
          tags,
          (value) => setTags(value),
          newTags,
          createNewTag
        )}
        <ScrapedImageRow
          title="Performer Image"
          className="performer-image"
          result={image}
          onChange={(value) => setImage(value)}
        />
      </>
    );
  }

  return (
    <ScrapeDialog
      title="Performer Scrape Results"
      renderScrapeRows={renderScrapeRows}
      onClose={(apply) => {
        props.onClose(apply ? makeNewScrapedItem() : undefined);
      }}
    />
  );
};
