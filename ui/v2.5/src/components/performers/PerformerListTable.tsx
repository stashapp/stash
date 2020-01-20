/* eslint-disable jsx-a11y/control-has-associated-label */

import React from "react";
import { Button, Table } from 'react-bootstrap';
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { Icon } from 'src/components/Shared';
import { NavUtils } from "src/utils";

interface IPerformerListTableProps {
  performers: GQL.PerformerDataFragment[];
}

export const PerformerListTable: React.FC<IPerformerListTableProps> = (props: IPerformerListTableProps) => {

  function maybeRenderFavoriteHeart(performer : GQL.PerformerDataFragment) {
    if (!performer.favorite) { return; }
    return (
      <Button disabled className="favorite">
        <Icon icon="heart" />
      </Button>
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
            <h5 className="text-truncate">
              {performer.name}
            </h5>
          </Link>
        </td>
        <td>
          {performer.aliases ? performer.aliases : ''}
        </td>
        <td>
          {maybeRenderFavoriteHeart(performer)}
        </td>
        <td>
          <Link to={NavUtils.makePerformerScenesUrl(performer)}>
            <h6>{performer.scene_count}</h6>
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
      <Table bordered striped>
        <thead>
          <tr>
            <th />
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
      </Table>
    </div>
    </>
  );
};

