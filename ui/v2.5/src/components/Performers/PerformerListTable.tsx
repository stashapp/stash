/* eslint-disable jsx-a11y/control-has-associated-label */

import React from "react";
import { useIntl } from "react-intl";
import { Button } from "react-bootstrap";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "../Shared/Icon";
import NavUtils from "src/utils/navigation";
import { faHeart } from "@fortawesome/free-solid-svg-icons";
import { usePerformerUpdate } from "src/core/StashService";
import { useTableColumns } from "src/hooks/useTableColumns";
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
import { IColumn, ListTable } from "../List/ListTable";

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

  const [updatePerformer] = usePerformerUpdate();

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
    <Link to={`/performers/${performer.id}`}>
      <img
        loading="lazy"
        className="image-thumbnail"
        alt={performer.name ?? ""}
        src={performer.image_path ?? ""}
      />
    </Link>
  );

  const NameCell = (performer: GQL.PerformerDataFragment) => (
    <Link to={`/performers/${performer.id}`}>
      <div className="ellips-data" title={performer.name}>
        {performer.name}
        {performer.disambiguation && (
          <span className="performer-disambiguation">
            {` (${performer.disambiguation})`}
          </span>
        )}
      </div>
    </Link>
  );

  const AliasesCell = (performer: GQL.PerformerDataFragment) => {
    let aliases = performer.alias_list ? performer.alias_list.join(", ") : "";
    return (
      <span className="ellips-data" title={aliases}>
        {aliases}
      </span>
    );
  };

  const GenderCell = (performer: GQL.PerformerDataFragment) => (
    <>
      {performer.gender
        ? intl.formatMessage({ id: "gender_types." + performer.gender })
        : ""}
    </>
  );

  const RatingCell = (performer: GQL.PerformerDataFragment) => (
    <RatingSystem
      value={performer.rating100}
      onSetRating={(value) => setRating(value, performer.id)}
      clickToRate
    />
  );

  const AgeCell = (performer: GQL.PerformerDataFragment) => (
    <span
      title={
        performer.birthdate
          ? TextUtils.formatFuzzyDate(intl, performer.birthdate ?? undefined)
          : ""
      }
    >
      {performer.birthdate
        ? TextUtils.age(performer.birthdate, performer.death_date)
        : ""}
    </span>
  );

  const DeathdateCell = (performer: GQL.PerformerDataFragment) => (
    <>{performer.death_date}</>
  );

  const FavoriteCell = (performer: GQL.PerformerDataFragment) => (
    <Button
      className={cx(
        "minimal",
        performer.favorite ? "favorite" : "not-favorite"
      )}
      onClick={() => setFavorite(!performer.favorite, performer.id)}
    >
      <Icon icon={faHeart} />
    </Button>
  );

  const CountryCell = (performer: GQL.PerformerDataFragment) => {
    const { locale } = useIntl();
    return (
      <span className="ellips-data">
        {getCountryByISO(performer.country, locale)}
      </span>
    );
  };

  const EthnicityCell = (performer: GQL.PerformerDataFragment) => (
    <>{performer.ethnicity}</>
  );

  const MeasurementsCell = (performer: GQL.PerformerDataFragment) => (
    <span className="ellips-data">{performer.measurements}</span>
  );

  const FakeTitsCell = (performer: GQL.PerformerDataFragment) => (
    <>{performer.fake_tits}</>
  );

  const PenisLengthCell = (performer: GQL.PerformerDataFragment) => (
    <>{FormatPenisLength(performer.penis_length)}</>
  );

  const CircumcisedCell = (performer: GQL.PerformerDataFragment) => (
    <>{FormatCircumcised(performer.circumcised)}</>
  );

  const HairColorCell = (performer: GQL.PerformerDataFragment) => (
    <span className="ellips-data">{performer.hair_color}</span>
  );

  const EyeColorCell = (performer: GQL.PerformerDataFragment) => (
    <>{performer.eye_color}</>
  );

  const HeightCell = (performer: GQL.PerformerDataFragment) => (
    <>{FormatHeight(performer.height_cm)}</>
  );

  const WeightCell = (performer: GQL.PerformerDataFragment) => (
    <>{FormatWeight(performer.weight)}</>
  );

  const CareerLengthCell = (performer: GQL.PerformerDataFragment) => (
    <span className="ellips-data">{performer.career_length}</span>
  );

  const SceneCountCell = (performer: GQL.PerformerDataFragment) => (
    <Link to={NavUtils.makePerformerScenesUrl(performer)}>
      <span>{performer.scene_count}</span>
    </Link>
  );

  const GalleryCountCell = (performer: GQL.PerformerDataFragment) => (
    <Link to={NavUtils.makePerformerGalleriesUrl(performer)}>
      <span>{performer.gallery_count}</span>
    </Link>
  );

  const ImageCountCell = (performer: GQL.PerformerDataFragment) => (
    <Link to={NavUtils.makePerformerImagesUrl(performer)}>
      <span>{performer.image_count}</span>
    </Link>
  );

  const OCounterCell = (performer: GQL.PerformerDataFragment) => (
    <>{performer.o_counter}</>
  );

  interface IColumnSpec {
    value: string;
    label: string;
    defaultShow?: boolean;
    mandatory?: boolean;
    render?: (
      scene: GQL.PerformerDataFragment,
      index: number
    ) => React.ReactNode;
  }

  const allColumns: IColumnSpec[] = [
    {
      value: "image",
      label: intl.formatMessage({ id: "image" }),
      defaultShow: true,
      render: ImageCell,
    },
    {
      value: "name",
      label: intl.formatMessage({ id: "name" }),
      mandatory: true,
      defaultShow: true,
      render: NameCell,
    },
    {
      value: "aliases",
      label: intl.formatMessage({ id: "aliases" }),
      defaultShow: true,
      render: AliasesCell,
    },
    {
      value: "gender",
      label: intl.formatMessage({ id: "gender" }),
      defaultShow: true,
      render: GenderCell,
    },
    {
      value: "rating",
      label: intl.formatMessage({ id: "rating" }),
      defaultShow: true,
      render: RatingCell,
    },
    {
      value: "age",
      label: intl.formatMessage({ id: "age" }),
      defaultShow: true,
      render: AgeCell,
    },
    {
      value: "death_date",
      label: intl.formatMessage({ id: "death_date" }),
      render: DeathdateCell,
    },
    {
      value: "favourite",
      label: intl.formatMessage({ id: "favourite" }),
      defaultShow: true,
      render: FavoriteCell,
    },
    {
      value: "country",
      label: intl.formatMessage({ id: "country" }),
      defaultShow: true,
      render: CountryCell,
    },
    {
      value: "ethnicity",
      label: intl.formatMessage({ id: "ethnicity" }),
      defaultShow: true,
      render: EthnicityCell,
    },
    {
      value: "hair_color",
      label: intl.formatMessage({ id: "hair_color" }),
      render: HairColorCell,
    },
    {
      value: "eye_color",
      label: intl.formatMessage({ id: "eye_color" }),
      render: EyeColorCell,
    },
    {
      value: "height_cm",
      label: intl.formatMessage({ id: "height_cm" }),
      render: HeightCell,
    },
    {
      value: "weight_kg",
      label: intl.formatMessage({ id: "weight_kg" }),
      render: WeightCell,
    },
    {
      value: "penis_length_cm",
      label: intl.formatMessage({ id: "penis_length_cm" }),
      render: PenisLengthCell,
    },
    {
      value: "circumcised",
      label: intl.formatMessage({ id: "circumcised" }),
      render: CircumcisedCell,
    },
    {
      value: "measurements",
      label: intl.formatMessage({ id: "measurements" }),
      render: MeasurementsCell,
    },
    {
      value: "fake_tits",
      label: intl.formatMessage({ id: "fake_tits" }),
      render: FakeTitsCell,
    },
    {
      value: "career_length",
      label: intl.formatMessage({ id: "career_length" }),
      defaultShow: true,
      render: CareerLengthCell,
    },
    {
      value: "scene_count",
      label: intl.formatMessage({ id: "scene_count" }),
      defaultShow: true,
      render: SceneCountCell,
    },
    {
      value: "gallery_count",
      label: intl.formatMessage({ id: "gallery_count" }),
      defaultShow: true,
      render: GalleryCountCell,
    },
    {
      value: "image_count",
      label: intl.formatMessage({ id: "image_count" }),
      defaultShow: true,
      render: ImageCountCell,
    },
    {
      value: "o_counter",
      label: intl.formatMessage({ id: "o_count" }),
      defaultShow: true,
      render: OCounterCell,
    },
  ];

  const defaultColumns = allColumns
    .filter((col) => col.defaultShow)
    .map((col) => col.value);

  const { selectedColumns, saveColumns } = useTableColumns(
    TABLE_NAME,
    defaultColumns
  );

  const columnRenderFuncs: Record<
    string,
    (scene: GQL.PerformerDataFragment, index: number) => React.ReactNode
  > = {};
  allColumns.forEach((col) => {
    if (col.render) {
      columnRenderFuncs[col.value] = col.render;
    }
  });

  function renderCell(
    column: IColumn,
    performer: GQL.PerformerDataFragment,
    index: number
  ) {
    const render = columnRenderFuncs[column.value];

    if (render) return render(performer, index);
  }

  return (
    <ListTable
      className="performer-table"
      items={props.performers}
      allColumns={allColumns}
      columns={selectedColumns}
      setColumns={(c) => saveColumns(c)}
      selectedIds={props.selectedIds}
      onSelectChange={props.onSelectChange}
      renderCell={renderCell}
    />
  );
};
