import React, { useState, useCallback, useEffect, useMemo } from "react";
import { Badge, Button, Col, Form, Row } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import {
  Icon,
  LoadingIndicator,
  PerformerSelect,
  StudioSelect,
  SuccessIcon,
  TagSelect,
  TruncatedText,
} from "src/components/Shared";
import { FormUtils } from "src/utils";
import { uniq } from "lodash";
import { OptionalField } from "src/components/Tagger/IncludeButton";
import { stringToGender } from "src/utils/gender";
import { blobToBase64 } from "base64-blob";
import { OperationButton } from "src/components/Shared/OperationButton";
import { IScrapedScene, SceneScraperStateContext } from "./context";
import { SceneScraperDialogsState } from "./modals";

const getDurationStatus = (
  scene: IScrapedScene,
  stashDuration: number | undefined | null
) => {
  if (!stashDuration) return "";

  const durations =
    scene.fingerprints
      ?.map((f) => f.duration)
      .map((d) => Math.abs(d - stashDuration)) ?? [];

  const sceneDuration = scene.duration ?? 0;

  if (!sceneDuration && durations.length === 0) return "";

  const matchCount = durations.filter((duration) => duration <= 5).length;

  let match;
  if (matchCount > 0)
    match = (
      <FormattedMessage
        id="component_tagger.results.fp_matches_multi"
        values={{ matchCount, durationsLength: durations.length }}
      />
    );
  else if (Math.abs(sceneDuration - stashDuration) < 5)
    match = <FormattedMessage id="component_tagger.results.fp_matches" />;

  if (match)
    return (
      <div className="font-weight-bold">
        <SuccessIcon className="mr-2" />
        {match}
      </div>
    );

  const minDiff = Math.min(
    Math.abs(sceneDuration - stashDuration),
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
  scene: IScrapedScene,
  stashScene: GQL.SlimSceneDataFragment
) => {
  const checksumMatch = scene.fingerprints?.some(
    (f) => f.hash === stashScene.checksum || f.hash === stashScene.oshash
  );
  const phashMatch = scene.fingerprints?.some(
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
  scene: IScrapedScene;
  stashScene: GQL.SlimSceneDataFragment;
  index: number;
  isActive: boolean;
}

const SceneScraperSceneEditor: React.FC<IStashSearchResultProps> = ({
  scene,
  stashScene,
  index,
  isActive,
}) => {
  const intl = useIntl();

  const {
    config,
    createNewTag,
    createNewPerformer,
    createNewStudio,
    resolveScene,
    currentSource,
    saveScene,
  } = React.useContext(SceneScraperStateContext);

  const performers = useMemo(
    () =>
      scene.performers?.filter((p) => {
        if (!config.showMales) {
          return (
            !p.gender || stringToGender(p.gender, true) !== GQL.GenderEnum.Male
          );
        }
        return true;
      }) ?? [],
    [config, scene]
  );

  const { createPerformerModal, createStudioModal } = React.useContext(
    SceneScraperDialogsState
  );

  const getInitialTags = useCallback(() => {
    const stashSceneTags = stashScene.tags.map((t) => t.id);
    if (!config.setTags) {
      return stashSceneTags;
    }

    const { tagOperation } = config;

    const newTags =
      scene.tags?.filter((t) => t.stored_id).map((t) => t.stored_id!) ?? [];

    if (tagOperation === "overwrite") {
      return newTags;
    }
    if (tagOperation === "merge") {
      return uniq(stashSceneTags.concat(newTags));
    }

    throw new Error("unexpected tagOperation");
  }, [stashScene, scene, config]);

  const getInitialPerformers = useCallback(() => {
    // default to override existing
    return performers.filter((t) => t.stored_id).map((t) => t.stored_id!) ?? [];
  }, [performers]);

  const getInitialStudio = useCallback(() => {
    return scene.studio?.stored_id ?? stashScene.studio?.id;
  }, [stashScene, scene]);

  const [loading, setLoading] = useState(false);
  const [excludedFields, setExcludedFields] = useState<Record<string, boolean>>(
    {}
  );
  const [tagIDs, setTagIDs] = useState<string[]>(getInitialTags());
  const [performerIDs, setPerformerIDs] = useState<string[]>(
    getInitialPerformers()
  );
  const [studioID, setStudioID] = useState<string | undefined>(
    getInitialStudio()
  );

  useEffect(() => {
    setTagIDs(getInitialTags());
  }, [getInitialTags]);

  useEffect(() => {
    setPerformerIDs(getInitialPerformers());
  }, [getInitialPerformers]);

  useEffect(() => {
    setStudioID(getInitialStudio());
  }, [getInitialStudio]);

  useEffect(() => {
    async function doResolveScene() {
      try {
        setLoading(true);
        await resolveScene(stashScene.id, index, scene);
      } finally {
        setLoading(false);
      }
    }

    if (isActive && !loading && !scene.resolved) {
      doResolveScene();
    }
  }, [isActive, loading, stashScene, index, resolveScene, scene]);

  // function getExcludedFields() {
  //   return Object.keys(excludedFields).filter((f) => excludedFields[f]);
  // }

  const setExcludedField = (name: string, value: boolean) =>
    setExcludedFields({
      ...excludedFields,
      [name]: value,
    });

  async function handleSave() {
    const excludedFieldList = Object.keys(excludedFields).filter(
      (f) => excludedFields[f]
    );

    function resolveField<T>(field: string, stashField: T, remoteField: T) {
      if (excludedFieldList.includes(field)) {
        return stashField;
      }

      return remoteField;
    }

    let imgData;
    if (!excludedFields.cover_image && config.setCoverImage) {
      const imgurl = scene.image;
      if (imgurl) {
        const img = await fetch(imgurl, {
          mode: "cors",
          cache: "no-store",
        });
        if (img.status === 200) {
          const blob = await img.blob();
          // Sanity check on image size since bad images will fail
          if (blob.size > 10000) imgData = await blobToBase64(blob);
        }
      }
    }

    const sceneCreateInput: GQL.SceneUpdateInput = {
      id: stashScene.id ?? "",
      title: resolveField("title", stashScene.title, scene.title),
      details: resolveField("details", stashScene.details, scene.details),
      date: resolveField("date", stashScene.date, scene.date),
      performer_ids:
        performerIDs.length === 0
          ? stashScene.performers.map((p) => p.id)
          : performerIDs,
      studio_id: studioID,
      cover_image: resolveField("cover_image", undefined, imgData),
      url: resolveField("url", stashScene.url, scene.url),
      tag_ids: tagIDs,
      stash_ids: stashScene.stash_ids ?? [],
    };

    if (currentSource?.stashboxEndpoint && scene.remote_site_id) {
      sceneCreateInput.stash_ids = [
        ...(stashScene?.stash_ids ?? []),
        {
          endpoint: currentSource.stashboxEndpoint,
          stash_id: scene.remote_site_id,
        },
      ];
    }

    await saveScene(sceneCreateInput);
  }

  function performerModalCallback(
    toCreate?: GQL.PerformerCreateInput | undefined
  ) {
    if (toCreate) {
      createNewPerformer(toCreate);
    }
  }

  function showPerformerModal(t: GQL.ScrapedPerformer) {
    createPerformerModal(t, performerModalCallback);
  }

  function studioModalCallback(toCreate?: GQL.StudioCreateInput | undefined) {
    if (toCreate) {
      createNewStudio(toCreate);
    }
  }

  function showStudioModal(t: GQL.ScrapedStudio) {
    createStudioModal(t, studioModalCallback);
  }

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

  // constants to get around dot-notation eslint rule
  const fields = {
    cover_image: "cover_image",
    title: "title",
    date: "date",
    url: "url",
    details: "details",
    studio: "studio",
  };

  const maybeRenderCoverImage = () => {
    if (scene.image) {
      return (
        <div className="scene-image-container">
          <OptionalField
            disabled={!config.setCoverImage}
            exclude={
              excludedFields[fields.cover_image] || !config.setCoverImage
            }
            setExclude={(v) => setExcludedField(fields.cover_image, v)}
          >
            <a href="TODO" target="_blank" rel="noopener noreferrer">
              <img
                src={scene.image}
                alt=""
                className="align-self-center scene-image"
              />
            </a>
          </OptionalField>
        </div>
      );
    }
  };

  const renderTitle = () => (
    <h4>
      <OptionalField
        exclude={excludedFields[fields.title]}
        setExclude={(v) => setExcludedField(fields.title, v)}
      >
        {sceneTitle}
      </OptionalField>
    </h4>
  );

  function renderStudioDate() {
    const text =
      scene.studio && scene.date
        ? `${scene.studio.name} • ${scene.date}`
        : `${scene.studio?.name ?? scene.date ?? ""}`;

    if (text) {
      return <h5>{text}</h5>;
    }
  }

  const renderPerformerList = () => {
    if (scene.performers?.length) {
      return (
        <div>
          {intl.formatMessage(
            { id: "countables.performers" },
            { count: scene?.performers?.length }
          )}
          : {scene?.performers?.map((p) => p.name).join(", ")}
        </div>
      );
    }
  };

  const maybeRenderDateField = () => {
    if (isActive && scene.date) {
      return (
        <h5>
          <OptionalField
            exclude={excludedFields[fields.date]}
            setExclude={(v) => setExcludedField(fields.date, v)}
          >
            {scene.date}
          </OptionalField>
        </h5>
      );
    }
  };

  const maybeRenderURL = () => {
    if (scene.url) {
      return (
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
      );
    }
  };

  const maybeRenderDetails = () => {
    if (scene.details) {
      return (
        <div className="scene-details">
          <OptionalField
            exclude={excludedFields[fields.details]}
            setExclude={(v) => setExcludedField(fields.details, v)}
          >
            <TruncatedText text={scene.details ?? ""} lineCount={3} />
          </OptionalField>
        </div>
      );
    }
  };

  const renderStudioField = () => (
    <div className="mt-2">
      <div>
        <Form.Group controlId="studio" as={Row}>
          {FormUtils.renderLabel({
            title: `${intl.formatMessage({ id: "studio" })}:`,
          })}
          <Col sm={9} xl={12}>
            <StudioSelect
              onSelect={(items) => {
                setStudioID(items[0]?.id);
              }}
              ids={studioID ? [studioID] : []}
            />
          </Col>
        </Form.Group>
      </div>
      {scene.studio && !scene.studio.stored_id && (
        <Badge
          className="tag-item"
          variant="secondary"
          key={scene.studio.name}
          onClick={() => {
            showStudioModal(scene.studio!);
          }}
        >
          {scene.studio.name}
          <Button className="minimal ml-2">
            <Icon className="fa-fw" icon="plus" />
          </Button>
        </Badge>
      )}
    </div>
  );

  const renderPerformerField = () => (
    <div className="mt-2">
      <div>
        <Form.Group controlId="performers" as={Row}>
          {FormUtils.renderLabel({
            title: `${intl.formatMessage({ id: "performers" })}:`,
          })}
          <Col sm={9} xl={12}>
            <PerformerSelect
              isMulti
              onSelect={(items) => {
                setPerformerIDs(items.map((i) => i.id));
              }}
              ids={performerIDs}
            />
          </Col>
        </Form.Group>
      </div>
      {performers
        ?.filter((p) => !p.stored_id)
        .map((p) => (
          <Badge
            className="tag-item"
            variant="secondary"
            key={p.name}
            onClick={() => {
              showPerformerModal(p);
            }}
          >
            {p.name}
            <Button className="minimal ml-2">
              <Icon className="fa-fw" icon="plus" />
            </Button>
          </Badge>
        ))}
    </div>
  );

  const renderTagsField = () => (
    <div className="mt-2">
      <div>
        <Form.Group controlId="tags" as={Row}>
          {FormUtils.renderLabel({
            title: `${intl.formatMessage({ id: "tags" })}:`,
          })}
          <Col sm={9} xl={12}>
            <TagSelect
              isMulti
              onSelect={(items) => {
                setTagIDs(items.map((i) => i.id));
              }}
              ids={tagIDs}
            />
          </Col>
        </Form.Group>
      </div>
      {scene.tags
        ?.filter((t) => !t.stored_id)
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
  );

  if (loading) {
    return <LoadingIndicator card />;
  }

  return (
    <>
      <div className="col-lg-6">
        <div className="row pl-3">
          {maybeRenderCoverImage()}
          <div className="d-flex flex-column justify-content-center scene-metadata">
            {renderTitle()}

            {!isActive && (
              <>
                {renderStudioDate()}
                {renderPerformerList()}
              </>
            )}

            {maybeRenderDateField()}
            {getDurationStatus(scene, stashScene.file?.duration)}
            {getFingerprintStatus(scene, stashScene)}
          </div>
        </div>
        {isActive && (
          <div className="d-flex flex-column">
            {maybeRenderURL()}
            {maybeRenderDetails()}
          </div>
        )}
      </div>
      {isActive && (
        <div className="col-lg-6">
          {renderStudioField()}
          {renderPerformerField()}
          {renderTagsField()}

          <div className="row no-gutters mt-2 align-items-center justify-content-end">
            <OperationButton operation={handleSave}>
              <FormattedMessage id="actions.save" />
            </OperationButton>
          </div>
        </div>
      )}
    </>
  );
};

export default SceneScraperSceneEditor;
