import * as React from "react";
import * as GQL from "src/core/generated-graphql";
import { useAllMoviesForFilter } from "src/core/StashService";
import { Form, Row, Col } from "react-bootstrap";

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
      <>
        {itemsFilter!.map((item, index: number) => (
          <Row key={item.toString()}>
            <Form.Label column xs={9}>
              {item.name}
            </Form.Label>
            <Col xs={3}>
              <Form.Control
                className="text-input"
                type="number"
                value={storeIdx[index] ? storeIdx[index]?.toString() : ""}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  updateFieldChanged(
                    item.id,
                    Number.parseInt(
                      e.currentTarget.value ? e.currentTarget.value : "0",
                      10
                    )
                  );
                }}
              />
            </Col>
          </Row>
        ))}
      </>
    );
  }

  if (props.movieSceneIndexes.size > 0) {
    return (
      <div className="movie-table">
        <Row>
          <Form.Label column xs={9}>
            Movie
          </Form.Label>
          <Form.Label column xs={3}>
            Scene #
          </Form.Label>
        </Row>
        {renderTableData()}
      </div>
    );
  }

  return <></>;
};
