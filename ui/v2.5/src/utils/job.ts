import { useEffect, useState } from "react";
import { getWSClient, useWSState } from "src/core/StashService";
import {
  Job,
  JobStatus,
  JobStatusUpdateType,
  useFindJobQuery,
  useJobsSubscribeSubscription,
} from "src/core/generated-graphql";

export type JobFragment = Pick<
  Job,
  "id" | "status" | "subTasks" | "description" | "progress" | "error"
>;

export const useMonitorJob = (
  jobID: string | undefined | null,
  onComplete?: (job?: JobFragment) => void
) => {
  const { state } = useWSState(getWSClient());

  const jobsSubscribe = useJobsSubscribeSubscription({
    skip: !jobID,
  });
  const {
    data: jobData,
    loading,
    startPolling,
    stopPolling,
  } = useFindJobQuery({
    variables: {
      input: { id: jobID ?? "" },
    },
    fetchPolicy: "network-only",
    skip: !jobID,
  });

  const [job, setJob] = useState<JobFragment | undefined>();

  useEffect(() => {
    if (!jobID) {
      return;
    }

    if (loading) {
      return;
    }

    const j = jobData?.findJob;
    if (j) {
      setJob(j);

      if (
        j.status === JobStatus.Finished ||
        j.status === JobStatus.Failed ||
        j.status === JobStatus.Cancelled
      ) {
        setJob(undefined);
        onComplete?.(j);
      }
    } else {
      // must've already finished
      setJob(undefined);
      onComplete?.();
    }
  }, [jobID, jobData, loading, onComplete]);

  // monitor job
  useEffect(() => {
    if (!jobID) {
      return;
    }

    if (!jobsSubscribe.data) {
      return;
    }

    const event = jobsSubscribe.data.jobsSubscribe;
    if (event.job.id !== jobID) {
      return;
    }

    if (event.type !== JobStatusUpdateType.Remove) {
      setJob(event.job);
    } else {
      setJob(undefined);
      onComplete?.(event.job);
    }
  }, [jobsSubscribe, jobID, onComplete]);

  // it's possible that the websocket connection isn't present
  // in that case, we'll just poll the server
  useEffect(() => {
    if (!jobID) {
      stopPolling();
      return;
    }

    if (state === "connected") {
      stopPolling();
    } else {
      const defaultPollInterval = 1000;
      startPolling(defaultPollInterval);
    }
  }, [jobID, state, startPolling, stopPolling]);

  return { job };
};
