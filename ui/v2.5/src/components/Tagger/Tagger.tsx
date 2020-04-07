import React, { useCallback, useEffect, useState, useRef, Dispatch, SetStateAction } from "react";
import ApolloClient from 'apollo-client';
import { ApolloLink } from 'apollo-link';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { HttpLink } from 'apollo-link-http';
import { Badge, Button, Form, InputGroup } from 'react-bootstrap';
import path from 'parse-filepath';
import { debounce } from "lodash";
import cx from 'classnames';
import { blobToBase64 } from 'base64-blob';
import localForage from "localforage";

import { BreastTypeEnum, FingerprintAlgorithm, GenderEnum } from 'src/definitions-box/globalTypes';
import * as GQL from 'src/core/generated-graphql';
import { FindPerformersDocument, FindStudioByStashIdDocument } from '../../core/generated-graphql';
import { Pagination } from "src/components/List/Pagination";
import { Icon, LoadingIndicator, Modal, PerformerSelect, StudioSelect } from 'src/components/Shared';
import { ValidTypes } from 'src/components/Shared/Select';

import {
  SearchSceneVariables,
  SearchScene,
  SearchScene_searchScene as SearchResult,
  SearchScene_searchScene_performers_performer as StashPerformer,
  SearchScene_searchScene_studio as StashStudio
} from 'src/definitions-box/SearchScene';
import {
  FindSceneByFingerprintVariables,
  FindSceneByFingerprint,
  FindSceneByFingerprint_findSceneByFingerprint as FingerprintResult
} from 'src/definitions-box/FindSceneByFingerprint';
import {
  SubmitFingerprintVariables,
  SubmitFingerprint
} from 'src/definitions-box/SubmitFingerprint';
import { loader } from 'graphql.macro';

const SubmitFingerprintMutation = loader('src/queries/submitFingerprint.gql');
const SearchSceneQuery = loader('src/queries/searchScene.gql');
const FindSceneByFingerprintQuery = loader('src/queries/searchFingerprint.gql');

const ApiKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI0N2ZkNzAwYi03ZmVlLTQyYTktYTBiYy1kMTUyYTQzMWQzYjkiLCJpYXQiOjE1ODQyMjI3MzgsInN1YiI6IkFQSUtleSJ9.FGpuM_4QxqA4iuMeioWGriqhpuVpTKrcpV2rTIyZ3wc';

const createClient = () => {
  const httpLink = new HttpLink({
    uri: 'http://localhost:9998/graphql',
    fetch
  });

	const middlewareLink = new ApolloLink((operation, forward) => {
		operation.setContext({ headers: { ApiKey } });
		if (forward)
				return forward(operation);
		return null;
	});

	const link = middlewareLink.concat(httpLink);
  const cache = new InMemoryCache();

	return new ApolloClient({
		name: 'stashdb',
		connectToDevTools: true,
		link,
		cache
	});
};

const client = createClient();

