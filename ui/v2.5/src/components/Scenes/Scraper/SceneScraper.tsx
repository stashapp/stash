import React, { useContext, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneQueue } from "src/models/sceneQueue";
import { Button, Card, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon, LoadingIndicator } from "src/components/Shared";
import { OperationButton } from "src/components/Shared/OperationButton";
import { SceneScraperScene } from "./SceneScraperScene";
import { SceneScraperStateContext } from "./context";
import { SceneSearchResults } from "./SceneScraperSearchResult";
import { SceneScraperModals } from "./modals";
import Config from "./Config";

interface IScraperProps {
  scenes: GQL.SlimSceneDataFragment[];
  queue?: SceneQueue;
}

export const SceneScraper: React.FC<IScraperProps> = ({ scenes, queue }) => {
  const {
    sources,
    setCurrentSource,
    currentSource,
    doSceneQuery,
    doSceneFragmentScrape,
    doMultiSceneFragmentScrape,
    stopMultiScrape,
    searchResults,
    loading,
    loadingMulti,
    multiError,
    submitFingerprints,
    pendingFingerprints,
  } = useContext(SceneScraperStateContext);

  const [showConfig, setShowConfig] = useState(false);
  const [hideUnmatched, setHideUnmatched] = useState(false);

  const intl = useIntl();

  function generateSceneLink(scene: GQL.SlimSceneDataFragment, index: number) {
    return queue
      ? queue.makeLink(scene.id, { sceneIndex: index })
      : `/scenes/${scene.id}`;
  }

  function handleSourceSelect(e: React.ChangeEvent<HTMLSelectElement>) {
    setCurrentSource(sources!.find((s) => s.id === e.currentTarget.value));
  }

  function renderSourceSelector() {
    return (
      <Form.Group controlId="scraper">
        <Form.Label>
          <FormattedMessage id="component_tagger.config.source" />
        </Form.Label>
        <div>
          <Form.Control
            as="select"
            value={currentSource?.id}
            className="input-control"
            disabled={loading || !sources.length}
            onChange={handleSourceSelect}
          >
            {!sources.length && <option>No scraper sources</option>}
            {sources.map((i) => (
              <option value={i.id} key={i.id}>
                {i.displayName}
              </option>
            ))}
          </Form.Control>
        </div>
      </Form.Group>
    );
  }

  function renderConfigButton() {
    return (
      <div className="ml-2">
        <Button onClick={() => setShowConfig(!showConfig)}>
          <Icon className="fa-fw" icon="cog" />
        </Button>
      </div>
    );
  }

  function renderScenes() {
    const filteredScenes = !hideUnmatched
      ? scenes
      : scenes.filter((s) => searchResults[s.id]?.results?.length);

    return filteredScenes.map((scene, index) => {
      const sceneLink = generateSceneLink(scene, index);
      let errorMessage: string | undefined;
      const searchResult = searchResults[scene.id];
      if (searchResult?.error) {
        errorMessage = searchResult.error;
      } else if (searchResult && searchResult.results?.length === 0) {
        errorMessage = intl.formatMessage({
          id: "component_tagger.results.match_failed_no_result",
        });
      }

      return (
        <SceneScraperScene
          key={scene.id}
          loading={loading}
          scene={scene}
          url={sceneLink}
          errorMessage={errorMessage}
          doSceneQuery={
            currentSource?.supportQuery
              ? async (v) => {
                  await doSceneQuery(scene.id, v);
                }
              : undefined
          }
          scrapeSceneFragment={
            currentSource?.supportFragment
              ? async () => {
                  await doSceneFragmentScrape(scene.id);
                }
              : undefined
          }
        >
          {searchResult && searchResult.results?.length ? (
            <SceneSearchResults scenes={searchResult.results} target={scene} />
          ) : undefined}
        </SceneScraperScene>
      );
    });
  }

  const toggleHideUnmatchedScenes = () => {
    setHideUnmatched(!hideUnmatched);
  };

  function maybeRenderShowHideUnmatchedButton() {
    if (Object.keys(searchResults).length) {
      return (
        <Button onClick={toggleHideUnmatchedScenes}>
          <FormattedMessage
            id="component_tagger.verb_toggle_unmatched"
            values={{
              toggle: (
                <FormattedMessage
                  id={`actions.${!hideUnmatched ? "hide" : "show"}`}
                />
              ),
            }}
          />
        </Button>
      );
    }
  }

  function maybeRenderSubmitFingerprintsButton() {
    if (pendingFingerprints.length) {
      return (
        <OperationButton
          className="ml-1"
          operation={submitFingerprints}
          disabled={loading || loadingMulti}
        >
          <span>
            <FormattedMessage
              id="component_tagger.verb_submit_fp"
              values={{ fpCount: pendingFingerprints.length }}
            />
          </span>
        </OperationButton>
      );
    }
  }

  function renderFragmentScrapeButton() {
    if (!currentSource?.supportFragment) {
      return;
    }

    if (loadingMulti) {
      return (
        <Button
          className="ml-1"
          variant="danger"
          onClick={() => {
            stopMultiScrape();
          }}
        >
          <LoadingIndicator message="" inline small />
          <span className="ml-2">
            {intl.formatMessage({ id: "actions.stop" })}
          </span>
        </Button>
      );
    }

    return (
      <div className="ml-1">
        <OperationButton
          disabled={loading}
          operation={async () => {
            await doMultiSceneFragmentScrape(scenes.map((s) => s.id));
          }}
        >
          {intl.formatMessage({ id: "component_tagger.verb_scrape_all" })}
        </OperationButton>
        {multiError && (
          <>
            <br />
            <b className="text-danger">{multiError}</b>
          </>
        )}
      </div>
    );
  }

  return (
    <SceneScraperModals>
      <Card className="tagger-container">
        <Card.Header className="tagger-container-header">
          <div className="d-flex justify-content-between align-items-center">
            {renderSourceSelector()}
            <div className="ml-auto d-flex">
              {maybeRenderShowHideUnmatchedButton()}
              {maybeRenderSubmitFingerprintsButton()}
              {renderFragmentScrapeButton()}
              {renderConfigButton()}
            </div>
          </div>
          <Config show={showConfig} />
        </Card.Header>
        <Card.Body>{renderScenes()}</Card.Body>
      </Card>
    </SceneScraperModals>
  );
};
