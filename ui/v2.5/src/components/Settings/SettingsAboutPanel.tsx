import React from "react";
import { Button, Table } from "react-bootstrap";
import { LoadingIndicator } from "src/components/Shared";
import { useLatestVersion } from "src/core/StashService";

export const SettingsAboutPanel: React.FC = () => {
  const gitHash = process.env.REACT_APP_GITHASH;
  const stashVersion = process.env.REACT_APP_STASH_VERSION;
  const buildTime = process.env.REACT_APP_DATE;

  const {
    data: dataLatest,
    error: errorLatest,
    loading: loadingLatest,
    refetch,
    networkStatus,
  } = useLatestVersion();

  function maybeRenderTag() {
    if (!stashVersion) {
      return;
    }
    return (
      <tr>
        <td>Version:</td>
        <td>{stashVersion}</td>
      </tr>
    );
  }

  function maybeRenderLatestVersion() {
    if (
      !dataLatest?.latestversion.shorthash ||
      !dataLatest?.latestversion.url
    ) {
      return;
    }

    if (gitHash !== dataLatest.latestversion.shorthash) {
      return (
        <>
          <strong>{dataLatest.latestversion.shorthash} [NEW] </strong>
          <a href={dataLatest.latestversion.url}>Download</a>
        </>
      );
    }

    return <>{dataLatest.latestversion.shorthash}</>;
  }

  function renderLatestVersion() {
    return (
      <Table>
        <tbody>
          <tr>
            <td>Latest Version Build Hash: </td>
            <td>{maybeRenderLatestVersion()} </td>
          </tr>
          <tr>
            <td>
              <Button onClick={() => refetch()}>Check for new version</Button>
            </td>
          </tr>
        </tbody>
      </Table>
    );
  }

  function renderVersion() {
    return (
      <>
        <Table>
          <tbody>
            {maybeRenderTag()}
            <tr>
              <td>Build hash:</td>
              <td>{gitHash}</td>
            </tr>
            <tr>
              <td>Build time:</td>
              <td>{buildTime}</td>
            </tr>
          </tbody>
        </Table>
      </>
    );
  }
  return (
    <>
      <h4>About</h4>
      <Table>
        <tbody>
          <tr>
            <td>
              Stash home at{" "}
              <a
                href="https://github.com/stashapp/stash"
                rel="noopener noreferrer"
                target="_blank"
              >
                Github
              </a>
            </td>
          </tr>
          <tr>
            <td>
              Stash{" "}
              <a
                href="https://github.com/stashapp/stash/wiki"
                rel="noopener noreferrer"
                target="_blank"
              >
                Wiki
              </a>{" "}
              page
            </td>
          </tr>
          <tr>
            <td>
              Join our{" "}
              <a
                href="https://discord.gg/2TsNFKt"
                rel="noopener noreferrer"
                target="_blank"
              >
                Discord
              </a>{" "}
              channel
            </td>
          </tr>
          <tr>
            <td>
              Support us through{" "}
              <a
                href="https://opencollective.com/stashapp"
                rel="noopener noreferrer"
                target="_blank"
              >
                Open Collective
              </a>
            </td>
          </tr>
        </tbody>
      </Table>
      {errorLatest && <span>{errorLatest.message}</span>}
      {renderVersion()}
      {!dataLatest || loadingLatest || networkStatus === 4 ? (
        <LoadingIndicator inline />
      ) : (
        renderLatestVersion()
      )}
    </>
  );
};
