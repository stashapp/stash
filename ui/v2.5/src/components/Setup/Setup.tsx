import React, { useEffect, useState, useContext } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import {
  Alert,
  Button,
  Card,
  Container,
  Form,
  InputGroup,
} from "react-bootstrap";
import * as GQL from "src/core/generated-graphql";
import {
  mutateSetup,
  useConfigureUI,
  useSystemStatus,
} from "src/core/StashService";
import { Link, useHistory } from "react-router-dom";
import { ConfigurationContext } from "src/hooks/Config";
import StashConfiguration from "../Settings/StashConfiguration";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { ModalComponent } from "../Shared/Modal";
import { FolderSelectDialog } from "../Shared/FolderSelect/FolderSelectDialog";
import {
  faEllipsisH,
  faExclamationTriangle,
  faQuestionCircle,
} from "@fortawesome/free-solid-svg-icons";
import { releaseNotes } from "src/docs/en/ReleaseNotes";
import { ExternalLink } from "../Shared/ExternalLink";

export const Setup: React.FC = () => {
  const { configuration, loading: configLoading } =
    useContext(ConfigurationContext);
  const [saveUI] = useConfigureUI();

  const [step, setStep] = useState(0);
  const [setupInWorkDir, setSetupInWorkDir] = useState(false);
  const [stashes, setStashes] = useState<GQL.StashConfig[]>([]);
  const [showStashAlert, setShowStashAlert] = useState(false);
  const [databaseFile, setDatabaseFile] = useState("");
  const [generatedLocation, setGeneratedLocation] = useState("");
  const [cacheLocation, setCacheLocation] = useState("");
  const [storeBlobsInDatabase, setStoreBlobsInDatabase] = useState(false);
  const [blobsLocation, setBlobsLocation] = useState("");
  const [loading, setLoading] = useState(false);
  const [setupError, setSetupError] = useState<string>();

  const intl = useIntl();
  const history = useHistory();

  const [showGeneratedSelectDialog, setShowGeneratedSelectDialog] =
    useState(false);
  const [showCacheSelectDialog, setShowCacheSelectDialog] = useState(false);
  const [showBlobsDialog, setShowBlobsDialog] = useState(false);

  const { data: systemStatus, loading: statusLoading } = useSystemStatus();
  const status = systemStatus?.systemStatus;

  const windows = status?.os === "windows";
  const pathSep = windows ? "\\" : "/";
  const homeDir = windows ? "%USERPROFILE%" : "$HOME";
  const pwd = windows ? "%CD%" : "$PWD";

  function pathJoin(...paths: string[]) {
    return paths.join(pathSep);
  }

  // simply returns everything preceding the last path separator
  function pathDir(path: string) {
    const lastSep = path.lastIndexOf(pathSep);
    if (lastSep === -1) return "";
    return path.slice(0, lastSep);
  }

  const workingDir = status?.workingDir ?? ".";

  // When running Stash.app, the working directory is (usually) set to /.
  // Assume that the user doesn't want to set up in / (it's usually mounted read-only anyway),
  // so in this situation disallow setting up in the working directory.
  const macApp = status?.os === "darwin" && workingDir === "/";

  const fallbackStashDir = pathJoin(homeDir, ".stash");
  const fallbackConfigPath = pathJoin(fallbackStashDir, "config.yml");

  const overrideConfig = status?.configPath;
  const overrideGenerated = configuration?.general.generatedPath;
  const overrideCache = configuration?.general.cachePath;
  const overrideBlobs = configuration?.general.blobsPath;
  const overrideDatabase = configuration?.general.databasePath;

  useEffect(() => {
    if (configuration) {
      const configStashes = configuration.general.stashes;
      if (configStashes.length > 0) {
        setStashes(
          configStashes.map((s) => {
            const { __typename, ...withoutTypename } = s;
            return withoutTypename;
          })
        );
      }
    }
  }, [configuration]);

  const discordLink = (
    <ExternalLink href="https://discord.gg/2TsNFKt">Discord</ExternalLink>
  );
  const githubLink = (
    <ExternalLink href="https://github.com/stashapp/stash/issues">
      <FormattedMessage id="setup.github_repository" />
    </ExternalLink>
  );

  function onConfigLocationChosen(inWorkDir: boolean) {
    setSetupInWorkDir(inWorkDir);
    next();
  }

  function goBack(n?: number) {
    let dec = n;
    if (!dec) {
      dec = 1;
    }
    setStep(Math.max(0, step - dec));
  }

  function next() {
    setStep(step + 1);
  }

  function confirmPaths() {
    if (stashes.length > 0) {
      next();
      return;
    }

    setShowStashAlert(true);
  }

  function maybeRenderStashAlert() {
    if (!showStashAlert) {
      return;
    }

    return (
      <ModalComponent
        show
        icon={faExclamationTriangle}
        accept={{
          text: intl.formatMessage({ id: "actions.confirm" }),
          variant: "danger",
          onClick: () => {
            setShowStashAlert(false);
            next();
          },
        }}
        cancel={{ onClick: () => setShowStashAlert(false) }}
      >
        <p>
          <FormattedMessage id="setup.paths.stash_alert" />
        </p>
      </ModalComponent>
    );
  }

  function renderWelcomeSpecificConfig() {
    return (
      <>
        <section>
          <h2 className="mb-5">
            <FormattedMessage id="setup.welcome_to_stash" />
          </h2>
          <p className="lead text-center">
            <FormattedMessage id="setup.welcome_specific_config.unable_to_locate_specified_config" />
          </p>
          <p>
            <FormattedMessage
              id="setup.welcome_specific_config.config_path"
              values={{
                path: overrideConfig,
                code: (chunks: string) => <code>{chunks}</code>,
              }}
            />
          </p>
          <p>
            <FormattedMessage id="setup.welcome_specific_config.next_step" />
          </p>
        </section>

        <section className="mt-5">
          <div className="d-flex justify-content-center">
            <Button variant="primary mx-2 p-5" onClick={() => next()}>
              <FormattedMessage id="actions.next_action" />
            </Button>
          </div>
        </section>
      </>
    );
  }

  function renderWelcome() {
    const homeDirPath = pathJoin(status?.homeDir ?? homeDir, ".stash");

    return (
      <>
        <section>
          <h2 className="mb-5">
            <FormattedMessage id="setup.welcome_to_stash" />
          </h2>
          <p className="lead text-center">
            <FormattedMessage id="setup.welcome.unable_to_locate_config" />
          </p>
          <p>
            <FormattedMessage
              id="setup.welcome.config_path_logic_explained"
              values={{
                code: (chunks: string) => <code>{chunks}</code>,
                fallback_path: fallbackConfigPath,
              }}
            />
          </p>
          <Alert variant="info text-center">
            <FormattedMessage
              id="setup.welcome.unexpected_explained"
              values={{
                code: (chunks: string) => <code>{chunks}</code>,
              }}
            />
          </Alert>
          <p>
            <FormattedMessage id="setup.welcome.next_step" />
          </p>
        </section>

        <section className="mt-5">
          <h3 className="text-center mb-5">
            <FormattedMessage id="setup.welcome.store_stash_config" />
          </h3>

          <div className="d-flex justify-content-center">
            <Button
              variant="secondary mx-2 p-5"
              onClick={() => onConfigLocationChosen(false)}
            >
              <FormattedMessage
                id="setup.welcome.in_current_stash_directory"
                values={{
                  code: (chunks: string) => <code>{chunks}</code>,
                  path: fallbackStashDir,
                }}
              />
              <br />
              <code>{homeDirPath}</code>
            </Button>
            <Button
              variant="secondary mx-2 p-5"
              onClick={() => onConfigLocationChosen(true)}
              disabled={macApp}
            >
              {macApp ? (
                <>
                  <FormattedMessage
                    id="setup.welcome.in_the_current_working_directory_disabled"
                    values={{
                      code: (chunks: string) => <code>{chunks}</code>,
                      path: pwd,
                    }}
                  />
                  <br />
                  <b>
                    <FormattedMessage
                      id="setup.welcome.in_the_current_working_directory_disabled_macos"
                      values={{
                        code: (chunks: string) => <code>{chunks}</code>,
                        br: () => <br />,
                      }}
                    />
                  </b>
                </>
              ) : (
                <>
                  <FormattedMessage
                    id="setup.welcome.in_the_current_working_directory"
                    values={{
                      code: (chunks: string) => <code>{chunks}</code>,
                      path: pwd,
                    }}
                  />
                  <br />
                  <code>{workingDir}</code>
                </>
              )}
            </Button>
          </div>
        </section>
      </>
    );
  }

  function onGeneratedSelectClosed(d?: string) {
    if (d) {
      setGeneratedLocation(d);
    }

    setShowGeneratedSelectDialog(false);
  }

  function maybeRenderGeneratedSelectDialog() {
    if (!showGeneratedSelectDialog) {
      return;
    }

    return <FolderSelectDialog onClose={onGeneratedSelectClosed} />;
  }

  function onBlobsClosed(d?: string) {
    if (d) {
      setBlobsLocation(d);
    }

    setShowBlobsDialog(false);
  }

  function maybeRenderBlobsSelectDialog() {
    if (!showBlobsDialog) {
      return;
    }

    return <FolderSelectDialog onClose={onBlobsClosed} />;
  }

  function maybeRenderDatabase() {
    if (overrideDatabase) return;

    return (
      <Form.Group id="database">
        <h3>
          <FormattedMessage id="setup.paths.where_can_stash_store_its_database" />
        </h3>
        <p>
          <FormattedMessage
            id="setup.paths.where_can_stash_store_its_database_description"
            values={{
              code: (chunks: string) => <code>{chunks}</code>,
            }}
          />
          <br />
          <FormattedMessage
            id="setup.paths.where_can_stash_store_its_database_warning"
            values={{
              strong: (chunks: string) => <strong>{chunks}</strong>,
            }}
          />
        </p>
        <Form.Control
          className="text-input"
          defaultValue={databaseFile}
          placeholder={intl.formatMessage({
            id: "setup.paths.database_filename_empty_for_default",
          })}
          onChange={(e) => setDatabaseFile(e.currentTarget.value)}
        />
      </Form.Group>
    );
  }

  function maybeRenderGenerated() {
    if (overrideGenerated) return;

    return (
      <Form.Group id="generated">
        <h3>
          <FormattedMessage id="setup.paths.where_can_stash_store_its_generated_content" />
        </h3>
        <p>
          <FormattedMessage
            id="setup.paths.where_can_stash_store_its_generated_content_description"
            values={{
              code: (chunks: string) => <code>{chunks}</code>,
            }}
          />
        </p>
        <InputGroup>
          <Form.Control
            className="text-input"
            value={generatedLocation}
            placeholder={intl.formatMessage({
              id: "setup.paths.path_to_generated_directory_empty_for_default",
            })}
            onChange={(e) => setGeneratedLocation(e.currentTarget.value)}
          />
          <InputGroup.Append>
            <Button
              variant="secondary"
              className="text-input"
              onClick={() => setShowGeneratedSelectDialog(true)}
            >
              <Icon icon={faEllipsisH} />
            </Button>
          </InputGroup.Append>
        </InputGroup>
      </Form.Group>
    );
  }

  function onCacheSelectClosed(d?: string) {
    if (d) {
      setCacheLocation(d);
    }

    setShowCacheSelectDialog(false);
  }

  function maybeRenderCacheSelectDialog() {
    if (!showCacheSelectDialog) {
      return;
    }

    return <FolderSelectDialog onClose={onCacheSelectClosed} />;
  }

  function maybeRenderCache() {
    if (overrideCache) return;

    return (
      <Form.Group id="cache">
        <h3>
          <FormattedMessage id="setup.paths.where_can_stash_store_cache_files" />
        </h3>
        <p>
          <FormattedMessage
            id="setup.paths.where_can_stash_store_cache_files_description"
            values={{
              code: (chunks: string) => <code>{chunks}</code>,
            }}
          />
        </p>
        <InputGroup>
          <Form.Control
            className="text-input"
            value={cacheLocation}
            placeholder={intl.formatMessage({
              id: "setup.paths.path_to_cache_directory_empty_for_default",
            })}
            onChange={(e) => setCacheLocation(e.currentTarget.value)}
          />
          <InputGroup.Append>
            <Button
              variant="secondary"
              className="text-input"
              onClick={() => setShowCacheSelectDialog(true)}
            >
              <Icon icon={faEllipsisH} />
            </Button>
          </InputGroup.Append>
        </InputGroup>
      </Form.Group>
    );
  }

  function maybeRenderBlobs() {
    if (overrideBlobs) return;

    return (
      <Form.Group id="blobs">
        <h3>
          <FormattedMessage id="setup.paths.where_can_stash_store_blobs" />
        </h3>
        <p>
          <FormattedMessage
            id="setup.paths.where_can_stash_store_blobs_description"
            values={{
              code: (chunks: string) => <code>{chunks}</code>,
            }}
          />
        </p>
        <p>
          <FormattedMessage
            id="setup.paths.where_can_stash_store_blobs_description_addendum"
            values={{
              code: (chunks: string) => <code>{chunks}</code>,
              strong: (chunks: string) => <strong>{chunks}</strong>,
            }}
          />
        </p>

        <p>
          <Form.Check
            id="store-blobs-in-database"
            checked={storeBlobsInDatabase}
            label={intl.formatMessage({
              id: "setup.paths.store_blobs_in_database",
            })}
            onChange={() => setStoreBlobsInDatabase(!storeBlobsInDatabase)}
          />
        </p>

        {!storeBlobsInDatabase && (
          <InputGroup>
            <Form.Control
              className="text-input"
              value={blobsLocation}
              placeholder={intl.formatMessage({
                id: "setup.paths.path_to_blobs_directory_empty_for_default",
              })}
              onChange={(e) => setBlobsLocation(e.currentTarget.value)}
              disabled={storeBlobsInDatabase}
            />
            <InputGroup.Append>
              <Button
                variant="secondary"
                className="text-input"
                onClick={() => setShowBlobsDialog(true)}
                disabled={storeBlobsInDatabase}
              >
                <Icon icon={faEllipsisH} />
              </Button>
            </InputGroup.Append>
          </InputGroup>
        )}
      </Form.Group>
    );
  }

  function renderSetPaths() {
    return (
      <>
        {maybeRenderStashAlert()}
        <section>
          <h2 className="mb-3">
            <FormattedMessage id="setup.paths.set_up_your_paths" />
          </h2>
          <p>
            <FormattedMessage id="setup.paths.description" />
          </p>
        </section>
        <section>
          <Form.Group id="stashes">
            <h3>
              <FormattedMessage id="setup.paths.where_is_your_porn_located" />
            </h3>
            <p>
              <FormattedMessage id="setup.paths.where_is_your_porn_located_description" />
            </p>
            <Card>
              <StashConfiguration
                stashes={stashes}
                setStashes={(s) => setStashes(s)}
              />
            </Card>
          </Form.Group>
          {maybeRenderDatabase()}
          {maybeRenderGenerated()}
          {maybeRenderCache()}
          {maybeRenderBlobs()}
        </section>
        <section className="mt-5">
          <div className="d-flex justify-content-center">
            <Button variant="secondary mx-2 p-5" onClick={() => goBack()}>
              <FormattedMessage id="actions.previous_action" />
            </Button>
            <Button variant="primary mx-2 p-5" onClick={() => confirmPaths()}>
              <FormattedMessage id="actions.next_action" />
            </Button>
          </div>
        </section>
      </>
    );
  }

  function maybeRenderExclusions(s: GQL.StashConfig) {
    if (!s.excludeImage && !s.excludeVideo) {
      return;
    }

    const excludes = [];
    if (s.excludeVideo) {
      excludes.push("videos");
    }
    if (s.excludeImage) {
      excludes.push("images");
    }

    return `(excludes ${excludes.join(" and ")})`;
  }

  async function onSave() {
    let configLocation = overrideConfig;
    if (!configLocation) {
      configLocation = setupInWorkDir ? "config.yml" : "";
    }

    try {
      setLoading(true);
      await mutateSetup({
        configLocation,
        databaseFile,
        generatedLocation,
        cacheLocation,
        storeBlobsInDatabase,
        blobsLocation,
        stashes,
      });
      // Set lastNoteSeen to hide release notes dialog
      await saveUI({
        variables: {
          input: {
            ...configuration?.ui,
            lastNoteSeen: releaseNotes[0].date,
          },
        },
      });
    } catch (e) {
      if (e instanceof Error && e.message) {
        setSetupError(e.message);
      } else {
        setSetupError(String(e));
      }
    } finally {
      setLoading(false);
      next();
    }
  }

  function renderConfirm() {
    let cfgDir: string;
    let config: string;
    if (overrideConfig) {
      cfgDir = pathDir(overrideConfig);
      config = overrideConfig;
    } else {
      cfgDir = setupInWorkDir ? pwd : fallbackStashDir;
      config = pathJoin(cfgDir, "config.yml");
    }

    function joinCfgDir(path: string) {
      if (cfgDir) {
        return pathJoin(cfgDir, path);
      } else {
        return path;
      }
    }

    return (
      <>
        <section>
          <h2 className="mb-3">
            <FormattedMessage id="setup.confirm.nearly_there" />
          </h2>
          <p>
            <FormattedMessage id="setup.confirm.almost_ready" />
          </p>
          <dl>
            <dt>
              <FormattedMessage id="setup.confirm.configuration_file_location" />
            </dt>
            <dd>
              <code>{config}</code>
            </dd>
          </dl>
          <dl>
            <dt>
              <FormattedMessage id="setup.confirm.stash_library_directories" />
            </dt>
            <dd>
              <ul>
                {stashes.map((s) => (
                  <li key={s.path}>
                    <code>{s.path} </code>
                    {maybeRenderExclusions(s)}
                  </li>
                ))}
              </ul>
            </dd>
          </dl>
          {!overrideDatabase && (
            <dl>
              <dt>
                <FormattedMessage id="setup.confirm.database_file_path" />
              </dt>
              <dd>
                <code>{databaseFile || joinCfgDir("stash-go.sqlite")}</code>
              </dd>
            </dl>
          )}
          {!overrideGenerated && (
            <dl>
              <dt>
                <FormattedMessage id="setup.confirm.generated_directory" />
              </dt>
              <dd>
                <code>{generatedLocation || joinCfgDir("generated")}</code>
              </dd>
            </dl>
          )}
          {!overrideCache && (
            <dl>
              <dt>
                <FormattedMessage id="setup.confirm.cache_directory" />
              </dt>
              <dd>
                <code>{cacheLocation || joinCfgDir("cache")}</code>
              </dd>
            </dl>
          )}
          {!overrideBlobs && (
            <dl>
              <dt>
                <FormattedMessage id="setup.confirm.blobs_directory" />
              </dt>
              <dd>
                <code>
                  {storeBlobsInDatabase ? (
                    <FormattedMessage id="setup.confirm.blobs_use_database" />
                  ) : (
                    blobsLocation || joinCfgDir("blobs")
                  )}
                </code>
              </dd>
            </dl>
          )}
        </section>
        <section className="mt-5">
          <div className="d-flex justify-content-center">
            <Button variant="secondary mx-2 p-5" onClick={() => goBack()}>
              <FormattedMessage id="actions.previous_action" />
            </Button>
            <Button variant="success mx-2 p-5" onClick={() => onSave()}>
              <FormattedMessage id="actions.confirm" />
            </Button>
          </div>
        </section>
      </>
    );
  }

  function renderError() {
    function onBackClick() {
      setSetupError(undefined);
      goBack(2);
    }

    return (
      <>
        <section>
          <h2>
            <FormattedMessage id="setup.errors.something_went_wrong" />
          </h2>
          <p>
            <FormattedMessage
              id="setup.errors.something_went_wrong_while_setting_up_your_system"
              values={{ error: <pre>{setupError}</pre> }}
            />
          </p>
          <p>
            <FormattedMessage
              id="setup.errors.something_went_wrong_description"
              values={{ githubLink, discordLink }}
            />
          </p>
        </section>
        <section className="mt-5">
          <div className="d-flex justify-content-center">
            <Button variant="secondary mx-2 p-5" onClick={onBackClick}>
              <FormattedMessage id="actions.previous_action" />
            </Button>
          </div>
        </section>
      </>
    );
  }

  function renderSuccess() {
    return (
      <>
        <section>
          <h2>
            <FormattedMessage id="setup.success.your_system_has_been_created" />
          </h2>
          <p>
            <FormattedMessage id="setup.success.next_config_step_one" />
          </p>
          <p>
            <FormattedMessage
              id="setup.success.next_config_step_two"
              values={{
                code: (chunks: string) => <code>{chunks}</code>,
                localized_task: intl.formatMessage({
                  id: "config.categories.tasks",
                }),
                localized_scan: intl.formatMessage({ id: "actions.scan" }),
              }}
            />
          </p>
        </section>
        <section>
          <h3>
            <FormattedMessage id="setup.success.getting_help" />
          </h3>
          <p>
            <FormattedMessage
              id="setup.success.in_app_manual_explained"
              values={{ icon: <Icon icon={faQuestionCircle} /> }}
            />
          </p>
          <p>
            <FormattedMessage
              id="setup.success.help_links"
              values={{ discordLink, githubLink }}
            />
          </p>
        </section>
        <section>
          <h3>
            <FormattedMessage id="setup.success.support_us" />
          </h3>
          <p>
            <FormattedMessage
              id="setup.success.open_collective"
              values={{
                open_collective_link: (
                  <ExternalLink href="https://opencollective.com/stashapp">
                    Open Collective
                  </ExternalLink>
                ),
              }}
            />
          </p>
          <p>
            <FormattedMessage id="setup.success.welcome_contrib" />
          </p>
        </section>
        <section>
          <p className="lead text-center">
            <FormattedMessage id="setup.success.thanks_for_trying_stash" />
          </p>
        </section>
        <section className="mt-5">
          <div className="d-flex justify-content-center">
            <Link to="/settings?tab=library">
              <Button variant="success mx-2 p-5">
                <FormattedMessage id="actions.finish" />
              </Button>
            </Link>
          </div>
        </section>
      </>
    );
  }

  function renderFinish() {
    if (setupError !== undefined) {
      return renderError();
    }

    return renderSuccess();
  }

  // only display setup wizard if system is not setup
  if (statusLoading || configLoading) {
    return <LoadingIndicator />;
  }

  if (step === 0 && status && status.status !== GQL.SystemStatusEnum.Setup) {
    // redirect to main page
    history.push("/");
    return <LoadingIndicator />;
  }

  const welcomeStep = overrideConfig
    ? renderWelcomeSpecificConfig
    : renderWelcome;
  const steps = [welcomeStep, renderSetPaths, renderConfirm, renderFinish];

  function renderCreating() {
    return (
      <Card>
        <LoadingIndicator
          message={intl.formatMessage({
            id: "setup.creating.creating_your_system",
          })}
        />
        <Alert variant="info text-center">
          <FormattedMessage
            id="setup.creating.ffmpeg_notice"
            values={{
              code: (chunks: string) => <code>{chunks}</code>,
            }}
          />
        </Alert>
      </Card>
    );
  }

  return (
    <Container>
      {maybeRenderGeneratedSelectDialog()}
      {maybeRenderCacheSelectDialog()}
      {maybeRenderBlobsSelectDialog()}
      <h1 className="text-center">
        <FormattedMessage id="setup.stash_setup_wizard" />
      </h1>
      {loading ? renderCreating() : <Card>{steps[step]()}</Card>}
    </Container>
  );
};

export default Setup;
