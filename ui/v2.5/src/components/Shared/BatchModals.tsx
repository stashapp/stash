import React, { useMemo, useRef, useState } from "react";
import { Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";

import { ModalComponent } from "src/components/Shared/Modal";
import { faStar, faTags } from "@fortawesome/free-solid-svg-icons";

interface IEntityWithStashIDs {
  stash_ids: { endpoint: string }[];
}

interface IBatchUpdateModalProps {
  entities: IEntityWithStashIDs[];
  isIdle: boolean;
  selectedEndpoint: { endpoint: string; index: number };
  allCount: number | undefined;
  onBatchUpdate: (queryAll: boolean, refresh: boolean) => void;
  onRefreshChange?: (refresh: boolean) => void;
  batchAddParents: boolean;
  setBatchAddParents: (addParents: boolean) => void;
  close: () => void;
  localePrefix: string;
  entityName: string;
  countVariableName: string;
}

export const BatchUpdateModal: React.FC<IBatchUpdateModalProps> = ({
  entities,
  isIdle,
  selectedEndpoint,
  allCount,
  onBatchUpdate,
  onRefreshChange,
  batchAddParents,
  setBatchAddParents,
  close,
  localePrefix,
  entityName,
  countVariableName,
}) => {
  const intl = useIntl();

  const [queryAll, setQueryAll] = useState(false);
  const [refresh, setRefreshState] = useState(false);

  const setRefresh = (value: boolean) => {
    setRefreshState(value);
    onRefreshChange?.(value);
  };

  const entityCount = useMemo(() => {
    const filteredStashIDs = entities.map((e) =>
      e.stash_ids.filter((s) => s.endpoint === selectedEndpoint.endpoint)
    );

    return queryAll
      ? allCount
      : filteredStashIDs.filter((s) =>
          refresh ? s.length > 0 : s.length === 0
        ).length;
  }, [queryAll, refresh, entities, allCount, selectedEndpoint.endpoint]);

  return (
    <ModalComponent
      show
      icon={faTags}
      header={intl.formatMessage({
        id: `${localePrefix}.update_${entityName}s`,
      })}
      accept={{
        text: intl.formatMessage({
          id: `${localePrefix}.update_${entityName}s`,
        }),
        onClick: () => onBatchUpdate(queryAll, refresh),
      }}
      cancel={{
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "danger",
        onClick: () => close(),
      }}
      disabled={!isIdle}
    >
      <Form.Group>
        <Form.Label>
          <h6>
            <FormattedMessage id={`${localePrefix}.${entityName}_selection`} />
          </h6>
        </Form.Label>
        <Form.Check
          id="query-page"
          type="radio"
          name={`${entityName}-query`}
          label={<FormattedMessage id={`${localePrefix}.current_page`} />}
          checked={!queryAll}
          onChange={() => setQueryAll(false)}
        />
        <Form.Check
          id="query-all"
          type="radio"
          name={`${entityName}-query`}
          label={intl.formatMessage({
            id: `${localePrefix}.query_all_${entityName}s_in_the_database`,
          })}
          checked={queryAll}
          onChange={() => setQueryAll(true)}
        />
      </Form.Group>
      <Form.Group>
        <Form.Label>
          <h6>
            <FormattedMessage id={`${localePrefix}.tag_status`} />
          </h6>
        </Form.Label>
        <Form.Check
          id={`untagged-${entityName}s`}
          type="radio"
          name={`${entityName}-refresh`}
          label={intl.formatMessage({
            id: `${localePrefix}.untagged_${entityName}s`,
          })}
          checked={!refresh}
          onChange={() => setRefresh(false)}
        />
        <Form.Text>
          <FormattedMessage
            id={`${localePrefix}.updating_untagged_${entityName}s_description`}
          />
        </Form.Text>
        <Form.Check
          id={`tagged-${entityName}s`}
          type="radio"
          name={`${entityName}-refresh`}
          label={intl.formatMessage({
            id: `${localePrefix}.refresh_tagged_${entityName}s`,
          })}
          checked={refresh}
          onChange={() => setRefresh(true)}
        />
        <Form.Text>
          <FormattedMessage
            id={`${localePrefix}.refreshing_will_update_the_data`}
          />
        </Form.Text>
      </Form.Group>
      <div className="mt-4">
        <Form.Check
          id="add-parent"
          checked={batchAddParents}
          label={intl.formatMessage({
            id: `${localePrefix}.create_or_tag_parent_${entityName}s`,
          })}
          onChange={() => setBatchAddParents(!batchAddParents)}
        />
      </div>
      <b>
        <FormattedMessage
          id={`${localePrefix}.number_of_${entityName}s_will_be_processed`}
          values={{
            [countVariableName]: entityCount,
          }}
        />
      </b>
    </ModalComponent>
  );
};

interface IBatchAddModalProps {
  isIdle: boolean;
  onBatchAdd: (input: string) => void;
  batchAddParents: boolean;
  setBatchAddParents: (addParents: boolean) => void;
  close: () => void;
  localePrefix: string;
  entityName: string;
}

export const BatchAddModal: React.FC<IBatchAddModalProps> = ({
  isIdle,
  onBatchAdd,
  batchAddParents,
  setBatchAddParents,
  close,
  localePrefix,
  entityName,
}) => {
  const intl = useIntl();

  const inputRef = useRef<HTMLTextAreaElement | null>(null);

  return (
    <ModalComponent
      show
      icon={faStar}
      header={intl.formatMessage({
        id: `${localePrefix}.add_new_${entityName}s`,
      })}
      accept={{
        text: intl.formatMessage({
          id: `${localePrefix}.add_new_${entityName}s`,
        }),
        onClick: () => {
          if (inputRef.current) {
            onBatchAdd(inputRef.current.value);
          } else {
            close();
          }
        },
      }}
      cancel={{
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "danger",
        onClick: () => close(),
      }}
      disabled={!isIdle}
    >
      <Form.Control
        className="text-input"
        as="textarea"
        ref={inputRef}
        placeholder={intl.formatMessage({
          id: `${localePrefix}.${entityName}_names_or_stashids_separated_by_comma`,
        })}
        rows={6}
      />
      <Form.Text>
        <FormattedMessage
          id={`${localePrefix}.any_names_entered_will_be_queried`}
        />
      </Form.Text>
      <div className="mt-2">
        <Form.Check
          id="add-parent"
          checked={batchAddParents}
          label={intl.formatMessage({
            id: `${localePrefix}.create_or_tag_parent_${entityName}s`,
          })}
          onChange={() => setBatchAddParents(!batchAddParents)}
        />
      </div>
    </ModalComponent>
  );
};
