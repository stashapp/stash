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
import { StudiosCriterion } from "src/models/list-filter/criteria/studios";
import { PerformersCriterion } from "src/models/list-filter/criteria/performers";

interface IPerformerListTableProps {
  performers: GQL.PerformerDataFragment[];
  extraCriteria?: StudiosCriterion;
  extraPerformerFilter?: PerformersCriterion;
  useFilteredCounts?: Boolean;
  filteredCounts?: GQL.FilteredCountsDataFragment[];
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

  const PerformerRow = (row: {
    performer: GQL.PerformerDataFragment;
    filteredCounts?: GQL.FilteredCountsDataFragment;
  }) => (
    <tr key={row.performer.id}>
      <td>
        <Link to={`/performers/${row.performer.id}`}>
          <img
            className="image-thumbnail"
            alt={row.performer.name ?? ""}
            src={row.performer.image_path ?? ""}
          />
        </Link>
      </td>
      <td className="text-left">
        <Link to={`/performers/${row.performer.id}`}>
          <h5>
            {row.performer.name}
            {row.performer.disambiguation && (
              <span className="performer-disambiguation">
                {` (${row.performer.disambiguation})`}
              </span>
            )}
          </h5>
        </Link>
      </td>
      <td>
        {row.performer.alias_list ? row.performer.alias_list.join(", ") : ""}
      </td>
      <td>
        {row.performer.favorite && (
          <Button disabled className="favorite">
            <Icon icon={faHeart} />
          </Button>
        )}
      </td>
      <td>
        <Link
          to={NavUtils.makePerformerScenesUrl(
            row.performer,
            props.useFilteredCounts ? props.extraCriteria : undefined,
            props.useFilteredCounts ? props.extraPerformerFilter : undefined
          )}
        >
          <h6>
            {props.useFilteredCounts &&
            row.filteredCounts?.scene_count_filtered != null
              ? row.filteredCounts.scene_count_filtered
              : row.performer.scene_count}
          </h6>
        </Link>
      </td>
      <td>
        <Link
          to={NavUtils.makePerformerImagesUrl(
            row.performer,
            props.useFilteredCounts ? props.extraCriteria : undefined,
            props.useFilteredCounts ? props.extraPerformerFilter : undefined
          )}
        >
          <h6>
            {props.useFilteredCounts &&
            row.filteredCounts?.image_count_filtered != null
              ? row.filteredCounts.image_count_filtered
              : row.performer.image_count}
          </h6>
        </Link>
      </td>
      <td>
        <Link
          to={NavUtils.makePerformerGalleriesUrl(
            row.performer,
            props.useFilteredCounts ? props.extraCriteria : undefined,
            props.useFilteredCounts ? props.extraPerformerFilter : undefined
          )}
        >
          <h6>
            {props.useFilteredCounts &&
            row.filteredCounts?.gallery_count_filtered != null
              ? row.filteredCounts.gallery_count_filtered
              : row.performer.gallery_count}
          </h6>
        </Link>
      </td>
      <td>
        <h6>{row.performer.o_counter}</h6>
      </td>
      <td>{row.performer.birthdate}</td>
      <td>
        {!!row.performer.height_cm && formatHeight(row.performer.height_cm)}
      </td>
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
        <tbody>
          {props.performers.map((p) => (
            <PerformerRow
              key={p.id}
              performer={p}
              filteredCounts={props.filteredCounts?.find((c) => c.id === p.id)}
            />
          ))}
        </tbody>
      </Table>
    </div>
  );
};
