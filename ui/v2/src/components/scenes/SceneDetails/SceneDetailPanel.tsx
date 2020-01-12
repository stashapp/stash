import {
  H1,
  H4,
  H6,
  Tag,
  HTMLTable
 } from "@blueprintjs/core";


import React, { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../../core/generated-graphql";
import { NavigationUtils } from "../../../utils/navigation";
import { TextUtils } from "../../../utils/text";
import { TagLink } from "../../Shared/TagLink";
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
      <TagLink key={tag.id} tag={tag} />
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
    <HTMLTable style={{width: "100%"}}>
      <tbody>
      <tr>
        <td>  {SceneHelpers.maybeRenderFrontDvd(props.scene, 200, false)}</td>
        <td> {SceneHelpers.maybeRenderBackDvd(props.scene, 200, false)}</td>
      </tr>
      </tbody>
    </HTMLTable>
    
      <H1 className="bp3-heading">
        {!!props.scene.title ? props.scene.title : TextUtils.fileNameFromPath(props.scene.path)}
      </H1>
      {!!props.scene.date ? <H4>{props.scene.date}</H4> : undefined}
      {!!props.scene.rating ? <H6>Rating: {props.scene.rating}</H6> : undefined}
      {!!props.scene.file.height ? <H6>Resolution: {TextUtils.resolution(props.scene.file.height)}</H6> : undefined}
      {renderDetails()}
      {renderTags()}
    </>
  );
};
