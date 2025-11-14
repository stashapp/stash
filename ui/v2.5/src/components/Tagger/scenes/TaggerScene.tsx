import React, { useState, useContext, PropsWithChildren, useMemo } from "react";
import * as GQL from "src/core/generated-graphql";
import { Link, useHistory } from "react-router-dom";
import { Button, Collapse, Form, InputGroup } from "react-bootstrap";
import { FormattedMessage } from "react-intl";

import { sortPerformers } from "src/core/performers";
import { Icon } from "src/components/Shared/Icon";
import { OperationButton } from "src/components/Shared/OperationButton";
import { StashIDPill } from "src/components/Shared/StashID";
import { PerformerLink, TagLink } from "src/components/Shared/TagLink";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { parsePath, prepareQueryString } from "src/components/Tagger/utils";
import { ScenePreview } from "src/components/Scenes/SceneCard";
import { TaggerStateContext } from "../context";
import {
  faChevronDown,
  faChevronUp,
  faImage,
} from "@fortawesome/free-solid-svg-icons";
import { objectPath, objectTitle } from "src/core/files";
import { useConfigurationContext } from "src/hooks/Config";
import { SceneQueue } from "src/models/sceneQueue";

interface ITaggerSceneDetails {
  scene: GQL.SlimSceneDataFragment;
}

const TaggerSceneDetails: React.FC<ITaggerSceneDetails> = ({ scene }) => {
  const [open, setOpen] = useState(false);
  const sorted = sortPerformers(scene.performers);

  return (
    <div className="original-scene-details">
      <Collapse in={open}>
        <div className="row">
          <div className="col col-lg-6">
            <h4>{objectTitle(scene)}</h4>
            <h5>
              {scene.studio?.name}
              {scene.studio?.name && scene.date && ` • `}
              {scene.date}
            </h5>
            <TruncatedText text={scene.details ?? ""} lineCount={3} />
          </div>
          <div className="col col-lg-6">
            <div>
              {sorted.map((performer) => (
                <div className="performer-tag-container row" key={performer.id}>
                  <Link
                    to={`/performers/${performer.id}`}
                    className="performer-tag col m-auto zoom-2"
                  >
                    <img
                      loading="lazy"
                      className="image-thumbnail"
                      alt={performer.name ?? ""}
                      src={performer.image_path ?? ""}
                    />
                  </Link>
                  <PerformerLink
                    key={performer.id}
                    performer={performer}
                    className="d-block"
                  />
                </div>
              ))}
            </div>
            <div>
              {scene.tags.map((tag) => (
                <TagLink key={tag.id} tag={tag} />
              ))}
            </div>
          </div>
        </div>
      </Collapse>
      <Button
        onClick={() => setOpen(!open)}
        className="minimal collapse-button"
        size="lg"
      >
        <Icon icon={open ? faChevronUp : faChevronDown} />
      </Button>
    </div>
  );
};

type StashID = Pick<GQL.StashId, "endpoint" | "stash_id">;

const StashIDs: React.FC<{ stashIDs: StashID[] }> = ({ stashIDs }) => {
  if (!stashIDs.length) {
    return null;
  }

  const stashLinks = stashIDs.map((stashID) => {
    const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
    const link = base ? (
      <StashIDPill stashID={stashID} linkType="scenes" />
    ) : (
      <span className="small">{stashID.stash_id}</span>
    );

    return <div key={stashID.stash_id}>{link}</div>;
  });

  return <div className="mt-2 sub-content text-right">{stashLinks}</div>;
};

interface ITaggerScene {
  scene: GQL.SlimSceneDataFragment;
  url: string;
  errorMessage?: string;
  doSceneQuery?: (queryString: string) => void;
  scrapeSceneFragment?: (scene: GQL.SlimSceneDataFragment) => void;
  loading?: boolean;
  showLightboxImage: (imagePath: string) => void;
  queue?: SceneQueue;
  index?: number;
}

