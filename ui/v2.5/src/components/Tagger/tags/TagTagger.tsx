import React, { useEffect, useMemo, useRef, useState } from "react";
import { Button, Card, Form, InputGroup, ProgressBar } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Link } from "react-router-dom";
import { HashLink } from "react-router-hash-link";

import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { ModalComponent } from "src/components/Shared/Modal";
import {
  stashBoxTagQuery,
  useJobsSubscribe,
  mutateStashBoxBatchTagTag,
  getClient,
} from "src/core/StashService";
import { Manual } from "src/components/Help/Manual";
import { useConfigurationContext } from "src/hooks/Config";

import StashSearchResult from "./StashSearchResult";
import TaggerConfig from "../TaggerConfig";
import { ITaggerConfig, TAG_FIELDS } from "../constants";
import { useUpdateTag } from "../queries";
import { faStar, faTags } from "@fortawesome/free-solid-svg-icons";
import { ExternalLink } from "src/components/Shared/ExternalLink";
import { mergeTagStashIDs } from "../utils";
import { separateNamesAndStashIds } from "src/utils/stashIds";
import { useTaggerConfig } from "../config";

type JobFragment = Pick<
  GQL.Job,
  "id" | "status" | "subTasks" | "description" | "progress"
>;

const CLASSNAME = "StudioTagger";

interface ITagBatchUpdateModal {
  tags: GQL.TagListDataFragment[];
  isIdle: boolean;
  selectedEndpoint: { endpoint: string; index: number };
  onBatchUpdate: (queryAll: boolean, refresh: boolean) => void;
  close: () => void;
}

