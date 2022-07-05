/* eslint-disable jsx-a11y/control-has-associated-label */

import React from "react";
import { useIntl } from "react-intl";
import { Button, Table } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "src/components/Shared";
import { NavUtils } from "src/utils";
import { faHeart } from "@fortawesome/free-solid-svg-icons";

interface IPerformerListTableProps {
  performers: GQL.PerformerDataFragment[];
}

export const PerformerListTable: React.FC<IPerformerListTableProps> = (
  props: IPerformerListTableProps
) => {
  const intl = useIntl();

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
          <h5>{performer.name}</h5>
        </Link>
      </td>
      <td>{performer.aliases ? performer.aliases : ""}</td>
      <td>
        {performer.favorite && (
          <Button disabled className="favorite">
            <Icon icon={faHeart} />
          </Button>
        )}
      </td>
      <td>
        <Link to={NavUtils.makePerformerScenesUrl(performer)}>
          <h6>{performer.scene_count}</h6>
        </Link>
      </td>
      <td>
        <Link to={NavUtils.makePerformerImagesUrl(performer)}>
          <h6>{performer.image_count}</h6>
        </Link>
      </td>
      <td>
        <Link to={NavUtils.makePerformerGalleriesUrl(performer)}>
          <h6>{performer.gallery_count}</h6>
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
            <th>{intl.formatMessage({ id: "name" })}</th>
            <th>{intl.formatMessage({ id: "aliases" })}</th>
            <th>{intl.formatMessage({ id: "favourite" })}</th>
            <th>{intl.formatMessage({ id: "scene_count" })}</th>
            <th>{intl.formatMessage({ id: "image_count" })}</th>
            <th>{intl.formatMessage({ id: "gallery_count" })}</th>
            <th>{intl.formatMessage({ id: "birthdate" })}</th>
            <th>{intl.formatMessage({ id: "height" })}</th>
          </tr>
        </thead>
        <tbody>{props.performers.map(renderPerformerRow)}</tbody>
      </Table>
    </div>
  );
};
