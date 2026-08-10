import { Form, Col, Row, Button, FormControl } from "react-bootstrap";
import React, { useEffect, useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { GallerySelect } from "../Shared/Select";
import * as FormUtils from "src/utils/form";
import ImageUtils from "src/utils/image";
import {
  mutateAudioMerge,
  queryFindFullAudiosByID,
} from "src/core/StashService";
import { FormattedMessage, useIntl } from "react-intl";
import { useToast } from "src/hooks/Toast";
import { faExchangeAlt, faSignInAlt } from "@fortawesome/free-solid-svg-icons";
import {
  ScrapeDialogRow,
  ScrapedCustomFieldRows,
  ScrapedImageRow,
  ScrapedInputGroupRow,
  ScrapedStringListRow,
  ScrapedTextAreaRow,
} from "../Shared/ScrapeDialog/ScrapeDialogRow";
import { ScrapeDialog } from "../Shared/ScrapeDialog/ScrapeDialog";
import { clone, uniq } from "lodash-es";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { ModalComponent } from "../Shared/Modal";
import {
  idToStoredID,
  sortStoredIdObjects,
  uniqIDStoredIDs,
} from "src/utils/data";
import {
  CustomFieldScrapeResults,
  ObjectListScrapeResult,
  ScrapeResult,
  ZeroableScrapeResult,
  hasScrapedValues,
} from "../Shared/ScrapeDialog/scrapeResult";
import {
  ScrapedGroupsRow,
  ScrapedPerformersRow,
  ScrapedStudioRow,
  ScrapedTagsRow,
} from "../Shared/ScrapeDialog/ScrapedObjectsRow";
import { Audio, AudioSelect } from "./AudioSelect";

type MergeOptions = {
  values: GQL.AudioUpdateInput;
  includeViewHistory: boolean;
  includeOHistory: boolean;
};

interface IAudioMergeDetailsProps {
  sources: GQL.AudioDataFragment[];
  dest: GQL.AudioDataFragment;
  onClose: (options?: MergeOptions) => void;
}

function groupToStoredID(o: { group: { id: string; name: string } }) {
  return {
    stored_id: o.group.id,
    name: o.group.name,
  };
}

const AudioMergeDetails: React.FC<IAudioMergeDetailsProps> = ({
  sources,
  dest,
  onClose,
}) => {
  const intl = useIntl();

  const [loading, setLoading] = useState(true);

  const [title, setTitle] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.title)
  );
  const [code, setCode] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.code)
  );
  const [url, setURL] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(dest.urls)
  );
  const [date, setDate] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.date)
  );

  const [rating, setRating] = useState(
    new ZeroableScrapeResult<number>(dest.rating100)
  );
  // audioUpdate cannot set these directly - they are only used to decide
  // whether the source histories are merged into the destination
  const [oCounter, setOCounter] = useState(
    new ScrapeResult<number>(dest.o_counter)
  );
  const [playCount, setPlayCount] = useState(
    new ScrapeResult<number>(dest.play_count)
  );

  const [studio, setStudio] = useState<ScrapeResult<GQL.ScrapedStudio>>(
    new ScrapeResult<GQL.ScrapedStudio>(
      dest.studio ? idToStoredID(dest.studio) : undefined
    )
  );

  function sortIdList(idList?: string[] | null) {
    if (!idList) {
      return;
    }

    const ret = clone(idList);
    // sort by id numerically
    ret.sort((a, b) => {
      return parseInt(a, 10) - parseInt(b, 10);
    });

    return ret;
  }

  const [performers, setPerformers] = useState<
    ObjectListScrapeResult<GQL.ScrapedPerformer>
  >(
    new ObjectListScrapeResult<GQL.ScrapedPerformer>(
      sortStoredIdObjects(dest.performers.map(idToStoredID))
    )
  );

  const [groups, setGroups] = useState<
    ObjectListScrapeResult<GQL.ScrapedGroup>
  >(
    new ObjectListScrapeResult<GQL.ScrapedGroup>(
      sortStoredIdObjects(dest.groups.map(groupToStoredID))
    )
  );

  const [tags, setTags] = useState<ObjectListScrapeResult<GQL.ScrapedTag>>(
    new ObjectListScrapeResult<GQL.ScrapedTag>(
      sortStoredIdObjects(dest.tags.map(idToStoredID))
    )
  );

  const [details, setDetails] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.details)
  );

  const [galleries, setGalleries] = useState<ScrapeResult<string[]>>(
    new ScrapeResult<string[]>(sortIdList(dest.galleries.map((p) => p.id)))
  );

  const [organized, setOrganized] = useState(
    new ZeroableScrapeResult<boolean>(dest.organized)
  );

  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.paths.screenshot)
  );

  const [customFields, setCustomFields] = useState<CustomFieldScrapeResults>(
    new Map()
  );

  // calculate the values for everything
  // uses the first set value for single value fields, and combines all
  useEffect(() => {
    async function loadImages() {
      const src = sources.find((s) => s.paths.screenshot);
      if (!dest.paths.screenshot || !src) return;

      setLoading(true);

      const destData = await ImageUtils.imageToDataURL(dest.paths.screenshot);
      const srcData = await ImageUtils.imageToDataURL(src.paths.screenshot!);

      // keep destination image by default
      const useNewValue = false;
      setImage(new ScrapeResult(destData, srcData, useNewValue));

      setLoading(false);
    }

    const all = sources.concat(dest);

    setTitle(
      new ScrapeResult(
        dest.title,
        sources.find((s) => s.title)?.title,
        !dest.title
      )
    );
    setCode(
      new ScrapeResult(dest.code, sources.find((s) => s.code)?.code, !dest.code)
    );
    setURL(new ScrapeResult(dest.urls, uniq(all.flatMap((s) => s.urls))));
    setDate(
      new ScrapeResult(dest.date, sources.find((s) => s.date)?.date, !dest.date)
    );

    const foundStudio = sources.find((s) => s.studio)?.studio;

    setStudio(
      new ScrapeResult<GQL.ScrapedStudio>(
        dest.studio ? idToStoredID(dest.studio) : undefined,
        foundStudio
          ? {
              stored_id: foundStudio.id,
              name: foundStudio.name,
            }
          : undefined,
        !dest.studio
      )
    );

    setPerformers(
      new ObjectListScrapeResult<GQL.ScrapedPerformer>(
        sortStoredIdObjects(dest.performers.map(idToStoredID)),
        uniqIDStoredIDs(all.flatMap((s) => s.performers.map(idToStoredID)))
      )
    );
    setTags(
      new ObjectListScrapeResult<GQL.ScrapedTag>(
        sortStoredIdObjects(dest.tags.map(idToStoredID)),
        uniqIDStoredIDs(all.flatMap((s) => s.tags.map(idToStoredID)))
      )
    );
    setDetails(
      new ScrapeResult(
        dest.details,
        sources.find((s) => s.details)?.details,
        !dest.details
      )
    );

    setGroups(
      new ObjectListScrapeResult<GQL.ScrapedGroup>(
        sortStoredIdObjects(dest.groups.map(groupToStoredID)),
        uniqIDStoredIDs(all.flatMap((s) => s.groups.map(groupToStoredID)))
      )
    );

    setGalleries(
      new ScrapeResult(
        dest.galleries.map((p) => p.id),
        uniq(all.flatMap((s) => s.galleries.map((p) => p.id)))
      )
    );

    setRating(
      new ScrapeResult(
        dest.rating100,
        sources.find((s) => s.rating100)?.rating100,
        !dest.rating100
      )
    );

    setOCounter(
      new ScrapeResult(
        dest.o_counter ?? 0,
        all.map((s) => s.o_counter ?? 0).reduce((pv, cv) => pv + cv, 0)
      )
    );

    setPlayCount(
      new ScrapeResult(
        dest.play_count ?? 0,
        all.map((s) => s.play_count ?? 0).reduce((pv, cv) => pv + cv, 0)
      )
    );

    setOrganized(
      new ScrapeResult(
        dest.organized ?? false,
        sources.every((s) => s.organized)
      )
    );

    const customFieldNames = new Set<string>(
      Object.keys(dest.custom_fields ?? {})
    );

    for (const s of sources) {
      for (const n of Object.keys(s.custom_fields ?? {})) {
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
        title,
        code,
        url,
        date,
        rating,
        oCounter,
        playCount,
        galleries,
        studio,
        performers,
        groups,
        tags,
        details,
        organized,
        image,
      ])
    );
  }, [
    title,
    code,
    url,
    date,
    rating,
    oCounter,
    playCount,
    galleries,
    studio,
    performers,
    groups,
    tags,
    details,
    organized,
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

    const trueString = intl.formatMessage({ id: "true" });
    const falseString = intl.formatMessage({ id: "false" });

    return (
      <>
        <ScrapedInputGroupRow
          field="title"
          title={intl.formatMessage({ id: "title" })}
          result={title}
          onChange={(value) => setTitle(value)}
        />
        <ScrapedInputGroupRow
          field="code"
          title={intl.formatMessage({ id: "audio_code" })}
          result={code}
          onChange={(value) => setCode(value)}
        />
        <ScrapedStringListRow
          field="urls"
          title={intl.formatMessage({ id: "urls" })}
          result={url}
          onChange={(value) => setURL(value)}
        />
        <ScrapedInputGroupRow
          field="date"
          title={intl.formatMessage({ id: "date" })}
          placeholder="YYYY-MM-DD"
          result={date}
          onChange={(value) => setDate(value)}
        />
        <ScrapeDialogRow
          field="rating"
          title={intl.formatMessage({ id: "rating" })}
          result={rating}
          originalField={<RatingSystem value={rating.originalValue} disabled />}
          newField={<RatingSystem value={rating.newValue} disabled />}
          onChange={(value) => setRating(value)}
        />
        <ScrapeDialogRow
          field="o_count"
          title={intl.formatMessage({ id: "o_count" })}
          result={oCounter}
          originalField={
            <FormControl
              value={oCounter.originalValue ?? 0}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          }
          newField={
            <FormControl
              value={oCounter.newValue ?? 0}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          }
          onChange={(value) => setOCounter(value)}
        />
        <ScrapeDialogRow
          field="play_count"
          title={intl.formatMessage({ id: "play_count" })}
          result={playCount}
          originalField={
            <FormControl
              value={playCount.originalValue ?? 0}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          }
          newField={
            <FormControl
              value={playCount.newValue ?? 0}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          }
          onChange={(value) => setPlayCount(value)}
        />
        <ScrapeDialogRow
          field="galleries"
          title={intl.formatMessage({ id: "galleries" })}
          result={galleries}
          originalField={
            <GallerySelect
              className="form-control react-select"
              ids={galleries.originalValue ?? []}
              onSelect={() => {}}
              isMulti
              isDisabled
            />
          }
          newField={
            <GallerySelect
              className="form-control react-select"
              ids={galleries.newValue ?? []}
              onSelect={() => {}}
              isMulti
              isDisabled
            />
          }
          onChange={(value) => setGalleries(value)}
        />
        <ScrapedStudioRow
          field="studio"
          title={intl.formatMessage({ id: "studios" })}
          result={studio}
          onChange={(value) => setStudio(value)}
        />
        <ScrapedPerformersRow
          field="performers"
          title={intl.formatMessage({ id: "performers" })}
          result={performers}
          onChange={(value) => setPerformers(value)}
          ageFromDate={date.useNewValue ? date.newValue : date.originalValue}
        />
        <ScrapedGroupsRow
          field="groups"
          title={intl.formatMessage({ id: "groups" })}
          result={groups}
          onChange={(value) => setGroups(value)}
        />
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
        <ScrapeDialogRow
          field="organized"
          title={intl.formatMessage({ id: "organized" })}
          result={organized}
          originalField={
            <FormControl
              value={organized.originalValue ? trueString : falseString}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          }
          newField={
            <FormControl
              value={organized.newValue ? trueString : falseString}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          }
          onChange={(value) => setOrganized(value)}
        />
        <ScrapedImageRow
          field="cover_image"
          title={intl.formatMessage({ id: "cover_image" })}
          className="audio-cover"
          result={image}
          onChange={(value) => setImage(value)}
        />
        {hasCustomFieldValues && (
          <ScrapedCustomFieldRows
            results={customFields}
            onChange={(newCustomFields) => setCustomFields(newCustomFields)}
          />
        )}
      </>
    );
  }

  function createValues(): MergeOptions {
    const all = [dest, ...sources];

    // only set the cover image if it's different from the existing cover image
    const coverImage = image.useNewValue ? image.getNewValue() : undefined;

    return {
      values: {
        id: dest.id,
        title: title.getNewValue(),
        code: code.getNewValue(),
        urls: url.getNewValue(),
        date: date.getNewValue(),
        rating100: rating.getNewValue(),
        gallery_ids: galleries.getNewValue(),
        studio_id: studio.getNewValue()?.stored_id,
        performer_ids: performers.getNewValue()?.map((p) => p.stored_id!),
        groups: groups.getNewValue()?.map((m) => {
          // find the equivalent group in the original audios
          const found = all
            .flatMap((s) => s.groups)
            .find((mm) => mm.group.id === m.stored_id);
          return {
            group_id: m.stored_id!,
            audio_index: found!.audio_index,
          };
        }),
        tag_ids: tags.getNewValue()?.map((t) => t.stored_id!),
        details: details.getNewValue(),
        organized: organized.getNewValue(),
        cover_image: coverImage,
        custom_fields: {
          partial: Object.fromEntries(
            Array.from(customFields.entries()).flatMap(([field, v]) =>
              v.useNewValue ? [[field, v.getNewValue()]] : []
            )
          ),
        },
      },
      includeViewHistory: playCount.getNewValue() !== undefined,
      includeOHistory: oCounter.getNewValue() !== undefined,
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
    : intl.formatMessage({ id: "dialogs.merge.combined" });

  return (
    <ScrapeDialog
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

interface IAudioMergeModalProps {
  show: boolean;
  onClose: (mergedID?: string) => void;
  audios: { id: string; title: string }[];
}

export const AudioMergeModal: React.FC<IAudioMergeModalProps> = ({
  show,
  onClose,
  audios,
}) => {
  const [sourceAudios, setSourceAudios] = useState<Audio[]>([]);
  const [destAudio, setDestAudio] = useState<Audio[]>([]);

  const [loadedSources, setLoadedSources] = useState<GQL.AudioDataFragment[]>(
    []
  );
  const [loadedDest, setLoadedDest] = useState<GQL.AudioDataFragment>();

  const [running, setRunning] = useState(false);
  const [secondStep, setSecondStep] = useState(false);

  const intl = useIntl();
  const Toast = useToast();

  const title = intl.formatMessage({
    id: "actions.merge",
  });

  const srcIDs = useMemo(() => sourceAudios.map((s) => s.id), [sourceAudios]);
  const destID = useMemo(
    () => (destAudio[0] ? [destAudio[0].id] : []),
    [destAudio]
  );

  useEffect(() => {
    if (audios.length > 0) {
      // set the first audio as the destination, others as source
      setDestAudio([audios[0]]);

      if (audios.length > 1) {
        setSourceAudios(audios.slice(1));
      }
    }
  }, [audios]);

  async function loadAudios() {
    const audioIDs = sourceAudios.map((s) => s.id);
    audioIDs.push(destAudio[0].id);
    const query = await queryFindFullAudiosByID(audioIDs);
    const { audios: loadedAudios } = query.data.findAudios;

    setLoadedDest(loadedAudios.find((s) => s.id === destAudio[0].id));
    setLoadedSources(loadedAudios.filter((s) => s.id !== destAudio[0].id));
    setSecondStep(true);
  }

  async function onMerge(options: MergeOptions) {
    const { values, includeViewHistory, includeOHistory } = options;
    try {
      setRunning(true);
      const result = await mutateAudioMerge(
        destAudio[0].id,
        sourceAudios.map((s) => s.id),
        values,
        includeViewHistory,
        includeOHistory
      );
      if (result.data?.audioMerge) {
        Toast.success(intl.formatMessage({ id: "toast.merged_audios" }));
        onClose(destAudio[0].id);
      }
      onClose();
    } catch (e) {
      Toast.error(e);
    } finally {
      setRunning(false);
    }
  }

  function canMerge() {
    return sourceAudios.length > 0 && destAudio.length !== 0;
  }

  function switchAudios() {
    if (sourceAudios.length && destAudio.length) {
      const newDest = sourceAudios[0];
      setSourceAudios([...sourceAudios.slice(1), destAudio[0]]);
      setDestAudio([newDest]);
    }
  }

  if (secondStep && destAudio.length > 0) {
    return (
      <AudioMergeDetails
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
      show={show}
      header={title}
      icon={faSignInAlt}
      accept={{
        text: intl.formatMessage({ id: "actions.next_action" }),
        onClick: () => loadAudios(),
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
              <AudioSelect
                isMulti
                onSelect={(items) => setSourceAudios(items)}
                values={sourceAudios}
                menuPortalTarget={document.body}
                excludeIds={destID}
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
              onClick={() => switchAudios()}
              disabled={!sourceAudios.length || !destAudio.length}
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
              <AudioSelect
                onSelect={(items) => setDestAudio(items)}
                values={destAudio}
                menuPortalTarget={document.body}
                excludeIds={srcIDs}
              />
            </Col>
          </Form.Group>
        </div>
      </div>
    </ModalComponent>
  );
};
