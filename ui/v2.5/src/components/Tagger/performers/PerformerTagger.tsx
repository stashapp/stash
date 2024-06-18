import React, { useEffect, useMemo, useRef, useState } from "react";
import { Button, Card, Form, InputGroup, ProgressBar } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Link } from "react-router-dom";
import { HashLink } from "react-router-hash-link";
import { useLocalForage } from "src/hooks/LocalForage";

import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { ModalComponent } from "src/components/Shared/Modal";
import {
  stashBoxPerformerQuery,
  useJobsSubscribe,
  mutateStashBoxBatchPerformerTag,
  getClient,
  evictQueries,
  performerMutationImpactedQueries,
} from "src/core/StashService";
import { Manual } from "src/components/Help/Manual";
import { ConfigurationContext } from "src/hooks/Config";

import StashSearchResult from "./StashSearchResult";
import PerformerConfig from "./Config";
import { LOCAL_FORAGE_KEY, ITaggerConfig, initialConfig } from "../constants";
import PerformerModal from "../PerformerModal";
import { useUpdatePerformer } from "../queries";
import { faStar, faTags } from "@fortawesome/free-solid-svg-icons";
import { mergeStashIDs } from "src/utils/stashbox";
import { ExternalLink } from "src/components/Shared/ExternalLink";

type JobFragment = Pick<
  GQL.Job,
  "id" | "status" | "subTasks" | "description" | "progress"
>;

const CLASSNAME = "PerformerTagger";

interface IPerformerBatchUpdateModal {
  performers: GQL.PerformerDataFragment[];
  isIdle: boolean;
  selectedEndpoint: { endpoint: string; index: number };
  onBatchUpdate: (queryAll: boolean, refresh: boolean) => void;
  close: () => void;
}

