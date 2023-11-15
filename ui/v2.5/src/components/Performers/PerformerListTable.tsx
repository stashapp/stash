/* eslint-disable jsx-a11y/control-has-associated-label */

import React from "react";
import { useIntl } from "react-intl";
import { Button, Table } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import NavUtils from "src/utils/navigation";
import { faHeart } from "@fortawesome/free-solid-svg-icons";
import { cmToImperial } from "src/utils/units";

interface IPerformerListTableProps {
  performers: GQL.PerformerDataFragment[];
}

export const PerformerListTable: React.FC<IPerformerListTableProps> = (
  props: IPerformerListTableProps
) => {
  const intl = useIntl();

  const formatHeight = (height?: number | null) => {
    if (!height) {
      return "";
    }

    const [feet, inches] = cmToImperial(height);

    return (
      <span className="performer-height">
        <span className="height-metric">
          {intl.formatNumber(height, {
            style: "unit",
            unit: "centimeter",
            unitDisplay: "short",
          })}
        </span>
        <span className="height-imperial">
          {intl.formatNumber(feet, {
            style: "unit",
            unit: "foot",
            unitDisplay: "narrow",
          })}
          {intl.formatNumber(inches, {
            style: "unit",
            unit: "inch",
            unitDisplay: "narrow",
          })}
        </span>
      </span>
    );
  };

  const renderPerformerRow = (performer: GQL.PerformerDataFragment) => (
    <tr key={performer.id}>
      <td>
        <Link to={`/performers/${performer.id}`}>
          <img
            loading="lazy"
            className="image-thumbnail"
            alt={performer.name ?? ""}
            src={performer.image_path ?? ""}
          />
        </Link>
      </td>
      <td className="text-left">
        <Link to={`/performers/${performer.id}`}>
          <h5>
            {performer.name}
            {performer.disambiguation && (
              <span className="performer-disambiguation">
                {` (${performer.disambiguation})`}
              </span>
            )}
          </h5>
        </Link>
      </td>
      <td>{performer.alias_list ? performer.alias_list.join(", ") : ""}</td>
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
      <td>
        <h6>{performer.o_counter}</h6>
      </td>
      <td>{performer.birthdate}</td>
      <td>{!!performer.height_cm && formatHeight(performer.height_cm)}</td>
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
            <th>{intl.formatMessage({ id: "o_counter" })}</th>
            <th>{intl.formatMessage({ id: "birthdate" })}</th>
            <th>{intl.formatMessage({ id: "height" })}</th>
          </tr>
        </thead>
        <tbody>{props.performers.map(renderPerformerRow)}</tbody>
      </Table>
    </div>
  );
};
