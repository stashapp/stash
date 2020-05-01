import React, { useCallback, useEffect, useState, useRef } from "react";
import { Badge, Button, Form, InputGroup } from 'react-bootstrap';
import path from 'parse-filepath';
import { debounce } from "lodash";
import localForage from "localforage";

import { FingerprintAlgorithm } from 'src/definitions-box/globalTypes';
import * as GQL from 'src/core/generated-graphql';
import { Pagination } from "src/components/List/Pagination";
import { Icon, LoadingIndicator } from 'src/components/Shared';

import {
  SearchSceneVariables,
  SearchScene
} from 'src/definitions-box/SearchScene';
import {
  FindSceneByFingerprintVariables,
  FindSceneByFingerprint,
  FindSceneByFingerprint_findSceneByFingerprint as FingerprintResult
} from 'src/definitions-box/FindSceneByFingerprint';
import { loader } from 'graphql.macro';
import StashSearchResult from './StashSearchResult';
import { client } from './client';

const SearchSceneQuery = loader('src/queries/searchScene.gql');
const FindSceneByFingerprintQuery = loader('src/queries/searchFingerprint.gql');

const DEFAULT_BLACKLIST = [' XXX ', '1080p', '720p', '2160p', 'KTR', 'RARBG', 'MP4', 'x264', '\\[', '\\]'];
const dateRegex = /\.(\d\d)\.(\d\d)\.(\d\d)\./;
function prepareQueryString(scene: Partial<GQL.Scene>, paths: string[], mode:ParseMode, blacklist: string[]) {
  if ((mode === 'auto' && scene.date && scene.studio) || mode === "metadata") {
    let str = [
      scene.date,
      scene.studio?.name ?? '',
      (scene?.performers ?? []).map(p => p.name).join(' '),
      scene?.title ? scene.title.replace(/[^a-zA-Z0-9 ]+/g, '') : ''
    ].filter(s => s !== '').join(' ');
    blacklist.forEach(b => { str = str.replace(new RegExp(b, 'gi'), '') });
    return str;
  }
  let s = '';
  if (mode === 'auto' || mode === 'filename') {
    s = paths[paths.length - 1];
  }
  else if (mode === 'path') {
    s = paths.join(' ');
  } else {
    s = paths[paths.length - 2];
  }
  blacklist.forEach(b => { s = s.replace(new RegExp(b, 'i'), '') });
  const date = s.match(dateRegex);
  s = s.replace(/-/g, ' ');
  if(date) {
    s = s.replace(date[0], ` 20${date[1]}-${date[2]}-${date[3]} `);
  }
  return s.replace(/\./g, ' ');
}

type ParseMode = 'auto'|'filename'|'dir'|'path'|'metadata';
const ModeDesc = {
  'auto': 'Uses metadata if present, or filename',
  'metadata': 'Only uses metadata',
  'filename': 'Only uses filename',
  'dir': 'Only uses parent directory of video file',
  'path': 'Uses entire file path'
};


interface ITaggerConfig {
  blacklist: string[];
  showMales: boolean;
  mode: ParseMode;
  setCoverImage: boolean;
}