const TagBatchUpdateModal: React.FC<ITagBatchUpdateModal> = ({
  tags,
  isIdle,
  selectedEndpoint,
  onBatchUpdate,
  close,
}) => {
  const intl = useIntl();

  const [queryAll, setQueryAll] = useState(false);

  const [refresh, setRefresh] = useState(false);
  const { data: allTags } = GQL.useFindTagsQuery({
    variables: {
      tag_filter: {
        stash_id_endpoint: {
          endpoint: selectedEndpoint.endpoint,
          modifier: refresh
            ? GQL.CriterionModifier.NotNull
            : GQL.CriterionModifier.IsNull,
        },
      },
      filter: {
        per_page: 0,
      },
    },
  });

  const tagCount = useMemo(() => {
    const filteredStashIDs = tags.map((t) =>
      t.stash_ids.filter((s) => s.endpoint === selectedEndpoint.endpoint)
    );

    return queryAll
      ? allTags?.findTags.count
      : filteredStashIDs.filter((s) =>
          refresh ? s.length > 0 : s.length === 0
        ).length;
  }, [queryAll, refresh, tags, allTags, selectedEndpoint.endpoint]);

  return (
    <ModalComponent
      show
      icon={faTags}
      header={intl.formatMessage({
        id: "tag_tagger.update_tags",
      })}
      accept={{
        text: intl.formatMessage({
          id: "tag_tagger.update_tags",
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
            <FormattedMessage id="tag_tagger.tag_selection" />
          </h6>
        </Form.Label>
        <Form.Check
          id="query-page"
          type="radio"
          name="tag-query"
          label={<FormattedMessage id="tag_tagger.current_page" />}
          checked={!queryAll}
          onChange={() => setQueryAll(false)}
        />
        <Form.Check
          id="query-all"
          type="radio"
          name="tag-query"
          label={intl.formatMessage({
            id: "tag_tagger.query_all_tags_in_the_database",
          })}
          checked={queryAll}
          onChange={() => setQueryAll(true)}
        />
      </Form.Group>
      <Form.Group>
        <Form.Label>
          <h6>
            <FormattedMessage id="tag_tagger.tag_status" />
          </h6>
        </Form.Label>
        <Form.Check
          id="untagged-tags"
          type="radio"
          name="tag-refresh"
          label={intl.formatMessage({
            id: "tag_tagger.untagged_tags",
          })}
          checked={!refresh}
          onChange={() => setRefresh(false)}
        />
        <Form.Text>
          <FormattedMessage id="tag_tagger.updating_untagged_tags_description" />
        </Form.Text>
        <Form.Check
          id="tagged-tags"
          type="radio"
          name="tag-refresh"
          label={intl.formatMessage({
            id: "tag_tagger.refresh_tagged_tags",
          })}
          checked={refresh}
          onChange={() => setRefresh(true)}
        />
        <Form.Text>
          <FormattedMessage id="tag_tagger.refreshing_will_update_the_data" />
        </Form.Text>
      </Form.Group>
      <b>
        <FormattedMessage
          id="tag_tagger.number_of_tags_will_be_processed"
          values={{
            tag_count: tagCount,
          }}
        />
      </b>
    </ModalComponent>
  );
};

interface ITagBatchAddModal {
  isIdle: boolean;
  onBatchAdd: (input: string) => void;
  close: () => void;
}

const TagBatchAddModal: React.FC<ITagBatchAddModal> = ({
  isIdle,
  onBatchAdd,
  close,
}) => {
  const intl = useIntl();

  const tagInput = useRef<HTMLTextAreaElement | null>(null);

  return (
    <ModalComponent
      show
      icon={faStar}
      header={intl.formatMessage({
        id: "tag_tagger.add_new_tags",
      })}
      accept={{
        text: intl.formatMessage({
          id: "tag_tagger.add_new_tags",
        }),
        onClick: () => {
          if (tagInput.current) {
            onBatchAdd(tagInput.current.value);
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
        ref={tagInput}
        placeholder={intl.formatMessage({
          id: "tag_tagger.tag_names_or_stashids_separated_by_comma",
        })}
        rows={6}
      />
      <Form.Text>
        <FormattedMessage id="tag_tagger.any_names_entered_will_be_queried" />
      </Form.Text>
    </ModalComponent>
  );
};

interface ITagTaggerListProps {
  tags: GQL.TagListDataFragment[];
  selectedEndpoint: { endpoint: string; index: number };
  isIdle: boolean;
  config: ITaggerConfig;
  onBatchAdd: (tagInput: string) => void;
  onBatchUpdate: (ids: string[] | undefined, refresh: boolean) => void;
}

const TagTaggerList: React.FC<ITagTaggerListProps> = ({
  tags,
  selectedEndpoint,
  isIdle,
  config,
  onBatchAdd,
  onBatchUpdate,
}) => {
  const intl = useIntl();

  const [loading, setLoading] = useState(false);
  const [searchResults, setSearchResults] = useState<
    Record<string, GQL.ScrapedSceneTagDataFragment[]>
  >({});
  const [searchErrors, setSearchErrors] = useState<
    Record<string, string | undefined>
  >({});
  const [taggedTags, setTaggedTags] = useState<
    Record<string, Partial<GQL.TagListDataFragment>>
  >({});
  const [queries, setQueries] = useState<Record<string, string>>({});

  const [showBatchAdd, setShowBatchAdd] = useState(false);
  const [showBatchUpdate, setShowBatchUpdate] = useState(false);

  const [error, setError] = useState<
    Record<string, { message?: string; details?: string } | undefined>
  >({});
  const [loadingUpdate, setLoadingUpdate] = useState<string | undefined>();

  const doBoxSearch = (tagID: string, searchVal: string) => {
    stashBoxTagQuery(searchVal, selectedEndpoint.endpoint)
      .then((queryData) => {
        const s = queryData.data?.scrapeSingleTag ?? [];
        setSearchResults({
          ...searchResults,
          [tagID]: s,
        });
        setSearchErrors({
          ...searchErrors,
          [tagID]: undefined,
        });
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
        const { [tagID]: unassign, ...results } = searchResults;
        setSearchResults(results);
        setSearchErrors({
          ...searchErrors,
          [tagID]: intl.formatMessage({
            id: "tag_tagger.network_error",
          }),
        });
      });

    setLoading(true);
  };

  const updateTag = useUpdateTag();

  const doBoxUpdate = (tagID: string, stashID: string, endpoint: string) => {
    setLoadingUpdate(stashID);
    setError({
      ...error,
      [tagID]: undefined,
    });
    stashBoxTagQuery(stashID, endpoint)
      .then(async (queryData) => {
        const data = queryData.data?.scrapeSingleTag ?? [];
        if (data.length > 0) {
          const stashboxTag = data[0];
          const updateData: GQL.TagUpdateInput = {
            id: tagID,
          };

          if (
            !(config.excludedTagFields ?? []).includes("name") &&
            stashboxTag.name
          ) {
            updateData.name = stashboxTag.name;
          }

          if (
            stashboxTag.description &&
            !(config.excludedTagFields ?? []).includes("description")
          ) {
            updateData.description = stashboxTag.description;
          }

          if (
            stashboxTag.alias_list &&
            stashboxTag.alias_list.length > 0 &&
            !(config.excludedTagFields ?? []).includes("aliases")
          ) {
            updateData.aliases = stashboxTag.alias_list;
          }

          if (stashboxTag.remote_site_id) {
            updateData.stash_ids = await mergeTagStashIDs(tagID, [
              {
                endpoint,
                stash_id: stashboxTag.remote_site_id,
              },
            ]);
          }

          const res = await updateTag(updateData);
          if (!res?.data?.tagUpdate) {
            setError({
              ...error,
              [tagID]: {
                message: `Failed to update tag`,
                details: res?.errors?.[0]?.message ?? "",
              },
            });
          }
        }
      })
      .finally(() => setLoadingUpdate(undefined));
  };

  async function handleBatchAdd(input: string) {
    onBatchAdd(input);
    setShowBatchAdd(false);
  }

  const handleBatchUpdate = (queryAll: boolean, refresh: boolean) => {
    onBatchUpdate(!queryAll ? tags.map((t) => t.id) : undefined, refresh);
    setShowBatchUpdate(false);
  };

  const handleTaggedTag = (
    tag: Pick<GQL.TagListDataFragment, "id"> &
      Partial<Omit<GQL.TagListDataFragment, "id">>
  ) => {
    setTaggedTags({
      ...taggedTags,
      [tag.id]: tag,
    });
  };

  const renderTags = () =>
    tags.map((tag) => {
      const isTagged = taggedTags[tag.id];

      const stashID = tag.stash_ids.find((s) => {
        return s.endpoint === selectedEndpoint.endpoint;
      });

      let mainContent;
      if (!isTagged && stashID !== undefined) {
        mainContent = (
          <div className="text-left">
            <h5 className="text-bold">
              <FormattedMessage id="tag_tagger.tag_already_tagged" />
            </h5>
          </div>
        );
      } else if (!isTagged && !stashID) {
        mainContent = (
          <InputGroup>
            <Form.Control
              className="text-input"
              defaultValue={tag.name ?? ""}
              onChange={(e) =>
                setQueries({
                  ...queries,
                  [tag.id]: e.currentTarget.value,
                })
              }
              onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) =>
                e.key === "Enter" &&
                doBoxSearch(tag.id, queries[tag.id] ?? tag.name ?? "")
              }
            />
            <InputGroup.Append>
              <Button
                disabled={loading}
                onClick={() =>
                  doBoxSearch(tag.id, queries[tag.id] ?? tag.name ?? "")
                }
              >
                <FormattedMessage id="actions.search" />
              </Button>
            </InputGroup.Append>
          </InputGroup>
        );
      } else if (isTagged) {
        mainContent = (
          <div className="d-flex flex-column text-left">
            <h5>
              <FormattedMessage id="tag_tagger.tag_successfully_tagged" />
            </h5>
          </div>
        );
      }

      let subContent;
      if (stashID !== undefined) {
        const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
        const link = base ? (
          <ExternalLink
            className="small d-block"
            href={`${base}tags/${stashID.stash_id}`}
          >
            {stashID.stash_id}
          </ExternalLink>
        ) : (
          <div className="small">{stashID.stash_id}</div>
        );

        subContent = (
          <div key={tag.id}>
            <InputGroup className="StudioTagger-box-link">
              <InputGroup.Text>{link}</InputGroup.Text>
              <InputGroup.Append>
                <Button
                  onClick={() =>
                    doBoxUpdate(tag.id, stashID.stash_id, stashID.endpoint)
                  }
                  disabled={!!loadingUpdate}
                >
                  {loadingUpdate === stashID.stash_id ? (
                    <LoadingIndicator inline small message="" />
                  ) : (
                    <FormattedMessage id="actions.refresh" />
                  )}
                </Button>
              </InputGroup.Append>
            </InputGroup>
            {error[tag.id] && (
              <div className="text-danger mt-1">
                <strong>
                  <span className="mr-2">Error:</span>
                  {error[tag.id]?.message}
                </strong>
                <div>{error[tag.id]?.details}</div>
              </div>
            )}
          </div>
        );
      } else if (searchErrors[tag.id]) {
        subContent = (
          <div className="text-danger font-weight-bold">
            {searchErrors[tag.id]}
          </div>
        );
      } else if (searchResults[tag.id]?.length === 0) {
        subContent = (
          <div className="text-danger font-weight-bold">
            <FormattedMessage id="tag_tagger.no_results_found" />
          </div>
        );
      }

      let searchResult;
      if (searchResults[tag.id]?.length > 0 && !isTagged) {
        searchResult = (
          <StashSearchResult
            key={tag.id}
            stashboxTags={searchResults[tag.id]}
            tag={tag}
            endpoint={selectedEndpoint.endpoint}
            onTagTagged={handleTaggedTag}
            excludedTagFields={config.excludedTagFields ?? []}
          />
        );
      }

      return (
        <div key={tag.id} className={`${CLASSNAME}-studio`}>
          <div className={`${CLASSNAME}-details`}>
            <div></div>
            <div>
              <Card className="studio-card">
                <img loading="lazy" src={tag.image_path ?? ""} alt="" />
              </Card>
            </div>
            <div className={`${CLASSNAME}-details-text`}>
              <Link to={`/tags/${tag.id}`} className={`${CLASSNAME}-header`}>
                <h2>{tag.name}</h2>
              </Link>
              {mainContent}
              <div className="sub-content text-left">{subContent}</div>
              {searchResult}
            </div>
          </div>
        </div>
      );
    });

  return (
    <Card>
      {showBatchUpdate && (
        <TagBatchUpdateModal
          close={() => setShowBatchUpdate(false)}
          isIdle={isIdle}
          selectedEndpoint={selectedEndpoint}
          tags={tags}
          onBatchUpdate={handleBatchUpdate}
        />
      )}

      {showBatchAdd && (
        <TagBatchAddModal
          close={() => setShowBatchAdd(false)}
          isIdle={isIdle}
          onBatchAdd={handleBatchAdd}
        />
      )}
      <div className="ml-auto mb-3">
        <Button onClick={() => setShowBatchAdd(true)}>
          <FormattedMessage id="tag_tagger.batch_add_tags" />
        </Button>
        <Button className="ml-3" onClick={() => setShowBatchUpdate(true)}>
          <FormattedMessage id="tag_tagger.batch_update_tags" />
        </Button>
      </div>
      <div className={CLASSNAME}>{renderTags()}</div>
    </Card>
  );
};

interface ITaggerProps {
  tags: GQL.TagListDataFragment[];
}

export const TagTagger: React.FC<ITaggerProps> = ({ tags }) => {
  const jobsSubscribe = useJobsSubscribe();
  const intl = useIntl();
  const { configuration: stashConfig } = useConfigurationContext();
  const { config, setConfig } = useTaggerConfig();
  const [showConfig, setShowConfig] = useState(false);
  const [showManual, setShowManual] = useState(false);

  const [batchJobID, setBatchJobID] = useState<string | undefined | null>();
  const [batchJob, setBatchJob] = useState<JobFragment | undefined>();

  useEffect(() => {
    if (!jobsSubscribe.data) {
      return;
    }

    const event = jobsSubscribe.data.jobsSubscribe;
    if (event.job.id !== batchJobID) {
      return;
    }

    if (event.type !== GQL.JobStatusUpdateType.Remove) {
      setBatchJob(event.job);
    } else {
      setBatchJob(undefined);
      setBatchJobID(undefined);

      const ac = getClient();
      ac.cache.evict({ fieldName: "findTags" });
      ac.cache.gc();
    }
  }, [jobsSubscribe, batchJobID]);

  if (!config) return <LoadingIndicator />;

  const savedEndpointIndex =
    stashConfig?.general.stashBoxes.findIndex(
      (s) => s.endpoint === config.selectedEndpoint
    ) ?? -1;
  const selectedEndpointIndex =
    savedEndpointIndex === -1 && stashConfig?.general.stashBoxes.length
      ? 0
      : savedEndpointIndex;
  const selectedEndpoint =
    stashConfig?.general.stashBoxes[selectedEndpointIndex];

  async function batchAdd(tagInput: string) {
    if (tagInput && selectedEndpoint) {
      const inputs = tagInput
        .split(",")
        .map((n) => n.trim())
        .filter((n) => n.length > 0);

      const { names, stashIds } = separateNamesAndStashIds(inputs);

      if (names.length > 0 || stashIds.length > 0) {
        const ret = await mutateStashBoxBatchTagTag({
          names: names.length > 0 ? names : undefined,
          stash_ids: stashIds.length > 0 ? stashIds : undefined,
          endpoint: selectedEndpointIndex,
          refresh: false,
          createParent: false,
          exclude_fields: config?.excludedTagFields ?? [],
        });

        setBatchJobID(ret.data?.stashBoxBatchTagTag);
      }
    }
  }

  async function batchUpdate(ids: string[] | undefined, refresh: boolean) {
    if (selectedEndpoint) {
      const ret = await mutateStashBoxBatchTagTag({
        ids: ids,
        endpoint: selectedEndpointIndex,
        refresh,
        createParent: false,
        exclude_fields: config?.excludedTagFields ?? [],
      });

      setBatchJobID(ret.data?.stashBoxBatchTagTag);
    }
  }

  function renderStatus() {
    if (batchJob) {
      const progress =
        batchJob.progress !== undefined && batchJob.progress !== null
          ? batchJob.progress * 100
          : undefined;
      return (
        <Form.Group className="px-4">
          <h5>
            <FormattedMessage id="tag_tagger.status_tagging_tags" />
          </h5>
          {progress !== undefined && (
            <ProgressBar
              animated
              now={progress}
              label={`${progress.toFixed(0)}%`}
            />
          )}
        </Form.Group>
      );
    }

    if (batchJobID !== undefined) {
      return (
        <Form.Group className="px-4">
          <h5>
            <FormattedMessage id="tag_tagger.status_tagging_job_queued" />
          </h5>
        </Form.Group>
      );
    }
  }

  const showHideConfigId = showConfig
    ? "actions.hide_configuration"
    : "actions.show_configuration";

  return (
    <>
      <Manual
        show={showManual}
        onClose={() => setShowManual(false)}
        defaultActiveTab="Tagger.md"
      />
      {renderStatus()}
      <div className="tagger-container mx-md-auto">
        {selectedEndpointIndex !== -1 && selectedEndpoint ? (
          <>
            <div className="row mb-2 no-gutters">
              <Button onClick={() => setShowConfig(!showConfig)} variant="link">
                {intl.formatMessage({ id: showHideConfigId })}
              </Button>
              <Button
                className="ml-auto"
                onClick={() => setShowManual(true)}
                title={intl.formatMessage({ id: "help" })}
                variant="link"
              >
                <FormattedMessage id="help" />
              </Button>
            </div>

            <TaggerConfig
              config={config}
              setConfig={setConfig}
              show={showConfig}
              excludedFields={config.excludedTagFields ?? []}
              onFieldsChange={(fields) =>
                setConfig({ ...config, excludedTagFields: fields })
              }
              fields={TAG_FIELDS}
              entityName="tags"
            />
            <TagTaggerList
              tags={tags}
              selectedEndpoint={{
                endpoint: selectedEndpoint.endpoint,
                index: selectedEndpointIndex,
              }}
              isIdle={batchJobID === undefined}
              config={config}
              onBatchAdd={batchAdd}
              onBatchUpdate={batchUpdate}
            />
          </>
        ) : (
          <div className="my-4">
            <h3 className="text-center mt-4">
              <FormattedMessage id="tag_tagger.to_use_the_tag_tagger" />
            </h3>
            <h5 className="text-center">
              Please see{" "}
              <HashLink
                to="/settings?tab=metadata-providers#stash-boxes"
                scroll={(el) =>
                  el.scrollIntoView({ behavior: "smooth", block: "center" })
                }
              >
                Settings.
              </HashLink>
            </h5>
          </div>
        )}
      </div>
    </>
  );
};
