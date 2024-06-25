import { Form, Col, Row, Button, FormControl } from "react-bootstrap";
import React, { useEffect, useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { StringListSelect, GallerySelect } from "../Shared/Select";
import * as FormUtils from "src/utils/form";
import ImageUtils from "src/utils/image";
import TextUtils from "src/utils/text";
import { mutateSceneMerge, queryFindScenesByID } from "src/core/StashService";
import { FormattedMessage, useIntl } from "react-intl";
import { useToast } from "src/hooks/Toast";
import { faExchangeAlt, faSignInAlt } from "@fortawesome/free-solid-svg-icons";
import {
  ScrapeDialog,
  ScrapeDialogRow,
  ScrapedImageRow,
  ScrapedInputGroupRow,
  ScrapedStringListRow,
  ScrapedTextAreaRow,
} from "../Shared/ScrapeDialog/ScrapeDialog";
import { clone, uniq } from "lodash-es";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { ModalComponent } from "../Shared/Modal";
import { IHasStoredID, sortStoredIdObjects } from "src/utils/data";
import {
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
import { Scene, SceneSelect } from "src/components/Scenes/SceneSelect";

interface IStashIDsField {
  values: GQL.StashId[];
}

const StashIDsField: React.FC<IStashIDsField> = ({ values }) => {
  return <StringListSelect value={values.map((v) => v.stash_id)} />;
};

type MergeOptions = {
  values: GQL.SceneUpdateInput;
  includeViewHistory: boolean;
  includeOHistory: boolean;
};

interface ISceneMergeDetailsProps {
  sources: GQL.SlimSceneDataFragment[];
  dest: GQL.SlimSceneDataFragment;
  onClose: (options?: MergeOptions) => void;
}

const SceneMergeDetails: React.FC<ISceneMergeDetailsProps> = ({
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
  // zero values can be treated as missing for these fields
  const [oCounter, setOCounter] = useState(
    new ScrapeResult<number>(dest.o_counter)
  );
  const [playCount, setPlayCount] = useState(
    new ScrapeResult<number>(dest.play_count)
  );
  const [playDuration, setPlayDuration] = useState(
    new ScrapeResult<number>(dest.play_duration)
  );

  function idToStoredID(o: { id: string; name: string }) {
    return {
      stored_id: o.id,
      name: o.name,
    };
  }

  function groupToStoredID(o: { movie: { id: string; name: string } }) {
    return {
      stored_id: o.movie.id,
      name: o.movie.name,
    };
  }

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

  function uniqIDStoredIDs<T extends IHasStoredID>(objs: T[]) {
    return objs.filter((o, i) => {
      return objs.findIndex((oo) => oo.stored_id === o.stored_id) === i;
    });
  }

  const [performers, setPerformers] = useState<
    ObjectListScrapeResult<GQL.ScrapedPerformer>
  >(
    new ObjectListScrapeResult<GQL.ScrapedPerformer>(
      sortStoredIdObjects(dest.performers.map(idToStoredID))
    )
  );

  const [groups, setGroups] = useState<
    ObjectListScrapeResult<GQL.ScrapedMovie>
  >(
    new ObjectListScrapeResult<GQL.ScrapedMovie>(
      sortStoredIdObjects(dest.movies.map(groupToStoredID))
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

  const [stashIDs, setStashIDs] = useState(new ScrapeResult<GQL.StashId[]>([]));

  const [organized, setOrganized] = useState(
    new ZeroableScrapeResult<boolean>(dest.organized)
  );

  const [image, setImage] = useState<ScrapeResult<string>>(
    new ScrapeResult<string>(dest.paths.screenshot)
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

    // append dest to all so that if dest has stash_ids with the same
    // endpoint, then it will be excluded first
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
    setURL(
      new ScrapeResult(
        dest.urls,
        sources.find((s) => s.urls)?.urls,
        !dest.urls?.length
      )
    );
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
        uniqIDStoredIDs(all.map((s) => s.performers.map(idToStoredID)).flat())
      )
    );
    setTags(
      new ObjectListScrapeResult<GQL.ScrapedTag>(
        sortStoredIdObjects(dest.tags.map(idToStoredID)),
        uniqIDStoredIDs(all.map((s) => s.tags.map(idToStoredID)).flat())
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
      new ObjectListScrapeResult<GQL.ScrapedMovie>(
        sortStoredIdObjects(dest.movies.map(groupToStoredID)),
        uniqIDStoredIDs(all.map((s) => s.movies.map(groupToStoredID)).flat())
      )
    );

    setGalleries(
      new ScrapeResult(
        dest.galleries.map((p) => p.id),
        uniq(all.map((s) => s.galleries.map((p) => p.id)).flat())
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

    setPlayDuration(
      new ScrapeResult(
        dest.play_duration ?? 0,
        all.map((s) => s.play_duration ?? 0).reduce((pv, cv) => pv + cv, 0)
      )
    );

    setOrganized(
      new ScrapeResult(
        dest.organized ?? false,
        sources.every((s) => s.organized)
      )
    );

    setStashIDs(
      new ScrapeResult(
        dest.stash_ids,
        all
          .map((s) => s.stash_ids)
          .flat()
          .filter((s, index, a) => {
            // remove entries with duplicate endpoints
            return index === a.findIndex((ss) => ss.endpoint === s.endpoint);
          }),
        !dest.stash_ids.length
      )
    );

    loadImages();
  }, [sources, dest]);

  // ensure this is updated if fields are changed
  const hasValues = useMemo(() => {
    return hasScrapedValues([
      title,
      code,
      url,
      date,
      rating,
      oCounter,
      galleries,
      studio,
      performers,
      groups,
      tags,
      details,
      organized,
      stashIDs,
      image,
    ]);
  }, [
    title,
    code,
    url,
    date,
    rating,
    oCounter,
    galleries,
    studio,
    performers,
    groups,
    tags,
    details,
    organized,
    stashIDs,
    image,
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
          title={intl.formatMessage({ id: "title" })}
          result={title}
          onChange={(value) => setTitle(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "scene_code" })}
          result={code}
          onChange={(value) => setCode(value)}
        />
        <ScrapedStringListRow
          title={intl.formatMessage({ id: "urls" })}
          result={url}
          onChange={(value) => setURL(value)}
        />
        <ScrapedInputGroupRow
          title={intl.formatMessage({ id: "date" })}
          placeholder="YYYY-MM-DD"
          result={date}
          onChange={(value) => setDate(value)}
        />
        <ScrapeDialogRow
          title={intl.formatMessage({ id: "rating" })}
          result={rating}
          renderOriginalField={() => (
            <RatingSystem value={rating.originalValue} disabled />
          )}
          renderNewField={() => (
            <RatingSystem value={rating.newValue} disabled />
          )}
          onChange={(value) => setRating(value)}
        />
        <ScrapeDialogRow
          title={intl.formatMessage({ id: "o_count" })}
          result={oCounter}
          renderOriginalField={() => (
            <FormControl
              value={oCounter.originalValue ?? 0}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          )}
          renderNewField={() => (
            <FormControl
              value={oCounter.newValue ?? 0}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          )}
          onChange={(value) => setOCounter(value)}
        />
        <ScrapeDialogRow
          title={intl.formatMessage({ id: "play_count" })}
          result={playCount}
          renderOriginalField={() => (
            <FormControl
              value={playCount.originalValue ?? 0}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          )}
          renderNewField={() => (
            <FormControl
              value={playCount.newValue ?? 0}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          )}
          onChange={(value) => setPlayCount(value)}
        />
        <ScrapeDialogRow
          title={intl.formatMessage({ id: "play_duration" })}
          result={playDuration}
          renderOriginalField={() => (
            <FormControl
              value={TextUtils.secondsToTimestamp(
                playDuration.originalValue ?? 0
              )}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          )}
          renderNewField={() => (
            <FormControl
              value={TextUtils.secondsToTimestamp(playDuration.newValue ?? 0)}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          )}
          onChange={(value) => setPlayDuration(value)}
        />
        <ScrapeDialogRow
          title={intl.formatMessage({ id: "galleries" })}
          result={galleries}
          renderOriginalField={() => (
            <GallerySelect
              className="form-control react-select"
              ids={galleries.originalValue ?? []}
              onSelect={() => {}}
              isMulti
              isDisabled
            />
          )}
          renderNewField={() => (
            <GallerySelect
              className="form-control react-select"
              ids={galleries.newValue ?? []}
              onSelect={() => {}}
              isMulti
              isDisabled
            />
          )}
          onChange={(value) => setGalleries(value)}
        />
        <ScrapedStudioRow
          title={intl.formatMessage({ id: "studios" })}
          result={studio}
          onChange={(value) => setStudio(value)}
        />
        <ScrapedPerformersRow
          title={intl.formatMessage({ id: "performers" })}
          result={performers}
          onChange={(value) => setPerformers(value)}
        />
        <ScrapedGroupsRow
          title={intl.formatMessage({ id: "groups" })}
          result={groups}
          onChange={(value) => setGroups(value)}
        />
        <ScrapedTagsRow
          title={intl.formatMessage({ id: "tags" })}
          result={tags}
          onChange={(value) => setTags(value)}
        />
        <ScrapedTextAreaRow
          title={intl.formatMessage({ id: "details" })}
          result={details}
          onChange={(value) => setDetails(value)}
        />
        <ScrapeDialogRow
          title={intl.formatMessage({ id: "organized" })}
          result={organized}
          renderOriginalField={() => (
            <FormControl
              value={organized.originalValue ? trueString : falseString}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          )}
          renderNewField={() => (
            <FormControl
              value={organized.newValue ? trueString : falseString}
              readOnly
              onChange={() => {}}
              className="bg-secondary text-white border-secondary"
            />
          )}
          onChange={(value) => setOrganized(value)}
        />
        <ScrapeDialogRow
          title={intl.formatMessage({ id: "stash_id" })}
          result={stashIDs}
          renderOriginalField={() => (
            <StashIDsField values={stashIDs?.originalValue ?? []} />
          )}
          renderNewField={() => (
            <StashIDsField values={stashIDs?.newValue ?? []} />
          )}
          onChange={(value) => setStashIDs(value)}
        />
        <ScrapedImageRow
          title={intl.formatMessage({ id: "cover_image" })}
          className="scene-cover"
          result={image}
          onChange={(value) => setImage(value)}
        />
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
        o_counter: oCounter.getNewValue(),
        play_count: playCount.getNewValue(),
        play_duration: playDuration.getNewValue(),
        gallery_ids: galleries.getNewValue(),
        studio_id: studio.getNewValue()?.stored_id,
        performer_ids: performers.getNewValue()?.map((p) => p.stored_id!),
        movies: groups.getNewValue()?.map((m) => {
          // find the equivalent movie in the original scenes
          const found = all
            .map((s) => s.movies)
            .flat()
            .find((mm) => mm.movie.id === m.stored_id);
          return {
            movie_id: m.stored_id!,
            scene_index: found!.scene_index,
          };
        }),
        tag_ids: tags.getNewValue()?.map((t) => t.stored_id!),
        details: details.getNewValue(),
        organized: organized.getNewValue(),
        stash_ids: stashIDs.getNewValue(),
        cover_image: coverImage,
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
    : intl.formatMessage({ id: "dialogs.merge.source" });

  return (
    <ScrapeDialog
      title={dialogTitle}
      existingLabel={destinationLabel}
      scrapedLabel={sourceLabel}
      renderScrapeRows={renderScrapeRows}
      onClose={(apply) => {
        if (!apply) {
          onClose();
        } else {
          onClose(createValues());
        }
      }}
    />
  );
};

interface ISceneMergeModalProps {
  show: boolean;
  onClose: (mergedID?: string) => void;
  scenes: { id: string; title: string }[];
}

export const SceneMergeModal: React.FC<ISceneMergeModalProps> = ({
  show,
  onClose,
  scenes,
}) => {
  const [sourceScenes, setSourceScenes] = useState<Scene[]>([]);
  const [destScene, setDestScene] = useState<Scene[]>([]);

  const [loadedSources, setLoadedSources] = useState<
    GQL.SlimSceneDataFragment[]
  >([]);
  const [loadedDest, setLoadedDest] = useState<GQL.SlimSceneDataFragment>();

  const [running, setRunning] = useState(false);
  const [secondStep, setSecondStep] = useState(false);

  const intl = useIntl();
  const Toast = useToast();

  const title = intl.formatMessage({
    id: "actions.merge",
  });

  useEffect(() => {
    if (scenes.length > 0) {
      // set the first scene as the destination, others as source
      setDestScene([scenes[0]]);

      if (scenes.length > 1) {
        setSourceScenes(scenes.slice(1));
      }
    }
  }, [scenes]);

  async function loadScenes() {
    const sceneIDs = sourceScenes.map((s) => parseInt(s.id));
    sceneIDs.push(parseInt(destScene[0].id));
    const query = await queryFindScenesByID(sceneIDs);
    const { scenes: loadedScenes } = query.data.findScenes;

    setLoadedDest(loadedScenes.find((s) => s.id === destScene[0].id));
    setLoadedSources(loadedScenes.filter((s) => s.id !== destScene[0].id));
    setSecondStep(true);
  }

  async function onMerge(options: MergeOptions) {
    const { values, includeViewHistory, includeOHistory } = options;
    try {
      setRunning(true);
      const result = await mutateSceneMerge(
        destScene[0].id,
        sourceScenes.map((s) => s.id),
        values,
        includeViewHistory,
        includeOHistory
      );
      if (result.data?.sceneMerge) {
        Toast.success(intl.formatMessage({ id: "toast.merged_scenes" }));
        // refetch the scene
        await queryFindScenesByID([parseInt(destScene[0].id)]);
        onClose(destScene[0].id);
      }
      onClose();
    } catch (e) {
      Toast.error(e);
    } finally {
      setRunning(false);
    }
  }

  function canMerge() {
    return sourceScenes.length > 0 && destScene.length !== 0;
  }

  function switchScenes() {
    if (sourceScenes.length && destScene.length) {
      const newDest = sourceScenes[0];
      setSourceScenes([...sourceScenes.slice(1), destScene[0]]);
      setDestScene([newDest]);
    }
  }

  if (secondStep && destScene.length > 0) {
    return (
      <SceneMergeDetails
        sources={loadedSources}
        dest={loadedDest!}
        onClose={(values) => {
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
        onClick: () => loadScenes(),
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
              <SceneSelect
                isMulti
                onSelect={(items) => setSourceScenes(items)}
                values={sourceScenes}
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
              onClick={() => switchScenes()}
              disabled={!sourceScenes.length || !destScene.length}
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
              <SceneSelect
                onSelect={(items) => setDestScene(items)}
                values={destScene}
              />
            </Col>
          </Form.Group>
        </div>
      </div>
    </ModalComponent>
  );
};
