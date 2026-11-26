import React, { useContext, useMemo } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { Form, Row, Col } from "react-bootstrap";
import { Group, GroupSelect } from "src/components/Groups/GroupSelect";
import cx from "classnames";
import { NumberField } from "src/utils/form";
import { LockContext } from "src/components/Shared/LockFieldButton";

export type GroupSceneIndexMap = Map<string, number | undefined>;

export interface IGroupEntry {
  group: Group;
  scene_index?: GQL.InputMaybe<number> | undefined;
}

export interface IProps {
  value: IGroupEntry[];
  onUpdate: (input: IGroupEntry[]) => void;
}

export const SceneGroupTable: React.FC<IProps> = (props) => {
  const { value, onUpdate } = props;

  const intl = useIntl();
  const lock = useContext(LockContext);
  const locked = !!lock?.locked;

  // only the first select shows the lock toggle
  function withLock(index: number, select: React.ReactNode) {
    return index === 0 ? (
      select
    ) : (
      <LockContext.Provider value={undefined}>{select}</LockContext.Provider>
    );
  }

  const groupIDs = useMemo(() => value.map((m) => m.group.id), [value]);

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

  function renderTableData() {
    return (
      <>
        {value.map((m, i) => (
          <Row key={m.group.id} className="group-row">
            <Col xs={9}>
              {withLock(
                i,
                <GroupSelect
                  onSelect={(items) => onGroupSet(i, items)}
                  values={[m.group!]}
                  excludeIds={groupIDs}
                  isDisabled={locked}
                />
              )}
            </Col>
            <Col xs={3}>
              <NumberField
                className="text-input"
                disabled={locked}
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
        {/* while locked, only shown to hold the lock toggle when there are no groups */}
        {(!locked || !value.length) && (
          <Row className="group-row">
            <Col xs={12}>
              {withLock(
                value.length,
                <GroupSelect
                  onSelect={(items) => onNewGroupSet(items)}
                  values={[]}
                  excludeIds={groupIDs}
                />
              )}
            </Col>
          </Row>
        )}
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
