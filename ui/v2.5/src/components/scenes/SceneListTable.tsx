import React from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";
import { TextUtils } from "../../utils/text";
import { NavigationUtils } from "../../utils/navigation";

import { Table } from 'react-bootstrap';
  
interface ISceneListTableProps {
  scenes: GQL.SlimSceneDataFragment[];
}

export const SceneListTable: React.FC<ISceneListTableProps> = (props: ISceneListTableProps) => {
  
  function renderSceneImage(scene : GQL.SlimSceneDataFragment) {
    const style: React.CSSProperties = {
      backgroundImage: `url('${scene.paths.screenshot}')`,
      lineHeight: 5,
      backgroundSize: "contain",
      display: "inline-block",
      backgroundPosition: "center",
      backgroundRepeat: "no-repeat",
    };

    return (
      <Link 
        className="scene-list-thumbnail"
        to={`/scenes/${scene.id}`}
        style={style}/>
    )
  }

  function renderDuration(scene : GQL.SlimSceneDataFragment) {
    if (scene.file.duration === undefined) { return; }
    return TextUtils.secondsToTimestamp(scene.file.duration);
  }

  function renderTags(tags : GQL.SlimSceneDataTags[]) {
    return tags.map((tag) => (
      <Link to={NavigationUtils.makeTagScenesUrl(tag)}>
        <h6>{tag.name}</h6>
      </Link>
    ));
  }

  function renderPerformers(performers : GQL.SlimSceneDataPerformers[]) {
    return performers.map((performer) => (
      <Link to={NavigationUtils.makePerformerScenesUrl(performer)}>
        <h6>{performer.name}</h6>
      </Link>
    ));
  }

  function renderStudio(studio : GQL.SlimSceneDataStudio | undefined) {
    if (!!studio) {
      return (
        <Link to={NavigationUtils.makeStudioScenesUrl(studio)}>
          <h6>{studio.name}</h6>
        </Link>
      );
    }
  }

  function renderSceneRow(scene : GQL.SlimSceneDataFragment) {
    return (
      <>
      <tr>
        <td>
          {renderSceneImage(scene)}
        </td>
        <td style={{textAlign: "left"}}>
          <Link to={`/scenes/${scene.id}`}>
            <h5 className="text-truncate">
              {!!scene.title ? scene.title : TextUtils.fileNameFromPath(scene.path)}
            </h5>
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
      <Table striped bordered>
        <thead>
          <tr>
            <th></th>
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
      </Table>
    </div>
    </>
  );
};
  
