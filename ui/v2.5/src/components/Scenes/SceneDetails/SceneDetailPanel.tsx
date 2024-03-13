import React, { useMemo } from "react";
import { Link } from "react-router-dom";
import { FormattedDate, FormattedMessage, useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { TagLink } from "src/components/Shared/TagLink";
import { TruncatedText } from "src/components/Shared/TruncatedText";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { RatingSystem } from "src/components/Shared/Rating/RatingSystem";
import { objectTitle } from "src/core/files";
import { DirectorLink } from "src/components/Shared/Link";

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
}

export const SceneDetailPanel: React.FC<ISceneDetailProps> = (props) => {
  const intl = useIntl();

  const file = useMemo(
    () => (props.scene.files.length > 0 ? props.scene.files[0] : undefined),
    [props.scene]
  );

  function renderDetails() {
    if (!props.scene.details || props.scene.details === "") return;
    return (
      <>
        <h6>
          <FormattedMessage id="details" />:{" "}
        </h6>
        <p className="pre">{props.scene.details}</p>
      </>
    );
  }

  function renderTags() {
    if (props.scene.tags.length === 0) return;
    const tags = props.scene.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} />
    ));
    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.tags"
            values={{ count: props.scene.tags.length }}
          />
        </h6>
        {tags}
      </>
    );
  }

  function renderPerformers() {
    if (props.scene.performers.length === 0) return;
    const performers = sortPerformers(props.scene.performers);
    const cards = performers.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={props.scene.date ?? undefined}
      />
    ));

    return (
      <>
        <h6>
          <FormattedMessage
            id="countables.performers"
            values={{ count: props.scene.performers.length }}
          />
        </h6>
        <div className="row justify-content-center scene-performers">
          {cards}
        </div>
      </>
    );
  }

  // filename should use entire row if there is no studio
  const sceneDetailsWidth = props.scene.studio ? "col-9" : "col-12";

  return (
    <>
      <div className="row">
        <div className={`${sceneDetailsWidth} col-xl-12 scene-details`}>
          <div className="scene-header d-xl-none">
            <h3>
              <TruncatedText text={objectTitle(props.scene)} />
            </h3>
          </div>
          {props.scene.date ? (
            <h5>
              <FormattedDate
                value={props.scene.date}
                format="long"
                timeZone="utc"
              />
            </h5>
          ) : undefined}
          {props.scene.rating100 ? (
            <h6>
              <FormattedMessage id="rating" />:{" "}
              <RatingSystem value={props.scene.rating100} disabled />
            </h6>
          ) : (
            ""
          )}
          {file?.width && file?.height && (
            <h6>
              <FormattedMessage id="resolution" />:{" "}
              {TextUtils.resolution(file.width, file.height)}
            </h6>
          )}
          <h6>
            <FormattedMessage id="created_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.scene.created_at)}{" "}
          </h6>
          <h6>
            <FormattedMessage id="updated_at" />:{" "}
            {TextUtils.formatDateTime(intl, props.scene.updated_at)}{" "}
          </h6>
          {props.scene.code && (
            <h6>
              <FormattedMessage id="scene_code" />: {props.scene.code}{" "}
            </h6>
          )}
          {props.scene.director && (
            <h6>
              <FormattedMessage id="director" />:{" "}
              <DirectorLink director={props.scene.director} linkType="scene" />
            </h6>
          )}
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
      <div className="row">
        <div className="col-12">
          {renderDetails()}
          {renderTags()}
          {renderPerformers()}
        </div>
      </div>
    </>
  );
};

export default SceneDetailPanel;
