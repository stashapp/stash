import React, { useState, useEffect, useRef } from "react";
import { initialConfig, ITaggerConfig } from "src/components/Tagger/constants";
import * as GQL from "src/core/generated-graphql";
import {
  queryFindPerformer,
  queryFindStudio,
  queryScrapeScene,
  queryScrapeSceneQuery,
  queryScrapeSceneQueryFragment,
  stashBoxSceneBatchQuery,
  useListSceneScrapers,
  usePerformerCreate,
  usePerformerUpdate,
  useSceneUpdate,
  useStudioCreate,
  useStudioUpdate,
  useTagCreate,
  useTagUpdate,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import { useConfigurationContext } from "src/hooks/Config";
import { ITaggerSource, SCRAPER_PREFIX, STASH_BOX_PREFIX } from "./constants";
import { errorToString } from "src/utils";
import { mergeStudioStashIDs } from "./utils";
import { useTaggerConfig } from "./config";

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
  createNewTag: (
    tag: GQL.ScrapedTag,
    toCreate: GQL.TagCreateInput
  ) => Promise<string | undefined>;
  createNewPerformer: (
    performer: GQL.ScrapedPerformer,
    toCreate: GQL.PerformerCreateInput
  ) => Promise<string | undefined>;
  linkPerformer: (
    performer: GQL.ScrapedPerformer,
    performerID: string
  ) => Promise<void>;
  createNewStudio: (
    studio: GQL.ScrapedStudio,
    toCreate: GQL.StudioCreateInput
  ) => Promise<string | undefined>;
  updateStudio: (studio: GQL.StudioUpdateInput) => Promise<void>;
  linkStudio: (studio: GQL.ScrapedStudio, studioID: string) => Promise<void>;
  updateTag: (
    tag: GQL.ScrapedTag,
    updateInput: GQL.TagUpdateInput
  ) => Promise<void>;
  resolveScene: (
    sceneID: string,
    index: number,
    scene: IScrapedScene
  ) => Promise<void>;
  submitFingerprints: () => Promise<void>;
  pendingFingerprints: string[];
  saveScene: (
    sceneCreateInput: GQL.SceneUpdateInput,
    queueFingerprint: boolean
  ) => Promise<void>;
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
  linkPerformer: dummyFn,
  createNewStudio: dummyValFn,
  updateStudio: dummyFn,
  linkStudio: dummyFn,
  updateTag: dummyFn,
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
  const [loading, setLoading] = useState(false);
  const [loadingMulti, setLoadingMulti] = useState(false);
  const [sources, setSources] = useState<ITaggerSource[]>([]);
  const [currentSource, setCurrentSource] = useState<ITaggerSource>();
  const [multiError, setMultiError] = useState<string | undefined>();
  const [searchResults, setSearchResults] = useState<
    Record<string, ISceneQueryResult>
  >({});

  const stopping = useRef(false);

  const { configuration: stashConfig } = useConfigurationContext();
  const { config, setConfig } = useTaggerConfig();

  const Scrapers = useListSceneScrapers();

  const Toast = useToast();
  const [createTag] = useTagCreate();
  const [createPerformer] = usePerformerCreate();
  const [updatePerformer] = usePerformerUpdate();
  const [createStudio] = useStudioCreate();
  const [updateStudio] = useStudioUpdate();
  const [updateScene] = useSceneUpdate();
  const [updateTag] = useTagUpdate();

  useEffect(() => {
    if (!stashConfig || !Scrapers.data) {
      return;
    }

    const { stashBoxes } = stashConfig.general;
    const scrapers = Scrapers.data.listScrapers;

    const stashboxSources: ITaggerSource[] = stashBoxes.map((s, i) => ({
      id: `${STASH_BOX_PREFIX}${s.endpoint}`,
      sourceInput: {
        stash_box_endpoint: s.endpoint,
      },
      displayName: `stash-box: ${s.name || `#${i + 1}`}`,
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
  }, [Scrapers.data, stashConfig]);

  // set the current source on load
  useEffect(() => {
    if (!sources.length || currentSource) {
      return;
    }
    // First, see if we have a saved endpoint.
    if (config.selectedEndpoint) {
      let source = sources.find(
        (s) => s.sourceInput.stash_box_endpoint == config.selectedEndpoint
      );
      if (source) {
        setCurrentSource(source);
        return;
      }
    }
    // Otherwise, just use the first source.
    setCurrentSource(sources[0]);
  }, [sources, currentSource, config]);

  // clear the search results when the source changes
  useEffect(() => {
    setSearchResults({});
  }, [currentSource]);

  // keep selected endpoint in config in sync with current source
  useEffect(() => {
    const selectedEndpoint = currentSource?.sourceInput.stash_box_endpoint;
    if (selectedEndpoint && selectedEndpoint !== config.selectedEndpoint) {
      setConfig({
        ...config,
        selectedEndpoint,
      });
    }
  }, [currentSource, config, setConfig]);

  function getPendingFingerprints() {
    const endpoint = currentSource?.sourceInput.stash_box_endpoint;
    if (!config || !endpoint) return [];

    return config.fingerprintQueue[endpoint] ?? [];
  }

  function clearSubmissionQueue() {
    const endpoint = currentSource?.sourceInput.stash_box_endpoint;
    if (!config || !endpoint) return;

    setConfig({
      ...config,
      fingerprintQueue: {
        ...config.fingerprintQueue,
        [endpoint]: [],
      },
    });
  }

  const [submitFingerprintsMutation] =
    GQL.useSubmitStashBoxFingerprintsMutation();

  async function submitFingerprints() {
    const endpoint = currentSource?.sourceInput.stash_box_endpoint;

    if (!config || !endpoint) return;

    try {
      setLoading(true);
      await submitFingerprintsMutation({
        variables: {
          input: {
            stash_box_endpoint: endpoint,
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
    const endpoint = currentSource?.sourceInput.stash_box_endpoint;
    if (!config || !endpoint) return;

    setConfig({
      ...config,
      fingerprintQueue: {
        ...config.fingerprintQueue,
        [endpoint]: [...(config.fingerprintQueue[endpoint] ?? []), sceneId],
      },
    });
  }

  function clearSearchResults(sceneID: string) {
    setSearchResults((current) => {
      const newSearchResults = { ...current };
      delete newSearchResults[sceneID];
      return newSearchResults;
    });
  }

  async function doSceneQuery(sceneID: string, searchVal: string) {
    if (!currentSource) {
      return;
    }

    try {
      setLoading(true);
      clearSearchResults(sceneID);

      const results = await queryScrapeSceneQuery(
        currentSource.sourceInput,
        searchVal
      );
      let newResult: ISceneQueryResult;
      // scenes are already resolved if they come from stash-box
      const resolved =
        currentSource.sourceInput.stash_box_endpoint !== undefined;

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

    clearSearchResults(sceneID);

    let newResult: ISceneQueryResult;

    try {
      const results = await queryScrapeScene(
        currentSource.sourceInput,
        sceneID
      );

      if (results.error) {
        newResult = { error: results.error.message };
      } else if (results.errors) {
        newResult = { error: results.errors.toString() };
      } else {
        newResult = {
          results: results.data.scrapeSingleScene.map((r) => ({
            ...r,
            // scenes are already resolved if they are scraped via fragment
            resolved: true,
          })),
        };
      }
    } catch (err: unknown) {
      newResult = { error: errorToString(err) };
    }

    setSearchResults((current) => {
      return { ...current, [sceneID]: newResult };
    });
  }

  async function doSceneFragmentScrape(sceneID: string) {
    if (!currentSource) {
      return;
    }

    clearSearchResults(sceneID);

    try {
      setLoading(true);
      await sceneFragmentScrape(sceneID);
    } catch (err) {
      Toast.error(err);
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

      const stashBoxEndpoint =
        currentSource.sourceInput.stash_box_endpoint ?? undefined;

      // if current source is stash-box, we can use the multi-scene
      // interface
      if (stashBoxEndpoint !== undefined) {
        const results = await stashBoxSceneBatchQuery(
          sceneIDs,
          stashBoxEndpoint
        );

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
    } catch (err) {
      Toast.error(err);
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
        urls: scene.urls,
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

  async function saveScene(
    sceneCreateInput: GQL.SceneUpdateInput,
    queueFingerprint: boolean
  ) {
    try {
      await updateScene({
        variables: {
          input: {
            ...sceneCreateInput,
            // only set organized if it is enabled in the config
            organized: config?.markSceneAsOrganizedOnSave || undefined,
          },
        },
      });

      if (queueFingerprint) {
        queueFingerprintSubmission(sceneCreateInput.id);
      }
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

  async function createNewTag(
    tag: GQL.ScrapedTag,
    toCreate: GQL.TagCreateInput
  ) {
    try {
      const result = await createTag({
        variables: {
          input: toCreate,
        },
      });

      const tagID = result.data?.tagCreate?.id;
      if (tagID === undefined) return undefined;

      const newSearchResults = mapResults((r) => {
        if (!r.tags) {
          return r;
        }

        return {
          ...r,
          tags: r.tags.map((t) => {
            if (t.name === tag.name) {
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

      Toast.success(
        <span>
          Created tag: <b>{toCreate.name}</b>
        </span>
      );

      return tagID;
    } catch (e) {
      Toast.error(e);
    }
  }

  async function createNewPerformer(
    performer: GQL.ScrapedPerformer,
    toCreate: GQL.PerformerCreateInput
  ) {
    try {
      const result = await createPerformer({
        variables: {
          input: toCreate,
        },
      });

      const performerID = result.data?.performerCreate?.id;
      if (performerID === undefined) return undefined;

      const newSearchResults = mapResults((r) => {
        if (!r.performers) {
          return r;
        }

        return {
          ...r,
          performers: r.performers.map((p) => {
            // Match by remote_site_id if available, otherwise fall back to name
            const matches = performer.remote_site_id
              ? p.remote_site_id === performer.remote_site_id
              : p.name === performer.name;

            if (matches) {
              return {
                ...p,
                stored_id: performerID,
              };
            }

            return p;
          }),
        };
      });

      setSearchResults(newSearchResults);

      Toast.success(
        <span>
          Created performer: <b>{toCreate.name}</b>
        </span>
      );

      return performerID;
    } catch (e) {
      Toast.error(e);
    }
  }

  async function linkPerformer(
    performer: GQL.ScrapedPerformer,
    performerID: string
  ) {
    if (
      !performer.remote_site_id ||
      !currentSource?.sourceInput.stash_box_endpoint
    )
      return;

    try {
      const queryResult = await queryFindPerformer(performerID);
      if (queryResult.data.findPerformer) {
        const target = queryResult.data.findPerformer;

        const stashIDs: GQL.StashIdInput[] = target.stash_ids.map((e) => {
          return {
            endpoint: e.endpoint,
            stash_id: e.stash_id,
            updated_at: e.updated_at,
          };
        });

        stashIDs.push({
          stash_id: performer.remote_site_id,
          endpoint: currentSource?.sourceInput.stash_box_endpoint,
          updated_at: new Date().toISOString(),
        });

        await updatePerformer({
          variables: {
            input: {
              id: performerID,
              stash_ids: stashIDs,
            },
          },
        });

        const newSearchResults = mapResults((r) => {
          if (!r.performers) {
            return r;
          }

          return {
            ...r,
            performers: r.performers.map((p) => {
              if (p.remote_site_id === performer.remote_site_id) {
                return {
                  ...p,
                  stored_id: performerID,
                };
              }

              return p;
            }),
          };
        });

        setSearchResults(newSearchResults);

        Toast.success(<span>Added stash-id to performer</span>);
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  async function createNewStudio(
    studio: GQL.ScrapedStudio,
    toCreate: GQL.StudioCreateInput
  ) {
    try {
      const result = await createStudio({
        variables: {
          input: toCreate,
        },
      });

      const studioID = result.data?.studioCreate?.id;
      if (studioID === undefined) return undefined;

      const newSearchResults = mapResults((r) => {
        if (!r.studio) {
          return r;
        }

        let resultStudio = r.studio;
        if (resultStudio.name === studio.name) {
          resultStudio = {
            ...resultStudio,
            stored_id: studioID,
          };
        }

        // #5821 - set the stored_id of the parent studio if it matches too
        if (resultStudio.parent?.name === studio.name) {
          resultStudio = {
            ...resultStudio,
            parent: {
              ...resultStudio.parent,
              stored_id: studioID,
            },
          };
        }

        return {
          ...r,
          studio: resultStudio,
        };
      });

      setSearchResults(newSearchResults);

      Toast.success(
        <span>
          Created studio: <b>{toCreate.name}</b>
        </span>
      );

      return studioID;
    } catch (e) {
      Toast.error(e);
    }
  }

  async function updateExistingStudio(input: GQL.StudioUpdateInput) {
    try {
      const inputCopy = { ...input };
      inputCopy.stash_ids = await mergeStudioStashIDs(
        input.id,
        input.stash_ids ?? []
      );
      const result = await updateStudio({
        variables: {
          input: input,
        },
      });

      const studioID = result.data?.studioUpdate?.id;

      const stashID = input.stash_ids?.find((e) => {
        return e.endpoint === currentSource?.sourceInput.stash_box_endpoint;
      })?.stash_id;

      if (stashID) {
        const newSearchResults = mapResults((r) => {
          if (!r.studio) {
            return r;
          }

          return {
            ...r,
            studio:
              r.remote_site_id === stashID
                ? {
                    ...r.studio,
                    stored_id: studioID,
                  }
                : r.studio,
          };
        });

        setSearchResults(newSearchResults);
      }

      Toast.success(
        <span>
          Created studio: <b>{input.name}</b>
        </span>
      );
    } catch (e) {
      Toast.error(e);
    }
  }

  async function linkStudio(studio: GQL.ScrapedStudio, studioID: string) {
    if (
      !studio.remote_site_id ||
      !currentSource?.sourceInput.stash_box_endpoint
    )
      return;

    try {
      const queryResult = await queryFindStudio(studioID);
      if (queryResult.data.findStudio) {
        const target = queryResult.data.findStudio;

        const stashIDs: GQL.StashIdInput[] = target.stash_ids.map((e) => {
          return {
            endpoint: e.endpoint,
            stash_id: e.stash_id,
            updated_at: e.updated_at,
          };
        });

        stashIDs.push({
          stash_id: studio.remote_site_id,
          endpoint: currentSource?.sourceInput.stash_box_endpoint,
          updated_at: new Date().toISOString(),
        });

        await updateStudio({
          variables: {
            input: {
              id: studioID,
              stash_ids: stashIDs,
            },
          },
        });

        const newSearchResults = mapResults((r) => {
          if (!r.studio) {
            return r;
          }

          return {
            ...r,
            studio:
              r.studio.remote_site_id === studio.remote_site_id
                ? {
                    ...r.studio,
                    stored_id: studioID,
                  }
                : r.studio,
          };
        });

        setSearchResults(newSearchResults);

        Toast.success(<span>Added stash-id to studio</span>);
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  async function updateExistingTag(
    tag: GQL.ScrapedTag,
    updateInput: GQL.TagUpdateInput
  ) {
    const hasRemoteID = !!tag.remote_site_id;

    try {
      await updateTag({
        variables: {
          input: updateInput,
        },
      });

      const newSearchResults = mapResults((r) => {
        if (!r.tags) {
          return r;
        }

        return {
          ...r,
          tags: r.tags.map((t) => {
            if (
              (hasRemoteID && t.remote_site_id === tag.remote_site_id) ||
              (!hasRemoteID && t.name === tag.name)
            ) {
              return {
                ...t,
                stored_id: updateInput.id,
              };
            }

            return t;
          }),
        };
      });

      setSearchResults(newSearchResults);

      Toast.success(<span>Updated tag</span>);
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
        linkPerformer,
        createNewStudio,
        updateStudio: updateExistingStudio,
        linkStudio,
        updateTag: updateExistingTag,
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
