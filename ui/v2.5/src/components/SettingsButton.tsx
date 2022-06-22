import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useEffect, useState } from "react";
import { Button } from "react-bootstrap";
import { useJobQueue, useJobsSubscribe } from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { useIntl } from "react-intl";
import { faCog } from "@fortawesome/free-solid-svg-icons";

type JobFragment = Pick<
  GQL.Job,
  "id" | "status" | "subTasks" | "description" | "progress"
>;

export const SettingsButton: React.FC = () => {
  const intl = useIntl();
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
    <Button
      className="minimal d-flex align-items-center h-100"
      title={intl.formatMessage({ id: "settings" })}
    >
      <FontAwesomeIcon icon={faCog} spin={queue.length > 0} />
    </Button>
  );
};
