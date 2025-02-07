import {
  faBan,
  faCheck,
  faCircle,
  faCircleExclamation,
  faCog,
  faHourglassStart,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import moment from "moment/min/moment-with-locales";
import React, { useEffect, useState } from "react";
import { Button, Card, ProgressBar } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { Icon } from "src/components/Shared/Icon";
import {
  mutateStopJob,
  useJobQueue,
  useJobsSubscribe,
} from "src/core/StashService";
import * as GQL from "src/core/generated-graphql";

type JobFragment = Pick<
  GQL.Job,
  | "id"
  | "status"
  | "subTasks"
  | "description"
  | "progress"
  | "error"
  | "startTime"
>;

interface IJob {
  job: JobFragment;
}

const Task: React.FC<IJob> = ({ job }) => {
  const [stopping, setStopping] = useState(false);
  const [className, setClassName] = useState("");

  useEffect(() => {
    setTimeout(() => setClassName("fade-in"));
  }, []);

  useEffect(() => {
    if (
      job.status === GQL.JobStatus.Cancelled ||
      job.status === GQL.JobStatus.Failed ||
      job.status === GQL.JobStatus.Finished
    ) {
      // fade out around 10 seconds
      setTimeout(() => {
        setClassName("fade-out");
      }, 9800);
    }
  }, [job]);

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

  function getStatusClass() {
    switch (job.status) {
      case GQL.JobStatus.Ready:
        return "ready";
      case GQL.JobStatus.Running:
        return "running";
      case GQL.JobStatus.Stopping:
        return "stopping";
      case GQL.JobStatus.Finished:
        return "finished";
      case GQL.JobStatus.Cancelled:
        return "cancelled";
      case GQL.JobStatus.Failed:
        return "failed";
    }
  }

  function getStatusIcon() {
    let icon = faCircle;
    let iconClass = "";
    switch (job.status) {
      case GQL.JobStatus.Ready:
        icon = faHourglassStart;
        break;
      case GQL.JobStatus.Running:
        icon = faCog;
        iconClass = "fa-spin";
        break;
      case GQL.JobStatus.Stopping:
        icon = faCog;
        iconClass = "fa-spin";
        break;
      case GQL.JobStatus.Finished:
        icon = faCheck;
        break;
      case GQL.JobStatus.Cancelled:
        icon = faBan;
        break;
      case GQL.JobStatus.Failed:
        icon = faCircleExclamation;
        break;
    }

    return <Icon icon={icon} className={`fa-fw ${iconClass}`} />;
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

  function maybeRenderETA() {
    if (
      job.status === GQL.JobStatus.Running &&
      job.startTime !== null &&
      job.startTime !== undefined &&
      job.progress !== null &&
      job.progress !== undefined &&
      job.progress > 0
    ) {
      const now = new Date();
      const start = new Date(job.startTime);
      const nowMS = now.valueOf();
      const startMS = start.valueOf();
      const estimatedLength = (nowMS - startMS) / job.progress;
      const estLenStr = moment.duration(estimatedLength).humanize();
      return (
        <span className="job-eta">
          <FormattedMessage id="eta" />: {estLenStr}
        </span>
      );
    }
  }

  function maybeRenderSubTasks() {
    if (
      job.status === GQL.JobStatus.Running ||
      job.status === GQL.JobStatus.Stopping
    ) {
      return (
        <div>
          {/* eslint-disable react/no-array-index-key */}
          {(job.subTasks ?? []).map((t, i) => (
            <div className="job-subtask" key={i}>
              {t}
            </div>
          ))}
          {/* eslint-enable react/no-array-index-key */}
        </div>
      );
    }

    if (job.status === GQL.JobStatus.Failed && job.error) {
      return <div className="job-error">{job.error}</div>;
    }
  }

  return (
    <li className={`job ${className}`}>
      <div>
        <Button
          className="minimal stop"
          size="sm"
          onClick={() => stopJob()}
          disabled={!canStop()}
        >
          <Icon icon={faTimes} />
        </Button>
        <div className={`job-status ${getStatusClass()}`}>
          <div className="job-description">
            <div>
              {getStatusIcon()}
              <span>{job.description}</span>
            </div>
            {maybeRenderETA()}
          </div>
          <div>{maybeRenderProgress()}</div>
          {maybeRenderSubTasks()}
        </div>
      </div>
    </li>
  );
};

export const JobTable: React.FC = () => {
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

  return (
    <Card className="job-table">
      <ul>
        {!queue?.length ? (
          <span className="empty-queue-message">
            {intl.formatMessage({ id: "config.tasks.empty_queue" })}
          </span>
        ) : undefined}
        {(queue ?? []).map((j) => (
          <Task job={j} key={j.id} />
        ))}
      </ul>
    </Card>
  );
};
