import React from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { useAllMoviesForFilter } from "src/core/StashService";
import { Form, Row, Col } from "react-bootstrap";

export type MovieSceneIndexMap = Map<string, number | undefined>;

export interface IProps {
  movieScenes: GQL.SceneMovieInput[];
  onUpdate: (value: GQL.SceneMovieInput[]) => void;
}

export const SceneMovieTable: React.FC<IProps> = (props) => {
  const intl = useIntl();
  const { data } = useAllMoviesForFilter();

  const items = !!data && !!data.allMovies ? data.allMovies : [];

  const movieEntries = props.movieScenes.map((m) => {
    return {
      movie: items.find((mm) => m.movie_id === mm.id),
      ...m,
    };
  });

  const updateFieldChanged = (movieId: string, value: number) => {
    const newValues = props.movieScenes.map((ms) => {
      if (ms.movie_id === movieId) {
        return {
          movie_id: movieId,
          scene_index: value,
        };
      }
      return ms;
    });
    props.onUpdate(newValues);
  };

  function renderTableData() {
    return (
      <>
        {movieEntries.map((m) => (
          <Row key={m.movie_id}>
            <Form.Label column xs={9}>
              {m.movie?.name ?? ""}
            </Form.Label>
            <Col xs={3}>
              <Form.Control
                className="text-input"
                type="number"
                value={m.scene_index ?? ""}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  updateFieldChanged(
                    m.movie_id,
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

  if (props.movieScenes.length > 0) {
    return (
      <div className="movie-table">
        <Row>
          <Form.Label column xs={9}>
            {intl.formatMessage({ id: "movie" })}
          </Form.Label>
          <Form.Label column xs={3}>
            {intl.formatMessage({ id: "movie_scene_number" })}
          </Form.Label>
        </Row>
        {renderTableData()}
      </div>
    );
  }

  return <></>;
};
