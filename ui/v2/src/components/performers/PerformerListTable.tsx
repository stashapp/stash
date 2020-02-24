import {
    HTMLTable,
    H5,
    H6,
    Button,
  } from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import * as GQL from "../../core/generated-graphql";
import { NavigationUtils } from "../../utils/navigation";
  
interface IPerformerListTableProps {
  performers: GQL.PerformerDataFragment[];
}
  
export const PerformerListTable: FunctionComponent<IPerformerListTableProps> = (props: IPerformerListTableProps) => {
  
  function maybeRenderFavoriteHeart(performer : GQL.PerformerDataFragment) {
    if (!performer.favorite) { return; }
    return (
      <Button
        icon="heart"
        disabled={true}
        className="favorite"
        minimal={true}
      />
    );
  }

  function renderPerformerImage(performer : GQL.PerformerDataFragment) {
    const style: React.CSSProperties = {
      backgroundImage: `url('${performer.image_path}')`,
      lineHeight: 5,
      backgroundSize: "contain",
      display: "inline-block",
      backgroundPosition: "center",
      backgroundRepeat: "no-repeat",
    };

    return (
      <Link 
        className="performer-list-thumbnail"
        to={`/performers/${performer.id}`}
        style={style}/>
    )
  }

  function renderPerformerRow(performer : GQL.PerformerDataFragment) {
    return (
      <>
      <tr>
        <td>
          {renderPerformerImage(performer)}
        </td>
        <td style={{textAlign: "left"}}>
          <Link to={`/performers/${performer.id}`}>
            <H5 style={{textOverflow: "ellipsis", overflow: "hidden"}}>
              {performer.name}
            </H5>
          </Link>
        </td>
        <td>
          {performer.aliases ? performer.aliases : ''}
        </td>
        <td>
          {maybeRenderFavoriteHeart(performer)}
        </td>
        <td>
          <Link to={NavigationUtils.makePerformerScenesUrl(performer)}>
            <H6>{performer.scene_count}</H6>
          </Link>
        </td>
        <td>
          {performer.birthdate}
        </td>
        <td>
          {performer.height}
        </td>
      </tr>
      </>
    )
  }

  return (
    <>
    <div className="grid">
      <HTMLTable className="bp3-html-table bp3-html-table-bordered bp3-html-table-condensed bp3-html-table-striped bp3-interactive">
        <thead>
          <tr>
            <th></th>
            <th>Name</th>
            <th>Aliases</th>
            <th>Favourite</th>
            <th>Scene Count</th>
            <th>Birthdate</th>
            <th>Height</th>
          </tr>
        </thead>
        <tbody>
          {props.performers.map(renderPerformerRow)}
        </tbody>
      </HTMLTable>
    </div>
    </>
  );
};
  