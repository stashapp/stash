import React from "react";
import { Button } from "react-bootstrap";
import { useIntl } from "react-intl";
import { useLatestVersion } from "src/core/StashService";
import { ConstantSetting, Setting, SettingGroup } from "./Inputs";
import { SettingSection } from "./SettingSection";

export const SettingsAboutPanel: React.FC = () => {
  const gitHash = import.meta.env.VITE_APP_GITHASH;
  const stashVersion = import.meta.env.VITE_APP_STASH_VERSION;
  const buildTime = import.meta.env.VITE_APP_DATE;

  const intl = useIntl();

  const {
    data: dataLatest,
    error: errorLatest,
    loading: loadingLatest,
    refetch,
    networkStatus,
  } = useLatestVersion();

  const hasNew = dataLatest && gitHash !== dataLatest.latestversion.shorthash;

  return (
    <>
      <SettingSection headingID="config.about.version">
        <SettingGroup
          settingProps={{
            heading: stashVersion,
          }}
        >
          <ConstantSetting
            headingID="config.about.build_hash"
            value={gitHash}
          />
          <ConstantSetting
            headingID="config.about.build_time"
            value={buildTime}
          />
        </SettingGroup>

        <SettingGroup
          settingProps={{
            headingID: "config.about.latest_version",
          }}
        >
          {errorLatest ? (
            <Setting heading={errorLatest.message} />
          ) : !dataLatest || loadingLatest || networkStatus === 4 ? (
            <Setting headingID="loading.generic" />
          ) : (
            <div className="setting">
              <div>
                <h3>
                  {intl.formatMessage({
                    id: "config.about.latest_version_build_hash",
                  })}
                </h3>
                <div className="value">
                  {dataLatest.latestversion.shorthash}{" "}
                  {hasNew
                    ? intl.formatMessage({
                        id: "config.about.new_version_notice",
                      })
                    : undefined}
                </div>
              </div>
              <div>
                <a href={dataLatest.latestversion.url}>
                  <Button>
                    {intl.formatMessage({ id: "actions.download" })}
                  </Button>
                </a>
                <Button onClick={() => refetch()}>
                  {intl.formatMessage({
                    id: "config.about.check_for_new_version",
                  })}
                </Button>
              </div>
            </div>
          )}
        </SettingGroup>
      </SettingSection>

      <SettingSection headingID="config.categories.about">
        <div className="setting">
          <div>
            <p>
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
            </p>
            <p>
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
            </p>
            <p>
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
            </p>
            <p>
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
            </p>
          </div>
          <div />
        </div>
      </SettingSection>
    </>
  );
};
