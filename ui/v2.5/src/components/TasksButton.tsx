import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useEffect, useState } from "react";
import { Button } from "react-bootstrap";
import { useJobQueue, useJobsSubscribe } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "./Shared";

type JobFragment = Pick<
  GQL.Job,
  "id" | "status" | "subTasks" | "description" | "progress"
>;

export const TasksButton: React.FC = () => {
  const jobStatus = useJobQueue();
  const jobsSubscribe = useJobsSubscribe();

  const [queue, setQueue] = useState<JobFragment[]>([]);

  useEffect(() => {
    setQueue(jobStatus.data?.jobQueue ?? []);
  }, [jobStatus]);

  useEffect(() => {
    if (!jobsSubscribe.data) {
      return;
    }

    const event = jobsSubscribe.data.jobsSubscribe;

    function updateJob() {
      setQueue((q) =>
        q.map((j) => {
          if (j.id === event.job.id) {
            return event.job;
          }

          return j;
        })
      );
    }

    switch (event.type) {
      case GQL.JobStatusUpdateType.Add:
        // add to the end of the queue
        setQueue((q) => q.concat([event.job]));
        break;
      case GQL.JobStatusUpdateType.Remove:
        setQueue((q) => q.filter((j) => j.id !== event.job.id));
        break;
      case GQL.JobStatusUpdateType.Update:
        updateJob();
        break;
    }
  }, [jobsSubscribe.data]);

  return (
    <Button className="minimal d-flex align-items-center h-100" title="Tasks">
      {queue.length > 0 ? (
        <FontAwesomeIcon icon="spinner" pulse />
      ) : (
        <Icon icon="bolt" />
      )}
    </Button>
  );
};
