import React, { useState, useEffect, useRef } from "react";
import {
  initialConfig,
  ITaggerConfig,
  LOCAL_FORAGE_KEY,
} from "src/components/Tagger/constants";
import * as GQL from "src/core/generated-graphql";
import {
  queryScrapeScene,
  queryScrapeSceneQuery,
  queryScrapeSceneQueryFragment,
  stashBoxSceneBatchQuery,
  useConfiguration,
  useListSceneScrapers,
  usePerformerCreate,
  useSceneUpdate,
  useStudioCreate,
  useTagCreate,
} from "src/core/StashService";
import { useLocalForage, useToast } from "src/hooks";
import { ITaggerSource, SCRAPER_PREFIX, STASH_BOX_PREFIX } from "./constants";

export interface ITaggerContextState {
  config: ITaggerConfig;
  setConfig: (c: ITaggerConfig) => void;
  loading: boolean;
  loadingMulti?: boolean;
  multiError?: string;
  sources: ITaggerSource[];
  currentSource?: ITaggerSource;
  searchResults: Record<string, ISceneQueryResult>;
  setCurrentSource: (src?: ITaggerSource) => void;
  doSceneQuery: (sceneID: string, searchStr: string) => Promise<void>;
  doSceneFragmentScrape: (sceneID: string) => Promise<void>;
  doMultiSceneFragmentScrape: (sceneIDs: string[]) => Promise<void>;
  stopMultiScrape: () => void;
  createNewTag: (toCreate: GQL.ScrapedTag) => Promise<string | undefined>;
  createNewPerformer: (
    toCreate: GQL.PerformerCreateInput
  ) => Promise<string | undefined>;
  createNewStudio: (
    toCreate: GQL.StudioCreateInput
  ) => Promise<string | undefined>;
  resolveScene: (
    sceneID: string,
    index: number,
    scene: IScrapedScene
  ) => Promise<void>;
  submitFingerprints: () => Promise<void>;
  pendingFingerprints: string[];
  saveScene: (sceneCreateInput: GQL.SceneUpdateInput) => Promise<void>;
}

const dummyFn = () => {
  return Promise.resolve();
};
const dummyValFn = () => {
  return Promise.resolve(undefined);
};

export const TaggerStateContext = React.createContext<ITaggerContextState>({
  config: initialConfig,
  setConfig: () => {},
  loading: false,
  sources: [],
  searchResults: {},
  setCurrentSource: () => {},
  doSceneQuery: dummyFn,
  doSceneFragmentScrape: dummyFn,
  doMultiSceneFragmentScrape: dummyFn,
  stopMultiScrape: () => {},
  createNewTag: dummyValFn,
  createNewPerformer: dummyValFn,
  createNewStudio: dummyValFn,
  resolveScene: dummyFn,
  submitFingerprints: dummyFn,
  pendingFingerprints: [],
  saveScene: dummyFn,
});

export type IScrapedScene = GQL.ScrapedScene & { resolved?: boolean };

export interface ISceneQueryResult {
  results?: IScrapedScene[];
  error?: string;
}

