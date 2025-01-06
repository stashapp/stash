import React, { useEffect, useMemo, useState } from "react";
import { Button, Card, Container, Form, ProgressBar } from "react-bootstrap";
import { useIntl, FormattedMessage } from "react-intl";
import { useHistory } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import {
  useSystemStatus,
  mutateMigrate,
  postMigrate,
  refetchSystemStatus,
} from "src/core/StashService";
import { migrationNotes } from "src/docs/en/MigrationNotes";
import { ExternalLink } from "../Shared/ExternalLink";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { MarkdownPage } from "../Shared/MarkdownPage";
import { JobFragment, useMonitorJob } from "src/utils/job";

export const Migrate: React.FC = () => {
  const intl = useIntl();
  const history = useHistory();

  const { data: systemStatus, loading } = useSystemStatus();

  const [backupPath, setBackupPath] = useState<string | undefined>();
  const [migrateLoading, setMigrateLoading] = useState(false);
  const [migrateError, setMigrateError] = useState("");

  const [jobID, setJobID] = useState<string | undefined>();

  function onJobFinished(finishedJob?: JobFragment) {
    setJobID(undefined);
    setMigrateLoading(false);

    if (finishedJob?.error) {
      setMigrateError(finishedJob.error);
    } else {
      postMigrate();
      // refetch the system status so that the we get redirected
      refetchSystemStatus();
    }
  }

  const { job } = useMonitorJob(jobID, onJobFinished);

  // if database path includes path separators, then this is passed through
  // to the migration path. Extract the base name of the database file.
  const databasePath = systemStatus
    ? systemStatus?.systemStatus.databasePath?.split(/[\\/]/).pop()
    : "";

  // make suffix based on current time
  const now = new Date()
    .toISOString()
    .replace(/T/g, "_")
    .replace(/-/g, "")
    .replace(/:/g, "")
    .replace(/\..*/, "");
  const defaultBackupPath = systemStatus
    ? `${databasePath}.${systemStatus.systemStatus.databaseSchema}.${now}`
    : "";

  const discordLink = (
    <ExternalLink href="https://discord.gg/2TsNFKt">Discord</ExternalLink>
  );
  const githubLink = (
    <ExternalLink href="https://github.com/stashapp/stash/issues">
      <FormattedMessage id="setup.github_repository" />
    </ExternalLink>
  );

  useEffect(() => {
    if (backupPath === undefined && defaultBackupPath) {
      setBackupPath(defaultBackupPath);
    }
  }, [defaultBackupPath, backupPath]);

  const status = systemStatus?.systemStatus;

  const maybeMigrationNotes = useMemo(() => {
    if (
      !status ||
      status.databaseSchema === undefined ||
      status.databaseSchema === null ||
      status.appSchema === undefined ||
      status.appSchema === null
    )
      return;

    const notes = [];
    for (let i = status.databaseSchema + 1; i <= status.appSchema; ++i) {
      const note = migrationNotes[i];
      if (note) {
        notes.push(note);
      }
    }

    if (notes.length === 0) return;

    return (
      <div className="migration-notes">
        <h2>
          <FormattedMessage id="setup.migrate.migration_notes" />
        </h2>
        <div>
          {notes.map((n, i) => (
            <div key={i}>
              <MarkdownPage page={n} />
            </div>
          ))}
        </div>
      </div>
    );
  }, [status]);

  // only display setup wizard if system is not setup
  if (loading || !systemStatus || !status) {
    return <LoadingIndicator />;
  }

  if (migrateLoading) {
    const progress =
      job && job.progress !== undefined && job.progress !== null
        ? job.progress * 100
        : undefined;

    return (
      <div className="migrate-loading-status">
        <h4>
          <LoadingIndicator inline small message="" />
          <span>
            <FormattedMessage id="setup.migrate.migrating_database" />
          </span>
        </h4>
        {progress !== undefined && (
          <ProgressBar
            animated
            now={progress}
            label={`${progress.toFixed(0)}%`}
          />
        )}
        {job?.subTasks?.map((subTask, i) => (
          <div key={i}>
            <p>{subTask}</p>
          </div>
        ))}
      </div>
    );
  }

  if (
    systemStatus.systemStatus.status !== GQL.SystemStatusEnum.NeedsMigration
  ) {
    // redirect to main page
    history.replace("/");
    return <LoadingIndicator />;
  }

  async function onMigrate() {
    try {
      setMigrateLoading(true);
      setMigrateError("");

      // migrate now uses the job manager
      const ret = await mutateMigrate({
        backupPath: backupPath ?? "",
      });

      setJobID(ret.data?.migrate);
    } catch (e) {
      if (e instanceof Error) setMigrateError(e.message ?? e.toString());
      setMigrateLoading(false);
    }
  }

  function maybeRenderError() {
    if (!migrateError) {
      return;
    }

    return (
      <section>
        <h2 className="text-danger">
          <FormattedMessage id="setup.migrate.migration_failed" />
        </h2>

        <p>
          <FormattedMessage id="setup.migrate.migration_failed_error" />
        </p>

        <Card>
          <pre>{migrateError}</pre>
        </Card>

        <p>
          <FormattedMessage
            id="setup.migrate.migration_failed_help"
            values={{ discordLink, githubLink }}
          />
        </p>
      </section>
    );
  }

  return (
    <Container>
      <h1 className="text-center mb-3">
        <FormattedMessage id="setup.migrate.migration_required" />
      </h1>
      <Card>
        <section>
          <p>
            <FormattedMessage
              id="setup.migrate.schema_too_old"
              values={{
                databaseSchema: <strong>{status.databaseSchema}</strong>,
                appSchema: <strong>{status.appSchema}</strong>,
                strong: (chunks: string) => <strong>{chunks}</strong>,
                code: (chunks: string) => <code>{chunks}</code>,
              }}
            />
          </p>

          <p className="lead text-center my-5">
            <FormattedMessage id="setup.migrate.migration_irreversible_warning" />
          </p>

          <p>
            <FormattedMessage
              id="setup.migrate.backup_recommended"
              values={{
                defaultBackupPath,
                code: (chunks: string) => <code>{chunks}</code>,
              }}
            />
          </p>
        </section>

        {maybeMigrationNotes}

        <section>
          <Form.Group id="migrate">
            <Form.Label>
              <FormattedMessage id="setup.migrate.backup_database_path_leave_empty_to_disable_backup" />
            </Form.Label>
            <Form.Control
              className="text-input"
              name="backupPath"
              defaultValue={backupPath}
              placeholder={intl.formatMessage({
                id: "setup.paths.database_filename_empty_for_default",
              })}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setBackupPath(e.currentTarget.value)
              }
            />
          </Form.Group>
        </section>

        <section>
          <div className="d-flex justify-content-center">
            <Button variant="primary mx-2 p-5" onClick={() => onMigrate()}>
              <FormattedMessage id="setup.migrate.perform_schema_migration" />
            </Button>
          </div>
        </section>

        {maybeRenderError()}
      </Card>
    </Container>
  );
};

export default Migrate;
