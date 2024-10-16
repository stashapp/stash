import React, { useContext, useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneQueue } from "src/models/sceneQueue";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";

import { Icon } from "src/components/Shared/Icon";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { OperationButton } from "src/components/Shared/OperationButton";
import { ISceneQueryResult, TaggerStateContext } from "../context";
import Config from "./Config";
import { TaggerScene } from "./TaggerScene";
import { SceneTaggerModals } from "./sceneTaggerModals";
import { SceneSearchResults } from "./StashSearchResult";
import { ConfigurationContext } from "src/hooks/Config";
import { faCog } from "@fortawesome/free-solid-svg-icons";
import { useLightbox } from "src/hooks/Lightbox/hooks";

const Scene: React.FC<{
  scene: GQL.SlimSceneDataFragment;
  searchResult?: ISceneQueryResult;
  queue?: SceneQueue;
  index: number;
  showLightboxImage: (imagePath: string) => void;
}> = ({ scene, searchResult, queue, index, showLightboxImage }) => {
  const intl = useIntl();
  const { currentSource, doSceneQuery, doSceneFragmentScrape, loading } =
    useContext(TaggerStateContext);
  const { configuration } = React.useContext(ConfigurationContext);

  const cont = configuration?.interface.continuePlaylistDefault ?? false;

  const sceneLink = useMemo(
    () =>
      queue
        ? queue.makeLink(scene.id, { sceneIndex: index, continue: cont })
        : `/scenes/${scene.id}`,
    [queue, scene.id, index, cont]
  );

  const errorMessage = useMemo(() => {
    if (searchResult?.error) {
      return searchResult.error;
    } else if (searchResult && searchResult.results?.length === 0) {
      return intl.formatMessage({
        id: "component_tagger.results.match_failed_no_result",
      });
    }
  }, [intl, searchResult]);

  return (
    <TaggerScene
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
};

interface ITaggerProps {
  scenes: GQL.SlimSceneDataFragment[];
  queue?: SceneQueue;
}

export const Tagger: React.FC<ITaggerProps> = ({ scenes, queue }) => {
  const {
    sources,
    setCurrentSource,
    currentSource,
    doMultiSceneFragmentScrape,
    stopMultiScrape,
    searchResults,
    loading,
    loadingMulti,
    multiError,
    submitFingerprints,
    pendingFingerprints,
  } = useContext(TaggerStateContext);
  const [showConfig, setShowConfig] = useState(false);
  const [hideUnmatched, setHideUnmatched] = useState(false);

  const intl = useIntl();

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

  const filteredScenes = useMemo(
    () =>
      !hideUnmatched
        ? scenes
        : scenes.filter((s) => searchResults[s.id]?.results?.length),
    [scenes, searchResults, hideUnmatched]
  );

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

  function maybeRenderSearchAllButton() {
    const searchButtons = document.querySelectorAll('.input-group-append .btn.btn-primary');
  
    if (searchButtons.length > 0) {
      return (
        <Button onClick={handleSearchAllClick}>
          <FormattedMessage id="component_tagger.verb_search_all" />
        </Button>
      );
    }
    return null;
  }

  function handleSearchAllClick() {
    const searchButtons = document.querySelectorAll('.input-group-append .btn.btn-primary');
  
    searchButtons.forEach((button, index) => {
      setTimeout(() => {
        console.log(`Clicking button ${index + 1}:`, button);
        
        if (button instanceof HTMLElement) {
          button.click();
        }
      }, index * 700); // 700ms delay between clicks... its a lot but if it goes too quick then the UI breaks and skips scenes.
    });
  }

  function maybeRenderCreateAllButton() {
    const CreateButtons = document.querySelectorAll('.row.no-gutters.align-items-center.mt-2 .btn-group .btn.btn-secondary');
  
    if (CreateButtons.length > 0) {
      return (
        <Button onClick={handleCreateAllClick}>
          <FormattedMessage id="component_tagger.verb_create_all" />
        </Button>
      );
    }
    return null;
  }

  function handleCreateAllClick() {
    const createButtons = document.querySelectorAll('.row.no-gutters.align-items-center.mt-2 .btn-group .btn.btn-secondary');
  
    createButtons.forEach((button, index) => {
      setTimeout(() => {
        console.log(`Clicking create button ${index + 1}:`, button);
        
        if (button instanceof HTMLElement) {
          button.click();

          setTimeout(() => {
            const saveButton = document.querySelector('.ModalFooter.modal-footer .ml-2.btn.btn-primary');
            if (saveButton instanceof HTMLElement) {
              saveButton.click();
            } else {
              console.warn(`Save button not found for create button ${index + 1}`);
            }
          }, 500); 
        }
      }, index * 800);
    });
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
              {maybeRenderCreateAllButton()}
              {maybeRenderSearchAllButton()}
              {maybeRenderShowHideUnmatchedButton()}
              {maybeRenderSubmitFingerprintsButton()}
              {renderFragmentScrapeButton()}
              {renderConfigButton()}
            </div>
          </div>
          <Config show={showConfig} />
        </div>
        <div>
          {filteredScenes.map((s, i) => (
            <Scene
              key={s.id}
              scene={s}
              searchResult={searchResults[s.id]}
              index={i}
              showLightboxImage={showLightboxImage}
              queue={queue}
            />
          ))}
        </div>
      </div>
    </SceneTaggerModals>
  );
};
