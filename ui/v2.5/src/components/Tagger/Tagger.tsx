import React, { useCallback, useState, useRef } from "react";
import ApolloClient from 'apollo-client';
import { ApolloLink } from 'apollo-link';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { HttpLink } from 'apollo-link-http';
import { Button, Form, InputGroup } from 'react-bootstrap';
import path from 'parse-filepath';
import { debounce } from "lodash";
import cx from 'classnames';

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

const blacklist = ['XXX', '1080p', '720p', '2160p', /-[A-Z]+\[rarbg\]/, 'MP4'];
const dateRegex = /\.(\d\d)\.(\d\d)\.(\d\d)\./;
function prepareQueryString(scene: Partial<GQL.Scene>, str: string) {
  if (scene.date && scene.studio) {
    return `${scene.date} ${scene.studio.name} ${(scene?.performers ?? []).map(p => p.name).join(' ')}`;
  }
  let s = str;
  blacklist.forEach(b => { s = s.replace(b, '') });
  const date = s.match(dateRegex);
  if(date) {
    s = s.replace(date[0], ` 20${date[1]}-${date[2]}-${date[3]} `);
  }
  return s.split(/(?=[A-Z])/).join(' ').replace(/\./g, ' ');
}

interface ISelectResult {
  label: string;
  value: string;
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

export const Tagger: React.FC = () => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [loading, setLoading] = useState(false);
  const [searchFilter, setSearchFilter] = useState('');
  const [page, setPage] = useState(1);
  const [searchResults, setSearchResults] = useState<Record<string, SearchScene|null>>({});
  const [queryString, setQueryString] = useState<Record<string, string>>({});
  const [selectedResult, setSelectedResult] = useState<Record<string, number>>();
  const [showMales, setShowMales] = useState(false);

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

  const scenes = sceneData?.findScenes?.scenes ?? [];

