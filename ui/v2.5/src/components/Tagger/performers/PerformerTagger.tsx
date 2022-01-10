import React, { useEffect, useRef, useState } from "react";
import { Button, Card, Form, InputGroup, ProgressBar } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Link } from "react-router-dom";
import { HashLink } from "react-router-hash-link";
import { useLocalForage } from "src/hooks";

import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator, Modal } from "src/components/Shared";
import {
  stashBoxPerformerQuery,
  useJobsSubscribe,
  mutateStashBoxBatchPerformerTag,
} from "src/core/StashService";
import { Manual } from "src/components/Help/Manual";
import { ConfigurationContext } from "src/hooks/Config";

import StashSearchResult from "./StashSearchResult";
import PerformerConfig from "./Config";
import { LOCAL_FORAGE_KEY, ITaggerConfig, initialConfig } from "../constants";
import PerformerModal from "../PerformerModal";
import { useUpdatePerformer } from "../queries";

type JobFragment = Pick<
  GQL.Job,
  "id" | "status" | "subTasks" | "description" | "progress"
>;

const CLASSNAME = "PerformerTagger";

interface IPerformerTaggerListProps {
  performers: GQL.PerformerDataFragment[];
  selectedEndpoint: { endpoint: string; index: number };
  isIdle: boolean;
  config: ITaggerConfig;
  stashBoxes?: GQL.StashBox[];
  onBatchAdd: (performerInput: string) => void;
  onBatchUpdate: (ids: string[] | undefined, refresh: boolean) => void;
}

