import React, { useMemo } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { Form, Row, Col } from "react-bootstrap";
import { Movie, MovieSelect } from "src/components/Movies/MovieSelect";
import cx from "classnames";

export type MovieSceneIndexMap = Map<string, number | undefined>;

export interface IMovieEntry {
  movie: Movie;
  scene_index?: GQL.InputMaybe<number> | undefined;
}

export interface IProps {
  value: IMovieEntry[];
  onUpdate: (input: IMovieEntry[]) => void;
}

export const SceneMovieTable: React.FC<IProps> = (props) => {
  const { value, onUpdate } = props;

  const intl = useIntl();

  const movieIDs = useMemo(() => value.map((m) => m.movie.id), [value]);

  const updateFieldChanged = (index: number, sceneIndex: number | null) => {
    const newValues = value.map((existing, i) => {
      if (i === index) {
        return {
          ...existing,
          scene_index: sceneIndex,
        };
      }
      return existing;
    });

    onUpdate(newValues);
  };

  function onMovieSet(index: number, movies: Movie[]) {
    if (!movies.length) {
      // remove this entry
      const newValues = value.filter((_, i) => i !== index);
      onUpdate(newValues);
      return;
    }

    const movie = movies[0];

    const newValues = value.map((existing, i) => {
      if (i === index) {
        return {
          ...existing,
          movie: movie,
        };
      }
      return existing;
    });

    onUpdate(newValues);
  }

  function onNewMovieSet(movies: Movie[]) {
    if (!movies.length) {
      return;
    }

    const movie = movies[0];

    const newValues = [
      ...value,
      {
        movie: movie,
        scene_index: null,
      },
    ];

    onUpdate(newValues);
  }

  function renderTableData() {
    return (
      <>
        {value.map((m, i) => (
          <Row key={m.movie.id} className="movie-row">
            <Col xs={9}>
              <MovieSelect
                onSelect={(items) => onMovieSet(i, items)}
                values={[m.movie!]}
                excludeIds={movieIDs}
              />
            </Col>
            <Col xs={3}>
              <Form.Control
                className="text-input"
                type="number"
                value={m.scene_index ?? ""}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  updateFieldChanged(
                    i,
                    e.currentTarget.value === ""
                      ? null
                      : Number.parseInt(e.currentTarget.value, 10)
                  );
                }}
              />
            </Col>
          </Row>
        ))}
        <Row className="movie-row">
          <Col xs={12}>
            <MovieSelect
              onSelect={(items) => onNewMovieSet(items)}
              values={[]}
              excludeIds={movieIDs}
            />
          </Col>
        </Row>
      </>
    );
  }

  return (
    <div className={cx("movie-table", { "no-movies": !value.length })}>
      <Row className="movie-table-header">
        <Col xs={9}></Col>
        <Form.Label column xs={3} className="movie-scene-number-header">
          {intl.formatMessage({ id: "movie_scene_number" })}
        </Form.Label>
      </Row>
      {renderTableData()}
    </div>
  );
};
