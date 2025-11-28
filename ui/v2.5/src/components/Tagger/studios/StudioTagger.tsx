import React, { useEffect, useMemo, useRef, useState } from "react";
import { Button, Card, Form, InputGroup, ProgressBar } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Link } from "react-router-dom";
import { HashLink } from "react-router-hash-link";

import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { ModalComponent } from "src/components/Shared/Modal";
import {
  stashBoxStudioQuery,
  useJobsSubscribe,
  mutateStashBoxBatchStudioTag,
  getClient,
  studioMutationImpactedQueries,
  useStudioCreate,
  evictQueries,
} from "src/core/StashService";
import { Manual } from "src/components/Help/Manual";
import { useConfigurationContext } from "src/hooks/Config";

import StashSearchResult from "./StashSearchResult";
import StudioConfig from "./Config";
import { ITaggerConfig } from "../constants";
import StudioModal from "../scenes/StudioModal";
import { useUpdateStudio } from "../queries";
import { apolloError } from "src/utils";
import { faStar, faTags } from "@fortawesome/free-solid-svg-icons";
import { ExternalLink } from "src/components/Shared/ExternalLink";
import { mergeStudioStashIDs } from "../utils";
import { separateNamesAndStashIds } from "src/utils/stashIds";
import { useTaggerConfig } from "../config";

type JobFragment = Pick<
  GQL.Job,
  "id" | "status" | "subTasks" | "description" | "progress"
>;

const CLASSNAME = "StudioTagger";

interface IStudioBatchUpdateModal {
  studios: GQL.StudioDataFragment[];
  isIdle: boolean;
  selectedEndpoint: { endpoint: string; index: number };
  onBatchUpdate: (queryAll: boolean, refresh: boolean) => void;
  batchAddParents: boolean;
  setBatchAddParents: (addParents: boolean) => void;
  close: () => void;
}