export const Tagger: React.FC = () => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [loading, setLoading] = useState(false);
  const [searchFilter, setSearchFilter] = useState('');
  const [page, setPage] = useState(1);
  const [searchResults, setSearchResults] = useState<Record<string, SearchScene|null>>({});
  const [queryString, setQueryString] = useState<Record<string, string>>({});
  const [selectedResult, setSelectedResult] = useState<Record<string, number>>();
  const [blacklistInput, setBlacklistInput] = useState<string>('');
  const [taggedScenes, setTaggedScenes] = useState<Record<string, Partial<GQL.Scene>>>({});
  const [fingerprints, setFingerprints] = useState<Record<string, FingerprintResult|null>>({});
  const [loadingFingerprints, setLoadingFingerprints] = useState(false);
  const [config, setConfig] = useState<ITaggerConfig>({
    blacklist: DEFAULT_BLACKLIST,
    showMales: false,
    mode: 'auto',
    setCoverImage: true,
  });

  useEffect(() => {
    localForage.getItem<ITaggerConfig>('tagger').then((data) => {;
      setConfig({
        blacklist: data?.blacklist ?? DEFAULT_BLACKLIST,
        showMales: data?.showMales ?? false,
        mode: data?.mode ?? 'auto',
        setCoverImage: data?.setCoverImage ?? true,
      });
    }
  )}, []);

  useEffect(() => {
    localForage.setItem('tagger', config);
  }, [config]);

  const { data: sceneData, loading: sceneLoading } = GQL.useFindScenesQuery({
    variables: {
      filter: {
        q: searchFilter,
        page,
        per_page: 20
      }
    }
  });

  const searchCallback = useCallback(
    debounce((value: string) => {
      setSearchFilter(value);
      setPage(1);
    }, 500),
    []
  );

  const doBoxSearch = (sceneID: string, searchVal: string) => {
    client.query<SearchScene, SearchSceneVariables>({
      query: SearchSceneQuery,
      variables: { term: searchVal}
    }).then(queryData => {
      setSearchResults({
        ...searchResults,
        [sceneID]: queryData.data.searchScene.length > 0 ? queryData.data : null
      });
      setLoading(false);
    });

    setLoading(true);
  };

  const handleTaggedScene = (scene: Partial<GQL.Scene>) => {
    setTaggedScenes({
      ...taggedScenes,
      [scene.id as string]: scene
    });
  };

  const removeBlacklist = (index: number) => {
    setConfig({
      ...config,
      blacklist: [...config.blacklist.slice(0, index), ...config.blacklist.slice(index+1)]
    });
  };

  const handleBlacklistAddition = () => {
    setConfig({
      ...config,
      blacklist: [...config.blacklist, blacklistInput]
    });
    setBlacklistInput('');
  };

  const scenes = sceneData?.findScenes?.scenes ?? [];

  const handleFingerprintSearch = async () => {
    setLoadingFingerprints(true);
    const newFingerprints = { ...fingerprints };

    await Promise.all(
      scenes
        .filter(s => fingerprints[s.id] === undefined)
        .map(s => (
          client.query<FindSceneByFingerprint, FindSceneByFingerprintVariables>({
              query: FindSceneByFingerprintQuery,
              variables: {
                fingerprint: {
                  hash: s.checksum,
                  algorithm: FingerprintAlgorithm.MD5
                }
              }
          }).then((res) => {
            newFingerprints[s.id] = res.data.findSceneByFingerprint.length > 0 ?
              res.data.findSceneByFingerprint[0] : null;
          })
        ))
    );

    setFingerprints(newFingerprints);
    setLoadingFingerprints(false);
  };

  const canFingerprintSearch = () => {
    return !scenes.some(s => (
      fingerprints[s.id] === undefined
    ));
  };

  return (
    <div className="tagger-container mx-auto">
      <h2>StashDB Tagger</h2>
      <hr />

      <div className="row mb-4 my-2">
        <div className="col-4">
          <Form.Group controlId="mode-select">
            <Form.Label>Mode: </Form.Label>
              <Form.Control as="select" value={config.mode} onChange={(e:React.FormEvent<HTMLSelectElement>) => setConfig({ ...config, mode: e.currentTarget.value as ParseMode})}>
                <option value="auto">Auto</option>
                <option value="filename">Filename</option>
                <option value="dir">Dir</option>
                <option value="path">Path</option>
                <option value="metadata">Metadata</option>
            </Form.Control>
            <span>{ ModeDesc[config.mode] }</span>
          </Form.Group>
        </div>
        <div className="col-4">
          <h5>Blacklist</h5>
          { config.blacklist.map((item, index) => (
              <Badge className="tag-item d-inline-block" variant="secondary" key={item}>
                { item.toString() }
                <Button className="minimal ml-2" onClick={() => removeBlacklist(index)}>
                  <Icon icon="times" />
                </Button>
              </Badge>
          ))}
        </div>
        <div className="col-4">
          <h5>Add Blacklist Item</h5>
          <InputGroup>
            <Form.Control
              value={blacklistInput}
              onChange={(e: React.FormEvent<HTMLInputElement>) => setBlacklistInput(e.currentTarget.value)} />
            <InputGroup.Append>
              <Button onClick={handleBlacklistAddition}>Add</Button>
            </InputGroup.Append>
          </InputGroup>
          <div>Note that all blacklist items are regular expressions and also case-insensitive. Certain characters must be escaped with a backslash: <code>[\^$.|?*+()</code></div>
        </div>
      </div>

      <div className="row">
        <div className="col text-right mr-2">{ sceneData?.findScenes?.count } results</div>
      </div>
      <div className="row mb-2">
        <input
          className="form-control col-2 ml-4"
          onChange={(event: React.FormEvent<HTMLInputElement>) => searchCallback(event.currentTarget.value)}
          ref={inputRef}
          placeholder="Search text"
          disabled={sceneLoading}
        />
        <Form.Group controlId="tag-males" className="mx-4 d-flex align-items-center mt-1">
          <Form.Check label="Show male performers" checked={config.showMales} onChange={(e: React.FormEvent<HTMLInputElement>) => setConfig({ ...config, showMales: e.currentTarget.checked })} />
        </Form.Group>
        <Form.Group controlId="set-cover" className="mx-4 d-flex align-items-center mt-1">
          <Form.Check label="Set scene cover image" checked={config.setCoverImage} onChange={(e: React.FormEvent<HTMLInputElement>) => setConfig({ ...config, setCoverImage: e.currentTarget.checked })} />
        </Form.Group>
        <div className="float-right mr-4 ml-auto">
          <Pagination
            currentPage={page}
            itemsPerPage={20}
            totalItems={sceneData?.findScenes?.count ?? 0}
            onChangePage={newPage => setPage(newPage)}
          />
        </div>
      </div>

      <div className="tagger-table card">
        <div className="tagger-table-header row mb-4">
          <div className="col-6"><b>Path</b></div>
          <div className="col-4"><b>StashDB Query</b></div>
          <div className="col-2 text-right">
            <Button onClick={handleFingerprintSearch} disabled={canFingerprintSearch() && !loadingFingerprints}>
              Search Fingerprints
              { loadingFingerprints && <LoadingIndicator message="" inline small /> }
            </Button>
          </div>
        </div>
          {scenes.map(scene => {
            const paths = scene.path.split('/');
            const parsedPath = path(scene.path);
            const { dir } = parsedPath;
            const defaultQueryString = prepareQueryString(scene, paths, config.mode, config.blacklist);
            const modifiedQuery = queryString[scene.id];
            const fingerprintMatch = fingerprints[scene.id];
            return (
              <div key={scene.id} className="mb-4">
                <div className="row">
                  <div className="col-6">
                    <a href={`/scenes/${scene.id}`} className="scene-link">{`${dir}/${parsedPath.base}`}</a>
                  </div>
                  <div className="col-6">
                    { !taggedScenes[scene.id] && (
                      <InputGroup>
                        <Form.Control
                          value={modifiedQuery || defaultQueryString}
                          onChange={(e: React.FormEvent<HTMLInputElement>) => setQueryString({ ...queryString, [scene.id]: e.currentTarget.value})} />
                        <InputGroup.Append>
                          <Button disabled={loading} onClick={() => doBoxSearch(scene.id, queryString[scene.id] || defaultQueryString)}>Search</Button>
                        </InputGroup.Append>
                      </InputGroup>
                    )}
                    { taggedScenes[scene.id] && (
                      <h5 className="text-right">
                        <b>Scene successfully tagged:</b>
                        <a className="ml-4" href={`/scenes/${scene.id}`}>{taggedScenes[scene.id].title}</a>
                      </h5>
                    )}
                  </div>
                </div>
                { scene?.stash_id && <div className="col-5 offset-6 text-right"><b>Scene already tagged</b></div> }
                { searchResults[scene.id] === null && <div>No results found.</div> }
                { fingerprintMatch && !scene?.stash_id && !taggedScenes[scene.id] && (
                    <StashSearchResult
                      showMales={config.showMales}
                      stashScene={scene}
                      isActive
                      setActive={() => {}}
                      setScene={handleTaggedScene}
                      scene={fingerprintMatch}
                      setCoverImage={config.setCoverImage}
                      isFingerprintMatch
                    />
                )}
                { searchResults[scene.id] && !taggedScenes[scene.id] && !fingerprintMatch && (
                  <div className="col mt-4">
                    <ul className="pl-0">
                      { searchResults[scene.id]?.searchScene.sort((a, b) => {
                        if(!a?.duration && !b?.duration) return 0;
                        if(a?.duration && !b?.duration) return -1;
                        if(!a?.duration && b?.duration) return 1;

                        const sceneDur = scene.file.duration;
                        if(!sceneDur) return 0;

                        const aDiff = Math.abs((a?.duration ?? 0) - sceneDur);
                        const bDiff = Math.abs((b?.duration ?? 0) - sceneDur);

                        if(aDiff < bDiff) return -1;
                        if(aDiff > bDiff) return 1;
                        return 0;
                      }).map((sceneResult, i) => (
                        sceneResult && (
                          <StashSearchResult
                            key={sceneResult.id}
                            showMales={config.showMales}
                            stashScene={scene}
                            scene={sceneResult}
                            isActive={(selectedResult?.[scene.id] ?? 0) === i}
                            setActive={() => setSelectedResult({ ...selectedResult, [scene.id]: i})}
                            setCoverImage={config.setCoverImage}
                            setScene={handleTaggedScene}
                          />
                        )
                      ))
                      }
                    </ul>
                  </div>
                )}
              </div>
            );
          })}
      </div>
    </div>
  );
};
