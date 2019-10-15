import {
    H4,
    HTMLTable,
    H5,
    H6,
  } from "@blueprintjs/core";
  import React, { FunctionComponent } from "react";
  import { Link } from "react-router-dom";
  import * as GQL from "../../core/generated-graphql";
  import { TextUtils } from "../../utils/text";
  import { TagLink } from "../Shared/TagLink";
import { NavigationUtils } from "../../utils/navigation";
  
  interface ISceneListTableProps {
    scenes: GQL.SlimSceneDataFragment[];
  }
  
  export const SceneListTable: FunctionComponent<ISceneListTableProps> = (props: ISceneListTableProps) => {
    
    function renderDuration(scene : GQL.SlimSceneDataFragment) {
      if (scene.file.duration === undefined) { return; }
      return TextUtils.secondsToTimestamp(scene.file.duration);
    }

    function renderTags(tags : GQL.SlimSceneDataTags[]) {
      return tags.map((tag) => (
        <Link to={NavigationUtils.makeTagScenesUrl(tag)}>
          <H6>{tag.name}</H6>
        </Link>
      ));
    }

    function renderPerformers(performers : GQL.SlimSceneDataPerformers[]) {
      return performers.map((performer) => (
        <Link to={NavigationUtils.makePerformerScenesUrl(performer)}>
          <H6>{performer.name}</H6>
        </Link>
      ));
    }

    function renderStudio(studio : GQL.SlimSceneDataStudio | undefined) {
      if (!!studio) {
        return (
          <Link to={NavigationUtils.makeStudioScenesUrl(studio)}>
            <H6>{studio.name}</H6>
          </Link>
        );
      }
    }

    function renderSceneRow(scene : GQL.SlimSceneDataFragment) {
      return (
        <>
        <tr>
          <td>
            <Link to={`/scenes/${scene.id}`}>
              <H5 style={{textOverflow: "ellipsis", overflow: "hidden"}}>
                {!!scene.title ? scene.title : TextUtils.fileNameFromPath(scene.path)}
              </H5>
            </Link>
          </td>
          <td>
            {scene.rating ? scene.rating : ''}
          </td>
          <td>
            {renderDuration(scene)}
          </td>
          <td>
            {renderTags(scene.tags)}
          </td>
          <td>
            {renderPerformers(scene.performers)}
          </td>
          <td>
            {renderStudio(scene.studio)}
          </td>
        </tr>
        </>
      )
    }
  
    return (
      <>
      <div className="grid">
        <HTMLTable>
          <thead>
            <tr>
              <th>Title</th>
              <th>Rating</th>
              <th>Duration</th>
              <th>Tags</th>
              <th>Performers</th>
              <th>Studio</th>
            </tr>
          </thead>
          <tbody>
            {props.scenes.map(renderSceneRow)}
          </tbody>
        </HTMLTable>
      </div>
      </>
    );
  };
  