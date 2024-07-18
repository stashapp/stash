import React, { useMemo } from "react";
import { FormattedMessage } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { Form, Row, Col } from "react-bootstrap";
import { Group, GroupSelect } from "src/components/Groups/GroupSelect";
import cx from "classnames";

export type GroupSceneIndexMap = Map<string, number | undefined>;

export interface IContainingGroupEntry {
  group: Group;
  description?: GQL.InputMaybe<string> | undefined;
}

export const ContainingGroupTable: React.FC<{
  value: IContainingGroupEntry[];
  onUpdate: (input: IContainingGroupEntry[]) => void;
}> = (props) => {
  const { value, onUpdate } = props;

  const groupIDs = useMemo(() => value.map((m) => m.group.id), [value]);

  const updateFieldChanged = (index: number, description: string | null) => {
    const newValues = value.map((existing, i) => {
      if (i === index) {
        return {
          ...existing,
          description,
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
          group: group,
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
        group: group,
        scene_index: null,
      },
    ];

    onUpdate(newValues);
  }

  return (
    <div className={cx("group-table", { "no-groups": !value.length })}>
      <Row className="group-table-header">
        <Col xs={9}></Col>
        <Form.Label column xs={3} className="group-scene-number-header">
          <FormattedMessage id="description" />
        </Form.Label>
      </Row>
      {value.map((m, i) => (
        <Row key={m.group.id} className="group-row">
          <Col xs={9}>
            <GroupSelect
              onSelect={(items) => onGroupSet(i, items)}
              values={[m.group!]}
              excludeIds={groupIDs}
            />
          </Col>
          <Col xs={3}>
            <Form.Control
              className="text-input"
              value={m.description ?? ""}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                updateFieldChanged(
                  i,
                  e.currentTarget.value === "" ? null : e.currentTarget.value
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
    </div>
  );
};
