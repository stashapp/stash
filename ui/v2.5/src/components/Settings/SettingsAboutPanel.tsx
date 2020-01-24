import React from "react";
import { Button, Table, Spinner } from "react-bootstrap";
import { StashService } from "src/core/StashService";

export const SettingsAboutPanel: React.FC = () => {
  const { data, error, loading } = StashService.useVersion();
  const { data: dataLatest, error: errorLatest, loading: loadingLatest, refetch, networkStatus } = StashService.useLatestVersion();

  function maybeRenderTag() {
    if (!data || !data.version || !data.version.version) {
      return;
    }
    return (
      <tr>
        <td>Version:</td>
        <td>{data.version.version}</td>
      </tr>
    );
  }

  function maybeRenderLatestVersion() {
    if (!dataLatest || !dataLatest.latestversion || !dataLatest.latestversion.shorthash || !dataLatest.latestversion.url) { return; }
    if (!data || !data.version || !data.version.hash) {
      return (
        <>{dataLatest.latestversion.shorthash}</>
      );
    }

    if (data.version.hash !== dataLatest.latestversion.shorthash) {
      return (
        <>
          <strong>{dataLatest.latestversion.shorthash} [NEW] </strong><a href={dataLatest.latestversion.url}>Download</a>
        </>
      );
    }

    return (
      <>{dataLatest.latestversion.shorthash}</>
    );
  }

  function renderLatestVersion() {
    if (!data || !data.version || !data.version.version) { return; } // if there is no "version" latest version check is obviously not supported
    return (
      <Table>
        <tbody>
          <tr>
						<td>Latest Version Build Hash: </td>
            <td>{maybeRenderLatestVersion()} </td>
          </tr>
          <tr>
            <td><Button onClick={() => refetch()}>Check for new version</Button></td>
          </tr>
        </tbody>
      </Table>
    );
  }

  function renderVersion() {
    if (!data || !data.version) {
      return;
    }
    return (
      <>
        <Table>
          <tbody>
            {maybeRenderTag()}
            <tr>
              <td>Build hash:</td>
              <td>{data.version.hash}</td>
            </tr>
            <tr>
              <td>Build time:</td>
              <td>{data.version.build_time}</td>
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
            <td>Stash home at <a href="https://github.com/stashapp/stash" rel="noopener noreferrer" target="_blank">Github</a></td>
          </tr>
          <tr>
            <td>Stash <a href="https://github.com/stashapp/stash/wiki" rel="noopener noreferrer" target="_blank">Wiki</a> page</td>
          </tr>
          <tr>
            <td>Join our <a href="https://discord.gg/2TsNFKt" rel="noopener noreferrer" target="_blank">Discord</a> channel</td>
          </tr>
          <tr>
            <td>Support us through <a href="https://opencollective.com/stashapp" rel="noopener noreferrer" target="_blank">Open Collective</a></td>
          </tr>
        </tbody>
      </Table>
      {!data || loading ? <Spinner animation="border" variant="light" /> : ""}
      {error && <span>{error.message}</span>}
      {errorLatest && <span>{errorLatest.message}</span>}
      {renderVersion()}
      {!dataLatest || loadingLatest || networkStatus === 4 ? <Spinner animation="border" variant="light" /> : <>{renderLatestVersion()}</>}
    </>
  );
};
