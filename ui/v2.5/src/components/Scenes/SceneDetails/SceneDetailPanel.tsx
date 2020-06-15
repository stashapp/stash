import React from "react";
import { Link } from "react-router-dom";
import { FormattedDate } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { TagLink } from "src/components/Shared";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { RatingStars } from "./RatingStars";

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
}

export const SceneDetailPanel: React.FC<ISceneDetailProps> = (props) => {
  function renderDetails() {
    if (!props.scene.details || props.scene.details === "") return;
    return (
      <>
        <h6>Details</h6>
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
        <h6>Tags</h6>
        {tags}
      </>
    );
  }

  function renderPerformers() {
    if (props.scene.performers.length === 0) return;
    const cards = props.scene.performers.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={props.scene.date ?? undefined}
      />
    ));

    return (
      <>
        <h6>Performers</h6>
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
            <h3 className="text-truncate">
              {props.scene.title ??
                TextUtils.fileNameFromPath(props.scene.path)}
            </h3>
          </div>
          {props.scene.date ? (
            <h5>
              <FormattedDate value={props.scene.date} format="long" />
            </h5>
          ) : undefined}
          {props.scene.rating ? (
            <h6>
              Rating: <RatingStars value={props.scene.rating} />
            </h6>
          ) : (
            ""
          )}
          {props.scene.file.height && (
            <h6>Resolution: {TextUtils.resolution(props.scene.file.height)}</h6>
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
