import React, { useCallback, useEffect, useState, useRef, Dispatch, SetStateAction } from "react";
import ApolloClient from 'apollo-client';
import { ApolloLink } from 'apollo-link';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { HttpLink } from 'apollo-link-http';
import { Button, Form, InputGroup } from 'react-bootstrap';
import path from 'parse-filepath';
import { debounce } from "lodash";
import cx from 'classnames';
import { blobToBase64 } from 'base64-blob';

import { BreastTypeEnum } from 'src/definitions-box/globalTypes';
import * as GQL from 'src/core/generated-graphql';
import { Pagination } from "src/components/List/Pagination";
import { Icon, Modal, PerformerSelect, StudioSelect } from 'src/components/Shared';
import { ValidTypes } from 'src/components/Shared/Select';

import {
  SearchSceneVariables,
  SearchScene,
  SearchScene_searchScene as SearchResult,
  SearchScene_searchScene_performers_performer as StashPerformer,
  SearchScene_searchScene_studio as StashStudio
} from 'src/definitions-box/SearchScene';
import { loader } from 'graphql.macro';

const SearchSceneQuery = loader('src/queries/searchScene.gql');

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

const DEFAULT_BLACKLIST = ['XXX', '1080p', '720p', '2160p', /-[A-Z]+\[rarbg\]/, 'MP4'];
const dateRegex = /\.(\d\d)\.(\d\d)\.(\d\d)\./;
function prepareQueryString(scene: Partial<GQL.Scene>, str: string, mode:ParseMode) {
  if ((scene.date && scene.studio) || mode === "metadata") {
    return `${scene.date} ${scene.studio?.name ?? ''} ${(scene?.performers ?? []).map(p => p.name).join(' ')}`;
  }
  let s = str;
  DEFAULT_BLACKLIST.forEach(b => { s = s.replace(b, '') });
  const date = s.match(dateRegex);
  if(date) {
    s = s.replace(date[0], ` 20${date[1]}-${date[2]}-${date[3]} `);
  }
  return s.split(/(?=[A-Z])/).join(' ').replace(/\./g, ' ');
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

export const Tagger: React.FC = () => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [loading, setLoading] = useState(false);
  const [searchFilter, setSearchFilter] = useState('');
  const [page, setPage] = useState(1);
  const [searchResults, setSearchResults] = useState<Record<string, SearchScene|null>>({});
  const [queryString, setQueryString] = useState<Record<string, string>>({});
  const [selectedResult, setSelectedResult] = useState<Record<string, number>>();
  const [showMales, setShowMales] = useState(false);
  const [advanced, setAdvanced] = useState(false);
  const [mode, setMode] = useState<ParseMode>('auto');
  const [blacklist, setBlacklist] = useState<(string|RegExp)[]>(DEFAULT_BLACKLIST);
  const [taggedScenes, setTaggedScenes] = useState<Record<string, Partial<GQL.Scene>>>({});

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

  const scenes = sceneData?.findScenes?.scenes ?? [];

  return (
    <div className="tagger-container mx-auto">
      <div className="row my-2">
        <input
          className="form-control col-2 ml-4"
          onChange={(event: React.FormEvent<HTMLInputElement>) => searchCallback(event.currentTarget.value)}
          ref={inputRef}
          placeholder="Search text"
          disabled={sceneLoading}
        />
        <Form.Group controlId="tag-males" className="mx-4 d-flex align-items-center mt-1">
          <Form.Check label="Show male performers" onChange={(e: React.FormEvent<HTMLInputElement>) => setShowMales(e.currentTarget.checked)} />
        </Form.Group>
        <Form.Group controlId="advanced-config" className="mr-auto d-flex align-items-center mt-1">
          <Form.Check label="Advanced configuration" onChange={(e: React.FormEvent<HTMLInputElement>) => setAdvanced(e.currentTarget.checked)} />
        </Form.Group>
        <div className="float-right mr-4">
          <Pagination
            currentPage={page}
            itemsPerPage={20}
            totalItems={sceneData?.findScenes?.count ?? 0}
            onChangePage={newPage => setPage(newPage)}
          />
        </div>
      </div>

      { advanced && (
        <div className="mx-4">
          <div className="row">
            <h5 className="col">Advanced settings</h5>
          </div>
          <Form.Group controlId="mode-select" className="col-4">
            <Form.Label>Mode: </Form.Label>
              <Form.Control as="select" selected={mode} onChange={(e:React.FormEvent<HTMLSelectElement>) => setMode(e.currentTarget.value as ParseMode)}>
                <option value="auto">Auto</option>
                <option value="filename">Filename</option>
                <option value="dir">Dir</option>
                <option value="path">Path</option>
                <option value="metadata">Metadata</option>
            </Form.Control>
          </Form.Group>
        </div>
      )}

      <div className="tagger-table card">
        <div className="tagger-table-header row mb-4">
          <div className="col-6"><b>Path</b></div>
          <div className="col-6"><b>StashDB Query</b></div>
        </div>
          {scenes.map(scene => {
            const parsedPath = path(scene.path);
            const dir = path(parsedPath.dir).base;
            const defaultQueryString = prepareQueryString(scene, dir, mode);
            return (
              <div key={scene.id} className="mb-4">
                <div className="row">
                  <div className="col-6">{`${dir}/${parsedPath.base}`}</div>
                  <div className="col-6">
                    <InputGroup>
                      <Form.Control
                        defaultValue={defaultQueryString}
                        onChange={(e: React.FormEvent<HTMLInputElement>) => setQueryString({ ...queryString, [scene.id]: e.currentTarget.value})} />
                      <InputGroup.Append>
                        <Button disabled={loading} onClick={() => doBoxSearch(scene.id, queryString[scene.id] || defaultQueryString)}>Search</Button>
                      </InputGroup.Append>
                    </InputGroup>
                  </div>
                </div>
                { scene?.stash_id && <div className="col-5 offset-6 text-right"><b>Scene already tagged</b></div> }
                { searchResults[scene.id] === null && <div>No results found.</div> }
                { taggedScenes[scene.id] && (
                  <div className="col-5 offset-6 text-right">
                    <b>Scene successfully tagged:</b>
                    <a href={`/scenes/${scene.id}`}>{taggedScenes[scene.id].title}</a>
                  </div>
                )}
                { searchResults[scene.id] && !taggedScenes[scene.id] && (
                  <div className="col mt-4">
                    <ul className="pl-0">
                      { searchResults[scene.id]?.searchScene.map((sceneResult, i) => (
                        sceneResult && (
                          <StashSearchResult
                            key={sceneResult.id}
                            showMales={showMales}
                            stashScene={scene}
                            scene={sceneResult}
                            isActive={(selectedResult?.[scene.id] ?? 0) === i}
                            setActive={() => setSelectedResult({ ...selectedResult, [scene.id]: i})}
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

  const handlePerformerCreate = () => {
    setSelectedSource('create');
    setPerformer({
      type: 'Create',
      data: performer
    });
    showModal(false);
  };

  if(stashLoading || loading)
    return <div>Loading studio</div>;

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
      <Modal
        show={modalVisible}
        accept={{ text: "Save", onClick: handlePerformerCreate }}
        cancel={{ onClick: () => showModal(false), variant: "secondary" }}
      >
        <div className="row">
          <div className="col-6">
            <div className="row">
              <div className="col-6">Name:</div>{ performer.name }
            </div>
            <div className="row">
              <div className="col-6">Gender:</div> { performer.gender }
            </div>
            <div className="row">
              <div className="col-6">Birthdate:</div> { performer.birthdate?.date ?? 'Unknown' }
            </div>
            <div className="row">
              <div className="col-6">Ethnicity:</div>{ performer.ethnicity }
            </div>
            <div className="row">
              <div className="col-6">Country:</div>{ performer.country }
            </div>
            <div className="row">
              <div className="col-6">Eye Color:</div>{ performer.eye_color }
            </div>
            <div className="row">
              <div className="col-6">Height:</div>{ performer.height }
            </div>
            <div className="row">
              <div className="col-6">Measurements:</div>{ `${performer.measurements.band_size}${performer.measurements.cup_size}-${performer.measurements.waist}-${performer.measurements.hip}` }
            </div>
            <div className="row">
              <div className="col-6">Fake Tits:</div>{ performer.breast_type === BreastTypeEnum.FAKE ? "Yes" : "No" }
            </div>
            <div className="row">
              <div className="col-6">Career Length:</div>{ `${performer.career_start_year} - ${ performer.career_end_year ?? ''}` }
            </div>
            <div className="row">
              <div className="col-6">Tattoos:</div>{ performer.tattoos?.join(', ') ?? '' }
            </div>
            <div className="row">
              <div className="col-6">Piercings:</div>{ performer.piercings?.join(', ') ?? '' }
            </div>
          </div>
          <div className="col-6">
            <img src={performer.urls?.[0]?.url} className="w-100" />
          </div>
        </div>
      </Modal>
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
}

interface IPerformerOperation {
  type: "Create"|"Existing"|"Update";
  data: StashPerformer|string;
}

interface IStudioOperation {
  type: "Create"|"Existing"|"Update";
  data: StashStudio|string;
}

const StashSearchResult: React.FC<IStashSearchResultProps> = ({ scene, stashScene, isActive, setActive, showMales, setScene }) => {
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
      const imgurl = scene.urls?.[scene.urls.length - 1]?.url;
      let imgData = null;
      if(imgurl) {
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
    }
  };

  const classname = cx('row mb-4 search-result', { 'selected-result': isActive });
  return (
    <li className={classname} key={scene?.id} onClick={() => !isActive && setActive()}>
      <div className="col-6 row">
        <img src={scene?.urls?.[0]?.url} alt="" className="align-self-center scene-image" />
        <div className="d-flex flex-column justify-content-center scene-metadata">
          <h4 className="text-truncate">{scene?.title}</h4>
          <h5>{scene?.studio?.name} • {scene?.date}</h5>
          <div>Performers: {scene?.performers?.map(p => p.performer.name).join(', ')}</div>
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
