/* eslint-disable jsx-a11y/control-has-associated-label */

import React from "react";
import { useIntl } from "react-intl";
import { Button, Form, Table } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import NavUtils from "src/utils/navigation";
import { faHeart } from "@fortawesome/free-solid-svg-icons";
import { usePerformerUpdate } from "src/core/StashService";
import { useTableColumns } from "src/hooks/useTableColumns";
import { ColumnSelector, IColumn } from "../Shared/ColumnSelector";
import { RatingSystem } from "../Shared/Rating/RatingSystem";
import cx from "classnames";
import {
  FormatCircumcised,
  FormatHeight,
  FormatPenisLength,
  FormatWeight,
} from "./PerformerList";
import TextUtils from "src/utils/text";
import { getCountryByISO } from "src/utils/country";

interface IPerformerListTableProps {
  performers: GQL.PerformerDataFragment[];
  selectedIds: Set<string>;
  onSelectChange: (id: string, selected: boolean, shiftKey: boolean) => void;
}

const TABLE_NAME = "performers";

export const PerformerListTable: React.FC<IPerformerListTableProps> = (
  props: IPerformerListTableProps
) => {
  const intl = useIntl();

  const imageCol = {
    value: "image",
    label: intl.formatMessage({ id: "image" }),
  };
  const nameCol = {
    value: "name",
    label: intl.formatMessage({ id: "name" }),
  };
  const aliasesCol = {
    value: "aliases",
    label: intl.formatMessage({ id: "aliases" }),
  };
  const genderCol = {
    value: "gender",
    label: intl.formatMessage({ id: "gender" }),
  };
  const ratingCol = {
    value: "rating",
    label: intl.formatMessage({ id: "rating" }),
  };
  const ageCol = {
    value: "age",
    label: intl.formatMessage({ id: "age" }),
  };
  const deathdateCol = {
    value: "death_date",
    label: intl.formatMessage({ id: "death_date" }),
  };
  const favoriteCol = {
    value: "favourite",
    label: intl.formatMessage({ id: "favourite" }),
  };
  const countryCol = {
    value: "country",
    label: intl.formatMessage({ id: "country" }),
  };
  const ethnicityCol = {
    value: "ethnicity",
    label: intl.formatMessage({ id: "ethnicity" }),
  };
  const hairColorCol = {
    value: "hair_color",
    label: intl.formatMessage({ id: "hair_color" }),
  };
  const eyeColorCol = {
    value: "eye_color",
    label: intl.formatMessage({ id: "eye_color" }),
  };
  const heightCol = {
    value: "height_cm",
    label: intl.formatMessage({ id: "height_cm" }),
  };
  const weightCol = {
    value: "weight_kg",
    label: intl.formatMessage({ id: "weight_kg" }),
  };
  const penisLengthCol = {
    value: "penis_length_cm",
    label: intl.formatMessage({ id: "penis_length_cm" }),
  };
  const circumcisedCol = {
    value: "circumcised",
    label: intl.formatMessage({ id: "circumcised" }),
  };
  const measurementsCol = {
    value: "measurements",
    label: intl.formatMessage({ id: "measurements" }),
  };
  const fakeTitsCol = {
    value: "fake_tits",
    label: intl.formatMessage({ id: "fake_tits" }),
  };
  const careerLengthCol = {
    value: "career_length",
    label: intl.formatMessage({ id: "career_length" }),
  };
  const sceneCountCol = {
    value: "scene_count",
    label: intl.formatMessage({ id: "scene_count" }),
  };
  const galleryCountCol = {
    value: "gallery_count",
    label: intl.formatMessage({ id: "gallery_count" }),
  };
  const imageCountCol = {
    value: "image_count",
    label: intl.formatMessage({ id: "image_count" }),
  };
  const oCounterCol = {
    value: "o_counter",
    label: intl.formatMessage({ id: "o_counter" }),
  };
  const columns = [
    imageCol,
    nameCol,
    aliasesCol,
    ratingCol,
    genderCol,
    ageCol,
    deathdateCol,
    favoriteCol,
    countryCol,
    ethnicityCol,
    hairColorCol,
    eyeColorCol,
    heightCol,
    weightCol,
    measurementsCol,
    fakeTitsCol,
    penisLengthCol,
    circumcisedCol,
    careerLengthCol,
    sceneCountCol,
    galleryCountCol,
    imageCountCol,
    oCounterCol,
  ];
  const defaultColumns = [
    imageCol,
    nameCol,
    aliasesCol,
    genderCol,
    ratingCol,
    ageCol,
    favoriteCol,
    countryCol,
    ethnicityCol,
    careerLengthCol,
    sceneCountCol,
    galleryCountCol,
    imageCountCol,
    oCounterCol,
  ].map((c) => c.value);

  const [updatePerformer] = usePerformerUpdate();
  const selectedColumns = useTableColumns(TABLE_NAME, defaultColumns);

  function maybeRenderColHead(column: IColumn) {
    if (selectedColumns[column.value]) {
      return <th className={`${column.value}-head`}>{column.label}</th>;
    }
  }

  const maybeRenderCell = (column: IColumn, cell: React.ReactNode) => {
    if (selectedColumns[column.value]) return cell;
  };

  function setRating(v: number | null, performerId: string) {
    if (performerId) {
      updatePerformer({
        variables: {
          input: {
            id: performerId,
            rating100: v,
          },
        },
      });
    }
  }

  function setFavorite(v: boolean, performerId: string) {
    if (performerId) {
      updatePerformer({
        variables: {
          input: {
            id: performerId,
            favorite: v,
          },
        },
      });
    }
  }

  const ImageCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${imageCol.value}-data`}>
      <Link to={`/performers/${performer.id}`}>
        <img
          loading="lazy"
          className="image-thumbnail"
          alt={performer.name ?? ""}
          src={performer.image_path ?? ""}
        />
      </Link>
    </td>
  );

  const NameCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${nameCol.value}-data`} title={performer.name}>
      <Link to={`/performers/${performer.id}`}>
        <div className="ellips-data">
          {performer.name}
          {performer.disambiguation && (
            <span className="performer-disambiguation">
              {` (${performer.disambiguation})`}
            </span>
          )}
        </div>
      </Link>
    </td>
  );

  const AliasesCell = (performer: GQL.PerformerDataFragment) => {
    let aliases = performer.alias_list ? performer.alias_list.join(", ") : "";
    return (
      <td className={`${nameCol.value}-data`} title={aliases}>
        <span className="ellips-data">{aliases}</span>
      </td>
    );
  };

  const GenderCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${genderCol.value}-data`}>
      {performer.gender
        ? intl.formatMessage({ id: "gender_types." + performer.gender })
        : ""}
    </td>
  );

  const RatingCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${ratingCol.value}-data`}>
      <RatingSystem
        value={performer.rating100}
        onSetRating={(value) => setRating(value, performer.id)}
      />
    </td>
  );

  const AgeCell = (performer: GQL.PerformerDataFragment) => (
    <td
      className={`${ageCol.value}-data`}
      title={
        performer.birthdate
          ? TextUtils.formatDate(intl, performer.birthdate ?? undefined)
          : ""
      }
    >
      <span>
        {performer.birthdate
          ? TextUtils.age(performer.birthdate, performer.death_date)
          : ""}
      </span>
    </td>
  );

  const DeathdateCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${deathdateCol.value}-data`}>{performer.death_date}</td>
  );

  const FavoriteCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${favoriteCol.value}-data`}>
      <Button
        className={cx(
          "minimal",
          performer.favorite ? "favorite" : "not-favorite"
        )}
        onClick={() => setFavorite(!performer.favorite, performer.id)}
      >
        <Icon icon={faHeart} />
      </Button>
    </td>
  );

  const CountryCell = (performer: GQL.PerformerDataFragment) => {
    const { locale } = useIntl();
    return (
      <td className={`${countryCol.value}-data`}>
        <span className="ellips-data">
          {getCountryByISO(performer.country, locale)}
        </span>
      </td>
    );
  };

  const EthnicityCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${ethnicityCol.value}-data`}>{performer.ethnicity}</td>
  );

  const MeasurementsCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${measurementsCol.value}-data`}>
      <span className="ellips-data">{performer.measurements}</span>
    </td>
  );

  const FakeTitsCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${fakeTitsCol.value}-data`}>{performer.fake_tits}</td>
  );

  const PenisLengthCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${penisLengthCol.value}-data`}>
      {FormatPenisLength(performer.penis_length)}
    </td>
  );

  const CircumcisedCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${circumcisedCol.value}-data`}>
      {FormatCircumcised(performer.circumcised)}
    </td>
  );

  const HairColorCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${hairColorCol.value}-data`}>
      <span className="ellips-data">{performer.hair_color}</span>
    </td>
  );

  const EyeColorCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${eyeColorCol.value}-data`}>{performer.eye_color}</td>
  );

  const HeightCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${heightCol.value}-data`}>
      {FormatHeight(performer.height_cm)}
    </td>
  );

  const WeightCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${weightCol.value}-data`}>
      {FormatWeight(performer.weight)}
    </td>
  );

  const CareerLengthCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${careerLengthCol.value}-data`}>
      <span className="ellips-data">{performer.career_length}</span>
    </td>
  );

  const SceneCountCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${sceneCountCol.value}-data`}>
      <Link to={NavUtils.makePerformerScenesUrl(performer)}>
        <span>{performer.scene_count}</span>
      </Link>
    </td>
  );

  const GalleryCountCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${galleryCountCol.value}-data`}>
      <Link to={NavUtils.makePerformerGalleriesUrl(performer)}>
        <span>{performer.gallery_count}</span>
      </Link>
    </td>
  );

  const ImageCountCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${imageCountCol.value}-data`}>
      <Link to={NavUtils.makePerformerImagesUrl(performer)}>
        <span>{performer.image_count}</span>
      </Link>
    </td>
  );

  const OCounterCell = (performer: GQL.PerformerDataFragment) => (
    <td className={`${oCounterCol.value}-data`}>{performer.o_counter}</td>
  );

  let shiftKey = false;

  const renderPerformerRow = (performer: GQL.PerformerDataFragment) => (
    <tr key={performer.id}>
      <td className="select-col">
        <label>
          <Form.Control
            type="checkbox"
            checked={props.selectedIds.has(performer.id)}
            onChange={() =>
              props.onSelectChange(
                performer.id,
                !props.selectedIds.has(performer.id),
                shiftKey
              )
            }
            onClick={(
              event: React.MouseEvent<HTMLInputElement, MouseEvent>
            ) => {
              shiftKey = event.shiftKey;
              event.stopPropagation();
            }}
          />
        </label>
      </td>
      {maybeRenderCell(imageCol, ImageCell(performer))}
      {maybeRenderCell(nameCol, NameCell(performer))}
      {maybeRenderCell(aliasesCol, AliasesCell(performer))}
      {maybeRenderCell(ratingCol, RatingCell(performer))}
      {maybeRenderCell(genderCol, GenderCell(performer))}
      {maybeRenderCell(ageCol, AgeCell(performer))}
      {maybeRenderCell(deathdateCol, DeathdateCell(performer))}
      {maybeRenderCell(favoriteCol, FavoriteCell(performer))}
      {maybeRenderCell(countryCol, CountryCell(performer))}
      {maybeRenderCell(ethnicityCol, EthnicityCell(performer))}
      {maybeRenderCell(hairColorCol, HairColorCell(performer))}
      {maybeRenderCell(eyeColorCol, EyeColorCell(performer))}
      {maybeRenderCell(heightCol, HeightCell(performer))}
      {maybeRenderCell(weightCol, WeightCell(performer))}
      {maybeRenderCell(measurementsCol, MeasurementsCell(performer))}
      {maybeRenderCell(fakeTitsCol, FakeTitsCell(performer))}
      {maybeRenderCell(penisLengthCol, PenisLengthCell(performer))}
      {maybeRenderCell(circumcisedCol, CircumcisedCell(performer))}
      {maybeRenderCell(careerLengthCol, CareerLengthCell(performer))}
      {maybeRenderCell(sceneCountCol, SceneCountCell(performer))}
      {maybeRenderCell(galleryCountCol, GalleryCountCell(performer))}
      {maybeRenderCell(imageCountCol, ImageCountCell(performer))}
      {maybeRenderCell(oCounterCol, OCounterCell(performer))}
    </tr>
  );

  return (
    <div className="row justify-content-center table-list">
      <Table striped bordered>
        <thead>
          <tr>
            <th className="select-col">
              <div
                className="d-inline-block"
                data-toggle="popover"
                data-trigger="focus"
              >
                <ColumnSelector
                  tableName={TABLE_NAME}
                  columns={columns}
                  defaultColumns={defaultColumns}
                />
              </div>
            </th>
            {maybeRenderColHead(imageCol)}
            {maybeRenderColHead(nameCol)}
            {maybeRenderColHead(aliasesCol)}
            {maybeRenderColHead(ratingCol)}
            {maybeRenderColHead(genderCol)}
            {maybeRenderColHead(ageCol)}
            {maybeRenderColHead(deathdateCol)}
            {maybeRenderColHead(favoriteCol)}
            {maybeRenderColHead(countryCol)}
            {maybeRenderColHead(ethnicityCol)}
            {maybeRenderColHead(hairColorCol)}
            {maybeRenderColHead(eyeColorCol)}
            {maybeRenderColHead(heightCol)}
            {maybeRenderColHead(weightCol)}
            {maybeRenderColHead(measurementsCol)}
            {maybeRenderColHead(fakeTitsCol)}
            {maybeRenderColHead(penisLengthCol)}
            {maybeRenderColHead(circumcisedCol)}
            {maybeRenderColHead(careerLengthCol)}
            {maybeRenderColHead(sceneCountCol)}
            {maybeRenderColHead(galleryCountCol)}
            {maybeRenderColHead(imageCountCol)}
            {maybeRenderColHead(oCounterCol)}
          </tr>
          <tr>
            <th className="border-row" colSpan={100}></th>
          </tr>
        </thead>
        <tbody>{props.performers.map(renderPerformerRow)}</tbody>
      </Table>
    </div>
  );
};