export const TaggerScene: React.FC<PropsWithChildren<ITaggerScene>> = ({
  scene,
  url,
  loading,
  doSceneQuery,
  scrapeSceneFragment,
  errorMessage,
  children,
  showLightboxImage,
  queue,
  index,
}) => {
  const { config } = useContext(TaggerStateContext);
  const [queryString, setQueryString] = useState<string>("");
  const [queryLoading, setQueryLoading] = useState(false);

  const { paths, file: basename } = parsePath(objectPath(scene));
  const defaultQueryString = prepareQueryString(
    scene,
    paths,
    basename,
    config.mode,
    config.blacklist
  );

  const file = useMemo(
    () => (scene.files.length > 0 ? scene.files[0] : undefined),
    [scene]
  );

  const width = file?.width ? file.width : 0;
  const height = file?.height ? file.height : 0;
  const isPortrait = height > width;

  const history = useHistory();

  const { configuration } = useConfigurationContext();
  const cont = configuration?.interface.continuePlaylistDefault ?? false;

  async function query() {
    if (!doSceneQuery) return;

    try {
      setQueryLoading(true);
      await doSceneQuery(queryString || defaultQueryString);
    } finally {
      setQueryLoading(false);
    }
  }

  function renderQueryForm() {
    if (!doSceneQuery) return;

    return (
      <InputGroup>
        <InputGroup.Prepend>
          <InputGroup.Text>
            <FormattedMessage id="component_tagger.noun_query" />
          </InputGroup.Text>
        </InputGroup.Prepend>
        <Form.Control
          className="text-input"
          value={queryString || defaultQueryString}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setQueryString(e.currentTarget.value);
          }}
          onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) =>
            e.key === "Enter" && query()
          }
        />
        <InputGroup.Append>
          <OperationButton
            disabled={loading}
            operation={query}
            loading={queryLoading}
            setLoading={setQueryLoading}
          >
            <FormattedMessage id="actions.search" />
          </OperationButton>
        </InputGroup.Append>
      </InputGroup>
    );
  }

  function onSpriteClick(ev: React.MouseEvent<HTMLElement>) {
    ev.preventDefault();
    showLightboxImage(scene.paths.sprite ?? "");
  }

  function maybeRenderSpriteIcon() {
    // If a scene doesn't have any files, or doesn't have a sprite generated, the
    // path will be http://localhost:9999/scene/_sprite.jpg
    if (scene.files.length > 0) {
      return (
        <Button
          className="sprite-button"
          variant="link"
          onClick={onSpriteClick}
        >
          <Icon icon={faImage} />
        </Button>
      );
    }
  }

  function onScrubberClick(timestamp: number) {
    const link = queue
      ? queue.makeLink(scene.id, {
          sceneIndex: index,
          continue: cont,
          start: timestamp,
        })
      : `/scenes/${scene.id}?t=${timestamp}`;

    history.push(link);
  }

  return (
    <div key={scene.id} className="mt-3 search-item">
      <div className="row">
        <div className="col col-lg-6 overflow-hidden align-items-center d-flex flex-column flex-sm-row">
          <div className="scene-card mr-3">
            <Link to={url}>
              <ScenePreview
                image={scene.paths.screenshot ?? undefined}
                video={scene.paths.preview ?? undefined}
                isPortrait={isPortrait}
                soundActive={false}
                vttPath={scene.paths.vtt ?? undefined}
                onScrubberClick={onScrubberClick}
              />
              {maybeRenderSpriteIcon()}
            </Link>
          </div>
          <Link to={url} className="scene-link overflow-hidden">
            <TruncatedText text={objectTitle(scene)} lineCount={2} />
          </Link>
        </div>
        <div className="col-md-6 my-1">
          <div>
            {renderQueryForm()}
            {scrapeSceneFragment ? (
              <div className="mt-2 text-right">
                <OperationButton
                  disabled={loading}
                  operation={async () => {
                    await scrapeSceneFragment(scene);
                  }}
                >
                  <FormattedMessage id="actions.scrape_scene_fragment" />
                </OperationButton>
              </div>
            ) : undefined}
          </div>
          {errorMessage ? (
            <div className="text-danger font-weight-bold">{errorMessage}</div>
          ) : undefined}
          <StashIDs stashIDs={scene.stash_ids} />
        </div>
        <TaggerSceneDetails scene={scene} />
      </div>
      {children}
    </div>
  );
};
