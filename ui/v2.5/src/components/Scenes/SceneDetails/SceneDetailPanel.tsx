import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { TagLink } from "src/components/Shared/TagLink";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { objectTitle } from "src/core/files";
import { DirectorLink } from "src/components/Shared/Link";
import { DetailItem } from "src/components/Shared/DetailItem";

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
}

export const SceneDetailPanel: React.FC<ISceneDetailProps> = (props) => {
  const intl = useIntl();

  // filename should use entire row if there is no studio
  const sceneDetailsWidth = props.scene.studio ? "col-9" : "col-12";

  const tags = useMemo(
    () => props.scene.tags.map((tag) => <TagLink key={tag.id} tag={tag} />),
    [props.scene.tags]
  );

  const performers = useMemo(() => {
    const sorted = sortPerformers(props.scene.performers);
    return sorted.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={props.scene.date ?? undefined}
      />
    ));
  }, [props.scene.performers, props.scene.date]);

  const details = useMemo(() => {
    return props.scene.details?.length ? (
      <p className="pre">{props.scene.details}</p>
    ) : undefined;
  }, [props.scene.details]);

  return (
    <>
      <div className="row">
        <div className={`${sceneDetailsWidth} col-xl-12 scene-details`}>
          <div className="scene-header d-xl-none">
            <h3>
              <TruncatedText text={objectTitle(props.scene)} />
            </h3>
          </div>

          <div className="detail-group">
            <DetailItem id="studio-code" value={props.scene.code} fullWidth />
            <DetailItem
              id="director"
              value={
                props.scene.director ? (
                  <DirectorLink
                    director={props.scene.director}
                    linkType="scene"
                  />
                ) : undefined
              }
              fullWidth
            />
            <DetailItem id="tags" value={tags.length ? tags : undefined} />
            <DetailItem id="details" value={details} />
            <DetailItem
              id="performers"
              value={performers.length ? performers : undefined}
            />
            <DetailItem
              id="created_at"
              value={TextUtils.formatDateTime(intl, props.scene.created_at)}
              fullWidth
            />
            <DetailItem
              id="updated_at"
              value={TextUtils.formatDateTime(intl, props.scene.updated_at)}
              fullWidth
            />
          </div>
        </div>
        {props.scene.studio && (
          <div className="col-3 d-xl-none">
            <Link to={`/studios/${props.scene.studio.id}`}>
              <img
                src={props.scene.studio.image_path ?? ""}
                alt={`${props.scene.studio.name} logo`}
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
