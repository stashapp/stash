import React, { useEffect, useState } from "react";
import { Button, Card, Container, Form } from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import { useSystemStatus, mutateMigrate } from "src/core/StashService";
import { LoadingIndicator } from "../Shared";

export const Migrate: React.FC = () => {
  const { data: systemStatus, loading } = useSystemStatus();
  const [backupPath, setBackupPath] = useState<string | undefined>();
  const [migrateLoading, setMigrateLoading] = useState(false);
  const [migrateError, setMigrateError] = useState("");

  // make suffix based on current time
  const now = new Date()
    .toISOString()
    .replace(/T/g, "_")
    .replace(/-/g, "")
    .replace(/:/g, "")
    .replace(/\..*/, "");
  const defaultBackupPath = systemStatus
    ? `${systemStatus.systemStatus.databasePath}.${systemStatus.systemStatus.databaseSchema}.${now}`
    : "";

  const discordLink = (
    <a href="https://discord.gg/2TsNFKt" target="_blank" rel="noreferrer">
      Discord
    </a>
  );
  const githubLink = (
    <a
      href="https://github.com/stashapp/stash/issues"
      target="_blank"
      rel="noreferrer"
    >
      Github repository
    </a>
  );

  useEffect(() => {
    if (backupPath === undefined && defaultBackupPath) {
      setBackupPath(defaultBackupPath);
    }
  }, [defaultBackupPath, backupPath]);

  // only display setup wizard if system is not setup
  if (loading || !systemStatus) {
    return <LoadingIndicator />;
  }

  if (migrateLoading) {
    return <LoadingIndicator message="Migrating database" />;
  }

  if (
    systemStatus.systemStatus.status !== GQL.SystemStatusEnum.NeedsMigration
  ) {
    // redirect to main page
    const newURL = new URL("/", window.location.toString());
    window.location.href = newURL.toString();
    return <LoadingIndicator />;
  }

  const status = systemStatus.systemStatus;

  async function onMigrate() {
    try {
      setMigrateLoading(true);
      setMigrateError("");
      await mutateMigrate({
        backupPath: backupPath ?? "",
      });

      const newURL = new URL("/", window.location.toString());
      window.location.href = newURL.toString();
    } catch (e) {
      setMigrateError(e.message ?? e.toString());
      setMigrateLoading(false);
    }
  }

  function maybeRenderError() {
    if (!migrateError) {
      return;
    }

    return (
      <section>
        <h2 className="text-danger">Migration failed</h2>

        <p>The following error was encountered while migrating the database:</p>

        <Card>
          <pre>{migrateError}</pre>
        </Card>

        <p>
          Please make any necessary corrections and try again. Otherwise, raise
          a bug on the {githubLink} or seek help in the {discordLink}.
        </p>
      </section>
    );
  }

  return (
    <Container>
      <h1 className="text-center mb-3">Migration required</h1>
      <Card>
        <section>
          <p>
            Your current stash database is schema version{" "}
            <strong>{status.databaseSchema}</strong> and needs to be migrated to
            version <strong>{status.appSchema}</strong>. This version of Stash
            will not function without migrating the database.
          </p>

          <p className="lead text-center my-5">
            The schema migration process is not reversible. Once the migration
            is performed, your database will be incompatible with previous
            versions of stash.
          </p>

          <p>
            It is recommended that you backup your existing database before you
            migrate. We can do this for you, making a copy of your writing a
            backup to <code>{defaultBackupPath}</code> if required.
          </p>
        </section>

        <section>
          <Form.Group id="migrate">
            <Form.Label>
              Backup database path (leave empty to disable backup):
            </Form.Label>
            <Form.Control
              className="text-input"
              name="backupPath"
              defaultValue={backupPath}
              placeholder="database filename (empty for default)"
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setBackupPath(e.currentTarget.value)
              }
            />
          </Form.Group>
        </section>

        <section>
          <div className="d-flex justify-content-center">
            <Button variant="primary mx-2 p-5" onClick={() => onMigrate()}>
              Perform schema migration
            </Button>
          </div>
        </section>

        {maybeRenderError()}
      </Card>
    </Container>
  );
};
