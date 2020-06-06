import * as React from "react";
import * as GQL from "src/core/generated-graphql";
import { useAllMoviesForFilter } from "src/core/StashService";
import { Form } from "react-bootstrap";

type ValidTypes = GQL.SlimMovieDataFragment;

export type MovieSceneIndexMap = Map<string, number | undefined>;

export interface IProps {
  movieSceneIndexes: MovieSceneIndexMap;
  onUpdate: (value: MovieSceneIndexMap) => void;
}

export const SceneMovieTable: React.FunctionComponent<IProps> = (
  props: IProps
) => {
  const { data } = useAllMoviesForFilter();

  const items = !!data && !!data.allMoviesSlim ? data.allMoviesSlim : [];
  let itemsFilter: ValidTypes[] = [];

  if (!!props.movieSceneIndexes && !!items) {
    props.movieSceneIndexes.forEach((_index, movieId) => {
      itemsFilter = itemsFilter.concat(items.filter((x) => x.id === movieId));
    });
  }

  const storeIdx = itemsFilter.map((movie) => {
    return props.movieSceneIndexes.get(movie.id);
  });

  const updateFieldChanged = (movieId: string, value: number) => {
    const newMap = new Map(props.movieSceneIndexes);
    newMap.set(movieId, value);
    props.onUpdate(newMap);
  };

  function renderTableData() {
    return (
      <tbody>
        {itemsFilter!.map((item, index: number) => (
          <tr key={item.toString()}>
            <td>{item.name} </td>
            <td />
            <td>Scene number:</td>
            <td>
              <Form.Control
                as="select"
                className="input-control"
                value={storeIdx[index] ? storeIdx[index]?.toString() : ""}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  updateFieldChanged(
                    item.id,
                    Number.parseInt(
                      e.currentTarget.value ? e.currentTarget.value : "0",
                      10
                    )
                  )
                }
              >
                {["", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"].map(
                  (opt) => (
                    <option value={opt} key={opt}>
                      {opt}
                    </option>
                  )
                )}
              </Form.Control>
            </td>
          </tr>
        ))}
      </tbody>
    );
  }

  return (
    <div>
      <table className="movie-table">{renderTableData()}</table>
    </div>
  );
};