export const TaggerContext: React.FC = ({ children }) => {
  const [{ data: config }, setConfig] = useLocalForage<ITaggerConfig>(
    LOCAL_FORAGE_KEY,
    initialConfig
  );

  const [loading, setLoading] = useState(false);
  const [loadingMulti, setLoadingMulti] = useState(false);
  const [sources, setSources] = useState<ITaggerSource[]>([]);
  const [currentSource, setCurrentSource] = useState<ITaggerSource>();
  const [multiError, setMultiError] = useState<string | undefined>();
  const [searchResults, setSearchResults] = useState<
    Record<string, ISceneQueryResult>
  >({});

  const stopping = useRef(false);

  const stashConfig = useConfiguration();
  const Scrapers = useListSceneScrapers();

  const Toast = useToast();
  const [createTag] = useTagCreate();
  const [createPerformer] = usePerformerCreate();
  const [createStudio] = useStudioCreate();
  const [updateScene] = useSceneUpdate();

  useEffect(() => {
    if (!stashConfig.data || !Scrapers.data) {
      return;
    }

    const { stashBoxes } = stashConfig.data.configuration.general;
    const scrapers = Scrapers.data.listSceneScrapers;

    const stashboxSources: ITaggerSource[] = stashBoxes.map((s, i) => ({
      id: `${STASH_BOX_PREFIX}${i}`,
      stashboxEndpoint: s.endpoint,
      sourceInput: {
        stash_box_index: i,
      },
      displayName: `stash-box: ${s.name}`,
      supportSceneFragment: true,
      supportSceneQuery: true,
    }));

    // filter scraper sources such that only those that can query scrape or
    // scrape via fragment are added
    const scraperSources: ITaggerSource[] = scrapers
      .filter((s) =>
        s.scene?.supported_scrapes.some(
          (t) => t === GQL.ScrapeType.Name || t === GQL.ScrapeType.Fragment
        )
      )
      .map((s) => ({
        id: `${SCRAPER_PREFIX}${s.id}`,
        sourceInput: {
          scraper_id: s.id,
        },
        displayName: s.name,
        supportSceneQuery: s.scene?.supported_scrapes.includes(
          GQL.ScrapeType.Name
        ),
        supportSceneFragment: s.scene?.supported_scrapes.includes(
          GQL.ScrapeType.Fragment
        ),
      }));

    setSources(stashboxSources.concat(scraperSources));
  }, [Scrapers.data, stashConfig.data]);

  useEffect(() => {
    if (sources.length && !currentSource) {
      setCurrentSource(sources[0]);
    }
  }, [sources, currentSource]);

  useEffect(() => {
    setSearchResults({});
  }, [currentSource]);

  function getPendingFingerprints() {
    const endpoint = currentSource?.stashboxEndpoint;
    if (!config || !endpoint) return [];

    return config.fingerprintQueue[endpoint] ?? [];
  }

  function clearSubmissionQueue() {
    const endpoint = currentSource?.stashboxEndpoint;
    if (!config || !endpoint) return;

    setConfig({
      ...config,
      fingerprintQueue: {
        ...config.fingerprintQueue,
        [endpoint]: [],
      },
    });
  }

  const [
    submitFingerprintsMutation,
  ] = GQL.useSubmitStashBoxFingerprintsMutation();

  async function submitFingerprints() {
    const endpoint = currentSource?.stashboxEndpoint;
    const stashBoxIndex =
      currentSource?.sourceInput.stash_box_index ?? undefined;

    if (!config || !endpoint || stashBoxIndex === undefined) return;

    try {
      setLoading(true);
      await submitFingerprintsMutation({
        variables: {
          input: {
            stash_box_index: stashBoxIndex,
            scene_ids: config.fingerprintQueue[endpoint],
          },
        },
      });

      clearSubmissionQueue();
    } catch (err) {
      Toast.error(err);
    } finally {
      setLoading(false);
    }
  }

  function queueFingerprintSubmission(sceneId: string) {
    const endpoint = currentSource?.stashboxEndpoint;
    if (!config || !endpoint) return;

    setConfig({
      ...config,
      fingerprintQueue: {
        ...config.fingerprintQueue,
        [endpoint]: [...(config.fingerprintQueue[endpoint] ?? []), sceneId],
      },
    });
  }

  async function doSceneQuery(sceneID: string, searchVal: string) {
    if (!currentSource) {
      return;
    }

    try {
      setLoading(true);

      const results = await queryScrapeSceneQuery(
        currentSource.sourceInput,
        searchVal
      );
      let newResult: ISceneQueryResult;
      // scenes are already resolved if they come from stash-box
      const resolved = currentSource.sourceInput.stash_box_index !== undefined;

      if (results.error) {
        newResult = { error: results.error.message };
      } else if (results.errors) {
        newResult = { error: results.errors.toString() };
      } else {
        newResult = {
          results: results.data.scrapeSingleScene.map((r) => ({
            ...r,
            resolved,
          })),
        };
      }

      setSearchResults({ ...searchResults, [sceneID]: newResult });
    } catch (err) {
      Toast.error(err);
    } finally {
      setLoading(false);
    }
  }

  async function sceneFragmentScrape(sceneID: string) {
    if (!currentSource) {
      return;
    }

    const results = await queryScrapeScene(currentSource.sourceInput, sceneID);
    let newResult: ISceneQueryResult;
    // scenes are already resolved if they come from stash-box
    const resolved = currentSource.sourceInput.stash_box_index !== undefined;

    if (results.error) {
      newResult = { error: results.error.message };
    } else if (results.errors) {
      newResult = { error: results.errors.toString() };
    } else {
      newResult = {
        results: results.data.scrapeSingleScene.map((r) => ({
          ...r,
          resolved,
        })),
      };
    }

    setSearchResults((current) => {
      return { ...current, [sceneID]: newResult };
    });
  }

  async function doSceneFragmentScrape(sceneID: string) {
    if (!currentSource) {
      return;
    }

    setSearchResults((current) => {
      const newResults = { ...current };
      delete newResults[sceneID];
      return newResults;
    });

    try {
      setLoading(true);
      await sceneFragmentScrape(sceneID);
    } finally {
      setLoading(false);
    }
  }

  async function doMultiSceneFragmentScrape(sceneIDs: string[]) {
    if (!currentSource) {
      return;
    }

    setSearchResults({});

    try {
      stopping.current = false;
      setLoading(true);
      setMultiError(undefined);

      const stashBoxIndex =
        currentSource.sourceInput.stash_box_index ?? undefined;

      // if current source is stash-box, we can use the multi-scene
      // interface
      if (stashBoxIndex !== undefined) {
        const results = await stashBoxSceneBatchQuery(sceneIDs, stashBoxIndex);

        if (results.error) {
          setMultiError(results.error.message);
        } else if (results.errors) {
          setMultiError(results.errors.toString());
        } else {
          const newSearchResults = { ...searchResults };
          sceneIDs.forEach((sceneID, index) => {
            const newResults = results.data.scrapeMultiScenes[index].map(
              (r) => ({
                ...r,
                resolved: true,
              })
            );

            newSearchResults[sceneID] = {
              results: newResults,
            };
          });

          setSearchResults(newSearchResults);
        }
      } else {
        setLoadingMulti(true);

        // do singular calls
        await sceneIDs.reduce(async (promise, id) => {
          await promise;
          if (!stopping.current) {
            await sceneFragmentScrape(id);
          }
        }, Promise.resolve());
      }
    } finally {
      setLoading(false);
      setLoadingMulti(false);
    }
  }

  function stopMultiScrape() {
    stopping.current = true;
  }

  async function resolveScene(
    sceneID: string,
    index: number,
    scene: IScrapedScene
  ) {
    if (!currentSource || scene.resolved || !searchResults[sceneID].results) {
      return Promise.resolve();
    }

    try {
      const sceneInput: GQL.ScrapedSceneInput = {
        date: scene.date,
        details: scene.details,
        remote_site_id: scene.remote_site_id,
        title: scene.title,
        url: scene.url,
      };

      const result = await queryScrapeSceneQueryFragment(
        currentSource.sourceInput,
        sceneInput
      );

      if (result.data.scrapeSingleScene.length) {
        const resolvedScene = result.data.scrapeSingleScene[0];

        // set the scene in the results and mark as resolved
        const newResult = [...searchResults[sceneID].results!];
        newResult[index] = { ...resolvedScene, resolved: true };
        setSearchResults({
          ...searchResults,
          [sceneID]: { ...searchResults[sceneID], results: newResult },
        });
      }
    } catch (err) {
      Toast.error(err);

      const newResult = [...searchResults[sceneID].results!];
      newResult[index] = { ...newResult[index], resolved: true };
      setSearchResults({
        ...searchResults,
        [sceneID]: { ...searchResults[sceneID], results: newResult },
      });
    }
  }

  function clearSearchResults(sceneID: string) {
    setSearchResults((current) => {
      const newSearchResults = { ...current };
      delete newSearchResults[sceneID];
      return newSearchResults;
    });
  }

  async function saveScene(sceneCreateInput: GQL.SceneUpdateInput) {
    try {
      await updateScene({
        variables: {
          input: sceneCreateInput,
        },
      });

      queueFingerprintSubmission(sceneCreateInput.id);
      clearSearchResults(sceneCreateInput.id);
    } catch (err) {
      Toast.error(err);
    } finally {
      setLoading(false);
    }
  }

  function mapResults(fn: (r: IScrapedScene) => IScrapedScene) {
    const newSearchResults = { ...searchResults };

    Object.keys(newSearchResults).forEach((k) => {
      const searchResult = searchResults[k];
      if (!searchResult.results) {
        return;
      }

      newSearchResults[k].results = searchResult.results.map(fn);
    });

    return newSearchResults;
  }

  async function createNewTag(toCreate: GQL.ScrapedTag) {
    const tagInput: GQL.TagCreateInput = { name: toCreate.name ?? "" };
    try {
      const result = await createTag({
        variables: {
          input: tagInput,
        },
      });

      const tagID = result.data?.tagCreate?.id;

      const newSearchResults = mapResults((r) => {
        if (!r.tags) {
          return r;
        }

        return {
          ...r,
          tags: r.tags.map((t) => {
            if (t.name === toCreate.name) {
              return {
                ...t,
                stored_id: tagID,
              };
            }

            return t;
          }),
        };
      });

      setSearchResults(newSearchResults);

      Toast.success({
        content: (
          <span>
            Created tag: <b>{toCreate.name}</b>
          </span>
        ),
      });

      return tagID;
    } catch (e) {
      Toast.error(e);
    }
  }

  async function createNewPerformer(toCreate: GQL.PerformerCreateInput) {
    try {
      const result = await createPerformer({
        variables: {
          input: toCreate,
        },
      });

      const performerID = result.data?.performerCreate?.id;

      const newSearchResults = mapResults((r) => {
        if (!r.performers) {
          return r;
        }

        return {
          ...r,
          performers: r.performers.map((t) => {
            if (t.name === toCreate.name) {
              return {
                ...t,
                stored_id: performerID,
              };
            }

            return t;
          }),
        };
      });

      setSearchResults(newSearchResults);

      Toast.success({
        content: (
          <span>
            Created performer: <b>{toCreate.name}</b>
          </span>
        ),
      });

      return performerID;
    } catch (e) {
      Toast.error(e);
    }
  }

  async function createNewStudio(toCreate: GQL.StudioCreateInput) {
    try {
      const result = await createStudio({
        variables: {
          input: toCreate,
        },
      });

      const studioID = result.data?.studioCreate?.id;

      const newSearchResults = mapResults((r) => {
        if (!r.studio) {
          return r;
        }

        return {
          ...r,
          studio:
            r.studio.name === toCreate.name
              ? {
                  ...r.studio,
                  stored_id: studioID,
                }
              : r.studio,
        };
      });

      setSearchResults(newSearchResults);

      Toast.success({
        content: (
          <span>
            Created studio: <b>{toCreate.name}</b>
          </span>
        ),
      });

      return studioID;
    } catch (e) {
      Toast.error(e);
    }
  }

  return (
    <TaggerStateContext.Provider
      value={{
        config: config ?? initialConfig,
        setConfig,
        loading: loading || loadingMulti,
        loadingMulti,
        multiError,
        sources,
        currentSource,
        searchResults,
        setCurrentSource: (src) => {
          setCurrentSource(src);
        },
        doSceneQuery,
        doSceneFragmentScrape,
        doMultiSceneFragmentScrape,
        stopMultiScrape,
        createNewTag,
        createNewPerformer,
        createNewStudio,
        resolveScene,
        saveScene,
        submitFingerprints,
        pendingFingerprints: getPendingFingerprints(),
      }}
    >
      {children}
    </TaggerStateContext.Provider>
  );
};
