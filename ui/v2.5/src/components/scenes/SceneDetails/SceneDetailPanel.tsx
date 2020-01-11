import React from "react";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { TagLink } from "src/components/Shared";
import { SceneHelpers } from "../helpers";

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
}

export const SceneDetailPanel: React.FC<ISceneDetailProps> = (props: ISceneDetailProps) => {
  function renderDetails() {
    if (!props.scene.details || props.scene.details === "") { return; }
    return (
      <>
        <h6>Details</h6>
        <p className="pre">{props.scene.details}</p>
      </>
    );
  }

  function renderTags() {
    if (props.scene.tags.length === 0) { return; }
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

  return (
    <>
    {SceneHelpers.maybeRenderStudio(props.scene, 70)}
      <h1>
        {!!props.scene.title ? props.scene.title : TextUtils.fileNameFromPath(props.scene.path)}
      </h1>
      {props.scene.date ? <h4>{props.scene.date}</h4> : ''}
      {props.scene.rating ? <h6>Rating: {props.scene.rating}</h6> : ''}
      {props.scene.file.height ? <h6>Resolution: {TextUtils.resolution(props.scene.file.height)}</h6> : ''}
      {renderDetails()}
      {renderTags()}
    </>
  );
};