  return (
    <div className="col-9 mx-auto">
      <div className="row my-2">
        <input
          className="form-control col-2 ml-4"
          onChange={(event: React.FormEvent<HTMLInputElement>) => searchCallback(event.currentTarget.value)}
          ref={inputRef}
          placeholder="Search text"
          disabled={sceneLoading}
        />
        <Form.Group controlId="tag-males" className="col-2 mr-auto d-flex align-items-center mt-1">
          <Form.Check label="Show male performers" onChange={(e: React.FormEvent<HTMLInputElement>) => setShowMales(e.currentTarget.checked)} />
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

      <div className="tagger-table card">
        <div className="tagger-table-header row mb-4">
          <div className="col-6"><b>Path</b></div>
          <div className="col-6"><b>StashDB Query</b></div>
        </div>
          {scenes.map(scene => {
            const parsedPath = path(scene.path);
            const dir = path(parsedPath.dir).base;
            const defaultQueryString = prepareQueryString(scene, dir);
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
                { searchResults[scene.id] === null && <div>No results found.</div> }
                { searchResults[scene.id] && (
                  <div className="col mt-4">
                    <ul className="pl-0">
                      { searchResults[scene.id]?.searchScene.map((sceneResult, i) => (
                        sceneResult && (
                          <StashSearchResult
                            showMales={showMales}
                            stashScene={scene}
                            scene={sceneResult}
                            isActive={(selectedResult?.[scene.id] ?? 0) === i}
                            setActive={() => setSelectedResult({ ...selectedResult, [scene.id]: i})}
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
}

const PerformerResult: React.FC<IPerformerResultProps> = ({ performer }) => {
  const [selectedPerformer, setSelectedPerformer] = useState();
  const [selectedSource, setSelectedSource] = useState<'create'|'existing'|undefined>();
  const [modalVisible, showModal] = useState(false);
  const [newPerformer, setNewPerformer] = useState<GQL.PerformerCreateInput|undefined>();
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
      const performer = data.findPerformers?.performers?.[0]?.id;
      if (performer) {
        setSelectedPerformer(performer);
        setSelectedSource('existing');
      }
    }
  });

  const handlePerformerSelect = (items: ValidTypes[]) => {
    if (items.length) {
      setSelectedSource('existing');
      setSelectedPerformer(items[0].id);
    }
    else {
      setSelectedSource(undefined);
      setSelectedPerformer(null);
    }
  };

  const handlePerformerCreate = () => {
    setSelectedSource('create');
    showModal(false);
  };

  if(stashLoading || loading)
    return <div>Loading studio</div>;

  if((stashData?.findPerformers.count ?? 0) > 0) {
    return (
      <div>
        <SuccessIcon />
        <span>{ stashData!.findPerformers.performers[0].name }</span>
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
      <div className="col-4">
        Performer:
        <b className="ml-2">{performer.name}</b>
      </div>
      <div className="col-2">
        <Button variant="secondary" className="mr-2" onClick={() => showModal(true)}>Create</Button>
        { selectedSource === 'create'
          ? <SuccessIcon />
          : <FailIcon />
        }
      </div>
      <div className="col-3">
        { selectedSource === 'existing'
          ? <SuccessIcon />
          : <FailIcon />
        }
        <span className="d-inline-block">Select existing:</span>
      </div>
      <div className="col-3">
        <PerformerSelect
          ids={selectedPerformer ? [selectedPerformer] : []}
          onSelect={handlePerformerSelect}
        />
      </div>
    </div>
  );
}

interface IStudioResultProps {
  studio: StashStudio|null;
}

const StudioResult: React.FC<IStudioResultProps> = ({ studio }) => {
  const [selectedStudio, setSelectedStudio] = useState();
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
    onCompleted: (data) => {
      setSelectedStudio(data.findStudios.count > 0
        ? data.findStudios.studios[0].id
        : null
      );
    }
  });

  if(loading || stashLoading)
    return <div>Loading performer</div>;

  if(stashData?.findStudioByStashID) {
    return <div>{ stashData.findStudioByStashID?.name }</div>
  }
  return <StudioSelect ids={selectedStudio ? [selectedStudio] : []} onSelect={(items) => setSelectedStudio(items.length ? items[0].id : null)} />
}

interface IStashSearchResultProps {
  scene: SearchResult;
  stashScene: Partial<GQL.Scene>;
  isActive: boolean;
  setActive: () => void;
  showMales: boolean;
}

type PerformerData = StashPerformer|string;
type StudioData = StashStudio|string;

const StashSearchResult: React.FC<IStashSearchResultProps> = ({ scene, stashScene, isActive, setActive, showMales }) => {
  const [studio, setStudio] = useState<StudioData>();
  const [performers, setPerformers] = useState<Record<string, PerformerData>>();

  const [updateScene] = GQL.useSceneUpdateMutation({ variables: {
    id: stashScene.id
  }});

  const setPerformer = (performerData: PerformerData, performerID: string) => (
    setPerformers({ ...performers, [performerID]: performerData })
  );

  const classname = cx('row mb-4 search-result', { 'selected-result': isActive });
  return (
    <li className={classname} key={scene?.id}>
      <div className="col-6 row">
        <label className="d-flex justify-content-center align-items-center col-2 scene-select">
          <input type="radio" checked={isActive} onChange={setActive} />
        </label>
        <div className="d-flex col-3">
          <img height={100} src={scene?.urls?.[0]?.url} alt="" className="align-self-center" />
        </div>
        <div className="col-7 d-flex flex-column justify-content-center">
          <h4 className="text-truncate">{scene?.title}</h4>
          <h5>{scene?.studio?.name} • {scene?.date}</h5>
          <div>Performers: {scene?.performers?.map(p => p.performer.name).join(', ')}</div>
        </div>
      </div>
      { isActive && (
        <div className="col-6">
          <StudioResult studio={scene.studio} />
          { scene.performers
            .filter(p => p.performer.gender !== 'MALE' || showMales)
            .map(performer => (
            <PerformerResult performer={performer.performer} />
            ))
          }
          <div className="row pr-3 mt-2">
            <Button className="col-1 offset-11">Save</Button>
          </div>
        </div>
      )}
    </li>
  );
};
