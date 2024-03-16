import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import { FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { TagLink } from "src/components/Shared/TagLink";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { objectTitle } from "src/core/files";
import { DirectorLink } from "src/components/Shared/Link";
import { DetailItem } from "src/components/Shared/DetailItem";
import { Button } from "react-bootstrap";
import { MovieCard } from "src/components/Movies/MovieCard";

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
  onClickFileDetails?: (fileID?: string) => void;
}

export const SceneDetailPanel: React.FC<ISceneDetailProps> = ({
  scene,
  onClickFileDetails,
}) => {
  const intl = useIntl();

  // filename should use entire row if there is no studio
  const sceneDetailsWidth = scene.studio ? "col-9" : "col-12";

  const tags = useMemo(
    () => scene.tags.map((tag) => <TagLink key={tag.id} tag={tag} />),
    [scene.tags]
  );

  const movies = useMemo(
    () =>
      scene.movies.map((sceneMovie) => (
        <MovieCard
          key={sceneMovie.movie.id}
          movie={sceneMovie.movie}
          sceneIndex={sceneMovie.scene_index ?? undefined}
        />
      )),
    [scene.movies]
  );

  const performers = useMemo(() => {
    const sorted = sortPerformers(scene.performers);
    return sorted.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={scene.date ?? null}
      />
    ));
  }, [scene.performers, scene.date]);

  const details = useMemo(() => {
    return scene.details?.length ? (
      <p className="pre">{scene.details}</p>
    ) : undefined;
  }, [scene.details]);

  const files = useMemo(() => {
    return (
      <ul>
        {scene.files.map((file) => (
          <li key={file.id}>
            <Button
              variant="link"
              size="sm"
              className="file-info-button"
              onClick={() => onClickFileDetails?.(file.id)}
            >
              <TruncatedText text={file.basename} />
            </Button>
          </li>
        ))}
      </ul>
    );
  }, [scene.files, onClickFileDetails]);

  return (
    <>
      <div className="row">
        <div className={`${sceneDetailsWidth} col-xl-12 scene-details`}>
          <div className="scene-header d-xl-none">
            <h3>
              <TruncatedText text={objectTitle(scene)} />
            </h3>
          </div>

          <div className="detail-group">
            <DetailItem id="studio-code" value={scene.code} fullWidth />
            <DetailItem
              id="director"
              value={
                scene.director ? (
                  <DirectorLink director={scene.director} linkType="scene" />
                ) : undefined
              }
              fullWidth
            />
            <DetailItem id="details" value={details} />
            <DetailItem
              id="movies"
              value={movies.length ? movies : undefined}
            />
            <DetailItem
              id="tags"
              heading={
                <>
                  <FormattedMessage id="tags" />
                  <Button variant="link" size="sm" className="add-tag-button">
                    <FormattedMessage id="actions.add" />
                  </Button>
                </>
              }
              value={tags.length ? tags : undefined}
            />
            <DetailItem
              id="performers"
              heading={
                <>
                  <FormattedMessage id="performers" />
                  <Button
                    variant="link"
                    size="sm"
                    className="add-performer-button"
                  >
                    <FormattedMessage id="actions.add" />
                  </Button>
                </>
              }
              value={performers.length ? performers : undefined}
            />
            <DetailItem
              id="files"
              value={scene.files.length ? files : undefined}
              fullWidth
            />
            <DetailItem
              id="created_at"
              value={TextUtils.formatDateTime(intl, scene.created_at)}
              fullWidth
            />
            <DetailItem
              id="updated_at"
              value={TextUtils.formatDateTime(intl, scene.updated_at)}
              fullWidth
            />
          </div>
        </div>
        {scene.studio && (
          <div className="col-3 d-xl-none">
            <Link to={`/studios/${scene.studio.id}`}>
              <img
                src={scene.studio.image_path ?? ""}
                alt={`${scene.studio.name} logo`}
                className="studio-logo float-right"
              />
            </Link>
          </div>
        )}
      </div>
    </>
  );
};

export default SceneDetailPanel;