const DEFAULT_BLACKLIST = [' XXX ', '1080p', '720p', '2160p', 'KTR', 'RARBG', 'MP4', 'x264', '\\[', '\\]'];
const dateRegex = /\.(\d\d)\.(\d\d)\.(\d\d)\./;
function prepareQueryString(scene: Partial<GQL.Scene>, paths: string[], mode:ParseMode, blacklist: string[]) {
  if ((mode === 'auto' && scene.date && scene.studio) || mode === "metadata") {
    let str = [
      scene.date,
      scene.studio?.name ?? '',
      (scene?.performers ?? []).map(p => p.name).join(' '),
      scene?.title ?? ''
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

interface IconProps {
  className?: string;
}
const SuccessIcon: React.FC<IconProps> = ({ className }) => (
  <Icon icon="check" className={cx("success mr-4", className)} color="#0f9960" />
);
const FailIcon: React.FC<IconProps> = ({ className }) => (
  <Icon icon="times" className={cx("secondary mr-4", className)} color="#394b59" />
);

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
        blacklist: data.blacklist ?? DEFAULT_BLACKLIST,
        showMales: data.showMales ?? false,
        mode: data.mode ?? 'auto',
        setCoverImage: data.setCoverImage ?? true,
      });
    }
  )}, []);

  useEffect(() => {
    localForage.setItem('tagger', config);
  }, [config]);

  const { data: sceneData, loading: sceneLoading } = GQL.useFindScenesQuery({
    variables: {
      scene_filter: {},
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

    for (const s of scenes) {
      if(fingerprints[s.id] !== undefined)
        continue;

      const res = await client.query<FindSceneByFingerprint, FindSceneByFingerprintVariables>({
          query: FindSceneByFingerprintQuery,
          variables: {
            fingerprint: {
              hash: s.checksum,
              algorithm: FingerprintAlgorithm.MD5
            }
          }
      });

      newFingerprints[s.id] = res.data.findSceneByFingerprint.length > 0 ?
        res.data.findSceneByFingerprint[0] : null;
    };

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
              <Badge className={`tag-item d-inline-block`} variant="secondary">
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

      <div className="row my-2">
        <input
          className="form-control col-2 ml-4"
          onChange={(event: React.FormEvent<HTMLInputElement>) => searchCallback(event.currentTarget.value)}
          ref={inputRef}
          placeholder="Search text"
          disabled={sceneLoading}
        />
        <Form.Group controlId="tag-males" className="mx-4 d-flex align-items-center mt-1">
          <Form.Check label="Show male performers" onChange={(e: React.FormEvent<HTMLInputElement>) => setConfig({ ...config, showMales: e.currentTarget.checked })} />
        </Form.Group>
        <Form.Group controlId="set-cover" className="mx-4 d-flex align-items-center mt-1">
          <Form.Check label="Set scene cover image" onChange={(e: React.FormEvent<HTMLInputElement>) => setConfig({ ...config, setCoverImage: e.currentTarget.checked })} />
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
            const dir = parsedPath.dir;
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
                { fingerprintMatch && (
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

interface IPerformerResultProps {
  performer: StashPerformer
  setPerformer: (data:IPerformerOperation) => void;
}

const PerformerResult: React.FC<IPerformerResultProps> = ({ performer, setPerformer }) => {
  const [selectedPerformer, setSelectedPerformer] = useState();
  const [selectedSource, setSelectedSource] = useState<'create'|'existing'|undefined>();
  const [modalVisible, showModal] = useState(false);
  const { data: stashData, loading: stashLoading } = GQL.useFindPerformersQuery({
    variables: {
      performer_filter: {
        stash_id: {
          value: performer.id,
          modifier: GQL.CriterionModifier.Equals
        }
      }
    }
  });
  const { loading } = GQL.useFindPerformersQuery({
    variables: {
      filter: {
        q: `"${performer.name}"`
      }
    },
    onCompleted: (data) => {
      const performerResult = data.findPerformers?.performers?.[0]?.id;
      if (performerResult) {
        setSelectedPerformer(performerResult);
        setSelectedSource('existing');
        setPerformer({
          type: 'Update',
          data: performerResult
        });
      }
    }
  });

  useEffect(() => {
    if(!stashData?.findPerformers.count)
      return;

    setPerformer({
      type: 'Existing',
      data: stashData.findPerformers.performers[0].id
    });
  }, [stashData]);

  const handlePerformerSelect = (items: ValidTypes[]) => {
    if (items.length) {
      setSelectedSource('existing');
      setSelectedPerformer(items[0].id);
      setPerformer({
        type: 'Update',
        data: items[0].id
      });
    }
    else {
      setSelectedSource(undefined);
      setSelectedPerformer(null);
    }
  };

  const handlePerformerCreate = (imageIndex: number) => {
    const images = sortImageURLs(performer.urls, 'portrait');
    const urls = images.length ? [{
      url: images[imageIndex].url,
      type: 'PHOTO'
    }] : [];
    setSelectedSource('create');
    setPerformer({
      type: 'Create',
      data: {
        ...performer,
        urls
      }
    });
    showModal(false);
  };

  if(stashLoading || loading)
    return <div>Loading performer</div>;

  if((stashData?.findPerformers.count ?? 0) > 0) {
    return (
      <div className="row my-2">
        <span className="ml-auto">
          <SuccessIcon />Performer matched:
        </span>
        <b className="col-3 text-right">{ stashData!.findPerformers.performers[0].name }</b>
      </div>
    );
  }
  return (
    <div className="row align-items-center mt-2">
      <PerformerModal
        showModal={showModal}
        modalVisible={modalVisible}
        performer={performer}
        handlePerformerCreate={handlePerformerCreate}
      />
      <div className="entity-name">
        Performer:
        <b className="ml-2">{performer.name}</b>
      </div>
      <div>
        <Button variant="secondary" className="mr-1" onClick={() => showModal(true)}>Create</Button>
        { selectedSource === 'create'
          ? <SuccessIcon />
          : <FailIcon />
        }
      </div>
      <div className="select-existing">
        { selectedSource === 'existing'
          ? <SuccessIcon />
          : <FailIcon />
        }
      </div>
      <PerformerSelect
        ids={selectedPerformer ? [selectedPerformer] : []}
        onSelect={handlePerformerSelect}
        className="performer-select"
      />
    </div>
  );
}

interface IPerformerModalProps {
  performer: StashPerformer;
  modalVisible: boolean;
  showModal: (show: boolean) => void;
  handlePerformerCreate: (imageIndex: number) => void;
};

const PerformerModal: React.FC<IPerformerModalProps> = ({ modalVisible, performer, handlePerformerCreate, showModal }) => {
  const [imageIndex, setImageIndex] = useState(0);

  const images = sortImageURLs(performer.urls, 'portrait');

  const setPrev = () => (
    setImageIndex(imageIndex === 0 ? images.length - 1 : imageIndex - 1)
  );
  const setNext = () => (
    setImageIndex(imageIndex === images.length - 1 ? 0 : imageIndex + 1)
  );

  return (
    <Modal
      show={modalVisible}
      accept={{ text: "Save", onClick: () => handlePerformerCreate(imageIndex) }}
      cancel={{ onClick: () => showModal(false), variant: "secondary" }}
      onHide={() => showModal(false)}
    >
      <div className="row">
        <div className="col-6">
          <div className="row no-gutters">
            <strong className="col-6">Name:</strong>
            <span className="col-6 text-truncate">{ performer.name }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Gender:</strong>
            <span className="col-6 text-truncate">{ performer.gender}</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Birthdate:</strong>
            <span className="col-6 text-truncate">{ performer.birthdate?.date ?? 'Unknown' }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Ethnicity:</strong>
            <span className="col-6 text-truncate">{ performer.ethnicity }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Country:</strong>
            <span className="col-6 text-truncate">{ performer.country }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Eye Color:</strong>
            <span className="col-6 text-truncate">{ performer.eye_color }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Height:</strong>
            <span className="col-6 text-truncate">{ performer.height }</span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Measurements:</strong>
            <span className="col-6 text-truncate">{
              (performer.measurements.cup_size && !performer.measurements.waist && !performer.measurements.hip) &&
              `${performer.measurements.band_size}${performer.measurements.cup_size}-${performer.measurements.waist}-${performer.measurements.hip}` }</span>
          </div>
          { performer?.gender !== GenderEnum.MALE && (
            <div className="row no-gutters">
              <strong className="col-6">Fake Tits:</strong>
              <span className="col-6 text-truncate">{ performer.breast_type === BreastTypeEnum.FAKE ? "Yes" : "No" }</span>
            </div>
          )}
          <div className="row no-gutters">
            <strong className="col-6">Career Length:</strong>
            <span className="col-6 text-truncate">{
              (performer.career_start_year) &&
              `${performer.career_start_year} - ${ performer.career_end_year ?? ''}`} </span>
          </div>
          <div className="row no-gutters">
            <strong className="col-6">Tattoos:</strong>
            <span className="col-6 text-truncate">{ performer.tattoos?.join(', ') ?? '' }</span>
          </div>
          <div className="row no-gutters ">
            <strong className="col-6">Piercings:</strong>
            <span className="col-6 text-truncate">{ performer.piercings?.join(', ') ?? '' }</span>
          </div>
        </div>
        { images.length > 0 && (
          <div className="col-6">
            <img src={images[imageIndex].url} alt='' className="w-100" />
            <div className="d-flex mt-2">
              <Button className="mr-auto" onClick={setPrev}>
                <Icon icon="arrow-left" />
              </Button>
              <h5>Image {imageIndex+1} of {images.length}</h5>
              <Button className="ml-auto" onClick={setNext}>
                <Icon icon="arrow-right" />
              </Button>
            </div>
          </div>
        )}
      </div>
    </Modal>
  );
}

interface IStudioResultProps {
  studio: StashStudio|null;
  setStudio: Dispatch<SetStateAction<IStudioOperation|undefined>>;
}

const StudioResult: React.FC<IStudioResultProps> = ({ studio, setStudio }) => {
  const [selectedStudio, setSelectedStudio] = useState();
  const [modalVisible, showModal] = useState(false);
  const [selectedSource, setSelectedSource] = useState<'create'|'existing'|undefined>();
  const { data: stashData, loading: stashLoading } = GQL.useFindStudioByStashIdQuery({
    variables: {
      id: studio?.id ?? ''
    }
  })
  const { loading } = GQL.useFindStudiosQuery({
    variables: {
      filter: {
        q: `"${studio?.name ?? ''}"`
      }
    },
    onCompleted: (data) => (
      handleStudioSelect(data.findStudios?.studios?.[0]?.id)
    )
  });

  useEffect(() => {
    if(!stashData?.findStudioByStashID)
      return;

    setStudio({
      type: 'Existing',
      data: stashData.findStudioByStashID.id
    });
  }, [stashData]);

  const handleStudioCreate = () => {
    if(!studio)
      return;
    setSelectedSource('create');
    setStudio({
      type: 'Create',
      data: studio
    });
    showModal(false);
  };

  const handleStudioSelect = (id?: string) => {
    if (id) {
      setSelectedStudio(id);
      setSelectedSource('existing');
      setStudio({
        type: 'Update',
        data: id
      });
    }
    else {
      setSelectedSource(undefined);
      setSelectedStudio(null);
    }
  };

  if(loading || stashLoading)
    return <div>Loading studio</div>;

  if(stashData?.findStudioByStashID) {
    return (
      <div className="row my-2">
        <span className="ml-auto">
          <SuccessIcon />Studio matched:
        </span>
        <b className="col-3 text-right">{ stashData.findStudioByStashID.name }</b>
      </div>
    );
  }


  return (
    <div className="row align-items-center mt-2">
      <Modal
        show={modalVisible}
        accept={{ text: "Save", onClick: handleStudioCreate }}
        cancel={{ onClick: () => showModal(false), variant: "secondary" }}
      >
        <div className="row">
          <div className="col-6">
            <div className="row">
              <div className="col-6">Name:</div>{ studio?.name }
            </div>
          </div>
        </div>
      </Modal>

      <div className="entity-name">
        Studio:
        <b className="ml-2">{studio?.name}</b>
      </div>
      <div>
        <Button variant="secondary" className="mr-1" onClick={() => showModal(true)}>Create</Button>
        { selectedSource === 'create'
          ? <SuccessIcon />
          : <FailIcon />
        }
      </div>
      <div className="select-existing">
        { selectedSource === 'existing'
          ? <SuccessIcon />
          : <FailIcon />
        }
      </div>
      <StudioSelect
        ids={selectedStudio ? [selectedStudio] : []}
        onSelect={(items) => handleStudioSelect(items.length ? items[0].id : undefined)}
        className="studio-select"
      />
    </div>

  );
}

interface IStashSearchResultProps {
  scene: SearchResult;
  stashScene: Partial<GQL.Scene>;
  isActive: boolean;
  setActive: () => void;
  showMales: boolean;
  setScene: (scene: Partial<GQL.Scene>) => void;
  isFingerprintMatch?: boolean;
  setCoverImage: boolean;
}

interface IPerformerOperation {
  type: "Create"|"Existing"|"Update";
  data: StashPerformer|string;
}

interface IStudioOperation {
  type: "Create"|"Existing"|"Update";
  data: StashStudio|string;
}

const StashSearchResult: React.FC<IStashSearchResultProps> = ({ scene, stashScene, isActive, setActive, showMales, setScene, setCoverImage }) => {
  const [studio, setStudio] = useState<IStudioOperation>();
  const [performers, setPerformers] = useState<Record<string, IPerformerOperation>>();

  const [createStudio] = GQL.useStudioCreateMutation();
  const [updateStudio] = GQL.useStudioUpdateMutation();
  const [updateScene] = GQL.useSceneUpdateMutation();
  const [createPerformer] = GQL.usePerformerCreateMutation();
  const [updatePerformer] = GQL.usePerformerUpdateMutation();

  const setPerformer = (performerData: IPerformerOperation, performerID: string) => (
    setPerformers({ ...performers, [performerID]: performerData })
  );

  const handleSave = async () => {
    if(!performers || !studio)
      return;

    let studioID:string;
    let performerIDs = [];

    if (studio.type === 'Update') {
      const studioUpdateResult = await updateStudio({
        variables: {
          id: studio.data as string,
          stash_id: scene.studio?.id ?? ''
        },
        update: (store, studio) => {
          if (!studio?.data?.studioUpdate)
            return;

          store.writeQuery({
            query: FindStudioByStashIdDocument,
            variables: {
              id: studio.data.studioUpdate.stash_id
            },
            data: {
              findStudioByStashID: studio.data.studioUpdate
            }
          });
        }
      });
      const id = studioUpdateResult.data?.studioUpdate?.id;
      if(studioUpdateResult.errors || !id)
        return;
      studioID = id;
    }
    else if(studio.type === 'Create') {
      const studioData = studio.data as StashStudio;
      const studioCreateResult = await createStudio({
        variables: {
          name: studioData.name,
          stash_id: studioData.id
        }
      });

      const id = studioCreateResult.data?.studioCreate?.id;
      if(studioCreateResult.errors || !id)
        return;
      studioID = id;
    }
    else {
      studioID = studio.data as string;
    }

    performerIDs = await Promise.all(Object.keys(performers).map(async (performerID) => {
      const performer = performers[performerID];
      if (performer.type === 'Update') {
        const res = await updatePerformer({
          variables: {
            id: performer.data as string,
            stash_id: performerID
          },
          update: (store, performer) => {
            if (!performer?.data?.performerUpdate)
              return;

            store.writeQuery({
              query: FindPerformersDocument,
              variables: {
                performer_filter: {
                  stash_id: {
                    value: performer.data.performerUpdate.stash_id,
                    modifier: GQL.CriterionModifier.Equals
                  }
                }
              },
              data: {
                findPerformers: {
                  performers: [performer.data.performerUpdate],
                  count: 1,
                  __typename: "FindPerformersResultType"
                }
              }
            });
          }
        });

        if(res.errors)
          return;

        return res?.data?.performerUpdate?.id ?? null;
      }
      if(performer.type === 'Create') {
        const performerData = performer.data as StashPerformer;
        const imgurl = performerData.urls?.[0]?.url;
        let imgData = null;
        if(imgurl) {
          const img = await fetch(imgurl);
          if(img.status === 200) {
            const blob = await img.blob();
            imgData = await blobToBase64(blob);
          }
        }

        const res = await createPerformer({
          variables: {
            name: performerData.name,
            country: performerData.country,
            height: performerData.height?.toString(),
            ethnicity: performerData.ethnicity,
            birthdate: performerData.birthdate?.date ?? null,
            eye_color: performerData.eye_color,
            fake_tits: performerData.breast_type === BreastTypeEnum.FAKE ? 'Yes' : 'No',
            measurements: `${performerData.measurements.band_size}${performerData.measurements.cup_size}-${performerData.measurements.waist}-${performerData.measurements.hip}`,
            image: imgData,
            stash_id: performerID
          }
        });

        if(res.errors)
          return;

        return res?.data?.performerCreate?.id ?? null;
      }
      return performer.data as string;
    }));

    if(studioID && !performerIDs.some(id => !id)) {
      const imgurl = getUrlByType(scene.urls, 'PHOTO', 'landscape');
      let imgData = null;
      if(imgurl && setCoverImage) {
        const img = await fetch(imgurl);
        if(img.status === 200) {
          const blob = await img.blob();
          imgData = await blobToBase64(blob);
        }
      }
      const sceneUpdateResult = await updateScene({
        variables: {
          id: stashScene.id ?? '',
          stash_id: scene.id,
          title: scene.title,
          details: scene.details,
          date: scene.date,
          performer_ids: performerIDs as string[],
          studio_id: studioID,
          cover_image: imgData
        }
      });
      if(sceneUpdateResult.data?.sceneUpdate)
        setScene(sceneUpdateResult.data.sceneUpdate);

      if(stashScene.checksum)
        client.mutate<SubmitFingerprint, SubmitFingerprintVariables>({
          mutation: SubmitFingerprintMutation,
          variables: {
            input: {
              scene_id: scene.id,
              fingerprint: {
                hash: stashScene.checksum,
                algorithm: FingerprintAlgorithm.MD5
              }
            }
          }
        });
    }
  };

  const classname = cx('row mb-4 search-result', { 'selected-result': isActive });
  return (
    <li className={classname} key={scene?.id} onClick={() => !isActive && setActive()}>
      <div className="col-6 row">
        <img src={getUrlByType(scene?.urls as URLInput[], 'PHOTO', 'landscape')} alt="" className="align-self-center scene-image" />
        <div className="d-flex flex-column justify-content-center scene-metadata">
          <h4 className="text-truncate">{scene?.title}</h4>
          <h5>{scene?.studio?.name} • {scene?.date}</h5>
          <div>Performers: {scene?.performers?.map(p => p.performer.name).join(', ')}</div>
          { getDurationStatus(scene.duration, stashScene.file?.duration) }
        </div>
      </div>
      { isActive && (
        <div className="col-6">
          <StudioResult studio={scene.studio} setStudio={setStudio} />
          { scene.performers
            .filter(p => p.performer.gender !== 'MALE' || showMales)
            .map(performer => (
              <PerformerResult performer={performer.performer} setPerformer={(data:IPerformerOperation) => setPerformer(data, performer.performer.id)} key={`${scene.id}${performer.performer.id}`} />
            ))
          }
          <div className="row pr-3 mt-2">
            <Button className="col-1 offset-11" onClick={handleSave}>Save</Button>
          </div>
        </div>
      )}
    </li>
  );
};

interface URLInput {
  url: string;
  type: string;
};

const getDurationStatus = (dbDuration: number|null, stashDuration: number|undefined|null) => {
  if(!dbDuration || !stashDuration) return '';
  const diff = Math.abs(dbDuration - stashDuration);
  if(diff < 5) {
    return <div><b>Duration is a match</b></div>;
  }
  return <div>Duration off by {Math.floor(diff)}s</div>;
};

const sortImageURLs = (urls: URLInput[], orientation: 'portrait'|'landscape') => {
  return urls.filter((u) => u.type === 'PHOTO').map((u:URLInput) => {
      const width = Number.parseInt(u.url.match(/width=(\d+)/)?.[1] ?? '0', 10);
      const height = Number.parseInt(u.url.match(/height=(\d+)/)?.[1] ?? '0', 10);
      return {
          url: u.url,
          width,
          height,
          aspect: orientation === 'portrait' ? (height / width > 1) : (width / height) > 1
      }
  }).sort((a, b) => {
      if (a.aspect > b.aspect) return -1;
      if (a.aspect < b.aspect) return 1;
      if (orientation === 'portrait' && a.height > b.height) return -1;
      if (orientation === 'portrait' && a.height < b.height) return 1;
      if (orientation === 'landscape' && a.width > b.width) return -1;
      if (orientation === 'landscape' && a.width < b.width) return 1;
      return 0;
  });
}

export const getUrlByType = (
    urls:URLInput[],
    type:string,
    orientation?: 'portrait'|'landscape'
) => {
  if (urls.length > 0 && type === 'PHOTO' && orientation)
    return sortImageURLs(urls, orientation)[0].url;
  return (urls && (urls.find((url) => url.type === type) || {}).url) || '';
};
