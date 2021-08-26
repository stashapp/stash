import React, { useState, useReducer, useEffect, useCallback } from "react";
import cx from "classnames";
import { Badge, Button, Col, Form, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import {
  Icon,
  LoadingIndicator,
  SuccessIcon,
  TagSelect,
  TruncatedText,
} from "src/components/Shared";
import { FormUtils } from "src/utils";
import { uniq } from "lodash";
import PerformerResult, { PerformerOperation } from "./PerformerResult";
import StudioResult, { StudioOperation } from "./StudioResult";
import { IStashBoxScene } from "./utils";
import { useTagScene } from "./taggerService";
import { TagOperation } from "./constants";
import { OptionalField } from "./IncludeButton";

const getDurationStatus = (
  scene: IStashBoxScene,
  stashDuration: number | undefined | null
) => {
  if (!stashDuration) return "";

  const durations = scene.fingerprints
    .map((f) => f.duration)
    .map((d) => Math.abs(d - stashDuration));
  const matchCount = durations.filter((duration) => duration <= 5).length;

  let match;
  if (matchCount > 0)
    match = (
      <FormattedMessage
        id="component_tagger.results.fp_matches_multi"
        values={{ matchCount, durationsLength: durations.length }}
      />
    );
  else if (Math.abs(scene.duration - stashDuration) < 5)
    match = <FormattedMessage id="component_tagger.results.fp_matches" />;

  if (match)
    return (
      <div className="font-weight-bold">
        <SuccessIcon className="mr-2" />
        {match}
      </div>
    );

  if (!scene.duration && durations.length === 0)
    return <FormattedMessage id="component_tagger.results.duration_unknown" />;

  const minDiff = Math.min(
    Math.abs(scene.duration - stashDuration),
    ...durations
  );
  return (
    <FormattedMessage
      id="component_tagger.results.duration_off"
      values={{ number: Math.floor(minDiff) }}
    />
  );
};

const getFingerprintStatus = (
  scene: IStashBoxScene,
  stashScene: GQL.SlimSceneDataFragment
) => {
  const checksumMatch = scene.fingerprints.some(
    (f) => f.hash === stashScene.checksum || f.hash === stashScene.oshash
  );
  const phashMatch = scene.fingerprints.some(
    (f) => f.hash === stashScene.phash
  );
  if (checksumMatch || phashMatch)
    return (
      <div className="font-weight-bold">
        <SuccessIcon className="mr-2" />
        <FormattedMessage
          id="component_tagger.results.hash_matches"
          values={{
            hash_type: (
              <FormattedMessage
                id={`media_info.${phashMatch ? "phash" : "checksum"}`}
              />
            ),
          }}
        />
      </div>
    );
};

interface IStashSearchResultProps {
  scene: IStashBoxScene;
  stashScene: GQL.SlimSceneDataFragment;
  isActive: boolean;
  setActive: () => void;
  showMales: boolean;
  setScene: (scene: GQL.SlimSceneDataFragment) => void;
  setCoverImage: boolean;
  tagOperation: TagOperation;
  setTags: boolean;
  endpoint: string;
  queueFingerprintSubmission: (sceneId: string, endpoint: string) => void;
  createNewTag: (toCreate: GQL.ScrapedTag) => void;
  excludedFields: Record<string, boolean>;
  setExcludedFields: (v: Record<string, boolean>) => void;
}

interface IPerformerReducerAction {
  id: string;
  data: PerformerOperation;
}

const performerReducer = (
  state: Record<string, PerformerOperation>,
  action: IPerformerReducerAction
) => ({ ...state, [action.id]: action.data });

const StashSearchResult: React.FC<IStashSearchResultProps> = ({
  scene,
  stashScene,
  isActive,
  setActive,
  showMales,
  setScene,
  setCoverImage,
  tagOperation,
  setTags,
  endpoint,
  queueFingerprintSubmission,
  createNewTag,
  excludedFields,
  setExcludedFields,
}) => {
  const getInitialTags = useCallback(() => {
    const stashSceneTags = stashScene.tags.map((t) => t.id);
    if (!setTags) {
      return stashSceneTags;
    }

    const newTags = scene.tags.filter((t) => t.id).map((t) => t.id!);
    if (tagOperation === "overwrite") {
      return newTags;
    }
    if (tagOperation === "merge") {
      return uniq(stashSceneTags.concat(newTags));
    }

    throw new Error("unexpected tagOperation");
  }, [stashScene, tagOperation, scene, setTags]);

  const [studio, setStudio] = useState<StudioOperation>();
  const [performers, dispatch] = useReducer(performerReducer, {});
  const [tagIDs, setTagIDs] = useState<string[]>(getInitialTags());
  const [saveState, setSaveState] = useState<string>("");
  const [error, setError] = useState<{ message?: string; details?: string }>(
    {}
  );

  const intl = useIntl();

  useEffect(() => {
    setTagIDs(getInitialTags());
  }, [setTags, tagOperation, getInitialTags]);

  const tagScene = useTagScene(
    {
      tagOperation,
      setCoverImage,
      setTags,
    },
    setSaveState,
    setError
  );

  function getExcludedFields() {
    return Object.keys(excludedFields).filter((f) => excludedFields[f]);
  }

  async function handleSave() {
    const updatedScene = await tagScene(
      stashScene,
      scene,
      studio,
      performers,
      tagIDs,
      getExcludedFields(),
      endpoint
    );

    if (updatedScene) setScene(updatedScene);

    queueFingerprintSubmission(stashScene.id, endpoint);
  }

  const setPerformer = (
    performerData: PerformerOperation,
    performerID: string
  ) => dispatch({ id: performerID, data: performerData });

  const setExcludedField = (name: string, value: boolean) =>
    setExcludedFields({
      ...excludedFields,
      [name]: value,
    });

  const classname = cx("row mx-0 mt-2 search-result", {
    "selected-result": isActive,
  });

  const sceneTitle = scene.url ? (
    <a
      href={scene.url}
      target="_blank"
      rel="noopener noreferrer"
      className="scene-link"
    >
      <TruncatedText text={scene?.title} />
    </a>
  ) : (
    <TruncatedText text={scene?.title} />
  );

  const saveEnabled =
    Object.keys(performers ?? []).length ===
      scene.performers.filter((p) => p.gender !== "MALE" || showMales).length &&
    Object.keys(performers ?? []).every((id) => performers?.[id].type) &&
    saveState === "";

  const endpointBase = endpoint.match(/https?:\/\/.*?\//)?.[0];
  const stashBoxURL = endpointBase
    ? `${endpointBase}scenes/${scene.stash_id}`
    : "";

  // constants to get around dot-notation eslint rule
  const fields = {
    cover_image: "cover_image",
    title: "title",
    date: "date",
    url: "url",
    details: "details",
  };

  return (
    // eslint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions
    <li
      className={`${classname} ${isActive && "active"}`}
      key={scene.stash_id}
      onClick={() => !isActive && setActive()}
    >
      <div className="col-lg-6">
        <div className="row">
          <div className="scene-image-container">
            <OptionalField
              exclude={excludedFields[fields.cover_image] || !setCoverImage}
              disabled={!setCoverImage}
              setExclude={(v) => setExcludedField(fields.cover_image, v)}
            >
              <a href={stashBoxURL} target="_blank" rel="noopener noreferrer">
                <img
                  src={scene.images[0]}
                  alt=""
                  className="align-self-center scene-image"
                />
              </a>
            </OptionalField>
          </div>
          <div className="d-flex flex-column justify-content-center scene-metadata">
            <h4>
              <OptionalField
                exclude={excludedFields[fields.title]}
                setExclude={(v) => setExcludedField(fields.title, v)}
              >
                {sceneTitle}
              </OptionalField>
            </h4>

            {!isActive && (
              <>
                <h5>
                  {scene?.studio?.name} • {scene?.date}
                </h5>
                <div>
                  {intl.formatMessage(
                    { id: "countables.performers" },
                    { count: scene?.performers?.length }
                  )}
                  : {scene?.performers?.map((p) => p.name).join(", ")}
                </div>
              </>
            )}

            {isActive && scene.date && (
              <h5>
                <OptionalField
                  exclude={excludedFields[fields.date]}
                  setExclude={(v) => setExcludedField(fields.date, v)}
                >
                  {scene.date}
                </OptionalField>
              </h5>
            )}
            {getDurationStatus(scene, stashScene.file?.duration)}
            {getFingerprintStatus(scene, stashScene)}
          </div>
        </div>
        {isActive && (
          <div className="d-flex flex-column">
            {scene.url && (
              <div className="scene-details">
                <OptionalField
                  exclude={excludedFields[fields.url]}
                  setExclude={(v) => setExcludedField(fields.url, v)}
                >
                  <a href={scene.url} target="_blank" rel="noopener noreferrer">
                    {scene.url}
                  </a>
                </OptionalField>
              </div>
            )}
            {scene.details && (
              <div className="scene-details">
                <OptionalField
                  exclude={excludedFields[fields.details]}
                  setExclude={(v) => setExcludedField(fields.details, v)}
                >
                  <TruncatedText text={scene.details ?? ""} lineCount={3} />
                </OptionalField>
              </div>
            )}
          </div>
        )}
      </div>
      {isActive && (
        <div className="col-lg-6">
          <StudioResult studio={scene.studio} setStudio={setStudio} />
          {scene.performers
            .filter((p) => p.gender !== "MALE" || showMales)
            .map((performer) => (
              <PerformerResult
                performer={performer}
                setPerformer={(data: PerformerOperation) =>
                  setPerformer(data, performer.stash_id)
                }
                key={`${scene.stash_id}${performer.stash_id}`}
                endpoint={endpoint}
              />
            ))}
          <div className="mt-2">
            <div>
              <Form.Group controlId="tags" as={Row}>
                {FormUtils.renderLabel({
                  title: `${intl.formatMessage({ id: "tags" })}:`,
                })}
                <Col sm={9} xl={12}>
                  <TagSelect
                    isDisabled={!setTags}
                    isMulti
                    onSelect={(items) => {
                      setTagIDs(items.map((i) => i.id));
                    }}
                    ids={tagIDs}
                  />
                </Col>
              </Form.Group>
            </div>
            {setTags &&
              scene.tags
                .filter((t) => !t.id)
                .map((t) => (
                  <Badge
                    className="tag-item"
                    variant="secondary"
                    key={t.name}
                    onClick={() => {
                      createNewTag(t);
                    }}
                  >
                    {t.name}
                    <Button className="minimal ml-2">
                      <Icon className="fa-fw" icon="plus" />
                    </Button>
                  </Badge>
                ))}
          </div>
          <div className="row no-gutters mt-2 align-items-center justify-content-end">
            {error.message && (
              <strong className="mt-1 mr-2 text-danger text-right">
                <abbr title={error.details} className="mr-2">
                  Error:
                </abbr>
                {error.message}
              </strong>
            )}
            {saveState && (
              <strong className="col-4 mt-1 mr-2 text-right">
                {saveState}
              </strong>
            )}
            <Button onClick={handleSave} disabled={!saveEnabled}>
              {saveState ? (
                <LoadingIndicator inline small message="" />
              ) : (
                <FormattedMessage id="actions.save" />
              )}
            </Button>
          </div>
        </div>
      )}
    </li>
  );
};

export default StashSearchResult;