const PerformerTaggerList: React.FC<IPerformerTaggerListProps> = ({
  performers,
  selectedEndpoint,
  isIdle,
  config,
  stashBoxes,
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
  const [queryAll, setQueryAll] = useState(false);

  const [refresh, setRefresh] = useState(false);
  const { data: allPerformers } = GQL.useFindPerformersQuery({
    variables: {
      performer_filter: {
        stash_id: {
          value: "",
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
  const [showBatchAdd, setShowBatchAdd] = useState(false);
  const [showBatchUpdate, setShowBatchUpdate] = useState(false);
  const performerInput = useRef<HTMLTextAreaElement | null>(null);

  const [error, setError] = useState<
    Record<string, { message?: string; details?: string } | undefined>
  >({});
  const [loadingUpdate, setLoadingUpdate] = useState<string | undefined>();
  const [modalPerformer, setModalPerformer] = useState<
    GQL.ScrapedPerformerDataFragment | undefined
  >();

  const doBoxSearch = (performerID: string, searchVal: string) => {
    stashBoxPerformerQuery(searchVal, selectedEndpoint.index)
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
          [performerID]: "Network Error",
        });
      });

    setLoading(true);
  };

  const doBoxUpdate = (
    performerID: string,
    stashID: string,
    endpointIndex: number
  ) => {
    setLoadingUpdate(stashID);
    setError({
      ...error,
      [performerID]: undefined,
    });
    stashBoxPerformerQuery(stashID, endpointIndex)
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

  async function handleBatchAdd() {
    if (performerInput.current) {
      onBatchAdd(performerInput.current.value);
    }
    setShowBatchAdd(false);
  }

  const handleBatchUpdate = () => {
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

  const handlePerformerUpdate = async (input: GQL.PerformerCreateInput) => {
    setModalPerformer(undefined);
    const performerID = modalPerformer?.stored_id;
    if (performerID) {
      const updateData: GQL.PerformerUpdateInput = {
        ...input,
        id: performerID,
      };

      const res = await updatePerformer(updateData);
      if (!res.data?.performerUpdate)
        setError({
          ...error,
          [performerID]: {
            message: `Failed to save performer "${modalPerformer?.name}"`,
            details:
              res?.errors?.[0].message ===
              "UNIQUE constraint failed: performers.checksum"
                ? "Name already exists"
                : res?.errors?.[0].message,
          },
        });
    }
  };

  const renderPerformers = () =>
    performers.map((performer) => {
      const isTagged = taggedPerformers[performer.id];
      const hasStashIDs = performer.stash_ids.length > 0;

      let mainContent;
      if (!isTagged && hasStashIDs) {
        mainContent = (
          <div className="text-left">
            <h5 className="text-bold">Performer already tagged</h5>
          </div>
        );
      } else if (!isTagged && !hasStashIDs) {
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
            <h5>Performer successfully tagged:</h5>
            <h6>
              <Link className="bold" to={`/performers/${performer.id}`}>
                {taggedPerformers[performer.id].name}
              </Link>
            </h6>
          </div>
        );
      }

      let subContent;
      if (performer.stash_ids.length > 0) {
        const stashLinks = performer.stash_ids.map((stashID) => {
          const base = stashID.endpoint.match(/https?:\/\/.*?\//)?.[0];
          const link = base ? (
            <a
              className="small d-block"
              href={`${base}performers/${stashID.stash_id}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              {stashID.stash_id}
            </a>
          ) : (
            <div className="small">{stashID.stash_id}</div>
          );

          const endpoint =
            stashBoxes?.findIndex((box) => box.endpoint === stashID.endpoint) ??
            -1;

          return (
            <div key={performer.id}>
              <InputGroup className="PerformerTagger-box-link">
                <InputGroup.Text>{link}</InputGroup.Text>
                <InputGroup.Append>
                  {endpoint !== -1 && (
                    <Button
                      onClick={() =>
                        doBoxUpdate(performer.id, stashID.stash_id, endpoint)
                      }
                      disabled={!!loadingUpdate}
                    >
                      {loadingUpdate === stashID.stash_id ? (
                        <LoadingIndicator inline small message="" />
                      ) : (
                        "Refresh"
                      )}
                    </Button>
                  )}
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
        });
        subContent = <>{stashLinks}</>;
      } else if (searchErrors[performer.id]) {
        subContent = (
          <div className="text-danger font-weight-bold">
            {searchErrors[performer.id]}
          </div>
        );
      } else if (searchResults[performer.id]?.length === 0) {
        subContent = (
          <div className="text-danger font-weight-bold">No results found.</div>
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
              onSave={handlePerformerUpdate}
              excludedPerformerFields={config.excludedPerformerFields}
              icon="tags"
              header="Update Performer"
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
              <h2>{performer.name}</h2>
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
      <Modal
        show={showBatchUpdate}
        icon="tags"
        header="Update Performers"
        accept={{ text: "Update Performers", onClick: handleBatchUpdate }}
        cancel={{
          text: intl.formatMessage({ id: "actions.cancel" }),
          variant: "danger",
          onClick: () => setShowBatchUpdate(false),
        }}
        disabled={!isIdle}
      >
        <Form.Group>
          <Form.Label>
            <h6>Performer selection</h6>
          </Form.Label>
          <Form.Check
            id="query-page"
            type="radio"
            name="performer-query"
            label="Current page"
            defaultChecked
            onChange={() => setQueryAll(false)}
          />
          <Form.Check
            id="query-all"
            type="radio"
            name="performer-query"
            label="All performers in the database"
            defaultChecked={false}
            onChange={() => setQueryAll(true)}
          />
        </Form.Group>
        <Form.Group>
          <Form.Label>
            <h6>Tag Status</h6>
          </Form.Label>
          <Form.Check
            id="untagged-performers"
            type="radio"
            name="performer-refresh"
            label="Untagged performers"
            defaultChecked
            onChange={() => setRefresh(false)}
          />
          <Form.Text>
            Updating untagged performers will try to match any performers that
            lack a stashid and update the metadata.
          </Form.Text>
          <Form.Check
            id="tagged-performers"
            type="radio"
            name="performer-refresh"
            label="Refresh tagged performers"
            defaultChecked={false}
            onChange={() => setRefresh(true)}
          />
          <Form.Text>
            Refreshing will update the data of any tagged performers from the
            stash-box instance.
          </Form.Text>
        </Form.Group>
        <b>{`${
          queryAll
            ? allPerformers?.findPerformers.count
            : performers.filter((p) =>
                refresh ? p.stash_ids.length > 0 : p.stash_ids.length === 0
              ).length
        } performers will be processed`}</b>
      </Modal>
      <Modal
        show={showBatchAdd}
        icon="star"
        header="Add New Performers"
        accept={{ text: "Add Performers", onClick: handleBatchAdd }}
        cancel={{
          text: intl.formatMessage({ id: "actions.cancel" }),
          variant: "danger",
          onClick: () => setShowBatchAdd(false),
        }}
        disabled={!isIdle}
      >
        <Form.Control
          className="text-input"
          as="textarea"
          ref={performerInput}
          placeholder="Performer names separated by comma"
          rows={6}
        />
        <Form.Text>
          Any names entered will be queried from the remote Stash-Box instance
          and added if found. Only exact matches will be considered a match.
        </Form.Text>
      </Modal>
      <div className="ml-auto mb-3">
        <Button onClick={() => setShowBatchAdd(true)}>
          Batch Add Performers
        </Button>
        <Button className="ml-3" onClick={() => setShowBatchUpdate(true)}>
          Batch Update Performers
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
          performer_names: names,
          endpoint: selectedEndpointIndex,
          refresh: false,
        });

        setBatchJobID(ret.data?.stashBoxBatchPerformerTag);
      }
    }
  }

  async function batchUpdate(ids: string[] | undefined, refresh: boolean) {
    if (config && selectedEndpoint) {
      const ret = await mutateStashBoxBatchPerformerTag({
        performer_ids: ids,
        endpoint: selectedEndpointIndex,
        refresh,
        exclude_fields: config.excludedPerformerFields ?? [],
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
          <h5>Status: Tagging performers</h5>
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
          <h5>Status: Tagging job queued</h5>
        </Form.Group>
      );
    }
  }

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
                {showConfig ? "Hide" : "Show"} Configuration
              </Button>
              <Button
                className="ml-auto"
                onClick={() => setShowManual(true)}
                title="Help"
                variant="link"
              >
                Help
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
              stashBoxes={stashConfig?.general.stashBoxes}
              onBatchAdd={batchAdd}
              onBatchUpdate={batchUpdate}
            />
          </>
        ) : (
          <div className="my-4">
            <h3 className="text-center mt-4">
              To use the performer tagger a stash-box instance needs to be
              configured.
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
