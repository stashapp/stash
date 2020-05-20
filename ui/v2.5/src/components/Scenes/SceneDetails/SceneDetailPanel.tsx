import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { TagLink } from "src/components/Shared";
import { PerformerCard } from "src/components/Performers/PerformerCard";

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
        <div className="row justify-content-center">{cards}</div>
      </>
    );
  }

  return (
    <>
    <div className="row">
      <h3 className="col scene-header text-truncate d-xl-none">
        {props.scene.title ?? TextUtils.fileNameFromPath(props.scene.path)}
      </h3>
      <div className="col-6 col-xl-12 scene-details">
        <h4>{props.scene.date ?? ""}</h4>
        {props.scene.rating ? <h6>Rating: {props.scene.rating}</h6> : ""}
        {props.scene.file.height && (
          <h6>Resolution: {TextUtils.resolution(props.scene.file.height)}</h6>
        )}
        {renderDetails()}
        {renderTags()}
      </div>
      <div className="col-4 offset-2 d-xl-none">
        {props.scene.studio && (
          <Link to={`/studios/${props.scene.studio.id}`}>
            <img
              src={props.scene.studio.image_path ?? ""}
              alt={`${props.scene.studio.name} logo`}
              className="studio-logo"
            />
          </Link>
        )}
      </div>
    </div>
    <div className="row">
      <div className="col-12">
          {renderPerformers()}
        </div>
    </div>
    </>
  );
};
