import React, { useEffect, useState } from "react";
import {
  Badge,
  Button,
  Card,
  Collapse,
  Form,
  InputGroup,
} from "react-bootstrap";
import localForage from "localforage";

import { FingerprintAlgorithm } from "src/definitions-box/globalTypes";
import * as GQL from "src/core/generated-graphql";
import { Icon, LoadingIndicator } from "src/components/Shared";
import { stashBoxQuery, useConfiguration } from "src/core/StashService";

import {
  FindSceneByFingerprintVariables,
  FindSceneByFingerprint,
  FindSceneByFingerprint_findSceneByFingerprint as FingerprintResult,
} from "src/definitions-box/FindSceneByFingerprint";
import { Me } from "src/definitions-box/Me";
import { loader } from "graphql.macro";
import StashSearchResult from "./StashSearchResult";
import { useStashBoxClient } from "./client";
import { parsePath, selectScenes, IStashBoxScene } from "./utils";

const FindSceneByFingerprintQuery = loader("src/queries/searchFingerprint.gql");
const MeQuery = loader("src/queries/me.gql");

const DEFAULT_BLACKLIST = [
  "\\sXXX\\s",
  "1080p",
  "720p",
  "2160p",
  "KTR",
  "RARBG",
  "com",
  "\\[",
  "\\]",
];
const dateRegex = /\.(\d\d)\.(\d\d)\.(\d\d)\./;
function prepareQueryString(
  scene: Partial<GQL.SlimSceneDataFragment>,
  paths: string[],
  filename: string,
  mode: ParseMode,
  blacklist: string[]
) {
  if ((mode === "auto" && scene.date && scene.studio) || mode === "metadata") {
    let str = [
      scene.date,
      scene.studio?.name ?? "",
      (scene?.performers ?? []).map((p) => p.name).join(" "),
      scene?.title ? scene.title.replace(/[^a-zA-Z0-9 ]+/g, "") : "",
    ]
      .filter((s) => s !== "")
      .join(" ");
    blacklist.forEach((b) => {
      str = str.replace(new RegExp(b, "gi"), " ");
    });
    return str;
  }
  let s = "";
  if (mode === "auto" || mode === "filename") {
    s = filename;
  } else if (mode === "path") {
    s = [...paths, filename].join(" ");
  } else {
    s = paths[paths.length - 1];
  }
  blacklist.forEach((b) => {
    s = s.replace(new RegExp(b, "i"), "");
  });
  const date = s.match(dateRegex);
  s = s.replace(/-/g, " ");
  if (date) {
    s = s.replace(date[0], ` 20${date[1]}-${date[2]}-${date[3]} `);
  }
  return s.replace(/\./g, " ");
}

type ParseMode = "auto" | "filename" | "dir" | "path" | "metadata";
const ModeDesc = {
  auto: "Uses metadata if present, or filename",
  metadata: "Only uses metadata",
  filename: "Only uses filename",
  dir: "Only uses parent directory of video file",
  path: "Uses entire file path",
};

interface ITaggerConfig {
  blacklist: string[];
  showMales: boolean;
  mode: ParseMode;
  setCoverImage: boolean;
  setTags: boolean;
  tagOperation: string;
  selectedEndpoint?: string;
}

interface ITaggerProps {
  scenes: GQL.SlimSceneDataFragment[];
}

