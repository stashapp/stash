/* eslint-disable jsx-a11y/control-has-associated-label */

import React from "react";
import { Button, Table } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { Icon, TruncatedText } from "src/components/Shared";
import { NavUtils } from "src/utils";

interface IPerformerListTableProps {
  performers: GQL.PerformerDataFragment[];
}

export const PerformerListTable: React.FC<IPerformerListTableProps> = (
  props: IPerformerListTableProps
) => {
  const renderPerformerRow = (performer: GQL.PerformerDataFragment) => (
    <tr key={performer.id}>
      <td>
        <Link to={`/performers/${performer.id}`}>
          <img
            className="image-thumbnail"
            alt={performer.name ?? ""}
            src={performer.image_path ?? ""}
          />
        </Link>
      </td>
      <td className="text-left">
        <Link to={`/performers/${performer.id}`}>
          <h5>
            <TruncatedText text={performer.name} />
          </h5>
        </Link>
      </td>
      <td>{performer.aliases ? performer.aliases : ""}</td>
      <td>
        {performer.favorite && (
          <Button disabled className="favorite">
            <Icon icon="heart" />
          </Button>
        )}
      </td>
      <td>
        <Link to={NavUtils.makePerformerScenesUrl(performer)}>
          <h6>{performer.scene_count}</h6>
        </Link>
      </td>
      <td>{performer.birthdate}</td>
      <td>{performer.height}</td>
    </tr>
  );

  return (
    <div className="row justify-content-center table-list">
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
        <tbody>{props.performers.map(renderPerformerRow)}</tbody>
      </Table>
    </div>
  );
};
