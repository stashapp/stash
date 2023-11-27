import { useEffect, useState } from "react";
import {
  Job,
  JobStatusUpdateType,
  useJobQueueQuery,
  useJobsSubscribeSubscription,
} from "src/core/generated-graphql";

export type JobFragment = Pick<
  Job,
  "id" | "status" | "subTasks" | "description" | "progress"
>;

export const useMonitorJob = (
  jobID: string | undefined | null,
  onComplete?: () => void
) => {
  const jobsSubscribe = useJobsSubscribeSubscription({
    skip: !jobID,
  });
  const { data: jobData, loading } = useJobQueueQuery({
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

    const j = jobData?.jobQueue?.find((jj) => jj.id === jobID);
    if (j) {
      setJob(j);
    } else {
      // must've already finished
      setJob(undefined);
      if (onComplete) {
        onComplete();
      }
    }
  }, [jobID, jobData, loading, onComplete]);

  // monitor batch operation
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
      if (onComplete) {
        onComplete();
      }
    }
  }, [jobsSubscribe, jobID, onComplete]);

  return { job };
};
