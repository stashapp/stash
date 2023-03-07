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

export const Setup: React.FC = () => {
  const { configuration, loading: configLoading } =
    useContext(ConfigurationContext);
  const [saveUI] = useConfigureUI();

  const [step, setStep] = useState(0);
  const [configLocation, setConfigLocation] = useState("");
  const [stashes, setStashes] = useState<GQL.StashConfig[]>([]);
  const [showStashAlert, setShowStashAlert] = useState(false);
  const [databaseFile, setDatabaseFile] = useState("");
  const [generatedLocation, setGeneratedLocation] = useState("");
  const [cacheLocation, setCacheLocation] = useState("");
  const [loading, setLoading] = useState(false);
  const [setupError, setSetupError] = useState("");

  const intl = useIntl();
  const history = useHistory();

  const [showGeneratedSelectDialog, setShowGeneratedSelectDialog] =
    useState(false);
  const [showCacheSelectDialog, setShowCacheSelectDialog] = useState(false);

  const { data: systemStatus, loading: statusLoading } = useSystemStatus();

  useEffect(() => {
    if (systemStatus?.systemStatus.configPath) {
      setConfigLocation(systemStatus.systemStatus.configPath);
    }
  }, [systemStatus]);

  useEffect(() => {
    if (configuration) {
      const { stashes: configStashes, generatedPath } = configuration.general;
      if (configStashes.length > 0) {
        setStashes(
          configStashes.map((s) => {
            const { __typename, ...withoutTypename } = s;
            return withoutTypename;
          })
        );
      }
      if (generatedPath) {
        setGeneratedLocation(generatedPath);
      }
    }
  }, [configuration]);

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
      <FormattedMessage id="setup.github_repository" />
    </a>
  );

  function onConfigLocationChosen(loc: string) {
    setConfigLocation(loc);
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
                path: configLocation,
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
              onClick={() => onConfigLocationChosen("")}
            >
              <FormattedMessage
                id="setup.welcome.in_current_stash_directory"
                values={{
                  code: (chunks: string) => <code>{chunks}</code>,
                }}
              />
            </Button>
            <Button
              variant="secondary mx-2 p-5"
              onClick={() => onConfigLocationChosen("config.yml")}
            >
              <FormattedMessage id="setup.welcome.in_the_current_working_directory" />
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

  function maybeRenderGenerated() {
    if (!configuration?.general.generatedPath) {
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
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setGeneratedLocation(e.currentTarget.value)
              }
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
    if (!configuration?.general.cachePath) {
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
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setCacheLocation(e.currentTarget.value)
              }
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
              onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                setDatabaseFile(e.currentTarget.value)
              }
            />
          </Form.Group>
          {maybeRenderGenerated()}
          {maybeRenderCache()}
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

  function renderConfigLocation() {
    if (configLocation === "config.yml") {
      return <code>&lt;current working directory&gt;/config.yml</code>;
    }

    if (configLocation === "") {
      return <code>$HOME/.stash/config.yml</code>;
    }

    return <code>{configLocation}</code>;
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

  function renderStashLibraries() {
    return (
      <ul>
        {stashes.map((s) => (
          <li key={s.path}>
            <code>{s.path} </code>
            {maybeRenderExclusions(s)}
          </li>
        ))}
      </ul>
    );
  }

  async function onSave() {
    try {
      setLoading(true);
      await mutateSetup({
        configLocation,
        databaseFile,
        generatedLocation,
        cacheLocation,
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
      if (e instanceof Error) setSetupError(e.message ?? e.toString());
    } finally {
      setLoading(false);
      next();
    }
  }

  function renderConfirm() {
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
            <dd>{renderConfigLocation()}</dd>
          </dl>
          <dl>
            <dt>
              <FormattedMessage id="setup.confirm.stash_library_directories" />
            </dt>
            <dd>{renderStashLibraries()}</dd>
          </dl>
          <dl>
            <dt>
              <FormattedMessage id="setup.confirm.database_file_path" />
            </dt>
            <dd>
              <code>
                {databaseFile !== ""
                  ? databaseFile
                  : intl.formatMessage({
                      id: "setup.confirm.default_db_location",
                    })}
              </code>
            </dd>
          </dl>
          <dl>
            <dt>
              <FormattedMessage id="setup.confirm.generated_directory" />
            </dt>
            <dd>
              <code>
                {generatedLocation !== ""
                  ? generatedLocation
                  : intl.formatMessage({
                      id: "setup.confirm.default_generated_content_location",
                    })}
              </code>
            </dd>
          </dl>
          <dl>
            <dt>
              <FormattedMessage id="setup.confirm.cache_directory" />
            </dt>
            <dd>
              <code>
                {cacheLocation !== ""
                  ? cacheLocation
                  : intl.formatMessage({
                      id: "setup.confirm.default_cache_location",
                    })}
              </code>
            </dd>
          </dl>
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
            <Button variant="secondary mx-2 p-5" onClick={() => goBack(2)}>
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
                  <a
                    href="https://opencollective.com/stashapp"
                    target="_blank"
                    rel="noreferrer"
                  >
                    {" "}
                    OpenCollective{" "}
                  </a>
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
    if (setupError) {
      return renderError();
    }

    return renderSuccess();
  }

  // only display setup wizard if system is not setup
  if (statusLoading || configLoading) {
    return <LoadingIndicator />;
  }

  if (
    step === 0 &&
    systemStatus &&
    systemStatus.systemStatus.status !== GQL.SystemStatusEnum.Setup
  ) {
    // redirect to main page
    history.push("/");
    return <LoadingIndicator />;
  }

  const welcomeStep =
    systemStatus && systemStatus.systemStatus.configPath !== ""
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
      <h1 className="text-center">
        <FormattedMessage id="setup.stash_setup_wizard" />
      </h1>
      {loading ? renderCreating() : <Card>{steps[step]()}</Card>}
    </Container>
  );
};

export default Setup;
