import React, { useContext, useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneQueue } from "src/models/sceneQueue";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";

import { Icon } from "src/components/Shared/Icon";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { OperationButton } from "src/components/Shared/OperationButton";
import { IScrapedScene, TaggerStateContext } from "../context";
import Config from "./Config";
import { TaggerScene } from "./TaggerScene";
import { SceneTaggerModals } from "./sceneTaggerModals";
import { SceneSearchResults } from "./StashSearchResult";
import { ConfigurationContext } from "src/hooks/Config";
import { faCog } from "@fortawesome/free-solid-svg-icons";
import { distance } from "src/utils/hamming";
import { useLightbox } from "src/hooks/Lightbox/hooks";

interface ITaggerProps {
  scenes: GQL.SlimSceneDataFragment[];
  queue?: SceneQueue;
}

export const Tagger: React.FC<ITaggerProps> = ({ scenes, queue }) => {
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
  } = useContext(TaggerStateContext);
  const { configuration } = React.useContext(ConfigurationContext);

  const [showConfig, setShowConfig] = useState(false);
  const [hideUnmatched, setHideUnmatched] = useState(false);

  const intl = useIntl();

  const cont = configuration?.interface.continuePlaylistDefault ?? false;

  function generateSceneLink(scene: GQL.SlimSceneDataFragment, index: number) {
    return queue
      ? queue.makeLink(scene.id, { sceneIndex: index, continue: cont })
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
          <Icon className="fa-fw" icon={faCog} />
        </Button>
      </div>
    );
  }

  function minDistance(hash: string, stashScene: GQL.SlimSceneDataFragment) {
    let ret = 9999;
    stashScene.files.forEach((cv) => {
      if (ret === 0) return;

      const stashHash = cv.fingerprints.find((fp) => fp.type === "phash");
      if (!stashHash) {
        return;
      }

      const d = distance(hash, stashHash.value);
      if (d < ret) {
        ret = d;
      }
    });

    return ret;
  }

  function calculatePhashComparisonScore(
    stashScene: GQL.SlimSceneDataFragment,
    scrapedScene: IScrapedScene
  ) {
    const phashFingerprints =
      scrapedScene.fingerprints?.filter((f) => f.algorithm === "PHASH") ?? [];
    const filteredFingerprints = phashFingerprints.filter(
      (f) => minDistance(f.hash, stashScene) <= 8
    );

    if (phashFingerprints.length == 0) return [0, 0];

    return [
      filteredFingerprints.length,
      filteredFingerprints.length / phashFingerprints.length,
    ];
  }

  function minDurationDiff(
    stashScene: GQL.SlimSceneDataFragment,
    duration: number
  ) {
    let ret = 9999;
    stashScene.files.forEach((cv) => {
      if (ret === 0) return;

      const d = Math.abs(duration - cv.duration);
      if (d < ret) {
        ret = d;
      }
    });

    return ret;
  }

  function calculateDurationComparisonScore(
    stashScene: GQL.SlimSceneDataFragment,
    scrapedScene: IScrapedScene
  ) {
    if (scrapedScene.fingerprints && scrapedScene.fingerprints.length > 0) {
      const durations = scrapedScene.fingerprints.map((f) => f.duration);
      const diffs = durations.map((d) => minDurationDiff(stashScene, d));
      const filteredDurations = diffs.filter((duration) => duration <= 5);

      const minDiff = Math.min(...diffs);

      return [
        filteredDurations.length,
        filteredDurations.length / durations.length,
        minDiff,
      ];
    }
    return [0, 0, 0];
  }

  function compareScenesForSort(
    stashScene: GQL.SlimSceneDataFragment,
    sceneA: IScrapedScene,
    sceneB: IScrapedScene
  ) {
    // Compare sceneA and sceneB to each other for sorting based on similarity to stashScene
    // Order of priority is: nb. phash match > nb. duration match > ratio duration match > ratio phash match

    // scenes without any fingerprints should be sorted to the end
    if (!sceneA.fingerprints?.length && sceneB.fingerprints?.length) {
      return 1;
    }
    if (!sceneB.fingerprints?.length && sceneA.fingerprints?.length) {
      return -1;
    }

    const [nbPhashMatchSceneA, ratioPhashMatchSceneA] =
      calculatePhashComparisonScore(stashScene, sceneA);
    const [nbPhashMatchSceneB, ratioPhashMatchSceneB] =
      calculatePhashComparisonScore(stashScene, sceneB);

    // If only one scene has matching phash, prefer that scene
    if (
      (nbPhashMatchSceneA != nbPhashMatchSceneB && nbPhashMatchSceneA === 0) ||
      nbPhashMatchSceneB === 0
    ) {
      return nbPhashMatchSceneB - nbPhashMatchSceneA;
    }

    // Prefer scene with highest ratio of phash matches
    if (ratioPhashMatchSceneA !== ratioPhashMatchSceneB) {
      return ratioPhashMatchSceneB - ratioPhashMatchSceneA;
    }

    // Same ratio of phash matches, check duration
    const [
      nbDurationMatchSceneA,
      ratioDurationMatchSceneA,
      minDurationDiffSceneA,
    ] = calculateDurationComparisonScore(stashScene, sceneA);
    const [
      nbDurationMatchSceneB,
      ratioDurationMatchSceneB,
      minDurationDiffSceneB,
    ] = calculateDurationComparisonScore(stashScene, sceneB);

    if (nbDurationMatchSceneA != nbDurationMatchSceneB) {
      return nbDurationMatchSceneB - nbDurationMatchSceneA;
    }

    // Same number of phash & duration, check duration ratio
    if (ratioDurationMatchSceneA != ratioDurationMatchSceneB) {
      return ratioDurationMatchSceneB - ratioDurationMatchSceneA;
    }

    // fall back to duration difference - less is better
    return minDurationDiffSceneA - minDurationDiffSceneB;
  }

  const [spriteImage, setSpriteImage] = useState<string | null>(null);
  const lightboxImage = useMemo(
    () => [{ paths: { thumbnail: spriteImage, image: spriteImage } }],
    [spriteImage]
  );
  const showLightbox = useLightbox({
    images: lightboxImage,
  });
  function showLightboxImage(imagePath: string) {
    setSpriteImage(imagePath);
    showLightbox();
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
      } else if (
        searchResult &&
        searchResult.results &&
        searchResult.results?.length >= 2
      ) {
        searchResult.results?.sort((scrapedSceneA, scrapedSceneB) =>
          compareScenesForSort(scene, scrapedSceneA, scrapedSceneB)
        );
      }

      return (
        <TaggerScene
          key={scene.id}
          loading={loading}
          scene={scene}
          url={sceneLink}
          errorMessage={errorMessage}
          doSceneQuery={
            currentSource?.supportSceneQuery
              ? async (v) => {
                  await doSceneQuery(scene.id, v);
                }
              : undefined
          }
          scrapeSceneFragment={
            currentSource?.supportSceneFragment
              ? async () => {
                  await doSceneFragmentScrape(scene.id);
                }
              : undefined
          }
          showLightboxImage={showLightboxImage}
        >
          {searchResult && searchResult.results?.length ? (
            <SceneSearchResults scenes={searchResult.results} target={scene} />
          ) : undefined}
        </TaggerScene>
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
    if (!currentSource?.supportSceneFragment) {
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
    <SceneTaggerModals>
      <div className="tagger-container mx-md-auto">
        <div className="tagger-container-header">
          <div className="d-flex justify-content-between align-items-center flex-wrap">
            <div className="w-auto">{renderSourceSelector()}</div>
            <div className="d-flex">
              {maybeRenderShowHideUnmatchedButton()}
              {maybeRenderSubmitFingerprintsButton()}
              {renderFragmentScrapeButton()}
              {renderConfigButton()}
            </div>
          </div>
          <Config show={showConfig} />
        </div>
        <div>{renderScenes()}</div>
      </div>
    </SceneTaggerModals>
  );
};
