import {
  H1,
  H4,
  H6,
  Tag,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import * as GQL from "../../../core/generated-graphql";
import { TextUtils } from "../../../utils/text";
import { SceneHelpers } from "../helpers";

interface ISceneDetailProps {
  scene: GQL.SceneDataFragment;
}

export const SceneDetailPanel: FunctionComponent<ISceneDetailProps> = (props: ISceneDetailProps) => {
  function renderDetails() {
    if (!props.scene.details || props.scene.details === "") { return; }
    return (
      <>
        <H6>Details</H6>
        <p className="pre">{props.scene.details}</p>
      </>
    );
  }

  function renderTags() {
    if (props.scene.tags.length === 0) { return; }
    const tags = props.scene.tags.map((tag) => (
      <Tag key={tag.id} className="tag-item">{tag.name}</Tag>
    ));
    return (
      <>
        <H6>Tags</H6>
        {tags}
      </>
    );
  }

  return (
    <>
    {SceneHelpers.maybeRenderStudio(props.scene, 70, false)}
      <H1 className="bp3-heading">
        {!!props.scene.title ? props.scene.title : TextUtils.fileNameFromPath(props.scene.path)}
      </H1>
      {!!props.scene.date ? <H4>{props.scene.date}</H4> : ""}
      {renderDetails()}
      {renderTags()}
    </>
  );
};