export const Tagger: React.FC<ITaggerProps> = ({ scenes }) => {
  const stashConfig = useConfiguration();
  const [loading, setLoading] = useState(false);
  const [searchResults, setSearchResults] = useState<
    Record<string, IStashBoxScene[]>
  >({});
  const [queryString, setQueryString] = useState<Record<string, string>>({});
  const [selectedResult, setSelectedResult] = useState<
    Record<string, number>
  >();
  const [blacklistInput, setBlacklistInput] = useState<string>("");
  const [taggedScenes, setTaggedScenes] = useState<
    Record<string, Partial<GQL.SlimSceneDataFragment>>
  >({});
  const [fingerprints, setFingerprints] = useState<
    Record<string, FingerprintResult | null>
  >({});
  const [loadingFingerprints, setLoadingFingerprints] = useState(false);
  const [showConfig, setShowConfig] = useState(false);
  const [user, setUser] = useState<Me | null | undefined>();
  const [credentials, setCredentials] = useState({ endpoint: "", api_key: "" });
  const authFailure = user === undefined;

  const [config, setConfig] = useState<ITaggerConfig>({
    blacklist: DEFAULT_BLACKLIST,
    showMales: false,
    mode: "auto",
    setCoverImage: true,
    setTags: false,
    tagOperation: "merge",
  });

  useEffect(() => {
    localForage.getItem<ITaggerConfig>("tagger").then((data) => {
      setConfig({
        blacklist: data?.blacklist ?? DEFAULT_BLACKLIST,
        showMales: data?.showMales ?? false,
        mode: data?.mode ?? "auto",
        setCoverImage: data?.setCoverImage ?? true,
        setTags: data?.setTags ?? false,
        tagOperation: data?.tagOperation ?? "merge",
        selectedEndpoint: data?.selectedEndpoint,
      });
    });
  }, []);

  useEffect(() => {
    localForage.setItem("tagger", config);
  }, [config]);

  useEffect(() => {
    if (!stashConfig.data?.configuration.general) return;
    const selectedEndpoint = stashConfig.data.configuration.general.stashBoxes.find(
      (i) => i.endpoint === config.selectedEndpoint
    );
    if (selectedEndpoint) {
      setCredentials(selectedEndpoint);
    } else {
      setCredentials(stashConfig.data.configuration.general.stashBoxes[0]);
    }
  }, [stashConfig, config]);

  const handleInstanceSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedEndpoint = e.currentTarget.value;
    const creds = stashConfig.data?.configuration.general.stashBoxes.find(
      (i) => i.endpoint === selectedEndpoint
    );
    if (creds) {
      setCredentials(creds);
      setConfig({
        ...config,
        selectedEndpoint,
      });
    }
  };

  const client = useStashBoxClient(credentials?.endpoint, credentials?.api_key);
  useEffect(() => {
    if (client)
      client
        .query<Me>({
          query: MeQuery,
          errorPolicy: "ignore",
          fetchPolicy: "no-cache",
        })
        .then((result) => setUser(result.data))
        .catch(() => setUser(null));
  }, [client]);


  const selectedEndpointIndex = stashConfig.data?.configuration.general.stashBoxes.findIndex(s => s.endpoint === credentials.endpoint);

  const doBoxSearch = (sceneID: string, searchVal: string) => {
    if (selectedEndpointIndex === undefined || selectedEndpointIndex === -1) return;

		stashBoxQuery(searchVal, selectedEndpointIndex)
      .then((queryData) => {
        const s = selectScenes(queryData.data.queryStashBoxScene);
        setSearchResults({
          ...searchResults,
          [sceneID]: s
        });
        setLoading(false);
      });

    setLoading(true);
  };

  const handleTaggedScene = (scene: Partial<GQL.SlimSceneDataFragment>) => {
    setTaggedScenes({
      ...taggedScenes,
      [scene.id as string]: scene,
    });
  };

  const removeBlacklist = (index: number) => {
    setConfig({
      ...config,
      blacklist: [
        ...config.blacklist.slice(0, index),
        ...config.blacklist.slice(index + 1),
      ],
    });
  };

  const handleBlacklistAddition = () => {
    setConfig({
      ...config,
      blacklist: [...config.blacklist, blacklistInput],
    });
    setBlacklistInput("");
  };

  const handleFingerprintSearch = async () => {
    setLoadingFingerprints(true);
    const newFingerprints = { ...fingerprints };

    await Promise.all(
      scenes
        .filter(
          (s) => fingerprints[s.id] === undefined && s.stash_ids.length === 0
        )
        .map((s) => {
          let hash: string;
          let algorithm: FingerprintAlgorithm;
          if (s.oshash) {
            hash = s.oshash;
            algorithm = FingerprintAlgorithm.OSHASH;
          } else if (s.checksum) {
            hash = s.checksum;
            algorithm = FingerprintAlgorithm.MD5;
          } else {
            return null;
          }
          return client
            ?.query<FindSceneByFingerprint, FindSceneByFingerprintVariables>({
              query: FindSceneByFingerprintQuery,
              variables: {
                fingerprint: {
                  hash,
                  algorithm,
                },
              },
            })
            .then((res) => {
              newFingerprints[s.id] =
                res.data.findSceneByFingerprint.length > 0
                  ? res.data.findSceneByFingerprint[0]
                  : null;
            });
        })
    );

    setFingerprints(newFingerprints);
    setLoadingFingerprints(false);
  };

  const canFingerprintSearch = () =>
    scenes.some(
      (s) => s.stash_ids.length === 0 && fingerprints[s.id] === undefined
    );
  const getFingerprintCount = () => {
    const count = scenes.filter(
      (s) => s.stash_ids.length === 0 && fingerprints[s.id]
    ).length;
    return `${count > 0 ? count : "No"} new fingerprint matches found`;
  };

  const stashBoxes = stashConfig.data?.configuration.general.stashBoxes ?? [];

  return (
    <div className="tagger-container mx-auto">
      <div className="row mb-2 no-gutters">
        <Button
          onClick={() => setShowConfig(!showConfig)}
          variant="link"
          disabled={authFailure}
        >
          {showConfig ? "Hide" : "Show"} Configuration
        </Button>
      </div>

      <Collapse in={showConfig || authFailure}>
        <Card>
          <div className="row">
            <h4 className="col-12">Configuration</h4>
            <hr className="w-100" />
            <Form className="col-6">
              <Form.Group controlId="tag-males" className="align-items-center">
                <Form.Check
                  label="Show male performers"
                  checked={config.showMales}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                    setConfig({ ...config, showMales: e.currentTarget.checked })
                  }
                />
                <Form.Text>
                  Toggle whether male performers will be available to tag.
                </Form.Text>
              </Form.Group>
              <Form.Group controlId="set-cover" className="align-items-center">
                <Form.Check
                  label="Set scene cover image"
                  checked={config.setCoverImage}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                    setConfig({
                      ...config,
                      setCoverImage: e.currentTarget.checked,
                    })
                  }
                />
                <Form.Text>Replace the scene cover if one is found.</Form.Text>
              </Form.Group>
              <Form.Group className="align-items-center">
                <div className="d-flex align-items-center">
                  <Form.Check
                    id="tag-mode"
                    label="Set tags"
                    className="mr-4"
                    checked={config.setTags}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      setConfig({ ...config, setTags: e.currentTarget.checked })
                    }
                  />
                  <Form.Control
                    id="tag-operation"
                    className="col-2"
                    as="select"
                    value={config.tagOperation}
                    onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
                      setConfig({
                        ...config,
                        tagOperation: e.currentTarget.value,
                      })
                    }
                    disabled={!config.setTags}
                  >
                    <option value="merge">Merge</option>
                    <option value="overwrite">Overwrite</option>
                  </Form.Control>
                </div>
                <Form.Text>
                  Attach tags to scene, either by overwriting or merging with
                  existing tags on scene.
                </Form.Text>
              </Form.Group>

              <Form.Group controlId="mode-select">
                <div className="row no-gutters">
                  <Form.Label className="mr-4 mt-1">Query Mode:</Form.Label>
                  <Form.Control
                    as="select"
                    className="col-2"
                    value={config.mode}
                    onChange={(e: React.ChangeEvent<HTMLSelectElement>) =>
                      setConfig({
                        ...config,
                        mode: e.currentTarget.value as ParseMode,
                      })
                    }
                  >
                    <option value="auto">Auto</option>
                    <option value="filename">Filename</option>
                    <option value="dir">Dir</option>
                    <option value="path">Path</option>
                    <option value="metadata">Metadata</option>
                  </Form.Control>
                </div>
                <Form.Text>{ModeDesc[config.mode]}</Form.Text>
              </Form.Group>
            </Form>
            <div className="col-6">
              <h5>Blacklist</h5>
              <InputGroup>
                <Form.Control
                  value={blacklistInput}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                    setBlacklistInput(e.currentTarget.value)
                  }
                />
                <InputGroup.Append>
                  <Button onClick={handleBlacklistAddition}>Add</Button>
                </InputGroup.Append>
              </InputGroup>
              <div>
                Blacklist items are excluded from queries. Note that they are
                regular expressions and also case-insensitive. Certain
                characters must be escaped with a backslash:{" "}
                <code>[\^$.|?*+()</code>
              </div>
              {config.blacklist.map((item, index) => (
                <Badge
                  className="tag-item d-inline-block"
                  variant="secondary"
                  key={item}
                >
                  {item.toString()}
                  <Button
                    className="minimal ml-2"
                    onClick={() => removeBlacklist(index)}
                  >
                    <Icon icon="times" />
                  </Button>
                </Badge>
              ))}

              <Form.Group
                controlId="stash-box-endpoint"
                className="align-items-center row no-gutters mt-4"
              >
                <Form.Label className="mr-4">
                  Active stash-box instance:
                </Form.Label>
                <Form.Control
                  as="select"
                  value={credentials?.endpoint}
                  className="col-4"
                  disabled={!stashBoxes.length}
                  onChange={handleInstanceSelect}
                >
                  {!stashBoxes.length && <option>No instances found</option>}
                  {stashConfig.data?.configuration.general.stashBoxes.map(
                    (i) => (
                      <option value={i.endpoint} key={i.endpoint}>
                        {i.endpoint}
                      </option>
                    )
                  )}
                </Form.Control>
              </Form.Group>

              <div className="row">
                {user?.me?.id ? (
                  <h5 className="text-success col">
                    Connection successful. You are logged in as{" "}
                    <b>{user.me.name}</b>.
                  </h5>
                ) : (
                  <h5 className="text-danger col">
                    Connection failed.{" "}
                    <a href="/settings?tab=configuration">
                      Please check that the endpoint and API key are correct.
                    </a>
                  </h5>
                )}
              </div>
            </div>
          </div>
        </Card>
      </Collapse>

      <Card className="tagger-table">
        <div className="tagger-table-header row mb-4">
          <div className="col-6">
            <b>Path</b>
          </div>
          <div className="col-2">
            <b>Query</b>
          </div>
          <div className="col-4 text-right">
            <Button
              onClick={handleFingerprintSearch}
              disabled={
                authFailure || (!canFingerprintSearch() && !loadingFingerprints)
              }
            >
              {canFingerprintSearch() && <span>Match Fingerprints</span>}
              {!canFingerprintSearch() && getFingerprintCount()}
              {loadingFingerprints && (
                <LoadingIndicator message="" inline small />
              )}
            </Button>
          </div>
        </div>
        {scenes.map((scene) => {
          const { paths, file, ext } = parsePath(scene.path);
          const originalDir = scene.path.slice(
            0,
            scene.path.length - file.length - ext.length
          );
          const defaultQueryString = prepareQueryString(
            scene,
            paths,
            file,
            config.mode,
            config.blacklist
          );
          const modifiedQuery = queryString[scene.id];
          const fingerprintMatch = fingerprints[scene.id];
          return (
            <div key={scene.id} className="my-2 search-item">
              <div className="row">
                <div className="col-6 text-truncate align-self-center">
                  <a
                    href={`/scenes/${scene.id}`}
                    className="scene-link"
                    title={scene.path}
                  >
                    {originalDir}
                    <wbr />
                    {`${file}.${ext}`}
                  </a>
                </div>
                <div className="col-6">
                  {!taggedScenes[scene.id] && scene?.stash_ids.length > 0 && (
                    <h5 className="text-right text-bold">
                      Scene already tagged
                    </h5>
                  )}
                  {!taggedScenes[scene.id] && !scene?.stash_ids.length && (
                    <InputGroup>
                      <Form.Control
                        value={modifiedQuery || defaultQueryString}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                          setQueryString({
                            ...queryString,
                            [scene.id]: e.currentTarget.value,
                          })
                        }
                        onKeyPress={(
                          e: React.KeyboardEvent<HTMLInputElement>
                        ) =>
                          e.key === "Enter" &&
                          doBoxSearch(
                            scene.id,
                            queryString[scene.id] || defaultQueryString
                          )
                        }
                      />
                      <InputGroup.Append>
                        <Button
                          disabled={authFailure || loading}
                          onClick={() =>
                            doBoxSearch(
                              scene.id,
                              queryString[scene.id] || defaultQueryString
                            )
                          }
                        >
                          Search
                        </Button>
                      </InputGroup.Append>
                    </InputGroup>
                  )}
                  {taggedScenes[scene.id] && (
                    <h5 className="row no-gutters">
                      <b className="col-4">Scene successfully tagged:</b>
                      <a
                        className="offset-1 col-7 text-right"
                        href={`/scenes/${scene.id}`}
                      >
                        {taggedScenes[scene.id].title}
                      </a>
                    </h5>
                  )}
                </div>
              </div>
              {searchResults[scene.id] === null && <div>No results found.</div>}
              { /*
              TODO
              {fingerprintMatch &&
                credentials.endpoint &&
                !scene?.stash_ids.length &&
                !taggedScenes[scene.id] && (
                  <StashSearchResult
                    showMales={config.showMales}
                    stashScene={scene}
                    isActive
                    setActive={() => {}}
                    setScene={handleTaggedScene}
                    scene={fingerprintMatch}
                    setCoverImage={config.setCoverImage}
                    tagOperation={config.tagOperation}
                    client={client}
                    endpoint={credentials.endpoint}
                  />
                )}
                 */}
              {searchResults[scene.id] &&
                !taggedScenes[scene.id] &&
                !fingerprintMatch && (
                  <ul className="pl-0 mt-4">
                    {searchResults[scene.id]
                      .sort((a, b) => {
                        const adur =
                          a?.duration ??
                          a?.fingerprints.map((f) => f.duration)?.[0] ??
                          null;
                        const bdur =
                          b?.duration ??
                          b?.fingerprints.map((f) => f.duration)?.[0] ??
                          null;
                        if (!adur && !bdur) return 0;
                        if (adur && !bdur) return -1;
                        if (!adur && bdur) return 1;

                        const sceneDur = scene.file.duration;
                        if (!sceneDur) return 0;

                        const aDiff = Math.abs((adur ?? 0) - sceneDur);
                        const bDiff = Math.abs((bdur ?? 0) - sceneDur);

                        if (aDiff < bDiff) return -1;
                        if (aDiff > bDiff) return 1;
                        return 0;
                      })
                      .map(
                        (sceneResult, i) =>
                          credentials.endpoint &&
                          sceneResult && (
                            <StashSearchResult
                              key={sceneResult.id}
                              showMales={config.showMales}
                              stashScene={scene}
                              scene={sceneResult}
                              isActive={(selectedResult?.[scene.id] ?? 0) === i}
                              setActive={() =>
                                setSelectedResult({
                                  ...selectedResult,
                                  [scene.id]: i,
                                })
                              }
                              setCoverImage={config.setCoverImage}
                              tagOperation={config.tagOperation}
                              setScene={handleTaggedScene}
                              client={client}
                              endpoint={credentials.endpoint}
                            />
                          )
                      )}
                  </ul>
                )}
            </div>
          );
        })}
      </Card>
    </div>
  );
};
