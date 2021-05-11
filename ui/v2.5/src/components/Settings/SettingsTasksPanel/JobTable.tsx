import React, { useState, useEffect } from "react";
import { Button, ProgressBar } from "react-bootstrap";
import {
  mutateStopJob,
  useJobQueue,
  useJobsSubscribe,
} from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";
import { Icon } from "src/components/Shared";
import { IconProp } from "@fortawesome/fontawesome-svg-core";

interface IJob {
  job: GQL.JobDataFragment;
}

const Task: React.FC<IJob> = ({ job }) => {
  const [stopping, setStopping] = useState(false);

  async function stopJob() {
    setStopping(true);
    await mutateStopJob(job.id);
  }

  function canStop() {
    return (
      !stopping &&
      (job.status === GQL.JobStatus.Ready ||
        job.status === GQL.JobStatus.Running)
    );
  }

  function getStatusIcon() {
    let icon: IconProp = "circle";
    let className = "";
    switch (job.status) {
      case GQL.JobStatus.Ready:
        icon = "hourglass-start";
        className = "ready";
        break;
      case GQL.JobStatus.Running:
        icon = "cog";
        className = "fa-spin running";
        break;
      case GQL.JobStatus.Stopping:
        icon = "cog";
        className = "fa-spin stopping";
        break;
      case GQL.JobStatus.Finished:
        icon = "check";
        className = "finished";
        break;
      case GQL.JobStatus.Cancelled:
        icon = "ban";
        className = "cancelled";
        break;
    }

    return <Icon icon={icon} className={`fa-fw ${className}`} />;
  }

  function maybeRenderProgress() {
    if (
      job.status === GQL.JobStatus.Running &&
      job.progress !== undefined &&
      job.progress !== null
    ) {
      const progress = job.progress * 100;
      return (
        <ProgressBar
          animated
          now={progress}
          label={`${progress.toFixed(0)}%`}
        />
      );
    }
  }

  return (
    <li className="job">
      <div>
        <Button
          className="minimal stop"
          size="sm"
          onClick={() => stopJob()}
          disabled={!canStop()}
        >
          <Icon icon="times" />
        </Button>
        <div className="job-status">
          <div>
            {getStatusIcon()}
            <span>{job.description}</span>
          </div>
          <div>{maybeRenderProgress()}</div>
        </div>
      </div>
    </li>
  );
};

export const JobTable: React.FC = () => {
  const jobStatus = useJobQueue();
  const jobsSubscribe = useJobsSubscribe();

  const [queue, setQueue] = useState<GQL.JobDataFragment[]>([]);

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
        // update the job then remove after a timeout
        updateJob();
        setTimeout(() => {
          setQueue((q) => q.filter((j) => j.id !== event.job.id));
        }, 10000);
        break;
      case GQL.JobStatusUpdateType.Update:
        updateJob();
        break;
    }
  }, [jobsSubscribe.data]);

  function getCurrentJob() {
    if (queue.length > 0) {
      return queue[0].description;
    }

    return "Idle";
  }

  return (
    <div className="job-table">
      <h5>Currently running: {getCurrentJob()}</h5>

      <ul>
        {(queue ?? []).map((j) => (
          <Task job={j} key={j.id} />
        ))}
      </ul>
    </div>
  );
};
