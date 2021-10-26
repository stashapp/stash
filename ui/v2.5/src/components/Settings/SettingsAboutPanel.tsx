import React from "react";
import { Button, Table } from "react-bootstrap";
import { useIntl } from "react-intl";
import { LoadingIndicator } from "src/components/Shared";
import { useLatestVersion } from "src/core/StashService";

export const SettingsAboutPanel: React.FC = () => {
  const gitHash = process.env.VITE_APP_GITHASH;
  const stashVersion = process.env.VITE_APP_STASH_VERSION;
  const buildTime = process.env.VITE_APP_DATE;

  const intl = useIntl();

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
        <td>{intl.formatMessage({ id: "config.about.version" })}:</td>
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
          <strong>
            {dataLatest.latestversion.shorthash}{" "}
            {intl.formatMessage({ id: "config.about.new_version_notice" })}{" "}
          </strong>
          <a href={dataLatest.latestversion.url}>
            {intl.formatMessage({ id: "actions.download" })}
          </a>
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
            <td>
              {intl.formatMessage({
                id: "config.about.latest_version_build_hash",
              })}{" "}
            </td>
            <td>{maybeRenderLatestVersion()} </td>
          </tr>
          <tr>
            <td>
              <Button onClick={() => refetch()}>
                {intl.formatMessage({
                  id: "config.about.check_for_new_version",
                })}
              </Button>
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
              <td>{intl.formatMessage({ id: "config.about.build_hash" })}</td>
              <td>{gitHash}</td>
            </tr>
            <tr>
              <td>{intl.formatMessage({ id: "config.about.build_time" })}</td>
              <td>{buildTime}</td>
            </tr>
          </tbody>
        </Table>
      </>
    );
  }
  return (
    <>
      <h4>{intl.formatMessage({ id: "config.categories.about" })}</h4>
      <Table>
        <tbody>
          <tr>
            <td>
              {intl.formatMessage(
                { id: "config.about.stash_home" },
                {
                  url: (
                    <a
                      href="https://github.com/stashapp/stash"
                      rel="noopener noreferrer"
                      target="_blank"
                    >
                      GitHub
                    </a>
                  ),
                }
              )}
            </td>
          </tr>
          <tr>
            <td>
              {intl.formatMessage(
                { id: "config.about.stash_wiki" },
                {
                  url: (
                    <a
                      href="https://github.com/stashapp/stash/wiki"
                      rel="noopener noreferrer"
                      target="_blank"
                    >
                      Wiki
                    </a>
                  ),
                }
              )}
            </td>
          </tr>
          <tr>
            <td>
              {intl.formatMessage(
                { id: "config.about.stash_discord" },
                {
                  url: (
                    <a
                      href="https://discord.gg/2TsNFKt"
                      rel="noopener noreferrer"
                      target="_blank"
                    >
                      Discord
                    </a>
                  ),
                }
              )}
            </td>
          </tr>
          <tr>
            <td>
              {intl.formatMessage(
                { id: "config.about.stash_open_collective" },
                {
                  url: (
                    <a
                      href="https://opencollective.com/stashapp"
                      rel="noopener noreferrer"
                      target="_blank"
                    >
                      Open Collective
                    </a>
                  ),
                }
              )}
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
