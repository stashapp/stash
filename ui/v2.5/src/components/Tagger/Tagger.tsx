import React, { useCallback, useState, useRef } from "react";
import ApolloClient from 'apollo-client';
import { ApolloLink } from 'apollo-link';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { HttpLink } from 'apollo-link-http';
import { Button, Table } from 'react-bootstrap';
import path from 'parse-filepath';
import { debounce } from "lodash";

import { StashService } from 'src/core/StashService';
import * as GQL from 'src/core/generated-graphql';
import { Pagination } from "src/components/List/Pagination";

import { SearchSceneVariables, SearchScene } from 'src/definitions-box/SearchScene';
import { loader } from 'graphql.macro';
const SearchSceneQuery = loader('src/queries/searchScene.gql');

const ApiKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxODM2ZWE0YS03NmQyLTRjODEtYjJiMy1mNGVjZGNjOTRmOWUiLCJpYXQiOjE1ODMxMzgzODMsInN1YiI6IkFQSUtleSJ9.Ac1PtkDWvIZuIBstqDuFQad_vYlZfHHyrKE-DXvGBgc';

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
function prepareQueryString(str: string) {
  let s = str;
  blacklist.forEach(b => { s = s.replace(b, '') });
  const date = s.match(dateRegex);
  if(date) {
    s = s.replace(date[0], ` 20${date[1]}-${date[2]}-${date[3]} `);
  }
  return s.split(/(?=[A-Z])/).join(' ').replace(/\./g, ' ');
}


export const Tagger: React.FC = () => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<SearchScene>();
  const [searchFilter, setSearchFilter] = useState('');
  const [page, setPage] = useState(1);
  const [searchResults, setSearchResults] = useState<Record<string, SearchScene>>({});

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
        [sceneID]: queryData.data
      });
      setLoading(false);
    });

    setLoading(true);
  };

  const scenes = sceneData?.findScenes?.scenes ?? [];

  return (
    <div>
      <div className="row">
        <input className="form-control col-2" onChange={(event: React.FormEvent<HTMLInputElement>) => searchCallback(event.currentTarget.value)} ref={inputRef} placeholder="Search text" disabled={sceneLoading} />
      </div>
      { !loading && (
        <div>
          <div>Name: {data?.searchScene?.[0]?.title}</div>
          <div>Date: {data?.searchScene?.[0]?.date}</div>
        </div>
      )}

      <Pagination
        currentPage={page}
        itemsPerPage={20}
        totalItems={sceneData?.findScenes?.count ?? 0}
        onChangePage={newPage => setPage(newPage)}
      />

      <Table>
        <thead>
          <tr>
            <th>Filename</th>
            <th>Path</th>
            <th>StashDB-query</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {scenes.map(scene => {
            const parsedPath = path(scene.path);
            const dir = path(parsedPath.dir).base;
            return (
              <>
                <tr>
                  <td>{parsedPath.base}</td>
                  <td>{dir}</td>
                  <td>
                    <Button disabled={loading} onClick={() => doBoxSearch(scene.id, prepareQueryString(dir))}>Search</Button>
                    {prepareQueryString(dir)}
                  </td>
                </tr>
                { searchResults[scene.id] && (
                  <tr>
                    <td colSpan={3}>
                      <ul>
                        { searchResults[scene.id].searchScene.map(scene => (
                          <li>Title: {scene?.title ?? 'Unknown'}, Date: {scene?.date ?? 'Unknown'}, Studio: {scene?.studio?.name ?? 'Unknown'}, Performer(s): { scene?.performers.map(p => p.performer.name).join(', ') }</li>
                        ))
                        }
                      </ul>
                    </td>
                  </tr>
                )}
              </>
            );
          })}
        </tbody>
      </Table>
    </div>
  );
};
