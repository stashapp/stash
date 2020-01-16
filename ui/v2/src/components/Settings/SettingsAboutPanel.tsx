import {
  Button,
  H4,
  HTMLTable,
  Spinner,
} from "@blueprintjs/core";
import React, { FunctionComponent } from "react";
import { StashService } from "../../core/StashService";

interface IProps { }

export const SettingsAboutPanel: FunctionComponent<IProps> = (props: IProps) => {
  const { data, error, loading } = StashService.useVersion();
  const { data: dataLatest, error: errorLatest, loading: loadingLatest, refetch, networkStatus } = StashService.useLatestVersion();

  function maybeRenderTag() {
    if (!data || !data.version || !data.version.version) { return; }
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
    if (!data || !data.version || !data.version.version) { return; } //if there is no "version" latest version check is obviously not supported
    return (
      <HTMLTable>
        <tbody>
          <tr>
            <td>Latest Version Build Hash: </td>
            <td>{maybeRenderLatestVersion()} </td>
          </tr>
          <tr>
            <td><Button onClick={() => refetch()} text="Check for new version" /></td>
          </tr>
        </tbody>
      </HTMLTable>
    );
  }

  function renderVersion() {
    if (!data || !data.version) { return; }
    return (
      <>
        <HTMLTable>
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
        </HTMLTable>
      </>
    );
  }
  return (
    <>
      <H4>About</H4>
      <HTMLTable>
        <tbody>
          <tr>
            <td>Stash home at <a href="https://github.com/stashapp/stash">Github</a></td>
          </tr>
          <tr>
            <td>Stash <a href="https://github.com/stashapp/stash/wiki">Wiki</a> page</td>
          </tr>
          <tr>
            <td>Join our <a href="https://discord.gg/2TsNFKt">Discord</a> channel</td>
          </tr>
          <tr>
            <td>Support us through <a href="https://opencollective.com/stashapp">Open Collective</a></td>
          </tr>
        </tbody>
      </HTMLTable>
      {!data || loading ? <Spinner size={Spinner.SIZE_LARGE} /> : undefined}
      {!!error ? <span>{error.message}</span> : undefined}
      {!!errorLatest ? <span>{errorLatest.message}</span> : undefined}
      {renderVersion()}
      {!dataLatest || loadingLatest || networkStatus === 4 ? <Spinner size={Spinner.SIZE_SMALL} /> : <>{renderLatestVersion()}</>}
    </>
  );
};