const StudioBatchUpdateModal: React.FC<IStudioBatchUpdateModal> = ({
  studios,
  isIdle,
  selectedEndpoint,
  onBatchUpdate,
  batchAddParents,
  setBatchAddParents,
  close,
}) => {
  const intl = useIntl();

  const [queryAll, setQueryAll] = useState(false);

  const [refresh, setRefresh] = useState(false);
  const { data: allStudios } = GQL.useFindStudiosQuery({
    variables: {
      studio_filter: {
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

  const studioCount = useMemo(() => {
    // get all stash ids for the selected endpoint
    const filteredStashIDs = studios.map((p) =>
      p.stash_ids.filter((s) => s.endpoint === selectedEndpoint.endpoint)
    );

    return queryAll
      ? allStudios?.findStudios.count
      : filteredStashIDs.filter((s) =>
          // if refresh, then we filter out the studios without a stash id
          // otherwise, we want untagged studios, filtering out those with a stash id
          refresh ? s.length > 0 : s.length === 0
        ).length;
  }, [queryAll, refresh, studios, allStudios, selectedEndpoint.endpoint]);

  return (
    <ModalComponent
      show
      icon={faTags}
      header={intl.formatMessage({
        id: "studio_tagger.update_studios",
      })}
      accept={{
        text: intl.formatMessage({
          id: "studio_tagger.update_studios",
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
            <FormattedMessage id="studio_tagger.studio_selection" />
          </h6>
        </Form.Label>
        <Form.Check
          id="query-page"
          type="radio"
          name="studio-query"
          label={<FormattedMessage id="studio_tagger.current_page" />}
          checked={!queryAll}
          onChange={() => setQueryAll(false)}
        />
        <Form.Check
          id="query-all"
          type="radio"
          name="studio-query"
          label={intl.formatMessage({
            id: "studio_tagger.query_all_studios_in_the_database",
          })}
          checked={queryAll}
          onChange={() => setQueryAll(true)}
        />
      </Form.Group>
      <Form.Group>
        <Form.Label>
          <h6>
            <FormattedMessage id="studio_tagger.tag_status" />
          </h6>
        </Form.Label>
        <Form.Check
          id="untagged-studios"
          type="radio"
          name="studio-refresh"
          label={intl.formatMessage({
            id: "studio_tagger.untagged_studios",
          })}
          checked={!refresh}
          onChange={() => setRefresh(false)}
        />
        <Form.Text>
          <FormattedMessage id="studio_tagger.updating_untagged_studios_description" />
        </Form.Text>
        <Form.Check
          id="tagged-studios"
          type="radio"
          name="studio-refresh"
          label={intl.formatMessage({
            id: "studio_tagger.refresh_tagged_studios",
          })}
          checked={refresh}
          onChange={() => setRefresh(true)}
        />
        <Form.Text>
          <FormattedMessage id="studio_tagger.refreshing_will_update_the_data" />
        </Form.Text>
        <div className="mt-4">
          <Form.Check
            id="add-parent"
            checked={batchAddParents}
            label={intl.formatMessage({
              id: "studio_tagger.create_or_tag_parent_studios",
            })}
            onChange={() => setBatchAddParents(!batchAddParents)}
          />
        </div>
      </Form.Group>
      <b>
        <FormattedMessage
          id="studio_tagger.number_of_studios_will_be_processed"
          values={{
            studio_count: studioCount,
          }}
        />
      </b>
    </ModalComponent>
  );
};

interface IStudioBatchAddModal {
  isIdle: boolean;
  onBatchAdd: (input: string) => void;
  batchAddParents: boolean;
  setBatchAddParents: (addParents: boolean) => void;
  close: () => void;
}

const StudioBatchAddModal: React.FC<IStudioBatchAddModal> = ({
  isIdle,
  onBatchAdd,
  batchAddParents,
  setBatchAddParents,
  close,
}) => {
  const intl = useIntl();

  const studioInput = useRef<HTMLTextAreaElement | null>(null);

  return (
    <ModalComponent
      show
      icon={faStar}
      header={intl.formatMessage({
        id: "studio_tagger.add_new_studios",
      })}
      accept={{
        text: intl.formatMessage({
          id: "studio_tagger.add_new_studios",
        }),
        onClick: () => {
          if (studioInput.current) {
            onBatchAdd(studioInput.current.value);
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
        ref={studioInput}
        placeholder={intl.formatMessage({
          id: "studio_tagger.studio_names_or_stashids_separated_by_comma",
        })}
        rows={6}
      />
      <Form.Text>
        <FormattedMessage id="studio_tagger.any_names_entered_will_be_queried" />
      </Form.Text>
      <div className="mt-2">
        <Form.Check
          id="add-parent"
          checked={batchAddParents}
          label={intl.formatMessage({
            id: "studio_tagger.create_or_tag_parent_studios",
          })}
          onChange={() => setBatchAddParents(!batchAddParents)}
        />
      </div>
    </ModalComponent>
  );
};

interface IStudioTaggerListProps {
  studios: GQL.StudioDataFragment[];
  selectedEndpoint: { endpoint: string; index: number };
  isIdle: boolean;
  config: ITaggerConfig;
  onBatchAdd: (studioInput: string, createParent: boolean) => void;
  onBatchUpdate: (
    ids: string[] | undefined,
    refresh: boolean,
    createParent: boolean
  ) => void;
}

const StudioTaggerList: React.FC<IStudioTaggerListProps> = ({
  studios,
  selectedEndpoint,
  isIdle,
  config,
  onBatchAdd,
  onBatchUpdate,
}) => {
  const intl = useIntl();

  const [loading, setLoading] = useState(false);
  const [searchResults, setSearchResults] = useState<
    Record<string, GQL.ScrapedStudioDataFragment[]>
  >({});
  const [searchErrors, setSearchErrors] = useState<
    Record<string, string | undefined>
  >({});
  const [taggedStudios, setTaggedStudios] = useState<
    Record<string, Partial<GQL.SlimStudioDataFragment>>
  >({});
  const [queries, setQueries] = useState<Record<string, string>>({});

  const [showBatchAdd, setShowBatchAdd] = useState(false);
  const [showBatchUpdate, setShowBatchUpdate] = useState(false);
  const [batchAddParents, setBatchAddParents] = useState(
    config.createParentStudios || false
  );

  const [error, setError] = useState<
    Record<string, { message?: string; details?: string } | undefined>
  >({});
  const [loadingUpdate, setLoadingUpdate] = useState<string | undefined>();
  const [modalStudio, setModalStudio] = useState<
    GQL.ScrapedStudioDataFragment | undefined
  >();

  const doBoxSearch = (studioID: string, searchVal: string) => {
    stashBoxStudioQuery(searchVal, selectedEndpoint.endpoint)
      .then((queryData) => {
        const s = queryData.data?.scrapeSingleStudio ?? [];
        setSearchResults({
          ...searchResults,
          [studioID]: s,
        });
        setSearchErrors({
          ...searchErrors,
          [studioID]: undefined,
        });
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
        // Destructure to remove existing result
        const { [studioID]: unassign, ...results } = searchResults;
        setSearchResults(results);
        setSearchErrors({
          ...searchErrors,
          [studioID]: intl.formatMessage({
            id: "studio_tagger.network_error",
          }),
        });
      });

    setLoading(true);
  };

  const doBoxUpdate = (studioID: string, stashID: string, endpoint: string) => {
    setLoadingUpdate(stashID);
    setError({
      ...error,
      [studioID]: undefined,
    });
    stashBoxStudioQuery(stashID, endpoint)
      .then((queryData) => {
        const data = queryData.data?.scrapeSingleStudio ?? [];
        if (data.length > 0) {
          setModalStudio({
            ...data[0],
            stored_id: studioID,
          });
        }
      })
      .finally(() => setLoadingUpdate(undefined));
  };

  async function handleBatchAdd(input: string) {
    onBatchAdd(input, batchAddParents);
    setShowBatchAdd(false);
  }

  const handleBatchUpdate = (queryAll: boolean, refresh: boolean) => {
    onBatchUpdate(
      !queryAll ? studios.map((p) => p.id) : undefined,
      refresh,
      batchAddParents
    );
    setShowBatchUpdate(false);
  };

  const handleTaggedStudio = (
    studio: Pick<GQL.SlimStudioDataFragment, "id"> &
      Partial<Omit<GQL.SlimStudioDataFragment, "id">>
  ) => {
    setTaggedStudios({
      ...taggedStudios,
      [studio.id]: studio,
    });
  };

  const [createStudio] = useStudioCreate();
  const updateStudio = useUpdateStudio();

  function handleSaveError(studioID: string, name: string, message: string) {
    setError({
      ...error,
      [studioID]: {
        message: intl.formatMessage(
          { id: "studio_tagger.failed_to_save_studio" },
          { studio: modalStudio?.name }
        ),
        details:
          message === "UNIQUE constraint failed: studios.name"
            ? intl.formatMessage({
                id: "studio_tagger.name_already_exists",
              })
            : message,
      },
    });
  }

  const handleStudioUpdate = async (
    input: GQL.StudioCreateInput,
    parentInput?: GQL.StudioCreateInput
  ) => {
    setModalStudio(undefined);
    const studioID = modalStudio?.stored_id;
    if (studioID) {
      if (parentInput) {
        try {
          // if parent id is set, then update the existing studio
          if (input.parent_id) {
            const parentUpdateData: GQL.StudioUpdateInput = {
              ...parentInput,
              id: input.parent_id,
            };
            parentUpdateData.stash_ids = await mergeStudioStashIDs(
              input.parent_id,
              parentInput.stash_ids ?? []
            );
            await updateStudio(parentUpdateData);
          } else {
            const parentRes = await createStudio({
              variables: { input: parentInput },
            });
            input.parent_id = parentRes.data?.studioCreate?.id;
          }
        } catch (e) {
          handleSaveError(studioID, parentInput.name, apolloError(e));
        }
      }

      const updateData: GQL.StudioUpdateInput = {
        ...input,
        id: studioID,
      };
      updateData.stash_ids = await mergeStudioStashIDs(
        studioID,
        input.stash_ids ?? []
      );

      const res = await updateStudio(updateData);
      if (!res.data?.studioUpdate)
        handleSaveError(
          studioID,
          modalStudio?.name ?? "",
          res?.errors?.[0]?.message ?? ""
        );
    }
  };

  const renderStudios = () =>
    studios.map((studio) => {
      const isTagged = taggedStudios[studio.id];

      const stashID = studio.stash_ids.find((s) => {
        return s.endpoint === selectedEndpoint.endpoint;
      });

      let mainContent;
      if (!isTagged && stashID !== undefined) {
        mainContent = (
          <div className="text-left">
            <h5 className="text-bold">
              <FormattedMessage id="studio_tagger.studio_already_tagged" />
            </h5>
          </div>
        );
      } else if (!isTagged && !stashID) {
        mainContent = (
          <InputGroup>
            <Form.Control
              className="text-input"
              defaultValue={studio.name ?? ""}
              onChange={(e) =>
                setQueries({
                  ...queries,
                  [studio.id]: e.currentTarget.value,
                })
              }
              onKeyPress={(e: React.KeyboardEvent<HTMLInputElement>) =>
                e.key === "Enter" &&
                doBoxSearch(studio.id, queries[studio.id] ?? studio.name ?? "")
              }
            />
            <InputGroup.Append>
              <Button
                disabled={loading}
                onClick={() =>
                  doBoxSearch(
                    studio.id,
                    queries[studio.id] ?? studio.name ?? ""
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
              <FormattedMessage id="studio_tagger.studio_successfully_tagged" />
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
            href={`${base}studios/${stashID.stash_id}`}
          >
            {stashID.stash_id}
          </ExternalLink>
        ) : (
          <div className="small">{stashID.stash_id}</div>
        );

        subContent = (
          <div key={studio.id}>
            <InputGroup className="StudioTagger-box-link">
              <InputGroup.Text>{link}</InputGroup.Text>
              <InputGroup.Append>
                <Button
                  onClick={() =>
                    doBoxUpdate(studio.id, stashID.stash_id, stashID.endpoint)
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
            {error[studio.id] && (
              <div className="text-danger mt-1">
                <strong>
                  <span className="mr-2">Error:</span>
                  {error[studio.id]?.message}
                </strong>
                <div>{error[studio.id]?.details}</div>
              </div>
            )}
          </div>
        );
      } else if (searchErrors[studio.id]) {
        subContent = (
          <div className="text-danger font-weight-bold">
            {searchErrors[studio.id]}
          </div>
        );
      } else if (searchResults[studio.id]?.length === 0) {
        subContent = (
          <div className="text-danger font-weight-bold">
            <FormattedMessage id="studio_tagger.no_results_found" />
          </div>
        );
      }

      let searchResult;
      if (searchResults[studio.id]?.length > 0 && !isTagged) {
        searchResult = (
          <StashSearchResult
            key={studio.id}
            stashboxStudios={searchResults[studio.id]}
            studio={studio}
            endpoint={selectedEndpoint.endpoint}
            onStudioTagged={handleTaggedStudio}
            excludedStudioFields={config.excludedStudioFields ?? []}
          />
        );
      }

      return (
        <div key={studio.id} className={`${CLASSNAME}-studio`}>
          {modalStudio && (
            <StudioModal
              closeModal={() => setModalStudio(undefined)}
              modalVisible={modalStudio.stored_id === studio.id}
              studio={modalStudio}
              handleStudioCreate={handleStudioUpdate}
              excludedStudioFields={config.excludedStudioFields}
              icon={faTags}
              header={intl.formatMessage({
                id: "studio_tagger.update_studio",
              })}
              endpoint={selectedEndpoint.endpoint}
            />
          )}
          <div className={`${CLASSNAME}-details`}>
            <div></div>
            <div>
              <Card className="studio-card">
                <img loading="lazy" src={studio.image_path ?? ""} alt="" />
              </Card>
            </div>
            <div className={`${CLASSNAME}-details-text`}>
              <Link
                to={`/studios/${studio.id}`}
                className={`${CLASSNAME}-header`}
              >
                <h2>{studio.name}</h2>
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
        <StudioBatchUpdateModal
          close={() => setShowBatchUpdate(false)}
          isIdle={isIdle}
          selectedEndpoint={selectedEndpoint}
          studios={studios}
          onBatchUpdate={handleBatchUpdate}
          batchAddParents={batchAddParents}
          setBatchAddParents={setBatchAddParents}
        />
      )}

      {showBatchAdd && (
        <StudioBatchAddModal
          close={() => setShowBatchAdd(false)}
          isIdle={isIdle}
          onBatchAdd={handleBatchAdd}
          batchAddParents={batchAddParents}
          setBatchAddParents={setBatchAddParents}
        />
      )}
      <div className="ml-auto mb-3">
        <Button onClick={() => setShowBatchAdd(true)}>
          <FormattedMessage id="studio_tagger.batch_add_studios" />
        </Button>
        <Button className="ml-3" onClick={() => setShowBatchUpdate(true)}>
          <FormattedMessage id="studio_tagger.batch_update_studios" />
        </Button>
      </div>
      <div className={CLASSNAME}>{renderStudios()}</div>
    </Card>
  );
};

interface ITaggerProps {
  studios: GQL.StudioDataFragment[];
}

export const StudioTagger: React.FC<ITaggerProps> = ({ studios }) => {
  const jobsSubscribe = useJobsSubscribe();
  const intl = useIntl();
  const { configuration: stashConfig } = useConfigurationContext();
  const { config, setConfig } = useTaggerConfig();
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

      // Once the studio batch is complete, refresh all local studio data
      const ac = getClient();
      evictQueries(ac.cache, studioMutationImpactedQueries);
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

  async function batchAdd(studioInput: string, createParent: boolean) {
    if (studioInput && selectedEndpoint) {
      const inputs = studioInput
        .split(",")
        .map((n) => n.trim())
        .filter((n) => n.length > 0);

      const { names, stashIds } = separateNamesAndStashIds(inputs);

      if (names.length > 0 || stashIds.length > 0) {
        const ret = await mutateStashBoxBatchStudioTag({
          names: names.length > 0 ? names : undefined,
          stash_ids: stashIds.length > 0 ? stashIds : undefined,
          endpoint: selectedEndpointIndex,
          refresh: false,
          exclude_fields: config?.excludedStudioFields ?? [],
          createParent: createParent,
        });

        setBatchJobID(ret.data?.stashBoxBatchStudioTag);
      }
    }
  }

  async function batchUpdate(
    ids: string[] | undefined,
    refresh: boolean,
    createParent: boolean
  ) {
    if (selectedEndpoint) {
      const ret = await mutateStashBoxBatchStudioTag({
        ids: ids,
        endpoint: selectedEndpointIndex,
        refresh,
        exclude_fields: config?.excludedStudioFields ?? [],
        createParent: createParent,
      });

      setBatchJobID(ret.data?.stashBoxBatchStudioTag);
    }
  }

  // const progress =
  //   jobStatus.data?.metadataUpdate.status ===
  //     "Stash-Box Studio Batch Operation" &&
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
            <FormattedMessage id="studio_tagger.status_tagging_studios" />
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
            <FormattedMessage id="studio_tagger.status_tagging_job_queued" />
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

            <StudioConfig
              config={config}
              setConfig={setConfig}
              show={showConfig}
            />
            <StudioTaggerList
              studios={studios}
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
              <FormattedMessage id="studio_tagger.to_use_the_studio_tagger" />
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