const PerformerBatchUpdateModal: React.FC<IPerformerBatchUpdateModal> = ({
  performers,
  isIdle,
  selectedEndpoint,
  onBatchUpdate,
  close,
}) => {
  const intl = useIntl();

  const [queryAll, setQueryAll] = useState(false);

  const [refresh, setRefresh] = useState(false);
  const { data: allPerformers } = GQL.useFindPerformersQuery({
    variables: {
      performer_filter: {
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

  const performerCount = useMemo(() => {
    // get all stash ids for the selected endpoint
    const filteredStashIDs = performers.map((p) =>
      p.stash_ids.filter((s) => s.endpoint === selectedEndpoint.endpoint)
    );

    return queryAll
      ? allPerformers?.findPerformers.count
      : filteredStashIDs.filter((s) =>
          // if refresh, then we filter out the performers without a stash id
          // otherwise, we want untagged performers, filtering out those with a stash id
          refresh ? s.length > 0 : s.length === 0
        ).length;
  }, [queryAll, refresh, performers, allPerformers, selectedEndpoint.endpoint]);

  return (
    <ModalComponent
      show
      icon={faTags}
      header={intl.formatMessage({
        id: "performer_tagger.update_performers",
      })}
      accept={{
        text: intl.formatMessage({
          id: "performer_tagger.update_performers",
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
            <FormattedMessage id="performer_tagger.performer_selection" />
          </h6>
        </Form.Label>
        <Form.Check
          id="query-page"
          type="radio"
          name="performer-query"
          label={<FormattedMessage id="performer_tagger.current_page" />}
          checked={!queryAll}
          onChange={() => setQueryAll(false)}
        />
        <Form.Check
          id="query-all"
          type="radio"
          name="performer-query"
          label={intl.formatMessage({
            id: "performer_tagger.query_all_performers_in_the_database",
          })}
          checked={queryAll}
          onChange={() => setQueryAll(true)}
        />
      </Form.Group>
      <Form.Group>
        <Form.Label>
          <h6>
            <FormattedMessage id="performer_tagger.tag_status" />
          </h6>
        </Form.Label>
        <Form.Check
          id="untagged-performers"
          type="radio"
          name="performer-refresh"
          label={intl.formatMessage({
            id: "performer_tagger.untagged_performers",
          })}
          checked={!refresh}
          onChange={() => setRefresh(false)}
        />
        <Form.Text>
          <FormattedMessage id="performer_tagger.updating_untagged_performers_description" />
        </Form.Text>
        <Form.Check
          id="tagged-performers"
          type="radio"
          name="performer-refresh"
          label={intl.formatMessage({
            id: "performer_tagger.refresh_tagged_performers",
          })}
          checked={refresh}
          onChange={() => setRefresh(true)}
        />
        <Form.Text>
          <FormattedMessage id="performer_tagger.refreshing_will_update_the_data" />
        </Form.Text>
      </Form.Group>
      <b>
        <FormattedMessage
          id="performer_tagger.number_of_performers_will_be_processed"
          values={{
            performer_count: performerCount,
          }}
        />
      </b>
    </ModalComponent>
  );
};

interface IPerformerBatchAddModal {
  isIdle: boolean;
  onBatchAdd: (input: string) => void;
  close: () => void;
}

const PerformerBatchAddModal: React.FC<IPerformerBatchAddModal> = ({
  isIdle,
  onBatchAdd,
  close,
}) => {
  const intl = useIntl();

  const performerInput = useRef<HTMLTextAreaElement | null>(null);

  return (
    <ModalComponent
      show
      icon={faStar}
      header={intl.formatMessage({
        id: "performer_tagger.add_new_performers",
      })}
      accept={{
        text: intl.formatMessage({
          id: "performer_tagger.add_new_performers",
        }),
        onClick: () => {
          if (performerInput.current) {
            onBatchAdd(performerInput.current.value);
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
        ref={performerInput}
        placeholder={intl.formatMessage({
          id: "performer_tagger.performer_names_separated_by_comma",
        })}
        rows={6}
      />
      <Form.Text>
        <FormattedMessage id="performer_tagger.any_names_entered_will_be_queried" />
      </Form.Text>
    </ModalComponent>
  );
};

interface IPerformerTaggerListProps {
  performers: GQL.PerformerDataFragment[];
  selectedEndpoint: { endpoint: string; index: number };
  isIdle: boolean;
  config: ITaggerConfig;
  onBatchAdd: (performerInput: string) => void;
  onBatchUpdate: (ids: string[] | undefined, refresh: boolean) => void;
}

const PerformerTaggerList: React.FC<IPerformerTaggerListProps> = ({
  performers,
  selectedEndpoint,
  isIdle,
  config,
  onBatchAdd,
  onBatchUpdate,
}) => {
  const intl = useIntl();
  const [loading, setLoading] = useState(false);
  const [searchResults, setSearchResults] = useState<
    Record<string, GQL.ScrapedPerformerDataFragment[]>
  >({});
  const [searchErrors, setSearchErrors] = useState<
    Record<string, string | undefined>
  >({});
  const [taggedPerformers, setTaggedPerformers] = useState<
    Record<string, Partial<GQL.SlimPerformerDataFragment>>
  >({});
  const [queries, setQueries] = useState<Record<string, string>>({});

  const [showBatchAdd, setShowBatchAdd] = useState(false);
  const [showBatchUpdate, setShowBatchUpdate] = useState(false);

  const [error, setError] = useState<
    Record<string, { message?: string; details?: string } | undefined>
  >({});
  const [loadingUpdate, setLoadingUpdate] = useState<string | undefined>();
  const [modalPerformer, setModalPerformer] = useState<
    GQL.ScrapedPerformerDataFragment | undefined
  >();

  const doBoxSearch = (performerID: string, searchVal: string) => {
    stashBoxPerformerQuery(searchVal, selectedEndpoint.endpoint)
      .then((queryData) => {
        const s = queryData.data?.scrapeSinglePerformer ?? [];
        setSearchResults({
          ...searchResults,
          [performerID]: s,
        });
        setSearchErrors({
          ...searchErrors,
          [performerID]: undefined,
        });
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
        // Destructure to remove existing result
        const { [performerID]: unassign, ...results } = searchResults;
        setSearchResults(results);
        setSearchErrors({
          ...searchErrors,
          [performerID]: intl.formatMessage({
            id: "performer_tagger.network_error",
          }),
        });
      });

    setLoading(true);
  };

  const doBoxUpdate = (
    performerID: string,
    stashID: string,
    endpoint: string
  ) => {
    setLoadingUpdate(stashID);
    setError({
      ...error,
      [performerID]: undefined,
    });
    stashBoxPerformerQuery(stashID, endpoint)
      .then((queryData) => {
        const data = queryData.data?.scrapeSinglePerformer ?? [];
        if (data.length > 0) {
          setModalPerformer({
            ...data[0],
            stored_id: performerID,
          });
        }
      })
      .finally(() => setLoadingUpdate(undefined));
  };

  async function handleBatchAdd(input: string) {
    onBatchAdd(input);
    setShowBatchAdd(false);
  }

  const handleBatchUpdate = (queryAll: boolean, refresh: boolean) => {
    onBatchUpdate(!queryAll ? performers.map((p) => p.id) : undefined, refresh);
    setShowBatchUpdate(false);
  };

  const handleTaggedPerformer = (
    performer: Pick<GQL.SlimPerformerDataFragment, "id"> &
      Partial<Omit<GQL.SlimPerformerDataFragment, "id">>
  ) => {
    setTaggedPerformers({
      ...taggedPerformers,
      [performer.id]: performer,
    });
  };

  const updatePerformer = useUpdatePerformer();

  function handleSaveError(performerID: string, name: string, message: string) {
    setError({
      ...error,
      [performerID]: {
        message: intl.formatMessage(
          { id: "performer_tagger.failed_to_save_performer" },
          { performer: modalPerformer?.name }
        ),
        details:
          message === "UNIQUE constraint failed: performers.name"
            ? intl.formatMessage({
                id: "performer_tagger.name_already_exists",
              })
            : message,
      },
    });
  }

  const handlePerformerUpdate = async (
    existing: GQL.PerformerDataFragment,
    input: GQL.PerformerCreateInput
  ) => {
    setModalPerformer(undefined);
    const performerID = modalPerformer?.stored_id;
    if (performerID) {
      // handle stash ids - we want to add, not set them
      if (input.stash_ids?.length) {
        input.stash_ids = mergeStashIDs(existing.stash_ids, input.stash_ids);
      }

      const updateData: GQL.PerformerUpdateInput = {
        ...input,
        id: performerID,
      };

      const res = await updatePerformer(updateData);
      if (!res.data?.performerUpdate)
        handleSaveError(
          performerID,
          modalPerformer?.name ?? "",
          res?.errors?.[0]?.message ?? ""
        );
    }
  };

  const renderPerformers = () =>
    performers.map((performer) => {
      const isTagged = taggedPerformers[performer.id];

      const stashID = performer.stash_ids.find((s) => {
        return s.endpoint === selectedEndpoint.endpoint;
      });

      let mainContent;
      if (!isTagged && stashID !== undefined) {
        mainContent = (
          <div className="text-left">
            <h5 className="text-bold">
              <FormattedMessage id="performer_tagger.performer_already_tagged" />
            </h5>
          </div>
        );
      } else if (!isTagged && !stashID) {
        mainContent = (
          <InputGroup>
            <Form.Control
              className="text-input"
              defaultValue={performer.name ?? ""}
              onChange={(e) =>
                setQueries({
                  ...queries,
                  [performer.id]: e.currentTarget.value,
                })
              }
              onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) =>
                e.key === "Enter" &&
                doBoxSearch(
                  performer.id,
                  queries[performer.id] ?? performer.name ?? ""
                )
              }
            />
            <InputGroup.Append>
              <Button
                disabled={loading}
                onClick={() =>
                  doBoxSearch(
                    performer.id,
                    queries[performer.id] ?? performer.name ?? ""
                  )
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
              <FormattedMessage id="performer_tagger.performer_successfully_tagged" />
            </h5>
            <h6>
              <Link className="bold" to={`/performers/${performer.id}`}>
                {taggedPerformers[performer.id].name}
              </Link>
            </h6>
          </div>
        );
      }

      let subContent;
      if (stashID !== undefined) {
        const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
        const link = base ? (
          <ExternalLink
            className="small d-block"
            href={`${base}performers/${stashID.stash_id}`}
          >
            {stashID.stash_id}
          </ExternalLink>
        ) : (
          <div className="small">{stashID.stash_id}</div>
        );

        subContent = (
          <div key={performer.id}>
            <InputGroup className="PerformerTagger-box-link">
              <InputGroup.Text>{link}</InputGroup.Text>
              <InputGroup.Append>
                <Button
                  onClick={() =>
                    doBoxUpdate(
                      performer.id,
                      stashID.stash_id,
                      stashID.endpoint
                    )
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
            {error[performer.id] && (
              <div className="text-danger mt-1">
                <strong>
                  <span className="mr-2">Error:</span>
                  {error[performer.id]?.message}
                </strong>
                <div>{error[performer.id]?.details}</div>
              </div>
            )}
          </div>
        );
      } else if (searchErrors[performer.id]) {
        subContent = (
          <div className="text-danger font-weight-bold">
            {searchErrors[performer.id]}
          </div>
        );
      } else if (searchResults[performer.id]?.length === 0) {
        subContent = (
          <div className="text-danger font-weight-bold">
            <FormattedMessage id="performer_tagger.no_results_found" />
          </div>
        );
      }

      let searchResult;
      if (searchResults[performer.id]?.length > 0 && !isTagged) {
        searchResult = (
          <StashSearchResult
            key={performer.id}
            stashboxPerformers={searchResults[performer.id]}
            performer={performer}
            endpoint={selectedEndpoint.endpoint}
            onPerformerTagged={handleTaggedPerformer}
            excludedPerformerFields={config.excludedPerformerFields ?? []}
          />
        );
      }

      return (
        <div key={performer.id} className={`${CLASSNAME}-performer`}>
          {modalPerformer && (
            <PerformerModal
              closeModal={() => setModalPerformer(undefined)}
              modalVisible={modalPerformer.stored_id === performer.id}
              performer={modalPerformer}
              onSave={(input) => {
                handlePerformerUpdate(performer, input);
              }}
              excludedPerformerFields={config.excludedPerformerFields}
              icon={faTags}
              header={intl.formatMessage({
                id: "performer_tagger.update_performer",
              })}
              endpoint={selectedEndpoint.endpoint}
            />
          )}
          <Card className="performer-card p-0 m-0">
            <img src={performer.image_path ?? ""} alt="" />
          </Card>
          <div className={`${CLASSNAME}-details`}>
            <Link
              to={`/performers/${performer.id}`}
              className={`${CLASSNAME}-header`}
            >
              <h2>
                {performer.name}
                {performer.disambiguation && (
                  <span className="performer-disambiguation">
                    {` (${performer.disambiguation})`}
                  </span>
                )}
              </h2>
            </Link>
            {mainContent}
            <div className="sub-content text-left">{subContent}</div>
            {searchResult}
          </div>
        </div>
      );
    });

  return (
    <Card>
      {showBatchUpdate && (
        <PerformerBatchUpdateModal
          close={() => setShowBatchUpdate(false)}
          isIdle={isIdle}
          selectedEndpoint={selectedEndpoint}
          performers={performers}
          onBatchUpdate={handleBatchUpdate}
        />
      )}

      {showBatchAdd && (
        <PerformerBatchAddModal
          close={() => setShowBatchAdd(false)}
          isIdle={isIdle}
          onBatchAdd={handleBatchAdd}
        />
      )}

      <div className="ml-auto mb-3">
        <Button onClick={() => setShowBatchAdd(true)}>
          <FormattedMessage id="performer_tagger.batch_add_performers" />
        </Button>
        <Button className="ml-3" onClick={() => setShowBatchUpdate(true)}>
          <FormattedMessage id="performer_tagger.batch_update_performers" />
        </Button>
      </div>
      <div className={CLASSNAME}>{renderPerformers()}</div>
    </Card>
  );
};

interface ITaggerProps {
  performers: GQL.PerformerDataFragment[];
}

export const PerformerTagger: React.FC<ITaggerProps> = ({ performers }) => {
  const jobsSubscribe = useJobsSubscribe();
  const intl = useIntl();
  const { configuration: stashConfig } = React.useContext(ConfigurationContext);
  const [{ data: config }, setConfig] = useLocalForage<ITaggerConfig>(
    LOCAL_FORAGE_KEY,
    initialConfig
  );
  const [showConfig, setShowConfig] = useState(false);
  const [showManual, setShowManual] = useState(false);

  const [batchJobID, setBatchJobID] = useState<string | undefined | null>();
  const [batchJob, setBatchJob] = useState<JobFragment | undefined>();

  // monitor batch operation
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

      // Once the performer batch is complete, refresh all local performer data
      const ac = getClient();
      evictQueries(ac.cache, performerMutationImpactedQueries);
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

  async function batchAdd(performerInput: string) {
    if (performerInput && selectedEndpoint) {
      const names = performerInput
        .split(",")
        .map((n) => n.trim())
        .filter((n) => n.length > 0);

      if (names.length > 0) {
        const ret = await mutateStashBoxBatchPerformerTag({
          names: names,
          endpoint: selectedEndpointIndex,
          refresh: false,
          createParent: false,
        });

        setBatchJobID(ret.data?.stashBoxBatchPerformerTag);
      }
    }
  }

  async function batchUpdate(ids: string[] | undefined, refresh: boolean) {
    if (config && selectedEndpoint) {
      const ret = await mutateStashBoxBatchPerformerTag({
        ids: ids,
        endpoint: selectedEndpointIndex,
        refresh,
        exclude_fields: config.excludedPerformerFields ?? [],
        createParent: false,
      });

      setBatchJobID(ret.data?.stashBoxBatchPerformerTag);
    }
  }

  // const progress =
  //   jobStatus.data?.metadataUpdate.status ===
  //     "Stash-Box Performer Batch Operation" &&
  //   jobStatus.data.metadataUpdate.progress >= 0
  //     ? jobStatus.data.metadataUpdate.progress * 100
  //     : null;

  function renderStatus() {
    if (batchJob) {
      const progress =
        batchJob.progress !== undefined && batchJob.progress !== null
          ? batchJob.progress * 100
          : undefined;
      return (
        <Form.Group className="px-4">
          <h5>
            <FormattedMessage id="performer_tagger.status_tagging_performers" />
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
            <FormattedMessage id="performer_tagger.status_tagging_job_queued" />
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

            <PerformerConfig
              config={config}
              setConfig={setConfig}
              show={showConfig}
            />
            <PerformerTaggerList
              performers={performers}
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
              <FormattedMessage id="performer_tagger.to_use_the_performer_tagger" />
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
