import React, { useMemo } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { Form, Row, Col } from "react-bootstrap";
import { Group, GroupSelect } from "src/components/Movies/MovieSelect";
import cx from "classnames";

export type MovieSceneIndexMap = Map<string, number | undefined>;

export interface IGroupEntry {
  movie: Group;
  scene_index?: GQL.InputMaybe<number> | undefined;
}

export interface IProps {
  value: IGroupEntry[];
  onUpdate: (input: IGroupEntry[]) => void;
}

export const SceneGroupTable: React.FC<IProps> = (props) => {
  const { value, onUpdate } = props;

  const intl = useIntl();

  const groupIDs = useMemo(() => value.map((m) => m.movie.id), [value]);

  const updateFieldChanged = (index: number, sceneIndex: number | null) => {
    const newValues = value.map((existing, i) => {
      if (i === index) {
        return {
          ...existing,
          scene_index: sceneIndex,
        };
      }
      return existing;
    });

    onUpdate(newValues);
  };

  function onGroupSet(index: number, groups: Group[]) {
    if (!groups.length) {
      // remove this entry
      const newValues = value.filter((_, i) => i !== index);
      onUpdate(newValues);
      return;
    }

    const group = groups[0];

    const newValues = value.map((existing, i) => {
      if (i === index) {
        return {
          ...existing,
          movie: group,
        };
      }
      return existing;
    });

    onUpdate(newValues);
  }

  function onNewGroupSet(groups: Group[]) {
    if (!groups.length) {
      return;
    }

    const group = groups[0];

    const newValues = [
      ...value,
      {
        movie: group,
        scene_index: null,
      },
    ];

    onUpdate(newValues);
  }

  function renderTableData() {
    return (
      <>
        {value.map((m, i) => (
          <Row key={m.movie.id} className="group-row">
            <Col xs={9}>
              <GroupSelect
                onSelect={(items) => onGroupSet(i, items)}
                values={[m.movie!]}
                excludeIds={groupIDs}
              />
            </Col>
            <Col xs={3}>
              <Form.Control
                className="text-input"
                type="number"
                value={m.scene_index ?? ""}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  updateFieldChanged(
                    i,
                    e.currentTarget.value === ""
                      ? null
                      : Number.parseInt(e.currentTarget.value, 10)
                  );
                }}
              />
            </Col>
          </Row>
        ))}
        <Row className="group-row">
          <Col xs={12}>
            <GroupSelect
              onSelect={(items) => onNewGroupSet(items)}
              values={[]}
              excludeIds={groupIDs}
            />
          </Col>
        </Row>
      </>
    );
  }

  return (
    <div className={cx("group-table", { "no-groups": !value.length })}>
      <Row className="group-table-header">
        <Col xs={9}></Col>
        <Form.Label column xs={3} className="group-scene-number-header">
          {intl.formatMessage({ id: "group_scene_number" })}
        </Form.Label>
      </Row>
      {renderTableData()}
    </div>
  );
};
